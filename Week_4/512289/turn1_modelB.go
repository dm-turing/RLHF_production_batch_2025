package main

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"math/big"
	"testing"
)

func generateRandomURL(minLength, maxLength int) string {
	// Generate a random length between minLength and maxLength
	length, _ := rand.Int(rand.Reader, big.NewInt(int64(maxLength-minLength+1)))
	length = length.Add(length, big.NewInt(int64(minLength)))

	// Generate random bytes
	bytes := make([]byte, length.Int64())
	_, _ = rand.Read(bytes)

	// Base64 encode the bytes to get a URL-safe string
	return base64.URLEncoding.EncodeToString(bytes)
}

func processURLs(urls []string) {
	// Simple processing: iterate over the URLs and print them
	for _, url := range urls {
		url = url
		// fmt.Println(url)
	}
}

func BenchmarkProcessURLs(b *testing.B) {
	// Generate a moderate number of URLs for each length range
	const numURLs = 1000

	// Length ranges to benchmark
	lengthRanges := []struct {
		min, max int
	}{
		{75, 100},
		{200, 255},
	}

	for _, lr := range lengthRanges {
		b.Run(fmt.Sprintf("%d-%d", lr.min, lr.max), func(b *testing.B) {
			urls := make([]string, numURLs)
			for i := 0; i < numURLs; i++ {
				urls[i] = generateRandomURL(lr.min, lr.max)
			}

			b.ResetTimer()
			for n := 0; n < b.N; n++ {
				processURLs(urls)
			}
		})
	}
}

func main() {
	// Generate and print a moderately large sample space of random URLs
	fmt.Println("Generating sample URLs...")
	sampleSize := 10000
	minLength := 75
	maxLength := 255
	urls := make([]string, sampleSize)
	for i := 0; i < sampleSize; i++ {
		urls[i] = generateRandomURL(minLength, maxLength)
	}
	fmt.Println("Sample URLs generated.")

	// Process the sample URLs to measure the average server processing time
	fmt.Println("Processing sample URLs...")
	processURLs(urls)
	fmt.Println("Sample URLs processed.")

	// Run the benchmarks
	fmt.Println("Running benchmarks...")
	testing.Benchmark(BenchmarkProcessURLs)
	fmt.Println("Benchmarks completed.")
}
