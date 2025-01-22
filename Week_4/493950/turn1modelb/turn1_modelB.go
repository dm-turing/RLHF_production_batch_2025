package main

import (
	"embed"
	"fmt"
	"html/template"
	"io/fs"
	"log"
	"net/http"
)

//go:embed static
var staticFiles embed.FS

func main() {
	http.HandleFunc("/", serveHome)
	http.HandleFunc("/profile", serveProfile)
	http.HandleFunc("/account", serveAccount)
	fsys, _ := fs.Sub(staticFiles, "static")

	// fs := http.FS(staticFiles)
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.FS(fsys))))

	fmt.Println("Server is running on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func serveHome(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.New("home").ParseFS(staticFiles, "static/templates/home.html")
	if err != nil {
		log.Println(err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	if err := tmpl.Execute(w, nil); err != nil {
		log.Println(err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}

func serveProfile(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.New("profile").ParseFS(staticFiles, "static/templates/profile.html")
	if err != nil {
		log.Println(err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	if err := tmpl.Execute(w, nil); err != nil {
		log.Println(err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}

func serveAccount(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.New("account").ParseFS(staticFiles, "static/templates/account.html")
	if err != nil {
		log.Println(err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	if err := tmpl.Execute(w, nil); err != nil {
		log.Println(err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}
