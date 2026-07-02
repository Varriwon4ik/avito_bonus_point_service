# Definition of Done

This document defines Team 01's shared, minimum completion standard for the
Bonus Points Ledger Service. It complements `Process_Requirements.md`, which is
the authoritative source for course-level Definition of Done expectations, and
the maintained quality assets it depends on:
[docs/quality-requirements.md](quality-requirements.md),
[docs/quality-requirement-tests.md](quality-requirement-tests.md), and
[docs/testing.md](testing.md).

A Product Backlog Item (PBI) may be marked **`Done`** only when **every**
applicable item below is satisfied. A user story is `Done` only when all linked
supporting PBIs required to satisfy its acceptance criteria are also `Done`.

## 1. Scope and acceptance

- [ ] All acceptance criteria on the PBI are met and individually verified.
- [ ] The implemented behaviour matches the user-story statement and agreed scope.
- [ ] Any scope change versus the original PBI is recorded in the issue history.

## 2. Implementation quality

- [ ] Code is merged into `main` through a reviewed, issue-linked PR.
- [ ] The branch follows the `<issue-number>-short-description` naming convention.
- [ ] The PR uses a merge-commit and closes/links its issue.
- [ ] Code builds successfully (`go build ./...`) with no new warnings.
- [ ] Code is formatted (`gofmt`) and passes static checks (`go vet ./...`).

## 3. Testing and quality gates

- [ ] New or changed behaviour is covered by automated tests (unit and, where
      product components interact, integration tests).
- [ ] The full test suite passes (`go test ./... -race`) against a real Postgres.
- [ ] Concurrency- or data-integrity-sensitive changes include race-safe tests.
- [ ] No known regression is introduced in previously passing tests.
- [ ] Relevant automated **quality requirement tests** pass, or are explicitly
      documented as not applicable
      ([docs/quality-requirement-tests.md](quality-requirement-tests.md)).
- [ ] Critical modules (`internal/data`, `internal/api`) keep **≥ 30% line
      coverage**; the CI per-module coverage gate passes.
- [ ] All required CI checks pass on the PR and before merge: `gofmt`,
      `go vet`, build, unit + integration tests, QRTs, the coverage gate,
      `govulncheck` (additional QA check), and the Lychee link check
      ([.github/workflows/ci.yml](../.github/workflows/ci.yml)).
- [ ] Testing evidence (CI run + coverage) is preserved in the PR, CI, or linked
      documentation.

## 4. Review and verification

- [ ] At least one other team member (the recorded reviewer) approved the PR.
- [ ] At least one meaningful review comment was left and resolved.
- [ ] Acceptance criteria were explicitly verified in the PR before merge.

## 5. Documentation and traceability

- [ ] `CHANGELOG.md` is updated for every user-visible change (Keep a Changelog
      format, SemVer), consistent with the repository requirements.
- [ ] `README.md` / `api/openapi.yaml` are updated when behaviour or the API
      contract changes.
- [ ] `docs/user-stories.md` reflects the PBI's current Work Status and Sprint.
- [ ] The issue links the PR(s), and the PR references the issue number.

## 6. Increment integrity

- [ ] The delivered increment runs via `docker compose up --build` with no
      manual patching.
- [ ] No secrets, credentials, or sensitive payloads are logged or committed.
- [ ] The change does not break the documented MVP v1 flows.

> A PBI that fails any applicable checkbox above remains `In Progress`,
> `In Review`, or `Blocked` — it is **not** `Done`.

## 7. Continuing governance

This Definition of Done continues to govern **all later project work**. The
Assignment 4 tests, automated QRTs, per-module coverage gate, and additional QA
check are maintained gates: later PBIs must keep them passing, and a gate may be
removed or narrowed only when replaced by a documented equivalent or stronger
check. When the product stack, critical modules, quality requirements, or CI
configuration change, update this document and the linked quality assets so they
keep describing the current completion standard.
