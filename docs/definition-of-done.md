# Definition of Done

This document defines Team 01's shared, minimum completion standard for the
Bonus Points Ledger Service. It complements `Process_Requirements.md`, which is
the authoritative source for course-level Definition of Done expectations.

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

## 3. Testing

- [ ] New or changed behaviour is covered by automated tests.
- [ ] The full test suite passes (`go test ./...`) against a real Postgres.
- [ ] Concurrency- or data-integrity-sensitive changes include race-safe tests.
- [ ] No known regression is introduced in previously passing tests.

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
