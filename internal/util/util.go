package util

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"time"
)

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

// SendLogToElasticsearch envía un log a Elasticsearch.
func SendLogToElasticsearch(logEntry LogEntry, index string) {

	esURL := os.Getenv("elastic_url")

	body, err := json.Marshal(logEntry)
	if err != nil {
		log.Printf("Error marshaling log entry: %v", err)
		return
	}

	req, err := http.NewRequest("POST", esURL+"/"+index+"/_doc", bytes.NewBuffer(body))
	if err != nil {
		log.Printf("Error creating request: %v", err)
		return
	}

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("Error sending log to Elasticsearch: %v", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated && resp.StatusCode != http.StatusOK {
		log.Printf("Unexpected response from Elasticsearch: %s", resp.Status)
	}
}
