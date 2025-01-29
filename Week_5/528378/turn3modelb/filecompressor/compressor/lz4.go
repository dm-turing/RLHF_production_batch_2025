package compressor

import (
	"io"
	"os"

	"github.com/pierrec/lz4/v4"
)

func CompressLz4(source, destination string) error {
	return lz4Operation(source, destination, true)
}

func DecompressLz4(source, destination string) error {
	return lz4Operation(source, destination, false)
}

func lz4Operation(source, destination string, compress bool) error {
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
		writer := lz4.NewWriter(outFile)
		defer writer.Close()
		_, err = io.Copy(writer, inFile)
	} else {
		reader := lz4.NewReader(inFile)
		_, err = io.Copy(outFile, reader)
	}
	return err
}
