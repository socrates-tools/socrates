package runlog

import "time"

type Event struct {
	At      time.Time      `json:"at"`
	Kind    string         `json:"kind"`
	Subject string         `json:"subject"`
	Fields  map[string]any `json:"fields,omitempty"`
}

type Recorder interface {
	Record(Event) error
}
