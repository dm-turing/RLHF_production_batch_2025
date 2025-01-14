package main  
import (  
    "database/sql"
    "encoding/json"
    "fmt"
    "log"
    "net/http"
    "time"

    "github.com/go-redis/redis"
	_ "github.com/ncruces/go-sqlite3/driver"
	_ "github.com/ncruces/go-sqlite3/embed"

)

// Define User struct
type User struct {  
    ID           int    `json:"id"`
    Username     string `json:"username"`
    Email        string `json:"email"`
    ConversionRate float32 `json:"conversion_rate"`
    SessionDuration time.Duration `json:"session_duration"`
}

var (  
    db     *sql.DB
    rdb    *redis.Client
    err    error
)

// initializes database connection
func initDB() {  
    db, err = sql.Open("sqlite3", "./mydata.db")
    if err != nil {
        log.Fatalf("Error opening database: %v", err)
    }
    // create table if not exists
    createTableSQL := `
        CREATE TABLE IF NOT EXISTS users (
            id INTEGER PRIMARY KEY AUTOINCREMENT,
            username TEXT,
            email TEXT,
            conversion_rate REAL,
            session_duration INTEGER
        )
    `
    _, err = db.Exec(createTableSQL)
    if err != nil {
        log.Fatalf("Error creating table: %v", err)
    }
}

// initializes redis connection
func initRedis() {  
    rdb = redis.NewClient(&redis.Options{  
        Addr:     "localhost:6379",   // Redis default port
        Password: "",                 // No password set
        DB:       0,                  // Default database
    })
    pong, err := rdb.Ping().Result()
    if err != nil {
        log.Fatalf("Error connecting to Redis: %v", err)
    }
    fmt.Printf("Redis Ping: %s\n", pong)
}

// close database and redis connections
func closeConnections() {  
    err := db.Close()
    if err != nil {
        log.Printf("Error closing database: %v", err)
    }
    err = rdb.Close()
    if err != nil {
        log.Printf("Error closing Redis: %v", err)
    }
}

// handler for POST /users
func createUser(w http.ResponseWriter, r *http.Request) {  
    var user User
    err := json.NewDecoder(r.Body).Decode(&user)
    if err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }
    insertSQL := "INSERT INTO users (username, email, conversion_rate, session_duration) VALUES (?, ?, ?, ?)"
    result, err := db.Exec(insertSQL, user.Username, user.Email, user.ConversionRate, user.SessionDuration.Nanoseconds())
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    userID, err := result.LastInsertId()
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    user.ID = int(userID)
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(user)
}

// handler for GET /users/{id}
func getUser(w http.ResponseWriter, r *http.Request) {  
    userID := r.FormValue("id")
    if userID == "" {
        http.Error(w, "Invalid user ID", http.StatusBadRequest)
        return
    }
    var user User
    row := db.QueryRow("SELECT id, username, email, conversion_rate, session_duration FROM users WHERE id = ?", userID)
    err := row.Scan(&user.ID, &user.Username, &user.Email, &user.ConversionRate, &user.SessionDuration)