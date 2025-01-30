package main

import (
	"database/sql"
	"encoding/xml"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	_ "github.com/ncruces/go-sqlite3/driver"
	_ "github.com/ncruces/go-sqlite3/embed"
)

type URL struct {
	Loc        string  `xml:"loc"`
	LastMod    string  `xml:"lastmod"`
	ChangeFreq string  `xml:"changefreq"`
	Priority   float64 `xml:"priority"`
	Image      string  `xml:"image,omitempty"`
	Video      string  `xml:"video,omitempty"`
}

type URLSet struct {
	XMLName xml.Name `xml:"urlset"`
	Urls    []URL    `xml:"url"`
}

var (
	db *sql.DB
)

func initDB() error {
	var err error
	db, err = sql.Open("sqlite3", "./sitemap.db")
	if err != nil {
		return err
	}

	createTableSQL := `
	CREATE TABLE IF NOT EXISTS sitemap (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		loc TEXT NOT NULL,
		lastmod TEXT,
		changefreq TEXT,
		priority REAL,
		image TEXT,
		video TEXT
	);`

	_, err = db.Exec(createTableSQL)
	if err != nil {
		return err
	}

	return nil
}

func insertSitemapToDB(urls []URL) error {
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	stmt, err := tx.Prepare("INSERT INTO sitemap(loc, lastmod, changefreq, priority, image, video) VALUES(?, ?, ?, ?, ?, ?)")
	if err != nil {
		return err
	}
	defer stmt.Close()

	for _, url := range urls {
		_, err = stmt.Exec(url.Loc, url.LastMod, url.ChangeFreq, url.Priority, url.Image, url.Video)
		if err != nil {
			return err
		}
	}

	return tx.Commit()
}

func generateRandomURLs(baseURL string, count int) []URL {
	urls := make([]URL, count)
	changeFreqOptions := []string{"always", "hourly", "daily", "weekly", "monthly", "yearly", "never"}

	for i := 0; i < count; i++ {
		urls[i] = URL{
			Loc:        fmt.Sprintf("%s/page%d", baseURL, i),
			LastMod:    randomDate(),
			ChangeFreq: changeFreqOptions[rand.Intn(len(changeFreqOptions))],
			Priority:   rand.Float64(),
			Image:      fmt.Sprintf("%s/image%d.jpg", baseURL, i),
			Video:      fmt.Sprintf("%s/video%d.mp4", baseURL, i),
		}
	}

	return urls
}

func randomDate() string {
	latestTime := time.Now().Unix()
	oldestTime := latestTime - 5*365*24*60*60
	randomTime := rand.Int63n(latestTime-oldestTime) + oldestTime
	return time.Unix(randomTime, 0).Format("2006-01-02")
}

func querySitemap(w http.ResponseWriter, r *http.Request) {
	params := r.URL.Query()
	changeFreq := params.Get("changefreq")
	sortBy := params.Get("sort")
	limit := params.Get("limit")

	query := "SELECT loc, lastmod, changefreq, priority, image, video FROM sitemap WHERE 1=1"

	if changeFreq != "" {
		query += " AND changefreq = ?"
	}

	if sortBy != "" {
		query += fmt.Sprintf(" ORDER BY %s", sortBy)
	}

	if limit != "" {
		query += fmt.Sprintf(" LIMIT %s", limit)
	}

	rows, err := db.Query(query, changeFreq)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var urls []URL
	for rows.Next() {
		var url URL
		if err := rows.Scan(&url.Loc, &url.LastMod, &url.ChangeFreq, &url.Priority, &url.Image, &url.Video); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		urls = append(urls, url)
	}

	w.Header().Set("Content-Type", "application/json")
	xml.NewEncoder(w).Encode(URLSet{Urls: urls})
}

func main() {
	rand.Seed(time.Now().UnixNano())

	err := initDB()
	if err != nil {
		log.Fatal("Could not initialize database: ", err)
	}

	baseURL := "https://example.com"
	sitemap := generateRandomURLs(baseURL, 10)

	err = insertSitemapToDB(sitemap)
	if err != nil {
		log.Fatal("Could not insert sitemap to database: ", err)
	}

	r := mux.NewRouter()
	r.HandleFunc("/sitemap", querySitemap).Methods("GET")

	fmt.Println("Server started at :8000")
	log.Fatal(http.ListenAndServe(":8000", r))
}
