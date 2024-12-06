package middleware

import (
	"fmt"
	"log/slog"
	"net/http"
	"time"
)

type Middleware struct {
	log *slog.Logger
}

func New(log *slog.Logger) *Middleware {
	return &Middleware{log: log}
}

func(m *Middleware) Logger(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log := m.log.With(
			slog.String("component", "middleware.logger"),
			slog.String("method", r.Method),
			slog.String("path", r.URL.Path),
			slog.String("remote_addr", r.RemoteAddr),
			slog.String("user-agent", r.UserAgent()), 
		)
		start := time.Now()
		next.ServeHTTP(w, r)
		log.Info("request completed", slog.String("duration", fmt.Sprint(time.Since(start))))
	})
}