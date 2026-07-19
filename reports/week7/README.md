# Week 7 Report — Assignment 6 (Sprint 5, MVP v3)

> Canonical public Week 7 report and final Assignment 6 submission index for
> Team 01 — Bonus Points Ledger Service. Private recording, rehearsal, consent,
> customer-identifying, and access evidence is kept only in the Moodle PDF.

## 1. Previous evidence and Sprint planning

- **Complete Week 6 evidence:** [reports/week6/README.md](../week6/README.md)
- **Product Backlog board:** <https://github.com/users/Varriwon4ik/projects/1>
- **Sprint 5 Backlog view:** the Product Backlog board filtered by the Sprint 5
  milestone; [platform milestone view](https://github.com/Varriwon4ik/avito_bonus_point_service/milestone/5)
- **Sprint 5 milestone:** [Sprint 5](https://github.com/Varriwon4ik/avito_bonus_point_service/milestone/5)
- **Dates:** 13–19 July 2026
- **Goal:** respond to Week 6 customer-trial and documentation-review feedback,
  complete the actual transition, and deliver final course version **MVP v3**,
  leaving the Customer able to use, verify, and operate the ledger with the
  maintained documentation set.
- **Scope:** #64 horizontal-scaling assessment and migration-startup hardening;
  final transition confirmation; maintained documentation/report updates;
  MVP v3 release preparation; sanitized demo and Demo Day preparation.
- **Total recorded Sprint size:** **0 Story Points visible on GitHub.** Sprint 5
  contains one inspectable PBI, #64, but its issue body left Story Points,
  assignee, and reviewer to Sprint planning and those fields were not recorded.
  PR #66 identifies Mikhail as author and Sanzhar as approving reviewer. This is
  disclosed as a traceability deviation, not backfilled with an invented value.

## 2. Week 7 follow-up and MVP v3

| Feedback / PBI | Result | Issue / PR | Implementer | Reviewer |
|---|---|---|---|---|
| Explicitly assess horizontal scaling, a Customer condition for final delivery | Architecture now states that the stateless API tier can scale horizontally with conditions; ADR-006 records the audit and caveats | [#64](https://github.com/Varriwon4ik/avito_bonus_point_service/issues/64) / [PR #66](https://github.com/Varriwon4ik/avito_bonus_point_service/pull/66) | [@Varriwon4ik](https://github.com/Varriwon4ik) | [@kadambaevsanzhar](https://github.com/kadambaevsanzhar) (approved) |
| Concurrent replica startup could race migrations | `data.Migrate` now serializes migrations with a PostgreSQL advisory lock | [#64](https://github.com/Varriwon4ik/avito_bonus_point_service/issues/64) / [PR #66](https://github.com/Varriwon4ik/avito_bonus_point_service/pull/66) | @Varriwon4ik | @kadambaevsanzhar |
| Missing Bulk accrual and Test-mode controls in Week 6 trial UI | Fixed, rebuilt, redeployed, and visible in the final product demonstrated to the Customer | [#60](https://github.com/Varriwon4ik/avito_bonus_point_service/issues/60) / [PR #61](https://github.com/Varriwon4ik/avito_bonus_point_service/pull/61) | @Varriwon4ik | @kadambaevsanzhar |

MVP v3 therefore contains the full Week 6 trial increment plus both completion
conditions from the Customer: corrected final UI delivery and explicit,
technically hardened horizontal-scaling guidance.

## 3. Final access and maintained documentation

- **Final product access:** `http://10.93.26.175:8080/` (University network or
  VPN; web UI, `/docs` Swagger, `/openapi.yaml`; kept available through grading)
- **Current run/access instructions:** [README → Running](../../README.md#running)
  and [→ Deployment](../../README.md#deployment)
- **Repository entry point:** [README.md](../../README.md)
- **Contributor guidance:** [CONTRIBUTING.md](../../CONTRIBUTING.md)
- **Agent guidance:** [AGENTS.md](../../AGENTS.md)
- **Customer handover:** [docs/customer-handover.md](../../docs/customer-handover.md)
- **Hosted documentation:** <https://varriwon4ik.github.io/avito_bonus_point_service/>
- **Roadmap/final course state:** [docs/roadmap.md](../../docs/roadmap.md)
- **Architecture/scaling:** [docs/architecture/README.md](../../docs/architecture/README.md)
  and [ADR-006](../../docs/architecture/adr/ADR-006-horizontal-scaling-stateless-api-over-single-postgres.md)
- **Quality/testing/UAT:** [quality requirements](../../docs/quality-requirements.md),
  [quality requirement tests](../../docs/quality-requirement-tests.md),
  [testing](../../docs/testing.md), [UAT](../../docs/user-acceptance-tests.md),
  and [Definition of Done](../../docs/definition-of-done.md)

## 4. Final transition outcome

- **Handover level:** **Ready for independent use**. The Customer independently
  used the Week 6 trial; the final product is self-hostable from the public
  repository and maintained documentation.
- **Customer-confirmation status:** **Accepted**. After the 19 July final
  demonstration, the Customer assessed the product as complete and working
  well. No new product change was requested.
- **Transferred/made available:** public MIT-licensed source and history,
  container setup, migrations, CI configuration, hosted and repository docs,
  configuration/secrets guidance, operational verification and recovery steps,
  API/OpenAPI artifacts, and the team-operated evaluation deployment. Exact
  boundaries are maintained in [customer-handover.md](../../docs/customer-handover.md#what-is-transferred-delegated-or-retained).
- **Not transferred:** original repository administration and the University VM
  remain team-operated through grading. No production secret or external
  service exists to transfer.
- **Not customer-side deployed:** the Customer's company infrastructure is
  confidential and has a higher operational/security bar; their interns will
  evaluate and deploy after delivery. This blocker is on the customer/external
  coordination side, not a missing self-hosting capability.
- **Support:** the team keeps the evaluation instance available through grading;
  no continuing product support was requested after delivery.

## 5. Customer feedback response

| Feedback | Sprint 5 response | Status |
|---|---|---|
| Fix the stale deployed UI controls (#60) | Same-day Week 6 merge/redeploy; final product demo completed without a new defect | Addressed |
| Explicitly say whether horizontal scaling is supported (#64) | Added architecture verdict, ADR-006, handover caveats, and advisory-locked migrations | Addressed |
| Complete final delivery | Maintained handover/report set finalized; MVP v3 release candidate prepared; Customer assessed product complete | Addressed; GitHub release publication remains an administrative step |

No new product or documentation follow-up item was requested in the final
review. The Customer's attendance at the 21 July final presentation remained
unconfirmed because the university/customer coordination contact had not yet
responded.

## 6. UAT and verification

The Customer observed the final product demo and accepted the overall final
increment as complete and working well. The condensed transcript does not map
individual steps to UAT-007/UAT-009, so this report does not fabricate separate
scenario-pass claims. Their Week 6 results and the transparent Week 7 evidence
boundary are maintained in [docs/user-acceptance-tests.md](../../docs/user-acceptance-tests.md).

All prior automated gates remain applicable: formatting, vet, build, real-
Postgres race-enabled tests, QRTs, module coverage gate, `govulncheck`, and link
checking. PR #66 records local verification and passed protected-branch CI.

## 7. Release, changelog, and public demo

- **Final release mapped to MVP v3:**
  [`v3.0.0`](https://github.com/Varriwon4ik/avito_bonus_point_service/releases/tag/v3.0.0)
  — release notes are prepared; publication must use the final protected
  `main` submission commit and link this report, Sprint 5, handover, run
  instructions, and the public demo.
- **Changelog:** [CHANGELOG.md](../../CHANGELOG.md)
- **Public sanitized MVP v3 demo video:**
  <https://drive.google.com/file/d/1OsT4TuF2ABjpJcbTyRDUEZU4pkpucZa7/view?usp=drive_link>

## 8. Demo Day preparation

The team completed the required Week 7 rehearsal preparation and prepared the
final narrative around project context, prioritized requirements, MVP roadmap,
handover status, engineering/quality evidence, contribution, reflection, and a
pre-recorded demo. The private rehearsal-video link and presentation materials
are kept out of this public repository and are submitted through Moodle.

## 9. Sprint Review and team reflection

- [Public sanitized Sprint Review transcript](sprint-review-transcript.md)
- [Sprint Review summary](sprint-review-summary.md)
- [Reflection](reflection.md)
- [Retrospective](retrospective.md)
- [LLM usage report](llm-report.md)

Recording permission was requested and granted. The private recording link,
consent timecode, and exact activity timecodes are in the Moodle wrapper only.

## 10. Contribution traceability

| Member | Roles | Sprint 5 contribution |
|---|---|---|
| Mikhail Ilin ([@Varriwon4ik](https://github.com/Varriwon4ik)) | PM; QA; Backend | Implemented #64 / PR #66: scaling audit, ADR-006, architecture/handover updates, migration advisory lock; coordinated final documentation, release candidate, transition meeting, and Moodle evidence |
| Sanzhar Kadambaev ([@kadambaevsanzhar](https://github.com/kadambaevsanzhar)) | Scrum Master; QA | Approved and merged PR #66; review/quality evidence and final presentation preparation |
| Nurislam Denisov ([@NurikDen](https://github.com/NurikDen)) | QA; Backend | Presented the final product demonstration; final regression/demo and presentation preparation |
| Sergey Chuenko ([@SergeiCh07](https://github.com/SergeiCh07)) | Scrum Master; QA | Final verification and presentation preparation; no Sprint 5 implementation PR attributed in repository history |

No listed team member is reported as wholly non-participating in Sprint 5;
repository-visible implementation/review work was concentrated in Mikhail and
Sanzhar, as disclosed above.

## 11. Final product status

The course outcome is **MVP v3**: all agreed `Must Have` and `Should Have`
stories are delivered, the Customer's Week 6 completion conditions are
resolved, maintained customer/engineering documentation is current, quality
gates remain active, and the product is **Ready for independent use / Accepted**.
Customer-side operation remains outside the reached level for the documented
customer-infrastructure reason.

## 12. Screenshot evidence and deviations

The required capture list is in [images/README.md](images/README.md). The public
repository currently provides direct platform links for the milestone, issue,
PR, CI, docs, deployment, and release target. Screenshots still need capture and
embedding before final submission.

Known evidence deviations:

1. Sprint 5 issue #64 has no recorded Story Point estimate or assignee; the PR
   provides author and independent approval evidence.
2. The `v3.0.0` release and milestone closure require GitHub-side actions after
   this report reaches protected `main`; they cannot be truthfully claimed as
   already published in a local draft.
3. Named Week 7 UAT-007/UAT-009 re-executions cannot be separated from the
   supplied condensed transcript; only overall final-demo acceptance is claimed.
4. Screenshots are listed but not yet captured.
