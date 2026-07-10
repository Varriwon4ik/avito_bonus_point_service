package api_test

import (
	"net/http"
	"testing"
	"time"
)

type lotsFixture struct {
	activeID    int64
	exhaustedID int64
	expiredID   int64
}

func mustAccrueLotID(t *testing.T, env *testEnv, user string, amount, ttlDays int, idem string) int64 {
	t.Helper()

	status, header, body := doJSON(t, http.MethodPost, env.Server.URL+"/v1/users/"+user+"/accruals",
		map[string]any{"amount": amount, "ttl_days": ttlDays, "idempotency_key": idem})
	if status != http.StatusCreated {
		t.Fatalf("accrue: status=%d body=%v", status, body)
	}
	assertJSONContentType(t, header)

	return int64(body["lot_id"].(float64))
}

func mustDebit(t *testing.T, env *testEnv, user string, amount int, idem string) {
	t.Helper()

	status, header, body := doJSON(t, http.MethodPost, env.Server.URL+"/v1/users/"+user+"/debits",
		map[string]any{"amount": amount, "idempotency_key": idem})
	if status != http.StatusOK {
		t.Fatalf("debit: status=%d body=%v", status, body)
	}
	assertJSONContentType(t, header)
}

func setLotExpiry(t *testing.T, env *testEnv, lotID int64, expiresAt time.Time) {
	t.Helper()

	if _, err := env.DB.Exec(`UPDATE points_lots SET expires_at = $1 WHERE id = $2`, expiresAt, lotID); err != nil {
		t.Fatalf("set lot %d expiry: %v", lotID, err)
	}
}

func seedLotsForAudit(t *testing.T, env *testEnv, user string) lotsFixture {
	t.Helper()

	activeID := mustAccrueLotID(t, env, user, 100, 30, "lots-active")
	exhaustedID := mustAccrueLotID(t, env, user, 40, 30, "lots-exhausted")

	now := time.Now().UTC()
	setLotExpiry(t, env, activeID, now.Add(48*time.Hour))
	setLotExpiry(t, env, exhaustedID, now.Add(24*time.Hour))
	mustDebit(t, env, user, 40, "lots-exhausted-debit")

	expiredID := mustAccrueLotID(t, env, user, 70, 30, "lots-expired")
	setLotExpiry(t, env, expiredID, now.Add(-24*time.Hour))

	return lotsFixture{
		activeID:    activeID,
		exhaustedID: exhaustedID,
		expiredID:   expiredID,
	}
}

func responseLots(t *testing.T, body map[string]any) []map[string]any {
	t.Helper()

	rawLots, ok := body["lots"].([]any)
	if !ok {
		t.Fatalf("unexpected lots payload: %T", body["lots"])
	}

	lots := make([]map[string]any, 0, len(rawLots))
	for _, item := range rawLots {
		lot, ok := item.(map[string]any)
		if !ok {
			t.Fatalf("unexpected lot item payload: %T", item)
		}
		lots = append(lots, lot)
	}
	return lots
}

