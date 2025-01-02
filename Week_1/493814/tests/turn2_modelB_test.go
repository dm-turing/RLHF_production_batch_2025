package main

import (
	"net/http"
	"testing"
	"time"
)

func TestClientCommunicationOverhead(t *testing.T) {
	// Start the configuration server
	go main()

	// Wait for the server to start
	time.Sleep(1 * time.Second)

	start := time.Now()
	for i := 0; i < 10; i++ {
		response, err := http.Get("http://localhost:8080/config")
		if err != nil {
			t.Fatal(err)
		}
		response.Body.Close()
	}
	elapsed := time.Since(start)
	// Check if the elapsed time is significantly more than expected (e.g., > 100ms) for each request
	if elapsed.Milliseconds() > 100*10 {
		t.Errorf("Excessive communication overhead: %v", elapsed)
	}
}
