package main

import (
	"bufio"
	"os"

	"github.com/haivision/srtgo"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/williamlsh/logging"
)

func main() {
	logging.Debug(true)

	srtgo.SrtSetLogLevel(srtgo.SrtLogLevelDebug)

	if err := Ingest(&log.Logger); err != nil {
		panic(err)
	}
}

func Ingest(logger *zerolog.Logger) error {
	options := map[string]string{
		"blocking":  "0",
		"transtype": "file",
	}

	hostname := "0.0.0.0"
	port := 8090

	logger.Info().Str("hostname", hostname).Int("port", port).Msg("Listening")

	srtSocket := srtgo.NewSrtSocket(hostname, uint16(port), options)
	if err := srtSocket.Listen(1); err != nil {
		logger.Err(err).Msg("failed to listen")
		return err
	}
	defer srtSocket.Close()

	for {
		if err := handleSRT(srtSocket); err != nil {
			return err
		}
	}
}

func handleSRT(srtSocket *srtgo.SrtSocket) error {
	s, _, err := srtSocket.Accept()
	if err != nil {
		return err
	}
	defer s.Close()

	f, err := os.Create("sample.ts")
	if err != nil {
		return err
	}
	defer f.Close()

	w := bufio.NewWriter(f)
	defer w.Flush()

	buf := make([]byte, 2048)
	for {
		n, err := s.Read(buf)
		if err != nil {
			return err
		}
		if n == 0 {
			break
		}

		if _, err := w.Write(buf[:n]); err != nil {
			return err
		}
	}

	return nil
}
