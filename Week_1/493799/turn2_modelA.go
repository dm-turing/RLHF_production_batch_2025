package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

// FileAuditor struct encapsulates the logic for file auditing
type FileAuditor struct {
	path      string
	auditPath string
	interval  time.Duration
	ctx       context.Context
	cancel    context.CancelFunc
}

// NewFileAuditor creates a new FileAuditor instance
func NewFileAuditor(path, auditPath string, interval time.Duration) *FileAuditor {
	ctx, cancel := context.WithCancel(context.Background())
	return &FileAuditor{
		path:      path,
		auditPath: auditPath,
		interval:  interval,
		ctx:       ctx,
		cancel:    cancel,
	}
}

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

// Stop stops the file auditing process
func (fa *FileAuditor) Stop() {
	fa.cancel()
}

// audit checks for file access and modification and logs changes
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
		entry := fmt.Sprintf("%s,%s,%s\n", time.Now().Format("2006-01-02 15:04:05"), file.Name(), modTime.Format("2006-01-02 15:04:05"))
		fa.logEntry(entry)
		fmt.Printf("%s last modified at %s\n", file.Name(), modTime.Format("2006-01-02 15:04:05"))
	}
}

// logEntry logs an audit entry to the audit file
func (fa *FileAuditor) logEntry(entry string) {
	file, err := os.OpenFile(fa.auditPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Printf("Error opening audit file %s: %v", fa.auditPath, err)
		return
	}
	defer file.Close()
	_, err = file.WriteString(entry)
	if err != nil {
		log.Printf("Error writing to audit file %s: %v", fa.auditPath, err)
	}
}

func main() {
	if len(os.Args) < 3 {
		fmt.Println("Usage: go run main.go <directory-path> <audit-file-path> [interval-in-seconds]")
		os.Exit(1)
	}

	path := os.Args[1]
	auditPath := os.Args[2]
	interval := 60 * time.Second // Default to 1 minute if no interval is provided

	if len(os.Args) > 3 {
		intervalStr := os.Args[3]
		intInterval, err := time.ParseDuration(intervalStr + "s")
		if err != nil {
			log.Fatalf("Invalid interval: %v", err)
		}
		interval = intInterval
	}

	auditor := NewFileAuditor(path, auditPath, interval)
	auditor.Start()

	// Waiting for a signal to stop the auditor
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
	<-c

	auditor.Stop()
}
