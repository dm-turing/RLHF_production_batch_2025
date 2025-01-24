package main

import (
	"database/sql"
	"log"
	"time"

	_ "github.com/ncruces/go-sqlite3/driver"
	_ "github.com/ncruces/go-sqlite3/embed"
	// "github.com/thomas11/atom10" // hypothetical package for Atom parsing
)

type FeedItem struct {
	Category string
	Link     string
	Title    string
	Updated  time.Time
}

func fetchFeeds() []FeedItem {
	// This function should retrieve and parse real Atom feeds.
	// For the purpose of this example, we'll simulate data.
	return []FeedItem{
		{Category: "sports", Link: "https://sports.example.com/item1", Title: "Sports News 1", Updated: time.Now()},
		{Category: "politics", Link: "https://politics.example.com/item1", Title: "Political News 1", Updated: time.Now()},
		{Category: "media", Link: "https://media.example.com/item1", Title: "Media News 1", Updated: time.Now()},
		{Category: "music", Link: "https://music.example.com/item1", Title: "Music News 1", Updated: time.Now()},
		{Category: "business", Link: "https://business.example.com/item1", Title: "Business News 1", Updated: time.Now()},
		// More feed items can be added here
	}
}

func insertFeedsIntoDB(db *sql.DB, feeds []FeedItem) {
	tx, err := db.Begin()
	if err != nil {
		log.Fatal(err)
	}
	stmt, err := tx.Prepare("INSERT INTO feeds(category, link, title, updated) VALUES (?, ?, ?, ?)")
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()

	for _, feed := range feeds {
		_, err = stmt.Exec(feed.Category, feed.Link, feed.Title, feed.Updated)
		if err != nil {
			log.Fatal(err)
		}
	}
	tx.Commit()
}

func main() {
	db, err := sql.Open("sqlite3", "./feeds.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	for range ticker.C {
		feeds := fetchFeeds()
		insertFeedsIntoDB(db, feeds)
	}
}
