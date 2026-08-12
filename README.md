# Letters, or how to parse emails in Go

[![Go Reference](https://img.shields.io/badge/Reference-blue?logo=Go&labelColor=black)](https://pkg.go.dev/github.com/mnako/letters)
[![Test](https://github.com/mnako/letters/actions/workflows/test.yml/badge.svg)](https://github.com/mnako/letters/actions/workflows/test.yml)

**Letters** is a minimal Go library that parses plain-text and MIME emails.

Letters parses plain-text, enriched-text, and HTML bodies. It also parses inline
and attached files. Letters decodes Base64 and Quoted-Printable content-transfer
encodings. It decodes all character encodings supported by
`golang.org/x/net/html/charset`.

The parser returns an `Email` struct with standard email headers; plain-text,
enriched-text, and HTML bodies; and inline and attached files.

## The User Guide

- [Installation](#installation)
- [Quick Start](#quick-start)
  - [Parse an Email](#parse-an-email)
  - [Parse Email Headers](#parse-email-headers)
- [Parser Options](#parser-options)
  - [The Email Parser](#the-email-parser)
  - [Skip Email Parts](#skip-email-parts)
  - [Customize Header Parsers](#customize-header-parsers)
  - [Customize Parsers for Extra Headers](#customize-parsers-for-extra-headers)

### Installation

Install Letters with this command:

```sh
go get github.com/mnako/letters@v0.4.0
```

### Quick Start

Parse emails with the default parser.

#### Parse an Email

Use `letters.ParseEmail()` to parse an email with the default parser. The
function accepts an `io.Reader` that contains the email and returns an `Email`
struct or an error:

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
    log.Fatal("error while reading email from file: ", err)
    return
  }

  defer func() {
    if err := rawEmail.Close(); err != nil {
      log.Fatal("error while closing rawEmail: ", err)
      return
    }
  }()

  email, err := letters.ParseEmail(rawEmail)
  if err != nil {
    log.Fatal(err)
  }
}
```

After Letters parses the email, you can get the RFC headers:

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

Get extra headers from `email.Headers.ExtraHeaders`:

```go
email.Headers.ExtraHeaders
// map[string][]string{
//    "X-Clacks-Overhead": {"GNU Terry Pratchett"},
// }
```

Get the decoded bodies from these fields:

```go
email.Text
// "The quick brown fox jumps over a lazy dog..."

email.HTML
// "<html><div dir="ltr"><p>The quick brown fox jumps over a lazy dog..."
```

Get inline files from `email.InlineFiles`:

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

Get attached files from `email.AttachedFiles`:

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

The default parser also parses messages written in other languages when they use
supported character encodings and content-transfer encodings:

```go
r := strings.NewReader(`Subject: =?ISO-2022-JP?Q?=1B=24=42=24=24=24=6D=24=4F=32=4E=1B=28=42?=
Content-Type: text/plain; charset=ISO-2022-JP


=1B$B?'$OFw$($I=1B(B
=1B$B;6$j$L$k$r=1B(B`)

email, _ := letters.ParseEmail(r)

email.Headers.Subject
// "いろは歌"

email.Text
// "色は匂えど散りぬるを..."
```

#### Parse Email Headers

Use `letters.ParseEmailHeaders()` to parse only email headers. The function
accepts a `mail.Header` and returns a `letters.Headers` struct or an error:

```go
msg, err := mail.ReadMessage(rawEmail)
if err != nil {
    log.Fatal("error while reading message from file: ", err)
    return
}

headers, err := letters.ParseEmailHeaders(msg.Header)
if err != nil {
    log.Fatal("error while parsing email headers: ", err)
    return
}

headers.Sender
// mail.Address{Name: "Alice Sender", Address: "alice.sender@example.com"}

headers.From
// []mail.Address{
//  {Name: "Alice Sender", Address: "alice.sender@example.com"},
//  {Name: "Alice Sender", Address: "alice.sender@example.net"},
// }

// ...
```

> [!TIP]
> `letters.ParseEmail()` and `letters.ParseEmailHeaders()` provide
> convenient ways to get started. To handle non-compliant headers or select
> which parts of an email to parse, configure an `EmailParser` as described in
> [Parser Options](#parser-options).

### Parser Options

Use `EmailParser` options to configure parsing.

#### The Email Parser

`letters.ParseEmail()` creates a parser with the default options and returns the
parsed email or an error.

Create an `EmailParser` directly when you need more control:

```go
defaultEmailParser := letters.NewEmailParser()
email, err := defaultEmailParser.Parse(rawEmail)
if err != nil {
    log.Fatal(err)
}
```

#### Skip Email Parts

By default, Letters parses all bodies and files.

A **body filter** controls the bodies that the parser parses or skips. A **file
filter** controls the files that the parser parses or skips.

A **body filter** receives the parsed Content-Type header for a body. The parser
parses the body when the filter returns `true` and skips it when the filter
returns `false`.

A **file filter** receives the parsed Content-Type and Content-Disposition
headers for a file. The parser parses the file when the filter returns `true`
and skips it when the filter returns `false`.

The `NoFiles` filter always returns false.

Skip all files with this filter:

```go
noFilesEmailParser := letters.NewEmailParser(
    letters.WithFileFilter(letters.NoFiles),
)
email, err := noFilesEmailParser.Parse(rawEmail)
if err != nil {
    log.Fatal(err)
}
```

Letters has these filters:

- `NoBodies` returns false for each body. `WithBodyFilter(NoBodies)` skips all
  bodies.
- `AllBodies` returns true for each body. It is the default body filter.
- `NoFiles` returns false for each file. `WithFileFilter(NoFiles)` skips all
  files.
- `AllFiles` returns true for each file. It is the default file filter.

`WithBodyFilter(NoBodies)` configures the email parser to parse only headers.
You can use other parser options with this configuration.

##### Examples

Parse only `.jpg` files with a file filter that checks the `name` parameter in
the Content-Type header:

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

Parse only inline files with a file filter that checks the Content-Disposition
header:

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

These filters can implement any conditions you need.

#### Customize Header Parsers

Letters closely follows email RFCs. Some real-world emails do not comply with
these RFCs. Use custom header parsers to handle non-compliant headers.

Configure a custom header parser with the corresponding
`With<HeaderName>HeaderParser()` option.

Customize the `Date` header parser with the `WithDateHeaderParser()` option:

```go
customDateHeaderEmailParser := letters.NewEmailParser(
    letters.WithDateHeaderParser(
        func(s string) time.Time {
            // Decode the raw Date header in s.
            // Parse the decoded header.
            // Return a time.Time value.
        },
    ),
)
```

The `letters.Headers` struct contains these headers. The table also gives the
options and parser signatures for these headers:

| Header              | Option                                                                | Parser Signature                                         |
|---------------------|-----------------------------------------------------------------------|----------------------------------------------------------|
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

#### Customize Parsers for Extra Headers

Letters provides dedicated fields and parser options for the RFC-defined headers
listed above.

It stores other headers in the `ExtraHeaders` map of the `letters.Headers`
struct.

Set an extra header parser with the
`WithExtraHeaderParser(headerName string, extraHeaderParserFn parseStringHeaderFn)`
option.

Customize an extra header parser with this code:

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

Set parsers for more than one extra header with this code:

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

## What Letters Does

- Letters parses plain-text emails.
- Letters parses `multipart/alternative`, `multipart/mixed`,
  `multipart/parallel`, `multipart/related`, and `multipart/signed` emails.
- Letters unfolds headers.
- Letters decodes non-ASCII email headers according to
  [RFC 2047](https://datatracker.ietf.org/doc/html/rfc2047).
- Letters supports the 7bit, 8bit, binary, Base64, and Quoted-Printable
  content-transfer encodings.
- Letters supports the character encodings provided by
  `golang.org/x/net/html/charset`. Examples include UTF-8, GB18030, GBK,
  ISO-8859-15, ISO-8859-1, ISO-2022-JP, EUC-KR, and ISO-8859-2.

The repository contains email examples and tests.

## Limits

- Letters parses only clear-signed S/MIME `multipart/signed` messages.
- Decryption and signature verification are outside the scope of Letters.

## Project State

Letters is feature-complete, and all tests pass.

We are gathering feedback and refactoring the code before releasing v1.0.0. The
public API, including struct fields, may still change.

## Release Policy

Letters follows [Go’s Release Policy](https://go.dev/doc/devel/release#policy).
Each Letters release supports at least the two most recent major Go releases.

Letters v0.4.0 supports Go 1.25 and Go 1.26.
