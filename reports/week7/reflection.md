# Week 7 Reflection (Assignment 6, Sprint 5)

## Learning points

- A documentation request can reveal a real production risk. The horizontal-
  scaling analysis did more than add an architecture paragraph: auditing each
  single-instance assumption exposed a concurrent migration-startup race. A
  PostgreSQL advisory lock made the documented scaling claim true in practice.
- Final delivery benefits from reducing scope. The team treated Sprint 5 as a
  maintenance and transition Sprint, resolved the Customer's explicit
  conditions, and avoided adding features that would destabilize the handover.
- Acceptance claims must match the evidence. The Customer called the final
  product complete and working well, which supports final product acceptance.
  The team still distinguishes this from customer-side operation, which did not
  occur because deployment into the Customer's confidential infrastructure is
  their responsibility after delivery.

## Validated assumptions

- The service's database-enforced concurrency and idempotency rules support
  multiple stateless API replicas; duplicate hold sweepers are safe because
  every resolution is rechecked under a row lock.
- A self-hostable repository, maintained documentation site, and explicit
  operational limitations are an appropriate transition outcome where direct
  customer-side deployment is not available during the course.
- Fixing the embedded-UI deployment gap was sufficient to restore confidence:
  no new defect was raised during the final demonstration.

## Friction and gaps

- The GitHub Sprint 5 issue lacks complete estimate/assignment/review metadata,
  so the report must disclose the tracking gap rather than invent values.
- The condensed final-meeting transcript preserves the acceptance outcome but
  not a scenario-by-scenario UAT trace.
- The university-VM deployment is useful for evaluation but remains private-
  network infrastructure with weaker controls than the Customer's operational
  environment.

## What we would do differently

We would freeze and smoke-test the deployed UI before every customer session,
record the exact named UAT scenario at the moment it is performed, and make
estimate/assignee/reviewer fields mandatory before an issue enters a Sprint.
Those three practices would align technical completion, customer evidence, and
platform traceability without last-day reconstruction.
