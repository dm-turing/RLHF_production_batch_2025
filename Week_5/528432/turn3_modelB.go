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

	"github.com/sirupsen/logrus"
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
	log      *logrus.Logger // logrus Global logger variable
	logStore []LogEntry
	mu       sync.Mutex
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

func init() {
	log = logrus.New()
	// Log as JSON instead of the default ASCII formatter.
	log.Formatter = &logrus.JSONFormatter{}

	// Output to stderr instead of stdout, could also be a file.
	log.Out = os.Stderr

	// Only log the warning severity or above.
	log.Level = logrus.WarnLevel
}

// String representation of severity levels
func (s Severity) String() string {
	return [...]string{"INFO", "DEBUG", "WARNING", "ERROR", "FATAL", "CRITICAL"}[s]
}

func LogError(severity Severity, message string, path string, userID int) {
	now := time.Now()
	logEntry := logrus.Fields{
		"timestamp": now,
		"severity":  severity,
		"message":   message,
		"path":      path,
		"user_id":   userID,
	}

	switch severity {
	case DEBUG:
		log.WithFields(logEntry).Debug("Log message")
	case INFO:
		log.WithFields(logEntry).Info("Log message")
	case WARNING:
		log.WithFields(logEntry).Warn("Log message")
	case ERROR:
		log.WithFields(logEntry).Error("Log message")
	case FATAL:
		log.WithFields(logEntry).Fatal("Log message")
	case CRITICAL:
		log.WithFields(logEntry).Panic("Log message")
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
	// log.SetFlags(log.LstdFlags)

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
