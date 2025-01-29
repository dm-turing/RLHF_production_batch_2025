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

func Decompress(algorithm, source, destination string) error {
	switch algorithm {
	case "gzip":
		return DecompressGzip(source, destination)
	case "deflate":
		return DecompressDeflate(source, destination)
	case "bzip2":
		return DecompressBzip2(source, destination)
	case "lz4":
		return DecompressLz4(source, destination)
	case "zstd":
		return DecompressZstd(source, destination)
	default:
		return errors.New("unsupported compression algorithm")
	}
}
