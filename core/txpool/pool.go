package txpool

import (
	"sync"

	"github.com/manishmeganathan/essensio/common"
	"github.com/manishmeganathan/essensio/core"
)

const (
	PoolSize = 100
	PoolPage = 20
)

type TxnPool interface {
	// Fetch collects Transactions from the pool.
	// They are not removed from the pool until Clear is called with that Transaction
	Fetch() core.Transactions
	// FetchFor collects Transactions from queue for a specific address
	// They are not removed from the pool until Clear is called with the Transaction
	FetchFor(common.Address) *TransactionSet

	// Insert inserts Transactions into the pool
	Insert(...*core.Transaction)
	// Contains returns whether a transaction exists for a given transaction hash.
	// Will return true only if the transaction exists in the active set.
	Contains(common.Hash) bool

	// Clear removes the given Transactions from the pending and active
	Clear(...*core.Transaction)
	// Restore attempts to restore the given Transaction from the pending set into to pool.
	Restore(...*core.Transaction)
	// Purge completely removes all Transactions pending or not from the pool
	Purge() int

	// Pending returns the number of transactions in the pending set
	// i.e, transaction that have been fetched from the pool but not cleared.
	Pending() int
	// Active returns the number of transactions in active set.
	// This does not include transactions in the pending set.
	Active() int
}

// TxnNoncePool is a Transaction Pool that implements the TxnPool interface.
// Transactions inserted into the pool are pushed into the queued, and are concurrently
// processed and added to the pending set of transactions and the lookup directory.
type TxnNoncePool struct {
	// thread safety mutex
	mu *sync.RWMutex

	// pool is the collection of Transactions in the TxnNoncePool
	// grouped by address and indexed by nonce
	pool map[common.Address]*TransactionSet

	// pending is a collection of Transactions in the TxnNoncePool
	// that have been collected but yet to be cleared from the pool.
	pending map[common.Hash]*core.Transaction

	// lookup is a flat collection of Transactions that can be used
	// to peek into the pool. Does not contain pending transactions.
	lookup map[common.Hash]*core.Transaction
}

// NewTxnNoncePool generates and returns a new TxnNoncePool object
func NewTxnNoncePool() *TxnNoncePool {
	return &TxnNoncePool{
		mu:      &sync.RWMutex{},
		pool:    make(map[common.Address]*TransactionSet),
		pending: make(map[common.Hash]*core.Transaction),
		lookup:  make(map[common.Hash]*core.Transaction),
	}
}

// Fetch implements the TxnPool interface for TxnNoncePool.
// Collects upto PoolPage number of Transactions from the pool's active set.
// The collected transactions are moved to the pending set and returned.
func (pool *TxnNoncePool) Fetch() core.Transactions {
	// Acquire the mutex
	pool.mu.Lock()
	defer pool.mu.Unlock()

	transactions := make(core.Transactions, 0, PoolPage)

	// Iterate over the transactions in lookup
	for _, txn := range pool.lookup {
		// Hash the transaction
		txnhash := txn.Hash()

		// Remove from the active set (pool and lookup)
		pool.pool[txn.From].Remove(txn.Nonce)
		delete(pool.lookup, txnhash)

		// Add the txns to the pending set
		pool.pending[txnhash] = txn

		// Add to output set of transactions
		transactions = append(transactions, txn)

		// End iteration if transactions Page is full
		if len(transactions) == PoolPage {
			break
		}
	}

	return transactions
}

// FetchFor implements the TxnPool interface for TxnNoncePool.
// Returns a TransactionSet object for address. (nil if no Transactions from address).
// All Transactions in the TransactionSet are moved from the active set into the pending set.
func (pool *TxnNoncePool) FetchFor(address common.Address) *TransactionSet {
	// Acquire Mutex
	pool.mu.Lock()
	defer pool.mu.Unlock()

	// Get transaction for Address from Pool
	txset := pool.pool[address]

	// Flatten the transaction set and iterate through it
	for _, txn := range txset.Flatten() {
		// Hash the transaction
		txnhash := txn.Hash()

		// Remove from lookup
		delete(pool.lookup, txnhash)
		// Add to pending
		pool.pending[txnhash] = txn
	}

	// Remove from transaction set from the active pool
	delete(pool.pool, address)
	return txset
}

