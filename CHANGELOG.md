# Changelog

All notable changes to this project are documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Added

- **Admin authentication for manual accrual (US-07).** The manual accrual
  endpoint (`POST /v1/users/{id}/accruals`) is now guarded by a static admin
  bearer token configured via `ADMIN_API_TOKEN` / `-admin-token`. Requests must
  present `Authorization: Bearer <token>`; missing or invalid credentials return
  `401 Unauthorized`. When no token is configured the check is disabled (open),
  and the server logs a startup warning. (#4)

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

[Unreleased]: https://github.com/Varriwon4ik/avito_bonus_point_service/compare/v1.0.0...HEAD
[1.0.0]: https://github.com/Varriwon4ik/avito_bonus_point_service/releases/tag/v1.0.0
