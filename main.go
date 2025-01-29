package main

import (
    "errors"
    "fmt"
    "log"
    "math/rand"
    "os"
    "time"
)

type Handler func(userID string) error

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
    return func(userID string) error {
        if rand.Intn(10) < 3 { // 30% chance to simulate an error
            return generateRandomError("Order Placement")
        }
        return next(userID)
    }
}

// Logging middleware
func loggingMiddleware(next Handler) Handler {
    return func(userID string) error {
        start := time.Now()
        err := next(userID)
        if err != nil {
            log.Printf("ERROR: user_id=%s timestamp=%s error=%s duration=%s\n", userID, start.Format(time.RFC3339), err.Error(), time.Since(start))
        }
        return err
    }
}

// Product Search stage
func productSearch(userID string) error {
    fmt.Println("Executing product search...")
    // Simulate Success
    return nil
}

// Order Placement stage
func orderPlacement(userID string) error {
    fmt.Println("Placing the order...")
    // Simulate Success
    return nil
}

// Shipping stage
func shipping(userID string) error {
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

    // Configure the logger to write to a file with a specific format
    logFile, err := os.OpenFile("shopping_cart.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
    if err != nil {
        fmt.Println("Failed to open log file:", err)
        return
    }
    defer logFile.Close()
    log.SetOutput(logFile)
    log.SetFlags(0) // Disable default timestamp and file location logging

    pipeline := []Handler{productSearch, orderPlacement, shipping}

    processOrder := func(userID string) {
        for _, handler := range pipeline {
            wrappedHandler := chainMiddleware(handler, errorSimulationMiddleware, loggingMiddleware)
            if err := wrappedHandler(userID); err != nil {
                fmt.Println("Error encountered:", err)
            }
        }
    }

    ticker := time.NewTicker(3 * time.Second)

    fmt.Println("Starting the shopping cart application...")

    userIDs := []string{"user1", "user2", "user3"} // Simulating multiple users

    for {
        select {
        case <-ticker.C:
            fmt.Println("\nInitiating new order process...")
            // Select a random user ID for each order process
            userID := userIDs[rand.Intn(len(userIDs))]
            processOrder(userID)
        }
    }
}