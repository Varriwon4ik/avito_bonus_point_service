// Package autotest holds the reusable engine behind the US-15 autotester: it
// drives a running Bonus Points API over HTTP and asserts that accrual and
// parallel/concurrent accrual behave correctly. Both the cmd/autotest console
// tool and the web "Autotester" tab (US-17) share this engine so the two front
// ends always run identical checks.
package autotest

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"

	"bonus-ledger/internal/data"
)

const (
	// DefaultLedgerLabel is the ledger note the autotester writes and asserts on.
	DefaultLedgerLabel = "test"
	// DefaultTTLDays is the fallback lifetime for scenarios that omit ttl_days.
	DefaultTTLDays = 365
	// DefaultParallelRequests is the fallback fan-out for the parallel test.
	DefaultParallelRequests = 5
	// UserIDPrefix guards the autotester against touching real user accounts.
	UserIDPrefix = "autotest-"
)

type errorEnvelope struct {
	Error   string `json:"error"`
	Message string `json:"message"`
}

type userState struct {
	Balance     data.BalanceResult
	LedgerTotal int
	Lots        []data.LotInfo
}

type parallelResult struct {
	Status int
	Body   string
	LotID  int64
}

// Runtime executes autotest scenarios against a running API instance.
type Runtime struct {
	BaseURL string
	Client  *http.Client
}

// NewRuntime builds a Runtime pointed at baseURL with a sane request timeout.
func NewRuntime(baseURL string) *Runtime {
	return &Runtime{
		BaseURL: strings.TrimRight(strings.TrimSpace(baseURL), "/"),
		Client: &http.Client{
			Timeout: 15 * time.Second,
		},
	}
}

// CheckHealth verifies the target API is reachable and healthy before a run.
func (rt *Runtime) CheckHealth() error {
	if rt.BaseURL == "" {
		return errors.New("base URL is required")
	}

	var health map[string]string
	status, apiErr, err := rt.doJSON(http.MethodGet, "/healthz", nil, &health)
	if err != nil {
		return fmt.Errorf("cannot reach application at %s: %w", rt.BaseURL, err)
	}
	if status != http.StatusOK {
		return fmt.Errorf("health check on %s returned status %d: %s %s", rt.BaseURL, status, apiErr.Error, apiErr.Message)
	}
	return nil
}

// A Check is a single named autotest that can pass or fail against a scenario.
type Check struct {
	Name string
	Run  func(*Runtime, data.AutotestScenario) error
}

// Checks returns the ordered set of checks the autotester runs, matching the
// options exposed by the cmd/autotest console tool.
func Checks() []Check {
	return []Check{
		{Name: "accrual correctness", Run: RunAccrualCorrectness},
		{Name: "parallel accrual", Run: RunParallelAccrual},
	}
}

// Validate rejects scenarios that would touch real accounts or exercise the
// API with nonsensical inputs.
func Validate(scn data.AutotestScenario) error {
	switch {
	case scn.Label == "":
		return errors.New("label is required")
	case !strings.HasPrefix(scn.UserID, UserIDPrefix):
		return errors.New("user_id must start with " + UserIDPrefix)
	case scn.Amount <= 0:
		return errors.New("amount must be a positive integer")
	case scn.TTLDays <= 0:
		return errors.New("ttl_days must be a positive integer")
	case scn.ParallelRequests < 2:
		return errors.New("parallel_requests must be at least 2")
	case scn.LedgerLabel != DefaultLedgerLabel:
		return fmt.Errorf("ledger_label must be %q", DefaultLedgerLabel)
	default:
		return nil
	}
}

// NormalizeScenario applies the defaults and slug rules the autotester relies
// on so callers can pass loosely-typed, user-entered values.
func NormalizeScenario(scn data.AutotestScenario) data.AutotestScenario {
	scn.Label = NormalizeLabel(scn.Label)
	scn.UserID = NormalizeTestUserID(scn.UserID)
	if scn.TTLDays <= 0 {
		scn.TTLDays = DefaultTTLDays
	}
	if scn.ParallelRequests < 2 {
		scn.ParallelRequests = DefaultParallelRequests
	}
	scn.LedgerLabel = DefaultLedgerLabel
	return scn
}

