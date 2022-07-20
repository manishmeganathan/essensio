package core

import (
	"bytes"
	"crypto/sha256"
	"fmt"
	"strings"
)

type Block struct {
	// Header
	PrioriHash []byte
	BlockData  []byte

	// Body
	BlockHash []byte
}

// GenerateHash generates a BlockHash = Hash(Header)
func (b *Block) GenerateHash() []byte {
	header := bytes.Join([][]byte{b.PrioriHash, b.BlockData}, []byte{})
	hash := sha256.Sum256(header)
	return hash[:]
}

func (b *Block) String() string {
	var s strings.Builder

	s.WriteString("Block Hash: ")
	s.WriteString(fmt.Sprintf("%x", b.BlockHash))
	s.WriteString("\n")

	s.WriteString("Block Data: ")
	s.WriteString(string(b.BlockData))
	s.WriteString("\n")

	s.WriteString("Prev Block Hash: ")
	s.WriteString(fmt.Sprintf("%x", b.PrioriHash))

	return s.String()
}

// NewBlock generates a new Block for some given
// data and the hash of the previous block.
func NewBlock(data string, priori []byte) *Block {
	block := &Block{
		PrioriHash: priori,
		BlockData:  []byte(data),
	}

	block.BlockHash = block.GenerateHash()
	return block
}
