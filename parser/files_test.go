package parser

import (
	"fmt"
	"strings"
	"testing"

	"github.com/mnako/letters/email"
)

/*
	ContentTypeMultipartPrefix   ContentType = "multipart/"
	ContentTypeMultipartMixed                = "multipart/mixed"
	ContentTypeMultipartParallel             = "multipart/parallel"
	ContentTypeMultipartRelated              = "multipart/related"
	ContentTypeTextPlain                     = "text/plain"
	ContentTypeTextEnriched                  = "text/enriched"
	ContentTypeTextHtml                      = "text/html"

type ContentDisposition string

const (
	Attachment ContentDisposition = "attachment"
	Inline                        = "inline"
	TextPlain                     = "text/plain" // this is not in the rfc but used by some mail user agents
)

type ContentTransferEncoding string

const (
	CTE7bit            ContentTransferEncoding = "7bit"
	CTE8bit                                    = "8bit"
	CTEBinary                                  = "binary"
	CTEQuotedPrintable                         = "quoted-printable"
	CTEBase64                                  = "base64"
	// note that ietf-<token> and x-<token> mechanisms may also be
	// encountered.
)

// FileType describes the file type of an inline or attached file
type FileType string

const (
	InlineFileType   FileType = "inline"
	AttachedFileType FileType = "attached"
)

*/

func TestIsInlineFile(t *testing.T) {

	tests := []struct {
		ct       email.ContentTypeHeader
		parentCT email.ContentTypeHeader
		cdh      email.ContentDispositionHeader
		expected bool
	}{
		{
			ct: email.ContentTypeHeader{
				// ContentType: string(email.ContentTypeMultipartRelated),
			},
			parentCT: email.ContentTypeHeader{
				ContentType: string(email.ContentTypeMultipartMixed),
			},
			cdh: email.ContentDispositionHeader{
				ContentDisposition: email.ContentDisposition("attachment"),
			},
			expected: false,
		},
		{
			ct: email.ContentTypeHeader{
				// ContentType: string(email.ContentTypeMultipartRelated),
			},
			parentCT: email.ContentTypeHeader{
				ContentType: string(email.ContentTypeMultipartMixed),
			},
			cdh: email.ContentDispositionHeader{
				ContentDisposition: email.ContentDisposition("inline"),
			},
			expected: true,
		},
		{
			ct: email.ContentTypeHeader{
				// ContentType: string(email.ContentTypeMultipartRelated),
			},
			parentCT: email.ContentTypeHeader{
				ContentType: string(email.ContentTypeMultipartRelated),
			},
			cdh: email.ContentDispositionHeader{
				// ContentDisposition: email.ContentDisposition("inline"),
			},
			expected: true,
		},
	}

	for i, tt := range tests {
		t.Run(fmt.Sprintf("test_%d", i), func(t *testing.T) {
			if got, want := isInlineFile(tt.ct, tt.parentCT, tt.cdh), tt.expected; got != want {
				t.Errorf("got %t want %t", got, want)
			}
		})
	}
}

func TestIsAttachedFile(t *testing.T) {

	tests := []struct {
		ct       email.ContentTypeHeader
		parentCT email.ContentTypeHeader
		cdh      email.ContentDispositionHeader
		expected bool
	}{
		{
			ct: email.ContentTypeHeader{
				// ContentType: string(email.ContentTypeMultipartRelated),
			},
			parentCT: email.ContentTypeHeader{
				// ContentType: string(email.ContentTypeMultipartMixed),
			},
			cdh: email.ContentDispositionHeader{
				ContentDisposition: email.ContentDisposition("inline"),
			},
			expected: false,
		},
		{
			ct: email.ContentTypeHeader{
				// ContentType: string(email.ContentTypeMultipartRelated),
			},
			parentCT: email.ContentTypeHeader{
				// ContentType: string(email.ContentTypeMultipartMixed),
			},
			cdh: email.ContentDispositionHeader{
				ContentDisposition: email.ContentDisposition("attached"),
			},
			expected: true,
		},
		{
			ct: email.ContentTypeHeader{
				// ContentType: string(email.ContentTypeMultipartRelated),
			},
			parentCT: email.ContentTypeHeader{
				ContentType: string(email.ContentTypeMultipartMixed),
			},
			cdh: email.ContentDispositionHeader{
				// ContentDisposition: email.ContentDisposition("inline"),
			},
			expected: true,
		},
		{
			ct: email.ContentTypeHeader{
				// ContentType: string(email.ContentTypeMultipartRelated),
			},
			parentCT: email.ContentTypeHeader{
				ContentType: string(email.ContentTypeMultipartParallel),
			},
			cdh: email.ContentDispositionHeader{
				// ContentDisposition: email.ContentDisposition("inline"),
			},
			expected: true,
		},
	}

	for i, tt := range tests {
		t.Run(fmt.Sprintf("test_%d", i), func(t *testing.T) {
			if got, want := isAttachedFile(tt.ct, tt.parentCT, tt.cdh), tt.expected; got != want {
				t.Errorf("got %t want %t", got, want)
			}
		})
	}
}

