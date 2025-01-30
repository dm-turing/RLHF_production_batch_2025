package main

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
)

// Operator is a function type representing a binary operation
type Operator func(a, b float64) float64

// Curried function to create closure over operator
func curryOperation(op Operator) func(float64) func(float64) float64 {
	return func(a float64) func(float64) float64 {
		return func(b float64) float64 {
			return op(a, b)
		}
	}
}

// Parse and evaluate the expression using closures
func evaluateExpression(expr string) (float64, error) {
	// Supported operators
	operators := map[string]Operator{
		"+": func(a, b float64) float64 { return a + b },
		"-": func(a, b float64) float64 { return a - b },
		"*": func(a, b float64) float64 { return a * b },
		"/": func(a, b float64) float64 { if b == 0 { return 0 } return a / b },
	}

	// Parse expression tokens
	tokens := strings.Fields(expr)

	// Validate the number of tokens
	if len(tokens)%2 == 0 {
		return 0, errors.New("invalid expression: must have an odd number of tokens")
	}

	// Initial value is the first number
	value, err := strconv.ParseFloat(tokens[0], 64)
	if err != nil {
		return 0, fmt.Errorf("invalid number: %s", tokens[0])
	}

	// Traverse tokens and apply operations
	for i := 1; i < len(tokens); i += 2 {
		opToken := tokens[i]
		numToken := tokens[i+1]

		// Validate operator
		opFunc, ok := operators[opToken]