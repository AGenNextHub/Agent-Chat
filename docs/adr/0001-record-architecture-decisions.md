# 1. Record architecture decisions

Date: 2026-06-06

## Status
Accepted

## Context
We need a durable, reviewable record of significant architectural decisions, consistent
with the "no hidden logic / operate with clarity" tenet and SDLC best practice.

## Decision
We use Architecture Decision Records (ADRs), one Markdown file per decision in `docs/adr/`,
numbered sequentially. Each records Status, Context, Decision, and Consequences. Format
after Michael Nygard's ADR pattern.

## Consequences
Architectural changes are proposed and reviewed as ADRs during the Define/Design/Align
stages before implementation. The history of *why* is preserved alongside the code.
