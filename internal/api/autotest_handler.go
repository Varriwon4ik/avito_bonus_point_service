package api

import (
	"net/http"

	"bonus-ledger/internal/autotest"
	"bonus-ledger/internal/data"
)

// autotestRunRequest is the payload the web "Autotester" tab (US-17) submits:
// the same items the cmd/autotest console tool asks for. amount is required;
// the rest fall back to sensible defaults.
type autotestRunRequest struct {
	Label            string `json:"label"`
	UserID           string `json:"user_id"`
	Amount           *int   `json:"amount"`
	TTLDays          *int   `json:"ttl_days,omitempty"`
	ParallelRequests *int   `json:"parallel_requests,omitempty"`
}

type autotestCheckResult struct {
	Name    string `json:"name"`
	Passed  bool   `json:"passed"`
	Message string `json:"message,omitempty"`
}

// autotestScenarioEcho is the normalized scenario returned to the caller. It
// omits the DB-only bookkeeping fields of data.AutotestScenario since a UI run
// is not persisted.
type autotestScenarioEcho struct {
	Label            string `json:"label"`
	UserID           string `json:"user_id"`
	Amount           int    `json:"amount"`
	TTLDays          int    `json:"ttl_days"`
	ParallelRequests int    `json:"parallel_requests"`
	LedgerLabel      string `json:"ledger_label"`
}

type autotestRunResponse struct {
	Scenario autotestScenarioEcho  `json:"scenario"`
	Passed   bool                  `json:"passed"`
	Results  []autotestCheckResult `json:"results"`
}

// handleAutotestRun runs the shared autotester engine against this very
// instance and returns a per-check pass/fail report, so an administrator can
// exercise the accrual and parallel-request behaviour straight from the UI.
func (s *Server) handleAutotestRun(w http.ResponseWriter, r *http.Request) {
	var req autotestRunRequest
	if err := readJSON(w, r, &req); err != nil {
		badRequest(w, err.Error())
		return
	}
	if req.Amount == nil {
		badRequest(w, "amount is required")
		return
	}

	scn := data.AutotestScenario{
		Label:  req.Label,
		UserID: req.UserID,
		Amount: *req.Amount,
	}
	if req.TTLDays != nil {
		scn.TTLDays = *req.TTLDays
	} else {
		scn.TTLDays = s.DefaultTTLDays
	}
	if req.ParallelRequests != nil {
		scn.ParallelRequests = *req.ParallelRequests
	}

	scn = autotest.NormalizeScenario(scn)
	if err := autotest.Validate(scn); err != nil {
		badRequest(w, err.Error())
		return
	}

	rt := autotest.NewRuntime(selfBaseURL(r))
	if err := rt.CheckHealth(); err != nil {
		s.Logger.Error("autotest health check failed", "err", err)
		internalServerError(w)
		return
	}

	resp := autotestRunResponse{
		Scenario: autotestScenarioEcho{
			Label:            scn.Label,
			UserID:           scn.UserID,
			Amount:           scn.Amount,
			TTLDays:          scn.TTLDays,
			ParallelRequests: scn.ParallelRequests,
			LedgerLabel:      scn.LedgerLabel,
		},
		Passed: true,
	}
	for _, check := range autotest.Checks() {
		result := autotestCheckResult{Name: check.Name, Passed: true}
		if err := check.Run(rt, scn); err != nil {
			result.Passed = false
			result.Message = err.Error()
			resp.Passed = false
		}
		resp.Results = append(resp.Results, result)
	}

	writeJSON(w, http.StatusOK, resp)
}

// selfBaseURL reconstructs the base URL the request arrived on so the
// autotester can drive the running instance through the full HTTP stack.
func selfBaseURL(r *http.Request) string {
	scheme := "http"
	if r.TLS != nil {
		scheme = "https"
	}
	if proto := r.Header.Get("X-Forwarded-Proto"); proto != "" {
		scheme = proto
	}
	return scheme + "://" + r.Host
}
