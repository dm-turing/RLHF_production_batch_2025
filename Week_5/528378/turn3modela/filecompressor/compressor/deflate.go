package compressor

import (
	"bytes"
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

	src, err := io.ReadAll(inFile)
	if err != nil {
		return err
	}
	dst, err := io.ReadAll(outFile)
	if err != nil {
		return err
	}
	return deflateOperation(bytes.NewBuffer(src), bytes.NewBuffer(dst), true)
}

func DecompressDeflate(source, destination string) error {
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
	return deflateOperation(bytes.NewBuffer(src), bytes.NewBuffer(dst), false)
}

func deflateOperation(src, dst io.Writer, compress bool) error {
	var err error
	srcBody, oksrc := src.(*bytes.Buffer)
	dstBody, okdst := dst.(*bytes.Buffer)
	if oksrc && okdst {
		if compress {
			writer, er := flate.NewWriter(dstBody, flate.DefaultCompression)
			if er != nil {
				return err
			}
			defer writer.Close()
			_, err = io.Copy(writer, srcBody)
		} else {
			reader := flate.NewReader(srcBody)
			defer reader.Close()
			_, err = io.Copy(dstBody, reader)
		}
	}
	return err
}
