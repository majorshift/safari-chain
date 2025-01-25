package crypto

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSignTransaction(t *testing.T) {
	// private key of sender
	privKeyFrom, err := GeneratePrivateKey()
	assert.NoError(t, err)
	// private key of receiver
	privKeyReceiver, err := GeneratePrivateKey()
	assert.NoError(t, err)

	// public key of sender
	pubKeyFrom := privKeyFrom.PublicKey()
	// public key of receiver
	pubKeyReceiver := privKeyReceiver.PublicKey()

	// data to be sent
	data := []byte("Hello, World")

	// create new transaction
	tx := NewTransaction(pubKeyFrom, pubKeyReceiver, data)

	// transaction is not signed yet, signature property should be nil
	assert.Nil(t, tx.Signature)
	// verifying at this stage fails since the transaction is not signed
	assert.Error(t, tx.Verify()) // error is: transaction has no signature

	//	Sign transaction
	tx.Sign(privKeyFrom)
	assert.True(t, tx.Signature.Verify(pubKeyFrom, data))
}

func TestVerifyTransaction(t *testing.T) {
	privKeyFrom, err := GeneratePrivateKey()
	assert.NoError(t, err)

	privKeyReceiver, err := GeneratePrivateKey()
	assert.NoError(t, err)

	pubKeyFrom := privKeyFrom.PublicKey()
	pubKeyReceiver := privKeyReceiver.PublicKey()

	data := []byte("Hello, World")
	tx := NewTransaction(pubKeyFrom, pubKeyReceiver, data)

	// sign transaction with appropriate private key(sender's)
	tx.Sign(privKeyFrom)
	// no error when verifying
	assert.NoError(t, tx.Verify())

	// should sign with sender's private key(the From property of Transaction)
	// here we sign with the wrong private key (receiver's private key)
	tx.Sign(privKeyReceiver)

	// verification should fail with the error: transaction has invalid signature
	assert.Error(t, tx.Verify())
}
