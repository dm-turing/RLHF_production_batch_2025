package main

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
)

// Role defines a user's role
type Role string

// User represents a user
type User struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Role  Role   `json:"role"`
	Token string `json:"token"`
}

// ValidRoles for authentication
var ValidRoles = []Role{"admin", "user"}

// Middleware functions

func corsMiddleware(next http.Handler) func(w http.ResponseWriter, r *http.Request) {
	allowedOrigins := []string{"http://example.com", "https://another-example.com"}
	return func(w http.ResponseWriter, r *http.Request) {
		origin := r.Header.Get("Origin")
		if strings.Contains(strings.Join(allowedOrigins, ","), origin) {
			w.Header().Set("Access-Control-Allow-Origin", origin)
			w.Header().Set("Access-Control-Allow-Methods", "GET,POST,PUT,DELETE")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type,Authorization")
		} else {
			http.Error(w, "Request not allowed", http.StatusForbidden)
			return
		}
		next.ServeHTTP(w, r)
	}
}

func rbacMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authToken := r.Header.Get("Authorization")
		if authToken == "" {
			http.Error(w, "Authorization token required", http.StatusUnauthorized)
			return
		}

		// Validate token (for demonstration purposes, we assume a valid token)
		var user User
		err := json.Unmarshal([]byte(authToken), &user)

		if err != nil {
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		if !containsRole(ValidRoles, user.Role) {
			http.Error(w, "Insufficient privileges", http.StatusForbidden)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func containsRole(roles []Role, role Role) bool {
	for _, r := range roles {
		if r == role {
			return true
		}
	}
	return false
}

func handler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Access granted!"))
}

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/api/resource", corsMiddleware(rbacMiddleware(http.HandlerFunc(handler)))).Methods("GET")
	http.ListenAndServe(":8080", r)
}
