package main

import (
	"net/http"

	"github.com/rs/zerolog/log"
	"github.com/williamlsh/logging"
)

func init() {
	logging.Debug(true)
}

func main() {
	log.Fatal().AnErr("err", http.ListenAndServe(":8080", serveMux())).Send()
}
