# LLM Usage Report — Assignment 5

The team used AI/LLM tools during Assignment 5 as follows.

## Coding

- **Feature work (US-08, US-16, US-17, US-18):** team members used LLM
  assistants (Claude, ChatGPT-class tools) as coding aides for the web UI
  JavaScript, Go handler changes, and test scaffolding — generating first
  drafts and debugging suggestions. All generated code was read, adapted,
  tested locally, and passed through the normal reviewed-PR workflow with the
  full CI gate set (build, race-enabled tests, QRTs, coverage gate,
  govulncheck) before merge. No code was merged unreviewed.

## Documentation and reporting

- **Claude Code (Anthropic, Claude agent in the terminal/IDE)** was used
  extensively for the Assignment 5 documentation set: drafting
  `docs/architecture/README.md` and the Mermaid view diagrams, the five ADRs,
  `docs/development-process.md` (including the `gitGraph`), the MkDocs/GitHub
  Pages setup, updates to the maintained docs (roadmap, user stories, quality
  requirements, testing, Definition of Done, UATs, changelog, README), and the
  assembly of this Week 5 report set.
- The agent worked from the repository's actual state (source code, PR/issue
  history, CI configuration) plus facts supplied by the team (Sprint dates,
  meeting outcomes, the raw Sprint Review transcript, links). The team reviews
  the generated documentation for factual accuracy against the repository
  before merge; the architecture decisions described in the ADRs are the
  team's own, made during Sprints 1–3 — the LLM documented them, it did not
  make them.

## Transcription and translation

- The Sprint Review recording was transcribed/condensed by the team and
  translated into English with LLM assistance; the sanitized result was
  reviewed by the meeting participants before publication in
  [sprint-review-transcript.md](sprint-review-transcript.md).

## Research and ideas

- LLMs were consulted for tooling comparisons (diagrams-as-code options that
  render on both GitHub and MkDocs; hosted-docs approaches) before the team
  chose Mermaid + MkDocs Material.

## What LLMs were not used for

- Sprint planning decisions, estimation, prioritization, customer
  communication, and the Sprint Review/UAT itself were done by the team and
  the customer without AI involvement.
