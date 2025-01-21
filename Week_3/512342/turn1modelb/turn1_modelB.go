package main

import (
	"fmt"
	"math/rand"
	"time"
)

type MouseEvent struct {
	Type string // e.g., "Click", "DoubleClick", "Move", "Drag"
	X    int
	Y    int
	// Add other relevant mouse event attributes as needed
}

type MouseEventCallback func(event MouseEvent)

type EventHandler struct {
	callbacks []MouseEventCallback
}

func NewEventHandler() *EventHandler {
	return &EventHandler{callbacks: []MouseEventCallback{}}
}

func (h *EventHandler) RegisterCallback(callback MouseEventCallback) {
	h.callbacks = append(h.callbacks, callback)
}

func (h *EventHandler) HandleEvent(event MouseEvent) {
	for _, callback := range h.callbacks {
		callback(event)
	}
}

type GuiApplication interface {
	Run() error
	RegisterMouseEventCallback(callback MouseEventCallback)
}

type PollingGuiApplication struct {
	eventHandler *EventHandler
	// Add other GUI-specific fields
}

func NewPollingGuiApplication() *PollingGuiApplication {
	app := &PollingGuiApplication{
		eventHandler: NewEventHandler(),
	}
	// Initialize GUI components here
	return app
}

func (a *PollingGuiApplication) Run() error {
	// Main event loop
	for {
		event := a.pollGuiEvent()
		if event.Type != "" {
			a.eventHandler.HandleEvent(event)
		}
		// Add sleep or yield here to prevent tight looping
	}
}

func (a *PollingGuiApplication) pollGuiEvent() MouseEvent {
	// Implement GUI-specific event polling logic
	// For this example, let's simulate polling with a timer
	event := MouseEvent{}
	time.Sleep(100 * time.Millisecond)
	event.Type = "Move"
	event.X = rand.Intn(100)
	event.Y = rand.Intn(100)
	return event
}

func (a *PollingGuiApplication) RegisterMouseEventCallback(callback MouseEventCallback) {
	a.eventHandler.RegisterCallback(callback)
}

// Callback to log mouse events
func logMouseEvent(event MouseEvent) {
	fmt.Printf("Mouse Event: %+v\n", event)
}

// Callback to handle drag events
func handleDragEvent(event MouseEvent) {
	if event.Type == "Drag" {
		fmt.Println("Dragging...")
	}
}

func main() {
	app := NewPollingGuiApplication()

	// Register callbacks
	app.RegisterMouseEventCallback(logMouseEvent)
	app.RegisterMouseEventCallback(handleDragEvent)

	// Start the event loop
	if err := app.Run(); err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
}
