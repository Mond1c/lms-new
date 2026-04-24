package repo

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"testing"
	"time"

	"github.com/Mond1c/lms/pkg/pg"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
)

var testDB *pgxpool.Pool

func TestMain(m *testing.M) {
	ctx := context.Background()

	_, thisFile, _, _ := runtime.Caller(0)
	migrationsDir := filepath.Join(filepath.Dir(thisFile), "..", "..", "migrations")

	pgc, err := postgres.Run(ctx,
		"postgres:16-alpine",
		postgres.WithDatabase("lms_test"),
		postgres.WithUsername("lms"),
		postgres.WithPassword("lms"),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2).
				WithStartupTimeout(30*time.Second),
		),
	)
	if err != nil {
		_, _ = fmt.Fprintln(os.Stderr, "start postgres:", err)
		os.Exit(1)
	}

	dsn, err := pgc.ConnectionString(ctx, "sslmode=disable")
	if err != nil {
		_, _ = fmt.Fprintln(os.Stderr, "dsn:", err)
		os.Exit(1)
	}

	if err := pg.Migrate(dsn, migrationsDir); err != nil {
		_, _ = fmt.Fprintln(os.Stderr, "migrate:", err)
		os.Exit(1)
	}

	pool, err := pg.Open(ctx, pg.Config{DSN: dsn, MaxConnections: 5})
	if err != nil {
		_, _ = fmt.Fprintln(os.Stderr, "open pool:", err)
		os.Exit(1)
	}
	testDB = pool

	code := m.Run()

	pool.Close()
	_ = pgc.Terminate(ctx)

	os.Exit(code)
}

func truncate(t *testing.T) {
	t.Helper()
	_, err := testDB.Exec(context.Background(), `
		TRUNCATE TABLE users RESTART IDENTITY CASCADE;
	`)
	if err != nil {
		t.Fatalf("truncate: %v", err)
	}
}
