package compressor

import (
	"io"
	"os"

	"github.com/dsnet/compress/bzip2"
)

func CompressBzip2(source, destination string) error {
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

	writer, err := bzip2.NewWriter(outFile, nil)
	if err != nil {
		return err
	}
	defer writer.Close()

	_, err = io.Copy(writer, inFile)
	return err
}
