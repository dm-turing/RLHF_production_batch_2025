package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sync"
	"time"
)

// Event represents a file system event.
type Event struct {
	Time time.Time
	Path string
	Op   string
}

// FileAuditor monitors file system events.
type FileAuditor struct {
	events   chan Event
	wg       sync.WaitGroup
	done     chan struct{}
	watchDir string
}

// NewFileAuditor creates a new FileAuditor instance.
func NewFileAuditor(watchDir string) *FileAuditor {
	a := &FileAuditor{
		events:   make(chan Event),
		done:     make(chan struct{}),
		watchDir: watchDir,
	}
	return a
}

// Start begins monitoring file system events.
func (a *FileAuditor) Start() {
	a.wg.Add(1)
	go a.watch()
}

// Stop stops monitoring file system events and closes the event channel.
func (a *FileAuditor) Stop() {
	close(a.done)
	a.wg.Wait()
	close(a.events)
}

// Events returns a channel of file system events.
func (a *FileAuditor) Events() <-chan Event {
	return a.events
}

func (a *FileAuditor) watch() {
	defer a.wg.Done()
	err := filepath.Walk(a.watchDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			log.Printf("Error walking %s: %v", path, err)
			return nil // continue walking
		}

		// Watch only regular files
		if info.Mode().IsRegular() {
			watcher, err := os.Open(path)
			if err != nil {
				log.Printf("Error watching %s: %v", path, err)
				return nil // continue walking
			}
			defer watcher.Close()

			a.wg.Add(1)
			go func() {
				defer a.wg.Done()
				a.monitorFile(watcher, path)
			}()
		}
		return nil
	})
	if err != nil {
		log.Printf("Error walking directory: %v", err)
	}

	<-a.done // wait for done signal to stop monitoring
}

func (a *FileAuditor) monitorFile(f *os.File, path string) {
	var prevInfo os.FileInfo
	var err error

	for {
		select {
		case <-a.done:
			return
		default:
			prevInfo, err = f.Stat()
			if err != nil {
				log.Printf("Error stat'ing %s: %v", path, err)
				return
			}

			time.Sleep(time.Second) // Wait for next check

			newInfo, err := f.Stat()
			if err != nil {
				log.Printf("Error stat'ing %s: %v", path, err)
				return
			}

			if !prevInfo.ModTime().Equal(newInfo.ModTime()) {
				a.events <- Event{Time: time.Now(), Path: path, Op: "Modified"}
			}
		}
	}
}

func main() {
	watchDir := "../493799" // Replace this with the directory you want to audit
	auditor := NewFileAuditor(watchDir)
	auditor.Start()
	defer auditor.Stop()

	// Consume events from the channel
	for event := range auditor.Events() {
		fmt.Printf("%s: %s %s\n", event.Time.Format("2006-01-02 15:04:05"), event.Op, event.Path)
	}
}
