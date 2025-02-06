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

	// the main email content info
	contentInfo *email.ContentInfo

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

	// extract content information
	p.contentInfo, err = email.ExtractContentInfo(p.msg.Header, nil)
	if err != nil {
		return nil, fmt.Errorf("cannot extract content: %w", err)
	}

	// parse headers
	err = p.parseHeaders()
	if err != nil {
		return nil, fmt.Errorf("cannot parse headers: %w", err)
	}
	if p.processType == headersOnly {
		return p.email, nil
	}

	switch ct := p.contentInfo.Type; { // true switch
	case ct == "text/plain", ct == "text/enriched", ct == "text/html":
		// parse body
		err = p.parseBody()
		if err != nil {
			return nil, err
		}
	case strings.HasPrefix(ct, "multipart/"):
		// parse parts
		err = p.parsePart(
			p.msg.Body,
			p.contentInfo,
			p.contentInfo.TypeParams["boundary"],
		)
		if err != nil {
			return nil, err
		}
	default:
		// parse attachment
		err = p.parseFile(p.msg.Body, p.contentInfo)
		if err != nil {
			return nil, err
		}
	}
	return p.email, err
}

// parsePart parses the parts of a multipart message
func (p *Parser) parsePart(msg io.Reader, parentCI *email.ContentInfo, boundary string) error {

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

		// extract content information
		contentInfo, err := email.ExtractContentInfo(part.Header, p.contentInfo)
		if err != nil {
			return fmt.Errorf("content extraction errror: %w", err)
		}

		// commence extraction of data with attached file
		if contentInfo.Disposition == "attachment" {
			err = p.parseFile(
				part,
				contentInfo,
			)
			if err != nil {
				return fmt.Errorf("cannot parse attached file: %w", err)
			}
			continue
		}

		// process text plain content
		if contentInfo.Type == "text/plain" {
			partTextBody, err := p.parseText(part, contentInfo)
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
		if contentInfo.Type == "text/enriched" {
			partEnrichedText, err := p.parseText(part, contentInfo)
			if err != nil {
				return fmt.Errorf("cannot parse enriched text: %w", err)
			}
			p.email.EnrichedText += partEnrichedText
			continue
		}

		// process html content
		if contentInfo.Type == "text/html" {
			partHtmlBody, err := p.parseText(part, contentInfo)
			if err != nil {
				return fmt.Errorf("cannot parse html text: %w", err)
			}
			p.email.HTML += partHtmlBody
			continue
		}

		// recursive call to parsePart
		if strings.HasPrefix(contentInfo.Type, "multipart") {
			err := p.parsePart(part, contentInfo, contentInfo.TypeParams["boundary"])
			if err != nil {
				return fmt.Errorf("cannot parse nested part: %w", err)
			}
			continue
		}

		// process inline file
		if contentInfo.IsInlineFile(contentInfo) {
			if p.processType != wholeEmail {
				continue
			}
			err = p.parseFile(part, contentInfo)
			if err != nil {
				return fmt.Errorf("cannot parse inline file: %w", err)
			}
			continue
		}

		// process attached file
		if contentInfo.IsAttachedFile(contentInfo) {
			if p.processType != wholeEmail {
				continue
			}
			err := p.parseFile(part, contentInfo)
			if err != nil {
				return fmt.Errorf("cannot parse attached file: %w", err)
			}
			continue
		}

		// fallthrough error
		return &UnknownContentTypeError{contentType: parentCI.Type}
	}

	return nil
}
