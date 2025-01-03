package main

import (
	"fmt"
	"net/http"
)

// ExtendedABTestManager uses cookies for persistent assignments
type ExtendedABTestManager struct {
}

func (m *ExtendedABTestManager) GetVersion(r *http.Request) string {
	// Example: look for a cookie
	if cookie, err := r.Cookie("ab_test_version"); err == nil {
		return cookie.Value
	}
	return "A"
}

func (m *ExtendedABTestManager) AssignUser(r *http.Request) (string, error) {
	version := m.GetVersion(r)
	if version == "" {
		version = "A" // Default to version A
		http.SetCookie(r.ResponseWriter, &http.Cookie{
			Name:     "ab_test_version",
			Value:    version,
			Path:     "/",
			MaxAge:   3600 * 24 * 7, // Set the cookie to last for one week
			HttpOnly: true,
		})
	}
	return version, nil
}

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		manager := ExtendedABTestManager{}
		version, err := manager.AssignUser(r)
		if err != nil {
			fmt.Fprintf(w, "Error assigning user: %v\n", err)
			return
		}
		fmt.Fprintf(w, "You are assigned to version %s\n", version)
	})

	http.ListenAndServe(":8080", nil)
}
