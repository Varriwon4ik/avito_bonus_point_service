package data

import (
	"context"
	"database/sql"
	"encoding/json"
)

// CreateHold reserves `amount` points for the user, consuming the
// soonest-to-expire lots first (see allocateLots). The points remain part
// of the user's balance but are no longer "available" until the hold is
// confirmed (permanently spent) or cancelled (released back).
func (s *Store) CreateHold(ctx context.Context, userID string, amount int, idempotencyKey string) (int, []byte, error) {
	return s.withIdempotency(ctx, idempotencyKey, "hold", func(tx *sql.Tx) (int, any, error) {
		if amount <= 0 {
			return 0, nil, ErrInvalidAmount
		}
		exists, err := userExists(ctx, tx, userID)
		if err != nil {
			return 0, nil, err
		}
		if !exists {
			return 0, nil, ErrUserNotFound
		}

		allocs, err := allocateLots(ctx, tx, userID, amount)
		if err != nil {
			return 0, nil, err
		}

		var holdID int64
		if err := tx.QueryRowContext(ctx, `
			INSERT INTO holds (user_id, amount, status) VALUES ($1, $2, 'active')
			RETURNING id`, userID, amount).Scan(&holdID); err != nil {
			return 0, nil, err
		}

		for _, a := range allocs {
			if _, err := tx.ExecContext(ctx, `
				INSERT INTO hold_allocations (hold_id, lot_id, amount) VALUES ($1, $2, $3)`,
				holdID, a.lotID, a.amount); err != nil {
				return 0, nil, err
			}
		}

		if _, err := tx.ExecContext(ctx, `
			INSERT INTO ledger_entries (user_id, type, amount, ref_type, ref_id)
			VALUES ($1, 'hold', $2, 'hold', $3)`, userID, -amount, holdID); err != nil {
			return 0, nil, err
		}

		return 201, HoldResult{HoldID: holdID, UserID: userID, Amount: amount, Status: "active"}, nil
	})
}

// ConfirmHold finalizes an active hold: the held points are permanently
// spent. Calling Confirm again on an already-confirmed hold is a no-op that
// returns the same result (idempotent by construction).
func (s *Store) ConfirmHold(ctx context.Context, holdID int64) (int, []byte, error) {
	tx, err := s.DB.BeginTx(ctx, nil)
	if err != nil {
		return 0, nil, err
	}
	defer tx.Rollback()

	var userID, status string
	var amount int
	err = tx.QueryRowContext(ctx, `
		SELECT user_id, amount, status FROM holds WHERE id = $1 FOR UPDATE`, holdID).
		Scan(&userID, &amount, &status)
	if err == sql.ErrNoRows {
		return 0, nil, ErrHoldNotFound
	}
	if err != nil {
		return 0, nil, err
	}

	if status == "confirmed" {
		if err := tx.Commit(); err != nil {
			return 0, nil, err
		}
		body, _ := json.Marshal(HoldResult{HoldID: holdID, UserID: userID, Amount: amount, Status: status})
		return 200, body, nil
	}
	if status != "active" {
		return 0, nil, ErrInvalidHoldStatus
	}

	if _, err := tx.ExecContext(ctx, `
		UPDATE holds SET status = 'confirmed', updated_at = now() WHERE id = $1`, holdID); err != nil {
		return 0, nil, err
	}
	if _, err := tx.ExecContext(ctx, `
		INSERT INTO ledger_entries (user_id, type, amount, ref_type, ref_id)
		VALUES ($1, 'confirm', 0, 'hold', $2)`, userID, holdID); err != nil {
		return 0, nil, err
	}

	if err := tx.Commit(); err != nil {
		return 0, nil, err
	}
	body, _ := json.Marshal(HoldResult{HoldID: holdID, UserID: userID, Amount: amount, Status: "confirmed"})
	return 200, body, nil
}

