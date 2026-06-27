// Command node runs an IEXA Member Node.
//
//	NODE_ID         member identifier (default "node-1")
//	NODE_NAME       human-readable name (default = NODE_ID)
//	NODE_ADDR       listen address (default ":8101")
//	NODE_ENDPOINT   externally reachable base URL (default "http://localhost<NODE_ADDR>")
//	REGISTRY_URL    registry base URL (default "http://localhost:8090")
//	NODE_DATASPACES comma-separated data spaces (default "invoicing,identity")
package main

import (
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/surdykbaba/iexa/internal/node"
	"github.com/surdykbaba/iexa/internal/registry"
)

func main() {
	id := env("NODE_ID", "node-1")
	name := env("NODE_NAME", id)
	addr := env("NODE_ADDR", ":8101")
	endpoint := env("NODE_ENDPOINT", "http://localhost"+addr)
	regURL := env("REGISTRY_URL", "http://localhost:8090")
	dataSpaces := strings.Split(env("NODE_DATASPACES", "invoicing,identity"), ",")

	reg := registry.NewClient(regURL)
	n, err := node.New(id, name, endpoint, reg)
	if err != nil {
		log.Fatalf("create node: %v", err)
	}

	// Connect once: publish our entry to the Registry (retry while it boots).
	registerWithRetry(n, dataSpaces)

	mux := http.NewServeMux()
	n.Routes(mux)
	log.Printf("IEXA node %q listening on %s (endpoint %s)", id, addr, endpoint)
	if err := http.ListenAndServe(addr, mux); err != nil {
		log.Fatal(err)
	}
}

func registerWithRetry(n *node.Node, dataSpaces []string) {
	for attempt := 1; attempt <= 10; attempt++ {
		if err := n.Register(dataSpaces); err != nil {
			log.Printf("registry not ready (attempt %d/10): %v", attempt, err)
			time.Sleep(500 * time.Millisecond)
			continue
		}
		log.Printf("registered %q with registry", n.ID)
		return
	}
	log.Printf("WARNING: could not register %q after retries; continuing", n.ID)
}

func env(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}
