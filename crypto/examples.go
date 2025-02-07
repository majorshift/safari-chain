package crypto

import (
	"time"
)

// NewTxWithSignature returns a signed transaction
func NewTxWithSignature(data []byte) *Transaction {
	fromPrivKey, _ := GeneratePrivateKey()
	receiverPrivKey, _ := GeneratePrivateKey()
	tx := NewTransaction(fromPrivKey.PublicKey(), receiverPrivKey.PublicKey(), data)
	tx.Sign(fromPrivKey)

	return tx
}

// ExampleBlock returns a block
func ExampleBlock(height uint32, prevBlockHash Hash) *Block {
	tx := NewTxWithSignature([]byte("hello world"))
	header := &Header{
		Version:       1,
		PrevBlockHash: prevBlockHash,
		Height:        height,
		Timestamp:     time.Now().UnixNano(),
	}

	return NewBlock(header, []*Transaction{tx})
}

// NewSignedBlockExample creates a new blocked that is signed by a validator
func NewSignedBlockExample(validator *PrivateKey, transactions []*Transaction, height uint32, prevBlockHash Hash) *Block {
	header := &Header{
		Version:       1,
		PrevBlockHash: prevBlockHash,
		Height:        height,
		Timestamp:     time.Now().UnixNano(),
	}

	b := NewBlock(header, transactions)
	b.Sign(validator)

	return b
}
