package compressor

import (
	"bytes"
	"compress/gzip"
	"io"
	"os"
)

func CompressGzip(source, destination string) error {
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
	return gzipOperation(bytes.NewBuffer(src), bytes.NewBuffer(dst), true)
}

func DecompressGzip(source, destination string) error {
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
	return gzipOperation(bytes.NewBuffer(src), bytes.NewBuffer(dst), false)
}

func gzipOperation(src, dst io.Writer, compress bool) error {
	srcBody, oksrc := src.(*bytes.Buffer)
	dstBody, okdst := dst.(*bytes.Buffer)
	if oksrc && okdst {
		if compress {
			writer := gzip.NewWriter(dstBody)
			defer writer.Close()
			_, err := io.Copy(writer, srcBody)
			return err
		} else {
			reader, err := gzip.NewReader(srcBody)
			if err != nil {
				return err
			}
			defer reader.Close()
			_, err = io.Copy(dstBody, reader)
			return err
		}
	}
	return nil
}
