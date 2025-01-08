package letters

import (
	"fmt"
	"io"
	"net/mail"
	"strings"

	"golang.org/x/net/html/charset"
)

// processType set the type of processing that will be done by
// parseEmail
type processType int

const (
	entireEmail        processType = iota // process the entire email (default)
	headersOnly                           // only process headers
	withoutAttachments                    // do not include attachments
)

// processSetting sets the processType for this run, by default the
// entire email
var processSetting processType = entireEmail

// ParseEmail parses all parts of an email.
func ParseEmail(r io.Reader) (Email, error) {
	return parseEmail(r)
}

// ParseEmailHeaders parses only the headers of an email.
func ParseEmailHeaders(r io.Reader) (Email, error) {
	processSetting = headersOnly
	return parseEmail(r)
}

// ParseEmailWithoutAttachments parses only the headers and inline
// attachments of an email.
func ParseEmailWithoutAttachments(r io.Reader) (Email, error) {
	processSetting = withoutAttachments
	return parseEmail(r)
}

// parseEmail is the main email parsing function. Depending on the
// processSetting, it may return early by processing only part of the
// email.
func parseEmail(r io.Reader) (Email, error) {
	var email Email

	msg, err := mail.ReadMessage(r)
	if err != nil {
		return email, fmt.Errorf("letters.ParseEmail: cannot read message: %w", err)
	}

	headers, err := parseHeaders(msg.Header)
	if err != nil {
		return email, fmt.Errorf("letters.ParseEmail: cannot parse headers: %w", err)
	}

	email = Email{
		Headers: headers,
	}
	if processSetting == headersOnly {
		return email, nil
	}

	encoding, _ := charset.Lookup(email.Headers.ContentType.Params["charset"])
	cte, err := parseContentTransferEncoding(msg.Header.Get("Content-Transfer-Encoding"))
	if err != nil {
		return email, fmt.Errorf("letters.ParseEmail: cannot parse Content-Transfer-Encoding: %w", err)
	}

	if email.Headers.ContentType.ContentType == contentTypeTextPlain {
		email.Text, err = parseText(msg.Body, encoding, cte)
		if err != nil {
			return email, fmt.Errorf("letters.ParseEmail: cannot parse plain text: %w", err)
		}

	} else if email.Headers.ContentType.ContentType == contentTypeTextEnriched {
		email.EnrichedText, err = parseText(msg.Body, encoding, cte)
		if err != nil {
			return email, fmt.Errorf("letters.ParseEmail: cannot parse enriched text: %w", err)
		}

	} else if email.Headers.ContentType.ContentType == contentTypeTextHtml {
		email.HTML, err = parseText(msg.Body, encoding, cte)
		if err != nil {
			return email, fmt.Errorf("letters.ParseEmail: cannot parse html text: %w", err)
		}

	} else if strings.HasPrefix(email.Headers.ContentType.ContentType, contentTypeMultipartPrefix) {
		boundary := email.Headers.ContentType.Params["boundary"]
		emailBodies, err := parsePart(msg.Body, email.Headers.ContentType, boundary)
		if err != nil {
			return email, fmt.Errorf(
				"letters.ParseEmail: cannot parse part %q with boundary %q: %w",
				email.Headers.ContentType.ContentType,
				boundary,
				err)
		}
		email.Text = emailBodies.text
		email.EnrichedText = emailBodies.enrichedText
		email.HTML = emailBodies.html
		email.InlineFiles = emailBodies.InlineFiles
		email.AttachedFiles = emailBodies.AttachedFiles

	} else {
		afl, err := decodeAttachmentFileFromBody(msg.Body, email.Headers, cte)
		if err != nil {
			return email, fmt.Errorf(
				"letters.decoders.ParseEmail: cannot decode attached file content from body: %w",
				err)
		}
		email.AttachedFiles = append(email.AttachedFiles, afl)
	}

	email.Text = normalizeMultilineString(email.Text)
	email.EnrichedText = normalizeMultilineString(email.EnrichedText)
	email.HTML = normalizeMultilineString(email.HTML)

	return email, nil
}
