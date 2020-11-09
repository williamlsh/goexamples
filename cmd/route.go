package cmd

import (
	"net"
	"strconv"

	"github.com/spf13/cobra"
	"github.com/williamlsh/goexamples/pkg/log"
	"github.com/williamlsh/goexamples/pkg/tracing"
	"github.com/williamlsh/goexamples/services/route"
	"go.uber.org/zap"
)

// routeCmd represents the route command
var routeCmd = &cobra.Command{
	Use:   "route",
	Short: "Starts Route service",
	Long:  `Starts Route service.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		zapLogger := logger.With(zap.String("service", "route"))
		logger := log.NewFactory(zapLogger)
		server := route.NewServer(
			net.JoinHostPort("0.0.0.0", strconv.Itoa(routePort)),
			tracing.Init("route", metricsFactory, logger),
			logger,
		)
		return logError(zapLogger, server.Run())
	},
}

func init() {
	RootCmd.AddCommand(routeCmd)
}
