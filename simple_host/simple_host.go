package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
)

func main() {
	var nFlag = flag.Int("port", 8080, "Listening port")
	flag.Parse() 

	router := http.NewServeMux()
	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello from port %v", *nFlag) 
	})

	server := http.Server{
		Addr:    fmt.Sprintf("localhost:%d", *nFlag), 
		Handler: router,
	}

	log.Println("Listening on port:", *nFlag) 
	if err := server.ListenAndServe(); err != nil {
		log.Panic(err.Error())
	}
}
