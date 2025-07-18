package sloki

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"strconv"
	"time"
)

type LogHandler interface {
	Handle(ctx context.Context, timestamp time.Time, level slog.Level, msg string, attr map[string]string) error
}

type Service struct {
	url          string
	service      string
	consoleLevel slog.Level
	lokiLevel    slog.Level
	enableLoki   bool
	httpClient   *http.Client
	handlers     []LogHandler
}

type Configuration struct {
	URL          string
	Service      string
	ConsoleLevel slog.Level
	LokiLevel    slog.Level
	EnableLoki   bool
	Handlers     []LogHandler
}

func NewService(cfg Configuration) *Service {
	if cfg.Handlers == nil {
		cfg.Handlers = []LogHandler{}
	}

	return &Service{
		url:          cfg.URL,
		service:      cfg.Service,
		consoleLevel: cfg.ConsoleLevel,
		lokiLevel:    cfg.LokiLevel,
		enableLoki:   cfg.EnableLoki,
		httpClient:   &http.Client{},
		handlers:     cfg.Handlers,
	}
}

func (s *Service) Enabled(_ context.Context, _ slog.Level) bool {
	return true
}

func (s *Service) printToConsole(level slog.Level) bool {
	return level >= s.consoleLevel
}

func (s *Service) sendToLoki(level slog.Level) bool {
	return s.enableLoki && level >= s.lokiLevel
}

func (s *Service) Handle(ctx context.Context, r slog.Record) error {
	attrs := map[string]string{}
	r.Attrs(func(a slog.Attr) bool {
		if a.Value.Kind() == slog.KindGroup {
			for _, gAttr := range a.Value.Group() {
				attrs[a.Key+"__"+gAttr.Key] = fmt.Sprint(gAttr.Value)
			}

			return true
		}

		attrs[a.Key] = fmt.Sprint(a.Value)
		return true
	})

	var attrJson []byte
	if len(attrs) > 0 {
		attrJson, _ = json.Marshal(attrs)
		attrJson = append([]byte(" "), attrJson...)
	}

	if s.printToConsole(r.Level) {
		fmt.Printf("%s [%s] %s%s\n",
			r.Time.Format("2006-01-02 15:04:05"),
			r.Level.String(),
			r.Message,
			string(attrJson),
		)
	}

	for _, h := range s.handlers {
		err := h.Handle(ctx, r.Time, r.Level, r.Message, attrs)
		if err != nil {
			fmt.Printf("Error in slog handler: %v\n", err)
			return err
		}
	}

	if !s.sendToLoki(r.Level) {
		return nil
	}

	unixTimestamp := strconv.FormatInt(r.Time.UnixNano(), 10)
	if err := s.pushLogToLoki(unixTimestamp, r.Level.String(), r.Message, attrs); err != nil {
		fmt.Printf("Failed to send log to Loki: %v\n", err)
		return err
	}

	return nil
}

func (s *Service) WithAttrs(_ []slog.Attr) slog.Handler {
	return s
}

func (s *Service) WithGroup(_ string) slog.Handler {
	return s
}

func (s *Service) pushLogToLoki(timestamp, level, message string, attrs map[string]string) error {
	labels := map[string]string{
		"service": s.service,
		"level":   level,
	}

	logEntry := map[string]interface{}{
		"timestamp": time.Now(),
	}
	for k, v := range attrs {
		logEntry[k] = v
	}

	req := PushLogsRequest{
		Streams: []Stream{
			{
				Labels: labels,
				Values: [][]any{
					{timestamp, message, logEntry},
				},
			},
		},
	}

	reqJson, err := json.Marshal(req)
	if err != nil {
		return fmt.Errorf("failed to marshal request: %w", err)
	}

	resp, err := s.httpClient.Post(s.url, "application/json", bytes.NewReader(reqJson))
	if err != nil {
		return fmt.Errorf("failed to send request to Loki: %w", err)
	}

	if resp.StatusCode != http.StatusNoContent {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("loki responded with status %d: %s", resp.StatusCode, string(body))
	}

	return nil
}
