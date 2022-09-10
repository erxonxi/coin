package cmd

import (
	"fmt"
	"log"
	"os"

	"github.com/spf13/cobra"

	"github.com/erxonxi/coin/blockchain"
	"github.com/erxonxi/coin/network"
	"github.com/erxonxi/coin/wallet"
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
	Run: func(cmd *cobra.Command, args []string) {
		nodeID := os.Getenv("NODE_ID")
		if nodeID == "" {
			log.Panic("Please provide a NODE_ID")
		}

		send(from, to, amount, nodeID, false)
	},
}

func init() {
	rootCmd.AddCommand(sendCmd)

	sendCmd.Flags().StringVarP((&from), "from", "f", "Xonxi", "Address from send amount")
	sendCmd.Flags().StringVarP((&to), "to", "t", "Juan", "Address to send amount")
	sendCmd.Flags().IntVarP((&amount), "amount", "a", 10, "Address to send amount")
}

func send(from, to string, amount int, nodeID string, mineNow bool) {
	if !wallet.ValidateAddress(to) {
		log.Panic("Address is not Valid")
	}
	if !wallet.ValidateAddress(from) {
		log.Panic("Address is not Valid")
	}
	chain := blockchain.ContinueBlockChain(nodeID)
	UTXOSet := blockchain.UTXOSet{chain}
	defer chain.Database.Close()

	wallets, err := wallet.CreateWallets(nodeID)
	if err != nil {
		log.Panic(err)
	}
	wallet := wallets.GetWallet(from)

	tx := blockchain.NewTransaction(&wallet, to, amount, &UTXOSet)
	if mineNow {
		cbTx := blockchain.CoinbaseTx(from, "")
		txs := []*blockchain.Transaction{cbTx, tx}
		block := chain.MineBlock(txs)
		UTXOSet.Update(block)
	} else {
		network.SendTx(network.KnownNodes[0], tx)
		fmt.Println("send tx")
	}

	fmt.Println("Success!")

}
