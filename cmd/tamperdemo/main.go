// Command tamperdemo proves IEXA's tamper detection end-to-end. It registers a
// sender identity, sends a correctly signed invoice (accepted, verified), then
// sends the SAME message with the payload mutated after signing — which the
// receiver rejects with 401 because the signature no longer matches.
package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/surdykbaba/iexa/internal/message"
	"github.com/surdykbaba/iexa/internal/registry"
	"github.com/surdykbaba/iexa/internal/trust"
)

func main() {
	regURL := env("REGISTRY_URL", "http://localhost:8090")
	target := env("TARGET_ENDPOINT", "http://localhost:8102") // the receiving node

	reg := registry.NewClient(regURL)
	pub, priv, err := trust.GenerateKeyPair()
	if err != nil {
		log.Fatal(err)
	}

	// Publish our public key so the receiver can verify what we send.
	if err := reg.Register(registry.Member{
		ID: "tamper-co", Name: "Tamper Co", Endpoint: "http://localhost:9999",
		PublicKey: trust.EncodePublic(pub), DataSpaces: []string{"invoicing"},
	}); err != nil {
		log.Fatalf("register: %v", err)
	}

	original := json.RawMessage(`{"invoiceNumber":"INV-T-1","currency":"NGN","totals":{"grossAmount":537500}}`)
	env, err := message.New("tamper-co", "globex-ke", "invoicing", "application/json", original)
	if err != nil {
		log.Fatal(err)
	}
	if err := env.Sign(priv); err != nil {
		log.Fatal(err)
	}

	fmt.Println("==> 1. sending a correctly SIGNED invoice  (grossAmount 537,500)")
	send(target, env)

	fmt.Println("\n==> 2. TAMPERING: rewriting the amount to 999,999,999 after signing")
	fmt.Println("        (the signature is left unchanged — exactly what an attacker would do)")
	env.Payload = json.RawMessage(`{"invoiceNumber":"INV-T-1","currency":"NGN","totals":{"grossAmount":999999999}}`)
	send(target, env)
}

func send(target string, env *message.Envelope) {
	body, _ := json.Marshal(env)
	resp, err := http.Post(target+"/v1/receive", "application/json", bytes.NewReader(body))
	if err != nil {
		log.Fatalf("send: %v", err)
	}
	defer resp.Body.Close()
	b, _ := io.ReadAll(resp.Body)
	verdict := "ACCEPTED"
	if resp.StatusCode != http.StatusAccepted {
		verdict = "REJECTED"
	}
	fmt.Printf("    -> HTTP %d %s  %s", resp.StatusCode, verdict, string(b))
}

func env(k, d string) string {
	if v := os.Getenv(k); v != "" {
		return v
	}
	return d
}
