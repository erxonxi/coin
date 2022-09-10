package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

var address string

var rootCmd = &cobra.Command{
	Use:   "coin",
	Short: "Command Interface for coin network communication.",
	Long:  `You can run flag "--help" for more information about the command interface.`,
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.Flags().StringVarP((&address), "address", "a", "15AfJY1BtvMsD5Zzd7mtBLyaxQavTESxaa", "Addres of blockchain")
}
