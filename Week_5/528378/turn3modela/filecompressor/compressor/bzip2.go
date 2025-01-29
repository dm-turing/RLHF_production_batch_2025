package compressor

import (
	"bytes"
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

	src, err := io.ReadAll(inFile)
	if err != nil {
		return err
	}
	dst, err := io.ReadAll(outFile)
	if err != nil {
		return err
	}
	return bzip2Operation(bytes.NewBuffer(src), bytes.NewBuffer(dst), true)
}

func DecompressBzip2(source, destination string) error {
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
	return bzip2Operation(bytes.NewBuffer(src), bytes.NewBuffer(dst), false)
}

func bzip2Operation(src, dst io.Writer, compress bool) error {
	var err error
	srcBody, oksrc := src.(*bytes.Buffer)
	dstBody, okdst := dst.(*bytes.Buffer)
	if oksrc && okdst {
		if compress {
			writer, err := bzip2.NewWriter(dstBody, nil)
			if err != nil {
				return err
			}
			defer writer.Close()
			_, err = io.Copy(writer, srcBody)
		} else {
			reader, err := bzip2.NewReader(srcBody, nil)
			if err != nil {
				return err
			}
			defer reader.Close()
			_, err = io.Copy(dstBody, reader)
		}
	}
	return err
}
