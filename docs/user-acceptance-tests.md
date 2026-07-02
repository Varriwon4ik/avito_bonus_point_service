# User Acceptance Tests

End-user-facing acceptance scenarios for the Bonus Points Ledger Service that the
customer (or a relevant stakeholder) can execute to confirm the product supports
its intended goals. UAT scenarios are **maintained product assets**; IDs are
stable and execution history is preserved (new results are appended, never
overwritten), per `Process_Requirements.md`.

The "end user" of this service is the integrating online store and its support
operators: they accrue points, read balances, run two-phase redemptions at
checkout, and audit a user's history. Scenarios are executed against the deployed
increment (web UI at `/`, Swagger UI at `/docs`) or via the HTTP API.

| ID | Title | Status | Verifies |
|---|---|---|---|
| UAT-001 | Read a user's points balance | Active | Core balance workflow |
| UAT-002 | Two-phase redemption at checkout (hold → confirm / cancel) | Active | Core debit safety |
| UAT-003 | Review paginated transaction history | Active | US-09 (Sprint 2) |

---

## UAT-001: Read a user's points balance

- **Status:** Active
- **User goal:** A store operator can see how many points a user can spend now,
  how many are reserved, and how many expire soon.
- **Preconditions:** The service is deployed and reachable; a test user exists or
  can be created by accruing points.

**Steps:**

1. Accrue points for a test user: `POST /v1/users/{id}/accruals` with
   `{ "amount": 500, "ttl_days": 180, "idempotency_key": "uat1-accrual" }`.
2. Read the balance: `GET /v1/users/{id}/balance?expiring_within_days=7`.

**Expected outcome:** Step 1 returns `201` with a `lot_id` and `expires_at`. Step
2 returns `200` with `available`, `held`, `total`, and `expiring_soon` consistent
with the accrual (e.g. `available = 500`, `held = 0`, `total = 500`). An unknown
user returns `404` with the standard error envelope.

### Execution history

**2026-06-28 (Sprint 2 review / UAT session)** — **Passed.** The customer
inspected the balance workflow and the underlying handler/query during the code
walkthrough and confirmed the returned fields and semantics match expectations.
No issues raised.

---

## UAT-002: Two-phase redemption at checkout (hold → confirm / cancel)

- **Status:** Active
- **User goal:** The store can safely reserve a user's points during checkout and
  either finalise the redemption or release the points if the order is abandoned,
  without ever overspending.
- **Preconditions:** A test user with a known available balance (e.g. 500).

**Steps:**

1. Create a hold: `POST /v1/users/{id}/holds` with
   `{ "amount": 200, "idempotency_key": "uat2-hold" }`.
2. Read the balance and confirm `held = 200`, `available = 300`.
3. Either confirm `POST /v1/holds/{hold_id}/confirm` (finalise) **or** cancel
   `POST /v1/holds/{hold_id}/cancel` (release).
4. Read the balance again.

**Expected outcome:** Step 1 returns `201` with `status: active`. After confirm,
the points are permanently spent (`available = 300`, `held = 0`). After cancel,
the points return to the user (`available = 500`, `held = 0`). Confirm/cancel are
idempotent. A hold that exceeds the available balance returns `409 Conflict`.

### Execution history

**2026-06-28 (Sprint 2 review / UAT session)** — **Passed.** The customer
reviewed the two-phase hold/confirm/cancel design and the concurrency safety
(row-level locking, no double-spend) and accepted the approach.

---

## UAT-003: Review paginated transaction history

- **Status:** Active
- **User goal:** A support operator can page through a user's ledger/audit history
  instead of receiving one unbounded list, to investigate a specific transaction.
- **Preconditions:** A test user with several ledger entries (accruals, holds,
  debits).

**Steps:**

1. Generate several transactions for the user (accruals + a debit or two).
2. Request the first page:
   `GET /v1/users/{id}/transactions?page=1&offset=20`.
3. Request the next page: `GET /v1/users/{id}/transactions?page=2&offset=20`.
4. Try an invalid value: `GET /v1/users/{id}/transactions?offset=0`.

**Expected outcome:** Steps 2–3 return `200` with the envelope
`{ user_id, page, offset, total, entries }`, newest entries first, with `entries`
bounded by `offset` (page size) and `total` reflecting the full count. Step 4
returns `400 Bad Request` with the standard error envelope (`offset must be
between 1 and 500`).

### Execution history

**2026-06-28 (Sprint 2 review / UAT session)** — **Passed.** The customer
specifically requested a demonstration of paginated history access; the team
showed the `page`/`offset` parameters and the response envelope. The customer was
satisfied that the Sprint 2 changes deliver the requested capability.

---

## Customer feedback and resulting backlog decisions (2026-06-28)

The customer's primary feedback was a request — carried over from the Sprint 1
review and reaffirmed — to provide **additional automated tests run against the
earlier version of the product** to objectively prove that the team's changes are
valid. This is addressed by the US-15 autotester (replayable accrual/concurrency
scenarios) and the automated QRTs/CI gates; further regression coverage of older
code paths remains a Sprint 3 follow-up. No UAT scenario failed; no defect PBIs
were opened from this session.

Summarised public results for the assignment are in
[reports/week4/README.md](../reports/week4/README.md) and
[reports/week4/customer-review-summary.md](../reports/week4/customer-review-summary.md).
