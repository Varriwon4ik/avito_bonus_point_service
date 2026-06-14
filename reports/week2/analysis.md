
# Week 2 Analysis

## Learning points

- Writing the backlog in user-story form helped the team separate end-user needs from support, operator, and developer needs instead of treating the product as only a CRUD API.
- MoSCoW prioritization made it easier to keep the course MVP focused on core accrual and spending flows while still preserving operational and observability needs for later iterations.
- The API-first prototype work showed that interface documentation needs to describe response codes and error behavior, not only endpoint paths and payloads.
- MVP v0 proved that a transactional ledger foundation, idempotency, and hold-based spending can be delivered early even before all customer-facing requirements are complete.
- Customer validation showed that automated tests are seen as part of product quality and delivery readiness, not just an internal engineering preference.

## Validated assumptions

- The assumption that an API-first interface fits stakeholder needs was confirmed during the customer review, because the customer accepted OpenAPI and Postman artifacts as the main Week 2 prototype.
- The assumption that database transactions and row locking are a suitable foundation for concurrent balance changes was supported by the current implementation and by the integration scenarios in [internal/api/integration_test.go](../../internal/api/integration_test.go).
- The assumption that support staff need lot-level visibility was confirmed and remains represented by `US-02` in [user-stories.md](./user-stories.md).
- The assumption that a simple transaction-history listing would be enough was rejected. The customer explicitly asked for paginated history, which reinforces `US-09`.
- The assumption that endpoint documentation could stay lightweight was rejected. The customer asked the team to document HTTP status codes clearly in the interface artifacts.

## Needs clarification

- The exact pagination contract still needs a final implementation decision because the story text mentions `page` and `offset`, while the current codebase uses a `limit` query parameter.
- `US-05` still needs a final design for timeout duration, scheduler behavior, and ledger wording for auto-released holds.
- The future authentication model for administrator and internal-service operations remains open because MVP v0 currently runs without auth.
- The final public submission plan still needs external publication work for a public deployment, a public demo video, and screenshot evidence.

## Planned response

- Keep the initial proposed MVP v1 scope centered on `US-03`, `US-04`, `US-05`, `US-06`, `US-07`, and `US-09` in [user-stories.md](./user-stories.md).
- Extend [api/openapi.yaml](../../api/openapi.yaml) and [api/postman_collection.json](../../api/postman_collection.json) first, then align the implementation with the documented pagination and response-code behavior.
- Expand automated verification around core ledger flows in line with `US-04`, using the current integration tests as the baseline.
- Design and implement stale-hold recovery for `US-05` so that future MVP versions can recover from downstream crashes without manual operator intervention.
