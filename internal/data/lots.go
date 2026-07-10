package data

import (
	"context"
	"database/sql"
	"time"
)

// Accrue credits the user with `amount` points that expire after ttlDays
// days. If idempotencyKey has been used before for the "accrual" endpoint,
// the original result is returned and no new points are created.
func (s *Store) Accrue(ctx context.Context, userID string, amount, ttlDays int, idempotencyKey string) (int, []byte, error) {
	return s.AccrueWithLabel(ctx, userID, amount, ttlDays, idempotencyKey, "")
}

// AccrueWithLabel stores the same accrual as Accrue, but attaches an optional
// user-facing label to the resulting ledger entry.
func (s *Store) AccrueWithLabel(ctx context.Context, userID string, amount, ttlDays int, idempotencyKey, label string) (int, []byte, error) {
	normalizedLabel, err := NormalizeTransactionLabel(label)
	if err != nil {
		return 0, nil, err
	}

	return s.withIdempotency(ctx, idempotencyKey, "accrual", func(tx *sql.Tx) (int, any, error) {
		if amount <= 0 {
			return 0, nil, ErrInvalidAmount
		}
		if ttlDays <= 0 {
			return 0, nil, ErrInvalidAmount
		}

		expiresAt := time.Now().UTC().AddDate(0, 0, ttlDays)

		var lotID int64
		err = tx.QueryRowContext(ctx, `
			INSERT INTO points_lots (user_id, amount, remaining, expires_at)
			VALUES ($1, $2, $2, $3)
			RETURNING id`, userID, amount, expiresAt).Scan(&lotID)
		if err != nil {
			return 0, nil, err
		}

		var labelArg any
		if normalizedLabel != "" {
			labelArg = normalizedLabel
		}

		if _, err := tx.ExecContext(ctx, `
			INSERT INTO ledger_entries (user_id, type, amount, ref_type, ref_id, label)
			VALUES ($1, 'accrual', $2, 'lot', $3, $4)`, userID, amount, lotID, labelArg); err != nil {
			return 0, nil, err
		}

		return 201, AccrualResult{LotID: lotID, UserID: userID, Amount: amount, ExpiresAt: expiresAt}, nil
	})
}

// Balance reports the user's available (spendable now), held (reserved by
// active holds) and total points, plus how many available points will
// expire within `expiringWithinDays` days.
func (s *Store) Balance(ctx context.Context, userID string, expiringWithinDays int) (BalanceResult, error) {
	res := BalanceResult{UserID: userID}

	exists, err := userExists(ctx, s.DB, userID)
	if err != nil {
		return res, err
	}
	if !exists {
		return res, ErrUserNotFound
	}

	if err := s.DB.QueryRowContext(ctx, `
		SELECT COALESCE(SUM(remaining), 0) FROM points_lots
		WHERE user_id = $1 AND expires_at > now()`, userID).Scan(&res.Available); err != nil {
		return res, err
	}

	if err := s.DB.QueryRowContext(ctx, `
		SELECT COALESCE(SUM(amount), 0) FROM holds
		WHERE user_id = $1 AND status = 'active'`, userID).Scan(&res.Held); err != nil {
		return res, err
	}

	if err := s.DB.QueryRowContext(ctx, `
		SELECT COALESCE(SUM(remaining), 0) FROM points_lots
		WHERE user_id = $1 AND expires_at > now()
		  AND expires_at <= now() + make_interval(days => $2)`, userID, expiringWithinDays).Scan(&res.ExpiringSoon); err != nil {
		return res, err
	}

	res.Total = res.Available + res.Held
	return res, nil
}

func classifyLotStatus(expiresAt time.Time, remaining int, now time.Time) string {
	if !expiresAt.After(now) {
		return LotStatusExpired
	}
	if remaining == 0 {
		return LotStatusExhausted
	}
	return LotStatusActive
}

