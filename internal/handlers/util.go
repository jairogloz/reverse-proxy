package handlers

import (
	"net/http"
	"net/http/httputil"
	"net/url"
)


func ReverseProxy(target string, prefix string) http.HandlerFunc {
	url, _ := url.Parse(target)
	proxy := httputil.NewSingleHostReverseProxy(url)

	return func(w http.ResponseWriter, r *http.Request) {
		// Preserve the rest of the URL path after the prefix
		r.URL.Path = r.URL.Path[len(prefix):]

		r.URL.Host = url.Host
		r.URL.Scheme = url.Scheme
		r.Header.Set("X-Forwarded-Host", r.Host)
		r.Host = url.Host
		proxy.ServeHTTP(w, r)
	}
}