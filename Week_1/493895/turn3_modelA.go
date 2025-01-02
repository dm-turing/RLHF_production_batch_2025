package main

import (
	"fmt"
	"log"
	"os"

	"github.com/jinzhu/gorm"
)

// User represents the users table
type User struct {
	gorm.Model
	Name   string  `gorm:"type:text;not null"`
	Email  string  `gorm:"type:text;unique;not null"`
	Orders []Order `gorm:"foreignkey:user_id"`
}

// Product represents the products table
type Product struct {
	gorm.Model
	Name   string  `gorm:"type:text;not null"`
	Price  float64 `gorm:"type:real;not null"`
	Stock  int     `gorm:"type:integer;not null"`
	Orders []Order `gorm:"foreignkey:product_id"`
}

// Order represents the orders table
type Order struct {
	gorm.Model
	UserID    int
	ProductID int
	Quantity  int     `gorm:"type:integer;not null"`
	User      User    `gorm:"foreignkey:user_id"`
	Product   Product `gorm:"foreignkey:product_id"`
}

// main function
func main() {
	// Connect to a SQLite database for this example
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		dsn = "file:./ecommerce.db?cache=shared&_fk=1"
	}

	db, err := gorm.Open("sqlite3", dsn)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	defer db.Close()

	// Create the necessary tables
	err = db.AutoMigrate(&User{}, &Product{}, &Order{})
	if err != nil {
		log.Fatal("Failed to create tables:", err)
	}

	// Create a new user and product
	user := User{Name: "Alice", Email: "alice@example.com"}
	product := Product{Name: "Widget", Price: 9.99, Stock: 100}

	err = db.Create(&user).Error
	if err != nil {
		log.Fatal("Error creating user:", err)
	}

	err = db.Create(&product).Error
	if err != nil {
		log.Fatal("Error creating product:", err)
	}

	// Create an order for the user and product
	order := Order{UserID: user.ID, ProductID: product.ID, Quantity: 2}
	err = db.Create(&order).Error
	if err != nil {
		log.Fatal("Error creating order:", err)
	}

	// Retrieve the user with associated orders
	var retrievedUser User
	err = db.Preload("Orders").First(&retrievedUser, "id = ?", user.ID).Error
	if err != nil {
		log.Fatal("Error retrieving user with orders:", err)
	}

	fmt.Println("Retrieved User:", retrievedUser)
	for _, o := range retrievedUser.Orders {
		fmt.Println("Order:", o)
	}
}
