package letters

type (
	EmailBodyFilter func(cth ContentTypeHeader) bool
	EmailFileFilter func(cth ContentTypeHeader) bool
)

func NoBodies(_ ContentTypeHeader) bool {
	return false
}

func AllBodies(_ ContentTypeHeader) bool {
	return true
}

func NoFiles(_ ContentTypeHeader) bool {
	return false
}

func AllFiles(_ ContentTypeHeader) bool {
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
