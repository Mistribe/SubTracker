package db

import (
	"context"
	"database/sql"
	"log/slog"
	"strings"

	"github.com/Oleexo/config-go"
	"github.com/go-jet/jet/v2/postgres"
	"github.com/go-jet/jet/v2/qrm"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/stdlib"
)

const (
	TransactionKey = "X-pgx-transaction"
)

type Context struct {
	pgxConfig *pgx.ConnConfig
	logger    *slog.Logger
}

func NewContext(
	cfg config.Configuration,
	logger *slog.Logger) *Context {
	dsn := cfg.GetString("DATABASE_DSN")

	pgxConfig, err := pgx.ParseConfig(dsn)
	if err != nil {
		panic(err)
	}

	return &Context{
		pgxConfig: pgxConfig,
		logger:    logger,
	}
}

// NewContextFromDSN creates a new db context from a raw DSN (used mainly for integration tests)
func NewContextFromDSN(dsn string, logger *slog.Logger) *Context {
	pgxConfig, err := pgx.ParseConfig(dsn)
	if err != nil {
		panic(err)
	}
	return &Context{pgxConfig: pgxConfig, logger: logger}
}

func (r *Context) logDebugSql(ctx context.Context, stmt postgres.Statement) {
	if r.logger.Enabled(ctx, slog.LevelDebug) {
		debugSql := stmt.DebugSql()
		debugSql = strings.ReplaceAll(debugSql, "\n", "")
		r.logger.Log(ctx, slog.LevelDebug, debugSql)
	}
}

func (r *Context) Query(ctx context.Context, stmt postgres.SelectStatement, s interface{}) error {
	var queryable qrm.Queryable
	tx, ok := ctx.Value(TransactionKey).(*sql.Tx)
	if ok {
		queryable = tx
	} else {
		sqlDB := stdlib.OpenDB(*r.pgxConfig) // returns *sql.DB
		defer sqlDB.Close()
		queryable = sqlDB
	}

	r.logDebugSql(ctx, stmt)

	return stmt.QueryContext(ctx, queryable, s)
}

func (r *Context) Execute(ctx context.Context, stmt postgres.Statement) (int64, error) {
	var queryable qrm.Executable
	tx, ok := ctx.Value(TransactionKey).(*sql.Tx)
	if ok {
		queryable = tx
	} else {
		sqlDB := stdlib.OpenDB(*r.pgxConfig) // returns *sql.DB
		defer sqlDB.Close()
		queryable = sqlDB
	}

	r.logDebugSql(ctx, stmt)

	result, err := stmt.ExecContext(ctx, queryable)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected()
}
