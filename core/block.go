package core

import (
	"fmt"
	"strings"

	"github.com/manishmeganathan/essensio/common"
)

// Block is a struct that represents a Block of data in the BlockChain
type Block struct {
	BlockHeader

	// Number of blocks preceding the current block
	BlockHeight int64
	// Raw data of the block (placeholder for transactions)
	BlockData []byte
	// Hash of the block header
	BlockHash common.Hash
}

// String implements the Stringer interface for Block
func (b *Block) String() string {
	var s strings.Builder

	s.WriteString(fmt.Sprintf("Block Hash: %x\n", b.BlockHash))
	s.WriteString(fmt.Sprintf("Priori Hash: %x\n", b.BlockHeader.Priori))
	s.WriteString(fmt.Sprintf("Data: %v\n", string(b.BlockData)))
	s.WriteString(fmt.Sprintf("Timestamp: %v\n", b.BlockHeader.Timestamp))
	s.WriteString(fmt.Sprintf("Nonce: %v\n", b.BlockHeader.Nonce))

	return s.String()
}

// NewBlock generates a new Block for some given data,
// the hash of the previous block and the block height
func NewBlock(data string, priori []byte, height int64) *Block {
	block := &Block{
		BlockData:   []byte(data),
		BlockHeight: height,
	}

	// Generate the hash of the data
	summary := common.Hash256(block.BlockData)
	// Create a BlockHeader with the priori and summary
	header := NewBlockHeader(priori, summary)
	block.BlockHeader = header

	// Mine the Block & set the block hash
	block.BlockHash = block.BlockHeader.Mint(height)

	return block
}
