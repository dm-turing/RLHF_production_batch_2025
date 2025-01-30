package main

import (
	"fmt"
	"hash/crc32"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"sync"
	"time"

	"github.com/cheggaaa/pb"
)

// ComputeCRC32 calculates the CRC32 checksum of a file content.
func ComputeCRC32(filePath string) (uint32, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return 0, err
	}
	defer file.Close()
	hash := crc32.NewIEEE()
	if _, err := io.Copy(hash, file); err != nil {
		return 0, err
	}
	return hash.Sum32(), nil
}

func copyFile(src, dst string) error {
	sourceFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer sourceFile.Close()

	destinationFile, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer destinationFile.Close()

	_, err = io.Copy(destinationFile, sourceFile)
	return err
}

func copyFileWithIntegrityCheck(src, dst string) error {
	// Copy the file first
	if err := copyFile(src, dst); err != nil {
		return err
	}

	// Compute and compare checksums
	srcCRC, err := ComputeCRC32(src)
	if err != nil {
		return fmt.Errorf("failed to compute CRC for src: %w", err)
	}

	dstCRC, err := ComputeCRC32(dst)
	if err != nil {
		return fmt.Errorf("failed to compute CRC for dst: %w", err)
	}

	if srcCRC != dstCRC {
		return fmt.Errorf("CRC mismatch for file: %s (src: %d, dst: %d)", src, srcCRC, dstCRC)
	}

	return nil
}

func worker(wg *sync.WaitGroup, jobs <-chan [2]string, bar *pb.ProgressBar) {
	defer wg.Done()
	for pair := range jobs {
		src := pair[0]
		dst := pair[1]
		if err := copyFileWithIntegrityCheck(src, dst); err != nil {
			log.Printf("Failed to copy %s to %s: %v", src, dst, err)
		}
		bar.Increment()
	}
}

func concurrentCopyFiles(srcDir, dstDir string, numWorkers int) error {
	files, err := ioutil.ReadDir(srcDir)
	if err != nil {
		return err
	}

	jobs := make(chan [2]string, len(files))
	var wg sync.WaitGroup

	// Initialize progress bar
	bar := pb.StartNew(len(files))

	for w := 0; w < numWorkers; w++ {
		wg.Add(1)
		go worker(&wg, jobs, bar)
	}

	for _, file := range files {
		srcPath := filepath.Join(srcDir, file.Name())
		dstPath := filepath.Join(dstDir, file.Name())
		jobs <- [2]string{srcPath, dstPath}
	}
	close(jobs)

	wg.Wait()
	bar.Finish()
	return nil
}

func main() {
	if len(os.Args) != 4 {
		fmt.Println("Usage: go run main.go <source_directory> <destination_directory> <number_of_workers>")
		return
	}

	srcDir := os.Args[1]
	dstDir := os.Args[2]
	numWorkers, err := strconv.Atoi(os.Args[3])
	if err != nil || numWorkers <= 0 {
		log.Fatalf("Invalid number of workers: %v", os.Args[3])
	}

	// Measure elapsed time for concurrent version
	start := time.Now()
	err = concurrentCopyFiles(srcDir, dstDir, numWorkers)
	if err != nil {
		log.Fatalf("Error in concurrent copying: %v", err)
	}
	fmt.Printf("Concurrent Copy Time: %v\n", time.Since(start))
}
