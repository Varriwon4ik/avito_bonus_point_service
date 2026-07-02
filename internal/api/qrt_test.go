package api_test

// Automated Quality Requirement Tests (QRTs) for the Bonus Points Ledger
// Service. These are maintained product assets — see docs/quality-requirements.md
// and docs/quality-requirement-tests.md for the linked QR scenarios.
//
//   QRT-001  ->  QR-001 Time behaviour: balance-read p95 latency budget.
//   QRT-002  ->  QR-002 Integrity:      no overspend under concurrent debits.
//
// Both reuse the integration harness in integration_test.go and therefore
// self-skip when TEST_DATABASE_URL is not set (e.g. local runs without a
// database). In CI a Postgres service container is provided, so they execute.

import (
	"fmt"
	"math"
	"net/http"
	"os"
	"sort"
	"strconv"
	"sync"
	"testing"
	"time"
)

// TestQRT001BalanceResponseTime verifies QR-001 (Time behaviour): when an
// internal client reads a user's balance under a warmed-up, production-like
// in-process deployment, 95% of reads complete within the latency budget
// (default 200 ms; override with QRT_BALANCE_P95_BUDGET_MS).
func TestQRT001BalanceResponseTime(t *testing.T) {
	env := newTestEnv(t)
	user := "user_qrt_latency"

	// Seed several lots so the balance aggregation has real work to do.
	for i := 0; i < 10; i++ {
		mustAccrue(t, env, user, 100, 365, fmt.Sprintf("qrt-lat-acc-%d", i))
	}

	const warmup = 20
	const samples = 200
	url := env.Server.URL + "/v1/users/" + user + "/balance"

	// Warm up connections and query plans so one-time setup cost does not
	// pollute the measured sample.
	for i := 0; i < warmup; i++ {
		if st := doRaw(t, http.MethodGet, url, "", nil).StatusCode; st != http.StatusOK {
			t.Fatalf("warmup balance read: status=%d", st)
		}
	}

	latencies := make([]time.Duration, samples)
	for i := 0; i < samples; i++ {
		start := time.Now()
		res := doRaw(t, http.MethodGet, url, "", nil)
		latencies[i] = time.Since(start)
		if res.StatusCode != http.StatusOK {
			t.Fatalf("balance read %d: status=%d", i, res.StatusCode)
		}
	}

	sort.Slice(latencies, func(i, j int) bool { return latencies[i] < latencies[j] })
	p95 := latencies[int(math.Ceil(0.95*float64(samples)))-1]

	budget := 200 * time.Millisecond
	if v := os.Getenv("QRT_BALANCE_P95_BUDGET_MS"); v != "" {
		if ms, err := strconv.Atoi(v); err == nil && ms > 0 {
			budget = time.Duration(ms) * time.Millisecond
		}
	}

	t.Logf("QRT-001 balance read latency over %d samples: p50=%v p95=%v max=%v (budget p95<%v)",
		samples, latencies[samples/2], p95, latencies[samples-1], budget)

	if p95 > budget {
		t.Fatalf("QRT-001 FAILED: p95 balance latency %v exceeds budget %v", p95, budget)
	}
}

// TestQRT002DebitIntegrityUnderConcurrency verifies QR-002 (Integrity): when
// many debit requests for the same user are fired concurrently and together
// request far more than the available balance, the ledger must never overspend
// (no lost updates / double-spend) and the final balance must be exact and
// non-negative.
func TestQRT002DebitIntegrityUnderConcurrency(t *testing.T) {
	env := newTestEnv(t)
	user := "user_qrt_integrity"

	const initialBalance = 1000
	const debitAmount = 100
	const workers = 40 // 40 x 100 = 4000 requested >> 1000 available
	mustAccrue(t, env, user, initialBalance, 365, "qrt-integrity-acc")

	statuses := make([]int, workers)
	var ready, done sync.WaitGroup
	ready.Add(workers)
	done.Add(workers)

	for i := 0; i < workers; i++ {
		go func(i int) {
			defer done.Done()
			ready.Done()
			ready.Wait() // start pistol: maximise contention
			status, _, _ := doJSON(t, http.MethodPost,
				env.Server.URL+"/v1/users/"+user+"/debits",
				map[string]any{"amount": debitAmount, "idempotency_key": fmt.Sprintf("qrt-integrity-%d", i)},
			)
			statuses[i] = status // distinct index per goroutine: no shared write
		}(i)
	}
	done.Wait()

	successes := 0
	for _, s := range statuses {
		switch s {
		case http.StatusOK:
			successes++
		case http.StatusConflict:
			// expected once the balance is exhausted
		default:
			t.Errorf("unexpected status %d (want 200 or 409)", s)
		}
	}

	spent := successes * debitAmount
	// Invariant 1: never spend more than the available balance.
	if spent > initialBalance {
		t.Fatalf("QRT-002 FAILED: spent %d exceeds initial balance %d (double-spend)", spent, initialBalance)
	}
	// With 1000 available and 100 per debit, exactly 10 must succeed.
	if want := initialBalance / debitAmount; successes != want {
		t.Errorf("QRT-002: expected exactly %d successful debits, got %d", want, successes)
	}

	// Invariant 2: final balance is exactly what is left and never negative.
	status, _, bal := doJSON(t, http.MethodGet, env.Server.URL+"/v1/users/"+user+"/balance", nil)
	if status != http.StatusOK {
		t.Fatalf("balance check: status=%d body=%v", status, bal)
	}
	gotAvailable := int(bal["available"].(float64))
	if gotAvailable < 0 {
		t.Fatalf("QRT-002 FAILED: negative balance %d", gotAvailable)
	}
	if want := initialBalance - spent; gotAvailable != want {
		t.Errorf("QRT-002: available=%d, want %d (1000 - %d spent)", gotAvailable, want, spent)
	}

	t.Logf("✓ QRT-002 integrity: %d/%d debits succeeded, spent=%d, available=%d (no overspend)",
		successes, workers, spent, gotAvailable)
}
