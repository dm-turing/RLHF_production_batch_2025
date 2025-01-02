package main

import (
	"fmt"
	"log"
	"math/rand"
	"os"
	"sync"
	"time"
)

var (
	file        *os.File
	metricsMap  = make(map[string]*Metrics)
	metricsLock sync.Mutex
)

type Metrics struct {
	Sucesses            int64
	Failures            int64
	AverageResponseTime time.Duration
}

func init() {
	file, err := os.OpenFile("abtest.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal("Error opening log file:", err)
	}
	log.SetOutput(file)
	log.SetFlags(log.Ldate | log.Ltime | log.Lmicroseconds)
}
func logEvent(testName string, event string, duration time.Duration) {
	metricsLock.Lock()
	defer metricsLock.Unlock()
	// Update the metrics
	metrics, ok := metricsMap[testName]
	if !ok {
		metrics = &Metrics{}
		metricsMap[testName] = metrics
	}
	if event == "success" {
		metrics.Sucesses++
	} else if event == "failure" {
		metrics.Failures++
	}
	metrics.AverageResponseTime += duration / 2

	log.Printf("Test: %s, Event: %s, Duration: %s\n", testName, event, duration)
}

// Run the A/B test
func simulateABTest() {
	testName := "my_test"
	for i := 0; i < 10; i++ {
		time.Sleep(time.Duration(rand.Intn(500)) * time.Millisecond)
		duration := time.Duration(rand.Intn(1000)) * time.Millisecond
		// Simulate test success or failure based on some random condition
		if rand.Intn(2) == 0 {
			logEvent(testName, "success", duration)
		} else {
			logEvent(testName, "failure", duration)
		}
	}
}

// Calculate and display the statistics
func showInsights() {
	metricsLock.Lock()
	defer metricsLock.Unlock()
	for testName, metrics := range metricsMap {
		successRate := float64(metrics.Sucesses) / float64(metrics.Sucesses+metrics.Failures)
		averageResponseTime := metrics.AverageResponseTime.Truncate(time.Millisecond)
		fmt.Printf("\nTest: %s\n", testName)
		fmt.Printf("====================\n")
		fmt.Printf("Success Rate: %.2f%%\n", successRate*100)
		fmt.Printf("Average Response Time: %s\n", averageResponseTime)
	}
}

func main() {
	defer file.Close()
	simulateABTest()
	// Wait for a bit to allow all test events to be logged
	time.Sleep(2 * time.Second)
	showInsights()
}
