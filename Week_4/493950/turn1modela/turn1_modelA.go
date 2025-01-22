package main

import (
	"embed"
	"io/fs"
	"log"
	"net/http"
)

//go:embed static/*
var staticFiles embed.FS

func main() {
	// Serve static files
	fsys, err := fs.Sub(staticFiles, "static")
	if err != nil {
		log.Fatal(err)
	}
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.FS(fsys))))

	// Handle profile and account routes
	http.HandleFunc("/profile", func(w http.ResponseWriter, r *http.Request) {
		data, err := staticFiles.ReadFile("static/html/profile.html")
		if err != nil {
			http.Error(w, "File not found.", http.StatusNotFound)
			return
		}
		w.Write(data)
	})

	http.HandleFunc("/account", func(w http.ResponseWriter, r *http.Request) {
		data, err := staticFiles.ReadFile("static/html/account.html")
		if err != nil {
			http.Error(w, "File not found.", http.StatusNotFound)
			return
		}
		w.Write(data)
	})

	log.Println("Server started at http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
