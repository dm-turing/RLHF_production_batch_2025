// Ensuring reliable file reading in Go
package main

import (
	"fmt"
	"io/ioutil"
	"os"
)

func readFile(filename string) (string, error) {
	content, err := ioutil.ReadFile(filename)
	if err != nil {
		return "", fmt.Errorf("failed to read file: %w", err)
	}
	return string(content), nil
}

func main() {
	fileContent, err := readFile("someFile.txt")
	if err != nil {
		fmt.Printf("Error reading file: %s\n", err)
		os.Exit(1)
	}
	// Process file content only when there are no errors
	fmt.Println(fileContent)
}
