package cmd

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/erxonxi/coin/blockchain"
	"github.com/erxonxi/coin/network"
	"github.com/erxonxi/coin/wallet"
	"github.com/spf13/cobra"
)

var create bool
var printChain bool

var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "A command to run a network",
	Long:  `This command will run a network for coin blockchain.`,
	Run:   serverFun,
}

func init() {
	rootCmd.AddCommand(serverCmd)

	serverCmd.Flags().StringVarP((&address), "address", "a", "17fX5v53fG1kHacpAPEw8BW8stqQh1L7ku", "Yout address to connect to blockchain")
	serverCmd.Flags().BoolVarP((&create), "create", "c", false, "Create a new blockchain")
	serverCmd.Flags().BoolVarP((&printChain), "print", "p", false, "Print blocks of blockchain")
}

func serverFun(cmd *cobra.Command, args []string) {
	nodeID := os.Getenv("NODE_ID")
	if nodeID == "" {
		log.Panic("Please provide a NODE_ID")
	}

	if printChain == true {
		printBlockChain(nodeID)
		return
	}

	if create == true {
		createBlockChain(nodeID, address)
	} else {
		StartNode(nodeID, address)
	}
}

func createBlockChain(nodeID string, address string) {
	if !wallet.ValidateAddress(address) {
		log.Panic("Address is not Valid")
	}
	chain := blockchain.InitBlockChain(address, nodeID)

	UTXOSet := blockchain.UTXOSet{chain}
	UTXOSet.Reindex()

	chain.Database.Close()
	fmt.Println("Finished!")
}

func StartNode(nodeID, minerAddress string) {
	fmt.Printf("Starting Node %s\n", nodeID)

	if len(minerAddress) > 0 {
		if wallet.ValidateAddress(minerAddress) {
			fmt.Println("Mining is on. Address to receive rewards: ", minerAddress)
		} else {
			log.Panic("Wrong miner address!")
		}
	}
	network.StartServer(nodeID, minerAddress)
}

func printBlockChain(nodeID string) {
	chain := blockchain.ContinueBlockChain(nodeID)
	defer chain.Database.Close()
	iter := chain.Iterator()

	for {
		block := iter.Next()

		fmt.Printf("Hash: %x\n", block.Hash)
		fmt.Printf("Prev. hash: %x\n", block.PrevHash)
		pow := blockchain.NewProof(block)
		fmt.Printf("PoW: %s\n", strconv.FormatBool(pow.Validate()))
		for _, tx := range block.Transactions {
			fmt.Println(tx)
		}
		fmt.Println()

		if len(block.PrevHash) == 0 {
			break
		}
	}
}
