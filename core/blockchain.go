package core

import (
	"fmt"

	"github.com/manishmeganathan/essensio/common"
	"github.com/manishmeganathan/essensio/db"
)

var (
	ChainHeadKey   = []byte("state-chainhead")
	ChainHeightKey = []byte("state-chainheight")
)

// BlockChain represents a blockchain as a set of Blocks
type BlockChain struct {
	// Represents the database of blockchain data
	// This contains the state and blocks of the blockchain
	db *db.Database

	// Represents the hash of the last Block
	head common.Hash
	// Represents the height of the chain. Last block height+1
	height int64
}

// String implements the Stringer interface for BlockChain
func (chain *BlockChain) String() string {
	return fmt.Sprintf("Chain Head: %x || Chain Height: %v", chain.head, chain.height)
}

// AddBlock generates and appends a Block to the chain for a given string data.
// The generated block is stored in the database. Any error that occurs is returned.
func (chain *BlockChain) AddBlock(data string) error {
	// Create a new Block with the given data
	block := NewBlock(data, chain.head, chain.height)

	// Serialize the Block
	blockData, err := block.Serialize()
	if err != nil {
		return fmt.Errorf("block serialize failed: %w", err)
	}

	// Add block to db
	if err := chain.db.SetEntry(block.BlockHash, blockData); err != nil {
		return fmt.Errorf("block store to db failed: %w", err)
	}

	// Update the chain head with the new block hash and increment chain height
	chain.head = block.BlockHash
	chain.height++

	// Sync the chain state into the DB
	if err := chain.syncState(); err != nil {
		return fmt.Errorf("chain state sync failed: %w", err)
	}

	return nil
}

// NewBlockChain returns a new BlockChain with an initialized
// Genesis Block with the provided genesis data.
func NewBlockChain() (*BlockChain, error) {
	// Create a new BlockChain object
	chain := new(BlockChain)

	// Check if the database already exists
	if db.Exists() {
		// Load blockchain state from database
		if err := chain.load(); err != nil {
			return nil, fmt.Errorf("failed to load existing blockchain: %w", err)
		}

	} else {
		// Initialize blockchain state and database
		if err := chain.init(); err != nil {
			return nil, fmt.Errorf("failed to initialize new blockchain: $%w", err)
		}
	}

	return chain, nil
}

// load restarts an existing BlockChain from the database.
// It updates its in-memory chain state chain information from the DB.
func (chain *BlockChain) load() (err error) {
	// Open the database
	if chain.db, err = db.Open(); err != nil {
		return err
	}

	// Get the chain head and set it
	if chain.head, err = chain.db.GetEntry(ChainHeadKey); err != nil {
		return fmt.Errorf("chain head retrieve failed: %w", err)
	}

	// Get the chain height
	height, err := chain.db.GetEntry(ChainHeightKey)
	if err != nil {
		return fmt.Errorf("chain height retrieve failed: %w", err)
	}

	// Deserialize the height into an int64
	object, err := common.GobDecode(height, new(int64))
	if err != nil {
		return fmt.Errorf("error deserializing chain height: %w", err)
	}

	// Cast the object into an int64 and set it
	chain.height = *object.(*int64)
	return nil
}

// init initializes a new BlockChain in the database.
// It generates a Genesis Block and adds it to DB and updates all chain state data.
func (chain *BlockChain) init() (err error) {
	// Open the database
	if chain.db, err = db.Open(); err != nil {
		return err
	}

	fmt.Println(">>>> New Blockchain Initialization. Creating Genesis Block <<<<")

	// Create Genesis Block & serialize it
	genesisBlock := NewBlock("genesis", []byte{}, 0)
	genesisData, err := genesisBlock.Serialize()
	if err != nil {
		return fmt.Errorf("block serialize failed: %w", err)
	}

	// Add Genesis Block to DB
	if err := chain.db.SetEntry(genesisBlock.BlockHash, genesisData); err != nil {
		return fmt.Errorf("genesis block store to db failed: %w", err)
	}

	// Set the chain height and head into struct
	chain.head, chain.height = genesisBlock.BlockHash, 1

	// Sync the chain state into the DB
	if err := chain.syncState(); err != nil {
		return fmt.Errorf("chain state sync failed: %w", err)
	}

	return nil
}

// syncState updates the chain head and height values into the DB at keys
// specified by the ChainHeadKey and ChainHeightKey respectively.
func (chain *BlockChain) syncState() error {
	// Sync chain head into the DB
	if err := chain.db.SetEntry(ChainHeadKey, chain.head); err != nil {
		return fmt.Errorf("error syncing chain head: %w", err)
	}

	// Serialize the chain height
	height, err := common.GobEncode(chain.height)
	if err != nil {
		return fmt.Errorf("error serializing chain height: %w", err)
	}

	// Sync the encoded height into the DB
	if err := chain.db.SetEntry(ChainHeightKey, height); err != nil {
		return fmt.Errorf("error syncing chain height: %w", err)
	}

	return nil
}
