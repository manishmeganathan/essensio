package chainmgr

import (
	"fmt"

	"github.com/manishmeganathan/essensio/common"
	"github.com/manishmeganathan/essensio/core"
	"github.com/manishmeganathan/essensio/db"
)

// ChainIterator is a struct that can iterate
// over each Block in a blockchain.
type ChainIterator struct {
	// Represents the hash of the current Block on the iterator
	cursor common.Hash
	// Represents the database containing all Block data indexed by their hash
	database *db.Database
}

// NewIterator constructs a new ChainIterator for the BlockChain.
func (chain *ChainManager) NewIterator() *ChainIterator {
	return &ChainIterator{chain.Head, chain.db}
}

// Next returns the next Block in the ChainIterator.
// Returns an error if a Block is not found or is invalid.
func (iter *ChainIterator) Next() (*core.Block, error) {
	// Find the Block with hash represented by the iterator cursor
	data, err := iter.database.GetEntry(iter.cursor.Bytes())
	if err != nil {
		return nil, fmt.Errorf("cannot finding block '%x': %w", iter.cursor, err)
	}

	// Create a new Block and deserialize the block data into it
	block := new(core.Block)
	if err := block.Deserialize(data); err != nil {
		return nil, fmt.Errorf("block deserialize failed: %w", err)
	}

	// Update the iterator cursor to the hash of the previous Block
	iter.cursor = block.Priori
	return block, nil
}

// Done returns whether the ChainIterator has
// reached the Genesis Block of the chain.
func (iter *ChainIterator) Done() bool {
	// If the cursor hash is null, the ChainIterator is done
	return iter.cursor == common.NullHash()
}
