package main

import (
	"errors"
	"net/http"
)

// app.go
func handleProductsPage(w http.ResponseWriter, r *http.Request, db ProductDatabase, ch chan int) error {
	select {
	case ch <- 1:
		// The rate limit window is not exceeded, handle the request
		products, err := db.GetAllProducts()
		if err != nil {
			return err
		}
		// Render the products on the webpage
		return nil
	default:
		// Rate limit threshold exceeded, respond with an error
		return errors.New("rate limit exceeded, try again later")
	}
}
