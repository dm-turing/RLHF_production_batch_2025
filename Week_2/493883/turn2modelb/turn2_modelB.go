package main

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sync"

	"github.com/rjeczalik/notify"
)

const (
	watchDirectory  = "folder1/"
	targetDirectory = "folder2/"
)

func main() {
	// Watch changes in the watchDirectory
	err := notify.Watch(watchDirectory, make(chan notify.EventInfo), notify.All)
	if err != nil {
		fmt.Printf("Error watching directory: %v\n", err)
		os.Exit(1)
	}
	defer notify.Stop(nil)

	var wg sync.WaitGroup
	wg.Add(1)

	go func() {
		defer wg.Done()
		for {
			select {
			case event := <-notify.Events:
				handleEvent(event)
			case err := <-notify.Errors:
				fmt.Printf("Error: %v\n", err)
			}
		}
	}()

	fmt.Println("File synchronization started. Press Ctrl+C to exit.")

	// Wait for Ctrl+C signal to stop the synchronization
	<-make(chan os.Signal, 1)

	fmt.Println("Stopping file synchronization.")
}

func handleEvent(event notify.EventInfo) {
	// Ignore events related to self
	if event.Path() == watchDirectory {
		return
	}

	filePath := event.Path()
	targetPath := filepath.Join(targetDirectory, filepath.Base(filePath))

	// Determine the type of event
	switch event.Event() {
	case notify.Create, notify.Write, notify.Remove:
		// Create, write, or remove a file
		if event.Event() == notify.Remove {
			// Handle file removal
			os.Remove(targetPath)
			fmt.Printf("Removed: %s\n", filePath)
		} else {
			// Handle file creation and writing
			copyFile(filePath, targetPath)
			fmt.Printf("Synced: %s\n", filePath)
		}
	case notify.Rename:
		// Handle file renaming
		oldName := event.Path()
		newName := event.RenamePath()
		newTargetPath := filepath.Join(targetDirectory, filepath.Base(newName))

		os.Rename(targetPath, newTargetPath)
		fmt.Printf("Renamed: %s -> %s\n", oldName, newName)
	default:
		fmt.Printf("Unsupported event: %s\n", event.Event())
	}
}

// Simple function to copy a file
func copyFile(src, dst string) {
	input, err := os.Open(src)
	if err != nil {
		fmt.Printf("Error copying file: %v\n", err)
		return
	}
	defer input.Close()

	output, err := os.Create(dst)
	if err != nil {
		fmt.Printf("Error copying file: %v\n", err)
		return
	}
	defer output.Close()

	_, err = io.Copy(output, input)
	if err != nil {
		fmt.Printf("Error copying file: %v\n", err)
		return
	}
}
