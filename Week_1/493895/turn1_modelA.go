package main

import (
	"fmt"
)

// User represents the user table
type User struct {
	ID    int
	Name  string
	Email string
}

// QueryUserByID performs a simple query by ID using a in-memory map
func QueryUserByID(id int) *User {
	users := map[int]User{
		1: {ID: 1, Name: "Alice", Email: "alice@example.com"},
		2: {ID: 2, Name: "Bob", Email: "bob@example.com"},
	}
	user, ok := users[id]
	if ok {
		return &user
	}
	return nil
}

// SaveUser updates or inserts a user in-memory
func SaveUser(user *User) error {
	users := map[int]User{
		1: {ID: 1, Name: "Alice", Email: "alice@example.com"},
		2: {ID: 2, Name: "Bob", Email: "bob@example.com"},
	}

	if user.ID != 0 {
		users[user.ID] = *user
	} else {
		var maxID int
		for id := range users {
			if id > maxID {
				maxID = id
			}
		}
		user.ID = maxID + 1
		users[user.ID] = *user
	}

	return nil
}

// main function
func main() {
	// Query a user by ID
	user := QueryUserByID(1)
	if user != nil {
		fmt.Println(user)
	} else {
		fmt.Println("User not found.")
	}

	// Save a new user
	newUser := &User{ID: 3, Name: "Charlie", Email: "charlie@example.com"}
	err := SaveUser(newUser)
	if err != nil {
		fmt.Println("Error saving user:", err)
	} else {
		fmt.Println("User saved successfully:", newUser)
	}

	// Query again to verify
	user = QueryUserByID(3)
	if user != nil {
		fmt.Println(user)
	} else {
		fmt.Println("User not found.")
	}
}
