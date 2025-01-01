package main

import (
	"fmt"
	"io"
	"os"
	"strings"
)

func main() {
	// Create multiple readers for different files
	file1, _ := os.Open("file1.txt")
	file2, _ := os.Open("file2.txt")
	defer file1.Close()
	defer file2.Close()

	r := io.MultiReader(file1, file2, strings.NewReader("Hello, World!"))

	// Read from the combined stream
	buf := make([]byte, 1024)
	n, err := r.Read(buf)
	fmt.Println(string(buf[:n])) // Output: Data from file1.txtData from file2.txtHello, World!
	if err != nil && err != io.EOF {
		fmt.Println("Error reading:", err)
	}
}
