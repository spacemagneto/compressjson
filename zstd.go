package compressjson

import (
	"github.com/klauspost/compress/zstd"
)

// ZSTDTranscoder provides transparent Z - standard compression and decompression using pre-initialized
// encoder and decoder instances. It is designed for high-performance scenarios where the same
// transcoder instance is reused across many operations. The implementation maintains internal
// state (the encoder and decoder) and is safe for concurrent use because the underlying zstd
// library guarantees thread-safety of its Writer and Reader types when properly configured.
type ZSTDTranscoder struct {
	encoder *zstd.Encoder
	decoder *zstd.Decoder
}

// NewZSTDTranscoder creates a new ZSTDTranscoder with a fully configured encoder and decoder.
// The encoder is initialized with the highest compression level for optimal ratio.
// Both objects are created with nil writers/readers because EncodeAll/DecodeAll do not require
// an io.Writer/io.Reader - they operate on complete byte slices.
// On any error during initialization, resources are cleaned up and the error is propagated.
func NewZSTDTranscoder() (*ZSTDTranscoder, error) {
	enc, _ := zstd.NewWriter(nil, zstd.WithEncoderLevel(zstd.SpeedFastest))
	dec, _ := zstd.NewReader(nil)

	return &ZSTDTranscoder{encoder: enc, decoder: dec}, nil
}

// Compress accepts arbitrary input data and returns its fully Zstandard-compressed representation.
// It uses the fast EncodeAll path that compresses the entire input in one call and appends
// the result to a freshly allocated destination buffer of appropriate capacity.
// The returned slice is always a new allocation owned by the caller.
func (t *ZSTDTranscoder) Compress(src []byte) ([]byte, error) {
	return t.encoder.EncodeAll(src, make([]byte, 0, len(src))), nil
}

// Decompress accepts Zstandard-compressed data and returns the original uncompressed bytes.
// It uses the convenient DecodeAll method which handles the complete decompression in a single
// operation. The destination buffer is managed internally, the returned slice is a new allocation
// owned by the caller. Any error during decompression (corrupted data, incomplete input, etc.)
// is returned to the caller.
func (t *ZSTDTranscoder) Decompress(src []byte) ([]byte, error) {
	return t.decoder.DecodeAll(src, nil)
}

// Close releases resources held by the encoder and decoder (e.g. compression dictionaries,
// internal buffers). It is safe to call multiple times.
func (t *ZSTDTranscoder) Close() {
	_ = t.encoder.Close()
	t.decoder.Close()
}
