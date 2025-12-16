package ratelimit

import (
	"errors"
	"log/slog"
	"net/http"
	"time"

	"github.com/OliverSchlueter/goutils/ratelimit/database/memory"
	"github.com/OliverSchlueter/goutils/sloki"
)

type DB interface {
	GetTokens(client string) (int, error)
	SetTokens(client string, tokens int) error
	GetLastRefill(client string) (time.Time, error)
	SetLastRefill(client string, t time.Time) error
}

type Service struct {
	db              DB
	tokensPerSecond float64
	maxTokens       int
	getIP           func(r *http.Request) string
}

type Configuration struct {
	DB              DB
	TokensPerSecond float64
	MaxTokens       int
	GetIP           func(r *http.Request) string
}

func NewService(config Configuration) *Service {
	if config.DB == nil {
		config.DB = memory.NewDB()
	}

	if config.GetIP == nil {
		config.GetIP = GetIP
	}

	return &Service{
		db:              config.DB,
		tokensPerSecond: config.TokensPerSecond,
		maxTokens:       config.MaxTokens,
		getIP:           config.GetIP,
	}
}

func (s *Service) CheckAndConsume(client string) error {
	tokens, err := s.db.GetTokens(client)
	if err != nil {
		return err
	}

	// Refill tokens based on time elapsed
	lastRefill, err := s.db.GetLastRefill(client)
	if err != nil {
		return err
	}

	now := time.Now()
	elapsed := now.Sub(lastRefill).Seconds()

	refillTokens := int(elapsed * s.tokensPerSecond)
	if refillTokens > 0 {
		tokens += refillTokens
		if tokens > s.maxTokens {
			tokens = s.maxTokens
		}
		if err := s.db.SetLastRefill(client, now); err != nil {
			return err
		}
		if err := s.db.SetTokens(client, tokens); err != nil {
			return err
		}
	}

	// Consume a token
	if tokens <= 0 {
		return ErrRateLimitExceeded
	}

	return s.db.SetTokens(client, tokens-1)
}

func (s *Service) CheckRequest(r *http.Request, resource string) error {
	client := s.getIP(r) + "--" + resource
	return s.CheckAndConsume(client)
}

func (s *Service) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if err := s.CheckRequest(r, "*"); err != nil {
			if errors.Is(err, ErrRateLimitExceeded) {
				RateLimitExceededProblem().WriteToHTTP(w)
				return
			}

			slog.Error("Rate limit check failed", sloki.WrapError(err))
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		next.ServeHTTP(w, r)
	})
}
