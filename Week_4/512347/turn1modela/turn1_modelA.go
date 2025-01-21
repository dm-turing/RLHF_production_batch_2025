package main

import (
	"database/sql"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	// _ "github.com/mattn/go-sqlite3"
	_ "github.com/ncruces/go-sqlite3/driver"
	_ "github.com/ncruces/go-sqlite3/embed"
)

// Function to perform the backup
func backup(sourceDir string, dbPath string) error {
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return err
	}
	defer db.Close()

	// Prepare the database
	sqlStmt := `
	CREATE TABLE IF NOT EXISTS files (
		id INTEGER PRIMARY KEY,
		path TEXT NOT NULL,
		data BLOB NOT NULL
	);
	DELETE FROM files;
	`
	_, err = db.Exec(sqlStmt)
	if err != nil {
		return fmt.Errorf("failed to execute SQL: %w", err)
	}

	tx, err := db.Begin()
	if err != nil {
		return err
	}

	stmt, err := tx.Prepare("INSERT INTO files (path, data) VALUES (?, ?)")
	if err != nil {
		return err
	}
	defer stmt.Close()

	err = filepath.Walk(sourceDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			data, err := ioutil.ReadFile(path)
			if err != nil {
				return err
			}
			relPath, err := filepath.Rel(sourceDir, path)
			if err != nil {
				return err
			}
			_, err = stmt.Exec(relPath, data)
			if err != nil {
				return err
			}
		}
		return nil
	})

	if err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit()
}

// Function to perform the restore
func restore(dbPath string, restoreDir string) error {
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return err
	}
	defer db.Close()

	rows, err := db.Query("SELECT path, data FROM files")
	if err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {
		var path string
		var data []byte
		err = rows.Scan(&path, &data)
		if err != nil {
			return err
		}
		fullPath := filepath.Join(restoreDir, path)
		err = os.MkdirAll(filepath.Dir(fullPath), os.ModePerm)
		if err != nil {
			return err
		}
		err = ioutil.WriteFile(fullPath, data, 0644)
		if err != nil {
			return err
		}
	}
	return nil
}

func main() {
	backupFlag := flag.Bool("backup", false, "set to perform a backup")
	restoreFlag := flag.Bool("restore", false, "set to restore from backup")
	sourceDir := flag.String("source", "", "source directory for backup/restore")
	dbPath := flag.String("db", "backup.db", "path to the SQLite database")
	restoreDir := flag.String("restore-to", "", "target directory for restore (required for restore)")

	flag.Parse()

	if *backupFlag == *restoreFlag {
		log.Fatalf("You must specify either -backup or -restore")
	}

	if *backupFlag {
		if *sourceDir == "" {
			log.Fatalf("You must specify a source directory to backup")
		}
		err := backup(*sourceDir, *dbPath)
		if err != nil {
			log.Fatalf("Backup failed: %v", err)
		}
		fmt.Println("Backup completed successfully")
	} else if *restoreFlag {
		if *restoreDir == "" {
			log.Fatalf("You must specify a restore directory")
		}
		err := restore(*dbPath, *restoreDir)
		if err != nil {
			log.Fatalf("Restore failed: %v", err)
		}
		fmt.Println("Restore completed successfully")
	}
}
