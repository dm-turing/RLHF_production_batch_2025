package main

import (
	"database/sql"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/mux"
	_ "github.com/ncruces/go-sqlite3/driver"
	_ "github.com/ncruces/go-sqlite3/embed"
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

var db *sql.DB

func initDB() {
	var err error
	db, err = sql.Open("sqlite3", "errors.db")
	if err != nil {
		log.Fatal(err)
	}

	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS logs (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		context TEXT,
		user_id INTEGER,
		error_code INTEGER,
		message TEXT,
		severity TEXT,
		time DATETIME
	)`)
	if err != nil {
		log.Fatal(err)
	}
}

func logErrorToSQLite(e *Error) {
	insertSQL := `INSERT INTO logs(context, user_id, error_code, message, severity, time) VALUES (?, ?, ?, ?, ?, ?)`
	_, err := db.Exec(insertSQL, e.Context, e.UserID, e.ErrorCode, e.Message, e.Severity, e.Time)
	if err != nil {
		log.Printf("Error inserting log into db: %v", err)
	}
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

var (
	contexts   = []string{"database", "network", "filesystem", "authentication"}
	messages   = []string{"connection timeout", "missing file", "permission denied", "invalid credentials", "unknown error"}
	severities = []Severity{DEBUG, INFO, WARNING, ERROR, CRITICAL, FATAL}
	mu         sync.Mutex
)

func getLogsHandler(w http.ResponseWriter, r *http.Request) {
	query := "SELECT context, user_id, error_code, message, severity, time FROM logs WHERE context = ? AND severity = ?"

	params := mux.Vars(r)
	rows, err := db.Query(query, params["context"], params["severity"])

	if err != nil {
		http.Error(w, "Failed to query logs", http.StatusInternalServerError)
		return
	}

	defer rows.Close()
	logResults := []Error{}
	for rows.Next() {
		var logRow Error
		var err error
		err = rows.Scan(&logRow.Context, &logRow.UserID, &logRow.ErrorCode, &logRow.Message, &logRow.Severity, &logRow.Time)
		if err == nil {
			logResults = append(logResults, logRow)
		} else {
			log.Printf("Failed to scan row: %v", err)
		}
	}
	for _, lg := range logResults {
		fmt.Fprintf(w, "%s: %s [%d]: %s\n", lg.Context, lg.Severity, lg.UserID, lg.Message)
	}
}

func printHealthAndStability() {
	ticker := time.NewTicker(1 * time.Minute)
	for range ticker.C {
		criticalCount, fatalCount, totalCount := 0, 0, 0

		mu.Lock()
		defer mu.Unlock()

		rows, err := db.Query(`SELECT severity FROM logs WHERE time > datetime('now', '-1 minute')`)
		if err != nil {
			log.Printf("Failed to query logs for health analysis: %v", err)
			continue
		}
		defer rows.Close()

		for rows.Next() {
			var severity Severity
			if err := rows.Scan(&severity); err != nil {
				log.Printf("Error scanning severity: %v", err)
				continue
			}
			totalCount++
			switch severity {
			case CRITICAL:
				criticalCount++
			case FATAL:
				fatalCount++
			}
		}

		printPercentageOfErrors(totalCount, criticalCount, fatalCount)
	}
}

func printPercentageOfErrors(totalCount, criticalCount, fatalCount int) {
	if totalCount == 0 {
		fmt.Printf("No errors in the past minute\n")
		return
	}

	fmt.Println("Error Severity Breakdown:")
	fmt.Printf("CRITICAL errors: %d (%.2f%%)\n", criticalCount, percentage(criticalCount, totalCount))
	fmt.Printf("FATAL errors: %d (%.2f%%)\n", fatalCount, percentage(fatalCount, totalCount))

	stability := computeStability(percentage(criticalCount+fatalCount, totalCount))
	fmt.Printf("System Stability: %s\n", stability)
}

func percentage(part, total int) float64 {
	if total == 0 {
		return 0
	}
	return (float64(part) / float64(total)) * 100
}

func computeStability(p float64) string {
	if p < 25 {
		return "Very Good"
	} else if p < 50 {
		return "Good"
	} else {
		return "Needs Attention"
	}
}

func startRESTServer() {
	router := mux.NewRouter()
	router.HandleFunc("/logs/{context}/{severity}", getLogsHandler).Methods("GET")
	http.Handle("/", router)
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func main() {
	rand.Seed(time.Now().UnixNano())
	initDB()

	go func() {
		for range time.Tick(time.Second) {
			newError := generateRandomError()
			mu.Lock()
			logErrorToSQLite(newError)
			mu.Unlock()
		}
	}()

	go printHealthAndStability()

	startRESTServer()
}
