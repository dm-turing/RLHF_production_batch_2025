package main

import (
	"database/sql"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"sync"

	_ "github.com/ncruces/go-sqlite3/driver"
	_ "github.com/ncruces/go-sqlite3/embed"
)

const workerCount = 10

type Job struct {
	Id   int
	Path string
}

func main() {
	if len(os.Args) != 3 {
		fmt.Println("Usage: mycli <folder_path> <database_name>")
		os.Exit(1)
	}

	folderPath, dbName := os.Args[1], os.Args[2]

	// Initialize SQLite database
	db, err := sql.Open("sqlite3", dbName)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	if err := initializeDB(db); err != nil {
		log.Fatal(err)
	}

	files, err := ioutil.ReadDir(folderPath)
	if err != nil {
		log.Fatalf("Failed to read directory: %s", err)
	}

	// Create jobs
	jobs := make(chan Job)
	var wg sync.WaitGroup

	// Start workers
	for i := 0; i < workerCount; i++ {
		wg.Add(1)
		go worker(i, db, jobs, &wg)
	}

	// Send jobs
	for id, file := range files {
		if !file.IsDir() {
			jobs <- Job{Id: id, Path: filepath.Join(folderPath, file.Name())}
		}
	}
	close(jobs)

	wg.Wait()
}

func initializeDB(db *sql.DB) error {
	query := `
    CREATE TABLE IF NOT EXISTS text_files (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        filename TEXT,
        content TEXT
    )`
	if _, err := db.Exec(query); err != nil {
		return fmt.Errorf("failed to create table: %w", err)
	}
	return nil
}

func worker(workerId int, db *sql.DB, jobs chan Job, wg *sync.WaitGroup) {
	defer wg.Done()
	for job := range jobs {
		fmt.Printf("Worker %d: Opening file handle %d (%s)\n", workerId, job.Id, job.Path)
		content, err := ioutil.ReadFile(job.Path)
		if err != nil {
			fmt.Printf("Worker %d: Error reading file %d: %s\n", workerId, job.Id, err)
			continue
		}

		if err := insertIntoDB(db, job.Path, string(content)); err != nil {
			fmt.Printf("Worker %d: Error writing file content %d to DB: %s\n", workerId, job.Id, err)
		}

		fmt.Printf("Worker %d: Closing file handle %d (%s)\n", workerId, job.Id, job.Path)
	}
}

func insertIntoDB(db *sql.DB, filename, content string) error {
	query := `
    INSERT INTO text_files (filename, content) VALUES (?, ?)`
	_, err := db.Exec(query, filename, content)
	return err
}