// NormalizeLabel turns free text into a stable scenario slug.
func NormalizeLabel(value string) string {
	token := sanitizeToken(value)
	if token == "" {
		return "demo"
	}
	return token
}

// NormalizeTestUserID guarantees the UserIDPrefix and a usable slug suffix.
func NormalizeTestUserID(value string) string {
	token := sanitizeToken(strings.TrimPrefix(strings.TrimSpace(strings.ToLower(value)), UserIDPrefix))
	if token == "" {
		token = "demo-user"
	}
	return UserIDPrefix + token
}

func sanitizeToken(value string) string {
	var b strings.Builder
	lastDash := false
	for _, r := range strings.ToLower(strings.TrimSpace(value)) {
		isAlphaNum := (r >= 'a' && r <= 'z') || (r >= '0' && r <= '9')
		if isAlphaNum {
			b.WriteRune(r)
			lastDash = false
			continue
		}
		if !lastDash {
			b.WriteByte('-')
			lastDash = true
		}
	}
	return strings.Trim(b.String(), "-")
}

// RunAccrualCorrectness accrues once and asserts the balance, lot and ledger
// all move consistently, then asserts an invalid accrual is rejected without
// side effects.
func RunAccrualCorrectness(rt *Runtime, scn data.AutotestScenario) error {
	before, err := rt.loadUserState(scn.UserID)
	if err != nil {
		return err
	}

	key := scenarioKey(scn.Label, "accrual")
	var created data.AccrualResult
	status, apiErr, err := rt.doJSON(http.MethodPost, accrualPath(scn.UserID), map[string]any{
		"amount":          scn.Amount,
		"ttl_days":        scn.TTLDays,
		"idempotency_key": key,
		"label":           scn.LedgerLabel,
	}, &created)
	if err != nil {
		return err
	}
	if status != http.StatusCreated {
		return fmt.Errorf("accrual request returned status %d: %s %s", status, apiErr.Error, apiErr.Message)
	}
	if created.UserID != scn.UserID || created.Amount != scn.Amount || created.LotID <= 0 {
		return fmt.Errorf("unexpected accrual payload: %+v", created)
	}

	after, err := rt.loadUserState(scn.UserID)
	if err != nil {
		return err
	}

	if after.Balance.Available != before.Balance.Available+scn.Amount {
		return fmt.Errorf("available balance mismatch: want %d, got %d", before.Balance.Available+scn.Amount, after.Balance.Available)
	}
	if after.Balance.Held != before.Balance.Held {
		return fmt.Errorf("held balance changed unexpectedly: want %d, got %d", before.Balance.Held, after.Balance.Held)
	}
	if after.Balance.Total != before.Balance.Total+scn.Amount {
		return fmt.Errorf("total balance mismatch: want %d, got %d", before.Balance.Total+scn.Amount, after.Balance.Total)
	}
	if after.LedgerTotal != before.LedgerTotal+1 {
		return fmt.Errorf("ledger total mismatch: want %d, got %d", before.LedgerTotal+1, after.LedgerTotal)
	}
	if len(after.Lots) != len(before.Lots)+1 {
		return fmt.Errorf("lot count mismatch: want %d, got %d", len(before.Lots)+1, len(after.Lots))
	}

	lot, ok := findLot(after.Lots, created.LotID)
	if !ok {
		return fmt.Errorf("new lot %d was not found in /lots", created.LotID)
	}
	if lot.Amount != scn.Amount || lot.Remaining != scn.Amount {
		return fmt.Errorf("new lot %d has wrong values: amount=%d remaining=%d", created.LotID, lot.Amount, lot.Remaining)
	}

	ledgerEntries, err := rt.loadAllLedgerEntries(scn.UserID)
	if err != nil {
		return err
	}
	if !containsAccrualEntry(ledgerEntries, created.LotID, scn.Amount, scn.LedgerLabel) {
		return fmt.Errorf("transactions endpoint did not expose the accrual for lot %d with label %q", created.LotID, scn.LedgerLabel)
	}

	invalidBefore := after
	status, apiErr, err = rt.doJSON(http.MethodPost, accrualPath(scn.UserID), map[string]any{
		"amount":          0,
		"ttl_days":        scn.TTLDays,
		"idempotency_key": scenarioKey(scn.Label, "invalid"),
		"label":           scn.LedgerLabel,
	}, nil)
	if err != nil {
		return err
	}
	if status != http.StatusBadRequest {
		return fmt.Errorf("invalid accrual should return 400, got %d", status)
	}
	if apiErr.Message != data.ErrInvalidAmount.Error() {
		return fmt.Errorf("invalid accrual returned unexpected message: %q", apiErr.Message)
	}

	invalidAfter, err := rt.loadUserState(scn.UserID)
	if err != nil {
		return err
	}
	if !sameState(invalidBefore, invalidAfter) {
		return errors.New("invalid accrual changed ledger state")
	}

	return nil
}

