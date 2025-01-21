package main

import (
	"fmt"
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
	fmt.Printf("Event: %s\n", event.Description)
	if event.Error != nil {
		fmt.Printf("Error: %v\n", event.Error)
	}
	elapsed := time.Since(start)
	event.TimeTaken = elapsed
	next(event)
}

// MouseClickHandler handles mouse click events
func MouseClickHandler(event MouseEvent, next func(MouseEvent)) {
	if event.Description == "MouseClick" {
		start := time.Now()
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

// MouseDoubleClickHandler handles mouse double click events
func MouseDoubleClickHandler(event MouseEvent, next func(MouseEvent)) {
	if event.Description == "MouseDoubleClick" {
		start := time.Now()
		fmt.Println("Handling Mouse Double Click Event")
		time.Sleep(150 * time.Millisecond) // Simulate processing time
		event.TimeTaken = time.Since(start)
	}
	next(event)
}

func main() {
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
