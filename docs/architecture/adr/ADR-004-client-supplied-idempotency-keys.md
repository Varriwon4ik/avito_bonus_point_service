# ADR-004: Client-supplied idempotency keys with cached first responses

- **Status:** Accepted
- **Date:** 2026-06-20 (documented 2026-07-05)
- **Quality requirements addressed:** [QR-002 — Ledger integrity under concurrency](../../quality-requirements.md#qr-002-ledger-integrity-under-concurrency)

## Context

Integrating store services call the ledger over the network and **will** retry:
timeouts, crashed checkout processes, and at-least-once message delivery all
produce duplicate requests. Without protection, a retried accrual doubles a
user's points and a retried debit spends them twice. The caller — not the
ledger — knows which business event a request belongs to (an order, a checkout
attempt), so deduplication must key on caller-provided identity, not on request
timing.

Alternatives considered:

- **Best-effort dedup by payload hash + time window.** Heuristic; legitimately
  identical requests (two 100-point accruals for the same user) would be
  wrongly merged.
- **Server-generated tokens fetched in a pre-step.** Correct but forces a
  two-round-trip protocol on every caller and still fails when the caller
  crashes between the steps.

## Decision

Every mutating endpoint (accrue, hold, confirm, cancel, debit) **requires** a
client-supplied `idempotency_key`. The first execution stores
`(idempotency_key, endpoint) → (status code, response body)` in the
`idempotency_keys` table **inside the same transaction** as the mutation
(ADR-001); any retry — sequential or concurrent — finds the record and gets the
original response replayed with no second side effect. A key that is currently
in progress produces `409 Conflict` rather than a second execution.

## Consequences and tradeoffs

- **(+)** At-most-once application of every balance mutation under retries and
  crashes — together with row locking, this is the second half of the QR-002
  no-double-spend guarantee, verified by
  [QRT-002](../../quality-requirement-tests.md#qrt-002-debit-integrity-under-concurrency)
  and the dedicated concurrent-duplicate tests (US-11).
- **(+)** Recovery after a caller crash is trivial: retry with the same key, or
  cancel the hold (two-phase debit) — no reconciliation protocol needed.
- **(+)** Committing the idempotency record atomically with the mutation means
  there is no window where the mutation applied but the key is unknown.
- **(−)** Callers must generate and persist good keys (e.g. `order_12345`);
  a caller that reuses a key across different intents gets a stale replay. This
  is an integration contract documented in the README and OpenAPI spec.
- **(−)** The `idempotency_keys` table grows with write traffic; entries are
  never read after callers stop retrying. A TTL-based cleanup is a candidate
  future decision.

## Links

- Verified by: [QRT-002](../../quality-requirement-tests.md#qrt-002-debit-integrity-under-concurrency),
  [`internal/api/concurrent_idempotency_test.go`](../../../internal/api/concurrent_idempotency_test.go)
- Related decision: [ADR-001 — row locking](ADR-001-postgres-row-locking-for-ledger-integrity.md)
- Follow-up: the customer asked (Sprint 3 review) for an autotester scenario
  exercising multiple distinct keys in parallel —
  [US-19 / #50](https://github.com/Varriwon4ik/avito_bonus_point_service/issues/50)
