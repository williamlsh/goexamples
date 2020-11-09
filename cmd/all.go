package cmd

import "github.com/spf13/cobra"

var allCmd = &cobra.Command{
	Use:   "all",
	Short: "Starts all services",
	Long:  `Starts all services.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		logger.Info("Starting all services")
		go customerCmd.RunE(customerCmd, args)
		go driverCmd.RunE(driverCmd, args)
		go routeCmd.RunE(routeCmd, args)
		return frontendCmd.RunE(frontendCmd, args)
	},
}

func init() {
	RootCmd.AddCommand(allCmd)
}
