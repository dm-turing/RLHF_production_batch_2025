package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

var logFile *os.File
var wg sync.WaitGroup
var actionHistory []string
var actionIndex int

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
		actionHistory = append(actionHistory, fmt.Sprintf("Navigate to %v", dir))
		actionIndex++
		logActivity("Navigation", fmt.Sprintf("Navigated to %v", dir))
	}()
}

func createFile(path string) {
	wg.Add(1)
	go func() {
		defer wg.Done()
		action := fmt.Sprintf("Create file %v", path)
		actionHistory = append(actionHistory, action)
		actionIndex++
		logActivity("File Creation", fmt.Sprintf("Created file %v", path))
		os.Create(path)
	}()
}

func modifyFile(path string) {
	wg.Add(1)
	go func() {
		defer wg.Done()
		actionHistory = append(actionHistory, fmt.Sprintf("Modify file %v", path))
		actionIndex++
		logActivity("File Modification", fmt.Sprintf("Modified file %v", path))
		// Here you might want to read the file, apply a change, and write it back.
	}()
}

func deleteFile(path string) {
	wg.Add(1)
	go func() {
		defer wg.Done()
		action := fmt.Sprintf("Delete file %v", path)
		actionHistory = append(actionHistory, action)
		actionIndex++
		logActivity("File Deletion", fmt.Sprintf("Deleted file %v", path))
		os.Remove(path)
	}()
}

func undoAction() {
	if actionIndex > 0 {
		actionIndex--
		action := actionHistory[actionIndex]
		logActivity("Undo", fmt.Sprintf("Undoing: %v", action))
		if strings.HasPrefix(action, "Create file ") {
			filePath := strings.TrimPrefix(action, "Create file ")
			os.Remove(filePath)
			logActivity("File Deletion", fmt.Sprintf("Deleted file %v", filePath))
		} else if strings.HasPrefix(action, "Modify file ") {
			// Implement undo for modification
		} else if strings.HasPrefix(action, "Delete file ") {
			filePath := strings.TrimPrefix(action, "Delete file ")
			// Create an empty file for demonstration purposes
			os.Create(filePath)
			logActivity("File Creation", fmt.Sprintf("Created file %v", filePath))
		}
	}
}

func redoAction() {
	if actionIndex < len(actionHistory)-1 {
		actionIndex++
		action := actionHistory[actionIndex]
		logActivity("Redo", fmt.Sprintf("Redoing: %v", action))
		if strings.HasPrefix(action, "Create file ") {
			filePath := strings.TrimPrefix(action, "Create file ")
			os.Create(filePath)
			logActivity("File Creation", fmt.Sprintf("Created file %v", filePath))
		} else if strings.HasPrefix(action, "Modify file ") {
			// Implement redo for modification
		} else if strings.HasPrefix(action, "Delete file ") {
			filePath := strings.TrimPrefix(action, "Delete file ")
			os.Remove(filePath)
			logActivity("File Deletion", fmt.Sprintf("Deleted file %v", filePath))
		}
	}
}

func displayNavigationPath() {
	fmt.Println("Navigation Path:")
	for i, action := range actionHistory {
		if strings.HasPrefix(action, "Navigate to ") {
			dir := strings.TrimPrefix(action, "Navigate to ")
			fmt.Printf("%d. %v\n", i+1, dir)
		}
	}
}

func main() {
	navigateToDirectory("/home/user/docs")
	createFile("new_file.txt")
	modifyFile("new_file.txt") // Modification stub
	deleteFile("new_file.txt")

	displayNavigationPath()
	undoAction() // Undo delete
	redoAction() // Redo delete
	undoAction() // Undo delete again (this should recreate the file)

	wg.Wait()
	logFile.Close()
}
