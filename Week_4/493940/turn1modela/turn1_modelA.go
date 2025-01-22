package turn1modela

import (
	"database/sql"
	"fmt"
	"log"
	"math/rand"
	"testing"
	"time"

	// _ "github.com/mattn/go-sqlite3"
	_ "github.com/ncruces/go-sqlite3/driver"
	_ "github.com/ncruces/go-sqlite3/embed"
)

type User struct {
	ID    int
	Name  string
	Email string
}

// Setup a backend instance with SQLite.
func setupDatabase() (*sql.DB, error) {
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		return nil, err
	}

	createTableSQL := `
	CREATE TABLE users (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT NOT NULL,
		email TEXT NOT NULL UNIQUE
	);`
	_, err = db.Exec(createTableSQL)
	if err != nil {
		return nil, err
	}

	return db, nil
}

// Populate the database with 100 random user profiles.
func populateDatabase(db *sql.DB, n int) ([]User, error) {
	users := make([]User, n)
	for i := 0; i < n; i++ {
		user := User{
			Name:  fmt.Sprintf("User%d", i+1),
			Email: fmt.Sprintf("user%d@example.com", i+1),
		}
		users[i] = user

		_, err := db.Exec("INSERT INTO users (name, email) VALUES (?, ?)", user.Name, user.Email)
		if err != nil {
			return nil, err
		}
	}
	return users, nil
}

// Retrieve the list of users from the database.
func getUsersFromDatabase(db *sql.DB) ([]User, error) {
	rows, err := db.Query("SELECT id, name, email FROM users")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	users := []User{}
	for rows.Next() {
		var user User
		err := rows.Scan(&user.ID, &user.Name, &user.Email)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	return users, nil
}

// Teardown the backend instance.
func teardownDatabase(db *sql.DB) error {
	return db.Close()
}

func main() {
	// Seed the random number generator
	rand.Seed(time.Now().UnixNano())

	// Run tests
	var t testing.T
	TestUserDatabase(&t)
	if t.Failed() {
		log.Println("Test failed")
	} else {
		log.Println("Test passed")
	}
}
