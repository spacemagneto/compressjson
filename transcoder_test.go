package compressjson

import (
	"encoding/base64"
	"strings"
	"testing"

	"github.com/davecgh/go-spew/spew"
	"github.com/goccy/go-json"
	"github.com/stretchr/testify/assert"
)

type user struct {
	ID    int    `json:"id"`
	Name  string `json:"name,omitempty"`
	Email string `json:"email,omitempty"`
	Age   int    `json:"age,omitempty"`
}

func TestPipelineTranscoderWithError(t *testing.T) {
	t.Parallel()

	tr := NewTranscoder[user]()

	cases := []struct {
		name          string
		input         user
		modifyEncoded func(string) string
		wantErr       bool
	}{
		{name: "Full struct", input: user{ID: 42, Name: "Alice", Email: "alice@example.com", Age: 30}, wantErr: false},
		{name: "Partial struct with omitted fields", input: user{ID: 1, Name: "Bob"}, wantErr: false},
		{name: "Empty struct", input: user{}, wantErr: false},
		{name: "Only email field", input: user{ID: 10, Email: "charlie@x.ai"}, wantErr: false},
		{
			name:          "Invalid Base64 input",
			modifyEncoded: func(_ string) string { return "!!! this is not base64 !!!" },
			wantErr:       true,
		},
		{
			name:          "Base64 with invalid padding",
			input:         user{ID: 999},
			modifyEncoded: func(s string) string { return s[:len(s)-2] + "%%" },
			wantErr:       true,
		},
		{
			name:  "Corrupted Zstd after Base64",
			input: user{ID: 2, Name: "Bob"},
			modifyEncoded: func(s string) string {
				dec, _ := base64.StdEncoding.DecodeString(s)
				if len(dec) > 20 {
					dec = dec[:20]
				}

				return base64.StdEncoding.EncodeToString(dec)
			},
			wantErr: true,
		},
		{
			name:  "Bit-flipped Zstd payload",
			input: user{ID: 3},
			modifyEncoded: func(s string) string {
				dec, _ := base64.StdEncoding.DecodeString(s)
				if len(dec) > 15 {
					dec = append([]byte(nil), dec...)
					dec[15] ^= 0xff
				}
				return base64.StdEncoding.EncodeToString(dec)
			},
			wantErr: true,
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {

			var encoded string
			var err error

			if tt.modifyEncoded != nil && strings.Contains(tt.name, "Invalid Base64 input") {
				encoded = tt.modifyEncoded("")
			} else {
				encoded, err = tr.Encode(tt.input)
				assert.NoError(t, err, "Encode must succeed on valid input")
				assert.NotEmpty(t, encoded, "Encode must return non-empty string")
			}

			if tt.modifyEncoded != nil {
				encoded = tt.modifyEncoded(encoded)
			}

			decoded, err := tr.Decode(encoded)

			if tt.wantErr {
				assert.Error(t, err, "Decode must return an error when input is invalid or corrupted")
				return
			}

			assert.NoError(t, err, "Decode must succeed when given valid encoded data")
			assert.Equal(t, tt.input, decoded, "Failed: decoded value does not match original input")
		})
	}
}

// TestPipelineTranscoderJSONMarshalError verifies that Encode returns a wrapped
// error when the JSON marshaller fails. The test uses a value that cannot be
// JSON-marshaled (it contains an unmarshalled function field) to trigger the
// JSON marshaling failure path inside PipelineTranscoder.Encode.
func TestPipelineTranscoderJSONMarshalError(t *testing.T) {
	pt := NewTranscoder[struct{ F func() }]()
	src := struct{ F func() }{F: func() {}}

	_, err := pt.Encode(src)

	assert.Error(t, err, "Expected error when JSON marshal fails")
}

func TestPipelineTranscoder(t *testing.T) {
	t.Parallel()

	t.Run("One", func(t *testing.T) {
		u := []user{
			{ID: 1, Name: "Name", Email: "email@gmail.com", Age: 11},
			{ID: 2, Name: "Name1", Email: "email1@gmail.com", Age: 22},
			{ID: 3, Name: "Name2", Email: "email2@gmail.com", Age: 33},
			{ID: 4, Name: "Name3", Email: "email3@gmail.com", Age: 44},
		}

		tr := NewTranscoder[[]user]()

		str, err := tr.Encode(u)
		assert.NoError(t, err)
		assert.NotEmpty(t, str)
		spew.Dump(str)

		data, _ := json.Marshal(u)
		spew.Dump(string(data))
	})
}
