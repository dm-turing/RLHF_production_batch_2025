package compressor

import (
	"io"
	"os"

	"github.com/pierrec/lz4/v4"
)

func CompressLz4(source, destination string) error {
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

	writer := lz4.NewWriter(outFile)
	defer writer.Close()

	_, err = io.Copy(writer, inFile)
	return err
}
