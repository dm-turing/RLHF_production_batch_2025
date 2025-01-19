package auth

import (
	"errors"
)

// User represents a user with a username and password.
type User struct {
	Username string
	Password string
	Email    string
}

// UserStore stores registered users.
type UserStore struct {
	users  map[string]User
	emails map[string]bool
}

// NewUserStore creates a new UserStore.
func NewUserStore() *UserStore {
	return &UserStore{users: make(map[string]User), emails: make(map[string]bool)}
}

// Register user with username, password and email.
func (store *UserStore) Register(username, password, email string, validUsernameCheck func(string) error, validPasswordCheck func(string) error, validEmailCheck func(string) error) error {
	if err := validUsernameCheck(username); err != nil {
		return err
	}

	if err := validEmailCheck(email); err != nil {
		return err
	}

	if err := validPasswordCheck(password); err != nil {
		return err
	}

	if _, exists := store.users[username]; exists {
		return errors.New("username already exists")
	}

	if _, exists := store.emails[email]; exists {
		return errors.New("email already exists")
	}

	store.users[username] = User{Username: username, Password: password, Email: email}
	store.emails[email] = true
	return nil
}

// Authenticate checks if the username and password match an existing user.
func (store *UserStore) Authenticate(username, password string) bool {
	user, exists := store.users[username]
	return exists && user.Password == password
}