func TestListLotsReturnsPaginatedEnvelope(t *testing.T) {
	env := newTestEnv(t)
	user := "user_lots_audit"
	fixture := seedLotsForAudit(t, env, user)

	status, header, body := doJSON(t, http.MethodGet, env.Server.URL+"/v1/users/"+user+"/lots", nil)
	if status != http.StatusOK {
		t.Fatalf("want 200, got %d body=%v", status, body)
	}
	assertJSONContentType(t, header)

	if body["user_id"] != user {
		t.Fatalf("want user_id=%q, got %v", user, body["user_id"])
	}
	if body["page"].(float64) != 1 {
		t.Fatalf("want page=1, got %v", body["page"])
	}
	if body["offset"].(float64) != 20 {
		t.Fatalf("want offset=20, got %v", body["offset"])
	}
	if body["total"].(float64) != 3 {
		t.Fatalf("want total=3, got %v", body["total"])
	}

	lots := responseLots(t, body)
	if len(lots) != 3 {
		t.Fatalf("want 3 lots, got %d", len(lots))
	}

	want := []struct {
		lotID     int64
		status    string
		remaining float64
	}{
		{lotID: fixture.expiredID, status: "expired", remaining: 70},
		{lotID: fixture.exhaustedID, status: "exhausted", remaining: 0},
		{lotID: fixture.activeID, status: "active", remaining: 100},
	}

	for i, w := range want {
		if got := lots[i]["lot_id"].(float64); got != float64(w.lotID) {
			t.Fatalf("lot[%d]: want lot_id=%d, got %v", i, w.lotID, got)
		}
		if got := lots[i]["user_id"]; got != user {
			t.Fatalf("lot[%d]: want user_id=%q, got %v", i, user, got)
		}
		if got := lots[i]["status"]; got != w.status {
			t.Fatalf("lot[%d]: want status=%q, got %v", i, w.status, got)
		}
		if got := lots[i]["remaining"].(float64); got != w.remaining {
			t.Fatalf("lot[%d]: want remaining=%v, got %v", i, w.remaining, got)
		}
	}
}

func TestListLotsPagination(t *testing.T) {
	env := newTestEnv(t)
	user := "user_lots_pages"
	seedLotsForAudit(t, env, user)

	status, _, p1 := doJSON(t, http.MethodGet, env.Server.URL+"/v1/users/"+user+"/lots?page=1&offset=2", nil)
	if status != http.StatusOK {
		t.Fatalf("page1: want 200, got %d body=%v", status, p1)
	}
	status, _, p2 := doJSON(t, http.MethodGet, env.Server.URL+"/v1/users/"+user+"/lots?page=2&offset=2", nil)
	if status != http.StatusOK {
		t.Fatalf("page2: want 200, got %d body=%v", status, p2)
	}

	if p1["total"].(float64) != 3 || p2["total"].(float64) != 3 {
		t.Fatalf("want total=3 on both pages, got page1=%v page2=%v", p1["total"], p2["total"])
	}
	if p1["page"].(float64) != 1 || p2["page"].(float64) != 2 {
		t.Fatalf("unexpected page metadata: page1=%v page2=%v", p1["page"], p2["page"])
	}
	if p1["offset"].(float64) != 2 || p2["offset"].(float64) != 2 {
		t.Fatalf("unexpected offset metadata: page1=%v page2=%v", p1["offset"], p2["offset"])
	}

	lots1 := responseLots(t, p1)
	lots2 := responseLots(t, p2)
	if len(lots1) != 2 {
		t.Fatalf("page1: want 2 lots, got %d", len(lots1))
	}
	if len(lots2) != 1 {
		t.Fatalf("page2: want 1 lot, got %d", len(lots2))
	}
	if lots1[0]["lot_id"] == lots2[0]["lot_id"] || lots1[1]["lot_id"] == lots2[0]["lot_id"] {
		t.Fatalf("pages should not overlap: page1=%v page2=%v", lots1, lots2)
	}
}

func TestListLotsStatusActive(t *testing.T) {
	env := newTestEnv(t)
	user := "user_lots_active"
	fixture := seedLotsForAudit(t, env, user)

	status, header, body := doJSON(t, http.MethodGet, env.Server.URL+"/v1/users/"+user+"/lots?status=active", nil)
	if status != http.StatusOK {
		t.Fatalf("want 200, got %d body=%v", status, body)
	}
	assertJSONContentType(t, header)

	lots := responseLots(t, body)
	if body["total"].(float64) != 1 || len(lots) != 1 {
		t.Fatalf("want 1 active lot, got total=%v len=%d", body["total"], len(lots))
	}
	if got := lots[0]["lot_id"].(float64); got != float64(fixture.activeID) {
		t.Fatalf("want active lot_id=%d, got %v", fixture.activeID, got)
	}
	if got := lots[0]["status"]; got != "active" {
		t.Fatalf("want active status, got %v", got)
	}
}

