package parser

import (
	"fmt"
	"io"
	"path/filepath"

	"github.com/mnako/letters/decoders"
	"github.com/mnako/letters/email"
)

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
func (p *Parser) parseFile(r io.Reader, ci *email.ContentInfo) error {

	var err error
	file := &email.File{
		FileType:    ci.Disposition,
		ContentInfo: ci,
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
	if name, ok := file.ContentInfo.DispositionParams["filename"]; ok {
		tmpFileName = name
	} else {
		// fallback to ContentTypeHeader
		if name, ok := file.ContentInfo.TypeParams["name"]; ok {
			tmpFileName = name
		} else {
			// Make up a unique name if none exists. Todo: Suffix ideally needed.
			tmpFileName = fmt.Sprintf("attachment_%d_%s", len(p.email.Files), ci.Disposition)
		}
	}
	file.Name = filepath.Base(filepath.Clean(tmpFileName))

	file.Reader = decoders.DecodeContent(r, ci)
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
