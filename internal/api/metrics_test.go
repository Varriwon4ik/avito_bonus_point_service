package api

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"bonus-ledger/internal/data"
)

func TestMetricsObserveAndRender(t *testing.T) {
	m := NewMetrics()
	m.ObserveRequest("GET", "/v1/users/{id}/balance", 200, 30*time.Millisecond)
	m.ObserveRequest("GET", "/v1/users/{id}/balance", 200, 80*time.Millisecond)
	m.ObserveRequest("POST", "/v1/users/{id}/accruals", 201, 5*time.Millisecond)

	var sb strings.Builder
	m.WriteHTTPMetrics(&sb)
	out := sb.String()

	wantSubstrings := []string{
		"# TYPE http_requests_total counter",
		`http_requests_total{method="GET",path="/v1/users/{id}/balance",status="200"} 2`,
		`http_requests_total{method="POST",path="/v1/users/{id}/accruals",status="201"} 1`,
		"# TYPE http_request_duration_seconds histogram",
		`http_request_duration_seconds_bucket{method="GET",path="/v1/users/{id}/balance",le="+Inf"} 2`,
		`http_request_duration_seconds_count{method="GET",path="/v1/users/{id}/balance"} 2`,
	}
	for _, w := range wantSubstrings {
		if !strings.Contains(out, w) {
			t.Errorf("metrics output missing %q\n--- got ---\n%s", w, out)
		}
	}

	// The 30ms and 80ms observations both fall at or below the 0.1s bucket but
	// only the 30ms one falls at or below 0.05s.
	if !strings.Contains(out, `http_request_duration_seconds_bucket{method="GET",path="/v1/users/{id}/balance",le="0.05"} 1`) {
		t.Errorf("expected le=0.05 bucket count of 1\n--- got ---\n%s", out)
	}
	if !strings.Contains(out, `http_request_duration_seconds_bucket{method="GET",path="/v1/users/{id}/balance",le="0.1"} 2`) {
		t.Errorf("expected le=0.1 bucket count of 2\n--- got ---\n%s", out)
	}
}

func TestMetricsNilSafe(t *testing.T) {
	var m *Metrics
	// Should not panic when the registry is nil.
	m.ObserveRequest("GET", "/healthz", 200, time.Millisecond)
}

func TestResolvePattern(t *testing.T) {
	srv := &Server{Mux: http.NewServeMux()}
	srv.routes()

	cases := []struct {
		method string
		path   string
		want   string
	}{
		{"GET", "/v1/users/42/balance", "/v1/users/{id}/balance"},
		{"POST", "/v1/users/abc/accruals", "/v1/users/{id}/accruals"},
		{"POST", "/v1/holds/7/confirm", "/v1/holds/{id}/confirm"},
		{"GET", "/healthz", "/healthz"},
		{"GET", "/metrics", "/metrics"},
		{"GET", "/totally/unknown", "other"},
	}
	for _, c := range cases {
		r := httptest.NewRequest(c.method, c.path, nil)
		if got := resolvePattern(srv.Mux, r); got != c.want {
			t.Errorf("resolvePattern(%s %s) = %q, want %q", c.method, c.path, got, c.want)
		}
	}
}

func TestLedgerGaugesRender(t *testing.T) {
	var sb strings.Builder
	writeLedgerGauges(&sb, data.LedgerStats{
		AvailablePoints: 1500,
		HeldPoints:      200,
		ActiveHolds:     3,
		Lots:            12,
		Users:           5,
	})
	out := sb.String()
	for _, w := range []string{
		"# TYPE bonus_points_available gauge",
		"bonus_points_available 1500",
		"bonus_points_held 200",
		"bonus_active_holds 3",
		"bonus_lots_total 12",
		"bonus_users_total 5",
	} {
		if !strings.Contains(out, w) {
			t.Errorf("gauge output missing %q\n--- got ---\n%s", w, out)
		}
	}
}

func TestEscapeLabelValue(t *testing.T) {
	if got := escapeLabelValue(`a"b\c`); got != `a\"b\\c` {
		t.Errorf("escapeLabelValue = %q", got)
	}
	if got := escapeLabelValue("/v1/users/{id}/balance"); got != "/v1/users/{id}/balance" {
		t.Errorf("unexpected escaping: %q", got)
	}
}
