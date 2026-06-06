# Capabilities — *the capability is the contract*

The atomic, composable unit of AGenNext Chat is the **Capability**. A feature is not code
that happens to run; **a feature is a contract**, and that contract **is** the capability.
If it isn't declared in a contract, the platform will not run it.

## Why capability-as-contract

- **Composable** — capabilities snap together because each declares exactly what it
  provides, requires, and is allowed to touch.
- **No hidden logic** — the contract is the complete, inspectable description of behaviour
  and authority. Nothing runs outside it.
- **Least privilege** — `define scope` is enforced *by* the contract, not bolted on after.
- **Portable** — a capability ships as an **OCI artifact**; the contract travels with it.

## The contract (shape)

Every capability declares:

| Field | Meaning |
|---|---|
| `name` / `version` | identity (semver) |
| `provides` | interface(s) it exposes to other layers/capabilities |
| `requires` | capabilities/services it depends on |
| `scope` | data, tools, and tenants it may access (least privilege) |
| `policy` | **OPA** rules that must pass for admission/execution |
| `authz` | **OpenFGA** relations governing who may invoke it |
| `artifact` | **OCI** reference (digest-pinned) |
| `sandbox` | execution profile (default: isolated, restricted) |

A capability whose contract is unsatisfied — missing dependency, failing OPA policy, or
out-of-scope access — is **rejected by the Kernel at admission**, never partially run.

## Lifecycle (foundation-first ordering)

1. **Author the contract** (protocol/interface + scope + policy) — *protocol first*.
2. **Package** the implementation as an OCI artifact, digest-pinned to the contract.
3. **Publish** to the registry.
4. **Admit**: Kernel validates contract, OPA policy, and OpenFGA relations.
5. **Schedule & sandbox**: Kernel places it; Runtime Core runs it in its declared profile.
6. **Observe**: every invocation is attributable to a contract (auditable, no hidden logic).

## Worked example (illustrative)

```yaml
capability: rag.retrieve
version: 0.1.0
provides: [ "retrieve(query) -> passages" ]
requires: [ "index.search", "pii.screen" ]
scope:
  tenants: [ "$caller.tenant" ]      # never cross-tenant
  data:    [ "tenant://$caller.tenant/kb/*" ]
  tools:   [ ]                        # no tool/network access
policy:    opa://policies/rag-retrieve.rego
authz:     openfga://type/capability/relation/can_invoke
artifact:  oci://registry/agennext/rag-retrieve@sha256:...
sandbox:   isolated
```

This is a design contract, not yet an implemented schema; the canonical machine-readable
form will live alongside the Kernel CRDs.
