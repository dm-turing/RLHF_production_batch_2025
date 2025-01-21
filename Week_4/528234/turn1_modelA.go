package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

const (
	targetURL   = "http://github.com"
	numRequests = 100
	concurrency = 10
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

	// Start monitoring go routine
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
			go func() {
				defer wg.Done()
				defer func() { <-requestCh }()
				performRequest(ctx)
			}()
		}
	}

	// Wait for all requests to complete
	wg.Wait()
	duration := time.Since(start)
	fmt.Printf("Completed %d requests in %v\n", numRequests, duration)
}

// performRequest performs a single HTTP request and logs the result
func performRequest(ctx context.Context) {
	req, err := http.NewRequestWithContext(ctx, "GET", targetURL, nil)
	if err != nil {
		log.Printf("Error creating request: %v", err)
		return
	}

	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	resp, err := client.Do(req)
	if err != nil {
		log.Printf("Request error: %v", err)
		return
	}
	defer resp.Body.Close()

	log.Printf("Response status: %s", resp.Status)
}
