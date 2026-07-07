# Changelog

All notable changes to this project are documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Fixed

- **Web autotester always reported "Autotester found issues" (US-17).** The
  Autotester tab read the pass/fail verdict from the HTTP wrapper object instead
  of the response body, so `passed` was always `undefined` and every run — even
  a fully successful one with code 200 — showed "✗ Some checks failed" with an
  empty per-check table. The tab now unwraps the response body before rendering
  the verdict and surfaces API validation errors (HTTP 4xx) as their own toast
  message. (#54)

## [2.0.0] - 2026-07-05

MVP v2 — the Sprint 3 increment (Assignment 5, 29 Jun–5 Jul 2026). Maps to the
[Sprint 3 milestone](https://github.com/Varriwon4ik/avito_bonus_point_service/milestone/3):
US-08, US-16, US-17, US-18, plus the maintained architecture and
development-process documentation.

### Added

- **Web autotester tab and run endpoint (US-17).** A new "Autotester" tab in the
  web UI lets an administrator fill in a scenario (label, test user, amount, TTL,
  parallel requests) and run it against the live instance, then see a per-check
  pass/fail report. It is backed by a new `POST /v1/autotest/run` endpoint and a
  shared `internal/autotest` engine extracted from the `cmd/autotest` console
  tool, so the web and console front ends run identical accrual-correctness and
  parallel-request checks. All requests target a dedicated `autotest-`-prefixed
  user so real accounts are never touched. (#40, PR #49)
- **Exact HTTP response codes in the web UI (US-16).** Accrual and debit
  operations in the web UI now display the exact HTTP status code and outcome of
  each request (e.g. green `200` on success, red `404` for an unknown user)
  alongside the other request details, so administrators can report issues
  precisely without backend tooling. (#39, PR #46)
- **Labels on transactions (US-18).** Accruals can carry a label — a preset
  value such as `test` or `real`, or a custom short label — which is validated
  (trimmed, ≤ 32 characters, no control characters), stored on the ledger entry
  (migration `0004`), and displayed in the web UI transactions view, so test and
  real traffic can be told apart. First landed in PR #44, reverted in PR #47
  after issues surfaced, and re-landed fixed in PR #49. (#41)
- **Configurable per-accrual TTL validation and bounds (US-08).** New
  `MIN_TTL_DAYS` / `MAX_TTL_DAYS` environment variables bound the accepted
  `ttl_days` on accruals; requests outside the configured range return
  `400 Bad Request` with a clear message, protecting expiry logic from absurd
  lifetimes. `DEFAULT_TTL_DAYS` behaviour is unchanged when `ttl_days` is
  omitted. (#5, PR #42)
- **Maintained architecture and process documentation (Assignment 5).**
  Architecture documentation with static, dynamic, and deployment views and five
  ADRs (`docs/architecture/`), the documented development process and
  configuration management (`docs/development-process.md`), and a hosted MkDocs
  documentation site published to GitHub Pages on every merge to `main`. (#51)

## [1.1.0] - 2026-06-28

Sprint 2 increment (Assignment 4, 22–28 Jun 2026). Maps to the
[Sprint 2 milestone](https://github.com/Varriwon4ik/avito_bonus_point_service/milestone/2):
US-09, US-14, US-15, plus the automated quality gates. US-07 was reverted.

### Added

- **Pagination for transaction history (US-09).** `GET /v1/users/{id}/transactions`
  now returns a structured page envelope (`user_id`, `page`, `offset`, `total`,
  `entries`) and accepts `page` (1-based page number) and `offset` (page size,
  `1`–`500`, default `20`) query parameters, so callers can page through a user's
  newest-first ledger/audit history instead of receiving an unbounded list.
  Invalid `page` or `offset` values return `400 Bad Request`. (#6, PR #30)
- **Automated autotester tool (US-15).** A new `cmd/autotest` console tool
  defines, stores (in a dedicated `autotest_scenarios` table, migration `0003`),
  and replays reusable accrual test scenarios — including concurrent/parallel
  requests — against a running instance and reports the observed balance and
  ledger outcomes, so newly committed code can be regression-checked end-to-end
  without hand-writing a bespoke test for each case. (#29, PR #31)
- **Continuous integration pipeline (US-14).** A GitHub Actions workflow
  (`.github/workflows/ci.yml`) runs on every push and every pull request to
  `main`. It provisions a Postgres service container, pins the Go toolchain to
  the version in `go.mod`, and runs `go mod verify`, `gofmt`, `go vet`,
  `go build` and the full test suite with the race detector. Failures produce a
  status check that can be made required in branch protection, so regressions
  are caught before merge and `main` stays releasable. (#28, PR #33)
- **Automated quality gates and quality documentation (QR-001/002/003).** CI now
  runs automated quality requirement tests (balance-read p95 latency; ledger
  integrity under concurrent debits, race-enabled), a per-module line-coverage
  gate (≥30% for `internal/data` and `internal/api`), and a `govulncheck`
  dependency/standard-library vulnerability scan (the additional QA check). Added
  maintained docs: `docs/quality-requirements.md`,
  `docs/quality-requirement-tests.md`, `docs/testing.md`, and
  `docs/user-acceptance-tests.md`. (US-14)

### Removed

- **Admin authentication for manual accrual (US-07).** An admin bearer-token
  guard around `POST /v1/users/{id}/accruals` was merged earlier in this cycle
  (PR #32) but then reverted (PR #34) due to bugs and integration issues found
  in review, so it never shipped in a tagged release; the endpoint behaves as it
  did in `v1.0.0`. US-07 is now marked `Removed`, and a different feature was
  prioritized in its place for the Sprint. (#4)

## [1.0.0] - 2026-06-21

MVP v1 — the first delivered increment (Sprint 1, 15–21 Jun 2026). Maps to the
`MVP v1` scope: US-05, US-10, US-11, US-12, US-13.

### Added

- **Hold timeout / automatic release sweep (US-05).** A background job
  (`HOLD_TIMEOUT_HOURS`, default 24h) periodically releases active holds left
  unconfirmed or uncancelled past the timeout, returning the reserved points to
  the user and writing an audit ledger entry. (#3, PR #16/#17)
- **Structured request logging and Prometheus `/metrics` endpoint (US-10).** An
  observability middleware logs every request in structured form (method, route
  pattern, status, latency, bytes, `user_id`) without reading request bodies,
  and exposes request counters/latency histograms plus ledger-level gauges. (#7, PR #18)
- **Concurrent idempotent-key deduplication tests (US-11).** Automated
  concurrency tests verify race-safety and idempotent deduplication of duplicate
  requests sharing an idempotency key. (#8, PR #14)
- **Points removal / expiry system (US-12).** Points expire automatically based
  on their configured lifetime and stop counting toward the balance. (#11, PR #13)

### Changed

- **HTTP response codes and OpenAPI docs (US-13).** Endpoints now return correct
  status codes (200/201/400/404/409/500) with a consistent JSON error envelope,
  and the OpenAPI specification was updated to match. (#12, PR #15)

[Unreleased]: https://github.com/Varriwon4ik/avito_bonus_point_service/compare/v2.0.0...HEAD
[2.0.0]: https://github.com/Varriwon4ik/avito_bonus_point_service/releases/tag/v2.0.0
[1.1.0]: https://github.com/Varriwon4ik/avito_bonus_point_service/releases/tag/v1.1.0
[1.0.0]: https://github.com/Varriwon4ik/avito_bonus_point_service/releases/tag/v1.0.0
