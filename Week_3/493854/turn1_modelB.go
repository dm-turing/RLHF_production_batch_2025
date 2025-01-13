package main

import (
	"fmt"
	"log"
	"net/http"

	// Import the generated virtual file system package
	"github.com/shurcooL/httpfs/html/vfstemplate"
)

func main() {
	// Use the Asset function to read the "index.html" file
	indexHTML, err := Asset("assets/index.html")
	if err != nil {
		log.Fatal(err)
	}

	// Process the index.html file content
	fmt.Println(string(indexHTML))

	// Create a virtual file system for serving static files
	fs := vfstemplate.FileSystem(Asset)

	// Serve static files from the virtual file system
	http.Handle("/assets/", http.StripPrefix("/assets/", http.FileServer(fs)))

	// Start the HTTP server
	log.Fatal(http.ListenAndServe(":8080", nil))
}
