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

	// Update the keys in the slice 10,000 times randomly
	fmt.Println("Updating keys in slice:")
	startUpdateSlice := time.Now()
	for i := 0; i < 10000; i++ {
		index := rand.Intn(len(keysSlice))
		newKey := "newKey" + strconv.Itoa(i)
		keysSlice[index] = newKey
	}
	elapsedUpdateSlice := time.Since(startUpdateSlice)
	fmt.Printf("Elapsed time for updating keys in slice: %s\n", elapsedUpdateSlice)

	// Update the values in the map 10,000 times randomly
	fmt.Println("Updating values in map:")
	startUpdateMap := time.Now()
	for i := 0; i < 10000; i++ {
		key := "key" + strconv.Itoa(rand.Intn(10000))
		newValue := rand.Intn(10000)
		dataMap[key] = newValue
	}
	elapsedUpdateMap := time.Since(startUpdateMap)
	fmt.Printf("Elapsed time for updating values in map: %s\n", elapsedUpdateMap)

	// Compare the elapsed times for updates
	fmt.Println("Comparison of elapsed times for updates:")
	fmt.Printf("Slice update: %s, Map update: %s\n", elapsedUpdateSlice, elapsedUpdateMap)
}
