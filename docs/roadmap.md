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

## Sprint 1 (MVP v1) — 15–21 Jun 2026 — delivered

[Sprint 1 milestone](https://github.com/Varriwon4ik/avito_bonus_point_service/milestone/1)

**Sprint Goal:** complete the initial core features requested by the customer.

Delivered increment (all `Must Have`, marked MVP v1):

- **US-05** — auto-release stale holds after a configurable timeout.
- **US-11** — verified idempotent-key deduplication under concurrent requests.
- **US-12** — points removal / expiry system.
- **US-13** — correct HTTP response codes and OpenAPI documentation.

Also delivered in Sprint 1 (`Should Have`):

- **US-10** — structured request logging and a Prometheus `/metrics` endpoint.

## Now — Sprint 2 (Assignment 4) — 22–28 Jun 2026

[Sprint 2 milestone](https://github.com/Varriwon4ik/avito_bonus_point_service/milestone/2)

**Sprint Goal:** strengthen the reliability and verifiability of the increment by
delivering an automated autotester and a CI pipeline that gate every change,
while demonstrating paginated transaction-history access to the customer.

Delivered increment:

- **US-14** (`Must Have`) — continuous integration pipeline that builds and tests
  every push and pull request to `main`.
- **US-15** (`Must Have`) — `cmd/autotest` autotester that defines, stores, and
  replays reusable accrual / concurrency scenarios against a running instance.
- **US-09** (`Should Have`) — pagination for transaction history.

Reverted this Sprint: **US-07** (manual accrual) was implemented then reverted
due to bugs and is now `Removed`.

## Next — Sprint 3 (candidate MVP v2)

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

- Authentication / authorization for admin-facing operations. A first attempt at
  US-07 (manual-accrual admin auth) was reverted in Sprint 2; revisit with a
  hardened, fully tested design before re-introducing it.

## Out of scope (current)

- US-07 — `Removed` (implemented then reverted in Sprint 2 due to bugs). US-09 is
  no longer out of scope — it was reconsidered and delivered in Sprint 2.
- US-03, US-04, US-06 — removed; already satisfied by MVP v0 base functionality.

## Release mapping

| Release | Scope | Status |
|---|---|---|
| `v1.0.0` | MVP v1 (US-05, US-10, US-11, US-12, US-13) | Sprint 1 — released |
| `v1.1.0` (planned) | Sprint 2 increment (US-09, US-14, US-15) | Sprint 2 — pending release |
| `v2.0.0` (planned) | MVP v2 (US-01, US-02, US-08 + review follow-ups) | Sprint 3 |
