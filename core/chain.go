package core

import "github.com/manishmeganathan/essensio/common"

// BlockChain represents a blockchain as a set of Blocks
type BlockChain struct {
	Blocks []*Block
}

// AddBlock generates and appends a Block to the chain for a given data.
// Returns the Block hash after appending it.
func (chain *BlockChain) AddBlock(data string) common.Hash {
	// Get the hash of the previous block
	prevBlock := chain.Blocks[len(chain.Blocks)-1]
	prevHash := prevBlock.BlockHash

	// Create a new Block and append it
	block := NewBlock(data, prevHash, prevBlock.BlockHeight+1)
	chain.Blocks = append(chain.Blocks, block)

	return block.BlockHash
}

// NewBlockChain returns a new BlockChain with an initialized
// Genesis Block with the provided genesis data.
func NewBlockChain(genesis string) *BlockChain {
	genesisblock := NewBlock(genesis, []byte{}, 0)
	return &BlockChain{[]*Block{genesisblock}}
}
