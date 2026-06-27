# Customer Review Summary — Sprint 2 (Assignment 4)

- **Date:** 28 June 2026
- **Sprint:** Sprint 2 / Assignment 4 (22–28 June 2026)
- **Milestone:** [Sprint 2](https://github.com/Varriwon4ik/avito_bonus_point_service/milestone/2)
- **Participants (roles):** Mikhail Ilin (Project Manager; QA Engineer; Backend
  Developer), Nurislam Denisov (QA Engineer; Backend Developer), Sanzhar
  Kadambaev (Scrum Master; QA Engineer), Sergey Chuenko (Scrum Master; QA
  Engineer), and the Customer.
- **Recording:** A single recorded session covers both the Sprint Review and the
  customer-executed UAT. Recording permission and transcript-publication
  permission were granted by the Customer at the start of the meeting. The
  private recording link and exact recording timecodes are submitted through
  Moodle; the customer-executed UAT segment begins at the timecode noted there.

## Sprint Goal reviewed

Strengthen the reliability and verifiability of the increment by delivering an
automated autotester and a CI pipeline that gate every change, while
demonstrating paginated transaction-history access to the customer.

## Delivered increment discussed

| PBI | Item | Outcome |
|---|---|---|
| US-14 | Continuous integration pipeline for every change | Accepted |
| US-15 | Automated autotester for points accrual / concurrency | Accepted |
| US-09 | Pagination for transaction history | Accepted |
| US-07 | Manual accrual admin auth | Reverted (bugs); reclassified `Removed` |

The Sprint Review opened with a file-by-file walkthrough of everything committed
during the Sprint. The Customer inspected the changes and the team's design
decisions directly.

## UAT results

The Customer executed the acceptance scenarios via a guided code-and-behaviour
walkthrough (the third of three UAT approaches the team prepared). All three
active scenarios passed:

- **UAT-001** Read a user's points balance — Passed.
- **UAT-002** Two-phase redemption (hold → confirm / cancel) — Passed.
- **UAT-003** Review paginated transaction history — Passed.

Details and execution history: [docs/user-acceptance-tests.md](../../docs/user-acceptance-tests.md).

## Quality evidence discussed

The automated quality gates were shown: the QRTs (balance-read latency, ledger
integrity under concurrency), the per-module coverage gate (≥30%), and the
`govulncheck` scan, all running in CI on every change. See
[docs/quality-requirements.md](../../docs/quality-requirements.md) and
[docs/testing.md](../../docs/testing.md).

## Customer feedback

- The Customer reiterated the request (from the Sprint 1 review) for **additional
  automated tests run against the earlier version of the product** to objectively
  prove that the team's changes are valid. The US-15 autotester and the new
  automated QRTs/CI gates address this; broader regression coverage of legacy
  code paths is queued for Sprint 3.

## Approvals and decisions

- The Customer was satisfied with the current version and approved the team's
  design decisions. No demonstrated item was rejected; no defect was found during
  UAT.

## Risks and action points

- **Risk:** the deployment is on a private University VM; graders need
  university-network/VPN access (private access details via Moodle).
- **Action:** extend automated regression coverage over older code paths in
  Sprint 3 (continuation of the customer's validity-of-changes request).

## Resulting Product Backlog / scope changes

- No new defect PBIs from this session.
- Sprint 3 candidates reaffirmed: US-01, US-02, US-08 plus the Sprint 1
  follow-ups; see [docs/roadmap.md](../../docs/roadmap.md).
