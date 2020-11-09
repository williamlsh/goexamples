package cmd

import (
	"net"
	"strconv"

	"github.com/spf13/cobra"
	"github.com/williamlsh/goexamples/pkg/log"
	"github.com/williamlsh/goexamples/pkg/tracing"
	"github.com/williamlsh/goexamples/services/frontend"
	"go.uber.org/zap"
)

// frontendCmd represents the frontend command
var frontendCmd = &cobra.Command{
	Use:   "frontend",
	Short: "Starts Frontend service",
	Long:  `Starts Frontend service.`,
	RunE: func(cmd *cobra.Command, args []string) error {

		options.FrontendHostPort = net.JoinHostPort("0.0.0.0", strconv.Itoa(frontendPort))
		options.DriverHostPort = net.JoinHostPort("0.0.0.0", strconv.Itoa(driverPort))
		options.CustomerHostPort = net.JoinHostPort("0.0.0.0", strconv.Itoa(customerPort))
		options.RouteHostPort = net.JoinHostPort("0.0.0.0", strconv.Itoa(routePort))
		options.Basepath = basepath
		options.JaegerUI = jaegerUI

		zapLogger := logger.With(zap.String("service", "frontend"))
		logger := log.NewFactory(zapLogger)
		server := frontend.NewServer(
			options,
			tracing.Init("frontend", metricsFactory, logger),
			logger,
		)
		return logError(zapLogger, server.Run())
	},
}

var options frontend.ConfigOptions

func init() {
	RootCmd.AddCommand(frontendCmd)

}
