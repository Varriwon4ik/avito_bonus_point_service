package api_test

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"testing"

	_ "github.com/lib/pq"

	"bonus-ledger/internal/api"
	"bonus-ledger/internal/data"
)

type testEnv struct {
	DB     *sql.DB
	Server *httptest.Server
}

// newTestEnv connects to TEST_DATABASE_URL, applies migrations and wipes all
// tables so each test starts from a clean slate. Tests are skipped if
// TEST_DATABASE_URL is not set.
func newTestEnv(t *testing.T) *testEnv {
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

	specPath := filepath.Join("..", "..", "api", "openapi.yaml")
	spec, err := os.ReadFile(specPath)
	if err != nil {
		t.Fatalf("read openapi spec: %v", err)
	}

	store := data.NewStore(db)
	logger := slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{Level: slog.LevelError}))
	os.Setenv("ADMIN_AUTH_TOKEN", testAdminToken)
	handler := api.NewAppHandler(api.NewServer(store, logger, 365), nil, spec)

	ts := httptest.NewServer(handler)
	t.Cleanup(ts.Close)

	return &testEnv{DB: db, Server: ts}
}

type httpResult struct {
	StatusCode int
	Header     http.Header
	Body       []byte
}

const testAdminToken = "test-admin-token"

func doRaw(t *testing.T, method, url, contentType string, body []byte) httpResult {
	return doRawWithHeaders(t, method, url, contentType, body, nil)
}

