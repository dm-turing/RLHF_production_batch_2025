package compressor

import (
	"errors"
)

func Compress(algorithm, source, destination string) error {
	switch algorithm {
	case "gzip":
		return CompressGzip(source, destination)
	case "deflate":
		return CompressDeflate(source, destination)
	case "bzip2":
		return CompressBzip2(source, destination)
	case "lz4":
		return CompressLz4(source, destination)
	case "zstd":
		return CompressZstd(source, destination)
	default:
		return errors.New("unsupported compression algorithm")
	}
}
