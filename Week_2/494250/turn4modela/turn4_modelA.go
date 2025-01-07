package main

import (
	"archive/zip"
	"fmt"
	"log"
	"os"
	"path/filepath"
)

type IncrementalBackup interface {
	CreateInitialBackup(sourceDir, archivePath string) error
	CreateIncrementalBackup(sourceDir, archivePath string) error
}

type BasicBackup struct{}

func (bb BasicBackup) CreateInitialBackup(sourceDir, archivePath string) error {
	return createZip(sourceDir, archivePath)
}

func (bb BasicBackup) CreateIncrementalBackup(sourceDir, archivePath string) error {
	return createZip(sourceDir, archivePath) // Note: Simple overwrite
}

func createZip(sourceDir, archivePath string) error {
	archive, err := os.Create(archivePath)
	if err != nil {
		return err
	}
	defer archive.Close()

	writer := zip.NewWriter(archive)

	defer writer.Close()

	return filepath.Walk(sourceDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		var header zip.FileHeader
		header.Name, _ = filepath.Rel(sourceDir, path)
		header.Modified = info.ModTime()
		header.Method = zip.Deflate

		reader, err := os.Open(path)
		if err != nil {
			return err
		}
		defer reader.Close()
		writer.CreateHeader(&header)
		_, err = archive.Write([]byte(path))
		if err != nil {
			return err
		}
		return nil
	})
}

func main() {
	sourceDir := "../../494250"
	archivePath := "./backup.zip"

	backupSystem := BasicBackup{}

	err := backupSystem.CreateInitialBackup(sourceDir, archivePath)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Initial backup completed successfully!")

	err = backupSystem.CreateIncrementalBackup(sourceDir, archivePath)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Incremental backup completed successfully!")
}
