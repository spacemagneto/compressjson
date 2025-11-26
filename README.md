# compressjson: High-Performance JSON Transcoder (Zstd + Base64)

## Overview

The compressjson library provides a high-throughput, type-safe solution for converting Go structs to a compact, text-safe string representation and vice-versa.
This tool is designed for low-latency, high-throughput scenarios and is ideal for:
1. High-Speed Caching (e.g., Redis, Memcached) where decompression speed is critical for read latency. 
2. Storage Optimization for large JSON objects in databases.
3. Efficient Data Transfer in high-load, distributed systems.


## Operating Pipeline
### The encode/decode process is an optimized three-stage chain:

| Step | Operation        | Library / Settings                     | Purpose                                           |
|------|------------------|----------------------------------------|---------------------------------------------------|
| 1    | Serialization    | `goccy/go-json`                        | Ultra-fast conversion of Go structs ↔ bytes       |
| 2    | Compression      | `klauspost/compress/zstd` <br> `SpeedFastest` | Maximum compression/decompression speed           |
| 3    | Encoding         | `encoding/base64` (std lib)            | Convert binary data to transport-safe ASCII string |


# Usage
### Installation
```bash
    go get github.com/spacemagneto/compressjson
```

## Example Code

```go
package main

import (
	"log"

	"github.com/spacemagneto/compressjson"
)

type User struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
	Role string `json:"role"`
}

func main() {
	transcoder := compressjson.NewTranscoder[[]User]()

	users := []User{
		{ID: 1, Name: "Alice", Role: "admin"},
		{ID: 2, Name: "Bob", Role: "user"},
		{ID: 3, Name: "Charlie", Role: "moderator"},
	}

	encoded, err := transcoder.Encode(users)
	if err != nil {
		log.Fatalf("encode failed: %v", err)
	}

	var decoded []User
	decoded, err = transcoder.Decode(encoded)
	if err != nil {
		log.Fatalf("decode failed: %v", err)
	}

	...
}
```

### Testing Support

For users integrating and testing compressjson in various scenarios (e.g., handling specific errors, simulating compression/decompression outcomes), the repository includes a dedicated transcoder mock pkg that provides the MockTranscoder implementation. This facilitates robust unit testing and integration testing without requiring actual Zstd operations.