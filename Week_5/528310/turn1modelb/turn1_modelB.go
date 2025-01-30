package main

import (
	"fmt"
	"math/rand"
	"net/url"
	"strconv"
	"time"
)

const (
	baseURL         = "https://example.com/"
	maxURLLength    = 25
	shortURLLength  = 6
	queryParamsSize = 5
)

var (
	queryParams = []string{"id", "page", "sort", "filter", "limit"}
	urlChars    = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")
)

func main() {
	rand.Seed(time.Now().UnixNano())
	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			generateAndPrintURL()
		}
	}
}

func generateAndPrintURL() {
	u, err := generateRandomURL()
	if err != nil {
		fmt.Printf("Error generating URL: %v\n", err)
		return
	}

	fmt.Println(u)
	if len(u.String()) > maxURLLength {
		shortURL := generateShortURL()
		fmt.Printf("Short URL: %s\n", shortURL)
	}
}

func generateRandomURL() (*url.URL, error) {
	u, err := url.Parse(baseURL)
	if err != nil {
		return nil, err
	}

	pathLength := rand.Intn(20) + 5 // Random path length between 5 and 25
	path := make([]rune, pathLength)
	for i := range path {
		path[i] = urlChars[rand.Intn(len(urlChars))]
	}
	u.Path = string(path)

	q := u.Query()
	for i := 0; i < queryParamsSize; i++ {
		param := queryParams[rand.Intn(len(queryParams))]
		value := strconv.Itoa(rand.Intn(1000))
		q.Set(param, value)
	}
	u.RawQuery = q.Encode()

	return u, nil
}

func generateShortURL() string {
	shortURL := make([]rune, shortURLLength)
	for i := range shortURL {
		shortURL[i] = urlChars[rand.Intn(len(urlChars))]
	}
	return string(shortURL)
}
