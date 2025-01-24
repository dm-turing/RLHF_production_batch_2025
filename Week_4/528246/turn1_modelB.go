package main

import (
	"fmt"
	"math/rand"
	"strings"
	"sync"
	"time"
)

var wg sync.WaitGroup

type user struct {
	Name  string
	Email string
}

func generatePassword() string {
	// Define the character sets for each required character class
	lowercaseChars := "abcdefghijklmnopqrstuvwxyz"
	uppercaseChars := "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	digitChars := "0123456789"
	specialChars := "!@#$%^&*()_+-=[]{}|;:,.<>?"

	allChars := lowercaseChars + uppercaseChars + digitChars + specialChars

	var password strings.Builder
	password.Grow(10) // Set the capacity of the builder to avoid reallocations

	// Generate a random password of length 10, ensuring each class is represented
	for i := 0; i < 4; i++ {
		password.WriteByte(lowercaseChars[rand.Intn(len(lowercaseChars))])
		password.WriteByte(uppercaseChars[rand.Intn(len(uppercaseChars))])
		password.WriteByte(digitChars[rand.Intn(len(digitChars))])
	}
	password.WriteByte(specialChars[rand.Intn(len(specialChars))])

	remaining := 10 - password.Len()
	for i := 0; i < remaining; i++ {
		password.WriteByte(allChars[rand.Intn(len(allChars))])
	}
	r := []rune(password.String())
	// Shuffle the password to randomize the character order
	for i, _ := range password.String() {
		j := rand.Intn(i + 1)
		r[i], r[j] = r[j], r[i]
	}
	return string(r)
}

func generateAndPrintPasswords(users []user) {
	for _, u := range users {
		wg.Add(1)
		go func(u user) {
			defer wg.Done()
			for {
				password := generatePassword()
				t := time.Now()
				fmt.Printf("User: %s, Email: %s, Password: %s, Timestamp: %s\n", u.Name, u.Email, password, t.Format("2006-01-02 15:04:05"))
				time.Sleep(2 * time.Second)
			}
		}(u)
	}
}

func main() {
	// Initialize the random number generator
	rand.Seed(time.Now().UnixNano())

	// Generate 5 users with metadata
	users := []user{
		{Name: "Alice", Email: "alice@example.com"},
		{Name: "Bob", Email: "bob@example.com"},
		{Name: "Charlie", Email: "charlie@example.com"},
		{Name: "David", Email: "david@example.com"},
		{Name: "Emma", Email: "emma@example.com"},
	}

	generateAndPrintPasswords(users)
	wg.Wait()
}
