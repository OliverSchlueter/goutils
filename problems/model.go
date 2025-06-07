package problems

import (
	"encoding/json"
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

func (p *Problem) Send(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/problem+json")
	w.WriteHeader(p.Status)

	data, err := json.Marshal(p)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	_, _ = w.Write(data)
}
