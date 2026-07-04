package main

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	_ "github.com/lib/pq"

	"bonus-ledger/internal/data"
)

const (
	defaultBaseURL     = "http://localhost:8080"
	defaultLedgerLabel = "test"
	defaultTTLDays     = 365
)

type config struct {
	dsn        string
	baseURL    string
	defaultTTL int
}

type tool struct {
	cfg config
	in  *bufio.Reader
	out io.Writer
}

type runtimeEnv struct {
	baseURL string
	client  *http.Client
}

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

func main() {
	var cfg config

	flag.StringVar(&cfg.dsn, "db-dsn", os.Getenv("DB_DSN"), "PostgreSQL DSN used to store autotest scenarios")
	flag.StringVar(&cfg.baseURL, "base-url", envString("AUTOTEST_BASE_URL", defaultBaseURL), "Base URL of the running API started from cmd/api")
	flag.IntVar(&cfg.defaultTTL, "default-ttl-days", envInt("DEFAULT_TTL_DAYS", defaultTTLDays), "default TTL for accrual requests that do not override ttl_days")
	flag.Parse()

	t := tool{
		cfg: cfg,
		in:  bufio.NewReader(os.Stdin),
		out: os.Stdout,
	}
	if err := t.run(); err != nil && !errors.Is(err, io.EOF) {
		fmt.Fprintf(os.Stderr, "autotest error: %v\n", err)
		os.Exit(1)
	}
}

func envString(name, fallback string) string {
	value := strings.TrimSpace(os.Getenv(name))
	if value == "" {
		return fallback
	}
	return value
}

func envInt(name string, fallback int) int {
	value := strings.TrimSpace(os.Getenv(name))
	if value == "" {
		return fallback
	}
	n, err := strconv.Atoi(value)
	if err != nil || n <= 0 {
		return fallback
	}
	return n
}

func (t *tool) run() error {
	for {
		t.printf("\nBonus Points Accrual Autotest\n")
		t.printf("1. Create a new accrual test scenario\n")
		t.printf("2. Run accrual service correctness test\n")
		t.printf("3. Run parallel transaction test\n")
		t.printf("4. Run generated/stored test scenarios\n")
		t.printf("5. Exit\n")

		choice, err := t.promptInt("Choose an option", 5, 1)
		if err != nil {
			return err
		}

		switch choice {
		case 1:
			scn, err := t.promptScenario()
			if err != nil {
				return err
			}
			stored, err := t.saveScenario(scn)
			if err != nil {
				return err
			}
			t.printf("Stored scenario %q in database for user %q with label %q\n", stored.Label, stored.UserID, stored.LedgerLabel)
		case 2:
			scn, ok, err := t.selectScenario()
			if err != nil {
				return err
			}
			if !ok {
				continue
			}
			t.runAndReport("accrual correctness", scn, runAccrualCorrectness)
		case 3:
			scn, ok, err := t.selectScenario()
			if err != nil {
				return err
			}
			if !ok {
				continue
			}
			t.runAndReport("parallel accrual", scn, runParallelAccrual)
		case 4:
			scenarios, err := t.loadScenarios()
			if err != nil {
				return err
			}
			if len(scenarios) == 0 {
				t.printf("No stored scenarios yet. Create one first.\n")
				continue
			}
			for _, scn := range scenarios {
				t.runAndReport("accrual correctness", scn, runAccrualCorrectness)
				t.runAndReport("parallel accrual", scn, runParallelAccrual)
			}
		case 5:
			t.printf("Bye.\n")
			return nil
		default:
			t.printf("Unknown choice: %d\n", choice)
		}
	}
}

func (t *tool) runAndReport(name string, scn data.AutotestScenario, fn func(*runtimeEnv, data.AutotestScenario) error) {
	t.printf("\nScenario %q (%s) for user %q\n", scn.Label, name, scn.UserID)
	err := t.withRuntime(func(rt *runtimeEnv) error {
		return fn(rt, scn)
	})
	if err != nil {
		t.printf("FAILED: %v\n", err)
		return
	}
	t.printf("PASSED\n")
}

