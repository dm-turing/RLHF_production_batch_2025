package compressor

import (
	"bytes"
	"testing"
)

func benchmarkCompression(b *testing.B, compressFunc func(src, dst *bytes.Buffer) error) {
	data := bytes.Repeat([]byte("Lorem ipsum dolor sit amet, consectetur adipiscing elit. "), 10000)
	src := bytes.NewBuffer(data)
	compressed := new(bytes.Buffer)

	b.ReportAllocs()
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		src.Reset()
		src.Write(data)
		compressed.Reset()

		if err := compressFunc(src, compressed); err != nil {
			b.Fatal(err)
		}
	}

	b.StopTimer()

	compressionRatio := float64(src.Len()) / float64(compressed.Len())
	b.ReportMetric(compressionRatio, "compression-ratio")
}

func benchmarkDecompression(b *testing.B, compressFunc func(src, dst *bytes.Buffer) error, decompressFunc func(src, dst *bytes.Buffer) error) {
	data := bytes.Repeat([]byte("Lorem ipsum dolor sit amet, consectetur adipiscing elit. "), 10000)
	src := bytes.NewBuffer(data)
	compressed := new(bytes.Buffer)

	if err := compressFunc(src, compressed); err != nil {
		b.Fatal(err)
	}

	decompressed := new(bytes.Buffer)

	b.ReportAllocs()
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		decompressed.Reset()
		compressedReader := bytes.NewBuffer(compressed.Bytes())
		if err := decompressFunc(compressedReader, decompressed); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkGzipCompression(b *testing.B) {
	benchmarkCompression(b, func(src, dst *bytes.Buffer) error {
		return gzipOperation(src, dst, true)
	})
}

func BenchmarkGzipDecompression(b *testing.B) {
	benchmarkDecompression(b, func(src, dst *bytes.Buffer) error {
		return gzipOperation(src, dst, true)
	}, func(src, dst *bytes.Buffer) error {
		return gzipOperation(src, dst, false)
	})
}

func BenchmarkDeflateCompression(b *testing.B) {
	benchmarkCompression(b, func(src, dst *bytes.Buffer) error {
		return deflateOperation(src, dst, true)
	})
}

func BenchmarkDeflateDecompression(b *testing.B) {
	benchmarkDecompression(b, func(src, dst *bytes.Buffer) error {
		return deflateOperation(src, dst, true)
	}, func(src, dst *bytes.Buffer) error {
		return deflateOperation(src, dst, false)
	})
}

func BenchmarkBzip2Compression(b *testing.B) {
	benchmarkCompression(b, func(src, dst *bytes.Buffer) error {
		return bzip2Operation(src, dst, true)
	})
}

func BenchmarkBzip2Decompression(b *testing.B) {
	benchmarkDecompression(b, func(src, dst *bytes.Buffer) error {
		return bzip2Operation(src, dst, true)
	}, func(src, dst *bytes.Buffer) error {
		return bzip2Operation(src, dst, false)
	})
}

func BenchmarkLZ4Compression(b *testing.B) {
	benchmarkCompression(b, func(src, dst *bytes.Buffer) error {
		return lz4Operation(src, dst, true)
	})
}

func BenchmarkLZ4Decompression(b *testing.B) {
	benchmarkDecompression(b, func(src, dst *bytes.Buffer) error {
		return lz4Operation(src, dst, true)
	}, func(src, dst *bytes.Buffer) error {
		return lz4Operation(src, dst, false)
	})
}

func BenchmarkZstdCompression(b *testing.B) {
	benchmarkCompression(b, func(src, dst *bytes.Buffer) error {
		return zstdOperation(src, dst, true)
	})
}

func BenchmarkZstdDecompression(b *testing.B) {
	benchmarkDecompression(b, func(src, dst *bytes.Buffer) error {
		return zstdOperation(src, dst, true)
	}, func(src, dst *bytes.Buffer) error {
		return zstdOperation(src, dst, false)
	})
}
