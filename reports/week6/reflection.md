# Week 6 Reflection (Assignment 6, Sprint 4)

## Learning points

- **"Merged" is not "delivered" — the demo proved it the hard way.** All the
  Sprint 4 backend work was merged and green in CI, yet at the 10 July
  customer session the deployed page did not show the new controls: the web
  UI is compiled into the binary (`//go:embed`), and the deployment served a
  stale interface. Two of three UAT scenarios failed in front of the
  Customer for that single reason
  ([#60](https://github.com/Varriwon4ik/avito_bonus_point_service/issues/60)),
  even though the underlying endpoints worked when shown via Swagger and the
  console. The failure mode is now called out as a hard rule in
  [AGENTS.md](../../AGENTS.md) and [CONTRIBUTING.md](../../CONTRIBUTING.md)
  ("rebuild after any HTML change"), and the same-day turnaround —
  fix merged ([PR #61](https://github.com/Varriwon4ik/avito_bonus_point_service/pull/61))
  and VM redeployed on 10 July — limited the damage.
- **A customer reviews documentation with operational questions, not
  editorial ones.** The verdict on our doc set was "complete overall", and
  the one gap was nothing we would have found proofreading: an explicit
  statement whether the service scales horizontally. Documentation for
  handover has to answer the questions an operator asks before running the
  system, not just describe what exists —
  [#64](https://github.com/Varriwon4ik/avito_bonus_point_service/issues/64)
  now makes that statement a `Must Have` for final delivery.
- **Transition does not necessarily mean "customer runs our deployment".**
  The Customer will not adopt the university-VM instance — their
  large-company environment has a security and deployment bar our trial
  setup does not pretend to meet, and their internal structure is
  confidential. The realistic handover shape is: we deliver a polished,
  well-documented, self-hostable project; their own people (a group of
  interns) evaluate and deploy it after delivery. That reframes what
  [docs/customer-handover.md](../../docs/customer-handover.md) must be good
  at: enabling strangers to deploy without us in the room.
- **Writing the handover document surfaced undocumented knowledge.** Filling
  in the transferred/delegated/retained table, the recovery steps, and the
  limitation list forced concrete answers that previously existed only as
  team habit — and it was exactly the artifact the Customer could then
  review and judge.

## Validated assumptions

- **Per-item results were the right bulk-accrual contract:** the 207
  Multi-Status design (one bad row never fails the campaign, per-item
  idempotency keys survive partial retries) was accepted by the Customer at
  the session — the endpoint behaviour itself drew no objections even while
  the UI was broken.
- **The shared autotest engine keeps paying off:** when the web tab was
  unavailable at the demo, the same multi-key check could still be shown
  from the console tool — one engine, two front ends saved the
  demonstration.
- **Cross-review inside the team is functioning:** all five Sprint PBIs were
  implemented and approved by different members (four distinct
  implementer/reviewer pairs), which is what the customer's earlier "work as
  a team" feedback asked the process to show.

## Friction and gaps

- **Integration debt surfaced at the worst moment — again.** Sprint 3's
  review ran on an undeployed build; Sprint 4's review ran on a stale one.
  The pattern is the same: the last mile between `main` and the running
  instance is manual and unverified, so it fails precisely when a customer
  is watching. The gap between "endpoint merged" and "feature usable by a
  person" stayed invisible until the session because nobody re-walked the
  deployed UI after the last merges.
- **A same-day merge burst:** four of five PRs merged on 10 July, the day of
  the review itself — leaving no slack between merge, deploy, and demo.
- **Review load was uneven this Sprint:** one member implemented a large
  story but recorded no PR review (see the
  [contribution table](README.md#14-contribution-traceability)), and his
  absence at the session meant others covered his demo cold; pairing intent
  from the Week 5 retrospective was only partially realized.

## Planned response

- Both [retrospective action points](retrospective.md#action-points) target
  exactly this failure mode in Sprint 5: a feature-freeze buffer before the
  final release, and a documentation-driven smoke test of the deployed
  instance after every deploy — someone who did not deploy re-walks the UI
  using only [docs/customer-handover.md](../../docs/customer-handover.md).
- Deliver the Customer's two completion conditions early in Sprint 5: the UI
  fix is already done (#60), and the horizontal-scaling assessment
  ([#64](https://github.com/Varriwon4ik/avito_bonus_point_service/issues/64))
  is scheduled as `Must Have`; UAT-007/009 are re-executed with the Customer
  at the Week 7 confirmation.
- Keep [docs/customer-handover.md](../../docs/customer-handover.md) current
  through Week 7 — it now carries the transition state the assignment
  grades, so every access/deployment/limitation change updates it in the
  same PR.
