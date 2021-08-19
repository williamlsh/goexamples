package main

import (
	"context"
	"encoding/json"
	"sync"

	"github.com/pion/mediadevices"
	"github.com/pion/mediadevices/pkg/frame"
	"github.com/pion/mediadevices/pkg/prop"
	"github.com/pion/webrtc/v3"
	"github.com/rs/zerolog/log"
	"nhooyr.io/websocket"
	"nhooyr.io/websocket/wsjson"

	// If you don't like x264, you can also use vpx by importing as below
	// "github.com/pion/mediadevices/pkg/codec/vpx" // This is required to use VP8/VP9 video encoder
	// or you can also use openh264 for alternative h264 implementation
	// "github.com/pion/mediadevices/pkg/codec/openh264"
	// or if you use a raspberry pi like, you can use mmal for using its hardware encoder
	// "github.com/pion/mediadevices/pkg/codec/mmal"
	"github.com/pion/mediadevices/pkg/codec/opus" // This is required to use opus audio encoder
	"github.com/pion/mediadevices/pkg/codec/x264" // This is required to use h264 video encoder

	// Note: If you don't have a camera or microphone or your adapters are not supported,
	//       you can always swap your adapters with our dummy adapters below.
	_ "github.com/pion/mediadevices/pkg/driver/audiotest"
	_ "github.com/pion/mediadevices/pkg/driver/videotest"
	// _ "github.com/pion/mediadevices/pkg/driver/camera"     // This is required to register camera adapter
	// _ "github.com/pion/mediadevices/pkg/driver/microphone" // This is required to register microphone adapter
)

type SignalFunc func(msg interface{}, event string) error

func createPeerConnection(signalPeer SignalFunc, answerCh <-chan *webrtc.SessionDescription, candidateCh <-chan *webrtc.ICECandidateInit) error {
	var candidatesMux sync.Mutex
	pendingCandidates := make([]*webrtc.ICECandidate, 0)

	config := webrtc.Configuration{
		ICEServers: []webrtc.ICEServer{
			{
				URLs:       []string{"turn:localhost:3478"},
				Username:   "user",
				Credential: "password",
			},
		},
	}

	x264Params, err := x264.NewParams()
	if err != nil {
		return err
	}
	x264Params.BitRate = 500_000 // 500kbps

	opusParams, err := opus.NewParams()
	if err != nil {
		return err
	}

	codecSelector := mediadevices.NewCodecSelector(
		mediadevices.WithVideoEncoders(&x264Params),
		mediadevices.WithAudioEncoders(&opusParams),
	)

	var mediaEngine webrtc.MediaEngine
	codecSelector.Populate(&mediaEngine)
	api := webrtc.NewAPI(webrtc.WithMediaEngine(&mediaEngine))
	peerConnection, err := api.NewPeerConnection(config)
	if err != nil {
		return err
	}

	peerConnection.OnICECandidate(func(c *webrtc.ICECandidate) {
		if c == nil {
			return
		}

		candidatesMux.Lock()
		defer candidatesMux.Unlock()

		if desc := peerConnection.RemoteDescription(); desc == nil {
			pendingCandidates = append(pendingCandidates, c)
			return
		}
		if err := signalPeer(c.ToJSON(), "candidate"); err != nil {
			log.Err(err).Msg("could not send candidate")
		}
	})

	peerConnection.OnICEConnectionStateChange(func(connectionState webrtc.ICEConnectionState) {
		log.Info().Str("state", connectionState.String()).Msg("connection state has changed")
	})

	s, err := mediadevices.GetUserMedia(mediadevices.MediaStreamConstraints{
		Video: func(c *mediadevices.MediaTrackConstraints) {
			c.FrameFormat = prop.FrameFormat(frame.FormatI420)
			c.Width = prop.Int(640)
			c.Height = prop.Int(480)
		},
		Audio: func(c *mediadevices.MediaTrackConstraints) {},
		Codec: codecSelector,
	})
	if err != nil {
		return err
	}

	for _, track := range s.GetTracks() {
		track.OnEnded(func(err error) {
			log.Err(err).Str("track_id", track.ID()).Msg("track ended with error")
		})

		_, err := peerConnection.AddTransceiverFromTrack(
			track,
			webrtc.RTPTransceiverInit{
				Direction: webrtc.RTPTransceiverDirectionSendonly,
			},
		)
		if err != nil {
			return err
		}
	}

	offer, err := peerConnection.CreateOffer(nil)
	if err != nil {
		panic(err)
	}

	if err = peerConnection.SetLocalDescription(offer); err != nil {
		panic(err)
	}

	if err := signalPeer(peerConnection.LocalDescription(), "offer"); err != nil {
		return err
	}

	answer := <-answerCh
	if err := peerConnection.SetRemoteDescription(*answer); err != nil {
		panic(err)
	}
	log.Info().Msg("set remote description")

	go func() {
		if err := addICECandidate(peerConnection, candidateCh); err != nil {
			log.Err(err).Msg("could not add candidate")
		}
	}()

	candidatesMux.Lock()
	defer candidatesMux.Unlock()

	for _, c := range pendingCandidates {
		if onICECandidateErr := signalPeer(c.ToJSON(), "candidate"); onICECandidateErr != nil {
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

// signalPeer sends offer or candidate to peer.
func signalPeer(ctx context.Context, conn *websocket.Conn) SignalFunc {
	return func(msg interface{}, event string) error {
		b, err := json.Marshal(msg)
		if err != nil {
			return err
		}
		if err := wsjson.Write(ctx, conn, &message{
			Event: event,
			Data:  string(b),
		}); err != nil {
			return err
		}
		log.Info().Str("event", event).Msg("signaled peer")
		return nil
	}
}
