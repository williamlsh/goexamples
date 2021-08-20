package main

import (
	"fmt"
	"net"
	"os"
	"sync"
	"time"

	"github.com/pion/rtcp"
	"github.com/pion/rtp"
	"github.com/pion/webrtc/v3"
	"github.com/rs/zerolog/log"
)

type SignalFunc func(msg interface{}, event string) error

type udpConn struct {
	conn        net.Conn
	port        int
	payloadType uint8
}

func createPeerConnection(offer *webrtc.SessionDescription, candidateCh <-chan *webrtc.ICECandidateInit, signalPeer SignalFunc) error {
	var candidatesMux sync.Mutex
	pendingCandidates := make([]*webrtc.ICECandidate, 0)

	// Prepare the configuration
	config := webrtc.Configuration{
		ICEServers: []webrtc.ICEServer{
			{
				URLs:       []string{"turn:localhost:3478"},
				Username:   "user",
				Credential: "password",
			},
		},
	}

	// Create a new RTCPeerConnection
	peerConnection, err := webrtc.NewPeerConnection(config)
	if err != nil {
		log.Err(err).Msg("could not create new peer connection")
		return err
	}

	// Allow us to receive 1 audio track, and 1 video track
	if _, err = peerConnection.AddTransceiverFromKind(webrtc.RTPCodecTypeAudio); err != nil {
		log.Err(err).Msg("could not add audio transceiver")
		return err
	}
	if _, err = peerConnection.AddTransceiverFromKind(webrtc.RTPCodecTypeVideo); err != nil {
		log.Err(err).Msg("could not add video transceiver")
		return err
	}

	// Prepare udp conns
	// Also update incoming packets with expected PayloadType, the browser may use
	// a different value. We have to modify so our stream matches what rtp-forwarder.sdp expects
	udpConns := map[string]*udpConn{
		"audio": {port: 4000, payloadType: 111},
		"video": {port: 4002, payloadType: 96},
	}
	for _, c := range udpConns {
		c.conn, err = net.Dial("udp", fmt.Sprintf("127.0.0.1:%d", c.port))
		if err != nil {
			return err
		}
	}

	peerConnection.OnICECandidate(func(c *webrtc.ICECandidate) {
		if c == nil {
			return
		}

		if desc := peerConnection.RemoteDescription(); desc == nil {
			candidatesMux.Lock()
			pendingCandidates = append(pendingCandidates, c)
			candidatesMux.Unlock()

			return
		}
		// Send candidate.
		if err := signalPeer(c.ToJSON(), "candidate"); err != nil {
			log.Err(err).Msg("could not send candidate")
		}
	})

	// Set a handler for when a new remote track starts, this handler will forward data to
	// our UDP listeners.
	// In your application this is where you would handle/process audio/video
	peerConnection.OnTrack(func(track *webrtc.TrackRemote, receiver *webrtc.RTPReceiver) {
		// Retrieve udp connection
		c, ok := udpConns[track.Kind().String()]
		if !ok {
			return
		}

		// Send a PLI on an interval so that the publisher is pushing a keyframe every rtcpPLIInterval
		go func() {
			ticker := time.NewTicker(time.Second * 2)
			for range ticker.C {
				if rtcpErr := peerConnection.WriteRTCP([]rtcp.Packet{
					&rtcp.PictureLossIndication{MediaSSRC: uint32(track.SSRC())},
				}); rtcpErr != nil {
					log.Err(rtcpErr).Msg("could not write rtcp")
				}
			}
		}()

		b := make([]byte, 1500)
		rtpPacket := &rtp.Packet{}
		for {
			// Read
			n, _, readErr := track.Read(b)
			if readErr != nil {
				log.Err(readErr).Msg("could not read track data")
				return
			}

			// Unmarshal the packet and update the PayloadType
			if err = rtpPacket.Unmarshal(b[:n]); err != nil {
				log.Err(err).Msg("could not unmarshal rtp packat")
				return
			}
			rtpPacket.PayloadType = c.payloadType

			// Marshal into original buffer with updated PayloadType
			if n, err = rtpPacket.MarshalTo(b); err != nil {
				log.Err(err).Msg("could not marshal rtp packat")
				return
			}

			// Write
			if _, err = c.conn.Write(b[:n]); err != nil {
				// For this particular example, third party applications usually timeout after a short
				// amount of time during which the user doesn't have enough time to provide the answer
				// to the browser.
				// That's why, for this particular example, the user first needs to provide the answer
				// to the browser then open the third party application. Therefore, we must not kill
				// the forward on "connection refused" errors
				if opError, ok := err.(*net.OpError); ok && opError.Err.Error() == "write: connection refused" {
					continue
				}
				log.Err(err).Msg("could not write date to udp connection")
				return
			}
		}
	})

	// Set the handler for ICE connection state
	// This will notify you when the peer has connected/disconnected
	peerConnection.OnICEConnectionStateChange(func(connectionState webrtc.ICEConnectionState) {
		log.Info().Str("state", connectionState.String()).Msg("ICE connection state has changed")
	})

	// Set the handler for Peer connection state
	// This will notify you when the peer has connected/disconnected
	peerConnection.OnConnectionStateChange(func(s webrtc.PeerConnectionState) {
		log.Info().Str("state", s.String()).Msg("peer connection state has changed")

		if s == webrtc.PeerConnectionStateFailed {
			// Wait until PeerConnection has had no network activity for 30 seconds or another failure. It may be reconnected using an ICE Restart.
			// Use webrtc.PeerConnectionStateDisconnected if you are interested in detecting faster timeout.
			// Note that the PeerConnection may come back from PeerConnectionStateDisconnected.
			os.Exit(0)
		}
	})

	// Set the remote SessionDescription
	if err = peerConnection.SetRemoteDescription(*offer); err != nil {
		log.Err(err).Msg("could not set remote description")
		return err
	}
	log.Info().Msg("set remote description")

	// Create answer
	answer, err := peerConnection.CreateAnswer(nil)
	if err != nil {
		return err
	}
	log.Info().Msg("created answer")

	// Sets the LocalDescription
	if err = peerConnection.SetLocalDescription(answer); err != nil {
		return err
	}
	log.Info().Msg("set local description")

	// Sends the answer back.
	if err := signalPeer(&answer, "answer"); err != nil {
		return err
	}

	go func() {
		if err := addICECandidate(peerConnection, candidateCh); err != nil {
			log.Err(err).Msg("could not add candidate")
		}
	}()

	candidatesMux.Lock()
	defer candidatesMux.Unlock()

	for _, c := range pendingCandidates {
		if err := signalPeer(c.ToJSON(), "candidate"); err != nil {
			log.Err(err).Msg("could not send candidate")
			return err
		}
	}

	return nil
}

func addICECandidate(peerConnection *webrtc.PeerConnection, candidateCh <-chan *webrtc.ICECandidateInit) error {
	for c := range candidateCh {
		if err := peerConnection.AddICECandidate(*c); err != nil {
			return err
		}
		log.Info().Msg("added a candidate")
	}
	return nil
}
