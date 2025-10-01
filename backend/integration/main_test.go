//go:build integration

package integration

import (
	"context"
	"database/sql"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"testing"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"

	"github.com/mistribe/subtracker/internal/adapters/persistence/db"
	"github.com/mistribe/subtracker/pkg/testx"
)

var (
	pgC           *postgres.PostgresContainer
	testDBContext *db.Context
	databaseReady bool
)

// resolveMigrationsDir returns the absolute path to the database migrations folder
func resolveMigrationsDir() string {
	_, thisFile, _, ok := runtime.Caller(0)
	if !ok {
		log.Fatalf("cannot determine caller file path for migrations resolution")
	}
	root := filepath.Dir(filepath.Dir(thisFile)) // go up from integration/ to repo root
	return filepath.Join(root, "database")
}

// TestMain sets up a single Postgres container (shared across all integration tests)
func TestMain(m *testing.M) {
	ctx := context.Background()

	goose.SetTableName("goose.migrations")

	pc, err := postgres.Run(ctx,
		"postgres:16-alpine",
		postgres.WithDatabase("app"),
		postgres.WithUsername("postgres"),
		postgres.WithPassword("postgres"),
		postgres.WithInitScripts(), // none but keeps pattern consistent
		postgres.WithWaitStrategy(wait.ForListeningPort("5432/tcp").WithStartupTimeout(60*time.Second)),
	)
	if err != nil {
		log.Fatalf("failed to start postgres container: %v", err)
	}
	pgC = pc

	dsn, err := pgC.ConnectionString(ctx, "sslmode=disable")
	if err != nil {
		log.Fatalf("failed to get connection string: %v", err)
	}

	dbConn, err := sql.Open("pgx", dsn)
	if err != nil {
		log.Fatalf("failed to open db: %v", err)
	}
	defer dbConn.Close()

	migrationsDir := resolveMigrationsDir()
	if err := goose.Up(dbConn, migrationsDir); err != nil {
		log.Fatalf("failed to run migrations (%s): %v", migrationsDir, err)
	}

	testDBContext = db.NewContextFromDSN(dsn, testx.DiscardLogger())
	databaseReady = true

	code := m.Run()

	if pgC != nil {
		_ = pgC.Terminate(ctx) // best effort
	}
	os.Exit(code)
}

// GetDBContext returns the initialized db context (panic if not ready)
func GetDBContext() *db.Context {
	if !databaseReady || testDBContext == nil {
		panic("database context not initialized - ensure tests are run with -tags=integration and TestMain executed")
	}
	return testDBContext
}
