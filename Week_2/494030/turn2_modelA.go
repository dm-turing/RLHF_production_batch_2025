package main

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"
)

// Event is a simple event structure.
type Event struct {
	Type    string `json:"type"`
	Payload string `json:"payload"`
}

// Middleware is an interface that defines the behavior of each middleware component.
type Middleware interface {
	Handle(event *Event) (*Event, error)
}

// LoggerMiddleware logs the event.
type LoggerMiddleware struct {
	next Middleware
}

func NewLoggerMiddleware(next Middleware) *LoggerMiddleware {
	return &LoggerMiddleware{next}
}

func (l *LoggerMiddleware) Handle(event *Event) (*Event, error) {
	log.Printf("Logging event: %+v\n", event)
	return l.next.Handle(event)
}

// ValidatorMiddleware validates the event payload.
type ValidatorMiddleware struct {
	next Middleware
}

func NewValidatorMiddleware(next Middleware) *ValidatorMiddleware {
	return &ValidatorMiddleware{next}
}

func (v *ValidatorMiddleware) Handle(event *Event) (*Event, error) {
	if event.Payload == "" {
		return nil, fmt.Errorf("Invalid payload: %s", event.Payload)
	}
	return v.next.Handle(event)
}

// ProcessorMiddleware processes the event payload.
type ProcessorMiddleware struct {
	next Middleware
}

func NewProcessorMiddleware(next Middleware) *ProcessorMiddleware {
	return &ProcessorMiddleware{next}
}

func (p *ProcessorMiddleware) Handle(event *Event) (*Event, error) {
	event.Payload = strings.ToUpper(event.Payload)
	return p.next.Handle(event)
}

// PersisterMiddleware persists the processed event.
type PersisterMiddleware struct {
	next Middleware
}

func NewPersisterMiddleware(next Middleware) *PersisterMiddleware {
	return &PersisterMiddleware{next}
}

func (p *PersisterMiddleware) Handle(event *Event) (*Event, error) {
	out, err := json.MarshalIndent(event, "", "  ")
	if err != nil {
		return nil, err
	}
	fmt.Println(string(out))
	return nil, nil
}

// Chain combines multiple middleware components into a single pipeline.
func Chain(middleware ...Middleware) Middleware {
	if len(middleware) == 0 {
		return nil
	}
	var last Middleware = middleware[len(middleware)-1]
	for i := len(middleware) - 2; i >= 0; i-- {
		last = NewLoggerMiddleware(last)
	}
	return middleware[0]
}

func main() {
	// Define the event to be processed.
	event := &Event{Type: "someEvent", Payload: "hello world"}

	// Create a chain of middleware.
	pipeline := Chain(
		NewValidatorMiddleware(
			NewProcessorMiddleware(
				NewPersisterMiddleware(nil),
			),
		),
	)

	// Process the event through the pipeline.
	if _, err := pipeline.Handle(event); err != nil {
		log.Fatalf("Error processing event: %v", err)
	}
}
