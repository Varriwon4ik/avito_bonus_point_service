package data

import "context"

// LedgerStats is a point-in-time snapshot of ledger-wide aggregates, exposed
// as Prometheus gauges by the /metrics endpoint.
type LedgerStats struct {
	// AvailablePoints is the sum of remaining points across all non-expired lots.
	AvailablePoints int64
	// HeldPoints is the sum of points reserved by currently active holds.
	HeldPoints int64
	// ActiveHolds is the number of holds still in 'active' status.
	ActiveHolds int64
	// Lots is the total number of points lots ever created.
	Lots int64
	// Users is the number of distinct users that appear in the ledger.
	Users int64
}

// LedgerStats computes ledger-wide gauges in a single round-trip. The
// sub-selects are cheap aggregates over indexed columns and run with whatever
// timeout the caller's context carries.
func (s *Store) LedgerStats(ctx context.Context) (LedgerStats, error) {
	var st LedgerStats
	err := s.DB.QueryRowContext(ctx, `
		SELECT
			(SELECT COALESCE(SUM(remaining), 0) FROM points_lots WHERE expires_at > now()),
			(SELECT COALESCE(SUM(amount), 0)    FROM holds       WHERE status = 'active'),
			(SELECT COUNT(*)                     FROM holds       WHERE status = 'active'),
			(SELECT COUNT(*)                     FROM points_lots),
			(SELECT COUNT(DISTINCT user_id)      FROM ledger_entries)
	`).Scan(&st.AvailablePoints, &st.HeldPoints, &st.ActiveHolds, &st.Lots, &st.Users)
	return st, err
}
