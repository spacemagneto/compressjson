package compressjson

import (
	"testing"
)

// Global transcoder reused across all benchmarks – this is the intended real-world usage pattern
var globalZSTD *ZSTDTranscoder

func init() {
	globalZSTD = NewZSTDTranscoder()
}

// BenchmarkZSTD_Compress_Small measures compression performance on small typical payloads.
// This represents the most common real-world case (events, messages, cache entries).
// SpeedBestCompression
// BenchmarkZSTD_Compress_Small-8            248517              4752 ns/op          23.57 MB/s         112 B/op          1 allocs/op
// SpeedBetterCompression
// BenchmarkZSTD_Compress_Small-8            567555              2089 ns/op          53.61 MB/s         112 B/op          1 allocs/op
// SpeedFastest
// BenchmarkZSTD_Compress_Small-8           5928451               183.2 ns/op       611.31 MB/s         336 B/op          2 allocs/op
func BenchmarkZSTD_Compress_Small(b *testing.B) {
	b.SetBytes(int64(len(smallPayload)))
	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_, _ = globalZSTD.Compress(smallPayload)
	}
}

// BenchmarkZSTD_Compress_Medium tests compression on medium-sized structured data.
// Common in analytics events, API responses, or moderate cache values.
// SpeedBestCompression
// BenchmarkZSTD_Compress_Medium-8           132940              9096 ns/op          36.83 MB/s         352 B/op          1 allocs/op
// SpeedBetterCompression
// BenchmarkZSTD_Compress_Medium-8           314161              3757 ns/op          89.16 MB/s         352 B/op          1 allocs/op
// SpeedFastest
// BenchmarkZSTD_Compress_Medium-8           431554              2768 ns/op         121.03 MB/s         352 B/op          1 allocs/op
func BenchmarkZSTD_Compress_Medium(b *testing.B) {
	b.SetBytes(int64(len(mediumPayload)))
	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_, _ = globalZSTD.Compress(mediumPayload)
	}
}

// BenchmarkZSTD_Compress_Large measures best-case compression throughput.
// Highly repetitive text achieves excellent ratios and very high speed.
// SpeedBestCompression
// BenchmarkZSTD_Compress_Large-8         7285            166550 ns/op         780.55 MB/s      131072 B/op          1 allocs/op
// SpeedBetterCompression
// BenchmarkZSTD_Compress_Large-8        79058             12804 ns/op        10152.82 MB/s     131072 B/op          1 allocs/op
// SpeedFastest
// BenchmarkZSTD_Compress_Large-8        84102             12812 ns/op        10146.44 MB/s     131072 B/op          1 allocs/op
func BenchmarkZSTD_Compress_Large(b *testing.B) {
	b.SetBytes(int64(len(largeCompressiblePayload)))
	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_, _ = globalZSTD.Compress(largeCompressiblePayload)
	}
}

// SpeedBestCompression
// BenchmarkZSTD_Decompress_Small-8         1132474              1047 ns/op         106.99 MB/s         224 B/op          1 allocs/op
// SpeedBetterCompression
// BenchmarkZSTD_Decompress_Small-8         1101944              1055 ns/op         106.13 MB/s         224 B/op          1 allocs/op
// SpeedFastest
// BenchmarkZSTD_Decompress_Small-8         9773276               109.5 ns/op      1022.91 MB/s         256 B/op          1 allocs/op
func BenchmarkZSTD_Decompress_Small(b *testing.B) {
	compressed, err := globalZSTD.Compress(smallPayload)
	if err != nil {
		b.Fatal(err)
	}

	b.SetBytes(int64(len(smallPayload)))
	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_, _ = globalZSTD.Decompress(compressed)
	}
}

// BenchmarkZSTD_Decompress_Medium tests decompression of typical medium compressed blobs.
// SpeedBestCompression
// BenchmarkZSTD_Decompress_Medium-8         694434              1729 ns/op         193.80 MB/s         352 B/op          1 allocs/op
// SpeedBetterCompression
// BenchmarkZSTD_Decompress_Medium-8         693674              1733 ns/op         193.34 MB/s         352 B/op          1 allocs/op
// SpeedFastest
// BenchmarkZSTD_Decompress_Medium-8         694903              1719 ns/op         194.91 MB/s         352 B/op          1 allocs/op
func BenchmarkZSTD_Decompress_Medium(b *testing.B) {
	compressed, err := globalZSTD.Compress(mediumPayload)
	if err != nil {
		b.Fatal(err)
	}
	b.SetBytes(int64(len(mediumPayload)))
	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_, _ = globalZSTD.Decompress(compressed)
	}
}

// BenchmarkZSTD_Decompress_Large measures peak decompression throughput.
// Zstandard decompression is extremely fast – this shows maximum achievable speed.
// SpeedBestCompression
// BenchmarkZSTD_Decompress_Large-8              30405             35964 ns/op        3614.70 MB/s      131075 B/op          1 allocs/op
// SpeedBetterCompression
// BenchmarkZSTD_Decompress_Large-8              31490             36259 ns/op        3585.28 MB/s      131086 B/op          1 allocs/op
// SpeedFastest
// BenchmarkZSTD_Decompress_Large-8              46920             24082 ns/op        5398.13 MB/s      131125 B/op          1 allocs/op
func BenchmarkZSTD_Decompress_Large(b *testing.B) {
	compressed, err := globalZSTD.Compress(largeCompressiblePayload)
	if err != nil {
		b.Fatal(err)
	}

	b.SetBytes(int64(len(largeCompressiblePayload)))
	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_, _ = globalZSTD.Decompress(compressed)
	}
}
