package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"net"
	"net/url"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/pion/webrtc/v3"
	"github.com/sourcegraph/jsonrpc2"
)

type Candidate struct {
	Target    int                  `json:"target"`
	Candidate *webrtc.ICECandidate `json:"candidate"`
}

type ResponseCandidate struct {
	Target    int                      `json:"target"`
	Candidate *webrtc.ICECandidateInit `json:"candidate"`
}

// SendOffer object to send to the sfu over Websockets
type SendOffer struct {
	SID   string                     `json:"sid"`
	Offer *webrtc.SessionDescription `json:"offer"`
}

// SendAnswer object to send to the sfu over Websockets
type SendAnswer struct {
	SID    string                     `json:"sid"`
	Answer *webrtc.SessionDescription `json:"answer"`
}

// TrickleResponse received from the sfu server
type TrickleResponse struct {
	Params ResponseCandidate `json:"params"`
	Method string            `json:"method"`
}

// Response received from the sfu over Websockets
type Response struct {
	Params *webrtc.SessionDescription `json:"params"`
	Result *webrtc.SessionDescription `json:"result"`
	Method string                     `json:"method"`
	Id     uint64                     `json:"id"`
}

var peerConnection *webrtc.PeerConnection
var connectionID uint64

var addr string

func main() {
	flag.StringVar(&addr, "a", "localhost:7000", "address to use")
	flag.Parse()

	u := url.URL{Scheme: "ws", Host: addr, Path: "/ws"}
	log.Printf("connecting to %s", u.String())

	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		log.Fatal("dial:", err)
	}
	defer c.Close()

	peerConnection, err := webrtc.NewPeerConnection(webrtc.Configuration{
		ICEServers: []webrtc.ICEServer{
			{
				URLs: []string{"stun:stun.l.google.com:19302"},
			},
		},
	})
	if err != nil {
		panic(err)
	}

	// Read incoming Websocket messages
	done := make(chan struct{})

	go readMessage(c, done)

	// Open a UDP Listener for RTP Packets on port 5004
	listener, err := net.ListenUDP("udp", &net.UDPAddr{IP: net.ParseIP("127.0.0.1"), Port: 9004})
	if err != nil {
		panic(err)
	}
	defer func() {
		if err = listener.Close(); err != nil {
			panic(err)
		}
	}()

	// Create a video track
	videoTrack, err := webrtc.NewTrackLocalStaticRTP(webrtc.RTPCodecCapability{MimeType: webrtc.MimeTypeVP8}, "video", "pion")
	if err != nil {
		panic(err)
	}
	rtpSender, err := peerConnection.AddTrack(videoTrack)
	if err != nil {
		panic(err)
	}

	// Read incoming RTCP packets
	// Before these packets are returned they are processed by interceptors. For things
	// like NACK this needs to be called.
	go func() {
		rtcpBuf := make([]byte, 1500)
		for {
			if _, _, rtcpErr := rtpSender.Read(rtcpBuf); rtcpErr != nil {
				return
			}
		}
	}()

	// Creating WebRTC offer
	offer, err := peerConnection.CreateOffer(nil)
	if err != nil {
		panic(err)
	}

	// Set the remote SessionDescription
	err = peerConnection.SetLocalDescription(offer)
	if err != nil {
		panic(err)
	}

	// Handling OnICECandidate event
	peerConnection.OnICECandidate(func(candidate *webrtc.ICECandidate) {
		if candidate != nil {
			candidateJSON, err := json.Marshal(&Candidate{
				Candidate: candidate,
				Target:    0,
			})

			params := (*json.RawMessage)(&candidateJSON)

			if err != nil {
				log.Fatal(err)
			}

			message := &jsonrpc2.Request{
				Method: "trickle",
				Params: params,
			}

			reqBodyBytes := new(bytes.Buffer)
			json.NewEncoder(reqBodyBytes).Encode(message)

			messageBytes := reqBodyBytes.Bytes()
			c.WriteMessage(websocket.TextMessage, messageBytes)
		}
	})

	peerConnection.OnICEConnectionStateChange(func(connectionState webrtc.ICEConnectionState) {
		fmt.Printf("Connection State has changed to %s \n", connectionState.String())
	})

	offerJSON, err := json.Marshal(&SendOffer{
		Offer: peerConnection.LocalDescription(),
		SID:   "test room",
	})
	if err != nil {
		panic(err)
	}

	params := (*json.RawMessage)(&offerJSON)

	connectionUUID := uuid.New()
	connectionID = uint64(connectionUUID.ID())

	offerMessage := &jsonrpc2.Request{
		Method: "join",
		Params: params,
		ID: jsonrpc2.ID{
			IsString: false,
			Str:      "",
			Num:      connectionID,
		},
	}

	reqBodyBytes := new(bytes.Buffer)
	json.NewEncoder(reqBodyBytes).Encode(offerMessage)

	messageBytes := reqBodyBytes.Bytes()
	c.WriteMessage(websocket.TextMessage, messageBytes)

	go func() {
		// Read RTP packets forever and send them to the WebRTC Client
		inboundRTPPacket := make([]byte, 1600) // UDP MTU
		for {
			n, _, err := listener.ReadFrom(inboundRTPPacket)
			if err != nil {
				panic(fmt.Sprintf("error during read: %s", err))
			}

			if _, err = videoTrack.Write(inboundRTPPacket[:n]); err != nil {
				panic(err)
			}
		}
	}()

	<-done
}

func readMessage(connection *websocket.Conn, done chan struct{}) {
	defer close(done)
	for {
		_, message, err := connection.ReadMessage()
		if err != nil || err == io.EOF {
			log.Fatal("Error reading: ", err)
			break
		}

		fmt.Printf("recv: %s", message)

		var response Response
		json.Unmarshal(message, &response)

		if response.Id == connectionID {
			result := *response.Result
			if err := peerConnection.SetRemoteDescription(result); err != nil {
				log.Fatal(err)
			}
		} else if response.Method == "trickle" {
			var trickleResponse TrickleResponse

			if err := json.Unmarshal(message, &trickleResponse); err != nil {
				log.Fatal(err)
			}

			err := peerConnection.AddICECandidate(*trickleResponse.Params.Candidate)

			if err != nil {
				log.Fatal(err)
			}
		}
	}
}
