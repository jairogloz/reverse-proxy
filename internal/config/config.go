package config

import (
	"encoding/json"
	"log"
	"net/http"
	"os"

	"github.com/AndresKenji/reverse-proxy/internal/handlers"
)

type Config struct {
	Prefix           string            `json:"prefix"`
	HeaderIdentifier string            `json:"header_identifier"`
	BackendUrls      map[string]string `json:"backend_urls"`
}

type ConfigFile struct {
	Port      string   `json:"port"`
	Endpoints []Config `json:"endpoints"`
}

func (c *Config) GenerateHandler() http.HandlerFunc {
	log.Println("Creating Handler for ", c.Prefix)

	return func(w http.ResponseWriter, r *http.Request) {
		hi := r.Header.Get(c.HeaderIdentifier)
		log.Println("Hit ", c.Prefix, hi)
		url, exist := c.BackendUrls[hi]
		if exist {
			handlers.ReverseProxy(url, c.Prefix).ServeHTTP(w, r)
		} else {
			http.Error(w, "Not Found", http.StatusNotFound)
		}
	}
}

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
