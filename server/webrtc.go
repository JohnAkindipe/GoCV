package main

import (
	"sync"

	"github.com/pion/webrtc/v4"
)

type WebRTCManager struct {
	mu              sync.RWMutex
	peerConnections map[string]*webrtc.PeerConnection
	api             *webrtc.API
}

func NewWebRTCManager() *WebRTCManager {
	// Configure media engine (what codecs we support)
	mediaEngine := &webrtc.MediaEngine{}

	// Support VP8 video codec
	if err := mediaEngine.RegisterCodec(webrtc.RTPCodecParameters{
		RTPCodecCapability: webrtc.RTPCodecCapability{
			MimeType:  webrtc.MimeTypeVP8,
			ClockRate: 90000,
			Channels:  0,
		},
		PayloadType: 96,
	}, webrtc.RTPCodecTypeVideo); err != nil {
		panic(err)
	}

	// Support Opus audio codec
	if err := mediaEngine.RegisterCodec(webrtc.RTPCodecParameters{
		RTPCodecCapability: webrtc.RTPCodecCapability{
			MimeType:  webrtc.MimeTypeOpus,
			ClockRate: 48000,
			Channels:  2,
		},
		PayloadType: 111,
	}, webrtc.RTPCodecTypeAudio); err != nil {
		panic(err)
	}

	// Configure NAT traversal settings
	settingEngine := webrtc.SettingEngine{}
	
	// Set public IP for NAT traversal
	settingEngine.SetNAT1To1IPs([]string{"161.35.36.3"}, webrtc.ICECandidateTypeHost)
	
	// Use a specific UDP port range (then open these ports in firewall)
	settingEngine.SetEphemeralUDPPortRange(10000, 20000)

	api := webrtc.NewAPI(
		webrtc.WithMediaEngine(mediaEngine),
		webrtc.WithSettingEngine(settingEngine),
	)

	return &WebRTCManager{
		peerConnections: make(map[string]*webrtc.PeerConnection),
		api:             api,
	}
}

func (m *WebRTCManager) CreatePeerConnection(id string) (*webrtc.PeerConnection, error) {
	config := webrtc.Configuration{
		ICEServers: []webrtc.ICEServer{
			{
				URLs: []string{
					"stun:stun.l.google.com:19302",
					"stun:stun1.l.google.com:19302",
				},
			},
			// Multiple TURN servers for better reliability
			{
				URLs: []string{
					"turn:openrelay.metered.ca:80",
					"turn:openrelay.metered.ca:443",
				},
				Username:   "openrelayproject",
				Credential: "openrelayproject",
			},
			{
				URLs: []string{
					"turn:relay.metered.ca:80",
					"turn:relay.metered.ca:443",
				},
				Username:   "9e95e1c078e4b8c20fc98f97",
				Credential: "iFhwFGiLLtc+Rq7a",
			},
		},
	}

	pc, err := m.api.NewPeerConnection(config)
	if err != nil {
		return nil, err
	}

	m.mu.Lock()
	m.peerConnections[id] = pc
	m.mu.Unlock()

	return pc, nil
}

func (m *WebRTCManager) GetPeerConnection(id string) *webrtc.PeerConnection {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.peerConnections[id]
}
