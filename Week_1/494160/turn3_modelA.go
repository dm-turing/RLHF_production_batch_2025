package main

import (
	"crypto/rand"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/go-redis/redis"
	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	"golang.org/x/crypto/nacl/secretbox"
)

// UserSession struct
type UserSession struct {
	UserID    string    `json:"user_id"`
	Username  string    `json:"username"`
	LoggedIn  bool      `json:"logged_in"`
	LastLogin time.Time `json:"last_login"`
}

// SessionManager struct
type SessionManager struct {
	store      *sessions.CookieStore
	redis      *redis.Client
	privateKey [32]byte
}

// NewSessionManager creates a new session manager
func NewSessionManager() *SessionManager {
	// Initialize Redis client
	redisClient := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})

	// Initialize Gorilla sessions store
	store := sessions.NewCookieStore([]byte("your-secure-cookie-secret"))

	// Set Redis as the cache store for sessions
	store.Options.MaxAge = 3600 // Session timeout in seconds
	store.Options.Secure = true // Use HTTPS
	store.Options.HttpOnly = true
	store.Options.Path = "/"

	return &SessionManager{
		store:      store,
		redis:      redisClient,
		privateKey: [32]byte{},
	}
}

// GenerateSessionToken generates a secure session token
func (sm *SessionManager) GenerateSessionToken() string {
	token := make([]byte, 32)
	if _, err := rand.Read(token); err != nil {
		log.Fatal(err)
	}
	return fmt.Sprintf("%x", token)
}

// EncryptSessionData encrypts session data
func (sm *SessionManager) EncryptSessionData(data []byte) ([]byte, error) {
	nonce := make([]byte, 24)
	if _, err := rand.Read(nonce); err != nil {
		return nil, err
	}
	return secretbox.Seal(nonce, data, nonce, &sm.privateKey), nil
}

// DecryptSessionData decrypts session data
func (sm *SessionManager) DecryptSessionData(data []byte) ([]byte, error) {
	nonce, data := data[:24], data[24:]
	return secretbox.Open(nil, data, nonce, &sm.privateKey)
}

// CreateSession creates a new session
func (sm *SessionManager) CreateSession(userID, username string) (string, error) {
	session := UserSession{
		UserID:    userID,
		Username:  username,
		LoggedIn:  true,
		LastLogin: time.Now(),
	}

	sessionData, err := json.Marshal(session)
	if err != nil {
		return "", err
	}

	encryptedData, err := sm.EncryptSessionData(sessionData)
	if err != nil {
		return "", err
	}

	sessionToken := sm.GenerateSessionToken()
	sm.redis.Set(sessionToken, encryptedData, 3600*time.Second).Err() // Cache for 1 hour

	return sessionToken, nil
}

// GetSession retrieves a session
func (sm *SessionManager) GetSession(sessionToken string) (*UserSession, error) {
	encryptedData, err := sm.redis.Get(sessionToken).Bytes()
	if err != nil {
		return nil, err
	}

	data, err := sm.DecryptSessionData(encryptedData)
	if err != nil {
		return nil, err
	}

	var session UserSession
	if err := json.Unmarshal(data, &session); err != nil {
		return nil, err
	}

	return &session, nil
}

// InvalidateSession invalidates a session
func (sm *SessionManager) InvalidateSession(sessionToken string) {
	sm.redis.Del(sessionToken).Err()
}

func main() {
	// Initialize session manager
	sm := NewSessionManager()

	// Generate a private key for encryption (should be done securely and stored)
	if _, err := rand.Read(sm.privateKey[:]); err != nil {
		log.Fatal(err)
	}

	// Create a simple HTTP handler
	r := mux.NewRouter()

	r.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) {
		sessionToken, err := sm.CreateSession("user123", "alice")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		http.SetCookie(w, &http.Cookie{Name: "session_token", Value: sessionToken})
		fmt.Fprintln(w, "Logged in successfully.")
	})

	r.HandleFunc("/profile", func(w http.ResponseWriter, r *http.Request) {
		sessionToken, err := r.Cookie("session_token")
		if err != nil || sessionToken.Value == "" {
			http.Error(w, "Not logged in.", http.StatusUnauthorized)
			return
		}

		session, err := sm.GetSession(sessionToken.Value)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		fmt.Fprintf(w, "Welcome, %s!", session.Username)
	})

	r.HandleFunc("/logout", func(w http.ResponseWriter, r *http.Request) {
		sessionToken, err := r.Cookie("session_token")
		if err != nil || sessionToken.Value == "" {
			http.Error(w, "Not logged in.", http.StatusUnauthorized)
			return
		}

		sm.InvalidateSession(sessionToken.Value)
		http.SetCookie(w, &http.Cookie{Name: "session_token", Value: "", Expires: time.Now().Add(-time.Hour)})
		fmt.Fprintln(w, "Logged out successfully.")
	})

	http.ListenAndServe(":8080", r)
}
