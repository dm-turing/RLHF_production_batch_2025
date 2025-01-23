package main

import (
	"math/rand"
	// "crypto/rand"
	"fmt"
	"time"
)

// User struct containing metadata
type User struct {
	ID   int
	Name string
}

// Generate a random password that meets the criteria
func generatePassword() (string, error) {
	const (
		lower          = "abcdefghijklmnopqrstuvwxyz"
		upper          = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
		digits         = "0123456789"
		special        = "!@#$%^&*()-_=+[]{}<>?,./"
		all            = lower + upper + digits + special
		passwordLength = 10
	)

	password := make([]byte, passwordLength)

	// Ensure at least one character from each required set
	sets := []string{lower, upper, digits, special}
	for i, set := range sets {
		c, err := randomCharFromSet(set)
		if err != nil {
			return "", err
		}
		password[i] = c
	}

	// Fill remaining spaces with random characters from all character sets
	for i := len(sets); i < passwordLength; i++ {
		c, err := randomCharFromSet(all)
		if err != nil {
			return "", err
		}
		password[i] = c
	}

	// Shuffle the password to ensure randomness
	shuffled, err := shuffle(password)
	if err != nil {
		return "", err
	}

	return string(shuffled), nil
}

// Generate a random character from a given set
func randomCharFromSet(set string) (byte, error) {
	index := rand.Int31n(int32(len(set)))
	return set[index], nil
}

// Shuffle the password to ensure randomness
func shuffle(password []byte) ([]byte, error) {
	perm := rand.Perm(len(password))
	shuffled := make([]byte, len(password))
	for i, v := range perm {
		shuffled[v] = password[i]
	}
	return shuffled, nil
}

func main() {
	users := []User{
		{ID: 1, Name: "Alice"},
		{ID: 2, Name: "Bob"},
		{ID: 3, Name: "Charlie"},
		{ID: 4, Name: "David"},
		{ID: 5, Name: "Eve"},
	}

	ticker := time.NewTicker(2 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case t := <-ticker.C:
			// Select a random user
			userIndex := rand.Int31n(int32(len(users)))
			user := users[userIndex]

			// Generate the password
			password, err := generatePassword()
			if err != nil {
				fmt.Println("Error generating password:", err)
				continue
			}

			// Print user details, password, and timestamp
			fmt.Printf("User: %+v, Password: %s, Timestamp: %s\n", user, password, t.Format(time.RFC3339))
		}
	}
}
