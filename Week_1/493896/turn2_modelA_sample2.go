package main

import (
	"fmt"
	"math/rand"
	"runtime"
	"sort"
	"sync"
	"time"
)

func sortStringsConcurrently(slice []string, numPartitions int) []string {
	var wg sync.WaitGroup
	result := make([]string, len(slice))
	segmentSize := len(slice) / numPartitions
	remainder := len(slice) % numPartitions

	for i := 0; i < numPartitions; i++ {
		start := i * segmentSize
		end := start + segmentSize
		if i < remainder {
			end++
		}

		wg.Add(1)
		go func() {
			partition := slice[start:end]
			sort.Strings(partition)
			result = partition
			wg.Done()
		}()
	}

	wg.Wait()
	return result
}

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	rand.Seed(time.Now().UnixNano())

	words := make([]string, 100000)
	for i := range words {
		words[i] = fmt.Sprintf("word%05d", i)
	}

	start := time.Now()
	fmt.Println("Original words:", words[:10])

	sortedWords := sortStringsConcurrently(words, runtime.NumCPU())
	fmt.Println("Sorted words:", sortedWords[:10])

	elapsed := time.Since(start)
	fmt.Printf("Concurrent sorting time: %s\n", elapsed)

	start = time.Now()
	sort.Strings(words)
	elapsed = time.Since(start)
	fmt.Printf("Sequential sorting time: %s\n", elapsed)
}
