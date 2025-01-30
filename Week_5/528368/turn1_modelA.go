package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	_ "github.com/ncruces/go-sqlite3/driver"
	_ "github.com/ncruces/go-sqlite3/embed"
)

type Message struct {
	ID        int       `json:"id"`
	UserID    string    `json:"userid"`
	Context   string    `json:"context"`
	EventType string    `json:"event_type"`
	Timestamp time.Time `json:"timestamp"`
}

func setupDatabase() *sql.DB {
	db, err := sql.Open("sqlite3", "./messages.db")
	if err != nil {
		log.Fatal(err)
	}

	createTable := `
	CREATE TABLE IF NOT EXISTS messages (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		userid TEXT,
		context TEXT,
		event_type TEXT,
		timestamp DATETIME
	);
	`
	_, err = db.Exec(createTable)
	if err != nil {
		log.Fatal(err)
	}

	return db
}

func generateRandomMessage() Message {
	userIDs := []string{"user1", "user2", "user3"}
	contexts := []string{"context1", "context2", "context3"}
	eventTypes := []string{"click", "keypress", "scroll"}

	return Message{
		UserID:    userIDs[rand.Intn(len(userIDs))],
		Context:   contexts[rand.Intn(len(contexts))],
		EventType: eventTypes[rand.Intn(len(eventTypes))],
		Timestamp: time.Now(),
	}
}

func insertMessage(db *sql.DB, msg Message) {
	query := `INSERT INTO messages (userid, context, event_type, timestamp) VALUES (?, ?, ?, ?)`
	_, err := db.Exec(query, msg.UserID, msg.Context, msg.EventType, msg.Timestamp)
	if err != nil {
		log.Fatal(err)
	}
}

func messageGenerator(db *sql.DB) {
	for {
		msg := generateRandomMessage()
		insertMessage(db, msg)
		time.Sleep(1 * time.Second)
	}
}

func filterMessages(db *sql.DB, filters map[string]string) ([]Message, error) {
	query := "SELECT id, userid, context, event_type, timestamp FROM messages WHERE 1=1"
	args := []interface{}{}

	for key, value := range filters {
		if value != "" {
			query += fmt.Sprintf(" AND %s = ?", key)
			args = append(args, value)
		}
	}

	rows, err := db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var messages []Message
	for rows.Next() {
		var msg Message
		if err := rows.Scan(&msg.ID, &msg.UserID, &msg.Context, &msg.EventType, &msg.Timestamp); err != nil {
			return nil, err
		}
		messages = append(messages, msg)
	}

	return messages, nil
}

func messagesHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		filters := map[string]string{
			"userid":     r.URL.Query().Get("userid"),
			"context":    r.URL.Query().Get("context"),
			"event_type": r.URL.Query().Get("event_type"),
		}

		messages, err := filterMessages(db, filters)
		if err != nil {
			http.Error(w, "Failed to retrieve messages", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(messages)
	}
}

func main() {
	db := setupDatabase()
	defer db.Close()

	go messageGenerator(db)

	r := mux.NewRouter()
	r.HandleFunc("/messages", messagesHandler(db)).Methods("GET")

	log.Fatal(http.ListenAndServe(":8080", r))
}
