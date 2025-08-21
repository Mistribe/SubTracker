package persistence

import (
	"context"
	"database/sql"
	"log/slog"

	"github.com/Oleexo/config-go"
	"github.com/go-jet/jet/v2/postgres"
	"github.com/go-jet/jet/v2/qrm"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/stdlib"
)

const (
	TransactionKey = "X-pgx-transaction"
)

type DatabaseContext struct {
	pgxConfig *pgx.ConnConfig
	logger    *slog.Logger
}

func NewDatabaseContext(
	cfg config.Configuration,
	logger *slog.Logger) *DatabaseContext {
	dsn := cfg.GetString("DATABASE_DSN")

	pgxConfig, err := pgx.ParseConfig(dsn)
	if err != nil {
		panic(err)
	}

	return &DatabaseContext{
		pgxConfig: pgxConfig,
		logger:    logger,
	}
}

func (r *DatabaseContext) Query(ctx context.Context, stmt postgres.SelectStatement, s interface{}) error {
	var queryable qrm.Queryable
	tx, ok := ctx.Value(TransactionKey).(*sql.Tx)
	if ok {
		queryable = tx
	} else {
		sqlDB := stdlib.OpenDB(*r.pgxConfig) // returns *sql.DB
		defer sqlDB.Close()
		queryable = sqlDB
	}

	if r.logger.Enabled(ctx, slog.LevelDebug) {
		r.logger.Log(ctx, slog.LevelDebug, stmt.DebugSql())
	}

	return stmt.QueryContext(ctx, queryable, s)
}

func (r *DatabaseContext) Execute(ctx context.Context, stmt postgres.Statement) (int64, error) {
	var queryable qrm.Executable
	tx, ok := ctx.Value(TransactionKey).(*sql.Tx)
	if ok {
		queryable = tx
	} else {
		sqlDB := stdlib.OpenDB(*r.pgxConfig) // returns *sql.DB
		defer sqlDB.Close()
		queryable = sqlDB
	}

	if r.logger.Enabled(ctx, slog.LevelDebug) {
		r.logger.Log(ctx, slog.LevelDebug, stmt.DebugSql())
	}

	result, err := stmt.ExecContext(ctx, queryable)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected()
}
