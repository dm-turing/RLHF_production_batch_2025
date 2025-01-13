package main

import (
	"fmt"
	"os"
)

func main() {
	file, err := os.Open("test.txt")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer file.Close()
	// Perform some file I/O operations here...
	// ...
	// Now, the application tries to reuse the same file descriptor
	// for another operation
	file, err = os.Open("test2.txt")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer file.Close()
	// Perform some more file I/O operations here...
}
