#!/usr/bin/env bash
# IEXA tamper-detection demo.
# Boots a Registry and one receiving Member Node, then sends a correctly signed
# invoice (accepted) followed by a payload-tampered copy (rejected with 401).
set -euo pipefail
cd "$(dirname "$0")/.."

REG_URL="http://localhost:8090"; B_URL="http://localhost:8102"

echo "==> building binaries"
go build -o ./bin/registry ./cmd/registry
go build -o ./bin/node ./cmd/node

pids=()
cleanup() { echo; echo "==> shutting down"; for p in "${pids[@]:-}"; do kill "$p" 2>/dev/null || true; done; }
trap cleanup EXIT

echo "==> starting registry and receiver node 'globex-ke'"
REGISTRY_ADDR=":8090" ./bin/registry & pids+=($!)
NODE_ID=globex-ke NODE_NAME="Globex Kenya PLC" NODE_ADDR=":8102" NODE_ENDPOINT="$B_URL" REGISTRY_URL="$REG_URL" ./bin/node & pids+=($!)

echo "==> waiting for services"
for url in "$REG_URL/healthz" "$B_URL/healthz"; do
  curl -s --fail --retry 20 --retry-delay 1 --retry-connrefused -o /dev/null "$url"
done

echo; echo "==> running tamper demonstration"; echo
REGISTRY_URL="$REG_URL" TARGET_ENDPOINT="$B_URL" go run ./cmd/tamperdemo
echo
