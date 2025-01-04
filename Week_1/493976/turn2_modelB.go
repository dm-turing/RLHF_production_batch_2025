package main  

import (  
    "fmt"
    "log"
    "os"
    "path/filepath"
    "sync"
    "time"
)

type LogEntry struct {  
    Time     time.Time
    Activity string
    Details  string
}

var logFile *os.File
var wg sync.WaitGroup
var logEntries []LogEntry
var currentIndex int

func init() {  
    logFilePath := filepath.Join(".", "file_history.log")
    var err error
    logFile, err = os.OpenFile(logFilePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
    if err != nil {
        log.Fatalf("Failed to open log file: %v", err)
    }
    currentIndex = -1
}

func logActivity(activity string, details string) {  
    entry := LogEntry{Time: time.Now(), Activity: activity, Details: details}
    logEntries = append(logEntries, entry)
    currentIndex++
    writeLog(entry)
}

func writeLog(entry LogEntry) {
    entryString := fmt.Sprintf("[%v] %v: %v\n", entry.Time.Format(time.RFC3339), entry.Activity, entry.Details)  
    logFile.WriteString(entryString)
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

func undo() {  
    if currentIndex < 0 {  
        fmt.Println("No actions to undo.")
        return
    }
    currentIndex--
    entry := logEntries[currentIndex]
    redoEntry := logEntries[currentIndex+1]
    handleUndoRedo(entry, redoEntry, true)
}

func redo() {  
    if currentIndex >= len(logEntries)-1 {  
        fmt.Println("No actions to redo.")
        return
    }
    currentIndex++
    entry := logEntries[currentIndex]
    undoEntry := logEntries[currentIndex-1]
    handleUndoRedo(undoEntry, entry, false)
}

func handleUndoRedo(undoEntry LogEntry, redoEntry LogEntry, isUndo bool) {
    fmt.Printf("\n")
    switch undoEntry.Activity {
    case "Navigation":
        if isUndo {
            fmt.Println("Undo navigation: " + undoEntry.Details)
        } else {
            fmt.Println("Redo navigation: " + redoEntry.Details)
        }
    case "File Creation":
        if isUndo {
            fmt.Println("Undo file creation: " + undoEntry.Details)
            os.Remove(undoEntry.Details)