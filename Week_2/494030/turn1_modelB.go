package main

import (
	"fmt"
	"strings"
)

// MiddlewareFunc defines the signature for a middleware function.
type MiddlewareFunc func(string) string

// Chain returns a new middleware function that chains together the given functions.
func Chain(functions ...MiddlewareFunc) MiddlewareFunc {
	return func(data string) string {
		result := data
		for _, f := range functions {
			result = f(result)
		}
		return result
	}
}

// UpperCase middleware converts the input string to uppercase.
func UpperCase(data string) string {
	return strings.ToUpper(data)
}

// Reverse middleware reverses the input string.
func Reverse(data string) string {
	runes := []rune(data)
	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
		runes[i], runes[j] = runes[j], runes[i]
	}
	return string(runes)
}

// Prefix middleware prepends the specified prefix to the input string.
func Prefix(prefix string) MiddlewareFunc {
	return func(data string) string {
		return prefix + data
	}
}

func main() {
	// Define the data processing pipeline using a chain of middleware functions.
	pipeline := Chain(UpperCase, Reverse, Prefix("PROCESSED:"))

	// Example data to process.
	data := "hello world"

	// Apply the pipeline to the data.
	result := pipeline(data)

	// Print the final result.
	fmt.Println(result) // Output: PROCESSED:DLROW OLLEH
}
