# Sprint Review Summary — Sprint 4 (Assignment 6, Week 6)

<!-- TODO(team): fill every «…» placeholder from the Week 6 session before
     submission. One recorded session is expected to cover the Sprint Review,
     the transition-readiness discussion, the customer documentation review,
     and the customer trial / UAT (per Assignment 6 Parts 5, 9, 10). -->

- **Date:** «date of the Week 6 session, July 2026»
- **Sprint:** Sprint 4 / Assignment 6 Week 6 (6 – 12 July 2026)
- **Milestone:** [Sprint 4](https://github.com/Varriwon4ik/avito_bonus_point_service/milestone/4)
- **Participants (roles):** Mikhail Ilin (Project Manager; QA Engineer;
  Backend Developer), Nurislam Denisov (QA Engineer; Backend Developer),
  Sanzhar Kadambaev (Scrum Master; QA Engineer), Sergey Chuenko (Scrum
  Master; QA Engineer), and the Customer.
- **Recording and permissions:** «record whether recording permission was
  asked and granted before recording started; whether public transcript
  publication was permitted; and, if publication was refused, whether private
  instructor sharing was permitted». The private recording link and exact
  private timecodes (trial, transition discussion, review, UAT segments) are
  submitted through Moodle only.
- **Format:** the customer trials the Week 6 release
  ([`v2.1.0`](https://github.com/Varriwon4ik/avito_bonus_point_service/releases/tag/v2.1.0))
  independently or with minimal guidance where practical, and reviews the
  customer-facing documentation set; the session also covers the
  transition-readiness discussion points of Assignment 6 Part 5.

## Sprint Goal reviewed

Deliver the Week 6 trial / handover-candidate release that completes the
customer's remaining `Should Have` scope — bulk accrual for promotional
campaigns, the lots audit API for support tooling, and the multi-key parallel
autotester requested at the Sprint 3 review — fix the defects found while
trialling the web UI, and put the customer in a position to try the product
independently and judge transition readiness against the reviewed
customer-facing documentation set.

## Delivered increment discussed (Week 6 trial release `v2.1.0`)

| PBI | Item | Demonstrated / trialled | Outcome |
|---|---|---|---|
| US-01 ([#1](https://github.com/Varriwon4ik/avito_bonus_point_service/issues/1)) | Bulk points accrual (endpoint + UI card) | «who/how» | «accepted / changes requested» |
| US-02 ([#2](https://github.com/Varriwon4ik/avito_bonus_point_service/issues/2)) | Paginated lots audit API with status filters | «who/how» | «…» |
| US-19 ([#50](https://github.com/Varriwon4ik/avito_bonus_point_service/issues/50)) | Multi-key parallel autotester (Test-mode selector) | «who/how» | «…» |
| Bugs [#54](https://github.com/Varriwon4ik/avito_bonus_point_service/issues/54) / [#60](https://github.com/Varriwon4ik/avito_bonus_point_service/issues/60) | Autotester verdict fix; restored UI controls | «who/how» | «…» |

## Customer-facing documentation review results

Reviewed set: [README.md](../../README.md),
[docs/customer-handover.md](../../docs/customer-handover.md), access/run and
deployment instructions, troubleshooting notes, known limitations,
[CONTRIBUTING.md](../../CONTRIBUTING.md), [AGENTS.md](../../AGENTS.md).

«What the customer found clear / unclear / missing; resulting issues.»

## Transition-readiness findings

«Per Assignment 6 Part 5.2: completeness for transition; parts ready vs.
needing changes; whether the customer already uses the product (and how) or
why not; whether it is deployed/operated on the customer side or what blocks
that; what must happen in Week 7 to complete transition; how to keep the
product useful after final delivery.»

Record explicitly whether the customer: confirmed readiness for independent
use after Week 7 work — «yes/no»; independently used the trial release —
«yes/no»; deployed or operated it on their side — «yes/no».

## Customer trial / UAT results

Relevant maintained scenarios: **UAT-007** (bulk accrual), **UAT-008** (lots
audit), **UAT-009** (multi-key autotester) —
[docs/user-acceptance-tests.md](../../docs/user-acceptance-tests.md).

«Passed / failed per scenario; the most important feedback points; the
resulting PBIs or issues. Append the same results to the execution history in
docs/user-acceptance-tests.md.»

## Feedback, risks, and action points

«New feature requests, documentation gaps, deployment/access problems — each
converted into a traceable PBI, issue, or explicit transition action; risks;
action points.»

## Resulting Product Backlog changes

«Issues created or reprioritized during the session, added to the
[Sprint 5 milestone](https://github.com/Varriwon4ik/avito_bonus_point_service/milestone/5)
where committed; mirror material changes in
[docs/roadmap.md](../../docs/roadmap.md).»
