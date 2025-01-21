// backup.go
package main

import (
	"crypto/md5"
	"database/sql"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
)

func fullBackup(db *sql.DB, sourcePath string) {
	err := filepath.Walk(sourcePath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			hash, content := readFile(path)
			storeFile(db, path, hash, content)
		}
		return nil
	})
	if err != nil {
		log.Fatalf("Error walking through files: %v", err)
	}
}

func incrementalBackup(db *sql.DB, sourcePath string) {
	existingFiles, err := fetchExistingFiles(db)
	if err != nil {
		log.Fatalf("Error fetching existing files: %v", err)
	}

	err = filepath.Walk(sourcePath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			hash, content := readFile(path)
			// Check if the file has changed
			if existingHash, exists := existingFiles[path]; !exists || existingHash != hash {
				storeFile(db, path, hash, content)
			}
		}
		return nil
	})
	if err != nil {
		log.Fatalf("Error walking through files: %v", err)
	}
}

func readFile(path string) (string, []byte) {
	content, err := ioutil.ReadFile(path)
	if err != nil {
		log.Fatalf("Could not read file %s: %v", path, err)
	}
	hash := md5.Sum(content)
	return string(hash[:]), content
}

func storeFile(db *sql.DB, path string, hash string, content []byte) {
	_, err := db.Exec(`INSERT OR REPLACE INTO files (path, hash, content) VALUES (?, ?, ?)`, path, hash, content)
	if err != nil {
		log.Fatalf("Could not store file %s: %v", path, err)
	}
}

func fetchExistingFiles(db *sql.DB) (map[string]string, error) {
	rows, err := db.Query(`SELECT path, hash FROM files`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	files := make(map[string]string)
	for rows.Next() {
		var path, hash string
		if err := rows.Scan(&path, &hash); err != nil {
			return nil, err
		}
		files[path] = hash
	}
	return files, nil
}
