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
	wg.Add(1)

	go func() {
		defer wg.Done()
		// Read the entire MultiReader into a buffer
		buf := make([]byte, 1024)
		for {
			n, err := mr.Read(buf)
			if err != nil {
				if err == io.EOF {
					break
				}
				fmt.Println(err)
				return
			}
			// Process the read data
			fmt.Print(string(buf[:n]))
		}
	}()
	wg.Wait()
}
