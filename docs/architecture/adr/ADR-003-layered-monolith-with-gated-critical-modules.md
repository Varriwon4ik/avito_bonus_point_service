# ADR-003: Layered monolith with coverage-gated critical modules

- **Status:** Accepted
- **Date:** 2026-06-27 (documented 2026-07-05)
- **Quality requirements addressed:** [QR-003 — Critical module testability](../../quality-requirements.md#qr-003-critical-module-testability)

## Context

The product is built by a four-person team on a weekly Sprint cadence. It needs
to stay easy to change, easy to review, and easy to test against a real
database, while still being deployable to a single university VM. A
microservice split would add network boundaries, deployment complexity, and
per-service CI for no current benefit; an unstructured single package would make
the blast radius of any change the whole codebase.

## Decision

Keep the product a **single Go module structured as a layered monolith** with
explicit, dependency-directed layers:

- `cmd/api`, `cmd/autotest` — thin entrypoints: flag/env parsing and process
  wiring only, no business logic.
- `internal/api` — the HTTP layer: routing, request validation, error envelope,
  observability middleware, metrics. Knows about HTTP, not SQL.
- `internal/data` — the persistence layer: transactions, locking, business
  rules, migrations. Knows about SQL, not HTTP.
- `internal/autotest` — the scenario engine, consuming the public HTTP API like
  any other client.
- `cmd/api/web` — a static, dependency-free HTML/JS UI served by the same
  binary, calling the same public API.

Declare `internal/data` and `internal/api` the **critical modules** and enforce
≥ 30% line coverage on each with a per-module CI gate
([`scripts/coverage_gate.sh`](../../../scripts/coverage_gate.sh)); the build
fails when either drops below the threshold. `cmd/*` shells are deliberately
not gated.

## Consequences and tradeoffs

- **(+)** High cohesion inside each layer and one-directional coupling
  (`api → data`, never the reverse) keep changes local: an API-contract change
  does not touch SQL, and a schema change does not touch handlers.
- **(+)** Both critical layers are independently testable against a real
  Postgres, which is what makes the QR-003 coverage scenario meaningful —
  enforced continuously by
  [QRT-003](../../quality-requirement-tests.md#qrt-003-critical-module-line-coverage).
- **(+)** One binary, one Dockerfile, one compose file — deployment and local
  reproduction stay trivial (see
  [ADR-005](ADR-005-single-binary-web-ui-and-compose-deployment.md)).
- **(−)** All features share one process and one release cadence; a crash or a
  hot loop affects every endpoint. Acceptable at current scale.
- **(−)** The 30% threshold is a floor, not a target — it prevents erosion but
  does not guarantee deep coverage; raising it is a candidate future decision.

## Links

- Verified by: [QRT-003](../../quality-requirement-tests.md#qrt-003-critical-module-line-coverage)
- Related decisions: [ADR-001](ADR-001-postgres-row-locking-for-ledger-integrity.md)
  (the `internal/data` transaction discipline),
  [ADR-005](ADR-005-single-binary-web-ui-and-compose-deployment.md) (deployment
  shape of the monolith)
- Evidence: [CI workflow](../../../.github/workflows/ci.yml), critical-module
  table in [docs/testing.md](../../testing.md#critical-modules-and-coverage)
