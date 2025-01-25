package crypto

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestSignBlock(t *testing.T) {
	// create private key
	privKey, _ := GeneratePrivateKey()

	// the prevHash in this case will be nil since the block being created
	// is the genesis block
	prevBlockHash := Hash{}
	b := ExampleBlock(1, prevBlockHash)

	// sign block
	// the private key used to sign the block belongs to the validator
	b.Sign(privKey)

	// after signing, the validator property f the block is set to
	// its public key
	assert.Equal(t, b.Validator, privKey.PublicKey())
	assert.NotNil(t, b.Signature)
}
