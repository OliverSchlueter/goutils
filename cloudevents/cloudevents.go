// Package cloudevents is based on the CloudEvents specification: https://github.com/cloudevents/spec
package cloudevents

import "time"

const SpecVersion = "1.0.2"

type CloudEvent struct {
	SpecVersion     string    `json:"specversion"`
	ID              string    `json:"id"`
	Type            string    `json:"type"`
	Subject         string    `json:"subject"`
	Source          string    `json:"source"`
	Time            time.Time `json:"time"`
	DataContentType string    `json:"datacontenttype"`
}
