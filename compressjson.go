package compressjson

// Transcoder defines a generic interface for bidirectional conversion between
// a value of type T and its string representation.
//
// Implementations may perform serialization, compression, encryption, encoding,
// or any combination thereof — the interface makes no assumptions about the
// internal steps. The only requirements are type safety, correct round-trip behavior
// when possible, and proper resource cleanup via Close.
type Transcoder[T any] interface {
	// Encode converts a value of type T into a string.
	// The resulting string should be safe for storage or transmission.
	Encode(T) (string, error)

	// Decode reconstructs a value of type T from a string previously produced by Encode.
	// Returns the zero value of T and an error if decoding fails.
	Decode(string) (T, error)
}
