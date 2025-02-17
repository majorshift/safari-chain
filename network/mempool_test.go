package network

import (
	"github.com/majorshift/safari-chain/crypto"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestTxPoolAdd(t *testing.T) {
	p := NewMempool(1)
	assert.Equal(t, 0, p.PendingTxCount())
	assert.Equal(t, 0, p.AllTxCount())

	tx := crypto.NewTxWithSignature([]byte("hello world"))
	p.Add(tx)

	assert.Equal(t, 1, p.PendingTxCount())
	assert.Equal(t, 1, p.AllTxCount())

	txHash := tx.Hash(&crypto.TxHash{})
	assert.True(t, p.allTransactions.Contains(txHash))
}

func TestTxPool_Prune(t *testing.T) {
	// maximum number of transactions in pool is 1
	p := NewMempool(1)

	tx1 := crypto.NewTxWithSignature([]byte("hello world"))
	p.Add(tx1)
	assert.Equal(t, 1, p.AllTxCount())

	txHash1 := tx1.Hash(&crypto.TxHash{})
	assert.True(t, p.allTransactions.Contains(txHash1))

	// add another transaction. expect original tx to be pruned
	tx2 := crypto.NewTxWithSignature([]byte("hello world 1"))
	p.Add(tx2)

	assert.Equal(t, 1, p.AllTxCount())

	// transaction count remains the at maximum
	assert.Equal(t, 1, p.allTransactions.Count())
}

func TestTxPool_ClearPending(t *testing.T) {
	p := NewMempool(1)

	tx1 := crypto.NewTxWithSignature([]byte("hello world"))
	p.Add(tx1)
	assert.Equal(t, 1, p.PendingTxCount())

	// clear pending
	p.ClearPendingList()
	assert.Equal(t, 0, p.PendingTxCount())
}
