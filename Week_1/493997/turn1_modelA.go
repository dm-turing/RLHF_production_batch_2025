package main

import (
	"fmt"
	"math/rand"
	"sort"
	"sync"
	"time"
)

type FileEntry struct {
	Name    string
	ModTime time.Time
	Access  int
}

type EntrySlice []FileEntry

func (s EntrySlice) Len() int      { return len(s) }
func (s EntrySlice) Swap(i, j int) { s[i], s[j] = s[j], s[i] }
func (s EntrySlice) Less(i, j int) bool {
	switch accessSorting {
	case "access":
		return s[i].Access > s[j].Access
	case "modtime":
		return s[i].ModTime.After(s[j].ModTime)
	default:
		return s[i].Name < s[j].Name
	}
}

var accessSorting string
var fileMap map[string]FileEntry
var fileLock sync.RWMutex

func main() {
	fileMap = make(map[string]FileEntry)

	// Simulate adding files
	files := []string{"file1.txt", "file2.txt", "dir/subdir/file3.txt", "file4.txt"}
	for _, fileName := range files {
		fileLock.Lock()
		fileMap[fileName] = FileEntry{Name: fileName, ModTime: time.Now(), Access: 0}
		fileLock.Unlock()
	}

	// Simulate accessing files and updating their metadata
	go simulateFileAccesses()

	// Periodically print the sorted file/directory list
	go func() {
		for range time.Tick(2 * time.Second) {
			sort.Sort(EntrySlice(getFileList()))
			printFiles()
		}
	}()

	select {} // Block indefinitely
}

func simulateFileAccesses() {
	for {
		time.Sleep(time.Duration(rand.Intn(100)) * time.Millisecond)
		fileName := generateRandomFileName(fileMap)
		accessFile(fileName)
	}
}

func generateRandomFileName(fileMap map[string]FileEntry) string {
	if len(fileMap) == 0 {
		return ""
	}
	keys, i := make([]string, len(fileMap)), 0
	for key := range fileMap {
		keys[i] = key
		i++
	}
	return keys[rand.Intn(len(keys))]
}

func accessFile(fileName string) {
	fileLock.RLock()
	defer fileLock.RUnlock()

	if entry, exists := fileMap[fileName]; exists {
		entry.Access++
		fileMap[fileName] = entry
	}
}

func getFileList() []FileEntry {
	fileLock.RLock()
	defer fileLock.RUnlock()
	var entries []FileEntry
	for _, entry := range fileMap {
		entries = append(entries, entry)
	}
	return entries
}

func printFiles() {
	entries := getFileList()
	sort.Sort(EntrySlice(entries))

	fmt.Printf("Sorted files and directories:\n")
	for _, entry := range entries {
		fmt.Printf("%s: %d accesses, last modified at %v\n", entry.Name, entry.Access, entry.ModTime)
	}
	fmt.Printf("------------------------------\n")
}

// Function to dynamically set the sorting criterion
func setSortingCriterion(criterion string) {
	accessSorting = criterion
}
