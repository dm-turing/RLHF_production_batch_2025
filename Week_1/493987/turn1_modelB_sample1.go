package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
)

func main() {
	http.HandleFunc("/", handler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func handler(w http.ResponseWriter, r *http.Request) {
	// Get the query parameters
	params := r.URL.Query()

	// Get the value of the 'version' parameter and convert it to an integer
	versionStr := params.Get("version")
	version, err := strconv.Atoi(versionStr)
	if err != nil {
		// Default to version A if the parameter is missing or invalid
		version = 1
	}

	// Serve the appropriate version of the webpage
	if version == 1 {
		serveVersionA(w, r)
	} else if version == 2 {
		serveVersionB(w, r)
	} else {
		// Invalid version, default to version A
		serveVersionA(w, r)
	}
}

func serveVersionA(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Welcome to Version A!")
}

func serveVersionB(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Welcome to Version B!")
}
