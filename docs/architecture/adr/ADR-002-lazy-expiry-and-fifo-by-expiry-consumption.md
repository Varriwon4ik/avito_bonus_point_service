# ADR-002: Lazy point expiry and FIFO-by-expiry consumption

- **Status:** Accepted
- **Date:** 2026-06-20 (documented 2026-07-05)
- **Quality requirements addressed:** [QR-001 — Balance read response time](../../quality-requirements.md#qr-001-balance-read-response-time)

## Context

Points have a configurable lifetime (US-12): after `ttl_days` they must stop
counting toward the balance. Two design questions follow: *how* expired points
are removed, and *which* points a debit consumes first.

Alternatives considered for expiry:

- **Eager expiry via a background sweeper** that periodically rewrites expired
  lots. Keeps tables "clean" but introduces a second writer that races with
  user traffic, needs its own scheduling/locking, and creates a window where a
  crashed sweeper silently makes balances wrong.
- **Lazy expiry (chosen):** expiry is a *query predicate* — an expired lot is
  simply never selected as spendable.

Alternatives considered for consumption order:

- **Accrual order (insertion FIFO):** simple, but lets soon-to-expire points rot
  while older long-lived points are spent — users lose value they could have
  used.
- **FIFO by `expires_at` (chosen):** always consume the soonest-to-expire
  points first.

## Decision

- A lot is spendable iff `remaining > 0 AND expires_at > now()`. Balance reads
  and debit/hold selection apply this predicate directly; nothing has to run for
  points to expire ("lazy expiry").
- Debits and holds consume lots in `ORDER BY expires_at ASC` — soonest-to-expire
  first — recording per-lot draws in `hold_allocations` so a cancel can return
  exactly what was taken.
- The balance endpoint additionally reports `expiring_soon` within a
  client-chosen window, so integrating stores can nudge users to spend points
  before they expire.

## Consequences and tradeoffs

- **(+)** Balance reads stay a single indexed aggregate query with no
  background-job coupling — this is what keeps the p95 ≤ 200 ms budget of
  QR-001 realistic, verified continuously by
  [QRT-001](../../quality-requirement-tests.md#qrt-001-balance-read-response-time).
- **(+)** No sweeper process to operate, monitor, or recover; correctness does
  not depend on a scheduler.
- **(+)** Users get maximum value from their points (soonest-to-expire spent
  first).
- **(−)** Expired lots remain in `points_lots` until archived; table growth is
  bounded only by accrual volume. Acceptable at current scale; revisit with an
  archival job if lot counts grow.
- **(−)** Expiry is currently **silent** — it is not written as an explicit
  ledger transaction, so an auditor must derive expirations from lot data. The
  customer flagged this in the Sprint 1 review; modelling expiry as an explicit
  ledger entry is a tracked backlog follow-up (see
  [docs/roadmap.md](../../roadmap.md)).

## Links

- Verified by: [QRT-001](../../quality-requirement-tests.md#qrt-001-balance-read-response-time)
- Related decision: [ADR-001 — row locking](ADR-001-postgres-row-locking-for-ledger-integrity.md)
  (lot selection and locking happen in the same `FOR UPDATE` query)
- Evidence: FIFO-by-expiry integration tests in
  [`internal/api/integration_test.go`](../../../internal/api/integration_test.go)
