package jsonrpc

import (
	"log"

	"github.com/manishmeganathan/essensio/core/chainmgr"
)

type API struct {
	chain *chainmgr.ChainManager
}

func NewAPI() *API {
	chain, err := chainmgr.NewChainManager()
	if err != nil {
		log.Fatalln("Failed to Start Blockchain:", err)
	}

	return &API{chain}
}

func (api *API) Stop() {
	api.chain.Stop()
}
