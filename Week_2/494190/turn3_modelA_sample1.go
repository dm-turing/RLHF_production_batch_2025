package main

import (
	"fmt"
	"log"
	"sync"
	"time"
)

func worker(id int, wg *sync.WaitGroup, done chan error) {
	defer wg.Done()
	if id%2 == 0 {
		done <- fmt.Errorf("worker %d failed", id)
	}
}

func main() {
	var wg sync.WaitGroup
	done := make(chan error)

	for i := 0; i < 10; i++ {
		wg.Add(1)
		go worker(i, &wg, done)
	}

	go func() {
		for {
			select {
			case err := <-done:
				log.Println("Error received:", err)
			case <-time.After(5 * time.Second):
				return
			}
		}
	}()

	wg.Wait()
	close(done)
}
