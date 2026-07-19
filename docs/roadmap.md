# Product Roadmap

Bonus Points Ledger Service — a REST-like service that manages an online
store's bonus-points program (configurable expiry, transactional balance
mutations, two-phase debits, idempotency, observability).

This roadmap gives direction across releases. It is intentionally lighter than
the Product Backlog: see [docs/user-stories.md](user-stories.md) and the GitHub
issues for the authoritative, current state. The course ends with the Week 7
delivery of **MVP v3**; this roadmap covers the remaining course work and the
state reached at course end — it deliberately does not plan speculative
post-course versions.

## Product Goal

Provide a reliable, auditable, and operationally observable bonus-points ledger
that downstream services can integrate with safely under concurrency, without
losing or double-counting points — and hand it over so the customer can use,
verify, and operate it independently.

## Sprint 1 (MVP v1) — 15–21 Jun 2026 — delivered

[Sprint 1 milestone](https://github.com/Varriwon4ik/avito_bonus_point_service/milestone/1)

**Sprint Goal:** complete the initial core features requested by the customer.

Delivered increment (all `Must Have`, marked MVP v1):

- **US-05** — auto-release stale holds after a configurable timeout.
- **US-11** — verified idempotent-key deduplication under concurrent requests.
- **US-12** — points removal / expiry system.
- **US-13** — correct HTTP response codes and OpenAPI documentation.

Also delivered in Sprint 1 (`Should Have`):

- **US-10** — structured request logging and a Prometheus `/metrics` endpoint.

## Sprint 2 (Assignment 4) — 22–28 Jun 2026 — delivered

[Sprint 2 milestone](https://github.com/Varriwon4ik/avito_bonus_point_service/milestone/2)

**Sprint Goal:** strengthen the reliability and verifiability of the increment by
delivering an automated autotester and a CI pipeline that gate every change,
while demonstrating paginated transaction-history access to the customer.

Delivered increment:

- **US-14** (`Must Have`) — continuous integration pipeline that builds and tests
  every push and pull request to `main`.
- **US-15** (`Must Have`) — `cmd/autotest` autotester that defines, stores, and
  replays reusable accrual / concurrency scenarios against a running instance.
- **US-09** (`Should Have`) — pagination for transaction history.

Reverted this Sprint: **US-07** (manual accrual) was implemented then reverted
due to bugs and is now `Removed`.

## Sprint 3 (MVP v2, Assignment 5) — 29 Jun–5 Jul 2026 — delivered

[Sprint 3 milestone](https://github.com/Varriwon4ik/avito_bonus_point_service/milestone/3)

**Sprint Goal:** deliver MVP v2 by making the ledger directly usable and
verifiable from the web UI — exact HTTP response feedback, a browser-based
autotester, transaction labels, and enforced per-accrual TTL bounds — while
documenting the architecture, key decisions (ADRs), and the development process
so the product can keep evolving safely.

Delivered increment (16 SP, mapped to release `v2.0.0`):

- **US-08** (`Should Have`) — configurable per-accrual TTL validation and bounds
  (`MIN_TTL_DAYS` / `MAX_TTL_DAYS`, out-of-range requests get `400`).
- **US-16** (`Should Have`) — exact HTTP response codes surfaced in the web UI
  for accrual/debit operations.
- **US-17** (`Should Have`) — web Autotester tab backed by
  `POST /v1/autotest/run` and the shared `internal/autotest` engine.
- **US-18** (`Could Have`) — labels on transactions (preset `test`/`real` or
  custom), shown in the transactions view. First landing was reverted and the
  feature was re-landed fixed within the Sprint.

Also delivered: maintained architecture documentation with three views
([docs/architecture/README.md](architecture/README.md)), five ADRs,
[docs/development-process.md](development-process.md), and the hosted
documentation site.

**Scope note:** the earlier candidate plan for MVP v2 listed US-01 and US-02.
During Sprint 3 planning the team re-prioritized toward the web-UI feature set
above, because the customer's Sprint 2 feedback centred on seeing and verifying
behaviour directly (responses, autotests, labels); US-01/US-02 moved to
Sprint 4, as agreed with the customer at the Sprint 3 review.

## Sprint 4 (Week 6 trial release, Assignment 6) — 6–12 Jul 2026 — delivered

[Sprint 4 milestone](https://github.com/Varriwon4ik/avito_bonus_point_service/milestone/4)

**Sprint Goal:** deliver the Week 6 trial / handover-candidate release that
completes the customer's remaining `Should Have` scope — bulk accrual for
promotional campaigns, the lots audit API for support tooling, and the
multi-key parallel autotester requested at the Sprint 3 review — fix the
defects found while trialling the web UI, and put the customer in a position to
try the product independently and judge transition readiness against the
reviewed customer-facing documentation set.

Delivered increment (24 SP, mapped to the trial release `v2.1.0`):

- **US-01** (`Should Have`, 8 SP) — bulk points accrual for promotional
  campaigns: `POST /v1/accruals/batch` with per-item results (HTTP 207) plus a
  Bulk accrual card in the web UI
  ([#1](https://github.com/Varriwon4ik/avito_bonus_point_service/issues/1)).
- **US-02** (`Should Have`, 5 SP) — paginated lots audit API with
  `status` filtering for support tooling
  ([#2](https://github.com/Varriwon4ik/avito_bonus_point_service/issues/2)).
- **US-19** (`Could Have`, 3 SP) — autotester scenarios with multiple distinct
  idempotency keys in parallel, selectable from the web Autotester tab
  ([#50](https://github.com/Varriwon4ik/avito_bonus_point_service/issues/50)).
- **Bug fixes** (8 SP) — web autotester false "found issues" verdict
  ([#54](https://github.com/Varriwon4ik/avito_bonus_point_service/issues/54));
  Test-mode selector and Bulk accrual card missing from the served UI
  ([#60](https://github.com/Varriwon4ik/avito_bonus_point_service/issues/60)).

Also delivered for Assignment 6 Week 6 (Course Task
[#62](https://github.com/Varriwon4ik/avito_bonus_point_service/issues/62)): the
maintained [customer handover documentation](customer-handover.md), the
contributor and agent guidance (`CONTRIBUTING.md`, `AGENTS.md`), the polished
repository entry point, and the Week 6 report set.

## Sprint 5 (MVP v3, Assignment 6) — 13–19 Jul 2026 — delivered

[Sprint 5 milestone](https://github.com/Varriwon4ik/avito_bonus_point_service/milestone/5)

**Sprint Goal:** respond to the Week 6 customer-trial and documentation-review
feedback, complete the actual transition of the product to the customer, and
deliver the final course version **MVP v3** as a release with higher SemVer
precedence than the Week 6 trial release, leaving the customer able to use,
verify, and operate the ledger with the maintained documentation set.

**Delivered from the Week 6 customer trial and documentation review
(10 Jul 2026):**

- **[#64](https://github.com/Varriwon4ik/avito_bonus_point_service/issues/64)**
  (`Must Have`) — assess horizontal-scaling suitability and state it
  explicitly in the architecture documentation. One of the customer's two
  conditions for considering the delivery complete; the other — the UI
  display fix ([#60](https://github.com/Varriwon4ik/avito_bonus_point_service/issues/60))
  — was resolved on the day of the review.
- Re-execute UAT-007 and UAT-009 with the customer at the Week 7 transition
  confirmation (they failed at the Week 6 session on the stale-UI defect
  fixed the same day).

**Remaining planned candidates:**

- Run the autotester/regression checks against a demo (earlier) version of the
  product to prove changes are valid (customer request carried from
  Sprint 2/3).
- Model point expiry as an explicit ledger transaction rather than silent
  removal, so expirations are auditable (follow-up on US-12; see
  [ADR-002](architecture/adr/ADR-002-lazy-expiry-and-fifo-by-expiry-consumption.md)).
- Evaluate deploy-on-merge automation for the university VM
  ([ADR-005](architecture/adr/ADR-005-single-binary-web-ui-and-compose-deployment.md);
  retrospective action point).
- Final transition work: confirm the handover level and customer acceptance of
  [docs/customer-handover.md](customer-handover.md), release MVP v3, record the
  public sanitized demo video, and prepare Demo Day.

## State at course end (reached)

By the end of Week 7 the product reached **MVP v3** and the `v3.0.0` release
candidate: all `Must Have` and `Should Have` stories delivered
(US-01, US-02, US-05, US-08–US-19 minus removed IDs), the maintained
documentation set (architecture, process, quality, testing, UAT, handover)
current, all CI quality gates active, and final status **Ready for independent
use / Accepted** stated in
[docs/customer-handover.md](customer-handover.md) and the Week 7 report.
Post-course version planning is intentionally out of scope.

## Out of scope (current)

- US-07 — `Removed` (implemented then reverted in Sprint 2 due to bugs).
- US-03, US-04, US-06 — removed; already satisfied by MVP v0 base functionality.
- Authentication / authorization for admin-facing operations — the reverted
  US-07 attempt is not being re-introduced within the course timeframe; the
  service stays internal-network-only per the original spec.
- Speculative post-course feature versions.

## Release mapping

| Release | Scope | Status |
|---|---|---|
| [`v1.0.0`](https://github.com/Varriwon4ik/avito_bonus_point_service/releases/tag/v1.0.0) | MVP v1 (US-05, US-10, US-11, US-12, US-13) | Sprint 1 — released |
| [`v1.1.0`](https://github.com/Varriwon4ik/avito_bonus_point_service/releases/tag/v1.1.0) | Sprint 2 increment (US-09, US-14, US-15) | Sprint 2 — released |
| [`v2.0.0`](https://github.com/Varriwon4ik/avito_bonus_point_service/releases/tag/v2.0.0) | MVP v2 (US-08, US-16, US-17, US-18 + architecture/process docs) | Sprint 3 — released |
| [`v2.1.0`](https://github.com/Varriwon4ik/avito_bonus_point_service/releases/tag/v2.1.0) | Week 6 trial / handover candidate (US-01, US-02, US-19 + fixes #54/#60) | Sprint 4 — released |
| `v3.0.0` (release candidate) | MVP v3 — final course version after Week 6 trial feedback and transition | Sprint 5 — delivered; GitHub release pending publication from final `main` commit |