// RunParallelAccrual fires ParallelRequests accruals simultaneously and asserts
// each produced a distinct lot and the aggregate balance/ledger totals are
// exactly right — the core "parallel request issues" check from US-17.
func RunParallelAccrual(rt *Runtime, scn data.AutotestScenario) error {
	before, err := rt.loadUserState(scn.UserID)
	if err != nil {
		return err
	}

	results := make([]parallelResult, scn.ParallelRequests)
	start := make(chan struct{})
	var wg sync.WaitGroup
	errCh := make(chan error, scn.ParallelRequests)

	runID := time.Now().UTC().UnixNano()
	for i := 0; i < scn.ParallelRequests; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			<-start

			var created data.AccrualResult
			status, apiErr, err := rt.doJSON(http.MethodPost, accrualPath(scn.UserID), map[string]any{
				"amount":          scn.Amount,
				"ttl_days":        scn.TTLDays,
				"idempotency_key": fmt.Sprintf("autotest-%s-parallel-%d-%d", scn.Label, runID, i),
				"label":           scn.LedgerLabel,
			}, &created)
			if err != nil {
				errCh <- err
				return
			}
			results[i] = parallelResult{
				Status: status,
				Body:   fmt.Sprintf("%s %s", apiErr.Error, apiErr.Message),
				LotID:  created.LotID,
			}
		}(i)
	}

	close(start)
	wg.Wait()
	close(errCh)

	for err := range errCh {
		if err != nil {
			return err
		}
	}

	uniqueLots := make(map[int64]struct{}, scn.ParallelRequests)
	for i, result := range results {
		if result.Status != http.StatusCreated {
			return fmt.Errorf("parallel request %d returned status %d: %s", i+1, result.Status, strings.TrimSpace(result.Body))
		}
		if result.LotID <= 0 {
			return fmt.Errorf("parallel request %d returned an invalid lot id", i+1)
		}
		if _, exists := uniqueLots[result.LotID]; exists {
			return fmt.Errorf("parallel request %d reused lot id %d", i+1, result.LotID)
		}
		uniqueLots[result.LotID] = struct{}{}
	}

	after, err := rt.loadUserState(scn.UserID)
	if err != nil {
		return err
	}

	expectedDelta := scn.Amount * scn.ParallelRequests
	if after.Balance.Available != before.Balance.Available+expectedDelta {
		return fmt.Errorf("available balance mismatch after parallel run: want %d, got %d", before.Balance.Available+expectedDelta, after.Balance.Available)
	}
	if after.Balance.Held != before.Balance.Held {
		return fmt.Errorf("held balance changed unexpectedly after parallel run: want %d, got %d", before.Balance.Held, after.Balance.Held)
	}
	if after.Balance.Total != before.Balance.Total+expectedDelta {
		return fmt.Errorf("total balance mismatch after parallel run: want %d, got %d", before.Balance.Total+expectedDelta, after.Balance.Total)
	}
	if after.LedgerTotal != before.LedgerTotal+scn.ParallelRequests {
		return fmt.Errorf("ledger total mismatch after parallel run: want %d, got %d", before.LedgerTotal+scn.ParallelRequests, after.LedgerTotal)
	}
	if len(after.Lots) != len(before.Lots)+scn.ParallelRequests {
		return fmt.Errorf("lot count mismatch after parallel run: want %d, got %d", len(before.Lots)+scn.ParallelRequests, len(after.Lots))
	}

	ledgerEntries, err := rt.loadAllLedgerEntries(scn.UserID)
	if err != nil {
		return err
	}
	for lotID := range uniqueLots {
		lot, ok := findLot(after.Lots, lotID)
		if !ok {
			return fmt.Errorf("parallel-created lot %d was not found in /lots", lotID)
		}
		if lot.Amount != scn.Amount || lot.Remaining != scn.Amount {
			return fmt.Errorf("parallel-created lot %d has wrong values: amount=%d remaining=%d", lotID, lot.Amount, lot.Remaining)
		}
		if !containsAccrualEntry(ledgerEntries, lotID, scn.Amount, scn.LedgerLabel) {
			return fmt.Errorf("transactions endpoint did not expose the parallel accrual for lot %d with label %q", lotID, scn.LedgerLabel)
		}
	}

	return nil
}

