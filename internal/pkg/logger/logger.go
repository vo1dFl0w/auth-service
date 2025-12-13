package logger

import (
	"log/slog"
	"os"
	"time"

	"github.com/lmittmann/tint"
)

var (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"

	logTimeFormat = "02/01/2006 15:04:05"
)

func LoadLogger(env string) *slog.Logger {
	var handler slog.Handler

	switch env {
	case envLocal:
		handler = tint.NewHandler(os.Stdout, &tint.Options{
			Level:      slog.LevelDebug,
			TimeFormat: logTimeFormat,
		})

	case envDev:
		handler = slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
			Level: slog.LevelDebug,
			ReplaceAttr: setLoggerOptions,
		})

	case envProd:
		handler = slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
			Level: slog.LevelInfo,
			ReplaceAttr: setLoggerOptions,
		})

	default:
		handler = slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo})
	}

	return slog.New(handler)
}

func setLoggerOptions(groups []string, a slog.Attr) slog.Attr {
	if a.Key == slog.TimeKey {
		if t, ok := a.Value.Any().(time.Time); ok {
			return slog.String(slog.TimeKey, t.Format(logTimeFormat))
		}
	}
	return a
}
