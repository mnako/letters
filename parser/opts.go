package parser

import (
	"fmt"
	"io"
	"net/mail"
	"os"
	"path/filepath"
	"time"

	"github.com/mnako/letters/email"
)

// WithVerbose is presently a no-op option for testing if options are
// working as anticipated.
func WithVerbose() Opt {
	return func(p *Parser) {
		p.verbose = true
	}
}

// WithHeadersOnly only parses email headers, which typically provides a
// substantial speedup in processing.
func WithHeadersOnly() Opt {
	return func(p *Parser) {
		p.processType = headersOnly
	}
}

// WithoutAttachments skips parsing email attachments, which often
// provides a speedup in processing.
func WithoutAttachments() Opt {
	return func(p *Parser) {
		p.processType = noAttachments
	}
}

// WithCustomDateFunc allows for the provision of a custom date parsing
// func.
func WithCustomDateFunc(df func(string) (time.Time, error)) Opt {
	return func(p *Parser) {
		p.dateFunc = df
	}
}

// WithCustomAddressFunc allows for the provision of a custom func for
// parsing an email name/address combination.
func WithCustomAddressFunc(af func(string) (*mail.Address, error)) Opt {
	return func(p *Parser) {
		p.addressFunc = af
	}
}

// WithCustomAddressesFunc allows for the provision of a custom func for
// parsing lists of email names and addresses from strings.
func WithCustomAddressesFunc(af func(list string) ([]*mail.Address, error)) Opt {
	return func(p *Parser) {
		p.addressesFunc = af
	}
}

// WithCustomFileFunc allows for the provision of a custom func for
// reading a file attachment io.Reader. Note that the io.Reader provided
// by the underlying net/mail package is not concurrent safe. The reader
// for each file should be drained in turn, otherwise it is likely to be
// unexpectedly truncated.
//
// Note that all the fields of Parser and email.File are
// available for custom uses, such as file filtering, file saving or
// sending files over the network.
func WithCustomFileFunc(ff func(*email.File) error) Opt {
	return func(p *Parser) {
		p.fileFunc = ff
	}
}

// WithSaveFilesToDirectory is an example WithCustomFileFunc showing how
// to inject a custom-defined func into the parser to save inline and
// attached files from an email to the supplied directory.
//
// Note that all the public fields in email.Email are available, as are
// the fields of email.File. For example the message ID of each email
// could be used to make a directory, and then each file saved in
// sequence using its file name suffix (if available).
//
// Caution should be used using the filename provided in
// internet-provided emails, although some cleaning is done by the
// parsers module -- see Parser.parseFile.
func WithSaveFilesToDirectory(dir string) Opt {
	return func(p *Parser) {
		// attach the inline func to p.fileFunc
		p.fileFunc = func(ef *email.File) error {
			f, err := os.Create(filepath.Join(dir, ef.Name))
			if err != nil {
				return fmt.Errorf("file creation error %w", err)
			}
			defer f.Close()
			_, err = io.Copy(f, ef.Reader)
			if err != nil {
				return fmt.Errorf("file saving error %w", err)
			}
			return nil
		}
	}
}
