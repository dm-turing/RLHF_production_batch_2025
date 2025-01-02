package main

import (
	"context"
	"net/http"
	"testing"
	"time"

	"github.com/gorilla/concurrent"
)

func TestConfigServerScalability(t *testing.T) {
	numClients := 100
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	wg := concurrent.NewWaitGroup()
	for i := 0; i < numClients; i++ {
		wg.Add()
		go func() {
			defer wg.Done()
			client := &http.Client{}
			resp, err := client.GetWithContext(ctx, "http://localhost:8080/config")
			if err != nil {
				t.Errorf("Request failed: %v", err)
				return
			}
			defer resp.Body.Close()

			if resp.StatusCode != http.StatusOK {
				t.Errorf("Expected status %d, got %d", http.StatusOK, resp.StatusCode)
				return
			}
		}()
	}

	wg.Wait()
}
