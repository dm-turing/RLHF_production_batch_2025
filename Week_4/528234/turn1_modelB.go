package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"sync"
	"time"
)

// Worker represents a single worker in the load test.
type Worker struct {
	id     int
	wg     *sync.WaitGroup
	ctx    context.Context
	cancel context.CancelFunc
}

// NewWorker creates a new Worker with the given ID and WaitGroup.
func NewWorker(id int, wg *sync.WaitGroup) *Worker {
	ctx, cancel := context.WithCancel(context.Background())
	return &Worker{
		id:     id,
		wg:     wg,
		ctx:    ctx,
		cancel: cancel,
	}
}

// Start initiates the worker's load testing process.
func (w *Worker) Start() {
	w.wg.Add(1)
	go w.run()
}

// Stop signals the worker to stop and gracefully clean up.
func (w *Worker) Stop() {
	w.cancel()
}

// run is the main loop of the worker, responsible for executing the load test.
func (w *Worker) run() {
	defer w.wg.Done()
	for {
		select {
		case <-w.ctx.Done():
			// Graceful cleanup: Close any resources here
			fmt.Printf("Worker %d: Stopped\n", w.id)
			return
		default:
			// Simulate work by making an HTTP request
			w.makeRequest()
			time.Sleep(time.Duration(rand.Intn(100)) * time.Millisecond)
		}
	}
}

func (w *Worker) makeRequest() {
	resp, err := http.Get("https://github.com")
	if err != nil {
		fmt.Printf("Worker %d: Error making request: %v\n", w.id, err)
		return
	}
	defer resp.Body.Close()
	_, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("Worker %d: Error reading response: %v\n", w.id, err)
	}
}

// LoadTest orchestrates the load test by starting and stopping workers.
func LoadTest(numWorkers int, duration time.Duration) {
	var wg sync.WaitGroup
	workers := make([]*Worker, numWorkers)

	// Create workers
	for i := 0; i < numWorkers; i++ {
		workers[i] = NewWorker(i+1, &wg)
	}

	// Start workers
	for _, worker := range workers {
		worker.Start()
	}

	// Wait for the specified duration
	time.Sleep(duration)

	// Stop workers gracefully
	for _, worker := range workers {
		worker.Stop()
	}

	// Wait for all workers to complete
	wg.Wait()
	fmt.Println("Load test completed.")
}

func main() {
	numWorkers := 10
	duration := 5 * time.Second
	LoadTest(numWorkers, duration)
}
