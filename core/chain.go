package core

type BlockChain struct {
	Blocks []*Block
}

func (chain *BlockChain) AddBlock(data string) {
	// Get the hash of the previous block
	prevBlock := chain.Blocks[len(chain.Blocks)-1]
	prevHash := prevBlock.BlockHash

	// Create a new Block and append it
	block := NewBlock(data, prevHash)
	chain.Blocks = append(chain.Blocks, block)
}

func NewBlockChain(genesis string) *BlockChain {
	genesisblock := NewBlock(genesis, []byte{})
	return &BlockChain{[]*Block{genesisblock}}
}
