package crypto

import (
	"fmt"
)

// Transaction structure
type Transaction struct {
	Data      []byte
	From      *PublicKey // public key of the one sending value/initiating transaction
	Receiver  *PublicKey // public key of the one receiving value
	Signature *Signature // signature verifying transaction's authenticity
}

func NewTransaction(from, receiver *PublicKey, data []byte) *Transaction {
	return &Transaction{
		From:     from,
		Receiver: receiver,
		Data:     data,
	}
}

// Sign signs a transaction
func (tx *Transaction) Sign(privKey *PrivateKey) {
	sig := privKey.Sign(tx.Data)

	tx.Signature = sig
}

// Hash hashes a transaction
func (tx *Transaction) Hash(hasher Hasher[*Transaction]) Hash {
	return hasher.Hash(tx)
}

// Verify checks the validity of the transaction signature
func (tx *Transaction) Verify() error {
	// transaction must be signed
	if tx.Signature == nil {
		return fmt.Errorf("transaction has no signature")
	}
	// if signature is set, it must be valid and not tampered with
	if !tx.Signature.Verify(tx.From, tx.Data) {
		return fmt.Errorf("transaction has invalid signature")
	}

	return nil
}
