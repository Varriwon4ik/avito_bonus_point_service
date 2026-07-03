package data

import (
	"context"
	"fmt"
)

// UpsertAutotestScenario stores or replaces a reusable autotest scenario.
func (s *Store) UpsertAutotestScenario(ctx context.Context, scenario AutotestScenario) (AutotestScenario, error) {
	if scenario.Label == "" {
		return AutotestScenario{}, fmt.Errorf("label is required")
	}
	if scenario.UserID == "" {
		return AutotestScenario{}, fmt.Errorf("user_id is required")
	}
	if scenario.Amount <= 0 {
		return AutotestScenario{}, ErrInvalidAmount
	}
	if scenario.TTLDays <= 0 {
		return AutotestScenario{}, ErrInvalidAmount
	}
	if scenario.ParallelRequests < 2 {
		return AutotestScenario{}, fmt.Errorf("parallel_requests must be at least 2")
	}
	if scenario.LedgerLabel == "" {
		scenario.LedgerLabel = "test"
	}

	err := s.DB.QueryRowContext(ctx, `
		INSERT INTO autotest_scenarios (label, user_id, amount, ttl_days, parallel_requests, ledger_label)
		VALUES ($1, $2, $3, $4, $5, $6)
		ON CONFLICT (label) DO UPDATE SET
			user_id = EXCLUDED.user_id,
			amount = EXCLUDED.amount,
			ttl_days = EXCLUDED.ttl_days,
			parallel_requests = EXCLUDED.parallel_requests,
			ledger_label = EXCLUDED.ledger_label,
			updated_at = now()
		RETURNING id, label, user_id, amount, ttl_days, parallel_requests, ledger_label, created_at, updated_at`,
		scenario.Label,
		scenario.UserID,
		scenario.Amount,
		scenario.TTLDays,
		scenario.ParallelRequests,
		scenario.LedgerLabel,
	).Scan(
		&scenario.ID,
		&scenario.Label,
		&scenario.UserID,
		&scenario.Amount,
		&scenario.TTLDays,
		&scenario.ParallelRequests,
		&scenario.LedgerLabel,
		&scenario.CreatedAt,
		&scenario.UpdatedAt,
	)
	if err != nil {
		return AutotestScenario{}, err
	}
	return scenario, nil
}

// ListAutotestScenarios returns every stored scenario in creation order.
func (s *Store) ListAutotestScenarios(ctx context.Context) ([]AutotestScenario, error) {
	rows, err := s.DB.QueryContext(ctx, `
		SELECT id, label, user_id, amount, ttl_days, parallel_requests, ledger_label, created_at, updated_at
		FROM autotest_scenarios
		ORDER BY created_at ASC, id ASC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var scenarios []AutotestScenario
	for rows.Next() {
		var scenario AutotestScenario
		if err := rows.Scan(
			&scenario.ID,
			&scenario.Label,
			&scenario.UserID,
			&scenario.Amount,
			&scenario.TTLDays,
			&scenario.ParallelRequests,
			&scenario.LedgerLabel,
			&scenario.CreatedAt,
			&scenario.UpdatedAt,
		); err != nil {
			return nil, err
		}
		scenarios = append(scenarios, scenario)
	}
	return scenarios, rows.Err()
}
