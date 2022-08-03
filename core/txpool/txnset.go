package txpool

import (
	"container/heap"
	"sort"

	"github.com/manishmeganathan/essensio/core"
)

// TransactionSet represents a group of nonce sorted Transactions
type TransactionSet struct {
	// Hashmap containing the Transactions indexed by their nonce
	items map[uint64]*core.Transaction
	// Heap containing all the nonces for which Transactions exist
	index *nonceHeap
}

// NewTransactionSet generates and returns a new TransactionSet
func NewTransactionSet() *TransactionSet {
	return &TransactionSet{
		items: make(map[uint64]*core.Transaction),
		index: new(nonceHeap),
	}
}

// Flatten returns a flat slice of nonce-sorted transactions as core.Transactions.
func (txset *TransactionSet) Flatten() core.Transactions {
	// Create a slice in which to collect transactions
	transactions := make(core.Transactions, 0, len(txset.items))

	// Iterate over the items and add them to the transactions slice
	for _, tx := range txset.items {
		transactions = append(transactions, tx)
	}

	// Wrap the Transactions into a TxnByNonce sorter and sort them
	sort.Sort(TxnByNonce(transactions))
	return transactions
}

// Put inserts the given txn into the set while maintaining the sorting
func (txset *TransactionSet) Put(txn *core.Transaction) {
	// Get the transaction nonce
	nonce := txn.Nonce
	// Add txn into the items
	txset.items[nonce] = txn
	// Push the txn nonce into the index heap
	heap.Push(txset.index, nonce)
}

// Remove removes the transaction with the given nonce from the set.
// Returns a bool indicating if the transaction was removed.
func (txset *TransactionSet) Remove(nonce uint64) bool {
	// Check if transactions exists in the items
	if _, ok := txset.items[nonce]; ok {
		// Delete from items
		delete(txset.items, nonce)

		// Search for index with nonce in index
		for i := 0; i < txset.index.Len(); i++ {
			if (*txset.index)[i] == nonce {
				// Remove the nonce from the heap
				heap.Remove(txset.index, i)
				break
			}
		}

		// Return that txn was removed
		return true
	}

	// No transaction was removed
	return false
}
