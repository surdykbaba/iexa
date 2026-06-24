// Package trust holds the cryptographic identity primitives of the IXEA
// reference implementation. Every Member Node owns an Ed25519 key pair; its
// public key is published to the Registry so that any counterparty can verify
// the messages it signs. There is no central key escrow — keys stay at the node.
package trust

import (
	"crypto/ed25519"
	"crypto/rand"
	"encoding/base64"
	"fmt"
)

// GenerateKeyPair returns a fresh Ed25519 key pair for a node identity.
func GenerateKeyPair() (ed25519.PublicKey, ed25519.PrivateKey, error) {
	pub, priv, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		return nil, nil, fmt.Errorf("generate key: %w", err)
	}
	return pub, priv, nil
}

// EncodePublic renders a public key as standard base64 for transport/registry.
func EncodePublic(pub ed25519.PublicKey) string {
	return base64.StdEncoding.EncodeToString(pub)
}

// DecodePublic parses a base64 public key as published by the Registry.
func DecodePublic(s string) (ed25519.PublicKey, error) {
	b, err := base64.StdEncoding.DecodeString(s)
	if err != nil {
		return nil, fmt.Errorf("decode public key: %w", err)
	}
	if len(b) != ed25519.PublicKeySize {
		return nil, fmt.Errorf("invalid public key length: %d", len(b))
	}
	return ed25519.PublicKey(b), nil
}
