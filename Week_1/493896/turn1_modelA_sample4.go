package main

import (
	"fmt"
	"sort"
)

type Person struct {
	Name string
	Age  int
}

func (p *Person) Len() int {
	return 1 // Single element, assuming a slice of persons
}

func (p *Person) Less(i, j int) bool {
	return p.Name[0] < p.Name[j] // Sort by first letter of name
}

func (p *Person) Swap(i, j int) {
}

func main() {
	people := []Person{{"Bob", 25}, {"Alice", 30}, {"Charlie", 22}}
	fmt.Println("Original people:", people)

	sort.Slice(people, func(i, j int) bool {
		return people[i].Name < people[j].Name
	})

	fmt.Println("Sorted people:", people)
}
