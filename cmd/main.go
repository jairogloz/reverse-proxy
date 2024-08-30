package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/AndresKenji/reverse-proxy/internal/config"
	"github.com/AndresKenji/reverse-proxy/internal/middleware"
)

func main() {

	cfgFile, err := config.NewConfig("config.json")
	if err != nil{
		log.Fatal(err.Error())
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("hello world!"))
	})

	for _, cfg := range cfgFile.Endpoints {
		log.Println(cfg)
		mux.HandleFunc(cfg.Prefix, middleware.RequestLoggerMiddleware(cfg.GenerateHandler().ServeHTTP) )
	}

	server := http.Server{
		Addr:    fmt.Sprintf(":%s", cfgFile.Port),
		Handler: mux,
	}

	log.Println("API GateWay running on port:", cfgFile.Port)

	log.Fatal(server.ListenAndServe())

}