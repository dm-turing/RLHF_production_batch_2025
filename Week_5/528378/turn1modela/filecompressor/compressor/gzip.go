package compressor

import (
	"compress/gzip"
	"io"
	"os"
)

func CompressGzip(source, destination string) error {
	inFile, err := os.Open(source)
	if err != nil {
		return err
	}
	defer inFile.Close()

	outFile, err := os.Create(destination)
	if err != nil {
		return err
	}
	defer outFile.Close()

	writer := gzip.NewWriter(outFile)
	defer writer.Close()

	_, err = io.Copy(writer, inFile)
	return err
}
