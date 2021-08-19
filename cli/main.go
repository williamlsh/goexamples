package main

import (
	"context"
	"encoding/json"

	"github.com/pion/webrtc/v3"
	"github.com/rs/zerolog/log"
	"github.com/williamlsh/logging"
	"nhooyr.io/websocket"
	"nhooyr.io/websocket/wsjson"
)

type message struct {
	Event string `json:"event"`
	Data  string `json:"data"`
}

func init() {
	logging.Debug(true)
}

func main() {
	log.Fatal().AnErr("err", negotiate()).Send()
}

func negotiate() error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	c, _, err := websocket.Dial(ctx, "ws://localhost:9090/webrtc", nil)
	if err != nil {
		return err
	}
	defer c.Close(websocket.StatusInternalError, "the sky is falling")

	var (
		answerCh    = make(chan *webrtc.SessionDescription, 1)
		candidateCh = make(chan *webrtc.ICECandidateInit, 4)
	)

	go func() {
		if err := createPeerConnection(signalPeer(ctx, c), answerCh, candidateCh); err != nil {
			log.Err(err).Msg("could not create peer connection")
		}
	}()

	for {
		var msg message
		if err := wsjson.Read(ctx, c, &msg); err != nil {
			log.Err(err).Str("event", msg.Event).Str("data", msg.Data).Msg("received a message")
			return err
		}
		log.Info().Str("event", msg.Event).Str("data", msg.Data).Msg("got message")

		switch msg.Event {
		case "answer":
			var sdp webrtc.SessionDescription
			if err := json.Unmarshal([]byte(msg.Data), &sdp); err != nil {
				log.Err(err).Msg("could not unmarshal message data")
				return err
			}
			answerCh <- &sdp
		case "candidate":
			var sdp webrtc.ICECandidateInit
			if err := json.Unmarshal([]byte(msg.Data), &sdp); err != nil {
				log.Err(err).Msg("could not unmarshal message data")
				return err
			}
			candidateCh <- &sdp
		default:
			log.Info().Msg("unknown event")
		}
	}
}
