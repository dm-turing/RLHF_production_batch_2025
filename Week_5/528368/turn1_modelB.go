package main

import (
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"time"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"github.com/labstack/echo/v4"
)

// Message struct represents a message to be logged
type Message struct {
	gorm.Model
	UserID    int    `json:"user_id"`
	Context   string `json:"context"`
	Event     string `json:"event"`
	Timestamp int64  `json:"timestamp"`
}

// NewDB creates a new SQLite database connection
func NewDB() *gorm.DB {
	db, err := gorm.Open("sqlite3", "messages.db")
	if err != nil {
		log.Fatalf("Error opening database: %v", err)
	}

	// Migrate the Message struct to the database
	db.AutoMigrate(&Message{})

	return db
}

// GenerateMessage generates a new random message
func GenerateMessage() *Message {
	return &Message{
		UserID:    rand.Intn(100) + 1,
		Context:   fmt.Sprintf("Context %d", rand.Intn(100)),
		Event:     fmt.Sprintf("Event %d", rand.Intn(10)),
		Timestamp: time.Now().Unix(),
	}
}

// LogMessage logs a message to the database
func LogMessage(db *gorm.DB, message *Message) {
	if err := db.Create(message).Error; err != nil {
		log.Printf("Error logging message: %v", err)
	}
}

func main() {
	db := NewDB()
	defer db.Close()

	// Start generating and logging messages in a goroutine
	go func() {
		ticker := time.NewTicker(time.Second)
		defer ticker.Stop()
		for range ticker.C {
			msg := GenerateMessage()
			LogMessage(db, msg)
		}
	}()

	// Start the REST API server
	e := echo.New()
	e.GET("/messages", func(c echo.Context) error {
		var messages []Message
		query := db
		userID := c.QueryParam("user_id")
		if userID != "" {
			query = query.Where("user_id = ?", userID)
		}
		context := c.QueryParam("context")
		if context != "" {
			query = query.Where("context LIKE ?", fmt.Sprintf("%%%s%%", context))
		}
		event := c.QueryParam("event")
		if event != "" {
			query = query.Where("event = ?", event)
		}

		if err := query.Find(&messages).Error; err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
		}

		return c.JSON(http.StatusOK, messages)
	})
	e.Logger.Fatal(e.Start(":1323"))
}
