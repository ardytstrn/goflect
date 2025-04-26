package middleware

import (
	"net/http"

	"github.com/ardytstrn/goflect/internal/handlers"
)

type responseWriter struct {
	http.ResponseWriter
	status int
}

func (rw *responseWriter) WriteHeader(status int) {
	rw.status = status
	rw.ResponseWriter.WriteHeader(status)
}

func Chain(next http.Handler, app *handlers.App) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ww := &responseWriter{ResponseWriter: w}

		w.Header().Set("Content-Security-Policy", "default-src 'self")
		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.Header().Set("X-Frame-Options", "DENY")

		next.ServeHTTP(ww, r)
	})
}
