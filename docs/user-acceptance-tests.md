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
| UAT-004 | See exact HTTP response codes in the web UI | Active | US-16 (Sprint 3) |
| UAT-005 | Run the autotester from the web UI | Active | US-17 (Sprint 3) |
| UAT-006 | Label a transaction and find it in the history | Active | US-18, US-08 (Sprint 3) |

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

## UAT-004: See exact HTTP response codes in the web UI

- **Status:** Active
- **User goal:** A system administrator sees the exact HTTP status code and
  outcome of every accrual/debit operation directly in the web UI, so problems
  can be reported and diagnosed precisely.
- **Verifies:** [US-16 / #39](https://github.com/Varriwon4ik/avito_bonus_point_service/issues/39)
  and its acceptance criteria (frontend displays the backend response code
  alongside the other request details).
- **Preconditions:** The service is deployed and reachable; the web UI opens at `/`.

**Steps:**

1. In the **Accrue** view, accrue points for a test user (valid amount and
   idempotency key) and submit.
2. In the **Debit** view, attempt to debit points from a non-existent user.
3. Attempt an accrual with an invalid amount (e.g. `0`).

**Expected outcome:** Step 1 shows a success status (`200`/`201`, green) next
to the operation result. Step 2 shows `404` (red) with the not-found message.
Step 3 shows `400` with the validation message. In every case the displayed
code matches the actual backend response (verifiable via Swagger UI or curl).

### Execution history

**2026-07-03 (Sprint 3 review / UAT session)** — **Passed.** The implementer
shared the screen and the customer directed the demonstration: a successful
accrual showed the green success code, and a debit against a non-existent user
showed the red `404`. The customer accepted the feature and left a non-blocking
suggestion to reconsider how the codes are presented ("you can change the
response codes here... it's not critical").

---

## UAT-005: Run the autotester from the web UI

- **Status:** Active
- **User goal:** A system administrator runs the built-in autotester scenarios
  against the live instance from the browser and reads a per-check pass/fail
  report, without needing the console tool.
- **Verifies:** [US-17 / #40](https://github.com/Varriwon4ik/avito_bonus_point_service/issues/40)
  and its acceptance criteria (Autotester tab drives the backend engine and
  displays its results).
- **Preconditions:** The service is deployed and reachable; the web UI opens at `/`.

**Steps:**

1. Open the **Autotester** tab.
2. Fill in a scenario: label, test user id, amount, TTL days, number of
   parallel requests.
3. Run the scenario and read the per-check report.

**Expected outcome:** The run returns a report listing each check (accrual
correctness, parallel accrual) with pass/fail per check and an overall verdict.
The target user is forced under the `autotest-` prefix, so no real account is
touched; the results match what `cmd/autotest` reports for the same scenario.

### Execution history

**2026-07-03 (Sprint 3 review / UAT session)** — **Passed.** The customer
directed the demonstration of the Autotester tab and probed the scenario
semantics (single vs. multiple idempotency keys). Feature accepted. Resulting
backlog item: extend the autotester with a parallel multi-idempotency-key
scenario —
[US-19 / #50](https://github.com/Varriwon4ik/avito_bonus_point_service/issues/50).

---

## UAT-006: Label a transaction and find it in the history

- **Status:** Active
- **User goal:** A system administrator marks accruals with a label (preset
  such as `test`/`real`, or a custom one) and later tells labelled transactions
  apart in the history, e.g. to separate test traffic from real traffic.
- **Verifies:** [US-18 / #41](https://github.com/Varriwon4ik/avito_bonus_point_service/issues/41)
  and its acceptance criteria; step 4 also exercises the US-08
  ([#5](https://github.com/Varriwon4ik/avito_bonus_point_service/issues/5))
  TTL bounds.
- **Preconditions:** The service is deployed and reachable; a test user exists
  or can be created by accruing points.

**Steps:**

1. In the **Accrue** view, accrue points choosing the preset label `test`.
2. Accrue again with a custom label (e.g. `promo-july`).
3. Open the **Transactions** view for the user and locate both entries.
4. Attempt an accrual with `ttl_days` outside the configured
   `MIN_TTL_DAYS`–`MAX_TTL_DAYS` range.

**Expected outcome:** Steps 1–2 succeed (`201`); labels are trimmed/validated
(≤ 32 chars, no control characters — invalid labels get `400`). Step 3 shows
both ledger entries with their labels displayed. Step 4 returns
`400 Bad Request` with a clear message, and no lot is created.

### Execution history

**2026-07-03 (Sprint 3 review / UAT session)** — **Passed.** The implementer
demonstrated preset and custom labels under the customer's direction; the
customer specifically praised custom labels ("you can assign a certain product
with a label to its transaction and have a convenient tool to confirm debits").
The TTL-validation behaviour (US-08) was demonstrated in the same session and
accepted. No defect PBI opened.

---

## Customer feedback and resulting backlog decisions (2026-07-03)

The Sprint 3 session was run screen-share style: each implementer demonstrated
their feature while the **customer directed the demonstration**. UAT-004,
UAT-005, and UAT-006 all passed; the older scenarios UAT-001–UAT-003 were not
formally re-executed but their core flows (accrual, debit, transaction history)
were exercised throughout the new demonstrations and behaved as previously
accepted. Feedback converted into backlog decisions:

- **Autotester with multiple parallel idempotency keys** → new story
  [US-19 / #50](https://github.com/Varriwon4ik/avito_bonus_point_service/issues/50)
  (added to the Product Backlog during the meeting).
- **Response-code presentation could be improved** → non-blocking suggestion on
  US-16; noted for future UI polish, no PBI opened at the customer's own
  assessment ("not critical").
- **Run tests against a demo/earlier version of the product** → still open;
  carried in [docs/roadmap.md](roadmap.md) as a Sprint 4 candidate.
- **Team-process feedback** ("work as a team, a whole") → taken into the
  Sprint 3 retrospective action points
  ([reports/week5/retrospective.md](../reports/week5/retrospective.md)).

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
