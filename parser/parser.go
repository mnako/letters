// Package parser provides the capabilities for parsing an email
// io.Reader into an [email.Email]. The parser can receive options of
// type `Opt` which alter the parsing process.
package parser

import (
	"fmt"
	"io"
	"mime/multipart"
	"net/mail"
	"strings"
	"time"

	"github.com/mnako/letters/email"
	"golang.org/x/net/html/charset"
	"golang.org/x/text/encoding"
)

// UnknownContentTypeError reports an unknown Content Type
type UnknownContentTypeError struct {
	contentType string
}

func (e *UnknownContentTypeError) Error() string {
	return fmt.Sprintf("unknown Content-Type %q", e.contentType)
}

// typeOfProcessing determines the type of processing used by the
type typeOfProcessing string

const (
	wholeEmail    typeOfProcessing = "wholeEmail"
	headersOnly   typeOfProcessing = "headersOnly"
	noAttachments typeOfProcessing = "noAttachments"
)

// Opt is a parser option type provided as a closure to add options to a
// parser default instance instantiated by NewParser. The options are
// held in the [opts] subpackage.
type Opt func(p *Parser)

type Parser struct {
	// what parts of the email to process (default all)
	processType typeOfProcessing
	// email to be returned, for incremental processing
	email *email.Email

	// msg is the net/mail.Message used for deriving parts to build
	// the output email.
	msg *mail.Message

	// funcs that can be overridden by the user; defaults are set
	// attached by NewParser.
	// addressFunc : the function for processing email header addresses
	addressFunc func(string) (*mail.Address, error)
	// addressesFunc: the functionfor processing a list of email header
	// addresses
	addressesFunc func(list string) ([]*mail.Address, error)
	// dateFunc : the function for processing the email header Date
	dateFunc func(string) (time.Time, error)
	// fileFunc : a function for processing inline and attached files
	fileFunc func(*email.File) error

	// the main email encoding and content transfer encoding
	encoding encoding.Encoding
	cte      email.ContentTransferEncoding

	// debugging, for future use
	verbose bool
}

// NewParser initialises a new Parser. The default parser can be
// changed using options.
func NewParser(options ...Opt) *Parser {
	p := &Parser{

		// initialise main fields
		processType: wholeEmail,
		email:       &email.Email{},
		msg:         &mail.Message{},

		// initialise overrideable funcs
		// use net/mail.ParseAddress and ParseAddressList  as default
		// address parsers
		addressFunc:   mail.ParseAddress,
		addressesFunc: mail.ParseAddressList,
		// use net/mail.ParseDate as the default date parser
		dateFunc: mail.ParseDate,
		// by default write file io.Readers to email.File.Data.
		// User-supplied funcs might write files directly to disk, for
		// example, bypassing this step.
		fileFunc: func(f *email.File) error {
			var err error
			f.Data, err = io.ReadAll(f.Reader)
			return err
		},

		// debugging
		verbose: false,
	}

	for _, opt := range options {
		opt(p)
	}
	return p
}

// Parse is the main entry point of letters
func (p *Parser) Parse(r io.Reader) (*email.Email, error) {
	var err error
	p.msg, err = mail.ReadMessage(r)
	if err != nil {
		return nil, fmt.Errorf("cannot read message: %w", err)
	}

	// parse headers
	err = p.parseHeaders()
	if err != nil {
		return nil, fmt.Errorf("cannot parse headers: %w", err)
	}
	if p.processType == headersOnly {
		return p.email, nil
	}

	h := p.email.Headers

	// determine encoding
	p.encoding, _ = charset.Lookup(h.ContentType.Params["charset"])
	p.cte, err = extractContentTransferEncoding(p.msg.Header.Get("Content-Transfer-Encoding"))
	if err != nil {
		return nil, fmt.Errorf("cannot parse Content-Transfer-Encoding: %w", err)
	}

	switch ct := string(h.ContentType.ContentType); { // true switch
	case
		ct == string(email.ContentTypeTextPlain),
		ct == string(email.ContentTypeTextEnriched),
		ct == string(email.ContentTypeTextHtml):
		// parse body
		err = p.parseBody()
		if err != nil {
			return nil, err
		}
	case strings.HasPrefix(ct, string(email.ContentTypeMultipartPrefix)):
		// parse parts
		err = p.parsePart(p.msg.Body, h.ContentType, h.ContentType.Params["boundary"])
		if err != nil {
			return nil, err
		}
	default:
		// parse attachment
		err = p.parseFile(
			p.msg.Body,
			email.AttachedFileType,
			p.email.Headers.ContentType,
			p.email.Headers.ContentDisposition,
			p.cte,
		)
		if err != nil {
			return nil, err
		}
	}
	return p.email, err
}

