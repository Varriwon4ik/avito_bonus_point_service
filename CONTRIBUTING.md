# Contributing

Thanks for contributing to the **Bonus Points Ledger Service**. This guide is
the entry point for making a change land in this repository — as a team
member, a customer-side developer after handover, or an external contributor.
The full description of how the team works is
[docs/development-process.md](docs/development-process.md); this file is the
practical summary you need before opening a pull request.

## Set up a development environment

Prerequisites: **Docker** (with compose) and **Go** (any version ≥ the version
pinned in [`go.mod`](go.mod) — it bootstraps the pinned toolchain
automatically).

```sh
git clone https://github.com/Varriwon4ik/avito_bonus_point_service.git
cd avito_bonus_point_service

# full product (API + Postgres), web UI on http://localhost:8080
docker compose up --build

# or: database only, API from source
docker compose up -d postgres
export DB_DSN="postgres://bonus:bonus@localhost:5432/bonus_ledger?sslmode=disable"
go run ./cmd/api
```

Sanitized example environment values: [`.env.example`](.env.example). Never
commit a real `.env`, credentials, or any other secret — see the
[configuration and secrets rules](docs/development-process.md#configuration-and-secrets-management).

## Verify your change

Run what CI runs before you push:

```sh
gofmt -l .          # must print nothing
go vet ./...
go build ./...

# tests need a running Postgres and TEST_DATABASE_URL
docker compose up -d postgres
export TEST_DATABASE_URL="postgres://bonus:bonus@localhost:5432/bonus_ledger?sslmode=disable"
go test ./... -race -count=1 -p 1 -coverpkg=./... -coverprofile=coverage.out

bash scripts/coverage_gate.sh coverage.out   # per-module ≥30% coverage gate
```

Note: the web UI (`cmd/api/web/index.html`) is **embedded into the binary** —
rebuild/restart the API after editing it; it is not served from disk. Tests
truncate tables in the target database, so point `TEST_DATABASE_URL` at a
disposable database only.

## Workflow: issue → branch → PR → review → merge

Every change — code, docs, CI config — lands through a reviewed pull request.
Direct pushes to `main` are disabled.

1. **Start from an issue.** Create it with the matching issue form (User
   Story, Other PBI, Bug Report, or Course Task; blank issues are disabled).
   Record the expected outcome, acceptance criteria, Story Points, the
   implementer, and a **different** reviewer.
2. **Branch from the issue**, named `<issue-number>-short-description`
   (e.g. `42-add-login-form`).
3. **Open a PR** using the template: summary, `Closes #NN`, testing performed,
   acceptance-criteria verification, and the changelog checklist — exactly one
   of: added/updated a user-visible entry in [`CHANGELOG.md`](CHANGELOG.md)
   (under `[Unreleased]`, [Keep a Changelog](https://keepachangelog.com/)
   categories, issue-linked), or "not user-visible".
4. **CI must be green.** The required checks: `gofmt`, `go vet`, build, unit +
   integration tests with the race detector against real Postgres, the
   automated quality requirement tests, the per-module ≥30% coverage gate,
   `govulncheck`, and Lychee link checking for Markdown. A red build blocks
   merge.
5. **Review** by the recorded reviewer — never the author; self-approval is
   disabled. The reviewer verifies the acceptance criteria against the change.
6. **Merge with a merge commit** (squash/rebase merging is disabled). The
   `Closes #NN` link closes the issue; it may be marked `Done` only when the
   [Definition of Done](docs/definition-of-done.md) holds.

If a merged change turns out broken, **revert it with a new reviewed PR** —
history is never rewritten.

## Keep the maintained docs current

If your change affects behaviour, structure, or process, update the affected
maintained documentation **in the same PR** (this is part of the Definition of
Done):

- API/user-visible behaviour → [`README.md`](README.md),
  [`api/openapi.yaml`](api/openapi.yaml), `CHANGELOG.md`, and the
  [UAT scenarios](docs/user-acceptance-tests.md) where user goals change
- Architecture, deployment, or an important decision →
  [docs/architecture/](docs/architecture/README.md) views and a new or
  superseding ADR
- Quality requirements or their tests →
  [docs/quality-requirements.md](docs/quality-requirements.md) /
  [docs/quality-requirement-tests.md](docs/quality-requirement-tests.md) /
  [docs/testing.md](docs/testing.md)
- Workflow or tooling → [docs/development-process.md](docs/development-process.md)
  and this file
- Access, deployment steps, or limitations →
  [docs/customer-handover.md](docs/customer-handover.md)

The `docs/` set is published automatically to the
[hosted documentation site](https://varriwon4ik.github.io/avito_bonus_point_service/)
on every merge to `main`.

## Conventions

- Go code is formatted with `gofmt`; keep handlers thin and put business rules
  and persistence behind `internal/data`'s locked transactions (see
  [ADR-001](docs/architecture/adr/ADR-001-postgres-row-locking-for-ledger-integrity.md)
  and [ADR-003](docs/architecture/adr/ADR-003-layered-monolith-with-gated-critical-modules.md)).
- Error responses use the standard JSON envelope
  (`{"error": ..., "message": ...}`); mutating endpoints take a client
  `idempotency_key`
  ([ADR-004](docs/architecture/adr/ADR-004-client-supplied-idempotency-keys.md)).
- Schema changes are new forward-only migration pairs in
  [`internal/data/migrations/`](internal/data/migrations/).
- AI coding agents working in this repository follow [`AGENTS.md`](AGENTS.md);
  AI-assisted contributions go through exactly the same review and CI gates as
  any other change, and AI usage is disclosed in the course LLM reports.

## Questions

Open an issue with the appropriate form, or start from the
[documentation site](https://varriwon4ik.github.io/avito_bonus_point_service/)
— the [development process](docs/development-process.md) and
[architecture](docs/architecture/README.md) pages answer most "why is it like
this?" questions.
