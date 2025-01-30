package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"

	"github.com/sirupsen/logrus"
)

type Severity string

const (
	DEBUG    Severity = "DEBUG"
	INFO     Severity = "INFO"
	WARNING  Severity = "WARNING"
	ERROR    Severity = "ERROR"
	CRITICAL Severity = "CRITICAL"
	FATAL    Severity = "FATAL"
)

type Error struct {
	Context   string
	UserID    int
	ErrorCode int
	Message   string
	Time      time.Time
	Severity  Severity
	Next      *Error
}

func (e *Error) Error() string {
	return fmt.Sprintf("Context: %s, UserID: %d, ErrorCode: %d, Message: %s, Severity: %s, Time: %s",
		e.Context, e.UserID, e.ErrorCode, e.Message, e.Severity, e.Time.Format(time.RFC3339))
}

var contexts = []string{"database", "network", "filesystem", "authentication"}
var messages = []string{
	"connection timeout", "missing file", "permission denied",
	"invalid credentials", "unknown error",
}
var severities = []Severity{DEBUG, INFO, WARNING, ERROR, CRITICAL, FATAL}

func generateRandomError() *Error {
	return &Error{
		Context:   contexts[rand.Intn(len(contexts))],
		UserID:    rand.Intn(1000),
		ErrorCode: rand.Intn(500),
		Message:   messages[rand.Intn(len(messages))],
		Time:      time.Now(),
		Severity:  severities[rand.Intn(len(severities))],
	}
}

func main() {
	rand.Seed(time.Now().UnixNano())
	errorMap := make(map[string][]*Error)
	mu := sync.Mutex{}

	log := logrus.New()
	log.SetFormatter(&logrus.TextFormatter{FullTimestamp: true})

	for range time.Tick(time.Second) {
		go func() {
			mu.Lock()
			defer mu.Unlock()

			newError := generateRandomError()
			currentContext := newError.Context

			if errors, exists := errorMap[currentContext]; exists {
				// Filter errors from the last 5 seconds
				filteredErrors := []*Error{}
				for _, err := range errors {
					if time.Since(err.Time).Seconds() <= 5 {
						filteredErrors = append(filteredErrors, err)
					}
				}
				errorMap[currentContext] = filteredErrors

				// Append this error to the last valid error in the same context within the last 5 seconds
				if len(filteredErrors) > 0 {
					lastError := filteredErrors[len(filteredErrors)-1]
					lastError.Next = newError
				}
			}

			// Add the new error to the map
			errorMap[currentContext] = append(errorMap[currentContext], newError)
			log.WithFields(logrus.Fields{
				"Context":   newError.Context,
				"UserID":    newError.UserID,
				"ErrorCode": newError.ErrorCode,
				"Message":   newError.Message,
				"Severity":  newError.Severity,
			}).Error(newError.Error()) // Log as error for demonstration
		}()
	}
}
