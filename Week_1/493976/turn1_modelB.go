package main  
import (  
    "encoding/json"
    "fmt"
    "log"
    "net/http"
    "strconv"
    "time"
    "github.com/gorilla/mux"
    "gorm.io/driver/sqlite"
    "gorm.io/gorm"
)

// Query represents a query with its associated parameters
type Query struct {
    ID           uint       `json:"id" gorm:"primaryKey"`
    CreatedAt    time.Time  `json:"created_at"`
    UpdatedAt    time.Time  `json:"updated_at"`
    DeletedAt    gorm.DeletedAt `json:"-" gorm:"index"`
    Parameters   string     `json:"parameters"`
}

// AuditTrail represents a change in query parameters
type AuditTrail struct {
    ID        uint       `json:"id" gorm:"primaryKey"`
    QueryID   uint       `json:"query_id"`
    CreatedAt time.Time  `json:"created_at"`
    Parameters string     `json:"parameters"`
}

func main() {
    db, err := gorm.Open(sqlite.Open("audit_trail.db"), &gorm.Config{})
    if err != nil {
        log.Fatalf("Error opening database: %v", err)
    }
    defer db.Close()

    // Create tables if they don't exist
    db.AutoMigrate(&Query{}, &AuditTrail{})

    r := mux.NewRouter()
    r.HandleFunc("/query", createQuery(db)).Methods("POST")
    r.HandleFunc("/query/{queryId}", getQuery(db)).Methods("GET")
    r.HandleFunc("/query/{queryId}/history", getQueryHistory(db)).Methods("GET")

    fmt.Println("Server running on http://localhost:8080")
    log.Fatal(http.ListenAndServe(":8080", r))
}

// createQuery creates a new query or updates an existing one based on the provided ID
func createQuery(db *gorm.DB) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        var query Query
        if err := json.NewDecoder(r.Body).Decode(&query); err != nil {
            http.Error(w, err.Error(), http.StatusBadRequest)
            return
        }

        defer r.Body.Close()

        // Find the query if it exists
        var existingQuery Query
        if result := db.First(&existingQuery, query.ID); result.Error == nil {
            // Update query parameters if they differ
            if existingQuery.Parameters != query.Parameters {
                if err := db.Model(&existingQuery).Update("parameters", query.Parameters).Error; err != nil {
                    http.Error(w, err.Error(), http.StatusInternalServerError)
                    return
                }

                // Create a new audit trail entry
                if err := db.Create(&AuditTrail{QueryID: existingQuery.ID, Parameters: query.Parameters}).Error; err != nil {
                    http.Error(w, err.Error(), http.StatusInternalServerError)
                    return
                }
            }
        } else {
            // Create a new query
            if err := db.Create(&query).Error; err != nil {
                http.Error(w, err.Error(), http.StatusInternalServerError)
                return
            }
        }

        json.NewEncoder(w).Encode(query)
    }