func accrualPath(userID string) string {
	return "/v1/users/" + url.PathEscape(userID) + "/accruals"
}

func scenarioKey(label, suffix string) string {
	return fmt.Sprintf("autotest-%s-%s-%d", label, suffix, time.Now().UTC().UnixNano())
}

func sameState(a, b userState) bool {
	return a.Balance.Available == b.Balance.Available &&
		a.Balance.Held == b.Balance.Held &&
		a.Balance.Total == b.Balance.Total &&
		a.Balance.ExpiringSoon == b.Balance.ExpiringSoon &&
		a.LedgerTotal == b.LedgerTotal &&
		len(a.Lots) == len(b.Lots)
}

func findLot(lots []data.LotInfo, lotID int64) (data.LotInfo, bool) {
	for _, lot := range lots {
		if lot.LotID == lotID {
			return lot, true
		}
	}
	return data.LotInfo{}, false
}

func containsAccrualEntry(entries []data.LedgerEntry, lotID int64, amount int, label string) bool {
	for _, entry := range entries {
		if entry.Type != "accrual" || entry.Amount != amount || entry.RefID == nil {
			continue
		}
		if *entry.RefID != lotID {
			continue
		}
		if label == "" {
			return true
		}
		if entry.Label != nil && *entry.Label == label {
			return true
		}
		if entry.Note != nil && *entry.Note == label {
			return true
		}
	}
	return false
}

func (rt *Runtime) loadUserState(userID string) (userState, error) {
	var state userState

	balance, exists, err := rt.loadBalance(userID)
	if err != nil {
		return state, err
	}
	if exists {
		state.Balance = balance
	} else {
		state.Balance = data.BalanceResult{UserID: userID}
	}

	lots, err := rt.loadLots(userID)
	if err != nil {
		return state, err
	}
	state.Lots = lots

	ledger, err := rt.loadLedgerPage(userID, 1, 1)
	if err != nil {
		return state, err
	}
	state.LedgerTotal = ledger.Total

	return state, nil
}

func (rt *Runtime) loadAllLedgerEntries(userID string) ([]data.LedgerEntry, error) {
	page := 1
	offset := 100
	var all []data.LedgerEntry

	for {
		ledger, err := rt.loadLedgerPage(userID, page, offset)
		if err != nil {
			return nil, err
		}
		all = append(all, ledger.Entries...)
		if len(all) >= ledger.Total || len(ledger.Entries) == 0 {
			return all, nil
		}
		page++
	}
}

