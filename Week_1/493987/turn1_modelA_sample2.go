package main

import (
	"fmt"
	"net/http"
	"net/url"
)

// ABTestManager interface
type ABTestManager interface {
	GetVersion(r *http.Request) string
	AssignUser(r *http.Request) (string, error)
}

// BasicABTestManager is a simple implementation
type BasicABTestManager struct {
}

func (m *BasicABTestManager) GetVersion(r *http.Request) string {
	values, _ := url.ParseQuery(r.URL.RawQuery)
	version := values.Get("ab_test")
	if version == "" {
		version = "A"
	}
	return version
}

func (m *BasicABTestManager) AssignUser(r *http.Request) (string, error) {
	// Simple assignment: choose randomly between A and B
	return m.GetVersion(r), nil
}

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		manager := BasicABTestManager{}
		version, err := manager.AssignUser(r)
		if err != nil {
			fmt.Fprintf(w, "Error assigning user: %v\n", err)
			return
		}
		fmt.Fprintf(w, "You are assigned to version %s\n", version)
	})

	http.ListenAndServe(":8080", nil)
}
