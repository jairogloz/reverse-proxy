package models

import "time"

// LogEntry representa la estructura del log que se enviará a Elasticsearch.
type LogEntry struct {
	Timestamp  time.Time `json:"@timestamp"`
	RemoteAddr string    `json:"remote_addr"`
	Method     string    `json:"method"`
	Path       string    `json:"path"`
	Proto      string    `json:"proto"`
	TargetURL  string    `json:"target_url"`
	Duration   string    `json:"duration"`
}
