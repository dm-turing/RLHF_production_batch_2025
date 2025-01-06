package main

import (
	"fmt"
	"os"
)

func readFile(filename string) ([]byte, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	return data, nil
}

func main() {
	filename := "example.txt"
	data, err := readFile(filename)

	if err != nil {
		fmt.Println("Error reading file:", err)
		return
	}

	fmt.Println("File content:", string(data))
}
