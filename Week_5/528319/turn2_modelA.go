package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"sync"
	"time"

	_ "github.com/ncruces/go-sqlite3/driver"
	_ "github.com/ncruces/go-sqlite3/embed"
)

// FileInfo struct holds file details
type FileInfo struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	Content   string    `json:"content,omitempty"`
	Timestamp time.Time `json:"timestamp"`
}

var db *sql.DB
var jobs = make(chan Job, 100)
var results = make(chan Result, 100)
var wg sync.WaitGroup

// Job struct to hold tasks for the worker pool
type Job struct {
	Action string
	File   FileInfo
}

// Result struct for processing results
type Result struct {
	Status  string
	Message string
}

func randomString(n int) string {
	const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(b)
}

func initializeDatabase() {
	var err error
	db, err = sql.Open("sqlite3", "./files.db")
	if err != nil {
		panic(err)
	}

	statement, _ := db.Prepare(
		`CREATE TABLE IF NOT EXISTS files (
			id INTEGER PRIMARY KEY,
			name TEXT,
			content TEXT,
			timestamp DATETIME DEFAULT CURRENT_TIMESTAMP
		)`)
	statement.Exec()
}

func insertFile(name, content string) (sql.Result, error) {
	statement, err := db.Prepare("INSERT INTO files (name, content) VALUES (?, ?)")
	if err != nil {
		return nil, err
	}
	return statement.Exec(name, content)
}

func deleteFile(name string) (sql.Result, error) {
	statement, err := db.Prepare("DELETE FROM files WHERE name = ?")
	if err != nil {
		return nil, err
	}
	return statement.Exec(name)
}

func listFiles() ([]FileInfo, error) {
	rows, err := db.Query("SELECT id, name, content, timestamp FROM files")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var files []FileInfo
	for rows.Next() {
		var file FileInfo
		err := rows.Scan(&file.ID, &file.Name, &file.Content, &file.Timestamp)
		if err != nil {
			return nil, err
		}
		files = append(files, file)
	}
	return files, nil
}

func fileOperations(job Job) {
	var result Result
	switch job.Action {
	case "create":
		content := randomString(1024) // 1KB of random text
		_, err := insertFile(job.File.Name, content)
		if err != nil {
			result = Result{Status: "error", Message: err.Error()}
		} else {
			result = Result{Status: "success", Message: "File created successfully"}
		}
	case "delete":
		_, err := deleteFile(job.File.Name)
		if err != nil {
			result = Result{Status: "error", Message: err.Error()}
		} else {
			result = Result{Status: "success", Message: "File deleted successfully"}
		}
	}
	results <- result
	wg.Done()
}

func workerPool(n int) {
	for i := 0; i < n; i++ {
		go func() {
			for job := range jobs {
				fileOperations(job)
			}
		}()
	}
}

func listFilesHandler(w http.ResponseWriter, r *http.Request) {
	files, err := listFiles()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(files)
}

func createFileHandler(w http.ResponseWriter, r *http.Request) {
	filename := r.URL.Query().Get("filename")
	if filename == "" {
		http.Error(w, "filename is required", http.StatusBadRequest)
		return
	}

	wg.Add(1)
	jobs <- Job{Action: "create", File: FileInfo{Name: filename}}
	wg.Wait()
	result := <-results

	if result.Status == "success" {
		w.WriteHeader(http.StatusCreated)
	} else {
		w.WriteHeader(http.StatusInternalServerError)
	}
	fmt.Fprintf(w, result.Message)
}

func deleteFileHandler(w http.ResponseWriter, r *http.Request) {
	filename := r.URL.Query().Get("filename")
	if filename == "" {
		http.Error(w, "filename is required", http.StatusBadRequest)
		return
	}

	wg.Add(1)
	jobs <- Job{Action: "delete", File: FileInfo{Name: filename}}
	wg.Wait()
	result := <-results

	if result.Status == "success" {
		w.WriteHeader(http.StatusOK)
	} else {
		w.WriteHeader(http.StatusInternalServerError)
	}
	fmt.Fprintf(w, result.Message)
}

func main() {
	initializeDatabase()

	rand.Seed(time.Now().UnixNano())

	workerPool(10) // Start worker pool with 10 workers

	http.HandleFunc("/files", listFilesHandler)
	http.HandleFunc("/create", createFileHandler)
	http.HandleFunc("/delete", deleteFileHandler)

	fmt.Println("Server starting on port 8080...")
	http.ListenAndServe(":8080", nil)
}
