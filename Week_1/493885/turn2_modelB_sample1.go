package main

import (
	"fmt"
	"io"
	"os"
	"sync"
)

func main() {
	// Define the number of files to read
	const numFiles = 3

	// Create a WaitGroup to track the number of readers
	var wg sync.WaitGroup
	wg.Add(numFiles)

	// Create a shared data structure to store the results
	results := make([]string, numFiles)

	for i := 0; i < numFiles; i++ {
		go readFile(i, &results, &wg)
	}

	// Wait for all readers to complete
	wg.Wait()

	// Display the combined results
	fmt.Println("Combined Results:")
	for _, result := range results {
		fmt.Println(result)
	}
}

func readFile(fileIndex int, results *[]string, wg *sync.WaitGroup) {
	// Defer the WaitGroup Done call to signal that the reader has finished
	defer wg.Done()

	fileName := fmt.Sprintf("file%d.txt", fileIndex+1)
	file, err := os.Open(fileName)
	if err != nil {
		fmt.Println("Error opening file:", fileName, err)
		return
	}
	defer file.Close()

	buf, err := io.ReadAll(file)
	if err != nil {
		fmt.Println("Error reading file:", fileName, err)
		return
	}

	// Lock the shared data structure to ensure thread safety
	// when updating the results slice.
	resultsMu.Lock()
	defer resultsMu.Unlock()
	results[fileIndex] = string(buf)
}

// Define a global mutex to protect the results slice
var resultsMu sync.Mutex
