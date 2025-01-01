package main

import (
	"fmt"
	"io"
	"os"
)

func main() {
	// Create multiple writers for different files
	file1, _ := os.Create("file1_copy.txt")
	file2, _ := os.Create("file2_copy.txt")
	defer file1.Close()
	defer file2.Close()

	w := io.MultiWriter(file1, file2, os.Stdout)

	// Write to the combined stream
	data := "This data will be written to all three destinations."
	_, err := w.Write([]byte(data))
	if err != nil {
		fmt.Println("Error writing:", err)
	}
}
