package main

import (
	"fmt"
	"sync"
)

// Callback defines the signature of a callback function.
type Callback func(string)

// CallbackRegistry manages a collection of callbacks for a given event.
type CallbackRegistry struct {
	callbacks []Callback
	mu        sync.Mutex
}

// Register adds a new callback to the registry.
func (r *CallbackRegistry) Register(cb Callback) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.callbacks = append(r.callbacks, cb)
}

// Unregister removes a callback from the registry.
func (r *CallbackRegistry) Unregister(cb Callback) {
	r.mu.Lock()
	defer r.mu.Unlock()
	for i, c := range r.callbacks {
		if fmt.Sprintf("%v", c) == fmt.Sprintf("%v", cb) {
			copy(r.callbacks[i:], r.callbacks[i+1:])
			r.callbacks = r.callbacks[:len(r.callbacks)-1]
			return
		}
	}
}

// Invoke calls all registered callbacks with the provided argument.
func (r *CallbackRegistry) Invoke(arg string) {
	r.mu.Lock()
	defer r.mu.Unlock()
	for _, cb := range r.callbacks {
		cb(arg)
	}
}

func main() {
	var registry = &CallbackRegistry{}

	// Register some callbacks
	registry.Register(func(msg string) { fmt.Println("Callback 1:", msg) })
	registry.Register(func(msg string) { fmt.Println("Callback 2:", msg) })

	// Invoke the callbacks
	registry.Invoke("Hello from the main function!")

	// Unregister a callback
	registry.Unregister(func(msg string) { fmt.Println("Callback 2:", msg) })

	// Invoke the remaining callbacks
	registry.Invoke("Hello from the main function again!")
}
