package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/fsnotify/fsnotify"
)

const (
	// Maximum number of concurrent synchronization operations.
	maxConcurrency = 10
)

var (
	wg      sync.WaitGroup
	syncQueue chan struct{}
)

func init() {
	syncQueue = make(chan struct{}, maxConcurrency)
}

func watchDir(dir string, targetDir string) error {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return err
	}
	defer watcher.Close()

	done := make(chan bool)

	wg.Add(1)
	go func() {
		defer wg.Done()
		for {
			select {
			case event := <-watcher.Events:
				fmt.Println("Event:", event)
				// Limit the number of concurrent synchronization operations.
				syncQueue <- struct{}{}
				go syncFile(dir, targetDir, event.Name)
				<-syncQueue
			case err := <-watcher.Errors:
				fmt.Println("Error:", err)
			case <-done:
				return
			}
		}
	}()

	err = watcher.Add(dir)
	if err != nil {
		return err
	}

	return <-done
}

func syncFile(srcDir string, targetDir string, relPath string) {
	srcPath := filepath.Join(srcDir, relPath)
	targetPath := filepath.Join(targetDir, relPath)

	// Acquire a lock to prevent race conditions during file operations.
	lock := sync.Mutex{}
	lock.Lock()
	defer lock.Unlock()

	if !filepath.IsAbs(srcPath) {
		return
	}

	fi, err := os.Stat(srcPath)
	if err != nil {
		log.Println("Error while stat'ing file:", err)
		return
	}

	if fi.IsDir() {
		os.MkdirAll(targetPath, fi.Mode())
	} else {
		// Check for file conflicts before synchronizing.
		if _, err := os.Stat(targetPath); err == nil {
			// Resolve conflict by comparing the modified times of the two files.
			srcModTime := fi.ModTime()
			targetModTime, err := os.Stat(targetPath)
			if err != nil {
				log.Println("Error while stat'ing target file:", err)
				return
			}

			if srcModTime.After(targetModTime.ModTime()) {
				log.Println("Source file is newer, synchronizing.")
				// Synchronize the source file to the target.
				err = syncFileContents(srcPath, targetPath)
				if err != nil {
					log.Println("Error while synchronizing file:", err)
				}
			} else {
				log.Println("Target file is newer, skipping synchronization.")
			}
		} else {
			// Synchronize the source file to the target since it doesn't exist.
			err = syncFileContents(srcPath, targetPath)
			if err != nil {
				log.Println("Error while synchronizing file:", err)
			}
		}
	}
}

func syncFileContents(srcPath string, targetPath string) error {
	file, err := os.Open(srcPath)
	if err != nil {
		return err
	}
	defer file.Close()

	outFile, err := os.Create(targetPath)
	if err != nil {
		return err