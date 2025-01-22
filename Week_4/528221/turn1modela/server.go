package main

import (
	"fmt"
	"log"
	"net/http"
)

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello from server %s", r.Host)
	log.Printf("Handled request from %s", r.Host)
}

func main() {
	http.HandleFunc("/", handler)
	port := ":8081" // Change to 8082, 8083 for other servers
	log.Printf("Starting server on %s", port)
	log.Fatal(http.ListenAndServe(port, nil))
}
