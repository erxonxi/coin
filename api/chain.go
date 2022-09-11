package api

import (
	"encoding/json"
	"net/http"

	"github.com/erxonxi/coin/blockchain"
)

type ChainServer struct{}

func (s *ChainServer) GetChain(w http.ResponseWriter, r *http.Request) {
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
