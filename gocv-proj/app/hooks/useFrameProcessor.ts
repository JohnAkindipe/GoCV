import { useRef, useState, useCallback } from 'react';

export function useFrameProcessor() {
  const [isProcessing, setIsProcessing] = useState(false);
  const [processedImage, setProcessedImage] = useState<string | null>(null);
  const [asciiText, setAsciiText] = useState<string | null>(null);
  const [expectingASCII, setExpectingASCII] = useState(false);
  const intervalRef = useRef<NodeJS.Timeout | null>(null);
  const expectingASCIIRef = useRef(false);

  const captureAndSend = useCallback(async (videoElement: HTMLVideoElement) => {
    const canvas = document.createElement('canvas');
    canvas.width = videoElement.videoWidth;
    canvas.height = videoElement.videoHeight;

    const ctx = canvas.getContext('2d');
    if (!ctx) return;

    ctx.drawImage(videoElement, 0, 0);

    const blob = await new Promise<Blob | null>((resolve) =>
      canvas.toBlob(resolve, 'image/jpeg', 0.8)
    );
    if (!blob) {
      console.error('Failed to create blob from canvas');
      return;
    }

    const mode = expectingASCIIRef.current ? 'ascii' : 'glitch';
    console.log(`Sending frame: ${blob.size} bytes, mode: ${mode}`);

    try {
      const response = await fetch(`http://localhost:4000/process-frame?mode=${mode}`, {
        method: 'POST',
        headers: { 'Content-Type': 'image/jpeg' },
        body: blob,
      });

      if (!response.ok) {
        const errorText = await response.text();
        console.error(`Server error ${response.status}: ${errorText}`);
        return;
      }

      const contentType = response.headers.get('Content-Type') || '';

      if (contentType.includes('text/html')) {
        const text = await response.text();
        setAsciiText(text);
        setProcessedImage(null);
      } else {
        const responseBlob = await response.blob();
        const imageUrl = URL.createObjectURL(responseBlob);
        setProcessedImage(imageUrl);
        setAsciiText(null);
      }
    } catch (err) {
      console.error('Failed to send frame:', err);
    }
  }, []);

  const startProcessing = useCallback((videoElement: HTMLVideoElement) => {
    setIsProcessing(true);

    // Send first frame immediately
    captureAndSend(videoElement);

    // Then every 20 seconds
    intervalRef.current = setInterval(() => {
      captureAndSend(videoElement);
    }, 20000);
  }, [captureAndSend]);

  const stopProcessing = useCallback(() => {
    if (intervalRef.current) {
      clearInterval(intervalRef.current);
      intervalRef.current = null;
    }
    setIsProcessing(false);
    setProcessedImage(null);
    setAsciiText(null);
  }, []);

  const toggleASCIIMode = useCallback(() => {
    setExpectingASCII((prev) => {
      const newValue = !prev;
      expectingASCIIRef.current = newValue;
      return newValue;
    });
  }, []);

  return { isProcessing, processedImage, asciiText, expectingASCII, startProcessing, stopProcessing, toggleASCIIMode };
}
