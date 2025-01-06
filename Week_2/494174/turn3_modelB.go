package main

import (
	"errors"
	"fmt"
	"sync"
	"time"
)

// Example function that can return an error
func processData(wg *sync.WaitGroup) error {
	defer wg.Done()
	time.Sleep(time.Second) // Simulate some work
	return errors.New("data processing failed")
}

func main() {
	var wg sync.WaitGroup
	errorChannel := make(chan error, 1)

	numRoutines := 5
	wg.Add(numRoutines)

	// Start multiple concurrent routines to process data
	for i := 0; i < numRoutines; i++ {
		go func() {
			err := processData(&wg)
			if err != nil {
				errorChannel <- err
			}
		}()
	}

	// Wait for all routines to finish
	wg.Wait()
	close(errorChannel)

	// Collect and handle all errors that occurred concurrently
	var errorsOccurred []error
	go func() {
		for {
			err := <-errorChannel
			errorsOccurred = append(errorsOccurred, err)
		}
	}()

	if len(errorsOccurred) > 0 {
		fmt.Println("Errors occurred concurrently:")
		for _, err := range errorsOccurred {
			fmt.Println("-", err)
		}
	}
}
