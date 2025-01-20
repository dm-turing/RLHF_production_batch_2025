package main

import (
    "math/rand"
    "strings"
    "time"
)

// characterSet defines the characters that can be used in the URL path
const characterSet = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789-_.~"

// generateRandomURL generates a random URL with a specified path length
func generateRandomURL(minLength, maxLength int) string {
    rand.Seed(time.Now().UnixNano())
    length := rand.Intn(maxLength-minLength+1) + minLength
    var path strings.Builder
    for i := 0; i < length; i++ {
        idx := rand.Intn(len(characterSet))
        path.WriteByte(characterSet[idx])
    }
    return "http://example.com/" + path.String()
}

func main() {
    // For demonstration, generate a single random URL
    randomURL := generateRandomURL(75, 255)
    println(randomURL)
}