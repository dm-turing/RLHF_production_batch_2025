package main

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sync"

	"github.com/fsnotify/fsnotify"
)

var wg sync.WaitGroup

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
				syncFile(dir, targetDir, event.Name)
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

	<-done
	return nil
}

func syncFile(srcDir string, targetDir string, relPath string) {
	srcPath := filepath.Join(srcDir, relPath)
	targetPath := filepath.Join(targetDir, relPath)

	if !filepath.IsAbs(srcPath) {
		return
	}
	fi, err := os.Stat(srcPath)
	if err != nil {
		fmt.Println("Error while stat'ing file:", err)
		return
	}

	if fi.IsDir() {
		os.MkdirAll(targetPath, fi.Mode())
	} else {
		file, err := os.Open(srcPath)
		if err != nil {
			fmt.Println("Error while opening file:", err)
			return
		}
		defer file.Close()

		outFile, err := os.Create(targetPath)
		if err != nil {
			fmt.Println("Error while creating file:", err)
			return
		}
		defer outFile.Close()

		_, err = io.Copy(outFile, file)
		if err != nil {
			fmt.Println("Error while copying file:", err)
			return
		}
	}
}

func main() {
	if len(os.Args) != 3 {
		fmt.Println("Usage:", os.Args[0], "<source-dir> <target-dir>")
		os.Exit(1)
	}

	sourceDir := os.Args[1]
	targetDir := os.Args[2]

	err := os.MkdirAll(targetDir, 0755)
	if err != nil {
		fmt.Println("Error creating target directory:", err)
		os.Exit(1)
	}

	err = watchDir(sourceDir, targetDir)
	if err != nil {
		fmt.Println("Error watching directory:", err)
	}

	wg.Wait()
}
