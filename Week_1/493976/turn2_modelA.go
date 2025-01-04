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
		actionHistory = append(actionHistory, fmt.Sprintf("Create file %v", path))
		actionIndex++
		logActivity("File Creation", fmt.Sprintf("Created file %v", path))
	}()
}

func modifyFile(path string) {
	wg.Add(1)
	go func() {
		defer wg.Done()
		actionHistory = append(actionHistory, fmt.Sprintf("Modify file %v", path))
		actionIndex++
		logActivity("File Modification", fmt.Sprintf("Modified file %v", path))
	}()
}

func deleteFile(path string) {
	wg.Add(1)
	go func() {
		defer wg.Done()
		actionHistory = append(actionHistory, fmt.Sprintf("Delete file %v", path))
		actionIndex++
		logActivity("File Deletion", fmt.Sprintf("Deleted file %v", path))
	}()
}

func undoAction() {
	if actionIndex > 0 {
		actionIndex--
		action := actionHistory[actionIndex]
		logActivity("Undo", fmt.Sprintf("Undoing: %v", action))
		// Implement the actual undo logic here
		// For example, if the action is "Create file %v", delete the file
		if strings.HasPrefix(action, "Create file ") {
			filePath := strings.TrimPrefix(action, "Create file ")
			os.Remove(filePath)
			logActivity("File Deletion", fmt.Sprintf("Deleted file %v", filePath))
		} else if strings.HasPrefix(action, "Modify file ") {
			// Implement undo for modification
		} else if strings.HasPrefix(action, "Delete file ") {
			// Implement undo for deletion
		}
	}
}

func redoAction() {
	if actionIndex < len(actionHistory)-1 {
		actionIndex++
		action := actionHistory[actionIndex]
		logActivity("Redo", fmt.Sprintf("Redoing: %v", action))
		// Implement the actual redo logic here
		// For example, if the action is "Create file %v", create the file
		if strings.HasPrefix(action, "Create file ") {
			filePath := strings.TrimPrefix(action, "Create file ")
			os.Create(filePath)
			logActivity("File Creation", fmt.Sprintf("Created file %v", filePath))
		} else if strings.HasPrefix(action, "Modify file ") {
			// Implement redo for modification
		} else if strings.HasPrefix(action, "Delete file ") {
			// Implement redo for deletion
		}
	}
}

func displayNavigationPath() {
	fmt.Println("Navigation Path:")
	for i, action := range actionHistory {
		fmt.Println(action)
		if strings.HasPrefix(action, "Navigate to ") {
			dir := strings.TrimPrefix(action, "Navigate to ")
			fmt.Printf("%d. %v\n", i+1, dir)
		}
	}
}

func main() {
	navigateToDirectory("/home/user/docs")
	createFile("new_file.txt")
	modifyFile("new_file.txt")
	deleteFile("new_file.txt")

	displayNavigationPath()
	undoAction()
	redoAction()

	wg.Wait()
	logFile.Close()
}
