// Package node implements an IXEA Member Node — the single secure gateway each
// participant runs. It connects once to the Registry, then signs and routes
// messages to any other participant the Registry can resolve. Messages travel
// directly node-to-node; the node verifies every inbound message against the
// sender's published key before accepting it. This is the "connect once, reach
// the whole network" four-corner model in miniature.
package node

import (
	"bytes"
	"crypto/ed25519"
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/surdykbaba/ixea/internal/message"
	"github.com/surdykbaba/ixea/internal/registry"
	"github.com/surdykbaba/ixea/internal/trust"
)

// Node is a running Member Node.
type Node struct {
	ID       string
	Name     string
	Endpoint string

	pub  ed25519.PublicKey
	priv ed25519.PrivateKey
	reg  *registry.Client

	mu    sync.Mutex
	inbox []message.Envelope
}

// New builds a node with a freshly generated identity key pair.
func New(id, name, endpoint string, reg *registry.Client) (*Node, error) {
	pub, priv, err := trust.GenerateKeyPair()
	if err != nil {
		return nil, err
	}
	return &Node{ID: id, Name: name, Endpoint: endpoint, pub: pub, priv: priv, reg: reg}, nil
}

// Register publishes this node's directory entry to the Registry.
func (n *Node) Register(dataSpaces []string) error {
	return n.reg.Register(registry.Member{
		ID:         n.ID,
		Name:       n.Name,
		Endpoint:   n.Endpoint,
		PublicKey:  trust.EncodePublic(n.pub),
		DataSpaces: dataSpaces,
	})
}

// sendRequest is the client-facing body for POST /v1/send.
type sendRequest struct {
	To          string          `json:"to"`
	DataSpace   string          `json:"dataSpace"`
	ContentType string          `json:"contentType"`
	Payload     json.RawMessage `json:"payload"`
}

// Routes wires the node HTTP API.
func (n *Node) Routes(mux *http.ServeMux) {
	mux.HandleFunc("GET /healthz", func(w http.ResponseWriter, _ *http.Request) {
		writeJSON(w, http.StatusOK, map[string]string{"status": "ok", "service": "ixea-node", "id": n.ID})
	})
	mux.HandleFunc("GET /v1/info", func(w http.ResponseWriter, _ *http.Request) {
		writeJSON(w, http.StatusOK, map[string]string{
			"id": n.ID, "name": n.Name, "endpoint": n.Endpoint,
			"publicKey": trust.EncodePublic(n.pub),
		})
	})
	mux.HandleFunc("POST /v1/send", n.handleSend)
	mux.HandleFunc("POST /v1/receive", n.handleReceive)
	mux.HandleFunc("GET /v1/inbox", n.handleInbox)
}

// handleSend resolves the recipient, signs an envelope and delivers it.
func (n *Node) handleSend(w http.ResponseWriter, r *http.Request) {
	var req sendRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request: "+err.Error())
		return
	}
	if req.To == "" || len(req.Payload) == 0 {
		writeError(w, http.StatusBadRequest, "to and payload are required")
		return
	}
	if req.ContentType == "" {
		req.ContentType = "application/json"
	}
	if req.DataSpace == "" {
		req.DataSpace = "invoicing"
	}

	recipient, err := n.reg.Resolve(req.To)
	if err != nil {
		writeError(w, http.StatusNotFound, err.Error())
		return
	}

	env, err := message.New(n.ID, req.To, req.DataSpace, req.ContentType, req.Payload)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	if err := env.Sign(n.priv); err != nil {
		writeError(w, http.StatusInternalServerError, "sign: "+err.Error())
		return
	}

	body, _ := json.Marshal(env)
	client := &http.Client{Timeout: 5 * time.Second}
	resp, err := client.Post(recipient.Endpoint+"/v1/receive", "application/json", bytes.NewReader(body))
	if err != nil {
		writeError(w, http.StatusBadGateway, "deliver: "+err.Error())
		return
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusAccepted {
		writeError(w, http.StatusBadGateway, fmt.Sprintf("recipient rejected message: %s", resp.Status))
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{
		"status": "delivered", "messageId": env.Header.MessageID, "deliveredTo": req.To,
	})
}

// handleReceive verifies the sender's signature, then accepts the message.
func (n *Node) handleReceive(w http.ResponseWriter, r *http.Request) {
	var env message.Envelope
	if err := json.NewDecoder(r.Body).Decode(&env); err != nil {
		writeError(w, http.StatusBadRequest, "invalid envelope: "+err.Error())
		return
	}
	sender, err := n.reg.Resolve(env.Header.From)
	if err != nil {
		writeError(w, http.StatusUnauthorized, "unknown sender: "+err.Error())
		return
	}
	pub, err := trust.DecodePublic(sender.PublicKey)
	if err != nil {
		writeError(w, http.StatusUnauthorized, "bad sender key: "+err.Error())
		return
	}
	if !env.Verify(pub) {
		writeError(w, http.StatusUnauthorized, "signature verification failed")
		return
	}
	n.mu.Lock()
	n.inbox = append(n.inbox, env)
	n.mu.Unlock()
	writeJSON(w, http.StatusAccepted, map[string]any{
		"status": "accepted", "verified": true, "messageId": env.Header.MessageID,
	})
}

// handleInbox lists verified messages this node has received.
func (n *Node) handleInbox(w http.ResponseWriter, _ *http.Request) {
	n.mu.Lock()
	defer n.mu.Unlock()
	writeJSON(w, http.StatusOK, map[string]any{"count": len(n.inbox), "messages": n.inbox})
}

func writeJSON(w http.ResponseWriter, code int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	_ = json.NewEncoder(w).Encode(v)
}

func writeError(w http.ResponseWriter, code int, msg string) {
	writeJSON(w, code, map[string]string{"error": msg})
}
