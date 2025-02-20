package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
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
	Timestamp time.Time
	Severity  Severity
	Message   string
	Path      string
	UserID    int
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

	// Format log message
	logMsg := fmt.Sprintf("%s [%s] %s - Path: %s, UserID: %d", now.Format("2006-01-02 15:04:05"), severity, message, path, userID)

	// Direct the output based on severity
	if severity == DEBUG || severity == INFO {
		log.Println(logMsg)
	} else {
		// Log to a file for WARNING, ERROR, FATAL, and CRITICAL
		f, err := os.OpenFile("errors.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
		if err != nil {
			log.Fatalf("Error opening log file: %v", err)
		}
		defer f.Close()
		fileLogger := log.New(f, "", log.LstdFlags)
		fileLogger.Println(logMsg)
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

	// Set up a logger to log to the console for DEBUG and INFO
	log.SetFlags(log.LstdFlags)

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
