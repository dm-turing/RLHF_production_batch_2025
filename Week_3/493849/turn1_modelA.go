package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"sync"
	"time"
)

func readFile(filename string, mu *sync.Mutex, ctx context.Context) error {
	mu.Lock()
	defer mu.Unlock()

	// Check if the context is canceled before proceeding
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}

	file, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	data, err := ioutil.ReadAll(file)
	if err != nil {
		return err
	}

	fmt.Println(string(data))

	return nil
}

func main() {
	filename := "example.txt"
	mu := &sync.Mutex{}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Simulating concurrent access to the file
	go readFile(filename, mu, ctx)
	go readFile(filename, mu, ctx)

	// Gracefully cancel the context after a delay
	time.Sleep(2 * time.Second)
	cancel()
}
