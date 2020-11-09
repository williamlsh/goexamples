package cmd

import (
	"net"
	"strconv"

	"github.com/spf13/cobra"
	"github.com/williamlsh/goexamples/pkg/log"
	"github.com/williamlsh/goexamples/pkg/tracing"
	"github.com/williamlsh/goexamples/services/customer"
	"go.uber.org/zap"
)

var customerCmd = &cobra.Command{
	Use:   "customer",
	Short: "Starts Customer service",
	Long:  `Starts Customer service.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		zapLogger := logger.With(zap.String("service", "customer"))
		logger := log.NewFactory(zapLogger)
		server := customer.NewServer(
			net.JoinHostPort("0.0.0.0", strconv.Itoa(customerPort)),
			tracing.Init("customer", metricsFactory, logger),
			metricsFactory,
			logger,
		)
		return logError(zapLogger, server.Run())
	},
}

func init() {
	RootCmd.AddCommand(customerCmd)
}
