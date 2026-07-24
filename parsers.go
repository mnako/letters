package letters

import (
	"errors"
	"fmt"
	"io"
	"mime"
	"mime/multipart"
	"net/mail"
	"strings"
	"time"

	"golang.org/x/net/html/charset"
	"golang.org/x/text/encoding"
)

func normalizeMultilineString(s string) string {
	s = strings.ReplaceAll(s, "\r\n", "\n")
	s = strings.Trim(s, "\n ")

	return s
}

func normalizeParametrizedAttributeValue(s string) string {
	s = strings.Trim(s, " ")
	s = strings.ToLower(s)

	return s
}

func getTimeLocationFromObsoleteDateFormat(dateHeader string) *time.Location {
	// From RFC5322 Section 4.3.  Obsolete Date and Time:
	// The remaining three character zones are the US time zones. The first
	// letter, "E", "C", "M", or "P" stands for "Eastern", "Central",
	// "Mountain", and "Pacific".  The second letter is either "S" for
	// "Standard" time, or "D" for "Daylight Savings" (or summer) time. Their
	// interpretations are as follows:
	//
	//   EDT is semantically equivalent to -0400
	//   EST is semantically equivalent to -0500
	//   CDT is semantically equivalent to -0500
	//   CST is semantically equivalent to -0600
	//   MDT is semantically equivalent to -0600
	//   MST is semantically equivalent to -0700
	//   PDT is semantically equivalent to -0700
	//   PST is semantically equivalent to -0800
	const (
		edtOffset = -4 * 60 * 60
		estOffset = -5 * 60 * 60
		cdtOffset = -5 * 60 * 60
		cstOffset = -6 * 60 * 60
		mdtOffset = -6 * 60 * 60
		mstOffset = -7 * 60 * 60
		pdtOffset = -7 * 60 * 60
		pstOffset = -8 * 60 * 60
	)

	equivalentTZs := map[string]*time.Location{
		"EDT": time.FixedZone("EDT", edtOffset),
		"EST": time.FixedZone("EST", estOffset),
		"CDT": time.FixedZone("CDT", cdtOffset),
		"CST": time.FixedZone("CST", cstOffset),
		"MDT": time.FixedZone("MDT", mdtOffset),
		"MST": time.FixedZone("MST", mstOffset),
		"PDT": time.FixedZone("PDT", pdtOffset),
		"PST": time.FixedZone("PST", pstOffset),
	}

	for name, location := range equivalentTZs {
		if strings.HasSuffix(dateHeader, name) {
			return location
		}
	}

	return nil
}

