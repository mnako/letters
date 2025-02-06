# Letters, or how to parse emails in Go

[![Test](https://github.com/mnako/letters/actions/workflows/test.yml/badge.svg)](https://github.com/mnako/letters/actions/workflows/test.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/mnako/letters)](https://goreportcard.com/report/github.com/mnako/letters)

**Letters** is a minimalistic Golang library for parsing plaintext and MIME
emails.

It correctly handles text and MIME mime-types, Base64 and Quoted-Printable 
Content-Transfer-Encoding, as well as any text encoding that Golang 
standard library is capable of handling. Letters will parse an email into 
a simple struct with standard headers and text, enriched text, and HTML 
content, and decode inline and attached files.

Letters also supports options for skipping parts of messages and
providing custom processing functions.

## Quickstart

Install

```
go get github.com/mnako/letters@v0.2.3
```

Parse a raw email from a Reader:

```go
email, err := letters.ParseEmail(r)
if err != nil {
    log.Fatal(err)
}
```

and you can access the common headers:

```go
email.Headers.Sender
// mail.Address{Name: "Alice Sender", Address: "alice.sender@example.com"}

email.Headers.From
// []mail.Address{
//  {Name: "Alice Sender", Address: "alice.sender@example.com"}, 
//  {Name: "Alice Sender", Address: "alice.sender@example.net"},
// }

email.Headers.Subject
// "📧 Test English Pangrams"

email.Headers.To
// []mail.Address{
//  {Name: "Bob Recipient", Address: "bob.recipient@example.com"}, 
//  {Name: "Carol Recipient", Address: "carol.recipient@example.com"},
// }

email.Headers.Cc
// []mail.Address{
//  {Name: "Dan Recipient", Address: "dan.recipient@example.com"}, 
//  {Name: "Eve Recipient", Address: "eve.recipient@example.com"},
// }

email.Headers.Bcc
// []mail.Address{
//  {Name: "Frank Recipient", Address: "frank.recipient@example.com"}, 
//  {Name: "Grace Recipient", Address: "grace.recipient@example.com"},
// }
```

get custom headers:

```go
email.Headers.ExtraHeaders
// map[string][]string{
//    "X-Clacks-Overhead": {"GNU Terry Pratchett"},
// }
```

get decoded bodies:

```go
email.Text
// "The quick brown fox jumps over a lazy dog..."

email.HTML
// "<html><div dir="ltr"><p>The quick brown fox jumps over a lazy dog..."
```

Both inline and attached files are stored in a slice. By default these
are read into a `Data` []byte slice but direct access can be made to the
underlying `io.Reader`.

```go
Files: []*email.File{
	{
		FileType: "inline",
		ContentTypeHeader: email.ContentTypeHeader{
			ContentType: "image/jpeg",
			Params: map[string]string{
				"name": "inline-jpg-image-name.jpg",
			},
		},
		ContentDispositionHeader: email.ContentDispositionHeader{
			ContentDisposition: "inline",
			Params: map[string]string{
				"filename": "inline-jpg-image-filename.jpg",
			},
		},
		Name: "inline-jpg-image-filename.jpg",
		Data: []byte{
			255, 216, 255, 219, 0, 67, 0, 3, 2, 2, 2, 2, 2, 3, 2, 2, 2, 3, 3,
			3, 3, 4, 6, 4, 4, 4, 4, 4, 8, 6, 6, 5, 6, 9, 8, 10, 10, 9, 8, 9, 9,
			10, 12, 15, 12, 10, 11, 14, 11, 9, 9, 13, 17, 13, 14, 15, 16, 16,
			17, 16, 10, 12, 18, 19, 18, 16, 19, 15, 16, 16, 16, 255, 201, 0,
			11, 8, 0, 1, 0, 1, 1, 1, 17, 0, 255, 204, 0, 6, 0, 16, 16, 5, 255,
			218, 0, 8, 1, 1, 0, 0, 63, 0, 210, 207, 32, 255, 217,
		},
	},
	{
		FileType: "attached",
		Name:     "attached-pdf-filename.pdf",
		ContentTypeHeader: email.ContentTypeHeader{
			ContentType: "application/pdf",
			Params: map[string]string{
				"name": "attached-pdf-name.pdf",
			},
		},
		ContentDispositionHeader: email.ContentDispositionHeader{
			ContentDisposition: "attachment",
			Params: map[string]string{
				"filename": "attached-pdf-filename.pdf",
			},
		},
		Data: []byte{
			37, 80, 68, 70, 45, 49, 46, 13, 116, 114, 97, 105, 108, 101, 114,
			60, 60, 47, 82, 111, 111, 116, 60, 60, 47, 80, 97, 103, 101, 115,
			60, 60, 47, 75, 105, 100, 115, 91, 60, 60, 47, 77, 101, 100, 105,
			97, 66, 111, 120, 91, 48, 32, 48, 32, 51, 32, 51, 93, 62, 62, 93,
			62, 62, 62, 62, 62, 62,
		},
	},
}
```

## Options

Various options are provided for customising the Parser, including:

```go
func WithCustomAddressFunc(af func(string) (*mail.Address, error)) Opt
func WithCustomAddressesFunc(af func(list string) ([]*mail.Address, error)) Opt
func WithCustomDateFunc(df func(string) (time.Time, error)) Opt
func WithCustomFileFunc(ff func(*email.File) error) Opt
func WithSaveFilesToDirectory(dir string) Opt
func WithHeadersOnly() Opt
func WithoutAttachments() Opt
func WithVerbose() Opt
```

The `WithoutAttachments` and `WithHeadersOnly` options determine if only part
of an email will be processed.

The date and address "With" options allow the provision of custom funcs to
override the [net/mail] funcs normally used. For example it might be necessary
to extend the date parsing capabilities to deal with poorly formatted date
strings produced by older SMTP servers.

The `WithCustomFileFunc` allows the provision of a custom func for saving,
filtering and/or processing of inline or attached files without reading them
first into an `email.File.Data` []byte slice first, which is the default
behaviour. The `WithSaveFilesToDirectory` option is an example of such a custom
func.

As shown in the [parser/optspkg_test.go](parser/optspkg_test.go) package
test, `WithCustomFileFunc` can be used to, for example, only process
`image/jpeg` files. More examples are shown in
[parser/opts_test.go](parser/opts_test.go), for example:

```go
opt := parser.WithHeadersOnly() // the headers only option
p := letters.NewParser(opt, parser.WithVerbose()) // options can be chained
parsedEmail, err := p.Parse(rawEmail)
if err != nil {
	return fmt.Errorf("error while parsing email headers: %s", err)
}
```

## Language and Encoding Support

The same parser and methods will work for other languages, text encodings, 
and transfer-encodings:

```go
r := strings.NewReader(```Subject: =?ISO-2022-JP?Q?=1B=24=42=24=24=24=6D=24=4F=32=4E=1B=28=42?=
Content-Type: text/plain; charset=ISO-2022-JP


=1B$B?'$OFw$($I=1B(B
=1B$B;6$j$L$k$r=1B(B```)

email, _ := letters.ParseEmail(r)

email.Headers.Subject
// "いろは歌"

email.Text
// "色は匂えど散りぬるを..."
```

## Current Scope and Features

* Parsing plaintext emails and recursively traversing multipart
  (`multipart/alternative`, `multipart/mixed`, `multipart/parallel`,
  `multipart/related`, `multipart/signed`) emails
* Unfolding headers
* Decoding non-US-ASCII email headers according to
  [RFC 2047](https://datatracker.ietf.org/doc/html/rfc2047)
* Decoding Base64 and Quoted-Printable Content-Transfer-Encodings
* Decoding any text encoding (e.g. UTF-8, Chinese GB18030 or GBK, Finnish 
  ISO-8859-15, Icelandic ISO-8859-1, Japanese ISO-2022-JP, Korean EUC-KR,
  Polish ISO-8859-2) in combination with any Transfer
  Encoding (e.g. ASCII-over-7bit, UTF-8-over-Base64,
  Japanese ISO-2022-JP-over-7bit, Polish ISO-8859-2-over-Quoted-Printable,
  etc.)
* Easy access to text, enriched text and HTML content of the email
* Easy access to inline attachments
* Easy access to attached files

All of that and more in a minimal Golang library with realistic email
examples and thorough test coverage.

## Current Limitations

* S/MIME `multipart/signed` email are limited to clear-signed messages
* The decryption and signature verification and any other
  cryptography-related tasks need to be performed outside of letters.

## Current Status

Feature-complete and tests passing. 

Currently, gathering feedback and refactoring code before releasing v1.0.0.
Fields and API are still subject to change.
