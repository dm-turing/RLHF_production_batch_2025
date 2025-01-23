package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"time"

	_ "github.com/ncruces/go-sqlite3/driver"
	_ "github.com/ncruces/go-sqlite3/embed"
)

type RequestData struct {
	RequestID   string
	URL         string
	Method      string
	GeoLocation string
	UserAgent   string
}

func initDB() *sql.DB {
	db, err := sql.Open("sqlite3", "./dataEnrichment.db")
	if err != nil {
		log.Fatal(err)
	}

	createTableSQL := `CREATE TABLE IF NOT EXISTS requests (
		"request_id" TEXT NOT NULL PRIMARY KEY,
		"url" TEXT,
		"method" TEXT,
		"geo_location" TEXT,
		"user_agent" TEXT
	  );`

	_, err = db.Exec(createTableSQL)
	if err != nil {
		log.Fatal(err)
	}

	return db
}

func generateRequestID() string {
	return fmt.Sprintf("%d", rand.Intn(100000))
}

func generateRandomGeoLocation() string {
	// Returns latitude, longitude format for simplicity
	return fmt.Sprintf("%f, %f", rand.Float64()*180-90, rand.Float64()*360-180)
}

func generateRandomUserAgent() string {
	userAgents := []string{
		"Mozilla/5.0",
		"Chrome/91.0",
		"Safari/537.36",
		"Opera/9.80",
	}
	return userAgents[rand.Intn(len(userAgents))]
}

func handler(w http.ResponseWriter, r *http.Request) {
	rand.Seed(time.Now().UnixNano())
	requestData := RequestData{
		RequestID:   generateRequestID(),
		URL:         r.URL.String(),
		Method:      r.Method,
		GeoLocation: generateRandomGeoLocation(),
		UserAgent:   generateRandomUserAgent(),
	}

	// Convert the data to JSON and send as response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(requestData)

	// Insert the record into the database
	insertRequestData(initDB(), requestData)
}

func insertRequestData(db *sql.DB, reqData RequestData) {
	insertSQL := `INSERT INTO requests (request_id, url, method, geo_location, user_agent) VALUES (?, ?, ?, ?, ?)`
	statement, err := db.Prepare(insertSQL)
	if err != nil {
		log.Fatal(err)
	}
	defer statement.Close()

	_, err = statement.Exec(reqData.RequestID, reqData.URL, reqData.Method, reqData.GeoLocation, reqData.UserAgent)
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	http.HandleFunc("/", handler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
