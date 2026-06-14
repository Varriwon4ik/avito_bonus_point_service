# Customer Meeting Summary

## Meeting details
- Participants: customer representative, product owner, backend developer, QA/developer
- Artifacts demonstrated: [user-stories.md](./user-stories.md), [api/openapi.yaml](../../api/openapi.yaml), [api/postman_collection.json](../../api/postman_collection.json), [mvp-v0-report.md](./mvp-v0-report.md), local MVP v0 at `http://localhost:8080`

## Permissions and publication status

This draft assumes that written customer consent to the public MIT-licensed development model had already been obtained before repository creation, and that the customer allowed recording, private instructor sharing, and publication of a sanitized English transcript in the repository. Because transcript publication is assumed to be permitted for this draft, [customer-meeting-transcript.md](./customer-meeting-transcript.md) is included and `customer-meeting-notes.md` is intentionally not created.

## Discussion points

- The team walked through the Week 2 user stories, the assigned MoSCoW priorities, and the initial proposed MVP v1 scope.
- The team demonstrated the existing MVP v0 foundation, including accrual, hold, confirm, cancel, balance, lots, and transaction-history capabilities.
- The customer reviewed the proposed API artifacts and asked for clearer documentation of HTTP status codes.
- The customer asked the team to keep automated testing.
- The customer requested pagination for transaction history as part of the planned MVP v1 work.

## Decisions and approvals

- The customer is satisfied with the proposed Week 2 direction, the user stories, and the MVP v0 foundation.
- The customer approved the documented user stories and MoSCoW priorities after the team captured the requested follow-up items.
- The customer approved the initial proposed MVP v1 scope in [user-stories.md](./user-stories.md), with the expectation that the listed technical-specification requirements will be implemented in future weeks.
- The customer did not require redesign of the current API-first direction.

## Action points

- Preserve automated testing as an explicit planned requirement in `US-04`.
- Document HTTP status codes directly in [api/openapi.yaml](../../api/openapi.yaml) and [api/postman_collection.json](../../api/postman_collection.json).
- Plan and implement page-based transaction history pagination for `US-09`.
- Continue technical design for stale hold timeout handling in `US-05`.

## Risks and follow-up concerns

- MVP v0 is technically solid but still incomplete relative to the full technical specification.
- Pagination and stale hold auto-release remain the most visible functional gaps in the current product foundation.

## Resulting changes after the meeting

- Week 2 user stories were formalized in English in [user-stories.md](./user-stories.md).
- API interface documentation was prepared in [api/openapi.yaml](../../api/openapi.yaml) and [api/postman_collection.json](../../api/postman_collection.json), including explicit response codes.
- The MVP v0 report was updated to reflect current limitations and the customer's expectation that the remaining technical-specification items be added in subsequent weeks.
