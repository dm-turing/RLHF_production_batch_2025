package main

type MouseError struct {
	ErrorType string // "invalid_position", "unsupported_event", etc.
	Message   string // Error details
}

type ErrorCallback func(MouseError)

type MouseErrorChain struct {
	Callbacks []ErrorCallback
}

func (chain *MouseErrorChain) AddCallback(callback ErrorCallback) {
	chain.Callbacks = append(chain.Callbacks, callback)
}

func (chain *MouseErrorChain) HandleError(err MouseError) {
	for _, callback := range chain.Callbacks {
		callback(err)
	}
}

func main() {
	errorChain := &MouseErrorChain{}

	// Example of adding error callbacks
	errorChain.AddCallback(func(err MouseError) {
		if err.ErrorType == "invalid_position" {
			// Handle invalid position error
			println("Error: Invalid Position -", err.Message)
		}
	})

	errorChain.AddCallback(func(err MouseError) {
		if err.ErrorType == "unsupported_event" {
			// Handle unsupported event error
			println("Error: Unsupported Event -", err.Message)
		}
	})

	// Dispatch an error
	errorChain.HandleError(MouseError{ErrorType: "invalid_position", Message: "Coordinates out of bounds"})
	errorChain.HandleError(MouseError{ErrorType: "unsupported_event", Message: "Event type not handled"})
}
