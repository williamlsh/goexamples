//go:build !js
// +build !js

package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/mux"
	"github.com/pion/rtcp"
	"github.com/pion/rtp"
	"github.com/pion/webrtc/v3"
	"github.com/rs/zerolog/log"
	"github.com/williamlsh/logging"
	"nhooyr.io/websocket"
	"nhooyr.io/websocket/wsjson"
)

func init() {
	logging.Debug(true)
}

type message struct {
	Event string `json:"event"`
	SDP   string `json:"sdp"`
}

type udpConn struct {
	conn        *net.UDPConn
	port        int
	payloadType uint8
}

func main() {
	log.Fatal().AnErr("err", http.ListenAndServe(":8080", serveMux())).Send()
}

func serveMux() http.Handler {
	r := mux.NewRouter()

	r.Handle("/", http.FileServer(http.Dir("static")))
	r.HandleFunc("/webrtc", handleWebrtc)

	return r
}

func handleWebrtc(w http.ResponseWriter, r *http.Request) {
	c, err := websocket.Accept(w, r, &websocket.AcceptOptions{
		OriginPatterns: []string{"*"},
	})
	if err != nil {
		log.Err(err).Msg("websocket could not accept connection")
		return
	}
	defer func() {
		c.Close(websocket.StatusAbnormalClosure, "an error occurred")
		os.Exit(0)
	}()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	for {
		var msg message
		if err := wsjson.Read(ctx, c, &msg); err != nil {
			panic(err)
		}
		log.Info().Str("event", msg.Event).Str("data", msg.SDP).Msg("got message")

		switch msg.Event {
		case "offer":
			var sdp webrtc.SessionDescription
			if err := json.Unmarshal([]byte(msg.SDP), &sdp); err != nil {
				log.Err(err).Msg("could not unmarshal offer")
				return
			}

			answer, err := createPeerConnection(sdp)
			if err != nil {
				log.Err(err).Msg("failed to create peer connection")
				return
			}
			b, err := json.Marshal(answer)
			if err != nil {
				panic(err)
			}
			if err := wsjson.Write(ctx, c, &message{
				Event: "answer",
				SDP:   string(b),
			}); err != nil {
				panic(err)
			}
		}
	}
}

func createPeerConnection(offer webrtc.SessionDescription) (*webrtc.SessionDescription, error) {
	// Prepare the configuration
	config := webrtc.Configuration{
		ICEServers: []webrtc.ICEServer{
			{
				URLs: []string{"stun:stun.l.google.com:19302"},
			},
		},
	}

	// Create a new RTCPeerConnection
	peerConnection, err := webrtc.NewPeerConnection(config)
	if err != nil {
		log.Err(err).Msg("could not create new peer connection")
		return nil, err
	}
	defer func() {
		if cErr := peerConnection.Close(); cErr != nil {
			fmt.Printf("cannot close peerConnection: %v\n", cErr)
		}
	}()

	// Allow us to receive 1 audio track, and 1 video track
	if _, err = peerConnection.AddTransceiverFromKind(webrtc.RTPCodecTypeAudio); err != nil {
		log.Err(err).Msg("could not add audio transceiver")
		return nil, err
	}
	if _, err = peerConnection.AddTransceiverFromKind(webrtc.RTPCodecTypeVideo); err != nil {
		log.Err(err).Msg("could not add video transceiver")
		return nil, err
	}

	// Create a local addr
	var laddr *net.UDPAddr
	if laddr, err = net.ResolveUDPAddr("udp", "127.0.0.1:"); err != nil {
		log.Err(err).Msg("could not resolve udp")
		return nil, err
	}

	// Prepare udp conns
	// Also update incoming packets with expected PayloadType, the browser may use
	// a different value. We have to modify so our stream matches what rtp-forwarder.sdp expects
	udpConns := map[string]*udpConn{
		"audio": {port: 4000, payloadType: 111},
		"video": {port: 4002, payloadType: 96},
	}
	for _, c := range udpConns {
		// Create remote addr
		var raddr *net.UDPAddr
		if raddr, err = net.ResolveUDPAddr("udp", fmt.Sprintf("127.0.0.1:%d", c.port)); err != nil {
			log.Err(err).Msg("could not resolve udp")
			return nil, err
		}

		// Dial udp
		if c.conn, err = net.DialUDP("udp", laddr, raddr); err != nil {
			log.Err(err).Msg("could not dial udp")
			return nil, err
		}
		defer func(conn net.PacketConn) {
			if closeErr := conn.Close(); closeErr != nil {
				log.Err(closeErr).Msg("could not close udp connection")
			}
		}(c.conn)
	}

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
		fmt.Printf("Connection State has changed %s \n", connectionState.String())

		if connectionState == webrtc.ICEConnectionStateConnected {
			fmt.Println("Ctrl+C the remote client to stop the demo")
		}
	})

	// Set the handler for Peer connection state
	// This will notify you when the peer has connected/disconnected
	peerConnection.OnConnectionStateChange(func(s webrtc.PeerConnectionState) {
		fmt.Printf("Peer Connection State has changed: %s\n", s.String())

		if s == webrtc.PeerConnectionStateFailed {
			// Wait until PeerConnection has had no network activity for 30 seconds or another failure. It may be reconnected using an ICE Restart.
			// Use webrtc.PeerConnectionStateDisconnected if you are interested in detecting faster timeout.
			// Note that the PeerConnection may come back from PeerConnectionStateDisconnected.
			fmt.Println("Done forwarding")
			os.Exit(0)
		}
	})

	// Set the remote SessionDescription
	if err = peerConnection.SetRemoteDescription(offer); err != nil {
		log.Err(err).Msg("could not set remote description")
		return nil, err
	}

	// Create answer
	answer, err := peerConnection.CreateAnswer(nil)
	if err != nil {
		return nil, err
	}

	// Create channel that is blocked until ICE Gathering is complete
	gatherComplete := webrtc.GatheringCompletePromise(peerConnection)

	// Sets the LocalDescription, and starts our UDP listeners
	if err = peerConnection.SetLocalDescription(answer); err != nil {
		return nil, err
	}

	// Block until ICE Gathering is complete, disabling trickle ICE
	// we do this because we only can exchange one signaling message
	// in a production application you should exchange ICE Candidates via OnICECandidate
	<-gatherComplete

	return peerConnection.LocalDescription(), nil
}
