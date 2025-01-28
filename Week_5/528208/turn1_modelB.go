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

const (
	dbFile         = "data.db"
	workerPoolSize = 10
)

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: go run main.go <folder_path>")
		os.Exit(1)
	}

	folderPath := os.Args[1]

	// Create the database and table
	db, err := sql.Open("sqlite3", dbFile)
	if err != nil {
		log.Fatalf("Error opening database: %v", err)
	}
	defer db.Close()

	if _, err := db.Exec("CREATE TABLE IF NOT EXISTS files (id INTEGER PRIMARY KEY, filename TEXT, content TEXT)"); err != nil {
		log.Fatalf("Error creating table: %v", err)
	}

	// Create a channel to send file paths to workers
	filePaths := make(chan string, 100)

	// Create a wait group to synchronize workers
	var wg sync.WaitGroup

	// Start the worker pool
	wg.Add(workerPoolSize)
	for i := 0; i < workerPoolSize; i++ {
		go worker(db, filePaths, &wg)
	}

	// Walk through the folder and send file paths to the workers
	err = filepath.Walk(folderPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() {
			filePaths <- path
		}
		return nil
	})
	if err != nil {
		log.Fatalf("Error walking through folder: %v", err)
	}

	// Close the file paths channel to signal workers to finish
	close(filePaths)

	// Wait for all workers to complete
	wg.Wait()

	fmt.Println("Database populated successfully!")
}

func worker(db *sql.DB, filePaths <-chan string, wg *sync.WaitGroup) {
	defer wg.Done()

	for filePath := range filePaths {
		func() {
			defer func() {
				if r := recover(); r != nil {
					log.Printf("Worker recovered from panic: %v", r)
				}
			}()

			// Read the file content
			content, err := ioutil.ReadFile(filePath)
			if err != nil {
				log.Printf("Error reading file %s: %v", filePath, err)
				return
			}

			// Extract the filename from the path
			filename := filepath.Base(filePath)

			// Insert the file data into the database
			_, err = db.Exec("INSERT INTO files (filename, content) VALUES (?, ?)", filename, string(content))
			if err != nil {
				log.Printf("Error inserting file %s into database: %v", filename, err)
			}

			fmt.Printf("Processed file: %s\n", filename)
		}()
	}
}
