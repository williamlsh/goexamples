package driver

import (
	"context"
	"net"

	otgrpc "github.com/opentracing-contrib/go-grpc"
	"github.com/opentracing/opentracing-go"
	"github.com/uber/jaeger-lib/metrics"
	"github.com/williamlsh/goexamples/pkg/log"
	"github.com/williamlsh/goexamples/services/driver/proto"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

// Server implements jaeger-demo-frontend service
type Server struct {
	proto.UnimplementedDriverServiceServer
	hostPort string
	tracer   opentracing.Tracer
	logger   log.Factory
	redis    *Redis
	server   *grpc.Server
}

var _ proto.DriverServiceServer = (*Server)(nil)

// NewServer creates a new driver.Server
func NewServer(hostPort string, tracer opentracing.Tracer, metricsFactory metrics.Factory, logger log.Factory) *Server {
	server := grpc.NewServer(
		grpc.UnaryInterceptor(otgrpc.OpenTracingServerInterceptor(tracer)),
		grpc.StreamInterceptor(otgrpc.OpenTracingStreamServerInterceptor(tracer)),
	)
	return &Server{
		hostPort: hostPort,
		tracer:   tracer,
		logger:   logger,
		server:   server,
		redis:    newRedis(metricsFactory, logger),
	}
}

// Run starts the Driver server
func (s *Server) Run() error {
	lis, err := net.Listen("tcp", s.hostPort)
	if err != nil {
		s.logger.Bg().Fatal("Unable to create http listener", zap.Error(err))
	}
	proto.RegisterDriverServiceServer(s.server, s)
	err = s.server.Serve(lis)
	if err != nil {
		s.logger.Bg().Fatal("Unable to start gRPC server", zap.Error(err))
	}
	return nil
}

// FindNearest implements gRPC driver interface
func (s *Server) FindNearest(ctx context.Context, location *proto.DriverLocationRequest) (*proto.DriverLocationResponse, error) {
	s.logger.For(ctx).Info("Searching for nearest drivers", zap.String("location", location.Location))
	driverIDs := s.redis.FindDriverIDs(ctx, location.Location)

	retMe := make([]*proto.DriverLocation, len(driverIDs))
	for i, driverID := range driverIDs {
		var drv Driver
		var err error
		for i := 0; i < 3; i++ {
			drv, err = s.redis.GetDriver(ctx, driverID)
			if err == nil {
				break
			}
			s.logger.For(ctx).Error("Retrying GetDriver after error", zap.Int("retry_no", i+1), zap.Error(err))
		}
		if err != nil {
			s.logger.For(ctx).Error("Failed to get driver after 3 attempts", zap.Error(err))
			return nil, err
		}
		retMe[i] = &proto.DriverLocation{
			DriverID: drv.DriverID,
			Location: drv.Location,
		}
	}

	s.logger.For(ctx).Info("Search successful", zap.Int("num_drivers", len(retMe)))
	return &proto.DriverLocationResponse{Locations: retMe}, nil
}
