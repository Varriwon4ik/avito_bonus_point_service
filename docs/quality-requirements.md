# Quality Requirements

This document defines the measurable, non-functional quality requirements (QRs)
for the Bonus Points Ledger Service. Each QR uses a distinct
[ISO/IEC 25010](https://iso25000.com/index.php/en/iso-25000-standards/iso-25010)
quality sub-characteristic, states a measurable scenario, explains why it
matters, and links to the automated quality requirement test(s) that verify it.

Quality requirements and their tests are **maintained product assets** (per
`Process_Requirements.md`). Later project work must keep them current and extend
them when product scope or risk changes, rather than disabling them after the
assignment is submitted.

Starting with Assignment 5, each quality requirement also links to the
[Architecture Decision Record(s)](architecture/README.md#architecture-decision-records-adr-index)
that address it, so the *requirement → decision → automated test* chain is
traceable in both directions.

- Linked tests: [docs/quality-requirement-tests.md](quality-requirement-tests.md)
- Architecture documentation and ADRs: [docs/architecture/README.md](architecture/README.md)
- Testing status & coverage: [docs/testing.md](testing.md)
- Completion standard: [docs/definition-of-done.md](definition-of-done.md)

## Scenario format

Each scenario follows the structure:

```text
When <source> <stimulus> under <environment>,
the <artifact> shall <response> within <response measure>.
```

## Summary

| QR | ISO/IEC 25010 sub-characteristic | Verified by | Addressed by ADR(s) |
|---|---|---|---|
| QR-001 | Time behaviour (Performance efficiency) | [QRT-001](quality-requirement-tests.md#qrt-001-balance-read-response-time) | [ADR-002](architecture/adr/ADR-002-lazy-expiry-and-fifo-by-expiry-consumption.md) |
| QR-002 | Integrity (Security) | [QRT-002](quality-requirement-tests.md#qrt-002-debit-integrity-under-concurrency) | [ADR-001](architecture/adr/ADR-001-postgres-row-locking-for-ledger-integrity.md), [ADR-004](architecture/adr/ADR-004-client-supplied-idempotency-keys.md) |
| QR-003 | Testability (Maintainability) | [QRT-003](quality-requirement-tests.md#qrt-003-critical-module-line-coverage) | [ADR-003](architecture/adr/ADR-003-layered-monolith-with-gated-critical-modules.md), [ADR-005](architecture/adr/ADR-005-single-binary-web-ui-and-compose-deployment.md) |

## QR-001: Balance read response time

**ISO/IEC 25010 sub-characteristic:** Time behaviour (Performance efficiency)

**Scenario:** When an internal client requests a user's points balance
(`GET /v1/users/{id}/balance`) under a warmed-up, production-like deployment,
the balance API shall respond within **200 ms for at least 95% of requests**
(p95), measured over a sample of 200 requests.

**Why this matters:** The balance endpoint is on the hot path of the customer's
storefront and checkout flow — the store reads a user's available points before
showing how many can be redeemed. Slow balance reads directly slow the purchase
experience for end users, so a bounded, measurable latency budget protects the
main user workflow rather than relying on a vague "fast" goal.

**Response measure:** p95 latency ≤ 200 ms over 200 sampled balance reads
(budget overridable in CI via `QRT_BALANCE_P95_BUDGET_MS`).

**Linked quality requirement tests:**
[QRT-001](quality-requirement-tests.md#qrt-001-balance-read-response-time)

**Linked architecture decisions:**
[ADR-002 — Lazy point expiry and FIFO-by-expiry consumption](architecture/adr/ADR-002-lazy-expiry-and-fifo-by-expiry-consumption.md)
keeps the balance read a single indexed query with no background-job coupling,
which is what makes this latency budget achievable.

## QR-002: Ledger integrity under concurrency

**ISO/IEC 25010 sub-characteristic:** Integrity (Security)

**Scenario:** When multiple debit requests for the **same user** are processed
concurrently and together request more points than are available, under the
standard transactional Postgres environment, the ledger shall **never allow the
total spent to exceed the available balance and shall never produce a negative
balance** — i.e. no lost updates and no double-spend. Exactly the number of
debits the balance can cover succeed; the rest are rejected with `409 Conflict`.

**Why this matters:** Bonus points are money-like. A lost update or double-spend
under concurrency would let a user redeem points they do not have, producing
financial and accounting loss for the customer and corrupting the audit ledger.
This is the core correctness guarantee of the service (row-level
`SELECT ... FOR UPDATE` locking), so it must be continuously verified — including
under the race detector.

**Response measure:** With an initial balance of 1000 and 40 concurrent debits of
100 each, exactly 10 succeed (`200`), the remainder return `409`, total spent
≤ 1000, and the final available balance equals `1000 − spent` and is ≥ 0.

**Linked quality requirement tests:**
[QRT-002](quality-requirement-tests.md#qrt-002-debit-integrity-under-concurrency)

**Linked architecture decisions:**
[ADR-001 — Serialize balance mutations with `SELECT ... FOR UPDATE`](architecture/adr/ADR-001-postgres-row-locking-for-ledger-integrity.md)
provides the no-lost-update / no-double-spend mechanism, and
[ADR-004 — Client-supplied idempotency keys](architecture/adr/ADR-004-client-supplied-idempotency-keys.md)
extends the same guarantee across client retries and crashes.

## QR-003: Critical module testability

**ISO/IEC 25010 sub-characteristic:** Testability (Maintainability)

**Scenario:** When a developer changes a critical product module
(`internal/data` — persistence, transactional balance mutations, FIFO-by-expiry
consumption, idempotency, business rules; or `internal/api` — HTTP handlers,
request/response contract, middleware) under the standard CI environment, that
module shall retain **at least 30% automated line coverage**, enforced by the CI
coverage gate which fails the build when a critical module drops below the
threshold.

**Why this matters:** The ledger's critical logic must be directly and
automatically verifiable so defects are detected before merge instead of in
production. A machine-enforced minimum keeps the safety net from silently eroding
as the product changes.

**Response measure:** Per-module line coverage ≥ 30% for `internal/data` and
`internal/api`, computed from the CI coverage profile and enforced by
[`scripts/coverage_gate.sh`](../scripts/coverage_gate.sh); the build fails
otherwise.

**Linked quality requirement tests:**
[QRT-003](quality-requirement-tests.md#qrt-003-critical-module-line-coverage)

**Linked architecture decisions:**
[ADR-003 — Layered monolith with coverage-gated critical modules](architecture/adr/ADR-003-layered-monolith-with-gated-critical-modules.md)
defines the module boundaries this requirement gates, and
[ADR-005 — One binary serves API, web UI, and autotester](architecture/adr/ADR-005-single-binary-web-ui-and-compose-deployment.md)
keeps all frontends on the same tested public API and a single shared autotest
engine.
