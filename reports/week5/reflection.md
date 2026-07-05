# Week 5 Reflection (Assignment 5)

## Learning points

- **Documenting the architecture changed how we talk about the product.**
  Writing the three views forced us to name things we had only held in our
  heads: that *every* mutation funnels through one locked-transaction
  choke-point, that the web UI is deliberately just another API client, and
  that the balance path's speed is a direct consequence of the lazy-expiry
  decision. The views made those claims inspectable instead of tribal.
- **ADRs are cheap now and expensive later.** Reconstructing the rationale for
  decisions made in Sprints 1–2 (row locking, lazy expiry, idempotency keys)
  took real archaeology through PRs and tests. Decisions we recorded the week
  they were made (the shared autotest engine, ADR-005) took minutes. We now
  understand why ADR history must be preserved rather than rewritten.
- **Linking QRs to ADRs closed a loop.** QR-002 ("no double-spend") used to be
  verified by a test but explained nowhere; now the requirement links to the
  decisions that implement it and the test that proves it, in both directions.
- **A documented process is a debuggable process.** Drawing our actual git
  workflow as a `gitGraph` — including the US-18 revert — made it obvious
  where the workflow works (revert + re-land kept `main` releasable) and where
  it leaks (issues found after merge, manual deploys lagging `main`).
- **Configuration management was mostly already right, but undocumented.**
  Env-vars-only runtime config, a sanitized `.env.example`, and a secretless CI
  existed since early Sprints; writing `docs/development-process.md` was the
  first time we could *verify* the whole story end to end.

## Validated assumptions

- **The customer values seeing behaviour over hearing about it.** The
  customer-directed screen-share UAT confirmed it: the most engaged feedback
  came while steering the demos (multi-key autotester idea, custom-label
  praise).
- **Keeping the UI on the public API was the right call:** every demo doubled
  as an API verification, and no UI feature needed private endpoints.
- **The maintained quality gates carry over without friction:** all Assignment 4
  gates (QRTs, coverage, govulncheck, Lychee) ran unchanged against Sprint 3's
  code and caught formatting/test issues on the feature branches, not on `main`.

## Friction and gaps

- **Manual deployment bit us at the worst moment** — the Sprint Review — and
  the customer had to see an undeployed build. The gap between "merged" and
  "deployed" is our biggest operational weakness (recorded in ADR-005 and the
  retrospective action points).
- **Merge-then-discover recurred** (US-18), showing that CI alone does not
  substitute for deeper pre-merge review of behaviour.
- **Team cohesion:** the customer's repeated "work as a team" feedback matches
  our own observation that knowledge silos formed around individual features.
- **Docs tooling had a learning curve:** making diagrams render identically on
  GitHub and the MkDocs site constrained tool choice (we picked Mermaid over
  PlantUML for exactly this reason) — a deviation from the assignment's
  recommendation that we justified in the Week 5 report.

## Planned response

- Pair implementers on the remaining stories (US-01, US-02) and deploy to the
  VM on the day of every merge — both are concrete retrospective action points.
- Keep the new maintained assets current under the updated Definition of Done:
  any structure/flow/deployment change must update the affected view and add or
  supersede an ADR in the same PR.
- Schedule the twice-carried demo-version regression task explicitly in
  Sprint 4 planning alongside US-19.
