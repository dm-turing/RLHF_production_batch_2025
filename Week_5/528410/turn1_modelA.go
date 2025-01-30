package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

type Error struct {
	Context  string
	UserID   int
	ErrorCode int
	Message  string
	Next     *Error // Pointer for chaining errors
}

var contexts = []string{"database", "network", "filesystem", "authentication"}
var messages = []string{
	"connection timeout", "missing file", "permission denied",
	"invalid credentials", "unknown error",
}

func generateRandomError() *Error {
	return &Error{
		Context:  contexts[rand.Intn(len(contexts))],
		UserID:   rand.Intn(1000),
		ErrorCode: rand.Intn(500), // Assumes error codes range from 0-499
		Message:  messages[rand.Intn(len(messages))],
	}
}

func printError(e *Error) {
	fmt.Printf("Error Context: %s, UserID: %d, ErrorCode: %d, Message: %s\n", e.Context, e.UserID, e.ErrorCode, e.Message)
	if e.Next != nil {
		fmt.Println("Chained Errors:")
		printError(e.Next)
	}
}

func main() {
	rand.Seed(time.Now().UnixNano())
	errorMap := make(map[string]*Error)
	mu := sync.Mutex{}

	for range time.Tick(time.Second) {
		go func() {
			mu.Lock()
			defer mu.Unlock()

			newError := generateRandomError()
			currentContext := newError.Context

			if existingError, exists := errorMap[currentContext]; exists {
				// Append this error to the last error generated with the same context
				lastError := existingError
				for lastError.Next != nil {
					lastError = lastError.Next
				}
				lastError.Next = newError
			} else {
				errorMap[currentContext] = newError
			}

			printError(newError)
			fmt.Println("-------")
		}()
	}
}