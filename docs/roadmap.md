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

## Sprint 2 (Assignment 4) — 22–28 Jun 2026 — delivered

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

## Now — Sprint 3 (MVP v2, Assignment 5) — 29 Jun–5 Jul 2026

[Sprint 3 milestone](https://github.com/Varriwon4ik/avito_bonus_point_service/milestone/3)

**Sprint Goal:** deliver MVP v2 by making the ledger directly usable and
verifiable from the web UI — exact HTTP response feedback, a browser-based
autotester, transaction labels, and enforced per-accrual TTL bounds — while
documenting the architecture, key decisions (ADRs), and the development process
so the product can keep evolving safely.

Delivered increment (16 SP, mapped to release `v2.0.0`):

- **US-08** (`Should Have`) — configurable per-accrual TTL validation and bounds
  (`MIN_TTL_DAYS` / `MAX_TTL_DAYS`, out-of-range requests get `400`).
- **US-16** (`Should Have`) — exact HTTP response codes surfaced in the web UI
  for accrual/debit operations.
- **US-17** (`Should Have`) — web Autotester tab backed by
  `POST /v1/autotest/run` and the shared `internal/autotest` engine.
- **US-18** (`Could Have`) — labels on transactions (preset `test`/`real` or
  custom), shown in the transactions view. First landing was reverted and the
  feature was re-landed fixed within the Sprint.

Also delivered: maintained architecture documentation with three views
([docs/architecture/README.md](architecture/README.md)), five ADRs,
[docs/development-process.md](development-process.md), and the hosted
documentation site.

**Scope note:** the earlier candidate plan for MVP v2 listed US-01 and US-02.
During Sprint 3 planning the team re-prioritized toward the web-UI feature set
above, because the customer's Sprint 2 feedback centred on seeing and verifying
behaviour directly (responses, autotests, labels); US-01/US-02 moved to
Sprint 4, as agreed with the customer at the Sprint 3 review.

## Next — Sprint 4 (candidate)

- **US-01** — bulk points accrual for promotional campaigns.
- **US-02** — list and audit a user's points lots (support tooling API).
- **US-19** — autotester scenarios with multiple idempotency keys in parallel
  ([#50](https://github.com/Varriwon4ik/avito_bonus_point_service/issues/50);
  requested by the customer at the Sprint 3 review).
- Customer-requested follow-ups still open:
  - Run the autotester/regression checks against a demo (earlier) version of
    the product to prove changes are valid (carried from Sprint 2/3).
  - Model point expiry as an explicit ledger transaction rather than silent
    removal, so expirations are auditable (follow-up on US-12; see
    [ADR-002](architecture/adr/ADR-002-lazy-expiry-and-fifo-by-expiry-consumption.md)).

## Later — backlog / under consideration

- Authentication / authorization for admin-facing operations. A first attempt at
  US-07 (manual-accrual admin auth) was reverted in Sprint 2; revisit with a
  hardened, fully tested design before re-introducing it.
- Architecture, quality, and process work that must continue: keep
  [docs/architecture/](architecture/README.md), the ADR set, and
  [docs/development-process.md](development-process.md) current as the product
  changes; keep all Assignment 4/5 CI gates active; consider raising the
  30% critical-module coverage floor and automating VM deployment
  ([ADR-005](architecture/adr/ADR-005-single-binary-web-ui-and-compose-deployment.md)).

## Out of scope (current)

- US-07 — `Removed` (implemented then reverted in Sprint 2 due to bugs). US-09 is
  no longer out of scope — it was reconsidered and delivered in Sprint 2.
- US-03, US-04, US-06 — removed; already satisfied by MVP v0 base functionality.

## Release mapping

| Release | Scope | Status |
|---|---|---|
| [`v1.0.0`](https://github.com/Varriwon4ik/avito_bonus_point_service/releases/tag/v1.0.0) | MVP v1 (US-05, US-10, US-11, US-12, US-13) | Sprint 1 — released |
| [`v1.1.0`](https://github.com/Varriwon4ik/avito_bonus_point_service/releases/tag/v1.1.0) | Sprint 2 increment (US-09, US-14, US-15) | Sprint 2 — released |
| [`v2.0.0`](https://github.com/Varriwon4ik/avito_bonus_point_service/releases/tag/v2.0.0) | MVP v2 (US-08, US-16, US-17, US-18 + architecture/process docs) | Sprint 3 — released |
