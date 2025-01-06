package main

import (
	"context"
	"fmt"
	"log"
	"time"
)

func worker(ctx context.Context, id int) error {
	select {
	case <-time.After(time.Second):
		if id%2 == 0 {
			return fmt.Errorf("worker %d failed", id)
		}
	case <-ctx.Done():
		return ctx.Err()
	}
	return nil
}

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	for i := 0; i < 10; i++ {
		go func(id int) {
			err := worker(ctx, id)
			if err != nil {
				log.Println("Worker error:", err)
			}
		}(i)
	}

	time.Sleep(2 * time.Second) // Simulate premature termination
}
