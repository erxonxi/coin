package cmd

import (
	"context"
	"fmt"
	"os"
	"runtime"
	"strconv"
	"syscall"

	"github.com/erxonxi/coin/blockchain"
	"github.com/erxonxi/coin/p2p"
	"github.com/libp2p/go-libp2p"
	"github.com/libp2p/go-libp2p-core/crypto"
	"github.com/vrecan/death"

	"github.com/spf13/cobra"

	ma "github.com/multiformats/go-multiaddr"
)

var port int
var hostAddress string
var pid string
var group string

var nodeCmd = &cobra.Command{
	Use:   "node",
	Short: "",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		done := make(chan bool, 1)

		node := makeNode(port, done)

		// autodiscover
		ctx := context.Background()
		peerChan := p2p.InitMDNS(node, group)
		peer := <-peerChan // will block untill we discover a peer

		if err := node.Connect(ctx, peer); err != nil {
			fmt.Println("Connection failed:", err)
		}

		for _, peer := range node.Network().Peers() {
			block := blockchain.Genesis(blockchain.CoinbaseTx(address, ""))
			node.SendBlock(peer, block)
		}

		for {
		}
	},
}

func init() {
	rootCmd.AddCommand(nodeCmd)

	nodeCmd.PersistentFlags().IntVarP((&port), "port", "p", 3000, "The source port")
	nodeCmd.Flags().StringVarP((&hostAddress), "host", "o", "0.0.0.0", "The host address")
	nodeCmd.Flags().StringVarP((&pid), "pid", "i", "/chain/1.1.0", "The host address")
	nodeCmd.Flags().StringVarP((&group), "group", "g", "main", "The group of peers name")
	nodeCmd.Flags().StringVarP((&address), "address", "a", "15AfJY1BtvMsD5Zzd7mtBLyaxQavTESxaa", "Address to mine")
}

// helper method - create a lib-p2p host to listen on a port
func makeNode(port int, done chan bool) *p2p.Node {
	priv, _, _ := crypto.GenerateKeyPair(crypto.Secp256k1, 256)
	listen, _ := ma.NewMultiaddr(fmt.Sprintf("/ip4/127.0.0.1/tcp/%d", port))
	host, _ := libp2p.New(
		libp2p.ListenAddrs(listen),
		libp2p.Identity(priv),
	)

	node := p2p.NewNode(host, done)

	chain := blockchain.ContinueBlockChain(strconv.Itoa(port))
	defer chain.Database.Close()
	go CloseDB(chain)

	node.BlockChain = chain

	return node
}

func CloseDB(chain *blockchain.BlockChain) {
	d := death.NewDeath(syscall.SIGINT, syscall.SIGTERM, os.Interrupt)

	d.WaitForDeathWithFunc(func() {
		defer os.Exit(1)
		defer runtime.Goexit()
		chain.Database.Close()
	})
}
