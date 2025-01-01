package main

import (
	"bufio"
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

	// Use a buffer to read lines efficiently
	br := bufio.NewReader(mr)

	for {
		line, isPrefix, err := br.ReadLine()
		if err != nil {
			if err == io.EOF {
				break
			}
			fmt.Println(err)
			return
		}
		if isPrefix {
			fmt.Println("Partial line:", string(line))
		} else {
			fmt.Println("Full line:", string(line))
		}
	}
}
