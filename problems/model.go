package problems

import "encoding/json"

// Problem represents a structured error response in compliance with RFC 7807.
type Problem struct {
	Type     string `json:"type"`
	Title    string `json:"title"`
	Detail   string `json:"detail"`
	Status   int    `json:"status"`
	Instance string `json:"instance"`
}

func (p *Problem) MarshalJSON() ([]byte, error) {
	return json.Marshal(p)
}
