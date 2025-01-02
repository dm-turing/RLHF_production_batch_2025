package main

import (
	"fmt"
	"math/rand"
	"runtime"
	"sort"
	"sync"
	"time"
)

func sortSliceConcurrently(slice []int, numPartitions int) []int {
	var wg sync.WaitGroup
	result := make([]int, len(slice))
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
			sort.Slice(partition, func(i, j int) bool { return partition[i] < partition[j] })
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

	numbers := make([]int, 100000)
	for i := range numbers {
		numbers[i] = rand.Intn(100000)
	}

	start := time.Now()
	fmt.Println("Original numbers:", numbers[:10])

	sortedNumbers := sortSliceConcurrently(numbers, runtime.NumCPU())
	fmt.Println("Sorted numbers:", sortedNumbers[:10])

	elapsed := time.Since(start)
	fmt.Printf("Concurrent sorting time: %s\n", elapsed)

	start = time.Now()
	sort.Slice(numbers, func(i, j int) bool { return numbers[i] < numbers[j] })
	elapsed = time.Since(start)
	fmt.Printf("Sequential sorting time: %s\n", elapsed)
}
