package main

import (
	"log"
	"os"
)

// With error handling
func foo() error {
	f, err := os.Open("example1.txt")
	if err != nil {
		log.Fatalln("Error processing file:", err) // Properly handles errors and ensures program termination
	}
	defer f.Close()
	// ... other operations
	return nil
}

func main() {
	foo()
}
