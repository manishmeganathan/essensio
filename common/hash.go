package common

import (
	"crypto/sha256"
	"math/big"
)

// HashLength is the expected length of a hash
const HashLength = 32

// Hash represents a 32 byte hash of some arbitrary data
type Hash [HashLength]byte

// BytesToHash converts a []byte into a Hash.
// If b is greater than HashLength, it will be cropped from the left
func BytesToHash(b []byte) (h Hash) {
	// If b is longer than HashLength, it is cropped
	if len(b) > len(h) {
		b = b[len(b)-HashLength:]
	}

	// Copy the value and return
	copy(h[HashLength-len(b):], b)
	return
}

// NullHash returns a zero Hash
func NullHash() Hash { return [32]byte{} }

// Bytes returns the byte representation of the Hash
func (h Hash) Bytes() []byte { return h[:] }

// Big returns the Hash as big integer
func (h Hash) Big() *big.Int { return new(big.Int).SetBytes(h[:]) }

// Hex returns the Hash as a hex string
func (h Hash) Hex() string { return HexEncode(h.Bytes()) }

// String implements the Stringer interface for Hash.
// Returns the Hash as hex string.
func (h Hash) String() string { return h.Hex() }

// Hash256 generates a 256-bit hash of some given data.
// The output of the given hash is equivalent to double hashing with the
// SHA2-256 hashing algorithm, rendering it safe from length extension attacks.
func Hash256(data []byte) Hash {
	var hash Hash

	hash = sha256.Sum256(data)
	hash = sha256.Sum256(hash[:])

	return hash
}
