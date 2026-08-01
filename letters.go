package letters

import (
	"fmt"
	"io"
	"net/mail"
	"strings"
)

// ParseEmailHeaders parses an email header using the default header parsers.
func ParseEmailHeaders(header mail.Header) (Headers, error) {
	defaultParser := NewEmailParser()

	return defaultParser.ParseHeaders(header)
}

// ParseHeaders parses an email header using the default header parsers.
//
// Deprecated: ParseHeaders exists for backwards compatibility and will be
// removed in the future. Use NewEmailParser().ParseHeaders or ParseEmailHeaders
// instead.
func ParseHeaders(header mail.Header) (Headers, error) {
	return ParseEmailHeaders(header)
}

// ParseEmail parses an email message using a parser with the default options.
func ParseEmail(r io.Reader) (Email, error) {
	defaultParser := NewEmailParser()

	return defaultParser.Parse(r)
}

// EmailParser parses email messages according to its configured options.
type EmailParser struct {
	bodyFilter     EmailBodyFilter
	fileFilter     EmailFileFilter
	headersParsers HeadersParsers
}

// EmailParserOption configures an EmailParser.
type EmailParserOption func(*EmailParser)

// DefaultHeadersParsers returns the default set of email header parsers.
func DefaultHeadersParsers() HeadersParsers {
	return HeadersParsers{
		Date:               ParseDateHeader,
		Sender:             ParseAddressHeader,
		From:               ParseAddressListHeader,
		ReplyTo:            ParseAddressListHeader,
		To:                 ParseAddressListHeader,
		Cc:                 ParseAddressListHeader,
		Bcc:                ParseAddressListHeader,
		MessageID:          ParseMessageIdHeader,
		InReplyTo:          ParseCommaSeparatedMessageIdHeader,
		References:         ParseCommaSeparatedMessageIdHeader,
		Subject:            ParseStringHeader,
		Comments:           ParseStringHeader,
		Keywords:           ParseCommaSeparatedStringHeader,
		ResentDate:         ParseDateHeader,
		ResentFrom:         ParseAddressListHeader,
		ResentSender:       ParseAddressHeader,
		ResentTo:           ParseAddressListHeader,
		ResentCc:           ParseAddressListHeader,
		ResentBcc:          ParseAddressListHeader,
		ResentMessageID:    ParseMessageIdHeader,
		ContentType:        ParseContentTypeHeader,
		ContentDisposition: ParseContentDisposition,
		ExtraHeaders:       make(map[string]parseStringHeaderFn),
	}
}

// NewEmailParser returns an EmailParser configured with the supplied options.
func NewEmailParser(options ...EmailParserOption) *EmailParser {
	ep := &EmailParser{
		bodyFilter:     AllBodies,
		fileFilter:     AllFiles,
		headersParsers: DefaultHeadersParsers(),
	}

	for _, option := range options {
		option(ep)
	}

	return ep
}

// Parse reads and parses an email message.
func (ep *EmailParser) Parse(r io.Reader) (Email, error) {
	var email Email

	msg, err := mail.ReadMessage(r)
	if err != nil {
		return email, fmt.Errorf(
			"letters.EmailParser.Parse: cannot read message: %w",
			err,
		)
	}

	headers, err := ep.ParseHeaders(msg.Header)
	if err != nil {
		return email, fmt.Errorf(
			"letters.EmailParser.Parse: cannot parse headers: %w",
			err,
		)
	}

	email.Headers = headers

	cte, err := ParseContentTransferEncoding(
		msg.Header.Get("Content-Transfer-Encoding"),
	)
	if err != nil {
		return email, fmt.Errorf(
			"letters.EmailParser.Parse: "+
				"cannot parse Content-Transfer-Encoding: %w",
			err,
		)
	}

	contentType := email.Headers.ContentType.ContentType

	switch {
	case contentType == contentTypeTextPlain:
		if ep.bodyFilter(email.Headers.ContentType) {
			email.Text, err = parseText(
				msg.Body,
				email.Headers.ContentType.Params["charset"],
				cte,
			)
			if err != nil {
				return email, fmt.Errorf(
					"letters.EmailParser.Parse: "+
						"cannot parse plain text: %w",
					err,
				)
			}
		}
	case contentType == contentTypeTextEnriched:
		if ep.bodyFilter(email.Headers.ContentType) {
			email.EnrichedText, err = parseText(
				msg.Body,
				email.Headers.ContentType.Params["charset"],
				cte,
			)
			if err != nil {
				return email,
					fmt.Errorf(
						"letters.EmailParser.Parse: "+
							"cannot parse enriched text: %w",
						err,
					)
			}
		}
	case contentType == contentTypeTextHTML:
		if ep.bodyFilter(email.Headers.ContentType) {
			email.HTML, err = parseText(
				msg.Body,
				email.Headers.ContentType.Params["charset"],
				cte,
			)
			if err != nil {
				return email,
					fmt.Errorf(
						"letters.EmailParser.Parse: "+
							"cannot parse html text: %w",
						err,
					)
			}
		}
	case strings.HasPrefix(contentType, contentTypeMultipartPrefix):
		boundary := email.Headers.ContentType.Params["boundary"]

		emailBodies, err := ep.parsePart(
			msg.Body,
			email.Headers.ContentType,
			boundary,
		)
		if err != nil {
			return email, fmt.Errorf(
				"letters.EmailParser.Parse: "+
					"cannot parse part %q with boundary %q: %w",
				email.Headers.ContentType.ContentType,
				boundary,
				err,
			)
		}

		email.Text = emailBodies.text
		email.EnrichedText = emailBodies.enrichedText
		email.HTML = emailBodies.html
		email.InlineFiles = emailBodies.InlineFiles
		email.AttachedFiles = emailBodies.AttachedFiles
	default:
		afl, err := decodeAttachmentFileFromBody(msg.Body, email.Headers, cte)
		if err != nil {
			return email, fmt.Errorf(
				"letters.EmailParser.Parse: "+
					"cannot decode attached file content from body: %w",
				err,
			)
		}

		email.AttachedFiles = append(email.AttachedFiles, afl)
	}

	email.Text = normalizeMultilineString(email.Text)
	email.EnrichedText = normalizeMultilineString(email.EnrichedText)
	email.HTML = normalizeMultilineString(email.HTML)

	return email, nil
}
