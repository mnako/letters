package letters

import (
	"encoding/base64"
	"fmt"
	"io"
	"mime"
	"mime/multipart"
	"mime/quotedprintable"
	"strings"

	"golang.org/x/net/html/charset"
	"golang.org/x/text/encoding"
	"golang.org/x/text/transform"

	"github.com/mnako/letters/base64toraw"
)

func decodeHeader(s string) (string, error) {
	CharsetReader := func(label string, input io.Reader) (io.Reader, error) {
		enc, _ := charset.Lookup(label)
		if enc == nil {
			normalizedLabel := strings.Replace(label, "windows-", "cp", -1)
			enc, _ = charset.Lookup(normalizedLabel)
		}
		if enc == nil {
			return nil, fmt.Errorf(
				"letters.decoders.decodeHeader.CharsetReader: cannot lookup encoding %s",
				label)
		}
		return enc.NewDecoder().Reader(input), nil
	}
	mimeDecoder := mime.WordDecoder{CharsetReader: CharsetReader}
	decodedHeader, err := mimeDecoder.DecodeHeader(s)
	if err != nil {
		return decodedHeader, fmt.Errorf(
			"letters.decoders.decodeHeader: cannot decode MIME-word-encoded header %q: %w",
			s,
			err)
	}
	return decodedHeader, nil
}

// decodeContent wraps the content io.Reader (from either a net/mail.Message.Body or
// mime/multipart.Part) in either a base64 or quoted printable decoder if applicable. Note that the
// base64 decoder "base64toraw.NewBase64ToRaw" decodes all base64 content to data that is
// base64.RawStdEncoding encoded, i.e. without "=" padding. The function further wraps the reader in
// a transform character decoder if an encoding is supplied.
func decodeContent(content io.Reader, e encoding.Encoding, cte ContentTransferEncoding) io.Reader {
	var contentReader io.Reader

	switch cte {
	case cteBase64:
		contentReader = base64.NewDecoder(
			base64.RawStdEncoding,
			base64toraw.NewBase64ToRaw(content), // normalise to raw encoding
		)
	case cteQuotedPrintable:
		contentReader = quotedprintable.NewReader(content)
	default:
		contentReader = content
	}
	if e == nil {
		return contentReader
	}
	return transform.NewReader(contentReader, e.NewDecoder())
}

func decodeInlineFile(part *multipart.Part, cte ContentTransferEncoding) (InlineFile, error) {
	var ifl InlineFile

	cid, err := decodeHeader(part.Header.Get("Content-Id"))
	if err != nil {
		return ifl, fmt.Errorf(
			"letters.decoders.decodeInlineFile: cannot decode Content-ID header for inline attachment: %w",
			err)
	}

	decoderReader := decodeContent(part, nil, cte)

	ifl.ContentID = strings.Trim(cid, "<>")

	ifl.ContentType, err = parseContentTypeHeader(part.Header.Get("Content-Type"))
	if err != nil {
		return ifl, fmt.Errorf(
			"letters.decoders.decodeInlineFile: cannot parse Content-Type of inline attachment: %w",
			err)
	}

	ifl.ContentDisposition, err = parseContentDisposition(part.Header.Get("Content-Disposition"))
	if err != nil {
		return ifl, fmt.Errorf(
			"letters.decoders.decodeInlineFile: cannot parse Content-Disposition of inline attachment: %w",
			err)
	}

	ifl.DataReader = decoderReader

	return ifl, nil
}

func decodeAttachmentFileFromBody(body io.Reader, headers Headers, cte ContentTransferEncoding) (AttachedFile, error) {
	var afl AttachedFile

	decoderReader := decodeContent(body, nil, cte)

	afl.ContentType = headers.ContentType
	afl.ContentDisposition = headers.ContentDisposition

	afl.DataReader = decoderReader

	return afl, nil
}

func decodeAttachedFileFromPart(part *multipart.Part, cte ContentTransferEncoding) (AttachedFile, error) {
	var afl AttachedFile

	decodedReader := decodeContent(part, nil, cte)

	var err error
	afl.ContentType, err = parseContentTypeHeader(part.Header.Get("Content-Type"))
	if err != nil {
		return afl, fmt.Errorf(
			"letters.decoders.decodeAttachedFileFromPart: cannot parse Content-Type of attached file: %w",
			err)
	}

	afl.ContentDisposition, err = parseContentDisposition(part.Header.Get("Content-Disposition"))
	if err != nil {
		return afl, fmt.Errorf(
			"letters.decoders.decodeAttachedFileFromPart: cannot parse Content-Disposition of attached file: %w",
			err)
	}

	afl.DataReader = decodedReader

	return afl, nil
}
