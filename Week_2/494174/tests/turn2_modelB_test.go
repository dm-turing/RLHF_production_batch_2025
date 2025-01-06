// Benchmarking error handling in Go
package main

import (
	"fmt"
	"testing"
)

func slowFunction(input int) (int, error) {
	if input < 0 {
		return 0, fmt.Errorf("input cannot be negative")
	}
	// Simulate some slow computation
	for i := 0; i < 100000; i++ {
	}
	return input * input, nil
}

func BenchmarkSlowFunctionWithErrorHandling(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, err := slowFunction(-1)
		if err != nil {
			b.Fatal(err)
		}
	}
}
