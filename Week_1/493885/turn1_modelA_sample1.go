package main

import (
	"fmt"
	"io"
	"os"
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

	// Read the entire MultiReader into a buffer
	buf := make([]byte, 1024)
	_, err = io.ReadFull(mr, buf)
	if err != nil && err != io.EOF {
		fmt.Println(err)
		return
	}

	// Print the combined content
	fmt.Println(string(buf))
}
