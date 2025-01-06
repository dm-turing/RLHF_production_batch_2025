package main

import (
	"fmt"
	"net/http"
	"os"
)

// CustomError is a struct representing a network error with a description, type, and underlying error.
type CustomError struct {
	Message string  `json:"message"`
	Type    string  `json:"type"`
	Errors  []error `json:"errors"`
	Cause   error   `json:"cause"`
}

// Error implementation for CustomError to allow it to be returned as an error.
func (ce *CustomError) Error() string {
	return fmt.Sprintf("CustomError { message: %s, type: %s, cause: %s }", ce.Message, ce.Type, ce.Cause.Error())
}

// NewCustomError returns a new CustomError instance.
func NewCustomError(message string, err error) *CustomError {
	return &CustomError{
		Message: message,
		Type:    "network_error",
		Errors:  []error{},
		Cause:   err,
	}
}

// NewAggregateCustomError returns a CustomError instance composed of multiple errors.
func NewAggregateCustomError(message string, errs []error) *CustomError {
	return &CustomError{
		Message: message,
		Type:    "aggregate_network_error",
		Errors:  errs,
		Cause:   nil,
	}
}

func main() {
	url := "https://example.com/api"

	// Function to simulate a network operation that can return an error.
	makeNetworkRequest := func(url string) (*http.Response, error) {
		if url == "https://example.com/api" {
			// Simulating a specific network error
			return nil, fmt.Errorf("timeout making request to %s", url)
		} else {
			// Simulating successful request
			return &http.Response{Status: "200 OK"}, nil
		}
	}

	// Making the first network request
	resp, err := makeNetworkRequest(url)
	if err != nil {
		// Handling a single network error
		networkError := NewCustomError("First network request failed", err)
		fmt.Fprintf(os.Stderr, "%v\n", networkError)
		os.Exit(1)
	}
	defer resp.Body.Close()

	// Making more network requests, simulating multiple errors
	urls := []string{"https://example.com/api", "https://invalid.url", "https://example.com/api"}
	errors := []error{}

	for _, u := range urls {
		resp, err := makeNetworkRequest(u)
		if err != nil {
			errors = append(errors, err)
		}
		defer resp.Body.Close()
	}

	if len(errors) > 0 {
		// Handling multiple network errors in an aggregate error
		aggregateError := NewAggregateCustomError("Multiple network requests failed", errors)
		fmt.Fprintf(os.Stderr, "%v\n", aggregateError)
		os.Exit(1)
	}

	fmt.Println("All network requests completed successfully!")
}
