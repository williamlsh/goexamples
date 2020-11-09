package driver

import (
	"context"
	"errors"
	"fmt"
	"math/rand"
	"sync"

	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	"github.com/uber/jaeger-lib/metrics"
	"github.com/williamlsh/goexamples/pkg/delay"
	"github.com/williamlsh/goexamples/pkg/log"
	"github.com/williamlsh/goexamples/pkg/tracing"
	"github.com/williamlsh/goexamples/services/config"
	"go.uber.org/zap"
)

// Redis is a simulator of remote Redis cache.
type Redis struct {
	tracer opentracing.Tracer // simulate redis as a separate process
	logger log.Factory
	errorSimulator
}

func newRedis(metricsFactory metrics.Factory, logger log.Factory) *Redis {
	return &Redis{
		tracer: tracing.Init("redis", metricsFactory, logger),
		logger: logger,
	}
}

// FindDriverIDs finds IDs of drivers who are near the location.
func (r *Redis) FindDriverIDs(ctx context.Context, location string) []string {
	if span := opentracing.SpanFromContext(ctx); span != nil {
		span := r.tracer.StartSpan("FindDriverIDs", opentracing.ChildOf(span.Context()))
		span.SetTag("param.location", location)
		ext.SpanKindRPCClient.Set(span)
		defer span.Finish()
		ctx = opentracing.ContextWithSpan(ctx, span)
	}
	// simulate RPC delay
	delay.Sleep(config.RedisFindDelay, config.RedisFindDelayStdDev)

	drivers := make([]string, 10)
	for i := range drivers {
		// #nosec
		drivers[i] = fmt.Sprintf("T7%05dC", rand.Int()%100000)
	}
	r.logger.For(ctx).Info("Found drivers", zap.Strings("drivers", drivers))
	return drivers
}

// GetDriver returns driver and the current car location
func (r *Redis) GetDriver(ctx context.Context, driverID string) (Driver, error) {
	if span := opentracing.SpanFromContext(ctx); span != nil {
		span := r.tracer.StartSpan("GetDriver", opentracing.ChildOf(span.Context()))
		span.SetTag("param.driverID", driverID)
		ext.SpanKindRPCClient.Set(span)
		defer span.Finish()
		ctx = opentracing.ContextWithSpan(ctx, span)
	}
	// simulate RPC delay
	delay.Sleep(config.RedisGetDelay, config.RedisGetDelayStdDev)
	if err := r.checkError(); err != nil {
		if span := opentracing.SpanFromContext(ctx); span != nil {
			ext.Error.Set(span, true)
		}
		r.logger.For(ctx).Error("redis timeout", zap.String("driver_id", driverID), zap.Error(err))
		return Driver{}, err
	}

	// #nosec
	return Driver{
		DriverID: driverID,
		Location: fmt.Sprintf("%d,%d", rand.Int()%1000, rand.Int()%1000),
	}, nil
}

var errTimeout = errors.New("redis timeout")

type errorSimulator struct {
	sync.Mutex
	countTillError int
}

func (es *errorSimulator) checkError() error {
	es.Lock()
	es.countTillError--
	if es.countTillError > 0 {
		es.Unlock()
		return nil
	}
	es.countTillError = 5
	es.Unlock()
	delay.Sleep(2*config.RedisFindDelay, 0) // add more delay for "timeout"
	return errTimeout
}
