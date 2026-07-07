# ADR-005: One binary serves API, web UI, and autotester; deployed with docker compose

- **Status:** Accepted
- **Date:** 2026-07-04 (documented 2026-07-05)
- **Quality requirements addressed:** [QR-003 — Critical module testability](../../quality-requirements.md#qr-003-critical-module-testability)

## Context

MVP v2 (Sprint 3) added user-facing surfaces on top of the API: exact HTTP
response feedback in the browser (US-16), a web autotester tab (US-17), and
transaction labels (US-18). The team had to decide where these live. A separate
frontend service (SPA + its own container/deploy) would double the deployment
surface on a single university VM and split CI. Separately, the US-15 console
autotester and the new US-17 web autotester risked becoming two diverging
implementations of the same checks.

## Decision

- The web UI stays a **static, dependency-free HTML/JS page** embedded in and
  served by the same Go binary (`cmd/api/web`, served at `/`). It talks to the
  service exclusively through the same public `/v1` API that integrating stores
  use — no private backdoor endpoints. Swagger UI (`/docs`) and the OpenAPI
  spec are served the same way.
- The autotester logic is **extracted into one shared engine**
  (`internal/autotest`) with two thin frontends: the `cmd/autotest` console
  tool and the `POST /v1/autotest/run` endpoint backing the web UI tab. Both
  run identical checks. Every autotest run forces the target user under an
  `autotest-` prefix so real accounts are never touched.
- Deployment stays **docker compose** on the university VM: an `api` container
  plus a `postgres:16-alpine` container with a named volume, configured via
  environment variables (`DB_DSN`, `DEFAULT_TTL_DAYS`, `MIN_TTL_DAYS`,
  `MAX_TTL_DAYS`, `HOLD_TIMEOUT_HOURS`). Deploys are manual
  (`git pull && docker compose up --build`).

## Consequences and tradeoffs

- **(+)** One artifact to build, test, deploy, and roll back; `docker compose
  up --build` reproduces the whole product anywhere, which keeps UAT and grading
  access simple.
- **(+)** Because the UI and both autotester frontends consume the public API,
  every demo and autotest run exercises the same code paths CI tests — the
  shared engine is tested once and cannot diverge between console and web
  (supports the QR-003 goal that critical behaviour stays automatically
  verifiable; verified via
  [QRT-003](../../quality-requirement-tests.md#qrt-003-critical-module-line-coverage)
  and the autotest engine tests).
- **(+)** The `autotest-` user-prefix sandbox makes it safe to expose the run
  endpoint on an internal deployment without auth.
- **(−)** No continuous delivery: deploys depend on a team member with VM
  access; the deployed version can lag `main` (this happened during the
  Sprint 3 review). Documented in
  [docs/development-process.md](../../development-process.md); automation is a
  candidate improvement.
- **(−)** The UI is intentionally plain (no framework); complex UI features
  cost more to build by hand.
- **(−)** `POST /v1/autotest/run` executes real requests against the live
  instance; on an internet-exposed deployment it would need auth/rate limiting
  before this decision could stand unchanged.

## Links

- Related decision: [ADR-003 — layered monolith](ADR-003-layered-monolith-with-gated-critical-modules.md)
- Delivered by: US-16 ([#39](https://github.com/Varriwon4ik/avito_bonus_point_service/issues/39)),
  US-17 ([#40](https://github.com/Varriwon4ik/avito_bonus_point_service/issues/40)),
  US-18 ([#41](https://github.com/Varriwon4ik/avito_bonus_point_service/issues/41))
- Evidence: [`internal/autotest`](../../../internal/autotest/),
  [`internal/api/autotest_handler.go`](../../../internal/api/autotest_handler.go),
  [`docker-compose.yml`](../../../docker-compose.yml)
