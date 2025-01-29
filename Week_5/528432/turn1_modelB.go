package main

import (
	"fmt"
	"log"
	"math/rand"
	"os"
	"time"
)

const (
	logFilePath = "application.log"
)

var severityLevels = []string{"INFO", "DEBUG", "WARNING", "ERROR", "FATAL", "CRITICAL"}

func main() {
	// Set up a logger for file writing
	file, err := os.OpenFile(logFilePath, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatalf("Error opening log file: %v", err)
	}
	defer file.Close()
	fileLogger := log.New(file, "", log.LstdFlags)

	rand.Seed(time.Now().UnixNano())

	// Keep generating random errors and logging them
	for {
		generateAndLogError(fileLogger)
		time.Sleep(time.Second)
	}
}

// generateAndLogError generates a random error with a severity level and logs it to the appropriate location
func generateAndLogError(fileLogger *log.Logger) {
	severity := severityLevels[rand.Intn(len(severityLevels))]
	message := fmt.Sprintf("Random %s error occurred", severity)

	// Log to console for DEBUG and INFO severity levels
	if severity == "DEBUG" || severity == "INFO" {
		log.Println(message)
	} else {
		// Log to file for other severity levels
		fileLogger.Println(message)
	}
}
