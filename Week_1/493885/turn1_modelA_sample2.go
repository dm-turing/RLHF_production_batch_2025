package main

import (
	"fmt"
	"io"
	"os"
)

func main() {
	file, err := os.Create("output.txt")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer file.Close()

	// Create a MultiWriter from os.Stdout and the file writer
	mw := io.MultiWriter(os.Stdout, file)

	// Write to the MultiWriter, which will go to both stdout and the file
	fmt.Fprintf(mw, "This message will be written to both stdout and the file.\n")
}
