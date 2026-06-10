package api

import (
	"encoding/json"
	"errors"
	"io"
	"log/slog"
	"net/http"

	"bonus-ledger/internal/data"
)

// Server wires the HTTP API to the persistent Store.
type Server struct {
	Store          *data.Store
	Logger         *slog.Logger
	DefaultTTLDays int
	Mux            *http.ServeMux
}

func NewServer(store *data.Store, logger *slog.Logger, defaultTTLDays int) *Server {
	s := &Server{
		Store:          store,
		Logger:         logger,
		DefaultTTLDays: defaultTTLDays,
		Mux:            http.NewServeMux(),
	}
	s.routes()
	return s
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.Mux.ServeHTTP(w, r)
}

type errorEnvelope struct {
	Error string `json:"error"`
}

func writeError(w http.ResponseWriter, status int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(errorEnvelope{Error: message})
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
		return err
	}

	if err := dec.Decode(&struct{}{}); !errors.Is(err, io.EOF) {
		return errors.New("request body must contain a single JSON object")
	}
	return nil
}

// respond maps a Store result/error pair to an HTTP response.
func (s *Server) respond(w http.ResponseWriter, status int, body []byte, err error) {
	if err != nil {
		switch {
		case errors.Is(err, data.ErrInsufficientFunds):
			writeError(w, http.StatusConflict, err.Error())
		case errors.Is(err, data.ErrInvalidAmount):
			writeError(w, http.StatusBadRequest, err.Error())
		case errors.Is(err, data.ErrNotFound):
			writeError(w, http.StatusNotFound, err.Error())
		case errors.Is(err, data.ErrInvalidHoldStatus):
			writeError(w, http.StatusConflict, err.Error())
		case errors.Is(err, data.ErrIdempotencyConflict):
			writeError(w, http.StatusConflict, err.Error())
		default:
			s.Logger.Error("internal error", "err", err)
			writeError(w, http.StatusInternalServerError, "internal server error")
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_, _ = w.Write(body)
}
