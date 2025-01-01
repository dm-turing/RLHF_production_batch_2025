package main

import (
	"fmt"
	"io"
	"os"
	"sync"
)

func main() {
	file1, err := os.Open("file1.txt")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer file1.Close()

	file2, err := os.Open("file2.txt")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer file2.Close()

	// Create a MultiReader from the two file readers
	mr := io.MultiReader(file1, file2)

	var wg sync.WaitGroup

	// Concurrently read from the MultiReader
	for {
		buf := make([]byte, 1024)
		n, err := mr.Read(buf)
		if err != nil && err != io.EOF {
			fmt.Println(err)
			return
		}
		if n == 0 {
			break
		}

		// Lock the output to avoid race conditions
		wg.Add(1)
		go func() {
			defer wg.Done()
			fmt.Println(string(buf))
		}()
	}

	wg.Wait()
}
