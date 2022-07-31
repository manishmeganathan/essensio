package jsonrpc

import (
	"log"
	"net/http"
	"time"
)

type ShowChainArgs struct{}

type ShowChainResult struct {
	ChainHead   string       `json:"chain_head"`
	ChainHeight uint64       `json:"chain_height"`
	Blocks      []ChainBlock `json:"blocks"`
}

type ChainBlock struct {
	Height    uint64 `json:"height"`
	Timestamp string `json:"timestamp"`

	BlockHash     string `json:"block_hash"`
	PrevBlockHash string `json:"prev_block_hash"`

	Nonce uint64 `json:"nonce"`
	Data  string `json:"data"`
}

func (api *API) ShowChain(r *http.Request, args *ShowChainArgs, result *ShowChainResult) error {
	log.Println("'ShowChain' Called")

	chainresult := ShowChainResult{
		ChainHead:   api.chain.Head.Hex(),
		ChainHeight: uint64(api.chain.Height),
	}

	iterator := api.chain.NewIterator()
	for !iterator.Done() {
		// Get the next block
		block, err := iterator.Next()
		if err != nil {
			log.Fatalln("Iterator Error:", err)
		}

		chainresult.Blocks = append(chainresult.Blocks, ChainBlock{
			Height:        uint64(block.BlockHeight),
			Timestamp:     time.Unix(block.Timestamp, 0).Format(time.RFC3339),
			BlockHash:     block.BlockHash.Hex(),
			PrevBlockHash: block.Priori.Hex(),
			Nonce:         uint64(block.Nonce),
			Data:          string(block.BlockData),
		})
	}

	*result = chainresult
	return nil
}
