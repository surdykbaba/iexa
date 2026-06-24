# Contributing to IXEA

IXEA is built in the open by the community. Whether you write code, run
infrastructure, or shape policy, there's a way in. This guide gets you from zero
to your first contribution.

## Ways to contribute

- **Code** — the reference Registry, Member Node, specs tooling and SDKs.
- **Specs & RFCs** — propose changes to the message envelope or a data-space profile.
- **Docs** — quickstarts, architecture notes, translations.
- **Triage** — reproduce issues, review pull requests, test the demo on your platform.

You do not need permission to start. Pick a [good first issue](https://github.com/surdykbaba/ixea/labels/good%20first%20issue), comment that you're on it, and open a PR.

## Local setup

Prerequisites: **Go 1.22+**.

```bash
git clone https://github.com/surdykbaba/ixea.git
cd ixea
go build ./...        # compile everything
bash scripts/demo.sh  # run a 2-node network end-to-end
```

See [docs/QUICKSTART.md](docs/QUICKSTART.md) and
[docs/ARCHITECTURE.md](docs/ARCHITECTURE.md) for how the pieces fit together.

## Before you open a pull request

Run the same checks CI runs:

```bash
gofmt -l .            # must print nothing (run `gofmt -w .` to fix)
go vet ./...          # must pass
go build ./...        # must compile
bash scripts/demo.sh  # must complete with a verified message
```

Then:

1. Branch from `main` (`git checkout -b my-change`).
2. Keep changes focused — one logical change per PR.
3. Write a clear PR description: what changed and why.
4. Reference the issue it closes (`Closes #123`).

## Coding conventions

- Standard library first; add a dependency only with a clear reason.
- Match the style of the surrounding code; keep packages small and documented.
- Public types and functions get a doc comment.
- No secrets, keys, or personal data in the repo or tests.

## Proposing a spec change (RFC)

Changes to the envelope or a data-space profile go through a lightweight RFC:

1. Open an issue describing the problem and proposed change.
2. If there's appetite, submit a PR updating the relevant file under `specs/`
   plus a short rationale in the PR body.
3. Maintainers and the relevant working group review in the open.

## Governance

IXEA is a neutral, community-governed project. Technical decisions are made in
the open through issues, PRs and working groups. The standard belongs to the
community — not to any company or individual.

## Code of Conduct

Be respectful, assume good faith, and help newcomers. Harassment or exclusionary
behaviour is not tolerated. Report concerns to the maintainers.

## License

By contributing, you agree your contributions are licensed under the
[Apache License 2.0](LICENSE).
