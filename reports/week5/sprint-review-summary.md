# Sprint Review Summary — Sprint 3 (Assignment 5)

- **Date:** 3 July 2026
- **Sprint:** Sprint 3 / Assignment 5 (29 June – 5 July 2026)
- **Milestone:** [Sprint 3](https://github.com/Varriwon4ik/avito_bonus_point_service/milestone/3)
- **Participants (roles):** Mikhail Ilin (Project Manager; QA Engineer; Backend
  Developer), Nurislam Denisov (QA Engineer; Backend Developer), Sanzhar
  Kadambaev (Scrum Master; QA Engineer), Sergey Chuenko (Scrum Master; QA
  Engineer), and the Customer.
- **Recording and permissions:** One recorded session covers both the Sprint
  Review and the customer-directed UAT. The Customer granted **recording
  permission** before recording started and **permitted public transcript
  publication** — see [sprint-review-transcript.md](sprint-review-transcript.md).
  The private recording link and exact private timecodes are submitted through
  Moodle only.
- **Format:** Each implementer shared their screen and demonstrated their
  feature while the **Customer directed the demonstration** (the team's chosen
  UAT execution mode for this Sprint).

## Sprint Goal reviewed

Deliver MVP v2 by making the ledger directly usable and verifiable from the web
UI — exact HTTP response feedback, a browser-based autotester, transaction
labels, and enforced per-accrual TTL bounds — while documenting the
architecture, key decisions (ADRs), and the development process so the product
can keep evolving safely.

## Delivered increment discussed

| PBI | Item | Demonstrated by | Outcome |
|---|---|---|---|
| US-16 ([#39](https://github.com/Varriwon4ik/avito_bonus_point_service/issues/39)) | Exact HTTP response codes in the web UI | Nurislam | Accepted; non-blocking presentation suggestion |
| US-17 ([#40](https://github.com/Varriwon4ik/avito_bonus_point_service/issues/40)) | Web autotester tab | Mikhail | Accepted; multi-key enhancement requested → US-19 |
| US-18 ([#41](https://github.com/Varriwon4ik/avito_bonus_point_service/issues/41)) | Labels on transactions (preset + custom) | Sanzhar | Accepted; custom labels specifically praised |
| US-08 ([#5](https://github.com/Varriwon4ik/avito_bonus_point_service/issues/5)) | Per-accrual TTL validation and bounds | Sergei | Accepted ("works great") |

Because of repository troubles at meeting time, the increment was shown from a
local (undeployed) build with the Customer's explicit agreement; the deployed
instance at `http://10.93.26.175:8080/` was updated afterwards.

## UAT results

Customer-directed execution of the new scenarios — all **passed**, no defect
PBI opened:

- **UAT-004** See exact HTTP response codes in the web UI (US-16) — Passed.
- **UAT-005** Run the autotester from the web UI (US-17) — Passed.
- **UAT-006** Label a transaction and find it in the history (US-18, US-08) — Passed.

UAT-001–UAT-003 were not formally re-executed; their core flows (accrual,
debit, history) were exercised throughout the demonstrations and behaved as
previously accepted. Full scenarios and execution history:
[docs/user-acceptance-tests.md](../../docs/user-acceptance-tests.md).

## Addressed customer feedback

- Sprint 2 feedback asked for more ways to verify the team's changes: the web
  autotester (US-17) makes the automated checks runnable and readable by a
  non-developer, and exact response codes (US-16) make behaviour inspectable in
  the UI.
- **Still open:** running the tests against a demo/earlier version of the
  product. The Customer remains satisfied with the team's workflow and expects
  the remaining stories (US-01, US-02) next week.

## Architecture evidence discussed

The team walked through the design decisions behind the demonstrated features
where the Customer probed them (autotester single- vs multi-key semantics,
label validation, TTL bounds). The Assignment 5 architecture documentation and
ADR set ([docs/architecture/README.md](../../docs/architecture/README.md),
five ADRs linked to QR-001..003) were completed alongside the Sprint and are
linked from the [Week 5 report](README.md); the customer-facing discussion
focused on the delivered behaviour.

## Feedback, risks, and action points

- **New feature request:** autotester scenario with multiple idempotency keys
  in parallel → created during the meeting as
  [US-19 / #50](https://github.com/Varriwon4ik/avito_bonus_point_service/issues/50).
- **Non-blocking suggestion:** reconsider how response codes are presented in
  the UI (US-16); the Customer marked it "not critical".
- **Process feedback:** "work as a team, as a whole" — the Customer repeated
  this twice; taken into the [retrospective](retrospective.md) action points.
- **Risk:** the deployed version lagged `main` at review time (manual deploys —
  see [ADR-005](../../docs/architecture/adr/ADR-005-single-binary-web-ui-and-compose-deployment.md));
  mitigated by deploying after the meeting and tracked as a process improvement
  candidate.

## Resulting Product Backlog changes

- **US-19** added to the Product Backlog
  ([#50](https://github.com/Varriwon4ik/avito_bonus_point_service/issues/50)).
- Sprint 4 candidates confirmed with the Customer: US-01, US-02, US-19, plus
  the demo-version regression follow-up — see
  [docs/roadmap.md](../../docs/roadmap.md).
