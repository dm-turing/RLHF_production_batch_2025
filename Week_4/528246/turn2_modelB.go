
package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"math/big"
	"net/http"
	"strings"
	"time"

	_ "github.com/ncruces/go-sqlite3/driver"
	_ "github.com/ncruces/go-sqlite3/embed"
)

// User struct containing metadata
type User struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

// Password struct containing password details
type Password struct {
	UserID   int       `json:"user_id"`
	Password string    `json:"password"`
	Reason   string    `json:"reason"`
	Timestamp time.Time `json:"timestamp"`
}

// database variable to hold the SQLite database connection
var db *sql.DB

func initDatabase() {
	// Initialize the SQLite database
	database, err := sql.Open("sqlite3", "./passwords.db")
	if err != nil {
		log.Fatalf("Error opening database: %v", err)
	}
	db = database

	// Create users table if it doesn't exist
	createUsersTableSQL := `
	CREATE TABLE IF NOT EXISTS users (
		id INTEGER PRIMARY KEY,
		name TEXT NOT NULL
	)
	`
	_, err = db.Exec(createUsersTableSQL)
	if err != nil {
		log.Fatalf("Error creating users table: %v", err)
	}

	// Create passwords table if it doesn't exist
	createPasswordsTableSQL := `
	CREATE TABLE IF NOT EXISTS passwords (
		id INTEGER PRIMARY KEY,
		user_id INTEGER,
		password TEXT NOT NULL,
		reason TEXT,
		timestamp DATETIME NOT NULL,
		FOREIGN KEY (user_id) REFERENCES users(id)
	)
	`
	_, err = db.Exec(createPasswordsTableSQL)
	if err != nil {
		log.Fatalf("Error creating passwords table: %v", err)
	}
}

// Generate a random password that meets the criteria
func generatePassword() (string, error) {
  // (Same function as before)
}

func addUser(user User) error {
	stmt, err := db.Prepare("INSERT INTO users (name) VALUES (?)")
	if err != nil {
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(user.Name)
	return err
}

func getUsers() ([]User, error) {
	var users []User
	rows, err := db.Query("SELECT id, name FROM users")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var user User
		if err := rows.Scan(&user.ID, &user.Name); err != nil {
			return nil, err
		}
		users = append(users, user)
	}
	return users, nil
}

func addPassword(password Password) error {
	stmt, err := db.Prepare("INSERT INTO passwords (user_id, password, reason, timestamp) VALUES (?, ?, ?, ?)")
	if err != nil {
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(password.UserID, password.Password, password.Reason, password.Timestamp)
	return err
}

func getPasswords(userID int, sortBy string) ([]Password, error) {
	var passwords []Password
	var rows *sql.Rows
	var err error
	
	switch sortBy {
	case "timestamp":
		rows, err = db.Query("SELECT user_id, password, reason, timestamp FROM passwords WHERE user_id = ? ORDER BY timestamp", userID)
	default:
		rows, err = db.Query("SELECT user_id, password, reason, timestamp FROM passwords WHERE user_id = ?", userID)
	}