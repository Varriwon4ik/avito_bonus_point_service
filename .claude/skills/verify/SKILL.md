---
name: verify
description: Build, launch, and drive this service (Go API + embedded web UI) to verify changes end-to-end.
---

# Verify: avito bonus point service

Single Go binary serving the REST API and the admin web UI (embedded from
`cmd/api/web/index.html` via `//go:embed` — **rebuild the binary after any
HTML change**, it is not served from disk).

## Build & launch

```bash
docker compose up -d postgres          # DB; wait for pg_isready -U bonus
go build -o /tmp/api.exe ./cmd/api
DB_DSN="postgres://bonus:bonus@localhost:5432/bonus_ledger?sslmode=disable" /tmp/api.exe -port 8091
curl -s http://localhost:8091/healthz  # {"status":"ok"}
```

Default port is 8080; pass `-port` to avoid clashes. The UI is at `/`,
OpenAPI docs at `/docs`.

## Driving the web UI

No Playwright browsers are cached on this machine, but system Edge works:

```bash
cd <scratchpad> && npm init -y && npm i playwright   # package only, fast
# in the script: chromium.launch({ channel: 'msedge', headless: true })
```

Useful element ids: nav buttons `#nav-accrue`/`#nav-autotester`; autotester
`#at-mode`, `#at-run-btn`, `#at-results-tbody`; bulk accrual `#bulk-rows`,
`#bulk-run-btn`, `#bulk-results-tbody`; toasts in `#toast`; raw responses in
`#ac-result` / `#at-result`.

## API surfaces worth curling

- `POST /v1/autotest/run` — `{label, user_id, amount, parallel_requests, mode}`;
  `mode` is `single` (default) or `multi_key`. Runs real accruals against a
  dedicated `autotest-`-prefixed user.
- `POST /v1/accruals/batch` — `{items:[{user_id, amount, idempotency_key, ...}]}`;
  always answers HTTP 207 with per-item `results`.
- Replaying an idempotency key returns the original cached body with
  `status: "created"` (not a conflict), even if the payload differs.

## Gotchas

- Tests need `TEST_DATABASE_URL` (see Makefile); they truncate tables in the
  target DB.
- Manual runs leave data in the dev DB (`bonus_ledger`); autotest users are
  `autotest-`-prefixed, other test users are whatever you typed.
