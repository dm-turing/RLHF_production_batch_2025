package auth

import (
	"regexp"
	"testing"
)

func TestRegisterAndAuthenticate(t *testing.T) {
	store := NewUserStore()

	// Testing callback functions for valid credentials
	validUsernameCheck := func(username string) error {
		if len(username) < 3 {
			return errors.New("username must be at least 3 characters long")
		}
		return nil
	}

	validPasswordCheck := func(password string) error {
		if len(password) < 8 {
			return errors.New("password must be at least 8 characters long")
		}
		hasNumber := regexp.MustCompile(`[0-9]`).MatchString(password)
		hasSpecial := regexp.MustCompile(`[!@#$%^&*(),.?":{}|<>]`).MatchString(password)
		if !hasNumber || !hasSpecial {
			return errors.New("password must contain at least one number and one special character")
		}
		return nil
	}
	
	validEmailCheck := func(email string) error {
		emailRegex := regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")
		if !emailRegex.MatchString(email) {
			return errors.New("invalid email format")
		}
		return nil
	}

	// Test Registration Success
	err := store.Register("testuser", "StrongPass1!", "testuser@example.com", validUsernameCheck, validPasswordCheck, validEmailCheck)
	if err != nil {
		t.Fatalf("register failed: %s", err)
	}

	// Test Authentication Success
	if !store.Authenticate("testuser", "StrongPass1!") {