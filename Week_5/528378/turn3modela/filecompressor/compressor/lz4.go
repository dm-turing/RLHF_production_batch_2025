package compressor

import (
	"bytes"
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

	src, err := io.ReadAll(inFile)
	if err != nil {
		return err
	}
	dst, err := io.ReadAll(outFile)
	if err != nil {
		return err
	}
	return lz4Operation(bytes.NewBuffer(src), bytes.NewBuffer(dst), true)
}

func DecompressLz4(source, destination string) error {
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

	src, err := io.ReadAll(inFile)
	if err != nil {
		return err
	}
	dst, err := io.ReadAll(outFile)
	if err != nil {
		return err
	}
	return lz4Operation(bytes.NewBuffer(src), bytes.NewBuffer(dst), false)
}

func lz4Operation(src, dst io.Writer, compress bool) error {
	var err error
	srcBody, oksrc := src.(*bytes.Buffer)
	dstBody, okdst := dst.(*bytes.Buffer)
	if oksrc && okdst {
		if compress {
			writer := lz4.NewWriter(dstBody)
			defer writer.Close()
			_, err = io.Copy(writer, srcBody)
		} else {
			reader := lz4.NewReader(srcBody)
			_, err = io.Copy(dstBody, reader)
		}
	}
	return err
}
