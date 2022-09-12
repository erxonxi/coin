package cmd

import (
	"context"
	"fmt"
	"time"

	"github.com/erxonxi/coin/p2p"
	"github.com/libp2p/go-libp2p"
	"github.com/libp2p/go-libp2p-core/crypto"

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

		for {
			time.Sleep(time.Second * 2)
			for _, peer := range node.Network().Peers() {
				node.Ping(peer)
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(nodeCmd)

	nodeCmd.PersistentFlags().IntVarP((&port), "port", "p", 3131, "The source port")
	nodeCmd.Flags().StringVarP((&hostAddress), "host", "o", "0.0.0.0", "The host address")
	nodeCmd.Flags().StringVarP((&pid), "pid", "i", "/chain/1.1.0", "The host address")
	nodeCmd.Flags().StringVarP((&group), "group", "g", "main", "The group of peers name")
}

// helper method - create a lib-p2p host to listen on a port
func makeNode(port int, done chan bool) *p2p.Node {
	priv, _, _ := crypto.GenerateKeyPair(crypto.Secp256k1, 256)
	listen, _ := ma.NewMultiaddr(fmt.Sprintf("/ip4/127.0.0.1/tcp/%d", port))
	host, _ := libp2p.New(
		libp2p.ListenAddrs(listen),
		libp2p.Identity(priv),
	)

	return p2p.NewNode(host, done)
}
