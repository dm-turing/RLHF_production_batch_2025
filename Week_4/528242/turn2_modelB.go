package main  
import (  
    "database/sql"
    "fmt"
    "log"
    "math/rand"
    "net/http"
    "strings"
    "time"

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
    ID        int
    Severity Severity
    Message  string
    Context  string
    CreatedAt time.Time
}
func initDatabase() *sql.DB {
    db, err := sql.Open("sqlite3", "./errors.db")
    if err != nil {
        log.Fatalf("Error opening SQLite database: %v", err)
    }

    // Create the errors table if it doesn't exist
    createTableSQL := `
        CREATE TABLE IF NOT EXISTS errors (
            id INTEGER PRIMARY KEY AUTOINCREMENT,
            severity INTEGER,
            message TEXT,
            context TEXT,
            created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
        )
    `
    _, err = db.Exec(createTableSQL)
    if err != nil {
        log.Fatalf("Error creating errors table: %v", err)
    }

    return db
}
func generateRandomErrorMessage() ErrorMessage {
    severities := []Severity{DEBUG, INFO, WARNING, ERROR, CRITICAL, FATAL}
    messages := []string{ 
        "Disk space low",        
        "Network latency detected", 
        "Application crashed",
        "User input error",
        "Configuration file missing",      
    }
 
    contextOptions := []string{       
        "on server A",
        "at node Y",
        "in region B",
        "with IP 192.168.1.1",
        "during task Z",
        "with node ID 3",  
    }
  
    randomSeverity := severities[rand.Intn(len(severities))]  
    randomMessage := messages[rand.Intn(len(messages))]    
    randomContext := contextOptions[rand.Intn(len(contextOptions))]    
  
    return ErrorMessage{ 
        Severity: randomSeverity,       
        Message:  randomMessage, 
        Context:  randomContext,      
    } 
} 

func addMetadata(err ErrorMessage) ErrorMessage {
    // No metadata is added in this updated function
    return err
}
func storeErrorInDB(db *sql.DB, err ErrorMessage) {  
    insertSQL := `INSERT INTO errors (severity, message, context) VALUES (?, ?, ?)`
    _, err = db.Exec(insertSQL, err.Severity, err.Message, err.Context)  
    if err != nil {  
        log.Printf("Error storing error in database: %v", err)
    } 
}
func handleErrors(db *sql.DB) {  
    rand.Seed(time.Now().UnixNano())      
    for {   
        errorMessage := generateRandomErrorMessage() 
        errorMessage = addMetadata(errorMessage)  
        // Reduce the frequency of CRITICAL or FATAL errors
        if errorMessage.Severity != CRITICAL && errorMessage.Severity != FATAL {
            // Ignore and discard CRITICAL or FATAL errors as per the requirement
            storeErrorInDB(db, errorMessage)
        }
        time.Sleep(time.Duration(rand.Intn(5)+1) * time.Second) // Sleep for 1-5 seconds
    } 
} 
func fetchErrorsBySeverity(db *sql.DB, severity string) ([]ErrorMessage, error) {
    var errors []ErrorMessage