// ParseDateHeader parses an RFC-compatible email date or returns the zero time.
func ParseDateHeader(dateHeader string) time.Time {
	// We follow date formats specified in RFC5322, RFC2822, and RFC822,
	// including obsolete but legal formats, but parse them according to
	// time.Parse Go implementation. The main difference is in parsing
	// two-digit years. While RFC5322 Section 4.3. prescribes that:
	//
	//   If a two digit year is encountered whose value is between 00 and 49,
	//   the year is interpreted by adding 2000, ending up with a value
	//   between 2000 and 2049.  If a two digit year is encountered with
	//   a value between 50 and 99, or any three digit year is encountered,
	//   the year is interpreted by adding 1900.
	//
	// time.Parse states that:
	//
	//   For layouts specifying the two-digit year 06, a value NN >= 69 will
	//   be treated as 19NN and a value NN < 69 will be treated as 20NN.
	//
	// We have chosen to unify the behaviour of letters with the Go implementation.
	//
	// Cf. parsers_test.TestParseDateHeader for test cases taken directly
	// from the specifications and appendices.
	var parsedDate time.Time

	obsLocation := getTimeLocationFromObsoleteDateFormat(dateHeader)

	formats := []string{
		"02 Jan 2006 1504 -0700",
		"02 Jan 2006 1504 MST",
		"02 Jan 2006 1504 (MST)",
		"02 Jan 2006 1504 -0700 (MST)",
		"02 Jan 2006 15:04 -0700",
		"02 Jan 2006 15:04 MST",
		"02 Jan 2006 15:04 (MST)",
		"02 Jan 2006 15:04 -0700 (MST)",
		"02 Jan 2006 15:04:05 -0700",
		"02 Jan 2006 15:04:05 MST",
		"02 Jan 2006 15:04:05 (MST)",
		"02 Jan 2006 15:04:05 -0700 (MST)",
		"02 Jan 06 1504 -0700",
		"02 Jan 06 1504 MST",
		"02 Jan 06 1504 (MST)",
		"02 Jan 06 1504 -0700 (MST)",
		"02 Jan 06 15:04 -0700",
		"02 Jan 06 15:04 MST",
		"02 Jan 06 15:04 (MST)",
		"02 Jan 06 15:04 -0700 (MST)",
		"02 Jan 06 15:04:05 -0700",
		"02 Jan 06 15:04:05 MST",
		"02 Jan 06 15:04:05 (MST)",
		"02 Jan 06 15:04:05 -0700 (MST)",
		"2 Jan 2006 1504 -0700",
		"2 Jan 2006 1504 MST",
		"2 Jan 2006 1504 (MST)",
		"2 Jan 2006 1504 -0700 (MST)",
		"2 Jan 2006 15:04 -0700",
		"2 Jan 2006 15:04 MST",
		"2 Jan 2006 15:04 (MST)",
		"2 Jan 2006 15:04 -0700 (MST)",
		"2 Jan 2006 15:04:05 -0700",
		"2 Jan 2006 15:04:05 MST",
		"2 Jan 2006 15:04:05 (MST)",
		"2 Jan 2006 15:04:05 -0700 (MST)",
		"2 Jan 06 1504 -0700",
		"2 Jan 06 1504 MST",
		"2 Jan 06 1504 (MST)",
		"2 Jan 06 1504 -0700 (MST)",
		"2 Jan 06 15:04 -0700",
		"2 Jan 06 15:04 MST",
		"2 Jan 06 15:04 (MST)",
		"2 Jan 06 15:04 -0700 (MST)",
		"2 Jan 06 15:04:05 -0700",
		"2 Jan 06 15:04:05 MST",
		"2 Jan 06 15:04:05 (MST)",
		"2 Jan 06 15:04:05 -0700 (MST)",
		"Mon, 02 Jan 2006 1504 -0700",
		"Mon, 02 Jan 2006 1504 MST",
		"Mon, 02 Jan 2006 1504 (MST)",
		"Mon, 02 Jan 2006 1504 -0700 (MST)",
		"Mon, 02 Jan 2006 15:04 -0700",
		"Mon, 02 Jan 2006 15:04 MST",
		"Mon, 02 Jan 2006 15:04 (MST)",
		"Mon, 02 Jan 2006 15:04 -0700 (MST)",
		"Mon, 02 Jan 2006 15:04:05 -0700",
		"Mon, 02 Jan 2006 15:04:05 MST",
		"Mon, 02 Jan 2006 15:04:05 (MST)",
		"Mon, 02 Jan 2006 15:04:05 -0700 (MST)",
		"Mon, 02 Jan 06 1504 -0700",
		"Mon, 02 Jan 06 1504 MST",
		"Mon, 02 Jan 06 1504 (MST)",
		"Mon, 02 Jan 06 1504 -0700 (MST)",
		"Mon, 02 Jan 06 15:04 -0700",
		"Mon, 02 Jan 06 15:04 MST",
		"Mon, 02 Jan 06 15:04 (MST)",
		"Mon, 02 Jan 06 15:04 -0700 (MST)",
		"Mon, 02 Jan 06 15:04:05 -0700",
		"Mon, 02 Jan 06 15:04:05 MST",
		"Mon, 02 Jan 06 15:04:05 (MST)",
		"Mon, 02 Jan 06 15:04:05 -0700 (MST)",
		"Mon, 2 Jan 2006 1504 -0700",
		"Mon, 2 Jan 2006 1504 MST",
		"Mon, 2 Jan 2006 1504 (MST)",
		"Mon, 2 Jan 2006 1504 -0700 (MST)",
		"Mon, 2 Jan 2006 15:04 -0700",
		"Mon, 2 Jan 2006 15:04 MST",
		"Mon, 2 Jan 2006 15:04 (MST)",
		"Mon, 2 Jan 2006 15:04 -0700 (MST)",
		"Mon, 2 Jan 2006 15:04:05 -0700",
		"Mon, 2 Jan 2006 15:04:05 MST",
		"Mon, 2 Jan 2006 15:04:05 (MST)",
		"Mon, 2 Jan 2006 15:04:05 -0700 (MST)",
		"Mon, 2 Jan 06 1504 -0700",
		"Mon, 2 Jan 06 1504 MST",
		"Mon, 2 Jan 06 1504 (MST)",
		"Mon, 2 Jan 06 1504 -0700 (MST)",
		"Mon, 2 Jan 06 15:04 -0700",
		"Mon, 2 Jan 06 15:04 MST",
		"Mon, 2 Jan 06 15:04 (MST)",
		"Mon, 2 Jan 06 15:04 -0700 (MST)",
		"Mon, 2 Jan 06 15:04:05 -0700",
		"Mon, 2 Jan 06 15:04:05 MST",
		"Mon, 2 Jan 06 15:04:05 (MST)",
		"Mon, 2 Jan 06 15:04:05 -0700 (MST)",
	}

	for _, format := range formats {
		if obsLocation == nil {
			t, err := time.Parse(format, dateHeader)
			if err == nil {
				return t
			}
		} else {
			t, err := time.ParseInLocation(format, dateHeader, obsLocation)
			if err == nil {
				return t
			}
		}
	}

	// Best-effort support for RFC5322 Appendix A.5. with obs-zone
	sWithoutObsZone := strings.Split(dateHeader, " (")[0]

	formatsWithObsZone := []string{
		"02 Jan 2006 1504 -0700",
		"02 Jan 2006 1504 MST",
		"02 Jan 2006 15:04 -0700",
		"02 Jan 2006 15:04 MST",
		"02 Jan 2006 15:04:05 -0700",
		"02 Jan 2006 15:04:05 MST",
		"02 Jan 06 1504 -0700",
		"02 Jan 06 1504 MST",
		"02 Jan 06 15:04 -0700",
		"02 Jan 06 15:04 MST",
		"02 Jan 06 15:04:05 -0700",
		"02 Jan 06 15:04:05 MST",
		"2 Jan 2006 1504 -0700",
		"2 Jan 2006 1504 MST",
		"2 Jan 2006 15:04 -0700",
		"2 Jan 2006 15:04 MST",
		"2 Jan 2006 15:04:05 -0700",
		"2 Jan 2006 15:04:05 MST",
		"2 Jan 06 1504 -0700",
		"2 Jan 06 1504 MST",
		"2 Jan 06 15:04 -0700",
		"2 Jan 06 15:04 MST",
		"2 Jan 06 15:04:05 -0700",
		"2 Jan 06 15:04:05 MST",
		"Mon, 02 Jan 2006 1504 -0700",
		"Mon, 02 Jan 2006 1504 MST",
		"Mon, 02 Jan 2006 15:04 -0700",
		"Mon, 02 Jan 2006 15:04 MST",
		"Mon, 02 Jan 2006 15:04:05 -0700",
		"Mon, 02 Jan 2006 15:04:05 MST",
		"Mon, 02 Jan 06 1504 -0700",
		"Mon, 02 Jan 06 1504 MST",
		"Mon, 02 Jan 06 15:04 -0700",
		"Mon, 02 Jan 06 15:04 MST",
		"Mon, 02 Jan 06 15:04:05 -0700",
		"Mon, 02 Jan 06 15:04:05 MST",
		"Mon, 2 Jan 2006 1504 -0700",
		"Mon, 2 Jan 2006 1504 MST",
		"Mon, 2 Jan 2006 15:04 -0700",
		"Mon, 2 Jan 2006 15:04 MST",
		"Mon, 2 Jan 2006 15:04:05 -0700",
		"Mon, 2 Jan 2006 15:04:05 MST",
		"Mon, 2 Jan 06 1504 -0700",
		"Mon, 2 Jan 06 1504 MST",
		"Mon, 2 Jan 06 15:04 -0700",
		"Mon, 2 Jan 06 15:04 MST",
		"Mon, 2 Jan 06 15:04:05 -0700",
		"Mon, 2 Jan 06 15:04:05 MST",
	}

	for _, formatWithObsZone := range formatsWithObsZone {
		if obsLocation == nil {
			parsedTime, err := time.Parse(formatWithObsZone, sWithoutObsZone)
			if err == nil {
				return parsedTime
			}
		} else {
			parsedTime, err := time.ParseInLocation(
				formatWithObsZone,
				sWithoutObsZone,
				obsLocation,
			)
			if err == nil {
				return parsedTime
			}
		}
	}

	return parsedDate
}

