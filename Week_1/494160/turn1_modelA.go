package main

import (
	"fmt"
	"time"
)

// UserSession represents a user session
type UserSession struct {
	UserID    string    `json:"user_id"`
	Username  string    `json:"username"`
	LoggedIn  bool      `json:"logged_in"`
	LastLogin time.Time `json:"last_login"`
}

// SessionManager is a struct to manage user sessions
type SessionManager struct {
	sessions map[string]UserSession // Map of session IDs to user sessions
}

// NewSessionManager creates a new session manager
func NewSessionManager() *SessionManager {
	return &SessionManager{
		sessions: make(map[string]UserSession),
	}
}

// CreateSession creates a new session for a user
func (sm *SessionManager) CreateSession(userID, username string) string {
	sessionID := generateSessionID() // Assume generateSessionID is a function that generates a unique session ID
	session := UserSession{
		UserID:    userID,
		Username:  username,
		LoggedIn:  true,
		LastLogin: time.Now(),
	}
	sm.sessions[sessionID] = session
	return sessionID
}

// GetSession retrieves a session by ID
func (sm *SessionManager) GetSession(sessionID string) *UserSession {
	session, ok := sm.sessions[sessionID]
	if !ok {
		return nil
	}
	return &session
}

// InvalidateSession invalidates a session
func (sm *SessionManager) InvalidateSession(sessionID string) {
	delete(sm.sessions, sessionID)
}

// generateSessionID generates a unique session ID (this is a placeholder)
func generateSessionID() string {
	return "example-session-id" // In a real application, you would use a cryptographically secure function
}

func main() {
	sm := NewSessionManager()

	// Create a session for user with ID "user123"
	sessionID := sm.CreateSession("user123", "alice")
	fmt.Println("Created session:", sessionID)

	// Retrieve the session
	session := sm.GetSession(sessionID)
	if session != nil {
		fmt.Printf("Session Details:\nUserID: %s\nUsername: %s\nLoggedIn: %t\nLastLogin: %s\n",
			session.UserID, session.Username, session.LoggedIn, session.LastLogin)
	} else {
		fmt.Println("Session not found.")
	}

	// Invalidate the session
	sm.InvalidateSession(sessionID)

	// Retrieve the session again to show it's invalidated
	session = sm.GetSession(sessionID)
	if session != nil {
		fmt.Println("Session Details:", *session)
	} else {
		fmt.Println("Session not found.")
	}
}
