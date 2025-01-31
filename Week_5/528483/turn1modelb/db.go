// database.go
package db

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/mattn/go-sqlite3"
	"go-e-commerce-cart/model"
)

// InitDB initializes the SQLite database and creates the products table if it doesn't exist.
func InitDB() (*sql.DB, error) {
	db, err := sql.Open("sqlite3", "./cart.db")
	if err != nil {
		return nil, err
	}

	createTableSQL := `
	CREATE TABLE IF NOT EXISTS products (
		id INTEGER PRIMARY KEY,
		name TEXT,
		description TEXT,
		quantity INTEGER,
		price REAL
	);
	`
	_, err = db.Exec(createTableSQL)
	if err != nil {
		return nil, err
	}

	return db, nil
}

// InsertProduct inserts a product into the database.
func InsertProduct(db *sql.DB, product *model.Product) error {
	stmt, err := db.Prepare("INSERT INTO products(id, name, description, quantity, price) VALUES(?, ?, ?, ?, ?)")
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(product.ID, product.Name, product.Description, product.Quantity, product.Price)
	return err
}

// GetProducts retrieves all products from the database.
func GetProducts(db *sql.DB) ([]*model.Product, error) {
	rows, err := db.Query("SELECT id, name, description, quantity, price FROM products")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var products []*model.Product
	for rows.Next() {
		var p model.Product
		if err := rows.Scan(&p.ID, &p.Name, &p.Description, &p.Quantity, &p.Price); err != nil {
