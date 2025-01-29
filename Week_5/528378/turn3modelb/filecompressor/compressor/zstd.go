package compressor

import (
	"io"
	"os"

	"github.com/klauspost/compress/zstd"
)

func CompressZstd(source, destination string) error {
	return zstdOperation(source, destination, true)
}

func DecompressZstd(source, destination string) error {
	return zstdOperation(source, destination, false)
}

func zstdOperation(source, destination string, compress bool) error {
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
		writer, err := zstd.NewWriter(outFile)
		if err != nil {
			return err
		}
		defer writer.Close()
		_, err = io.Copy(writer, inFile)
	} else {
		reader, err := zstd.NewReader(inFile)
		if err != nil {
			return err
		}
		defer reader.Close()
		_, err = io.Copy(outFile, reader)
	}
	return err
}
