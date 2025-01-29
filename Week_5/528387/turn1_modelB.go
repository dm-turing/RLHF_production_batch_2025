package main

import (
	"fmt"
	"math/rand"
	"net/http"
)

type middleware func(http.Handler) http.Handler

// Simulate error middleware
func errorMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Simulate error every 3 seconds
		if rand.Intn(3) == 0 {
			errors := []string{
				"Product not found in stock",
				"Payment gateway timeout",
				"Shipping address not valid",
				"Order exceeded maximum limit",
			}
			errMessage := errors[rand.Intn(len(errors))]
			http.Error(w, errMessage, http.StatusInternalServerError)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func searchProducts(w http.ResponseWriter, r *http.Request) {
	// Handle product search logic
	fmt.Fprintln(w, "Products found successfully.")
}

func placeOrder(w http.ResponseWriter, r *http.Request) {
	// Handle order placement logic
	fmt.Fprintln(w, "Order placed successfully.")
}

func shipOrder(w http.ResponseWriter, r *http.Request) {
	// Handle shipping logic
	fmt.Fprintln(w, "Order shipped successfully.")
}

func main() {
	r := http.NewServeMux()

	// Add error middleware to all routes
	r.Handle("/search", errorMiddleware(http.HandlerFunc(searchProducts)))
	r.Handle("/placeOrder", errorMiddleware(http.HandlerFunc(placeOrder)))
	r.Handle("/shipOrder", errorMiddleware(http.HandlerFunc(shipOrder)))

	// Start the server
	fmt.Println("Shopping cart application running on port 8080")
	err := http.ListenAndServe(":8080", r)
	if err != nil {
		panic(err)
	}
}
