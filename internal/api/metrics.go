package api

import (
	"io"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"
)

// durationBuckets are the cumulative histogram upper bounds (in seconds) used
// for http_request_duration_seconds. They follow Prometheus' default-ish set,
// trimmed to the latencies a ledger API realistically sees.
var durationBuckets = []float64{
	0.005, 0.01, 0.025, 0.05, 0.1, 0.25, 0.5, 1, 2.5, 5, 10,
}

// Metrics collects per-request counters and latency histograms in memory and
// renders them in the Prometheus text exposition format. It is safe for
// concurrent use. Labels are the request method, the matched route *pattern*
// (never the raw path), and the response status, which keeps label cardinality
// bounded regardless of how many distinct user IDs hit the service.
type Metrics struct {
	mu        sync.Mutex
	requests  map[requestKey]uint64
	durations map[routeKey]*histogram
}

type requestKey struct {
	method string
	path   string
	status int
}

type routeKey struct {
	method string
	path   string
}

// histogram holds per-bucket counts (aligned with durationBuckets) plus the
// running sum and total observation count.
type histogram struct {
	bucketCounts []uint64
	sum          float64
	count        uint64
}

func newHistogram() *histogram {
	return &histogram{bucketCounts: make([]uint64, len(durationBuckets))}
}

// NewMetrics returns an empty, ready-to-use metrics registry.
func NewMetrics() *Metrics {
	return &Metrics{
		requests:  make(map[requestKey]uint64),
		durations: make(map[routeKey]*histogram),
	}
}

// ObserveRequest records a single completed HTTP request.
func (m *Metrics) ObserveRequest(method, path string, status int, latency time.Duration) {
	if m == nil {
		return
	}
	m.mu.Lock()
	defer m.mu.Unlock()

	m.requests[requestKey{method: method, path: path, status: status}]++

	rk := routeKey{method: method, path: path}
	h := m.durations[rk]
	if h == nil {
		h = newHistogram()
		m.durations[rk] = h
	}
	secs := latency.Seconds()
	h.sum += secs
	h.count++
	for i, ub := range durationBuckets {
		if secs <= ub {
			h.bucketCounts[i]++
		}
	}
}

// WriteHTTPMetrics renders the request-count and latency-histogram families in
// Prometheus text format. Output is sorted so scrapes are stable/diffable.
func (m *Metrics) WriteHTTPMetrics(w io.Writer) {
	m.mu.Lock()
	defer m.mu.Unlock()

	io.WriteString(w, "# HELP http_requests_total Total number of HTTP requests handled, by method, route and status.\n")
	io.WriteString(w, "# TYPE http_requests_total counter\n")
	reqKeys := make([]requestKey, 0, len(m.requests))
	for k := range m.requests {
		reqKeys = append(reqKeys, k)
	}
	sort.Slice(reqKeys, func(i, j int) bool {
		if reqKeys[i].path != reqKeys[j].path {
			return reqKeys[i].path < reqKeys[j].path
		}
		if reqKeys[i].method != reqKeys[j].method {
			return reqKeys[i].method < reqKeys[j].method
		}
		return reqKeys[i].status < reqKeys[j].status
	})
	for _, k := range reqKeys {
		labels := labelSet(
			"method", k.method,
			"path", k.path,
			"status", strconv.Itoa(k.status),
		)
		io.WriteString(w, "http_requests_total"+labels+" "+strconv.FormatUint(m.requests[k], 10)+"\n")
	}

	io.WriteString(w, "# HELP http_request_duration_seconds HTTP request latency in seconds, by method and route.\n")
	io.WriteString(w, "# TYPE http_request_duration_seconds histogram\n")
	routeKeys := make([]routeKey, 0, len(m.durations))
	for k := range m.durations {
		routeKeys = append(routeKeys, k)
	}
	sort.Slice(routeKeys, func(i, j int) bool {
		if routeKeys[i].path != routeKeys[j].path {
			return routeKeys[i].path < routeKeys[j].path
		}
		return routeKeys[i].method < routeKeys[j].method
	})
	for _, rk := range routeKeys {
		h := m.durations[rk]
		for i, ub := range durationBuckets {
			labels := labelSet(
				"method", rk.method,
				"path", rk.path,
				"le", strconv.FormatFloat(ub, 'g', -1, 64),
			)
			io.WriteString(w, "http_request_duration_seconds_bucket"+labels+" "+strconv.FormatUint(h.bucketCounts[i], 10)+"\n")
		}
		infLabels := labelSet("method", rk.method, "path", rk.path, "le", "+Inf")
		io.WriteString(w, "http_request_duration_seconds_bucket"+infLabels+" "+strconv.FormatUint(h.count, 10)+"\n")

		base := labelSet("method", rk.method, "path", rk.path)
		io.WriteString(w, "http_request_duration_seconds_sum"+base+" "+strconv.FormatFloat(h.sum, 'g', -1, 64)+"\n")
		io.WriteString(w, "http_request_duration_seconds_count"+base+" "+strconv.FormatUint(h.count, 10)+"\n")
	}
}

// labelSet builds a Prometheus label block like {a="1",b="2"} from alternating
// name/value pairs, escaping values as the exposition format requires.
func labelSet(pairs ...string) string {
	if len(pairs) == 0 {
		return ""
	}
	var b strings.Builder
	b.WriteByte('{')
	for i := 0; i+1 < len(pairs); i += 2 {
		if i > 0 {
			b.WriteByte(',')
		}
		b.WriteString(pairs[i])
		b.WriteString(`="`)
		b.WriteString(escapeLabelValue(pairs[i+1]))
		b.WriteByte('"')
	}
	b.WriteByte('}')
	return b.String()
}

// escapeLabelValue escapes backslashes, double quotes and newlines per the
// Prometheus text format spec.
func escapeLabelValue(v string) string {
	if !strings.ContainsAny(v, "\\\"\n") {
		return v
	}
	r := strings.NewReplacer(`\`, `\\`, `"`, `\"`, "\n", `\n`)
	return r.Replace(v)
}
