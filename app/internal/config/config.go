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
	AllowedMethods   []string          `json:"allowed_methods" bson:"allowed_methods"`
	AllowedIps       []string          `json:"allowed_ips" bson:"allowed_ips"`
}

type ConfigFile struct {
	CreatedAt time.Time `json:"created_at" bson:"created_at"`
	Endpoints []Config  `json:"endpoints" bson:"endpoints"`
}

// GenerateHandler Create a reverproxy HandleFunc for each Endpoint on the ConfigFile
func (c *Config) GenerateProxyHandler() http.HandlerFunc {
	log.Println("Creating Handler for ", c.Prefix)

	return func(w http.ResponseWriter, r *http.Request) {
		 // Verificar si la IP está permitida
		 clientIP := r.RemoteAddr
		 if !isIPAllowed(clientIP, c.AllowedIps) {
			 http.Error(w, "Access Denied", http.StatusForbidden)
			 return
		 }
 
		 // Verificar si el método está permitido
		 if !isMethodAllowed(r.Method, c.AllowedMethods) {
			 http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
			 return
		 }
		hi := r.Header.Get(c.HeaderIdentifier)
		url, exist := c.BackendUrls[hi]
		if exist {
			handlers.CreateReverseProxyHandler(url, c.Prefix, c.Secure).ServeHTTP(w, r)
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
		log.Println("Error at opening the configuración file:", err.Error())
		return &cfgFile, err
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	err = decoder.Decode(&cfgFile)
	if err != nil {
		log.Println("Error Decoding the json file:", err.Error())
		return &cfgFile, err
	}

	return &cfgFile, nil
}

// DefaultConfig returns a default ConfigFile object
func DefaultConfig() *ConfigFile {
	cfgFile, err := NewConfig("/app/config.json")
	if err != nil {
		log.Panic("Could not read the defaiult config file")
	}
	return cfgFile
}

// isIPAllowed Verifica si la IP está permitida
func isIPAllowed(clientIP string, allowedIPs []string) bool {
    if allowedIPs[0] == "*" {
        return true
    }
    
    for _, allowed := range allowedIPs {
        if allowed == clientIP {
            return true
        }
    }
    return false
}

// isMethodAllowed Verifica si el método está permitido
func isMethodAllowed(method string, allowedMethods []string) bool {
    for _, allowed := range allowedMethods {
        if allowed == "*" || allowed == method {
            return true
        }
    }
    return false
}