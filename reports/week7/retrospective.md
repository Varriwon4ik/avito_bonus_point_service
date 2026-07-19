# Sprint 5 Retrospective (Assignment 6, Week 7)

- **Date:** 19 July 2026, after the final customer review
- **Sprint:** Sprint 5 (13–19 July 2026)

## What went well

- Both completion conditions from Week 6 were resolved before final delivery:
  the embedded UI was rebuilt and redeployed, and #64 added an explicit,
  customer-readable scaling assessment.
- The scaling review found one concrete multi-replica startup risk and removed
  it with a PostgreSQL advisory lock rather than limiting the work to prose.
- The final customer demo ended with a clear assessment that the product was
  complete and working well; no additional product changes were requested.
- The final scope stayed focused on maintenance, transition, documentation, and
  release readiness instead of adding speculative post-course features.

## What did not go well

- Sprint 5 platform tracking was incomplete: #64 did not retain a visible Story
  Point estimate, assignee, or recorded reviewer even though the implementation
  and merge are inspectable. This weakens process evidence despite a sound
  technical result.
- The final meeting's condensed transcript does not identify UAT-007 and
  UAT-009 as separate executions. The overall demo and acceptance are recorded,
  but scenario-level evidence is weaker than planned.
- The final release, screenshots, and report publication necessarily occur at
  the end of the Sprint, leaving little buffer for link verification.

## Previous action points and results

- **Feature-freeze buffer:** achieved for product features. Week 7 concentrated
  on the already-requested scaling work and final evidence rather than new
  feature scope.
- **Documentation-driven deployment smoke test:** partially evidenced. The
  corrected UI was demonstrated successfully to the Customer, but the condensed
  record does not preserve a separate operator checklist or named verifier.

## Action points

1. Before submission, a team member other than the author must verify every
   Week 7 public link, the deployed UI, Swagger, `/healthz`, and the public demo
   video from a clean browser session; record failures as issues.
2. For future projects, create every Sprint PBI with estimate, assignee, and
   reviewer before implementation starts, and record named UAT scenario results
   during the meeting rather than reconstructing them afterward.
