#!/usr/bin/env bash
# IEXA reference network demo.
# Boots a Registry and two Member Nodes, then sends a signed invoice from one to
# the other and shows it arriving verified. Tears everything down on exit.
set -euo pipefail
cd "$(dirname "$0")/.."

REG_ADDR=":8090"; REG_URL="http://localhost:8090"
A_ADDR=":8101";   A_URL="http://localhost:8101"
B_ADDR=":8102";   B_URL="http://localhost:8102"

echo "==> building binaries"
go build -o ./bin/registry ./cmd/registry
go build -o ./bin/node ./cmd/node

pids=()
cleanup() { echo; echo "==> shutting down"; for p in "${pids[@]:-}"; do kill "$p" 2>/dev/null || true; done; }
trap cleanup EXIT

echo "==> starting registry on $REG_ADDR"
REGISTRY_ADDR="$REG_ADDR" ./bin/registry & pids+=($!)

echo "==> starting node 'acme-ng' on $A_ADDR"
NODE_ID=acme-ng NODE_NAME="Acme Nigeria Ltd" NODE_ADDR="$A_ADDR" NODE_ENDPOINT="$A_URL" REGISTRY_URL="$REG_URL" ./bin/node & pids+=($!)

echo "==> starting node 'globex-ke' on $B_ADDR"
NODE_ID=globex-ke NODE_NAME="Globex Kenya PLC" NODE_ADDR="$B_ADDR" NODE_ENDPOINT="$B_URL" REGISTRY_URL="$REG_URL" ./bin/node & pids+=($!)

echo "==> waiting for services"
for url in "$REG_URL/healthz" "$A_URL/healthz" "$B_URL/healthz"; do
  curl -s --fail --retry 20 --retry-delay 1 --retry-connrefused -o /dev/null "$url"
done

echo; echo "==> registry directory"
curl -s "$REG_URL/v1/members" | sed 's/},{/},\n  {/g'

echo; echo "==> acme-ng sends a signed invoice to globex-ke"
curl -s -X POST "$A_URL/v1/send" -H 'content-type: application/json' -d '{
  "to": "globex-ke",
  "dataSpace": "invoicing",
  "payload": {
    "invoiceNumber": "INV-2026-0001",
    "issueDate": "2026-06-22",
    "currency": "NGN",
    "supplier": { "name": "Acme Nigeria Ltd", "memberId": "acme-ng", "country": "NG" },
    "customer": { "name": "Globex Kenya PLC", "memberId": "globex-ke", "country": "KE" },
    "lines": [ { "description": "Interoperability consulting", "quantity": 10, "unitPrice": 50000, "taxPercent": 7.5 } ],
    "totals": { "netAmount": 500000, "taxAmount": 37500, "grossAmount": 537500 },
    "payment": { "requestToPay": true, "reference": "INV-2026-0001" },
    "fiscal": { "taxId": "NG-12345678", "lotteryReceiptId": "LOT-0000001" }
  }
}'

echo; echo; echo "==> globex-ke inbox (should show the invoice, verified)"
curl -s "$B_URL/v1/inbox"
echo; echo; echo "==> demo complete."
