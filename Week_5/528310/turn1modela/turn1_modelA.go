package main

import (
	"fmt"
	"math/rand"
	"net/url"
	"time"
)

// Function to generate a random string
func randomString(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	var seededRand *rand.Rand = rand.New(rand.NewSource(time.Now().UnixNano()))
	result := make([]byte, length)
	for i := range result {
		result[i] = charset[seededRand.Intn(len(charset))]
	}
	return string(result)
}

// Function to generate a random URL
func generateRandomURL() string {
	baseURL := "https://www.example.com"
	path := "/" + randomString(rand.Intn(10)+5) // Random path of length 5 to 15

	params := url.Values{}
	params.Add("id", randomString(5))
	params.Add("page", fmt.Sprintf("%d", rand.Intn(100)+1))
	params.Add("sort", randomString(6))
	params.Add("filter", randomString(4))
	params.Add("limit", fmt.Sprintf("%d", rand.Intn(50)+1))

	return baseURL + path + "?" + params.Encode()
}

// Function to shorten a URL if it's longer than 25 characters
func shortenURL(longURL string) string {
	if len(longURL) > 25 {
		return "https://short.url/" + randomString(7) // Short URL with random 7 character string
	}
	return longURL
}

func main() {
	for {
		fullURL := generateRandomURL()
		shortURL := shortenURL(fullURL)

		fmt.Println("Full URL: ", fullURL)
		if fullURL != shortURL {
			fmt.Println("Shortened URL: ", shortURL)
		}

		time.Sleep(1 * time.Second)
	}
}
