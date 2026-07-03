# Sprint 2 Retrospective (Assignment 4)

- **Date:** 28 June 2026 (after the Sprint Review)
- **Sprint:** Sprint 2 / Assignment 4 (22–28 June 2026)

## What went well

- **CI was a big win.** Gating every push and pull request with build, tests, and
  the race detector saved a lot of time and space during code review — reviewers
  no longer re-check formatting, build, and tests by hand. The team considers it
  one of the most valuable additions to the workflow so far.
- **The autotester paid off.** Being able to define, store, and replay accrual and
  concurrency scenarios against a running instance — including parallel requests —
  turned out to be an effective way to validate newly committed code without
  hand-writing a bespoke test for each case.
- **The increment was accepted.** The customer inspected every file changed during
  the Sprint and was satisfied with the version and the team's design decisions.
- **Measurable quality requirements** gave us concrete, automatable targets
  (latency, integrity, coverage) instead of vague quality goals.

## What did not go well

- **US-07 had to be reverted.** The manual-accrual admin auth was merged and then
  reverted because of bugs and integration issues found afterwards. We spent
  effort on a change that did not ship and had to be unwound and reclassified
  `Removed`.
- **Regression coverage of older code paths** is still thinner than the customer
  wants; we addressed the request partly (autotester + QRTs) but not fully.
- **Deployment access friction** — the private VM requires the university
  network/VPN, which complicates remote access for reviewers.

## What we changed compared to the previous Sprint

- Acting on the previous retrospective's pull toward more automation, we **stood
  up a real CI pipeline** and made it a required merge gate — the previous Sprint
  relied on local, manual verification.
- We **added measurable quality requirements and automated QRTs**, plus a
  per-module coverage gate and a dependency vulnerability scan, so quality is now
  enforced by the pipeline rather than by reviewer diligence alone.

## Process improvements for the next Sprint

1. **Test before merge, not after.** Treat a change as un-mergeable until its
   tests (including the new QRTs) are green locally and in CI — the US-07 revert
   showed the cost of merging first and finding issues later.
2. **Plan dedicated regression tasks for legacy code paths** in Sprint 3 so the
   customer's "prove the changes are valid" request is fully closed, with a
   linked PBI and acceptance criteria.
