
package main  
import (  
    "fmt"
    "io"
    "log"
    "os"
    "path/filepath"
    "sync"
    "time"
    "crypto/sha256"
    "encoding/hex"
)

const (
    progressInterval = time.Second
)

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

    hash := sha256.New()
    multiWriter := io.MultiWriter(destinationFile, hash)

    _, err = io.Copy(multiWriter, sourceFile)
    if err != nil {
        return err
    }

    fileHash := hex.EncodeToString(hash.Sum(nil))
    return verifyFileHash(src, fileHash)
}

func verifyFileHash(path string, expectedHash string) error {
    file, err := os.Open(path)
    if err != nil {
        return err
    }
    defer file.Close()

    hash := sha256.New()
    _, err = io.Copy(hash, file)
    if err != nil {
        return err
    }

    calculatedHash := hex.EncodeToString(hash.Sum(nil))
    if calculatedHash != expectedHash {
        return fmt.Errorf("Hash mismatch for %s: expected %s, got %s", path, expectedHash, calculatedHash)
    }

    return nil
}

func worker(wg *sync.WaitGroup, jobs <-chan [2]string, progress chan<- struct{}) {
    defer wg.Done()
    for pair := range jobs {
        src := pair[0]
        dst := pair[1]
        err := copyFile(src, dst)
        if err != nil {
            log.Printf("Failed to copy %s to %s: %v", src, dst, err)
        }
        progress <- struct{}{}
    }
}

func progressBar(totalFiles int, progress <-chan struct{}) {
    copiedFiles := 0
    ticker := time.NewTicker(progressInterval)
    defer ticker.Stop()

    for {
        select {
        case <-progress:
            copiedFiles++
            fmt.Printf("\rCopied files: %d/%d", copiedFiles, totalFiles)
        case <-ticker.C:
            fmt.Printf("\rCopied files: %d/%d", copiedFiles, totalFiles)
        }
    }
}

func main() {  
    if len(os.Args) != 4 {  
        fmt.Println("Usage: go run main.go <source_directory> <destination_directory> <number_of_workers>")  
        return  
    }  

    srcDir := os.Args[1]  
    dstDir := os.Args[2]
    numWorkers := atoi(os.Args[3])

    // Create destination directory if it doesn't exist
    if err := os.MkdirAll(dstDir, 0755); err != nil {
        log.Fatalf("Error creating destination directory: %v", err)
    }

    files, err := ioutil.ReadDir(srcDir)  
    if err != nil {  
        log.Fatalf("Error reading source directory: %v", err)  
    }  

    totalFiles := len(files)
    jobs := make(chan [2]string, totalFiles)
    progress := make(chan struct{}, totalFiles)
    var wg sync.WaitGroup  

    for w := 0; w < numWorkers; w++ {  
        wg.Add(1)  
        go worker(&wg, jobs, progress)  
    }  

    go progressBar(totalFiles, progress)

    for _, file := range files {  
        srcPath := filepath.Join(srcDir, file.Name())
