// restore.go
package main

import (
	"database/sql"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
)

func restoreBackup(db *sql.DB, destinationPath string) {
	rows, err := db.Query(`SELECT path, content FROM files`)
	if err != nil {
		log.Fatalf("Error querying files: %v", err)
	}
	defer rows.Close()

	for rows.Next() {
		var path string
		var content []byte
		if err := rows.Scan(&path, &content); err != nil {
			log.Fatalf("Could not scan row: %v", err)
		}

		// Create necessary directories
		fullPath := filepath.Join(destinationPath, path)
		if err := os.MkdirAll(filepath.Dir(fullPath), os.ModePerm); err != nil {
			log.Fatalf("Could not create directories for %s: %v", fullPath, err)
		}

		// Write the file
		if err := ioutil.WriteFile(fullPath, content, os.ModePerm); err != nil {
			log.Fatalf("Could not write file %s: %v", fullPath, err)
		}
	}
}
