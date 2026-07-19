# Sprint 5 Review Summary (Assignment 6, Week 7)

- **Date:** 19 July 2026
- **Sprint:** Sprint 5 (13–19 July 2026)
- **Sprint Goal:** respond to Week 6 trial feedback, complete the transition,
  and deliver the final course version, MVP v3.
- **Increment reviewed:** the final Bonus Points Ledger Service, including the
  corrected embedded UI from #60 and the horizontal-scaling assessment and
  concurrent-migration hardening from #64 / PR #66.

## Review outcome

The team demonstrated the final product to the Customer. After the demo, the
Customer assessed the product as working well and complete. This closes the two
conditions raised during the Week 6 review: the missing UI controls were fixed
and redeployed, and the architecture now explicitly explains that the stateless
API tier can scale horizontally under documented conditions.

The reached handover level is **Ready for independent use**. The
customer-confirmation status is **Accepted** for the demonstrated final product
and the reached handover scope. Customer-side deployment was not claimed: the
Customer previously explained that their infrastructure is confidential and
that their own interns would evaluate and deploy the project after delivery.

## UAT and feedback

The final product demonstration provided customer-facing acceptance evidence
for the overall increment. No defect or requested product change was raised in
the final review. The supplied condensed transcript does not separate named
UAT-007/UAT-009 steps, so the team does not misrepresent them as individually
executed in this session; the automated regression and CI evidence remains
current, and the Customer accepted the demonstrated final behavior as complete.

## Remaining risks and limitations

- The public university-VM instance is an evaluation deployment on a private
  network, not a customer production deployment.
- Authentication and authorization remain outside the agreed internal-network
  product scope.
- PostgreSQL is the shared write bottleneck and single point of failure in the
  shipped topology; per-replica metrics must be scraped directly when scaling.
- The Customer's attendance at the 21 July presentation was not confirmed
  because university coordination was still pending.

The recording link and exact timecodes are kept in the private Week 7 Moodle
PDF. The sanitized public transcript is
[sprint-review-transcript.md](sprint-review-transcript.md).
