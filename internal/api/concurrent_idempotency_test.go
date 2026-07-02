package api_test

// US-11: Verify idempotent key deduplication under concurrent requests.
//
// Scenario (exactly as described in Issue #8):
//   - User starts with 1 000 points.
//   - Two debit requests are fired simultaneously:
//       • Request A — 500 points  (idempotency key "debit-conc-A")
//       • Request B — 800 points  (idempotency key "debit-conc-B")
//   - Together they exceed the balance (500 + 800 = 1 300 > 1 000).
//   - The database-level SELECT … FOR UPDATE inside allocateLots serialises
//     the two concurrent transactions, so only one can win.
//   - Expected outcome: exactly one 200 OK and exactly one 409 Conflict.
//   - Balance after: 500 OR 800 points were spent (whichever request won),
//     NOT 1 300 (i.e. no "double-spend" / last-write-wins corruption).
//
// A second sub-test (TestConcurrentSameKeyDeduplication) fires the *same*
// idempotency key from two goroutines at the same time and asserts that the
// key is processed exactly once (no duplicate accrual / double-spend).

import (
	"fmt"
	"net/http"
	"sync"
	"testing"
)

// TestConcurrentDebitRaceCondition is the canonical demonstration for
// Alexander: two simultaneous debits that together exceed the balance must
// result in one success and one 409, with the final balance reflecting only
// the winning debit — never a double-spend.
func TestConcurrentDebitRaceCondition(t *testing.T) {
	env := newTestEnv(t)
	user := "user_concurrent_debit"

	// Give the user exactly 1 000 points.
	mustAccrue(t, env, user, 1000, 30, "acc-concurrent")

	type result struct {
		status int
		body   map[string]any
	}

	resultCh := make(chan result, 2)

	// Fire both debit requests simultaneously using a sync.WaitGroup as a
	// starting pistol so they hit the server as concurrently as Go allows.
	var ready sync.WaitGroup
	ready.Add(2)

	fire := func(amount int, key string) {
		ready.Done() // signal this goroutine is ready
		ready.Wait() // wait until both goroutines are ready before sending

		status, _, body := doJSON(t, http.MethodPost,
			env.Server.URL+"/v1/users/"+user+"/debits",
			map[string]any{"amount": amount, "idempotency_key": key},
		)
		resultCh <- result{status, body}
	}

	go fire(500, "debit-conc-A") // Request A: 500 pts
	go fire(800, "debit-conc-B") // Request B: 800 pts

	r1 := <-resultCh
	r2 := <-resultCh

	// ── Assertions ────────────────────────────────────────────────────────

	statuses := [2]int{r1.status, r2.status}
	successes := 0
	conflicts := 0
	for _, s := range statuses {
		switch s {
		case http.StatusOK:
			successes++
		case http.StatusConflict:
			conflicts++
		default:
			t.Errorf("unexpected HTTP status %d (want 200 or 409)", s)
		}
	}

	if successes != 1 {
		t.Errorf("want exactly 1 successful debit, got %d (statuses: %v)", successes, statuses)
	}
	if conflicts != 1 {
		t.Errorf("want exactly 1 conflict (409), got %d (statuses: %v)", conflicts, statuses)
	}

	// Determine the winner's amount so we can verify the final balance.
	var winnerAmount float64
	for _, r := range []result{r1, r2} {
		if r.status == http.StatusOK {
			if amt, ok := r.body["amount"].(float64); ok {
				winnerAmount = amt
			}
		}
	}

	// The balance must equal 1000 − winnerAmount. Under the broken
	// last-write-wins scenario it would be 1000 − 1300 = −300 (or
	// the DB would have prevented that and left an inconsistent state).
	expectedAvailable := 1000 - int(winnerAmount)

	status, _, bal := doJSON(t, http.MethodGet, env.Server.URL+"/v1/users/"+user+"/balance", nil)
	if status != http.StatusOK {
		t.Fatalf("balance check failed: status=%d body=%v", status, bal)
	}

	gotAvailable := int(bal["available"].(float64))
	if gotAvailable != expectedAvailable {
		t.Errorf(
			"balance mismatch: want available=%d (1000 − %d), got %d — possible double-spend!",
			expectedAvailable, int(winnerAmount), gotAvailable,
		)
	}

	t.Logf("✓ concurrent debit race: winner spent %d pts, loser got 409, balance=%d (no double-spend)",
		int(winnerAmount), gotAvailable)
}

