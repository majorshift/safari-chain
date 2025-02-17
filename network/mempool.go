package network

import (
	"github.com/majorshift/safari-chain/crypto"
	"github.com/majorshift/safari-chain/types"
	"sync"
)

// MemPool is the structure for the mempool
type MemPool struct {
	allTransactions     *TxMap // Stores all transactions in the pool
	pendingTransactions *TxMap // Stores only pending transactions
	maxSize             int    // Maximum number of transactions in the pool
}

func NewMempool(maxLength int) *MemPool {
	return &MemPool{
		allTransactions:     NewTxMap(),
		pendingTransactions: NewTxMap(),
		maxSize:             maxLength,
	}
}

// Add inserts a new transaction to the mempool
func (p *MemPool) Add(tx *crypto.Transaction) {
	// prune the oldest transaction that is sitting in the allTransactions pool
	if p.allTransactions.Count() == p.maxSize {
		oldest := p.allTransactions.First()
		p.allTransactions.Remove(oldest.Hash(crypto.TxHash{}))
	}

	// prevent duplicate inclusion of transactions to mempool
	if !p.allTransactions.Contains(tx.Hash(crypto.TxHash{})) {
		p.allTransactions.Add(tx)
		p.pendingTransactions.Add(tx)
	}
}

// Contains checks if a transaction already exists in the mempool
func (p *MemPool) Contains(hash crypto.Hash) bool {
	return p.allTransactions.Contains(hash)
}

// GetPendingTx returns a slice of transactions that are in the pending pool
func (p *MemPool) GetPendingTx() []*crypto.Transaction {
	return p.pendingTransactions.transactionList.Data
}

// ClearPendingList deletes all transactions in the pending list of the mempool
func (p *MemPool) ClearPendingList() {
	p.pendingTransactions.Clear()
}

// PendingTxCount returns the number of transactions in the pending list
func (p *MemPool) PendingTxCount() int {
	return p.pendingTransactions.Count()
}

// AllTxCount returns the number of transactions in ever handled by the mempool
func (p *MemPool) AllTxCount() int {
	return p.allTransactions.Count()
}

// TxMap The TxMap struct acts as a thread-safe transaction storage
// with both fast lookup and ordered access capabilities.
// It efficiently manages transactions by combining a hashmap
// for quick retrieval and a list for ordered storage while ensuring
// safe concurrent access using a read-write mutex.
type TxMap struct {
	// Lock for concurrent access
	lock sync.RWMutex
	// Quick lookup for transactions by hash
	transactionsByHash map[crypto.Hash]*crypto.Transaction
	// Ordered list of transactions
	transactionList *types.List[*crypto.Transaction]
}

func NewTxMap() *TxMap {
	return &TxMap{
		transactionsByHash: make(map[crypto.Hash]*crypto.Transaction),
		transactionList:    types.NewList[*crypto.Transaction](),
	}
}

// First returns the first transaction in the pool
func (t *TxMap) First() *crypto.Transaction {
	t.lock.RLock()
	defer t.lock.RUnlock()
	first := t.transactionList.Get(0)
	return t.transactionsByHash[first.Hash(crypto.TxHash{})]
}

// Get returns transaction in the pool matching hash
func (t *TxMap) Get(h crypto.Hash) *crypto.Transaction {
	t.lock.RLock()
	defer t.lock.RUnlock()
	return t.transactionsByHash[h]
}

// Add inserts a new transaction to the pool
func (t *TxMap) Add(tx *crypto.Transaction) {
	hash := tx.Hash(crypto.TxHash{})

	t.lock.Lock()
	defer t.lock.Unlock()
	if _, ok := t.transactionsByHash[hash]; !ok {
		t.transactionsByHash[hash] = tx
		t.transactionList.Insert(tx)
	}
}

// Remove deletes transaction matching hash from the pool
func (t *TxMap) Remove(h crypto.Hash) {
	t.lock.Lock()
	defer t.lock.Unlock()
	t.transactionList.Remove(t.transactionsByHash[h])
	delete(t.transactionsByHash, h)
}

// Count returns the number of transactions in the transactionsByHash map
func (t *TxMap) Count() int {
	t.lock.RLock()
	defer t.lock.RUnlock()
	return len(t.transactionsByHash)
}

// Contains checks if a transaction matching hash is contained in the mempool of a node
func (t *TxMap) Contains(h crypto.Hash) bool {
	t.lock.RLock()
	defer t.lock.RUnlock()
	_, ok := t.transactionsByHash[h]
	return ok
}

// Clear removes all transactions in the transactionsByHash
func (t *TxMap) Clear() {
	t.lock.Lock()
	defer t.lock.Unlock()
	t.transactionsByHash = make(map[crypto.Hash]*crypto.Transaction)
	t.transactionList.Clear()
}