// ParseStringHeader decodes and trims a string-valued email header.
func ParseStringHeader(s string) string {
	decodedHeader, _ := decodeHeader(s)

	return strings.Trim(decodedHeader, " ")
}

// ParseCommaSeparatedStringHeader decodes a comma-separated email header.
func ParseCommaSeparatedStringHeader(headerValue string) []string {
	var values []string

	normalizedS := normalizeMultilineString(headerValue)
	if normalizedS == "" {
		return values
	}

	for _, value := range strings.Split(headerValue, ",") {
		values = append(values, ParseStringHeader(value))
	}

	return values
}

// ParseAddressHeader parses a named single-address email header.
func ParseAddressHeader(
	header mail.Header,
	name string,
) (*mail.Address, error) {
	var address *mail.Address

	ss, ok := header[name]
	if !ok {
		return address, nil
	}

	addressHeader := strings.Join(ss, ", ")

	normalizedS := normalizeMultilineString(addressHeader)
	if normalizedS == "" {
		return address, nil
	}

	decodedHeader, err := decodeHeader(normalizedS)
	if err != nil {
		return address, fmt.Errorf(
			"letters.parsers.parseAddressHeader: "+
				"cannot decode address header %q: %w",
			addressHeader,
			err,
		)
	}

	address, err = mail.ParseAddress(decodedHeader)
	if err != nil {
		return address, fmt.Errorf(
			"letters.parsers.parseAddressHeader: "+
				"cannot parse address header %q: %w",
			addressHeader,
			err,
		)
	}

	return address, nil
}

