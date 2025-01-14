package main

import (
	"database/sql"
	"fmt"
	"sort"

	_ "github.com/mattn/go-sqlite3"
)

// User struct to represent a user in the database
type User struct {
	ID        int
	FirstName string
	LastName  string
	Age       int
}

// UserSlice defines a slice of User structs to support sorting
type UserSlice []User

func (us UserSlice) Len() int {
	return len(us)
}

func (us UserSlice) Swap(i, j int) {
	us[i], us[j] = us[j], us[i]
}

// Less implements the stable sorting logic for users by age and then by first name
func (us UserSlice) Less(i, j int) bool {
	if us[i].Age != us[j].Age {
		return us[i].Age < us[j].Age
	}
	return us[i].FirstName < us[j].FirstName
}

func main() {
	db, err := sql.Open("sqlite3", "./users.db")
	if err != nil {
		panic(err)
	}
	defer db.Close()

	// Function to retrieve all users from the database
	getAllUsers := func() ([]User, error) {
		var users []User
		rows, err := db.Query("SELECT id, first_name, last_name, age FROM users")
		if err != nil {
			return nil, err
		}
		defer rows.Close()
		for rows.Next() {
			var user User
			if err := rows.Scan(&user.ID, &user.FirstName, &user.LastName, &user.Age); err != nil {
				return nil, err
			}
			users = append(users, user)
		}
		return users, nil
	}

	// Retrieve users from the database
	users, err := getAllUsers()
	if err != nil {
		panic(err)
	}

	// Sort the users using the defined stable sorting order
	sort.Stable(UserSlice(users))

	// Display the sorted list of users
	fmt.Println("Sorted Users:")
	for _, user := range users {
		fmt.Printf("%s %s - Age: %d\n", user.FirstName, user.LastName, user.Age)
	}
}
