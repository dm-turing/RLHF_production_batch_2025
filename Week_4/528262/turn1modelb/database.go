// database.go
package main

import (
	"database/sql"

	_ "github.com/ncruces/go-sqlite3/driver"
	_ "github.com/ncruces/go-sqlite3/embed"
)

type ProductDatabase interface {
	GetAllProducts() ([]Product, error)
	GetProduct(id int) (*Product, error)
	// Add more database operations as needed
}

type sqliteDatabase struct {
	db *sql.DB
}

func NewSQLiteDatabase(dbFile string) (ProductDatabase, error) {
	db, err := sql.Open("sqlite3", dbFile)
	if err != nil {
		return nil, err
	}
	return &sqliteDatabase{db: db}, nil
}

// Implement the GetAllProducts and GetProduct functions
