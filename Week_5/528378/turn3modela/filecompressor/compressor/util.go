package compressor

import (
	"hash/crc32"
	"io"
	"os"
)

// FileSize returns the size of the file at the given path
func FileSize(path string) (int64, error) {
	info, err := os.Stat(path)
	if err != nil {
		return 0, err
	}
	return info.Size(), nil
}

// VerifyCRC checks if two files have the same CRC32 checksum
func VerifyCRC(file1Path, file2Path string) bool {
	file1, err := os.Open(file1Path)
	if err != nil {
		return false
	}
	defer file1.Close()

	file2, err := os.Open(file2Path)
	if err != nil {
		return false
	}
	defer file2.Close()

	crc1 := crc32.NewIEEE()
	crc2 := crc32.NewIEEE()

	_, err = io.Copy(crc1, file1)
	if err != nil {
		return false
	}

	_, err = io.Copy(crc2, file2)
	if err != nil {
		return false
	}

	return crc1.Sum32() == crc2.Sum32()
}
