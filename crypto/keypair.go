package crypto

import (
	"crypto/ed25519"
	"crypto/rand"
	"io"
)

const (
	privKeyLen   = 64 // length of the private key
	seedLen      = 32
	pubKeyLen    = 32 // length of the public key
	signatureLen = 64
)

type PrivateKey struct {
	key ed25519.PrivateKey
}

// ToBytes converts private key to byte array
func (p *PrivateKey) ToBytes() []byte {
	return p.key
}

// Sign uses private key to create a signature
func (p *PrivateKey) Sign(msg []byte) *Signature {
	return &Signature{
		Value: ed25519.Sign(p.key, msg),
	}
}

// GeneratePrivateKey creates a new private key
func GeneratePrivateKey() (*PrivateKey, error) {
	seed := make([]byte, seedLen)
	_, err := io.ReadFull(rand.Reader, seed)
	if err != nil {
		return nil, err
	}

	return &PrivateKey{
		key: ed25519.NewKeyFromSeed(seed),
	}, nil
}

// PublicKey returns public key from the private key
func (p *PrivateKey) PublicKey() *PublicKey {
	b := make([]byte, pubKeyLen)
	copy(b, p.key[32:])

	return &PublicKey{
		Key: b,
	}
}

type PublicKey struct {
	Key ed25519.PublicKey
}

func (p *PublicKey) ToBytes() []byte {
	return p.Key
}

type Signature struct {
	Value []byte
}

// ToBytes returns the byte array value of a public key
func (s *Signature) ToBytes() []byte {
	return s.Value
}

// BytesToSignature converts byte array to Signature type
func BytesToSignature(b []byte) *Signature {
	if len(b) != signatureLen {
		panic("invalid signature length, must be 64")
	}

	return &Signature{Value: b}
}

// Verify checks if the signature is valid
func (s *Signature) Verify(pubKey *PublicKey, msg []byte) bool {
	return ed25519.Verify(pubKey.Key, msg, s.Value)
}
