package main

import (
	"fmt"
	"net/http"
	"net/url"
)

func getABTestVersion(r *http.Request) string {
	// Parse the URL query string
	values, _ := url.ParseQuery(r.URL.RawQuery)

	// Get the 'ab_test' parameter or return a default
	version := values.Get("ab_test")
	if version == "" {
		version = "A" // Default to version A if no parameter is provided
	}

	return version
}

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		version := getABTestVersion(r)
		fmt.Fprintf(w, "You are assigned to version %s\n", version)
	})

	http.ListenAndServe(":8080", nil)
}
