package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
)

type Person struct {
	Name string `json:"name"`
}

func main() {
	var nFlag = flag.Int("port", 8080, "Listening port")
	flag.Parse() 

	router := http.NewServeMux()
	router.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello from port %v", *nFlag) 
	})

	router.HandleFunc("GET /PING", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("PONG"))
	})

	router.HandleFunc("GET /users/{id}", func(w http.ResponseWriter, r *http.Request) {
		UserId := r.PathValue("id")
		w.Write([]byte("User id is :"+ UserId))
	})

	router.HandleFunc("GET /printname", func(w http.ResponseWriter, r *http.Request) {
		params := r.URL.Query()
		name := params.Get("name")
		w.Write([]byte("params ="+name))
	})

	router.HandleFunc("POST /person", func(w http.ResponseWriter, r *http.Request) {
		var p Person
		err := json.NewDecoder(r.Body).Decode(&p)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
		}

		fmt.Fprintf(w, "Persona: %s",p.Name)
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
