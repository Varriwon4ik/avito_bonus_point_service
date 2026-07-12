# Week 6 Report — Assignment 6 (Sprint 4, Trial Release)

> Canonical public Week 6 report and submission index for Assignment 6,
> Team 01 — Bonus Points Ledger Service. All public Week 6 evidence for the
> assignment is indexed here. Sections marked **"to be completed after the
> Week 6 customer session"** are filled in from that meeting before
> submission.

## 1. Project

**Bonus Points Ledger Service** — a REST-like service for managing an online
store's bonus-points program: configurable point expiry with validated TTL
bounds, transactional balance mutations with row-level locking, two-phase
debits (`hold` / `confirm` / `cancel`), idempotent operations (single and
bulk), FIFO-by-expiry consumption, paginated history and lots audit with
transaction labels, observability (structured logging + Prometheus
`/metrics`), a web UI with exact response feedback and a built-in autotester
(single- and multi-key modes), and Swagger/OpenAPI docs.

- Repository: <https://github.com/Varriwon4ik/avito_bonus_point_service>
- License: [LICENSE](../../LICENSE)

## 2. Planning and Sprint

- **Product Backlog board:** <https://github.com/users/Varriwon4ik/projects/1>
- **Sprint 4 Backlog platform board/view:** the same GitHub Projects board
  filtered to the Sprint 4 milestone (status, assignee, priority, Story-Point
  estimate) — milestone view:
  <https://github.com/Varriwon4ik/avito_bonus_point_service/milestone/4>
