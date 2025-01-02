package main

import (
	"fmt"
	"sort"
)

func main() {
	scores := map[string]int{"Alice": 80, "Bob": 90, "Charlie": 70}
	fmt.Println("Original map:", scores)

	keys := make([]string, 0, len(scores))
	for key := range scores {
		keys = append(keys, key)
	}

	sort.Strings(keys)

	fmt.Println("Sorted keys:", keys)
	for _, key := range keys {
		fmt.Println(key, ":", scores[key])
	}
}
