package crypto

import (
	"bytes"
	"encoding/gob"
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
