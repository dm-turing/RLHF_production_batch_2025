package main

import (
	"log"
	"net/http"
	"strconv"
)

type ABTestManager interface {
	ServeVersion(w http.ResponseWriter, r *http.Request)
	RecordInteraction(r *http.Request)
}

func main() {
	testManager := NewSimpleABTestManager()
	http.HandleFunc("/", testManager.ServeVersion)
	log.Fatal(http.ListenAndServe(":8080", nil))
}

type SimpleABTestManager struct {
}

func (m *SimpleABTestManager) ServeVersion(w http.ResponseWriter, r *http.Request) {
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
