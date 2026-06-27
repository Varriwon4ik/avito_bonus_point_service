package api

import (
	"encoding/json"
	"errors"
	"io"
	"log/slog"
	"net/http"
	"strings"

	"bonus-ledger/internal/data"
)

// Server wires the HTTP API to the persistent Store.
type Server struct {
	Store          *data.Store
	Logger         *slog.Logger
	DefaultTTLDays int
	Mux            *http.ServeMux
	Metrics        *Metrics
	// AdminToken is the bearer token required for privileged admin operations
	// (e.g. manual accrual, US-07). When empty, admin authentication is
	// disabled and those endpoints are open.
	AdminToken string
}

func NewServer(store *data.Store, logger *slog.Logger, defaultTTLDays int) *Server {
	s := &Server{
		Store:          store,
		Logger:         logger,
		DefaultTTLDays: defaultTTLDays,
		Mux:            http.NewServeMux(),
		Metrics:        NewMetrics(),
	}
	s.routes()
	return s
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.Mux.ServeHTTP(w, r)
}

type errorEnvelope struct {
	Error   string `json:"error"`
	Message string `json:"message"`
}

func writeJSON(w http.ResponseWriter, status int, body any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if body == nil {
		return
	}
	_ = json.NewEncoder(w).Encode(body)
}

func writeJSONBytes(w http.ResponseWriter, status int, body []byte) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if len(body) == 0 {
		return
	}
	_, _ = w.Write(body)
}

func errorCode(status int) string {
	switch status {
	case http.StatusBadRequest:
		return "bad_request"
	case http.StatusUnauthorized:
		return "unauthorized"
	case http.StatusNotFound:
		return "not_found"
	case http.StatusConflict:
		return "conflict"
	case http.StatusMethodNotAllowed:
		return "method_not_allowed"
	default:
		return "internal_server_error"
	}
}

func writeError(w http.ResponseWriter, status int, message string) {
	writeJSON(w, status, errorEnvelope{
		Error:   errorCode(status),
		Message: message,
	})
}

func badRequest(w http.ResponseWriter, message string) {
	writeError(w, http.StatusBadRequest, message)
}

func unauthorized(w http.ResponseWriter, message string) {
	writeError(w, http.StatusUnauthorized, message)
}

func notFound(w http.ResponseWriter, message string) {
	writeError(w, http.StatusNotFound, message)
}

func conflict(w http.ResponseWriter, message string) {
	writeError(w, http.StatusConflict, message)
}

func internalServerError(w http.ResponseWriter) {
	writeError(w, http.StatusInternalServerError, "internal server error")
}

// readJSON decodes a JSON request body, rejecting unknown fields and bodies
// over 1MB.
func readJSON(w http.ResponseWriter, r *http.Request, dst any) error {
	r.Body = http.MaxBytesReader(w, r.Body, 1<<20)
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()

	if err := dec.Decode(dst); err != nil {
		if errors.Is(err, io.EOF) {
			return errors.New("request body must not be empty")
		}
		return explainJSONError(err)
	}

	if err := dec.Decode(&struct{}{}); !errors.Is(err, io.EOF) {
		return errors.New("request body must contain a single JSON object")
	}
	return nil
}

func explainJSONError(err error) error {
	var syntaxErr *json.SyntaxError
	var typeErr *json.UnmarshalTypeError
	var maxBytesErr *http.MaxBytesError

	switch {
	case errors.As(err, &syntaxErr):
		return errors.New("request body contains malformed JSON")
	case errors.Is(err, io.ErrUnexpectedEOF):
		return errors.New("request body contains malformed JSON")
	case errors.As(err, &typeErr):
		if typeErr.Field != "" {
			return errors.New("request body contains an invalid value for " + typeErr.Field)
		}
		return errors.New("request body must be a JSON object")
	case strings.HasPrefix(err.Error(), "json: unknown field "):
		field := strings.TrimPrefix(err.Error(), "json: unknown field ")
		return errors.New("request body contains unknown field " + field)
	case errors.As(err, &maxBytesErr):
		return errors.New("request body must not be larger than 1 MB")
	default:
		return errors.New("request body contains malformed JSON")
	}
}

// respond maps a Store result/error pair to an HTTP response.
func (s *Server) respond(w http.ResponseWriter, status int, body []byte, err error) {
	if err != nil {
		switch {
		case errors.Is(err, data.ErrInsufficientFunds):
			conflict(w, err.Error())
		case errors.Is(err, data.ErrInvalidAmount):
			badRequest(w, err.Error())
		case errors.Is(err, data.ErrUserNotFound):
			notFound(w, err.Error())
		case errors.Is(err, data.ErrHoldNotFound):
			notFound(w, err.Error())
		case errors.Is(err, data.ErrInvalidHoldStatus):
			conflict(w, err.Error())
		case errors.Is(err, data.ErrIdempotencyConflict):
			conflict(w, err.Error())
		default:
			s.Logger.Error("internal error", "err", err)
			internalServerError(w)
		}
		return
	}

	writeJSONBytes(w, status, body)
}
