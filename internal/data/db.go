package data

import (
	"context"
	"database/sql"
	"embed"
	"fmt"
	"sort"
	"strings"
	"time"
)

//go:embed migrations/*.sql
var migrationsFS embed.FS

// OpenDB opens a connection pool to Postgres and verifies connectivity.
func OpenDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, err
	}

	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(25)
	db.SetConnMaxIdleTime(5 * time.Minute)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := db.PingContext(ctx); err != nil {
		db.Close()
		return nil, err
	}

	return db, nil
}

// migrateLockID is an arbitrary application-chosen key for the session-level
// Postgres advisory lock that serializes concurrent Migrate calls.
const migrateLockID = 792101

// Migrate applies all *.up.sql migrations in order. It is safe to call on
// every startup since the migrations only use IF NOT EXISTS / IF EXISTS.
// A Postgres advisory lock serializes concurrent callers, so several API
// replicas starting at the same time against a fresh database do not race
// (concurrent CREATE TABLE IF NOT EXISTS can otherwise fail in Postgres).
func Migrate(db *sql.DB) error {
	// The advisory lock is session-scoped, so acquisition and release must
	// happen on the same pooled connection.
	conn, err := db.Conn(context.Background())
	if err != nil {
		return err
	}
	defer conn.Close()

	if _, err := conn.ExecContext(context.Background(),
		`SELECT pg_advisory_lock($1)`, migrateLockID); err != nil {
		return err
	}
	defer conn.ExecContext(context.Background(),
		`SELECT pg_advisory_unlock($1)`, migrateLockID)

	entries, err := migrationsFS.ReadDir("migrations")
	if err != nil {
		return err
	}

	var names []string
	for _, e := range entries {
		if !e.IsDir() {
			names = append(names, e.Name())
		}
	}
	sort.Strings(names)

	for _, name := range names {
		if !strings.HasSuffix(name, ".up.sql") {
			continue
		}
		contents, err := migrationsFS.ReadFile("migrations/" + name)
		if err != nil {
			return err
		}
		if _, err := db.Exec(string(contents)); err != nil {
			return fmt.Errorf("migration %s: %w", name, err)
		}
	}
	return nil
}
