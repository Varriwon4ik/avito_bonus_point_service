# Week 4 Reflection — Sprint 2 (Assignment 4)

## Learning points

- **Quality work is deliverable work.** Spending a Sprint on CI, an autotester,
  and pagination — rather than a pile of new features — produced a more reliable,
  verifiable increment. Responding to the customer's "prove the changes are valid"
  feedback was best answered with automation, not more manual demos.
- **Quality requirements are easier to defend when measurable.** Writing QR-001
  (p95 balance latency ≤ 200 ms), QR-002 (no overspend under concurrency), and
  QR-003 (≥30% critical-module coverage) forced us to turn vague goals into
  thresholds a test can check.
- **CI is leverage.** Once the pipeline gated every push, code review stopped
  having to re-verify formatting, build, and tests by hand. The team felt the
  time savings immediately.
- **Running UAT as a code walkthrough worked for this customer.** Because the
  product is a backend service, the customer preferred inspecting the actual
  implementation and decisions over a scripted click-through, and that built
  confidence in the design.
- **Public/private separation needs deliberate handling.** The private VM
  address, the recording link, university emails, and exact timecodes all had to
  be kept out of the public repo and routed through Moodle.

## Validated assumptions

- **Validated:** Database-level `SELECT ... FOR UPDATE` serialization prevents
  overspend under concurrency — confirmed by QRT-002 under the race detector.
- **Validated:** The balance endpoint is fast enough for the hot path — QRT-001
  keeps p95 within budget.
- **Validated:** The existing integration suite already exercises the critical
  modules well above the 30% floor, so the coverage gate is a safety net rather
  than a blocker.
- **Rejected:** Adding admin authentication (US-07) this Sprint was premature —
  the first implementation had bugs and integration issues, so we reverted it and
  reclassified US-07 as `Removed`, to be revisited with a hardened design.

## Friction and gaps

- **US-07 revert cost time** that could have gone to the planned scope; we merged
  then reverted rather than catching the issues before merge.
- **Regression coverage of older code paths** is still thinner than the customer
  would like — the autotester helps, but legacy paths need dedicated tests.
- **Coverage numbers are not yet pinned in the docs** — `docs/testing.md` records
  the enforced ≥30% floor; exact per-module percentages come from the first green
  CI run and should be filled in.
- **Deployment reachability** depends on the university network/VPN, which adds
  friction for remote graders.
- **`govulncheck` is pinned to a version**; new advisories or a stdlib finding
  could fail CI until the toolchain is bumped — expected, but needs monitoring.

## Planned response

- Extend automated regression coverage over older code paths in **Sprint 3**
  (continuation of the customer's validity-of-changes request) — see
  [docs/roadmap.md](../../docs/roadmap.md).
- Revisit admin auth (former [US-07](https://github.com/Varriwon4ik/avito_bonus_point_service/issues/4))
  with a fully tested design before re-introducing it.
- Pull the planned Sprint 3 stories
  [US-01](https://github.com/Varriwon4ik/avito_bonus_point_service/issues/1),
  [US-02](https://github.com/Varriwon4ik/avito_bonus_point_service/issues/2),
  [US-08](https://github.com/Varriwon4ik/avito_bonus_point_service/issues/5).
- Keep the Assignment 4 gates (QRTs, coverage gate, `govulncheck`) green and
  update [docs/testing.md](../../docs/testing.md) with exact coverage figures
  after the first green CI run on `main`.
