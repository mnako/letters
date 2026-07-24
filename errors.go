package letters

import "errors"

var (
	// ErrUnknownCharset indicates that a MIME header uses an unsupported charset.
	ErrUnknownCharset = errors.New(
		"letters.decoders.decodeHeader.CharsetReader: cannot lookup encoding",
	)

	// ErrUnknownContentDisposition indicates an unsupported Content-Disposition.
	ErrUnknownContentDisposition = errors.New(
		"letters.parsers.parseContentDisposition: unknown Content-Disposition",
	)

	// ErrUnknownContentTransferEncoding indicates an unsupported transfer encoding.
	ErrUnknownContentTransferEncoding = errors.New(
		"letters.parsers.parseContentTransferEncoding: unknown Content-Transfer-Encoding",
	)
)
