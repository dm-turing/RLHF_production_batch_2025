package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"regexp"
	"strings"

	"github.com/gorilla/mux"
	_ "github.com/ncruces/go-sqlite3/driver"
	_ "github.com/ncruces/go-sqlite3/embed"
)

var db *sql.DB

func init() {
	var err error
	db, err = sql.Open("sqlite3", "./email_checker.db")
	if err != nil {
		log.Fatal(err)
	}

	createTableQuery := `
	CREATE TABLE IF NOT EXISTS users (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		username TEXT NOT NULL,
		email TEXT NOT NULL
	);
	`
	_, err = db.Exec(createTableQuery)
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	router := mux.NewRouter()
	router.HandleFunc("/check/{username}/{email}", checkEmailHandler).Methods("GET")

	fmt.Println("Server is listening on port 8080...")
	log.Fatal(http.ListenAndServe(":8080", router))
}

func checkEmailHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	username := vars["username"]
	email := vars["email"]

	if isValidUserEmail(username, email) {
		saveUser(username, email)
		fmt.Fprintf(w, "User and email saved: %s, %s\n", username, email)
	} else {
		http.Error(w, "Email validation failed", http.StatusBadRequest)
	}
}

func isValidUserEmail(username, email string) bool {
	firstChar := strings.ToUpper(string(username[0])) // Normalize to uppercase

	switch {
	case firstChar >= "A" && firstChar <= "J" && isEmailDomain(email, "gmail.com"):
		log.Printf("Error: User %s is not allowed to have a Gmail account: %s\n", username, email)
		return false
	case firstChar >= "K" && firstChar <= "Q" && isEmailDomain(email, "yahoo.com"):
		log.Printf("Error: User %s is not allowed to have a Yahoo account: %s\n", username, email)
		return false
	case firstChar >= "R" && firstChar <= "Z" && isEmailDomain(email, "hotmail.com"):
		log.Printf("Error: User %s is not allowed to have a Hotmail account: %s\n", username, email)
		return false
	}
	return true
}

func isEmailDomain(email, domain string) bool {
	regexPattern := fmt.Sprintf(`^[a-zA-Z0-9._%%+-]+@%s$`, domain)
	match, _ := regexp.MatchString(regexPattern, email)
	return match
}

func saveUser(username, email string) {
	_, err := db.Exec("INSERT INTO users (username, email) VALUES (?, ?)", username, email)
	if err != nil {
		log.Printf("Error saving user %s with email %s: %v\n", username, email, err)
	}
}
