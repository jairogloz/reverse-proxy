package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
)

func main() {

	file, err := os.Open("config.json")
	if err != nil {
		log.Fatal("Error al abrir el archivo de configuración:", err)
	}
	defer file.Close()

	var cfgFile ConfigFile
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&cfgFile)
	if err != nil {
		log.Fatal("Error al deserializar el archivo JSON:", err)
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("hello world!"))
	})

	for _, cfg := range cfgFile.Endpoints {
		log.Println(cfg)
		mux.HandleFunc(cfg.Prefix, cfg.generateHandler().ServeHTTP)
	}

	server := http.Server{
		Addr: fmt.Sprintf(":%s",cfgFile.Port),
		Handler: mux,
	}

	fmt.Println("API GateWay running on port:", cfgFile.Port)

	log.Fatal(server.ListenAndServe())

}

func reverseProxy(target string) http.HandlerFunc {
	url, _ := url.Parse(target)
	proxy := httputil.NewSingleHostReverseProxy(url)

	return func(w http.ResponseWriter, r *http.Request) {
		r.URL.Host = url.Host
		r.URL.Scheme = url.Scheme
		r.Header.Set("X-Forwarded-Host", r.Host)
		r.Host = url.Host
		proxy.ServeHTTP(w,r)
	}
}


func (c *Config) generateHandler() http.HandlerFunc {
	log.Println("Creating Handler for ",c.Prefix)

	return func(w http.ResponseWriter, r *http.Request) {
		hi := r.Header.Get(c.HeaderIdentifier)
		log.Println("Hit ",c.Prefix, hi)
		url, exist := c.BackendUrls[hi]
		if exist {
			reverseProxy(url).ServeHTTP(w,r)
		}else {
			http.Error(w, "Not Found", http.StatusNotFound)
		}
	}
}

type Config struct {
	Prefix string `json:"prefix"`
	HeaderIdentifier string `json:"header_identifier"`
	BackendUrls  map[string]string  `json:"backend_urls"`
}

type ConfigFile struct {
	Port string `json:"port"`
	Endpoints []Config `json:"endpoints"`
}