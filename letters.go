// # Letters, or how to parse emails in Go
//
// Letters is a minimalistic Golang library for parsing plaintext and
// MIME emails.
//
// It correctly handles text and MIME mime-types, Base64 and Quoted-Printable
// Content-Transfer-Encoding, as well as any text encoding that Golang
// standard library is capable of handling. Letters will parse an email into
// a simple struct with standard headers and text, enriched text, and HTML
// content, and decode inline and attached files.
//
// Letters also supports options for skipping processing parts of
// messages and providing custom processing functions.
//
// # Quickstart
//
// Install
//
//	go get github.com/mnako/letters@v0.2.3
//
// Parse a raw email from a Reader:
//
//	email, err := letters.ParseEmail(r)
//
//	if err != nil {
//	    log.Fatal(err)
//	}
//
// and you can access the common headers:
//
//	email.Headers.Sender
//	// mail.Address{Name: "Alice Sender", Address: "alice.sender@example.com"}
//
//	email.Headers.From
//	// []mail.Address{
//	//  {Name: "Alice Sender", Address: "alice.sender@example.com"},
//	//  {Name: "Alice Sender", Address: "alice.sender@example.net"},
//	// }
//
//	email.Headers.Subject
//	// "📧 Test English Pangrams"
//
//	email.Headers.To
//	// []mail.Address{
//	//  {Name: "Bob Recipient", Address: "bob.recipient@example.com"},
//	//  {Name: "Carol Recipient", Address: "carol.recipient@example.com"},
//	// }
//
//	email.Headers.Cc
//	// []mail.Address{
//	//  {Name: "Dan Recipient", Address: "dan.recipient@example.com"},
//	//  {Name: "Eve Recipient", Address: "eve.recipient@example.com"},
//	// }
//
//	email.Headers.Bcc
//	// []mail.Address{
//	//  {Name: "Frank Recipient", Address: "frank.recipient@example.com"},
//	//  {Name: "Grace Recipient", Address: "grace.recipient@example.com"},
//	// }
//
// get custom headers:
//
//	email.Headers.ExtraHeaders
//	// map[string][]string{
//	//    "X-Clacks-Overhead": {"GNU Terry Pratchett"},
//	// }
//
// get decoded bodies:
//
//	email.Text
//	// "The quick brown fox jumps over a lazy dog..."
//
//	email.HTML
//	// "<html><div dir="ltr"><p>The quick brown fox jumps over a lazy dog..."
//
// Both inline and attached files are stored in a slice. By default these
// are read into a `Data` []byte slice but direct access can be made to the
// underlying `io.Reader` by using a custom file processing func.
//
//	Files: []*email.File{
//		{
//			FileType: "inline",
//			ContentTypeHeader: email.ContentTypeHeader{
//				ContentType: "image/jpeg",
//				Params: map[string]string{
//					"name": "inline-jpg-image-name.jpg",
//				},
//			},
//			ContentDispositionHeader: email.ContentDispositionHeader{
//				ContentDisposition: "inline",
//				Params: map[string]string{
//					"filename": "inline-jpg-image-filename.jpg",
//				},
//			},
//			Name: "inline-jpg-image-filename.jpg",
//			Data: []byte{
//				255, 216, 255, 219, 0, 67, 0, 3, 2, 2, 2, 2, 2, 3, 2, 2, 2, 3, 3,
//				3, 3, 4, 6, 4, 4, 4, 4, 4, 8, 6, 6, 5, 6, 9, 8, 10, 10, 9, 8, 9, 9,
//				10, 12, 15, 12, 10, 11, 14, 11, 9, 9, 13, 17, 13, 14, 15, 16, 16,
//				17, 16, 10, 12, 18, 19, 18, 16, 19, 15, 16, 16, 16, 255, 201, 0,
//				11, 8, 0, 1, 0, 1, 1, 1, 17, 0, 255, 204, 0, 6, 0, 16, 16, 5, 255,
//				218, 0, 8, 1, 1, 0, 0, 63, 0, 210, 207, 32, 255, 217,
//			},
//		},
//		{
//			FileType: "attached",
//			Name:     "attached-pdf-filename.pdf",
//			ContentTypeHeader: email.ContentTypeHeader{
//				ContentType: "application/pdf",
//				Params: map[string]string{
//					"name": "attached-pdf-name.pdf",
//				},
//			},
//			ContentDispositionHeader: email.ContentDispositionHeader{
//				ContentDisposition: "attachment",
//				Params: map[string]string{
//					"filename": "attached-pdf-filename.pdf",
//				},
//			},
//			Data: []byte{
//				37, 80, 68, 70, 45, 49, 46, 13, 116, 114, 97, 105, 108, 101, 114,
//				60, 60, 47, 82, 111, 111, 116, 60, 60, 47, 80, 97, 103, 101, 115,
//				60, 60, 47, 75, 105, 100, 115, 91, 60, 60, 47, 77, 101, 100, 105,
//				97, 66, 111, 120, 91, 48, 32, 48, 32, 51, 32, 51, 93, 62, 62, 93,
//				62, 62, 62, 62, 62, 62,
//			},
//		},
//	}
//
// # Options
//
// Various options are provided for customising the Parser, including:
//
//	func WithCustomAddressFunc(af func(string) (*mail.Address, error)) Opt
//	func WithCustomAddressesFunc(af func(list string) ([]*mail.Address, error)) Opt
//	func WithCustomDateFunc(df func(string) (time.Time, error)) Opt
//	func WithCustomFileFunc(ff func(*email.File) error) Opt
//	func WithSaveFilesToDirectory(dir string) Opt
//	func WithHeadersOnly() Opt
//	func WithoutAttachments() Opt
//	func WithVerbose() Opt
//
// The `WithoutAttachments` and `WithHeadersOnly` options determine if
// only part of an email will be processed.
//
// The date and address "With" options allow the provision of custom
// funcs to override the [net/mail] funcs normally used. For example it
// might be necessary to extend the date parsing capabilities to deal
// with poorly formatted date strings produced by older SMTP servers.
//
// The `WithCustomFileFunc` allows the provision of a custom func for
// saving, filtering and/or processing of inline or attached files
// without reading them first into an `email.File.Data` []byte slice
// first, which is the default behaviour. The `WithSaveFilesToDirectory`
// option is an example of such a custom func.
//
// An example:
//
//	opt := parser.WithHeadersOnly() // pass the headers only option
//	p := letters.NewParser(opt, parser.WithVerbose()) // options can be chained
//	parsedEmail, err := p.Parse(rawEmail)
//	if err != nil {
//		return fmt.Errorf("error while parsing email headers: %s", err)
//	}
//
// See the [parser] package and tests for more details.
//
// # Language and Encoding Support
//
// The same parser and methods will work for other languages, text encodings,
// and transfer-encodings.
//
// [net/mail]: https://pkg.go.dev/net/mail
package letters

import (
	"github.com/mnako/letters/parser"
)

func NewParser(options ...parser.Opt) *parser.Parser {
	return parser.NewParser(options...)
}
