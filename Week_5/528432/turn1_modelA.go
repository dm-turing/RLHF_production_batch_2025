package main

import (
	"fmt"
	"log"
	"math/rand"
	"os"
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

// String representation of severity levels
func (s Severity) String() string {
	return [...]string{"INFO", "DEBUG", "WARNING", "ERROR", "FATAL", "CRITICAL"}[s]
}

// LogError logs a message with a given severity
func LogError(severity Severity, message string) {
	now := time.Now().Format("2006-01-02 15:04:05")
	logMsg := fmt.Sprintf("%s [%s] %s", now, severity, message)

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

// Randomly generates an error with a severity level
func generateRandomError() {
	severity := Severity(rand.Intn(6))
	message := "This is a " + severity.String() + " level error."
	LogError(severity, message)
}

func main() {
	// Seed the random number generator
	rand.Seed(time.Now().UnixNano())

	// Set up a logger to log to the console for DEBUG and INFO
	log.SetFlags(log.LstdFlags)

	// Generate random errors every second
	for {
		generateRandomError()
		time.Sleep(1 * time.Second)
	}
}
