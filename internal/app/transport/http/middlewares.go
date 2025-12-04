package http

import (
	"encoding/json"
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

func (h *Handler) CorsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		next.ServeHTTP(w, r)
	})
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

func errorResponse(w http.ResponseWriter, r *http.Request, code int, err any) {
	response(w, r, code, err)
}

func response(w http.ResponseWriter, r *http.Request, code int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)

	if data != nil {
		json.NewEncoder(w).Encode(data)
	}
}
