package main

import (
	"fmt"
	"net/http"
	"sync"
)

func makeRequest(wg *sync.WaitGroup, url string) {
	defer wg.Done()
	resp, err := http.Get(url)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	fmt.Println("Response Status:", resp.Status)
}

func main() {
	var wg sync.WaitGroup
	url := "http://localhost:8080/" // Replace with your URL
	for i := 0; i < 110; i++ {      // Adjust number of requests
		wg.Add(1)
		go makeRequest(&wg, url)
	}
	wg.Wait()
}
