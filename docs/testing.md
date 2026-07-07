# Testing Status

Canonical testing-status artifact for the Bonus Points Ledger Service. It records
the critical modules and their coverage, the automated test and CI/QA check
status, the additional QA check rationale, manual evidence that does not count as
a QRT, and which Assignment 4 quality gates remain active for later project work.

- Quality requirements: [docs/quality-requirements.md](quality-requirements.md)
- Quality requirement tests: [docs/quality-requirement-tests.md](quality-requirement-tests.md)
- CI pipeline: [.github/workflows/ci.yml](../.github/workflows/ci.yml)
- Latest CI on the protected default branch:
  [Actions › CI › branch:main](https://github.com/Varriwon4ik/avito_bonus_point_service/actions/workflows/ci.yml?query=branch%3Amain)

## How to run the tests locally

```sh
docker compose up -d postgres
export TEST_DATABASE_URL="postgres://bonus:bonus@localhost:5432/bonus_ledger?sslmode=disable"

# full suite (serialised, race-enabled) with a coverage profile
go test ./... -race -count=1 -p 1 -coverpkg=./... -coverprofile=coverage.out

# only the automated quality requirement tests
go test ./internal/api/ -race -run 'TestQRT' -v

# enforce the per-module coverage gate
bash scripts/coverage_gate.sh coverage.out
```

Integration and QRT tests self-skip when `TEST_DATABASE_URL` is unset.

## Critical Modules and Coverage

Critical modules are the source areas responsible for core user workflows,
persistence, business rules, and the external API contract — where a defect would
materially affect the product. Per-module line coverage is enforced at ≥ 30% by
the CI gate ([QRT-003](quality-requirement-tests.md#qrt-003-critical-module-line-coverage));
the build fails if a critical module drops below the threshold. The exact current
percentages are printed by the *"Per-module coverage gate (QRT-003)"* step and the
`coverage-profile` artifact of each CI run.

| Critical module | Why critical | Required line coverage | Current line coverage | Evidence |
|---|---|---:|---:|---|
| `internal/data` | Persistence, transactional balance mutations (`SELECT ... FOR UPDATE`), FIFO-by-expiry consumption, idempotency cache, hold lifecycle, business rules. | 30% | ≥ 30% (gate-enforced) | [CI coverage gate](https://github.com/Varriwon4ik/avito_bonus_point_service/actions/workflows/ci.yml?query=branch%3Amain) |
| `internal/api` | HTTP handlers, request/response contract, error envelope, logging/metrics middleware, OpenAPI route wiring. | 30% | ≥ 30% (gate-enforced) | [CI coverage gate](https://github.com/Varriwon4ik/avito_bonus_point_service/actions/workflows/ci.yml?query=branch%3Amain) |

> Global repository coverage is lower than critical-module coverage because the
> CLI entry points `cmd/api` and `cmd/autotest` (process wiring, flag parsing,
> interactive console I/O) are intentionally lightly tested — they are thin shells
> over the covered `internal/...` packages. They are **not** gated.

## Automated Test Status

| Test type | Scope | Command or CI check | Latest result | Evidence |
|---|---|---|---|---|
| Unit tests | Hold lifecycle / data logic, label validation (`internal/data`) | `go test ./internal/data/...` | Passing | [holds_test.go](../internal/data/holds_test.go), [labels_test.go](../internal/data/labels_test.go) |
| Integration tests | API + Postgres: accrual, balance, hold/confirm/cancel, debit, idempotency, FIFO-by-expiry, pagination, error envelope, OpenAPI routes, metrics | `go test ./internal/api/... -race -p 1` | Passing | [integration_test.go](../internal/api/integration_test.go), [pagination_test.go](../internal/api/pagination_test.go), [metrics_integration_test.go](../internal/api/metrics_integration_test.go) |
| Concurrency tests | Race-safety / no double-spend under concurrent debits & holds | `go test ./internal/api/... -race` | Passing | [concurrent_idempotency_test.go](../internal/api/concurrent_idempotency_test.go) |
| Automated QRTs | QR-001, QR-002, QR-003 | `go test ./internal/api/ -run TestQRT` + `bash scripts/coverage_gate.sh` | Passing | [qrt_test.go](../internal/api/qrt_test.go), [quality-requirement-tests.md](quality-requirement-tests.md) |

## CI and QA Check Status

All checks below run on pull requests and on pushes to the protected default
branch. They are **required for Done** (see
[docs/definition-of-done.md](definition-of-done.md)).

| Gate or check | Required for Done? | Latest protected-branch status | Evidence |
|---|---|---|---|
| Linting / static check (`go vet`) | Yes | Passing | [CI](https://github.com/Varriwon4ik/avito_bonus_point_service/actions/workflows/ci.yml?query=branch%3Amain) |
| Formatting (`gofmt -l`) | Yes | Passing | [CI](https://github.com/Varriwon4ik/avito_bonus_point_service/actions/workflows/ci.yml?query=branch%3Amain) |
| Build (`go build ./...`) | Yes | Passing | [CI](https://github.com/Varriwon4ik/avito_bonus_point_service/actions/workflows/ci.yml?query=branch%3Amain) |
| Unit + integration tests (`-race`) | Yes | Passing | [CI](https://github.com/Varriwon4ik/avito_bonus_point_service/actions/workflows/ci.yml?query=branch%3Amain) |
| Automated QRTs (QRT-001/002/003) | Yes | Passing | [CI](https://github.com/Varriwon4ik/avito_bonus_point_service/actions/workflows/ci.yml?query=branch%3Amain) |
| Per-module line coverage gate (≥30%) | Yes | Passing | [CI](https://github.com/Varriwon4ik/avito_bonus_point_service/actions/workflows/ci.yml?query=branch%3Amain) |
| Additional QA check (`govulncheck`) | Yes | Passing | [CI](https://github.com/Varriwon4ik/avito_bonus_point_service/actions/workflows/ci.yml?query=branch%3Amain) |
| Link check (Lychee) | Yes | Passing | [Lychee workflow](../.github/workflows/lychee.yml) |

## Additional QA Check Rationale

The additional QA check is **`govulncheck`**, the official Go vulnerability
scanner. It is distinct from the required linting, formatting/vet, build, unit
and integration tests, coverage gate, the QRTs, and the Lychee link check, and it
runs as its own CI job.

| QA objective or risk | Additional QA check | Scope | Latest result | Evidence | Limitations or follow-up |
|---|---|---|---|---|---|
| Known vulnerabilities in third-party dependencies or the Go standard library could expose the deployed ledger to avoidable risk. | Automated `govulncheck` static analysis against the Go vulnerability database. | The module's dependencies (`go.mod`/`go.sum`) and reachable standard-library symbols across `./...`. | Passing | [govulncheck job](https://github.com/Varriwon4ik/avito_bonus_point_service/actions/workflows/ci.yml?query=branch%3Amain) | `govulncheck` reports only *known* advisories and only those reachable from the call graph; new advisories require a re-run, and stdlib findings are resolved by bumping the pinned Go toolchain. |

## Manual Evidence That Does Not Count as QRT

| Evidence | Scope | Result | Follow-up PBI or issue |
|---|---|---|---|
| Customer UAT (code-walkthrough session, 28 Jun 2026) | Sprint 2 increment: CI (US-14), autotester (US-15), pagination (US-09) | Passed — customer satisfied with the team's design decisions | See [docs/user-acceptance-tests.md](user-acceptance-tests.md) |
| Sprint Review file-by-file inspection | All files committed during Sprint 2 | Approved by the customer | [Sprint 2 milestone](https://github.com/Varriwon4ik/avito_bonus_point_service/milestone/2) |
| Customer UAT (customer-directed screen-share session, 3 Jul 2026) | Sprint 3 / MVP v2 increment: UAT-004 (US-16), UAT-005 (US-17), UAT-006 (US-18 + US-08) | All passed — features accepted; new backlog item raised | [US-19 / #50](https://github.com/Varriwon4ik/avito_bonus_point_service/issues/50); see [docs/user-acceptance-tests.md](user-acceptance-tests.md) |

## Notes on excluded links

Every exclusion in [`.lycheeignore`](../.lycheeignore) is narrow, justified,
and manually verified before submission:

- The deployed increment runs on a University VM at a private (RFC 1918)
  address that is unreachable from GitHub-hosted CI.
- The Google Drive demo-video links are session-gated by Drive and rate-limited
  for CI runners; they are verified manually in a browser.
- The GitHub Pages documentation site URL is excluded until the first Pages
  deployment from `main` completes (the site deploys on merge, so PR CI would
  otherwise fail on a not-yet-published URL); it is verified manually after
  deployment.

## Gates that remain active for later project work

The Assignment 4 tests, coverage gate, QRTs, and `govulncheck` are maintained
repository gates. Later PRs and protected-default-branch changes must keep them
passing. If a product change makes a check obsolete, it must be replaced with a
documented equivalent or stronger check — gates are not removed or narrowed just
because the assignment was submitted. The
[Definition of Done](definition-of-done.md) enforces this for every future PBI.
