package parser

import (
	"fmt"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/mnako/letters/email"
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
			cth, err := extractContentTypeHeader(tt.input)
			if err != nil {
				t.Fatalf("cannot parse part Content-Type: %s", err)
			}

			if got, want := cth.ContentType, tt.contentType; got != want {
				t.Errorf("got %s want %s", got, want)
			}
			got, want := cth.Params, tt.params
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
			cd, err := extractContentDisposition(tt.input)
			if err != nil {
				t.Fatalf("cannot parse part Content-Disposition: %s", err)
			}

			if got, want := string(cd.ContentDisposition), tt.contentDisposition; got != want {
				t.Errorf("got %s want %s", got, want)
			}
			got, want := cd.Params, tt.params
			if diff := cmp.Diff(want, got); diff != "" {
				t.Errorf("params are not equal\n%s", diff)
			}
		})
	}
}

func TestExtractContentTransferEncoding(t *testing.T) {

	tests := []struct {
		input string
		cte   email.ContentTransferEncoding
	}{
		{
			input: `base64`,
			cte:   email.CTEBase64,
		},
		{
			input: `QUOTeD-PriNTABLE`,
			cte:   email.CTEQuotedPrintable,
		},
		{
			input: `7bit`,
			cte:   email.CTE7bit,
		},
	}
	for i, tt := range tests {
		t.Run(fmt.Sprintf("test_%d", i), func(t *testing.T) {
			cte, err := extractContentTransferEncoding(tt.input)
			if err != nil {
				t.Fatalf("cannot parse part Content-Transfer-Encoding: %s", err)
			}

			if got, want := cte, tt.cte; got != want {
				t.Errorf("got %s want %s", got, want)
			}
		})
	}
}
