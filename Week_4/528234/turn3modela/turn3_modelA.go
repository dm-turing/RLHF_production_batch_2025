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
				performRequest(ctx, &http.Client{}, targetURL, requestID)
			}(i)
		}
	}

	// Wait for all requests to complete
	wg.Wait()
	duration := time.Since(start)
	fmt.Printf("Completed %d requests in %v\n", numRequests, duration)
}

func performRequest(ctx context.Context, client *http.Client, url string, requestID int) error {
	// Set a timeout for the request context
	reqCtx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()

	req, err := http.NewRequestWithContext(reqCtx, "GET", url, nil)
	if err != nil {
		return fmt.Errorf("error creating request: %v", err)
	}

	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("request %d error: %v", requestID, err)
	}
	defer resp.Body.Close()

	fmt.Printf("Request %d response status: %s\n", requestID, resp.Status)
	return nil
}
