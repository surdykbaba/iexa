# IEXA — Interoperability and Exchange Alliance

**Together towards a federated and sovereign data infrastructure for Africa.**

IEXA is an open, community-governed framework for secure, federated data exchange
across Africa — invoicing, identity and trusted data. It is a neutral standard, a
trust framework, and a community. Not a product, not owned by any single company.

Inspired by the world's most successful interoperability initiatives —
[X-Road](https://x-road.global), [Gaia-X](https://gaia-x.eu),
[Peppol](https://peppol.org) and [MOSIP](https://mosip.io).

🌍 **Website:** https://surdykbaba.github.io/iexa/

## The reference implementation

A small, runnable realisation of the four-corner, federated model — written in Go,
standard-library only. It demonstrates the core trust properties; see
[docs/ARCHITECTURE.md](docs/ARCHITECTURE.md).

```bash
bash scripts/demo.sh
```

Boots a Registry and two Member Nodes, sends a **signed** invoice from one to the
other, and shows it arriving **verified**. Full walkthrough in
[docs/QUICKSTART.md](docs/QUICKSTART.md).

## Layout

| Path | What it is |
|------|-----------|
| `cmd/registry` | The shared **Registry** (directory / trust list) — metadata + public keys only. |
| `cmd/node` | A **Member Node** — the single gateway each participant runs. |
| `internal/message` | The signed **envelope** (build / sign / verify). |
| `internal/trust` | Ed25519 key primitives. |
| `specs/` | JSON Schema: the envelope + invoicing & identity profiles. |
| `scripts/demo.sh` | End-to-end local network demo. |
| `index.html` | The project website (self-contained). |
| `IEXA-Whitepaper.docx` | Concept & Founding Framework (v0.1 working draft). |

## The framework

- **Open Standards** — public, versioned specifications via an open RFC process.
- **Trust Framework** — verification, accreditation and conformance testing.
- **Federation** — decentralised member nodes; no single point of control, no central data store.

## Status

`v0.1` — a working draft and reference build, published openly for community
review. Developers, providers, regulators and young technologists all welcome.

## License

Free and open-source under the [Apache License 2.0](LICENSE).
