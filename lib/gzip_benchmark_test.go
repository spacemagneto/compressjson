package lib

import "testing"

var globalGzip *GZIPTranscoder

func init() {
	globalGzip = NewGZIPTranscoderWithPool()
}

// BenchmarkGzip_Compress_Small measures compression performance on small typical payloads.
// This represents the most common real-world case (events, messages, cache entries).
// BestCompression
// BenchmarkGzip_Compress_Small-8             84250             13008 ns/op           8.61 MB/s         240 B/op          3 allocs/op
// BestSpeed
// BenchmarkGzip_Compress_Small-8            865723              1384 ns/op          80.90 MB/s         498 B/op          4 allocs/op
// ConstantCompression
// BenchmarkGzip_Compress_Small-8            845385              1388 ns/op          80.71 MB/s         496 B/op          4 allocs/op
// HuffmanOnly
// BenchmarkGzip_Compress_Small-8            824084              1398 ns/op          80.14 MB/s         496 B/op          4 allocs/op
func BenchmarkGzip_Compress_Small(b *testing.B) {
	b.SetBytes(int64(len(smallPayload)))
	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_, _ = globalGzip.Compress(smallPayload)
	}
}

// BenchmarkGzip_Compress_Medium tests compression on medium-sized structured data.
// Common in analytics events, API responses, or moderate cache values.
// BestCompression
// BenchmarkGzip_Compress_Medium-8            69074              16122 ns/op          20.78 MB/s         384 B/op          3 allocs/op
// BestSpeed
// BenchmarkGzip_Compress_Medium-8           261393              4701 ns/op          71.27 MB/s         886 B/op          4 allocs/op
// ConstantCompression
// BenchmarkGzip_Compress_Medium-8           364970              3275 ns/op         102.30 MB/s         401 B/op          3 allocs/op
// HuffmanOnly
// BenchmarkGzip_Compress_Medium-8           345676              3450 ns/op          97.10 MB/s         401 B/op          3 allocs/op
func BenchmarkGzip_Compress_Medium(b *testing.B) {
	b.SetBytes(int64(len(mediumPayload)))
	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_, _ = globalGzip.Compress(mediumPayload)
	}
}

// BenchmarkGzip_Compress_Large measures best-case compression throughput.
// Highly repetitive text achieves excellent ratios and very high speed.
// BestCompression
// BenchmarkGzip_Compress_Large-8         4957            242 034 ns/op         537.11 MB/s        2135 B/op          5 allocs/op
// BestSpeed
// BenchmarkGzip_Compress_Large-8        43628             26219 ns/op        4958.23 MB/s        1923 B/op          5 allocs/op
// ConstantCompression
// BenchmarkGzip_Compress_Large-8              8356            142211 ns/op         914.14 MB/s      162499 B/op         11 allocs/op
// HuffmanOnly
// BenchmarkGzip_Compress_Large-8              8113            148644 ns/op         874.57 MB/s      162500 B/op         11 allocs/op
func BenchmarkGzip_Compress_Large(b *testing.B) {
	b.SetBytes(int64(len(largeCompressiblePayload)))
	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_, _ = globalGzip.Compress(largeCompressiblePayload)
	}
}

// BestCompression
// BenchmarkGZIP_Decompress_Small-8         9982855               108.9 ns/op      1028.72 MB/s         256 B/op          1 allocs/op
// BestSpeed
// BenchmarkGZIP_Decompress_Small-8         9333729               112.7 ns/op       993.47 MB/s         256 B/op          1 allocs/op
// ConstantCompression
// BenchmarkGZIP_Decompress_Small-8        10041343               108.4 ns/op      1033.19 MB/s         256 B/op          1 allocs/op
// HuffmanOnly
// BenchmarkGZIP_Decompress_Small-8         9808491               110.2 ns/op      1016.44 MB/s         256 B/op          1 allocs/op
func BenchmarkGZIP_Decompress_Small(b *testing.B) {
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

// BestCompression
// BenchmarkGZIP_Decompress_Medium-8         677472              1765 ns/op         189.83 MB/s         352 B/op          1 allocs/op
// BestSpeed
// BenchmarkGZIP_Decompress_Medium-8         668968              1750 ns/op         191.46 MB/s         352 B/op          1 allocs/op
// ConstantCompression
// BenchmarkGZIP_Decompress_Medium-8         640872              1832 ns/op         182.86 MB/s         352 B/op          1 allocs/op
// HuffmanOnly
// BenchmarkGZIP_Decompress_Medium-8         676502              1827 ns/op         183.34 MB/s         352 B/op          1 allocs/op
func BenchmarkGZIP_Decompress_Medium(b *testing.B) {
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

// BestCompression
// BenchmarkGZIP_Decompress_Large-8           45886             23791 ns/op        5464.26 MB/s      131129 B/op          1 allocs/op
// BestSpeed
// BenchmarkGZIP_Decompress_Large-8           48830             24478 ns/op        5310.80 MB/s      131126 B/op          1 allocs/op
// ConstantCompression
// BenchmarkGZIP_Decompress_Large-8           48054             24260 ns/op        5358.55 MB/s      131125 B/op          1 allocs/op
// HuffmanOnly
// BenchmarkGZIP_Decompress_Large-8           45960             24296 ns/op        5350.73 MB/s      131122 B/op          1 allocs/op
func BenchmarkGZIP_Decompress_Large(b *testing.B) {
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
