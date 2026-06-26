package api_test

// US-09: Transaction history pagination
//
// Acceptance criteria:
//   - GET /v1/users/{id}/transactions?page=N&offset=M returns only the N-th
//     page of M entries, newest-first.
//   - Default page=1, default offset=20 when the params are omitted.
//   - Invalid params (page<1, offset<1 or >500, non-integer) → 400 Bad Request.
//   - Pagination metadata (page, offset, total) is always present in the
//     response alongside the entries slice.
//   - An empty page (page beyond last) returns 200 with an empty entries slice,
//     not a 404.

import (
	"fmt"
	"net/http"
	"testing"
)

// seedTransactions accrues `n` separate lots so there are `n` ledger entries.
func seedTransactions(t *testing.T, env *testEnv, user string, n int) {
	t.Helper()
	for i := 0; i < n; i++ {
		mustAccrue(t, env, user, 10, 30, fmt.Sprintf("acc-%d", i))
	}
}

// TestPaginationDefaults verifies that omitting page and offset still returns
// a well-formed paginated response (page=1, offset=20).
func TestPaginationDefaults(t *testing.T) {
	env := newTestEnv(t)
	user := "user_pg_defaults"
	seedTransactions(t, env, user, 5)

	status, _, body := doJSON(t, http.MethodGet,
		env.Server.URL+"/v1/users/"+user+"/transactions", nil)

	if status != http.StatusOK {
		t.Fatalf("want 200, got %d body=%v", status, body)
	}

	if body["page"].(float64) != 1 {
		t.Errorf("want page=1, got %v", body["page"])
	}
	if body["offset"].(float64) != 20 {
		t.Errorf("want offset=20, got %v", body["offset"])
	}
	if body["total"].(float64) != 5 {
		t.Errorf("want total=5, got %v", body["total"])
	}

	entries := body["entries"].([]any)
	if len(entries) != 5 {
		t.Errorf("want 5 entries on first page, got %d", len(entries))
	}
}

// TestPaginationSecondPage checks that page=2 returns the correct slice.
func TestPaginationSecondPage(t *testing.T) {
	env := newTestEnv(t)
	user := "user_pg_page2"
	seedTransactions(t, env, user, 7)

	// page=1, offset=5 → entries 1-5
	status, _, p1 := doJSON(t, http.MethodGet,
		env.Server.URL+"/v1/users/"+user+"/transactions?page=1&offset=5", nil)
	if status != http.StatusOK {
		t.Fatalf("page1: want 200, got %d body=%v", status, p1)
	}
	entries1 := p1["entries"].([]any)
	if len(entries1) != 5 {
		t.Errorf("page1: want 5 entries, got %d", len(entries1))
	}

	// page=2, offset=5 → entries 6-7
	status, _, p2 := doJSON(t, http.MethodGet,
		env.Server.URL+"/v1/users/"+user+"/transactions?page=2&offset=5", nil)
	if status != http.StatusOK {
		t.Fatalf("page2: want 200, got %d body=%v", status, p2)
	}
	entries2 := p2["entries"].([]any)
	if len(entries2) != 2 {
		t.Errorf("page2: want 2 entries, got %d", len(entries2))
	}

	// Verify total is consistent across pages.
	if p1["total"] != p2["total"] {
		t.Errorf("total mismatch between pages: %v vs %v", p1["total"], p2["total"])
	}
	if p2["total"].(float64) != 7 {
		t.Errorf("want total=7, got %v", p2["total"])
	}

	// Entries on page 2 must not appear on page 1 (no overlap).
	ids1 := map[any]bool{}
	for _, e := range entries1 {
		ids1[e.(map[string]any)["id"]] = true
	}
	for _, e := range entries2 {
		if ids1[e.(map[string]any)["id"]] {
			t.Errorf("entry id=%v appears on both page 1 and page 2", e.(map[string]any)["id"])
		}
	}

	t.Logf("✓ page1=%d entries, page2=%d entries, total=7, no overlap",
		len(entries1), len(entries2))
}

// TestPaginationBeyondLastPage checks that requesting a page past the end
// returns 200 with an empty entries slice, not an error.
func TestPaginationBeyondLastPage(t *testing.T) {
	env := newTestEnv(t)
	user := "user_pg_empty"
	seedTransactions(t, env, user, 3)

	status, _, body := doJSON(t, http.MethodGet,
		env.Server.URL+"/v1/users/"+user+"/transactions?page=99&offset=20", nil)

	if status != http.StatusOK {
		t.Fatalf("want 200, got %d body=%v", status, body)
	}

	entries := body["entries"].([]any)
	if len(entries) != 0 {
		t.Errorf("want 0 entries past last page, got %d", len(entries))
	}
	if body["total"].(float64) != 3 {
		t.Errorf("want total=3, got %v", body["total"])
	}

	t.Logf("✓ page beyond end returns empty entries with correct total")
}

