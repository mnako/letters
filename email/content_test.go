package email

import (
	"fmt"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestExtractContentTypeHeader(t *testing.T) {
	tests := []struct {
		input       string
		contentType string
		params      map[string]string
	}{
		{
			input:       `MULtiPARt/mIXed; Charset="ascII"; bouNDARY="MixedBoundaryString"`,
			contentType: "multipart/mixed",
			params: map[string]string{
				"boundary": "MixedBoundaryString",
				"charset":  "ascii",
			},
		},
		{
			input:       `TEXT/HtmL; CharSEt="ascii"`,
			contentType: "text/html",
			params: map[string]string{
				"charset": "ascii",
			},
		},
		{
			input:       `image/jpeg; nAME="inline-jpg-image-without-disposition.jpg"`,
			contentType: "image/jpeg",
			params: map[string]string{
				"name": "inline-jpg-image-without-disposition.jpg",
			},
		},
		{
			input: `MUlTIpart/signed;
              cHarSET="iso-8859-2";
              proToCOL="APPLICATIOn/pkcs7-SiGnAture";
              micalg=sha1;
              boUNDARY=SignedBoundaryString`,
			contentType: "multipart/signed",
			params: map[string]string{
				"charset":  "iso-8859-2",
				"micalg":   "sha1",
				"protocol": "application/pkcs7-signature",
				"boundary": "SignedBoundaryString",
			},
		},
	}

	for i, tt := range tests {
		t.Run(fmt.Sprintf("test_%d", i), func(t *testing.T) {
			c := &ContentInfo{}

			err := c.extractType(tt.input)
			if err != nil {
				t.Fatalf("cannot parse part Content-Type: %s", err)
			}

			if got, want := c.Type, tt.contentType; got != want {
				t.Errorf("got %s want %s", got, want)
			}
			got, want := c.TypeParams, tt.params
			if diff := cmp.Diff(want, got); diff != "" {
				t.Errorf("params are not equal\n%s", diff)
			}
		})
	}
}

func TestExtractContentDisposition(t *testing.T) {
	tests := []struct {
		input              string
		contentDisposition string
		params             map[string]string
	}{
		{
			input:              `ATTACHment; filename=smime.p7s`,
			contentDisposition: "attachment",
			params: map[string]string{
				"filename": "smime.p7s",
			},
		},
		{
			input:              `INLinE; Filename="inline-jpg-image-filename.jpg"`,
			contentDisposition: "inline",
			params: map[string]string{
				"filename": "inline-jpg-image-filename.jpg",
			},
		},
	}
	for i, tt := range tests {
		t.Run(fmt.Sprintf("test_%d", i), func(t *testing.T) {
			c := &ContentInfo{}
			err := c.extractDisposition(tt.input)
			if err != nil {
				t.Fatalf("cannot parse part Content-Disposition: %s", err)
			}

			if got, want := c.Disposition, tt.contentDisposition; got != want {
				t.Errorf("got %s want %s", got, want)
			}
			got, want := c.DispositionParams, tt.params
			if diff := cmp.Diff(want, got); diff != "" {
				t.Errorf("params are not equal\n%s", diff)
			}
		})
	}
}

func TestExtractContentTransferEncoding(t *testing.T) {
	tests := []struct {
		input string
		cte   string
	}{
		{
			input: ``, // empty
			cte:   "7bit",
		},
		{
			input: `base64`,
			cte:   "base64",
		},
		{
			input: `QUOTeD-PriNTABLE `,
			cte:   "quoted-printable",
		},
	}
	for i, tt := range tests {
		t.Run(fmt.Sprintf("test_%d", i), func(t *testing.T) {
			c := &ContentInfo{}
			err := c.extractTransferEncoding(tt.input)
			if err != nil {
				t.Fatalf("cannot parse part Content-Transfer-Encoding: %s", err)
			}

			if got, want := c.TransferEncoding, tt.cte; got != want {
				t.Errorf("got %s want %s", got, want)
			}
		})
	}
}

func TestExtractCharset(t *testing.T) {
	tests := []struct {
		input       string
		parentCI    *ContentInfo
		charset     string
		hasEncoding bool
	}{
		{
			input:       "xyz",
			parentCI:    nil,
			charset:     "xyz",
			hasEncoding: false,
		},
		{
			input: "",
			parentCI: &ContentInfo{
				TypeParams: map[string]string{"charset": "UTF-8"},
			},
			charset:     "UTF-8",
			hasEncoding: true,
		},
	}
	for i, tt := range tests {
		t.Run(fmt.Sprintf("test_%d", i), func(t *testing.T) {
			c := &ContentInfo{}
			c.TypeParams = map[string]string{"charset": tt.input}
			c.extractCharset(tt.parentCI)
			if got, want := c.Charset, tt.charset; got != want {
				t.Errorf("charset got %s want %s", got, want)
			}
			c.ExtractEncoding()
			if got, want := !(c.Encoding == nil), tt.hasEncoding; got != want {
				t.Errorf("encoding got %t want %t", got, want)
			}
		})
	}
}

func TestIsInlineFile(t *testing.T) {
	tests := []struct {
		ci       *ContentInfo
		parentCI *ContentInfo
		expected bool
	}{
		{
			ci:       &ContentInfo{Disposition: "attachment"},
			parentCI: &ContentInfo{Type: "multipart/mixed"},
			expected: false,
		},
		{
			ci:       &ContentInfo{Disposition: "inline"},
			parentCI: &ContentInfo{Type: "multipart/mixed"},
			expected: true,
		},
		{
			ci:       &ContentInfo{},
			parentCI: &ContentInfo{Type: "multipart/related"},
			expected: true,
		},
	}
	for i, tt := range tests {
		t.Run(fmt.Sprintf("test_%d", i), func(t *testing.T) {
			if got, want := tt.ci.IsInlineFile(tt.parentCI), tt.expected; got != want {
				t.Errorf("got %t want %t", got, want)
			}
		})
	}
}

func TestIsAttachedFile(t *testing.T) {
	tests := []struct {
		ci       *ContentInfo
		parentCI *ContentInfo
		expected bool
	}{
		{
			ci:       &ContentInfo{Disposition: "inline"},
			parentCI: nil,
			expected: false,
		},
		{
			ci:       &ContentInfo{Disposition: "attached"},
			parentCI: nil,
			expected: true,
		},
		{
			ci:       &ContentInfo{},
			parentCI: &ContentInfo{Type: "multipart/mixed"},
			expected: true,
		},
		{
			ci:       &ContentInfo{},
			parentCI: &ContentInfo{Type: "multipart/parallel"},
			expected: true,
		},
	}

	for i, tt := range tests {
		t.Run(fmt.Sprintf("test_%d", i), func(t *testing.T) {
			if got, want := tt.ci.IsAttachedFile(tt.parentCI), tt.expected; got != want {
				t.Errorf("got %t want %t", got, want)
			}
		})
	}
}
