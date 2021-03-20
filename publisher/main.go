// +build !js

package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/deepch/vdk/av"
	"github.com/deepch/vdk/codec/h264parser"
	"github.com/deepch/vdk/format/rtspv2"
	"github.com/pion/webrtc/v3"
	"github.com/pion/webrtc/v3/pkg/media"
)

func main() {
	peerConnection, err := webrtc.NewPeerConnection(webrtc.Configuration{})
	if err != nil {
		panic(err)
	}

	// Create a video track
	videoTrack, err := webrtc.NewTrackLocalStaticSample(webrtc.RTPCodecCapability{MimeType: webrtc.MimeTypeH264}, "video", "pion")
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

	// Set the handler for ICE connection state
	// This will notify you when the peer has connected/disconnected
	peerConnection.OnICEConnectionStateChange(func(connectionState webrtc.ICEConnectionState) {
		fmt.Printf("Connection State has changed %s \n", connectionState.String())
	})

	offer, err := peerConnection.CreateOffer(nil)
	if err != nil {
		panic(err)
	}

	gatherComplete := webrtc.GatheringCompletePromise(peerConnection)
	if err = peerConnection.SetLocalDescription(offer); err != nil {
		panic(err)
	}
	<-gatherComplete

	offerJSON, err := json.Marshal(*peerConnection.LocalDescription())
	if err != nil {
		panic(err)
	}

	response, err := http.Post("http://localhost:8080/createPublisher", "application/json", bytes.NewBuffer(offerJSON))
	if err != nil {
		panic(err)
	}
	defer response.Body.Close()

	answer := webrtc.SessionDescription{}
	if err = json.NewDecoder(response.Body).Decode(&answer); err != nil {
		panic(err)
	}

	if err = peerConnection.SetRemoteDescription(answer); err != nil {
		panic(err)
	}

	if err := RTSPConsume(videoTrack); err != nil {
		panic(err)
	}
}

// Connect to an RTSP URL and pull media.
// Convert H264 to Annex-B, then write to outboundVideoTrack which sends to all PeerConnections
func RTSPConsume(videoTrack *webrtc.TrackLocalStaticSample) error {
	annexbNALUStartCode := func() []byte { return []byte{0x00, 0x00, 0x00, 0x01} }

	// Use a loop incase RTSP stream is stopped os we can retry.
	for {
		session, err := rtspv2.Dial(rtspv2.RTSPClientOptions{
			URL:              "rtsp://123.0.0.1",
			Debug:            false,
			DialTimeout:      3 * time.Second,
			ReadWriteTimeout: 3 * time.Second,
			DisableAudio:     true,
		})
		if err != nil {
			panic(err)
		}

		codecs := session.CodecData
		for i, t := range codecs {
			log.Println("Stream", i, "is of type", t.Type().String())
		}
		if codecs[0].Type() != av.H264 {
			panic("RTSP feed must begin with a H264 codec")
		}
		if len(codecs) != 1 {
			log.Println("Ignoring all but the first stream.")
		}

		var previousTime time.Duration
		for {
			pkt := <-session.OutgoingPacketQueue

			if pkt.Idx != 0 {
				//audio or other stream, skip it
				continue
			}

			pkt.Data = pkt.Data[4:]

			// For every key-frame pre-pend the SPS and PPS
			if pkt.IsKeyFrame {
				pkt.Data = append(annexbNALUStartCode(), pkt.Data...)
				pkt.Data = append(codecs[0].(h264parser.CodecData).PPS(), pkt.Data...)
				pkt.Data = append(annexbNALUStartCode(), pkt.Data...)
				pkt.Data = append(codecs[0].(h264parser.CodecData).SPS(), pkt.Data...)
				pkt.Data = append(annexbNALUStartCode(), pkt.Data...)
			}

			bufferDuration := pkt.Time - previousTime
			previousTime = pkt.Time
			if err = videoTrack.WriteSample(media.Sample{Data: pkt.Data, Duration: bufferDuration}); err != nil && err != io.ErrClosedPipe {
				panic(err)
			}
		}
	}
}
