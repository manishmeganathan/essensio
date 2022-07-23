package common

import "crypto/sha256"

type Hash []byte

// Hash256 generates a 256-bit hash of some given data.
// The output of the given hash is equivalent to double hashing with the
// SHA2-256 hashing algorithm, rendering it safe from length extension attacks.
func Hash256(data []byte) Hash {
	var hash [32]byte

	hash = sha256.Sum256(data)
	hash = sha256.Sum256(hash[:])

	return hash[:]
}
