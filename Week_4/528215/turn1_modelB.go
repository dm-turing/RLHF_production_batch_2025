```go
package main  

import (
	"database/sql"
	"fmt"
	"log"
	"math/rand"
	"sync"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

// upgradeInfo struct represents the data to be stored
type upgradeInfo struct {
	timestamp time.Time
	upgrade   string
}

// Ensure the upgradeInfo struct implements the sql.Scanner interface
// This allows us to read data from the database into our struct easily.
func (ui *upgradeInfo) Scan(value interface{}) error {
	str, ok := value.(string)
	if !ok {
		return fmt.Errorf("expected string, got %T", value)
	}

	// Parse the timestamp from the string stored in the database
	t, err := time.Parse(time.RFC3339, str)
	if err != nil {
		return fmt.Errorf("error parsing timestamp: %w", err)
	}
	ui.timestamp = t
	return nil
}

// Define the channel to pass data from the inserter goroutine to the main loop
var insertCh = make(chan upgradeInfo)

// func for Inserting data in background.
func insertUpgrades(db *sql.DB, wg *sync.WaitGroup) {
	defer wg.Done()
	// Create a prepared statement for inserting data
	stmt, err := db.Prepare("INSERT INTO upgrades(timestamp, upgrade) VALUES (?, ?)")
	if err != nil {
		log.Fatalf("Error preparing insert statement: %v", err)
	}
	defer stmt.Close()

	for {
		// Generate random upgrade information
		upgrade := fmt.Sprintf("Upgrade %d", rand.Intn(100))
		// Add one minute to the current time
		t := time.Now().Add(time.Minute)

		// Send the upgrade information to the main loop through the channel
		insertCh <- upgradeInfo{timestamp: t, upgrade: upgrade}

		// Execute the prepared statement to insert the data
		_, err = stmt.Exec(t.Format(time.RFC3339), upgrade)
		if err != nil {
			log.Fatalf("Error inserting data: %v", err)
		}

		// Insert new data every 30 seconds
		time.Sleep(30 * time.Second)
	}
}  

func main() {  
    //Open the database
	db, err := sql.Open("sqlite3", "./upgrades.db")
	if err != nil {
		log.Fatalf("Error opening database: %v", err)
	}
	defer db.Close()

	// Create the table if it doesn't exist
	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS upgrades (timestamp DATETIME, upgrade TEXT)`)
	if err != nil {
		log.Fatalf("Error creating table: %v", err)
	}
  
    //Start the Inserter goroutine
	var wg sync.WaitGroup
	wg.Add(1)
	go insertUpgrades(db, &wg)

	// Print and delete upgrade data when the oldest timestamp arrives
	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	for range ticker.C {
		// Query for the oldest timestamp
		var oldest upgradeInfo
		row := db.QueryRow("SELECT timestamp FROM upgrades ORDER BY timestamp ASC LIMIT 1")
		if err := row.Scan(&oldest); err != nil {