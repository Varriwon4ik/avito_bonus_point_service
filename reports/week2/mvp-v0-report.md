# MVP v0 Report

## Purpose and description

MVP v0 is a runnable technical foundation for the bonus points ledger service. It already provides persistent point accrual, balance lookup, lots inspection, hold creation, hold confirmation, hold cancellation, one-shot debit, idempotent request handling, and a small embedded web UI for smoke-check usage. The goal of this version is not to deliver the full course MVP, but to prove that the core ledger model, database transactions, and basic operator workflow already run end to end.

## Deployment URL or runnable artifact

- Runnable artifact: [docker-compose.yml](../../docker-compose.yml)
- Container build definition: [Dockerfile](../../Dockerfile)
- Default local access URL after startup: `http://10.93.26.175:8080/`
- Health endpoint: `http://10.93.26.175:8080/healthz`

## Customer feedback context

The customer is satisfied with the current product foundation and with the demonstrated MVP v0 direction. At the same time, the customer expects the team to add the remaining technical-specification items in the coming weeks, especially basic automated tests, explicit HTTP status-code documentation for the API, and paginated transaction history.

## Current limitations, placeholders, and mocks

- MVP v0 is a local runnable artifact in this repository snapshot and is not yet published as a public internet deployment.
- The service is intentionally unauthenticated in MVP v0 because it currently targets an internal-network usage model.
- Transaction history currently supports a `limit` parameter in code, while the proposed MVP v1 interface documents page-based pagination as a planned improvement.
- Hold timeout auto-release is a planned requirement and not yet implemented.
- The embedded web UI is a smoke-check and demonstration surface, not a production-ready customer interface.

## Local setup instructions

Local setup is documented in the root [README.md](../../README.md). The repository also includes a sanitized [`.env.example`](../../.env.example) file for the main environment values used by local runs and tests.

## Repeatable smoke-check scenario

### Access instructions

1. Open `http://10.93.26.175:8080/` in a browser.

### Steps

1. Open `http://10.93.26.175:8080/healthz` and confirm that the API returns `{"status":"ok"}`.
2. In the embedded web UI, keep the default user ID or enter a new test user.
3. Use the `Accrue Points` screen to add a positive amount with an idempotency key.
4. Return to `Dashboard` and refresh the balance.
5. Open `Holds & Debit` and create a hold for part of the balance.
6. Confirm or cancel the hold and refresh the dashboard again.

### Expected results

- The health endpoint returns HTTP `200 OK`.
- The UI loads successfully and can communicate with the API.
- After accrual, the user's available balance increases.
- After creating a hold, `available` decreases and `held` increases.
- After confirming a hold, the held amount is permanently spent.
- After cancelling a hold, the reserved points return to the available balance.
