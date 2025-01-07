
package main

import (  
    "archive/zip"
    "errors"
    "fmt"
    "io"
    "log"
    "os"
    "path/filepath"
    "strings"
    "github.com/ulikunitz/xz" // For advanced compression
)

// Error types for the Archiver
var (
    ErrInvalidSourceDir = errors.New("invalid source directory")
    ErrInvalidArchivePath = errors.New("invalid archive path")
    ErrCreatingArchive = errors.New("error creating archive")
    ErrWalkingDirectory = errors.New("error walking directory")
    ErrAddingFileToArchive = errors.New("error adding file to archive")
    ErrCreatingCompressor = errors.New("error creating compressor")
)

// Archiver interface defines the methods for archiving files
type Archiver interface {
    Archive(sourceDir, archivePath string) error
}

// xzArchiver implements the Archiver interface using XZ compression
type xzArchiver struct{}

func (xa xzArchiver) Archive(sourceDir, archivePath string) error {
    // Validate input directories
    if _, err := os.Stat(sourceDir); os.IsNotExist(err) {
        return ErrInvalidSourceDir
    }
    archiveDir := filepath.Dir(archivePath)
    if _, err := os.Stat(archiveDir); os.IsNotExist(err) {
        if err := os.MkdirAll(archiveDir, 0755); err != nil {
            return err
        }
    }

    // Create a new XZ archive
    archive, err := os.Create(archivePath)
    if err != nil {
        return ErrCreatingArchive
    }
    defer archive.Close()

    // Wrap the archive in a XZ compressor
    compressor, err := xz.NewWriter(archive)
    if err != nil {
        return ErrCreatingCompressor
    }
    defer compressor.Close()

    zipWriter := zip.NewWriter(compressor)
    defer zipWriter.Close()

    return filepath.Walk(sourceDir, func(path string, info os.FileInfo, err error) error {
        if err != nil {
            return err
        }

        // Skip the source directory itself
        if path == sourceDir {
            return nil
        }
 
        // Proceed with archiving the file
        return xa.archiveFile(path, zipWriter)
    })
}

func (xa xzArchiver) archiveFile(path string, zipWriter *zip.Writer) error {
    // Relative path within the archive
    relativePath, err := filepath.Rel(filepath.Dir(path), path)
    if err != nil {
        return err
    }

    // Add the file to the archive
    fw, err := zipWriter.Create(relativePath)
    if err != nil {
        return ErrAddingFileToArchive
    }

    // Open the file to read
    fr, err := os.Open(path)
    if err != nil {
        return err
    }
    defer fr.Close()

    // Copy the file content to the archive
    _, err = io.Copy(fw, fr)
    return err
} 
 
func main() {
   // Define the source directory and archive path
   sourceDir := "/path/to/source/directory"
   archivePath := "/path/to/backup.xz" // Use ".xz" extension for XZ compression

   // Create an instance of the XZ archiver
   archiver := xzArchiver{}

   if err := archiver.Archive(sourceDir, archivePath); err != nil {
        switch err {
        case ErrInvalidSourceDir:
            log.Fatal("Source directory does not exist.")
        case ErrInvalidArchivePath:
            log.Fatal("Invalid archive path provided.")