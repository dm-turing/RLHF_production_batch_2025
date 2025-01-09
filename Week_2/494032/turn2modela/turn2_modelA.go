package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

type FilterRule struct {
	Keyword   string   `json:"keyword"`
	Replace   string   `json:"replace"`
	Blocked   bool     `json:"blocked"`
	UserRoles []string `json:"user_roles"`
}

type FilterService struct {
	rules []FilterRule
}

func NewFilterService(rules []FilterRule) *FilterService {
	return &FilterService{rules: rules}
}

func (s *FilterService) FilterContent(content string, userRole string, callback func(string, error)) {
	filteredContent := content
	for _, rule := range s.rules {
		if strings.Contains(filteredContent, rule.Keyword) {
			// Check if the rule applies to the current user role
			if containsRole(rule.UserRoles, userRole) {
				if rule.Blocked {
					filteredContent = ""
					callback(filteredContent, fmt.Errorf("Content contains blocked keyword '%s' for user with role '%s'", rule.Keyword, userRole))
					return
				} else {
					filteredContent = strings.Replace(filteredContent, rule.Keyword, rule.Replace, -1)
				}
			}
		}
	}
	callback(filteredContent, nil)
}

func containsRole(roles []string, role string) bool {
	for _, r := range roles {
		if r == role {
			return true
		}
	}
	return false
}

func main() {
	rules := []FilterRule{
		{Keyword: "badword", Replace: "******", Blocked: false},
		{Keyword: "evil", Blocked: true, UserRoles: []string{"admin", "moderator"}},
		{Keyword: "sensitive", Replace: "REDACTED", Blocked: false, UserRoles: []string{"admin"}},
	}

	filterService := NewFilterService(rules)

	// Hypothetical API gateway handler
	http.HandleFunc("/filter", func(w http.ResponseWriter, r *http.Request) {
		type Request struct {
			Content  string `json:"content"`
			UserRole string `json:"user_role"`
		}

		var req Request
		err := json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		filterService.FilterContent(req.Content, req.UserRole, func(filtered string, err error) {
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