// CancelHold releases an active hold, returning its points to the lots they
// were taken from (so they remain subject to their original expiry date).
// Calling Cancel again on an already-cancelled hold is a no-op.
func (s *Store) CancelHold(ctx context.Context, holdID int64) (int, []byte, error) {
	tx, err := s.DB.BeginTx(ctx, nil)
	if err != nil {
		return 0, nil, err
	}
	defer tx.Rollback()

	var userID, status string
	var amount int
	err = tx.QueryRowContext(ctx, `
		SELECT user_id, amount, status FROM holds WHERE id = $1 FOR UPDATE`, holdID).
		Scan(&userID, &amount, &status)
	if err == sql.ErrNoRows {
		return 0, nil, ErrHoldNotFound
	}
	if err != nil {
		return 0, nil, err
	}

	if status == "cancelled" {
		if err := tx.Commit(); err != nil {
			return 0, nil, err
		}
		body, _ := json.Marshal(HoldResult{HoldID: holdID, UserID: userID, Amount: amount, Status: status})
		return 200, body, nil
	}
	if status != "active" {
		return 0, nil, ErrInvalidHoldStatus
	}

	rows, err := tx.QueryContext(ctx, `
		SELECT lot_id, amount FROM hold_allocations WHERE hold_id = $1`, holdID)
	if err != nil {
		return 0, nil, err
	}
	type alloc struct {
		lotID  int64
		amount int
	}
	var allocs []alloc
	for rows.Next() {
		var a alloc
		if err := rows.Scan(&a.lotID, &a.amount); err != nil {
			rows.Close()
			return 0, nil, err
		}
		allocs = append(allocs, a)
	}
	if err := rows.Err(); err != nil {
		return 0, nil, err
	}
	rows.Close()

	for _, a := range allocs {
		if _, err := tx.ExecContext(ctx, `
			UPDATE points_lots SET remaining = remaining + $1 WHERE id = $2`, a.amount, a.lotID); err != nil {
			return 0, nil, err
		}
	}

	if _, err := tx.ExecContext(ctx, `
		UPDATE holds SET status = 'cancelled', updated_at = now() WHERE id = $1`, holdID); err != nil {
		return 0, nil, err
	}
	if _, err := tx.ExecContext(ctx, `
		INSERT INTO ledger_entries (user_id, type, amount, ref_type, ref_id)
		VALUES ($1, 'cancel', $2, 'hold', $3)`, userID, amount, holdID); err != nil {
		return 0, nil, err
	}

	if err := tx.Commit(); err != nil {
		return 0, nil, err
	}
	body, _ := json.Marshal(HoldResult{HoldID: holdID, UserID: userID, Amount: amount, Status: "cancelled"})
	return 200, body, nil
}

// Debit is a convenience one-shot debit: it performs a hold and an
// immediate confirm in a single transaction. It is idempotent the same way
// Accrue and CreateHold are.
func (s *Store) Debit(ctx context.Context, userID string, amount int, idempotencyKey string) (int, []byte, error) {
	return s.withIdempotency(ctx, idempotencyKey, "debit", func(tx *sql.Tx) (int, any, error) {
		if amount <= 0 {
			return 0, nil, ErrInvalidAmount
		}
		exists, err := userExists(ctx, tx, userID)
		if err != nil {
			return 0, nil, err
		}
		if !exists {
			return 0, nil, ErrUserNotFound
		}

		allocs, err := allocateLots(ctx, tx, userID, amount)
		if err != nil {
			return 0, nil, err
		}

		var holdID int64
		if err := tx.QueryRowContext(ctx, `
			INSERT INTO holds (user_id, amount, status) VALUES ($1, $2, 'confirmed')
			RETURNING id`, userID, amount).Scan(&holdID); err != nil {
			return 0, nil, err
		}

		for _, a := range allocs {
			if _, err := tx.ExecContext(ctx, `
				INSERT INTO hold_allocations (hold_id, lot_id, amount) VALUES ($1, $2, $3)`,
				holdID, a.lotID, a.amount); err != nil {
				return 0, nil, err
			}
		}

		if _, err := tx.ExecContext(ctx, `
			INSERT INTO ledger_entries (user_id, type, amount, ref_type, ref_id)
			VALUES ($1, 'debit', $2, 'hold', $3)`, userID, -amount, holdID); err != nil {
			return 0, nil, err
		}

		return 200, DebitResult{HoldID: holdID, UserID: userID, Amount: amount, Status: "confirmed"}, nil
	})
}
