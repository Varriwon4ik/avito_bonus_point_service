# LLM / AI Usage Report — Week 4 (Assignment 4)

AI/LLM tools were used during Assignment 4. This report discloses how, in line
with the course's AI-usage policy. All AI output was reviewed, edited, and
verified by the team before being committed; the analysis and decisions are the
team's own.

## Where AI was used

- **Enhancing the UAT scenarios.** An LLM was used to enhance and improve the
  team's proposed user acceptance test scenarios — refining the wording,
  structure, preconditions, and step/expected-outcome clarity of the candidate
  scenarios in [docs/user-acceptance-tests.md](../../docs/user-acceptance-tests.md)
  before they were executed with the customer.
- **Drafting Week 4 quality assets and report scaffolding.** An AI coding
  assistant helped draft and structure the new quality documentation
  (`docs/quality-requirements.md`, `docs/quality-requirement-tests.md`,
  `docs/testing.md`), the automated quality requirement tests
  (`internal/api/qrt_test.go`), the CI additions (per-module coverage gate and
  the `govulncheck` job), and the `reports/week4/` report files. The team
  reviewed every file, confirmed it matched the actual code and CI behaviour, and
  is responsible for the content.

## How AI output was validated

- QRT tests and CI changes are verified by the CI pipeline itself — they only
  count when they build and pass against the real Postgres service and the race
  detector on GitHub Actions.
- Documentation was checked against the actual implementation (handlers, store,
  workflow files) and against the customer-review and UAT records.
- The Lychee link check verifies that report and documentation links resolve.

## What AI was not used for

- It was not used to fabricate customer feedback, UAT results, or review content;
  those reflect the real recorded session on 28 June 2026.
- It did not make product or scope decisions — Sprint scope, the US-07 revert, and
  acceptance were team/customer decisions.
