package main

import (
	"fmt"
	"os"
	"sync"
)

func main() {
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		file, err := os.Open("test.txt")
		if err != nil {
			fmt.Println(err)
			return
		}
		defer file.Close()
		// Perform some file I/O operations here...
		// ...
	}()
	// Wait for the first file I/O operation to complete
	wg.Wait()
	// Now, the application can safely reuse the file descriptor
	// for another operation
	file, err := os.Open("test2.txt")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer file.Close()
	// Perform some more file I/O operations here...
}
