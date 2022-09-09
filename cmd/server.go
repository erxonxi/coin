package cmd

import (
	"log"
	"net/http"

	"github.com/erxonxi/coin/blockchain"
	"github.com/gorilla/websocket"
	"github.com/spf13/cobra"
)

var upgrader = websocket.Upgrader{} // use default options

var runCmd = &cobra.Command{
	Use:   "server",
	Short: "A command to run a network",
	Long:  `This command will run a network for coin blockchain.`,
	Run:   serverFun,
}

func init() {
	chain := blockchain.InitBlockChain()
	defer chain.Database.Close()

	chain.AddBlock("Hello World")
	chain.AddBlock("Other World")
	chain.AddBlock("Three World")

	chain.PrintChain()

	rootCmd.AddCommand(runCmd)
}

func serverFun(cmd *cobra.Command, args []string) {
	http.HandleFunc("/echo", echo)
	log.Fatal(http.ListenAndServe("localhost:8000", nil))
}

func echo(w http.ResponseWriter, r *http.Request) {
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Print("upgrade:", err)
		return
	}
	defer c.Close()

	for {
		mt, message, err := c.ReadMessage()
		if err != nil {
			log.Println("read:", err)
			break
		}

		log.Printf("recv: %s", message)

		err = c.WriteMessage(mt, message)
		if err != nil {
			log.Println("write:", err)
			break
		}
	}
}
