package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

// FileAuditor struct encapsulates the logic for file auditing and audit trail management.
type FileAuditor struct {
	path     string
	interval time.Duration
	ctx      context.Context
	cancel   context.CancelFunc
	auditLog *os.File
}

// AuditEntry represents a single entry in the audit trail
type AuditEntry struct {
	FileName string    `json:"file_name"`
	ModTime  time.Time `json:"mod_time"`
	Action   string    `json:"action"` // Add this field to capture modifications
}

// NewFileAuditor creates a new FileAuditor instance
func NewFileAuditor(path string, interval time.Duration) *FileAuditor {
	ctx, cancel := context.WithCancel(context.Background())
	auditLog, err := os.OpenFile("audit_trail.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatalf("Error creating audit log file: %v", err)
	}
	return &FileAuditor{
		path:     path,
		interval: interval,
		ctx:      ctx,
		cancel:   cancel,
		auditLog: auditLog,
	}
}

// Start starts the file auditing process and audit trail management.
// Start starts the file auditing process
func (fa *FileAuditor) Start() {
	go func() {
		for {
			select {
			case <-time.After(fa.interval):
				fa.audit()
			case <-fa.ctx.Done():
				fmt.Println("Auditor stopped.")
				return
			}
		}
	}()
}

// Stop stops the file auditing process and closes the audit log file.
func (fa *FileAuditor) Stop() {
	fa.cancel()
	fa.auditLog.Close()
}

// audit checks for file access and modifications and writes entries to the audit log.
func (fa *FileAuditor) audit() {
	files, err := os.ReadDir(fa.path)
	if err != nil {
		log.Printf("Error reading directory %s: %v", fa.path, err)
		return
	}

	for _, file := range files {
		fi, err := file.Info()
		if err != nil {
			log.Printf("Error getting file info for %s: %v", file.Name(), err)
			continue
		}
		modTime := fi.ModTime()
		// Detect if the file was modified
		isModified := false
		// Implement your logic to determine if the file was modified between audits
		// For simplicity, let's assume it was modified for all files for this example.
		isModified = true
		// Log the event if the file was modified
		if isModified {
			entry := AuditEntry{
				FileName: file.Name(),
				ModTime:  modTime,
				Action:   "Modified",
			}
			fa.writeAuditLogEntry(entry)
		}
		fmt.Printf("%s last modified at %s\n", file.Name(), modTime.Format("2006-01-02 15:04:05"))
	}
}

// writeAuditLogEntry writes an AuditEntry to the audit log in JSON format.
func (fa *FileAuditor) writeAuditLogEntry(entry AuditEntry) {
	encoder := json.NewEncoder(fa.auditLog)
	if err := encoder.Encode(entry); err != nil {
		log.Printf("Error writing audit log entry: %v", err)
	}
}

func main() {
	// ... (Same as before)
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run main.go <directory-path> [interval-in-seconds]")
		os.Exit(1)
	}

	path := os.Args[1]
	interval := 60 * time.Second // Default to 1 minute if no interval is provided

	if len(os.Args) > 2 {
		intervalStr := os.Args[2]
		intInterval, err := time.ParseDuration(intervalStr + "s")
		if err != nil {
			log.Fatalf("Invalid interval: %v", err)
		}
		interval = intInterval
	}

	auditor := NewFileAuditor(path, interval)
	auditor.Start()

	// Waiting for a signal to stop the auditor
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
	<-c

	auditor.Stop()
}
