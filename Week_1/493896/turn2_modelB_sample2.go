package main

import (
	"fmt"
	"math/rand"
	"sort"
	"sync"
	"time"
)

func concurrentSortStrings(words []string, wg *sync.WaitGroup) {
	defer wg.Done()
	sort.Strings(words)
}

func main() {
	sizes := []int{1_000_000, 5_000_000, 10_000_000}
	for _, size := range sizes {
		words := generateRandomStrings(size)
  
		// Sequential sorting
		start := time.Now()
		sort.Strings(words)
		sequentialTime := time.Since(start)

		// Concurrent sorting
		start = time.Now()
		const numGoroutines = 8
		var wg sync.WaitGroup
		wg.Add(numGoroutines)
		for i := 0; i < numGoroutines; i++ {
			startIndex := i * size / numGoroutines
			endIndex := (i + 1) * size / numGoroutines
			go concurrentSortStrings(words[startIndex:endIndex], &wg)