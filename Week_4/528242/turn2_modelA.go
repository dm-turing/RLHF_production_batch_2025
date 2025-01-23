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

type Severity int

const (
	DEBUG Severity = iota
	INFO
	WARNING
	ERROR
	CRITICAL
	FATAL
)

func (s Severity) String() string {
	return []string{"DEBUG", "INFO", "WARNING", "ERROR", "CRITICAL", "FATAL"}[s]
}

type ErrorMessage struct {
	Severity Severity
	Message  string
	Metadata string
}

// Store the error message in the SQLite database.
func storeError(db *sql.DB, err ErrorMessage) {
	_, dbErr := db.Exec("INSERT INTO errors(severity, message, metadata, created_at) VALUES (?, ?, ?, datetime('now'))", err.Severity.String(), err.Message, err.Metadata)
	if dbErr != nil {
		log.Printf("Error storing message: %v\n", dbErr)
	}
}

func generateRandomErrorMessage() ErrorMessage {
	severities := []Severity{DEBUG, INFO, DEBUG, INFO, WARNING, ERROR, WARNING} // increased lower severity frequency
	messages := []string{
		"Disk space low",
		"Network latency detected",
		"Application crashed",
		"User input error",
		"Configuration file missing",
	}

	randomSeverity := severities[rand.Intn(len(severities))]
	randomMessage := messages[rand.Intn(len(messages))]

	// Add a reduced probability for CRITICAL or FATAL errors
	if rand.Intn(100) < 5 { // 5% chance for CRITICAL or less
		randomSeverity = CRITICAL
	}
	if rand.Intn(100) < 2 { // 2% chance for FATAL or less
		randomSeverity = FATAL
	}

	return ErrorMessage{
		Severity: randomSeverity,
		Message:  randomMessage,
	}
}

func addMetadata(err ErrorMessage) ErrorMessage {
	// More dynamic context information as metadata
	contextTypes := []string{"user", "server", "session", "process"}
	contextValues := []string{"X", "A", "12345", "backup"}

	contextType := contextTypes[rand.Intn(len(contextTypes))]
	contextValue := contextValues[rand.Intn(len(contextValues))]

	err.Metadata = fmt.Sprintf("%s-%s", contextType, contextValue)
	return err
}

func handleErrors(db *sql.DB) {
	// Timer tasks setup
	errorChannel := make(chan ErrorMessage, 1)
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()
	exit := make(chan bool)

	go func() {
		for {
			errorMessage := generateRandomErrorMessage()
			errorMessage = addMetadata(errorMessage)
			storeError(db, errorMessage)
			errorChannel <- errorMessage
			time.Sleep(1 * time.Second)
		}
	}()

	go func() {
		infoErrors := []ErrorMessage{}
		for {
			select {
			case err := <-errorChannel:
				switch err.Severity {
				case DEBUG, INFO:
					infoErrors = append(infoErrors, err)
				case WARNING, ERROR:
					fmt.Printf("WARNING/ERROR: [%s] %s - %s\n", err.Severity, err.Message, err.Metadata)
				case CRITICAL, FATAL:
					fmt.Printf("CRITICAL/FATAL: [%s] %s - %s\n", err.Severity, err.Message, err.Metadata)
					exit <- true
					return
				}

			case <-ticker.C:
				for _, err := range infoErrors {
					fmt.Printf("DEBUG/INFO: [%s] %s - %s\n", err.Severity, err.Message, err.Metadata)
				}
				infoErrors = nil
			}
		}
	}()

	<-exit
	fmt.Println("Session ended due to a CRITICAL or FATAL error.")
}

func getErrorsBySeverity(w http.ResponseWriter, r *http.Request) {
	severityParam := r.URL.Query().Get("severity")

	db, err := sql.Open("sqlite3", "./errors.db")
	if err != nil {
		http.Error(w, "Database connection error", http.StatusInternalServerError)
		return
	}
	defer db.Close()

	rows, err := db.Query(`SELECT severity, message, metadata, created_at FROM errors WHERE severity = ? ORDER BY created_at ASC`, severityParam)
	if err != nil {
		http.Error(w, "Query execution failed", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var errors []ErrorMessage
	for rows.Next() {
		var severity, message, metadata, createdAt string
		rows.Scan(&severity, &message, &metadata, &createdAt)
		fmt.Println(severity, message, metadata) // In debugging it's a great idea to print

		errors = append(errors, ErrorMessage{
			Severity: STRING_TO_SEVERITY[severity],
			Message:  fmt.Sprintf("%s (%s)", message, createdAt),
			Metadata: metadata,
		})
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(errors)
}

func setupDatabase() (*sql.DB, error) {
	db, err := sql.Open("sqlite3", "./errors.db")
	if err != nil {
		return nil, fmt.Errorf("error opening database: %w", err)
	}

	statement, err := db.Prepare(`CREATE TABLE IF NOT EXISTS errors (id INTEGER PRIMARY KEY, severity TEXT, message TEXT, metadata TEXT, created_at TIMESTAMP)`)
	if err != nil {
		return nil, fmt.Errorf("error preparing the create table statement: %w", err)
	}

	_, err = statement.Exec()
	if err != nil {
		return nil, fmt.Errorf("error executing the create table statement: %w", err)
	}

	return db, nil
}

var STRING_TO_SEVERITY = map[string]Severity{
	"DEBUG":    DEBUG,
	"INFO":     INFO,
	"WARNING":  WARNING,
	"ERROR":    ERROR,
	"CRITICAL": CRITICAL,
	"FATAL":    FATAL,
}

func main() {
	db, err := setupDatabase()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	go handleErrors(db)

	r := mux.NewRouter()
	r.HandleFunc("/errors", getErrorsBySeverity).Methods("GET")

	http.Handle("/", r)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
