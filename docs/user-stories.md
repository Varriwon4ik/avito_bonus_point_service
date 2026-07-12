# User Story Index

This file is the authoritative current registry of stable user-story IDs and
current user-story membership for the Bonus Points Ledger Service. It is kept
synchronized with the GitHub issue tracker. It is **not** the detailed Sprint
execution plan and does not duplicate full mutable story content — see each
linked issue for the live statement, acceptance criteria, and discussion.

- Historical Assignment 2 source: [reports/week2/user-stories.md](../reports/week2/user-stories.md)
- Process semantics (statuses, MoSCoW, Work Status, traceability): `Process_Requirements.md`
- Sprint 1 milestone: [Sprint 1](https://github.com/Varriwon4ik/avito_bonus_point_service/milestone/1)
- Sprint 2 milestone (Assignment 4): [Sprint 2](https://github.com/Varriwon4ik/avito_bonus_point_service/milestone/2)
- Sprint 3 milestone (Assignment 5): [Sprint 3](https://github.com/Varriwon4ik/avito_bonus_point_service/milestone/3)
- Sprint 4 milestone (Assignment 6, Week 6): [Sprint 4](https://github.com/Varriwon4ik/avito_bonus_point_service/milestone/4)
- Sprint 5 milestone (Assignment 6, Week 7): [Sprint 5](https://github.com/Varriwon4ik/avito_bonus_point_service/milestone/5)

Active stories are ordered by MoSCoW priority, then Sprint, then stable ID.
Removed stories are listed after all active stories.

| ID | Short title | MoSCoW priority | Issue | Requirement status | Work Status | Sprint |
|---|---|---|---|---|---|---|
| US-05 | Auto-release stale holds | Must Have | [#3](https://github.com/Varriwon4ik/avito_bonus_point_service/issues/3) | Active | Done | [Sprint 1](https://github.com/Varriwon4ik/avito_bonus_point_service/milestone/1) |
| US-11 | Concurrent idempotent-key deduplication tests | Must Have | [#8](https://github.com/Varriwon4ik/avito_bonus_point_service/issues/8) | Active | Done | [Sprint 1](https://github.com/Varriwon4ik/avito_bonus_point_service/milestone/1) |
| US-12 | Points removal / expiry system | Must Have | [#11](https://github.com/Varriwon4ik/avito_bonus_point_service/issues/11) | Active | Done | [Sprint 1](https://github.com/Varriwon4ik/avito_bonus_point_service/milestone/1) |
| US-13 | HTTP response codes & OpenAPI docs | Must Have | [#12](https://github.com/Varriwon4ik/avito_bonus_point_service/issues/12) | Active | Done | [Sprint 1](https://github.com/Varriwon4ik/avito_bonus_point_service/milestone/1) |
| US-14 | Continuous integration pipeline for every change | Must Have | [#28](https://github.com/Varriwon4ik/avito_bonus_point_service/issues/28) | Active | Done | [Sprint 2](https://github.com/Varriwon4ik/avito_bonus_point_service/milestone/2) |
| US-15 | Automated tests for points accrue (autotester) | Must Have | [#29](https://github.com/Varriwon4ik/avito_bonus_point_service/issues/29) | Active | Done | [Sprint 2](https://github.com/Varriwon4ik/avito_bonus_point_service/milestone/2) |
| US-10 | Structured request logging & metrics endpoint | Should Have | [#7](https://github.com/Varriwon4ik/avito_bonus_point_service/issues/7) | Active | Done | [Sprint 1](https://github.com/Varriwon4ik/avito_bonus_point_service/milestone/1) |
| US-09 | Pagination for transaction history | Should Have | [#6](https://github.com/Varriwon4ik/avito_bonus_point_service/issues/6) | Active | Done | [Sprint 2](https://github.com/Varriwon4ik/avito_bonus_point_service/milestone/2) |
| US-08 | Configurable per-accrual TTL validation and bounds | Should Have | [#5](https://github.com/Varriwon4ik/avito_bonus_point_service/issues/5) | Active | Done | [Sprint 3](https://github.com/Varriwon4ik/avito_bonus_point_service/milestone/3) |
| US-16 | Frontend HTTP responses implementation | Should Have | [#39](https://github.com/Varriwon4ik/avito_bonus_point_service/issues/39) | Active | Done | [Sprint 3](https://github.com/Varriwon4ik/avito_bonus_point_service/milestone/3) |
| US-17 | Frontend autotester implementation | Should Have | [#40](https://github.com/Varriwon4ik/avito_bonus_point_service/issues/40) | Active | Done | [Sprint 3](https://github.com/Varriwon4ik/avito_bonus_point_service/milestone/3) |
| US-18 | An option to put labels on transactions | Could Have | [#41](https://github.com/Varriwon4ik/avito_bonus_point_service/issues/41) | Active | Done | [Sprint 3](https://github.com/Varriwon4ik/avito_bonus_point_service/milestone/3) |
| US-01 | Bulk points accrual for promotional campaigns | Should Have | [#1](https://github.com/Varriwon4ik/avito_bonus_point_service/issues/1) | Active | Done | [Sprint 4](https://github.com/Varriwon4ik/avito_bonus_point_service/milestone/4) |
| US-02 | List and audit a user's points lots | Should Have | [#2](https://github.com/Varriwon4ik/avito_bonus_point_service/issues/2) | Active | Done | [Sprint 4](https://github.com/Varriwon4ik/avito_bonus_point_service/milestone/4) |
| US-19 | Autotester scenarios with multiple idempotency keys in parallel | Could Have | [#50](https://github.com/Varriwon4ik/avito_bonus_point_service/issues/50) | Active | Done | [Sprint 4](https://github.com/Varriwon4ik/avito_bonus_point_service/milestone/4) |
| US-07 | Manual bonus point accrual | — | [#4](https://github.com/Varriwon4ik/avito_bonus_point_service/issues/4) | Removed | — | — |
| US-03 | Earn bonus points after purchase | — | — | Removed | — | — |
| US-04 | Maintain automated regression coverage | — | — | Removed | — | — |
| US-06 | Confirm or cancel reserved points | — | — | Removed | — | — |

## Notes on removed and descoped stories

- **US-03, US-04, US-06** were proposed in Assignment 2 but **removed** during
  Assignment 3 refinement because they are already satisfied by the delivered
  MVP v0 base functionality (purchase-driven accrual, the existing regression
  suite, and the two-phase `hold` / `confirm` / `cancel` flow). Their stable IDs
  are preserved here for traceability; they are not re-issued or reused.
- **US-09** (transaction-history pagination) was `Won't Have` after the
  Assignment 3 customer negotiation, then **reconsidered and delivered in
  Sprint 2** (Assignment 4) in response to the customer's request to demonstrate
  paginated access. It is now `Should Have`, `Active`, and `Done`.
- **US-07** (manual bonus point accrual) was implemented during Sprint 2
  ([PR #32](https://github.com/Varriwon4ik/avito_bonus_point_service/pull/32))
  but **reverted** ([PR #34](https://github.com/Varriwon4ik/avito_bonus_point_service/pull/34))
  after bugs and integration issues surfaced in review, so it never shipped in a
  tagged release. Its requirement status is now **`Removed`**; the stable ID is
  preserved for traceability and is not re-issued or reused. The team
  prioritized a different feature in its place for the Sprint.
- **US-18** (transaction labels) was merged during Sprint 3
  ([PR #44](https://github.com/Varriwon4ik/avito_bonus_point_service/pull/44)),
  **reverted** ([PR #47](https://github.com/Varriwon4ik/avito_bonus_point_service/pull/47))
  after issues surfaced, and then **re-landed fixed** within the same Sprint
  ([PR #49](https://github.com/Varriwon4ik/avito_bonus_point_service/pull/49),
  together with US-17). Unlike US-07, the story shipped in `v2.0.0` and remains
  `Active`/`Done`; the revert history is preserved in the linked PRs.
- **US-19** was raised by the customer during the Sprint 3 Review / UAT session
  (3 Jul 2026): extend the autotester to run parallel scenarios with multiple
  distinct idempotency keys. It was estimated (3 SP), selected into Sprint 4,
  and **delivered** in the Week 6 trial release
  ([PR #55](https://github.com/Varriwon4ik/avito_bonus_point_service/pull/55);
  its first landing ran the multi-key check on every autotester run — the
  intended opt-in mode switch was added with the
  [#60](https://github.com/Varriwon4ik/avito_bonus_point_service/issues/60)
  fix in [PR #61](https://github.com/Varriwon4ik/avito_bonus_point_service/pull/61)).

## Stable ID allocation

The highest allocated stable ID is **US-19** (allocated in Sprint 3). New user
stories discovered after migration receive the next unused ID starting at
**US-20**.
