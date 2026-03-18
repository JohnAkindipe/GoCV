import { useRef, useState, useCallback } from 'react';

export function useWebRTC() {
  const [isConnected, setIsConnected] = useState(false);
  const [processedImage, setProcessedImage] = useState<string | null>(null);
  const [asciiText, setAsciiText] = useState<string | null>(null);
  const [expectingASCII, setExpectingASCII] = useState(false);
  const pcRef = useRef<RTCPeerConnection | null>(null);
  const dataChannelRef = useRef<RTCDataChannel | null>(null);
  const snapshotIntervalRef = useRef<NodeJS.Timeout | null>(null);
  const expectingASCIIRef = useRef(false);

  const captureAndSendSnapshot = useCallback((videoElement: HTMLVideoElement, dataChannel: RTCDataChannel) => {
    if (dataChannel.readyState !== 'open') {
      console.log('Data channel not open, skipping snapshot');
      return;
    }

    console.log('Capturing snapshot...');
    
    // Create canvas to capture frame
    const canvas = document.createElement('canvas');
    canvas.width = videoElement.videoWidth;
    canvas.height = videoElement.videoHeight;
    
    const ctx = canvas.getContext('2d');
    if (!ctx) return;
    
    ctx.drawImage(videoElement, 0, 0);
    
    // Convert to JPEG blob
    canvas.toBlob((blob) => {
      if (!blob) {
        console.error('Failed to create blob from canvas');
        return;
      }
      
      console.log(`Sending snapshot: ${blob.size} bytes`);
      
      // Send blob as ArrayBuffer
      blob.arrayBuffer().then((buffer) => {
        dataChannel.send(buffer);
      });
    }, 'image/jpeg', 0.8);
  }, []);

  const startStreaming = async (stream: MediaStream, videoElement: HTMLVideoElement) => {
    try {
      // Create peer connection with STUN server for NAT traversal
      const pc = new RTCPeerConnection({
        iceServers: [
          { 
            urls: [
              'stun:stun.l.google.com:19302',
              'stun:stun1.l.google.com:19302'
            ]
          },
          // Multiple TURN servers for symmetric NAT
          {
            urls: [
              'turn:openrelay.metered.ca:80',
              'turn:openrelay.metered.ca:443'
            ],
            username: 'openrelayproject',
            credential: 'openrelayproject',
          },
          {
            urls: [
              'turn:relay.metered.ca:80',
              'turn:relay.metered.ca:443'
            ],
            username: '9e95e1c078e4b8c20fc98f97',
            credential: 'iFhwFGiLLtc+Rq7a',
          },
        ],
        iceCandidatePoolSize: 10,
      });

      // Create data channel for sending/receiving images
      const dataChannel = pc.createDataChannel('images', {
        ordered: true,
      });

      dataChannel.binaryType = 'arraybuffer';

      dataChannel.onopen = () => {
        console.log('Data channel opened');
        
        // Start capturing and sending snapshots every 15 seconds
        snapshotIntervalRef.current = setInterval(() => {
          captureAndSendSnapshot(videoElement, dataChannel);
        }, 15000);

        // Send first snapshot immediately
        captureAndSendSnapshot(videoElement, dataChannel);
      };

      dataChannel.onmessage = (event) => {
        console.log('Received data from server, expecting ASCII:', expectingASCIIRef.current);
        
        if (expectingASCIIRef.current) {
          // Decode ArrayBuffer as text
          const decoder = new TextDecoder('utf-8');
          const text = decoder.decode(event.data);
          console.log('Received ASCII text:', text.length, 'characters');
          setAsciiText(text);
          setProcessedImage(null);
        } else {
          // Convert received ArrayBuffer to Blob and create URL
          const blob = new Blob([event.data], { type: 'image/jpeg' });
          const imageUrl = URL.createObjectURL(blob);
          setProcessedImage(imageUrl);
          setAsciiText(null);
        }
      };

      dataChannel.onerror = (error) => {
        console.error('Data channel error:', error);
      };

      dataChannel.onclose = () => {
        console.log('Data channel closed');
        if (snapshotIntervalRef.current) {
          clearInterval(snapshotIntervalRef.current);
          snapshotIntervalRef.current = null;
        }
      };

      dataChannelRef.current = dataChannel;

      // Add all tracks from the stream to the connection
      stream.getTracks().forEach((track) => {
        pc.addTrack(track, stream);
      });

      // Listen for ICE candidates (network paths)
      const iceCandidates: RTCIceCandidate[] = [];
      
      pc.onicecandidate = (event) => {
        if (event.candidate) {
          console.log('Generated ICE candidate:', event.candidate);
          iceCandidates.push(event.candidate);
        } else {
          console.log('ICE gathering complete, sending all candidates');
          // Send all candidates at once after gathering completes
          iceCandidates.forEach(async (candidate) => {
            try {
              await fetch('http://161.35.36.3:4000/webrtc/candidate', {
                method: 'POST',
                headers: { 'Content-Type': 'application/json' },
                body: JSON.stringify({ candidate }),
              });
            } catch (err) {
              console.error('Failed to send ICE candidate:', err);
            }
          });
        }
      };

      // Monitor connection state
      pc.onconnectionstatechange = () => {
        console.log('Connection state:', pc.connectionState);
        setIsConnected(pc.connectionState === 'connected');
      };

      // Create offer (SDP description of what we're sending)
      const offer = await pc.createOffer();
      await pc.setLocalDescription(offer);

      console.log('Sending offer to server:', { sdp: offer.sdp, type: offer.type });

      // Send offer to server and get answer
      const response = await fetch('http://161.35.36.3:4000/webrtc/offer', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ sdp: offer.sdp, type: offer.type }),
      });

      console.log('Offer response status:', response.status);

      if (!response.ok) {
        const errorText = await response.text();
        console.error('Offer error response:', errorText);
        throw new Error(`Server responded with ${response.status}: ${errorText}`);
      }

      const answer = await response.json();
      console.log('Received answer from server:', answer);
      
      await pc.setRemoteDescription(new RTCSessionDescription(answer));

      pcRef.current = pc;
    } catch (err) {
      console.error('WebRTC setup failed:', err);
      throw err;
    }
  };

  const stopStreaming = () => {
    if (snapshotIntervalRef.current) {
      clearInterval(snapshotIntervalRef.current);
      snapshotIntervalRef.current = null;
    }
    
    if (dataChannelRef.current) {
      dataChannelRef.current.close();
      dataChannelRef.current = null;
    }
    
    if (pcRef.current) {
      pcRef.current.close();
      pcRef.current = null;
      setIsConnected(false);
    }
    
    setProcessedImage(null);
    setAsciiText(null);
  };

  const toggleASCIIMode = () => {
    const newValue = !expectingASCII;
    setExpectingASCII(newValue);
    expectingASCIIRef.current = newValue;
  };

  return { isConnected, processedImage, asciiText, expectingASCII, startStreaming, stopStreaming, toggleASCIIMode };
}
