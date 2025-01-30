package main

import (
	"fmt"
	// "time"
)

func main() {
	var funcs []func() string

	// Incorrect - loop variable `i` is captured by reference
	for i := 0; i < 3; i++ {
		funcs = append(funcs, func() string {
			return fmt.Sprintf("Value: %d", i)
		})
	}

	for _, fn := range funcs {
		fmt.Println(fn())
	}
}
