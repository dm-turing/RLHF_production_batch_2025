package main

import (
	"fmt"
	"log"
	"os"
)

// define a Logger type to encapsulate the logging functionality
type Logger struct {
	file *os.File
}

// NewLogger creates a new Logger and opens the log file
func NewLogger(filename string) (*Logger, error) {
	file, err := os.OpenFile(filename, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return nil, err
	}
	return &Logger{file: file}, nil
}

// Close the log file when the Logger is no longer needed
func (l *Logger) Close() {
	l.file.Close()
}

// Logf writes a formatted log message to the file
func (l *Logger) Logf(format string, args ...interface{}) {
	log.SetOutput(l.file)
	log.Printf(format, args...)
}

func main() {
	// create a new logger instance
	logger, err := NewLogger("file_history.log")
	if err != nil {
		fmt.Println("Error creating logger:", err)
		return
	}
	defer logger.Close()

	// Example usage:
	userName := "exampleUser"
	filePath := "/path/to/the/file"

	// Log user navigation
	logger.Logf("User %s navigated to directory: %s", userName, filePath)

	// Log file operations
	logger.Logf("User %s opened file: %s", userName, filePath)
	logger.Logf("User %s saved file: %s", userName, filePath)
	logger.Logf("User %s deleted file: %s", userName, filePath)

}
