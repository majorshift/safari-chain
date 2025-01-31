package crypto

import (
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"testing"
)

// helper function to get previous block's hash
func getPrevBlockHash(t *testing.T, bc *Blockchain, height uint32) Hash {
	prevHeader, err := bc.GetHeaderByHeight(height - 1)
	assert.Nil(t, err)

	return BlockHash{}.Hash(prevHeader)
}

// helper function initializes a new blockchain with a genesis block
func newBlockchainWithGenesisExample() *Blockchain {
	tx := NewTxWithSignature([]byte("Hello, World"))
	prevBlockHash := Hash{}

	validatorPrivateKey, _ := GeneratePrivateKey()
	genesisBlock := NewSignedBlockExample(validatorPrivateKey, []*Transaction{tx}, 0, prevBlockHash)

	log := logrus.New()
	return NewBlockchain(log, genesisBlock)
}

func TestNewBlockchain_NewBlockAddedSuccessfully(t *testing.T) {
	blockchain := newBlockchainWithGenesisExample()

	// confirm blockchain has height of 0
	assert.Equal(t, uint32(0), blockchain.GetBlockchainHeight())

	// add another block

	// block validator private key
	validatorPrivateKey, _ := GeneratePrivateKey()
	// get the block hash of the previous block in the blockchain
	prevBlockHash := getPrevBlockHash(t, blockchain, 1)
	tx1 := NewTxWithSignature([]byte("Hello, World"))
	// create a new block
	block1 := NewSignedBlockExample(validatorPrivateKey, []*Transaction{tx1}, 1, prevBlockHash)
	// add block1 to the blockchain
	err := blockchain.AddBlock(block1)
	assert.Nil(t, err)

	// blockchain height should be updated to 1
	assert.Equal(t, uint32(1), blockchain.GetBlockchainHeight())

	header, err := blockchain.GetHeaderByHeight(1)
	assert.Nil(t, err)

	assert.Equal(t, block1.Header, header)
}

func TestNewBlockchain_AddBlockFailsWithInvalidPreviousBlockHash(t *testing.T) {
	blockchain := newBlockchainWithGenesisExample()

	// confirm blockchain has height of 0
	assert.Equal(t, uint32(0), blockchain.GetBlockchainHeight())

	// add another block

	// block validator private key
	validatorPrivateKey, _ := GeneratePrivateKey()
	//set previous block hash to an invalid one
	prevBlockHash := Hash{}
	tx1 := NewTxWithSignature([]byte("Hello, World"))
	// create a new block
	block1 := NewSignedBlockExample(validatorPrivateKey, []*Transaction{tx1}, 1, prevBlockHash)
	// add block1 to the blockchain
	err := blockchain.AddBlock(block1)
	assert.NotNil(t, err)

	// blockchain height remains as 0
	assert.Equal(t, uint32(0), blockchain.GetBlockchainHeight())
}

func TestBlockTooHigh(t *testing.T) {
	blockchain := newBlockchainWithGenesisExample()
	assert.Equal(t, uint32(0), blockchain.GetBlockchainHeight())

	// add block at a height higher than the next available slot
	// block validator private key
	validatorPrivateKey, _ := GeneratePrivateKey()
	// get the block hash of the previous block in the blockchain
	prevBlockHash := getPrevBlockHash(t, blockchain, 1)
	tx1 := NewTxWithSignature([]byte("Hello, World"))
	// create a new block with height higher than the blockchain's height
	block1 := NewSignedBlockExample(validatorPrivateKey, []*Transaction{tx1}, 2, prevBlockHash)

	err := blockchain.AddBlock(block1)
	// addition of block1 to blockchain fails
	assert.Error(t, err)
	// blockchain height remains at 0
	assert.Equal(t, uint32(0), blockchain.GetBlockchainHeight())
}
