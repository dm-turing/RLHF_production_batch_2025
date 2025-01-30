package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

type Operator func(int, int) int

var operators map[string]Operator = map[string]Operator{
	"+": func(a, b int) int { return a + b },
	"-": func(a, b int) int { return a - b },
	"*": func(a, b int) int { return a * b },
	"/": func(a, b int) int { return a / b },
}

func parseExpr(expr string) (int, error) {
	tokens := strings.Split(expr, " ")
	if len(tokens) < 3 {
		return 0, fmt.Errorf("Invalid expression: %s", expr)
	}

	left, err := strconv.Atoi(tokens[0])
	if err != nil {
		return 0, fmt.Errorf("Invalid left operand: %s", tokens[0])
	}

	right, err := strconv.Atoi(tokens[2])
	if err != nil {
		return 0, fmt.Errorf("Invalid right operand: %s", tokens[2])
	}

	op := operators[tokens[1]]
	if op == nil {
		return 0, fmt.Errorf("Invalid operator: %s", tokens[1])
	}

	return op(left, right), nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: go run main.go <expression>")
		return
	}

	expr := os.Args[1]
	result, err := parseExpr(expr)
	if err != nil {
		fmt.Printf("Error: %s\n", err)
		return
	}

	fmt.Printf("Result: %d\n", result)
}
