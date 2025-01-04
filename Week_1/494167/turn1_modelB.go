package main

import (
	"fmt"
	"log"
	"net/http"
)

var (
	// Key used to decrypt the token
	secretKey = "exampleSecretKey"
	// Map to store authenticated users
	authenticatedUsers = make(map[string]bool)
)

func authenticate(h http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Get the token from the URL query parameters
		token := r.URL.Query().Get("token")
		if token == "" {
			http.Error(w, "Missing token parameter", http.StatusUnauthorized)
			return
		}
		//Add your token validation logic here
		// For simplicity, we will check if the token is present in our map
		isAuthenticated, ok := authenticatedUsers[token]
		if !ok || !isAuthenticated {
			http.Error(w, "Invalid or expired token", http.StatusUnauthorized)
			return
		}

		// If the token is valid, we call the next handler
		h(w, r)
	})
}
func main() {
	fmt.Println("Server started on port 8080")

	http.HandleFunc("/protected", authenticate(protectedHandler))
	log.Fatal(http.ListenAndServe(":8080", nil))
}
func protectedHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Welcome to the protected area!")
}
