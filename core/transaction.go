package core

import (
	"bytes"
	"fmt"

	"github.com/manishmeganathan/essensio/common"
)

// Transactions is a group of Transaction objects
type Transactions []*Transaction

// Transaction represents a transaction between two addresses.
// It contains a nonce value to make it unique for transactions
// between the same account with the same value.
type Transaction struct {
	// Represents the amount of tokens transferred in Nubs
	Value uint64
	// Represents the sender account nonce
	Nonce uint64

	// Represents the address of the sender
	From common.Address
	// Represents the address of the receiver
	To common.Address
}

// NewTransaction generates a new Transaction between from and to for the given value and nonce.
func NewTransaction(from, to common.Address, nonce, value uint64) *Transaction {
	return &Transaction{value, nonce, from, to}
}

// newCoinbaseTransaction generates a new coinbase transaction that mints tokens for the given address.
// The value of the transaction is the default Block Reward for mining a block.
func newCoinbaseTransaction(address common.Address) *Transaction {
	return &Transaction{BlockReward, 0, address, common.NullAddress()}
}

// Serialize implements the common.Serializable interface for Transaction.
// Converts the Transaction into a stream of bytes encoded using common.GobEncode.
func (txn *Transaction) Serialize() ([]byte, error) {
	return common.GobEncode(txn)
}

// Deserialize implements the common.Serializable interface for Transaction.
// Converts the given data into Transaction and sets it the method's receiver using common.GobDecode.
func (txn *Transaction) Deserialize(data []byte) error {
	// Decode the data into a *Transaction
	object, err := common.GobDecode(data, new(Transaction))
	if err != nil {
		return err
	}

	// Cast the object into a *Transaction and
	// set it to the method receiver
	*txn = *object.(*Transaction)
	return nil
}

// Hash returns the SHA-256	hash of the Transaction's serialized representation.
func (txn *Transaction) Hash() common.Hash {
	data, err := txn.Serialize()
	if err != nil {
		return common.NullHash()
	}

	return common.Hash256(data)
}

// GenerateSummary generates a summary hash for a given set of Transactions.
// Currently, concatenates the hash of all given transactions and hashes that data to obtain the summary.
// This is a valid method of summary generation but does not for allow tamper detection or inclusivity checks
func GenerateSummary(txns Transactions) (common.Hash, error) {
	// Iterate over each transaction, obtain
	// its hash and append it into the buffer
	var buffer bytes.Buffer
	for _, txn := range txns {
		hash := txn.Hash()
		if hash == common.NullHash() {
			return hash, fmt.Errorf("got null hash")
		}

		buffer.Write(hash.Bytes())
	}

	// Generate the hash of the buffer bytes
	txnsum := common.Hash256(buffer.Bytes())
	return txnsum, nil
}
