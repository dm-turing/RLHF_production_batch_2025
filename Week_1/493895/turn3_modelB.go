package main

import (
	"fmt"
	"log"

	"github.com/jinzhu/gorm"
	_ "github.com/mattn/go-sqlite3" // Use the appropriate driver for your database
)

type User struct {
	gorm.Model
	Name  string
	Email string
	Posts []Post `gorm:"foreignkey:AuthorID"` // Define the one-to-many relationship
}

type Post struct {
	gorm.Model
	Title    string
	Content  string
	AuthorID uint      `gorm:"not null"`
	User     User      `gorm:"association_autoupdate:false"` // The association_autoupdate:false flag prevents infinite recursion
	Comments []Comment `gorm:"foreignkey:PostID"`            // Define the one-to-many relationship
}

type Comment struct {
	gorm.Model
	Content  string
	PostID   uint `gorm:"not null"`
	Post     Post `gorm:"association_autoupdate:false"`
	AuthorID uint `gorm:"not null"`
	Author   User `gorm:"association_autoupdate:false"`
}

func main() {
	db, err := gorm.Open("sqlite3", "./mydatabase.db")
	if err != nil {
		log.Fatal("Error connecting to database:", err)
	}
	defer db.Close()

	// Create tables if they don't exist
	db.AutoMigrate(&User{}, &Post{}, &Comment{})

	// Let's create some dummy data
	user := User{Name: "Alice", Email: "alice@example.com"}
	db.Create(&user)

	post := Post{Title: "Hello, GORM!", Content: "This is a test post.", AuthorID: user.ID}
	db.Create(&post)

	comment1 := Comment{Content: "Great post!", PostID: post.ID, AuthorID: user.ID}
	comment2 := Comment{Content: "Thanks!", PostID: post.ID, AuthorID: user.ID}
	db.Create(&comment1)
	db.Create(&comment2)

	// Query a post along with its author and comments
	var fetchedPost Post
	db.Preload("Author").Preload("Comments").First(&fetchedPost, post.ID)

	fmt.Printf("Post: %+v\n", fetchedPost)
}
