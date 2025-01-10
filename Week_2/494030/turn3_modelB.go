package main

import (
	"fmt"
	"time"
)

func main() {
	//Example without callbacks
	data := make([]int, 10000000)
	for i := 0; i < len(data); i++ {
		data[i] = i
	}
	start := time.Now()
	result1 := processDataWithoutCallback(data)
	end := time.Now()
	fmt.Println("Time taken without callbacks:", end.Sub(start))
	fmt.Println("Result 1: ", result1)

	//Example with callbacks
	start = time.Now()
	result2 := processDataWithCallback(data)
	end = time.Now()
	fmt.Println("Time taken with callbacks:", end.Sub(start))
	fmt.Println("Result 2: ", result2)
}

func processDataWithoutCallback(data []int) int {
	total := 0
	for _, value := range data {
		total += value
	}
	return total
}

func processDataWithCallback(data []int) int {
	total := 0
	processor := func(value int) {
		total += value
	}
	for _, value := range data {
		processor(value)
	}
	return total
}
