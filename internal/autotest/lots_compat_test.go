package autotest

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
	"time"

	"bonus-ledger/internal/data"
)

func testLot(id int64, remaining int) data.LotInfo {
	now := time.Now().UTC()
	return data.LotInfo{
		LotID:     id,
		UserID:    "user_compat",
		Amount:    100,
		Remaining: remaining,
		Status:    data.LotStatusActive,
		ExpiresAt: now.Add(24 * time.Hour),
		CreatedAt: now,
	}
}

func TestLoadLotsLegacyArrayResponse(t *testing.T) {
	lots := []data.LotInfo{testLot(1, 100), testLot(2, 40)}

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(lots)
	}))
	defer srv.Close()

	got, err := NewRuntime(srv.URL).loadLots("user_compat")
	if err != nil {
		t.Fatalf("loadLots against legacy array: %v", err)
	}
	if len(got) != 2 {
		t.Fatalf("want 2 lots, got %d", len(got))
	}
	if got[0].LotID != 1 || got[1].LotID != 2 {
		t.Fatalf("unexpected lot ids: %v", got)
	}
}

func TestLoadLotsPaginatedEnvelopeResponse(t *testing.T) {
	pages := map[int][]data.LotInfo{
		1: {testLot(1, 100), testLot(2, 40)},
		2: {testLot(3, 70)},
	}

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		page, err := strconv.Atoi(r.URL.Query().Get("page"))
		if err != nil {
			t.Errorf("bad page param %q: %v", r.URL.Query().Get("page"), err)
			page = 1
		}
		lots, ok := pages[page]
		if !ok {
			lots = []data.LotInfo{}
		}
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(data.PaginatedLots{
			UserID: "user_compat",
			Page:   page,
			Offset: 100,
			Total:  3,
			Lots:   lots,
		})
	}))
	defer srv.Close()

	got, err := NewRuntime(srv.URL).loadLots("user_compat")
	if err != nil {
		t.Fatalf("loadLots against paginated envelope: %v", err)
	}
	if len(got) != 3 {
		t.Fatalf("want 3 lots across pages, got %d", len(got))
	}
	for i, wantID := range []int64{1, 2, 3} {
		if got[i].LotID != wantID {
			t.Fatalf("lot[%d]: want lot_id=%d, got %d", i, wantID, got[i].LotID)
		}
	}
}

func TestLoadLotsUnknownUser(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		_ = json.NewEncoder(w).Encode(map[string]string{
			"error":   "not_found",
			"message": data.ErrUserNotFound.Error(),
		})
	}))
	defer srv.Close()

	got, err := NewRuntime(srv.URL).loadLots("no-such-user")
	if err != nil {
		t.Fatalf("loadLots against missing user: %v", err)
	}
	if got != nil {
		t.Fatalf("want nil lots for missing user, got %v", got)
	}
}
