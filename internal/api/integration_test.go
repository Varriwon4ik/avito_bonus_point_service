package api_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"os"
	"sync"
	"testing"

	_ "github.com/lib/pq"

	"bonus-ledger/internal/api"
	"bonus-ledger/internal/data"
)

// newTestServer connects to TEST_DATABASE_URL, applies migrations and wipes
// all tables so each test starts from a clean slate. Tests are skipped if
// TEST_DATABASE_URL is not set.
func newTestServer(t *testing.T) *httptest.Server {
	t.Helper()

	dsn := os.Getenv("TEST_DATABASE_URL")
	if dsn == "" {
		t.Skip("TEST_DATABASE_URL not set, skipping integration test")
	}

	db, err := data.OpenDB(dsn)
	if err != nil {
		t.Fatalf("open db: %v", err)
	}
	t.Cleanup(func() { db.Close() })

	if err := data.Migrate(db); err != nil {
		t.Fatalf("migrate: %v", err)
	}

	for _, table := range []string{"hold_allocations", "holds", "ledger_entries", "points_lots", "idempotency_keys"} {
		if _, err := db.Exec("TRUNCATE TABLE " + table + " RESTART IDENTITY CASCADE"); err != nil {
			t.Fatalf("truncate %s: %v", table, err)
		}
	}

	store := data.NewStore(db)
	logger := slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{Level: slog.LevelError}))
	srv := api.NewServer(store, logger, 365)

	ts := httptest.NewServer(srv)
	t.Cleanup(ts.Close)
	return ts
}

func doJSON(t *testing.T, method, url string, payload any) (int, map[string]any) {
	t.Helper()

	var body *bytes.Reader
	if payload != nil {
		b, err := json.Marshal(payload)
		if err != nil {
			t.Fatalf("marshal: %v", err)
		}
		body = bytes.NewReader(b)
	} else {
		body = bytes.NewReader(nil)
	}

	req, err := http.NewRequest(method, url, body)
	if err != nil {
		t.Fatalf("new request: %v", err)
	}
	if payload != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("do request: %v", err)
	}
	defer resp.Body.Close()

	var out map[string]any
	if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	return resp.StatusCode, out
}

func TestAccrualIdempotency(t *testing.T) {
	ts := newTestServer(t)
	user := "user_accrual"

	payload := map[string]any{"amount": 100, "ttl_days": 30, "idempotency_key": "order-1"}

	status, body := doJSON(t, http.MethodPost, ts.URL+"/v1/users/"+user+"/accruals", payload)
	if status != http.StatusCreated {
		t.Fatalf("first accrual: status=%d body=%v", status, body)
	}
	firstLotID := body["lot_id"]

	// repeat the exact same request: must not create a second lot
	status, body2 := doJSON(t, http.MethodPost, ts.URL+"/v1/users/"+user+"/accruals", payload)
	if status != http.StatusCreated {
		t.Fatalf("second accrual: status=%d body=%v", status, body2)
	}
	if body2["lot_id"] != firstLotID {
		t.Fatalf("expected idempotent replay to return same lot_id, got %v vs %v", firstLotID, body2["lot_id"])
	}

	status, bal := doJSON(t, http.MethodGet, ts.URL+"/v1/users/"+user+"/balance", nil)
	if status != http.StatusOK {
		t.Fatalf("balance: status=%d body=%v", status, bal)
	}
	if bal["available"].(float64) != 100 {
		t.Fatalf("expected available=100 after idempotent retry, got %v", bal["available"])
	}
}

func TestHoldConfirmCancel(t *testing.T) {
	ts := newTestServer(t)
	user := "user_hold"

	mustAccrue(t, ts, user, 1000, 30, "acc-1")

	status, hold := doJSON(t, http.MethodPost, ts.URL+"/v1/users/"+user+"/holds",
		map[string]any{"amount": 400, "idempotency_key": "hold-1"})
	if status != http.StatusCreated {
		t.Fatalf("create hold: status=%d body=%v", status, hold)
	}
	holdID := int64(hold["hold_id"].(float64))

	status, bal := doJSON(t, http.MethodGet, ts.URL+"/v1/users/"+user+"/balance", nil)
	if status != http.StatusOK || bal["available"].(float64) != 600 || bal["held"].(float64) != 400 {
		t.Fatalf("balance after hold: status=%d body=%v", status, bal)
	}

	// cancel releases the hold back to available
	status, cancel := doJSON(t, http.MethodPost, fmt.Sprintf("%s/v1/holds/%d/cancel", ts.URL, holdID), nil)
	if status != http.StatusOK || cancel["status"] != "cancelled" {
		t.Fatalf("cancel hold: status=%d body=%v", status, cancel)
	}

	status, bal = doJSON(t, http.MethodGet, ts.URL+"/v1/users/"+user+"/balance", nil)
	if status != http.StatusOK || bal["available"].(float64) != 1000 || bal["held"].(float64) != 0 {
		t.Fatalf("balance after cancel: status=%d body=%v", status, bal)
	}

	// a fresh hold, this time confirmed, permanently removes the points
	status, hold2 := doJSON(t, http.MethodPost, ts.URL+"/v1/users/"+user+"/holds",
		map[string]any{"amount": 250, "idempotency_key": "hold-2"})
	if status != http.StatusCreated {
		t.Fatalf("create hold 2: status=%d body=%v", status, hold2)
	}
	holdID2 := int64(hold2["hold_id"].(float64))

	status, confirm := doJSON(t, http.MethodPost, fmt.Sprintf("%s/v1/holds/%d/confirm", ts.URL, holdID2), nil)
	if status != http.StatusOK || confirm["status"] != "confirmed" {
		t.Fatalf("confirm hold: status=%d body=%v", status, confirm)
	}

	// confirming again is a no-op
	status, confirm2 := doJSON(t, http.MethodPost, fmt.Sprintf("%s/v1/holds/%d/confirm", ts.URL, holdID2), nil)
	if status != http.StatusOK || confirm2["status"] != "confirmed" {
		t.Fatalf("re-confirm hold: status=%d body=%v", status, confirm2)
	}

	status, bal = doJSON(t, http.MethodGet, ts.URL+"/v1/users/"+user+"/balance", nil)
	if status != http.StatusOK || bal["available"].(float64) != 750 || bal["held"].(float64) != 0 {
		t.Fatalf("balance after confirm: status=%d body=%v", status, bal)
	}
}