// ListLots returns a single page of lot-audit rows for a user ordered by
// expiry date, oldest first. The optional status filter uses the support/API
// semantics:
//   - active: non-expired lots with remaining points
//   - expired: any lot whose expiry has passed
//   - exhausted: non-expired lots with zero remaining points
func (s *Store) ListLots(ctx context.Context, userID string, page, offset int, status string) (PaginatedLots, error) {
	result := PaginatedLots{
		UserID: userID,
		Page:   page,
		Offset: offset,
		Lots:   []LotInfo{},
	}

	exists, err := userExists(ctx, s.DB, userID)
	if err != nil {
		return result, err
	}
	if !exists {
		return result, ErrUserNotFound
	}

	now := time.Now().UTC()

	if err := s.DB.QueryRowContext(ctx, `
		SELECT COUNT(*) FROM points_lots
		WHERE user_id = $1
		  AND (
			$2 = ''
			OR ($2 = 'active' AND expires_at > $3 AND remaining > 0)
			OR ($2 = 'expired' AND expires_at <= $3)
			OR ($2 = 'exhausted' AND expires_at > $3 AND remaining = 0)
		  )`, userID, status, now).Scan(&result.Total); err != nil {
		return result, err
	}

	skip := (page - 1) * offset
	rows, err := s.DB.QueryContext(ctx, `
		SELECT id, user_id, amount, remaining, expires_at, created_at
		FROM points_lots
		WHERE user_id = $1
		  AND (
			$2 = ''
			OR ($2 = 'active' AND expires_at > $3 AND remaining > 0)
			OR ($2 = 'expired' AND expires_at <= $3)
			OR ($2 = 'exhausted' AND expires_at > $3 AND remaining = 0)
		  )
		ORDER BY expires_at ASC, id ASC
		LIMIT $4 OFFSET $5`, userID, status, now, offset, skip)
	if err != nil {
		return result, err
	}
	defer rows.Close()

	for rows.Next() {
		var l LotInfo
		if err := rows.Scan(&l.LotID, &l.UserID, &l.Amount, &l.Remaining, &l.ExpiresAt, &l.CreatedAt); err != nil {
			return result, err
		}
		l.Status = classifyLotStatus(l.ExpiresAt, l.Remaining, now)
		result.Lots = append(result.Lots, l)
	}
	return result, rows.Err()
}

// ListLedger returns a single page of ledger entries for a user, newest first.
// page is 1-based; offset is the number of entries per page (1–500).
func (s *Store) ListLedger(ctx context.Context, userID string, page, offset int) (PaginatedLedger, error) {
	result := PaginatedLedger{UserID: userID, Page: page, Offset: offset}

	exists, err := userExists(ctx, s.DB, userID)
	if err != nil {
		return result, err
	}
	if !exists {
		return result, ErrUserNotFound
	}

	// Count total entries for this user so callers can compute page counts.
	if err := s.DB.QueryRowContext(ctx,
		`SELECT COUNT(*) FROM ledger_entries WHERE user_id = $1`, userID,
	).Scan(&result.Total); err != nil {
		return result, err
	}

	skip := (page - 1) * offset
	rows, err := s.DB.QueryContext(ctx, `
		SELECT id, user_id, type, amount, ref_type, ref_id, label, note, created_at
		FROM ledger_entries
		WHERE user_id = $1
		ORDER BY created_at DESC, id DESC
		LIMIT $2 OFFSET $3`, userID, offset, skip)
	if err != nil {
		return result, err
	}
	defer rows.Close()

	result.Entries = []LedgerEntry{}
	for rows.Next() {
		var e LedgerEntry
		if err := rows.Scan(&e.ID, &e.UserID, &e.Type, &e.Amount, &e.RefType, &e.RefID, &e.Label, &e.Note, &e.CreatedAt); err != nil {
			return result, err
		}
		result.Entries = append(result.Entries, e)
	}
	return result, rows.Err()
}
