package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"sync"
	"time"
)

const numberOfWorkers = 10

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

func worker(wg *sync.WaitGroup, jobs <-chan [2]string) {
	defer wg.Done()
	for pair := range jobs {
		src := pair[0]
		dst := pair[1]
		err := copyFile(src, dst)
		if err != nil {
			log.Printf("Failed to copy %s to %s: %v", src, dst, err)
		}
	}
}

func serialCopyFiles(srcDir, dstDir string) error {
	files, err := ioutil.ReadDir(srcDir)
	if err != nil {
		return err
	}

	for _, file := range files {
		srcPath := filepath.Join(srcDir, file.Name())
		dstPath := filepath.Join(dstDir, file.Name())

		err := copyFile(srcPath, dstPath)
		if err != nil {
			return err
		}
	}

	return nil
}

func concurrentCopyFiles(srcDir, dstDir string) error {
	files, err := ioutil.ReadDir(srcDir)
	if err != nil {
		return err
	}

	jobs := make(chan [2]string, len(files))
	var wg sync.WaitGroup

	for w := 0; w < numberOfWorkers; w++ {
		wg.Add(1)
		go worker(&wg, jobs)
	}

	for _, file := range files {
		srcPath := filepath.Join(srcDir, file.Name())
		dstPath := filepath.Join(dstDir, file.Name())
		jobs <- [2]string{srcPath, dstPath}
	}
	close(jobs)

	wg.Wait()
	return nil
}

func main() {
	if len(os.Args) != 3 {
		fmt.Println("Usage: go run main.go <source_directory> <destination_directory>")
		return
	}

	srcDir := os.Args[1]
	dstDir := os.Args[2]

	// Measure elapsed time for serial version
	start := time.Now()
	err := serialCopyFiles(srcDir, dstDir)
	if err != nil {
		log.Fatalf("Error in serial copying: %v", err)
	}
	fmt.Printf("Serial Copy Time: %v\n", time.Since(start))

	// Measure elapsed time for concurrent version
	start = time.Now()
	err = concurrentCopyFiles(srcDir, dstDir)
	if err != nil {
		log.Fatalf("Error in concurrent copying: %v", err)
	}
	fmt.Printf("Concurrent Copy Time: %v\n", time.Since(start))
}
