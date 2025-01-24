package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/juju/ratelimit"
	_ "github.com/ncruces/go-sqlite3/driver"
	_ "github.com/ncruces/go-sqlite3/embed"
	"github.com/rcrowley/go-metrics"
)

const (
	windowSize = 5 * time.Second  // Sliding window size
	burstSize  = 10               // Burst size for each bucket
	refillRate = float64(1) / 0.5 // Refill rate for each bucket (2 requests per second)
)

var (
	rateLimiter   *ratelimit.Bucket
	db            *sql.DB
	mutex         sync.Mutex
	requestsMeter metrics.Meter
)

type Product struct {
	ID    int     `json:"id"`
	Name  string  `json:"name"`
	Price float64 `json:"price"`
}

func init() {
	// Initialize metrics and rate limiter
	requestsMeter = metrics.NewMeter()
	rateLimiter = ratelimit.NewBucketWithRate(refillRate, burstSize)

	// Start a goroutine to adapt rate limiting based on traffic
	go adaptRateLimiter()
}

func adaptRateLimiter() {
	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	for range ticker.C {
		// Calculate the average rate of requests in the last windowSize
		averageRate := requestsMeter.Rate1() * float64(windowSize) / float64(time.Second)

		// Adjust the rate limit based on the average rate
		rateLimiter = ratelimit.NewBucketWithRate(averageRate*2, burstSize) // Double the rate for simplicity
		// rateLimiter.SetRate(ratelimit.Limit(averageRate * 2))

		// Log the current window size, available tokens, and errors
		log.Printf("Window Size: %s, Available Tokens: %d, Average Rate: %f",
			windowSize, rateLimiter.Available(), averageRate)
	}
}

func productHandler(w http.ResponseWriter, r *http.Request) {
	if rateLimiter.TakeAvailable(1) <= 0 {
		log.Println("Rate limit exceeded!")
		http.Error(w, "Too Many Requests", http.StatusTooManyRequests)
		return
	}

	products, err := getProducts()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	for _, product := range products {
		fmt.Fprintf(w, "ID: %d, Name: %s, Price: %f\n", product.ID, product.Name, product.Price)
	}

	// Mark the request as processed
	requestsMeter.Mark(1)
}

func getProducts() ([]Product, error) {
	mutex.Lock()
	defer mutex.Unlock()

	rows, err := db.Query("SELECT id, name, price FROM products")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var products []Product
	for rows.Next() {
		var product Product
		if err := rows.Scan(&product.ID, &product.Name, &product.Price); err != nil {
			return nil, err
		}
		products = append(products, product)
	}
	return products, nil
}

func main() {
	var err error
	db, err = sql.Open("sqlite3", "./db/products.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	http.HandleFunc("/products", productHandler)
	log.Println("Server starting on :8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}
