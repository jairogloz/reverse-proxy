package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/resources", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello from Service 1!\n")
	})

	fmt.Println("Service 1 running on port 8081")
	log.Fatal(http.ListenAndServe(":8081", nil))
}
