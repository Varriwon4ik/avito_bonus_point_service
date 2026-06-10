# Bonus Points Ledger Service

A REST-like service for managing an online store's bonus-points program,
addressing the problems of the original prototype:

- Points have a configurable **lifetime** and expire automatically (lazy
  expiry: expired points simply stop counting toward the balance).
- All balance mutations happen inside Postgres transactions with
  `SELECT ... FOR UPDATE` row locking, so concurrent accruals/holds/debits
  for the same user are serialized at the **database** level (no lost
  updates / "last write wins").
- Debits use **two-phase commit**: `hold` reserves points, `confirm`
  permanently spends them, `cancel` releases them back. If a calling
  service crashes after holding but before confirming, the hold can simply
  be cancelled later and the points are returned to the user.
- Every accrual, hold, confirm, cancel and debit is **idempotent** via a
  client-supplied `idempotency_key`: retried requests return the original
  result instead of being applied twice.
- Debits/holds always consume the **soonest-to-expire points first** (FIFO
  by `expires_at`), not in accrual order.

## Architecture

- `points_lots` — one row per accrual ("batch" of points), tracks
  `amount`, `remaining` and `expires_at`.
- `holds` — two-phase reservations (`active` / `confirmed` / `cancelled`).
- `hold_allocations` — which lots a hold drew points from (used to release
  points on cancel).
- `ledger_entries` — append-only audit log of every balance-affecting event.
- `idempotency_keys` — caches the (status, body) of the first response for
  a given `(idempotency_key, endpoint)` pair.

## Running

```sh
docker compose up --build
```

This starts Postgres and the API on `http://localhost:8080`. The same
address serves a small web UI (dashboard / accrue / holds & debit / lots /
transactions).

To run locally without Docker:

```sh
export DB_DSN="postgres://bonus:bonus@localhost:5432/bonus_ledger?sslmode=disable"
export DEFAULT_TTL_DAYS=365   # optional, default lifetime of accrued points
go run ./cmd/api
```

## API

All endpoints are unauthenticated (internal network only, per spec).

### Accrue points
```
POST /v1/users/{id}/accruals
{
  "amount": 500,
  "ttl_days": 180,            // optional, defaults to DEFAULT_TTL_DAYS
  "idempotency_key": "order_12345"
}
-> 201 { "lot_id": 1, "user_id": "user_123", "amount": 500, "expires_at": "..." }
```

### Get balance
```
GET /v1/users/{id}/balance?expiring_within_days=7
-> 200 {
  "user_id": "user_123",
  "available": 1200,   // spendable now
  "held": 300,         // reserved by active holds
  "total": 1500,       // available + held
  "expiring_soon": 80  // available points expiring within `expiring_within_days`
}
```

### Create a hold (phase 1 of two-phase debit)
```
POST /v1/users/{id}/holds
{ "amount": 200, "idempotency_key": "checkout_98765" }
-> 201 { "hold_id": 42, "user_id": "user_123", "amount": 200, "status": "active" }
-> 409 if available balance < amount
```

### Confirm / cancel a hold (phase 2)
```
POST /v1/holds/{hold_id}/confirm  -> 200 { ..., "status": "confirmed" }
POST /v1/holds/{hold_id}/cancel   -> 200 { ..., "status": "cancelled" }   // points returned
```
Both are idempotent: calling confirm/cancel again on an already
confirmed/cancelled hold returns the same result without side effects.

### One-shot debit (hold + confirm in a single call)
```
POST /v1/users/{id}/debits
{ "amount": 100, "idempotency_key": "loyalty_redeem_1" }
-> 200 { "hold_id": 43, "user_id": "user_123", "amount": 100, "status": "confirmed" }
-> 409 if available balance < amount
```

### Lots and ledger (for the UI / debugging)
```
GET /v1/users/{id}/lots          -> [ { lot_id, amount, remaining, expires_at, created_at }, ... ]
GET /v1/users/{id}/transactions  -> [ { id, type, amount, ref_type, ref_id, created_at }, ... ]
```

## Tests

Integration tests run against a real Postgres instance and cover
idempotency, hold/confirm/cancel, insufficient funds, FIFO-by-expiry
ordering and concurrent holds (race-safety):

```sh
docker compose up -d postgres
export TEST_DATABASE_URL="postgres://bonus:bonus@localhost:5432/bonus_ledger?sslmode=disable"
go test ./...
```
