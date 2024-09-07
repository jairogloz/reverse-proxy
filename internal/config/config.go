package config

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/AndresKenji/reverse-proxy/internal/handlers"
)

type Config struct {
	Prefix           string            `json:"prefix"`
	HeaderIdentifier string            `json:"header_identifier"`
	BackendUrls      map[string]string `json:"backend_urls"`
	Secure           bool              `json:"secure"`
}

type ConfigFile struct {
	Endpoints []Config `json:"endpoints"`
}

// GenerateHandler Create a reverproxy HandleFunc for each Endpoint on the ConfigFile
func (c *Config) GenerateHandler() http.HandlerFunc {
	log.Println("Creating Handler for ", c.Prefix)

	return func(w http.ResponseWriter, r *http.Request) {
		hi := r.Header.Get(c.HeaderIdentifier)
		log.Println("Hit ", c.Prefix, hi)
		url, exist := c.BackendUrls[hi]
		if exist {
			handlers.ReverseProxy(url, c.Prefix, c.Secure).ServeHTTP(w, r)
		} else {
			http.Error(w, "Not Found", http.StatusNotFound)
		}
	}
}

// NewConfig read a json config file and return the ConfigFile object or error
func NewConfig(filePath string) (*ConfigFile, error) {
	var cfgFile ConfigFile
	file, err := os.Open(filePath)
	if err != nil {
		log.Println("Error al abrir el archivo de configuración:", err.Error())
		return &cfgFile, err
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	err = decoder.Decode(&cfgFile)
	if err != nil {
		log.Println("Error al deserializar el archivo JSON:", err.Error())
		return &cfgFile, err
	}

	return &cfgFile, nil
}

// DefaultConfig returns a default ConfigFile object
func DefaultConfig() *ConfigFile {
	return &ConfigFile{
		Endpoints: []Config{
			{
				Prefix:           "/default/",
				HeaderIdentifier: "country",
				BackendUrls: map[string]string{
					"col": "http://localhost:8001",
					"mex": "http://localhost:8002",
				},
				Secure: false,
			},
		},
	}
}

// LoadConfigFromElasticsearch retrieves a configuration file from an Elasticsearch endpoint
func LoadConfigFromElasticsearch() (*ConfigFile, error) {
	elasticURL := os.Getenv("elastic_url")
	indexName := "gw_config"
	var cfgFile ConfigFile

	// Step 1: Search for the latest document
	searchURL := fmt.Sprintf("%s/%s/_search?sort=@timestamp:desc&size=1", elasticURL, indexName)
	resp, err := http.Get(searchURL)
	if err != nil {
		log.Println("Error making request to Elasticsearch:", err.Error())
		return &cfgFile, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Printf("Error: received non-200 response code %d", resp.StatusCode)
		return &cfgFile, fmt.Errorf("received non-200 response code %d", resp.StatusCode)
	}

	// Read and decode the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Println("Error reading response from Elasticsearch:", err.Error())
		return &cfgFile, err
	}

	var searchResult map[string]interface{}
	if err := json.Unmarshal(body, &searchResult); err != nil {
		log.Println("Error unmarshalling search result JSON:", err.Error())
		return &cfgFile, err
	}

	// Extract the latest document
	hits, ok := searchResult["hits"].(map[string]interface{})
	if !ok {
		log.Println("Invalid response format: missing 'hits' field")
		return &cfgFile, fmt.Errorf("invalid response format: missing 'hits' field")
	}

	hitArray, ok := hits["hits"].([]interface{})
	if !ok || len(hitArray) == 0 {
		log.Println("No documents found in the search results")
		return &cfgFile, fmt.Errorf("no documents found")
	}

	latestDoc := hitArray[0].(map[string]interface{})
	source, ok := latestDoc["_source"].(map[string]interface{})
	if !ok {
		log.Println("Invalid document format: missing '_source' field")
		return &cfgFile, fmt.Errorf("invalid document format: missing '_source' field")
	}

	// Convert the source document to JSON and unmarshal into ConfigFile
	sourceData, err := json.Marshal(source)
	if err != nil {
		log.Println("Error marshalling source document to JSON:", err.Error())
		return &cfgFile, err
	}

	if err := json.Unmarshal(sourceData, &cfgFile); err != nil {
		log.Println("Error unmarshalling source document JSON:", err.Error())
		return &cfgFile, err
	}

	return &cfgFile, nil
}
