package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

type File struct {
	Name    string    `json:"name"`
	Timestamp time.Time `json:"timestamp"`
}

func listFilesHandler(w http.ResponseWriter, r *http.Request) {
	// Get the path from the URL query parameter
	path := r.URL.Query().Get("path")
	if path == "" {
		// Default path to the current working directory
		path, _ = os.Getwd()
	}

	files, err := ioutil.ReadDir(path)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var fileInfos []File
	for _, file := range files {
		fi, err := os.Stat(filepath.Join(path, file.Name()))
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		fileInfos = append(fileInfos, File{
			Name:    file.Name(),
			Timestamp: fi.ModTime(),
		})
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(fileInfos)
}

func createFileHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Parse the form data to get the file name and content
	if err := r.ParseForm(); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	fileName := r.FormValue("filename")
	fileContent := r.FormValue("content")

	if fileName == "" || fileContent == "" {
		http.Error(w, "File name and content are required.", http.StatusBadRequest)
		return
	}

	// Create the file in the current working directory
	f, err := os.Create(fileName)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer f.Close()

	if _, err := f.WriteString(fileContent); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(w, "File %s created successfully.", fileName)
}

func deleteFileHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "DELETE" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Get the file name from the URL parameter
	fileName := r.URL.Query().Get("filename")
	if fileName == "" {
		http.Error(w, "File name is required.", http.StatusBadRequest)
		return
	}

	if err := os.Remove(fileName); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	