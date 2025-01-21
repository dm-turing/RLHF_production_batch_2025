package main

import (
	"fmt"
)

type MouseError struct {
	EventType string
	Error     error
}
type MouseErrorCallback func(MouseError)

type ErrorGroupingCategory struct {
	Name    string
	Filters []func(MouseError) bool
}

type MouseErrorHandler struct {
	ErrorCallbacks []MouseErrorCallback
	ErrorGroups    []ErrorGroupingCategory
}

func (h *MouseErrorHandler) AddErrorCallback(callback MouseErrorCallback) {
	h.ErrorCallbacks = append(h.ErrorCallbacks, callback)
}

func (h *MouseErrorHandler) AddErrorGroup(name string, filters ...func(MouseError) bool) {
	h.ErrorGroups = append(h.ErrorGroups, ErrorGroupingCategory{Name: name, Filters: filters})
}

func (h *MouseErrorHandler) HandleError(err MouseError) {
	// Call all error callbacks regardless of group
	for _, callback := range h.ErrorCallbacks {
		callback(err)
	}

	// Call error callbacks based on groups
	for _, group := range h.ErrorGroups {
		for _, filter := range group.Filters {
			if filter(err) {
				fmt.Println("Error group:", group.Name)
				break
			}
		}
	}
}

// Define example error filtering functions
func isClientError(err MouseError) bool {
	return err.EventType == "client_error"
}

func isSystemError(err MouseError) bool {
	return err.EventType == "system_error"
}
func main() {
	errorHandler := &MouseErrorHandler{}

	// Register generic error callback
	errorHandler.AddErrorCallback(func(err MouseError) {
		fmt.Printf("Generic Error: %v\n", err.Error)
	})

	// Register error callbacks for specific groups
	errorHandler.AddErrorGroup("Client Errors", isClientError)
	errorHandler.AddErrorGroup("System Errors", isSystemError)

	// Simulate handling errors
	errorHandler.HandleError(MouseError{EventType: "client_error", Error: fmt.Errorf("Invalid input")})
	errorHandler.HandleError(MouseError{EventType: "system_error", Error: fmt.Errorf("Out of memory")})
	errorHandler.HandleError(MouseError{EventType: "unexpected_error", Error: fmt.Errorf("Something went wrong")})
}
