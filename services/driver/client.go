package driver

import (
	"context"
	"time"

	otgrpc "github.com/opentracing-contrib/go-grpc"
	"github.com/opentracing/opentracing-go"
	"github.com/williamlsh/goexamples/pkg/log"
	"github.com/williamlsh/goexamples/services/driver/proto"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

// Client is a remote client that implements driver.Interface
type Client struct {
	tracer opentracing.Tracer
	logger log.Factory
	client proto.DriverServiceClient
}

// NewClient creates a new driver.Client
func NewClient(tracer opentracing.Tracer, logger log.Factory, hostPort string) *Client {
	conn, err := grpc.Dial(
		hostPort,
		grpc.WithInsecure(),
		grpc.WithChainUnaryInterceptor(otgrpc.OpenTracingClientInterceptor(tracer)),
		grpc.WithStreamInterceptor(otgrpc.OpenTracingStreamClientInterceptor(tracer)),
	)
	if err != nil {
		logger.Bg().Fatal("Cannot create gRPC connection", zap.Error(err))
	}

	client := proto.NewDriverServiceClient(conn)
	return &Client{
		tracer: tracer,
		logger: logger,
		client: client,
	}
}

// FindNearest implements driver.Interface#FindNearest as an RPC
func (c *Client) FindNearest(ctx context.Context, location string) ([]Driver, error) {
	c.logger.For(ctx).Info("Finding nearest drivers", zap.String("location", location))
	ctx, cancel := context.WithTimeout(ctx, 1*time.Second)
	defer cancel()

	response, err := c.client.FindNearest(ctx, &proto.DriverLocationRequest{Location: location})
	if err != nil {
		return nil, err
	}

	return fromProto(response), nil
}

func fromProto(response *proto.DriverLocationResponse) []Driver {
	retMe := make([]Driver, len(response.Locations))
	for i, result := range response.Locations {
		retMe[i] = Driver{
			DriverID: result.DriverID,
			Location: result.Location,
		}
	}
	return retMe
}
