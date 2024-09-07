package config

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/AndresKenji/reverse-proxy/internal/handlers"
)

type Config struct {
	Prefix           string            `json:"prefix" bson:"prefix"`
	HeaderIdentifier string            `json:"header_identifier" bson:"header_identifier"`
	BackendUrls      map[string]string `json:"backend_urls" bson:"backend_urls"`
	Secure           bool              `json:"secure" bson:"secure"`
}

type ConfigFile struct {
	CreatedAt time.Time `json:"created_at" bson:"created_at"`
	Endpoints []Config `json:"endpoints" bson:"endpoints"`
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

