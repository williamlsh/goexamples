package cmd

import (
	"net"
	"strconv"

	"github.com/spf13/cobra"
	"github.com/williamlsh/goexamples/pkg/log"
	"github.com/williamlsh/goexamples/pkg/tracing"
	"github.com/williamlsh/goexamples/services/driver"
	"go.uber.org/zap"
)

// driverCmd represents the driver command
var driverCmd = &cobra.Command{
	Use:   "driver",
	Short: "Starts Driver service",
	Long:  `Starts Driver service.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		zapLogger := logger.With(zap.String("service", "driver"))
		logger := log.NewFactory(zapLogger)
		server := driver.NewServer(
			net.JoinHostPort("0.0.0.0", strconv.Itoa(driverPort)),
			tracing.Init("driver", metricsFactory, logger),
			metricsFactory,
			logger,
		)
		return logError(zapLogger, server.Run())
	},
}

func init() {
	RootCmd.AddCommand(driverCmd)
}
