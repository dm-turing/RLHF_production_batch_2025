package main

import (
	"fmt"
	"time"
)

func doWork() (result int, err error) {
	time.Sleep(time.Second)
	return 42, nil
}

func main() {
	start := time.Now()
	result, err := doWork()
	if err != nil {
		fmt.Println("Error:", err)
	} else {
		fmt.Println("Result:", result)
	}
	fmt.Println("Execution time:", time.Since(start))
}
