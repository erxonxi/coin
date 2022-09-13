package cmd

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/erxonxi/coin/wallet"
	"github.com/spf13/cobra"
)

var walletCmd = &cobra.Command{
	Use:   "wallet",
	Short: "Command to manage your wallets.",
	Long:  `Command to manage your wallets.`,
	Run:   walletFunc,
}

func init() {
	rootCmd.AddCommand(walletCmd)

	walletCmd.Flags().BoolVarP((&create), "create", "c", false, "Create new wallet")
	walletCmd.Flags().StringVarP((&address), "address", "a", "", "Address to get information about")
}

func walletFunc(cmd *cobra.Command, args []string) {
	defer os.Exit(0)

	nodeID := os.Getenv("NODE_ID")
	if nodeID == "" {
		log.Panic("Please provide a NODE_ID")
	}

	if create == true {
		wallets, _ := wallet.CreateWallets(nodeID)
		address := wallets.AddWallet()
		wallets.SaveFile(nodeID)

		fmt.Printf("New address is: %s\n", address)
	} else {
		if address != "" {
			wallets, _ := wallet.InitializeWallets(nodeID)
			if !wallet.ValidateAddress(address) {
				log.Panic("Invalid address")
			}
			w := wallets.GetWallet(address)
			PrintWalletAddress(address, w)
		} else {
			wallets, _ := wallet.CreateWallets(nodeID)
			addresses := wallets.GetAllAddresses()

			for _, address := range addresses {
				fmt.Println(address)
			}
		}
	}
}

func PrintWalletAddress(address string, w wallet.Wallet) {
	var lines []string
	lines = append(lines, fmt.Sprintf("======ADDRESS:======\n %s ", address))
	lines = append(lines, fmt.Sprintf("======PUBLIC KEY:======\n %x", w.PublicKey))
	lines = append(lines, fmt.Sprintf("======PRIVATE KEY:======\n %x", wallet.DecodePriveteKey(w.PrivateKey).D.Bytes()))
	fmt.Println(strings.Join(lines, "\n"))
}
