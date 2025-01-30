package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"sync"
	"sync/atomic"
	"time"
)

const workerCount = 10

var destDir string

func main() {
	if len(os.Args) != 3 {
		fmt.Println("Usage: go run main.go <source_dir> <dest_dir>")
		os.Exit(1)
	}
	sourceDir := os.Args[1]
	destDir = os.Args[2]
	// Create the destination directory if it doesn't exist
	if err := os.MkdirAll(destDir, 0755); err != nil {
		log.Fatalf("Error creating destination directory: %v", err)
	}
	// Start the timer
	startTime := time.Now()

	// Channel to send file paths to workers
	filePaths := make(chan string, 1000)
	// WaitGroup to synchronize the main goroutine with workers
	var wg sync.WaitGroup
	// Counter to keep track of the number of files copied
	var copiedFiles uint64

	// Launch worker pool
	wg.Add(workerCount)
	for i := 0; i < workerCount; i++ {
		go worker(filePaths, &copiedFiles, &wg)
	}

	// Walk the source directory and send file paths to the workers
	err := filepath.Walk(sourceDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			filePaths <- path
		}
		return nil
	})
	if err != nil {
		log.Fatalf("Error walking source directory: %v", err)
	}

	// Close the channel to signal workers to stop
	close(filePaths)

	// Wait for all workers to finish
	wg.Wait()

	// Stop the timer and calculate the elapsed time
	elapsedTime := time.Since(startTime)
	fmt.Printf("Copied %d files in %s\n", copiedFiles, elapsedTime)
}

func worker(filePaths <-chan string, copiedFiles *uint64, wg *sync.WaitGroup) {
	defer wg.Done()
	for path := range filePaths {
		// Copy the file
		if err := copyFile(path, filepath.Join(destDir, filepath.Base(path))); err != nil {
			log.Printf("Error copying file %s: %v", path, err)
		} else {
			atomic.AddUint64(copiedFiles, 1)
		}
	}
}

func copyFile(src, dst string) error {
	srcFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer srcFile.Close()

	dstFile, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer dstFile.Close()

	_, err = io.Copy(dstFile, srcFile)
	return err
}
