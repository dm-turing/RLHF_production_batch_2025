package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

// ErrorLog defines the structure of an error log.
type ErrorLog struct {
	Errors []string `json:"errors"`
}

// Middleware is a function type that takes an http.Handler and returns an http.Handler.
type Middleware func(http.Handler) http.Handler

// LoggerMiddleware appends a generic error log message.
func LoggerMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		r.Header.Set("X-Error-Log", r.Header.Get("X-Error-Log")+"Logger middleware error detected: Unable to connect to DB. ")
		next.ServeHTTP(w, r)
	})
}

// UserInputMiddleware appends user-specific error messages.
func UserInputMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var el ErrorLog
		err := json.NewDecoder(r.Body).Decode(&el)
		if err != nil {
			http.Error(w, "Invalid input", http.StatusBadRequest)
			return
		}

		for _, errMsg := range el.Errors {
			currentLog := r.Header.Get("X-Error-Log")
			r.Header.Set("X-Error-Log", currentLog+"User error: "+errMsg+". ")
		}

		next.ServeHTTP(w, r)
	})
}

// FinalHandler generates the final error log message.
func FinalHandler(w http.ResponseWriter, r *http.Request) {
	finalLog := r.Header.Get("X-Error-Log")
	response := map[string]string{"final_error_log": finalLog}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// ChainMiddleware applies a chain of middleware to a handler.
func ChainMiddleware(h http.Handler, middlewares ...Middleware) http.Handler {
	for _, middleware := range middlewares {
		h = middleware(h)
	}
	return h
}

func main() {
	http.Handle("/compose-error-log", ChainMiddleware(http.HandlerFunc(FinalHandler), LoggerMiddleware, UserInputMiddleware))

	fmt.Println("Server is running on port 8080...")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatalf("Could not start server: %s\n", err)
	}
}
