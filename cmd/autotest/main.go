package main

import (
	"bufio"
	"context"
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"

	_ "github.com/lib/pq"

	"bonus-ledger/internal/autotest"
	"bonus-ledger/internal/data"
)

const defaultBaseURL = "http://localhost:8080"

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

func main() {
	var cfg config

	flag.StringVar(&cfg.dsn, "db-dsn", os.Getenv("DB_DSN"), "PostgreSQL DSN used to store autotest scenarios")
	flag.StringVar(&cfg.baseURL, "base-url", envString("AUTOTEST_BASE_URL", defaultBaseURL), "Base URL of the running API started from cmd/api")
	flag.IntVar(&cfg.defaultTTL, "default-ttl-days", envInt("DEFAULT_TTL_DAYS", autotest.DefaultTTLDays), "default TTL for accrual requests that do not override ttl_days")
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
		t.printf("4. Run multi-key parallel accrual test (US-19)\n")
		t.printf("5. Run generated/stored test scenarios\n")
		t.printf("6. Exit\n")

		choice, err := t.promptInt("Choose an option", 6, 1)
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
			t.runAndReport(autotest.Check{Name: "accrual correctness", Run: autotest.RunAccrualCorrectness}, scn)
		case 3:
			scn, ok, err := t.selectScenario()
			if err != nil {
				return err
			}
			if !ok {
				continue
			}
			t.runAndReport(autotest.Check{Name: "parallel accrual", Run: autotest.RunParallelAccrual}, scn)
		case 4:
			scn, ok, err := t.selectScenario()
			if err != nil {
				return err
			}
			if !ok {
				continue
			}
			t.runAndReport(autotest.Check{Name: "multi-key parallel accrual", Run: autotest.RunMultiKeyParallelAccrual}, scn)
		case 5:
			scenarios, err := t.loadScenarios()
			if err != nil {
				return err
			}
			if len(scenarios) == 0 {
				t.printf("No stored scenarios yet. Create one first.\n")
				continue
			}
			for _, scn := range scenarios {
				for _, check := range autotest.Checks() {
					t.runAndReport(check, scn)
				}
			}
		case 6:
			t.printf("Bye.\n")
			return nil
		default:
			t.printf("Unknown choice: %d\n", choice)
		}
	}
}

func (t *tool) runAndReport(check autotest.Check, scn data.AutotestScenario) {
	t.printf("\nScenario %q (%s) for user %q\n", scn.Label, check.Name, scn.UserID)
	err := t.withRuntime(func(rt *autotest.Runtime) error {
		return check.Run(rt, scn)
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

func (t *tool) withRuntime(fn func(*autotest.Runtime) error) error {
	rt := autotest.NewRuntime(t.cfg.baseURL)
	if rt.BaseURL == "" {
		return errors.New("base URL is required; set AUTOTEST_BASE_URL or pass -base-url")
	}
	if err := rt.CheckHealth(); err != nil {
		return fmt.Errorf("%w. Start it with `go run ./cmd/api` or `docker compose up`", err)
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
	if err := autotest.Validate(scn); err != nil {
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
	parallelRequests, err := t.promptInt("Parallel requests", autotest.DefaultParallelRequests, 2)
	if err != nil {
		return data.AutotestScenario{}, err
	}

	scn := autotest.NormalizeScenario(data.AutotestScenario{
		Label:            label,
		UserID:           userInput,
		Amount:           amount,
		TTLDays:          ttlDays,
		ParallelRequests: parallelRequests,
	})
	if err := autotest.Validate(scn); err != nil {
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
	if err := autotest.Validate(scn); err != nil {
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
