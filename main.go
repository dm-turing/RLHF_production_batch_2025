package main

import (
    "errors"
    "fmt"
    "math/rand"
    "time"
)

type Handler func() error

// Middleware is a function that wraps a handler
type Middleware func(Handler) Handler

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

    pipeline := []Handler{productSearch, orderPlacement, shipping}

    processOrder := func() {
        for _, handler := range pipeline {
            wrappedHandler := chainMiddleware(handler, errorSimulationMiddleware)
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
            fmt.Println("\nInitiating new order process...")
            processOrder()
        }
    }
}