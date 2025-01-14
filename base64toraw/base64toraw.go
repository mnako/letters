// package base64translator translates base64.StdEncoded data to
// base64.RawStdEncoded to allow both to be read by
// base64.RawStdEncoding decoding.
//
// For example:
//
//	b64 := NewBase64ToRaw(bytes.NewReader(encodedBytes))
//	b, err := io.ReadAll(base64.NewDecoder(base64.RawStdEncoding, b64))
package base64toraw

import (
	"io"
)

// Base64ToRaw translates base64.StdEncoded into base64.RawStdDecoded
// data.
type Base64ToRaw struct {
	wrapped io.Reader
}

func NewBase64ToRaw(r io.Reader) *Base64ToRaw {
	return &Base64ToRaw{
		wrapped: r,
	}
}

// Read is takend from encoding/base64/base64.go
func (b *Base64ToRaw) Read(p []byte) (int, error) {
	n, err := b.wrapped.Read(p)
	for n > 0 {
		offset := 0
		for i, b := range p[:n] {
			if b != '\r' && b != '\n' && b != '=' {
				if i != offset {
					p[offset] = b
				}
				offset++
			}
		}
		if offset > 0 {
			return offset, err
		}
		// Previous buffer entirely whitespace, read again
		n, err = b.wrapped.Read(p)
	}
	return n, err
}
