package crypto

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
)

const (
	hashLen = 32
)

type Hash [32]uint8

// IsHashed checks if hash has been set
func (h Hash) IsHashed() bool {
	for i := 0; i < hashLen; i++ {
		if h[i] != 0 {
			return false
		}
	}

	return true
}

// ToByteArray converts hash to []byte
func (h Hash) ToByteArray() []byte {
	b := make([]byte, hashLen)
	for i := 0; i < hashLen; i++ {
		b[i] = h[i]
	}

	return b
}

// ToString encodes hash to string
func (h Hash) ToString() string {
	return hex.EncodeToString(h.ToByteArray())
}

// BytesToHash receives []byte and converts to Hash
func BytesToHash(b []byte) (Hash, error) {
	if len(b) != hashLen {
		return Hash{}, fmt.Errorf("expected bytes of length 32, got %d", len(b))
	}

	var value [hashLen]uint8
	for i := 0; i < hashLen; i++ {
		value[i] = b[i]
	}

	return Hash(value), nil
}

// Hasher interface
type Hasher[T any] interface {
	Hash(T) Hash
}

// BlockHash to implement the Hasher interface for blocks
type BlockHash struct{}

// Hash hashes the header of the block
func (BlockHash) Hash(header *Header) Hash {
	h := sha256.Sum256(header.ToBytes())
	return Hash(h)
}

// TxHash to implement the Hasher interface for transactions
type TxHash struct{}

// Hash hashes the Data property of a Transaction
func (TxHash) Hash(tx *Transaction) Hash {
	return Hash(sha256.Sum256(tx.Data))
}
