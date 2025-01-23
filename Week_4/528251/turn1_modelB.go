

package main  
import (  
    "database/sql"
    "fmt"
    "log"
    "net/http"
    "strings"
    "regexp"

	_ "github.com/ncruces/go-sqlite3/driver"
	_ "github.com/ncruces/go-sqlite3/embed"
)

// User struct to store user information
type User struct {
    Username string `json:"username"`
    Email    string `json:"email"`
}

// Create a global SQLite database pointer
var db *sql.DB

// Initialize the database and create the user table
func initDB() {
    var err error
    db, err = sql.Open("sqlite3", "./email_checker.db")
    if err != nil {
        log.Fatalf("Error opening database: %q", err)
    }
    _, err = db.Exec("CREATE TABLE IF NOT EXISTS users (username TEXT PRIMARY KEY, email TEXT)")
    if err != nil {
        log.Fatalf("Error creating user table: %q", err)
    }
}

// Helper function to check if the email address is valid
func isValidEmail(email string) bool {
    re := regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")
    return re.MatchString(email)
}

// Helper function to check if the user's email provider is allowed based on username
func isAllowedEmailProvider(username string, email string) bool {
    usernamePrefix := strings.ToUpper(string(username[0]))
    emailProvider := strings.Split(email, "@")[1]

    if usernamePrefix >= "A" && usernamePrefix <= "J" && emailProvider == "gmail.com" {
        return false
    }

    if usernamePrefix >= "K" && usernamePrefix <= "Q" && emailProvider == "yahoo.com" {
        return false
    }

    if usernamePrefix >= "R" && usernamePrefix <= "Z" && emailProvider == "hotmail.com" {
        return false
    }

    return true
}

// Handle the /check endpoint for user email validation
func checkEmailHandler(w http.ResponseWriter, r *http.Request) {
    if r.Method != "GET" {
        http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
        return
    }

    username := r.URL.Query().Get("username")
    email := r.URL.Query().Get("email")

    // Check if username and email are provided
    if username == "" || email == "" {
        http.Error(w, "Username and email are required", http.StatusBadRequest)
        return
    }

    // Validate the email address
    if !isValidEmail(email) {
        fmt.Fprintf(w, "Invalid email format for user: %s - Email: %s\n", username, email)
        return
    }

    // Check if the email provider is allowed based on the username
    if !isAllowedEmailProvider(username, email) {
        fmt.Fprintf(w, "Email provider not allowed for user: %s - Email: %s\n", username, email)
        return
    }

    // Insert the valid user into the database
    _, err := db.Exec("INSERT OR REPLACE INTO users (username, email) VALUES (?, ?)", username, email)
    if err != nil {
        log.Printf("Error inserting user into database: %q", err)
        http.Error(w, "Internal server error", http.StatusInternalServerError)
        return
    }

    fmt.Fprintf(w, "User %s added successfully with email: %s\n", username, email)