- **Sprint 4 milestone:**
  [Sprint 4](https://github.com/Varriwon4ik/avito_bonus_point_service/milestone/4)
  (kept separate from the `v2.1.0` release mapped to the trial increment)
- **Sprint 5 milestone (planned next Sprint):**
  [Sprint 5](https://github.com/Varriwon4ik/avito_bonus_point_service/milestone/5)
  — Goal and explicitly marked expected follow-up scope recorded in the
  milestone description and [docs/roadmap.md](../../docs/roadmap.md).
- **Sprint 4 Goal:** deliver the Week 6 trial / handover-candidate release
  that completes the customer's remaining `Should Have` scope — bulk accrual
  for promotional campaigns, the lots audit API for support tooling, and the
  multi-key parallel autotester requested at the Sprint 3 review — fix the
  defects found while trialling the web UI, and put the customer in a
  position to try the product independently and judge transition readiness
  against the reviewed customer-facing documentation set.
- **Sprint dates:** 6 – 12 July 2026 (Mon–Sun).
- **Scope summary:** the two long-planned `Should Have` stories (US-01 bulk
  accrual, US-02 lots audit API), the customer-requested US-19 multi-key
  autotester, and two bugs found by the team while trialling its own UI
  (#54 false autotester verdict, #60 missing UI controls) — plus the
  Assignment 6 customer-facing documentation set (customer handover,
  contributor/agent guidance, README polish).
- **Total Sprint 4 size:** **24 Story Points**
  (US-01 = 8, US-02 = 5, US-19 = 3, #54 = 3, #60 = 5).

## 3. Delivered Week 6 trial-release changes

| PBI | Item | SP | Issue | PR | Implementer | Reviewer |
|---|---|---|---|---|---|---|
| US-01 | Bulk points accrual (`POST /v1/accruals/batch`, 207 per-item results) + web UI bulk card | 8 | [#1](https://github.com/Varriwon4ik/avito_bonus_point_service/issues/1) | [#57](https://github.com/Varriwon4ik/avito_bonus_point_service/pull/57) (UI display via [#61](https://github.com/Varriwon4ik/avito_bonus_point_service/pull/61)) | @SergeiCh07 | @NurikDen |
| US-02 | Paginated lots audit API with `status` filtering | 5 | [#2](https://github.com/Varriwon4ik/avito_bonus_point_service/issues/2) | [#59](https://github.com/Varriwon4ik/avito_bonus_point_service/pull/59) | @kadambaevsanzhar | @Varriwon4ik |
| US-19 | Autotester multi-key parallel accrual scenarios | 3 | [#50](https://github.com/Varriwon4ik/avito_bonus_point_service/issues/50) | [#55](https://github.com/Varriwon4ik/avito_bonus_point_service/pull/55) | @NurikDen | @kadambaevsanzhar |
| Bug | Web autotester always reported "found issues" | 3 | [#54](https://github.com/Varriwon4ik/avito_bonus_point_service/issues/54) | [#58](https://github.com/Varriwon4ik/avito_bonus_point_service/pull/58) | @Varriwon4ik | @NurikDen |
| Bug | Test-mode selector and bulk-accrual card missing from served UI | 5 | [#60](https://github.com/Varriwon4ik/avito_bonus_point_service/issues/60) | [#61](https://github.com/Varriwon4ik/avito_bonus_point_service/pull/61) | @Varriwon4ik | @kadambaevsanzhar |

Plus the Assignment 6 Week 6 maintained assets (Course Task
[#62](https://github.com/Varriwon4ik/avito_bonus_point_service/issues/62)):
the maintained customer handover documentation, contributor and agent
guidance, the polished repository entry point, the Sprint 4/5 milestones, the
trial release, and this report set.

## 4. Deployment and access (Week 6 product access artifact)

- **Product access artifact (deployed trial):** `http://10.93.26.175:8080/`
  (web UI; Swagger UI at `/docs`; API) serving the Week 6 trial release
  [`v2.1.0`](https://github.com/Varriwon4ik/avito_bonus_point_service/releases/tag/v2.1.0).
  Hosted on the University VM — a private (RFC 1918) address reachable only on
  the university network/VPN. Private access details for graders are submitted
  through Moodle. The artifact stays accessible until grading is complete.
  <!-- TODO(team): confirm the VM is redeployed to v2.1.0 (docker compose up
       --build -d at the release commit) — the Sprint 3 retrospective
       action point is deploy-on-merge-day. -->
- **Access/run instructions:** [root README → Running](../../README.md#running)
  and [→ Deployment](../../README.md#deployment); self-hosting steps in
  [docs/customer-handover.md](../../docs/customer-handover.md).

## 5. Public repository entry point and customer-facing documentation set

The Assignment 6 Part 3 documentation set reviewed with the customer this
week:

- **Repository entry point:** [README.md](../../README.md) — quick-access
  block (deployed trial, hosted docs, handover, contributing/agents), current
  setup/run/API documentation including the new bulk-accrual endpoint and
  autotester test modes.
- **Contributor guidance:** [CONTRIBUTING.md](../../CONTRIBUTING.md) — **new**
  this week: environment setup, verification commands, the
  issue → branch → PR → review → merge workflow, and the
  docs-stay-current rules.
- **Agent guidance:** [AGENTS.md](../../AGENTS.md) — **new** this week:
  project shape, verification commands, and hard constraints for AI coding
  agents working in the repository.
- **Customer handover:** [docs/customer-handover.md](../../docs/customer-handover.md)
  — **new** this week: what is transferred/delegated/retained, configuration
  and secrets the customer must know, setup/deployment/recovery/verification
  steps, main documentation entry points, known limitations, and the current
  handover status.
- **Access / usage instructions:** [README → Running](../../README.md#running),
  [README → API](../../README.md#api), Swagger UI at `/docs` on any instance.
- **Deployment / installation:** [README → Deployment](../../README.md#deployment),
  [customer-handover → Setup, deployment, recovery, and verification](../../docs/customer-handover.md#setup-deployment-recovery-and-verification).
- **Troubleshooting / support notes and known limitations:**
  [customer-handover → Main documentation entry points](../../docs/customer-handover.md#main-documentation-entry-points)
  and [→ Known limitations](../../docs/customer-handover.md#known-limitations).
- **Hosted documentation site:**
  <https://varriwon4ik.github.io/avito_bonus_point_service/> — republished on
  every merge to `main`; the handover page is part of the site navigation.

### Customer review of the documentation set

> **To be completed after the Week 6 customer session.**
<!-- TODO(team): summarize what the customer found clear, unclear, or missing
     in README.md, docs/customer-handover.md, access/deployment instructions,
     troubleshooting notes, and known limitations; link the issues opened for
     each gap. -->

## 6. Transition-readiness summary

**Standing state going into the Week 6 meeting** (from the repository and
[docs/customer-handover.md](../../docs/customer-handover.md)):

- The full source, documentation, and CI configuration are public under MIT —
  already usable by the customer without the team.
- The trial instance is team-operated on the university VM; the product is
  fully self-hostable with Docker, but no customer-side deployment exists yet
  — that (or an explicit self-hosting agreement) is the main remaining
  transition action.
- There are no external services, accounts, or secrets to hand over; the API
  is internal-network-only by design.
- Remaining known Week 7 work regardless of meeting outcome: respond to trial
  feedback, complete the transition arrangement, release `MVP v3` with higher
  SemVer precedence, record the public sanitized demo video, confirm customer
  acceptance of the handover document, and prepare Demo Day.

**Meeting outcomes** (readiness for transition, parts still needing changes,
whether the customer already uses/deploys the product and if not why, what
must happen in Week 7, how to keep the product useful after final delivery):

> **To be completed after the Week 6 customer session.**
<!-- TODO(team): record the answers to each Part 5.2 discussion point and
     whether the customer (a) confirmed readiness for independent use after
     Week 7 work, (b) independently used the trial release, (c) deployed or
     operated it on their side. Convert identified problems into PBIs/issues
     and link them here. -->

## 7. Customer feedback response

Feedback standing after the Sprint 3 review (3 Jul), and how Sprint 4
responded:

| Feedback point | Resulting PBI or issue | Status | Response |
|---|---|---|---|
| Autotester should support multiple idempotency keys in parallel (raised 3 Jul review). | [US-19 / #50](https://github.com/Varriwon4ik/avito_bonus_point_service/issues/50) | **Addressed in Sprint 4** | Delivered in `v2.1.0` ([PR #55](https://github.com/Varriwon4ik/avito_bonus_point_service/pull/55)): N parallel accruals with N distinct keys, each applying exactly once; selectable "Test mode" in the web UI (opt-in fixed via [#60](https://github.com/Varriwon4ik/avito_bonus_point_service/issues/60)). |
| Deliver the remaining agreed stories US-01 / US-02 (Sprint 3 review expectation). | [US-01 / #1](https://github.com/Varriwon4ik/avito_bonus_point_service/issues/1), [US-02 / #2](https://github.com/Varriwon4ik/avito_bonus_point_service/issues/2) | **Addressed in Sprint 4** | Both shipped in `v2.1.0`: bulk accrual with per-item results + UI card; paginated lots audit API with status filters. |
| Run the autotester/regression checks against a demo (earlier) version of the product (carried from Sprint 2/3). | Tracked in [docs/roadmap.md](../../docs/roadmap.md) (Sprint 5 candidate) | Added to backlog for later | Sprint 4 capacity went to the three stories plus the trial-blocking bug fixes; the check is a Sprint 5 planning candidate alongside the transition work. |
| Model point expiry as an explicit auditable ledger transaction (Sprint 1 follow-up). | Tracked in [docs/roadmap.md](../../docs/roadmap.md); analysed in [ADR-002](../../docs/architecture/adr/ADR-002-lazy-expiry-and-fifo-by-expiry-consumption.md) | Deferred | Unchanged tradeoff, documented in ADR-002 and listed as a known limitation in [docs/customer-handover.md](../../docs/customer-handover.md); Sprint 5 candidate. |
| "Work as a team, as a whole" (process feedback, repeated). | [Week 5 retrospective action point 1](../week5/retrospective.md#action-points) | Addressed in process | Sprint 4 work was spread across all four members with cross-review pairs (see §14); results assessed in this week's [retrospective](retrospective.md). |

**New feedback from the Week 6 trial / documentation review:**

> **To be completed after the Week 6 customer session** — each new feedback
> point gets a row: feedback → resulting PBI/issue → status
> (addressed / partially / backlog / rejected with rationale).
<!-- TODO(team): add the Week 6 meeting feedback rows. -->

**Feedback not addressed this Sprint:** the demo-version regression run and
the expiry-as-explicit-transaction model — both deferred with the rationale
above, tracked in [docs/roadmap.md](../../docs/roadmap.md), and both explicit
Sprint 5 planning candidates.

## 8. Maintained documentation updated during Sprint 4

- [docs/roadmap.md](../../docs/roadmap.md) — Sprint 4 outcome, Sprint 5 plan,
  MVP v3 and end-of-course state (no speculative post-course planning)
- [docs/customer-handover.md](../../docs/customer-handover.md) — **new**
- [docs/user-stories.md](../../docs/user-stories.md) — US-01/US-02/US-19 →
  `Done` (Sprint 4)
- [docs/user-acceptance-tests.md](../../docs/user-acceptance-tests.md) — new
  scenarios UAT-007/008/009 for the Sprint 4 increment
- [docs/testing.md](../../docs/testing.md) — integration-test scope extended
  (batch accrual, lots audit)
- [docs/development-process.md](../../docs/development-process.md) — links the
  new contributor/agent guidance
- Unchanged this Sprint but current:
  [docs/quality-requirements.md](../../docs/quality-requirements.md),
  [docs/quality-requirement-tests.md](../../docs/quality-requirement-tests.md),
  [docs/definition-of-done.md](../../docs/definition-of-done.md),
  [docs/architecture/README.md](../../docs/architecture/README.md) (no
  structural, deployment, or decision changes in Sprint 4 — the new endpoints
  live inside the existing `internal/api` / `internal/data` boundaries per
  ADR-003)

## 9. Testing and CI status for the trial increment

- All Assignment 4/5 gates remain active and passing for the Sprint 4 work:
  `gofmt`, `go vet`, build, unit + integration tests with the race detector
  against real Postgres, automated QRTs (QRT-001/002/003), the per-module
  ≥30% coverage gate, `govulncheck`, and Lychee link checking.
- Automated verification was extended where the increment changed the
  product: batch-accrual per-item results and validation
  ([integration_test.go](../../internal/api/integration_test.go)), lots
  pagination/filtering ([lots_test.go](../../internal/api/lots_test.go)),
  multi-key parallel accrual in the shared autotest engine
  ([internal/autotest](../../internal/autotest/autotest.go)), and legacy
  lots-response compatibility
  ([lots_compat_test.go](../../internal/autotest/lots_compat_test.go)).
- New UAT scenarios UAT-007/008/009 cover the new user-facing behaviour
  ([docs/user-acceptance-tests.md](../../docs/user-acceptance-tests.md)).
- Full status: [docs/testing.md](../../docs/testing.md).
- **CI pipeline:** [.github/workflows/ci.yml](../../.github/workflows/ci.yml)
- **Latest protected-default-branch CI run:**
  [Actions › CI › branch:main](https://github.com/Varriwon4ik/avito_bonus_point_service/actions/workflows/ci.yml?query=branch%3Amain)

## 10. Release and changelog

- **Week 6 SemVer trial release:**
  [v2.1.0](https://github.com/Varriwon4ik/avito_bonus_point_service/releases/tag/v2.1.0)
  — tagged on the protected default branch; the release description
  identifies it as the Week 6 trial / handover-candidate release for
  Assignment 6 and links the Sprint 4 milestone, run/access instructions,
  [docs/customer-handover.md](../../docs/customer-handover.md), and this
  report. The final `MVP v3` release follows in Week 7 with higher SemVer
  precedence.
- **Changelog:** [CHANGELOG.md](../../CHANGELOG.md) (the `2.1.0` section lists
  every user-visible trial change, issue-linked).

## 11. UAT and customer-trial results (public, sanitized)

The relevant maintained scenarios for the changed behaviour are
**UAT-007** (bulk accrual, US-01), **UAT-008** (lots audit, US-02), and
**UAT-009** (multi-key autotester, US-19) —
[docs/user-acceptance-tests.md](../../docs/user-acceptance-tests.md).

> **To be completed after the Week 6 customer session:** which scenarios
> passed, which failed or still need changes, the most important feedback
> points, and the resulting PBIs or issues.
<!-- TODO(team): fill from the trial session; append the execution history in
     docs/user-acceptance-tests.md at the same time. -->

## 12. Sprint Review artifacts

- **Summary:** [sprint-review-summary.md](sprint-review-summary.md)
- **Transcript:** [sprint-review-transcript.md](sprint-review-transcript.md)
  <!-- TODO(team): after the session, either (a) publish the sanitized
       English transcript there if the customer permits public publication —
       recording permission, publication permission, and private-sharing
       permission must each be asked separately — or (b) replace this line
       with the refusal statement and share the transcript/notes only through
       Moodle, per the assignment's fallback rules. -->
- The private recording link and the exact private timecodes (customer trial,
  transition-readiness discussion, Sprint Review, and UAT segments of the
  session) are submitted only through Moodle.

## 13. Reflection, retrospective, LLM report

- [reflection.md](reflection.md)
- [retrospective.md](retrospective.md)
- [llm-report.md](llm-report.md)

## 14. Contribution traceability

| Member | GitHub | Roles | Issues / PRs (implementer) | Review activity | Testing / docs / transition / deployment |
|---|---|---|---|---|---|
| Mikhail Ilin | [@Varriwon4ik](https://github.com/Varriwon4ik) | PM; QA; Backend | Bug #54 ([PR #58](https://github.com/Varriwon4ik/avito_bonus_point_service/pull/58)); bug #60 ([PR #61](https://github.com/Varriwon4ik/avito_bonus_point_service/pull/61)) | Reviewed US-02 ([PR #59](https://github.com/Varriwon4ik/avito_bonus_point_service/pull/59)) | Sprint 4/5 milestones, `v2.1.0` trial release, [docs/customer-handover.md](../../docs/customer-handover.md), [CONTRIBUTING.md](../../CONTRIBUTING.md), [AGENTS.md](../../AGENTS.md), README polish, Week 6 report set (Course Task [#62](https://github.com/Varriwon4ik/avito_bonus_point_service/issues/62)), VM deployment |
| Nurislam Denisov | [@NurikDen](https://github.com/NurikDen) | QA; Backend | US-19 ([#50](https://github.com/Varriwon4ik/avito_bonus_point_service/issues/50), [PR #55](https://github.com/Varriwon4ik/avito_bonus_point_service/pull/55)) | Reviewed US-01 ([PR #57](https://github.com/Varriwon4ik/avito_bonus_point_service/pull/57)) and the #54 fix ([PR #58](https://github.com/Varriwon4ik/avito_bonus_point_service/pull/58)) | Multi-key check in the shared autotest engine + engine tests |
| Sanzhar Kadambaev | [@kadambaevsanzhar](https://github.com/kadambaevsanzhar) | Scrum Master; QA | US-02 ([#2](https://github.com/Varriwon4ik/avito_bonus_point_service/issues/2), [PR #59](https://github.com/Varriwon4ik/avito_bonus_point_service/pull/59)) | Reviewed US-19 ([PR #55](https://github.com/Varriwon4ik/avito_bonus_point_service/pull/55)) and the #60 fix ([PR #61](https://github.com/Varriwon4ik/avito_bonus_point_service/pull/61)) | Lots audit tests ([lots_test.go](../../internal/api/lots_test.go)), autotester legacy-response compatibility, coverage-gate dedupe fix |
| Sergey Chuenko | [@SergeiCh07](https://github.com/SergeiCh07) | Scrum Master; QA | US-01 ([#1](https://github.com/Varriwon4ik/avito_bonus_point_service/issues/1), [PR #57](https://github.com/Varriwon4ik/avito_bonus_point_service/pull/57)) | — (no recorded PR review this Sprint; see [retrospective](retrospective.md)) | Batch accrual TTL-bounds and invalid-label error mapping + integration tests |

## 15. Current product status and expected Week 7 follow-up

The Week 6 trial / handover-candidate release `v2.1.0` is delivered: every
`Must Have` and `Should Have` story in the Product Backlog is now `Done`, the
core ledger guarantees stay gate-verified, and the product is accompanied by
a reviewed customer-facing documentation set (entry point, handover,
contributor/agent guidance) so the customer can try it independently.

Expected Week 7 (Sprint 5) follow-up — see the
[Sprint 5 milestone](https://github.com/Varriwon4ik/avito_bonus_point_service/milestone/5)
and [docs/roadmap.md](../../docs/roadmap.md): respond to trial and
documentation-review feedback, complete the actual transition and confirm the
handover level and acceptance status, release `MVP v3`, record the public
sanitized demo video, and prepare the Demo Day presentation.

## 16. Example reviewed, issue-linked PR

[PR #59 — US-02: add paginated lots audit API with status filters](https://github.com/Varriwon4ik/avito_bonus_point_service/pull/59):
issue-linked (`Closes #2`), reviewed and approved by a different team member
(@Varriwon4ik), merged into the protected default branch via a merge commit
with all CI gates green.

## 17. Screenshots

Embedded PNG evidence for artifacts where public links may not be reliably
inspectable is collected in [images/](images/README.md): the Sprint 4
milestone, the `v2.1.0` release, the example reviewed issue-linked PR (#59),
the Sprint board view, the latest protected-branch CI run, the hosted docs
handover page, and the deployed trial UI (bulk accrual + test-mode selector).

**Deviation note (same as Week 5):** at first publication the PNG set is
being captured; every listed item is directly inspectable at its live link
throughout this report until the screenshots land in `images/` before
submission. Screenshots of the recorded customer session contain
customer-identifying information and are submitted only through Moodle.
<!-- TODO(team): capture the PNGs listed in images/README.md and embed the
     key ones here before submission. -->

## 18. Deviations from expected defaults

- **Sprint Review evidence form:** finalized after the Week 6 session —
  transcript if the customer permits publication (as in Weeks 3–5), otherwise
  the documented fallback (§12).
- **Screenshots:** pending at first publication — see §17.
- No other artifact-form deviations this week: boards, milestones, release
  mapping, changelog, CI gates, and report structure follow the shared
  requirements.
