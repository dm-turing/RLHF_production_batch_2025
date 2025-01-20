package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sync"

	// _ "github.com/mattn/go-sqlite3"
	_ "github.com/ncruces/go-sqlite3/driver"
	_ "github.com/ncruces/go-sqlite3/embed"
)

// FileMetadata holds the basic metadata about a file
type FileMetadata struct {
	Path string
	Size int64
}

// initializeDatabase sets up the SQLite database and creates the required table
func initializeDatabase(db *sql.DB) error {
	_, err := db.Exec(`CREATE TABLE IF NOT EXISTS file_metadata (
		id INTEGER PRIMARY KEY,
		path TEXT,
		size INTEGER
	)`)
	return err
}

// saveMetadata stores the file metadata in the SQLite database
func saveMetadata(db *sql.DB, metadata FileMetadata) error {
	_, err := db.Exec(`INSERT INTO file_metadata (path, size) VALUES (?, ?)`, metadata.Path, metadata.Size)
	return err
}

// processFile gathers metadata from the given file path
func processFile(filePath string) (FileMetadata, error) {
	fileInfo, err := os.Stat(filePath)
	if err != nil {
		return FileMetadata{}, err
	}
	return FileMetadata{Path: filePath, Size: fileInfo.Size()}, nil
}

// worker is the function that each goroutine will execute
func worker(files <-chan string, db *sql.DB, wg *sync.WaitGroup) {
	defer wg.Done()

	for filePath := range files {
		metadata, err := processFile(filePath)
		if err != nil {
			log.Printf("Failed to process file %s: %v", filePath, err)
			continue
		}

		if err := saveMetadata(db, metadata); err != nil {
			log.Printf("Failed to save metadata for file %s: %v", filePath, err)
		}
	}
}

// processFiles concurrently processes files and stores their metadata
func processFiles(db *sql.DB, files []string) {
	fileChan := make(chan string)
	var wg sync.WaitGroup

	numWorkers := 10 // Adjust the number of workers as needed
	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go worker(fileChan, db, &wg)
	}

	for _, filePath := range files {
		fileChan <- filePath
	}

	close(fileChan)
	wg.Wait()
}

func main() {
	db, err := sql.Open("sqlite3", "file_metadata.db")
	if err != nil {
		log.Fatalf("Failed to open the database: %v", err)
	}
	defer db.Close()

	if err := initializeDatabase(db); err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}

	var files []string
	err = filepath.Walk("allFiles", func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			files = append(files, path)
		}
		return nil
	})

	if err != nil {
		log.Fatalf("Error walking the file path: %v", err)
	}

	processFiles(db, files)

	fmt.Println("Metadata processing completed.")
}
