# Letters, or how to parse emails in Go

[![Go Reference](https://img.shields.io/badge/Reference-blue?logo=Go&labelColor=black)](https://pkg.go.dev/github.com/mnako/letters)
[![Test](https://github.com/mnako/letters/actions/workflows/test.yml/badge.svg)](https://github.com/mnako/letters/actions/workflows/test.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/mnako/letters)](https://goreportcard.com/report/github.com/mnako/letters)

**Letters** is a minimalistic Golang library for parsing plaintext and MIME
emails.

It correctly handles text and MIME mime-types, Base64 and Quoted-Printable
Content-Transfer-Encoding, as well as any text encoding that Golang standard
library is capable of handling. Letters will parse an email into a simple struct
with standard headers and text, enriched text, and HTML content, and decode
inline and attached files.

## The User Guide

- [Installation](#installation)
- [Quickstart](#quickstart)
  - [Parse an Email](#parse-an-email)
  - [Parse Email Headers](#parse-email-headers)
- [Advanced Usage](#advanced-usage)
  - [The Email Parser](#the-email-parser)
  - [Skipping Parts of the Email](#skipping-parts-of-the-email)
  - [Customising Headers Parsers](#customising-headers-parsers)
  - [Customising Extra Headers Parsers](#customising-extra-headers-parsers)

### Installation

Install Letters:

```sh
go get github.com/mnako/letters@v0.2.8
```

### Quickstart

This section documents the simplified API of Letters.

#### Parse an Email

The quickest way to get started with parsing emails with Letters is to use the
`letters.ParseEmail()` helper function. It takes an `io.Reader` with the
contents of the email, parses it using the default parser, and returns an
`Email` struct or an error:

```go
package main

import (
  "log"
  "os"

  "github.com/mnako/letters"
)

func main() {
  rawEmail, err := os.Open("email.eml")
  if err != nil {
    log.Fatal("error while reading email from file: %w", err)
    return
  }

  defer func() {
    if err := rawEmail.Close(); err != nil {
      log.Fatal("error while closing rawEmail: %w", err)
      return
    }
  }()

  email, err := letters.ParseEmail(rawEmail)
  if err != nil {
    log.Fatal(err)
  }
}
```

Now, you can access the common headers:

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

and custom headers:

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

inline files:

```go
email.InlineFiles
// []InlineFile{
//    {
//        ContentType: ContentTypeHeader{
//            ContentType: "image/jpeg",
//            Params: map[string]string{
//                "name": "inline-jpg-image-without-disposition.jpg",
//            },
//        },
//        ContentDisposition: ContentDispositionHeader{
//            ContentDisposition: "",
//            Params:             map[string]string(nil),
//        },
//        Data: []byte{255, ...},
//    },
//    {
//        ContentID: "inline-jpg-image.jpg@example.com",
//        ContentType: ContentTypeHeader{
//            ContentType: "image/jpeg",
//            Params: map[string]string{
//                "name": "inline-jpg-image-name.jpg",
//            },
//        },
//        ContentDisposition: ContentDispositionHeader{
//            ContentDisposition: inline,
//            Params: map[string]string{
//                "filename": "inline-jpg-image-filename.jpg",
//            },
//        },
//        Data: []byte{255, ...},
//    },
// }
```

and attached files:

```go
email.AttachedFiles
// []AttachedFile{
//    {
//        ContentType: ContentTypeHeader{
//            ContentType: "application/pdf",
//            Params: map[string]string{
//                "name": "attached-pdf-name.pdf",
//            },
//        },
//        ContentDisposition: ContentDispositionHeader{
//            ContentDisposition: attachment,
//            Params: map[string]string{
//                "filename": "attached-pdf-filename.pdf",
//            },
//        },
//        Data: []byte{37, ...},
//    },
//    {
//        ContentType: ContentTypeHeader{
//            ContentType: "application/pdf",
//            Params: map[string]string{
//                "name": "attached-pdf-without-disposition.pdf",
//            },
//        },
//        ContentDisposition: ContentDispositionHeader{
//            ContentDisposition: "",
//            Params:             map[string]string(nil),
//        },
//        Data: []byte{37, ...},
// }
```

The same default parser and methods will work for other languages, text
encodings, and transfer-encodings:

````go
r := strings.NewReader(```Subject: =?ISO-2022-JP?Q?=1B=24=42=24=24=24=6D=24=4F=32=4E=1B=28=42?=
Content-Type: text/plain; charset=ISO-2022-JP


=1B$B?'$OFw$($I=1B(B
=1B$B;6$j$L$k$r=1B(B```)

email, _ := letters.ParseEmail(r)

email.Headers.Subject
// "いろは歌"

email.Text
// "色は匂えど散りぬるを..."
````

#### Parse Email Headers

If you already have a `mail.Message` and you are only interested in the headers
of the email, you can use the `letters.ParseEmailHeaders()` helper function. It
takes a `mail.Header` with the whole head section of an email, parses it using
the default parser, and returns a `letters.Headers` struct or an error:

```go
msg, err := mail.ReadMessage(rawEmail)
if err != nil {
    log.Fatal("error while reading message from file: %s", err)
    return
}

headers, err := letters.ParseEmailHeaders(msg.Header)
if err != nil {
    log.Fatal("error while parsing email headers: %s", err)
    return
}

headers.Sender
// mail.Address{Name: "Alice Sender", Address: "alice.sender@example.com"}

headers.Headers.From
// []mail.Address{
//  {Name: "Alice Sender", Address: "alice.sender@example.com"},
//  {Name: "Alice Sender", Address: "alice.sender@example.net"},
// }

// ...
```

> [!TIP] 
> The `letters.ParseEmail()` and `letters.ParseEmailHeaders()` helpers
> exist for developers’ convenience and are a good entrypoint to get started
> quickly. However, if you find yourself in need of customising the parser to,
> i.a. parse non-compliant headers or conditionally parse only some parts of the
> email, we recommend following the [Advanced Usage](#advanced-usage) section
> below for a more flexible and maintainable approach.

### Advanced Usage

This section documents come of Letters more advanced features.

#### The Email Parser

The `letters.ParseEmail()` documented above is a helper function that creates a
default email parser and returns the parsed email and error.

You can replace it with the following code to have more control over the parser:

```go
defaultEmailParser := letters.NewEmailParser()
email, err := defaultEmailParser.Parse(rawEmail)
if err != nil {
    log.Fatal(err)
}
```

#### Skipping Parts of the Email

By default, letters parses all bodies and files.

You can configure the parser to parse all, some, or none bodies, and files using
functional filters.

A **functional body filter** is a function that takes the Content-Type header of
the part and returns true or false. Only bodies for which the filter returns
true will be parsed. Parts for which the filter returned false, will be skipped.

Similar to the body filter, a **functional file filter** is a function that
takes the Content-Type and Content-Disposition headers of the part and returns
true or false. Only files for which the filter returns true will be parsed.
Files for which the filter returned false, will be skipped.

For example, if you do not want to parse any files, you can configure the Email
Parser with a file filter that always returns false. For convenience, letters
includes a `NoFiles` filter that does precisely that:

```go
noFilesEmailParser := letters.NewEmailParser(
    letters.WithFileFilter(NoFiles),
)
email, err := noFilesEmailParser.Parse(rawEmail)
if err != nil {
    log.Fatal(err)
}
```

Letters includes the following convenience filters:

- `NoBodies`, a function that always returns false, that can be used with
  `WithBodyFilter()`, to skip parsing all bodies of the email;
- `AllBodies`, a function that always returns true, that can be used with
  `WithBodyFilter()`, to parse all bodies of the email. This is the default
  behaviour;
- `NoFiles`, a function that always returns false, that can be used with
  `WithFileFilter()`, to skip parsing all attachments of the email. This option
  can speed up parsing in use cases where attachments are not needed;
- `AllFiles`, a function that always returns true, that can be used with
  `WithFileFilter()`, to parse all attachments of the email. This is the default
  behaviour;

This means that `WithBodyFilter(NoBodies)` is the more flexible equivalent of
the `ParseEmailHeaders()` helper function from the
[Quickstart](#parse-email-headers) section above, as it allows you to parse only
the headers of the email.

More interestingly, bodies and files can be skipped conditionally: bodies can be
skipped based on the Content-Type header of the part, and files can be skipped
based on the Content-Type and the Content-Disposition headers of the part.

For example, to only parse files with a filename that ends with `.jpg`, you can
pass a custom File Filter that checks the `name` Param of the Content-Type
header:

```go
customJPGOnlyEmailParser := letters.NewEmailParser(
    letters.WithFileFilter(
        func(
            cth letters.ContentTypeHeader,
            _ letters.ContentDispositionHeader,
        ) bool {
            return strings.HasSuffix(
                strings.ToLower(cth.Params["name"]),
                ".jpg",
            )
        },
    ),
)
email, err := customJPGOnlyEmailParser.Parse(rawEmail)
```

Files can be filtered based on the Content-Disposition header as well. For
example, to parse only inline files and skip attachments, you can pass a custom
File Filter that checks the Content-Disposition header:

```go
inlineFilesOnlyParser := letters.NewEmailParser(
    letters.WithFileFilter(
        func(
            _ letters.ContentTypeHeader,
            cdh letters.ContentDispositionHeader,
        ) bool {
            return cdh.ContentDisposition == letters.ContentDispositionInline
        },
    ),
)
```

You can implement arbitrarily complex conditions with those filter.

#### Customising Headers Parsers

Letters aims to be spec-compliant and follow email-related RFCs closely. Not all
emails that you will encounter in the wild will be compliant with the RFCs,
though.

To allow parsing such emails without adding exceptions to the library itself,
you can freely customise the parsers according to your needs with
`With<HeaderName>HeaderParser()` options.

For example, to customise the parser for the `Date` header, you can pass a
custom parser function to the `WithDateHeaderParser()` option:

```go
customDateHeaderEmailParser := letters.NewEmailParser(
    letters.WithDateHeaderParser(
        func(s string) time.Time {
            // Decode and parse the raw Date header from the s string here
            // and return a time.Time object
        },
    ),
)
```

Custom parsers for every header can be specified separately. E.g., if you want
to customise parsing of all dates, you need to specify a custom parser for the
Date header and for the ResentDate header.

The following headers are included in the `letters.Headers` struct and the
following options and types are available:

| Header              | Option                                                                | Expected Parser Signature                                |
| ------------------- | --------------------------------------------------------------------- | -------------------------------------------------------- |
| Date                | `WithDateHeaderParser(parseDateHeaderFn)`                             | `func(string) time.Time`                                 |
| Sender              | `WithSenderHeaderParser(parseAddressHeaderFn)`                        | `func(mail.Header, string) (*mail.Address, error)`       |
| From                | `WithFromHeaderParser(parseAddressListHeaderFn)`                      | `func(mail.Header, string) ([]*mail.Address, error)`     |
| Reply-To            | `WithReplyToHeaderParser(parseAddressListHeaderFn)`                   | `func(mail.Header, string) ([]*mail.Address, error)`     |
| To                  | `WithToHeaderParser(parseAddressListHeaderFn)`                        | `func(mail.Header, string) ([]*mail.Address, error)`     |
| Cc                  | `WithCcHeaderParser(parseAddressListHeaderFn)`                        | `func(mail.Header, string) ([]*mail.Address, error)`     |
| Bcc                 | `WithBccHeaderParser(parseAddressListHeaderFn)`                       | `func(mail.Header, string) ([]*mail.Address, error)`     |
| Message-ID          | `WithMessageIdHeaderParser(parseMessageIdHeaderFn)`                   | `func(string) letters.MessageId`                         |
| In-Reply-To         | `WithInReplyToHeaderParser(parseCommaSeparatedMessageIdHeaderFn)`     | `func(string) []letters.MessageId`                       |
| References          | `WithReferencesHeaderParser(parseCommaSeparatedMessageIdHeaderFn)`    | `func(string) []letters.MessageId`                       |
| Subject             | `WithSubjectHeaderParser(parseStringHeaderFn)`                        | `func(string) string`                                    |
| Comments            | `WithCommentsHeaderParser(parseStringHeaderFn)`                       | `func(string) string`                                    |
| Keywords            | `WithKeywordsHeaderParser(parseCommaSeparatedStringHeaderFn)`         | `func(string) []string`                                  |
| Resent-Date         | `WithResentDateHeaderParser(parseDateHeaderFn)`                       | `func(string) time.Time`                                 |
| Resent-From         | `WithResentFromHeaderParser(parseAddressListHeaderFn)`                | `func(mail.Header, string) ([]*mail.Address, error)`     |
| Resent-Sender       | `WithResentSenderHeaderParser(parseAddressHeaderFn)`                  | `func(mail.Header, string) (*mail.Address, error)`       |
| Resent-To           | `WithResentToHeaderParser(parseAddressListHeaderFn)`                  | `func(mail.Header, string) ([]*mail.Address, error)`     |
| Resent-Cc           | `WithResentCcHeaderParser(parseAddressListHeaderFn)`                  | `func(mail.Header, string) ([]*mail.Address, error)`     |
| Resent-Bcc          | `WithResentBccHeaderParser(parseAddressListHeaderFn)`                 | `func(mail.Header, string) ([]*mail.Address, error)`     |
| Resent-Message-ID   | `WithResentMessageIdHeaderParser(parseMessageIdHeaderFn)`             | `func(string) letters.MessageId`                         |
| Content-Type        | `WithContentTypeHeaderParser(parseContentTypeHeaderFn)`               | `func(string) (letters.ContentTypeHeader, error)`        |
| Content-Disposition | `WithContentDispositionHeaderParser(parseContentDispositionHeaderFn)` | `func(string) (letters.ContentDispositionHeader, error)` |

#### Customising Extra Headers Parsers

The RFC-documented headers, their corresponding options and parser signatures
are strictly typed.

Letters allows parsing custom headers to the `ExtraHeaders` map in the
`letters.Headers` struct. You can freely customise the parsers according to your
needs with the
`WithExtraHeaderParser(headerName string, extraHeaderParserFn parseStringHeaderFn)`
option.

For example, you can customise parsing the `X-CrossPremisesHeadersPromoted`
header:

```go
xCrossPremisesHeadersPromotedHeaderParser := func(s string) string {
    // ...
    return xCrossPremisesHeadersPromotedHeader
}

customEmailParser := letters.NewEmailParser(
    letters.WithExtraHeaderParser(
        "X-CrossPremisesHeadersPromoted",
        xCrossPremisesHeadersPromotedHeaderParser,
      ),
)
```

You can also specify custom parsers for many extra headers:

```go
xCrossPremisesHeadersPromotedHeaderParser := func(s string) string {
    // ...
    return xCrossPremisesHeadersPromotedHeader
}

xMSExchangeOrganizationProcessedByJournalingHeaderParser := func(s string) string {
    // ...
    return xMSExchangeOrganizationProcessedByJournalingHeader
}

xMSExchangeOrganizationOriginalEnvelopeRecipientsHeaderParser := func(s string) string {
    // ...
    return xMSExchangeOrganizationOriginalEnvelopeRecipientsHeader
}

customEmailParser := letters.NewEmailParser(
    letters.WithExtraHeaderParser(
        "X-CrossPremisesHeadersPromoted",
        xCrossPremisesHeadersPromotedHeaderParser,
    ),
    letters.WithExtraHeaderParser(
        "X-MS-Exchange-Organization-Processed-By-Journaling",
        xMSExchangeOrganizationProcessedByJournalingHeaderParser,
    ),
    letters.WithExtraHeaderParser(
        "X-MS-Exchange-Organization-OriginalEnvelopeRecipients",
        xMSExchangeOrganizationOriginalEnvelopeRecipientsHeaderParser,
    ),
)
```

## Current Scope and Features

- Parsing plaintext emails and recursively traversing multipart
  (`multipart/alternative`, `multipart/mixed`, `multipart/parallel`,
  `multipart/related`, `multipart/signed`) emails
- Unfolding headers
- Decoding non-US-ASCII email headers according to
  [RFC 2047](https://datatracker.ietf.org/doc/html/rfc2047)
- Decoding Base64 and Quoted-Printable Content-Transfer-Encodings
- Decoding any text encoding (e.g. UTF-8, Chinese GB18030 or GBK, Finnish
  ISO-8859-15, Icelandic ISO-8859-1, Japanese ISO-2022-JP, Korean EUC-KR, Polish
  ISO-8859-2) in combination with any Transfer Encoding (e.g. ASCII-over-7bit,
  UTF-8-over-Base64, Japanese ISO-2022-JP-over-7bit, Polish
  ISO-8859-2-over-Quoted-Printable, etc.)
- Easy access to text, enriched text and HTML content of the email
- Easy access to inline attachments
- Easy access to attached files

All of that and more in a minimal Golang library with realistic email examples
and thorough test coverage.

## Current Limitations

- S/MIME `multipart/signed` email are limited to clear-signed messages
- The decryption and signature verification and any other cryptography-related
  tasks need to be performed outside of letters.

## Current Status

Feature-complete and tests passing.

Currently, gathering feedback and refactoring code before releasing v1.0.0.
Fields and API are still subject to change.

# Release Policy

We follow [Go’s Release Policy](https://go.dev/doc/devel/release#policy) 
and commit to supporting at least the two most recent major versions of Go.

Letter v0.2.8 supports Go versions 1.24 through 1.25.
