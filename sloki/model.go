package sloki

type PushLogsRequest struct {
	Streams []Stream `json:"streams"`
}

type Stream struct {
	Labels map[string]string `json:"stream"`
	Values [][]any           `json:"values"`
}
