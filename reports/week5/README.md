# Week 5 Report — Assignment 5 (Sprint 3, MVP v2)

> Canonical public Week 5 report and submission index for Assignment 5,
> Team 01 — Bonus Points Ledger Service. All public evidence for the assignment
> is indexed here.

## 1. Project

**Bonus Points Ledger Service** — a REST-like service for managing an online
store's bonus-points program: configurable point expiry with validated TTL
bounds, transactional balance mutations with row-level locking, two-phase
debits (`hold` / `confirm` / `cancel`), idempotent operations, FIFO-by-expiry
consumption, paginated history with transaction labels, observability
(structured logging + Prometheus `/metrics`), a web UI with exact response
feedback and a built-in autotester, and Swagger/OpenAPI docs.

- Repository: <https://github.com/Varriwon4ik/avito_bonus_point_service>
- License: [LICENSE](../../LICENSE)

## 2. Planning and Sprint

- **Product Backlog board:** <https://github.com/users/Varriwon4ik/projects/1>
- **Sprint Backlog platform board/view:** the same GitHub Projects board
  filtered to the Sprint 3 milestone (showing status, assignee, priority, and
  Story-Point estimate) — milestone view:
  <https://github.com/Varriwon4ik/avito_bonus_point_service/milestone/3>
- **Sprint 3 milestone:**
  [Sprint 3](https://github.com/Varriwon4ik/avito_bonus_point_service/milestone/3)
  (kept separate from the `v2.0.0` release mapped to MVP v2)
- **Sprint Goal:** deliver MVP v2 by making the ledger directly usable and
  verifiable from the web UI — exact HTTP response feedback, a browser-based
  autotester, transaction labels, and enforced per-accrual TTL bounds — while
  documenting the architecture, key decisions (ADRs), and the development
  process so the product can keep evolving safely.
- **Sprint dates:** 29 June – 5 July 2026 (Mon–Sun).
- **Scope summary:** a UI-facing, verification-focused Sprint — US-16 (exact
  HTTP responses in the UI), US-17 (web autotester), US-18 (transaction
  labels), US-08 (TTL validation bounds) — plus the Assignment 5 maintained
  documentation (architecture views, ADRs, development process, hosted docs).
- **Total Sprint size:** **16 Story Points**
  (US-16 = 5, US-17 = 5, US-08 = 3, US-18 = 3).

## 3. Delivered MVP v2 changes

| PBI | Item | SP | Issue | PR | Implementer | Reviewer |
|---|---|---|---|---|---|---|
| US-16 | Exact HTTP response codes in the web UI | 5 | [#39](https://github.com/Varriwon4ik/avito_bonus_point_service/issues/39) | [#46](https://github.com/Varriwon4ik/avito_bonus_point_service/pull/46) | @NurikDen | @Varriwon4ik |
| US-17 | Web autotester tab (`POST /v1/autotest/run`) | 5 | [#40](https://github.com/Varriwon4ik/avito_bonus_point_service/issues/40) | [#49](https://github.com/Varriwon4ik/avito_bonus_point_service/pull/49) | @Varriwon4ik | @NurikDen |
| US-18 | Labels on transactions (preset + custom) | 3 | [#41](https://github.com/Varriwon4ik/avito_bonus_point_service/issues/41) | [#44](https://github.com/Varriwon4ik/avito_bonus_point_service/pull/44) → revert [#47](https://github.com/Varriwon4ik/avito_bonus_point_service/pull/47) → re-landed in [#49](https://github.com/Varriwon4ik/avito_bonus_point_service/pull/49) | @kadambaevsanzhar | @SergeiCh07 |
| US-08 | Per-accrual TTL validation and bounds | 3 | [#5](https://github.com/Varriwon4ik/avito_bonus_point_service/issues/5) | [#42](https://github.com/Varriwon4ik/avito_bonus_point_service/pull/42) | @SergeiCh07 | @kadambaevsanzhar |

Plus the Assignment 5 maintained assets (Course Task
[#51](https://github.com/Varriwon4ik/avito_bonus_point_service/issues/51)):
architecture documentation with three views and five ADRs, the documented
development process and configuration management, the hosted documentation
site, and this report set.

## 4. Deployment and access

- **Product access artifact (deployed product):** `http://10.93.26.175:8080/`
  (web UI; Swagger UI at `/docs`; API). Hosted on the University VM — a private
  (RFC 1918) address reachable only on the university network/VPN. Private
  access details for graders are submitted through Moodle. The artifact stays
  accessible until grading is complete.
- **Access/run instructions:** [root README → Running](../../README.md#running)
  and [→ Deployment](../../README.md#deployment).

## 5. Customer feedback response

Feedback on `MVP v1` and the Sprint 2 review, and how Sprint 3 responded:

| Feedback point | Resulting PBI or issue | Status | Response |
|---|---|---|---|
| Provide additional automated tests / more ways to verify that the team's changes are valid (carried from Sprint 1–2). | [US-17 / #40](https://github.com/Varriwon4ik/avito_bonus_point_service/issues/40) (web autotester), [US-16 / #39](https://github.com/Varriwon4ik/avito_bonus_point_service/issues/39) (exact response codes) | Partially addressed | The autotester is now runnable and readable from the browser by a non-developer (incl. the customer), and every UI operation shows its exact backend response. Running the checks against a **demo/earlier version** of the product is still open — carried to Sprint 4 ([docs/roadmap.md](../../docs/roadmap.md)). |
| (Sprint 1 follow-up) Model point expiry as an explicit ledger transaction so expirations are auditable. | Tracked in [docs/roadmap.md](../../docs/roadmap.md); analysed in [ADR-002](../../docs/architecture/adr/ADR-002-lazy-expiry-and-fifo-by-expiry-consumption.md) | Deferred | Sprint 3 capacity went to the customer-requested UI verification features; the tradeoff is now documented in ADR-002 and stays a Sprint 4 candidate. |
| (New, 3 Jul review) Autotester should support multiple idempotency keys in parallel requests. | [US-19 / #50](https://github.com/Varriwon4ik/avito_bonus_point_service/issues/50) | Added to backlog | Created during the meeting; refined with acceptance criteria; Sprint 4 candidate. |
| (New, 3 Jul review) Response-code presentation in the UI could be changed. | — (no PBI) | Rejected/deferred with rationale | The customer marked it "not critical" himself; noted for future UI polish, no PBI opened. |
| (New, 3 Jul review, process) "Work as a team, as a whole." | [Retrospective action point 1](retrospective.md#action-points) | Addressed in process | Next Sprint pairs a named co-worker with every implementer. |

**Feedback not addressed this Sprint:** the demo-version regression run and
the expiry-as-explicit-transaction model — both deferred with the rationale
above and tracked in [docs/roadmap.md](../../docs/roadmap.md). MVP v2 does
address the core of the customer's standing verification feedback (US-16 +
US-17), so no scope-substitution justification is needed.

## 6. Maintained documentation

- [docs/roadmap.md](../../docs/roadmap.md)
- [docs/definition-of-done.md](../../docs/definition-of-done.md) — updated for
  Assignment 5 (architecture-documentation completion criterion)
- [docs/testing.md](../../docs/testing.md)
- [docs/quality-requirements.md](../../docs/quality-requirements.md) — each QR
  now links its ADR(s)
- [docs/quality-requirement-tests.md](../../docs/quality-requirement-tests.md)
- [docs/user-acceptance-tests.md](../../docs/user-acceptance-tests.md)
- [docs/development-process.md](../../docs/development-process.md) — **new**:
  development process + configuration management, with the Mermaid `gitGraph`
- [docs/architecture/README.md](../../docs/architecture/README.md) — **new**:
  the maintained architecture documentation
- [docs/user-stories.md](../../docs/user-stories.md)

## 7. Architecture and ADRs

- **Architecture documentation:**
  [docs/architecture/README.md](../../docs/architecture/README.md)
- **View artifacts (diagrams-as-code sources):**
  - Static view: [component-diagram.mmd](../../docs/architecture/static-view/component-diagram.mmd)
  - Dynamic view: [two-phase-redemption-sequence.mmd](../../docs/architecture/dynamic-view/two-phase-redemption-sequence.mmd)
  - Deployment view: [deployment-diagram.mmd](../../docs/architecture/deployment-view/deployment-diagram.mmd)
- **ADR directory:** [docs/architecture/adr/](../../docs/architecture/adr/) —
  indexed in the
  [ADR index](../../docs/architecture/README.md#architecture-decision-records-adr-index)
  (ADR-001 row locking, ADR-002 lazy expiry/FIFO, ADR-003 layered monolith +
  coverage gates, ADR-004 idempotency keys, ADR-005 single binary + compose
  deployment — all `Accepted`).

**Architecture summary:** the product is a layered monolith in one Go binary —
web UI, Swagger UI, HTTP layer (`internal/api`), shared autotest engine
(`internal/autotest`), and persistence layer (`internal/data`) over PostgreSQL.
Every state change funnels through `internal/data`'s locked transactions, the
UI and autotester consume only the public `/v1` API, and deployment is docker
compose on a university VM. This shape gives the product its core guarantees:
no double-spend under concurrency, fast single-query balance reads, and
precisely gateable critical modules.

**How quality requirements link to architecture decisions:** each QR links its
ADR(s) and vice versa — QR-001 (balance p95 ≤ 200 ms) is achievable because
ADR-002 made expiry a query predicate instead of a background job; QR-002 (no
overspend) is implemented by ADR-001 (row locking) plus ADR-004 (idempotency
keys) and verified by QRT-002; QR-003 (≥30% critical-module coverage) is
enforceable because ADR-003/ADR-005 give the codebase exactly the module
boundaries the CI gate targets. See
[docs/quality-requirements.md](../../docs/quality-requirements.md).

## 8. Testing and CI status for the delivered increment

- All Assignment 4 gates remain active and passed for the Sprint 3 work:
  `gofmt`, `go vet`, build, unit + integration tests with the race detector
  against real Postgres, automated QRTs (QRT-001/002/003), the per-module
  ≥30% coverage gate, `govulncheck`, and Lychee link checking.
- Automated verification was extended where MVP v2 changed the product: label
  validation unit tests
  ([internal/data/labels_test.go](../../internal/data/labels_test.go)), and the
  shared autotest engine exercising accrual correctness and parallel requests
  end-to-end from both frontends.
- New UAT scenarios UAT-004/005/006 cover the new user-facing behaviour
  ([docs/user-acceptance-tests.md](../../docs/user-acceptance-tests.md)).
- Full status: [docs/testing.md](../../docs/testing.md).
- **CI pipeline:** [.github/workflows/ci.yml](../../.github/workflows/ci.yml)
- **Latest protected-default-branch CI run:**
  [Actions › CI › branch:main](https://github.com/Varriwon4ik/avito_bonus_point_service/actions/workflows/ci.yml?query=branch%3Amain)

## 9. Release and changelog

- **SemVer release mapped to MVP v2:**
  [v2.0.0](https://github.com/Varriwon4ik/avito_bonus_point_service/releases/tag/v2.0.0)
  — tagged on the protected default branch; the release description identifies
  the MVP v2 mapping and links the Sprint 3 milestone, run/access instructions,
  the public demo video, and this report.
- **Changelog:** [CHANGELOG.md](../../CHANGELOG.md) (the `2.0.0` section lists
  every user-visible MVP v2 change, issue-linked).

## 10. Demo video

- **Public sanitized demo video (< 2 min):**
  <https://drive.google.com/file/d/1BQiw7kSVUJNK1O_8oTDC2VBUupXqE8yK/view?usp=sharing>
  — explains the current state of MVP v2 (what was improved, fixed, and added),
  sanitized demo data only.

## 11. UAT results (public, sanitized)

Executed with the customer on **3 July 2026**: each implementer shared their
screen and the **customer directed the demonstration** of their feature.

- **UAT-004** See exact HTTP response codes in the web UI (US-16) — **Passed**
- **UAT-005** Run the autotester from the web UI (US-17) — **Passed**
- **UAT-006** Label a transaction and find it in the history (US-18 + US-08) — **Passed**

Nothing failed and no defect PBI was opened. UAT-001–003 were not formally
re-executed; their core flows were exercised throughout the demos and behaved
as previously accepted. **Still to fix in the product:** nothing from this
session; open items are enhancements — the multi-key autotester
([US-19 / #50](https://github.com/Varriwon4ik/avito_bonus_point_service/issues/50))
and the carried demo-version regression request. **Most important feedback:**
the multi-key autotester idea, the non-blocking response-code presentation
suggestion, and the repeated "work as a team" process feedback (→
[retrospective](retrospective.md)). Full scenarios and execution history:
[docs/user-acceptance-tests.md](../../docs/user-acceptance-tests.md).

## 12. Sprint Review artifacts

- **Transcript (public publication permitted by the customer):**
  [sprint-review-transcript.md](sprint-review-transcript.md) — one recorded
  session covers both the Sprint Review and the customer-directed UAT;
  recording permission was asked and granted before recording started.
- **Summary:** [sprint-review-summary.md](sprint-review-summary.md)
- The private recording link and the exact private timecodes (including where
  the customer-directed UAT segments occur) are submitted only through Moodle.

## 13. Hosted documentation site

- **Hosted docs:** <https://varriwon4ik.github.io/avito_bonus_point_service/>
  — the maintained `docs/` set (architecture with rendered diagrams,
  development process, quality/testing docs, roadmap, UATs) published with
  MkDocs Material to GitHub Pages on every merge to `main`
  ([workflow](../../.github/workflows/docs.yml)). Linked from the
  [root README](../../README.md#documentation) and the `v2.0.0` release.

## 14. Reflection, retrospective, LLM report

- [reflection.md](reflection.md)
- [retrospective.md](retrospective.md)
- [llm-report.md](llm-report.md)

## 15. Current product status

MVP v2 is delivered and released as `v2.0.0`: the ledger's core guarantees
(transactional integrity, idempotency, expiry, two-phase debits) are unchanged
and still gate-verified, and the product is now directly usable and verifiable
from the web UI — exact response codes, an autotester tab, transaction labels,
and validated TTL bounds. The architecture, its decisions, and the team's
process are documented as maintained assets and published as a browsable docs
site. The increment is deployed on the University VM; the customer accepted all
demonstrated features and is satisfied with the team's workflow.

## 16. Next steps

- Sprint 4: US-01 (bulk accrual), US-02 (lots audit API), US-19 (multi-key
  autotester), and the twice-carried demo-version regression task — see
  [docs/roadmap.md](../../docs/roadmap.md).
- Process: pair a co-worker with every implementer and deploy on the day of
  every merge ([retrospective action points](retrospective.md#action-points)).
- Keep the new maintained assets current under the updated
  [Definition of Done](../../docs/definition-of-done.md).

## 17. Contribution traceability

| Member | GitHub | Roles | Issues / PRs (implementer) | Review activity | Testing / quality / automation / architecture / docs |
|---|---|---|---|---|---|
| Mikhail Ilin | [@Varriwon4ik](https://github.com/Varriwon4ik) | PM; QA; Backend | US-17 ([#40](https://github.com/Varriwon4ik/avito_bonus_point_service/issues/40), [PR #49](https://github.com/Varriwon4ik/avito_bonus_point_service/pull/49)); revert [PR #47](https://github.com/Varriwon4ik/avito_bonus_point_service/pull/47) | Reviewed US-16 ([PR #46](https://github.com/Varriwon4ik/avito_bonus_point_service/pull/46)) and US-18 ([PR #44](https://github.com/Varriwon4ik/avito_bonus_point_service/pull/44)) | Sprint 3 milestone, architecture docs + ADRs, development-process doc, hosted docs site, Week 5 report, release v2.0.0, US-17 demo at review |
| Nurislam Denisov | [@NurikDen](https://github.com/NurikDen) | QA; Backend | US-16 ([#39](https://github.com/Varriwon4ik/avito_bonus_point_service/issues/39), [PR #46](https://github.com/Varriwon4ik/avito_bonus_point_service/pull/46)) | Reviewed US-17 ([PR #49](https://github.com/Varriwon4ik/avito_bonus_point_service/pull/49)) | Frontend HTTP-status display + ledger table work, US-16 demo at review |
| Sanzhar Kadambaev | [@kadambaevsanzhar](https://github.com/kadambaevsanzhar) | Scrum Master; QA | US-18 ([#41](https://github.com/Varriwon4ik/avito_bonus_point_service/issues/41), [PR #44](https://github.com/Varriwon4ik/avito_bonus_point_service/pull/44), re-landed via [PR #49](https://github.com/Varriwon4ik/avito_bonus_point_service/pull/49)) | Reviewed US-08 ([PR #42](https://github.com/Varriwon4ik/avito_bonus_point_service/pull/42)) | Label validation + tests (labels_test.go), migration 0004, US-18 demo at review |
| Sergey Chuenko | [@SergeiCh07](https://github.com/SergeiCh07) | Scrum Master; QA | US-08 ([#5](https://github.com/Varriwon4ik/avito_bonus_point_service/issues/5), [PR #42](https://github.com/Varriwon4ik/avito_bonus_point_service/pull/42)) | Reviewed US-18 revert ([PR #47](https://github.com/Varriwon4ik/avito_bonus_point_service/pull/47)) and approved [PR #46](https://github.com/Varriwon4ik/avito_bonus_point_service/pull/46) | TTL bounds env config + validation, US-08 demo at review |

## 18. Screenshots

The required Week 5 screenshot set (Sprint milestone, board view, latest
protected-branch CI run, `v2.0.0` release, example reviewed issue-linked PR,
hosted docs site, deployed product) is being captured and added to
[images/](images/README.md), which lists each planned file. **Deviation note:**
at the time this report was first published the screenshots were not yet
embedded — all listed evidence is directly inspectable at the live links
throughout this report (milestone, board, CI runs, release, PRs, hosted docs);
the PNG set is added to `images/` before submission so the evidence also
survives as repository-resident artifacts. Screenshots of the review/UAT video
session contain customer-identifying information and are submitted only through
Moodle.

## 19. Example reviewed, issue-linked PR

[PR #49 — web autotester + re-landed labels (US-17, US-18)](https://github.com/Varriwon4ik/avito_bonus_point_service/pull/49):
issue-linked (`Closes #40`), reviewed and approved by a different team member
(@NurikDen), merged into the protected default branch via a merge commit with
all CI gates green.

## 20. Deviations from expected defaults

- **Diagrams-as-code tool:** the assignment recommends PlantUML; the team uses
  **Mermaid** (explicitly allowed). Rationale: Mermaid renders natively both on
  GitHub and on the MkDocs Material site, so the maintained architecture
  diagrams are readable in context in *both* required places without a separate
  render pipeline or committed binary artifacts; the sources live in the
  required `docs/architecture/*-view/` directories.
- **Screenshots:** pending at first publication — see §18.
- **Undeployed demo at the review:** repository troubles at meeting time meant
  the increment was demonstrated from a local build with the customer's
  explicit consent; the VM deployment was updated after the meeting.
- **`v2.0.0` tag placement:** the release tag points at the protected-branch
  commit containing the complete MVP v2 *product* code; the Assignment 5
  documentation set (this report, architecture docs) merged to `main`
  immediately afterwards via the reviewed PR for
  [#51](https://github.com/Varriwon4ik/avito_bonus_point_service/issues/51).
