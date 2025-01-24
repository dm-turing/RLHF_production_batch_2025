package main

import (
	"database/sql"
	"fmt"
	"log"
	"sync"

	_ "github.com/ncruces/go-sqlite3/driver"
	_ "github.com/ncruces/go-sqlite3/embed"
)

var categories = []string{"sports", "politics", "media", "music", "business"}

type Feed struct {
	Category string
	Link     string
	Title    string
	Updated  string
}

func worker(db *sql.DB, category string, wg *sync.WaitGroup) {
	defer wg.Done()

	rows, err := db.Query("SELECT category, link, title, updated FROM feeds WHERE category = ?", category)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	for rows.Next() {
		var f Feed
		err := rows.Scan(&f.Category, &f.Link, &f.Title, &f.Updated)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("[%s] %s: %s (%s)\n", category, f.Title, f.Link, f.Updated)
	}
}

func main() {
	db, err := sql.Open("sqlite3", "./feeds.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	var wg sync.WaitGroup
	for _, category := range categories {
		wg.Add(1)
		go worker(db, category, &wg)
	}

	wg.Wait()
}