// ParseAddressListHeader parses a named address-list email header.
func ParseAddressListHeader(
	header mail.Header,
	name string,
) ([]*mail.Address, error) {
	var addresses []*mail.Address

	ss, ok := header[name]
	if !ok {
		return addresses, nil
	}

	addressListHeader := strings.Join(ss, ", ")

	normalizedS := normalizeMultilineString(addressListHeader)
	if normalizedS == "" {
		return addresses, nil
	}

	decodedHeader, err := decodeHeader(normalizedS)
	if err != nil {
		return addresses, fmt.Errorf(
			"letters.parsers.parseAddressListHeader: "+
				"cannot decode address list header %q: %w",
			addressListHeader,
			err,
		)
	}

	addresses, err = mail.ParseAddressList(decodedHeader)
	if err != nil {
		return addresses, fmt.Errorf(
			"letters.parsers.parseAddressListHeader: "+
				"cannot parse address list header %q: %w",
			addressListHeader,
			err,
		)
	}

	return addresses, nil
}

// ParseMessageIdHeader trims delimiters and whitespace from a message ID.
//
//nolint:revive // The exported name is kept for backward compatibility.
func ParseMessageIdHeader(s string) MessageId {
	return MessageId(strings.Trim(s, "<> \n"))
}

// ParseCommaSeparatedMessageIdHeader parses a space-separated list of message IDs.
//
//nolint:revive // The exported name is kept for backward compatibility.
func ParseCommaSeparatedMessageIdHeader(s string) []MessageId {
	var values []MessageId

	for _, value := range strings.Split(s, " ") {
		messageID := ParseMessageIdHeader(value)
		if messageID != "" {
			values = append(values, messageID)
		}
	}

	return values
}

