# Letters, or how to parse emails in Go

[![Test](https://github.com/mnako/letters/actions/workflows/test.yml/badge.svg)](https://github.com/mnako/letters/actions/workflows/test.yml)
[![Lint](https://github.com/mnako/letters/actions/workflows/lint.yml/badge.svg)](https://github.com/mnako/letters/actions/workflows/lint.yml)

**Letters** is a minimalistic Golang library for parsing plaintext and MIME
emails.

It correctly handles text and MIME mime-types, Base64 and Quoted-Printable 
Content-Transfer-Encoding, as well as any text encoding that Golang 
standard library is capable of handling. Letters will parse an email into 
a simple struct with standard headers and text, enriched text, and HTML 
content, and decode inline and attached files.

## Quickstart

Install

```
go get github.com/mnako/letters
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

The same parser and methods will work for other languages, text encodings, 
and transfer encodings:

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
* Decoding Base64 and Quoted-Printable Content Transfer Encodings
* Decoding any text encoding (e.g. UTF-8, Japanese ISO-2022-JP,
  Polish ISO-8859-2, Finnish ISO-8859-15) in combination with any Transfer
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

Currently gethering feedback and refactoring code before releasing v1.0.0.
Fields and API are still subject to change.
