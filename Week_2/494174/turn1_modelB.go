package main

import (
	"fmt"
)

func divide(x, y float64) (float64, error) {
	if y == 0 {
		return 0, fmt.Errorf("division by zero")
	}
	return x / y, nil
}

func main() {
	result, err := divide(10.0, 5.0)
	if err != nil {
		fmt.Println("Error:", err)
	} else {
		fmt.Println("Result:", result)
	}
	result, err = divide(10.0, 0.0)
	if err != nil {
		fmt.Println("Error:", err)
	} else {
		fmt.Println("Result:", result)
	}
}
