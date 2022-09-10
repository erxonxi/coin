package cmd

import (
	"fmt"
	"os"

	"github.com/erxonxi/coin/wallet"
	"github.com/spf13/cobra"
)

var walletCmd = &cobra.Command{
	Use:   "wallet",
	Short: "Command to manage your wallet.",
	Long: `Command to manage your wallet.
For example:
`,
	Run: walletFunc,
}

func init() {
	rootCmd.AddCommand(walletCmd)

	walletCmd.Flags().BoolVarP((&create), "create", "c", false, "Create new wallet")
}

func walletFunc(cmd *cobra.Command, args []string) {
	defer os.Exit(0)

	if create == true {
		wallets, _ := wallet.CreateWallets()
		address := wallets.AddWallet()
		wallets.SaveFile()

		fmt.Printf("New address is: %s\n", address)
	} else {
		wallets, _ := wallet.CreateWallets()
		addresses := wallets.GetAllAddresses()

		for _, address := range addresses {
			fmt.Println(address)
		}
	}
}
