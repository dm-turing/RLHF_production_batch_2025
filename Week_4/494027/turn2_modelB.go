package main

import (
	"fmt"
	"log"
	"math/rand"
	"sync"
	"time"
)

// MouseEvent represents a generic mouse event
type MouseEvent struct {
	Description string
	TimeTaken   time.Duration
	Error       error
}

// EventHandler is a type for event handling functions
type EventHandler func(event MouseEvent, next func(MouseEvent))

// logger is a custom logger that groups logs by mouse event
type logger struct {
	mu  sync.Mutex
	log *log.Logger
}

func newLogger() *logger {
	return &logger{log: log.New(log.Writer(), "", log.LstdFlags)}
}

func (l *logger) LogEvent(event MouseEvent) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.log.Printf("[%s] Event: %s", event.Description, event.Description)
	if event.Error != nil {
		l.log.Printf("[%s] Error: %v", event.Description, event.Error)
	}
}

// ProcessEvent processes the event through a chain of callbacks, using the provided logger
func ProcessEvent(event MouseEvent, handlers ...EventHandler) {
	if len(handlers) == 0 {
		return
	}
	logger := newLogger()
	handlers[0](event, func(e MouseEvent) {
		logger.LogEvent(e)
		ProcessEvent(e, handlers[1:]...)
	})
}

// LogEvent logs the event details using the custom logger
func LogEvent(event MouseEvent, next func(MouseEvent)) {
	start := time.Now()
	// Simulate some processing
	time.Sleep(time.Duration(100+rand.Intn(100)) * time.Millisecond)
	elapsed := time.Since(start)
	event.TimeTaken = elapsed
	next(event)
}

// MouseClickHandler handles mouse click events
func MouseClickHandler(event MouseEvent, next func(MouseEvent)) {
	if event.Description == "MouseClick" {
		start := time.Now()
		// Simulate an error in this handler
		if rand.Intn(2) == 0 {
			event.Error = fmt.Errorf("mouse click error occurred")
		}
		fmt.Println("Handling Mouse Click Event")
		time.Sleep(100 * time.Millisecond) // Simulate processing time
		event.TimeTaken = time.Since(start)
	}
	next(event)
}

// MouseMoveHandler handles mouse move events
func MouseMoveHandler(event MouseEvent, next func(MouseEvent)) {
	if event.Description == "MouseMove" {
		start := time.Now()
		fmt.Println("Handling Mouse Move Event")
		time.Sleep(50 * time.Millisecond) // Simulate processing time
		event.TimeTaken = time.Since(start)
	}
	next(event)
}

func main() {
	// Create some sample mouse events
	events := []MouseEvent{
		{Description: "MouseClick"},
		{Description: "MouseMove"},
		{Description: "MouseClick"},
		{Description: "MouseMove"},
		{Description: "MouseDoubleClick"},
	}

	// Register the handlers
	handlers := []EventHandler{
		LogEvent,
		MouseClickHandler,
		MouseMoveHandler,
	}

	// Process each event
	for _, event := range events {
		fmt.Println("Processing new event:")
		ProcessEvent(event, handlers...)
		fmt.Println()
	}
}
