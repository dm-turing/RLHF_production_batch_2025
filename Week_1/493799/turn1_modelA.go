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
	path     string
	interval time.Duration
	ctx      context.Context
	cancel   context.CancelFunc
}

// NewFileAuditor creates a new FileAuditor instance
func NewFileAuditor(path string, interval time.Duration) *FileAuditor {
	ctx, cancel := context.WithCancel(context.Background())
	return &FileAuditor{
		path:     path,
		interval: interval,
		ctx:      ctx,
		cancel:   cancel,
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

// audit checks for file access and modification
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
		fmt.Printf("%s last modified at %s\n", file.Name(), modTime.Format("2006-01-02 15:04:05"))
	}
}

func main() {
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
