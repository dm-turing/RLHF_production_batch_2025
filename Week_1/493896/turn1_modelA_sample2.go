package main

import (
	"fmt"
	"sort"
)

func main() {
	words := []string{"banana", "apple", "orange", "grape"}
	fmt.Println("Original words:", words)

	sort.Strings(words)

	fmt.Println("Sorted words:", words)
}
