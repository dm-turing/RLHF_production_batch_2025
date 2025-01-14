package main

import (
	"database/sql"
	"fmt"

	_ "github.com/ncruces/go-sqlite3/driver"
	_ "github.com/ncruces/go-sqlite3/embed"
)

// User represents a user in the database
type User struct {
	ID    int
	First string
	Last  string
	Age   int
}

// Users is a slice of User
type Users []User

// Len returns the length of the slice
func (us Users) Len() int {
	return len(us)
}

// Less compares two Users based on multiple fields
func (us Users) Less(i, j int) bool {
	// First sort by age
	if us[i].Age < us[j].Age {
		return true
	} else if us[i].Age > us[j].Age {
		return false
	}
	// If age is the same, sort by last name
	if us[i].Last < us[j].Last {
		return true
	} else if us[i].Last > us[j].Last {
		return false
	}
	// If last name is the same, sort by first name
	if us[i].First < us[j].First {
		return true
	}
	return false
}

// Swap swaps two Users in the slice
func (us Users) Swap(i, j int) {
	us[i], us[j] = us[j], us[i]
}

func main() {
	// Connect to the SQLite3 database
	db, err := sql.Open("sqlite3", "users_db.db")
	if err != nil {
		panic(err)
	}
	defer db.Close()

	stmt, err := db.Prepare("INSERT INTO users (first_name, last_name, age) VALUES (?, ?, ?)")
	if err != nil {
		panic(err)
	}
	defer stmt.Close()

	// Pre-allocate the capacity for the slice of users to insert
	usersToInsert := make([]User, 0, 100)
	users := make([]User, 100)
	for i := range users {
		users[i] = User{
			First: fmt.Sprintf("First%d", i),
			Last:  fmt.Sprintf("Last%d", i),
			Age:   i%100 + 1,
		}
	}
	// Add users to the slice
	for _, user := range users {
		usersToInsert = append(usersToInsert, user)
	}

	// Batch insert the users
	for _, user := range usersToInsert {
		_, err := stmt.Exec(user.First, user.Last, user.Age)
		if err != nil {
			panic(err)
		}
	}
}
