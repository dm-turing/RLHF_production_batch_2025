// auth.go
package auth

import (
	"errors"
	"regexp"
)

// User represents a user with a username and password.
type User struct {
	Username string
	Password string
}

// UserStore stores registered users.
type UserStore struct {
	users map[string]User
}

// NewUserStore creates a new UserStore.
func NewUserStore() *UserStore {
	return &UserStore{users: make(map[string]User)}
}

// Register user with username and password.
func (store *UserStore) Register(username, password string) error {
	if _, exists := store.users[username]; exists {
		return errors.New("username already exists")
	}

	if len(password) < 8 {
		return errors.New("password must be at least 8 characters long")
	}

	if !isValidPassword(password) {
		return errors.New("password must contain at least one number and one special character")
	}

	store.users[username] = User{Username: username, Password: password}
	return nil
}

// Authenticate checks if the username and password match an existing user.
func (store *UserStore) Authenticate(username, password string) bool {
	user, exists := store.users[username]
	return exists && user.Password == password
}

// isValidPassword checks for password strength.
func isValidPassword(password string) bool {
	hasNumber := regexp.MustCompile(`[0-9]`).MatchString(password)
	hasSpecial := regexp.MustCompile(`[!@#$%^&*(),.?":{}|<>]`).MatchString(password)
	return hasNumber && hasSpecial && len(password) >= 8
}
