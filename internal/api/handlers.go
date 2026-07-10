package api

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"strings"

	"bonus-ledger/internal/data"
)

func (s *Server) routes() {
	s.Mux.HandleFunc("POST /v1/accruals/batch", s.handleBatchAccrue)
	s.Mux.HandleFunc("POST /v1/users/{id}/accruals", s.handleAccrue)
	s.Mux.HandleFunc("GET /v1/users/{id}/balance", s.handleBalance)
	s.Mux.HandleFunc("GET /v1/users/{id}/lots", s.handleListLots)
	s.Mux.HandleFunc("GET /v1/users/{id}/transactions", s.handleListLedger)
	s.Mux.HandleFunc("POST /v1/users/{id}/holds", s.handleCreateHold)
	s.Mux.HandleFunc("POST /v1/users/{id}/debits", s.handleDebit)
	s.Mux.HandleFunc("POST /v1/holds/{id}/confirm", s.handleConfirmHold)
	s.Mux.HandleFunc("POST /v1/holds/{id}/cancel", s.handleCancelHold)
	s.Mux.HandleFunc("POST /v1/autotest/run", s.handleAutotestRun)
	s.Mux.HandleFunc("GET /healthz", s.handleHealthz)
}

func (s *Server) handleHealthz(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

func parsePageAndOffset(w http.ResponseWriter, r *http.Request) (page, offset int, ok bool) {
	page = 1
	if v := r.URL.Query().Get("page"); v != "" {
		n, err := strconv.Atoi(v)
		if err != nil || n < 1 {
			badRequest(w, "page must be a positive integer")
			return 0, 0, false
		}
		page = n
	}

	offset = 20
	if v := r.URL.Query().Get("offset"); v != "" {
		n, err := strconv.Atoi(v)
		if err != nil || n < 1 || n > 500 {
			badRequest(w, "offset must be between 1 and 500")
			return 0, 0, false
		}
		offset = n
	}

	return page, offset, true
}

type accrueRequest struct {
	Amount         *int    `json:"amount"`
	TTLDays        *int    `json:"ttl_days,omitempty"`
	IdempotencyKey *string `json:"idempotency_key"`
	Label          *string `json:"label,omitempty"`
}

type batchAccrualRequest struct {
	Items []batchAccrualItem `json:"items"`
}

type batchAccrualItem struct {
	UserID         *string `json:"user_id"`
	Amount         *int    `json:"amount"`
	TTLDays        *int    `json:"ttl_days,omitempty"`
	IdempotencyKey *string `json:"idempotency_key"`
	Label          *string `json:"label,omitempty"`
}

// handleAccrue implements "по каждому пользователю можно добавить бонусные
// баллы" with a configurable (and optionally per-request) lifetime, and is
// idempotent via idempotency_key.
func (s *Server) handleAccrue(w http.ResponseWriter, r *http.Request) {
	userID := r.PathValue("id")

	var req accrueRequest
	if err := readJSON(w, r, &req); err != nil {
		badRequest(w, err.Error())
		return
	}
	if req.Amount == nil {
		badRequest(w, "amount is required")
		return
	}
	if req.IdempotencyKey == nil || strings.TrimSpace(*req.IdempotencyKey) == "" {
		badRequest(w, "idempotency_key is required")
		return
	}
	if *req.Amount <= 0 {
		badRequest(w, "amount must be a positive integer")
		return
	}

	ttl, err := s.resolveTTLDays(req.TTLDays)
	if err != nil {
		badRequest(w, err.Error())
		return
	}

	label := ""
	if req.Label != nil {
		label = strings.TrimSpace(*req.Label)
	}

	status, body, err := s.Store.AccrueWithLabel(r.Context(), userID, *req.Amount, ttl, strings.TrimSpace(*req.IdempotencyKey), label)
	s.respond(w, status, body, err)
}

func (s *Server) handleBatchAccrue(w http.ResponseWriter, r *http.Request) {
	var req batchAccrualRequest
	if err := readJSON(w, r, &req); err != nil {
		badRequest(w, err.Error())
		return
	}
	if len(req.Items) == 0 {
		badRequest(w, "items must not be empty")
		return
	}

	results := make([]map[string]any, 0, len(req.Items))
	for idx, item := range req.Items {
		result := map[string]any{"index": idx}
		if item.UserID == nil || strings.TrimSpace(*item.UserID) == "" {
			result["status"] = "error"
			result["error"] = "bad_request"
			result["message"] = "user_id is required"
			results = append(results, result)
			continue
		}
		if item.Amount == nil {
			result["status"] = "error"
			result["error"] = "bad_request"
			result["message"] = "amount is required"
			results = append(results, result)
			continue
		}
		if item.IdempotencyKey == nil || strings.TrimSpace(*item.IdempotencyKey) == "" {
			result["status"] = "error"
			result["error"] = "bad_request"
			result["message"] = "idempotency_key is required"
			results = append(results, result)
			continue
		}
		if *item.Amount <= 0 {
			result["status"] = "error"
			result["error"] = "bad_request"
			result["message"] = "amount must be a positive integer"
			results = append(results, result)
			continue
		}

		ttl, err := s.resolveTTLDays(item.TTLDays)
		if err != nil {
			result["status"] = "error"
			result["error"] = "bad_request"
			result["message"] = err.Error()
			results = append(results, result)
			continue
		}

		label := ""
		if item.Label != nil {
			label = strings.TrimSpace(*item.Label)
		}

		_, body, err := s.Store.AccrueWithLabel(r.Context(), strings.TrimSpace(*item.UserID), *item.Amount, ttl, strings.TrimSpace(*item.IdempotencyKey), label)
		if err != nil {
			result["status"] = "error"
			result["error"] = s.batchErrorCode(err)
			result["message"] = err.Error()
			results = append(results, result)
			continue
		}

		var payload any
		if len(body) > 0 {
			if err := json.Unmarshal(body, &payload); err != nil {
				payload = string(body)
			}
		}
		result["status"] = "created"
		result["result"] = payload
		result["user_id"] = strings.TrimSpace(*item.UserID)
		results = append(results, result)
	}

	writeJSON(w, http.StatusMultiStatus, map[string]any{"results": results})
}

func (s *Server) batchErrorCode(err error) string {
	var invalidLabelErr *data.InvalidLabelError
	switch {
	case err == nil:
		return ""
	case errors.Is(err, data.ErrInsufficientFunds), errors.Is(err, data.ErrIdempotencyConflict):
		return "conflict"
	case errors.Is(err, data.ErrInvalidAmount), errors.As(err, &invalidLabelErr):
		return "bad_request"
	case errors.Is(err, data.ErrUserNotFound), errors.Is(err, data.ErrHoldNotFound):
		return "not_found"
	default:
		return "internal_server_error"
	}
}

// handleBalance implements "по каждому пользователю можно узнать сколько у
// него баллов и сколько баллов сгорит в ближайшие дни".
func (s *Server) handleBalance(w http.ResponseWriter, r *http.Request) {
	userID := r.PathValue("id")

	days := 7
	if v := r.URL.Query().Get("expiring_within_days"); v != "" {
		n, err := strconv.Atoi(v)
		if err != nil || n < 0 {
			badRequest(w, "expiring_within_days must be a non-negative integer")
			return
		}
		days = n
	}

	res, err := s.Store.Balance(r.Context(), userID, days)
	if err != nil {
		s.respond(w, 0, nil, err)
		return
	}
	writeJSON(w, http.StatusOK, res)
}

func (s *Server) handleListLots(w http.ResponseWriter, r *http.Request) {
	userID := r.PathValue("id")

	page, offset, ok := parsePageAndOffset(w, r)
	if !ok {
		return
	}

	status := strings.TrimSpace(r.URL.Query().Get("status"))
	switch status {
	case "", data.LotStatusActive, data.LotStatusExpired, data.LotStatusExhausted:
	default:
		badRequest(w, "status must be one of: active, expired, exhausted")
		return
	}

	lots, err := s.Store.ListLots(r.Context(), userID, page, offset, status)
	if err != nil {
		s.respond(w, 0, nil, err)
		return
	}
	writeJSON(w, http.StatusOK, lots)
}

func (s *Server) handleListLedger(w http.ResponseWriter, r *http.Request) {
	userID := r.PathValue("id")

	page, offset, ok := parsePageAndOffset(w, r)
	if !ok {
		return
	}

	result, err := s.Store.ListLedger(r.Context(), userID, page, offset)
	if err != nil {
		s.respond(w, 0, nil, err)
		return
	}
	writeJSON(w, http.StatusOK, result)
}

type amountRequest struct {
	Amount         *int    `json:"amount"`
	IdempotencyKey *string `json:"idempotency_key"`
}

// handleCreateHold implements the first phase of "двухфазное списание":
// reserve points without permanently spending them.
func (s *Server) handleCreateHold(w http.ResponseWriter, r *http.Request) {
	userID := r.PathValue("id")

	var req amountRequest
	if err := readJSON(w, r, &req); err != nil {
		badRequest(w, err.Error())
		return
	}
	if req.Amount == nil {
		badRequest(w, "amount is required")
		return
	}
	if req.IdempotencyKey == nil || strings.TrimSpace(*req.IdempotencyKey) == "" {
		badRequest(w, "idempotency_key is required")
		return
	}
	if *req.Amount <= 0 {
		badRequest(w, "amount must be a positive integer")
		return
	}

	status, body, err := s.Store.CreateHold(r.Context(), userID, *req.Amount, strings.TrimSpace(*req.IdempotencyKey))
	s.respond(w, status, body, err)
}

// handleDebit implements a one-shot ("списать баллы") debit by performing a
// hold + confirm atomically.
func (s *Server) handleDebit(w http.ResponseWriter, r *http.Request) {
	userID := r.PathValue("id")

	var req amountRequest
	if err := readJSON(w, r, &req); err != nil {
		badRequest(w, err.Error())
		return
	}
	if req.Amount == nil {
		badRequest(w, "amount is required")
		return
	}
	if req.IdempotencyKey == nil || strings.TrimSpace(*req.IdempotencyKey) == "" {
		badRequest(w, "idempotency_key is required")
		return
	}
	if *req.Amount <= 0 {
		badRequest(w, "amount must be a positive integer")
		return
	}

	status, body, err := s.Store.Debit(r.Context(), userID, *req.Amount, strings.TrimSpace(*req.IdempotencyKey))
	s.respond(w, status, body, err)
}

func (s *Server) parseHoldID(w http.ResponseWriter, r *http.Request) (int64, bool) {
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil || id <= 0 {
		badRequest(w, "invalid hold id")
		return 0, false
	}
	return id, true
}

// handleConfirmHold implements the second phase: permanently spend a held
// amount. This is what the calling service should invoke once its own
// transaction has committed successfully.
func (s *Server) handleConfirmHold(w http.ResponseWriter, r *http.Request) {
	holdID, ok := s.parseHoldID(w, r)
	if !ok {
		return
	}
	status, body, err := s.Store.ConfirmHold(r.Context(), holdID)
	s.respond(w, status, body, err)
}

// handleCancelHold releases a hold's points back to the user's balance.
// This is the fix for "если вызывающий сервис ... аварийно завершился ...
// баллы не возвращаются владельцу": holds that are never confirmed can be
// cancelled (e.g. by a timeout/reconciliation job) to release the points.
func (s *Server) handleCancelHold(w http.ResponseWriter, r *http.Request) {
	holdID, ok := s.parseHoldID(w, r)
	if !ok {
		return
	}
	status, body, err := s.Store.CancelHold(r.Context(), holdID)
	s.respond(w, status, body, err)
}
