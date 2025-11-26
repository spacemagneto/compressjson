package lib

import (
	"github.com/klauspost/compress/zstd"
)

var (
	// Shared global Z - standard encoder instance used by all ZSTDTranscoder objects.
	// Initialized once at startup with maximum speed settings and high parallelism.
	// Thread-safe and optimized for extremely high compression throughput.
	encoder, _ = zstd.NewWriter(nil, zstd.WithEncoderLevel(zstd.SpeedFastest), zstd.WithEncoderConcurrency(10))

	// Shared global Z - standard decoder instance used by all ZSTDTranscoder objects.
	// Pre-configured with multiple worker threads to achieve peak decompression performance.
	// Thread-safe and designed for ultra-fast decompression in hot paths.
	decoder, _ = zstd.NewReader(nil, zstd.WithDecoderConcurrency(4))
)

// ZSTDTranscoder provides zero-allocation, high-throughput Zstandard compression and decompression
// by reusing globally pre-configured encoder and decoder instances.
// It has no internal state - all heavy lifting is done by the shared, thread-safe global objects.
// This design eliminates per-instance initialization overhead and maximizes performance
// in hot paths (caching, messaging, logging, etc.) while remaining safe for concurrent use.
type ZSTDTranscoder struct{}

// NewZSTDTranscoder returns a lightweight transcoder instance that operates on the global
// pre-initialized encoder and decoder. No allocation or setup is performed - the instance
// is immediately ready for use and can be safely shared across the entire application.
func NewZSTDTranscoder() *ZSTDTranscoder {
	return &ZSTDTranscoder{}
}

// Compress compresses the input data in a single fast operation using the shared global encoder.
// It uses EncodeAll which is optimized for complete in-memory buffers and produces
// a fully framed, independently decompression output.
// The result is allocated once and returned - no internal buffers are reused.
func (t *ZSTDTranscoder) Compress(src []byte) ([]byte, error) {
	// EncodeAll appends to the provided dest buffer; we pass a zero-length slice with capacity
	// to avoid extra allocations while still getting a fresh result slice.
	// Docs: https://github.com/klauspost/compress/tree/master/zstd#blocks
	return encoder.EncodeAll(src, make([]byte, 0, len(src))), nil
}

// Decompress accepts Z - standard-compressed data and returns the original uncompressed bytes.
// It uses the convenient DecodeAll method which handles the complete decompression in a single
// operation. The destination buffer is managed internally, the returned slice is a new allocation
// owned by the caller. Any error during decompression (corrupted data, incomplete input, etc.)
// is returned to the caller.
func (t *ZSTDTranscoder) Decompress(src []byte) ([]byte, error) {
	return decoder.DecodeAll(src, nil)
}
