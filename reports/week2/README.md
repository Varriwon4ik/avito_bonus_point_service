# Week 2 Report Index

## Project

- Project name: Bonus Points Ledger Service
- Short description: API-first bonus points ledger with point accrual, balance tracking, lot expiry, holds, debits, and a local smoke-check web UI
- License: [MIT License](../../LICENSE)

## Week 2 repository files

- User stories: [user-stories.md](./user-stories.md)
- MVP v0 report: [mvp-v0-report.md](./mvp-v0-report.md)
- Customer meeting transcript: [customer-meeting-transcript.md](./customer-meeting-transcript.md)
- Customer meeting summary: [customer-meeting-summary.md](./customer-meeting-summary.md)
- Week 2 analysis: [analysis.md](./analysis.md)
- LLM usage report: [llm-report.md](./llm-report.md)

## Interface prototype artifacts

- OpenAPI specification: [api/openapi.yaml](../../api/openapi.yaml)
- Postman collection: [api/postman_collection.json](../../api/postman_collection.json)
- Current local MVP v0 interface surface: [embedded web UI](../../cmd/api/web/index.html)

No separate `docs/interface.md` file is needed because the product's externally used interface is the API documented above.

## MVP v0 access

- Runnable artifact: [docker-compose.yml](../../docker-compose.yml)
- Container build: [Dockerfile](../../Dockerfile)
- Local setup instructions: [root README](../../README.md)
- Smoke-check scenario: [mvp-v0-report.md](./mvp-v0-report.md)
- Default local URL after startup: `http://localhost:8080`
- Public video demonstration: not available in this local repository snapshot yet

## PR/MR workflow and link checking

- Minimal PR template: [.github/pull_request_template.md](../../.github/pull_request_template.md)
- Lychee workflow: [.github/workflows/lychee.yml](../../.github/workflows/lychee.yml)
- Lychee exclusions: [.lycheeignore](../../.lycheeignore)
- Reviewed PR/MR links: not available in this local repository snapshot
- Latest successful protected-default-branch Lychee run: not available in this local repository snapshot

### Excluded Lychee links and manual verification

- `http://localhost:8080` is excluded because it is only available after starting the local MVP v0 stack. It should be manually verified during the smoke check.
- `http://localhost:8080/healthz` is excluded for the same reason and should also be manually verified during the smoke check.

## Coverage

- The API prototype in [api/openapi.yaml](../../api/openapi.yaml) represents `US-03`, `US-06`, `US-07`, and `US-09`.
- The Postman collection in [api/postman_collection.json](../../api/postman_collection.json) demonstrates representative success and error flows for `US-03`, `US-06`, `US-07`, and `US-09`.
- MVP v0 currently provides the technical foundation for `US-03`, `US-04` (partial through tests), `US-06`, and `US-07`, and it also supports support-oriented inspection related to `US-02`.
- The remaining planned gaps most relevant to MVP v1 are `US-05` and the paginated implementation of `US-09`, both discussed in [analysis.md](./analysis.md) and [mvp-v0-report.md](./mvp-v0-report.md).

## Screenshots

The current task explicitly required that no `reports/week2/images/` directory be created, so screenshot evidence is not committed in this repository snapshot.

## Customer review evidence

- Published transcript: [customer-meeting-transcript.md](./customer-meeting-transcript.md)
- Meeting summary: [customer-meeting-summary.md](./customer-meeting-summary.md)
- `customer-meeting-notes.md` was intentionally not created because this draft assumes transcript publication permission.
