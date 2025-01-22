package main

import (
	"context"
	"net/http"
	"net/http/httptest"
	"os"
	"sync"
	"testing"
	"time"
)

func TestMain(m *testing.M) {
	// Set up a test server to handle requests
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	// Override the targetURL with the test server URL
	// targetURL = ts.URL

	os.Exit(m.Run())
}

func TestLoadTest_Success(t *testing.T) {
	// We can specify a smaller number of requests for the test to run faster
	numRequests := 5
	// concurrency := 2

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	var wg sync.WaitGroup
	var successfulRequests int

	start := time.Now()
	for i := 0; i < numRequests; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			if performRequest(ctx, i) {
				successfulRequests++
			}
		}()
	}

	wg.Wait()

	duration := time.Since(start)

	// Check if all requests were successful
	if successfulRequests != numRequests {
		t.Fatalf("Expected %d successful requests, got %d", numRequests, successfulRequests)
	}

	// Verify that the test completed within a reasonable time
	if duration > 2*time.Second {
		t.Fatalf("Test took too long to complete: %v", duration)
	}
}
