package main

import (
	"log"
	"net/http"

	// "github.com/dm-turing/RLHF_production_batch_2025/Week_3/turn1modelb"

	"github.com/shurcooL/vfsgen"
	// Import the generated virtual file system package
)

func main() {
	// Create a virtual file system for serving static files
	var assets http.FileSystem = http.Dir("./data/assets")

	err := vfsgen.Generate(assets, vfsgen.Options{
		PackageName:  "data",
		BuildTags:    "!dev",
		VariableName: "Assets",
	})
	if err != nil {
		log.Fatalln(err)
	}
	// Process the index.html file content
	file, err := assets.Open("./index.html")
	if err != nil {
		return err
	}
	defer file.Close()
	// fmt.Println(string(contents))

	// Serve static files from the virtual file system
	http.Handle("/assets/", http.StripPrefix("/assets/", http.FileServer(assets)))

	// Start the HTTP server
	log.Fatal(http.ListenAndServe(":8080", nil))
}
