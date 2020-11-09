package customer

import (
	"context"
	"fmt"
	"net/http"

	"github.com/opentracing-contrib/go-stdlib/nethttp"
	"github.com/opentracing/opentracing-go"
	"github.com/williamlsh/goexamples/pkg/log"
	"github.com/williamlsh/goexamples/pkg/tracing"
	"go.uber.org/zap"
)

// Client is a remote client that implements customer.Interface
type Client struct {
	tracer   opentracing.Tracer
	logger   log.Factory
	client   *tracing.HTTPClient
	hostPort string
}

// NewClient creates a new customer.Client.
func NewClient(tracer opentracing.Tracer, logger log.Factory, hostPort string) *Client {
	return &Client{
		tracer: tracer,
		logger: logger,
		client: &tracing.HTTPClient{
			Client: &http.Client{Transport: &nethttp.Transport{}},
			Tracer: tracer,
		},
		hostPort: hostPort,
	}
}

// Get implements customer.Interface#Get as an RPC.
func (c *Client) Get(ctx context.Context, customerID string) (*Customer, error) {
	c.logger.For(ctx).Info("Getting customer", zap.String("customer_id", customerID))

	url := fmt.Sprintf("http://"+c.hostPort+"/customer?customer=%s", customerID)
	fmt.Println(url)
	var customer Customer
	if err := c.client.GetJSON(ctx, "/customer", url, &customer); err != nil {
		return nil, err
	}
	return &customer, nil
}