func (t *tool) withStore(fn func(*data.Store) error) error {
	if strings.TrimSpace(t.cfg.dsn) == "" {
		return errors.New("DB_DSN is required to store/load autotest scenarios in the database")
	}

	db, err := data.OpenDB(t.cfg.dsn)
	if err != nil {
		return fmt.Errorf("open database: %w", err)
	}
	defer db.Close()

	if err := data.Migrate(db); err != nil {
		return fmt.Errorf("run migrations: %w", err)
	}

	return fn(data.NewStore(db))
}

func (t *tool) withRuntime(fn func(*runtimeEnv) error) error {
	baseURL := strings.TrimRight(strings.TrimSpace(t.cfg.baseURL), "/")
	if baseURL == "" {
		return errors.New("base URL is required; set AUTOTEST_BASE_URL or pass -base-url")
	}

	rt := &runtimeEnv{
		baseURL: baseURL,
		client: &http.Client{
			Timeout: 15 * time.Second,
		},
	}

	var health map[string]string
	status, apiErr, err := rt.doJSON(http.MethodGet, "/healthz", nil, &health)
	if err != nil {
		return fmt.Errorf("cannot reach application at %s: %w. Start it with `go run ./cmd/api` or `docker compose up`", baseURL, err)
	}
	if status != http.StatusOK {
		return fmt.Errorf("health check on %s returned status %d: %s %s", baseURL, status, apiErr.Error, apiErr.Message)
	}

	return fn(rt)
}

func (t *tool) loadScenarios() ([]data.AutotestScenario, error) {
	var scenarios []data.AutotestScenario
	err := t.withStore(func(store *data.Store) error {
		var err error
		scenarios, err = store.ListAutotestScenarios(context.Background())
		return err
	})
	if err != nil {
		return nil, err
	}
	return scenarios, nil
}

func (t *tool) saveScenario(scn data.AutotestScenario) (data.AutotestScenario, error) {
	if err := validateScenario(scn); err != nil {
		return data.AutotestScenario{}, err
	}

	var stored data.AutotestScenario
	err := t.withStore(func(store *data.Store) error {
		var err error
		stored, err = store.UpsertAutotestScenario(context.Background(), scn)
		return err
	})
	if err != nil {
		return data.AutotestScenario{}, err
	}
	return stored, nil
}

func (t *tool) promptScenario() (data.AutotestScenario, error) {
	label, err := t.promptString("Scenario label", "demo")
	if err != nil {
		return data.AutotestScenario{}, err
	}
	userInput, err := t.promptString("Test user ID or suffix", "demo-user")
	if err != nil {
		return data.AutotestScenario{}, err
	}
	amount, err := t.promptInt("Accrual amount", 100, 1)
	if err != nil {
		return data.AutotestScenario{}, err
	}
	ttlDays, err := t.promptInt("TTL days", t.cfg.defaultTTL, 1)
	if err != nil {
		return data.AutotestScenario{}, err
	}
	parallelRequests, err := t.promptInt("Parallel requests", 5, 2)
	if err != nil {
		return data.AutotestScenario{}, err
	}

	scn := data.AutotestScenario{
		Label:            normalizeLabel(label),
		UserID:           normalizeTestUserID(userInput),
		Amount:           amount,
		TTLDays:          ttlDays,
		ParallelRequests: parallelRequests,
		LedgerLabel:      defaultLedgerLabel,
	}
	if err := validateScenario(scn); err != nil {
		return data.AutotestScenario{}, err
	}
	return scn, nil
}

func (t *tool) selectScenario() (data.AutotestScenario, bool, error) {
	scenarios, err := t.loadScenarios()
	if err != nil {
		return data.AutotestScenario{}, false, err
	}
	if len(scenarios) == 0 {
		t.printf("No stored scenarios yet. Create one first.\n")
		return data.AutotestScenario{}, false, nil
	}

	t.printf("\nStored scenarios:\n")
	for i, scn := range scenarios {
		t.printf("%d. %s | user=%s | amount=%d | ttl=%d | parallel=%d | label=%s\n",
			i+1, scn.Label, scn.UserID, scn.Amount, scn.TTLDays, scn.ParallelRequests, scn.LedgerLabel)
	}

	index, err := t.promptInt("Select scenario number", 1, 1)
	if err != nil {
		return data.AutotestScenario{}, false, err
	}
	if index > len(scenarios) {
		t.printf("Scenario %d does not exist.\n", index)
		return data.AutotestScenario{}, false, nil
	}

	scn := scenarios[index-1]
	if err := validateScenario(scn); err != nil {
		return data.AutotestScenario{}, false, fmt.Errorf("stored scenario %q is invalid: %w", scn.Label, err)
	}
	return scn, true, nil
}

