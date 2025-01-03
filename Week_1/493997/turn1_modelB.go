package main

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"sync"
	"time"
)

// FileInfo represents metadata about a file or directory.
type FileInfo struct {
	os.FileInfo
	AccessTime time.Time
}

// ByAccessTime implements sort.Interface for []FileInfo based on the AccessTime field.
type ByAccessTime []FileInfo

func (a ByAccessTime) Len() int           { return len(a) }
func (a ByAccessTime) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByAccessTime) Less(i, j int) bool { return a[i].AccessTime.Before(a[j].AccessTime) }

// FileSystem represents a sorted file system.
type FileSystem struct {
	mu            sync.RWMutex
	root          string
	entries       map[string]*FileInfo
	sortedEntries []*FileInfo
	sortCriterion string
}

// NewFileSystem creates a new FileSystem with the given root directory.
func NewFileSystem(root string) *FileSystem {
	return &FileSystem{
		root:          root,
		entries:       make(map[string]*FileInfo),
		sortedEntries: []*FileInfo{},
	}
}

// Update updates the access time of the specified file or directory in the file system.
func (fs *FileSystem) Update(path string) {
	fs.mu.Lock()
	defer fs.mu.Unlock()

	absPath, err := filepath.Abs(path)
	if err != nil {
		fmt.Printf("Error updating path: %v\n", err)
		return
	}

	if _, ok := fs.entries[absPath]; !ok {
		fmt.Printf("Path not found: %v\n", absPath)
		return
	}

	fi, err := os.Stat(absPath)
	if err != nil {
		fmt.Printf("Error updating path: %v\n", err)
		return
	}

	fs.entries[absPath] = &FileInfo{
		FileInfo:   fi,
		AccessTime: time.Now(),
	}

	fs.reorderEntries()
}

// List lists the files and directories in the file system sorted by the current sort criterion.
func (fs *FileSystem) List() []string {
	fs.mu.RLock()
	defer fs.mu.RUnlock()

	var paths []string
	for _, entry := range fs.sortedEntries {
		paths = append(paths, entry.Name())
	}
	return paths
}

// SetSortCriterion sets the sort criterion for the file system.
func (fs *FileSystem) SetSortCriterion(criterion string) {
	fs.mu.Lock()
	defer fs.mu.Unlock()

	fs.sortCriterion = criterion
	fs.reorderEntries()
}

// reorderEntries reorders the entries in the file system based on the current sort criterion.
func (fs *FileSystem) reorderEntries() {
	switch fs.sortCriterion {
	case "access_time":
		sort.Sort(ByAccessTime(fs.sortedEntries))
	default:
		// Implement other sort criteria as needed
		sort.Sort(sort.StringSlice(fs.sortedEntries))
	}
}

// scanEntries scans the file system and adds new entries to the map.
func (fs *FileSystem) scanEntries() error {
	return filepath.Walk(fs.root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		absPath, err := filepath.Abs(path)
		if err != nil {
			return err
		}

		fs.entries[absPath] = &FileInfo{
			FileInfo:   info,
			AccessTime: info.ModTime(),
		}

		return nil
	})
}

// Start starts the file system by scanning its contents and setting up a watcher.
func (fs *FileSystem) Start() error {
	err := fs.scanEntries()
	if err != nil {
		return err
	}

	fs.reorderEntries()

	return nil
}

func main() {
	fs := NewFileSystem("/path/to/directory")
	err := fs.Start()
	if err != nil {
		fmt.Printf("Error starting file system: %v\n", err)
		return
	}

	// Update access time of a file
	fs.Update("/path/to/directory/file.txt")

	// List files sorted by access time
	paths := fs.List()
	fmt.Println("Files sorted by access time:")
	for _, path := range paths {
		fmt.Println(path)
	}

	// Set sort criterion to "name"
	fs.SetSortCriterion("name")

	// List files sorted by name
	paths = fs.List()
	fmt.Println("\nFiles sorted by name:")
	for _, path := range paths {
		fmt.Println(path)
	}
}
