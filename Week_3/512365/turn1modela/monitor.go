package main

import (
	"log"
	"time"

	"github.com/fsnotify/fsnotify"
)

type FileChange struct {
	Event    string    `json:"event"`
	FilePath string    `json:"file_path"`
	Time     time.Time `json:"time"`
}

var changes []FileChange

func watchFileSystem(path string) {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}
	defer watcher.Close()

	err = watcher.Add(path)
	if err != nil {
		log.Fatal(err)
	}

	for {
		select {
		case event, ok := <-watcher.Events:
			if !ok {
				return
			}

			change := FileChange{
				Event:    event.Op.String(),
				FilePath: event.Name,
				Time:     time.Now(),
			}

			changes = append(changes, change)
			log.Println("Event:", event)

		case err, ok := <-watcher.Errors:
			if !ok {
				return
			}
			log.Println("Error:", err)
		}
	}
}
