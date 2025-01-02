package main

import (
	"fmt"
	"math/rand"
	"sort"
	"sync"
	"time"
)

func concurrentSortInts(nums []int, wg *sync.WaitGroup) {
	defer wg.Done()
	sort.Slice(nums, func(i, j int) bool {
		return nums[i] < nums[j]
	})
}

func main() {
	sizes := []int{1_000_000, 5_000_000, 10_000_000}
	for _, size := range sizes {
		nums := generateRandomInts(size)

		// Sequential sorting
		start := time.Now()
		sort.Slice(nums, func(i, j int) bool {
			return nums[i] < nums[j]
		})
		sequentialTime := time.Since(start)

		// Concurrent sorting
		start = time.Now()
		const numGoroutines = 8
		var wg sync.WaitGroup
		wg.Add(numGoroutines)
		for i := 0; i < numGoroutines; i++ {
			startIndex := i * size / numGoroutines
			endIndex := (i + 1) * size / numGoroutines
			go concurrentSortInts(nums[startIndex:endIndex], &wg)
		}
		wg.Wait()
		concurrentTime := time.Since(start)

		fmt.Printf("Size: %d, Sequential Time: %s, Concurrent Time: %s, Speedup: %.2fx\n",
			size, sequentialTime, concurrentTime, sequentialTime.Seconds()/concurrentTime.Seconds())
	}
}

func generateRandomInts(size int) []int {
	rand.Seed(time.Now().UnixNano())
	nums := make([]int, size)
	for i := 0; i < size; i++ {
		nums[i] = rand.Intn(size)
	}
	return nums
}
