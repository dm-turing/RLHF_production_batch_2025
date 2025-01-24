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
)

type SlidingWindowLimiter struct {
	maxTokens      int
	refillRate     int
	windowDuration time.Duration
	tokens         map[string]ratelimit.Bucket
	mu             sync.Mutex
}

type Product struct {
	ID    int     `json:"id"`
	Name  string  `json:"name"`
	Price float64 `json:"price"`
}

func NewSlidingWindowLimiter(maxTokens, refillRate int, windowDuration time.Duration) *SlidingWindowLimiter {
	return &SlidingWindowLimiter{
		maxTokens:      maxTokens,
		refillRate:     refillRate,
		windowDuration: windowDuration,
		tokens:         make(map[string]ratelimit.Bucket),
	}
}

func (swl *SlidingWindowLimiter) Allow(key string) bool {
	swl.mu.Lock()
	defer swl.mu.Unlock()

	bucket, exists := swl.tokens[key]
	if !exists {
		bucket = *ratelimit.NewBucketWithRate(float64(swl.refillRate), int64(swl.maxTokens))
		swl.tokens[key] = bucket
	}

	if bucket.TakeAvailable(1) < 1 {
		return false
	}

	go func() {
		time.Sleep(swl.windowDuration)
		swl.mu.Lock()
		defer swl.mu.Unlock()
		delete(swl.tokens, key)
	}()

	return true
}

func (swl *SlidingWindowLimiter) PrintStatus() {
	swl.mu.Lock()
	defer swl.mu.Unlock()

	fmt.Printf("Number of Buckets: %d\n", len(swl.tokens))
	for key, bucket := range swl.tokens {
		fmt.Printf("Key: %s, Available Tokens: %d\n", key, bucket.Available())
	}
}

const (
	maxTokens      = 10          // Maximum number of tokens
	refillRate     = 2           // Refill rate (tokens per second)
	windowDuration = time.Minute // Window duration
)

var (
	swl = NewSlidingWindowLimiter(maxTokens, refillRate, windowDuration)
	db  *sql.DB
)

func productHandler(w http.ResponseWriter, r *http.Request) {
	key := r.RemoteAddr // This could be improved to consider IP + path for session-based limit
	if !swl.Allow(key) {
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
}

func getProducts() ([]Product, error) {
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

	go func() {
		for {
			time.Sleep(5 * time.Second)
			swl.PrintStatus()
		}
	}()

	http.HandleFunc("/products", productHandler)

	log.Println("Server starting on :8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}
