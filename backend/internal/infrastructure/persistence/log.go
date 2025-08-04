package persistence

import (
	"context"
	"log/slog"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/tracelog"
)

func overridePgxLogger(config *pgxpool.Config,
	logger *slog.Logger) {
	config.ConnConfig.Tracer = &tracelog.TraceLog{
		Logger:   newpgxSlogLogger(logger),
		LogLevel: tracelog.LogLevelTrace,
	}
}

type pgxSlogLogger struct {
	logger *slog.Logger
}

func (p *pgxSlogLogger) Log(ctx context.Context, level tracelog.LogLevel, msg string, data map[string]any) {
	slogLevel := p.convertLogLevel(level)
	var args []interface{}
	for key, value := range data {
		args = append(args, slog.Any(key, value))
	}
	p.logger.Log(ctx, slogLevel, msg, args...)
}

func (p *pgxSlogLogger) convertLogLevel(level tracelog.LogLevel) slog.Level {
	switch level {
	case tracelog.LogLevelTrace:
		return slog.LevelDebug
	case tracelog.LogLevelDebug:
		return slog.LevelDebug
	case tracelog.LogLevelInfo:
		return slog.LevelInfo
	case tracelog.LogLevelWarn:
		return slog.LevelWarn
	case tracelog.LogLevelError:
		return slog.LevelError
	default:
		return slog.LevelDebug
	}
}

func newpgxSlogLogger(logger *slog.Logger) tracelog.Logger {
	return &pgxSlogLogger{
		logger: logger,
	}
}
