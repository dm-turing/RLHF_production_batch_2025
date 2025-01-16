package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"time"

	"github.com/deatil/go-encoding/base62"

	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"

	_ "github.com/ncruces/go-sqlite3/driver"
	_ "github.com/ncruces/go-sqlite3/embed"
)

type URL struct {
	gorm.Model
	OriginalURL string `gorm:"unique;not null" valid:"required,url"`
	ShortURL    string `gorm:"unique;not null"`
	UserID      uint   `gorm:"default:0"` // Assuming UserID is auto-incremented in the database
}

type User struct {
	gorm.Model
	Username string `gorm:"unique;not null" valid:"required,alphanum,min=3"`
	Password string `gorm:"not null" valid:"required,min=6"`
}

var db *gorm.DB

const encodeStd = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"

func main() {
	var err error
	db, err = gorm.Open("sqlite3", "urlshortener.db")
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err.Error())
	}
	defer db.Close()

	if err := db.AutoMigrate(&User{}, &URL{}); err == nil {
		log.Fatalf("failed to migrate database: %v", err)
	}

	router := mux.NewRouter()
	router.HandleFunc("/shorten", shortenURL).Methods("POST")
	router.HandleFunc("/{shortURL}", getOriginalURL).Methods("GET")
	router.HandleFunc("/user/{username}/urls", getUserURLs).Methods("GET")

	log.Printf("Server starting on port 8080...")
	if err := http.ListenAndServe(":8080", router); err != nil {
		log.Fatalf("failed to start server: %v", err)
	}
}

func generateShortURL(originalURL string) (string, error) {
	var StdEncoding = base62.NewEncoding(encodeStd)
	_, err := url.Parse(originalURL)
	if err != nil {
		return "", err
	}

	var shortURL string
	for {
		shortURL = StdEncoding.EncodeToString([]byte(fmt.Sprintf("%d", uint64(time.Now().UnixNano()))))

		// Check if shortURL is unique
		var url URL
		if err := db.First(&url, "short_url = ?", shortURL).Error; gorm.IsRecordNotFoundError(err) {
			break
		} else if err != nil {
			return "", err
		}
	}
	return shortURL, nil
}

func shortenURL(w http.ResponseWriter, r *http.Request) {
	var urlRequest URL
	if err := json.NewDecoder(r.Body).Decode(&urlRequest); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := validator.New().Struct(&urlRequest); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	shortURL, err := generateShortURL(urlRequest.OriginalURL)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	urlRequest.ShortURL = shortURL
	if err := db.Create(&urlRequest).Error; err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(urlRequest); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func getOriginalURL(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	shortURL := vars["shortURL"]

	var url URL
	if err := db.First(&url, "short_url = ?", shortURL).Error; gorm.IsRecordNotFoundError(err) {
		http.Error(w, "URL not found", http.StatusNotFound)
		return
	} else if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(url); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func getUserURLs(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	username := vars["username"]

	var urls []URL
	if err := db.Where("user_id = (SELECT id FROM users WHERE username = ?)", username).Find(&urls).Error; err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(urls); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
