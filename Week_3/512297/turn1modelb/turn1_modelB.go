package main

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	_ "github.com/ncruces/go-sqlite3/driver"
	_ "github.com/ncruces/go-sqlite3/embed"
)

// Metric types
const (
	DistinctParamCount = "distinct_param_count"
	ParamFrequency     = "param_frequency"
	ErrorRate          = "error_rate"
	ParamLength        = "param_length"
	ResponseTime       = "response_time"
	ServerLoad         = "server_load"
)

// Metric struct
type Metric struct {
	Type   string  `json:"type"`
	Value  float64 `json:"value"`
	Time   string  `json:"time"` // ISO 8601 format
	Params string  `json:"params"`
}

var db *sql.DB

func init() {
	// Initialize the database connection
	var err error
	db, err = sql.Open("sqlite3", "./metrics.db")
	if err != nil {
		log.Fatalf("Error opening database: %v", err)
	}
	// Create tables if they don't exist
	createTables()
}

func main() {
	// Initialize Gin router
	r := gin.Default()

	// Register API endpoints
	r.POST("/metrics", collectMetric)
	r.GET("/metrics", getMetrics)

	// Start the server
	log.Fatal(r.Run(":8080"))
}

// createTables creates the tables in the SQLite database
func createTables() {
	createTableSQL := `
	CREATE TABLE IF NOT EXISTS metrics (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		metric_type TEXT,
		value REAL,
		time TEXT,
		params TEXT
	);
	`
	_, err := db.Exec(createTableSQL)
	if err != nil {
		log.Fatalf("Error creating table: %v", err)
	}
}

// collectMetric handles POST requests to collect metrics
func collectMetric(c *gin.Context) {
	var metric Metric
	if err := c.ShouldBindJSON(&metric); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Validate metric type
	switch metric.Type {
	case DistinctParamCount, ParamFrequency, ErrorRate, ParamLength, ResponseTime, ServerLoad:
	default:
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid metric type"})
		return
	}

	// Insert metric into the database
	insertSQL := `
	INSERT INTO metrics (metric_type, value, time, params) 
	VALUES (?, ?, ?, ?)
	`
	_, err := db.Exec(insertSQL, metric.Type, metric.Value, metric.Time, metric.Params)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Metric collected successfully"})
}
