package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/erxonxi/coin/blockchain"
)

var from string
var to string
var amount int

var sendCmd = &cobra.Command{
	Use:   "send",
	Short: "Command to send amount to address",
	Long: `You can send amount to address whith address
For example:
coin send --from "Xonxi" --to "Juan" --amount 10
`,
	Run: sendFunc,
}

func init() {
	rootCmd.AddCommand(sendCmd)

	sendCmd.Flags().StringVarP((&from), "from", "f", "Xonxi", "Address from send amount")
	sendCmd.Flags().StringVarP((&to), "to", "t", "Juan", "Address to send amount")
	sendCmd.Flags().IntVarP((&amount), "amount", "a", 10, "Address to send amount")
}

func sendFunc(cmd *cobra.Command, args []string) {
	chain := blockchain.ContinueBlockChain(from)
	defer chain.Database.Close()

	tx := blockchain.NewTransaction(from, to, amount, chain)
	chain.AddBlock([]*blockchain.Transaction{tx})
	fmt.Println("Send successfuly!")
}