func TestInsufficientFunds(t *testing.T) {
	ts := newTestServer(t)
	user := "user_poor"

	mustAccrue(t, ts, user, 50, 30, "acc-1")

	status, body := doJSON(t, http.MethodPost, ts.URL+"/v1/users/"+user+"/holds",
		map[string]any{"amount": 100, "idempotency_key": "hold-1"})
	if status != http.StatusConflict {
		t.Fatalf("expected 409 insufficient funds, got status=%d body=%v", status, body)
	}
}

// TestExpiryOrdering verifies that debits consume the lot with the nearest
// expiry date first, even though it was accrued after a longer-lived lot.
func TestExpiryOrdering(t *testing.T) {
	ts := newTestServer(t)
	user := "user_fifo"

	// accrued first, expires later (60 days)
	mustAccrue(t, ts, user, 100, 60, "acc-long")
	// accrued second, but expires sooner (5 days) -> must be spent first
	mustAccrue(t, ts, user, 100, 5, "acc-short")

	status, debit := doJSON(t, http.MethodPost, ts.URL+"/v1/users/"+user+"/debits",
		map[string]any{"amount": 60, "idempotency_key": "debit-1"})
	if status != http.StatusOK {
		t.Fatalf("debit: status=%d body=%v", status, debit)
	}

	resp, err := http.Get(ts.URL + "/v1/users/" + user + "/lots")
	if err != nil {
		t.Fatalf("get lots: %v", err)
	}
	defer resp.Body.Close()
	var lotList []data.LotInfo
	if err := json.NewDecoder(resp.Body).Decode(&lotList); err != nil {
		t.Fatalf("decode lots: %v", err)
	}

	var shortLotRemaining, longLotRemaining int
	for _, l := range lotList {
		if l.Amount == 100 {
			daysToExpiry := l.ExpiresAt.Sub(l.CreatedAt).Hours() / 24
			if daysToExpiry < 30 {
				shortLotRemaining = l.Remaining
			} else {
				longLotRemaining = l.Remaining
			}
		}
	}

	if shortLotRemaining != 40 {
		t.Fatalf("expected the soon-to-expire lot to be drawn down to 40, got %d", shortLotRemaining)
	}
	if longLotRemaining != 100 {
		t.Fatalf("expected the long-lived lot to be untouched (100), got %d", longLotRemaining)
	}
}

// TestConcurrentHolds fires many concurrent hold requests for the same user
// and asserts that the database-level row locking (SELECT ... FOR UPDATE)
// prevents the same points from being held twice: exactly enough holds
// should succeed to exhaust the balance, and the rest must fail with 409.
func TestConcurrentHolds(t *testing.T) {
	ts := newTestServer(t)
	user := "user_concurrent"

	mustAccrue(t, ts, user, 1000, 30, "acc-1")

	const workers = 20
	const holdAmount = 100 // exactly 10 of these should succeed (1000/100)

	var wg sync.WaitGroup
	results := make([]int, workers)
	for i := 0; i < workers; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			status, _ := doJSON(t, http.MethodPost, ts.URL+"/v1/users/"+user+"/holds",
				map[string]any{"amount": holdAmount, "idempotency_key": fmt.Sprintf("hold-conc-%d", i)})
			results[i] = status
		}(i)
	}
	wg.Wait()

	succeeded := 0
	for _, s := range results {
		if s == http.StatusCreated {
			succeeded++
		} else if s != http.StatusConflict {
			t.Fatalf("unexpected status %d", s)
		}
	}
	if succeeded != 10 {
		t.Fatalf("expected exactly 10 successful holds, got %d", succeeded)
	}

	status, bal := doJSON(t, http.MethodGet, ts.URL+"/v1/users/"+user+"/balance", nil)
	if status != http.StatusOK || bal["available"].(float64) != 0 || bal["held"].(float64) != 1000 {
		t.Fatalf("balance after concurrent holds: status=%d body=%v", status, bal)
	}
}

func mustAccrue(t *testing.T, ts *httptest.Server, user string, amount, ttlDays int, idem string) {
	t.Helper()
	status, body := doJSON(t, http.MethodPost, ts.URL+"/v1/users/"+user+"/accruals",
		map[string]any{"amount": amount, "ttl_days": ttlDays, "idempotency_key": idem})
	if status != http.StatusCreated {
		t.Fatalf("accrue: status=%d body=%v", status, body)
	}
}
