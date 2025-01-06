package main

import (
	"fmt"
	"os"
)

func readFile(filePath string) ([]byte, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("error opening file: %v", err)
	}
	defer file.Close()
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("error reading file: %v", err)
	}
	return data, nil
}

func main() {
	data, err := readFile("example.txt")
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(string(data))
}
