package main  
import (  
    "encoding/json"
    "fmt"
    "log"
    "math/rand"
    "net/http"
    "os"
    "strings"
    "time"
)

type Severity int
const (
    INFO Severity = iota
    DEBUG
    WARNING
    ERROR
    FATAL
    CRITICAL
)

type LogEntry struct {
    Timestamp string    `json:"timestamp"`
    Severity  Severity  `json:"severity"`
    Message   string    `json:"message"`
    Path      string    `json:"path"`
    UserID    string    `json:"user_id"`
}

func (s Severity) String() string {
    return [...]string{"INFO", "DEBUG", "WARNING", "ERROR", "FATAL", "CRITICAL"}[s]
}

func generateRandomError() LogEntry {
    severity := Severity(rand.Intn(6))
    message := "This is a " + severity.String() + " level error."
    path := fmt.Sprintf("/api/v1/user/%d/action/%s", rand.Intn(100), randomString(10))
    userID := randomString(8)
    return LogEntry{
        Timestamp: time.Now().Format(time.RFC3339),
        Severity:  severity,
        Message:   message,
        Path:      path,
        UserID:    userID,
    }
}

func randomString(n int) string {
    var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")
    b := make([]rune, n)
    for i := range b {
        b[i] = letters[rand.Intn(len(letters))]
    }
    return string(b)
}

func logEntryToJSON(entry LogEntry) string {
    jsonBytes, err := json.Marshal(entry)
    if err != nil {
        log.Fatalf("Error converting log entry to JSON: %v", err)
    }
    return string(jsonBytes)
}

func LogError(entry LogEntry) {
    logMsg := logEntryToJSON(entry)
    if entry.Severity == DEBUG || entry.Severity == INFO {
        log.Println(logMsg)
    } else {
        // Log to a file for WARNING, ERROR, FATAL, and CRITICAL
        f, err := os.OpenFile("errors.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
        if err != nil {
            log.Fatalf("Error opening log file: %v", err)
        }
        defer f.Close()
        fileLogger := log.New(f, "", log.LstdFlags)
        fileLogger.Println(logMsg)
    }
}

func handleLogQuery(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodGet {
        http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
        return
    }

    severity, ok := r.URL.Query()["severity"]
    if !ok {
        http.Error(w, "Missing 'severity' parameter", http.StatusBadRequest)
        return
    }

    // Simple SQL-like query filtering using substring matching
    querySeverity := severity[0]
    filter := func(entry LogEntry) bool {
        return strings.Contains(entry.Severity.String(), querySeverity)
    }

    w.Header().Set("Content-Type", "application/json")
    if err := json.NewEncoder(w).Encode(readLogs(filter)); err != nil {
        http.Error(w, "Error processing query", http.StatusInternalServerError)
    }