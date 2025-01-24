// product.go
package main

import (
	"log"
	"net/http"
	"time"
)

type Product struct {
	ID    int     `json:"id"`
	Name  string  `json:"name"`
	Price float64 `json:"price"`
}

const (
	// Rate limit constants
	MaxRequests     = 100
	RateLimitPeriod = 1 * time.Second
)

func main() {
	rateLimitCh := make(chan int, MaxRequests)
	go manageRateLimit(rateLimitCh)

	db, err := NewSQLiteDatabase("example.db") // Replace with your actual SQLite file
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	http.HandleFunc("/products", func(w http.ResponseWriter, r *http.Request) {
		if err := handleProductsPage(w, r, db, rateLimitCh); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	})

	log.Fatal(http.ListenAndServe(":8080", nil))
}
func manageRateLimit(ch chan int) {
	ticker := time.NewTicker(RateLimitPeriod)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			// Refresh the rate limit window every period
		case <-ch:
			// Decrement request count
		}
		// TODO: Implement threshold logic and alert message to console
	}
}
