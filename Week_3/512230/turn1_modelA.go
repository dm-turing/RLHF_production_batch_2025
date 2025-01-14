package main

import (
	"database/sql"
	"fmt"
	"sort"

	_ "github.com/mattn/go-sqlite3"
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

	// Query to retrieve all users
	rows, err := db.Query("SELECT id, first_name, last_name, age FROM users")
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	// Slice to hold user data
	var users Users

	// Scan rows into the users slice
	for rows.Next() {
		var user User
		if err := rows.Scan(&user.ID, &user.First, &user.Last, &user.Age); err != nil {
			panic(err)
		}
		users = append(users, user)
	}
	if err := rows.Err(); err != nil {
		panic(err)
	}

	// Sort users
	sort.Stable(users)

	// Print the sorted users
	fmt.Println("Sorted Users:")
	for _, user := range users {
		fmt.Printf("ID: %d, First Name: %s, Last Name: %s, Age: %d\n", user.ID, user.First, user.Last, user.Age)
	}
}
