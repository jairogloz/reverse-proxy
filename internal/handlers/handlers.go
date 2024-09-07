package handlers

import (
	"crypto/tls"
	"net/http"
	"net/http/httputil"
	"net/url"
	"time"

	"github.com/AndresKenji/reverse-proxy/internal/util"
)

// ReverseProxy function creates an HTTP handler that acts as a reverse proxy to forward incoming requests to a specified target server.
func ReverseProxy(target string, prefix string, secure bool) http.HandlerFunc {
	url, _ := url.Parse(target)
	proxy := httputil.NewSingleHostReverseProxy(url)
	if !secure {
		proxy.Transport = &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true}, // Desactivar la verificación del certificado TLS
		}
	}

	return func(w http.ResponseWriter, r *http.Request) {
		// Preserve the rest of the URL path after the prefix
		r.URL.Path = r.URL.Path[len(prefix):]

		r.URL.Host = url.Host
		r.URL.Scheme = url.Scheme
		r.Header.Set("X-Forwarded-Host", r.Host)
		r.Host = url.Host

		// Registrar la solicitud proxy
		start := time.Now()
		logEntry := util.LogEntry{
			Timestamp:  time.Now().UTC(),
			RemoteAddr: r.RemoteAddr,
			Method:     r.Method,
			Path:       r.URL.Path,
			Proto:      r.Proto,
			TargetURL:  target,
			Duration:   time.Since(start).String(),
		}
		// Enviar el log a Elasticsearch
		go util.SendLogToElasticsearch(logEntry, "api_gateway")

		proxy.ServeHTTP(w, r)
	}
}
