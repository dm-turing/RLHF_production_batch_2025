package main

import (
	"fmt"
	"math/rand"
	"runtime"
	"sort"
	"sync"
	"time"
)

type Person struct {
	Name string
	Age  int
}

func sortPeopleConcurrently(people []Person, numPartitions int) []Person {
	var wg sync.WaitGroup
	result := make([]Person, len(people))
	segmentSize := len(people) / numPartitions
	remainder := len(people) % numPartitions

	for i := 0; i < numPartitions; i++ {
		start := i * segmentSize
		end := start + segmentSize
		if i < remainder {
			end++
		}

		wg.Add(1)
		go func() {
			partition := people[start:end]
			sort.Slice(partition, func(i, j int) bool { return partition[i].Name < partition[j].Name })
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

	people := make([]Person, 100000)
	for i := range people {
		people[i].Name = fmt.Sprintf("person%05d", i)
		people[i].Age = rand.Intn(100)
	}

	start := time.Now()
	fmt.Println("Original people:", people[:10])

	sortedPeople := sortPeopleConcurrently(people, runtime.NumCPU())
	fmt.Println("Sorted people:", sortedPeople[:10])

	elapsed := time.Since(start)
	fmt.Printf("Concurrent sorting time: %s\n", elapsed)

	start = time.Now()
	sort.Slice(people, func(i, j int) bool { return people[i].Name < people[j].Name })
	elapsed = time.Since(start)
	fmt.Printf("Sequential sorting time: %s\n", elapsed)
}
