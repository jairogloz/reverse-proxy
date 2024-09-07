package server

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/AndresKenji/reverse-proxy/internal/config"
	"github.com/AndresKenji/reverse-proxy/internal/database"
	"github.com/AndresKenji/reverse-proxy/internal/middleware"
)

type Server struct {
	Port    string
	Mux     *http.ServeMux
	Context context.Context
	Database *database.Database
}

func NewServer(ctx context.Context) *Server {
	server := &Server{}
	server.Context = ctx
	port := os.Getenv("port")
	if port == "" {
		log.Println("There is no environment variable for server port, server will use port 8080 ")
		server.Port = "8080"
	} else {
		server.Port = fmt.Sprintf(":%s", port)
	}
	server.Database = database.NewDatabase()

	return server
}

func (s *Server) SetServerMux(cfgFile *config.ConfigFile) {
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Server listening"))
	})
	for _, cfg := range cfgFile.Endpoints {
		log.Println(cfg)
		mux.HandleFunc(cfg.Prefix, middleware.RequestLoggerMiddleware(cfg.GenerateHandler()))
	}

	s.Mux = mux

}

func (s *Server) StartServer() error {
	srv := &http.Server{
		Addr:    ":" + s.Port,
		Handler: s.Mux,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Error starting server: %v", err)
		}
	}()

	log.Printf("Server started on port %s\n", s.Port)

	<-s.Context.Done()

	log.Println("Shutting down server...")

	if err := srv.Shutdown(context.Background()); err != nil {
		return err
	}

	log.Println("Server stopped gracefully.")
	return nil
}
