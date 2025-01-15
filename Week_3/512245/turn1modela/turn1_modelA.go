package main

import (
	"database/sql"
	"fmt"
	"math/rand"
	"time"

	_ "github.com/ncruces/go-sqlite3/driver"
	_ "github.com/ncruces/go-sqlite3/embed"
)

// DBCP is a simple database connection pool manager
type DBCP struct {
	pool    *sql.DB
	maxOpen int
	maxIdle int
}

// NewDBCP creates a new instance of the database connection pool manager
func NewDBCP(dbPath string, maxOpen, maxIdle int) *DBCP {
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		panic(err)
	}
	db.SetMaxOpenConns(maxOpen)
	db.SetMaxIdleConns(maxIdle)

	return &DBCP{
		pool:    db,
		maxOpen: maxOpen,
		maxIdle: maxIdle,
	}
}

// Close gracefully closes all the database connections
func (cp *DBCP) Close() {
	fmt.Println("Closing database connections...")
	cp.pool.Close()
}

// ExecuteTxn executes a transaction in the database
func (cp *DBCP) ExecuteTxn(id int, dur time.Duration) {
	defer func() {
		fmt.Printf("Closing connection id %d\n", id)
	}()

	tx, err := cp.pool.Begin()
	if err != nil {
		fmt.Println("Failed to start transaction:", err)
		return
	}

	fmt.Printf("Executing transaction on connection id %d\n", id)
	for i := 0; i < 3; i++ { // Simulating random operations
		select {
		case <-time.After(dur):
			fmt.Printf("Connection id %d, Operation %d\n", id, i+1)
			_, err := tx.Exec("CREATE TABLE IF NOT EXISTS test (id INTEGER PRIMARY KEY AUTOINCREMENT)")
			if err != nil {
				fmt.Println("Error executing SQL:", err)
				return
			}
		case <-time.After(10 * time.Second): // Simulate long-running operation with timeout
			fmt.Println("Timeout for operation, rolling back...")
			err := tx.Rollback()
			if err != nil {
				fmt.Println("Error rolling back transaction:", err)
			}
			return
		}
	}

	err = tx.Commit()
	if err != nil {
		fmt.Println("Error committing transaction:", err)
		return
	}
}

func main() {
	dbPath := ":memory:" // In-memory database for simplicity
	cp := NewDBCP(dbPath, 5, 2)
	defer cp.Close()

	maxIdleDuration := 30 * time.Second
	maxOpenDuration := 60 * time.Second

	for i := 0; i < 10; i++ { // Number of goroutines to spawn
		randIdleDuration := maxIdleDuration * time.Duration(rand.Intn(10)) / 10
		randOpenDuration := maxOpenDuration * time.Duration(rand.Intn(10)) / 10

		go func(id int) {
			fmt.Printf("Creating connection id %d\n", id)
			cp.ExecuteTxn(id, randIdleDuration)
			time.Sleep(randOpenDuration)
		}(i)
	}

	time.Sleep(30 * time.Second) // Waiting for all goroutines to complete
}
