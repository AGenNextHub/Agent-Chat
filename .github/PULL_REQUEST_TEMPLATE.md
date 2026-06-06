<!-- Keep PRs focused. Link the issue/ADR. -->

## What & why
<!-- What does this change and why? Link the issue/ADR. -->

## Pipeline stage
<!-- Which stage does this complete? Define / Design / Act / Hardening -->

## Checklist
- [ ] `make check` passes locally (vet + lint + `go test -race`)
- [ ] New/changed behaviour has tests; coverage maintained
- [ ] Inspectable — no hidden logic; traces/decisions remain attributable
- [ ] If a dependency was added: followed `docs/SUPPLY_CHAIN.md` (SBOM, scan, license)
- [ ] Docs updated (incl. ADR if architectural)
- [ ] Commits are Conventional + signed off (`-s`)
