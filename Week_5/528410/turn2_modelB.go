package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Error struct {
	Context   string
	UserID    int
	ErrorCode int
	Message   string
	CreatedAt time.Time
	Next      *Error
}

var contexts = []string{"database", "network", "filesystem", "authentication"}
var messages = []string{
	"connection timeout", "missing file", "permission denied",
	"invalid credentials", "unknown error",
}

var (
	// Define severity levels
	severities = []zapcore.Level{zapcore.DebugLevel, zapcore.InfoLevel, zapcore.ErrorLevel, zapcore.WarnLevel, zapcore.FatalLevel}
	// Create a logger with zap
	logger, _ = zap.NewProduction()
)

func generateRandomError() *Error {
	return &Error{
		Context:   contexts[rand.Intn(len(contexts))],
		UserID:    rand.Intn(1000),
		ErrorCode: rand.Intn(500),
		Message:   messages[rand.Intn(len(messages))],
		CreatedAt: time.Now(),
	}
}

// WrapError function adds context to the error
func WrapError(err error, context string) error {
	return fmt.Errorf("%s: %w", context, err)
}

func logError(e *Error, severity zapcore.Level) {
	logger.Sugar().Log(
		severity,
		// zap.String("severity", severity.String()),
		zap.String("context", e.Context),
		zap.Int("user_id", e.UserID),
		zap.Int("error_code", e.ErrorCode),
		zap.String("message", e.Message),
	)

	if e.Next != nil {
		logError(e.Next, severity)
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

			// Chain errors with the same context in the last 5 seconds
			chainTimeout := 5 * time.Second
			if existingError, exists := errorMap[currentContext]; exists {
				if time.Since(existingError.CreatedAt) < chainTimeout {
					lastError := existingError
					for lastError.Next != nil {
						lastError = lastError.Next
					}
					lastError.Next = newError
				} else {
					errorMap[currentContext] = newError
				}
			} else {
				errorMap[currentContext] = newError
			}
			// Log the error with random severity
			severity := severities[rand.Intn(len(severities))]
			logError(newError, severity)
			fmt.Println("-------")
		}()
	}
}
