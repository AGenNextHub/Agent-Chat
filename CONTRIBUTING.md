# Contributing to AGenNext Chat

Thanks for your interest. This project follows a deliberate delivery pipeline and a few
non-negotiable principles; please read this before opening a PR.

## Delivery pipeline

> Define → Design → Align → Take Approval → Act (one pass) → Self-review → Sanity →
> End-to-End Check → Confirm

Non-trivial changes start with an issue (Define/Design) and, where architectural, an ADR
in [`docs/adr/`](docs/adr/). Align on the approach before large implementation.

## Ground rules (enforced)

- **No hidden logic / operate with clarity** — behaviour must be inspectable; prefer
  explicit, traceable code.
- **Capability is a contract** — new platform features are modeled as contract-bound
  capabilities (see [`docs/CAPABILITIES.md`](docs/CAPABILITIES.md)).
- **Supply chain** — adding a dependency requires the procedure in
  [`docs/SUPPLY_CHAIN.md`](docs/SUPPLY_CHAIN.md). The core stays pure-stdlib by default.
- **No license lock** — permissive/OSI licenses only.

## Local checks (make verification accessible)

```bash
make check     # tidy + vet + lint + test  (mirrors CI)
make run       # end-to-end demo
make cover     # coverage summary
```

All of `go vet`, `golangci-lint`, and `go test -race` must pass. New code needs tests; aim
to keep or improve coverage.

## Commits & PRs

- Use [Conventional Commits](https://www.conventionalcommits.org/) (`feat:`, `fix:`,
  `docs:`, `refactor:`, `test:`, `chore:`).
- Sign off commits (DCO): `git commit -s`.
- Keep PRs focused; fill in the PR template; link the issue/ADR.
- By contributing you agree your work is licensed under Apache-2.0.

## Code style

Idiomatic Go, `gofmt`-clean, exported identifiers documented. Keep packages small and
single-purpose; keep maintenance simple.
