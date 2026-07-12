# LLM Usage Report — Assignment 6, Week 6

The team used AI/LLM tools during Week 6 (Sprint 4) as follows.

## Coding

- **Feature and fix work (US-01, US-02, US-19, #54, #60):** team members used
  LLM assistants (Claude, ChatGPT-class tools) as coding aides for Go handler
  and autotest-engine changes, the web UI JavaScript (bulk-accrual card,
  test-mode selector, verdict rendering), and test scaffolding — first drafts
  and debugging suggestions. All generated code was read, adapted, tested
  locally, and passed through the normal reviewed-PR workflow with the full
  CI gate set (build, race-enabled tests against real Postgres, QRTs,
  per-module coverage gate, `govulncheck`) before merge. No code was merged
  unreviewed.

## Documentation, planning wiring, and reporting

- **Claude Code (Anthropic, Claude agent in the terminal/IDE)** was used
  extensively for the Assignment 6 Week 6 asset set: drafting the new
  maintained artifacts ([docs/customer-handover.md](../../docs/customer-handover.md),
  [CONTRIBUTING.md](../../CONTRIBUTING.md), [AGENTS.md](../../AGENTS.md)),
  updating the maintained docs (roadmap with the Sprint 4/5 plan, user-story
  registry, UAT scenarios UAT-007/008/009, testing status, changelog cut for
  `v2.1.0`, README polish including the bulk-accrual and autotester-mode API
  documentation), preparing the Sprint 4 and Sprint 5 milestone
  descriptions and the `v2.1.0` release notes, and assembling this Week 6
  report set.
- The agent worked from the repository's actual state (source code, issue/PR
  history, CI configuration, merged Sprint 4 changes) plus facts supplied by
  the team (Sprint dates, roles, meeting logistics). The team reviews all
  generated documentation for factual accuracy against the repository before
  merge. Sections that depend on the Week 6 customer session (documentation
  review results, transition-readiness findings, UAT outcomes, Sprint Review
  transcript) are **not** LLM-generated — they are recorded by the team from
  the actual meeting.

## Transcription and translation

- If the customer permits publication, the Week 6 session transcript is
  transcribed/condensed by the team and translated into English with LLM
  assistance, then reviewed by the meeting participants before publication —
  the same process as Weeks 3–5.
  <!-- TODO(team): confirm or update this section after the session. -->

## Research and ideas

- LLMs were consulted on conventions for repository contributor/agent
  guidance files (`CONTRIBUTING.md`, `AGENTS.md`) and on handover-document
  structure before the team settled the final content.

## What LLMs were not used for

- Sprint planning decisions, estimation, prioritization, customer
  communication, the customer trial, and the Sprint Review/UAT session
  itself were done by the team and the customer without AI involvement.
