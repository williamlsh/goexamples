// Reference: https://medium.com/@dlagoza/playing-with-multiple-contexts-in-go-9f72cbcff56e

package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"time"

	"github.com/justinas/alice"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/hlog"
	"github.com/rs/zerolog/log"
)

func main() {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: time.RFC3339})

	handler := alice.New(
		hlog.NewHandler(log.Logger),
		hlog.URLHandler("url"),
		hlog.RemoteAddrHandler("ip"),
		hlog.UserAgentHandler("user_agent"),
		hlog.RefererHandler("referer"),
	).ThenFunc(TaskHandler)

	srv := http.Server{
		Addr:    "8080",
		Handler: handler,
	}

	go func() {
		log.Info().Msg("Listening HTTP on: 8080")
		if err := srv.ListenAndServe(); err != nil {
			log.Fatal().Err(err).Msg("error when running http server")
		}
	}()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, os.Kill)
	<-c

	// We received an interrupt signal, shut down.
	log.Info().Msg("Shutting down...")
	ctx, _ := context.WithTimeout(context.Background(), time.Second*10)
	if err := srv.Shutdown(ctx); err != nil {
		log.Error().Err(err).Msg("HTTP server shutdown error")
	}
	log.Info().Msg("Server has been stopped")
}

func TaskHandler(w http.ResponseWriter, r *http.Request) {
	l := hlog.FromRequest(r)
	l.Info().Msgf("%s %s", r.Method, r.URL.RequestURI())
	t := r.FormValue("time")
	if t == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	duration, err := strconv.Atoi(t)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	go taskManager(l.WithContext(context.Background()), duration)
	w.Write([]byte("started"))
}

func taskManager(ctx context.Context, duration int) {
	ctx, cancell := context.WithTimeout(ctx, 1*time.Minute)
	defer cancell()
	task(ctx, duration)
}

func task(ctx context.Context, duration int) {
	l := log.Ctx(ctx)
	l.Info().Msgf("Task %d second(s): STARTED", duration)
	select {
	case <-ctx.Done():
		l.Info().Msgf("Task %d second(s): CANCELED", duration)
	case <-time.After(time.Duration(duration) * time.Second):
		l.Info().Msgf("Task %d second(s): FINISHED", duration)
	}
}
