package main

import (
	"database/sql"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	_ "github.com/ncruces/go-sqlite3/driver"
	_ "github.com/ncruces/go-sqlite3/embed"
)

// User defines the user structure
type User struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

// PasswordRecord defines the structure for a password change record
type PasswordRecord struct {
	UserID    int       `json:"user_id"`
	Password  string    `json:"password"`
	Timestamp time.Time `json:"timestamp"`
	Reason    string    `json:"reason"`
}

var db *sql.DB

func init() {
	var err error
	db, err = sql.Open("sqlite3", "./passwordManager.db")
	if err != nil {
		log.Fatal(err)
	}

	// Create tables if they don't exist
	createUserTable := `
	CREATE TABLE IF NOT EXISTS users (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT NOT NULL
	);`

	createPasswordTable := `
	CREATE TABLE IF NOT EXISTS passwords (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		user_id INTEGER,
		password TEXT,
		timestamp DATETIME,
		reason TEXT,
		FOREIGN KEY (user_id) REFERENCES users(id)
	);`

	_, err = db.Exec(createUserTable)
	if err != nil {
		log.Fatal(err)
	}

	_, err = db.Exec(createPasswordTable)
	if err != nil {
		log.Fatal(err)
	}
}

// Generate a random password that meets the criteria
func generatePassword() (string, error) {
	const (
		lower          = "abcdefghijklmnopqrstuvwxyz"
		upper          = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
		digits         = "0123456789"
		special        = "!@#$%^&*()-_=+[]{}<>?,./"
		all            = lower + upper + digits + special
		passwordLength = 10
	)

	password := make([]byte, passwordLength)

	// Ensure at least one character from each required set
	sets := []string{lower, upper, digits, special}
	for i, set := range sets {
		c, err := randomCharFromSet(set)
		if err != nil {
			return "", err
		}
		password[i] = c
	}

	// Fill remaining spaces with random characters from all character sets
	for i := len(sets); i < passwordLength; i++ {
		c, err := randomCharFromSet(all)
		if err != nil {
			return "", err
		}
		password[i] = c
	}

	// Shuffle the password to ensure randomness
	shuffled, err := shuffle(password)
	if err != nil {
		return "", err
	}

	return string(shuffled), nil
}

// Select a random character from a given set
func randomCharFromSet(set string) (byte, error) {
	index := rand.Intn(len(set))
	return set[index], nil
}

// Shuffle the password
func shuffle(password []byte) ([]byte, error) {
	perm := rand.Perm(len(password))
	shuffled := make([]byte, len(password))
	for i, v := range perm {
		shuffled[v] = password[i]
	}
	return shuffled, nil
}

func main() {
	// Insert users into the database
	users := []User{
		{Name: "Alice"},
		{Name: "Bob"},
		{Name: "Charlie"},
		{Name: "David"},
		{Name: "Eve"},
	}

	for _, user := range users {
		_, err := db.Exec("INSERT INTO users (name) VALUES (?)", user.Name)
		if err != nil {
			log.Fatalf("Error while inserting user %s: %v", user.Name, err)
		}
	}

	// Set up a ticker for password generation
	ticker := time.NewTicker(2 * time.Second)
	defer ticker.Stop()

	go func() {
		for range ticker.C {
			userIndex := rand.Intn(len(users))
			userID := userIndex + 1 // userIDs start from 1 in DB

			password, err := generatePassword()
			if err != nil {
				fmt.Println("Error generating password:", err)
				continue
			}

			reason := "Routine password change"
			timestamp := time.Now()

			_, err = db.Exec(`INSERT INTO passwords (user_id, password, timestamp, reason) 
							  VALUES (?, ?, ?, ?)`, userID, password, timestamp, reason)
			if err != nil {
				fmt.Println("Error inserting password record:", err)
				continue
			}
		}
	}()

	// Set up the HTTP server
	router := gin.Default()

	router.GET("/users", listUsers)
	router.GET("/passwords/:userID", listPasswords)

	router.Run(":8080")
}

// listUsers responds with the list of users
func listUsers(c *gin.Context) {
	rows, err := db.Query("SELECT id, name FROM users")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer rows.Close()

	var users []User
	for rows.Next() {
		var user User
		if err := rows.Scan(&user.ID, &user.Name); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		users = append(users, user)
	}

	c.JSON(http.StatusOK, users)
}

// listPasswords responds with the password info for a user, with optional sorting by timestamp
func listPasswords(c *gin.Context) {
	userID := c.Param("userID")
	sortOrder := c.DefaultQuery("sort", "asc")

	rows, err := db.Query(`
		SELECT user_id, password, timestamp, reason
		FROM passwords
		WHERE user_id = ?
		ORDER BY timestamp `+sortOrder, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer rows.Close()

	var passwordRecords []PasswordRecord
	for rows.Next() {
		var record PasswordRecord
		if err := rows.Scan(&record.UserID, &record.Password, &record.Timestamp, &record.Reason); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		passwordRecords = append(passwordRecords, record)
	}

	c.JSON(http.StatusOK, passwordRecords)
}
