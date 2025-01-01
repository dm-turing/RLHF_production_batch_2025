package main

import (
	"bytes"
	"fmt"
	"io"
	"os"
)

func main() {
	// Create a buffer to accumulate data
	var buf bytes.Buffer

	// Define the writers (file1 and file2)
	file1, err := os.Create("file1.txt")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer file1.Close()

	file2, err := os.Create("file2.txt")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer file2.Close()

	// Create a MultiWriter from the buffer and the file writers
	mw := io.MultiWriter(&buf, file1, file2)