// ParseContentDisposition parses a supported Content-Disposition header.
func ParseContentDisposition(
	contentDispositionValue string,
) (ContentDispositionHeader, error) {
	var cdh ContentDispositionHeader

	label, params, err := mime.ParseMediaType(contentDispositionValue)
	if label == "" {
		return cdh, nil
	}

	if err != nil {
		return cdh, fmt.Errorf(
			"letters.parsers.parseContentDisposition: "+
				"cannot parse Content-Disposition %q: %w",
			contentDispositionValue,
			err,
		)
	}

	cd := ContentDisposition(label)

	switch cd {
	case ContentDispositionAttachment, ContentDispositionInline:
	default:
		return cdh, fmt.Errorf("%w %q", ErrUnknownContentDisposition, label)
	}

	return ContentDispositionHeader{
		ContentDisposition: cd,
		Params:             params,
	}, nil
}

// ParseContentTransferEncoding parses a Content-Transfer-Encoding header.
func ParseContentTransferEncoding(s string) (ContentTransferEncoding, error) {
	label := normalizeParametrizedAttributeValue(s)
	if label == "" {
		return cte7bit, nil
	}

	cte := ContentTransferEncoding(label)

	switch cte {
	case cte7bit, cte8bit, cteBinary, cteQuotedPrintable, cteBase64:
	default:
		return cte, fmt.Errorf(
			"%w %q",
			ErrUnknownContentTransferEncoding,
			label,
		)
	}

	return cte, nil
}

// ParseDefaultMediaType parses a media type, defaulting an empty value to text/plain.
func ParseDefaultMediaType(
	contentTypeValue string,
) (string, map[string]string, error) {
	if contentTypeValue == "" {
		contentTypeValue = "text/plain"
	}

	mediatype, params, err := mime.ParseMediaType(contentTypeValue)
	if err != nil {
		return mediatype, params, fmt.Errorf(
			"letters.parsers.parseDefaultMediaType: "+
				"cannot parse Content-Type %q: %w",
			contentTypeValue,
			err,
		)
	}

	return mediatype, params, nil
}

// ParseContentTypeHeader parses and normalizes a Content-Type header.
func ParseContentTypeHeader(
	contentTypeValue string,
) (ContentTypeHeader, error) {
	var cth ContentTypeHeader

	mediaType, mediaTypeParams, err := ParseDefaultMediaType(contentTypeValue)
	if err != nil {
		return cth, fmt.Errorf(
			"letters.parsers.parseContentTypeHeader: "+
				"cannot parse Content-Type %q: %w",
			contentTypeValue,
			err,
		)
	}

	for _, param := range []string{"charset", "micalg", "protocol"} {
		if mediaTypeParams[param] != "" {
			mediaTypeParams[param] = normalizeParametrizedAttributeValue(
				mediaTypeParams[param],
			)
		}
	}

	return ContentTypeHeader{
		ContentType: mediaType,
		Params:      mediaTypeParams,
	}, nil
}

func isKnownHeader(name string) bool {
	switch name {
	case "Date",
		"Sender",
		"From",
		"Reply-To",
		"To",
		"Cc",
		"Bcc",
		"Message-Id",
		"In-Reply-To",
		"References",
		"Subject",
		"Comments",
		"Keywords",
		"Resent-Date",
		"Resent-From",
		"Resent-Sender",
		"Resent-To",
		"Resent-Cc",
		"Resent-Bcc",
		"Resent-Message-Id",
		"Content-Transfer-Encoding",
		"Content-Type",
		"Content-Disposition":
		return true
	default:
		return false
	}
}

