package api_test

import (
	"strings"
	"testing"
)

// TestMetricsEndpoint exercises the full /metrics endpoint against a live
// database: it accrues some points, drives a couple of requests through the
// API so the HTTP counters are non-zero, and then asserts the exposition
// output contains both the request metrics and the ledger gauges. Skipped when
// TEST_DATABASE_URL is unset (see newTestEnv).
func TestMetricsEndpoint(t *testing.T) {
	env := newTestEnv(t)
	base := env.Server.URL

	// Generate some traffic across a couple of routes.
	doJSONWithHeaders(t, "POST", base+"/v1/users/u-metrics/accruals", map[string]any{
		"amount":          500,
		"idempotency_key": "metrics-accrue-1",
	}, map[string]string{"Authorization": "Bearer " + testAdminToken})
	doJSON(t, "GET", base+"/v1/users/u-metrics/balance", nil)
	doJSON(t, "POST", base+"/v1/users/u-metrics/holds", map[string]any{
		"amount":          120,
		"idempotency_key": "metrics-hold-1",
	})

	resp := doRaw(t, "GET", base+"/metrics", "", nil)
	if resp.StatusCode != 200 {
		t.Fatalf("expected 200 from /metrics, got %d", resp.StatusCode)
	}
	if ct := resp.Header.Get("Content-Type"); !strings.HasPrefix(ct, "text/plain") {
		t.Fatalf("expected text/plain content type, got %q", ct)
	}

	out := string(resp.Body)
	wantFamilies := []string{
		"# TYPE http_requests_total counter",
		"# TYPE http_request_duration_seconds histogram",
		`path="/v1/users/{id}/accruals"`,
		`path="/v1/users/{id}/balance"`,
		"# TYPE bonus_points_available gauge",
		"# TYPE bonus_active_holds gauge",
		"bonus_points_held 120",
		"bonus_active_holds 1",
	}
	for _, w := range wantFamilies {
		if !strings.Contains(out, w) {
			t.Errorf("metrics output missing %q\n--- got ---\n%s", w, out)
		}
	}
}
