package main

import (
	"fmt"
	"math/rand"
	"strconv"
	"time"
)

func main() {
	// Seed the random number generator
	rand.Seed(time.Now().UnixNano())

	// Generate a map of type `map[string]int` with 10,000 random values
	dataMap := make(map[string]int)
	keysSlice := make([]string, 0, 10000)

	for i := 0; i < 10000; i++ {
		key := "key" + strconv.Itoa(i)
		value := rand.Intn(10000)
		dataMap[key] = value
		keysSlice = append(keysSlice, key)
	}

	// Iterate through the keys using the slice and print to the console
	fmt.Println("Iterating through keys in slice:")

	startSlice := time.Now()
	for _, key := range keysSlice {
		fmt.Println(key, dataMap[key])
	}
	elapsedSlice := time.Since(startSlice)
	fmt.Printf("Elapsed time for iterating through slice: %s\n", elapsedSlice)

	// Iterate through the keys in the map and print to the console
	fmt.Println("Iterating through keys in map:")

	startMap := time.Now()
	for key, value := range dataMap {
		fmt.Println(key, value)
	}
	elapsedMap := time.Since(startMap)
	fmt.Printf("Elapsed time for iterating through map: %s\n", elapsedMap)

	// Compare the elapsed times
	fmt.Println("Comparison of elapsed times:")
	fmt.Printf("Slice iteration: %s, Map iteration: %s\n", elapsedSlice, elapsedMap)
}