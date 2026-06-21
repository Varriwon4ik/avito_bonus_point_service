# Week 3 Reflection

## Learning points

- **Backlog migration is not a copy-paste.** Moving Assignment 2 stories into
  issues forced us to re-examine each one. Several stories (US-03, US-04, US-06)
  turned out to be already satisfied by the MVP v0 base, so migration became an
  act of refinement rather than transcription.
- **MoSCoW is a living decision.** Stories we had marked `Must Have` in Week 2
  (US-03/04/06) were removed and others (US-07, US-09) were reprioritized to
  `Won't Have` after negotiation, which sharpened the MVP v1 scope down to four
  core `Must Have` stories.
- **Decomposition reaches the qualifying-PBI bar.** Eight user stories were not
  enough on their own; decomposing the delivered work into supporting technical
  PBIs (migrations, metrics wiring, CI, deployment, test harness) both reached
  the required backlog size and made the work assignable without clarification.
- **Estimation with `effort` labels + Story Points** gave us a shared vocabulary
  for relative sizing and Sprint capacity.
- **Workflow enforcement matters.** Issue-linked branches, an extended PR
  template, acceptance-criteria verification before merge, and merge-commit PRs
  made the increment auditable and review-friendly.

## Validated assumptions

- **Confirmed:** Database-level `SELECT ... FOR UPDATE` serialization prevents
  lost updates under concurrency — validated by the US-11 concurrent tests.
- **Confirmed:** Two-phase holds plus a timeout sweep (US-05) close the
  "crashed caller locks points forever" gap from MVP v0.
- **Confirmed:** Structured logging + `/metrics` (US-10) is sufficient for basic
  operability without logging sensitive payloads.
- **Rejected / revised:** Our assumption that silently expiring points was
  acceptable — the customer asked for expiry to be an explicit, auditable ledger
  transaction instead.

## Friction and gaps

- US-05 could not be demonstrated live because its owner (N. Nuriev) could not
  connect; verification fell back to repository/PR evidence.
- US-05 implementation had to be completed by Mikhail after the original assignee
  hit technical issues — a single-owner risk.
- The customer's US-11 and US-12 follow-ups are real scope we have not yet built.
- Story Points were only partially tracked as `effort:` labels before this Sprint;
  we standardized them mid-assignment.

## Planned response

- Add Sprint 2 backlog items for the two customer follow-ups (concurrency
  regression on legacy paths; expiry-as-transaction). See
  [docs/roadmap.md](../../docs/roadmap.md).
- Demonstrate US-05 ([#3](https://github.com/Varriwon4ik/avito_bonus_point_service/issues/3))
  live at the next review.
- Spread ownership to avoid single-owner blockers; pair on critical PBIs.
- Pull refined `Should Have` stories (US-01, US-02, US-08) into Sprint 2 planning.
