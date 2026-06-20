package main

import (
	"context"
	"embed"
	"errors"
	"flag"
	"io/fs"
	"log/slog"
	"net/http"
	"os"
	"strconv"
	"time"

	_ "github.com/lib/pq"

	"bonus-ledger/internal/api"
	"bonus-ledger/internal/data"
)

//go:embed web
var webFS embed.FS

// holdSweepInterval is how often the background job checks for stale holds.
const holdSweepInterval = 5 * time.Minute

type config struct {
	port             int
	dsn              string
	defaultTTLDays   int
	holdTimeoutHours int
}

func main() {
	var cfg config

	flag.IntVar(&cfg.port, "port", 8080, "API server port")
	flag.StringVar(&cfg.dsn, "db-dsn", os.Getenv("DB_DSN"), "PostgreSQL DSN")
	defaultTTL := 365
	if v := os.Getenv("DEFAULT_TTL_DAYS"); v != "" {
		if n, err := strconv.Atoi(v); err == nil && n > 0 {
			defaultTTL = n
		}
	}
	flag.IntVar(&cfg.defaultTTLDays, "default-ttl-days", defaultTTL, "default lifetime of accrued points, in days")
	holdTimeoutHours := 24
	if v := os.Getenv("HOLD_TIMEOUT_HOURS"); v != "" {
		if n, err := strconv.Atoi(v); err == nil && n > 0 {
			holdTimeoutHours = n
		}
	}
	flag.IntVar(&cfg.holdTimeoutHours, "hold-timeout-hours", holdTimeoutHours,
		"holds left active/unresolved for longer than this are automatically released")
	flag.Parse()

	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	if cfg.dsn == "" {
		logger.Error("missing database DSN: set DB_DSN or pass -db-dsn")
		os.Exit(1)
	}

	db, err := data.OpenDB(cfg.dsn)
	if err != nil {
		logger.Error("failed to connect to database", "err", err)
		os.Exit(1)
	}
	defer db.Close()

	if err := data.Migrate(db); err != nil {
		logger.Error("failed to run migrations", "err", err)
		os.Exit(1)
	}
	logger.Info("database ready")

	store := data.NewStore(db)
	apiServer := api.NewServer(store, logger, cfg.defaultTTLDays)

	go runHoldSweep(context.Background(), store, cfg.holdTimeoutHours, logger)

	webRoot, err := fs.Sub(webFS, "web")
	if err != nil {
		logger.Error("failed to load embedded web assets", "err", err)
		os.Exit(1)
	}

	openAPISpec, err := loadOpenAPISpec()
	if err != nil {
		logger.Error("failed to load OpenAPI specification", "err", err)
		os.Exit(1)
	}

	srv := &http.Server{
		Addr:         ":" + strconv.Itoa(cfg.port),
		Handler:      api.NewAppHandler(apiServer, webRoot, openAPISpec),
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	logger.Info("starting server", "port", cfg.port, "default_ttl_days", cfg.defaultTTLDays,
		"hold_timeout_hours", cfg.holdTimeoutHours)
	err = srv.ListenAndServe()
	if !errors.Is(err, http.ErrServerClosed) {
		logger.Error("server error", "err", err)
		os.Exit(1)
	}
}

func loadOpenAPISpec() ([]byte, error) {
	paths := []string{
		"api/openapi.yaml",
		"/api/openapi.yaml",
	}

	for _, path := range paths {
		spec, err := os.ReadFile(path)
		if err == nil {
			return spec, nil
		}
		if !errors.Is(err, os.ErrNotExist) {
			return nil, err
		}
	}

	return nil, os.ErrNotExist
}

// runHoldSweep periodically releases active holds that have been left
// unresolved (never confirmed or cancelled) for longer than timeoutHours, so
// a crashed or buggy calling service can't permanently lock a user's points.
func runHoldSweep(ctx context.Context, store *data.Store, timeoutHours int, logger *slog.Logger) {
	ticker := time.NewTicker(holdSweepInterval)
	defer ticker.Stop()

	for {
		if n, err := store.ExpireStaleHolds(ctx, timeoutHours); err != nil {
			logger.Error("hold timeout sweep failed", "err", err)
		} else if n > 0 {
			logger.Info("auto-released stale holds", "count", n, "hold_timeout_hours", timeoutHours)
		}

		select {
		case <-ticker.C:
		case <-ctx.Done():
			return
		}
	}
}
