package main

import (
	"testing"
	"time"
)

// simulateProcessing simulates the processing of a URL by adding a delay
func simulateProcessing(url string) {
	// Simulate processing time based on URL length
	time.Sleep(time.Millisecond * time.Duration(len(url)/5))
}

func BenchmarkURLProcessing(b *testing.B) {
	minLengthUrls := make([]string, 0)
	maxLengthUrls := make([]string, 0)

	// Generate test URLs
	for i := 0; i < 1000; i++ {
		minLengthUrls = append(minLengthUrls, generateRandomURL(75, 100))
		maxLengthUrls = append(maxLengthUrls, generateRandomURL(200, 255))
	}

	// Benchmark URLs of length 75-100
	b.Run("MinLengthURLs", func(b *testing.B) {
		for _, url := range minLengthUrls {
			simulateProcessing(url)
		}
	})

	// Benchmark URLs of length 200-255
	b.Run("MaxLengthURLs", func(b *testing.B) {
		for _, url := range maxLengthUrls {
			simulateProcessing(url)
		}
	})
}
