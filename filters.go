package letters

type (
	EmailBodyFilter func(cth ContentTypeHeader) bool
	EmailFileFilter func(cth ContentTypeHeader, cdh ContentDispositionHeader) bool
)

func NoBodies(_ ContentTypeHeader) bool {
	return false
}

func AllBodies(_ ContentTypeHeader) bool {
	return true
}

func NoFiles(_ ContentTypeHeader, __ ContentDispositionHeader) bool {
	return false
}

func AllFiles(_ ContentTypeHeader, __ ContentDispositionHeader) bool {
	return true
}

func WithBodyFilter(bodyFilter EmailBodyFilter) EmailParserOption {
	return func(ep *EmailParser) {
		ep.bodyFilter = bodyFilter
	}
}

func WithFileFilter(fileFilter EmailFileFilter) EmailParserOption {
	return func(ep *EmailParser) {
		ep.fileFilter = fileFilter
	}
}

// WithUnquotedAtInDisplayName allows to parse email with @ in display name
// this provides a flexibility in case some clients send email with @ in display name
// which is not allowed by rfc5322
func WithUnquotedAtInDisplayName(allow bool) EmailParserOption {
	return func(ep *EmailParser) {
		ep.allowUnquotedAtInDisplayName = allow
	}
}
