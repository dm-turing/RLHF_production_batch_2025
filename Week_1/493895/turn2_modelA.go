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
	Name  string `gorm:"type:varchar(100);not null"`
	Email string `gorm:"type:varchar(100);not null;unique"`
}

// main function
func main() {
	// Connect to a SQLite database for this example
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		dsn = "file:./users.db?cache=shared&_fk=1"
	}

	db, err := gorm.Open("sqlite3", dsn)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	defer db.Close()

	// Create the users table
	db = db.AutoMigrate(&User{})
	if db == nil {
		log.Fatal("Failed to create users table:", db)
	}

	// Query a user by ID
	var user User
	err = db.First(&user, "id = ?", 1).Error
	if err != nil {
		if gorm.IsRecordNotFoundError(err) {
			fmt.Println("User not found.")
		} else {
			log.Fatal("Error querying user:", err)
		}
	} else {
		fmt.Println(user)
	}

	// Save a new user
	newUser := &User{Name: "Charlie", Email: "charlie@example.com"}
	err = db.Create(newUser).Error
	if err != nil {
		log.Fatal("Error saving user:", err)
	}
	fmt.Println("User saved successfully:", newUser)

	// Query again to verify
	err = db.First(&user, "id = ?", newUser.ID).Error
	if err != nil {
		if gorm.IsRecordNotFoundError(err) {
			fmt.Println("User not found.")
		} else {
			log.Fatal("Error querying user:", err)
		}
	} else {
		fmt.Println(user)
	}
}
