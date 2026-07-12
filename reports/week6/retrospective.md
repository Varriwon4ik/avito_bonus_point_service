# Sprint 4 Retrospective (Assignment 6, Week 6)

- **Date:** «after the Week 6 Sprint Review, July 2026»
  <!-- TODO(team): set the actual date; hold the retrospective after the
       Sprint Review and before Sprint 5 planning. -->
- **Sprint:** Sprint 4 / Assignment 6 Week 6 (6 – 12 July 2026)

## What went well

- **The whole planned scope shipped — 24 SP.** US-01, US-02, and US-19 plus
  both bug fixes were implemented, cross-reviewed, and merged into protected
  `main` with all CI gates green; the trial release `v2.1.0` was cut on time
  for the customer to try before the Week 7 transition.
- **We found our own bugs before the customer did.** Trialling the product
  the way a customer would surfaced
  [#54](https://github.com/Varriwon4ik/avito_bonus_point_service/issues/54)
  and [#60](https://github.com/Varriwon4ik/avito_bonus_point_service/issues/60)
  inside the Sprint — both were fixed and reviewed before the trial handover
  instead of during it.
- **The long-carried stories are finally done.** US-01 and US-02 had been
  planned since Assignment 3 and re-scheduled twice; delivering them closes
  the last `Should Have` items in the Product Backlog.
- **Review pairs spread across the whole team:** four distinct
  implementer/reviewer pairings this Sprint (see the
  [contribution table](README.md#14-contribution-traceability)) — no single
  person was the bottleneck reviewer.

## What did not go well

- **Integration debt hid until the end of the Sprint.** Backend features
  merged early (US-19 on 7 Jul), but the served UI didn't actually expose
  them until the #60 fix on 10 Jul — the "endpoint done ≠ feature usable"
  gap repeated the US-18 lesson from Sprint 3 in a new form.
- **Four of five PRs merged on one day (10 Jul).** The end-of-Sprint burst
  left no buffer for deep review findings; it worked out, but only because
  nothing serious was found.
- **Review load was uneven:** one member recorded no PR review this Sprint
  while others reviewed twice. The Week 5 "pair on every story" action point
  was only partially realized — reviewers were distinct, but named
  co-workers who participate *before* the PR opens mostly did not happen.
- **Deploy cadence still depends on remembering.** «TODO(team): record
  honestly whether the VM was updated on the day of each merge, per the
  Week 5 action point, or lagged again.»

## What the team changed based on the previous Sprint Retrospective, and the results

- **"Pair on every story (named co-worker, not just reviewer)":** partially
  adopted. Every PBI had a distinct reviewer and the pairs rotated across all
  four members, but pre-PR co-working was thin; the late discovery of #60
  is the kind of gap earlier co-working should catch. Carried forward in
  sharpened form (see action point 2).
- **"Deploy immediately after every merge to `main`":** «TODO(team): result —
  kept / partially kept / not kept, with the dates.» The structural fix
  (deploy automation) is a Sprint 5 planning candidate either way.

## Action points

1. **Feature-freeze buffer before the final release.** In Sprint 5, no
   feature PRs merge after Friday 17 Jul; the last two days are reserved for
   transition work, fixes, documentation, and release/deploy verification —
   so MVP v3 is not another Sunday burst.
2. **Documentation-driven smoke test of every deployment.** After each deploy
   (and before the Week 7 confirmation), one team member who did not deploy
   verifies the running instance following **only**
   [docs/customer-handover.md](../../docs/customer-handover.md) — if a step
   fails or is unclear, that is a documentation bug and gets an issue the
   same day. This tests the handover artifact itself, not just the product.
