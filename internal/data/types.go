package data

import (
	"errors"
	"time"
)

var (
	ErrInvalidAmount       = errors.New("amount must be a positive integer")
	ErrInsufficientFunds   = errors.New("insufficient available points")
	ErrUserNotFound        = errors.New("user not found")
	ErrHoldNotFound        = errors.New("hold not found")
	ErrInvalidHoldStatus   = errors.New("hold is not in 'active' status")
	ErrIdempotencyConflict = errors.New("a request with this idempotency key is already in progress")
)

// AccrualResult is returned after points are credited to a user's account.
type AccrualResult struct {
	LotID     int64     `json:"lot_id"`
	UserID    string    `json:"user_id"`
	Amount    int       `json:"amount"`
	ExpiresAt time.Time `json:"expires_at"`
}

// BalanceResult describes a user's current points position.
type BalanceResult struct {
	UserID       string `json:"user_id"`
	Available    int    `json:"available"`
	Held         int    `json:"held"`
	Total        int    `json:"total"`
	ExpiringSoon int    `json:"expiring_soon"`
}

// HoldResult describes a two-phase hold (reservation) of points.
type HoldResult struct {
	HoldID int64  `json:"hold_id"`
	UserID string `json:"user_id"`
	Amount int    `json:"amount"`
	Status string `json:"status"`
}

// DebitResult describes a one-shot (already confirmed) debit of points.
type DebitResult struct {
	HoldID int64  `json:"hold_id"`
	UserID string `json:"user_id"`
	Amount int    `json:"amount"`
	Status string `json:"status"`
}

// LedgerEntry is a single immutable row of the points ledger / audit log.
type LedgerEntry struct {
	ID        int64     `json:"id"`
	UserID    string    `json:"user_id"`
	Type      string    `json:"type"`
	Amount    int       `json:"amount"`
	RefType   *string   `json:"ref_type,omitempty"`
	RefID     *int64    `json:"ref_id,omitempty"`
	Note      *string   `json:"note,omitempty"`
	CreatedAt time.Time `json:"created_at"`
}

// PaginatedLedger wraps a page of ledger entries with pagination metadata.
type PaginatedLedger struct {
	UserID  string        `json:"user_id"`
	Page    int           `json:"page"`
	Offset  int           `json:"offset"`
	Total   int           `json:"total"`
	Entries []LedgerEntry `json:"entries"`
}

// LotInfo describes a single batch of accrued points and its remaining balance.
type LotInfo struct {
	LotID     int64     `json:"lot_id"`
	Amount    int       `json:"amount"`
	Remaining int       `json:"remaining"`
	ExpiresAt time.Time `json:"expires_at"`
	CreatedAt time.Time `json:"created_at"`
}
