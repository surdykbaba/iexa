// Package registry implements the IXEA shared Registry (the directory / trust
// list). It holds only metadata about Member Nodes — their ID, name, network
// endpoint, published public key and the data spaces they participate in. It
// never stores or relays business data: messages travel node-to-node. This is
// what makes the network federated rather than a central data lake.
package registry

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"sort"
	"sync"
	"time"
)

// Member is a directory entry for a participating node.
type Member struct {
	ID         string   `json:"id"`
	Name       string   `json:"name"`
	Endpoint   string   `json:"endpoint"`   // base URL of the node
	PublicKey  string   `json:"publicKey"`  // base64 Ed25519 public key
	DataSpaces []string `json:"dataSpaces"` // e.g. ["invoicing","identity"]
	UpdatedAt  string   `json:"updatedAt"`
}

// Store is an in-memory directory. Production deployments would back this with
// a database and a signed trust list; the interface stays the same.
type Store struct {
	mu      sync.RWMutex
	members map[string]Member
}

// NewStore returns an empty directory.
func NewStore() *Store {
	return &Store{members: make(map[string]Member)}
}

// Put inserts or updates a member entry.
func (s *Store) Put(m Member) {
	s.mu.Lock()
	defer s.mu.Unlock()
	m.UpdatedAt = time.Now().UTC().Format(time.RFC3339)
	s.members[m.ID] = m
}

// Get resolves a member by ID.
func (s *Store) Get(id string) (Member, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	m, ok := s.members[id]
	return m, ok
}

// List returns all members sorted by ID.
func (s *Store) List() []Member {
	s.mu.RLock()
	defer s.mu.RUnlock()
	out := make([]Member, 0, len(s.members))
	for _, m := range s.members {
		out = append(out, m)
	}
	sort.Slice(out, func(i, j int) bool { return out[i].ID < out[j].ID })
	return out
}

// Routes wires the Registry HTTP API onto a mux (Go 1.22 pattern routing).
func (s *Store) Routes(mux *http.ServeMux) {
	mux.HandleFunc("GET /healthz", func(w http.ResponseWriter, _ *http.Request) {
		writeJSON(w, http.StatusOK, map[string]string{"status": "ok", "service": "ixea-registry"})
	})
	mux.HandleFunc("POST /v1/register", func(w http.ResponseWriter, r *http.Request) {
		var m Member
		if err := json.NewDecoder(r.Body).Decode(&m); err != nil {
			writeError(w, http.StatusBadRequest, "invalid member: "+err.Error())
			return
		}
		if m.ID == "" || m.Endpoint == "" || m.PublicKey == "" {
			writeError(w, http.StatusBadRequest, "id, endpoint and publicKey are required")
			return
		}
		s.Put(m)
		writeJSON(w, http.StatusOK, map[string]any{"status": "registered", "id": m.ID})
	})
	mux.HandleFunc("GET /v1/members", func(w http.ResponseWriter, _ *http.Request) {
		writeJSON(w, http.StatusOK, map[string]any{"members": s.List()})
	})
	mux.HandleFunc("GET /v1/members/{id}", func(w http.ResponseWriter, r *http.Request) {
		m, ok := s.Get(r.PathValue("id"))
		if !ok {
			writeError(w, http.StatusNotFound, "no such member")
			return
		}
		writeJSON(w, http.StatusOK, m)
	})
}

// Client is the node-side view of the Registry.
type Client struct {
	BaseURL string
	HTTP    *http.Client
}

// NewClient builds a registry client with a sane timeout.
func NewClient(baseURL string) *Client {
	return &Client{BaseURL: baseURL, HTTP: &http.Client{Timeout: 5 * time.Second}}
}

// Register publishes (or refreshes) this node's directory entry.
func (c *Client) Register(m Member) error {
	body, _ := json.Marshal(m)
	resp, err := c.HTTP.Post(c.BaseURL+"/v1/register", "application/json", bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("register: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("register: registry returned %s", resp.Status)
	}
	return nil
}

// Resolve looks up a counterparty by member ID.
func (c *Client) Resolve(id string) (Member, error) {
	var m Member
	resp, err := c.HTTP.Get(c.BaseURL + "/v1/members/" + id)
	if err != nil {
		return m, fmt.Errorf("resolve %q: %w", id, err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return m, fmt.Errorf("resolve %q: registry returned %s", id, resp.Status)
	}
	if err := json.NewDecoder(resp.Body).Decode(&m); err != nil {
		return m, fmt.Errorf("resolve %q: %w", id, err)
	}
	return m, nil
}

func writeJSON(w http.ResponseWriter, code int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	_ = json.NewEncoder(w).Encode(v)
}

func writeError(w http.ResponseWriter, code int, msg string) {
	writeJSON(w, code, map[string]string{"error": msg})
}
