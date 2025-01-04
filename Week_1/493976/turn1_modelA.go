package main

import (
	"fmt"
	"log"
	"net/http"
	"time"
)

// Audit represents a single audit entry for a query parameter change.
type Audit struct {
	Parameter string    `json:"parameter"`
	OldValue  string    `json:"old_value"`
	NewValue  string    `json:"new_value"`
	Timestamp time.Time `json:"timestamp"`
}

// Audits is a list of audit entries.
type Audits []Audit

// auditLogs is a slice that holds all the audit logs.
var auditLogs Audits

var oldParams map[string]string

func main() {
	http.HandleFunc("/query-params", handleQueryParams)
	log.Println("Server listening on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

// handleQueryParams processes an incoming HTTP request and logs changes to query parameters.
func handleQueryParams(w http.ResponseWriter, r *http.Request) {
	name := r.URL.Query().Get("name")
	age := r.URL.Query().Get("age")

	// Capture current query parameters.
	if _, ok := oldParams[name]; !ok {
		oldParams[name] = age
	}

	// Process the query parameters in the request.
	newParams := make(map[string]string)
	newParams[name] = age

	// Log changes.
	for param, oldValue := range oldParams {
		newValue := newParams[param]
		if oldValue != newValue {
			auditLogs = append(auditLogs, Audit{
				Parameter: param,
				OldValue:  oldValue,
				NewValue:  newValue,
				Timestamp: time.Now(),
			})
			oldParams[param] = newValue
			log.Printf("Parameter %s changed from %s to %s\n", param, oldValue, newValue)
		}
	}

	// Respond to the client.
	fmt.Fprintln(w, "Query parameters processed successfully.")
}

func init() {
	oldParams = make(map[string]string)
}
