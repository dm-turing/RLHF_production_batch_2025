package compressor

import (
	"compress/flate"
	"io"
	"os"
)

func CompressDeflate(source, destination string) error {
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

	writer, err := flate.NewWriter(outFile, flate.DefaultCompression)
	if err != nil {
		return err
	}
	defer writer.Close()

	_, err = io.Copy(writer, inFile)
	return err
}
