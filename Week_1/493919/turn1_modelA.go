package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
)

// A custom context value for holding user credentials.
type userCredentials struct {
	username string
	token    string
}

// SetCredentials adds user credentials to a context.
func setCredentials(ctx context.Context, credentials userCredentials) context.Context {
	return context.WithValue(ctx, "userCredentials", credentials)
}

// GetCredentials retrieves user credentials from a context.
func getCredentials(ctx context.Context) (*userCredentials, error) {
	value := ctx.Value("userCredentials")
	if value == nil {
		return nil, fmt.Errorf("credentials not found in context")
	}
	creds, ok := value.(userCredentials)
	if !ok {
		return nil, fmt.Errorf("unexpected value type %T for credentials in context", value)
	}
	return &creds, nil
}

// A handler to process user requests.
func processRequest(w http.ResponseWriter, r *http.Request) {
	// Extract credentials from the context
	creds, err := getCredentials(r.Context())
	if err != nil {
		http.Error(w, "No credentials provided", http.StatusUnauthorized)
		return
	}

	log.Printf("Processing request for user %s with token %s\n", creds.username, creds.token)
	// Process the request using the provided credentials
	// For example, make a call to an external service with the API key

	// Serve a simple response
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Hello, %s!\n", creds.username)
}

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// Authentication logic: validate the request and obtain credentials
		credentials := userCredentials{username: "alice", token: "secretToken"}

		// Create a new context with credentials
		ctx := setCredentials(r.Context(), credentials)

		// Pass the credentialed context to the request handler
		r = r.WithContext(ctx)

		processRequest(w, r)
	})

	log.Println("Server starting on :8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
