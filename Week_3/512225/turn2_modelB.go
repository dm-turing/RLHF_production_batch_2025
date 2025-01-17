package main

import (
 "encoding/json"
 "fmt"
 "log"
 "net/http"
 "net/url"
 "regexp"
 "time"

 "github.com/deatil/go-encoding/base62"
 "github.com/go-playground/validator/v10"
 "github.com/gorilla/mux"
 "github.com/jinzhu/gorm"
 _ "github.com/mattn/go-sqlite3"
)

// URL model
type URL struct {
 gorm.Model
 OriginalURL string `gorm:"unique;not null" valid:"required,url"`
 ShortURL    string `gorm:"unique;not null"`
 UserID      uint   `gorm:"default:0"`
}

// User model
type User struct {
 gorm.Model
 Username string `gorm:"unique;not null" valid:"required,alphanum,min=3"`
 Password string `gorm:"not null" valid:"required,min=6"`
}

var db *gorm.DB
var validator *validator.Validate

const encodeStd = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"

var urlRegexp = regexp.MustCompile(`^(?:http(s)?://)?[\w.-]+(?:\.[\w\.-]+)+[\w\-\._~:/?#[\]@!\$&'\(\)\*\+,;=.]+$`)

func main() {
 var err error
 db, err = gorm.Open("sqlite3", "urlshortener.db")
 if err != nil {
  log.Fatalf("failed to connect to database: %v", err.Error())
 }
 defer db.Close()

 // Migrate the database
 if err := db.AutoMigrate(&User{}, &URL{}); err != nil {
  log.Fatalf("failed to migrate database: %v", err)
 }

 // Initialize the validator
 validator = validator.New()

 router := mux.NewRouter()
 router.HandleFunc("/shorten", shortenURL).Methods("POST")
 router.HandleFunc("/{shortURL}", getOriginalURL).Methods("GET")