func TestListLotsStatusExpired(t *testing.T) {
	env := newTestEnv(t)
	user := "user_lots_expired"
	fixture := seedLotsForAudit(t, env, user)

	status, header, body := doJSON(t, http.MethodGet, env.Server.URL+"/v1/users/"+user+"/lots?status=expired", nil)
	if status != http.StatusOK {
		t.Fatalf("want 200, got %d body=%v", status, body)
	}
	assertJSONContentType(t, header)

	lots := responseLots(t, body)
	if body["total"].(float64) != 1 || len(lots) != 1 {
		t.Fatalf("want 1 expired lot, got total=%v len=%d", body["total"], len(lots))
	}
	if got := lots[0]["lot_id"].(float64); got != float64(fixture.expiredID) {
		t.Fatalf("want expired lot_id=%d, got %v", fixture.expiredID, got)
	}
	if got := lots[0]["status"]; got != "expired" {
		t.Fatalf("want expired status, got %v", got)
	}
}

func TestListLotsStatusExhausted(t *testing.T) {
	env := newTestEnv(t)
	user := "user_lots_exhausted"
	fixture := seedLotsForAudit(t, env, user)

	status, header, body := doJSON(t, http.MethodGet, env.Server.URL+"/v1/users/"+user+"/lots?status=exhausted", nil)
	if status != http.StatusOK {
		t.Fatalf("want 200, got %d body=%v", status, body)
	}
	assertJSONContentType(t, header)

	lots := responseLots(t, body)
	if body["total"].(float64) != 1 || len(lots) != 1 {
		t.Fatalf("want 1 exhausted lot, got total=%v len=%d", body["total"], len(lots))
	}
	if got := lots[0]["lot_id"].(float64); got != float64(fixture.exhaustedID) {
		t.Fatalf("want exhausted lot_id=%d, got %v", fixture.exhaustedID, got)
	}
	if got := lots[0]["status"]; got != "exhausted" {
		t.Fatalf("want exhausted status, got %v", got)
	}
}

func TestListLotsInvalidStatus(t *testing.T) {
	env := newTestEnv(t)
	user := "user_lots_invalid_status"
	mustAccrue(t, env, user, 100, 30, "lots-invalid-status-seed")

	status, header, body := doJSON(t, http.MethodGet, env.Server.URL+"/v1/users/"+user+"/lots?status=depleted", nil)
	assertErrorResponse(t, status, header, body,
		http.StatusBadRequest, "bad_request", "status must be one of: active, expired, exhausted")
}

func TestListLotsExistingUserWithoutLots(t *testing.T) {
	env := newTestEnv(t)
	user := "user_lots_empty"

	if _, err := env.DB.Exec(`INSERT INTO ledger_entries (user_id, type, amount) VALUES ($1, 'audit', 0)`, user); err != nil {
		t.Fatalf("seed ledger-only user: %v", err)
	}

	status, header, body := doJSON(t, http.MethodGet, env.Server.URL+"/v1/users/"+user+"/lots", nil)
	if status != http.StatusOK {
		t.Fatalf("want 200, got %d body=%v", status, body)
	}
	assertJSONContentType(t, header)

	lots := responseLots(t, body)
	if body["total"].(float64) != 0 {
		t.Fatalf("want total=0, got %v", body["total"])
	}
	if len(lots) != 0 {
		t.Fatalf("want no lots, got %d", len(lots))
	}
}

func TestListLotsUnknownUser(t *testing.T) {
	env := newTestEnv(t)

	status, header, body := doJSON(t, http.MethodGet, env.Server.URL+"/v1/users/no-such-user/lots", nil)
	assertErrorResponse(t, status, header, body,
		http.StatusNotFound, "not_found", "user not found")
}
