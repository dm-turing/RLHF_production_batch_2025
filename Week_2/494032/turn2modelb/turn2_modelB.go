package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
)

// ... (same code as before)
type FilterRule struct {
	Keyword string `json:"keyword"`
	Replace string `json:"replace"`
	Blocked bool   `json:"blocked"`
}

// User represents a simple user entity
type User struct {
	ID   int
	Role string // User role can be "admin", "user", etc.
}

// AuthenticationHandler is a callback to handle user authentication.
type AuthenticationHandler func(http.ResponseWriter, *http.Request) (*User, error)

// AuthorizationHandler is a callback to handle user authorization based on the user role.
type AuthorizationHandler func(http.ResponseWriter, *http.Request, *User) bool

// FilterService contains all filtering-related data and methods
type FilterService struct {
	rules          []FilterRule
	authentication AuthenticationHandler
	authorization  AuthorizationHandler
}

// NewFilterServiceWithHandlers creates a FilterService with authentication and authorization callbacks.
func NewFilterServiceWithHandlers(rules []FilterRule, auth AuthenticationHandler, authz AuthorizationHandler) *FilterService {
	return &FilterService{rules: rules, authentication: auth, authorization: authz}
}

// WithAuthentication sets the authentication callback for the FilterService
func (s *FilterService) WithAuthentication(auth AuthenticationHandler) *FilterService {
	s.authentication = auth
	return s
}

// WithAuthorization sets the authorization callback for the FilterService.
func (s *FilterService) WithAuthorization(authz AuthorizationHandler) *FilterService {
	s.authorization = authz
	return s
}

func (s *FilterService) FilterContent(content string, user *User, callback func(string, error)) {
	filteredContent := content

	// First, check if the user has access to perform content filtering based on their role
	if user != nil && s.authorization != nil && !s.authorization(nil, nil, user) {
		callback("", fmt.Errorf("User '%d' is not authorized to perform content filtering", user.ID))
		return
	}

	for _, rule := range s.rules {
		if strings.Contains(filteredContent, rule.Keyword) {
			if rule.Blocked {
				filteredContent = ""
				callback(filteredContent, fmt.Errorf("Content contains blocked keyword '%s'", rule.Keyword))
				return
			} else {
				filteredContent = strings.Replace(filteredContent, rule.Keyword, rule.Replace, -1)
			}
		}
	}
	callback(filteredContent, nil)
}

// Example implementation of AuthenticationHandler for a basic authorization header
func BasicAuthHandler(w http.ResponseWriter, r *http.Request) (*User, error) {
	_, _, ok := r.BasicAuth()
	if !ok {
		w.Header().Set("WWW-Authenticate", `Basic realm="Restricted"`)
		http.Error(w, "Not authorized", http.StatusUnauthorized)
		return nil, fmt.Errorf("authentication failed")
	}

	// Validate username/password against your user store here
	user := &User{ID: 1, Role: "user"} // For demonstration, we'll set a simple role
	return user, nil
}

func main() {
	// ... (same code as before)

	rules := []FilterRule{{Keyword: "badword", Replace: "******", Blocked: false},
		{Keyword: "evil", Blocked: true}}

	filterService := NewFilterServiceWithHandlers(rules, nil, nil)
	filterService = filterService.WithAuthentication(BasicAuthHandler)

	// Define your user-role-based authorization logic
	filterService = filterService.WithAuthorization(func(w http.ResponseWriter, r *http.Request, user *User) bool {
		if user == nil {
			return false // No user, deny access
		}
		return user.Role == "admin" // Only "admin" users can filter content
	})

	// Hypothetical API gateway handler
	http.HandleFunc("/filter", func(w http.ResponseWriter, r *http.Request) {
		type Request struct {
			Content string `json:"content"`
			User    string `json:"user"`
		}

		var req Request
		err := json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		parts := strings.Split(req.User, ":")
=		id, _ := strconv.Atoi(parts[0])
		user := &User{ID: id, Role: parts[1]}
		filterService.FilterContent(req.Content, user, func(filtered string, err error) {
			if err != nil {
				http.Error(w, err.Error(), http.StatusForbidden)
			} else {
				w.WriteHeader(http.StatusOK)
				w.Write([]byte(filtered))
			}
		})
	})

	fmt.Println("Filter service listening on :8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		fmt.Println("Error starting server:", err)
	}
}
