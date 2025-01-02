package main

import (
	"fmt"
	"sort"
)

func main() {
	numbers := []int{4, 2, 7, 1, 3}
	fmt.Println("Original numbers:", numbers)

	sort.Slice(numbers, func(i, j int) bool {
		return numbers[i] < numbers[j]
	})

	fmt.Println("Sorted numbers:", numbers)
}
