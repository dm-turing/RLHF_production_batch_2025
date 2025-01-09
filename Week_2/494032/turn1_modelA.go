package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

type FilterRule struct {
	Keyword  string  `json:"keyword"`
	Replace  string  `json:"replace"`
	Blocked bool    `json:"blocked"`
}

type FilterService struct {
	rules []FilterRule
}

func NewFilterService(rules []FilterRule) *FilterService {
	return &FilterService{rules: rules}
}

func (s *FilterService) FilterContent(content string, callback func(string, error)) {
	filteredContent := content
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

func main() {
	rules := []FilterRule{
		{Keyword: "badword", Replace: "******", Blocked: false},
		{Keyword: "evil", Blocked: true},
	}

	filterService := NewFilterService(rules)

	// Hypothetical API gateway handler
	http.HandleFunc("/filter", func(w http.ResponseWriter, r *http.Request) {
		type Request struct {
			Content string `json:"content"`
		}

		var req Request
		err := json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		filterService.FilterContent(req.Content, func(filtered string, err error) {
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