package api

import (
	"log/slog"
	"net/http"
	"strings"
	"time"
)

// statusRecorder wraps an http.ResponseWriter to capture the status code and
// the number of bytes written, which the middleware needs after the handler
// has run. It defaults to 200 because a handler that writes a body without
// calling WriteHeader implies a 200.
type statusRecorder struct {
	http.ResponseWriter
	status      int
	bytes       int
	wroteHeader bool
}

func (r *statusRecorder) WriteHeader(code int) {
	if !r.wroteHeader {
		r.status = code
		r.wroteHeader = true
	}
	r.ResponseWriter.WriteHeader(code)
}

func (r *statusRecorder) Write(b []byte) (int, error) {
	if !r.wroteHeader {
		r.status = http.StatusOK
		r.wroteHeader = true
	}
	n, err := r.ResponseWriter.Write(b)
	r.bytes += n
	return n, err
}

// observe wraps next with structured request logging and metrics collection.
// It is intentionally a thin wrapper that does not touch the handlers
// themselves: it times the request, captures the response status, resolves the
// matched route pattern (so metric/log labels never embed user-specific path
// segments), and records the result.
//
// The request body is never read or logged, so no sensitive payload data
// leaks into logs.
func observe(next http.Handler, apiMux *http.ServeMux, logger *slog.Logger, metrics *Metrics) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		rec := &statusRecorder{ResponseWriter: w, status: http.StatusOK}

		next.ServeHTTP(rec, r)

		latency := time.Since(start)
		path := resolvePattern(apiMux, r)

		metrics.ObserveRequest(r.Method, path, rec.status, latency)

		attrs := []any{
			"method", r.Method,
			"path", path,
			"status", rec.status,
			"latency_ms", latency.Milliseconds(),
			"bytes", rec.bytes,
		}
		// user_id is only meaningful on /v1/users/{id}/... routes; the {id}
		// path value is populated by the inner mux during ServeHTTP above.
		if strings.HasPrefix(path, "/v1/users/{id}") {
			if uid := r.PathValue("id"); uid != "" {
				attrs = append(attrs, "user_id", uid)
			}
		}
		logger.Info("http_request", attrs...)
	})
}

// resolvePattern returns the registered route pattern for r (e.g.
// "/v1/users/{id}/balance") rather than the concrete request path. Falling back
// to the inner API mux lets us recover the templated pattern even though the
// outer mux only sees the "/v1/" subtree. Anything that isn't a known route
// collapses to a small, fixed set of labels to keep cardinality bounded.
func resolvePattern(apiMux *http.ServeMux, r *http.Request) string {
	if apiMux != nil {
		if _, pat := apiMux.Handler(r); pat != "" {
			return stripMethod(pat)
		}
	}
	switch r.URL.Path {
	case "/metrics", "/docs", "/openapi.yaml", "/healthz":
		return r.URL.Path
	default:
		return "other"
	}
}

// stripMethod removes a leading "METHOD " token from a Go 1.22 mux pattern
// such as "GET /v1/users/{id}/balance", leaving just the path template.
func stripMethod(pattern string) string {
	if i := strings.IndexByte(pattern, ' '); i >= 0 {
		return pattern[i+1:]
	}
	return pattern
}
