package main

import (
	"fmt"
	"io/ioutil"
	"math/rand"
	"os"
	"path/filepath"
	"sync"
	"time"
)

const (
	folderPath     = "./test_files"
	fileSystemSize = 1024 * 1024 // 1MB in bytes
)

var (
	currentSize       int64
	wg                sync.WaitGroup
	lock              sync.Mutex
	messageThresholds = []int{50, 75, 90}
)

func main() {
	defer cleanup()
	os.MkdirAll(folderPath, 0755)

	wg.Add(1)
	go generateRandomData()

	wg.Wait()
}

func generateRandomData() {
	defer wg.Done()
	for {
		data := make([]byte, 1024) // Generate 1KB of random data
		rand.Read(data)
		filePath := filepath.Join(folderPath, fmt.Sprintf("%d.txt", time.Now().UnixNano()))

		err := ioutil.WriteFile(filePath, data, 0644)
		if err != nil {
			fmt.Printf("Error writing file: %v\n", err)
			continue
		}

		updateCurrentSize(len(data))
		printCurrentUsage()

		if checkQuotaReached() {
			fmt.Println("File system quota hit, gracefully exiting...")
			return
		}
	}
}

func updateCurrentSize(size int) {
	lock.Lock()
	defer lock.Unlock()
	currentSize += int64(size)
}

func checkQuotaReached() bool {
	lock.Lock()
	defer lock.Unlock()
	return float64(currentSize)/float64(fileSystemSize)*100 >= 99
}

func printCurrentUsage() {
	lock.Lock()
	defer lock.Unlock()
	usagePercent := float64(currentSize) / float64(fileSystemSize) * 100

	for _, threshold := range messageThresholds {
		if int(usagePercent) == threshold {
			fmt.Printf("File system usage reached %d%%\n", threshold)
		}
	}

	fmt.Printf("Current file system size: %.2fKB / %.2fKB\n", float64(currentSize)/1024, float64(fileSystemSize)/1024)
}

func cleanup() {
	err := os.RemoveAll(folderPath)
	if err != nil {
		fmt.Printf("Error cleaning up: %v\n", err)
	}
}
