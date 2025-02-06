package parser

/*
The contents of extract.go are the general purpose extraction functions for:
	extractDefaultMediaType        : extracting media type
	extractContentTypeHeader       : extracting content type
	extractContentDisposition      : extracting the content disposition
	extractContentTransferEncoding : extracting the content transfer encoding
*/

import (
	"fmt"
	"mime"
	"strings"

	"github.com/mnako/letters/email"
)

// contentDispositionMap is a convenience map for mapping strings to
// their email package typed equivalent.
var contentDispositionMap = map[string]email.ContentDisposition{
	"attachment": email.Attachment,
	"inline":     email.Inline,
}

// contentTransferEncodingMap is a convenience map for mapping strings
// to their email package typed equivalent.
var contentTransferEncodingMap = map[string]email.ContentTransferEncoding{
	"7bit":             email.CTE7bit,
	"8bit":             email.CTE8bit,
	"binary":           email.CTEBinary,
	"quoted-printable": email.CTEQuotedPrintable,
	"base64":           email.CTEBase64,
}

func extractContentTypeHeader(s string) (email.ContentTypeHeader, error) {
	var cth email.ContentTypeHeader
	if s == "" {
		s = "text/plain"
	}
	var err error
	cth.ContentType, cth.Params, err = mime.ParseMediaType(s)
	if err != nil {
		return cth, fmt.Errorf("cannot parse Content-Type %q: %w", s, err)
	}
	for _, param := range []string{"charset", "micalg", "protocol"} {
		if v, ok := cth.Params[param]; ok {
			cth.Params[param] = strings.ToLower(v)
		}
	}
	return cth, nil
}

func extractContentDisposition(s string) (email.ContentDispositionHeader, error) {
	var cdh email.ContentDispositionHeader
	if s == "" {
		return cdh, nil
	}
	var label string
	var err error
	label, cdh.Params, err = mime.ParseMediaType(s)
	if err != nil {
		return cdh, fmt.Errorf("cannot parse Content-Disposition %q: %w", s, err)
	}
	var ok bool
	cdh.ContentDisposition, ok = contentDispositionMap[label]
	if !ok {
		return cdh, fmt.Errorf("unknown Content-Disposition %q", label)
	}
	return cdh, nil
}

func extractContentTransferEncoding(s string) (email.ContentTransferEncoding, error) {
	label := strings.ToLower(strings.TrimSpace(s))
	if label == "" {
		return email.CTE7bit, nil
	}

	cte, ok := contentTransferEncodingMap[label]
	if !ok {
		return cte, fmt.Errorf("unknown Content-Transfer-Encoding %q", label)
	}
	return cte, nil
}
