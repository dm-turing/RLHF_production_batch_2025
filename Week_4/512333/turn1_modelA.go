package main

import "fmt"

// ApplicationState encapsulates all the global state of the application
type ApplicationState struct {
	counter int
	text    string
}

// StateHandler defines a function type for stateful callbacks
type StateHandler func(*ApplicationState) error

// NewApplicationState initializes the application state
func NewApplicationState() *ApplicationState {
	return &ApplicationState{counter: 0, text: ""}
}

// IncrementCounter creates a callback for incrementing a counter
func IncrementCounter(delta int) StateHandler {
	return func(state *ApplicationState) error {
		state.counter += delta
		fmt.Printf("Counter incremented: %d\n", state.counter)
		return nil
	}
}

// UpdateText creates a callback for updating the text
func UpdateText(newText string) StateHandler {
	return func(state *ApplicationState) error {
		state.text = newText
		fmt.Printf("Text updated: %s\n", state.text)
		return nil
	}
}

// RunStateHandler executes a callback to modify the state
func RunStateHandler(handler StateHandler, state *ApplicationState) {
	if err := handler(state); err != nil {
		fmt.Println("Error executing handler:", err)
	}
}

// Application logic simulation
func main() {
	appState := NewApplicationState()

	// Define handlers
	incrementByFive := IncrementCounter(5)
	updateText := UpdateText("Hello, Go!")

	// Using handlers to change application state
	RunStateHandler(incrementByFive, appState)
	RunStateHandler(updateText, appState)

	// Additional modular functionality
	newIncrement := IncrementCounter(3)
	RunStateHandler(newIncrement, appState)
}
