package letters

import (
	"net/mail"
	"strings"
	"time"
)

type (
	parseDateHeaderFn                    func(string) time.Time
	parseStringHeaderFn                  func(string) string
	parseCommaSeparatedStringHeaderFn    func(string) []string
	parseAddressHeaderFn                 func(mail.Header, string) (*mail.Address, error)
	parseAddressListHeaderFn             func(mail.Header, string) ([]*mail.Address, error)
	parseMessageIDHeaderFn               func(string) MessageId
	parseCommaSeparatedMessageIDHeaderFn func(string) []MessageId
	parseContentDispositionHeaderFn      func(string) (ContentDispositionHeader, error)
	parseContentTypeHeaderFn             func(string) (ContentTypeHeader, error)
)

// HeadersParsers contains the parser functions used for individual headers.
type HeadersParsers struct {
	Date               parseDateHeaderFn
	Sender             parseAddressHeaderFn
	From               parseAddressListHeaderFn
	ReplyTo            parseAddressListHeaderFn
	To                 parseAddressListHeaderFn
	Cc                 parseAddressListHeaderFn
	Bcc                parseAddressListHeaderFn
	MessageID          parseMessageIDHeaderFn
	InReplyTo          parseCommaSeparatedMessageIDHeaderFn
	References         parseCommaSeparatedMessageIDHeaderFn
	Subject            parseStringHeaderFn
	Comments           parseStringHeaderFn
	Keywords           parseCommaSeparatedStringHeaderFn
	ResentDate         parseDateHeaderFn
	ResentFrom         parseAddressListHeaderFn
	ResentSender       parseAddressHeaderFn
	ResentTo           parseAddressListHeaderFn
	ResentCc           parseAddressListHeaderFn
	ResentBcc          parseAddressListHeaderFn
	ResentMessageID    parseMessageIDHeaderFn
	ContentType        parseContentTypeHeaderFn
	ContentDisposition parseContentDispositionHeaderFn
	ExtraHeaders       map[string]parseStringHeaderFn
}

// WithDateHeaderParser configures the parser used for the Date header.
func WithDateHeaderParser(
	dateHeaderParserFn parseDateHeaderFn,
) EmailParserOption {
	return func(ep *EmailParser) {
		ep.headersParsers.Date = dateHeaderParserFn
	}
}

// WithSenderHeaderParser configures the parser used for the Sender header.
func WithSenderHeaderParser(
	senderHeaderParserFn parseAddressHeaderFn,
) EmailParserOption {
	return func(ep *EmailParser) {
		ep.headersParsers.Sender = senderHeaderParserFn
	}
}

// WithFromHeaderParser configures the parser used for the From header.
func WithFromHeaderParser(
	fromHeaderParserFn parseAddressListHeaderFn,
) EmailParserOption {
	return func(ep *EmailParser) {
		ep.headersParsers.From = fromHeaderParserFn
	}
}

// WithReplyToHeaderParser configures the parser used for the Reply-To header.
func WithReplyToHeaderParser(
	replyToHeaderParserFn parseAddressListHeaderFn,
) EmailParserOption {
	return func(ep *EmailParser) {
		ep.headersParsers.ReplyTo = replyToHeaderParserFn
	}
}

// WithToHeaderParser configures the parser used for the To header.
func WithToHeaderParser(
	toHeaderParserFn parseAddressListHeaderFn,
) EmailParserOption {
	return func(ep *EmailParser) {
		ep.headersParsers.To = toHeaderParserFn
	}
}

// WithCcHeaderParser configures the parser used for the Cc header.
func WithCcHeaderParser(
	ccHeaderParserFn parseAddressListHeaderFn,
) EmailParserOption {
	return func(ep *EmailParser) {
		ep.headersParsers.Cc = ccHeaderParserFn
	}
}

// WithBccHeaderParser configures the parser used for the Bcc header.
func WithBccHeaderParser(
	bccHeaderParserFn parseAddressListHeaderFn,
) EmailParserOption {
	return func(ep *EmailParser) {
		ep.headersParsers.Bcc = bccHeaderParserFn
	}
}

// WithMessageIdHeaderParser configures the parser used for the Message-ID header.
//
//nolint:revive // The exported name is kept for backward compatibility.
func WithMessageIdHeaderParser(
	messageIDHeaderParserFn parseMessageIDHeaderFn,
) EmailParserOption {
	return func(ep *EmailParser) {
		ep.headersParsers.MessageID = messageIDHeaderParserFn
	}
}

// WithInReplyToHeaderParser configures the parser used for the In-Reply-To header.
func WithInReplyToHeaderParser(
	inReplyToHeaderParserFn parseCommaSeparatedMessageIDHeaderFn,
) EmailParserOption {
	return func(ep *EmailParser) {
		ep.headersParsers.InReplyTo = inReplyToHeaderParserFn
	}
}

// WithReferencesHeaderParser configures the parser used for the References header.
func WithReferencesHeaderParser(
	referencesHeaderParserFn parseCommaSeparatedMessageIDHeaderFn,
) EmailParserOption {
	return func(ep *EmailParser) {
		ep.headersParsers.References = referencesHeaderParserFn
	}
}

