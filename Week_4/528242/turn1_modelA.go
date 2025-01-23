package main

import (
	"fmt"
	"math/rand"
	"time"
)

type Severity int

const (
	DEBUG Severity = iota
	INFO
	WARNING
	ERROR
	CRITICAL
	FATAL
)

func (s Severity) String() string {
	return []string{"DEBUG", "INFO", "WARNING", "ERROR", "CRITICAL", "FATAL"}[s]
}

type ErrorMessage struct {
	Severity Severity
	Message  string
	Metadata string
}

func generateRandomErrorMessage() ErrorMessage {
	severities := []Severity{DEBUG, INFO, WARNING, ERROR, CRITICAL, FATAL}
	messages := []string{
		"Disk space low",
		"Network latency detected",
		"Application crashed",
		"User input error",
		"Configuration file missing",
	}

	randomSeverity := severities[rand.Intn(len(severities))]
	randomMessage := messages[rand.Intn(len(messages))]

	return ErrorMessage{
		Severity: randomSeverity,
		Message:  randomMessage,
	}
}

func addMetadata(err ErrorMessage) ErrorMessage {
	metadataOptions := []string{
		"on server A",
		"at 12:00 PM",
		"by user X",
		"during backup process",
		"on cluster node 2",
	}
	randomMetadata := metadataOptions[rand.Intn(len(metadataOptions))]
	err.Metadata = randomMetadata
	return err
}

func main() {
	rand.Seed(time.Now().UnixNano())

	errorChannel := make(chan ErrorMessage, 1)
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()
	exit := make(chan bool)

	go func() {
		for {
			errorMessage := generateRandomErrorMessage()
			errorMessage = addMetadata(errorMessage)
			errorChannel <- errorMessage
			time.Sleep(1 * time.Second)
		}
	}()

	go func() {
		infoErrors := []ErrorMessage{}
		for {
			select {
			case err := <-errorChannel:
				switch err.Severity {
				case DEBUG, INFO:
					infoErrors = append(infoErrors, err)
				case WARNING, ERROR:
					fmt.Printf("WARNING/ERROR: [%s] %s - %s\n", err.Severity, err.Message, err.Metadata)
				case CRITICAL, FATAL:
					fmt.Printf("CRITICAL/FATAL: [%s] %s - %s\n", err.Severity, err.Message, err.Metadata)
					exit <- true
					return
				}

			case <-ticker.C:
				for _, err := range infoErrors {
					fmt.Printf("DEBUG/INFO: [%s] %s - %s\n", err.Severity, err.Message, err.Metadata)
				}
				infoErrors = nil
			}
		}
	}()

	<-exit
	fmt.Println("Session ended due to a CRITICAL or FATAL error.")
}
