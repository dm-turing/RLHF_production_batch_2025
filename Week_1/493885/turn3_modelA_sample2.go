package main

import (
	"fmt"
	"io"
	"os"
	"sync"
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

	var wg sync.WaitGroup

	// Read file1 first
	wg.Add(1)
	go func() {
		defer wg.Done()
		buf := make([]byte, 1024)
		n, err := io.ReadFull(file1, buf)
		if err != nil && err != io.EOF {
			fmt.Println(err)
			return
		}
		fmt.Println("File 1:", string(buf[:n]))
	}()

	// Read file2 after file1
	wg.Add(1)
	go func() {
		defer wg.Done()
		buf := make([]byte, 1024)
		n, err := io.ReadFull(file2, buf)
		if err != nil && err != io.EOF {
			fmt.Println(err)
			return
		}
		fmt.Println("File 2:", string(buf[:n]))
	}()

	wg.Wait()
}
