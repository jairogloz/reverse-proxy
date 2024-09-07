package server

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/AndresKenji/reverse-proxy/internal/config"
	"github.com/AndresKenji/reverse-proxy/internal/middleware"
)

type Server struct {
	Port        string
	Mux         *http.ServeMux
	Middlewares []middleware.Middleware // TODO: implementar arreglo de middlewares y aplicarlos a los endpoint
}

func NewServer() *Server {
	server := &Server{}
	port := os.Getenv("port")
	if port == "" {
		log.Println("There is no environment variable for server port, server will use port 8080 ")
		server.Port = ":8080"
	} else {
		server.Port = fmt.Sprintf(":%s", port)
	}

	return server
}

func (s *Server) SetServerMux(cfgFile *config.ConfigFile) {
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Server listening"))
	})
	for _, cfg := range cfgFile.Endpoints {
		log.Println(cfg)
		mux.HandleFunc(cfg.Prefix, cfg.GenerateHandler())
	}

	s.Mux = mux

}

func MiddlewareChain(middlewares ...middleware.Middleware) middleware.Middleware {
	return func(next http.Handler) http.HandlerFunc {
		for i := len(middlewares) - 1; i >= 0; i-- {
			next = middlewares[i](next)
		}
		return next.ServeHTTP
	}
}