// WithSubjectHeaderParser configures the parser used for the Subject header.
func WithSubjectHeaderParser(
	subjectHeaderParserFn parseStringHeaderFn,
) EmailParserOption {
	return func(ep *EmailParser) {
		ep.headersParsers.Subject = subjectHeaderParserFn
	}
}

// WithCommentsHeaderParser configures the parser used for the Comments header.
func WithCommentsHeaderParser(
	commentsHeaderParserFn parseStringHeaderFn,
) EmailParserOption {
	return func(ep *EmailParser) {
		ep.headersParsers.Comments = commentsHeaderParserFn
	}
}

// WithKeywordsHeaderParser configures the parser used for the Keywords header.
func WithKeywordsHeaderParser(
	keywordsHeaderParserFn parseCommaSeparatedStringHeaderFn,
) EmailParserOption {
	return func(ep *EmailParser) {
		ep.headersParsers.Keywords = keywordsHeaderParserFn
	}
}

// WithResentDateHeaderParser configures the parser used for the Resent-Date header.
func WithResentDateHeaderParser(
	resentDateHeaderParserFn parseDateHeaderFn,
) EmailParserOption {
	return func(ep *EmailParser) {
		ep.headersParsers.ResentDate = resentDateHeaderParserFn
	}
}

// WithResentFromHeaderParser configures the parser used for the Resent-From header.
func WithResentFromHeaderParser(
	resentFromHeaderParserFn parseAddressListHeaderFn,
) EmailParserOption {
	return func(ep *EmailParser) {
		ep.headersParsers.ResentFrom = resentFromHeaderParserFn
	}
}

// WithResentFromParser configures the parser used for the Resent-From header.
func WithResentFromParser(
	resentFromHeaderParserFn parseAddressListHeaderFn,
) EmailParserOption {
	return func(ep *EmailParser) {
		ep.headersParsers.ResentFrom = resentFromHeaderParserFn
	}
}

// WithResentSenderHeaderParser configures the parser used for the Resent-Sender header.
func WithResentSenderHeaderParser(
	resentSenderHeaderParserFn parseAddressHeaderFn,
) EmailParserOption {
	return func(ep *EmailParser) {
		ep.headersParsers.ResentSender = resentSenderHeaderParserFn
	}
}

// WithResentToHeaderParser configures the parser used for the Resent-To header.
func WithResentToHeaderParser(
	resentToHeaderParserFn parseAddressListHeaderFn,
) EmailParserOption {
	return func(ep *EmailParser) {
		ep.headersParsers.ResentTo = resentToHeaderParserFn
	}
}

// WithResentCcHeaderParser configures the parser used for the Resent-Cc header.
func WithResentCcHeaderParser(
	resentCcHeaderParserFn parseAddressListHeaderFn,
) EmailParserOption {
	return func(ep *EmailParser) {
		ep.headersParsers.ResentCc = resentCcHeaderParserFn
	}
}

// WithResentBccHeaderParser configures the parser used for the Resent-Bcc header.
func WithResentBccHeaderParser(
	resentBccHeaderParserFn parseAddressListHeaderFn,
) EmailParserOption {
	return func(ep *EmailParser) {
		ep.headersParsers.ResentBcc = resentBccHeaderParserFn
	}
}

// WithResentMessageIdHeaderParser configures the parser used for the Resent-Message-ID header.
//
//nolint:revive // The exported name is kept for backward compatibility.
func WithResentMessageIdHeaderParser(
	resentMessageIDHeaderParserFn parseMessageIDHeaderFn,
) EmailParserOption {
	return func(ep *EmailParser) {
		ep.headersParsers.ResentMessageID = resentMessageIDHeaderParserFn
	}
}

// WithContentTypeHeaderParser configures the parser used for the Content-Type header.
func WithContentTypeHeaderParser(
	contentTypeHeaderParserFn parseContentTypeHeaderFn,
) EmailParserOption {
	return func(ep *EmailParser) {
		ep.headersParsers.ContentType = contentTypeHeaderParserFn
	}
}

// WithContentDispositionHeaderParser configures the parser used for the Content-Disposition header.
func WithContentDispositionHeaderParser(
	contentDispositionHeaderParserFn parseContentDispositionHeaderFn,
) EmailParserOption {
	return func(ep *EmailParser) {
		ep.headersParsers.ContentDisposition = contentDispositionHeaderParserFn
	}
}

// WithExtraHeaderParser configures a parser for an additional named header.
func WithExtraHeaderParser(
	headerName string,
	extraHeaderParserFn parseStringHeaderFn,
) EmailParserOption {
	return func(ep *EmailParser) {
		ep.headersParsers.ExtraHeaders[strings.ToLower(headerName)] = extraHeaderParserFn
	}
}

// WithHeadersParsers replaces all header parsers used by an EmailParser.
func WithHeadersParsers(headersParsers HeadersParsers) EmailParserOption {
	return func(ep *EmailParser) {
		ep.headersParsers = headersParsers
	}
}
