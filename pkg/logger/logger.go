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
)

type Logger interface {
	With(args ...any) Logger
	Info(msg string, args ...any)
	Warn(msg string, args ...any)
	Error(msg string, args ...any)
}

type logger struct {
	slogger *slog.Logger
}

func (l *logger) With(args ...any) Logger {
	return &logger{slogger: l.slogger.With(args...)}
}

func (l *logger) Info(msg string, args ...any) {
	l.slogger.Info(msg, args...)
}

func (l *logger) Warn(msg string, args ...any) {
	l.slogger.Warn(msg, args...)
}

func (l *logger) Error(msg string, args ...any) {
	l.slogger.Error(msg, args...)
}

func LoadLogger(cfg *Config) *logger {
	var handler slog.Handler

	switch cfg.env {
	case envLocal:
		handler = tint.NewHandler(os.Stdout, &tint.Options{
			Level:      slog.LevelDebug,
			TimeFormat: cfg.loggerTimeFormat,
			AddSource:  false,
		})
	case envDev:
		handler = slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
			Level:       slog.LevelDebug,
			ReplaceAttr: setLoggerOptions(cfg.loggerTimeFormat),
			AddSource:   false,
		})
	case envProd:
		handler = slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
			Level:       slog.LevelInfo,
			ReplaceAttr: setLoggerOptions(cfg.loggerTimeFormat),
			AddSource:   false,
		})
	default:
		handler = slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
			Level:       slog.LevelInfo,
			ReplaceAttr: setLoggerOptions(cfg.loggerTimeFormat),
			AddSource:   false,
		})
	}

	return &logger{
		slogger: slog.New(handler),
	}
}

func setLoggerOptions(loggerTimeFormat string) func(groups []string, a slog.Attr) slog.Attr {
	return func(groups []string, a slog.Attr) slog.Attr {

		if a.Key == slog.TimeKey {
			if t, ok := a.Value.Any().(time.Time); ok {
				return slog.String(slog.TimeKey, t.Format(loggerTimeFormat))
			}
		}
		return a
	}
}