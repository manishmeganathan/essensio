package txpool

import "github.com/manishmeganathan/essensio/core"

// nonceHeap is a slice of sorted nonce numbers (uint64)
// It implements the heap.Interface type.
type nonceHeap []uint64

// Len returns the number of elements in the nonceHeap.
// Implements the sort.Interface type for heap.Interface.
func (n nonceHeap) Len() int { return len(n) }

// Less returns whether the value at index i is less than the value at index j in the nonceHeap.
// Implements the sort.Interface type for heap.Interface.
func (n nonceHeap) Less(i, j int) bool { return n[i] < n[j] }

// Swap switches the elements at index i and j in the nonceHeap.
// Implements the sort.Interface type for heap.Interface.
func (n nonceHeap) Swap(i, j int) { n[i], n[j] = n[j], n[i] }

// Push pushes x into the nonceHeap.
// Implements the heap.Interface type for nonceHeap.
func (n *nonceHeap) Push(x any) {
	*n = append(*n, x.(uint64))
}

// Pop pops the last element from the nonceHeap (the lowest nonce).
// Implements the heap.Interface type for nonceHeap
func (n *nonceHeap) Pop() any {
	// Dereference the slice and get its length
	old := *n
	count := len(old)

	// Get the element at len-1
	x := old[count-1]
	// Update n as a slice that does not contain the popped element
	*n = old[0 : count-1]

	// Return the popped element
	return x
}

// TxnByNonce implements the sort interface to allow sorting a list of transactions
// by their nonces. This is usually only useful for sorting transactions from a
// single account, otherwise a nonce comparison doesn't make much sense.
type TxnByNonce core.Transactions

// Len returns the number of elements in the TxnByNonce.
// Implements the sort.Interface type
func (s TxnByNonce) Len() int { return len(s) }

// Less returns whether the value at index i is less than the value at index j in the TxnByNonce.
// Implements the sort.Interface type.
func (s TxnByNonce) Less(i, j int) bool { return s[i].Nonce < s[j].Nonce }

// Swap switches the elements at index i and j in the TxnByNonce.
// Implements the sort.Interface type.
func (s TxnByNonce) Swap(i, j int) { s[i], s[j] = s[j], s[i] }
