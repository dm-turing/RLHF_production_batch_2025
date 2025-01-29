package main  
import (  
    "compress/flate"
    "compress/gzip"
    "flag"
    "fmt"
    "io"
    "log"
    "os"
    "strings"
    "github.com/ulikunitz/xz"
    "github.com/klauspost/compress/zstd"
    "github.com/pierrec/lz4"
)

func main() {  
    // Define command-line flags
    inputFile := flag.String("i", "", "Input file path")
    outputFile := flag.String("o", "", "Output file path")
    compressionMethod := flag.String("m", "gzip", "Compression method (gzip, deflate, bzip2, lz4, zstd)")
    flag.Parse()

    // Validate input and output file paths
    if *inputFile == "" {
        flag.PrintDefaults()
        os.Exit(1)
    }
    if _, err := os.Stat(*inputFile); os.IsNotExist(err) {
        fmt.Printf("Error: Input file not found: %s\n", *inputFile)
        os.Exit(1)
    }
    if *outputFile == "" {
        *outputFile = *inputFile + ".compressed"
    }

    // Create a compressor function map
    compressors := map[string]func(io.Writer, io.Reader) error{
        "gzip":    gzipCompress,
        "deflate": deflateCompress,
        "bzip2":   bzip2Compress,
        "lz4":    lz4Compress,
        "zstd":    zstdCompress,
    }

    // Check if the chosen compression method is valid
    if _, ok := compressors[*compressionMethod]; !ok {
        fmt.Printf("Error: Invalid compression method: %s. Valid methods: %s\n", *compressionMethod, strings.Join(getCompressorNames(compressors), ", "))
        os.Exit(1)
    }

    // Perform compression
    err := compress(*inputFile, *outputFile, compressors[*compressionMethod])
    if err != nil {
        log.Fatalf("Compression failed: %v", err)
    }
    fmt.Println("Compression successful!")
}

func getCompressorNames(compressors map[string]func(io.Writer, io.Reader) error) []string {
    var names []string
    for name := range compressors {
        names = append(names, name)
    }
    return names
}

// Compression functions
func gzipCompress(w io.Writer, r io.Reader) error {
    gzw := gzip.NewWriter(w)
    defer gzw.Close()