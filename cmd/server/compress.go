package main

import (
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/klauspost/compress/zstd"
)

// CompressToFile reads from src and writes zstd-compressed data to the given file path.
// The compression level can be tuned via zstd.WithEncoderLevel (default: SpeedDefault).
func CompressToFile(src io.Reader, path string, opts ...zstd.EOption) error {
	err := os.MkdirAll(filepath.Dir(path), 0755)
	if err != nil {
		return fmt.Errorf("making target directory: %w", err)
	}

	// Create (or truncate) the destination file
	f, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("creating destination file: %w", err)
	}
	defer f.Close()

	// Wrap the file writer in a zstd encoder
	encoder, err := zstd.NewWriter(f, opts...)
	if err != nil {
		return fmt.Errorf("creating zstd encoder: %w", err)
	}

	// Stream src → zstd encoder → file
	if _, err := io.Copy(encoder, src); err != nil {
		encoder.Close()
		return fmt.Errorf("compressing data: %w", err)
	}

	// Flush and finalize the zstd frame
	if err := encoder.Close(); err != nil {
		return fmt.Errorf("finalizing zstd stream: %w", err)
	}

	return nil
}

// DecompressFromFile opens a zstd-compressed file and returns a reader for the decompressed data.
// Caller is responsible for closing the returned ReadCloser.
func DecompressFromFile(srcPath string) (io.ReadCloser, error) {
	f, err := os.Open(srcPath)
	if err != nil {
		return nil, fmt.Errorf("opening source file: %w", err)
	}

	decoder, err := zstd.NewReader(f)
	if err != nil {
		f.Close()
		return nil, fmt.Errorf("creating zstd decoder: %w", err)
	}

	// Wrap both so closing the ReadCloser cleans up decoder + file
	return &zstdReadCloser{decoder: decoder, file: f}, nil
}

type zstdReadCloser struct {
	decoder *zstd.Decoder
	file    *os.File
}

func (z *zstdReadCloser) Read(p []byte) (int, error) { return z.decoder.Read(p) }
func (z *zstdReadCloser) Close() error {
	z.decoder.Close()
	return z.file.Close()
}
