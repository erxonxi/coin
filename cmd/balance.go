package cmd

import (
	"fmt"
	"log"

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

	balanceCmd.Flags().StringVarP((&address), "address", "a", "15AfJY1BtvMsD5Zzd7mtBLyaxQavTESxaa", "Addres of blockchain")
}

func getBalance() {
	if !wallet.ValidateAddress(address) {
		log.Panic("Address is not Valid")
	}
	chain := blockchain.ContinueBlockChain(address)
	UTXOSet := blockchain.UTXOSet{chain}
	defer chain.Database.Close()

	balance := 0
	pubKeyHash := wallet.Base58Decode([]byte(address))
	pubKeyHash = pubKeyHash[1 : len(pubKeyHash)-4]
	UTXOs := UTXOSet.FindUTXO(pubKeyHash)

	for _, out := range UTXOs {
		balance += out.Value
	}

	fmt.Printf("Balance of %s: %d\n", address, balance)
}
