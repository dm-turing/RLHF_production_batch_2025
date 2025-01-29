package main

import (
	"errors"
	"fmt"
	"log"
	"math/rand"
	"os"
	"time"
)

type Handler func() error
type Middleware func(Handler) Handler

// Logger struct to hold logger and user information
type Logger struct {
	Logger *log.Logger
	User   string
}

// Function to generate random errors
func generateRandomError(stage string) error {
	errorsList := []string{
		stage + " network failure",
		stage + " invalid response data",
		stage + " server timeout",
		stage + " unauthorized access",
	}
	return errors.New(errorsList[rand.Intn(len(errorsList))])
}

// Error simulator middleware
func errorSimulationMiddleware(next Handler) Handler {
	return func() error {
		if rand.Intn(10) < 3 { // 30% chance to simulate an error
			return generateRandomError("Order Placement")
		}
		return next()
	}
}

// Logging middleware
func loggingMiddleware(l *Logger) Middleware {
	return func(next Handler) Handler {
		return func() error {
			err := next()
			if err != nil {
				// Capture error trace
				trace := fmt.Sprintf("Error Trace: %+v", err)
				// Log the error with user, timestamp, and error trace
				l.Logger.Printf("User: %s, Timestamp: %s, Error: %s", l.User, time.Now().Format(time.RFC3339), trace)
			}
			return err
		}
	}
}

// Product Search stage
func productSearch() error {
	fmt.Println("Executing product search...")
	// Simulate Success
	return nil
}

// Order Placement stage
func orderPlacement() error {
	fmt.Println("Placing the order...")
	// Simulate Success
	return nil
}

// Shipping stage
func shipping() error {
	fmt.Println("Processing shipping...")
	// Simulate Success
	return nil
}

// Execute a chain of middlewares with an endpoint
func chainMiddleware(h Handler, middlewares ...Middleware) Handler {
	for _, middleware := range middlewares {
		h = middleware(h)
	}
	return h
}

func main() {
	rand.Seed(time.Now().UnixNano()) // Seed random number generator
	logger := log.New(os.Stdout, "Shopping Cart: ", log.LstdFlags)
	user := "user123" // Replace this with the actual user ID

	pipeline := []Handler{productSearch, orderPlacement, shipping}

	processOrder := func() {
		l := &Logger{Logger: logger, User: user}
		for _, handler := range pipeline {
			wrappedHandler := chainMiddleware(handler, errorSimulationMiddleware, loggingMiddleware(l))
			if err := wrappedHandler(); err != nil {
				fmt.Println("Error encountered:", err)
			}
		}
	}

	ticker := time.NewTicker(3 * time.Second)

	fmt.Println("Starting the shopping cart application...")

	for {
		select {
		case <-ticker.C:
			fmt.Println("\nInitiating new order process for user:", user)
			processOrder()
		}
	}
}
