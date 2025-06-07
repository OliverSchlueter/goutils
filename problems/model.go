package problems

import (
	"encoding/json"
	"github.com/OliverSchlueter/goutils/broker"
	"github.com/OliverSchlueter/goutils/sloki"
	"log/slog"
	"net/http"
	"time"
)

// Problem represents a structured error response in compliance with RFC 7807.
type Problem struct {
	Type      string    `json:"type"`
	Title     string    `json:"title"`
	Detail    string    `json:"detail"`
	Status    int       `json:"status"`
	Timestamp time.Time `json:"timestamp"`
}

func (p *Problem) WriteToHTTP(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/problem+json")
	w.WriteHeader(p.Status)

	data, err := json.Marshal(p)
	if err != nil {
		slog.Warn("failed to marshal problem response", sloki.WrapError(err))
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	_, _ = w.Write(data)
}

func (p *Problem) WriteToBroker(b broker.Broker, subj string) {
	data, err := json.Marshal(p)
	if err != nil {
		slog.Warn("failed to marshal problem response", sloki.WrapError(err))
		return
	}

	if err := b.Publish(subj, data); err != nil {
		slog.Error("failed to publish problem response", sloki.WrapError(err), "subject", subj)
		return
	}
}
