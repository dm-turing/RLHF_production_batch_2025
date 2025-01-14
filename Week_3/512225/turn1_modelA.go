package main

import (
	"database/sql"
	"encoding/base64"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	_ "github.com/ncruces/go-sqlite3/driver"
	_ "github.com/ncruces/go-sqlite3/embed"
)

var (
	db *sql.DB
)

func init() {
	dbName := "url_shortener.db"

	// Connect to the SQLite database
	conn, err := sql.Open("sqlite3", dbName)
	if err != nil {
		log.Fatal(err)
	}
	db = conn

	// Create the URLs table if it doesn't exist
	_, err = db.Exec("CREATE TABLE IF NOT EXISTS URLs (id INTEGER PRIMARY KEY, original_url TEXT, short_url TEXT)")
	if err != nil {
		log.Fatal(err)
	}
}

func shortURL(originalURL string) string {
	// Use a simple hashing algorithm and base64 encode it for simplicity
	// For production, use a strong cryptographic hash and handle collisions
	hash := base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%d:%s", os.Getpid(), originalURL)))
	return strings.ReplaceAll(hash, "/", "-") // Replace '/' with '-' to make it URL-safe
}

func handleShorten(w http.ResponseWriter, r *http.Request) {
	originalURL := r.URL.Query().Get("url")
	if originalURL == "" {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintln(w, "URL parameter is required")
		return
	}

	shortURL := shortURL(originalURL)

	_, err := db.Exec("INSERT INTO URLs (original_url, short_url) VALUES (?, ?)", originalURL, shortURL)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintln(w, err.Error())
		return
	}

	w.WriteHeader(http.StatusCreated)
	fmt.Fprintf(w, "Shortened URL: http://localhost:8080/%s\n", shortURL)
}

func handleRetrieve(w http.ResponseWriter, r *http.Request) {
	shortURL := r.URL.Path[1:]

	row := db.QueryRow("SELECT original_url FROM URLs WHERE short_url = ?", shortURL)
	var originalURL string
	if err := row.Scan(&originalURL); err != nil {
		if err == sql.ErrNoRows {
			w.WriteHeader(http.StatusNotFound)
			fmt.Fprintln(w, "Short URL not found")
		} else {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintln(w, err.Error())
		}
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Original URL: %s\n", originalURL)
}

func main() {
	http.HandleFunc("/shorten", handleShorten)
	http.HandleFunc("/", handleRetrieve)

	log.Println("URL Shortener server running on http://localhost:8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}
