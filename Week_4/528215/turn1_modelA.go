package main

import (
	"database/sql"
	"fmt"
	"log"
	"math/rand"
	"sync"
	"time"

	_ "github.com/ncruces/go-sqlite3/driver"
	_ "github.com/ncruces/go-sqlite3/embed"
)

type UpgradeInfo struct {
	ID        int
	Info      string
	Timestamp time.Time
}

func main() {
	db, err := sql.Open("sqlite3", "./upgradeinfo.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	createTable(db)

	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		defer wg.Done()
		insertDataRoutine(db)
	}()

	go func() {
		defer wg.Done()
		processDataRoutine(db)
	}()

	wg.Wait()
}

func createTable(db *sql.DB) {
	query := `
	CREATE TABLE IF NOT EXISTS upgradeInfo (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		info TEXT,
		timestamp DATETIME
	);
	`
	_, err := db.Exec(query)
	if err != nil {
		log.Fatalf("Create table failed: %s", err)
	}
}

func insertDataRoutine(db *sql.DB) {
	for {
		info := generateRandomUpgradeInfo()
		timestamp := time.Now().Add(time.Minute)

		_, err := db.Exec("INSERT INTO upgradeInfo(info, timestamp) VALUES(?, ?)", info, timestamp)
		if err != nil {
			log.Printf("Failed to insert data: %s", err)
		}
		time.Sleep(30 * time.Second)
	}
}

func generateRandomUpgradeInfo() string {
	return fmt.Sprintf("Upgrade-%d", rand.Intn(1000))
}

func processDataRoutine(db *sql.DB) {
	for {
		var ui UpgradeInfo
		err := db.QueryRow("SELECT id, info, timestamp FROM upgradeInfo ORDER BY timestamp ASC LIMIT 1").Scan(&ui.ID, &ui.Info, &ui.Timestamp)
		if err != nil {
			if err != sql.ErrNoRows {
				log.Printf("Query failed: %s", err)
			}
			time.Sleep(1 * time.Second)
			continue
		}

		if time.Now().Round(time.Second) == ui.Timestamp.Round(time.Second) {
			fmt.Printf("Processing upgrade: %s at %s\n", ui.Info, ui.Timestamp)
			_, err = db.Exec("DELETE FROM upgradeInfo WHERE id = ?", ui.ID)
			if err != nil {
				log.Printf("Failed to delete record: %s", err)
			}
		}
		time.Sleep(1 * time.Second)
	}
}
