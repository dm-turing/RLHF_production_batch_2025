package main

import (
	"fmt"
	"log"
	"net/http"
	// Import the generated virtual file system package
)

func main() {
	// Create a virtual file system for serving static files
	var assets http.FileSystem = http.Dir("./assets")
	file, err := assets.Open("/index.html")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	// Process the index.html file content
	contents := make([]byte, 4096)
	_, err = file.Read(contents)
	if err != nil {
		fmt.Println("Error reading file:", err)
		return
	}
	fmt.Println(string(contents))

	// Serve static files from the virtual file system
	http.Handle("/assets/", http.StripPrefix("/assets/", http.FileServer(assets)))

	// Start the HTTP server
	log.Fatal(http.ListenAndServe(":8080", nil))
}
