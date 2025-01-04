package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sync"
	"time"
)

var logFile *os.File
var wg sync.WaitGroup

func init() {
	logFilePath := filepath.Join(".", "file_history.log")
	var err error
	logFile, err = os.OpenFile(logFilePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("Failed to open log file: %v", err)
	}
}

func logActivity(activity string, details string) {
	entry := fmt.Sprintf("[%v] %v: %v\n", time.Now().Format(time.RFC3339), activity, details)
	logFile.WriteString(entry)
	logFile.Sync()
}

func navigateToDirectory(dir string) {
	wg.Add(1)
	go func() {
		defer wg.Done()
		logActivity("Navigation", fmt.Sprintf("Navigated to %v", dir))
	}()
}

func createFile(path string) {
	wg.Add(1)
	go func() {
		defer wg.Done()
		logActivity("File Creation", fmt.Sprintf("Created file %v", path))
	}()
}

func modifyFile(path string) {
	wg.Add(1)
	go func() {
		defer wg.Done()
		logActivity("File Modification", fmt.Sprintf("Modified file %v", path))
	}()
}

func deleteFile(path string) {
	wg.Add(1)
	go func() {
		defer wg.Done()
		logActivity("File Deletion", fmt.Sprintf("Deleted file %v", path))
	}()
}

func main() {
	navigateToDirectory("/home/user/docs")
	createFile("new_file.txt")
	modifyFile("new_file.txt")
	deleteFile("new_file.txt")

	wg.Wait()
	logFile.Close()
}
