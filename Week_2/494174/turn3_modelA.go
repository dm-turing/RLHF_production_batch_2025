package main

import (
	"fmt"
	"os"
	"sync"
)

func processFileConcurrently(filenames []string, wg *sync.WaitGroup) {
	defer wg.Done()

	for _, filename := range filenames {
		_, err := os.ReadFile(filename)
		if err != nil {
			fmt.Println("Error reading file:", err)
			continue
		}
		// processData(data)
	}
}

func main() {
	filenames := []string{"file1.txt", "file2.txt", "file3.txt"}
	wg := &sync.WaitGroup{}

	for _, filename := range filenames {
		wg.Add(1)
		go processFileConcurrently([]string{filename}, wg)
	}

	wg.Wait()
}
