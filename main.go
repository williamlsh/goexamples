package main

import (
	"errors"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func main() {
	// UNIX Time is faster and smaller than most timestamps
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix

	err := errors.New("a repo man spends his life getting into tense situations")
	service := "myservice"

	log.Fatal().Err(err).Str("service", service).Msgf("Cannot start %s", service)
}
