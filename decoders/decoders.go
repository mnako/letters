// Package decoders provides two functions for decoding parts of an
// email.
package decoders

import (
	"encoding/base64"
	"fmt"
	"io"
	"mime"
	"mime/quotedprintable"
	"strings"

	"golang.org/x/net/html/charset"
	"golang.org/x/text/transform"

	"github.com/mnako/letters/email"
	"github.com/rorycl/base64toraw"
)

func DecodeHeader(s string) (string, error) {
	charsetReader := func(label string, input io.Reader) (io.Reader, error) {
		enc, _ := charset.Lookup(label)
		if enc == nil {
			normalizedLabel := strings.Replace(label, "windows-", "cp", -1)
			enc, _ = charset.Lookup(normalizedLabel)
		}
		if enc == nil {
			return nil, fmt.Errorf("encoding lookup failed %s", label)
		}
		return enc.NewDecoder().Reader(input), nil
	}
	mimeDecoder := mime.WordDecoder{CharsetReader: charsetReader}
	decodedHeader, err := mimeDecoder.DecodeHeader(s)
	if err != nil {
		return "", fmt.Errorf("cannot decode MIME-word-encoded header %q: %w", s, err)
	}
	return decodedHeader, nil
}

// DecodeContent wraps the content io.Reader (from an email.Body or
// mime/multipart.Part) in either a base64 or quoted printable decoder
// if applicable. The function further wraps the reader in a transform
// character decoder if an encoding is supplied.
//
// Note that the base64 decoder "base64toraw.NewBase64ToRaw" decodes all
// base64 content to data that is base64.RawStdEncoding encoded, i.e.
// without "=" padding.
func DecodeContent(content io.Reader, ci *email.ContentInfo) io.Reader {
	var contentReader io.Reader
	switch ci.TransferEncoding {
	case "base64":
		contentReader = base64.NewDecoder(base64.RawStdEncoding, base64toraw.NewBase64ToRaw(content))
	case "quoted-printable":
		contentReader = quotedprintable.NewReader(content)
	default:
		contentReader = content
	}
	if ci.Encoding == nil {
		ci.ExtractEncoding() // lazy load
		if ci.Encoding == nil {
			return contentReader
		}
	}
	return transform.NewReader(contentReader, ci.Encoding.NewDecoder())
}
