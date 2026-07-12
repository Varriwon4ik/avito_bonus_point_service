# Sprint 4 Retrospective (Assignment 6, Week 6)

- **Date:** 12 July 2026 (after the Sprint Review of 10 July)
- **Sprint:** Sprint 4 / Assignment 6 Week 6 (6 – 12 July 2026)

The session focused on the root causes of the UI display bug that broke the
review demo, and on the approach for the horizontal-scaling assessment the
Customer requested; the team also used it to rehearse the presentation and
work on the presentation slides.

## What went well

- **The whole planned scope shipped — 24 SP.** US-01, US-02, and US-19 plus
  both bug fixes were implemented, cross-reviewed, and merged into protected
  `main` with all CI gates green; the trial release `v2.1.0` was cut for the
  customer trial.
- **Same-day recovery under pressure.** When the review demo hit the stale
  UI, the team turned it around within hours: issue
  [#60](https://github.com/Varriwon4ik/avito_bonus_point_service/issues/60)
  filed, fix reviewed and merged
  ([PR #61](https://github.com/Varriwon4ik/avito_bonus_point_service/pull/61)),
  VM redeployed — all on 10 July, meeting the Customer's "fix it by next
  week" ask a week early.
- **The shared autotest engine saved the demo:** with the web tab broken,
  the multi-key check (US-19) could still be shown passing from the console
  tool — one engine, two front ends.
- **The long-carried stories are finally done.** US-01 and US-02 had been
  planned since Assignment 3 and re-scheduled twice; delivering them closes
  the last `Should Have` items, and the Customer confirmed no parts are
  missing.
- **Review pairs spread across the whole team:** four distinct
  implementer/reviewer pairings this Sprint (see the
  [contribution table](README.md#14-contribution-traceability)).

## What did not go well

- **The demo broke in front of the Customer — second Sprint in a row.**
  Sprint 3 showed an undeployed build; Sprint 4 showed a deployment serving
  a stale embedded UI, failing UAT-007 and UAT-009 live. Root cause
  discussed at this retrospective: the web UI is compiled into the binary
  (`//go:embed`), so a deploy that does not rebuild the binary after HTML
  changes silently serves the old interface — and nobody re-walked the
  deployed UI between the final merges and the meeting.
- **Four of five PRs merged on 10 July, the review day itself.** The
  end-of-Sprint burst left zero slack between merge, deploy, and demo — the
  stale-UI failure is a direct consequence.
- **Review load was uneven:** one member recorded no PR review this Sprint,
  and his absence at the session meant teammates covered his demo cold. The
  Week 5 "pair on every story" action point was only partially realized —
  reviewers were distinct and rotated, but named co-workers participating
  *before* the PR opens mostly did not happen.

## What the team changed based on the previous Sprint Retrospective, and the results

- **"Pair on every story (named co-worker, not just reviewer)":** partially
  adopted. Every PBI had a distinct reviewer and the pairs rotated across
  the team, but pre-PR co-working was thin; a co-worker walking the deployed
  UI before the review is exactly what would have caught #60 early.
  Carried forward in sharpened form (action point 2).
- **"Deploy immediately after every merge to `main`":** formally kept — the
  VM was redeployed on 10 July, the same day the last Sprint 4 PRs merged —
  but it did not prevent the demo failure, because the meeting happened
  before the final merges and redeploy, and the interim deployment served a
  stale UI. Conclusion: deploy *timing* alone is not enough; every deploy
  needs a post-deploy walk-through of the served interface (action point 2),
  and customer-facing sessions must not race the Sprint's final merges
  (action point 1).

## Action points

1. **Feature-freeze buffer before the final release.** In Sprint 5, no
   feature PRs merge after Friday 17 July; the last two days are reserved
   for transition work, fixes, documentation, the horizontal-scaling
   write-up ([#64](https://github.com/Varriwon4ik/avito_bonus_point_service/issues/64)),
   and release/deploy verification — so the `MVP v3` delivery and the Week 7
   confirmation never race last-minute merges the way the 10 July review
   did.
2. **Documentation-driven smoke test of every deployment.** After each
   deploy (and before the Week 7 confirmation), one team member who did not
   deploy re-walks the running instance — including the web UI tabs —
   following **only** [docs/customer-handover.md](../../docs/customer-handover.md);
   a failed or unclear step is a documentation or deployment bug and gets an
   issue the same day. This would have caught the stale embedded UI before
   the Customer did, and it doubles as verification of the handover
   artifact itself.
