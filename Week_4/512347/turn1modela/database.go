// database.go
package main

import (
	"database/sql"
	"log"

	_ "github.com/ncruces/go-sqlite3/driver"
	_ "github.com/ncruces/go-sqlite3/embed"
)

func initDB() *sql.DB {
	db, err := sql.Open("sqlite3", "filebackup.db")
	if err != nil {
		log.Fatal(err)
	}
	createTables(db)
	return db
}

func createTables(db *sql.DB) {
	query := `
	CREATE TABLE IF NOT EXISTS files (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		path TEXT UNIQUE,
		hash TEXT,
		content BLOB
	);

	CREATE TABLE IF NOT EXISTS metadata (
		key TEXT PRIMARY KEY,
		value TEXT
	);`

	if _, err := db.Exec(query); err != nil {
		log.Fatalf("Error creating tables: %v\n", err)
	}
}
