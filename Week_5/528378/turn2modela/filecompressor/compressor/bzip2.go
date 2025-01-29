package compressor

import (
	"io"
	"os"

	"github.com/dsnet/compress/bzip2"
)

func CompressBzip2(source, destination string) error {
	return bzip2Operation(source, destination, true)
}

func DecompressBzip2(source, destination string) error {
	return bzip2Operation(source, destination, false)
}

func bzip2Operation(source, destination string, compress bool) error {
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
		writer, err := bzip2.NewWriter(outFile, nil)
		if err != nil {
			return err
		}
		defer writer.Close()
		_, err = io.Copy(writer, inFile)
	} else {
		reader, err := bzip2.NewReader(inFile, nil)
		if err != nil {
			return err
		}
		defer reader.Close()
		_, err = io.Copy(outFile, reader)
	}
	return err
}
