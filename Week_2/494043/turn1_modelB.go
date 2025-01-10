package main

import "fmt"

// CallbackRegistry is a registry to manage callbacks for specific events.
type CallbackRegistry struct {
	callbacks map[string][]func()
}

// NewCallbackRegistry creates a new CallbackRegistry.
func NewCallbackRegistry() *CallbackRegistry {
	return &CallbackRegistry{
		callbacks: make(map[string][]func()),
	}
}

// Subscribe registers a callback for a specific event.
func (r *CallbackRegistry) Subscribe(event string, callback func()) {
	r.callbacks[event] = append(r.callbacks[event], callback)
}

// Unsubscribe unregisters a callback for a specific event.
func (r *CallbackRegistry) Unsubscribe(event string, callback func()) {
	if callbacks, ok := r.callbacks[event]; ok {
		for i, c := range callbacks {
			if fmt.Sprintf("%v", c) == fmt.Sprintf("%v", callback) {
				r.callbacks[event] = append(callbacks[:i], callbacks[i+1:]...)
				return
			}
		}
	}
}

// Publish triggers all callbacks registered for a specific event.
func (r *CallbackRegistry) Publish(event string) {
	if callbacks, ok := r.callbacks[event]; ok {
		for _, c := range callbacks {
			c()
		}
	}
}

func main() {
	registry := NewCallbackRegistry()

	// Subscribe to the "user_input" event
	registry.Subscribe("user_input", func() {
		fmt.Println("User input received!")
	})

	// Subscribe to the "button_click" event
	registry.Subscribe("button_click", func() {
		fmt.Println("Button clicked!")
	})

	// Trigger the "user_input" event
	registry.Publish("user_input")

	// Trigger the "button_click" event
	registry.Publish("button_click")
}
