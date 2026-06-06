# Security Policy

## Reporting a vulnerability

Please report security issues **privately**. Do not open a public issue for vulnerabilities.

- Use GitHub's **private vulnerability reporting** (Security tab → "Report a vulnerability"),
  or
- Email **security@agennext.dev** with details and reproduction steps.

We aim to acknowledge within 3 business days and to provide a remediation timeline after
triage. Coordinated disclosure is appreciated.

## Scope

This platform's threat model (see [`docs/THREAT_MODEL.md`](docs/THREAT_MODEL.md)) centers on
prompt injection (direct and indirect), multi-tenant isolation, capability/scope escalation,
and supply-chain integrity. Reports in these areas are especially valuable.

## Supply chain

Dependencies and images follow [`docs/SUPPLY_CHAIN.md`](docs/SUPPLY_CHAIN.md): pure-stdlib
core, trusted sources, SBOMs, `govulncheck`/Grype scanning, and digest-pinned images. CI
fails on known High/Critical vulnerabilities.

## Supported versions

Pre-1.0: only the latest `0.x` minor receives security fixes.
