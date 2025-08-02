package persistence

import (
	"context"
	"errors"
	"log/slog"
	"time"

	"github.com/Oleexo/config-go"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/oleexo/subtracker/internal/infrastructure/persistence/sql"
)

const (
	TransactionKey = "X-pgx-transaction"
)

type DatabaseContext struct {
	pool    *pgxpool.Pool
	queries *sql.Queries
}

func NewDatabaseContext(
	cfg config.Configuration,
	logger *slog.Logger) *DatabaseContext {
	dsn := cfg.GetString("DATABASE_DSN")

	// Parse pgxConfig and set pool settings
	pgxConfig, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		panic(err)
	}

	// Configure connection pool settings
	pgxConfig.MaxConns = 25                      // Maximum number of connections
	pgxConfig.MinConns = 5                       // Minimum number of connections to maintain
	pgxConfig.MaxConnLifetime = time.Hour        // Maximum connection lifetime
	pgxConfig.MaxConnIdleTime = 30 * time.Minute // Maximum connection idle time

	// Create connection pool
	pool, err := pgxpool.NewWithConfig(context.Background(), pgxConfig)
	if err != nil {
		panic(err)
	}

	// Create sqlc queries instance with the pool
	queries := sql.New(pool)

	return &DatabaseContext{
		pool:    pool,
		queries: queries,
	}
}

func (r *DatabaseContext) GetPool() *pgxpool.Pool {
	return r.pool
}

func (r *DatabaseContext) GetQueries(ctx context.Context) *sql.Queries {
	tx, ok := ctx.Value(TransactionKey).(pgx.Tx)
	if ok {
		return r.queries.WithTx(tx)
	}
	return r.queries
}

func (r *DatabaseContext) Close() {
	if r.pool != nil {
		r.pool.Close()
	}
}

func (r *DatabaseContext) Ping(ctx context.Context) error {
	if r.pool != nil {
		return r.pool.Ping(ctx)
	}
	return errors.New("database pool is nil")
}
