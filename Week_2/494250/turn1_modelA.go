package main

import (
	"archive/tar"
	"archive/zip"
	"fmt"
	"io"
	"log"
	"os"
	"path"
	"path/filepath"
	"strings"
)

const tarSuffix = ".tar"
const tarGzSuffix = ".tar.gz"

func archiveFiles(srcDir, dstDir string, isCompressed bool) error {
	srcDir, err := filepath.Abs(srcDir)
	if err != nil {
		return err
	}

	dstFile := path.Join(dstDir, strings.Split(srcDir, "/")[-1]+tarSuffix)
	if isCompressed {
		dstFile += tarGzSuffix
	}

	srcFile, err := os.Open(srcDir)
	if err != nil {
		return err
	}
	defer srcFile.Close()

	var writer io.Writer
	var tarWriter *tar.Writer

	if isCompressed {
		gzFile, err := os.Create(dstFile)
		if err != nil {
			return err
		}
		defer gzFile.Close()

		gzWriter := zip.NewWriter(gzFile)
		defer gzWriter.Close()

		tarWriter = tar.NewWriter(gzWriter)
		writer = gzWriter
	} else {
		tarWriter = tar.NewWriter(srcFile)
		writer = tarWriter
	}

	err = filepath.Walk(srcDir, func(p string, fi os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		h := &tar.Header{
			Name:  strings.TrimPrefix(p, srcDir),
			Mode:  int64(fi.Mode()),
			Uid:   fi.Uid(),
			Gid:   fi.Gid(),
			Size:  fi.Size(),
			Mtime: fi.ModTime(),
		}

		err = tarWriter.WriteHeader(h)
		if err != nil {
			return err
		}

		if fi.IsDir() {
			return nil
		}

		file, err := os.Open(p)
		if err != nil {
			return err
		}
		defer file.Close()

		_, err = io.Copy(tarWriter, file)
		return err
	})

	if err != nil {
		return err
	}

	err = tarWriter.Close()
	if err != nil {
		return err
	}

	fmt.Printf("Archived to: %s\n", dstFile)
	return nil
}

func main() {
	if len(os.Args) != 4 {
		log.Fatal("Usage: go run main.go <source directory> <destination directory> <compress: yes/no>")
	}

	srcDir := os.Args[1]
	dstDir := os.Args[2]
	compressStr := os.Args[3]

	isCompressed := strings.ToLower(compressStr) == "yes"

	err := archiveFiles(srcDir, dstDir, isCompressed)
	if err != nil {
		log.Fatalf("Error archiving: %v", err)
	}
}
