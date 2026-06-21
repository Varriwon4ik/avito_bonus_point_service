package api

import (
	"context"
	"io"
	"net/http"
	"strconv"
	"time"

	"bonus-ledger/internal/data"
)

// handleMetrics serves the Prometheus exposition endpoint. It is intentionally
// unauthenticated so a scraper on the internal network can reach it. It always
// emits the HTTP request metrics; the ledger gauges are best-effort, so a
// transient database error degrades the response (gauges omitted) instead of
// failing the scrape.
func (s *Server) handleMetrics(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	stats, statsErr := s.Store.LedgerStats(ctx)
	if statsErr != nil {
		s.Logger.Error("metrics: failed to read ledger stats", "err", statsErr)
	}

	w.Header().Set("Content-Type", "text/plain; version=0.0.4; charset=utf-8")
	w.WriteHeader(http.StatusOK)

	s.Metrics.WriteHTTPMetrics(w)
	if statsErr == nil {
		writeLedgerGauges(w, stats)
	}
}

// writeLedgerGauges renders the ledger-level gauges in Prometheus text format.
func writeLedgerGauges(w io.Writer, st data.LedgerStats) {
	gauges := []struct {
		name string
		help string
		val  int64
	}{
		{"bonus_points_available", "Spendable points across all non-expired lots.", st.AvailablePoints},
		{"bonus_points_held", "Points reserved by currently active holds.", st.HeldPoints},
		{"bonus_active_holds", "Number of holds currently in 'active' status.", st.ActiveHolds},
		{"bonus_lots_total", "Total number of points lots created.", st.Lots},
		{"bonus_users_total", "Number of distinct users present in the ledger.", st.Users},
	}
	for _, g := range gauges {
		io.WriteString(w, "# HELP "+g.name+" "+g.help+"\n")
		io.WriteString(w, "# TYPE "+g.name+" gauge\n")
		io.WriteString(w, g.name+" "+strconv.FormatInt(g.val, 10)+"\n")
	}
}
