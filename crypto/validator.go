package crypto

import "fmt"

type Validator interface {
	ValidateBlock(*Block) error
}

// BlockValidator implements the Validator interface
type BlockValidator struct {
	bc *Blockchain
}

// NewBlockValidator initializes a new BlockValidator
func NewBlockValidator(bc *Blockchain) *BlockValidator {
	return &BlockValidator{bc: bc}
}

// ValidateBlock validates a new block before being added to the blockchain
// to ensure correctness
func (v *BlockValidator) ValidateBlock(b *Block) error {
	// no error means that we are attempting to add a block at a height in the blockchain
	// already occupied by an existing block
	if _, err := v.bc.GetHeaderByHeight(b.Header.Height); err == nil {
		return fmt.Errorf("the blockchain already contains block of height: %d with hash %s", b.Header.Height, b.headerHash)
	}

	// the new block to be added must occupy the next available slot i.e. blockchain height + 1
	// if the block height is greater than that, return error
	if b.Header.Height > v.bc.GetBlockchainHeight()+1 {
		return fmt.Errorf("the height of block with hash (%s) is too high", b.Hash(BlockHash{}))
	}

	// get header of the previous block in the blockchain
	prevHeader, err := v.bc.GetHeaderByHeight(b.Header.Height - 1)
	if err != nil {
		return err
	}

	// recalculate hash of previous block header
	prevHash := BlockHash{}.Hash(prevHeader)

	// the recalculated previous hash must match the PrevBlockHash of the new block
	// being added to the chain, otherwise treated as a fraudulent case
	// raise an error
	if prevHash != b.Header.PrevBlockHash {
		return fmt.Errorf("hash mismatch. expected previous hash: %s actual hash: %s", prevHeader.PrevBlockHash, prevHash)
	}

	// verify the new block
	if err := b.Verify(); err != nil {
		return err
	}

	return nil
}
