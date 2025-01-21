
package main

import (
	"database/sql"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/fsnotify/fsnotify"
	// _ "github.com/mattn/go-sqlite3"
	_ "github.com/ncruces/go-sqlite3/driver"
	_ "github.com/ncruces/go-sqlite3/embed"
)

// Backup represents a backup entry in the SQLite database
type Backup struct {
	Path     string `db:"path"`
	ModTime  time.Time `db:"mod_time"`
	FileSize int64 `db:"file_size"`
}

// SQLiteBackupStore implements a BackupStore using SQLite
type SQLiteBackupStore struct {
	db *sql.DB
}

// NewSQLiteBackupStore creates a new SQLiteBackupStore
func NewSQLiteBackupStore(dbPath string) (*SQLiteBackupStore, error) {
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, err
	}
	// Create the backups table if it doesn't exist
	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS backups (
		path TEXT PRIMARY KEY,
		mod_time TEXT,
		file_size INTEGER
	)`)
	if err != nil {
		return nil, err
	}
	return &SQLiteBackupStore{db: db}, nil
}

// Close closes the underlying database connection
func (s *SQLiteBackupStore) Close() error {
	return s.db.Close()
}

// InsertBackup inserts a new backup entry into the database
func (s *SQLiteBackupStore) InsertBackup(b Backup) error {
	_, err := s.db.Exec("INSERT INTO backups (path, mod_time, file_size) VALUES (?, ?, ?)", b.Path, b.ModTime.Format(time.RFC3339), b.FileSize)
	return err
}

// GetBackups retrieves all backup entries from the database
func (s *SQLiteBackupStore) GetBackups() ([]Backup, error) {
	rows, err := s.db.Query("SELECT path, mod_time, file_size FROM backups")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var backups []Backup
	for rows.Next() {
		var b Backup
		if err := rows.Scan(&b.Path, &b.ModTime, &b.FileSize); err != nil {
			return nil, err
		}
		backups = append(backups, b)
	}
	return backups, rows.Err()
}

// DeleteBackup deletes a backup entry from the database
func (s *SQLiteBackupStore) DeleteBackup(path string) error {
	_, err := s.db.Exec("DELETE FROM backups WHERE path = ?", path)
	return err
}
 
func main() {
	// Parse command-line arguments
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run backup.go <source_path> [--restore <restore_path>]")
		return
	}
	sourcePath := os.Args[1]
	var restorePath string
	if len(os.Args) > 2 && os.Args[2] == "--restore" {
		restorePath = os.Args[3]
	}
	
	// Create a new SQLite backup store
	dbPath := filepath.Join(sourcePath, ".backup.db")
	store, err := NewSQLiteBackupStore(dbPath)
	if err != nil {
		log.Fatalf("Error opening backup database: %v", err)
	}
	defer store.Close()