package handlers

import (
	"crypto/tls"
	"io"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"time"

	"github.com/AndresKenji/reverse-proxy/internal/models"
	//"github.com/AndresKenji/reverse-proxy/internal/util"
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
		logEntry := models.LogEntry{
			Timestamp:  time.Now().UTC(),
			RemoteAddr: r.RemoteAddr,
			Method:     r.Method,
			Path:       r.URL.Path,
			Proto:      r.Proto,
			TargetURL:  target,
			Duration:   time.Since(start).String(),
		}
		log.Println(logEntry)
		//go util.SendLogToElasticsearch(logEntry,"api_gateway")

		proxy.ServeHTTP(w, r)
	}
}

func RedirectRequest(target string, prefix string, secure bool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		remainingPath := r.URL.Path[len(prefix):]
		targetURL, err := url.Parse(target)
		if err != nil {
			http.Error(w, "Invalid target URL", http.StatusInternalServerError)
			return
		}
		// Reconstruir la URL con el resto del path y query params
		targetURL.Path += remainingPath
		targetURL.RawQuery = r.URL.RawQuery

		// Crear una nueva solicitud hacia la URL destino
		req, err := http.NewRequest(r.Method, targetURL.String(), r.Body)
		if err != nil {
			http.Error(w, "Error creating request", http.StatusInternalServerError)
			return
		}

		// Copiar headers del request original
		for key, values := range r.Header {
			for _, value := range values {
				req.Header.Add(key, value)
			}
		}

		// Si la conexión no es segura, desactivar la verificación TLS
		client := &http.Client{}
		if !secure {
			client.Transport = &http.Transport{
				TLSClientConfig: &tls.Config{InsecureSkipVerify: true}, // Desactivar la verificación del certificado TLS
			}
		}

		// Hacer la solicitud
		start := time.Now()
		resp, err := client.Do(req)
		if err != nil {
			http.Error(w, "Error making request", http.StatusInternalServerError)
			return
		}
		defer resp.Body.Close()

		// Registrar la solicitud
		logEntry := models.LogEntry{
			Timestamp:  time.Now().UTC(),
			RemoteAddr: req.RemoteAddr,
			Method:     req.Method,
			Path:       req.URL.Path,
			Proto:      req.Proto,
			TargetURL:  target,
			Duration:   time.Since(start).String(),
		}
		log.Println(logEntry)
		//go util.SendLogToElasticsearch(logEntry,"api_gateway")

		// Copiar el código de estado y headers de la respuesta
		w.WriteHeader(resp.StatusCode)
		for key, values := range resp.Header {
			for _, value := range values {
				w.Header().Add(key, value)
			}
		}

		// Escribir el cuerpo de la respuesta
		io.Copy(w, resp.Body)
	}
}
