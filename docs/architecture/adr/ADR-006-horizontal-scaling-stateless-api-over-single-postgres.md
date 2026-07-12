# ADR-006: The API tier is stateless and horizontally scalable; PostgreSQL stays the single coordination point

- **Status:** Accepted
- **Date:** 2026-07-12
- **Quality requirements addressed:** [QR-001 — Balance read response time](../../quality-requirements.md#qr-001-balance-read-response-time), [QR-002 — Ledger integrity under concurrency](../../quality-requirements.md#qr-002-ledger-integrity-under-concurrency)

## Context

At the Week 6 sprint review (10 Jul 2026) the customer asked, as a condition
for a complete delivery
([#64](https://github.com/Varriwon4ik/avito_bonus_point_service/issues/64)),
that the architecture documentation state **explicitly** whether the service
can be scaled horizontally, with the reasoning. The delivered deployment runs
one `api` container and one Postgres via docker compose, and the docs named
missing horizontal scaling as "the main structural risk if traffic grows"
without saying whether it is *possible*.

The team audited every component for single-instance assumptions:

- **HTTP handlers and web UI** (`internal/api`, `cmd/api/web`): no sessions,
  no in-process caches of ledger state; every request is served from
  PostgreSQL. The web UI is static and talks to the public `/v1` API.
- **Concurrency control** (ADR-001): `SELECT ... FOR UPDATE` row locks live in
  the database, so two mutations for the same user serialize correctly whether
  they arrive at one replica or two.
- **Idempotency** (ADR-004): keys and cached first responses are rows in
  `idempotency_keys`, reserved and committed atomically with the mutation.
  A retry landing on a *different* replica hits the same unique constraint and
  gets the replayed response or `409` — replica-safe by construction.
- **Point expiry** (ADR-002): lazy, a query predicate — no background job to
  coordinate between replicas.
- **Hold-timeout sweeper** (`runHoldSweep` in `cmd/api`): runs in every
  replica. Shown harmless: each release re-checks the hold's status under a
  `FOR UPDATE` row lock and treats "already confirmed/cancelled by someone
  else" as a skip, not an error. N replicas mean duplicate scans and log
  lines, never a double release.
- **Startup migrations** (`data.Migrate`): the one real race found. Postgres
  can fail concurrent `CREATE TABLE IF NOT EXISTS` on a fresh database, so
  several replicas first-booting simultaneously could crash-loop.
- **`/metrics`** : the registry is in-memory and therefore **per replica**.
- **Autotest engine**: scenarios live in the `autotest_scenarios` table; runs
  drive the public API, so fan-out across replicas behind a load balancer is
  safe.

## Decision

- **The service is horizontally scalable at the API tier, by design, with
  conditions.** Multiple identical `api` replicas may run against one
  PostgreSQL; correctness (no double-spend, idempotent retries, FIFO-by-expiry,
  hold auto-release) is enforced entirely in the database and does not depend
  on replica count.
- **Close the migration race now:** `data.Migrate` takes a Postgres
  **advisory lock** (`pg_advisory_lock`) for the duration of the migration
  pass, so concurrent replica startups serialize instead of racing.
- **PostgreSQL intentionally remains a single instance** and the single
  coordination point. Read replicas, failover, or sharding are out of scope
  for this product stage; the customer's operators can add standard Postgres
  HA underneath without touching the application, provided all writes go to
  the primary.
- The **conditions and caveats** for running more than one replica are
  documented in the architecture README (load balancer required, per-replica
  `/metrics` scraping, per-user throughput still serialized by row locks) and
  mirrored in the customer handover's known limitations.

## Consequences and tradeoffs

- **(+)** The delivery answers the customer's question with evidence: scaling
  out is a deployment change (N replicas + a load balancer), not a code
  change.
- **(+)** The QR-002 integrity guarantee is replica-independent — it was
  designed into the database layer (ADR-001/ADR-004), and this analysis
  confirms no code path weakens it under multiple instances.
- **(+)** The sweeper needs no leader election; duplicate execution is
  idempotent, which keeps the operational model simple.
- **(−)** Scaling the API tier does not scale the database: Postgres remains
  the write bottleneck and the single point of failure, and mutations for one
  user always serialize on that user's row locks regardless of replica count.
  Replicas add throughput across *different* users and availability of the
  API tier only.
- **(−)** Prometheus must scrape each replica directly (not through the load
  balancer), since metrics are per-process; this is standard Prometheus
  practice but must be configured.
- **(−)** The bundled docker compose still describes the single-replica
  trial topology; a multi-replica deployment (load balancer, replica count)
  is the operator's composition and is not shipped or load-tested by the
  team.

## Links

- Triggering request: [#64](https://github.com/Varriwon4ik/avito_bonus_point_service/issues/64)
  (customer condition, Week 6 review)
- Related decisions: [ADR-001](ADR-001-postgres-row-locking-for-ledger-integrity.md),
  [ADR-002](ADR-002-lazy-expiry-and-fifo-by-expiry-consumption.md),
  [ADR-004](ADR-004-client-supplied-idempotency-keys.md),
  [ADR-005](ADR-005-single-binary-web-ui-and-compose-deployment.md)
- Evidence: [`internal/data/db.go`](../../../internal/data/db.go) (advisory-locked
  migrations), [`internal/data/store.go`](../../../internal/data/store.go)
  (DB-side locking and idempotency),
  [`internal/data/holds.go`](../../../internal/data/holds.go) (`ExpireStaleHolds`
  concurrency tolerance), [`cmd/api/main.go`](../../../cmd/api/main.go)
  (per-replica sweeper)
