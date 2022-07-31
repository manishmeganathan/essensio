package jsonrpc

import (
	"fmt"
	"log"
	"net/http"

	"github.com/manishmeganathan/essensio/common"
	"github.com/manishmeganathan/essensio/core"
)

type AddBlockArgs struct {
	Transactions []TransactionInput `json:"transactions"`
}

type TransactionInput struct {
	To    string `json:"to"`
	From  string `json:"from"`
	Value uint64 `json:"value"`
}

type AddBlockResult struct {
	BlockHeight uint64 `json:"block_height"`
	BlockHash   string `json:"block_hash"`
}

func (api *API) AddBlock(r *http.Request, args *AddBlockArgs, result *AddBlockResult) error {
	log.Println("'AddBlock' Called")

	if len(args.Transactions) == 0 {
		return fmt.Errorf("no transactions receieved")
	}

	transactions := make(core.Transactions, 0, len(args.Transactions))
	for _, txn := range args.Transactions {
		newtxn := core.NewTransaction(common.Address(txn.From), common.Address(txn.To), 0, txn.Value)
		transactions = append(transactions, newtxn)
	}

	if err := api.chain.AddBlock(transactions); err != nil {
		return fmt.Errorf("failed to add block: %w", err)
	}

	*result = AddBlockResult{
		BlockHeight: uint64(api.chain.Height - 1),
		BlockHash:   api.chain.Head.Hex(),
	}

	return nil
}
