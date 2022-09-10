package cmd

import (
	"github.com/erxonxi/coin/blockchain"
	"github.com/spf13/cobra"
)

var create bool
var address string

var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "A command to run a network",
	Long:  `This command will run a network for coin blockchain.`,
	Run:   serverFun,
}

func init() {
	rootCmd.AddCommand(serverCmd)

	serverCmd.Flags().StringVarP((&address), "address", "a", "Xonxi", "Create new blockchain")
	serverCmd.Flags().BoolVarP((&create), "create", "c", false, "Create new blockchain")
}

func serverFun(cmd *cobra.Command, args []string) {
	if create == true {
		createBlockChain(address)
	}

	printBlockChain(address)
}

func createBlockChain(address string) {
	chain := blockchain.InitBlockChain(address)
	defer chain.Database.Close()
}

func printBlockChain(address string) {
	chain := blockchain.ContinueBlockChain(address)
	chain.PrintChain()
	defer chain.Database.Close()
}
