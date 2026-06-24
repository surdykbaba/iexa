# Architecture

IXEA's reference implementation is a deliberately small realisation of the
four-corner, federated model. It demonstrates the core trust properties; it is
not yet production infrastructure.

## Components

```
                        ┌──────────────────────┐
                        │       Registry       │   directory / trust list
                        │  (metadata + keys)   │   — no business data
                        └───────▲──────▲───────┘
                  register/     │      │   resolve
                  resolve       │      │
        ┌─────────────────┐     │      │     ┌─────────────────┐
        │  Member Node A  │─────┘      └─────│  Member Node B  │
        │  (acme-ng)      │                  │  (globex-ke)    │
        └───────▲─────────┘                  └────────▲────────┘
                │   signed envelope, node → node       │
        sender ─┘   (Ed25519, verified on receipt)     └─ receiver
```

- **Registry** — the shared directory. Holds only each member's ID, name,
  endpoint, **public key** and data spaces. It never stores or relays payloads.
- **Member Node** — the single gateway a participant runs. It connects once to
  the Registry, then signs outbound messages and verifies inbound ones.
- **Envelope** — the signed unit of exchange: a routing header plus an opaque
  payload (the data-space profile), with an Ed25519 signature over both.

## Trust model

1. Each node generates an Ed25519 key pair at startup; the private key never
   leaves the node.
2. The node publishes its public key to the Registry (`POST /v1/register`).
3. To send, a node resolves the recipient's endpoint from the Registry, signs
   an envelope, and POSTs it straight to the recipient — **no central relay**.
4. On receipt, a node resolves the *sender's* public key from the Registry and
   verifies the signature before accepting. A failed verification is rejected
   with `401`.

This yields the key federated properties: **data stays at the source**, every
exchange is **authenticated and tamper-evident**, and there is **no central
honeypot** of business data.

## What a production deployment adds

The reference build keeps things in-memory and HTTP-only so the model is easy to
read. A real network would add, without changing the shape:

- mutual TLS between nodes and a signed, distributed trust list
- a durable, replicated Registry backend and member accreditation workflow
- key rotation, revocation and hardware-backed key storage
- message persistence, delivery receipts and replay protection
- conformance testing against the published [specs](../specs) for certification

## Layout

```
cmd/registry      Registry entrypoint
cmd/node          Member Node entrypoint
internal/registry directory store, HTTP API and node-side client
internal/node     Member Node logic and HTTP API
internal/message  the signed envelope (build / sign / verify)
internal/trust    Ed25519 key primitives
specs/            JSON Schema: envelope + invoicing & identity profiles
scripts/demo.sh   end-to-end local network demo
```
