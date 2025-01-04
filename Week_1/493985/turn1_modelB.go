package main

import (
	"log"
	"net/http"
	"sync"
)

// User represents a user in the application
type User struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

// UserStore is a map-based store to manage users
type UserStore struct {
	mu    sync.RWMutex
	users map[string]*User
}

func (s *UserStore) Get(id string) *User {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.users[id]
}

func (s *UserStore) Set(user *User) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.users[user.ID] = user
}

func (s *UserStore) Delete(id string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.users, id)
}

var userStore = &UserStore{users: make(map[string]*User)}

func main() {
	http.HandleFunc("/users", usersHandler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func usersHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		handleGetUsers(w, r)
	case "POST":
		handlePostUser(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func handleGetUsers(w http.ResponseWriter, r *http.Request) {
	users := make([]*User, 0)
	for _, user := range userStore.users {
		users = append(users, user)
	}
	// Encode users to JSON and write to response
	// ...
}

func handlePostUser(w http.ResponseWriter, r *http.Request) {
	// Decode JSON request body to a User struct
	// ...

	// Validate user data
	// ...

	userStore.Set(user)
	// Encode user to JSON and write to response
	// ...
}
