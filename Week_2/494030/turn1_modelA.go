package main

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"
)

// Data is a simple data structure to represent the data flowing through the pipeline.
type Data struct {
	Name    string
	Age     int
	Message string
}

// Middleware is an interface that defines the behavior of each middleware component.
type Middleware interface {
	Handle(in *Data) (*Data, error)
}

// LoggerMiddleware logs the data at each step.
type LoggerMiddleware struct {
	next Middleware
}

func NewLoggerMiddleware(next Middleware) *LoggerMiddleware {
	return &LoggerMiddleware{next}
}

func (l *LoggerMiddleware) Handle(in *Data) (*Data, error) {
	log.Printf("Logging: %+v\n", in)
	return l.next.Handle(in)
}

// ValidatorMiddleware validates the data, ensuring that the age is a positive integer.
type ValidatorMiddleware struct {
	next Middleware
}

func NewValidatorMiddleware(next Middleware) *ValidatorMiddleware {
	return &ValidatorMiddleware{next}
}

func (v *ValidatorMiddleware) Handle(in *Data) (*Data, error) {
	if in.Age <= 0 {
		return nil, fmt.Errorf("Invalid age: %d", in.Age)
	}
	return v.next.Handle(in)
}

// UppercaseMiddleware converts the message to uppercase.
type UppercaseMiddleware struct {
	next Middleware
}

func NewUppercaseMiddleware(next Middleware) *UppercaseMiddleware {
	return &UppercaseMiddleware{next}
}

func (u *UppercaseMiddleware) Handle(in *Data) (*Data, error) {
	in.Message = strings.ToUpper(in.Message)
	return u.next.Handle(in)
}

// EndMiddleware is the final middleware that prints the processed data.
type EndMiddleware struct {
	next Middleware
}

func NewEndMiddleware(next Middleware) *EndMiddleware {
	return &EndMiddleware{next}
}

func (e *EndMiddleware) Handle(in *Data) (*Data, error) {
	out, err := json.MarshalIndent(in, "", "  ")
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
		last = NewEndMiddleware(last)
	}
	return middleware[0]
}

func main() {
	// Define the data to be processed.
	data := &Data{Name: "Alice", Age: 25, Message: "Hello World"}

	// Create a chain of middleware.
	pipeline := Chain(
		NewLoggerMiddleware(
			NewValidatorMiddleware(
				NewUppercaseMiddleware(
					NewEndMiddleware(nil),
				),
			),
		),
	)

	// Process the data through the pipeline.
	if _, err := pipeline.Handle(data); err != nil {
		log.Fatalf("Error processing data: %v", err)
	}
}
