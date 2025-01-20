package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sync"
	"time"

	// _ "github.com/mattn/go-sqlite3"
	_ "github.com/ncruces/go-sqlite3/driver"
	_ "github.com/ncruces/go-sqlite3/embed"
)

var numWorkers = 3

// Define your Metadata struct
type Metadata struct {
	Size     int64
	ModTime  time.Time
	Filename string
	// Add other metadata fields as needed
}

func extractMetadata(filePath string) (*Metadata, error) {
	// Implement your metadata extraction logic here
	return &Metadata{}, nil
}

func storeMetadata(db *sql.DB, filePath string, metadata *Metadata) error {
	// Implement your SQLite database insertion logic here
	stmt, err := db.Prepare("INSERT INTO metadata(filename, size, modtime) VALUES (?, ?, ?)")
	if err != nil {
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(filepath.Base(filePath), metadata.Size, metadata.ModTime)
	return err
}

func main() {
	var wg sync.WaitGroup
	filePaths := make(chan string, 100) // Buffered channel with capacity 100
	var allFiles []string
	filepath.Walk("allFiles", func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			allFiles = append(allFiles, path)
		}
		return nil
	})

	// Spawn worker goroutines
	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go worker(i, filePaths, &wg)
	}

	// Generate and send file paths to the channel
	go func() {
		for _, filePath := range allFiles {
			filePaths <- filePath
		}
		close(filePaths) // Close the channel once all file paths are sent
	}()

	// Wait for all workers to complete
	wg.Wait()
	fmt.Println("All files processed.")
}

func worker(id int, filePaths chan string, wg *sync.WaitGroup) {
	defer wg.Done()
	db, err := sql.Open("sqlite3", "./metadata.db")
	if err != nil {
		log.Fatalf("Error opening database: %v", err)
	}
	defer db.Close()

	for filePath := range filePaths {
		metadata, err := extractMetadata(filePath)
		if err != nil {
			log.Printf("Error extracting metadata for %s: %v", filePath, err)
			continue
		}
		err = storeMetadata(db, filePath, metadata)
		if err != nil {
			log.Printf("Error storing metadata for %s: %v", filePath, err)
		}
		fmt.Printf("Worker %d processed %s\n", id, filePath)
	}
}
