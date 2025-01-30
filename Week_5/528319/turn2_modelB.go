package main  
import (  
    "database/sql"
    "encoding/json"
    "fmt"
    "io/ioutil"
    "math/rand"
    "net/http"
    "os"
    "path/filepath"
    "sync"
    "time"

    _ "github.com/mattn/go-sqlite3"
)

const (
    dbName = "files.db"
    workerPoolSize = 10
    contentSize    = 1024 // 1KB
)

// FileInfo is a struct to hold file details  
type FileInfo struct {  
    Name      string    `json:"name"`
    Timestamp time.Time `json:"timestamp"`
    Content   string    `json:"content"`
}

// Worker represents a worker in the worker pool  
type Worker struct {  
    ID      int
    WorkQueue chan FileInfo
    Done    chan bool
}

// WorkerPool manages the worker pool  
type WorkerPool struct {  
    Workers []Worker
    WorkQueue chan FileInfo
    Done    chan bool
}

var (
    db *sql.DB
    wg sync.WaitGroup
)

func init() {
    var err error
    db, err = sql.Open("sqlite3", dbName)
    if err != nil {
        panic(err)
    }
    _, err = db.Exec("CREATE TABLE IF NOT EXISTS files (name TEXT PRIMARY KEY, timestamp TEXT, content TEXT)")
    if err != nil {
        panic(err)
    }
}

func generateRandomContent() string {
    bytes := make([]byte, contentSize)
    for i := range bytes {
        bytes[i] = byte(rand.Intn(256))
    }
    return string(bytes)
}

func (w *Worker) Start() {
    go func() {
        for {
            select {
            case fileInfo := <-w.WorkQueue:
                wg.Add(1)
                defer wg.Done()
                err := saveFileToDB(fileInfo)
                if err != nil {
                    fmt.Printf("Worker %d: Error saving file %s: %v\n", w.ID, fileInfo.Name, err)
                }
            case <-w.Done:
                return
            }
        }
    }()
}

func (wp *WorkerPool) Start() {
    for i := 0; i < len(wp.Workers); i++ {
        wp.Workers[i].Start()
    }
}

func (wp *WorkerPool) Stop() {
    for i := 0; i < len(wp.Workers); i++ {
        wp.Workers[i].Done <- true
    }
    wp.Done <- true
}

func saveFileToDB(fileInfo FileInfo) error {
    _, err := db.Exec("INSERT OR REPLACE INTO files (name, timestamp, content) VALUES (?, ?, ?)",
        fileInfo.Name, fileInfo.Timestamp.Format(time.RFC3339), fileInfo.Content)
    return err
}

// ListFilesHandler handles listing files  
func ListFilesHandler(w http.ResponseWriter, r *http.Request) {
    rows, err := db.Query("SELECT name, timestamp, content FROM files")
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    defer rows.Close()

    var fileInfos []FileInfo
    for rows.Next() {
        var fileInfo FileInfo
        if err := rows.Scan(&fileInfo.Name, &fileInfo.Timestamp, &fileInfo.Content); err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)
            return
