package http

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"time"
)

type ctxKey string

const (
	CtxKeyUserID       ctxKey = "user_id"
	CtxKeyRefreshToken ctxKey = "refresh_token"
)

func (h *Handler) TimeoutMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(r.Context(), time.Second*time.Duration(h.cfg.Server.RequestDuration))
		defer cancel()

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (h *Handler) CorsMiddleware(next http.Handler) http.Handler {
	return h.cors.Handler(next)
}

func (h *Handler) LoggerMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		log := h.log.With(
			"remote_addr", r.RemoteAddr,
			"http-method", r.Method,
			"path", r.URL.Path,
		)

		log.Info("started")

		rw := &responseWriter{w, http.StatusOK}

		next.ServeHTTP(rw, r)

		var level slog.Level
		switch {
		case rw.code >= 500:
			level = slog.LevelError
		case rw.code >= 400:
			level = slog.LevelWarn
		default:
			level = slog.LevelInfo
		}

		complited := time.Since(start)
		complitedStr := fmt.Sprintf("%.3fms", float64(complited.Microseconds())/1000)

		log.Info(
			"completed",
			slog.Any(slog.LevelKey, level),
			slog.Int("code", rw.code),
			slog.String("status-text", http.StatusText(rw.code)),
			slog.String("time", complitedStr),
		)
	})
}

type responseWriter struct {
	http.ResponseWriter
	code int
}

func (rw *responseWriter) WriteHeader(statusCode int) {
	rw.code = statusCode
	rw.ResponseWriter.WriteHeader(statusCode)
}
