package logfx

import (
	"log/slog"
	"os"

	cfg "github.com/Oleexo/config-go"
	"go.uber.org/fx"
)

const (
	ModKey       = "LOGGER_MOD"
	PlainTextMod = "plain"
	JsonMod      = "json"
	DefaultMod   = JsonMod
)

func newLogger(config cfg.Configuration) *slog.Logger {
	mod := config.GetStringOrDefault(ModKey, DefaultMod)
	lvl := new(slog.LevelVar)
	lvl.Set(slog.LevelDebug)
	var handler slog.Handler
	if mod == PlainTextMod {
		handler = slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
			Level: lvl,
		})
	} else {
		handler = slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
			Level: lvl,
		})
	}
	logger := slog.New(handler)
	return logger
}

func BuildLoggerModule() fx.Option {
	return fx.Module("loggerfx",
		fx.Provide(newLogger),
	)
}
