package adapters

import (
	"fmt"

	"github.com/pion/rtcp"
	"github.com/pion/webrtc/v3"
)

type WebRTCAdapter struct {
	peerConnection *webrtc.PeerConnection
	onTrack        func(track *webrtc.TrackRemote)
}

func (a *WebRTCAdapter) RequestKeyframe(ssrc uint32) {
	if a.peerConnection == nil {
		return
	}
	fmt.Printf("Backend: Solicitando Keyframe (PLI) para SSRC %d...\n", ssrc)
	a.peerConnection.WriteRTCP([]rtcp.Packet{
		&rtcp.PictureLossIndication{MediaSSRC: ssrc},
	})
}

func NewWebRTCAdapter() (*WebRTCAdapter, error) {
	// Configuración básica (sin ICE servers externos necesarios para localhost)
	config := webrtc.Configuration{
		ICEServers: []webrtc.ICEServer{
			{
				URLs: []string{"stun:stun.l.google.com:19302"},
			},
		},
	}

	pc, err := webrtc.NewPeerConnection(config)
	if err != nil {
		return nil, err
	}

	return &WebRTCAdapter{
		peerConnection: pc,
	}, nil
}

func (a *WebRTCAdapter) ProcessOffer(sdp string, onTrack func(track *webrtc.TrackRemote)) (string, error) {
	a.onTrack = onTrack

	// Configurar el manejador de tracks
	a.peerConnection.OnTrack(func(track *webrtc.TrackRemote, receiver *webrtc.RTPReceiver) {
		fmt.Printf("Recibiendo track: %s (Payload Type: %d)\n", track.Kind(), track.PayloadType())
		if a.onTrack != nil {
			a.onTrack(track)
		}
	})

	// Set remote description
	offer := webrtc.SessionDescription{
		Type: webrtc.SDPTypeOffer,
		SDP:  sdp,
	}
	if err := a.peerConnection.SetRemoteDescription(offer); err != nil {
		return "", err
	}

	// Create answer
	answer, err := a.peerConnection.CreateAnswer(nil)
	if err != nil {
		return "", err
	}

	// Sets the LocalDescription, and starts our UDP listeners
	if err := a.peerConnection.SetLocalDescription(answer); err != nil {
		return "", err
	}

	return answer.SDP, nil
}

func (a *WebRTCAdapter) Close() error {
	if a.peerConnection != nil {
		return a.peerConnection.Close()
	}
	return nil
}
