package main

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	log "github.com/sirupsen/logrus"
)

// Define the Severity type with available levels of severity
type Severity int

const (
	INFO Severity = iota
	DEBUG
	WARNING
	ERROR
	FATAL
	CRITICAL
)

// LogEntry represents a structured log entry
type LogEntry struct {
	Timestamp time.Time `json:"timestamp"`
	Severity  Severity  `json:"severity"`
	Message   string    `json:"message"`
	Path      string    `json:"path"`
	UserID    int       `json:"user_id"`
}

// Global in-memory log store and mutex for synchronization
var (
	logStore []LogEntry
	mu       sync.Mutex
)

// String representation of severity levels
func (s Severity) String() string {
	return [...]string{"INFO", "DEBUG", "WARNING", "ERROR", "FATAL", "CRITICAL"}[s]
}

// LogError creates and stores a log entry
func LogError(severity Severity, message string, path string, userID int) {
	now := time.Now()
	logEntry := LogEntry{
		Timestamp: now,
		Severity:  severity,
		Message:   message,
		Path:      path,
		UserID:    userID,
	}

	// Append to the in-memory log store
	mu.Lock()
	logStore = append(logStore, logEntry)
	mu.Unlock()

	// Structured logging with logrus
	logger := log.WithFields(log.Fields{
		"time":     logEntry.Timestamp,
		"severity": logEntry.Severity.String(),
		"message":  logEntry.Message,
		"path":     logEntry.Path,
		"userId":   logEntry.UserID,
	})

	// Direct the output based on severity
	switch severity {
	case DEBUG, INFO:
		logger.Info(message)
	default:
		logger.Warn(message) // Using Warn for simplification, you can map specific severities to appropriate log levels.
		fileWriter, err := os.OpenFile("errors.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
		if err == nil {
			defer fileWriter.Close()
			log.SetOutput(fileWriter)
			logger.Warn(message)
			log.SetOutput(os.Stdout) // Reset output to stdout
		} else {
			logger.Error("Could not open log file")
		}
	}
}

// Randomly generates an error with a severity level and context
func generateRandomError() {
	severity := Severity(rand.Intn(6))
	message := "This is a " + severity.String() + " level error."
	path := fmt.Sprintf("/api/resource/%d", rand.Intn(100))
	userID := rand.Intn(10)
	LogError(severity, message, path, userID)
}

// ServeHTTP handles querying logs via HTTP
func serveHTTP(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()

	// Example of simple query handling
	severityFilter := query.Get("severity")
	pathFilter := query.Get("path")
	userIDFilter := query.Get("userid")

	mu.Lock()
	defer mu.Unlock()

	filteredLogs := []LogEntry{}

	for _, logEntry := range logStore {
		if (severityFilter == "" || logEntry.Severity.String() == severityFilter) &&
			(pathFilter == "" || strings.Contains(logEntry.Path, pathFilter)) &&
			(userIDFilter == "" || strconv.Itoa(logEntry.UserID) == userIDFilter) {
			filteredLogs = append(filteredLogs, logEntry)
		}
	}

	response, _ := json.Marshal(filteredLogs)
	w.Header().Set("Content-Type", "application/json")
	w.Write(response)
}

func main() {
	// Seed the random number generator
	rand.Seed(time.Now().UnixNano())

	// Configure logrus for structured logging with JSON formatter
	log.SetFormatter(&log.JSONFormatter{})
	log.SetOutput(os.Stdout)

	// Start log generation in a separate goroutine
	go func() {
		for {
			generateRandomError()
			time.Sleep(1 * time.Second)
		}
	}()

	// HTTP handler for querying logs
	http.HandleFunc("/search", serveHTTP)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
