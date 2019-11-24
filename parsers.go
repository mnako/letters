package letters

import (
	"fmt"
	"io"
	"io/ioutil"
	"mime"
	"mime/multipart"
	"net/mail"
	"strings"
	"time"

	"golang.org/x/net/html/charset"
	"golang.org/x/text/encoding"
)

func normalizeMultilineString(s string) string {
	s = strings.Replace(s, "\r\n", "\n", -1)
	s = strings.Trim(s, "\n ")
	return s
}

func normalizeParametrizedAttributeValue(s string) string {
	s = strings.Trim(s, " ")
	s = strings.ToLower(s)
	return s
}

func parseDateHeader(s string) time.Time {
	var t time.Time

	formats := []string{
		time.RFC1123Z,
		"Mon, 2 Jan 2006 15:04:05 -0700",
		time.RFC1123Z + " (MST)",
		"Mon, 2 Jan 2006 15:04:05 -0700 (MST)",
	}

	for _, format := range formats {
		t, err := time.Parse(format, s)
		if err == nil {
			return t
		}
	}

	return t
}

func parseStringHeader(s string) string {
	decodedHeader, _ := decodeHeader(s)
	return strings.Trim(decodedHeader, " ")
}

func parseCommaSeparatedStringHeader(s string) []string {
	var values []string

	normalizedS := normalizeMultilineString(s)
	if normalizedS == "" {
		return values
	}

	for _, value := range strings.Split(s, ",") {
		values = append(values, parseStringHeader(value))
	}
	return values
}

func parseAddressHeader(s string) *mail.Address {
	var address *mail.Address

	normalizedS := normalizeMultilineString(s)
	if normalizedS == "" {
		return address
	}

	decodedHeader, _ := decodeHeader(normalizedS)
	address, _ = mail.ParseAddress(decodedHeader)
	return address
}

func parseAddressListHeader(s string) []*mail.Address {
	var addresses []*mail.Address
	for _, value := range strings.Split(s, ",") {
		normalizedValue := normalizeMultilineString(value)
		if normalizedValue != "" {
			addresses = append(addresses, parseAddressHeader(value))
		}
	}
	return addresses
}

func parseMessageIdHeader(s string) MessageId {
	return MessageId(strings.Trim(s, "<> \n"))
}

func parseCommaSeparatedMessageIdHeader(s string) []MessageId {
	var values []MessageId

	for _, value := range strings.Split(s, " ") {
		messageId := parseMessageIdHeader(value)
		if messageId != "" {
			values = append(values, messageId)
		}
	}

	return values
}

func parseContentDisposition(s string) (ContentDispositionHeader, error) {
	var cdh ContentDispositionHeader

	label, params, err := mime.ParseMediaType(s)
	if label == "" {
		return cdh, nil
	}
	if err != nil {
		return cdh, fmt.Errorf(
			"letters.parsers.parseContentDisposition: cannot parse Content-Disposition %q: %w",
			s,
			err)
	}

	cd, ok := cdMap[label]
	if !ok {
		return cdh, fmt.Errorf("letters.parsers.parseContentDisposition: unknown Content-Disposition %q", label)
	}
	return ContentDispositionHeader{
		ContentDisposition: cd,
		Params:             params,
	}, nil
}

func parseContentTransferEncoding(s string) (ContentTransferEncoding, error) {
	label := normalizeParametrizedAttributeValue(s)
	if label == "" {
		return cte7bit, nil
	}

	cte, ok := cteMap[label]
	if !ok {
		return cte, fmt.Errorf("letters.parsers.parseContentTransferEncoding: unknown Content-Transfer-Encoding %q", label)
	}
	return cte, nil
}

func parseDefaultMediaType(s string) (string, map[string]string, error) {
	if s == "" {
		s = "text/plain"
	}
	mediatype, params, err := mime.ParseMediaType(s)
	if err != nil {
		return mediatype, params, fmt.Errorf(
			"letters.parsers.parseDefaultMediaType: cannot parse Content-Type %q: %w",
			s,
			err)
	}
	return mediatype, params, nil
}

