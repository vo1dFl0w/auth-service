package http

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/google/uuid"
)

type ctxKey string

const (
	CtxKeyUserID       ctxKey = "user_id"
	CtxKeyRefreshToken ctxKey = "refresh_token"
	CtxKeyRequestID    ctxKey = "request_id"
)

func (h *Handler) CorsMiddleware(next http.Handler) http.Handler {
	return h.cors.Handler(next)
}

func (h *Handler) LoggerMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		log := h.log.With(
			"request_id", r.Context().Value(CtxKeyRequestID),
			"remote_addr", r.RemoteAddr,
			"http-method", r.Method,
			"path", r.URL.Path,
		)

		log.Info("started")

		rw := &responseWriter{w, http.StatusOK}

		next.ServeHTTP(rw, r)

		completed := time.Since(start)
		completedStr := fmt.Sprintf("%.3fms", float64(completed.Microseconds())/1000)

		attrs := []any{
			"code", rw.code,
			"status-text", http.StatusText(rw.code),
			"duration_ms", completedStr,
		}

		switch {
		case rw.code >= 500:
			log.Error("failed", attrs...)
		case rw.code >= 400:
			log.Warn("failed", attrs...)
		default:
			log.Info("completed", attrs...)
		}
	})
}

func (h *Handler) RequestIDMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		reqID := r.Header.Get("X-Request-ID")
		if reqID == "" {
			reqID = uuid.NewString()
		}

		w.Header().Set("X-Request-ID", reqID)

		next.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), CtxKeyRequestID, reqID)))
	})
}

func (h *Handler) TimeoutMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(r.Context(), time.Second*time.Duration(h.cfg.Server.RequestDuration))
		defer cancel()

		next.ServeHTTP(w, r.WithContext(ctx))
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
