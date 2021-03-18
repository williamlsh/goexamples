package main

import (
	"fmt"
	"net"

	"github.com/pion/webrtc/v3"
	"github.com/williamlsh/goexamples/signalling"
)

func main() {
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

	go func() {
		sdpChan := signalling.HTTPSDPServer(8081)

		// Answer from broadcasting peer.
		answer := webrtc.SessionDescription{}
		signalling.Decode(<-sdpChan, &answer)
		fmt.Println("Recv answer from broadcasting peer")

		if peerConnection.SetRemoteDescription(answer); err != nil {
			panic(err)
		}
	}()

	// Open a UDP Listener for RTP Packets on port 5004
	listener, err := net.ListenUDP("udp", &net.UDPAddr{IP: net.ParseIP("127.0.0.1"), Port: 5004})
	if err != nil {
		panic(err)
	}
	defer func() {
		if err = listener.Close(); err != nil {
			panic(err)
		}
	}()
	fmt.Println("Started UDP listener")

	// Create a video track
	videoTrack, err := webrtc.NewTrackLocalStaticRTP(webrtc.RTPCodecCapability{MimeType: webrtc.MimeTypeVP8}, "video", "pion")
	if err != nil {
		panic(err)
	}
	rtpSender, err := peerConnection.AddTrack(videoTrack)
	if err != nil {
		panic(err)
	}
	processRTCP(rtpSender)

	// Set the handler for ICE connection state
	// This will notify you when the peer has connected/disconnected
	peerConnection.OnICEConnectionStateChange(func(connectionState webrtc.ICEConnectionState) {
		fmt.Printf("Connection State has changed %s \n", connectionState.String())
	})

	peerConnection.OnICEGatheringStateChange(func(gathererState webrtc.ICEGathererState) {
		fmt.Printf("Gathering State has changed %s \n", gathererState.String())
	})

	// Wait for the offer to be pasted
	// offer := webrtc.SessionDescription{}
	// signalling.Decode(signalling.MustReadStdin(), &offer)

	// Set the remote SessionDescription
	// if err = peerConnection.SetRemoteDescription(offer); err != nil {
	// 	panic(err)
	// }

	// Create answer
	// answer, err := peerConnection.CreateAnswer(nil)
	// if err != nil {
	// 	panic(err)
	// }

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

	offerStr := signalling.Encode(*peerConnection.LocalDescription())
	signalling.HTTPSDPClient(offerStr, 8080)
	fmt.Println("Sent offer to broadcasting peer")

	// Create channel that is blocked until ICE Gathering is complete
	gatherComplete := webrtc.GatheringCompletePromise(peerConnection)

	// Sets the LocalDescription, and starts our UDP listeners
	// if err = peerConnection.SetLocalDescription(answer); err != nil {
	// 	panic(err)
	// }

	// Block until ICE Gathering is complete, disabling trickle ICE
	// we do this because we only can exchange one signaling message
	// in a production application you should exchange ICE Candidates via OnICECandidate
	<-gatherComplete

	// Output the answer in base64 so we can paste it in browser
	// fmt.Println(signalling.Encode(*peerConnection.LocalDescription()))

	// Read RTP packets forever and send them to the WebRTC Client
	inboundRTPPacket := make([]byte, 1600) // UDP MTU
	fmt.Println("Sending RTP packets")
	for {
		n, _, err := listener.ReadFrom(inboundRTPPacket)
		if err != nil {
			panic(fmt.Sprintf("error during read: %s", err))
		}

		if _, err = videoTrack.Write(inboundRTPPacket[:n]); err != nil {
			panic(err)
		}
	}
}

// Read incoming RTCP packets
// Before these packets are retuned they are processed by interceptors. For things
// like NACK this needs to be called.
func processRTCP(rtpSender *webrtc.RTPSender) {
	go func() {
		rtcpBuf := make([]byte, 1500)

		for {
			if _, _, rtcpErr := rtpSender.Read(rtcpBuf); rtcpErr != nil {
				return
			}
		}
	}()
}
