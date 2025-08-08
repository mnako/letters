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
	parseMessageIdHeaderFn               func(string) MessageId
	parseCommaSeparatedMessageIdHeaderFn func(string) []MessageId
	parseContentDispositionHeaderFn      func(string) (ContentDispositionHeader, error)
	parseContentTypeHeaderFn             func(string) (ContentTypeHeader, error)
)

type HeadersParsers struct {
	Date               parseDateHeaderFn
	Sender             parseAddressHeaderFn
	From               parseAddressListHeaderFn
	ReplyTo            parseAddressListHeaderFn
	To                 parseAddressListHeaderFn
	Cc                 parseAddressListHeaderFn
	Bcc                parseAddressListHeaderFn
	MessageID          parseMessageIdHeaderFn
	InReplyTo          parseCommaSeparatedMessageIdHeaderFn
	References         parseCommaSeparatedMessageIdHeaderFn
	Subject            parseStringHeaderFn
	Comments           parseStringHeaderFn
	Keywords           parseCommaSeparatedStringHeaderFn
	ResentDate         parseDateHeaderFn
	ResentFrom         parseAddressListHeaderFn
	ResentSender       parseAddressHeaderFn
	ResentTo           parseAddressListHeaderFn
	ResentCc           parseAddressListHeaderFn
	ResentBcc          parseAddressListHeaderFn
	ResentMessageID    parseMessageIdHeaderFn
	ContentType        parseContentTypeHeaderFn
	ContentDisposition parseContentDispositionHeaderFn
	ExtraHeaders       map[string]parseStringHeaderFn
}

func WithDateHeaderParser(
	dateHeaderParserFn parseDateHeaderFn,
) EmailParserOption {
	return func(ep *EmailParser) {
		ep.headersParsers.Date = dateHeaderParserFn
	}
}

func WithSenderHeaderParser(
	senderHeaderParserFn parseAddressHeaderFn,
) EmailParserOption {
	return func(ep *EmailParser) {
		ep.headersParsers.Sender = senderHeaderParserFn
	}
}

func WithFromHeaderParser(
	fromHeaderParserFn parseAddressListHeaderFn,
) EmailParserOption {
	return func(ep *EmailParser) {
		ep.headersParsers.From = fromHeaderParserFn
	}
}

func WithReplyToHeaderParser(
	replyToHeaderParserFn parseAddressListHeaderFn,
) EmailParserOption {
	return func(ep *EmailParser) {
		ep.headersParsers.ReplyTo = replyToHeaderParserFn
	}
}

func WithToHeaderParser(
	toHeaderParserFn parseAddressListHeaderFn,
) EmailParserOption {
	return func(ep *EmailParser) {
		ep.headersParsers.To = toHeaderParserFn
	}
}

func WithCcHeaderParser(
	ccHeaderParserFn parseAddressListHeaderFn,
) EmailParserOption {
	return func(ep *EmailParser) {
		ep.headersParsers.Cc = ccHeaderParserFn
	}
}

func WithBccHeaderParser(
	bccHeaderParserFn parseAddressListHeaderFn,
) EmailParserOption {
	return func(ep *EmailParser) {
		ep.headersParsers.Bcc = bccHeaderParserFn
	}
}

func WithMessageIdHeaderParser(
	messageIDHeaderParserFn parseMessageIdHeaderFn,
) EmailParserOption {
	return func(ep *EmailParser) {
		ep.headersParsers.MessageID = messageIDHeaderParserFn
	}
}

func WithInReplyHeaderParser(
	inReplyHeaderParserFn parseCommaSeparatedMessageIdHeaderFn,
) EmailParserOption {
	return func(ep *EmailParser) {
		ep.headersParsers.InReplyTo = inReplyHeaderParserFn
	}
}

func WithReferencesHeaderParser(
	referencesHeaderParserFn parseCommaSeparatedMessageIdHeaderFn,
) EmailParserOption {
	return func(ep *EmailParser) {
		ep.headersParsers.References = referencesHeaderParserFn
	}
}

