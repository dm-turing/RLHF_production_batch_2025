
package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"math"
	"math/rand"
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	_ "github.com/mattn/go-sqlite3"
	"github.com/sirupsen/logrus"
)

type Severity string

const (
	DEBUG    Severity = "DEBUG"
	INFO     Severity = "INFO"
	WARNING  Severity = "WARNING"
	ERROR    Severity = "ERROR"
	CRITICAL Severity = "CRITICAL"
	FATAL    Severity = "FATAL"
)

type Error struct {
	Context   string
	UserID    int
	ErrorCode int
	Message   string
	Time      time.Time
	Severity  Severity
	Next      *Error
}

func (e *Error) Error() string {
	return fmt.Sprintf("Context: %s, UserID: %d, ErrorCode: %d, Message: %s, Severity: %s, Time: %s",
		e.Context, e.UserID, e.ErrorCode, e.Message, e.Severity, e.Time.Format(time.RFC3339))
}

var contexts = []string{"database", "network", "filesystem", "authentication"}
var messages = []string{"connection timeout", "missing file", "permission denied", "invalid credentials", "unknown error"}
var severities = []Severity{DEBUG, INFO, WARNING, ERROR, CRITICAL, FATAL}

var db *sql.DB
var errorMap = make(map[string][]*Error)
var mu sync.Mutex

var logger = logrus.New()

func init() {
	var err error
	db, err = sql.Open("sqlite3", "./errors.db")
	if err != nil {
		log.Fatal(err)
	}

	if _, err = db.Exec(`CREATE TABLE IF NOT EXISTS errors (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		context TEXT,
		user_id INTEGER,
		error_code INTEGER,
		message TEXT,
		time TEXT,
		severity TEXT
	)`); err != nil {
		log.Fatal(err)
	}

	logger.SetOutput(gin.DefaultWriter)
	logger.SetFormatter(&logrus.TextFormatter{FullTimestamp: true})
}

func generateRandomError() *Error {
	return &Error{
		Context:   contexts[rand.Intn(len(contexts))],
		UserID:    rand.Intn(1000),
		ErrorCode: rand.Intn(500),
		Message:   messages[rand.Intn(len(messages))],
		Time:      time.Now(),
		Severity:  severities[rand.Intn(len(severities))],
	}
}

func logErrorToDB(err *Error) {
	stmt, err := db.Prepare("INSERT INTO errors (context, user_id, error_code, message, time, severity) VALUES (?, ?, ?, ?, ?, ?)")
	if err != nil {
		logger.Errorf("Error preparing statement: %v", err)
		return
	}
	defer stmt.Close()

	_, err = stmt.Exec(err.Context, err.UserID, err.ErrorCode, err.Message, err.Time, err.Severity)
	if err != nil {
		logger.Errorf("Error executing statement: %v", err)
		return
	}
}

func getErrorsFromDB(context, severity string) ([]*Error, error) {
	rows, err := db.Query("SELECT context, user_id, error_code, message, time, severity FROM errors WHERE context LIKE ? AND severity LIKE ? ORDER BY time DESC", context+"%", severity+"%")