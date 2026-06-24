// Command registry runs the IXEA shared Registry (directory / trust list).
//
//	REGISTRY_ADDR   listen address (default ":8090")
package main

import (
	"log"
	"net/http"
	"os"

	"github.com/surdykbaba/ixea/internal/registry"
)

func main() {
	addr := env("REGISTRY_ADDR", ":8090")
	store := registry.NewStore()
	mux := http.NewServeMux()
	store.Routes(mux)

	log.Printf("IXEA registry listening on %s", addr)
	if err := http.ListenAndServe(addr, mux); err != nil {
		log.Fatal(err)
	}
}

func env(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}
