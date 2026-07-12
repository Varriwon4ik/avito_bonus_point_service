# Bonus Points Ledger Service

[![CI](https://github.com/Varriwon4ik/avito_bonus_point_service/actions/workflows/ci.yml/badge.svg)](https://github.com/Varriwon4ik/avito_bonus_point_service/actions/workflows/ci.yml)

> **Quick access** —
> **Try it:** [deployed trial instance](#deployment) (university network) or
> `docker compose up --build` locally ([Running](#running)) ·
> **Docs:** [hosted documentation site](https://varriwon4ik.github.io/avito_bonus_point_service/) ·
> **Handover:** [docs/customer-handover.md](./docs/customer-handover.md) ·
> **Contributing:** [CONTRIBUTING.md](./CONTRIBUTING.md) ·
> **AI agents:** [AGENTS.md](./AGENTS.md)

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
- Points can be accrued **in bulk** for promotional campaigns:
  `POST /v1/accruals/batch` processes each item independently and returns
  per-item results (HTTP 207), so one bad row never fails the whole campaign.
- Every request is **logged in structured form** (method, route, status,
  latency, and `user_id` where applicable) and a Prometheus `/metrics`
  endpoint exposes request counts/latencies plus ledger-level gauges.

## Architecture

- `points_lots` — one row per accrual ("batch" of points), tracks
  `amount`, `remaining` and `expires_at`.
- `holds` — two-phase reservations (`active` / `confirmed` / `cancelled`).
- `hold_allocations` — which lots a hold drew points from (used to release
  points on cancel).
- `ledger_entries` — append-only audit log of every balance-affecting event,
  with an optional user-facing `label` for accruals and an optional
  service-side `note` for internal annotations.
- `idempotency_keys` — caches the (status, body) of the first response for
  a given `(idempotency_key, endpoint)` pair.

## Running

```sh
docker compose up --build
```

This starts Postgres and the API on `http://localhost:8080`. The same
address serves:

- the web UI at `http://localhost:8080/`
- Swagger UI at `http://localhost:8080/docs`
- the OpenAPI spec at `http://localhost:8080/openapi.yaml`

To run locally without Docker:

```sh
export DB_DSN="postgres://bonus:bonus@localhost:5432/bonus_ledger?sslmode=disable"
export DEFAULT_TTL_DAYS=365   # optional, default lifetime of accrued points
export MIN_TTL_DAYS=1         # optional, lower bound for per-accrual ttl_days (US-08)
export MAX_TTL_DAYS=1825      # optional, upper bound for per-accrual ttl_days (US-08)
export HOLD_TIMEOUT_HOURS=24  # optional, holds unresolved for longer than this are auto-released
go run ./cmd/api
```

Accruals whose `ttl_days` falls outside the configured
`MIN_TTL_DAYS`–`MAX_TTL_DAYS` range are rejected with `400 Bad Request`.

Sanitized example environment values are available in [`.env.example`](./.env.example).

## Deployment

The current increment — the Week 6 trial / handover-candidate release
[`v2.1.0`](https://github.com/Varriwon4ik/avito_bonus_point_service/releases/tag/v2.1.0)
— is deployed on the University VM and is reachable at
`http://10.93.26.175:8080/` — serving the web UI, Swagger UI at `/docs`, and the
API. The address is on the university private network, so access requires the
university network/VPN; exact private access details for graders are provided
through Moodle. Self-hosting steps and the current handover state:
[docs/customer-handover.md](./docs/customer-handover.md).

## Documentation

- **Hosted documentation site (browsable):**
  <https://varriwon4ik.github.io/avito_bonus_point_service/> — the maintained
  `docs/` set published with MkDocs on every merge to `main`.
- **Architecture:** [docs/architecture/README.md](./docs/architecture/README.md)
  — static, dynamic, and deployment views plus the
  [ADR index](./docs/architecture/README.md#architecture-decision-records-adr-index).
- **Development process & configuration management:**
  [docs/development-process.md](./docs/development-process.md).
- Quality and testing: [docs/quality-requirements.md](./docs/quality-requirements.md),
  [docs/quality-requirement-tests.md](./docs/quality-requirement-tests.md),
  [docs/testing.md](./docs/testing.md),
  [docs/user-acceptance-tests.md](./docs/user-acceptance-tests.md),
  [docs/definition-of-done.md](./docs/definition-of-done.md).
- Planning: [docs/roadmap.md](./docs/roadmap.md),
  [docs/user-stories.md](./docs/user-stories.md).
- **Customer handover (access, transition state, self-hosting, limitations):**
  [docs/customer-handover.md](./docs/customer-handover.md).
- Contributing and agent guidance: [CONTRIBUTING.md](./CONTRIBUTING.md),
  [AGENTS.md](./AGENTS.md).

## Submissions

- **Week 6 (Assignment 6, Sprint 4 / trial release):** public report index at
  [reports/week6/README.md](./reports/week6/README.md). Week 6 trial /
  handover-candidate release:
  [v2.1.0](https://github.com/Varriwon4ik/avito_bonus_point_service/releases/tag/v2.1.0).
  New maintained assets: [docs/customer-handover.md](./docs/customer-handover.md),
  [CONTRIBUTING.md](./CONTRIBUTING.md), [AGENTS.md](./AGENTS.md).
- **Week 5 (Assignment 5, Sprint 3 / MVP v2):** public report index at
  [reports/week5/README.md](./reports/week5/README.md). Release mapped to MVP v2:
  [v2.0.0](https://github.com/Varriwon4ik/avito_bonus_point_service/releases/tag/v2.0.0).
  New maintained assets: [docs/architecture/](./docs/architecture/README.md) and
  [docs/development-process.md](./docs/development-process.md).
- **Week 4 (Assignment 4, Sprint 2):** public report index at
  [reports/week4/README.md](./reports/week4/README.md). Quality assets:
  [docs/quality-requirements.md](./docs/quality-requirements.md),
  [docs/quality-requirement-tests.md](./docs/quality-requirement-tests.md),
  [docs/testing.md](./docs/testing.md), and
  [docs/user-acceptance-tests.md](./docs/user-acceptance-tests.md).
- **Week 3 (Assignment 3, MVP v1):** public report index at
  [reports/week3/README.md](./reports/week3/README.md). Current backlog registry:
  [docs/user-stories.md](./docs/user-stories.md). Release mapped to MVP v1:
  [v1.0.0](https://github.com/Varriwon4ik/avito_bonus_point_service/releases/tag/v1.0.0).
  See also [CHANGELOG.md](./CHANGELOG.md), [docs/roadmap.md](./docs/roadmap.md),
  and [docs/definition-of-done.md](./docs/definition-of-done.md).
- **Week 2:** public index at [reports/week2/README.md](./reports/week2/README.md),
  MVP v0 write-up at [reports/week2/mvp-v0-report.md](./reports/week2/mvp-v0-report.md).

## API

All endpoints are unauthenticated (internal network only, per spec).

### Important HTTP status codes

- `200 OK` for successful reads, updates, confirms, cancels, and one-shot debits
- `201 Created` only when a new accrual or hold is created
- `400 Bad Request` for malformed JSON, missing required fields, invalid amounts, invalid hold IDs, or invalid query parameters
- `404 Not Found` for unknown users on read/debit/hold operations and unknown hold IDs
- `409 Conflict` for insufficient funds, invalid hold state, or duplicate in-progress idempotency keys
- `500 Internal Server Error` for unexpected server-side failures

Error responses use a consistent JSON envelope:

```json
{
  "error": "bad_request",
  "message": "amount must be a positive integer"
}
```

### Accrue points
```http
POST /v1/users/{id}/accruals
Content-Type: application/json

{
  "amount": 500,
  "ttl_days": 180,            // optional, defaults to DEFAULT_TTL_DAYS
  "idempotency_key": "order_12345",
  "label": "test"             // optional: "test", "real", or a custom short label
}
-> 201 { "lot_id": 1, "user_id": "user_123", "amount": 500, "expires_at": "..." }
-> 400 { "error": "bad_request", "message": "idempotency_key is required" }
```

If `label` is provided, the backend trims it, accepts predefined values such
as `test` and `real`, and rejects labels longer than 32 characters or labels
containing control/unsafe characters.

### Bulk accrual (US-01)
```http
POST /v1/accruals/batch
Content-Type: application/json

{
  "items": [
    { "user_id": "user_a", "amount": 100, "ttl_days": 30, "idempotency_key": "promo-1-a" },
    { "user_id": "user_b", "amount": 200, "idempotency_key": "promo-1-b", "label": "promo-july" },
    { "user_id": "",       "amount": 50,  "idempotency_key": "promo-1-c" }
  ]
}
-> 207 {
  "results": [
    { "index": 0, "status": "created", "user_id": "user_a", "result": { "lot_id": 7, ... } },
    { "index": 1, "status": "created", "user_id": "user_b", "result": { "lot_id": 8, ... } },
    { "index": 2, "status": "error", "error": "bad_request", "message": "user_id is required" }
  ]
}
```

Accrues points to many users in one request — for promotional campaigns that
would otherwise need thousands of individual calls. Each item carries its own
`idempotency_key` (so partial retries never double-apply) and its own optional
`ttl_days` and `label` with the same validation as single accruals. The
response is always `207 Multi-Status` with per-item `created`/`error` results:
one bad row does not fail the rest of the batch. The web UI's **Accrue Points**
tab has a matching "Bulk accrual" row editor. An empty `items` array returns
`400 Bad Request`.

### Get balance
```http
GET /v1/users/{id}/balance?expiring_within_days=7
-> 200 {
  "user_id": "user_123",
  "available": 1200,   // spendable now
  "held": 300,         // reserved by active holds
  "total": 1500,       // available + held
  "expiring_soon": 80  // available points expiring within `expiring_within_days`
}
-> 404 { "error": "not_found", "message": "user not found" }
```

### Create a hold (phase 1 of two-phase debit)
```http
POST /v1/users/{id}/holds
Content-Type: application/json

{ "amount": 200, "idempotency_key": "checkout_98765" }
-> 201 { "hold_id": 42, "user_id": "user_123", "amount": 200, "status": "active" }
-> 404 { "error": "not_found", "message": "user not found" }
-> 409 { "error": "conflict", "message": "insufficient available points" }
```

### Confirm / cancel a hold (phase 2)
```http
POST /v1/holds/{hold_id}/confirm  -> 200 { ..., "status": "confirmed" }
POST /v1/holds/{hold_id}/cancel   -> 200 { ..., "status": "cancelled" }   // points returned
POST /v1/holds/{hold_id}/confirm  -> 404 { "error": "not_found", "message": "hold not found" }
```
Both are idempotent: calling confirm/cancel again on an already
confirmed/cancelled hold returns the same result without side effects.

### One-shot debit (hold + confirm in a single call)
```http
POST /v1/users/{id}/debits
Content-Type: application/json

{ "amount": 100, "idempotency_key": "loyalty_redeem_1" }
-> 200 { "hold_id": 43, "user_id": "user_123", "amount": 100, "status": "confirmed" }
-> 404 { "error": "not_found", "message": "user not found" }
-> 409 { "error": "conflict", "message": "insufficient available points" }
```

### Lots and ledger (support tooling / UI)
```http
GET /v1/users/{id}/lots?page=1&offset=20&status=active
-> 200 {
  "user_id": "user_123",
  "page": 1,
  "offset": 20,
  "total": 3,
  "lots": [
    {
      "lot_id": 1,
      "user_id": "user_123",
      "amount": 500,
      "remaining": 300,
      "status": "active",
      "expires_at": "...",
      "created_at": "..."
    }
  ]
}

GET /v1/users/{id}/transactions?page=1&offset=20
-> 200 {
  "user_id": "user_123",
  "page": 1,
  "offset": 20,
  "total": 42,
  "entries": [
    {
      "id": 10,
      "user_id": "user_123",
      "type": "accrual",
      "amount": 500,
      "ref_type": "lot",
      "ref_id": 1,
      "label": "test",
      "created_at": "..."
    }
  ]
}
```

Lots and transaction history are paginated with the same contract:
`page` is the 1-based page number (default 1) and `offset` is the page size
(`1`–`500`, default `20`). `GET /v1/users/{id}/lots` also accepts an optional
`status` filter: `active`, `expired`, or `exhausted`. Invalid `page`,
`offset`, or `status` values return `400 Bad Request`. Transaction entries may
also include a service-side `note` field when the platform records an internal
annotation such as `auto-released: timeout`.

### Autotester (US-15 / US-17 / US-19)
```http
POST /v1/autotest/run
Content-Type: application/json

{ "label": "demo", "user_id": "demo-user", "amount": 100, "parallel_requests": 5, "mode": "single" }
-> 200 {
     "scenario": { "label": "demo", "user_id": "autotest-demo-user", "amount": 100,
                   "ttl_days": 365, "parallel_requests": 5, "ledger_label": "test" },
     "passed": true,
     "results": [
       { "name": "accrual correctness", "passed": true },
       { "name": "parallel accrual", "passed": true }
     ]
   }
```

Runs the built-in autotester against the live instance and returns a per-check
pass/fail report. It verifies accrual correctness and that a burst of parallel
accrual requests each produce a distinct lot with a consistent balance and
ledger. Only `amount` is required; `ttl_days` defaults to the server's configured
TTL and `parallel_requests` defaults to 5. The `user_id` is always forced under an
`autotest-` prefix so real accounts are never touched. This backs the **Autotester**
tab in the web UI and shares the `internal/autotest` engine with the `cmd/autotest`
console tool.

The optional `mode` field selects the test set (US-19): the default `single`
runs the original correctness and parallel-burst checks with one shared
idempotency key, while `multi_key` fires N parallel accruals with N **distinct**
idempotency keys and verifies each key applies exactly once with a consistent
balance and ledger. In the web UI this is the **Test mode** selector on the
Autotester tab.

## Observability

Every request passes through a logging/metrics middleware that never reads or
logs the request body, so no sensitive payload data ends up in logs. Each
request emits one structured log line:

```text
level=INFO msg=http_request method=GET path=/v1/users/{id}/balance status=200 latency_ms=3 bytes=112 user_id=user_123
```

The route `path` is the templated pattern (e.g. `/v1/users/{id}/balance`), not
the concrete path, which keeps log/metric label cardinality bounded.

### Metrics endpoint

```http
GET /metrics
```

Returns Prometheus text exposition format (unauthenticated, for an internal
scraper). It exposes:

- `http_requests_total{method,path,status}` — request counter
- `http_request_duration_seconds{method,path}` — latency histogram
  (`_bucket` / `_sum` / `_count`)
- `bonus_points_available`, `bonus_points_held`, `bonus_active_holds`,
  `bonus_lots_total`, `bonus_users_total` — ledger-level gauges

## Tests

Integration tests run against a real Postgres instance and cover
idempotency, hold/confirm/cancel, malformed JSON, missing required fields,
not-found handling, invalid amounts, invalid pagination, FIFO-by-expiry
ordering, concurrent holds (race-safety), and OpenAPI route availability.
Automated **quality requirement tests** ([QRT-001/002/003](docs/quality-requirement-tests.md))
verify balance-read latency, ledger integrity under concurrency, and per-module
coverage:

```sh
docker compose up -d postgres
export TEST_DATABASE_URL="postgres://bonus:bonus@localhost:5432/bonus_ledger?sslmode=disable"

go test ./... -race -count=1 -p 1 -coverpkg=./... -coverprofile=coverage.out
go test ./internal/api/ -race -run 'TestQRT' -v   # quality requirement tests only
bash scripts/coverage_gate.sh coverage.out         # per-module coverage gate (>=30%)
```

Full testing status, critical modules, and coverage: [docs/testing.md](docs/testing.md).

## Continuous integration

Every push and every pull request to `main` is gated by the
[CI workflow](.github/workflows/ci.yml) (US-14). On GitHub's hosted runners it
spins up a Postgres service container (matching `docker-compose.yml`), pins the
Go toolchain to the version declared in `go.mod`, and runs `go mod verify`,
`gofmt`, `go vet`, `go build ./...`, the full test suite with the race detector,
the automated quality requirement tests, a **per-module line-coverage gate**
(≥30% for `internal/data` and `internal/api`), and **`govulncheck`** (dependency
and standard-library vulnerability scan — the additional QA check). The default
branch is protected by a ruleset requiring a pull request and passing status
checks, so broken changes cannot be merged and `main` stays releasable.
