package compressor

import (
	"compress/gzip"
	"io"
	"os"
)

func CompressGzip(source, destination string) error {
	return gzipOperation(source, destination, true)
}

func DecompressGzip(source, destination string) error {
	return gzipOperation(source, destination, false)
}

func gzipOperation(source, destination string, compress bool) error {
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

	if compress {
		writer := gzip.NewWriter(outFile)
		defer writer.Close()
		_, err = io.Copy(writer, inFile)
	} else {
		reader, err := gzip.NewReader(inFile)
		if err != nil {
			return err
		}
		defer reader.Close()
		_, err = io.Copy(outFile, reader)
	}
	return err
}
