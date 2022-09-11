package cmd

import (
	"bufio"
	"context"
	"crypto/rand"
	"fmt"

	"github.com/erxonxi/coin/network"
	"github.com/libp2p/go-libp2p"
	"github.com/libp2p/go-libp2p-core/crypto"
	"github.com/libp2p/go-libp2p-core/protocol"
	"github.com/multiformats/go-multiaddr"
	"github.com/spf13/cobra"
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
		fmt.Printf("[*] Listening on: %s with port: %d\n", hostAddress, port)

		ctx := context.Background()
		r := rand.Reader

		// Creates a new RSA key pair for this host.
		prvKey, _, err := crypto.GenerateKeyPairWithReader(crypto.RSA, 2048, r)
		if err != nil {
			panic(err)
		}

		// 0.0.0.0 will listen on any interface device.
		sourceMultiAddr, _ := multiaddr.NewMultiaddr(fmt.Sprintf("/ip4/%s/tcp/%d", hostAddress, port))

		// libp2p.New constructs a new libp2p Host.
		// Other options can be added here.
		host, err := libp2p.New(
			libp2p.ListenAddrs(sourceMultiAddr),
			libp2p.Identity(prvKey),
		)
		if err != nil {
			panic(err)
		}

		// Set a function as stream handler.
		// This function is called when a peer initiates a connection and starts a stream with this peer.
		host.SetStreamHandler(protocol.ID(pid), network.HandleStream)

		fmt.Printf("\n[*] Your Multiaddress Is: /ip4/%s/tcp/%v/p2p/%s\n", hostAddress, port, host.ID().Pretty())

		peerChan := network.InitMDNS(host, group)

		peer := <-peerChan // will block untill we discover a peer
		fmt.Println("Found peer:", peer, ", connecting")

		if err := host.Connect(ctx, peer); err != nil {
			fmt.Println("Connection failed:", err)
		}

		// open a stream, this stream will be handled by handleStream other end
		stream, err := host.NewStream(ctx, peer.ID, protocol.ID(pid))

		if err != nil {
			fmt.Println("Stream open failed", err)
		} else {
			rw := bufio.NewReadWriter(bufio.NewReader(stream), bufio.NewWriter(stream))

			go network.WriteData(rw)
			go network.ReadData(rw)
			fmt.Println("Connected to:", peer)
		}

		select {} //wait here

	},
}

func init() {
	rootCmd.AddCommand(nodeCmd)

	nodeCmd.PersistentFlags().IntVarP((&port), "port", "p", 3131, "The source port")
	nodeCmd.Flags().StringVarP((&hostAddress), "host", "o", "0.0.0.0", "The host address")
	nodeCmd.Flags().StringVarP((&pid), "pid", "i", "/chain/1.1.0", "The host address")
	nodeCmd.Flags().StringVarP((&group), "group", "g", "main", "The group of peers name")
}
