package main

// ... (Previous code)

import (
    "github.com/golang-collections/collections/slist"
)

var sortedList slist.SkipList

func init() {
    // Initialize skip sorted list
    sortedList.Func = func(a, b interface{}) int {
		entryA := a.(FileEntry)
		entryB := b.(FileEntry)
        switch accessSorting {
        case "access":
			return entryA.Access - entryB.Access
        case "modtime":
			return int(entryA.ModTime.Sub(entryB.ModTime).Milliseconds())
        default:
			if entryA.Name < entryB.Name {
				return -1
			}
			if entryA.Name > entryB.Name {
				return 1
			}
			return 0
        }
    }
}

// ... (Rest of the code)

func accessFile(fileName string) {
	// ... (Previous code)

	if entry, exists := fileMap[fileName]; exists {
		sortedList.Remove(entry)
		entry.Access++
		fileMap[fileName] = entry
		sortedList.Insert(entry)
	}
}

func getFileList() []FileEntry {
	var entries []FileEntry
	// Iterate through the skip sorted list in ascending order
	sortedList.Ascend(func(e interface{}) {
		entries = append(entries, e.(FileEntry))
	})