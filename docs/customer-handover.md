# Customer Handover

Maintained handover documentation for the **Bonus Points Ledger Service**
(Team 01). This document describes the **current actual handover state** of
the product — what the customer has access to today, what the team still
operates, and what the customer must know and be able to do to use, run, and
verify the ledger without the team. It is updated whenever access details,
deployment steps, limitations, or transition status change.

> **Where we are now (Week 6, trial stage):** the Week 6 trial /
> handover-candidate release
> [`v2.1.0`](https://github.com/Varriwon4ik/avito_bonus_point_service/releases/tag/v2.1.0)
> is available for the customer to try independently. The team still operates
> the deployed instance. The final transition, the reached handover level, and
> the customer-confirmation status are completed and recorded here during
> Week 7.

## What is transferred, delegated, or retained

| Arrangement | Current state |
|---|---|
| Source code and history | **Transferred (public).** The full monorepo — product code, migrations, documentation, CI configuration, reports — is public at <https://github.com/Varriwon4ik/avito_bonus_point_service> under the [MIT License](../LICENSE). The customer (and anyone) can clone, fork, build, and redistribute team-created content without asking the team. |
| Repository ownership / administration | **Retained by the team** until the course is graded: the repository, branch-protection rules, milestones, releases, and CI configuration stay under the team's GitHub account so grading evidence is preserved. Nothing blocks the customer from forking today; ownership transfer of the original repository can be arranged after grading if the customer wants it. |
| Deployed trial instance | **Operated by the team (delegated use).** The trial runs on a university VM at `http://10.93.26.175:8080/` (university network/VPN only). The team deploys releases to it and keeps it running until grading is complete. It is a trial/evaluation environment, not a production commitment. |
| Customer-side deployment | **Not yet in place — by the customer's own plan.** The product is fully self-hostable with Docker (see below). At the Week 6 review (10 Jul 2026) the customer stated they will not run the university-VM instance: their large-company environment has a security and deployment bar the trial setup does not target, and their internal deployment structure is confidential. After the final delivery, the customer's own staff (a group of interns) will evaluate the project and perform the customer-side deployment. The team's job is therefore to make this document and the repository sufficient for that team to deploy without us. |
| External services and accounts | **None to hand over.** The product depends only on PostgreSQL, which runs in the customer's own Docker environment; there are no third-party SaaS accounts, API keys, or paid services involved. |
| Documentation | **Transferred (public).** The maintained documentation set lives in the repository and is published as a browsable site at <https://varriwon4ik.github.io/avito_bonus_point_service/>; it stays available as long as the public repository exists. |

## Configuration, environment variables, and secrets

Runtime configuration is **environment variables only** — there are no runtime
config files. Everything the customer must know:

| Variable | Meaning | Default |
|---|---|---|
| `DB_DSN` | PostgreSQL connection string for the API. | — (required; docker compose sets it for the bundled Postgres) |
| `DEFAULT_TTL_DAYS` | Default lifetime of accrued points when a request omits `ttl_days`. | `365` |
| `MIN_TTL_DAYS` / `MAX_TTL_DAYS` | Accepted bounds for per-accrual `ttl_days`; out-of-range requests get `400`. | `1` / `1825` |
| `HOLD_TIMEOUT_HOURS` | Active holds unresolved for longer than this are auto-released. | `24` |
| `TEST_DATABASE_URL` | Used **only** by the test suite; tests self-skip without it. | — |

Secrets handling:

- The repository contains **no production secrets**. The only credentials in
  the repo are the throwaway Postgres credentials (`bonus:bonus`) used by the
  local/CI docker compose setup — replace them (set your own values in
  `docker-compose.yml` / `DB_DSN`) for any real deployment.
- A sanitized [`.env.example`](../.env.example) documents safe local values;
  real `.env` files are git-ignored and must never be committed.
- The API itself is **unauthenticated by design** (internal-network use per
  the original spec): deploy it only on a private network, never on the public
  internet. This is the product's most important operational constraint.
- Private access details for the university trial VM (network/VPN specifics,
  grader credentials) are intentionally **not** in this public document; they
  are shared through the private submission channel.

## Setup, deployment, recovery, and verification

Steps the customer must be able to follow on their own machine or server
(details and API examples: [root README](../README.md)):

1. **Run it:** `docker compose up --build` — starts PostgreSQL and the API on
   `http://localhost:8080` (web UI at `/`, Swagger UI at `/docs`, OpenAPI spec
   at `/openapi.yaml`, Prometheus metrics at `/metrics`). Database migrations
   apply automatically at startup.
2. **Verify it works:** open `/healthz` (expect `{"status":"ok"}`), then open
   the web UI **Autotester** tab and run a scenario (single-key or the
   multi-key US-19 mode) — the built-in autotester performs real accruals
   against a dedicated `autotest-`-prefixed user and reports per-check
   pass/fail results. The same checks run from the console:
   `go run ./cmd/autotest`.
3. **Run the full test suite** (optional, requires Go):
   `docker compose up -d postgres`, set `TEST_DATABASE_URL` (value in
   [`.env.example`](../.env.example)), then
   `go test ./... -race -count=1 -p 1`. The suite covers idempotency,
   two-phase debits, concurrency/race-safety, pagination, and the automated
   quality requirement tests.
4. **Recover / reset:** the ledger state lives in the `postgres` Docker
   volume. `docker compose down` + `up` restarts without data loss;
   `docker compose down -v` resets to a clean database. Holds left behind by a
   crashed calling service need no manual repair — they auto-release after
   `HOLD_TIMEOUT_HOURS` and the points return to the user with an audit ledger
   entry.
5. **Upgrade:** pull the desired release tag, then
   `docker compose up --build -d`; migrations are forward-only and applied
   automatically.

## Main documentation entry points

- **Normal use (integration and administration):**
  [root README](../README.md) — setup, full API reference with examples, web
  UI overview; interactive [Swagger UI](http://10.93.26.175:8080/docs) on the
  deployed trial instance (or `/docs` on any self-hosted instance).
- **Operation:** [Development process → configuration and secrets](development-process.md#configuration-and-secrets-management)
  and [Architecture → deployment view](architecture/README.md) — how the
  single-binary + compose deployment works;
  [observability](../README.md#observability) — structured logs and
  Prometheus metrics.
- **Troubleshooting and verification:** the web UI shows the exact HTTP
  response code of every operation (US-16); the
  [Autotester](../README.md#autotester-us-15--us-17--us-19) verifies core
  behaviour end-to-end; [docs/user-acceptance-tests.md](user-acceptance-tests.md) holds
  the customer-executable acceptance scenarios;
  [docs/testing.md](testing.md) records the automated test and CI gate
  status.
- **Known limitations:** see the next section.
- **Everything, browsable:**
  <https://varriwon4ik.github.io/avito_bonus_point_service/>.

## Known limitations

- **No authentication or authorization** — internal-network deployment only
  (per spec; an admin-auth attempt, US-07, was reverted and intentionally not
  re-introduced within the course).
- **Point expiry is lazy** — expired points stop counting immediately but are
  not written as explicit ledger transactions; expiry is therefore not
  itemized in the audit history
  ([ADR-002](architecture/adr/ADR-002-lazy-expiry-and-fifo-by-expiry-consumption.md);
  an explicit-expiry-transaction follow-up is a Sprint 5 candidate).
- **Deployment to the university VM is manual** — the deployed trial can
  briefly lag `main` between merge and deploy
  ([ADR-005](architecture/adr/ADR-005-single-binary-web-ui-and-compose-deployment.md)).
- **The trial VM is university-network-only** and is kept only until grading
  is complete — it is not a long-term hosting commitment; long-term operation
  means self-hosting from the repository.
- **Horizontal-scaling suitability is not yet explicitly documented.**
  Requested by the customer at the Week 6 documentation review; the
  assessment and an explicit statement in the architecture documentation are
  scheduled as `Must Have` for Sprint 5
  ([#64](https://github.com/Varriwon4ik/avito_bonus_point_service/issues/64)).

## Handover status

<!-- Finalized during Week 7 (Assignment 6 Part 8): after the transition
     confirmation, set the reached level and confirmation status below and
     mirror them in reports/week7/README.md. -->

- **Handover level reached:** trial stage — `v2.1.0` was trialled by the
  customer at the Week 6 session (10 Jul 2026); the customer confirmed the
  product will be **ready for independent use** after the Week 7 work and
  **independently used the trial release** during the session; it is **not**
  deployed or operated on the customer side, by their own plan (their staff
  deploy after final delivery). The final level
  (`Ready for independent use` / `Independently used by customer` /
  `Deployed or operated on customer side`) is confirmed and recorded here at
  the Week 7 transition.
- **Customer-confirmation status:** pending the Week 7 confirmation
  (`Accepted` / `Accepted with follow-up items` / `Not yet accepted`). The
  customer's stated conditions for a complete delivery: the UI display fix
  ([#60](https://github.com/Varriwon4ik/avito_bonus_point_service/issues/60)
  — resolved 10 Jul) and the horizontal-scaling assessment
  ([#64](https://github.com/Varriwon4ik/avito_bonus_point_service/issues/64)).
- **Is the documentation sufficient for this stage?** Reviewed by the
  customer at the Week 6 session: judged **complete overall** ("the READMEs
  cover what I need"), with one requested addition — an explicit
  horizontal-scaling statement in the architecture documentation, tracked as
  [#64](https://github.com/Varriwon4ik/avito_bonus_point_service/issues/64)
  and reflected in the known-limitations list above until resolved.
- **Support that remains necessary from the team (current stage):** operating
  the university trial VM until grading, deploying releases to it, and
  completing the remaining Week 7 actions (the #64 scaling write-up, final
  release `MVP v3`, and the acceptance confirmation of this document). After
  final delivery the customer takes care of the project on their own — no
  ongoing team support was requested.
