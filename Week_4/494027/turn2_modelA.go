package main

import (
	"fmt"
	"time"

	"github.com/sirupsen/logrus"
)

// MouseEvent represents a generic mouse event
type MouseEvent struct {
	Description string
	TimeTaken   time.Duration
	Error       error
}

// EventHandler is a type for event handling functions
type EventHandler func(event MouseEvent, next func(MouseEvent))

// ProcessEvent processes the event through a chain of callbacks
func ProcessEvent(event MouseEvent, handlers ...EventHandler) {
	if len(handlers) == 0 {
		return
	}

	handlers[0](event, func(e MouseEvent) {
		ProcessEvent(e, handlers[1:]...)
	})
}

// LogEvent logs the event details
func LogEvent(event MouseEvent, next func(MouseEvent)) {
	start := time.Now()
	logrus.WithFields(logrus.Fields{
		"event": event.Description,
		"error": event.Error,
	}).Info("Processing event")
	elapsed := time.Since(start)
	event.TimeTaken = elapsed
	next(event)
}

// MouseClickHandler handles mouse click events
func MouseClickHandler(event MouseEvent, next func(MouseEvent)) {
	if event.Description == "MouseClick" {
		start := time.Now()
		logrus.Info("Handling Mouse Click Event")
		time.Sleep(100 * time.Millisecond) // Simulate processing time
		event.TimeTaken = time.Since(start)
	}
	next(event)
}

// MouseMoveHandler handles mouse move events
func MouseMoveHandler(event MouseEvent, next func(MouseEvent)) {
	if event.Description == "MouseMove" {
		start := time.Now()
		logrus.Info("Handling Mouse Move Event")
		time.Sleep(50 * time.Millisecond) // Simulate processing time
		event.TimeTaken = time.Since(start)
	}
	next(event)
}

// MouseDoubleClickHandler handles mouse double click events
func MouseDoubleClickHandler(event MouseEvent, next func(MouseEvent)) {
	if event.Description == "MouseDoubleClick" {
		start := time.Now()
		logrus.Info("Handling Mouse Double Click Event")
		// Simulate an error situation
		event.Error = fmt.Errorf("simulated error in MouseDoubleClick event")
		logrus.WithFields(logrus.Fields{
			"event": event.Description,
			"error": event.Error,
		}).Error("Error occurred during event handling")
		time.Sleep(150 * time.Millisecond) // Simulate processing time
		event.TimeTaken = time.Since(start)
	}
	next(event)
}

func main() {
	// Set logrus to show the time and order logs appropriately
	logrus.SetFormatter(&logrus.TextFormatter{
		FullTimestamp: true,
		TimestampFormat: time.RFC3339,
	})

	// Create some sample mouse events
	events := []MouseEvent{
		{Description: "MouseClick"},
		{Description: "MouseMove"},
		{Description: "MouseDoubleClick"},
	}

	// Register the handlers
	handlers := []EventHandler{
		LogEvent,
		MouseClickHandler,
		MouseMoveHandler,
		MouseDoubleClickHandler,
	}

	// Process each event
	for _, event := range events {
		fmt.Println("Processing new event:")
		ProcessEvent(event, handlers...)
		fmt.Println()
	}
}
