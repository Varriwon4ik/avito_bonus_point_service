<!-- Extended PR template — Assignment 3 workflow. -->

## Summary

- **What changed?**
- **Why was it needed?**

## Linked issues

- Closes #
- Related user-story IDs: US-

## Type of change

- [ ] User story (user-facing PBI)
- [ ] Other PBI (technical / infra / docs / testing / deployment)
- [ ] Bug fix
- [ ] Course task / administration

## Acceptance-criteria verification

<!-- Copy the acceptance criteria from the linked issue and tick each one,
     describing how it was verified. Required before merge. -->
- [ ] AC1 — verified by:
- [ ] AC2 — verified by:
- [ ] AC3 — verified by:

## Verification

- [ ] `go build ./...` succeeds
- [ ] `go vet ./...` / `gofmt` clean
- [ ] `go test ./...` passes against Postgres
- [ ] Local smoke check via `docker compose up --build`
- [ ] No secrets or sensitive payloads logged/committed

## Definition of Done

- [ ] Meets [docs/definition-of-done.md](../docs/definition-of-done.md)
- [ ] `CHANGELOG.md` updated for any user-visible change
- [ ] Docs / OpenAPI / README updated if behaviour changed
- [ ] `docs/user-stories.md` Work Status updated

## Reviewer

- **Assigned reviewer:** @
- **Branch name follows** `<issue-number>-short-description`: [ ]

## Screenshots / evidence (optional)
