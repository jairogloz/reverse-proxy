package util

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"os"

	"github.com/AndresKenji/reverse-proxy/internal/models"
)

// SendLogToElasticsearch envía un log a Elasticsearch.
func SendLogToElasticsearch(logEntry models.LogEntry, index string) {

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