// TestConcurrentSameKeyDeduplication fires the *identical* idempotency key
// from two goroutines at the same time. Only one transaction should actually
// execute; the other must receive the cached response (or a 409 idempotency-
// conflict while the first is still in flight). Either way, the accrual must
// happen exactly once.
func TestConcurrentSameKeyDeduplication(t *testing.T) {
	env := newTestEnv(t)
	user := "user_same_key"

	type result struct {
		status int
		body   map[string]any
	}

	resultCh := make(chan result, 2)

	var ready sync.WaitGroup
	ready.Add(2)

	fireAccrue := func(goroutineID int) {
		ready.Done()
		ready.Wait()

		status, _, body := doJSON(t, http.MethodPost,
			env.Server.URL+"/v1/users/"+user+"/accruals",
			map[string]any{
				"amount":          500,
				"ttl_days":        30,
				"idempotency_key": "same-key-accrual", // identical key on both goroutines
			},
		)
		t.Logf("goroutine %d → status=%d body=%v", goroutineID, status, body)
		resultCh <- result{status, body}
	}

	go fireAccrue(1)
	go fireAccrue(2)

	r1 := <-resultCh
	r2 := <-resultCh

	// ── Assertions ────────────────────────────────────────────────────────

	// Both responses must be either 201 Created (idempotent replay) or
	// 409 Conflict (key in-flight). A 409 here is the
	// ErrIdempotencyConflict sentinel — not an insufficient-funds error.
	for _, r := range []result{r1, r2} {
		if r.status != http.StatusCreated && r.status != http.StatusConflict {
			t.Errorf("unexpected status %d; want 201 or 409", r.status)
		}
	}

	// The user's balance must be exactly 500, never 1000 (no double-credit).
	status, _, bal := doJSON(t, http.MethodGet, env.Server.URL+"/v1/users/"+user+"/balance", nil)
	if status != http.StatusOK {
		t.Fatalf("balance check failed: status=%d body=%v", status, bal)
	}

	gotAvailable := int(bal["available"].(float64))
	if gotAvailable != 500 {
		t.Errorf("want available=500 (accrued once), got %d — possible double-credit!", gotAvailable)
	}

	t.Logf("✓ same-key deduplication: balance=%d (accrued exactly once)", gotAvailable)
}

// TestConcurrentHoldsExhaustBalance is a more exhaustive stress variant:
// N goroutines race to hold the same pool of points; the sum of all
// successful holds must never exceed the initial balance.
func TestConcurrentHoldsExhaustBalance(t *testing.T) {
	env := newTestEnv(t)
	user := "user_exhaust"

	const initialBalance = 1000
	const holdAmount = 300 // 3 full holds fit (3×300=900), 4th must fail
	const workers = 6

	mustAccrue(t, env, user, initialBalance, 30, "acc-exhaust")

	type result struct {
		status int
		amount int
	}

	results := make([]result, workers)
	var wg sync.WaitGroup

	for i := 0; i < workers; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			status, _, body := doJSON(t, http.MethodPost,
				env.Server.URL+"/v1/users/"+user+"/holds",
				map[string]any{
					"amount":          holdAmount,
					"idempotency_key": fmt.Sprintf("hold-exhaust-%d", i),
				},
			)
			results[i] = result{status: status, amount: holdAmount}
			_ = body
		}(i)
	}
	wg.Wait()

	// Count how many succeeded and verify the total does not exceed the balance.
	totalHeld := 0
	for _, r := range results {
		if r.status == http.StatusCreated {
			totalHeld += r.amount
		} else if r.status != http.StatusConflict {
			t.Errorf("unexpected status %d", r.status)
		}
	}

	if totalHeld > initialBalance {
		t.Errorf("double-spend detected: total held=%d exceeds initial balance=%d", totalHeld, initialBalance)
	}

	status, _, bal := doJSON(t, http.MethodGet, env.Server.URL+"/v1/users/"+user+"/balance", nil)
	if status != http.StatusOK {
		t.Fatalf("balance check: status=%d body=%v", status, bal)
	}

	gotHeld := int(bal["held"].(float64))
	gotAvailable := int(bal["available"].(float64))

	if gotHeld != totalHeld {
		t.Errorf("ledger held=%d does not match sum of successful holds=%d", gotHeld, totalHeld)
	}
	if gotHeld+gotAvailable != initialBalance {
		t.Errorf("held(%d)+available(%d) != initialBalance(%d)", gotHeld, gotAvailable, initialBalance)
	}

	t.Logf("✓ %d workers × %d pts: %d pts held, %d pts still available (no double-spend)",
		workers, holdAmount, gotHeld, gotAvailable)
}
