package problems

import (
	"encoding/json"
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

func (p *Problem) MarshalJSON() ([]byte, error) {
	return json.Marshal(p)
}
