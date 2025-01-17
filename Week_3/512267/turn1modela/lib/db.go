package lib

import (
	"fmt"
	"log"

	"github.com/jmoiron/sqlx"
	_ "github.com/ncruces/go-sqlite3/driver"
	_ "github.com/ncruces/go-sqlite3/embed"
)

var db *sqlx.DB

func InitDB(dbPath string) {
	var err error
	db, err = sqlx.Open("sqlite3", fmt.Sprintf(":memory:")) // Replace with the file path
	if err != nil {
		log.Fatalf("Error connecting to DB: %v", err)
	}

	// Create users table
	_, err = db.Exec(`
	CREATE TABLE IF NOT EXISTS users (
		id INTEGER PRIMARY KEY,
		name TEXT NOT NULL,
		email TEXT NOT NULL UNIQUE,
		conversion_rate REAL,
		session_duration REAL
	);
	`)
	if err != nil {
		log.Fatalf("Error creating users table: %v", err)
	}

	// Populate user data with random values
	populateUsers()
}

func CloseDB() {
	if err := db.Close(); err != nil {
		log.Fatalf("Error closing DB: %v", err)
	}
}

func GetAllUsers() ([]User, error) {
	var users []User
	rows, err := db.Queryx("SELECT * FROM users")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var user User
		err := rows.StructScan(&user)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	return users, nil
}

func GetUserById(userId string) (User, error) {
	var user User
	err := db.Getx(&user, "SELECT * FROM users WHERE id = ?", userId)
	if err != nil {
		return user, err
	}
	return user, nil
}

func AddUser(user User) error {
	_, err := db.Execx(
		"INSERT INTO users (name, email, conversion_rate, session_duration) VALUES (?, ?, ?, ?)",
		user.Name, user.Email, user.ConversionRate, user.SessionDuration,
	)
	return err
}

func populateUsers() {
	for i := 1; i <= 10; i++ { // Randomly generate 10 users
		name := fmt.Sprintf("User%d", i)
		email := fmt.Sprintf("user%d@example.com", i)
		err := AddUser(User{Name: name, Email: email})
		if err != nil {
			log.Printf("Error adding user %s: %v", name, err)
		}
	}
}
