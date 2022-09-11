package api

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/erxonxi/coin/blockchain"
	"github.com/gorilla/mux"
)

func StartServer() {
	app := &httpServer{}

	r := mux.NewRouter()
	r.HandleFunc("/chain", app.GetChain).Methods("GET")

	server := &http.Server{
		Addr:    ":8080",
		Handler: r,
	}

	log.Panic(server.ListenAndServe())
}

type httpServer struct {
}

func (s *httpServer) GetChain(w http.ResponseWriter, r *http.Request) {
	fmt.Println("GET CHAIN")
	chain := blockchain.ContinueBlockChain("3000")
	defer chain.Database.Close()
	iter := chain.Iterator()

	var blocks []*blockchain.Block

	for {
		block := iter.Next()
		blocks = append(blocks, block)

		if len(block.PrevHash) == 0 {
			break
		}
	}

	err := json.NewEncoder(w).Encode(&blocks)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