func TestParseFile(t *testing.T) {

	tests := []struct {
		content    string
		contentLen int
		fileType   email.FileType
		fileName   string
		cth        email.ContentTypeHeader
		cdh        email.ContentDispositionHeader
		cte        email.ContentTransferEncoding
	}{
		{
			content:    "hi",
			contentLen: 2,
			fileType:   email.InlineFileType,
			fileName:   "test_file.jpg",
			cth:        email.ContentTypeHeader{
				// ContentType: string(email.ContentTypeMultipartRelated),
			},
			cdh: email.ContentDispositionHeader{
				ContentDisposition: email.ContentDisposition("inline"),
				Params:             map[string]string{"filename": "test_file.jpg"},
			},
			cte: email.CTE8bit,
		},
		{
			content: `
/9j/2wBDAAMCAgICAgMCAgIDAwMDBAYEBAQEBAgGBgUGCQgKCgkICQkKDA8MCgsOCwkJDRENDg8Q
EBEQCgwSExIQEw8QEBD/yQALCAABAAEBAREA/8wABgAQEAX/2gAIAQEAAD8A0s8g/9k=
`,
			contentLen: 107,
			fileType:   email.InlineFileType,
			fileName:   "inline-jpg-image-without-disposition.jpg",
			cth:        email.ContentTypeHeader{
				// ContentType: string(email.ContentTypeMultipartRelated),
			},
			cdh: email.ContentDispositionHeader{
				ContentDisposition: "",
				Params:             map[string]string{"filename": "inline-jpg-image-without-disposition.jpg"},
			},
			cte: email.CTEBase64,
		},
		{
			content: `
/9j/2wBDAAMCAgICAgMCAgIDAwMDBAYEBAQEBAgGBgUGCQgKCgkICQkKDA8MCgsOCwkJDRENDg8Q
EBEQCgwSExIQEw8QEBD/yQALCAABAAEBAREA/8wABgAQEAX/2gAIAQEAAD8A0s8g/9k=
`,
			contentLen: 107,
			fileType:   email.InlineFileType,
			fileName:   "attachment_0_inline",
			cth:        email.ContentTypeHeader{
				// ContentType: string(email.ContentTypeMultipartRelated),
			},
			cdh: email.ContentDispositionHeader{
				//
			},
			cte: email.CTEBase64,
		},
		{
			content: `
VGV4dC9wbGFpbiBjb250ZW50IGFzIGFuIGF0dGFjaGVkIC50eHQgZmlsZS4=
`,
			contentLen: 44,
			fileType:   email.AttachedFileType,
			fileName:   "attached-text-plain-filename.txt",
			cth: email.ContentTypeHeader{
				ContentType: string(email.ContentTypeTextPlain),
			},
			cdh: email.ContentDispositionHeader{
				ContentDisposition: "attachment",
				Params:             map[string]string{"filename": "attached-text-plain-filename.txt"},
			},
			cte: email.CTEBase64,
		},
		{
			content: `
VGV4dC9wbGFpbiBjb250ZW50IGFzIGFuIGF0dGFjaGVkIC50eHQgZmlsZS4=
`,
			contentLen: 44,
			fileType:   email.AttachedFileType,
			fileName:   "sshd_config",
			cth: email.ContentTypeHeader{
				ContentType: string(email.ContentTypeTextPlain),
			},
			cdh: email.ContentDispositionHeader{
				ContentDisposition: "attachment",
				Params:             map[string]string{"filename": "/etc/ssh/sshd_config"},
			},
			cte: email.CTEBase64,
		},
	}

	for i, tt := range tests {
		t.Run(fmt.Sprintf("test_%d", i), func(t *testing.T) {
			p := NewParser()
			err := p.parseFile(
				strings.NewReader(tt.content),
				tt.fileType,
				tt.cth,
				tt.cdh,
				tt.cte,
			)
			if err != nil {
				t.Fatal(err)
			}
			if got, want := len(p.email.Files), 1; got != want {
				t.Errorf("got %d want %d files", got, want)
			}
			filer := p.email.Files[0]
			if got, want := filer.Name, tt.fileName; got != want {
				t.Errorf("got name %s want %s ", got, want)
			}
			if got, want := len(filer.Data), tt.contentLen; got != want {
				t.Errorf("got data len %d want %d ", got, want)
			}
		})
	}

}
