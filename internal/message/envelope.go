// Package message defines the IEXA message envelope — the signed unit of
// exchange that travels between Member Nodes. The envelope carries a routing
// header and an opaque payload (the data-space profile, e.g. an invoice or an
// identity attestation). The signature covers the header and payload together,
// so neither routing metadata nor content can be tampered with in transit.
package message

import (
	"crypto/ed25519"
	"crypto/rand"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"time"
)

// Header is the routing and trust metadata of an envelope.
type Header struct {
	MessageID   string `json:"messageId"`
	From        string `json:"from"`        // sender member ID
	To          string `json:"to"`          // recipient member ID
	DataSpace   string `json:"dataSpace"`   // e.g. "invoicing", "identity"
	ContentType string `json:"contentType"` // e.g. "application/json"
	Timestamp   string `json:"timestamp"`   // RFC3339
}

// Envelope is the signed message exchanged between nodes.
type Envelope struct {
	Header    Header          `json:"header"`
	Payload   json.RawMessage `json:"payload"`
	Signature string          `json:"signature,omitempty"`
}

// New builds an unsigned envelope with a fresh message ID and timestamp.
func New(from, to, dataSpace, contentType string, payload json.RawMessage) (*Envelope, error) {
	id, err := newID()
	if err != nil {
		return nil, err
	}
	return &Envelope{
		Header: Header{
			MessageID:   id,
			From:        from,
			To:          to,
			DataSpace:   dataSpace,
			ContentType: contentType,
			Timestamp:   time.Now().UTC().Format(time.RFC3339),
		},
		Payload: payload,
	}, nil
}

// signingInput is the deterministic byte string that is signed and verified.
// It deliberately excludes the Signature field.
func (e *Envelope) signingInput() ([]byte, error) {
	return json.Marshal(struct {
		Header  Header          `json:"header"`
		Payload json.RawMessage `json:"payload"`
	}{e.Header, e.Payload})
}

// Sign attaches an Ed25519 signature over the header and payload.
func (e *Envelope) Sign(priv ed25519.PrivateKey) error {
	b, err := e.signingInput()
	if err != nil {
		return fmt.Errorf("signing input: %w", err)
	}
	e.Signature = base64.StdEncoding.EncodeToString(ed25519.Sign(priv, b))
	return nil
}

// Verify checks the signature against the sender's published public key.
func (e *Envelope) Verify(pub ed25519.PublicKey) bool {
	if e.Signature == "" {
		return false
	}
	b, err := e.signingInput()
	if err != nil {
		return false
	}
	sig, err := base64.StdEncoding.DecodeString(e.Signature)
	if err != nil {
		return false
	}
	return ed25519.Verify(pub, b, sig)
}

func newID() (string, error) {
	buf := make([]byte, 16)
	if _, err := rand.Read(buf); err != nil {
		return "", fmt.Errorf("message id: %w", err)
	}
	return "msg_" + hex.EncodeToString(buf), nil
}
