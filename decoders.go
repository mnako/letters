package letters

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"io"
	"io/ioutil"
	"mime"
	"mime/multipart"
	"mime/quotedprintable"
	"strings"

	"golang.org/x/net/html/charset"
	"golang.org/x/text/encoding"
	"golang.org/x/text/transform"
)

func decodeHeader(s string) (string, error) {
	CharsetReader := func(label string, input io.Reader) (io.Reader, error) {
		normalized := strings.Replace(label, "windows-", "cp", -1)
		enc, _ := charset.Lookup(normalized)
		if enc == nil {
			// no encoder found for normalized label,
			// try to lookup using the original value
			enc, _ = charset.Lookup(label)
			if enc == nil {
				return nil, fmt.Errorf(
					"letters.decoders.decodeHeader: cannot find MIME-word-encoded for %s", label)
			}
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

func decodeContent(
	content io.Reader,
	e encoding.Encoding,
	cte ContentTransferEncoding,
) (io.Reader, error) {
	var contentReader io.Reader
	contentBytes, err := ioutil.ReadAll(content)
	if err != nil && err != io.ErrUnexpectedEOF {
		return nil, fmt.Errorf(
			"letters.decoders.decodeContent: cannot decode content: %w",
			err)
	}

	switch cte {
	case cteBase64:
		decoded := base64.NewDecoder(base64.StdEncoding, bytes.NewReader(contentBytes))
		b, err := ioutil.ReadAll(decoded)
		if err == io.ErrUnexpectedEOF {
			decoded = base64.NewDecoder(base64.RawStdEncoding, bytes.NewReader(contentBytes))
			b, err = ioutil.ReadAll(decoded)
			if err != nil {
				return nil, fmt.Errorf(
					"letters.decoders.decodeContent: cannot decode raw-std-base64-encoded content: %w",
					err)
			}
		} else if err != nil {
			return nil, fmt.Errorf(
				"letters.decoders.decodeContent: cannot decode std-base64-encoded content: %w",
				err)
		}
		contentReader = bytes.NewReader(b)
	case cteQuotedPrintable:
		decoded := quotedprintable.NewReader(bytes.NewReader(contentBytes))
		b, err := ioutil.ReadAll(decoded)
		if err != nil {
			return nil, fmt.Errorf(
				"letters.decoders.decodeContent: cannot decode quoted-printable-encoded content: %w",
				err)
		}
		contentReader = bytes.NewReader(b)
	default:
		contentReader = bytes.NewReader(contentBytes)
	}

	if e != nil {
		contentReader = transform.NewReader(
			contentReader,
			e.NewDecoder(),
		)
	}

	return contentReader, nil
}

func decodeInlineFile(part *multipart.Part, cte ContentTransferEncoding) (InlineFile, error) {
	var ifl InlineFile

	cid, err := decodeHeader(part.Header.Get("Content-Id"))
	if err != nil {
		return ifl, fmt.Errorf(
			"letters.decoders.decodeInlineFile: cannot decode Content-ID header for inline attachment: %w",
			err)
	}

	decoded, err := decodeContent(part, nil, cte)
	if err != nil {
		return ifl, fmt.Errorf(
			"letters.decoders.decodeInlineFile: cannot decode inline attachment content: %w",
			err)
	}

	ifl.ContentID = strings.Trim(cid, "<>")
	ifl.Data, err = ioutil.ReadAll(decoded)
	if err != nil {
		return ifl, fmt.Errorf(
			"letters.decoders.decodeInlineFile: cannot read inline attachment data: %w",
			err)
	}

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

	return ifl, nil
}

func decodeAttachmentFileFromBody(body io.Reader, headers Headers, cte ContentTransferEncoding) (AttachedFile, error) {
	var afl AttachedFile

	decoded, err := decodeContent(body, nil, cte)
	if err != nil {
		return afl, fmt.Errorf(
			"letters.decoders.decodeAttachmentFileFromBody: cannot decode attached file content: %w",
			err)
	}

	afl.ContentType = headers.ContentType
	afl.ContentDisposition = headers.ContentDisposition
	afl.Data, err = ioutil.ReadAll(decoded)
	if err != nil {
		return afl, fmt.Errorf(
			"letters.decoders.decodeAttachmentFileFromBody: cannot read attached file data: %w",
			err)
	}

	return afl, nil
}

func decodeAttachedFileFromPart(part *multipart.Part, cte ContentTransferEncoding) (AttachedFile, error) {
	var afl AttachedFile

	decoded, err := decodeContent(part, nil, cte)
	if err != nil {
		return afl, fmt.Errorf(
			"letters.decoders.decodeAttachedFileFromPart: cannot decode attached file content: %w",
			err)
	}

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

	afl.Data, err = ioutil.ReadAll(decoded)
	if err != nil {
		return afl, fmt.Errorf(
			"letters.decoders.decodeAttachedFileFromPart: cannot read attached file data: %w",
			err)
	}

	return afl, nil
}
