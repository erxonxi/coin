package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/erxonxi/coin/blockchain"
)

// balanceCmd represents the balance command
var balanceCmd = &cobra.Command{
	Use:   "balance",
	Short: "Command to get balance.",
	Long: `Command to get balance.
For example:
coin balance --address "Xonxi"
`,
	Run: func(cmd *cobra.Command, args []string) {
		getBalance()
	},
}

func init() {
	rootCmd.AddCommand(balanceCmd)

	balanceCmd.Flags().StringVarP((&address), "address", "a", "Xonxi", "Create new blockchain")
}

func getBalance() {
	chain := blockchain.ContinueBlockChain(address)
	defer chain.Database.Close()

	balance := 0
	UTXOs := chain.FindUTXO(address)

	for _, out := range UTXOs {
		balance += out.Value
	}

	fmt.Printf("Balance of %s: %d\n", address, balance)
}
