package parser

import (
	"fmt"
	"io"
	"path/filepath"

	"github.com/mnako/letters/decoders"
	"github.com/mnako/letters/email"
)

// isInlineFile reports if the content type of the part describes an
// inline file.
func isInlineFile(
	ct email.ContentTypeHeader,
	parentCT email.ContentTypeHeader,
	cdh email.ContentDispositionHeader,
) bool {
	switch {
	case string(cdh.ContentDisposition) == string(email.InlineFileType):
		return true
	case string(ct.ContentType) == string(email.ContentTypeTextPlain):
		return false
	case string(ct.ContentType) == string(email.ContentTypeTextEnriched):
		return false
	case string(ct.ContentType) == string(email.ContentTypeTextHtml):
		return false
	case string(parentCT.ContentType) == string(email.ContentTypeMultipartRelated):
		return true
	}
	return false
}

// isAttachedFile reports if the content type of the part describes an
// attached file.
func isAttachedFile(
	ct email.ContentTypeHeader,
	parentCT email.ContentTypeHeader,
	cdh email.ContentDispositionHeader,
) bool {
	switch {
	case string(cdh.ContentDisposition) == string(email.AttachedFileType):
		return true
	case string(ct.ContentType) == string(email.ContentTypeTextPlain):
		return false
	case string(ct.ContentType) == string(email.ContentTypeTextEnriched):
		return false
	case string(ct.ContentType) == string(email.ContentTypeTextHtml):
		return false
	case string(parentCT.ContentType) == string(email.ContentTypeMultipartMixed):
		return true
	case string(parentCT.ContentType) == string(email.ContentTypeMultipartParallel):
		return true
	}
	return false
}

// parseFile parses inline and attached files from email parts, using
// the parser.fileFunc to process the io.Reader returned by
// decoders.DecodeContent. By default this func will write the reader
// into file.Data. User-supplied funcs, wich might be closures, can be
// used to avoid the allocation of bytes here. However users should be
// careful to process the entire file before continuing since the end of
// this function will terminate the underlying io.Reader with unexpected
// results for the consumer.
//
// Files that are successfully parsed are added to parser.email.Files.
func (p *Parser) parseFile(
	r io.Reader,
	fileType email.FileType,
	cth email.ContentTypeHeader,
	cdh email.ContentDispositionHeader,
	cte email.ContentTransferEncoding,
) error {

	var err error
	file := &email.File{
		FileType:                 fileType,
		ContentTypeHeader:        cth,
		ContentDispositionHeader: cdh,
	}

	// extract file name from filename or name field
	// RFC 2183 limits filenames to the US-ASCII printable range only.
	// However, modern standards accommodate encoding the parameter in
	// transport while representing Unicode characters properly at
	// destination. To safely transmit symbols and international
	// characters, the modern RFC 6266 standard actually specifies using
	// RFC 5987 encoding paired with the filename* parameter:
	//  Content-Disposition: attachment;
	//                  filename="fallback.txt";
	//                  filename*=UTF-8''my-%C3%BC-file.txt
	// --
	// File path security traversal guidelines
	// https://www.stackhawk.com/blog/golang-path-traversal-guide-examples-and-prevention/
	var tmpFileName string
	if name, ok := file.ContentDispositionHeader.Params["filename"]; ok {
		tmpFileName = name
	} else {
		// fallback to ContentTypeHeader
		if name, ok := file.ContentTypeHeader.Params["name"]; ok {
			tmpFileName = name
		} else {
			// Make up a unique name if none exists. Todo: Suffix ideally needed.
			tmpFileName = fmt.Sprintf("attachment_%d_%s", len(p.email.Files), fileType)
		}
	}
	file.Name = filepath.Base(filepath.Clean(tmpFileName))

	file.Reader = decoders.DecodeContent(r, nil, cte)
	// parser.fileFunc is a pluggable file reader with the signature
	// func(*email.File) error.
	// The fileFunc may be customised through parser.NewParser(...opts).
	err = p.fileFunc(file)
	if err != nil {
		return fmt.Errorf("could not read attachment data: %w", err)
	}

	p.email.Files = append(p.email.Files, file)
	return nil

}
