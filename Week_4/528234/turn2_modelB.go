package main

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

const (
	targetURL   = "http://example.com"
	numRequests = 100
	concurrency = 10
	timeout     = 2 * time.Second
)

func main() {
	// Context to manage resource cleanup
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	// Set up signal capturing to handle graceful shutdown
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	// WaitGroup to wait for all requests to finish
	var wg sync.WaitGroup
	requestCh := make(chan struct{}, concurrency)
	// Start monitoring go routine to gracefully terminate the test
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
			go func(reqID int) {
				defer wg.Done()
				defer func() { <-requestCh }()
				performRequest(ctx, reqID)
			}(i)
		}
	}
	// Wait for all requests to complete
	wg.Wait()
	duration := time.Since(start)
	fmt.Printf("Completed %d requests in %v\n", numRequests, duration)
}

// performRequest performs a single HTTP request and logs the result
func performRequest(ctx context.Context, reqID int) {
	// Introduce random delay to simulate failures
	delay := time.Duration(rand.Intn(3000)) * time.Millisecond
	time.Sleep(delay)
	// Create a new context with a timeout
	reqCtx, cancelReq := context.WithTimeout(ctx, timeout)
	defer cancelReq()
	req, err := http.NewRequestWithContext(reqCtx, "GET", targetURL, nil)
	if err != nil {
		log.Printf("Error creating request for reqID %d: %v", reqID, err)
		return
	}
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		// Check if the error is due to a timeout
		if ok := (err == context.DeadlineExceeded); ok {
			log.Printf("Request with reqID %d timed out after %v\n", reqID, timeout)
			return
		}
		log.Printf("Request error for reqID %d: %v", reqID, err)
		return
	}
	defer resp.Body.Close()
	log.Printf("Response status for reqID %d: %s", reqID, resp.Status)
}