// Insert implements the TxnPool interface for TxnNoncePool.
// Accepts a variadic number of Transactions and adds each one to the active set.
func (pool *TxnNoncePool) Insert(transactions ...*core.Transaction) {
	// Acquire Mutex
	pool.mu.Lock()
	defer pool.mu.Unlock()

	// Iterate over transactions
	for _, txn := range transactions {
		// Get the ref to TransactionSet for the sender address
		txset := pool.pool[txn.From]
		if txset == nil {
			// If no TransactionSet for address, create new and set into pool
			txset = NewTransactionSet()
			pool.pool[txn.From] = txset
		}

		// Add the transaction into the set for the address
		txset.Put(txn)
		// Add the transaction into the lookup
		pool.lookup[txn.Hash()] = txn
	}
}

// Contains implements the TxnPool interface for TxnNoncePool.
// Returns whether the Transaction with given Hash is present in the active set of the pool
func (pool *TxnNoncePool) Contains(hash common.Hash) bool {
	// Acquire RLock
	pool.mu.RLock()
	defer pool.mu.RUnlock()

	// Check if transaction exists in the lookup
	_, ok := pool.lookup[hash]
	return ok
}

// Clear implements the TxnPool interface for TxnNoncePool.
// Accepts a variadic number of Transactions and removes each of them from the pool.
// The transaction is removed from both the active and pending set, wherever it may exist.
func (pool *TxnNoncePool) Clear(transactions ...*core.Transaction) {
	// Acquire Mutex
	pool.mu.Lock()
	defer pool.mu.Unlock()

	// Iterate over transactions
	for _, txn := range transactions {
		// Hash the transaction
		txnhash := txn.Hash()

		// Remove from active set
		if _, ok := pool.lookup[txnhash]; ok {
			delete(pool.lookup, txnhash)
			pool.pool[txn.From].Remove(txn.Nonce)
		}

		// Remove from pending set
		delete(pool.pending, txnhash)
	}
}

// Restore implements the TxnPool interface for TxnNoncePool.
// Accepts a variadic number of transactions and restores
// each of them from the pending set into the active set.
func (pool *TxnNoncePool) Restore(transactions ...*core.Transaction) {
	// Acquire Mutex
	pool.mu.Lock()
	defer pool.mu.Unlock()

	// Iterate over the transactions
	for _, txn := range transactions {
		// Get the transaction hash
		txnhash := txn.Hash()

		// Check if transaction is available in the pending set
		if _, ok := pool.pending[txnhash]; ok {
			// Get the ref to TransactionSet for the sender address
			txset := pool.pool[txn.From]
			if txset == nil {
				// If no TransactionSet for address, create new and set into pool
				txset = NewTransactionSet()
				pool.pool[txn.From] = txset
			}

			// Add the transaction into the set
			txset.Put(txn)
			// Add to lookup of active set
			pool.lookup[txnhash] = txn

			// Remove from pending set
			delete(pool.pending, txnhash)
		}
	}
}

// Purge implements the TxnPool interface for TxnNoncePool.
// Resets the pool and removes all transactions from it.
func (pool *TxnNoncePool) Purge() int {
	// Acquire Mutex
	pool.mu.Lock()
	defer pool.mu.Unlock()

	// Collect total number of transactions in the pool
	count := pool.Pending() + pool.Active()

	pool.pool = make(map[common.Address]*TransactionSet)
	pool.pending = make(map[common.Hash]*core.Transaction)
	pool.lookup = make(map[common.Hash]*core.Transaction)

	return count
}

// Pending implements the TxnPool interface for TxnNoncePool.
// Returns the number of transactions in the pending set.
func (pool *TxnNoncePool) Pending() int {
	// Acquire RLock
	pool.mu.RLock()
	defer pool.mu.RUnlock()

	return len(pool.pending)
}

// Active implements the TxnPool interface for TxnNoncePool.
// Returns the number of transactions in the active set.
func (pool *TxnNoncePool) Active() int {
	// Acquire RLock
	pool.mu.RLock()
	defer pool.mu.RUnlock()

	return len(pool.lookup)
}
