# Quality Requirement Tests

This document defines the automated quality requirement tests (QRTs) that verify
the measurable scenarios in [docs/quality-requirements.md](quality-requirements.md).
Each QRT is an automated test or CI check that directly verifies one or more QR
scenarios; all three run in [CI](../.github/workflows/ci.yml) on every pull
request and on every push to the protected default branch.

QRTs are **maintained product assets**. Later project work must keep them
passing or replace them with documented equivalent or stronger checks when the
product changes.

| QRT | Verifies | Kind | Automated command / CI check |
|---|---|---|---|
| QRT-001 | QR-001 Time behaviour | Integration test (latency sample) | `go test ./internal/api/ -run TestQRT001BalanceResponseTime` |
| QRT-002 | QR-002 Integrity | Concurrency test (race-enabled) | `go test ./internal/api/ -run TestQRT002DebitIntegrityUnderConcurrency -race` |
| QRT-003 | QR-003 Testability | CI coverage gate | `bash scripts/coverage_gate.sh coverage.out` |

## QRT-001: Balance read response time

**Linked quality requirement:** [QR-001](quality-requirements.md#qr-001-balance-read-response-time)

**Verification method:** Automated Go integration test that measures the latency
distribution of the balance endpoint and asserts the p95 budget.

**Test data, setup, or environment:** An in-process `httptest` server backed by
the real Postgres service container (`TEST_DATABASE_URL`). A test user is seeded
with 10 accrual lots; the test issues 20 warm-up reads followed by 200 measured
`GET /v1/users/{id}/balance` requests and computes p50/p95/max. Self-skips when
`TEST_DATABASE_URL` is unset.

**Automated command or CI check:**
`go test ./internal/api/ -run TestQRT001BalanceResponseTime` — runs in the CI
step *"Quality requirement tests (QRT-001, QRT-002)"*. Source:
[`internal/api/qrt_test.go`](../internal/api/qrt_test.go).

**Expected measurable result:** p95 latency ≤ 200 ms over the 200-request sample
(budget overridable via `QRT_BALANCE_P95_BUDGET_MS`); the test fails if the p95
exceeds the budget.

**Evidence location:** The *"Quality requirement tests (QRT-001, QRT-002)"* step
of the latest [CI run on `main`](https://github.com/Varriwon4ik/avito_bonus_point_service/actions/workflows/ci.yml?query=branch%3Amain)
(the log line prints `p50/p95/max`).

## QRT-002: Debit integrity under concurrency

**Linked quality requirement:** [QR-002](quality-requirements.md#qr-002-ledger-integrity-under-concurrency)

**Verification method:** Automated Go concurrency test, executed under the Go
race detector, asserting the no-overspend / no-negative-balance invariants.

**Test data, setup, or environment:** `httptest` server backed by the Postgres
service container. A user is accrued 1000 points; 40 goroutines then fire
one-shot debits of 100 each (distinct idempotency keys) released simultaneously
via a `sync.WaitGroup` "start pistol" to maximise contention. Self-skips when
`TEST_DATABASE_URL` is unset.

**Automated command or CI check:**
`go test ./internal/api/ -run TestQRT002DebitIntegrityUnderConcurrency -race` —
runs in the CI step *"Quality requirement tests (QRT-001, QRT-002)"* and again
under `-race` in *"Test with coverage"*. Source:
[`internal/api/qrt_test.go`](../internal/api/qrt_test.go).

**Expected measurable result:** Exactly 10 debits succeed (`200`), the remaining
30 return `409 Conflict`, total spent ≤ 1000, final `available` = `1000 − spent`
and ≥ 0, and the race detector reports no data races. The test fails on any
overspend, negative balance, or lost update.

**Evidence location:** The QRT and coverage steps of the latest
[CI run on `main`](https://github.com/Varriwon4ik/avito_bonus_point_service/actions/workflows/ci.yml?query=branch%3Amain).

## QRT-003: Critical module line coverage

**Linked quality requirement:** [QR-003](quality-requirements.md#qr-003-critical-module-testability)

**Verification method:** Automated CI gate that parses the suite's coverage
profile and enforces a minimum per-module line-coverage threshold.

**Test data, setup, or environment:** The full test suite runs in CI with
`-coverpkg=./...` producing `coverage.out`, so coverage of a module includes
exercise through other packages (e.g. `internal/data` exercised via the
`internal/api` integration tests). [`scripts/coverage_gate.sh`](../scripts/coverage_gate.sh)
aggregates covered/total statements per package and per critical module.

**Automated command or CI check:** `bash scripts/coverage_gate.sh coverage.out`
— runs in the CI step *"Per-module coverage gate (QRT-003)"*.

**Expected measurable result:** `internal/data` ≥ 30% and `internal/api` ≥ 30%
line coverage (threshold overridable via `COVERAGE_THRESHOLD`); the step exits
non-zero and fails the build if any critical module is below the threshold.

**Evidence location:** The *"Per-module coverage gate (QRT-003)"* and
*"Coverage summary (total)"* steps of the latest
[CI run on `main`](https://github.com/Varriwon4ik/avito_bonus_point_service/actions/workflows/ci.yml?query=branch%3Amain),
plus the uploaded `coverage-profile` artifact.

## Notes on QRT classification

- QRT-001 and QRT-002 are integration/concurrency tests that **directly verify
  measurable QR scenarios**, so they qualify as QRTs (not merely generic unit
  tests). QRT-003 is a coverage CI check that directly verifies the QR-003
  Testability scenario.
- Per Assignment 4 rules, no single CI check is counted both as a QRT and as the
  *additional QA check*. The additional QA check is the separate
  [`govulncheck`](../.github/workflows/ci.yml) dependency/stdlib vulnerability
  scan, which is **not** one of these QRTs — see
  [docs/testing.md](testing.md#additional-qa-check-rationale).
