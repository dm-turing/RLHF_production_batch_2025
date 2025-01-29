package compressor

import (
	"io"
	"os"

	"github.com/klauspost/compress/zstd"
)

func CompressZstd(source, destination string) error {
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

	writer, err := zstd.NewWriter(outFile)
	if err != nil {
		return err
	}
	defer writer.Close()

	_, err = io.Copy(writer, inFile)
	return err
}
