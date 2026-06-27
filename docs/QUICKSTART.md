# Quickstart

Run a tiny IEXA network on your machine — a Registry and two Member Nodes — and
send a signed invoice between them.

## Prerequisites

- Go 1.22+

## One-command demo

```bash
bash scripts/demo.sh
```

This builds the binaries, starts a Registry (`:8090`) and two nodes
(`acme-ng` on `:8101`, `globex-ke` on `:8102`), sends a signed invoice from
`acme-ng` to `globex-ke`, and prints the recipient's inbox showing the message
arrived **verified**. Everything is torn down on exit.

## Tamper-detection demo

```bash
bash scripts/tamper-demo.sh
```

Sends a correctly signed invoice (accepted, `HTTP 202`), then the **same message
with its amount rewritten after signing** — which the receiver rejects with
`HTTP 401 signature verification failed`. Proof that no field can be altered in
transit without invalidating the signature.

## Run the pieces by hand

```bash
# 1. Registry (the shared directory / trust list)
go run ./cmd/registry            # listens on :8090

# 2. A node — connects once to the registry, then can reach the whole network
NODE_ID=acme-ng NODE_NAME="Acme Nigeria Ltd" NODE_ADDR=":8101" \
  NODE_ENDPOINT="http://localhost:8101" go run ./cmd/node

# 3. Another node
NODE_ID=globex-ke NODE_NAME="Globex Kenya PLC" NODE_ADDR=":8102" \
  NODE_ENDPOINT="http://localhost:8102" go run ./cmd/node
```

Send a message:

```bash
curl -X POST http://localhost:8101/v1/send -H 'content-type: application/json' -d '{
  "to": "globex-ke",
  "dataSpace": "invoicing",
  "payload": { "invoiceNumber": "INV-1", "currency": "NGN", "...": "see specs/invoice" }
}'

curl http://localhost:8102/v1/inbox     # verified message appears here
```

## API summary

**Registry**

| Method & path | Purpose |
|---|---|
| `POST /v1/register` | Publish/refresh a member entry |
| `GET /v1/members` | List the directory |
| `GET /v1/members/{id}` | Resolve one member |

**Member Node**

| Method & path | Purpose |
|---|---|
| `GET /v1/info` | This node's ID, endpoint and public key |
| `POST /v1/send` | Sign and deliver a message to another member |
| `POST /v1/receive` | Accept an inbound message (verifies signature) |
| `GET /v1/inbox` | List received, verified messages |

See [ARCHITECTURE.md](ARCHITECTURE.md) for how it fits together.
