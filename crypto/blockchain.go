package crypto

import (
	"fmt"
	"sync"

	"github.com/sirupsen/logrus"
)

type Blockchain struct {
	lock      sync.RWMutex   // Mutex to manage concurrent read and write access
	headers   []*Header      // Slice to store block headers instead of full blocks for efficiency
	logger    *logrus.Logger // Logger to track blockchain activity and debugging
	validator Validator      // Validator to verify block and transaction validity
}

func NewBlockchain(log *logrus.Logger, genesis *Block) *Blockchain {
	bc := &Blockchain{
		headers: []*Header{},
		logger:  log,
	}

	bc.validator = NewBlockValidator(bc)
	bc.addBlock(genesis)
	return bc
}

// AddBlock adds a new block to the blockchain
func (bc *Blockchain) AddBlock(b *Block) error {
	if err := bc.validator.ValidateBlock(b); err != nil {
		return err
	}

	// add block
	bc.addBlock(b)

	return nil
}

// HasBlock checks if a block of a given height exists in the blockchain
func (bc *Blockchain) HasBlock(height uint32) bool {
	return height <= bc.GetBlockchainHeight() // less than/equal to since height begins at 0
}

// GetHeaderByHeight returns header at given height
// or error is height is higher than the blockchain height
func (bc *Blockchain) GetHeaderByHeight(height uint32) (*Header, error) {
	if height > bc.GetBlockchainHeight() {
		return nil, fmt.Errorf("given height (%d) is greater than the blockchain height", height)
	}

	bc.lock.Lock()
	defer bc.lock.Unlock()

	return bc.headers[height], nil
}

// GetBlockchainHeight returns the height of the entire blockchain
// height is calculated similar to array indices hence
// height = (number of blocks in the blockchain) - 1
func (bc *Blockchain) GetBlockchainHeight() uint32 {
	bc.lock.RLock()
	defer bc.lock.RUnlock()

	return uint32(len(bc.headers) - 1)
}

// addBlock a new block to the blockchain by appending it's
// header to the blockchain header list
func (bc *Blockchain) addBlock(b *Block) {
	bc.lock.Lock()
	defer bc.lock.Unlock()

	bc.headers = append(bc.headers, b.Header)
	bc.logger.WithFields(logrus.Fields{
		"height":                 b.Header.Height,
		"hash":                   b.Hash(BlockHash{}),
		"number of transactions": len(b.Transactions),
	}).Info("adding new block")
}
