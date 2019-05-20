package addendpoint

import (
	"context"
	"time"

	"goexamples/pkg/addservice"

	"github.com/go-kit/kit/circuitbreaker"
	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/metrics"
	"github.com/go-kit/kit/ratelimit"
	"github.com/go-kit/kit/tracing/opentracing"
	"github.com/go-kit/kit/tracing/zipkin"
	stdopentracing "github.com/opentracing/opentracing-go"
	stdzipkin "github.com/openzipkin/zipkin-go"
	"github.com/sony/gobreaker"
	"golang.org/x/time/rate"
)

// Set collects all of the endpoints that compose an add service. It's meant to
// be used as a helper struct, to collect all of the endpoints into a single
// parameter.
type Set struct {
	SumEndpoint    endpoint.Endpoint
	ConcatEndpoint endpoint.Endpoint
}

// New returns a Set that wraps the provided server, and wires in all of the
// expected endpoint middlewares via the various parameters.
func New(svc addservice.Service, logger log.Logger, duration metrics.Histogram, otTracer stdopentracing.Tracer, zipkinTracer *stdzipkin.Tracer) Set {
	var sumEndpoint endpoint.Endpoint
	{
		sumEndpoint = MakeSumEndpoint(svc)
		sumEndpoint = ratelimit.NewErroringLimiter(rate.NewLimiter(rate.Every(time.Second), 1))(sumEndpoint)
		sumEndpoint = circuitbreaker.Gobreaker(gobreaker.NewCircuitBreaker(gobreaker.Settings{}))(sumEndpoint)
		sumEndpoint = opentracing.TraceServer(otTracer, "sum")(sumEndpoint)
		sumEndpoint = zipkin.TraceEndpoint(zipkinTracer, "sum")(sumEndpoint)
		sumEndpoint = LoggingMiddleware(log.With(logger, "method", "sum"))(sumEndpoint)
		sumEndpoint = InstrumentingMiddleware(duration.With("method", "sum"))(sumEndpoint)
	}
	var concatEndpoint endpoint.Endpoint
	{
		concatEndpoint = MakeConcatEndpoint(svc)
		concatEndpoint = ratelimit.NewErroringLimiter(rate.NewLimiter(rate.Every(time.Second), 100))(concatEndpoint)
		concatEndpoint = circuitbreaker.Gobreaker(gobreaker.NewCircuitBreaker(gobreaker.Settings{}))(concatEndpoint)
		concatEndpoint = opentracing.TraceServer(otTracer, "Concat")(concatEndpoint)
		concatEndpoint = zipkin.TraceEndpoint(zipkinTracer, "Concat")(concatEndpoint)
		concatEndpoint = LoggingMiddleware(log.With(logger, "method", "Concat"))(concatEndpoint)
		concatEndpoint = InstrumentingMiddleware(duration.With("method", "Concat"))(concatEndpoint)
	}
	return Set{
		SumEndpoint:    sumEndpoint,
		ConcatEndpoint: concatEndpoint,
	}
}

// MakeSumEndpoint constructs a Sum endpoint wrapping the service.
func MakeSumEndpoint(s addservice.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(SumRequest)
		i, err := s.Sum(ctx, req.A, req.B)
		return SumResponse{V: i, Err: err}, nil
	}
}

// MakeConcatEndpoint constructs a Concat endpoint wrapping the service.
func MakeConcatEndpoint(s addservice.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(ConcatRequest)
		string1, err := s.Concat(ctx, req.A, req.B)
		return ConcatResponse{V: string1, Err: err}, nil
	}
}

// compile time assertions for our response types implementing endpoint.Failer.
var (
	_ endpoint.Failer = SumResponse{}
	_ endpoint.Failer = ConcatResponse{}
)

type SumRequest struct {
	A, B int
}

type SumResponse struct {
	V   int   `json:"v"`
	Err error `json:"-"` // should be intercepted by Failed/errorEncoder
}

// Failed implements endpoint.Failer.
func (r SumResponse) Failed() error { return r.Err }

type ConcatRequest struct {
	A, B string
}

type ConcatResponse struct {
	V   string `json:"v"`
	Err error  `json:"-"` // should be intercepted by Failed/errorEncoder
}

// Failed implements endpoint.Failer.
func (r ConcatResponse) Failed() error { return r.Err }
