package core

import (
	"fmt"
	"log"
	"math"
	"math/big"

	"github.com/manishmeganathan/essensio/common"
)

// Difficulty represents the number of bits that need to be 0 for the Proof Of Work Algorithm.
// Currently Static, but can eventually be adjusted based on the total hash rate of the network,
// to achieve a block time of n minutes.
const Difficulty uint8 = 18

// GenerateTarget returns a big.Int with the target hash value for the current difficulty
func GenerateTarget() *big.Int {
	// Generate a new big Integer and left shift to match difficulty
	target := big.NewInt(1)
	target.Lsh(target, 256-uint(Difficulty))

	return target
}

// Mint is the Proof of Work routine that generates a nonce
// that is valid for the Target difficulty of the header.
func (header *BlockHeader) Mint() common.Hash {
	var hash common.Hash

	// Reset Nonce
	header.Nonce = 0

	for header.Nonce < math.MaxInt64 {
		// Serialize the Header
		data, err := header.Serialize()
		if err != nil {
			log.Fatalln("header serialization failed during PoW:", err)
		}

		// Hash the Header data
		hash = common.Hash256(data)

		// Print the hash mining process
		fmt.Printf("\rMining Block [%v]: %v", header.Nonce, hash.Hex())

		// Compare the hash with target
		if hash.Big().Cmp(header.Target) == -1 {
			break // Block Mined!
		} else {
			// Increment Nonce & Repeat
			header.Nonce++
		}
	}

	fmt.Println()
	return hash
}

// Validate is the Proof of Work validation routine.
// Returns a boolean indicating if the hash of the block is valid for its target.
func (header *BlockHeader) Validate() bool {
	// Serialize the Header
	data, err := header.Serialize()
	if err != nil {
		log.Fatalln("header serialization failed during PoW:", err)
	}

	// Hash the Header data
	hash := common.Hash256(data)

	// Compare hash with target
	return hash.Big().Cmp(header.Target) == -1
}
