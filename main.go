package main

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
)

// reverseProxy sets up a reverse proxy for a given target URL
func reverseProxy(target string) http.HandlerFunc {
	url, _ := url.Parse(target)
	proxy := httputil.NewSingleHostReverseProxy(url)

	return func(w http.ResponseWriter, r *http.Request) {
		r.URL.Host = url.Host
		r.URL.Scheme = url.Scheme
		r.Header.Set("X-Forwarded-Host", r.Host)
		r.Host = url.Host
		proxy.ServeHTTP(w, r)
	}
}

// routeHandler handles routing based on URL path prefix
func routeHandler(w http.ResponseWriter, r *http.Request) {
	country := r.Header.Get("Country")
	country = strings.ToLower(country)
	switch country {
	case "mx":
		reverseProxy("http://localhost:8081").ServeHTTP(w, r)
	case "ar":
		reverseProxy("http://localhost:8082").ServeHTTP(w, r)
	default:
		// You could also route to a default service here
		http.Error(w, "Not Found", http.StatusNotFound)
	}
}

func main() {
	http.HandleFunc("/", routeHandler)

	fmt.Println("API Gateway running on port 8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
