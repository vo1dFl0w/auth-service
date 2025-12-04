package logger

import (
	"log/slog"
	"os"
	"time"
)

var (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"
)

var (
	logTimeFormat = "02/01/2006 15:04:05"
)

func LoadLogger(env string) *slog.Logger {
	var log *slog.Logger

	opts := slog.HandlerOptions{
		ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
			if v := a.Value.Any(); v != nil {
				if a.Key == slog.TimeKey {
					if t, ok := v.(time.Time); ok {
						return slog.String(slog.TimeKey, t.Format(logTimeFormat))
					}
				}
			}
			return a
		},
	}

	switch env {
	case envLocal:
		log = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug, ReplaceAttr: opts.ReplaceAttr}))
	case envDev:
		log = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug, ReplaceAttr: opts.ReplaceAttr}))
	case envProd:
		log = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo, ReplaceAttr: opts.ReplaceAttr}))
	default:
		log = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo, ReplaceAttr: opts.ReplaceAttr}))
	}

	return log
}
