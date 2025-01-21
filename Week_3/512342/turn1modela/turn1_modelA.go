package main

type MouseEvent struct {
	EventType string // "click", "move", "drag", etc.
	X, Y      int    // Position of the mouse event
}

type MouseCallback func(MouseEvent)

type MouseEventChain struct {
	Callbacks []MouseCallback
}

func (chain *MouseEventChain) AddCallback(callback MouseCallback) {
	chain.Callbacks = append(chain.Callbacks, callback)
}

func (chain *MouseEventChain) HandleEvent(event MouseEvent) {
	for _, callback := range chain.Callbacks {
		callback(event)
	}
}

func main() {
	chain := &MouseEventChain{}

	// Example of adding callbacks
	chain.AddCallback(func(event MouseEvent) {
		if event.EventType == "click" {
			// Handle click event
			println("Click at:", event.X, event.Y)
		}
	})

	chain.AddCallback(func(event MouseEvent) {
		if event.EventType == "move" {
			// Handle move event
			println("Move at:", event.X, event.Y)
		}
	})

	// Dispatch an event
	chain.HandleEvent(MouseEvent{EventType: "click", X: 100, Y: 150})
	chain.HandleEvent(MouseEvent{EventType: "move", X: 110, Y: 160})
}
