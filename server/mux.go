package main

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/pion/webrtc/v3"
	"github.com/rs/zerolog/log"
	"nhooyr.io/websocket"
	"nhooyr.io/websocket/wsjson"
)

type message struct {
	Event string `json:"event"`
	Data  string `json:"data"`
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
	defer c.Close(websocket.StatusAbnormalClosure, "an error occurred")

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	var candidateCh = make(chan string, 4)
	defer close(candidateCh)

	for {
		var msg message
		if err := wsjson.Read(ctx, c, &msg); err != nil {
			log.Err(err).Msg("could not read message")
			return
		}
		log.Info().Str("event", msg.Event).Str("data", msg.Data).Msg("got message")

		switch msg.Event {
		case "offer":
			var sdp webrtc.SessionDescription
			if err := json.Unmarshal([]byte(msg.Data), &sdp); err != nil {
				log.Err(err).Msg("could not unmarshal offer")
				return
			}

			// Initiate webrtc peer connection.
			createPeerConnection(
				&sdp,
				candidateCh,
				signalPeer(ctx, c),
			)
		case "candidate":
			candidateCh <- msg.Data
		default:
			log.Info().Msg("unknown event")
		}
	}
}

// signalPeer sends offer or candidate to peer.
func signalPeer(ctx context.Context, conn *websocket.Conn) SignalFunc {
	return func(msg interface{}, event string) error {
		switch event {
		case "answer":
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
		case "candidate":
			if err := wsjson.Write(ctx, conn, &message{
				Event: event,
				Data:  msg.(string),
			}); err != nil {
				return err
			}
		default:
			log.Info().Msg("unknown event")
		}
		log.Info().Str("event", event).Msg("signaled peer")
		return nil
	}
}
