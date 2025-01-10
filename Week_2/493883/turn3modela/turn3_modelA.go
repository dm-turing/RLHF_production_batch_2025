package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"sync"
	"sync/atomic"
	"time"

	"github.com/fsnotify/fsnotify"
)

type conflictResolver interface {
	Resolve(sourceDir, targetDir, fileName string) error
}

func watchDir(dir string, targetDir string, resolver conflictResolver, tickerPeriod time.Duration) error {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return err
	}
	defer watcher.Close()

	done := make(chan bool)

	var eventCount int64
	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer wg.Done()
		ticker := time.NewTicker(tickerPeriod)
		for {
			select {
			case event := <-watcher.Events:
				// Debug
				fmt.Printf("Event: %v, %v %s\n", event.Op&fsnotify.Write, event.Op&fsnotify.Create, event.Name)

				if event.Op&(fsnotify.Write|fsnotify.Create) == 0 {
					continue
				}

				atomic.AddInt64(&eventCount, 1)
				syncFile(dir, targetDir, event.Name, resolver)
			case err := <-watcher.Errors:
				log.Println("Error:", err)
			case <-ticker.C:
				fmt.Printf("%d events in the last %s\n", eventCount, tickerPeriod)
				eventCount = 0
			case <-done:
				ticker.Stop()
				return
			}
		}
	}()

	err = watcher.Add(dir)
	if err != nil {
		return err
	}

	<-done
	return nil
}

func syncFile(srcDir string, targetDir string, relPath string, resolver conflictResolver) {
	srcPath := filepath.Join(srcDir, relPath)
	targetPath := filepath.Join(targetDir, relPath)

	srcMutex := &sync.Mutex{}
	targetMutex := &sync.Mutex{}

	if !filepath.IsAbs(srcPath) {
		return
	}
	fi, err := os.Stat(srcPath)
	if err != nil {
		log.Printf("Error while stat'ing file: %v", err)
		return
	}

	if fi.IsDir() {
		os.MkdirAll(targetPath, fi.Mode())
	} else {
		srcMutex.Lock()
		defer srcMutex.Unlock()

		targetMutex.Lock()
		defer targetMutex.Unlock()

		srcFile, err := os.Open(srcPath)
		if err != nil {
			log.Printf("Error while opening source file: %v", err)
			return
		}
		defer srcFile.Close()

		targetFile, err := os.OpenFile(targetPath, os.O_RDWR|os.O_CREATE|os.O_TRUNC, fi.Mode())
		if err != nil {
			log.Printf("Error while opening target file: %v", err)
			return
		}
		defer targetFile.Close()

		stat, err := srcFile.Stat()
		if err != nil {
			log.Printf("Error getting source file stats: %v", err)
			return
		}
		fi, err = targetFile.Stat()
		if err != nil {
			log.Printf("Error getting target file stats: %v", err)
			return
		}

		if fi.Size() == stat.Size() && fi.ModTime().Unix() >= stat.ModTime().Unix() {
			log.Printf("Target file is up to date or newer: %s", targetPath)
			return
		}

		_, err = io.Copy(targetFile, srcFile)
		if err != nil {
			log.Printf("Error while copying file: %v", err)
			return
		}

		err = resolver.Resolve(srcDir, targetDir, relPath)
		if err != nil {
			log.Printf("Conflict resolution failed: %v", err)
		}
	}
}

func (r *defaultConflictResolver) Resolve(srcDir, targetDir, fileName string) error {
	srcPath := filepath.Join(srcDir, fileName)
	targetPath := filepath.Join(targetDir, fileName)

	srcFile, err := os.Open(srcPath)
	if err != nil {
		return err
	}
	defer srcFile.Close()

	targetFile, err := os.OpenFile(targetPath, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}
	defer targetFile.Close()

	var conflictFileName string
	if _, err := fmt.Fscanf(srcFile, "%s\n", &conflictFileName); err != nil {
		return err
	}

	sourceConflict := filepath.Join(srcDir, conflictFileName)
	targetConflict := filepath.Join(targetDir, conflictFileName)

	err = os.Rename(srcPath, sourceConflict)
	if err != nil {
		return err
	}

	_, err = io.Copy(targetFile, srcFile)
	if err != nil {
		return err
	}

	return os.Rename(targetPath, targetConflict)
}

type defaultConflictResolver struct{}

func main() {
	if len(os.Args) != 3 {
		fmt.Println("Usage:", os.Args[0], "<source-dir> <target-dir>")
		os.Exit(1)
	}

	sourceDir := os.Args[1]
	targetDir := os.Args[2]

	resolver := &defaultConflictResolver{}

	err := os.MkdirAll(targetDir, 0755)
	if err != nil {
		log.Printf("Error creating target directory: %v", err)
		os.Exit(1)
	}

	tickerPeriod := time.Second * 10
	err = watchDir(sourceDir, targetDir, resolver, tickerPeriod)
	if err != nil {
		log.Printf("Error watching directory: %v", err)
	}
}
