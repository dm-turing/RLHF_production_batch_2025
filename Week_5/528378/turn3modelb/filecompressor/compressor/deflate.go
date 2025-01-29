package compressor

import (
	"compress/flate"
	"io"
	"os"
)

func CompressDeflate(source, destination string) error {
	return deflateOperation(source, destination, true)
}

func DecompressDeflate(source, destination string) error {
	return deflateOperation(source, destination, false)
}

func deflateOperation(source, destination string, compress bool) error {
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
		writer, err := flate.NewWriter(outFile, flate.DefaultCompression)
		if err != nil {
			return err
		}
		defer writer.Close()
		_, err = io.Copy(writer, inFile)
	} else {
		reader := flate.NewReader(inFile)
		defer reader.Close()
		_, err = io.Copy(outFile, reader)
	}
	return err
}
