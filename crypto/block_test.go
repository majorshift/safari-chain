package crypto

import (
	"github.com/stretchr/testify/assert"
	"strconv"
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

func TestVerifyBlock_Success(t *testing.T) {
	// create a validator private key
	validator, _ := GeneratePrivateKey()
	// the new block in this case is the genesis block
	// so no previous hash
	prevBlockHash := Hash{}

	var transactions []*Transaction
	// create seven signed transactions adding them to the transactiosn list
	for i := 0; i < 7; i++ {
		// append the index to Hello, World to generate unique transactions
		transactions = append(transactions, NewTxWithSignature([]byte("Hello, World"+strconv.Itoa(i))))
	}

	// create a signed block
	b := NewSignedBlockExample(validator, transactions, 0, prevBlockHash)

	// error is not because transaction is not signed
	assert.Nil(t, b.Verify())
}

func TestVerifyBlock_FailsForUnsignedTransaction(t *testing.T) {
	validator, _ := GeneratePrivateKey()
	prevBlockHash := Hash{}

	// create a transaction without a signature
	unsignedTx := &Transaction{Data: []byte("Hello, World")}
	b := NewSignedBlockExample(validator, []*Transaction{unsignedTx}, 0, prevBlockHash)

	// error is not because transaction is not signed
	assert.NotNil(t, b.Verify())
}

func TestComputeMerkleRoot(t *testing.T) {
	var transactions []*Transaction
	for i := 0; i < 7; i++ {
		// append the index to Hello, World to generate unique transactions
		transactions = append(transactions, NewTxWithSignature([]byte("Hello, World"+strconv.Itoa(i))))
	}

	validator, _ := GeneratePrivateKey()
	// create a new block: the merkle root is created during block creation
	b := NewSignedBlockExample(validator, transactions, 0, Hash{})

	// recalculate merkle root
	originalMerkleRoot := ComputeMerkleRoot(transactions)
	assert.Equal(t, b.MerkleRoot, originalMerkleRoot)

	// alter the first transaction
	b.Transactions[0] = NewTxWithSignature([]byte("Altered Transaction"))

	// recalculate merkle root: comparison to original merkle root
	// should show a mismatch
	alteredMerkleRoot := ComputeMerkleRoot(b.Transactions)
	assert.NotEqual(t, originalMerkleRoot, alteredMerkleRoot)

	// verification fails
	assert.Error(t, b.Verify())
}