// ParseHeaders parses an email header using the parser's configured header parsers.
func (ep *EmailParser) ParseHeaders(header mail.Header) (Headers, error) {
	contentType, err := ep.headersParsers.ContentType(
		header.Get("Content-Type"),
	)
	if err != nil {
		return Headers{}, fmt.Errorf(
			"letters.parsers.ParseHeaders: "+
				"cannot parse Content-Type: %w",
			err)
	}

	contentDisposition, _ := ep.headersParsers.ContentDisposition(
		header.Get("Content-Disposition"),
	)

	extraHeaders := make(map[string][]string)

	for key, value := range header {
		if isKnownHeader(key) {
			continue
		}

		normalisedVals := []string{}
		extraHeaderParserFn,
			hasExtraHeaderParseFn := ep.headersParsers.ExtraHeaders[strings.ToLower(key)]

		for _, val := range value {
			if !hasExtraHeaderParseFn {
				decodedHeader, _ := decodeHeader(val)
				normalisedVals = append(normalisedVals, decodedHeader)
			} else {
				decodedHeader := extraHeaderParserFn(val)
				normalisedVals = append(normalisedVals, decodedHeader)
			}
		}

		extraHeaders[key] = normalisedVals
	}

	sender, err := ep.headersParsers.Sender(header, "Sender")
	if err != nil {
		return Headers{}, fmt.Errorf(
			"letters.parsers.ParseHeaders: "+
				"cannot parse Sender header: %w",
			err,
		)
	}

	from, err := ep.headersParsers.From(header, "From")
	if err != nil {
		return Headers{}, fmt.Errorf(
			"letters.parsers.ParseHeaders: "+
				"cannot parse From header: %w",
			err,
		)
	}

	replyTo, err := ep.headersParsers.ReplyTo(header, "Reply-To")
	if err != nil {
		return Headers{}, fmt.Errorf(
			"letters.parsers.ParseHeaders: "+
				"cannot parse Reply-To header: %w",
			err,
		)
	}

	to, err := ep.headersParsers.To(header, "To")
	if err != nil {
		return Headers{}, fmt.Errorf(
			"letters.parsers.ParseHeaders: "+
				"cannot parse To header: %w",
			err,
		)
	}

	cc, err := ep.headersParsers.Cc(header, "Cc")
	if err != nil {
		return Headers{}, fmt.Errorf(
			"letters.parsers.ParseHeaders: "+
				"cannot parse Cc header: %w",
			err,
		)
	}

	bcc, err := ep.headersParsers.Bcc(header, "Bcc")
	if err != nil {
		return Headers{}, fmt.Errorf(
			"letters.parsers.ParseHeaders: "+
				"cannot parse Bcc header: %w",
			err,
		)
	}

	resentFrom, err := ep.headersParsers.ResentFrom(header, "Resent-From")
	if err != nil {
		return Headers{}, fmt.Errorf(
			"letters.parsers.ParseHeaders: "+
				"cannot parse Resent-From header: %w",
			err,
		)
	}

	resentSender, err := ep.headersParsers.ResentSender(header, "Resent-Sender")
	if err != nil {
		return Headers{}, fmt.Errorf(
			"letters.parsers.ParseHeaders: "+
				"cannot parse Resent-Sender header: %w",
			err,
		)
	}

	resentTo, err := ep.headersParsers.ResentTo(header, "Resent-To")
	if err != nil {
		return Headers{}, fmt.Errorf(
			"letters.parsers.ParseHeaders: "+
				"cannot parse Resent-To header: %w",
			err,
		)
	}

	resentCc, err := ep.headersParsers.ResentCc(header, "Resent-Cc")
	if err != nil {
		return Headers{}, fmt.Errorf(
			"letters.parsers.ParseHeaders: "+
				"cannot parse Resent-Cc header: %w",
			err,
		)
	}

	resentBcc, err := ep.headersParsers.ResentBcc(header, "Resent-Bcc")
	if err != nil {
		return Headers{}, fmt.Errorf(
			"letters.parsers.ParseHeaders: "+
				"cannot parse Resent-Bcc header: %w",
			err,
		)
	}

	return Headers{
		Date:    ep.headersParsers.Date(header.Get("Date")),
		Sender:  sender,
		From:    from,
		ReplyTo: replyTo,
		To:      to,
		Cc:      cc,
		Bcc:     bcc,
		MessageID: ep.headersParsers.MessageID(
			header.Get("Message-ID"),
		),
		InReplyTo: ep.headersParsers.InReplyTo(
			header.Get("In-Reply-To"),
		),
		References: ep.headersParsers.References(
			header.Get("References"),
		),
		Subject:  ep.headersParsers.Subject(header.Get("Subject")),
		Comments: ep.headersParsers.Comments(header.Get("Comments")),
		Keywords: ep.headersParsers.Keywords(header.Get("Keywords")),
		ResentDate: ep.headersParsers.ResentDate(
			header.Get("Resent-Date"),
		),
		ResentFrom:   resentFrom,
		ResentSender: resentSender,
		ResentTo:     resentTo,
		ResentCc:     resentCc,
		ResentBcc:    resentBcc,
		ResentMessageID: ep.headersParsers.ResentMessageID(
			header.Get("Resent-Message-ID"),
		),
		ContentType:        contentType,
		ContentDisposition: contentDisposition,
		ExtraHeaders:       extraHeaders,
	}, nil
}

