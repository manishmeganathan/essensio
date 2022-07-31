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
	Nonce     uint64 `json:"nonce"`
	Timestamp string `json:"timestamp"`

	BlockHash     string `json:"block_hash"`
	PrevBlockHash string `json:"prev_block_hash"`

	TxnCount     int                `json:"txn_count"`
	Transactions []BlockTransaction `json:"transactions"`
}

type BlockTransaction struct {
	To    string `json:"to"`
	From  string `json:"from"`
	Value uint64 `json:"value"`
	Nonce uint64 `json:"nonce"`
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

		transactions := make([]BlockTransaction, 0, block.TxnCount())
		for _, txn := range block.BlockTxns {
			transactions = append(transactions, BlockTransaction{
				string(txn.To), string(txn.From),
				txn.Value, txn.Nonce,
			})
		}

		chainresult.Blocks = append(chainresult.Blocks, ChainBlock{
			Height:        uint64(block.BlockHeight),
			Timestamp:     time.Unix(block.Timestamp, 0).Format(time.RFC3339),
			BlockHash:     block.BlockHash.Hex(),
			PrevBlockHash: block.Priori.Hex(),
			Nonce:         uint64(block.Nonce),
			Transactions:  transactions,
			TxnCount:      block.TxnCount(),
		})
	}

	*result = chainresult
	return nil
}
