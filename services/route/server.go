package route

import (
	"context"
	"encoding/json"
	"expvar"
	"math"
	"math/rand"
	"net/http"
	"time"

	"github.com/opentracing/opentracing-go"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/williamlsh/goexamples/pkg/delay"
	"github.com/williamlsh/goexamples/pkg/httperr"
	"github.com/williamlsh/goexamples/pkg/log"
	"github.com/williamlsh/goexamples/pkg/tracing"
	"github.com/williamlsh/goexamples/services/config"
	"go.uber.org/zap"
)

// Server implements Route service
type Server struct {
	hostPort string
	tracer   opentracing.Tracer
	logger   log.Factory
}

// NewServer creates a new route.Server
func NewServer(hostPort string, tracer opentracing.Tracer, logger log.Factory) *Server {
	return &Server{
		hostPort: hostPort,
		tracer:   tracer,
		logger:   logger,
	}
}

// Run starts the Route server
func (s *Server) Run() error {
	mux := s.createServeMux()
	s.logger.Bg().Info("Starting", zap.String("address", "http://"+s.hostPort))
	return http.ListenAndServe(s.hostPort, mux)
}

func (s *Server) createServeMux() http.Handler {
	mux := tracing.NewServeMux(s.tracer)
	mux.Handle("/route", http.HandlerFunc(s.route))
	mux.Handle("/debug/vars", expvar.Handler()) // expvar
	mux.Handle("/metrics", promhttp.Handler())  // Prometheus
	return mux
}

func (s *Server) route(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	s.logger.For(ctx).Info("HTTP request received", zap.String("method", r.Method), zap.Stringer("url", r.URL))
	if err := r.ParseForm(); httperr.HandleError(w, err, http.StatusBadRequest) {
		s.logger.For(ctx).Error("bad request", zap.Error(err))
		return
	}

	pickup := r.Form.Get("pickup")
	if pickup == "" {
		http.Error(w, "Missing required 'pickup' parameter", http.StatusBadRequest)
		return
	}

	dropoff := r.Form.Get("dropoff")
	if dropoff == "" {
		http.Error(w, "Missing required 'dropoff' parameter", http.StatusBadRequest)
		return
	}

	response := computeRoute(ctx, pickup, dropoff)

	data, err := json.Marshal(response)
	if httperr.HandleError(w, err, http.StatusInternalServerError) {
		s.logger.For(ctx).Error("cannot marshal response", zap.Error(err))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(data)
}

func computeRoute(ctx context.Context, pickup, dropoff string) *Route {
	start := time.Now()
	defer func() {
		updateCalcStats(ctx, time.Since(start))
	}()

	// Simulate expensive calculation
	delay.Sleep(config.RouteCalcDelay, config.RouteCalcDelayStdDev)

	eta := math.Max(2, rand.NormFloat64()*3+5)
	return &Route{
		Pickup:  pickup,
		Dropoff: dropoff,
		ETA:     time.Duration(eta) * time.Minute,
	}
}
