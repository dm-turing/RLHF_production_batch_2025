package main

import (
	"log"
	"net/http"

	"github.com/DataDog/datadog-go/statsd"
	"github.com/gorilla/mux"
)

// Our application's logic for handling resource management requests
func handleResourceManagement(w http.ResponseWriter, r *http.Request) {
	// Simulate some business logic
	log.Println("Processing resource management request")
	// Add some sleep to simulate workload
	//time.Sleep(time.Second * 2)
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Resource management response"))
}

// Middleware function to trace and instrument the request using Datadog
func datadogMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Create a new timer to track the request duration
		timer := client.NewTimer("http.request.duration")
		defer timer.Stop()

		// Increment the request count
		client.Incr("http.request.count", []string{"route:" + mux.CurrentRoute(r).GetName()}, 1)

		// Execute the next handler in the chain
		next.ServeHTTP(w, r)
	})
}

var client *statsd.Client

func init() {
	// Create a Datadog StatsD client
	var err error
	client, err = statsd.New("127.0.0.1:8125", statsd.DefaultNamespace)
	if err != nil {
		log.Fatalf("Error creating Datadog client: %v", err)
	}
}

func main() {
	r := mux.NewRouter()

	// Define the routes and attach the datadogMiddleware as a middleware
	r.HandleFunc("/", handleResourceManagement).Methods("GET").Name("root")
	r.Use(datadogMiddleware)

	log.Println("Resource management tool started on port 8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}
