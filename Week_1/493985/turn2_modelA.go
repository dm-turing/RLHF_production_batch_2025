package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"
)

type Session struct {
	UserID   int       `json:"user_id"`
	Username string    `json:"username"`
	Expires  time.Time `json:"expires"`
}

var sessions = make(map[string]Session)
var sessionMutex = sync.RWMutex{}

func generateSessionID() string {
	// Simple session ID generation using time; in practice, you'd use a more secure method.
	return fmt.Sprintf("%d", time.Now().UnixNano())
}

func isSessionExpired(session Session) bool {
	return time.Now().After(session.Expires)
}

func handleLogin(w http.ResponseWriter, r *http.Request) {
	userID := 1            // Example user ID
	username := "testuser" // Example username
	sessionID := generateSessionID()

	session := Session{
		UserID:   userID,
		Username: username,
		Expires:  time.Now().Add(1 * time.Hour),
	}

	sessionMutex.Lock()
	defer sessionMutex.Unlock()
	sessions[sessionID] = session

	// Set a session cookie
	cookie := http.Cookie{
		Name:     "session",
		Value:    sessionID,
		Path:     "/",
		HttpOnly: true,
		Secure:   false, // Set to true in production
		MaxAge:   int(time.Until(session.Expires).Seconds()),
	}
	http.SetCookie(w, &cookie)

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Login successful"))
}

func handleLogout(w http.ResponseWriter, r *http.Request) {
	sessionMutex.Lock()
	defer sessionMutex.Unlock()

	cookie, err := r.Cookie("session")
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte("No session found"))
		return
	}

	delete(sessions, cookie.Value)

	// Delete the session cookie
	http.SetCookie(w, &http.Cookie{
		Name:     "session",
		Value:    "",
		Path:     "/",
		HttpOnly: true,
		MaxAge:   -1,
	})

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Logout successful"))
}

func handleProfile(w http.ResponseWriter, r *http.Request) {
	sessionMutex.RLock()
	defer sessionMutex.RUnlock()

	cookie, err := r.Cookie("session")
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte("No session found"))
		return
	}

	session, ok := sessions[cookie.Value]
	if !ok || isSessionExpired(session) {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte("Session expired"))
		return
	}

	userData, err := json.Marshal(session)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Error marshalling user data"))
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(userData)
}

func main() {
	http.HandleFunc("/login", handleLogin)
	http.HandleFunc("/logout", handleLogout)
	http.HandleFunc("/profile", handleProfile)

	log.Fatal(http.ListenAndServe(":8080", nil))
}
