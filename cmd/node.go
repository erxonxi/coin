package cmd

import (
	"context"
	"fmt"
	"log"

	"github.com/erxonxi/coin/p2p"
	"github.com/libp2p/go-libp2p"
	"github.com/libp2p/go-libp2p-core/crypto"
	"github.com/libp2p/go-libp2p-core/peerstore"

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
		fmt.Println("Found peer:", peer, ", connecting")

		if err := node.Connect(ctx, peer); err != nil {
			fmt.Println("Connection failed:", err)
		}

		log.Printf("NODE_ID: %s\n", node.ID())

		for {
			fmt.Println(node.Network().Peers())
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

func run(h1, h2 *p2p.Node, done <-chan bool) {
	// connect peers
	h1.Peerstore().AddAddrs(h2.ID(), h2.Addrs(), peerstore.PermanentAddrTTL)
	h2.Peerstore().AddAddrs(h1.ID(), h1.Addrs(), peerstore.PermanentAddrTTL)

	// send messages using the protocols
	h1.Ping(h2.Host)
	h2.Ping(h1.Host)

	// block until all responses have been processed
	for i := 0; i < 8; i++ {
		<-done
	}
}