func (t *tool) promptString(label, defaultValue string) (string, error) {
	for {
		t.printf("%s [%s]: ", label, defaultValue)
		line, err := t.readLine()
		if err != nil {
			return "", err
		}
		if line == "" {
			return defaultValue, nil
		}
		return line, nil
	}
}

func (t *tool) promptInt(label string, defaultValue, min int) (int, error) {
	for {
		t.printf("%s [%d]: ", label, defaultValue)
		line, err := t.readLine()
		if err != nil {
			return 0, err
		}
		if line == "" {
			return defaultValue, nil
		}
		n, err := strconv.Atoi(line)
		if err != nil || n < min {
			t.printf("Enter an integer >= %d.\n", min)
			continue
		}
		return n, nil
	}
}

func (t *tool) readLine() (string, error) {
	line, err := t.in.ReadString('\n')
	if err != nil {
		if errors.Is(err, io.EOF) && line != "" {
			return strings.TrimSpace(line), nil
		}
		return "", err
	}
	return strings.TrimSpace(line), nil
}

func (t *tool) printf(format string, args ...any) {
	fmt.Fprintf(t.out, format, args...)
}

func validateScenario(scn data.AutotestScenario) error {
	switch {
	case scn.Label == "":
		return errors.New("label is required")
	case !strings.HasPrefix(scn.UserID, "autotest-"):
		return errors.New("user_id must start with autotest-")
	case scn.Amount <= 0:
		return errors.New("amount must be a positive integer")
	case scn.TTLDays <= 0:
		return errors.New("ttl_days must be a positive integer")
	case scn.ParallelRequests < 2:
		return errors.New("parallel_requests must be at least 2")
	case scn.LedgerLabel != defaultLedgerLabel:
		return fmt.Errorf("ledger_label must be %q", defaultLedgerLabel)
	default:
		return nil
	}
}

func normalizeLabel(value string) string {
	token := sanitizeToken(value)
	if token == "" {
		return "demo"
	}
	return token
}

func normalizeTestUserID(value string) string {
	token := sanitizeToken(strings.TrimPrefix(strings.TrimSpace(strings.ToLower(value)), "autotest-"))
	if token == "" {
		token = "demo-user"
	}
	return "autotest-" + token
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

func runAccrualCorrectness(rt *runtimeEnv, scn data.AutotestScenario) error {
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

func runParallelAccrual(rt *runtimeEnv, scn data.AutotestScenario) error {
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
		if entry.Note != nil && *entry.Note == label {
			return true
		}
	}
	return false
}

func (rt *runtimeEnv) loadUserState(userID string) (userState, error) {
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

func (rt *runtimeEnv) loadAllLedgerEntries(userID string) ([]data.LedgerEntry, error) {
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

func (rt *runtimeEnv) loadBalance(userID string) (data.BalanceResult, bool, error) {
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

func (rt *runtimeEnv) loadLots(userID string) ([]data.LotInfo, error) {
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

func (rt *runtimeEnv) loadLedgerPage(userID string, page, offset int) (data.PaginatedLedger, error) {
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

func (rt *runtimeEnv) doJSON(method, path string, payload any, dest any) (int, errorEnvelope, error) {
	var body io.Reader
	if payload != nil {
		raw, err := json.Marshal(payload)
		if err != nil {
			return 0, errorEnvelope{}, err
		}
		body = bytes.NewReader(raw)
	}

	req, err := http.NewRequest(method, rt.baseURL+path, body)
	if err != nil {
		return 0, errorEnvelope{}, err
	}
	if payload != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	resp, err := rt.client.Do(req)
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
