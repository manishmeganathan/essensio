package common

// Address represents the address for an Account
// Placeholder for [20]byte type Addresses.
type Address string

// Bytes returns the byte representation of the Address
func (addr Address) Bytes() []byte {
	return []byte(addr)
}

// NullAddress returns a zero Address
func NullAddress() Address {
	return ""
}

// MinerAddress returns the Address to use for Coinbase Transactions
func MinerAddress() Address {
	return "manish"
}
