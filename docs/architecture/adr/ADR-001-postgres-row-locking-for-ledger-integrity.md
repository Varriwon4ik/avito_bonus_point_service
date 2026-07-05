# ADR-001: Serialize balance mutations in Postgres with `SELECT ... FOR UPDATE`

- **Status:** Accepted
- **Date:** 2026-06-20 (documented 2026-07-05)
- **Quality requirements addressed:** [QR-002 — Ledger integrity under concurrency](../../quality-requirements.md#qr-002-ledger-integrity-under-concurrency)

## Context

Bonus points are money-like: a lost update or double-spend directly costs the
customer money and corrupts the audit trail. The original prototype mutated
balances with read-modify-write logic in application code, so two concurrent
debits for the same user could both read the same balance and both succeed
("last write wins"). The service is deliberately small (one Go binary), may
later run as several replicas behind a load balancer, and all state already
lives in PostgreSQL.

Alternatives considered:

- **Application-level locking (mutex per user).** Free of DB round-trips, but
  only correct inside a single process — it silently breaks the moment a second
  replica is started, and it cannot protect against other writers (migrations,
  manual SQL).
- **Optimistic concurrency (version column + retry).** Correct, but pushes
  retry loops into every handler and makes the failure mode (retry storms under
  contention) harder to reason about than simple blocking.
- **Serializable isolation level.** Correct but coarser: it aborts transactions
  on false positives and still requires retry handling everywhere.

## Decision

Every balance-affecting operation (accrue, hold, confirm, cancel, debit) runs
inside **one PostgreSQL transaction** that first locks the affected rows with
`SELECT ... FOR UPDATE` (the user's spendable lots, or the hold row), then
applies all mutations (lot updates, hold state, allocations, ledger entries,
idempotency record) and commits atomically. Concurrency control lives in the
**database**, not in the application: concurrent mutations for the same user
serialize on the row locks, and mutations for different users do not block each
other. The implementation is concentrated in [`internal/data`](../../../internal/data/).

## Consequences and tradeoffs

- **(+)** No lost updates and no double-spend, regardless of how many API
  processes run — the guarantee holds at the storage layer. This is the direct
  mechanism behind QR-002 and is continuously verified by
  [QRT-002](../../quality-requirement-tests.md#qrt-002-debit-integrity-under-concurrency)
  under the Go race detector.
- **(+)** Simple mental model for the team: a handler either commits a fully
  consistent state or nothing.
- **(−)** Mutations for one *hot* user serialize; throughput per user is bounded
  by lock wait time. Acceptable for the product's traffic profile (a person
  checks out one order at a time).
- **(−)** Tests that share one database must not run in parallel across
  packages — CI runs `go test -p 1` for this reason.
- **(−)** Long transactions would hold locks; handlers must keep transaction
  bodies short (no external calls inside a transaction).

## Links

- Verified by: [QRT-002](../../quality-requirement-tests.md#qrt-002-debit-integrity-under-concurrency)
- Related decision: [ADR-004 — idempotency keys](ADR-004-client-supplied-idempotency-keys.md)
  (the idempotency record is written inside the same transaction)
- Evidence: [`internal/api/concurrent_idempotency_test.go`](../../../internal/api/concurrent_idempotency_test.go),
  [`internal/api/qrt_test.go`](../../../internal/api/qrt_test.go)
