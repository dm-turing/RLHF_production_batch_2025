package main

import (
	"database/sql"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"time"

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

var db *sql.DB

func init() {
	var err error
	db, err = sql.Open("sqlite3", "./sitemap.db")
	if err != nil {
		log.Fatal(err)
	}

	createTableSQL := `
	CREATE TABLE IF NOT EXISTS sitemap (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		loc TEXT NOT NULL,
		lastmod TEXT NOT NULL,
		changefreq TEXT NOT NULL,
		priority REAL NOT NULL,
		image TEXT,
		video TEXT
	)
	`
	_, err = db.Exec(createTableSQL)
	if err != nil {
		log.Fatal(err)
	}

	rand.Seed(time.Now().UnixNano())
}

func main() {
	defer db.Close()

	http.HandleFunc("/generate", generateSitemapHandler)
	http.HandleFunc("/sitemap", sitemapHandler)

	fmt.Println("Server is running on port 8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func generateSitemapHandler(w http.ResponseWriter, r *http.Request) {
	baseURL := "https://example.com"
	sitemap := URLSet{
		Urls: generateRandomURLs(baseURL, 100),
	}

	insertSitemap(sitemap)

	w.WriteHeader(http.StatusCreated)
	fmt.Fprintf(w, "Sitemap generated and stored in database.")
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

func insertSitemap(sitemap URLSet) {
	tx, err := db.Begin()
	if err != nil {
		log.Fatal(err)
	}
	stmt, err := tx.Prepare("INSERT INTO sitemap (loc, lastmod, changefreq, priority, image, video) VALUES (?, ?, ?, ?, ?, ?)")
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()

	for _, url := range sitemap.Urls {
		_, err := stmt.Exec(url.Loc, url.LastMod, url.ChangeFreq, url.Priority, url.Image, url.Video)
		if err != nil {
			log.Fatal(err)
		}
	}

	if err := tx.Commit(); err != nil {
		log.Fatal(err)
	}
}

func sitemapHandler(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()

	var sortBy string
	var limit int
	var offset int

	sortBy = query.Get("sort")
	limit, _ = strconv.Atoi(query.Get("limit"))
	offset, _ = strconv.Atoi(query.Get("offset"))

	if limit == 0 {
		limit = 20
	}

	if offset < 0 {
		offset = 0
	}

	rows, err := db.Query(`
	SELECT loc, lastmod, changefreq, priority, image, video
	FROM sitemap
	ORDER BY ?
	LIMIT ?
	OFFSET ?
	`, sortBy, limit, offset)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	var urls []URL
	for rows.Next() {
		var url URL
		if err := rows.Scan(&url.Loc, &url.LastMod, &url.ChangeFreq, &url.Priority, &url.Image, &url.Video); err != nil {
			log.Fatal(err)
		}
		urls = append(urls, url)
	}
	if err := rows.Err(); err != nil {
		log.Fatal(err)
	}

	sitemap := URLSet{
		Urls: urls,
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(sitemap); err != nil {
		log.Fatal(err)
	}
}

func randomDate() string {
	latestTime := time.Now().Unix()
	oldestTime := latestTime - 5*365*24*60*60
	randomTime := rand.Int63n(latestTime-oldestTime) + oldestTime
	return time.Unix(randomTime, 0).Format("2006-01-02")
}
