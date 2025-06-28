package middleware

import (
	"github.com/OliverSchlueter/goutils/sloki"
	"log/slog"
	"net/http"
	"time"
)

// LogLevel defines the logging level for request logging.
var LogLevel = slog.LevelInfo

// OnlyLogStatusAbove defines the threshold for status codes above which logging will not occur.
var OnlyLogStatusAbove = 0

type StatusRecorder struct {
	http.ResponseWriter
	Status int
}

func (s *StatusRecorder) WriteHeader(code int) {
	s.Status = code
	s.ResponseWriter.WriteHeader(code)
}

func RequestLogging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		startTime := time.Now()

		sr := &StatusRecorder{
			ResponseWriter: w,
			Status:         http.StatusOK,
		}

		next.ServeHTTP(sr, r)

		if sr.Status < OnlyLogStatusAbove {
			// If the status code is above the threshold, do not log
			return
		}

		elapsedTime := time.Since(startTime)

		slog.Log(
			r.Context(),
			LogLevel,
			"RequestLogging received",
			sloki.WrapRequest(r),
			slog.Int("status", sr.Status),
			slog.Int64("elapsed_time", elapsedTime.Milliseconds()),
		)
	})
}
