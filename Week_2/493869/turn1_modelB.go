package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

// Define the roles for access control
const (
	RoleAdmin = "admin"
	RoleUser  = "user"
	RoleGuest = "guest"
)

// Define a custom HTTP handler function that checks for CORS and role-based access control
func corsAndRBACHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Check for CORS
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS, PUT, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")

		// Check for role-based access control
		role := r.URL.Query().Get("role")
		if role == "" {
			role = RoleGuest // Default role to guest if not provided
		}

		// Validate the role
		validRoles := []string{RoleAdmin, RoleUser, RoleGuest}
		if !contains(validRoles, role) {
			http.Error(w, "Invalid role", http.StatusBadRequest)
			return
		}

		// Create a context with the role value and pass it to the next handler
		ctx := context.WithValue(r.Context(), "role", role)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func contains(slice []string, element string) bool {
	for _, elem := range slice {
		if elem == element {
			return true
		}
	}
	return false
}

// Example protected endpoint handler
func protectedEndpoint(w http.ResponseWriter, r *http.Request) {
	// Retrieve the role from the context
	role, ok := r.Context().Value("role").(string)
	if !ok {
		http.Error(w, "Role not found", http.StatusInternalServerError)
		return
	}

	// Perform role-based access control check
	if role != RoleAdmin {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}

	// Handle the request if the role is admin
	fmt.Fprintf(w, "Hello, admin! You can access this protected endpoint.")
}

func main() {
	r := mux.NewRouter()

	// Register the protected endpoint with CORS and RBAC middleware
	r.Handle("/protected", corsAndRBACHandler(http.HandlerFunc(protectedEndpoint))).Methods("GET")

	log.Println("Server is running on port 8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}
