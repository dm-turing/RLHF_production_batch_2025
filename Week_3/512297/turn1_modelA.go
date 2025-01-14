package main

import (
	"database/sql"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"strings"
	"time"

	_ "github.com/ncruces/go-sqlite3/driver"
	_ "github.com/ncruces/go-sqlite3/embed"
)

var db *sql.DB

func init() {
	var err error
	db, err = sql.Open("sqlite3", "metrics.db")
	if err != nil {
		log.Fatalf("Error opening database: %v", err)
	}
	defer db.Close()

	_, err = db.Exec(`
	CREATE TABLE IF NOT EXISTS metrics (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		parameter TEXT UNIQUE,
		count INTEGER,
		frequency INTEGER,
		error_rate REAL,
		avg_length INTEGER,
		avg_response_time INTEGER,
		server_load REAL
	);
	`)
	if err != nil {
		log.Fatalf("Error creating table: %v", err)
	}
}

func handleLogRequest(w http.ResponseWriter, r *http.Request) {
	db, err := sql.Open("sqlite3", "metrics.db")
	if err != nil {
		log.Fatalf("Error opening database: %v", err)
	}
	defer db.Close()

	params := r.URL.Query()
	if len(params) == 0 {
		http.Error(w, "No query parameters provided", http.StatusBadRequest)
		return
	}

	for param, values := range params {
		for _, value := range values {
			count, err := db.Exec("INSERT INTO metrics (parameter, count, frequency, error_rate, avg_length, avg_response_time, server_load) VALUES (?, 1, 1, 0.0, ?, 0, 0.0) ON CONFLICT (parameter) DO UPDATE SET count = count + 1, frequency = frequency + 1, avg_length = (avg_length * (frequency - 1) + ?) / frequency", param, len(value), len(value))
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			_, err = count.LastInsertId()
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		}
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Metrics logged successfully\n")
}

func handleGetMetrics(w http.ResponseWriter, r *http.Request) {
	db, err := sql.Open("sqlite3", "metrics.db")
	if err != nil {
		log.Fatalf("Error opening database: %v", err)
	}
	defer db.Close()

	rows, err := db.Query("SELECT parameter, count, frequency, error_rate, avg_length, avg_response_time, server_load FROM metrics")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	w.WriteHeader(http.StatusOK)
	for rows.Next() {
		var parameter, count, frequency, errorRate, avgLength, avgResponseTime, serverLoad string
		err = rows.Scan(&parameter, &count, &frequency, &errorRate, &avgLength, &avgResponseTime, &serverLoad)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		fmt.Fprintf(w, "Parameter: %s, Count: %s, Frequency: %s, Error Rate: %s, Avg Length: %s, Avg Response Time: %s, Server Load: %s\n", parameter, count, frequency, errorRate, avgLength, avgResponseTime, serverLoad)
	}
}

func populateDatabase() {
	var db *sql.DB
	var err error
	db, err = sql.Open("sqlite3", "metrics.db")
	if err != nil {
		log.Fatalf("Error opening database: %v", err)
	}
	defer db.Close()

	params := []string{"param1", "param2", "param3", "param4", "param5"}

	for i := 0; i < 100; i++ {
		param := params[rand.Intn(len(params))]
		length := rand.Intn(100) + 1
		value := strings.Repeat(strconv.Itoa(rand.Intn(1000)), length)

		_, err = db.Exec("INSERT INTO metrics (parameter, count, frequency, error_rate, avg_length, avg_response_time, server_load) VALUES (?, 1, 1, 0.0, ?, 0, 0.0) ON CONFLICT (parameter) DO UPDATE SET count = count + 1, frequency = frequency + 1, avg_length = (avg_length * (frequency - 1) + ?) / frequency", param, len(value), len(value))
		if err != nil {
			log.Fatalf("Error inserting data: %v", err)
		}
	}
}

func main() {
	rand.Seed(time.Now().UnixNano())
	populateDatabase()
	http.HandleFunc("/log-request", handleLogRequest)
	http.HandleFunc("/get-metrics", handleGetMetrics)
	log.Println("Server listening on port 8080")
	if err := http.ListenAndServe(":8080", http.DefaultServeMux); err != nil {
		log.Fatal(err)
	}
}