func parseContentTypeHeader(s string) (ContentTypeHeader, error) {
	var cth ContentTypeHeader

	mediaType, mediaTypeParams, err := parseDefaultMediaType(s)
	if err != nil {
		return cth, fmt.Errorf(
			"letters.parsers.parseContentTypeHeader: cannot parse Content-Type %q: %w",
			s,
			err)
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

func parseHeaders(header mail.Header) (Headers, error) {
	contentType, err := parseContentTypeHeader(header.Get("Content-Type"))
	if err != nil {
		return Headers{}, fmt.Errorf(
			"letters.parsers.parseHeaders: cannot parse Content-Type: %w",
			err)
	}

	var extraHeaders = make(map[string][]string)
	for key, value := range header {
		_, isKnownHeader := knownHeaders[key]
		if !isKnownHeader {
			normalisedVals := []string{}
			for _, val := range value {
				decodedHeader, _ := decodeHeader(val)
				normalisedVals = append(normalisedVals, decodedHeader)
			}
			extraHeaders[key] = normalisedVals
		}
	}

	return Headers{
		Date:            parseDateHeader(header.Get("Date")),
		Sender:          parseAddressHeader(header.Get("Sender")),
		From:            parseAddressListHeader(header.Get("From")),
		ReplyTo:         parseAddressListHeader(header.Get("Reply-To")),
		To:              parseAddressListHeader(header.Get("To")),
		Cc:              parseAddressListHeader(header.Get("Cc")),
		Bcc:             parseAddressListHeader(header.Get("Bcc")),
		MessageID:       parseMessageIdHeader(header.Get("Message-ID")),
		InReplyTo:       parseCommaSeparatedMessageIdHeader(header.Get("In-Reply-To")),
		References:      parseCommaSeparatedMessageIdHeader(header.Get("References")),
		Subject:         parseStringHeader(header.Get("Subject")),
		Comments:        parseStringHeader(header.Get("Comments")),
		Keywords:        parseCommaSeparatedStringHeader(header.Get("Keywords")),
		ResentDate:      parseDateHeader(header.Get("Resent-Date")),
		ResentFrom:      parseAddressListHeader(header.Get("Resent-From")),
		ResentSender:    parseAddressHeader(header.Get("Resent-Sender")),
		ResentTo:        parseAddressListHeader(header.Get("Resent-To")),
		ResentCc:        parseAddressListHeader(header.Get("Resent-Cc")),
		ResentBcc:       parseAddressListHeader(header.Get("Resent-Bcc")),
		ResentMessageID: parseMessageIdHeader(header.Get("Resent-Message-ID")),
		ContentType:     contentType,
		ExtraHeaders:    extraHeaders,
	}, nil
}

func parseText(t io.Reader, e encoding.Encoding, cte ContentTransferEncoding) (string, error) {
	reader, err := decodeContent(t, e, cte)
	if err != nil {
		return "", fmt.Errorf(
			"letters.parsers.parseText: cannot decode plain text content: %w",
			err)
	}

	textBody, err := ioutil.ReadAll(reader)
	if err != nil {
		return "", fmt.Errorf(
			"letters.parsers.parseText: cannot read plain text content: %w",
			err)
	}

	return strings.TrimSuffix(string(textBody), "\n"), nil
}

func isInlineFile(contentType ContentTypeHeader, parentContentType ContentTypeHeader, p *multipart.Part) (bool, error) {
	cdh, err := parseContentDisposition(p.Header.Get("Content-Disposition"))
	if err != nil {
		return false, fmt.Errorf(
			"letters.parsers.isInlineFile: cannot parse Content-Disposition: %w",
			err)
	}
	if cdh.ContentDisposition == inline {
		return true, nil
	}
	if contentType.ContentType == contentTypeTextPlain || contentType.ContentType == contentTypeTextEnriched || contentType.ContentType == contentTypeTextHtml {
		return false, nil
	}
	return parentContentType.ContentType == contentTypeMultipartRelated, nil
}

func isAttachedFile(contentType ContentTypeHeader, parentContentType ContentTypeHeader, part *multipart.Part) (bool, error) {
	cdh, err := parseContentDisposition(part.Header.Get("Content-Disposition"))
	if err != nil {
		return false, fmt.Errorf(
			"letters.parsers.isAttachedFile: cannot parse Content-Disposition: %w",
			err)
	}
	if cdh.ContentDisposition == attachment {
		return true, nil
	}
	if contentType.ContentType != contentTypeTextPlain && contentType.ContentType != contentTypeTextEnriched && contentType.ContentType != contentTypeTextHtml {
		return true, nil
	}
	return parentContentType.ContentType == contentTypeMultipartMixed || parentContentType.ContentType == contentTypeMultipartParallel, nil
}

func parsePart(msg io.Reader, parentContentType ContentTypeHeader, boundary string) (emailBodies, error) {
	var emailBodies emailBodies

	multipartReader := multipart.NewReader(msg, boundary)
	if multipartReader == nil {
		return emailBodies, nil
	}

	for {
		part, err := multipartReader.NextPart()
		if err == io.EOF {
			break
		} else if err != nil {
			if strings.Contains(err.Error(), "EOF") {
				break
			}
			return emailBodies, fmt.Errorf(
				"letters.parsers.parsePart: cannot read part: %w",
				err)
		}

		partContentType, err := parseContentTypeHeader(part.Header.Get("Content-Type"))
		if err != nil {
			return emailBodies, fmt.Errorf(
				"letters.parsers.parsePart: cannot parse Content-Type: %w",
				err)
		}

		charsetLabel := partContentType.Params["charset"]
		if charsetLabel == "" {
			charsetLabel = parentContentType.Params["charset"]
		}

		enc, _ := charset.Lookup(charsetLabel)
		cte, err := parseContentTransferEncoding(part.Header.Get("Content-Transfer-Encoding"))
		if err != nil {
			return emailBodies, fmt.Errorf(
				"letters.parsers.parsePart: cannot parse Content-Transfer-Encoding: %w",
				err)
		}

		if partContentType.ContentType == contentTypeTextPlain {
			partTextBody, err := parseText(part, enc, cte)
			if err != nil {
				return emailBodies, fmt.Errorf(
					"letters.parsers.parsePart: cannot parse plain text: %w",
					err)
			}
			emailBodies.text += partTextBody
			emailBodies.text += "\n\n"

		} else if partContentType.ContentType == contentTypeTextEnriched {
			partEnrichedText, err := parseText(part, enc, cte)
			if err != nil {
				return emailBodies, fmt.Errorf(
					"letters.parsers.parsePart: cannot parse enriched text: %w",
					err)
			}
			emailBodies.enrichedText += partEnrichedText

		} else if partContentType.ContentType == contentTypeTextHtml {
			partHtmlBody, err := parseText(part, enc, cte)
			if err != nil {
				return emailBodies, fmt.Errorf(
					"letters.parsers.parsePart: cannot parse html text: %w",
					err)
			}
			emailBodies.html += partHtmlBody

		} else if strings.HasPrefix(partContentType.ContentType, contentTypeMultipartPrefix) {
			nestedEmailBodies, err := parsePart(part, partContentType, partContentType.Params["boundary"])
			if err != nil {
				return emailBodies, fmt.Errorf(
					"letters.parsers.parsePart: cannot parse nested part: %w",
					err)
			}

			emailBodies.extend(nestedEmailBodies)
		} else {

			isInlFile, err := isInlineFile(partContentType, parentContentType, part)
			if err != nil {
				return emailBodies, fmt.Errorf(
					"letters.parsers.parsePart: cannot read part: %w",
					err)
			}

			isAttFile, err := isAttachedFile(partContentType, parentContentType, part)
			if err != nil {
				return emailBodies, fmt.Errorf(
					"letters.parsers.parsePart: cannot check attached file: %w",
					err)
			}

			if isInlFile {
				inlineFile, err := decodeInlineFile(part, cte)
				if err != nil {
					return emailBodies, fmt.Errorf(
						"letters.parsers.parsePart: cannot decode inline file: %w",
						err)
				}
				emailBodies.InlineFiles = append(emailBodies.InlineFiles, inlineFile)
			} else if isAttFile {
				attachedFile, err := decodeAttachedFile(part, cte)
				if err != nil {
					return emailBodies, fmt.Errorf(
						"letters.parsers.parsePart: cannot decode attached file: %w",
						err)
				}
				emailBodies.AttachedFiles = append(emailBodies.AttachedFiles, attachedFile)
			} else {
				return emailBodies, &UnknownContentTypeError{contentType: parentContentType.ContentType}
			}
		}
	}

	emailBodies.text = strings.Trim(emailBodies.text, "\n")
	emailBodies.enrichedText = strings.Trim(emailBodies.enrichedText, "\n")
	emailBodies.html = strings.Trim(emailBodies.html, "\n")

	return emailBodies, nil
}