func WithSubjectHeaderParser(
	subjectHeaderParserFn parseStringHeaderFn,
) EmailParserOption {
	return func(ep *EmailParser) {
		ep.headersParsers.Subject = subjectHeaderParserFn
	}
}

func WithCommentsHeaderParser(
	commentsHeaderParserFn parseStringHeaderFn,
) EmailParserOption {
	return func(ep *EmailParser) {
		ep.headersParsers.Comments = commentsHeaderParserFn
	}
}

func WithKeywordsHeaderParser(
	keywordsHeaderParserFn parseCommaSeparatedStringHeaderFn,
) EmailParserOption {
	return func(ep *EmailParser) {
		ep.headersParsers.Keywords = keywordsHeaderParserFn
	}
}

func WithResentDateHeaderParser(
	resentDateHeaderParserFn parseDateHeaderFn,
) EmailParserOption {
	return func(ep *EmailParser) {
		ep.headersParsers.ResentDate = resentDateHeaderParserFn
	}
}

func WithResentFromHeaderParser(
	resentFromHeaderParserFn parseAddressListHeaderFn,
) EmailParserOption {
	return func(ep *EmailParser) {
		ep.headersParsers.ResentFrom = resentFromHeaderParserFn
	}
}

func WithResentFromParser(
	resentFromHeaderParserFn parseAddressListHeaderFn,
) EmailParserOption {
	return func(ep *EmailParser) {
		ep.headersParsers.ResentFrom = resentFromHeaderParserFn
	}
}

func WithResentSenderHeaderParser(
	resentSenderHeaderParserFn parseAddressHeaderFn,
) EmailParserOption {
	return func(ep *EmailParser) {
		ep.headersParsers.ResentSender = resentSenderHeaderParserFn
	}
}

func WithResentToHeaderParser(
	resentToHeaderParserFn parseAddressListHeaderFn,
) EmailParserOption {
	return func(ep *EmailParser) {
		ep.headersParsers.ResentTo = resentToHeaderParserFn
	}
}

func WithResentCcHeaderParser(
	resentCcHeaderParserFn parseAddressListHeaderFn,
) EmailParserOption {
	return func(ep *EmailParser) {
		ep.headersParsers.ResentCc = resentCcHeaderParserFn
	}
}

func WithResentBccHeaderParser(
	resentBccHeaderParserFn parseAddressListHeaderFn,
) EmailParserOption {
	return func(ep *EmailParser) {
		ep.headersParsers.ResentBcc = resentBccHeaderParserFn
	}
}

func WithResentMessageIdHeaderParser(
	resentMessageIDHeaderParserFn parseMessageIdHeaderFn,
) EmailParserOption {
	return func(ep *EmailParser) {
		ep.headersParsers.ResentMessageID = resentMessageIDHeaderParserFn
	}
}

func WithContentTypeHeaderParser(
	contentTypeHeaderParserFn parseContentTypeHeaderFn,
) EmailParserOption {
	return func(ep *EmailParser) {
		ep.headersParsers.ContentType = contentTypeHeaderParserFn
	}
}

func WithContentDispositionHeaderParser(
	contentDispositionHeaderParserFn parseContentDispositionHeaderFn,
) EmailParserOption {
	return func(ep *EmailParser) {
		ep.headersParsers.ContentDisposition = contentDispositionHeaderParserFn
	}
}

func WithExtraHeaderParser(
	headerName string,
	extraHeaderParserFn parseStringHeaderFn,
) EmailParserOption {
	return func(ep *EmailParser) {
		ep.headersParsers.ExtraHeaders[strings.ToLower(headerName)] = extraHeaderParserFn
	}
}

func WithHeadersParsers(headersParsers HeadersParsers) EmailParserOption {
	return func(ep *EmailParser) {
		ep.headersParsers = headersParsers
	}
}