// parsePart parses the parts of a multipart message
func (p *Parser) parsePart(msg io.Reader, parentContentType email.ContentTypeHeader, boundary string) error {

	multipartReader := multipart.NewReader(msg, boundary)
	if multipartReader == nil {
		return nil
	}

	for {
		part, err := multipartReader.NextPart()
		if err == io.EOF {
			break
		}
		if err != nil {
			return fmt.Errorf("cannot read part: %w", err)
		}

		// extract content-type
		partContentType, err := extractContentTypeHeader(part.Header.Get("Content-Type"))
		if err != nil {
			return fmt.Errorf("cannot parse part Content-Type: %w", err)
		}

		// extract charset
		charsetLabel := partContentType.Params["charset"]
		if charsetLabel == "" {
			charsetLabel = parentContentType.Params["charset"]
		}
		enc, _ := charset.Lookup(charsetLabel)

		// extract content-transfer-encoding
		cte, err := extractContentTransferEncoding(part.Header.Get("Content-Transfer-Encoding"))
		if err != nil {
			return fmt.Errorf("cannot parse part Content-Transfer-Encoding: %w", err)
		}

		// extract content-disposition
		cdh, err := extractContentDisposition(part.Header.Get("Content-Disposition"))
		if err != nil {
			return fmt.Errorf("cannot parse part Content-Disposition: %w", err)
		}

		// commence extraction of data with attached file
		if string(cdh.ContentDisposition) == "attachment" && string(email.AttachedFileType) == "attached" {
			err = p.parseFile(
				part,
				email.AttachedFileType,
				partContentType,
				cdh,
				cte,
			)
			if err != nil {
				return fmt.Errorf("cannot parse attached file: %w", err)
			}
			continue
		}

		// process text plain content
		if string(partContentType.ContentType) == string(email.ContentTypeTextPlain) {
			partTextBody, err := p.parseText(part, enc, cte)
			if err != nil {
				return fmt.Errorf("cannot parse plain text: %w", err)
			}
			if len(p.email.Text) > 0 { // add separator
				p.email.Text += "\n\n"
			}
			p.email.Text += partTextBody
			continue
		}

		// process text enriched content
		if string(partContentType.ContentType) == string(email.ContentTypeTextEnriched) {
			partEnrichedText, err := p.parseText(part, enc, cte)
			if err != nil {
				return fmt.Errorf("cannot parse enriched text: %w", err)
			}
			p.email.EnrichedText += partEnrichedText
			continue
		}

		// process html content
		if string(partContentType.ContentType) == string(email.ContentTypeTextHtml) {
			partHtmlBody, err := p.parseText(part, p.encoding, cte)
			if err != nil {
				return fmt.Errorf("cannot parse html text: %w", err)
			}
			p.email.HTML += partHtmlBody
			continue
		}

		// recursive call to parsePart
		if strings.HasPrefix(partContentType.ContentType, string(email.ContentTypeMultipartPrefix)) {
			err := p.parsePart(part, partContentType, partContentType.Params["boundary"])
			if err != nil {
				return fmt.Errorf("cannot parse nested part: %w", err)
			}
			continue
		}

		// process inline file
		if isInlineFile(partContentType, parentContentType, cdh) {
			switch p.processType {
			case headersOnly:
				continue
			case noAttachments:
				continue
			}
			err = p.parseFile(
				part,
				email.InlineFileType,
				partContentType,
				cdh,
				cte,
			)
			if err != nil {
				return fmt.Errorf("cannot parse inline file: %w", err)
			}
			continue
		}

		// process attached file
		if isAttachedFile(partContentType, parentContentType, cdh) {
			switch p.processType {
			case headersOnly:
				continue
			case noAttachments:
				continue
			}
			err := p.parseFile(
				part,
				email.AttachedFileType,
				partContentType,
				cdh,
				cte,
			)
			if err != nil {
				return fmt.Errorf("cannot parse attached file: %w", err)
			}
			continue
		}

		// fallthrough error
		return &UnknownContentTypeError{contentType: parentContentType.ContentType}
	}

	return nil
}
