package handlers

import (
	"net/http"
	"net/http/httputil"
	"net/url"
)

//  ReverseProxy function creates an HTTP handler that acts as a reverse proxy to forward incoming requests to a specified target server.
func ReverseProxy(target string) http.HandlerFunc {
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