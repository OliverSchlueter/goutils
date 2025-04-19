package middleware

import (
	"github.com/OliverSchlueter/goutils/sloki"
	"log/slog"
	"net/http"
	"time"
)

type StatusRecorder struct {
	http.ResponseWriter
	Status int
}

func (s *StatusRecorder) WriteHeader(code int) {
	s.Status = code
	s.ResponseWriter.WriteHeader(code)
}

func Logging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		startTime := time.Now()

		sr := &StatusRecorder{
			ResponseWriter: w,
			Status:         http.StatusOK,
		}

		next.ServeHTTP(sr, r)

		elapsedTime := time.Since(startTime)

		slog.Info(
			"Request received",
			sloki.WrapRequest(r),
			slog.Int("status", sr.Status),
			slog.Int64("elapsed_time", elapsedTime.Milliseconds()),
		)
	})
}
