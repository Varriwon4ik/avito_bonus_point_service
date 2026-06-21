# Week 3 Report — Assignment 3 (MVP v1)

> Canonical public Week 3 report and submission index for Assignment 3,
> Team 01 — Bonus Points Ledger Service.

## 1. Project

**Bonus Points Ledger Service** — a REST-like service for managing an online
store's bonus-points program: configurable point expiry, transactional balance
mutations with row-level locking, two-phase debits (`hold` / `confirm` /
`cancel`), idempotent operations, FIFO-by-expiry consumption, and observability
(structured logging + Prometheus `/metrics`).

- License: [LICENSE](../../LICENSE)
- Run / access instructions: [root README → Running](../../README.md#running)

## 2. Scope since Assignment 2

The Assignment 2 stories were migrated into GitHub issues and refined. The
authoritative current registry is [docs/user-stories.md](../../docs/user-stories.md);
the historical Assignment 2 source is preserved at
[reports/week2/user-stories.md](../week2/user-stories.md).

- **Delivered (Sprint 1 / MVP v1):** US-05, US-11, US-12, US-13 (`Must Have`) and
  US-10 (`Should Have`).
- **Refined out:** US-03, US-04, US-06 were **removed** (already covered by MVP v0
  base functionality); US-07 and US-09 were reprioritized to **`Won't Have`**.
- **Remaining active backlog:** US-01, US-02, US-08 (`Should Have`, not yet
  scheduled).

## 3. Assignment 2 customer feedback addressed in MVP v1

| Assignment 2 feedback | Addressed by |
|---|---|
| Holds could lock points forever if a caller crashed (MVP v0 gap) | US-05 — auto-release sweep ([#3](https://github.com/Varriwon4ik/avito_bonus_point_service/issues/3)) |
| Need proven safety under concurrent/duplicate requests | US-11 — concurrent idempotency tests ([#8](https://github.com/Varriwon4ik/avito_bonus_point_service/issues/8)) |
| Points must expire reliably | US-12 — points removal/expiry ([#11](https://github.com/Varriwon4ik/avito_bonus_point_service/issues/11)) |
| Correct, predictable API responses | US-13 — HTTP codes + OpenAPI ([#12](https://github.com/Varriwon4ik/avito_bonus_point_service/issues/12)) |
| Operability / monitoring | US-10 — logging + `/metrics` ([#7](https://github.com/Varriwon4ik/avito_bonus_point_service/issues/7)) |

## 4–11. Backlog, Sprint, MVP, and size

- Historical user stories: [reports/week2/user-stories.md](../week2/user-stories.md)
- Current user-story index: [docs/user-stories.md](../../docs/user-stories.md)
- **Product Backlog board/view:** <!-- TODO: insert GitHub Project URL after creation -->
- **Current Sprint Backlog board/view:** <!-- TODO: insert Sprint-filtered Project view URL -->
- **Sprint 1 milestone (authoritative Sprint Goal / dates / scope):**
  [Sprint 1](https://github.com/Varriwon4ik/avito_bonus_point_service/milestone/1)
- **Total Product Backlog size:** 71 SP (qualifying PBIs; excludes `Won't Have`
  US-07/US-09 and removed US-03/US-04/US-06)
- **Total Sprint 1 size:** 55 SP (5 user stories = 29 SP + 8 supporting PBIs = 26 SP)
- **MVP version view (MVP v1 scope):** <!-- TODO: insert grouped/filtered Project view URL -->

## 12. Selected MVP v1 scope

MVP v1 is the set of PBIs marked `MVP v1`. It is intentionally small and limited
to delivered `Must Have` user stories plus their supporting technical PBIs:

- **US-05** — auto-release stale holds.
- **US-11** — concurrent idempotent-key deduplication tests.
- **US-12** — points removal / expiry system.
- **US-13** — HTTP response codes and OpenAPI docs.
- Supporting technical PBIs (migrations, metrics wiring, CI, deployment, test
  harness) required to deliver and verify the above.

## 13. PBI types, statuses, priorities, Sprint & MVP tracking

We follow the shared definitions in
[Process_Requirements.md](../../Process_Requirements.md):

- **Types:** User Story, Other PBI (technical/infra/docs/testing/deployment),
  Course Task, Bug Report — see the issue templates in
  [.github/ISSUE_TEMPLATE/](../../.github/ISSUE_TEMPLATE/).
- **Work Status:** canonical values (To Do, In Progress, In Review, Done, Blocked,
  Won't Do), mirrored in `docs/user-stories.md`.
- **MoSCoW priority:** label-based (`Must Have` / `Should Have` / `Could Have` /
  `Won't Have`).
- **Sprint membership:** the Sprint 1 **milestone** is the authoritative container
  for Sprint-selected PBIs.
- **MVP version:** tracked with the `mvp-v1` label / MVP-version field.
- **Task decomposition:** user stories are split into smaller linked technical
  PBIs so developers can start without further clarification.

## 14. Roadmap direction

Sprint 1 delivered the core ledger-safety and observability features. Sprint 2
targets the remaining `Should Have` stories (US-01, US-02, US-08) plus the two
customer follow-ups (concurrency regression on legacy paths; expiry-as-an-explicit
-transaction). Full roadmap: [docs/roadmap.md](../../docs/roadmap.md).

## 15. Verification evidence for completed MVP v1 PBIs

| PBI | Issue | PR(s) | Evidence |
|---|---|---|---|
| US-12 | [#11](https://github.com/Varriwon4ik/avito_bonus_point_service/issues/11) | [#13](https://github.com/Varriwon4ik/avito_bonus_point_service/pull/13) | Expiry showcase; 3 tests pass |
| US-11 | [#8](https://github.com/Varriwon4ik/avito_bonus_point_service/issues/8) | [#14](https://github.com/Varriwon4ik/avito_bonus_point_service/pull/14) | 3 concurrency tests (race, dedup, exhaustion) |
| US-13 | [#12](https://github.com/Varriwon4ik/avito_bonus_point_service/issues/12) | [#15](https://github.com/Varriwon4ik/avito_bonus_point_service/pull/15) | HTTP codes + OpenAPI updated |
| US-05 | [#3](https://github.com/Varriwon4ik/avito_bonus_point_service/issues/3) | [#16](https://github.com/Varriwon4ik/avito_bonus_point_service/pull/16) / [#17](https://github.com/Varriwon4ik/avito_bonus_point_service/pull/17) | Hold timeout sweep + tests |
| US-10 | [#7](https://github.com/Varriwon4ik/avito_bonus_point_service/issues/7) | [#18](https://github.com/Varriwon4ik/avito_bonus_point_service/pull/18) | Structured logs + `/metrics` |

## 16. Current product status

The MVP v1 increment runs via `docker compose up --build` and serves the API, web
UI, Swagger UI, and `/metrics`. All four MVP v1 `Must Have` stories plus US-10 are
implemented; the customer accepted the demonstrated increment with two follow-up
improvement requests.

## 17. Next steps

- Land the US-11 and US-12 customer follow-ups.
- Demonstrate US-05 live at the next review.
- Plan Sprint 2 around US-01, US-02, US-08.

## 18. Contribution traceability

| Member | GitHub | Role | Issues | PR(s) authored | Reviewed |
|---|---|---|---|---|---|
| Mikhail Ilin | [@Varriwon4ik](https://github.com/Varriwon4ik) | PM / Scrum Master / backend | US-05, US-10 | [#17](https://github.com/Varriwon4ik/avito_bonus_point_service/pull/17), [#18](https://github.com/Varriwon4ik/avito_bonus_point_service/pull/18) | — |
| Sergey Chuenko | [@SergeiCh07](https://github.com/SergeiCh07) | QA | US-12 | [#13](https://github.com/Varriwon4ik/avito_bonus_point_service/pull/13) | US-11, US-10 |
| Nurislam Denisov | [@NurikDen](https://github.com/NurikDen) | backend / QA | US-11 | [#14](https://github.com/Varriwon4ik/avito_bonus_point_service/pull/14) | US-13, US-05 |
| Sanzhar Kadambaev | [@kadambaevsanzhar](https://github.com/kadambaevsanzhar) | QA | US-13 | [#15](https://github.com/Varriwon4ik/avito_bonus_point_service/pull/15) | US-12 |
| N. Nuriev | [@farcimin](https://github.com/farcimin) | backend | US-05 (assignee) | — | — |

> **Participation note:** N. Nuriev was the assignee for US-05, but due to
> technical issues the implementation was completed by Mikhail and reviewed by
> Nurislam. The reduced individual participation for this member was discussed
> with and approved by the course instructor.

## 19–28. Required artifacts and links

- **SemVer release mapped to MVP v1:**
  [v1.0.0](https://github.com/Varriwon4ik/avito_bonus_point_service/releases/tag/v1.0.0)
- **Root changelog:** [CHANGELOG.md](../../CHANGELOG.md)
- **Process requirements:** [Process_Requirements.md](../../Process_Requirements.md)
- **Roadmap:** [docs/roadmap.md](../../docs/roadmap.md)
- **Definition of Done:** [docs/definition-of-done.md](../../docs/definition-of-done.md)
- **Issue templates:** [.github/ISSUE_TEMPLATE/](../../.github/ISSUE_TEMPLATE/)
- **Extended PR template:** [.github/pull_request_template.md](../../.github/pull_request_template.md)
- **Reviewed issue-linked PRs (Week 3):** [#13](https://github.com/Varriwon4ik/avito_bonus_point_service/pull/13),
  [#14](https://github.com/Varriwon4ik/avito_bonus_point_service/pull/14),
  [#15](https://github.com/Varriwon4ik/avito_bonus_point_service/pull/15),
  [#17](https://github.com/Varriwon4ik/avito_bonus_point_service/pull/17),
  [#18](https://github.com/Varriwon4ik/avito_bonus_point_service/pull/18)
- **Delivered MVP v1 access point:** [root README → Running](../../README.md#running)
  (`docker compose up --build`, API on `http://localhost:8080`)
- **Public sanitized video demonstration (< 2 min):**
  <!-- TODO: insert public link to the sanitized < 2 min demo video here -->
- **Sprint Review recording (customer review, for instructor/Moodle):**
  - Part 1: <!-- TODO: insert recording link -->
  - Part 2: <!-- TODO: insert recording link -->

## 29. Screenshots

<!-- Add the PNGs to reports/week3/images/ and embed them here. -->
- Product Backlog view — `images/product-backlog.png`
- Sprint Backlog view — `images/sprint-backlog.png`
- Sprint milestone — `images/sprint-milestone.png`
- MVP version grouped/filtered view — `images/mvp-version-view.png`
- SemVer release — `images/release.png`
- Delivered MVP v1 — `images/mvp-v1.png`
- Example reviewed issue-linked PR — `images/example-pr.png`

## 30–34. Review and reflection artifacts

- **Customer review transcript:** [customer-review-transcript.md](customer-review-transcript.md)
- **Customer review summary:** [customer-review-summary.md](customer-review-summary.md)
- **Reflection:** [reflection.md](reflection.md)
- **Retrospective:** [retrospective.md](retrospective.md)
- **LLM report:** [llm-report.md](llm-report.md)
