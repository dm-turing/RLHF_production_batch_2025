package main

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/go-redis/redis"
)

// UserSession represents a user session
type UserSession struct {
	UserID    string    `json:"user_id"`
	Username  string    `json:"username"`
	LoggedIn  bool      `json:"logged_in"`
	LastLogin time.Time `json:"last_login"`
}

// Connect to Redis for caching sessions
var client *redis.Client = redis.NewClient(&redis.Options{
	Addr:     "localhost:6379",
	Password: "", // no password set
	DB:       0,  // use default DB
})

// StoreSession caches a session in Redis
func StoreSession(sessionID string, userSession UserSession, expiration time.Duration) error {
	return client.Set(sessionID, userSession, expiration).Err()
}

// GetSession retrieves a session from Redis
func GetSession(sessionID string) (*UserSession, error) {
	rawSession, err := client.Get(sessionID).Bytes()
	if err != nil {
		return nil, err
	}

	var session UserSession
	if err := json.Unmarshal(rawSession, &session); err != nil {
		return nil, err
	}

	return &session, nil
}

func main() {
	// Store a session in Redis
	sessionID := "example-session-id"
	session := UserSession{
		UserID:    "user123",
		Username:  "alice",
		LoggedIn:  true,
		LastLogin: time.Now(),
	}
	err := StoreSession(sessionID, session, 1*time.Hour) // Cache session for 1 hour
	if err != nil {
		fmt.Println("Error storing session:", err)
	}

	// Retrieve a session from Redis
	retrievedSession, err := GetSession(sessionID)
	if err != nil {
		fmt.Println("Error retrieving session:", err)
	} else if retrievedSession != nil {
		fmt.Println("Retrieved Session:", *retrievedSession)
	}
}
