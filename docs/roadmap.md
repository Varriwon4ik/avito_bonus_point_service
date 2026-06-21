# Product Roadmap

Bonus Points Ledger Service — a REST-like service that manages an online
store's bonus-points program (configurable expiry, transactional balance
mutations, two-phase debits, idempotency, observability).

This roadmap gives direction across releases. It is intentionally lighter than
the Product Backlog: see [docs/user-stories.md](user-stories.md) and the GitHub
issues for the authoritative, current state.

## Product Goal

Provide a reliable, auditable, and operationally observable bonus-points ledger
that downstream services can integrate with safely under concurrency, without
losing or double-counting points.

## Now — Sprint 1 (MVP v1) — 15–21 Jun 2026

**Sprint Goal:** complete the initial core features requested by the customer.

Delivered increment (all `Must Have`, marked MVP v1):

- **US-05** — auto-release stale holds after a configurable timeout.
- **US-11** — verified idempotent-key deduplication under concurrent requests.
- **US-12** — points removal / expiry system.
- **US-13** — correct HTTP response codes and OpenAPI documentation.

Also delivered in Sprint 1 (`Should Have`):

- **US-10** — structured request logging and a Prometheus `/metrics` endpoint.

## Next — Sprint 2 (candidate MVP v2)

Refine and pull from the remaining `Should Have` backlog:

- **US-01** — bulk points accrual for promotional campaigns.
- **US-02** — list and audit a user's points lots (support tooling API).
- **US-08** — configurable per-accrual TTL validation and bounds.
- Customer-requested follow-ups from the Sprint 1 review:
  - Apply concurrency tests to older code paths to show measurable impact
    (follow-up on US-11).
  - Model point expiry as an explicit ledger transaction rather than silent
    removal, so expirations are auditable (follow-up on US-12).

## Later — backlog / under consideration

- Authentication / authorization for admin-facing operations.
- Reconsider `Won't Have` stories US-07 (manual accrual) and US-09
  (transaction-history pagination) if customer priorities change.

## Out of scope (current)

- US-07, US-09 — `Won't Have` for the current increment.
- US-03, US-04, US-06 — removed; already satisfied by MVP v0 base functionality.

## Release mapping

| Release | Scope | Status |
|---|---|---|
| `v1.0.0` | MVP v1 (US-05, US-10, US-11, US-12, US-13) | Sprint 1 |
| `v2.0.0` (planned) | MVP v2 (US-01, US-02, US-08 + review follow-ups) | Sprint 2 |
