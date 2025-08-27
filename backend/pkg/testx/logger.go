package testx

import (
	"io"
	"log/slog"
)

func DiscardLogger() *slog.Logger {
	return slog.New(slog.NewTextHandler(io.Discard, &slog.HandlerOptions{
		Level: slog.LevelDebug, // or LevelInfo
	}))
}
