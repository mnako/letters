package parser

import (
	"bytes"
	"fmt"
	"testing"

	"github.com/mnako/letters/email"
	"golang.org/x/net/html/charset"
)

func TestParseBody(t *testing.T) {

	tests := []struct {
		bodyInfo           []byte
		charsetForEncoding string
		ct                 string // not email.ContentType?
		cte                email.ContentTransferEncoding
		textLen            int
		enrichedLen        int
		htmlLen            int
	}{

		{
			bodyInfo: []byte(`
The quick brown fox jumps over a lazy dog.
Glib jocks quiz nymph to vex dwarf.
Sphinx of black quartz, judge my vow.
How vexingly quick daft zebras jump!
The five boxing wizards jump quickly.
Jackdaws love my big sphinx of quartz.
Pack my box with five dozen liquor jugs.
`),
			charsetForEncoding: "UTF-8",
			ct:                 "text/plain",
			cte:                email.ContentTransferEncoding("quoted-printable"),
			textLen:            271,
			enrichedLen:        0,
			htmlLen:            0,
		},
		{
			bodyInfo: []byte(`
<bold>The quick brown fox jumps over a lazy dog.</bold>
<italic>Glib jocks quiz nymph to vex dwarf.</italic>
<fixed>Sphinx of black quartz, judge my vow.</fixed>
<underline>How vexingly quick daft zebras jump!</underline>
The five boxing wizards jump quickly.
Jackdaws love my big sphinx of quartz.
Pack my box with five dozen liquor jugs.
`),
			charsetForEncoding: "UTF-8",
			ct:                 "text/enriched",
			cte:                email.ContentTransferEncoding("quoted-printable"),
			textLen:            0,
			enrichedLen:        339,
			htmlLen:            0,
		},
		{
			bodyInfo: []byte(`
<html>
<div dir=3D"ltr">
<p>The quick brown fox jumps over a lazy dog.</p>
<p>Glib jocks quiz nymph to vex dwarf.</p>
<p>Sphinx of black quartz, judge my vow.</p>
<p>How vexingly quick daft zebras jump!</p>
<p>The five boxing wizards jump quickly.</p>
<p>Jackdaws love my big sphinx of quartz.</p>
<p>Pack my box with five dozen liquor jugs.</p>
</div>
</html>
`),
			charsetForEncoding: "UTF-8",
			ct:                 "text/html",
			cte:                email.ContentTransferEncoding("quoted-printable"),
			textLen:            0,
			enrichedLen:        0,
			htmlLen:            358,
		},
	}

	for i, tt := range tests {
		t.Run(fmt.Sprintf("test_%d", i), func(t *testing.T) {
			p := NewParser()
			p.encoding, _ = charset.Lookup(tt.charsetForEncoding)
			p.cte = tt.cte
			p.email.Headers.ContentType = email.ContentTypeHeader{
				ContentType: tt.ct,
			}
			p.msg.Body = bytes.NewReader(tt.bodyInfo)

			err := p.parseBody()
			if err != nil {
				t.Fatal(err)
			}
			if got, want := len(p.email.Text), tt.textLen; got != want {
				t.Errorf("text len got %d want %d", got, want)
			}
			if got, want := len(p.email.EnrichedText), tt.enrichedLen; got != want {
				t.Errorf("enriched len got %d want %d", got, want)
			}
			if got, want := len(p.email.HTML), tt.htmlLen; got != want {
				t.Errorf("html len got %d want %d", got, want)
			}

		})
	}

}