func doRawWithHeaders(t *testing.T, method, url, contentType string, body []byte, headers map[string]string) httpResult {
	t.Helper()

	var reader io.Reader
	if body != nil {
		reader = bytes.NewReader(body)
	}

	req, err := http.NewRequest(method, url, reader)
	if err != nil {
		t.Fatalf("new request: %v", err)
	}
	if contentType != "" {
		req.Header.Set("Content-Type", contentType)
	}
	for key, value := range headers {
		req.Header.Set(key, value)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("do request: %v", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("read response: %v", err)
	}

	return httpResult{
		StatusCode: resp.StatusCode,
		Header:     resp.Header.Clone(),
		Body:       respBody,
	}
}

func doJSON(t *testing.T, method, url string, payload any) (int, http.Header, map[string]any) {
	return doJSONWithHeaders(t, method, url, payload, nil)
}

func doJSONWithHeaders(t *testing.T, method, url string, payload any, headers map[string]string) (int, http.Header, map[string]any) {
	t.Helper()

	var body []byte
	contentType := ""
	if payload != nil {
		var err error
		body, err = json.Marshal(payload)
		if err != nil {
			t.Fatalf("marshal: %v", err)
		}
		contentType = "application/json"
	}

	resp := doRawWithHeaders(t, method, url, contentType, body, headers)
	return resp.StatusCode, resp.Header, decodeJSONMap(t, resp.Body)
}

func decodeJSONMap(t *testing.T, body []byte) map[string]any {
	t.Helper()

	var out map[string]any
	if err := json.Unmarshal(body, &out); err != nil {
		t.Fatalf("decode response: %v; body=%s", err, string(body))
	}
	return out
}

func assertJSONContentType(t *testing.T, header http.Header) {
	t.Helper()
	if got := header.Get("Content-Type"); !strings.Contains(got, "application/json") {
		t.Fatalf("expected application/json content type, got %q", got)
	}
}

func assertErrorResponse(t *testing.T, status int, header http.Header, body map[string]any, wantStatus int, wantCode, wantMessage string) {
	t.Helper()
	if status != wantStatus {
		t.Fatalf("expected status=%d, got %d body=%v", wantStatus, status, body)
	}
	assertJSONContentType(t, header)
	if body["error"] != wantCode {
		t.Fatalf("expected error=%q, got %v", wantCode, body["error"])
	}
	if body["message"] != wantMessage {
		t.Fatalf("expected message=%q, got %v", wantMessage, body["message"])
	}
}

func mustAccrue(t *testing.T, env *testEnv, user string, amount, ttlDays int, idem string) {
	t.Helper()
	status, header, body := doJSONWithHeaders(t, http.MethodPost, env.Server.URL+"/v1/users/"+user+"/accruals",
		map[string]any{"amount": amount, "ttl_days": ttlDays, "idempotency_key": idem}, map[string]string{
			"Authorization": "Bearer " + testAdminToken,
		})
	if status != http.StatusCreated {
		t.Fatalf("accrue: status=%d body=%v", status, body)
	}
	assertJSONContentType(t, header)
}

func TestAdminAccrualRequiresAuthentication(t *testing.T) {
	env := newTestEnv(t)

	status, header, body := doJSON(t, http.MethodPost, env.Server.URL+"/v1/users/user_123/accruals",
		map[string]any{"amount": 100, "ttl_days": 30, "idempotency_key": "order-1"})
	assertErrorResponse(t, status, header, body, http.StatusUnauthorized, "unauthorized", "admin authentication required")
}

func TestAccrualIdempotency(t *testing.T) {
	env := newTestEnv(t)
	user := "user_accrual"

	payload := map[string]any{"amount": 100, "ttl_days": 30, "idempotency_key": "order-1"}

	status, header, body := doJSONWithHeaders(t, http.MethodPost, env.Server.URL+"/v1/users/"+user+"/accruals", payload, map[string]string{
		"Authorization": "Bearer " + testAdminToken,
	})
	if status != http.StatusCreated {
		t.Fatalf("first accrual: status=%d body=%v", status, body)
	}
	assertJSONContentType(t, header)
	firstLotID := body["lot_id"]

	status, header, body2 := doJSONWithHeaders(t, http.MethodPost, env.Server.URL+"/v1/users/"+user+"/accruals", payload, map[string]string{
		"Authorization": "Bearer " + testAdminToken,
	})
	if status != http.StatusCreated {
		t.Fatalf("second accrual: status=%d body=%v", status, body2)
	}
	assertJSONContentType(t, header)
	if body2["lot_id"] != firstLotID {
		t.Fatalf("expected idempotent replay to return same lot_id, got %v vs %v", firstLotID, body2["lot_id"])
	}

	status, header, bal := doJSON(t, http.MethodGet, env.Server.URL+"/v1/users/"+user+"/balance", nil)
	if status != http.StatusOK {
		t.Fatalf("balance: status=%d body=%v", status, bal)
	}
	assertJSONContentType(t, header)
	if bal["available"].(float64) != 100 {
		t.Fatalf("expected available=100 after idempotent retry, got %v", bal["available"])
	}
}

func TestHoldConfirmCancel(t *testing.T) {
	env := newTestEnv(t)
	user := "user_hold"

	mustAccrue(t, env, user, 1000, 30, "acc-1")

	status, header, hold := doJSON(t, http.MethodPost, env.Server.URL+"/v1/users/"+user+"/holds",
		map[string]any{"amount": 400, "idempotency_key": "hold-1"})
	if status != http.StatusCreated {
		t.Fatalf("create hold: status=%d body=%v", status, hold)
	}
	assertJSONContentType(t, header)
	holdID := int64(hold["hold_id"].(float64))

	status, _, bal := doJSON(t, http.MethodGet, env.Server.URL+"/v1/users/"+user+"/balance", nil)
	if status != http.StatusOK || bal["available"].(float64) != 600 || bal["held"].(float64) != 400 {
		t.Fatalf("balance after hold: status=%d body=%v", status, bal)
	}

	status, _, cancel := doJSON(t, http.MethodPost, fmt.Sprintf("%s/v1/holds/%d/cancel", env.Server.URL, holdID), nil)
	if status != http.StatusOK || cancel["status"] != "cancelled" {
		t.Fatalf("cancel hold: status=%d body=%v", status, cancel)
	}

	status, _, bal = doJSON(t, http.MethodGet, env.Server.URL+"/v1/users/"+user+"/balance", nil)
	if status != http.StatusOK || bal["available"].(float64) != 1000 || bal["held"].(float64) != 0 {
		t.Fatalf("balance after cancel: status=%d body=%v", status, bal)
	}

	status, _, hold2 := doJSON(t, http.MethodPost, env.Server.URL+"/v1/users/"+user+"/holds",
		map[string]any{"amount": 250, "idempotency_key": "hold-2"})
	if status != http.StatusCreated {
		t.Fatalf("create hold 2: status=%d body=%v", status, hold2)
	}
	holdID2 := int64(hold2["hold_id"].(float64))

	status, _, confirm := doJSON(t, http.MethodPost, fmt.Sprintf("%s/v1/holds/%d/confirm", env.Server.URL, holdID2), nil)
	if status != http.StatusOK || confirm["status"] != "confirmed" {
		t.Fatalf("confirm hold: status=%d body=%v", status, confirm)
	}

	status, _, confirm2 := doJSON(t, http.MethodPost, fmt.Sprintf("%s/v1/holds/%d/confirm", env.Server.URL, holdID2), nil)
	if status != http.StatusOK || confirm2["status"] != "confirmed" {
		t.Fatalf("re-confirm hold: status=%d body=%v", status, confirm2)
	}

	status, _, bal = doJSON(t, http.MethodGet, env.Server.URL+"/v1/users/"+user+"/balance", nil)
	if status != http.StatusOK || bal["available"].(float64) != 750 || bal["held"].(float64) != 0 {
		t.Fatalf("balance after confirm: status=%d body=%v", status, bal)
	}
}

func TestMalformedJSONReturnsBadRequest(t *testing.T) {
	env := newTestEnv(t)

	resp := doRawWithHeaders(t, http.MethodPost, env.Server.URL+"/v1/users/user_123/accruals", "application/json", []byte(`{"amount":`), map[string]string{"Authorization": "Bearer " + testAdminToken})
	body := decodeJSONMap(t, resp.Body)
	assertErrorResponse(t, resp.StatusCode, resp.Header, body, http.StatusBadRequest, "bad_request", "request body contains malformed JSON")
}

func TestMissingRequiredFieldReturnsBadRequest(t *testing.T) {
	env := newTestEnv(t)

	status, header, body := doJSONWithHeaders(t, http.MethodPost, env.Server.URL+"/v1/users/user_123/accruals",
		map[string]any{"amount": 100}, map[string]string{"Authorization": "Bearer " + testAdminToken})
	assertErrorResponse(t, status, header, body, http.StatusBadRequest, "bad_request", "idempotency_key is required")
}

func TestInvalidAmountReturnsBadRequest(t *testing.T) {
	for _, amount := range []int{0, -10} {
		t.Run(fmt.Sprintf("amount=%d", amount), func(t *testing.T) {
			env := newTestEnv(t)

			status, header, body := doJSONWithHeaders(t, http.MethodPost, env.Server.URL+"/v1/users/user_123/accruals",
				map[string]any{"amount": amount, "idempotency_key": "order-1"}, map[string]string{"Authorization": "Bearer " + testAdminToken})
			assertErrorResponse(t, status, header, body, http.StatusBadRequest, "bad_request", "amount must be a positive integer")
		})
	}
}

func TestInvalidPaginationReturnsBadRequest(t *testing.T) {
	env := newTestEnv(t)
	user := "user_pagination"
	mustAccrue(t, env, user, 100, 30, "acc-1")

	status, header, body := doJSON(t, http.MethodGet, env.Server.URL+"/v1/users/"+user+"/transactions?offset=999", nil)
	assertErrorResponse(t, status, header, body, http.StatusBadRequest, "bad_request", "offset must be between 1 and 500")
}

func TestUnknownUserReturnsNotFound(t *testing.T) {
	env := newTestEnv(t)

	t.Run("balance", func(t *testing.T) {
		status, header, body := doJSON(t, http.MethodGet, env.Server.URL+"/v1/users/missing-user/balance", nil)
		assertErrorResponse(t, status, header, body, http.StatusNotFound, "not_found", "user not found")
	})

	t.Run("hold", func(t *testing.T) {
		status, header, body := doJSON(t, http.MethodPost, env.Server.URL+"/v1/users/missing-user/holds",
			map[string]any{"amount": 50, "idempotency_key": "hold-1"})
		assertErrorResponse(t, status, header, body, http.StatusNotFound, "not_found", "user not found")
	})
}

func TestUnknownHoldReturnsNotFound(t *testing.T) {
	env := newTestEnv(t)

	status, header, body := doJSON(t, http.MethodPost, env.Server.URL+"/v1/holds/999/confirm", nil)
	assertErrorResponse(t, status, header, body, http.StatusNotFound, "not_found", "hold not found")
}

func TestDuplicateIdempotencyKeyReturnsConflict(t *testing.T) {
	env := newTestEnv(t)

	if _, err := env.DB.Exec(`INSERT INTO idempotency_keys (key, endpoint) VALUES ($1, $2)`, "order-1", "accrual"); err != nil {
		t.Fatalf("seed idempotency key: %v", err)
	}

	status, header, body := doJSONWithHeaders(t, http.MethodPost, env.Server.URL+"/v1/users/user_123/accruals",
		map[string]any{"amount": 100, "ttl_days": 30, "idempotency_key": "order-1"}, map[string]string{"Authorization": "Bearer " + testAdminToken})
	assertErrorResponse(t, status, header, body, http.StatusConflict, "conflict", "a request with this idempotency key is already in progress")
}

func TestOpenAPIRouteExists(t *testing.T) {
	env := newTestEnv(t)

	resp := doRaw(t, http.MethodGet, env.Server.URL+"/openapi.yaml", "", nil)
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected status=200, got %d body=%s", resp.StatusCode, string(resp.Body))
	}
	if got := resp.Header.Get("Content-Type"); !strings.Contains(got, "application/yaml") {
		t.Fatalf("expected yaml content type, got %q", got)
	}
	if !strings.Contains(string(resp.Body), "openapi: 3.0.3") {
		t.Fatalf("expected openapi document body, got %s", string(resp.Body))
	}

	docsResp := doRaw(t, http.MethodGet, env.Server.URL+"/docs", "", nil)
	if docsResp.StatusCode != http.StatusOK {
		t.Fatalf("expected docs status=200, got %d body=%s", docsResp.StatusCode, string(docsResp.Body))
	}
	if got := docsResp.Header.Get("Content-Type"); !strings.Contains(got, "text/html") {
		t.Fatalf("expected html content type, got %q", got)
	}
}

