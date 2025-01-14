package main

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
)

// URLMap is a simple map to store short URLs and their corresponding original URLs.
var URLMap = map[string]string{}

// ShortenURLHandler handles the request to shorten a URL.
func ShortenURLHandler(w http.ResponseWriter, r *http.Request) {
	originalURL := r.FormValue("url")
	if originalURL == "" {
		http.Error(w, "URL parameter is required", http.StatusBadRequest)
		return
	}

	shortURL := generateShortURL(originalURL)

	URLMap[shortURL] = originalURL

	response := map[string]string{"shortURL": shortURL}
	json.NewEncoder(w).Encode(response)
}

// GetOriginalURLHandler handles the request to retrieve the original URL from a short URL.
func GetOriginalURLHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	shortURL := vars["shortURL"]

	originalURL, ok := URLMap[shortURL]
	if !ok {
		http.Error(w, "Short URL not found", http.StatusNotFound)
		return
	}

	response := map[string]string{"originalURL": originalURL}
	json.NewEncoder(w).Encode(response)
}

func generateShortURL(originalURL string) string {
	// Base64 encode the original URL
	shortURL := base64.StdEncoding.EncodeToString([]byte(originalURL))

	// Replace any '/' characters with '-' to make the URL more readable
	shortURL = strings.ReplaceAll(shortURL, "/", "-")

	return shortURL
}

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/shorten", ShortenURLHandler).Methods("GET")
	r.HandleFunc("/{shortURL}", GetOriginalURLHandler).Methods("GET")

	fmt.Println("Server is running on port 8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}
