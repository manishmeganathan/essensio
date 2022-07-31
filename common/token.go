package common

const (
	// Nub is the smallest unit of a token in the Essencio Blockchain
	Nub uint64 = 1

	// Pith is 1,000 Nub
	Pith = 1000 * Nub
	// Esse is 1,000 Pith or 1,000,000 Nub
	Esse = 1000 * Pith

	// Essence is 1,000 Esse or 1,000,000 Pith or 1,000,000,000 Nub.
	// An Essence is the default unit of tokens in the Essencio Blockchain.
	Essence = 1000 * Esse

	// Quintessence is 5 Essence and is the fixed
	// block reward in the Essencio Blockchain
	Quintessence = 5 * Essence
)
