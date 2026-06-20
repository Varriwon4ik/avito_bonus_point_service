package data_test

import (
	"context"
	"encoding/json"
	"os"
	"testing"
	"time"

	_ "github.com/lib/pq"

	"bonus-ledger/internal/data"
)

// newTestStore connects to TEST_DATABASE_URL, applies migrations and wipes
// all tables so each test starts from a clean slate. Tests are skipped if
// TEST_DATABASE_URL is not set.
func newTestStore(t *testing.T) *data.Store {
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

	return data.NewStore(db)
}

// backdateHold sets a hold's created_at into the past, simulating a hold
// that has been sitting unresolved for a while.
func backdateHold(t *testing.T, store *data.Store, holdID int64, age time.Duration) {
	t.Helper()
	if _, err := store.DB.Exec(
		`UPDATE holds SET created_at = now() - make_interval(secs => $1) WHERE id = $2`,
		age.Seconds(), holdID); err != nil {
		t.Fatalf("backdate hold: %v", err)
	}
}

// TestExpireStaleHolds verifies AC1-AC3 of the hold timeout sweep: holds
// older than the configured timeout are auto-released and restore the
// user's available balance, a ledger entry documents the release, and holds
// younger than the timeout are left untouched.
func TestExpireStaleHolds(t *testing.T) {
	store := newTestStore(t)
	ctx := context.Background()
	user := "user_timeout"

	if _, _, err := store.Accrue(ctx, user, 1000, 30, "acc-1"); err != nil {
		t.Fatalf("accrue: %v", err)
	}

	_, staleBody, err := store.CreateHold(ctx, user, 300, "hold-stale")
	if err != nil {
		t.Fatalf("create stale hold: %v", err)
	}
	staleID := mustHoldID(t, staleBody)
	backdateHold(t, store, staleID, 48*time.Hour)

	_, freshBody, err := store.CreateHold(ctx, user, 200, "hold-fresh")
	if err != nil {
		t.Fatalf("create fresh hold: %v", err)
	}
	freshID := mustHoldID(t, freshBody)

	released, err := store.ExpireStaleHolds(ctx, 24)
	if err != nil {
		t.Fatalf("expire stale holds: %v", err)
	}
	if released != 1 {
		t.Fatalf("expected exactly 1 released hold, got %d", released)
	}

	bal, err := store.Balance(ctx, user, 7)
	if err != nil {
		t.Fatalf("balance: %v", err)
	}
	// stale hold (300) released back to available, fresh hold (200) still held
	if bal.Available != 800 || bal.Held != 200 {
		t.Fatalf("expected available=800 held=200, got available=%d held=%d", bal.Available, bal.Held)
	}

	entries, err := store.ListLedger(ctx, user, 10)
	if err != nil {
		t.Fatalf("list ledger: %v", err)
	}
	var found bool
	for _, e := range entries {
		if e.RefID != nil && *e.RefID == staleID && e.Type == "cancel" {
			if e.Note == nil || *e.Note != "auto-released: timeout" {
				t.Fatalf("expected ledger entry for hold %d to be annotated, got note=%v", staleID, e.Note)
			}
			found = true
		}
		if e.RefID != nil && *e.RefID == freshID && e.Type == "cancel" {
			t.Fatalf("fresh hold %d should not have been released", freshID)
		}
	}
	if !found {
		t.Fatalf("expected an annotated cancel ledger entry for stale hold %d", staleID)
	}

	// running the sweep again must be a no-op: the stale hold is already cancelled.
	released, err = store.ExpireStaleHolds(ctx, 24)
	if err != nil {
		t.Fatalf("expire stale holds (second run): %v", err)
	}
	if released != 0 {
		t.Fatalf("expected second sweep to release 0 holds, got %d", released)
	}
}

func mustHoldID(t *testing.T, body []byte) int64 {
	t.Helper()
	var res struct {
		HoldID int64 `json:"hold_id"`
	}
	if err := json.Unmarshal(body, &res); err != nil {
		t.Fatalf("unmarshal hold result: %v", err)
	}
	return res.HoldID
}
