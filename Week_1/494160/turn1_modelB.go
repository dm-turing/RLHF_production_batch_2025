package main

import (
	"fmt"
	"sync"
	"time"
)

type UserSession struct {
	UserID       string
	LastActiveAt time.Time
	MetaData     map[string]interface{}
}

// Declare a global map to store user sessions
var (
	userSessions = map[string]UserSession{}
	sessionMutex sync.RWMutex
)

// LoadUserSession loads a user session from the session store
func LoadUserSession(sessionID string) (*UserSession, error) {
	sessionMutex.RLock()
	defer sessionMutex.RUnlock()
	if session, ok := userSessions[sessionID]; ok {
		return &session, nil
	}
	return nil, fmt.Errorf("session not found")
}

// SaveUserSession saves a user session to the session store
func SaveUserSession(sessionID string, session UserSession) error {
	sessionMutex.Lock()
	defer sessionMutex.Unlock()
	userSessions[sessionID] = session
	return nil
}

// DeleteUserSession deletes a user session from the session store
func DeleteUserSession(sessionID string) error {
	sessionMutex.Lock()
	defer sessionMutex.Unlock()
	delete(userSessions, sessionID)
	return nil
}
func main() {
	// Example usage:
	sessionID1 := "session-id-1"
	userSession1 := UserSession{
		UserID:       "user1",
		LastActiveAt: time.Now(),
		MetaData: map[string]interface{}{
			"cartItems": []string{"Book", "Pencil", "Ruler"},
		},
	}
	err := SaveUserSession(sessionID1, userSession1)
	if err != nil {
		fmt.Println("Error saving session:", err)
		return
	}
	// Load and update session
	loadedSession, err := LoadUserSession(sessionID1)
	if err != nil {
		fmt.Println("Error loading session:", err)
		return
	}
	loadedSession.MetaData["newKey"] = "newValue"
	err = SaveUserSession(sessionID1, *loadedSession)
	if err != nil {
		fmt.Println("Error saving session:", err)
		return
	}
	// Delete session
	DeleteUserSession(sessionID1)
}
