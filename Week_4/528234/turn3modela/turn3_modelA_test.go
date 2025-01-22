package main

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestPerformRequest(t *testing.T) {
	// Set up a test HTTP server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(1 * time.Second) // simulate processing time
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	}))
	defer server.Close()

	// Test context and client
	ctx := context.Background()
	client := &http.Client{}

	tests := []struct {
		name      string
		delay     time.Duration
		expectErr bool
	}{
		{
			name:      "Successful request",
			delay:     0,
			expectErr: false,
		},
		{
			name:      "Request timeout",
			delay:     3 * time.Second, // simulate delay greater than the timeout
			expectErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.delay > 0 {
				// Use a context with a shorter timeout to induce a timeout error
				timeoutCtx, cancel := context.WithTimeout(ctx, 0*time.Second)
				defer cancel()

				err := performRequest(timeoutCtx, client, server.URL, 1)
				if (err != nil) != tt.expectErr {
					t.Errorf("Expected error: %v, got: %v", tt.expectErr, err)
				}
			} else {
				err := performRequest(ctx, client, server.URL, 1)
				if (err != nil) != tt.expectErr {
					t.Errorf("Expected error: %v, got: %v", tt.expectErr, err)
				}
			}
		})
	}
}
