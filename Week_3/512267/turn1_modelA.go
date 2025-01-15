package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"

	"github.com/go-redis/redis"
	_ "github.com/ncruces/go-sqlite3/driver"
	_ "github.com/ncruces/go-sqlite3/embed"
)

type User struct {
	ID    int    `json:"id"`
	User  string `json:"username"`
	Email string `json:"email"`
}

var (
	db          *sql.DB
	redisClient *redis.Client
)

func init() {
	var err error
	db, err = sql.Open("sqlite3", "users.db")
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	// Setting up a signal handler
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-signals
		log.Println("Shutting down the database instance gracefully...")
		db.Close() // Close the database connection
		os.Exit(0)
	}()

	redisClient = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	http.HandleFunc("/users", handleUsers)
	http.HandleFunc("/users/{id}", handleUser)
	log.Println("Server running on :8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
	select {}
}

func handleUsers(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		users, err := getAllUsers()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		json.NewEncoder(w).Encode(users)
	} else if r.Method == http.MethodPost {
		user := new(User)
		if err := json.NewDecoder(r.Body).Decode(user); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		if err := saveUser(user); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusCreated)
	} else {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	}
}

func handleUser(w http.ResponseWriter, r *http.Request) {
	var userID int
	parts := strings.Split(r.URL.Path, "/")
	if len(parts) > 2 {
		_userID, err := strconv.Atoi(parts[2])
		userID = _userID
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
	}

	user, err := getUserByID(userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(user)
}

func saveUser(user *User) error {
	stmt, err := db.Prepare("INSERT INTO users (username, email, password) VALUES (?, ?, ?)")
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(user.User, user.Email, user.Email) // Simplified for demonstration
	if err != nil {
		return err
	}

	redisKey := fmt.Sprintf("user:%d", user.ID)
	// Create a map to hold the interface values
	interfaceMap := make(map[string]interface{})

	// Copy values from the string map to the interface map
	for key, value := range map[string]string{
		"username": user.User,
		"email":    user.Email,
	} {
		interfaceMap[key] = value
	}

	// Now call HMSet with the correct type
	return redisClient.HMSet(redisKey, interfaceMap).Err()
}

func getAllUsers() ([]*User, error) {
	rows, err := db.Query("SELECT id, username, email FROM users")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	users := make([]*User, 0)
	for rows.Next() {
		user := new(User)
		if err := rows.Scan(&user.ID, &user.User, &user.Email); err != nil {
			return nil, err
		}
		users = append(users, user)
	}
	return users, nil
}

func getUserByID(id int) (*User, error) {
	user := new(User)
	row := db.QueryRow("SELECT id, username, email FROM users WHERE id = ?", id)
	if row == nil {
		return nil, sql.ErrNoRows
	}

	if err := row.Scan(&user.ID, &user.User, &user.Email); err != nil {
		return nil, err
	}

	redisKey := fmt.Sprintf("user:%d", user.ID)
	userData, err := redisClient.HGetAll(redisKey).Result()
	if err == nil {
		fmt.Printf("User data from Redis: %+v\n", userData)
	} else {
		fmt.Printf("Error fetching user data from Redis: %v\n", err)
	}

	return user, nil
}
