package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/buildkite/terminal-to-html/v3"
	"github.com/pion/webrtc/v4"
)

var webrtcManager = NewWebRTCManager()

type OfferRequest struct {
	SDP  string `json:"sdp"`
	Type string `json:"type"`
}

type AnswerResponse struct {
	SDP  string `json:"sdp"`
	Type string `json:"type"`
}

type CandidateRequest struct {
	Candidate webrtc.ICECandidateInit `json:"candidate"`
}

func testHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Test endpoint reached"))
}

func offerHandler(w http.ResponseWriter, r *http.Request) {
	// Parse offer from client
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Failed to read body", http.StatusBadRequest)
		return
	}

	var offer OfferRequest
	if err := json.Unmarshal(body, &offer); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	log.Printf("Received offer from client")

	// Create new peer connection for this client
	clientID := "client-1" // In production, generate unique ID per client
	pc, err := webrtcManager.CreatePeerConnection(clientID)
	if err != nil {
		log.Printf("Failed to create peer connection: %v", err)
		http.Error(w, fmt.Sprintf("Failed to create peer connection: %v", err), http.StatusInternalServerError)
		return
	}

	// Handle data channel for image processing
	pc.OnDataChannel(func(dc *webrtc.DataChannel) {
		log.Printf("Data channel '%s' opened", dc.Label())

		dc.OnMessage(func(msg webrtc.DataChannelMessage) {
			timestamp := time.Now().Format("15:04:05")
			log.Printf("[%s] Received image data: %d bytes (binary: %v)", timestamp, len(msg.Data), msg.IsString == false)
			
			// Use faceToASCIIText to return ASCII string as bytes (much smaller payload)
			imageBytes, err := imageToASCII(msg.Data)
			if err != nil {
				log.Printf("Error processing image: %v", err)
				return
			}
			HTMLFormat := terminal.Render(imageBytes)
			dc.Send([]byte(HTMLFormat))
			log.Printf("[%s] Sent processed image data: %d bytes", timestamp, len(imageBytes))
		})

		dc.OnError(func(err error) {
			log.Printf("Data channel error: %v", err)
		})

		dc.OnClose(func() {
			log.Printf("Data channel closed")
		})
	})

	// Set up track handlers (receive video/audio)
	pc.OnTrack(func(track *webrtc.TrackRemote, receiver *webrtc.RTPReceiver) {
		log.Printf("Track received: %s (type: %s)", track.Kind().String(), track.Codec().MimeType)

		// Just drain the packets - we're using data channel for processing
		go func() {
			for {
				_, _, err := track.ReadRTP()
				if err != nil {
					return
				}
			}
		}()
	})

	// Handle ICE connection state changes
	pc.OnICEConnectionStateChange(func(state webrtc.ICEConnectionState) {
		// log.Printf("ICE connection state changed: %s", state.String())
	})

	// Handle connection state changes
	pc.OnConnectionStateChange(func(state webrtc.PeerConnectionState) {
		// log.Printf("Connection state changed: %s", state.String())
		if state == webrtc.PeerConnectionStateFailed {
			// log.Printf("Connection failed - check ICE candidates and network connectivity")
		}
	})

	// Log ICE candidates as they're gathered
	pc.OnICECandidate(func(candidate *webrtc.ICECandidate) {
		if candidate != nil {
			// log.Printf("New ICE candidate: %s", candidate.String())
		} else {
			// log.Printf("ICE gathering complete")
		}
	})

	// Set remote description (the offer from client)
	if err := pc.SetRemoteDescription(webrtc.SessionDescription{
		Type: webrtc.SDPTypeOffer,
		SDP:  offer.SDP,
	}); err != nil {
		http.Error(w, fmt.Sprintf("Failed to set remote description: %v", err), http.StatusInternalServerError)
		return
	}

	// Create answer
	answer, err := pc.CreateAnswer(nil)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to create answer: %v", err), http.StatusInternalServerError)
		return
	}

	// Set local description - this will start ICE gathering
	if err := pc.SetLocalDescription(answer); err != nil {
		http.Error(w, fmt.Sprintf("Failed to set local description: %v", err), http.StatusInternalServerError)
		return
	}

	// Wait for ICE gathering to complete so all candidates are in the SDP
	<-webrtc.GatheringCompletePromise(pc)

	// log.Printf("ICE gathering complete, sending answer with candidates")

	// Send answer back to client (now includes all ICE candidates)
	response := AnswerResponse{
		SDP:  pc.LocalDescription().SDP,
		Type: pc.LocalDescription().Type.String(),
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func candidateHandler(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Failed to read body", http.StatusBadRequest)
		return
	}

	var req CandidateRequest
	if err := json.Unmarshal(body, &req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	// Get peer connection
	clientID := "client-1"
	pc := webrtcManager.GetPeerConnection(clientID)
	if pc == nil {
		http.Error(w, "No peer connection found", http.StatusNotFound)
		return
	}

	// Add ICE candidate
	// log.Printf("Adding ICE candidate from client: %s", req.Candidate.Candidate)
	if err := pc.AddICECandidate(req.Candidate); err != nil {
		log.Printf("Failed to add ICE candidate: %v", err)
		http.Error(w, fmt.Sprintf("Failed to add ICE candidate: %v", err), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}