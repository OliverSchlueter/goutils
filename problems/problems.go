package problems

import (
	"encoding/json"
)

func UnmarshalJSON(data []byte) *Problem {
	if data == nil || len(data) == 0 {
		return nil
	}

	var problem Problem
	if err := json.Unmarshal(data, &problem); err != nil {
		return nil
	}

	return &problem
}
