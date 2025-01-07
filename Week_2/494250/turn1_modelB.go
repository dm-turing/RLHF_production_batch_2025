package main

import (
	"archive/zip"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
)

func createZip(sourceDir, archivePath string) error {
	archive, err := os.Create(archivePath)
	if err != nil {
		return err
	}
	defer archive.Close()

	zipWriter := zip.NewWriter(archive)
	defer zipWriter.Close()

	return filepath.Walk(sourceDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Skip the source directory itself
		if path == sourceDir {
			return nil
		}

		// Relative path within the archive
		relativePath, err := filepath.Rel(sourceDir, path)
		if err != nil {
			return err
		}

		// Add the file to the archive
		fw, err := zipWriter.Create(relativePath)
		if err != nil {
			return err
		}

		// Open the file to read
		fr, err := os.Open(path)
		if err != nil {
			return err
		}
		defer fr.Close()

		// Copy the file content to the archive
		_, err = io.Copy(fw, fr)
		if err != nil {
			return err
		}
		return nil
	})
}

func main() {
	sourceDir := "../494250"
	archivePath := "../494250/backup.zip"

	err := createZip(sourceDir, archivePath)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Backup completed successfully!")
}
