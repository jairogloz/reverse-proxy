package main

import (
	"log"
	"net/http"

	"github.com/AndresKenji/reverse-proxy/internal/config"
	"github.com/AndresKenji/reverse-proxy/internal/server"
)

func main() {

	cfgFile, err := config.NewConfig("config.json")
	if err != nil {
		log.Fatal(err.Error())
	}

	srv := server.NewServer()
	srv.SetServerMux(cfgFile)

	log.Println("API GateWay running on port:", srv.Port)
	err = http.ListenAndServe(srv.Port, srv.Mux)
	if err != nil {
		log.Panic(err)
	}

}
