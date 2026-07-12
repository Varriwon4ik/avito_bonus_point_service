# AGENTS.md

Guidance for AI coding agents (and their operators) working in this
repository. Human contributors: read [`CONTRIBUTING.md`](CONTRIBUTING.md)
first; this file adds the agent-specific context, commands, and constraints.

## Project shape

- **Single Go module, layered monolith, one binary.** `cmd/api` serves the
  REST API, the admin web UI, Swagger UI (`/docs`), and Prometheus `/metrics`.
  `cmd/autotest` is a console front end for the shared `internal/autotest`
  engine.
- `internal/api` — HTTP handlers, middleware, request/response contract.
  `internal/data` — persistence and **all** business rules inside Postgres
  transactions with `SELECT ... FOR UPDATE` row locking. Every state change
  must go through `internal/data`; never bypass it from a handler.
- `internal/data/migrations/` — forward-only SQL migration pairs, applied
  automatically at startup. Add a new numbered pair; never edit an applied
  migration.
- **The web UI is embedded:** `cmd/api/web/index.html` is compiled in via
  `//go:embed`. After editing it, rebuild and restart the binary — changes do
  not appear on disk-serve. (A missed rebuild is exactly how bug #60
  happened.)
- The UI and autotester consume only the public `/v1` API — do not add
  private endpoints for them.

## Setup and verification commands

```sh
docker compose up -d postgres      # local DB (bonus:bonus, disposable)
export DB_DSN="postgres://bonus:bonus@localhost:5432/bonus_ledger?sslmode=disable"
go run ./cmd/api                   # API + UI on :8080 (-port to change)

# what CI enforces — run all of it before declaring work done:
gofmt -l .                         # must print nothing
go vet ./...
go build ./...
export TEST_DATABASE_URL="postgres://bonus:bonus@localhost:5432/bonus_ledger?sslmode=disable"
go test ./... -race -count=1 -p 1 -coverpkg=./... -coverprofile=coverage.out
bash scripts/coverage_gate.sh coverage.out   # ≥30% for internal/data, internal/api
```

- Tests **truncate tables** in `TEST_DATABASE_URL` — point it only at a
  disposable database, and self-skip when it is unset.
- Keep tests serialized (`-p 1`): they share one database.
- Quick end-to-end check: `curl http://localhost:8080/healthz` →
  `{"status":"ok"}`; `POST /v1/autotest/run` runs the built-in checks against
  a dedicated `autotest-`-prefixed user (`mode`: `single` or `multi_key`).

## Hard constraints

- **Never commit secrets, PII, credentials, recordings, or recording links.**
  The only allowed in-repo credentials are the throwaway local/CI
  `bonus:bonus` Postgres pair. `.env` is git-ignored; the sanitized example
  is `.env.example`.
- **Never push to `main`.** Every change goes through an issue-linked branch
  (`<issue-number>-short-description`) and a PR reviewed by a person other
  than the author. Merge commits only.
- **Never rewrite history, delete PRs/issues/reviews, or move/delete release
  tags and milestones** — they are graded course evidence. To undo a merged
  change, open a reverting PR.
- **Do not edit past week reports** (`reports/week2` … `reports/week5`) —
  they are submitted assignment evidence. Only the current week's report set
  is fair game.
- **Do not weaken the gates:** tests, the coverage gate, QRTs, `govulncheck`,
  and Lychee stay active; a check may only be replaced by a documented
  equivalent or stronger one.
- The API is intentionally unauthenticated (internal-network product) — do
  not add auth ad hoc; that is an explicit product decision (reverted US-07).

## Definition of Done for any change

Issue acceptance criteria verified, CI green, a user-visible change gets an
issue-linked `[Unreleased]` entry in `CHANGELOG.md` (PR template checklist),
and the affected **maintained docs are updated in the same PR** — README/
OpenAPI for API changes, architecture views + ADR for structural decisions,
quality/testing docs for gate changes,
[docs/customer-handover.md](docs/customer-handover.md) for anything touching
access, deployment, or limitations. Full standard:
[docs/definition-of-done.md](docs/definition-of-done.md).

## Where to look things up

| Question | Source |
|---|---|
| How the team works, git/review flow | [docs/development-process.md](docs/development-process.md) |
| Why the system is shaped this way | [docs/architecture/README.md](docs/architecture/README.md) + ADRs |
| API contract and examples | [README.md](README.md), [api/openapi.yaml](api/openapi.yaml) |
| Test/CI status, critical modules | [docs/testing.md](docs/testing.md) |
| Quality requirements and their tests | [docs/quality-requirements.md](docs/quality-requirements.md), [docs/quality-requirement-tests.md](docs/quality-requirement-tests.md) |
| Customer-facing acceptance scenarios | [docs/user-acceptance-tests.md](docs/user-acceptance-tests.md) |
| Handover / operations state | [docs/customer-handover.md](docs/customer-handover.md) |
| Sprint plan and release mapping | [docs/roadmap.md](docs/roadmap.md) |

AI usage on this project is disclosed per assignment in
`reports/weekN/llm-report.md`; agent-assisted changes get no special
treatment — same review, same gates, same Definition of Done.
