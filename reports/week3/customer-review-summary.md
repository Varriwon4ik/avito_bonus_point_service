# Customer Review Summary — Sprint 1 (MVP v1)

- **Date:** 21 Jun 2026
- **Sprint:** Sprint 1 (15–21 Jun 2026)
- **Participants (roles):** Product Owner / Scrum Master & backend engineer
  (Mikhail Ilin), backend engineer (Nurislam Denisov), QA engineer
  (Sergey Chuenko), QA engineer (Sanzhar Kadambaev), Customer (Alexander).
  Backend engineer N. Nuriev was unable to connect.

## Artifacts demonstrated

- Live service walkthrough of the MVP v1 increment.
- US-11 — concurrent idempotent-key deduplication / race-safe balance handling.
- US-12 — points removal / expiry system.
- US-13 — HTTP response codes and OpenAPI documentation.
- US-05 — auto-release of stale holds was **not** demonstrated live (owner could
  not connect); deferred to the next review with customer agreement.

## Scope reviewed vs. implemented

| Planned MVP v1 PBI | Demonstrated | Outcome |
|---|---|---|
| US-11 | Yes | Accepted, with a follow-up request |
| US-12 | Yes | Accepted, with a follow-up request |
| US-13 | Yes | Accepted |
| US-05 | No (technical absence of owner) | Deferred demo to next review |
| US-10 (Should Have, in Sprint 1) | Not separately demoed | Delivered in repo |

## Customer feedback

1. **US-11 follow-up:** apply the concurrency tests to older code paths to show
   that the new solution measurably improves the project.
2. **US-12 follow-up:** model point expiry as an explicit ledger transaction
   (a removal transaction recorded with expiration-date handling) rather than
   silently marking/removing points, so expirations are auditable.
3. **US-13:** positive feedback on the code decisions; no changes requested.

## Approvals and decisions

- The customer reacted positively to the demonstrated increment and approved
  showing US-05 at the next review.
- No demonstrated item was rejected. The two follow-up requests are improvements,
  not blockers, and are captured as backlog updates for Sprint 2.

## Risks and action points

- **Risk:** US-05 was not demonstrated live; verification evidence is provided in
  the repository/PR instead. Schedule a live demo next review.
- **Action:** create/refine backlog follow-ups for the US-11 and US-12 requests
  (see [docs/roadmap.md](../../docs/roadmap.md), Sprint 2 candidates).

## Resulting Product Backlog / scope changes

- Two new follow-up items queued for Sprint 2 (concurrency-regression coverage on
  legacy paths; expiry-as-transaction model).
- No change to the accepted MVP v1 scope.
