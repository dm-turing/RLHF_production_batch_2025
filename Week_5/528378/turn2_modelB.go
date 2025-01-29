package main  
import (  
    "fmt"
    "io"
    "log"
    "os"
    "path/filepath"
    "strings"
    "time"

    "github.com/pkg/errors"
    "golang.org/x/crypto/crc32"
    "file-compressor/compressor"
)

func main() {
    if len(os.Args) != 4 {
        fmt.Println("Usage: file-compressor <algorithm> <source_file_path> <destination_file_path>")
        fmt.Println("Supported algorithms: gzip, deflate, bzip2, lz4, zstd, decompress")
        return
    }

    algorithm := os.Args[1]
    source := os.Args[2]
    destination := os.Args[3]

    err := performOperation(algorithm, source, destination)
    if err != nil {
        log.Fatalf("Error: %v\n", err)
    }
}

func performOperation(algorithm, source, destination string) error {
    // Check if the destination file already exists and remove it if it does
    if _, err := os.Stat(destination); err == nil {
        os.Remove(destination)
    }

    // Display file sizes before operation
    sourceFileInfo, err := os.Stat(source)
    if err != nil {
        return errors.Wrap(err, "error getting source file info")
    }
    fmt.Printf("File size before operation: %d bytes\n", sourceFileInfo.Size())

    // Perform compression/decompression based on the algorithm
    var err error
    switch algorithm {
    case "gzip", "deflate", "bzip2", "lz4", "zstd":
        err = compressor.Compress(algorithm, source, destination)
    case "decompress":
        err = compressor.Decompress(destination)
    default:
        return errors.New("unsupported operation")
    }

    if err != nil {
        return errors.Wrap(err, "error performing operation")
    }

    // Display file sizes after operation
    destinationFileInfo, err := os.Stat(destination)
    if err != nil {
        return errors.Wrap(err, "error getting destination file info")
    }
    fmt.Printf("File size after operation: %d bytes\n", destinationFileInfo.Size())

    // Perform CRC check
    err = performCRCCheck(source, destination)
    if err != nil {
        return errors.Wrap(err, "CRC check failed")
    }
    fmt.Println("CRC check passed!")

    return nil
}

func performCRCCheck(source, destination string) error {
    sourceFile, err := os.Open(source)
    if err != nil {
        return errors.Wrap(err, "error opening source file for CRC check")
    }
    defer sourceFile.Close()

    destinationFile, err := os.Open(destination)
    if err != nil {
        return errors.Wrap(err, "error opening destination file for CRC check")
    }
    defer destinationFile.Close()

    sourceCRC := crc32.ChecksumIEEE(readFile(sourceFile))
    destinationCRC := crc32.ChecksumIEEE(readFile(destinationFile))

    if sourceCRC != destinationCRC {
        return errors.New("CRC checksum mismatch")
    }

    return nil
}

func readFile(file *os.File) []byte {
    fileInfo, err := file.Stat()
    if err != nil {