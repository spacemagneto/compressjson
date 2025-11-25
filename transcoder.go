package compressjson

import "errors"

// transcoder is a concrete, high-performance implementation of transcoder[T]
// designed for low-latency, high-throughput scenarios. It is safe for concurrent use
// by multiple goroutines and reuses internal buffers and native resources across calls.
//
// Call Close when the transcoder is no longer needed to free memory.
type transcoder[T any] struct {
	jsonTranscoder     *JSONTranscoder[T]
	standardTranscoder *ZSTDTranscoder
	binaryTranscoder   *Base64Transcoder
}

// NewTranscoder creates a ready-to-use transcoder for type T.
// All internal components are initialized once and reused forever.
// The returned value satisfies Transcoder[T] and can be shared globally.
func NewTranscoder[T any]() Transcoder[T] {
	return &transcoder[T]{
		jsonTranscoder:     NewJSONTranscoder[T](),
		standardTranscoder: NewZSTDTranscoder(),
		binaryTranscoder:   NewBase64Transcoder(),
	}
}

// Encode converts a value of type T into a compact, text-safe string.
// The value is first marshaled to JSON, then compressed with Z - standard,
// and finally encoded to standard Base64. Any error aborts the process
// and returns a wrapped error with context.
func (t *transcoder[T]) Encode(src T) (string, error) {
	jsonBytes, err := t.jsonTranscoder.Marshal(src)
	if err != nil {
		return "", errors.Join(errors.New("failed to marshal JSON"), err)
	}

	compressedBytes, err := t.standardTranscoder.Compress(jsonBytes)
	if err != nil {
		return "", errors.Join(errors.New("failed to compress with Zstd"), err)
	}

	return t.binaryTranscoder.Encode(compressedBytes)
}

// Decode reconstructs the original value from the string produced by Encode.
// The process reverses the encoding steps: Base64 decoding, Z - standard decompression,
// and JSON unmarshalling. On success the original value is returned; on failure
// the zero value of T is returned along with a descriptive wrapped error.
func (t *transcoder[T]) Decode(src string) (T, error) {
	var entry T

	compressedBytes, err := t.binaryTranscoder.Decode(src)
	if err != nil {
		return entry, errors.Join(errors.New("failed to decode Base64"), err)
	}

	jsonBytes, err := t.standardTranscoder.Decompress(compressedBytes)
	if err != nil {
		return entry, errors.Join(errors.New("failed to decompress Zstd"), err)
	}

	return t.jsonTranscoder.Unmarshal(jsonBytes)
}
