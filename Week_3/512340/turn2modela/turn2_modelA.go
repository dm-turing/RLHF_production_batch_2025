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
	users               map[string]User
	validateCredentials func(username, password string) error
	validatePassword    func(password string) error
}

// NewUserStore creates a new UserStore with optional validation callbacks.
func NewUserStore(validateCredentials func(username, password string) error, validatePassword func(password string) error) *UserStore {
	return &UserStore{
		users:               make(map[string]User),
		validateCredentials: validateCredentials,
		validatePassword:    validatePassword,
	}
}

// Register user with username and password.
func (store *UserStore) Register(username, password string) error {
	if store.validateCredentials != nil {
		if err := store.validateCredentials(username, password); err != nil {
			return err
		}
	}

	if _, exists := store.users[username]; exists {
		return errors.New("username already exists")
	}

	if store.validatePassword != nil {
		if err := store.validatePassword(password); err != nil {
			return err
		}
	}

	store.users[username] = User{Username: username, Password: password}
	return nil
}

// Authenticate checks if the username and password match an existing user.
func (store *UserStore) Authenticate(username, password string) bool {
	user, exists := store.users[username]
	return exists && user.Password == password
}

// Default password validation function
func defaultPasswordValidation(password string) error {
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

// Default credentials validation function
func defaultCredentialsValidation(username, _ string) error {
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	if !emailRegex.MatchString(username) {
		return errors.New("username must be a valid email address")
	}
	return nil
}
