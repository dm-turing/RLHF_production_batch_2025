package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
)

func main() {
	file, err := os.Open("large_file.txt")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer file.Close()

	bufReader := bufio.NewReaderSize(file, 4096)

	for {
		line, err := bufReader.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				break
			}
			fmt.Println(err)
			return
		}
		// Process the line
		fmt.Println(line)
	}
}
