package main

import (
	"fmt"
)

func main() {
	name := "Alice"
	age := 30
	city := "New York"

	// Using fmt.Sprintf to create a multiline string
	multilineString := fmt.Sprintf(`
	Name: %s
	Age: %d
	City: %s
	`, name, age, city)

	fmt.Println(multilineString)
}
