package cmd

import (
	"context"
	"crypto/rand"
	"io"
	"log"
	mrand "math/rand"

	"github.com/spf13/cobra"

	"github.com/erxonxi/coin/network"
)

var port int
var dest string
var debug bool

var nodeCmd = &cobra.Command{
	Use:   "node",
	Short: "",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		var r io.Reader
		if debug {
			r = mrand.New(mrand.NewSource(int64(port)))
		} else {
			r = rand.Reader
		}

		h, err := network.MakeHost(port, r)
		if err != nil {
			log.Println(err)
			return
		}

		if dest == "" {
			network.StartPeer(ctx, h, network.HandleStream)
		} else {
			rw, err := network.StartPeerAndConnect(ctx, h, dest)
			if err != nil {
				log.Println(err)
				return
			}

			// Create a thread to read and write data.
			go network.WriteData(rw)
			go network.ReadData(rw)
		}

		select {}
	},
}

func init() {
	rootCmd.AddCommand(nodeCmd)

	nodeCmd.PersistentFlags().IntVarP((&port), "port", "p", 3131, "The source port")
	nodeCmd.Flags().StringVarP((&dest), "dest", "d", "", "The addres of destination multiaddress")
	nodeCmd.Flags().BoolVarP((&debug), "debug", "g", false, "Debug mode")
}