// TestPaginationNewestFirst verifies entries are returned newest-first.
func TestPaginationNewestFirst(t *testing.T) {
	env := newTestEnv(t)
	user := "user_pg_order"
	seedTransactions(t, env, user, 3)

	status, _, body := doJSON(t, http.MethodGet,
		env.Server.URL+"/v1/users/"+user+"/transactions?page=1&offset=10", nil)
	if status != http.StatusOK {
		t.Fatalf("want 200, got %d body=%v", status, body)
	}

	entries := body["entries"].([]any)
	if len(entries) < 2 {
		t.Fatalf("need at least 2 entries to check ordering, got %d", len(entries))
	}

	// IDs should be descending (newest = highest ID first).
	for i := 0; i < len(entries)-1; i++ {
		idA := entries[i].(map[string]any)["id"].(float64)
		idB := entries[i+1].(map[string]any)["id"].(float64)
		if idA <= idB {
			t.Errorf("entries not newest-first: entry[%d].id=%v <= entry[%d].id=%v", i, idA, i+1, idB)
		}
	}

	t.Logf("✓ entries are newest-first")
}

// TestPaginationInvalidParams covers every invalid-parameter combination.
func TestPaginationInvalidParams(t *testing.T) {
	env := newTestEnv(t)
	user := "user_pg_invalid"
	mustAccrue(t, env, user, 100, 30, "acc-seed")

	cases := []struct {
		query       string
		wantMessage string
	}{
		{"page=0", "page must be a positive integer"},
		{"page=-1", "page must be a positive integer"},
		{"page=abc", "page must be a positive integer"},
		{"offset=0", "offset must be between 1 and 500"},
		{"offset=-5", "offset must be between 1 and 500"},
		{"offset=501", "offset must be between 1 and 500"},
		{"offset=abc", "offset must be between 1 and 500"},
	}

	for _, tc := range cases {
		t.Run(tc.query, func(t *testing.T) {
			status, header, body := doJSON(t, http.MethodGet,
				env.Server.URL+"/v1/users/"+user+"/transactions?"+tc.query, nil)
			assertErrorResponse(t, status, header, body,
				http.StatusBadRequest, "bad_request", tc.wantMessage)
		})
	}
}

// TestPaginationUnknownUser verifies that a missing user returns 404.
func TestPaginationUnknownUser(t *testing.T) {
	env := newTestEnv(t)

	status, header, body := doJSON(t, http.MethodGet,
		env.Server.URL+"/v1/users/no-such-user/transactions", nil)
	assertErrorResponse(t, status, header, body,
		http.StatusNotFound, "not_found", "user not found")
}

// TestPaginationExactPageBoundary seeds exactly offset*N entries and checks
// that each page is full and the last+1 page is empty.
func TestPaginationExactPageBoundary(t *testing.T) {
	env := newTestEnv(t)
	user := "user_pg_boundary"
	const perPage = 4
	const pages = 3
	seedTransactions(t, env, user, perPage*pages) // 12 entries total

	for pg := 1; pg <= pages; pg++ {
		status, _, body := doJSON(t, http.MethodGet,
			fmt.Sprintf("%s/v1/users/%s/transactions?page=%d&offset=%d",
				env.Server.URL, user, pg, perPage), nil)
		if status != http.StatusOK {
			t.Fatalf("page%d: want 200, got %d body=%v", pg, status, body)
		}
		got := len(body["entries"].([]any))
		if got != perPage {
			t.Errorf("page%d: want %d entries, got %d", pg, perPage, got)
		}
	}

	// Page pages+1 must be empty.
	status, _, body := doJSON(t, http.MethodGet,
		fmt.Sprintf("%s/v1/users/%s/transactions?page=%d&offset=%d",
			env.Server.URL, user, pages+1, perPage), nil)
	if status != http.StatusOK {
		t.Fatalf("overflow page: want 200, got %d body=%v", status, body)
	}
	if got := len(body["entries"].([]any)); got != 0 {
		t.Errorf("overflow page: want 0 entries, got %d", got)
	}

	t.Logf("✓ %d pages × %d entries each, overflow page empty", pages, perPage)
}
