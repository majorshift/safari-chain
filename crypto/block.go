package crypto

import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"
	"errors"
	"fmt"
)

// Header structure
type Header struct {
	Version       uint32 // current block version
	PrevBlockHash Hash   // hash of the previous block
	MerkleRoot    Hash   // hash of transactions in the block
	Timestamp     int64  // when the block was created
	Height        uint32 // number of blocks in the blockchain - 1
}

// ToBytes converts Header to a byte slice
func (h *Header) ToBytes() []byte {
	buf := &bytes.Buffer{}
	enc := gob.NewEncoder(buf)
	enc.Encode(h)

	return buf.Bytes()
}

// Block structure
type Block struct {
	*Header
	Transactions []*Transaction // a list of transactions within the block
	Validator    *PublicKey     // public key of the validator that will add the block to the chain
	Signature    *Signature     // signature verifying block's authenticity
	headerHash   Hash           // cached hash value of the block's header; for quick access
}

func NewBlock(h *Header, txs []*Transaction) *Block {
	merkleRoot := ComputeMerkleRoot(txs)
	h.MerkleRoot = merkleRoot

	return &Block{
		Header:       h,
		Transactions: txs,
	}
}

// Sign uses the private key to sign the block header
func (b *Block) Sign(privKey *PrivateKey) {
	sig := privKey.Sign(b.Header.ToBytes())

	b.Validator = privKey.PublicKey()
	b.Signature = sig
}

// Hash hashes a block header
func (b *Block) Hash(hasher Hasher[*Header]) Hash {
	return hasher.Hash(b.Header)
}

// Verify checks the validity of a block
func (b *Block) Verify() error {
	// check that the header is signed
	if b.Signature == nil {
		return errors.New("block header has no signature")
	}

	// verify that the signature is valid i.e. no mismatch
	if !b.Signature.Verify(b.Validator, b.Header.ToBytes()) {
		return errors.New("block header has invalid signature")
	}

	// check the validity of all transaction signatures; all must be valid
	for _, tx := range b.Transactions {
		if err := tx.Verify(); err != nil {
			return err
		}
	}

	// verify merkle root
	merkleRoot := ComputeMerkleRoot(b.Transactions)
	if b.Header.MerkleRoot != merkleRoot {
		return fmt.Errorf("merkle root does not match")
	}

	return nil
}

// ComputeMerkleRoot calculates the merkle root of a given list of transactions
func ComputeMerkleRoot(transactions []*Transaction) Hash {
	if len(transactions) == 0 {
		return Hash{}
	}

	// Step 1: Create a list of hashes from the transactions
	var txHashes []Hash
	for _, tx := range transactions {
		txHashes = append(txHashes, tx.Hash(TxHash{}))
	}

	// calculate merkle root by repeatedly pairing hashes until one hash is left
	for len(txHashes) > 1 {
		var nextLevel []Hash
		for i := 0; i < len(txHashes); i += 2 {
			// If there's an odd number of hashes, duplicate the last one
			if i+1 >= len(txHashes) {
				nextLevel = append(nextLevel, hashPair(txHashes[i], txHashes[i]))
			} else {
				nextLevel = append(nextLevel, hashPair(txHashes[i], txHashes[i+1]))
			}
		}

		txHashes = nextLevel
	}

	// this is the merkle root
	return txHashes[0]
}

// hashPair combines 2 hashes
// returns a hash of the paired hashes
func hashPair(hash1, hash2 Hash) Hash {
	combined := append(hash1[:], hash2[:]...)
	h := sha256.Sum256(combined)
	return h
}
