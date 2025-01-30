package main

import (
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
		"/": func(a, b float64) float64 {
			if b != 0 {
				return a / b
			} else {
				fmt.Println("Cannot divide by 0")
				return 0
			}
		},
	}

	// Parse expression tokens
	tokens := strings.Fields(expr)

	// Initial value is the first number
	value, err := strconv.ParseFloat(tokens[0], 64)
	if err != nil {
		return 0, fmt.Errorf("invalid number: %s", tokens[0])
	}

	// Traverse tokens and apply operations
	for i := 1; i < len(tokens); i += 2 {
		opToken := tokens[i]
		numToken := tokens[i+1]

		num, err := strconv.ParseFloat(numToken, 64)
		if err != nil {
			return 0, fmt.Errorf("invalid number: %s", numToken)
		}

		opFunc, ok := operators[opToken]
		if !ok {
			return 0, fmt.Errorf("invalid operator: %s", opToken)
		}

		// Curried application of operation
		carryOperation := curryOperation(opFunc)(value)
		value = carryOperation(num)
	}

	return value, nil
}

// Entry point
func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run main.go <expression>")
		return
	}
	expression := os.Args[1]
	result, err := evaluateExpression(expression)
	if err != nil {
		fmt.Println("Error:", err)
	} else {
		fmt.Println("Result:", result)
	}
}
