package main

import (
	"github.com/haivision/srtgo"
	"github.com/rs/zerolog/log"
	"github.com/williamlsh/logging"
)

func init() {
	logging.Debug(true)
}

var allowedStreamIDs = map[string]bool{
	"foo":    true,
	"foobar": true,
}

func main() {
	srtgo.SrtSetLogLevel(srtgo.SrtLogLevelDebug)

	if err := run(); err != nil {
		panic(err)
	}
}

func run() error {
	options := map[string]string{
		"blocking":  "0",
		"transtype": "file",
		"latency":   "300",
	}

	hostname := "0.0.0.0"
	port := 8090

	log.Info().Str("hostname", hostname).Int("port", port).Msg("Listening")

	ss := srtgo.NewSrtSocket(hostname, uint16(port), options)
	if err := ss.Listen(1); err != nil {
		log.Err(err).Msg("failed to listen")
		return err
	}
	defer ss.Close()

	// ss.SetListenCallback(listenCallback)

	for {
		if err := handleSRT(ss); err != nil {
			return err
		}
	}
}

// func listenCallback(socket *srtgo.SrtSocket, version int, addr *net.UDPAddr, streamID string) bool {
// 	log.Info().Int("version", version).Str("stream_id", streamID).Msg("socket will connect")

// 	// socket not in allowed ids -> reject
// 	if _, found := allowedStreamIDs[streamID]; !found {
// 		// set custom reject reason
// 		socket.SetRejectReason(srtgo.RejectionReasonUnauthorized)
// 		return false
// 	}

// 	// allow connection
// 	return true
// }

func handleSRT(ss *srtgo.SrtSocket) error {
	s, _, err := ss.Accept()
	if err != nil {
		return err
	}
	defer s.Close()

	buf := make([]byte, 1500)
	for {
		n, err := s.Read(buf)
		if err != nil {
			return err
		}
		if n == 0 {
			break
		}

		_, err = s.Write(buf[:n])
		if err != nil {
			return err
		}
	}

	return nil
}
