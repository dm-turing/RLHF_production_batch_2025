package main

import (
	"fmt"
	"log"
	"os"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

// User represents the user table
type User struct {
	gorm.Model
	Name  string `gorm:"not null"`
	Email string `gorm:"unique;not null"`
}

// Initialize the database
func initDB() *gorm.DB {
	dbFile := "example.db"
	db, err := gorm.Open("sqlite3", dbFile)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	db.AutoMigrate(&User{}) // Create the user table if it doesn't exist
	return db
}

// QueryUserByID performs a query to fetch a user by ID
func QueryUserByID(db *gorm.DB, id uint) (*User, error) {
	var user User
	result := db.First(&user, id)
	if result.Error != nil {
		return nil, result.Error
	}
	return &user, nil
}

// SaveUser saves or updates a user in the database
func SaveUser(db *gorm.DB, user *User) error {
	return db.Save(user).Error
}

// main function
func main() {
	db := initDB()
	defer db.Close()

	// Query a user by ID
	user, err := QueryUserByID(db, 1)
	if err != nil {
		fmt.Println("Error querying user:", err)
	} else {
		fmt.Println("Found user:", user)
	}

	// Save a new user
	newUser := &User{Name: "Charlie", Email: "charlie@example.com"}
	err = SaveUser(db, newUser)
	if err != nil {
		fmt.Println("Error saving user:", err)
	} else {
		fmt.Println("User saved successfully:", newUser)
	}

	// Query again to verify the new user
	user, err = QueryUserByID(db, newUser.ID)
	if err != nil {
		fmt.Println("Error querying user:", err)
	} else {
		fmt.Println("Found user:", user)
	}

	// Clean up the database file