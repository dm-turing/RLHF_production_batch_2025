package main

import (
	"log"
	"time"

	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"github.com/sirupsen/logrus"
)

func main() {
	// Set up Logger
	log := logrus.New()
	log.SetOutput(getLogWriter())
	log.SetLevel(logrus.InfoLevel)

	// Write logs
	for i := 0; i < 100000; i++ {
		log.Info("Hello, World!")
	}
}
func getLogWriter() *rotatelogs.RotateLogs {
	// Generate the log file path
	baseLogPath := "application.log"
	linkName := "application.log"

	// Set up rotation options
	maxAge := 24 * time.Hour  // Keep logs for 24 hours
	rotationTime := time.Hour // Rotate logs every hour

	// Create a rotater
	writer, err := rotatelogs.New(
		baseLogPath+".%Y%m%d%H",                   // This is the actual log file format
		rotatelogs.WithLinkName(linkName),         // Create a link named application.log to the lastest log file
		rotatelogs.WithMaxAge(maxAge),             // Rotate logs based on the max age
		rotatelogs.WithRotationTime(rotationTime), // Rotate logs based on the rotation time
	)
	if err != nil {
		log.Fatal(err)
	}

	return writer
}
