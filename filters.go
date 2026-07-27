package letters

type (
	// EmailBodyFilter decides whether an email body should be parsed.
	EmailBodyFilter func(cth ContentTypeHeader) bool

	// EmailFileFilter decides whether an email file should be parsed.
	EmailFileFilter func(cth ContentTypeHeader, cdh ContentDispositionHeader) bool
)

// NoBodies excludes all email bodies from parsing.
func NoBodies(_ ContentTypeHeader) bool {
	return false
}

// AllBodies includes all email bodies in parsing.
func AllBodies(_ ContentTypeHeader) bool {
	return true
}

// NoFiles excludes all email files from parsing.
func NoFiles(_ ContentTypeHeader, _ ContentDispositionHeader) bool {
	return false
}

// AllFiles includes all email files in parsing.
func AllFiles(_ ContentTypeHeader, _ ContentDispositionHeader) bool {
	return true
}

// WithBodyFilter configures the filter used to select email bodies for parsing.
func WithBodyFilter(bodyFilter EmailBodyFilter) EmailParserOption {
	return func(ep *EmailParser) {
		ep.bodyFilter = bodyFilter
	}
}

// WithFileFilter configures the filter used to select email files for parsing.
func WithFileFilter(fileFilter EmailFileFilter) EmailParserOption {
	return func(ep *EmailParser) {
		ep.fileFilter = fileFilter
	}
}
