package middleware

import (
	"log/slog"
	"net/http"
	"time"
)

// responseWriter adalah wrapper untuk http.ResponseWriter
// agar bisa capture status code
type responseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}

// LoggingMiddleware mencatat setiap request yang masuk
func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		// Wrap response writer untuk capture status code
		wrapped := &responseWriter{
			ResponseWriter: w,
			statusCode:     http.StatusOK, // default
		}

		// Jalankan handler
		next.ServeHTTP(wrapped, r)

		// Log setelah request selesai
		duration := time.Since(start)

		slog.Info("Request processed",
			"method", r.Method,
			"path", r.URL.Path,
			"status", wrapped.statusCode,
			"duration_ms", duration.Milliseconds(),
			"ip", r.RemoteAddr,
			"user_agent", r.UserAgent(),
		)
	})
}
