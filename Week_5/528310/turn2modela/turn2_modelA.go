package main

import (
	"database/sql"
	"fmt"
	"math/rand"
	"time"

	_ "github.com/ncruces/go-sqlite3/driver"
	_ "github.com/ncruces/go-sqlite3/embed"
)

// Struct to store URL information
type URLInfo struct {
	FullURL     string
	ShortURL    string
	LoadTime    time.Duration
	ShortenTime time.Duration
}

// Function to initialize SQLite database
func initDB() (*sql.DB, error) {
	db, err := sql.Open("sqlite3", "./urlstore.db")
	if err != nil {
		return nil, err
	}

	query := `
	CREATE TABLE IF NOT EXISTS urls (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		full_url TEXT NOT NULL,
		short_url TEXT
	);`
	_, err = db.Exec(query)
	if err != nil {
		return nil, err
	}

	return db, nil
}

// Function to store URL in the database
func storeURL(db *sql.DB, fullURL, shortURL string) error {
	query := `
	INSERT INTO urls (full_url, short_url)
	VALUES (?, ?);`
	_, err := db.Exec(query, fullURL, shortURL)
	return err
}

// Function to generate a random string
func randomString(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	var seededRand *rand.Rand = rand.New(rand.NewSource(time.Now().UnixNano()))
	result := make([]byte, length)
	for i := range result {
		result[i] = charset[seededRand.Intn(len(charset))]
	}
	return string(result)
}

// Function to generate a random URL
func generateRandomURL() string {
	baseURL := "https://www.example.com"
	path := "/" + randomString(rand.Intn(10)+5)

	params := fmt.Sprintf("id=%s&page=%d&sort=%s&filter=%s&limit=%d",
		randomString(5), rand.Intn(100)+1, randomString(6), randomString(4), rand.Intn(50)+1)

	return baseURL + path + "?" + params
}

// Function to shorten a URL if it's longer than 25 characters
func shortenURL(longURL string) string {
	if len(longURL) > 25 {
		return "https://short.url/" + randomString(7)
	}
	return longURL
}

// Function to simulate loading a URL and return a random load time
func simulateLoadTime(url string) time.Duration {
	return time.Duration(rand.Intn(300)+100) * time.Millisecond
}

func main() {
	db, err := initDB()
	if err != nil {
		fmt.Println("Error initializing database:", err)
		return
	}
	defer db.Close()

	var urlInfoList []URLInfo

	for i := 0; i < 10000; i++ {
		fullURL := generateRandomURL()
		shortURL := shortenURL(fullURL)

		// Simulate loading the URLs
		loadTime := simulateLoadTime(fullURL)
		shortenLoadTime := simulateLoadTime(shortURL)

		urlInfoList = append(urlInfoList, URLInfo{FullURL: fullURL, ShortURL: shortURL, LoadTime: loadTime, ShortenTime: shortenLoadTime})

		err := storeURL(db, fullURL, shortURL)
		if err != nil {
			fmt.Println("Error storing URL:", err)
			return
		}
	}

	// Calculating averages
	var totalLoadTime, totalShortenTime time.Duration
	for _, info := range urlInfoList {
		totalLoadTime += info.LoadTime
		totalShortenTime += info.ShortenTime
	}

	avgLoadTime := totalLoadTime / time.Duration(len(urlInfoList))
	avgShortenTime := totalShortenTime / time.Duration(len(urlInfoList))

	fmt.Printf("Averages over 10000 URLs:\n")
	fmt.Printf("Average Load Time for full URLs: %v\n", avgLoadTime)
	fmt.Printf("Average Load Time for shortened URLs: %v\n", avgShortenTime)
}
