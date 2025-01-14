package main

import (
	"database/sql"
	"fmt"
	"sort"
	"time"

	_ "github.com/ncruces/go-sqlite3/driver"
	_ "github.com/ncruces/go-sqlite3/embed"
)

// User represents a user in the database
type User struct {
	ID    int
	First string
	Last  string
	Age   int
}

const batchSize = 100 // Adjust the batch size as per your need

func main() {
	db, err := sql.Open("sqlite3", "users_db.db")
	if err != nil {
		panic(err)
	}
	defer db.Close()

	// Query to retrieve the count of users
	var count int
	err = db.QueryRow("SELECT COUNT(*) FROM users").Scan(&count)
	if err != nil {
		panic(err)
	}
	//Pre-allocate the slice with a sufficient capacity
	users := make(Users, 0, count)

	//retrieve data
	rows, err := db.Query("SELECT id, first_name, last_name, age FROM users")
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	for rows.Next() {
		var user User
		if err := rows.Scan(&user.ID, &user.First, &user.Last, &user.Age); err != nil {
			panic(err)
		}
		users = append(users, user)
	}
	if err := rows.Err(); err != nil {
		panic(err)
	}

	//Sort the slice
	sort.Stable(users)

	//Batch the inserts
	start := time.Now()
	err = batchInsertUsers(db, users)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Batch Inserted users. Time taken: %s\n", time.Since(start))
}

func batchInsertUsers(db *sql.DB, users Users) error {
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	stmt, err := tx.Prepare("INSERT INTO users (first_name, last_name, age) VALUES (?, ?, ?)")
	if err != nil {
		return err
	}
	defer stmt.Close()
	for i := 0; i < len(users); i += batchSize {
		end := min(i+batchSize, len(users))
		for _, user := range users[i:end] {
			_, err := stmt.Exec(user.First, user.Last, user.Age)
			if err != nil {
				tx.Rollback()
				return err
			}
		}
	}
	return tx.Commit()
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
