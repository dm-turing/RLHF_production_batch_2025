package main

import (
	"fmt"
)

func IsValidCredentials(username string, password string) bool {
	// Simple password validation: Check if password length is at least 6
	validPassword := len(password) >= 6
	return username != "" && validPassword
}
func MeetPasswordRequirements(password string) bool {
	// Simple password requirement check: Check if the password contains at least 1 uppercase letter, 1 lowercase letter, 1 digit, and 1 special character.
	upper, lower, digit, special := false, false, false, false
	for _, char := range password {
		switch {
		case 'A' <= char && char <= 'Z':
			upper = true
		case 'a' <= char && char <= 'z':
			lower = true
		case '0' <= char && char <= '9':
			digit = true
		default:
			special = true
		}
	}
	return upper && lower && digit && special
}
func main() {
	fmt.Println("Simple Authentication System")
}
