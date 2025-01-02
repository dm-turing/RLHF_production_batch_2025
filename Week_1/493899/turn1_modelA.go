package main

import (
	"encoding/json"
	"log"
	"os"
	"time"
)

// ABTestLog represents a log entry for an A/B test.
type ABTestLog struct {
	Timestamp  time.Time `json:"timestamp"`
	UserID     string    `json:"user_id"`
	Version    string    `json:"version"`
	Event      string    `json:"event"`
	DurationMs int       `json:"duration_ms"`
	Success    bool      `json:"success"`
}

// LogABTestEvent logs an event for an A/B test.
func LogABTestEvent(userID string, version string, event string, durationMs int, success bool) {
	logEntry := ABTestLog{
		Timestamp:  time.Now(),
		UserID:     userID,
		Version:    version,
		Event:      event,
		DurationMs: durationMs,
		Success:    success,
	}

	jsonData, err := json.Marshal(logEntry)
	if err != nil {
		log.Printf("Error marshaling log entry: %v", err)
		return
	}

	log.Printf("Logging AB test event: %s", string(jsonData))

	// In a production environment, you would write to a file or send to a logging service
	file, err := os.OpenFile("ab_test_logs.txt", os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		log.Printf("Error opening log file: %v", err)
		return
	}
	defer file.Close()

	_, err = file.Write(append([]byte("\n"), jsonData...))
	if err != nil {
		log.Printf("Error writing log to file: %v", err)
		return
	}
}

func main() {
	// Example usage
	LogABTestEvent("user123", "versionA", "page_load", 2500, true)
	LogABTestEvent("user456", "versionB", "page_load", 3000, false)
	LogABTestEvent("user789", "versionA", "button_click", 100, true)
	LogABTestEvent("user101", "versionB", "button_click", 150, false)
}