func TestInsufficientFunds(t *testing.T) {
	env := newTestEnv(t)
	user := "user_poor"

	mustAccrue(t, env, user, 50, 30, "acc-1")

	status, header, body := doJSON(t, http.MethodPost, env.Server.URL+"/v1/users/"+user+"/holds",
		map[string]any{"amount": 100, "idempotency_key": "hold-1"})
	assertErrorResponse(t, status, header, body, http.StatusConflict, "conflict", "insufficient available points")
}

func TestExpiryOrdering(t *testing.T) {
	env := newTestEnv(t)
	user := "user_fifo"

	mustAccrue(t, env, user, 100, 60, "acc-long")
	mustAccrue(t, env, user, 100, 5, "acc-short")

	status, _, debit := doJSON(t, http.MethodPost, env.Server.URL+"/v1/users/"+user+"/debits",
		map[string]any{"amount": 60, "idempotency_key": "debit-1"})
	if status != http.StatusOK {
		t.Fatalf("debit: status=%d body=%v", status, debit)
	}

	resp, err := http.Get(env.Server.URL + "/v1/users/" + user + "/lots")
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

func TestConcurrentHolds(t *testing.T) {
	env := newTestEnv(t)
	user := "user_concurrent"

	mustAccrue(t, env, user, 1000, 30, "acc-1")

	const workers = 20
	const holdAmount = 100

	var wg sync.WaitGroup
	results := make([]int, workers)
	errs := make(chan error, workers)
	for i := 0; i < workers; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			payload, err := json.Marshal(map[string]any{"amount": holdAmount, "idempotency_key": fmt.Sprintf("hold-conc-%d", i)})
			if err != nil {
				errs <- err
				return
			}

			req, err := http.NewRequest(http.MethodPost, env.Server.URL+"/v1/users/"+user+"/holds", bytes.NewReader(payload))
			if err != nil {
				errs <- err
				return
			}
			req.Header.Set("Content-Type", "application/json")

			resp, err := http.DefaultClient.Do(req)
			if err != nil {
				errs <- err
				return
			}
			defer resp.Body.Close()

			results[i] = resp.StatusCode
		}(i)
	}
	wg.Wait()
	close(errs)

	for err := range errs {
		if err != nil {
			t.Fatalf("concurrent hold request failed: %v", err)
		}
	}

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

	status, _, bal := doJSON(t, http.MethodGet, env.Server.URL+"/v1/users/"+user+"/balance", nil)
	if status != http.StatusOK || bal["available"].(float64) != 0 || bal["held"].(float64) != 1000 {
		t.Fatalf("balance after concurrent holds: status=%d body=%v", status, bal)
	}
}
