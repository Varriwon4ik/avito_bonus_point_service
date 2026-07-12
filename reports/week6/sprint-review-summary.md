# Sprint Review Summary — Sprint 4 (Assignment 6, Week 6)

- **Date:** 10 July 2026
- **Sprint:** Sprint 4 / Assignment 6 Week 6 (6 – 12 July 2026)
- **Milestone:** [Sprint 4](https://github.com/Varriwon4ik/avito_bonus_point_service/milestone/4)
- **Participants (roles):** Mikhail Ilin (Project Manager; QA Engineer;
  Backend Developer), Nurislam Denisov (QA Engineer; Backend Developer),
  Sanzhar Kadambaev (Scrum Master; QA Engineer), and the Customer. Sergey
  Chuenko (Scrum Master; QA Engineer) could not attend this session; his
  feature slot (US-01) was covered by the teammates.
- **Recording and permissions:** One recorded session covers the Sprint
  Review, the customer trial / UAT, the documentation review, and the
  transition-readiness discussion. The Customer granted **recording
  permission** before recording started and **permitted public transcript
  publication** — see [sprint-review-transcript.md](sprint-review-transcript.md).
  The private recording link and exact private timecodes are submitted
  through Moodle only.
- **Format:** each present team member drove one of the new UAT scenarios in
  front of the Customer against the deployed trial instance: Sanzhar —
  UAT-007 (bulk accrual), Mikhail — UAT-008 (lots audit, covering for the
  absent Sergey), Nurislam — UAT-009 (multi-key autotester).

## Sprint Goal reviewed

Deliver the Week 6 trial / handover-candidate release that completes the
customer's remaining `Should Have` scope — bulk accrual for promotional
campaigns, the lots audit API for support tooling, and the multi-key parallel
autotester requested at the Sprint 3 review — fix the defects found while
trialling the web UI, and put the customer in a position to try the product
independently and judge transition readiness against the reviewed
customer-facing documentation set.

## Delivered increment discussed (Week 6 trial release `v2.1.0`)

| PBI | Item | Driven by | Outcome |
|---|---|---|---|
| US-19 ([#50](https://github.com/Varriwon4ik/avito_bonus_point_service/issues/50)) | Multi-key parallel autotester | Nurislam (UAT-009) | Engine check accepted (shown via console); **UI demo failed** — the Test-mode selector was missing from the deployed page ([#60](https://github.com/Varriwon4ik/avito_bonus_point_service/issues/60)) |
| US-02 ([#2](https://github.com/Varriwon4ik/avito_bonus_point_service/issues/2)) | Paginated lots audit API with status filters | Mikhail (UAT-008) | **Accepted** — "works as I would expect" |
| US-01 ([#1](https://github.com/Varriwon4ik/avito_bonus_point_service/issues/1)) | Bulk points accrual (endpoint + UI card) | Sanzhar (UAT-007) | Endpoint behaviour accepted (shown via Swagger); **UI demo failed** — the bulk card was missing from the deployed page ([#60](https://github.com/Varriwon4ik/avito_bonus_point_service/issues/60)) |
| Bug [#54](https://github.com/Varriwon4ik/avito_bonus_point_service/issues/54) | Autotester false "found issues" verdict fix | — | Superseded by the #60 discussion; verified after the same-day redeploy |

The deployed instance was serving a stale embedded UI ("improper
deployment"), so the new interface controls were absent during the demo. The
Customer asked the team to **fix everything by next week**; the fix was
tracked as [#60](https://github.com/Varriwon4ik/avito_bonus_point_service/issues/60),
merged the same day
([PR #61](https://github.com/Varriwon4ik/avito_bonus_point_service/pull/61)),
and the VM was redeployed the same day (10 Jul) to the commit later tagged
[`v2.1.0`](https://github.com/Varriwon4ik/avito_bonus_point_service/releases/tag/v2.1.0).

## Customer-facing documentation review results

Reviewed set: [README.md](../../README.md),
[docs/customer-handover.md](../../docs/customer-handover.md), access/run and
deployment instructions, troubleshooting notes, known limitations,
[CONTRIBUTING.md](../../CONTRIBUTING.md), [AGENTS.md](../../AGENTS.md).

- **Clear / complete:** the Customer called the documentation complete
  overall — "the READMEs cover what I need"; nothing was flagged as unclear
  or missing for normal use, setup, or deployment.
- **Requested addition:** the architecture documentation must state
  **explicitly whether the service can be scaled horizontally or not**, with
  the reasoning → converted into
  [#64](https://github.com/Varriwon4ik/avito_bonus_point_service/issues/64)
  (`Must Have`, Sprint 5).

## Transition-readiness findings

- **Complete enough for transition?** Yes — the Customer is satisfied with
  the current version; no parts are missing and the main features are
  implemented. The two conditions for considering the delivery complete:
  fix the UI display issue (done same day, #60 / PR #61) and assess
  horizontal-scaling suitability
  ([#64](https://github.com/Varriwon4ik/avito_bonus_point_service/issues/64)).
- **Is the customer already using the product?** No. Their production
  context is a large-company environment: a simple deployment on the
  university VM with its weak security posture is not acceptable there, so
  using the trial instance directly is not an option for them.
- **Deployed / operated on the customer side?** No, and not planned before
  delivery: the customer's internal deployment structure is confidential;
  after receiving the final delivery, a group of the customer's interns will
  evaluate the project and perform the customer-side deployment.
- **What must happen in Week 7:** fix the UI bug (already resolved), assess
  and explicitly document horizontal scaling, polish, and complete the final
  delivery (`MVP v3`) — at which point the Customer considers the project
  complete.
- **Keeping the product useful after delivery:** the Customer will take care
  of the project on his own after the handover; no ongoing team involvement
  was requested.

Recorded per Assignment 6 Part 5.5 — the Customer:

- confirmed the product will be **ready for independent use** after the
  Week 7 work: **yes**;
- **independently used the trial release** during the session: **yes**;
- **deployed or operated it on their side**: **no** (by their own choice and
  timeline — see above).

## Customer trial / UAT results

Executed with the Customer on 10 July 2026
([docs/user-acceptance-tests.md](../../docs/user-acceptance-tests.md)):

- **UAT-007** Run a bulk accrual campaign with per-item results (US-01) —
  **Failed** at the session: the bulk card was missing from the deployed UI
  (#60); the endpoint behaviour itself was accepted via Swagger.
- **UAT-008** Audit a user's points lots with status filtering (US-02) —
  **Passed**.
- **UAT-009** Run the multi-key parallel autotester scenario (US-19) —
  **Failed** at the session: the Test-mode selector was missing from the
  deployed UI (#60); the engine check was shown passing from the console.
- Both failures share one root cause — the stale embedded UI in the deployed
  build — fixed and redeployed the same day; UAT-007/009 are re-executed at
  the Week 7 transition confirmation.

## Feedback, risks, and action points

- **Change request:** UI controls missing from the deployed page →
  [#60](https://github.com/Varriwon4ik/avito_bonus_point_service/issues/60),
  fixed same day ([PR #61](https://github.com/Varriwon4ik/avito_bonus_point_service/pull/61)),
  VM redeployed 10 Jul.
- **New requirement:** horizontal-scaling assessment with an explicit
  statement in the architecture docs →
  [#64](https://github.com/Varriwon4ik/avito_bonus_point_service/issues/64)
  (`Must Have`, Sprint 5).
- **Risk:** a stale deployment broke a customer-facing demo for the second
  Sprint running (Sprint 3: undeployed build; Sprint 4: stale embedded UI) —
  addressed by the [retrospective](retrospective.md) action points
  (feature-freeze buffer, documentation-driven post-deploy smoke test).

## Resulting Product Backlog changes

- [#64](https://github.com/Varriwon4ik/avito_bonus_point_service/issues/64)
  created from the documentation-review feedback and added to the
  [Sprint 5 milestone](https://github.com/Varriwon4ik/avito_bonus_point_service/milestone/5)
  as `Must Have` (a customer condition for final delivery); mirrored in
  [docs/roadmap.md](../../docs/roadmap.md).
- No other scope changes: the remaining Sprint 5 plan (final transition,
  `MVP v3`, demo video, Demo Day preparation) stands as planned.
