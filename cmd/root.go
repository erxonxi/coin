package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

var address string

var rootCmd = &cobra.Command{
	Use:   "coin",
	Short: "Command line interface for coin blockchain.",
	Long:  `You can run flag "--help" for more information about the command interface.`,
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
}
