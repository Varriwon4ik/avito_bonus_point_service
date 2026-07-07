# Sprint 3 Retrospective (Assignment 5)

- **Date:** 5 July 2026 (after the Sprint Review of 3 July)
- **Sprint:** Sprint 3 / Assignment 5 (29 June – 5 July 2026)

## What went well

- **The whole planned scope shipped.** All four Sprint stories (US-08, US-16,
  US-17, US-18 — 16 SP) were delivered, demonstrated to the customer, and
  accepted; the three new UAT scenarios all passed with no defect PBI.
- **The customer-directed screen-share UAT format worked.** Letting the
  customer steer each demonstration surfaced concrete, actionable feedback
  (the multi-key autotester request became US-19 during the meeting) instead
  of a passive walkthrough.
- **The revert muscle works.** When US-18's first landing (PR #44) showed
  problems, the team reverted it cleanly (PR #47) and re-landed it fixed
  (PR #49) within the same Sprint — `main` stayed releasable throughout,
  exactly what the CI-gated workflow is for.
- **Extracting the shared autotest engine paid off immediately:** the web tab
  (US-17) and the console tool run identical checks with no duplicated logic.

## What did not go well

- **A feature had to be reverted again.** US-18 repeated the Sprint 2 US-07
  pattern (merge → discover issues → revert), even though this time it was
  successfully re-landed. Our "test before merge" resolution from the previous
  retrospective is not fully sticking.
- **The deployed instance lagged `main` at review time.** Repository troubles
  before the meeting meant the customer had to be shown an undeployed build;
  manual VM deployment is a recurring friction point.
- **Solo work pockets.** The customer explicitly told us — twice — to "work as
  a team, as a whole." Some features were effectively built by one person with
  others only watching the outcome, because routing knowledge mid-Sprint felt
  too slow.
- **The demo-version regression request is still open** (running our tests
  against an earlier product version); it has now been carried across two
  Sprints.

## What the team changed based on the previous Sprint Retrospective, and the results

- **"Test before merge, not after":** partially adopted. CI gates every PR and
  caught nothing broken on `main`, but US-18's issues were still found after
  merge rather than in review — the improvement needs a stronger review
  checklist, not just CI.
- **"Plan dedicated regression tasks for legacy code paths":** not done as a
  dedicated PBI; the Sprint's capacity went to the four UI-facing stories.
  Partially compensated by US-17 making the existing autotests runnable by
  anyone (including the customer), but the explicit demo-version regression
  task remains open and is carried into Sprint 4 planning.

## Action points

1. **Pair on every story next Sprint.** Each of US-01 and US-02 gets a named
   implementer *and* a named co-worker (not just a reviewer) who participates
   in design and testing before the PR opens — directly addressing the
   customer's team-cohesion feedback and the after-merge-discovery problem.
2. **Deploy immediately after every merge to `main`** (a named team member
   runs the VM update the same day), so a review never again has to fall back
   to an undeployed build; evaluate automating this in Sprint 4.
