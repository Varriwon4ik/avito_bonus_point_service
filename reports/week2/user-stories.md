



# Week 2 User Stories

The product serves several stakeholder groups, so the Week 2 backlog is written from multiple personas:

- Customer
- Administrator
- Customer support engineer
- Marketing service
- Platform operator
- Developer
- System administrator

## US-01: Run bulk promotional accruals

**Requirement status:** Active  
**MoSCoW priority:** Should Have

As a marketing service integrating with the bonus points ledger,  
I want to accrue points to many users in a single batch request,  
so that I can run loyalty campaigns without making thousands of individual API calls.

### Notes and constraints

Each batch item should have its own `idempotency_key` so retries do not duplicate accruals. The endpoint should return per-item results instead of failing the whole batch because of one invalid user. A maximum batch size is still open for clarification.

## US-02: Audit a user's points lots

**Requirement status:** Active  
**MoSCoW priority:** Must Have

As a support engineer investigating a balance discrepancy,  
I want to view all individual points lots for a user,  
so that I can explain the current balance and identify which lots are close to expiry.

### Notes and constraints

The existing web UI already shows lots for manual inspection, but the Week 2 requirement formalizes this need as a stable API capability for support tooling. The final API should support pagination and filtering by lot status such as `active`, `expired`, and `exhausted`.

## US-03: Earn bonus points after purchase

**Requirement status:** Active  
**MoSCoW priority:** Must Have

As a customer,  
I want to earn bonus points after making a purchase,  
so that I can use them for future discounts.

### Notes and constraints

Points are calculated from the purchase amount and are credited only after the purchase is confirmed. Reserved or failed purchases must not create a final accrual.

## US-04: Maintain automated regression coverage

**Requirement status:** Active  
**MoSCoW priority:** Must Have

As a developer,  
I want automated tests for all core bonus-point operations,  
so that the service can be verified, maintained, and deployed reliably.

### Notes and constraints

The minimum regression suite should cover successful accrual, successful spending, insufficient balance handling, duplicate order rejection, balance retrieval, transaction history retrieval, and concurrent request processing. Concurrent tests must confirm that parallel accrual and spending requests for the same user do not corrupt balance state or duplicate processing.

## US-05: Auto-release stale holds

**Requirement status:** Active  
**MoSCoW priority:** Must Have

As a platform operator,  
I want active holds that are not confirmed or cancelled within a configurable timeout to be released automatically,  
so that failed downstream services do not lock user points forever.

### Notes and constraints

This closes a known MVP v0 gap in the two-phase hold flow. The likely implementation is a background sweep that cancels stale holds, restores the original lots, and writes an audit entry such as `auto-released: timeout`. The timeout should be configurable through an environment variable.

## US-06: Confirm or cancel reserved points

**Requirement status:** Active  
**MoSCoW priority:** Must Have

As an administrator,  
I want reserved points to be either confirmed or canceled after a transaction,  
so that the balance reflects the final purchase outcome.

### Notes and constraints

Points must be reserved first and then either finalized or returned. Repeated confirmation or cancellation requests must stay idempotent, and stale reservations should eventually auto-cancel under the timeout policy described in [US-05](#us-05-auto-release-stale-holds).

## US-07: Manually accrue bonus points

**Requirement status:** Active  
**MoSCoW priority:** Must Have

As an administrator,  
I want to add bonus points to a user manually,  
so that the user can be rewarded for a purchase or service adjustment.

### Notes and constraints

The API call must include the target user ID and amount. Zero or negative values must be rejected. The long-term product should protect this operation with an admin-facing authentication mechanism even though MVP v0 currently runs on an internal, unauthenticated network.

## US-08: Validate accrual TTL bounds

**Requirement status:** Active  
**MoSCoW priority:** Should Have

As a system administrator,  
I want to enforce configurable minimum and maximum `ttl_days` values for accruals,  
so that calling services cannot create unrealistic expiration dates that break expiry logic or bloat storage.

### Notes and constraints

The product should introduce `MIN_TTL_DAYS` and `MAX_TTL_DAYS` environment variables while preserving `DEFAULT_TTL_DAYS` when `ttl_days` is omitted. Requests outside the allowed range should return a clear `400 Bad Request` response.

## US-09: Browse transaction history with pagination

**Requirement status:** Active  
**MoSCoW priority:** Must Have

As a customer support agent,  
I want to retrieve transaction history in pages,  
so that large transaction histories can be reviewed efficiently without loading all records at once.

### Notes and constraints

The planned API contract uses `page` and `offset` query parameters, for example `GET /v1/users/{user_id}/transactions?page=1&offset=20`. Invalid pagination parameters should return an explicit client error, and the response should contain only the transactions that belong to the requested page.

## US-10: Monitor requests and ledger health

**Requirement status:** Active  
**MoSCoW priority:** Should Have

As a platform operator,  
I want structured request logging and a basic metrics endpoint,  
so that I can monitor service health, debug production issues, and alert on anomalies.

### Notes and constraints

The preferred direction is structured logs with method, path, status code, latency, and user ID where applicable, plus a Prometheus-compatible `/metrics` endpoint. Request bodies should not be logged once they might contain sensitive data.

## Initial proposed MVP v1 scope

The initial proposed MVP v1 scope is intentionally small and includes only active `Must Have` stories:

- `US-03`
- `US-04`
- `US-05`
- `US-06`
- `US-07`
- `US-09`

This scope keeps the first delivery focused on core accrual and spending flows, operational confidence through tests, timeout safety for reservations, and scalable transaction-history access.
