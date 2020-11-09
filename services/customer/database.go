package customer

import (
	"context"
	"errors"

	"github.com/opentracing/opentracing-go"
	tags "github.com/opentracing/opentracing-go/ext"
	"github.com/williamlsh/goexamples/pkg/delay"
	"github.com/williamlsh/goexamples/pkg/log"
	"github.com/williamlsh/goexamples/pkg/tracing"
	"github.com/williamlsh/goexamples/services/config"
	"go.uber.org/zap"
)

// database simulates Customer repository implemented on top of an SQL database
type database struct {
	tracer    opentracing.Tracer
	logger    log.Factory
	customers map[string]*Customer
	lock      *tracing.Mutex
}

func newDatabase(tracer opentracing.Tracer, logger log.Factory) *database {
	return &database{
		tracer: tracer,
		logger: logger,
		lock: &tracing.Mutex{
			SessionBaggageKey: "request",
		},
		customers: map[string]*Customer{
			"123": {
				ID:       "123",
				Name:     "Rachel's Floral Designs",
				Location: "115,277",
			},
			"567": {
				ID:       "567",
				Name:     "Amazing Coffee Roasters",
				Location: "211,653",
			},
			"392": {
				ID:       "392",
				Name:     "Trom Chocolatier",
				Location: "577,322",
			},
			"731": {
				ID:       "731",
				Name:     "Japanese Desserts",
				Location: "728,326",
			},
		},
	}
}

func (d *database) Get(ctx context.Context, customerID string) (*Customer, error) {
	d.logger.For(ctx).Info("Loading customer", zap.String("customer_id", customerID))

	// simulate opentracing instrumentation of an SQL query
	if span := opentracing.SpanFromContext(ctx); span != nil {
		span := d.tracer.StartSpan("SQL SELECT", opentracing.ChildOf(span.Context()))
		tags.SpanKindRPCClient.Set(span)
		tags.PeerService.Set(span, "mysql")
		// #nosec
		span.SetTag("sql.query", "SELECT * FROM customer WHERE customer_id="+customerID)
		defer span.Finish()
		ctx = opentracing.ContextWithSpan(ctx, span)
	}

	if !config.MySQLMutexDisabled {
		// simulate misconfigured connection pool that only gives one connection at a time
		d.lock.Lock(ctx)
		defer d.lock.Unlock()
	}

	// simulate RPC delay
	delay.Sleep(config.MySQLGetDelay, config.MySQLGetDelayStdDev)

	if customer, ok := d.customers[customerID]; ok {
		return customer, nil
	}
	return nil, errors.New("invalid customer ID")
}
