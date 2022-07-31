package core

import (
	"fmt"
	"strings"
	"time"

	"github.com/manishmeganathan/essensio/common"
)

// Block is a struct that represents a Block of data in the BlockChain
type Block struct {
	BlockHeader

	// Transactions in the block
	BlockTxns Transactions
	// Number of blocks preceding the current block
	BlockHeight int64
	// Hash of the block header
	BlockHash common.Hash
}

// String implements the Stringer interface for Block
func (block *Block) String() string {
	var s strings.Builder

	s.WriteString(fmt.Sprintf("=======[%v][%v]\n", block.BlockHeight, time.Unix(block.Timestamp, 0)))
	s.WriteString(fmt.Sprintf("Block Hash: %v\n", block.BlockHash.Hex()))
	s.WriteString(fmt.Sprintf("Priori Hash: %v\n", block.Priori.Hex()))
	s.WriteString(fmt.Sprintf("Txn Count: %v\n", block.TxnCount()))
	s.WriteString(fmt.Sprintf("Nonce: %v\n", block.Nonce))
	s.WriteString("=========================================\n")

	return s.String()
}

// NewBlock generates a new Block for a given set of Transactions,
// the hash of the previous block and the block height
func NewBlock(txns Transactions, priori common.Hash, height int64) (*Block, error) {
	block := &Block{
		BlockTxns:   txns,
		BlockHeight: height,
	}

	// Generate the hash of the transactions
	summary, err := GenerateSummary(txns)
	if err != nil {
		return nil, fmt.Errorf("failed to generate transaction summary: %w", err)
	}

	// Create a BlockHeader with the priori and summary
	header := NewBlockHeader(priori, summary)
	block.BlockHeader = header

	// Mine the Block & set the block hash
	block.BlockHash = block.BlockHeader.Mint()

	return block, nil
}

// GenesisBlock returns a Block that represents a Genesis Block with just a Coinbase Transaction.
func GenesisBlock() (*Block, error) {
	return NewBlock(Transactions{newCoinbaseTransaction(common.MinerAddress())}, common.NullHash(), 0)
}

// TxnCount returns the number of Transaction items in the Block
func (block Block) TxnCount() int {
	return len(block.BlockTxns)
}

// Serialize implements the common.Serializable interface for Block.
// Converts the Block into a stream of bytes encoded using common.GobEncode.
func (block *Block) Serialize() ([]byte, error) {
	return common.GobEncode(block)
}

// Deserialize implements the common.Serializable interface for Block.
// Converts the given data into Block and sets it the method's receiver using common.GobDecode.
func (block *Block) Deserialize(data []byte) error {
	// Decode the data into a *Block
	object, err := common.GobDecode(data, new(Block))
	if err != nil {
		return err
	}

	// Cast the object into a *Block and
	// set it to the method receiver
	*block = *object.(*Block)
	return nil
}
