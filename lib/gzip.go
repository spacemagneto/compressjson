package lib

import (
	"bytes"
	"io"
	"sync"

	kgzip "github.com/klauspost/compress/gzip"
)

var writerPool = sync.Pool{
	New: func() interface{} {
		w, _ := kgzip.NewWriterLevel(nil, kgzip.HuffmanOnly)
		return w
	},
}

var readerPool = sync.Pool{
	New: func() interface{} {
		r, _ := kgzip.NewReader(nil)
		return r
	},
}

// GZIPTranscoder provides transparent GZIP compression and decompression by pooling
// and reusing underlying gzip.Writer and gzip.Reader instances.
// This design ensures thread-safety and high performance by avoiding object allocation
// and contention on a single shared object's state.
type GZIPTranscoder struct{}

// NewGZIPTranscoderWithPool creates a new GZIPTranscoder. Since the actual compression objects
// are managed by global pools, this function is primarily for structural initialization.
func NewGZIPTranscoderWithPool() *GZIPTranscoder {
	// Since pools manage object creation, no complex setup or cleanup is needed here.
	return &GZIPTranscoder{}
}

// Compress accepts arbitrary input data and returns its fully GZIP-compressed representation.
// It acquires a gzip.Writer from the pool, resets it to a new buffer, compresses the data,
// and returns the writer to the pool.
func (t *GZIPTranscoder) Compress(src []byte) ([]byte, error) {
	// Acquire a writer from the pool
	writer := writerPool.Get().(*kgzip.Writer)
	// Ensure the writer is returned to the pool, even if an error occurs
	defer writerPool.Put(writer)

	var buf bytes.Buffer
	writer.Reset(&buf)

	// Write the uncompressed data to the GZIP stream
	if _, err := writer.Write(src); err != nil {
		return nil, err
	}

	// Close the writer to flush any buffered data and write the GZIP footer
	if err := writer.Close(); err != nil {
		return nil, err
	}

	// The buffer now holds the compressed data
	return buf.Bytes(), nil
}

// Decompress accepts GZIP-compressed data and returns the original uncompressed bytes.
// It acquires a gzip.Reader from the pool, resets it to the source buffer, decompresses the data,
// and returns the reader to the pool.
func (t *GZIPTranscoder) Decompress(src []byte) ([]byte, error) {
	// Acquire a reader from the pool
	reader := readerPool.Get().(*kgzip.Reader)
	// Ensure the reader is returned to the pool, even if an error occurs
	defer readerPool.Put(reader)

	// Use bytes.NewReader to treat the source slice as a stream
	srcReader := bytes.NewReader(src)

	// Reset the internal reader to use the source stream
	if err := reader.Reset(srcReader); err != nil {
		return nil, err
	}

	// Read all uncompressed data from the GZIP stream
	decompressed, err := io.ReadAll(reader)
	if err != nil {
		return nil, err
	}

	// NOTE: We do not call reader.Close() here as the klauspost/compress/gzip library
	// is optimized for reuse via Reset(), and Close() is not strictly required after a full read.
	// We rely on the pool to manage the lifecycle.

	return decompressed, nil
}
