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

const (
	maxRequests = 10 // Maximum number of requests allowed
	refillRate  = 1  // Number of tokens to refill per second
	burstSize   = 5  // Burst size
)

var (
	rateLimiter = ratelimit.NewBucketWithRate(refillRate, burstSize)
	db          *sql.DB
	mutex       sync.Mutex
)

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

	// Ensure the rate limiter is refilled constantly
	go func() {
		for {
			rateLimiter.Put(1)
			time.Sleep(time.Second / 2)
		}
	}()

	http.HandleFunc("/products", productHandler)
	log.Println("Server starting on :8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}
