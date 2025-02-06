package parser

import (
	"fmt"
	"strings"
	"testing"

	"github.com/mnako/letters/email"
)

func TestParseFile(t *testing.T) {

	tests := []struct {
		content     string
		contentInfo *email.ContentInfo
		contentLen  int
		fileName    string
	}{
		{
			content:    "hi",
			contentLen: 2,
			fileName:   "test_file.jpg",
			contentInfo: &email.ContentInfo{
				Disposition:       "inline",
				DispositionParams: map[string]string{"filename": "test_file.jpg"},
				TransferEncoding:  "8bit",
			},
		},
		{
			content: `
/9j/2wBDAAMCAgICAgMCAgIDAwMDBAYEBAQEBAgGBgUGCQgKCgkICQkKDA8MCgsOCwkJDRENDg8Q
EBEQCgwSExIQEw8QEBD/yQALCAABAAEBAREA/8wABgAQEAX/2gAIAQEAAD8A0s8g/9k=`,
			contentLen: 107,
			fileName:   "inline-jpg-image-without-disposition.jpg",
			contentInfo: &email.ContentInfo{
				Disposition:       "",
				DispositionParams: map[string]string{"filename": "inline-jpg-image-without-disposition.jpg"},
				TransferEncoding:  "base64",
			},
		},
		{
			content: `
/9j/2wBDAAMCAgICAgMCAgIDAwMDBAYEBAQEBAgGBgUGCQgKCgkICQkKDA8MCgsOCwkJDRENDg8Q
EBEQCgwSExIQEw8QEBD/yQALCAABAAEBAREA/8wABgAQEAX/2gAIAQEAAD8A0s8g/9k=`,
			contentLen: 107,
			fileName:   "attachment_0_inline",
			contentInfo: &email.ContentInfo{
				Disposition:      "inline",
				TransferEncoding: "base64",
			},
		},
		{
			content: `
VGV4dC9wbGFpbiBjb250ZW50IGFzIGFuIGF0dGFjaGVkIC50eHQgZmlsZS4=`,
			contentLen: 44,
			fileName:   "attached-text-plain-filename.txt",
			contentInfo: &email.ContentInfo{
				Type:              "text/plain",
				Disposition:       "attachment",
				DispositionParams: map[string]string{"filename": "attached-text-plain-filename.txt"},
				TransferEncoding:  "base64",
			},
		},
		{
			content:    `VGV4dC9wbGFpbiBjb250ZW50IGFzIGFuIGF0dGFjaGVkIC50eHQgZmlsZS4=`,
			contentLen: 44,
			fileName:   "sshd_config",
			contentInfo: &email.ContentInfo{
				Type:              "text/plain",
				Disposition:       "attachment",
				DispositionParams: map[string]string{"filename": "/etc/ssh/sshd_config"},
				TransferEncoding:  "base64",
			},
		},
	}

	for i, tt := range tests {
		t.Run(fmt.Sprintf("test_%d", i), func(t *testing.T) {
			p := NewParser()
			err := p.parseFile(
				strings.NewReader(tt.content),
				tt.contentInfo,
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
