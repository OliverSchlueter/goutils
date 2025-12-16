package ratelimit

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/OliverSchlueter/goutils/problems"
)

var (
	ErrRateLimitExceeded = errors.New("rate limit exceeded")
)

func RateLimitExceededProblem() *problems.Problem {
	return &problems.Problem{
		Type:      "RateLimitExceeded",
		Title:     "Rate limit exceeded",
		Detail:    fmt.Sprintf("You have exceeded your rate limit."),
		Status:    http.StatusTooManyRequests,
		Timestamp: time.Now(),
	}
}
