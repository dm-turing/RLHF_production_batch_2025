package main

import (
	"fmt"
	"strings"
)

// Using strings.Join
func main() {
	lines := []string{
		"This is a multiline string",
		"that can be used for logging purposes",
		"or any other reason where a multi-line",
		"string is required.",
	}
	multilineString := strings.Join(lines, "\n")
	fmt.Println(multilineString)
}