func (rt *Runtime) loadBalance(userID string) (data.BalanceResult, bool, error) {
	var balance data.BalanceResult
	status, apiErr, err := rt.doJSON(http.MethodGet, "/v1/users/"+url.PathEscape(userID)+"/balance", nil, &balance)
	if err != nil {
		return balance, false, err
	}
	switch status {
	case http.StatusOK:
		return balance, true, nil
	case http.StatusNotFound:
		if apiErr.Message == data.ErrUserNotFound.Error() {
			return data.BalanceResult{UserID: userID}, false, nil
		}
		return balance, false, fmt.Errorf("balance lookup returned 404: %s", apiErr.Message)
	default:
		return balance, false, fmt.Errorf("balance lookup returned status %d: %s %s", status, apiErr.Error, apiErr.Message)
	}
}

func (rt *Runtime) loadLots(userID string) ([]data.LotInfo, error) {
	var lots []data.LotInfo
	status, apiErr, err := rt.doJSON(http.MethodGet, "/v1/users/"+url.PathEscape(userID)+"/lots", nil, &lots)
	if err != nil {
		return nil, err
	}
	switch status {
	case http.StatusOK:
		return lots, nil
	case http.StatusNotFound:
		if apiErr.Message == data.ErrUserNotFound.Error() {
			return nil, nil
		}
		return nil, fmt.Errorf("lots lookup returned 404: %s", apiErr.Message)
	default:
		return nil, fmt.Errorf("lots lookup returned status %d: %s %s", status, apiErr.Error, apiErr.Message)
	}
}

func (rt *Runtime) loadLedgerPage(userID string, page, offset int) (data.PaginatedLedger, error) {
	var ledger data.PaginatedLedger
	path := fmt.Sprintf("/v1/users/%s/transactions?page=%d&offset=%d", url.PathEscape(userID), page, offset)
	status, apiErr, err := rt.doJSON(http.MethodGet, path, nil, &ledger)
	if err != nil {
		return ledger, err
	}
	switch status {
	case http.StatusOK:
		return ledger, nil
	case http.StatusNotFound:
		if apiErr.Message == data.ErrUserNotFound.Error() {
			return data.PaginatedLedger{UserID: userID, Page: page, Offset: offset}, nil
		}
		return ledger, fmt.Errorf("transactions lookup returned 404: %s", apiErr.Message)
	default:
		return ledger, fmt.Errorf("transactions lookup returned status %d: %s %s", status, apiErr.Error, apiErr.Message)
	}
}

func (rt *Runtime) doJSON(method, path string, payload any, dest any) (int, errorEnvelope, error) {
	var body io.Reader
	if payload != nil {
		raw, err := json.Marshal(payload)
		if err != nil {
			return 0, errorEnvelope{}, err
		}
		body = bytes.NewReader(raw)
	}

	req, err := http.NewRequestWithContext(context.Background(), method, rt.BaseURL+path, body)
	if err != nil {
		return 0, errorEnvelope{}, err
	}
	if payload != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	resp, err := rt.Client.Do(req)
	if err != nil {
		return 0, errorEnvelope{}, err
	}
	defer resp.Body.Close()

	rawBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return 0, errorEnvelope{}, err
	}

	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		if dest != nil && len(rawBody) > 0 {
			if err := json.Unmarshal(rawBody, dest); err != nil {
				return 0, errorEnvelope{}, fmt.Errorf("decode success response: %w", err)
			}
		}
		return resp.StatusCode, errorEnvelope{}, nil
	}

	var apiErr errorEnvelope
	if len(rawBody) > 0 {
		if err := json.Unmarshal(rawBody, &apiErr); err != nil {
			apiErr.Message = strings.TrimSpace(string(rawBody))
		}
	}
	return resp.StatusCode, apiErr, nil
}
