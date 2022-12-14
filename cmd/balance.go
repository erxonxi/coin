package cmd

import (
	"fmt"
	"log"
	"os"

	"github.com/spf13/cobra"

	"github.com/erxonxi/coin/blockchain"
	"github.com/erxonxi/coin/wallet"
)

// balanceCmd represents the balance command
var balanceCmd = &cobra.Command{
	Use:   "balance",
	Short: "Command to get balance.",
	Long: `Command to get balance.
For example:
coin balance --address "15AfJY1BtvMsD5Zzd7mtBLyaxQavTESxaa"
`,
	Run: func(cmd *cobra.Command, args []string) {
		getBalance()
	},
}

func init() {
	rootCmd.AddCommand(balanceCmd)

	balanceCmd.Flags().StringVarP((&address), "address", "a", "15AfJY1BtvMsD5Zzd7mtBLyaxQavTESxaa", "Your address of blockchain")
}

func getBalance() {
	nodeID := os.Getenv("NODE_ID")
	if nodeID == "" {
		log.Panic("Please provide a NODE_ID")
	}

	if !wallet.ValidateAddress(address) {
		log.Panic("Address is not Valid")
	}

	chain := blockchain.ContinueBlockChain(nodeID)
	UTXOSet := blockchain.UTXOSet{Blockchain: chain}
	defer chain.Database.Close()

	balance := 0
	pubKeyHash := wallet.Base58Decode([]byte(address))
	pubKeyHash = pubKeyHash[1 : len(pubKeyHash)-4]
	UTXOs := UTXOSet.FindUnspentTransactions(pubKeyHash)

	for _, out := range UTXOs {
		balance += out.Value
	}

	fmt.Printf("Balance of %s: %d\n", address, balance)
}
