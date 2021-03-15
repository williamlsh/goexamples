package main

import (
	"context"

	"github.com/rs/zerolog/log"
)

func main() {
	logger := log.With().Str("component", "module").Logger()
	ctx := logger.WithContext(context.Background())

	log.Ctx(ctx).Info().Msg("hello world")
}