func parseText(
	t io.Reader,
	e encoding.Encoding,
	cte ContentTransferEncoding,
) (string, error) {
	reader, err := decodeContent(t, e, cte)
	if err != nil {
		return "", fmt.Errorf(
			"letters.parsers.parseText: "+
				"cannot decode plain text content: %w",
			err,
		)
	}

	textBody, err := io.ReadAll(reader)
	if err != nil {
		return "", fmt.Errorf(
			"letters.parsers.parseText: "+
				"cannot read plain text content: %w",
			err,
		)
	}

	return strings.TrimSuffix(string(textBody), "\n"), nil
}

func isInlineFile(
	contentType ContentTypeHeader,
	parentContentType ContentTypeHeader,
	cdh ContentDispositionHeader,
) bool {
	if cdh.ContentDisposition == ContentDispositionInline {
		return true
	}

	if contentType.ContentType == contentTypeTextPlain ||
		contentType.ContentType == contentTypeTextEnriched ||
		contentType.ContentType == contentTypeTextHTML {
		return false
	}

	return parentContentType.ContentType == contentTypeMultipartRelated
}

func isAttachedFile(
	contentType ContentTypeHeader,
	parentContentType ContentTypeHeader,
) bool {
	if contentType.ContentType != contentTypeTextPlain &&
		contentType.ContentType != contentTypeTextEnriched &&
		contentType.ContentType != contentTypeTextHTML {
		return true
	}

	return parentContentType.ContentType == contentTypeMultipartMixed ||
		parentContentType.ContentType == contentTypeMultipartParallel
}

