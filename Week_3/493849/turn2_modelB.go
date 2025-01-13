package main

import (
	"fmt"
	"io/ioutil"
	"os"
)

const maxFileDescriptors = 100

var filePool = make(chan *os.File, maxFileDescriptors)

func init() {
	// Initialize the file pool with empty file descriptors
	for i := 0; i < maxFileDescriptors; i++ {
		filePool <- nil
	}
}
func openFileFromPool(filename string) (*os.File, error) {
	// Wait for an available file descriptor from the pool
	file := <-filePool
	if file == nil {
		// If the file is nil, it means we need to open a new file
		var err error
		file, err = os.Open(filename)
		if err != nil {
			return nil, err
		}
	}
	return file, nil
}
func closeFileToPool(file *os.File) {
	// Close the file and return it to the pool
	file.Close()
	filePool <- file
}
func main() {
	filename := "example.txt"
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		fmt.Println("Error reading file:", err)
		return
	}
	fmt.Println(string(data))
}
