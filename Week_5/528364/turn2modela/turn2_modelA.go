package main

import (
	"fmt"
	"math/rand"
	"os"
	"path/filepath"
	"time"
)

// Total allowable size of our "file system" (1 MB)
const totalSize = 1 * 1024 * 1024 // 1 MB

// Folder name for our simulated file system
const folderName = "virtual_fs"
const fileName = "data.txt"

// Function to create a new folder
func createFolder(name string) error {
	return os.MkdirAll(name, os.ModePerm)
}

// Function to calculate the size of all files in the folder
func getFolderSize(folder string) (int64, error) {
	var size int64
	err := filepath.Walk(folder, func(_ string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		size += info.Size()
		return nil
	})
	return size, err
}

// Function to generate random text
func generateRandomText(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[rand.Intn(len(charset))]
	}
	return string(b)
}

// Main function to handle writing and monitoring
func main() {
	rand.Seed(time.Now().UnixNano())

	// Channel to communicate between the main application and the watcher process
	statusChannel := make(chan float64)
	done := make(chan bool)

	// Start the watcher process
	go watcher(statusChannel, done)

	// Create the folder for our "file system"
	if err := createFolder(folderName); err != nil {
		fmt.Println("Error creating folder:", err)
		return
	}

	// Open a file to write randomized data
	file, err := os.Create(filepath.Join(folderName, fileName))
	if err != nil {
		fmt.Println("Error creating file:", err)
		return
	}
	defer file.Close()

writeLoop:
	for {
		select {
		case <-done:
			fmt.Println("Watcher requested shutdown. Ending process gracefully.")
			break writeLoop
		default:
			// Write random data to the file
			data := generateRandomText(1024) // Write 1KB at a time
			if _, err := file.WriteString(data); err != nil {
				fmt.Println("Error writing to file:", err)
				break
			}

			// Check the folder size
			currentSize, err := getFolderSize(folderName)
			if err != nil {
				fmt.Println("Error getting folder size:", err)
				break
			}

			currentUsage := float64(currentSize) / float64(totalSize) * 100
			fmt.Printf("Current Size: %.2f%% (%d bytes)\n", currentUsage, currentSize)

			// Send current usage to the watcher
			statusChannel <- currentUsage

			if currentUsage >= 99 {
				fmt.Println("Reached 99% of the capacity. Ending process gracefully.")
				break writeLoop
			}
		}
	}

	close(statusChannel)
	fmt.Println("Final folder size monitoring complete.")
}

// Watcher function to monitor the file system usage
func watcher(statusChannel chan float64, done chan bool) {
	var lastPrintedThreshold float64
	for usage := range statusChannel {
		switch {
		case usage >= 90 && lastPrintedThreshold < 90:
			fmt.Println("Watcher: Reached 90% of the capacity.")
			lastPrintedThreshold = 90
		case usage >= 75 && lastPrintedThreshold < 75:
			fmt.Println("Watcher: Reached 75% of the capacity.")
			lastPrintedThreshold = 75
		case usage >= 50 && lastPrintedThreshold < 50:
			fmt.Println("Watcher: Reached 50% of the capacity.")
			lastPrintedThreshold = 50
		}

		// Example condition to interrupt the process (could be any logic or external trigger)
		if usage >= 95 {
			fmt.Println("Watcher: Requesting shutdown due to high usage.")
			done <- true
			return
		}
	}
}