func (ep *EmailParser) parsePart(
	msg io.Reader,
	parentContentType ContentTypeHeader,
	boundary string,
) (emailBodies, error) {
	var emailBodies emailBodies

	multipartReader := multipart.NewReader(msg, boundary)
	if multipartReader == nil {
		return emailBodies, nil
	}

	for {
		part, err := multipartReader.NextPart()
		if errors.Is(err, io.EOF) {
			break
		} else if err != nil {
			if strings.Contains(err.Error(), "EOF") {
				break
			}

			return emailBodies, fmt.Errorf(
				"letters.parsers.parsePart: cannot read part: %w",
				err,
			)
		}

		partContentType, err := ParseContentTypeHeader(
			part.Header.Get("Content-Type"),
		)
		if err != nil {
			return emailBodies, fmt.Errorf(
				"letters.parsers.parsePart: "+
					"cannot parse Content-Type: %w",
				err,
			)
		}

		charsetLabel := partContentType.Params["charset"]
		if charsetLabel == "" {
			charsetLabel = parentContentType.Params["charset"]
		}

		enc, _ := charset.Lookup(charsetLabel)

		cte, err := ParseContentTransferEncoding(
			part.Header.Get("Content-Transfer-Encoding"),
		)
		if err != nil {
			return emailBodies, fmt.Errorf(
				"letters.parsers.parsePart: "+
					"cannot parse Content-Transfer-Encoding: %w",
				err,
			)
		}

		cdh, err := ParseContentDisposition(
			part.Header.Get("Content-Disposition"),
		)
		if err != nil {
			return emailBodies, fmt.Errorf(
				"letters.parsers.parsePart: "+
					"cannot parse Content-Disposition: %w",
				err,
			)
		}

		if cdh.ContentDisposition == ContentDispositionAttachment {
			if !ep.fileFilter(partContentType, cdh) {
				continue
			}

			attachedFile, err := decodeAttachedFileFromPart(part, cte)
			if err != nil {
				return emailBodies, fmt.Errorf(
					"letters.parsers.parsePart: "+
						"cannot decode attached file: %w",
					err,
				)
			}

			emailBodies.AttachedFiles = append(
				emailBodies.AttachedFiles,
				attachedFile,
			)

			continue
		}

		if partContentType.ContentType == contentTypeTextPlain {
			if !ep.bodyFilter(partContentType) {
				continue
			}

			partTextBody, err := parseText(part, enc, cte)
			if err != nil {
				return emailBodies, fmt.Errorf(
					"letters.parsers.parsePart: "+
						"cannot parse plain text: %w",
					err,
				)
			}

			emailBodies.text += partTextBody
			emailBodies.text += "\n\n"

			continue
		}

		if partContentType.ContentType == contentTypeTextEnriched {
			if !ep.bodyFilter(partContentType) {
				continue
			}

			partEnrichedText, err := parseText(part, enc, cte)
			if err != nil {
				return emailBodies, fmt.Errorf(
					"letters.parsers.parsePart: "+
						"cannot parse enriched text: %w",
					err,
				)
			}

			emailBodies.enrichedText += partEnrichedText

			continue
		}

		if partContentType.ContentType == contentTypeTextHTML {
			if !ep.bodyFilter(partContentType) {
				continue
			}

			partHTMLBody, err := parseText(part, enc, cte)
			if err != nil {
				return emailBodies, fmt.Errorf(
					"letters.parsers.parsePart: "+
						"cannot parse html text: %w",
					err,
				)
			}

			emailBodies.html += partHTMLBody

			continue
		}

		if strings.HasPrefix(
			partContentType.ContentType,
			contentTypeMultipartPrefix,
		) {
			nestedEmailBodies, err := ep.parsePart(
				part,
				partContentType,
				partContentType.Params["boundary"],
			)
			if err != nil {
				return emailBodies, fmt.Errorf(
					"letters.parsers.parsePart: "+
						"cannot parse nested part: %w",
					err,
				)
			}

			emailBodies.extend(nestedEmailBodies)

			continue
		}

		if isInlineFile(partContentType, parentContentType, cdh) {
			if !ep.fileFilter(partContentType, cdh) {
				continue
			}

			inlineFile, err := decodeInlineFile(part, cte)
			if err != nil {
				return emailBodies, fmt.Errorf(
					"letters.parsers.parsePart: "+
						"cannot decode inline file: %w",
					err,
				)
			}

			emailBodies.InlineFiles = append(
				emailBodies.InlineFiles,
				inlineFile,
			)

			continue
		}

		if isAttachedFile(partContentType, parentContentType) {
			if !ep.fileFilter(partContentType, cdh) {
				continue
			}

			attachedFile, err := decodeAttachedFileFromPart(part, cte)
			if err != nil {
				return emailBodies, fmt.Errorf(
					"letters.parsers.parsePart: "+
						"cannot decode attached file: %w",
					err,
				)
			}

			emailBodies.AttachedFiles = append(
				emailBodies.AttachedFiles,
				attachedFile,
			)

			continue
		}

		return emailBodies, &UnknownContentTypeError{
			contentType: parentContentType.ContentType,
		}
	}

	emailBodies.text = strings.Trim(emailBodies.text, "\n")
	emailBodies.enrichedText = strings.Trim(emailBodies.enrichedText, "\n")
	emailBodies.html = strings.Trim(emailBodies.html, "\n")

	return emailBodies, nil
}
