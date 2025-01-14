package main

import (
	"database/sql"
	"fmt"
	"testing"

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

func BenchmarkInsertUsersWithoutBatch(b *testing.B) {
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		b.Fatal(err)
	}
	defer db.Close()

	_, err = db.Exec(`CREATE TABLE users (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		first_name TEXT NOT NULL,
		last_name TEXT NOT NULL,
		age INTEGER NOT NULL
	)`)
	if err != nil {
		b.Fatal(err)
	}

	stmt, err := db.Prepare("INSERT INTO users (first_name, last_name, age) VALUES (?, ?, ?)")
	if err != nil {
		b.Fatal(err)
	}
	defer stmt.Close()

	users := make([]User, 10000)
	for i := range users {
		users[i] = User{
			First: fmt.Sprintf("First%d", i),
			Last:  fmt.Sprintf("Last%d", i),
			Age:   i%100 + 1,
		}
	}

	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		for _, user := range users {
			_, err = stmt.Exec(user.First, user.Last, user.Age)
			if err != nil {
				b.Fatal(err)
			}
		}
	}
}

func BenchmarkInsertUsersWithBatch(b *testing.B) {
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		b.Fatal(err)
	}
	defer db.Close()

	_, err = db.Exec(`CREATE TABLE users (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		first_name TEXT NOT NULL,
		last_name TEXT NOT NULL,
		age INTEGER NOT NULL
	)`)
	if err != nil {
		b.Fatal(err)
	}

	stmt, err := db.Prepare("INSERT INTO users (first_name, last_name, age) VALUES (?, ?, ?)")
	if err != nil {
		b.Fatal(err)
	}
	defer stmt.Close()

	usersToInsert := make([]User, 0, 10000)
	for i := range usersToInsert {
		usersToInsert[i] = User{
			First: fmt.Sprintf("First%d", i),
			Last:  fmt.Sprintf("Last%d", i),
			Age:   i%100 + 1,
		}
	}

	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		for _, user := range usersToInsert {
			_, err = stmt.Exec(user.First, user.Last, user.Age)
			if err != nil {
				b.Fatal(err)
			}
		}
	}
}
