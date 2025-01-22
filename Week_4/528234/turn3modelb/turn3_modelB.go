package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

const (
	targetURL   = "http://github.com"
	numRequests = 20
	concurrency = 5
)

func main() {
	// Context to manage resource cleanup
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Setting up signal capturing to handle graceful shutdown
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	// WaitGroup to wait for all requests to finish
	var wg sync.WaitGroup
	requestCh := make(chan struct{}, concurrency)

	// Start monitoring goroutine
	go func() {
		<-sigs
		fmt.Println("\nReceived interrupt, terminating load test...")
		cancel()
	}()

	start := time.Now()
	for i := 0; i < numRequests; i++ {
		select {
		case <-ctx.Done():
			break
		case requestCh <- struct{}{}:
			wg.Add(1)
			go func(requestID int) {
				defer wg.Done()
				defer func() { <-requestCh }()
				performRequest(ctx, requestID)
			}(i)
		}
	}

	// Wait for all requests to complete
	wg.Wait()
	duration := time.Since(start)
	fmt.Printf("Completed %d requests in %v\n", numRequests, duration)
}

// This is a helper function to simulate a successful request
func performRequest(ctx context.Context, requestID int) bool {
	req, err := http.NewRequestWithContext(ctx, "GET", targetURL, nil)
	if err != nil {
		return false
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return false
	}
	defer resp.Body.Close()

	// Simulate a short delay
	time.Sleep(100 * time.Millisecond)

	return resp.StatusCode == http.StatusOK
}
