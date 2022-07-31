package jsonrpc

import (
	"fmt"
	"log"
	"net/http"
)

type AddBlockArgs struct {
	BlockData string `json:"block_data"`
}

type AddBlockResult struct {
	BlockHeight uint64 `json:"block_height"`
	BlockHash   string `json:"block_hash"`
}

func (api *API) AddBlock(r *http.Request, args *AddBlockArgs, result *AddBlockResult) error {
	log.Println("'AddBlock' Called")

	if args.BlockData == "" {
		return fmt.Errorf("no input data for block")
	}

	if err := api.chain.AddBlock(args.BlockData); err != nil {
		return fmt.Errorf("failed to add block: %w", err)
	}

	*result = AddBlockResult{
		BlockHeight: uint64(api.chain.Height - 1),
		BlockHash:   api.chain.Head.Hex(),
	}

	return nil
}
