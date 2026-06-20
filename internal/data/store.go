package data

import (
	"context"
	"database/sql"
	"encoding/json"

	"github.com/lib/pq"
)

// Store provides persistent access to the bonus-points ledger.
type Store struct {
	DB *sql.DB
}

func NewStore(db *sql.DB) *Store {
	return &Store{DB: db}
}

type rowQuerier interface {
	QueryRowContext(context.Context, string, ...any) *sql.Row
}

func userExists(ctx context.Context, q rowQuerier, userID string) (bool, error) {
	var exists bool
	err := q.QueryRowContext(ctx, `
		SELECT EXISTS (
			SELECT 1 FROM points_lots WHERE user_id = $1
			UNION ALL
			SELECT 1 FROM holds WHERE user_id = $1
			UNION ALL
			SELECT 1 FROM ledger_entries WHERE user_id = $1
		)`, userID).Scan(&exists)
	return exists, err
}

// withIdempotency runs fn inside a transaction. If a non-empty idempotency
// key was already used for this endpoint, the previously stored response is
// returned without re-running fn (so retries of the same logical request are
// safe and do not double-apply side effects).
func (s *Store) withIdempotency(ctx context.Context, key, endpoint string, fn func(tx *sql.Tx) (int, any, error)) (int, []byte, error) {
	tx, err := s.DB.BeginTx(ctx, nil)
	if err != nil {
		return 0, nil, err
	}
	defer tx.Rollback()

	if key != "" {
		var status int
		var body []byte
		err := tx.QueryRowContext(ctx,
			`SELECT response_status, response_body FROM idempotency_keys WHERE key = $1 AND endpoint = $2`,
			key, endpoint).Scan(&status, &body)
		switch {
		case err == nil:
			if status == 0 || len(body) == 0 {
				return 0, nil, ErrIdempotencyConflict
			}
			if err := tx.Commit(); err != nil {
				return 0, nil, err
			}
			return status, body, nil
		case err == sql.ErrNoRows:
			// not seen before, reserve the key
			if _, err := tx.ExecContext(ctx,
				`INSERT INTO idempotency_keys (key, endpoint) VALUES ($1, $2)`, key, endpoint); err != nil {
				if pqErr, ok := err.(*pq.Error); ok && pqErr.Code == "23505" {
					return 0, nil, ErrIdempotencyConflict
				}
				return 0, nil, err
			}
		default:
			return 0, nil, err
		}
	}

	status, value, err := fn(tx)
	if err != nil {
		return 0, nil, err
	}

	body, err := json.Marshal(value)
	if err != nil {
		return 0, nil, err
	}

	if key != "" {
		if _, err := tx.ExecContext(ctx,
			`UPDATE idempotency_keys SET response_status = $1, response_body = $2 WHERE key = $3 AND endpoint = $4`,
			status, body, key, endpoint); err != nil {
			return 0, nil, err
		}
	}

	if err := tx.Commit(); err != nil {
		return 0, nil, err
	}
	return status, body, nil
}

type lotAllocation struct {
	lotID  int64
	amount int
}

// allocateLots locks every non-expired lot with remaining points for the
// user (ORDER BY expires_at ASC ... FOR UPDATE), then greedily consumes the
// lots that expire soonest first until `amount` points have been allocated.
// The row locks are held until the caller's transaction commits or rolls
// back, which is what makes concurrent holds/debits for the same user
// consistent (pessimistic concurrency control at the database level).
func allocateLots(ctx context.Context, tx *sql.Tx, userID string, amount int) ([]lotAllocation, error) {
	rows, err := tx.QueryContext(ctx, `
		SELECT id, remaining FROM points_lots
		WHERE user_id = $1 AND remaining > 0 AND expires_at > now()
		ORDER BY expires_at ASC, id ASC
		FOR UPDATE`, userID)
	if err != nil {
		return nil, err
	}

	type lotRow struct {
		id        int64
		remaining int
	}
	var lots []lotRow
	for rows.Next() {
		var l lotRow
		if err := rows.Scan(&l.id, &l.remaining); err != nil {
			rows.Close()
			return nil, err
		}
		lots = append(lots, l)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	rows.Close()

	need := amount
	var allocs []lotAllocation
	for _, l := range lots {
		if need == 0 {
			break
		}
		take := l.remaining
		if take > need {
			take = need
		}
		allocs = append(allocs, lotAllocation{lotID: l.id, amount: take})
		need -= take
	}
	if need > 0 {
		return nil, ErrInsufficientFunds
	}

	for _, a := range allocs {
		if _, err := tx.ExecContext(ctx,
			`UPDATE points_lots SET remaining = remaining - $1 WHERE id = $2`, a.amount, a.lotID); err != nil {
			return nil, err
		}
	}
	return allocs, nil
}
