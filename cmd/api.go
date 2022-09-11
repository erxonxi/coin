package cmd

import (
	"github.com/spf13/cobra"

	"github.com/erxonxi/coin/api"
)

var apiCmd = &cobra.Command{
	Use:   "api",
	Short: "Command to run the API server.",
	Long:  `Command to run the API server.`,
	Run: func(cmd *cobra.Command, args []string) {
		api.StartServer()
	},
}

func init() {
	rootCmd.AddCommand(apiCmd)
}
