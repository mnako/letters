package parser

import (
	"bytes"
	"fmt"
	"io"
	"strings"

	"golang.org/x/text/encoding"

	"github.com/mnako/letters/decoders"
	"github.com/mnako/letters/email"
)

// parseBody parses the body of an email
func (p *Parser) parseBody() error {

	var err error
	switch string(p.email.Headers.ContentType.ContentType) {
	case string(email.ContentTypeTextPlain):
		p.email.Text, err = p.parseText(p.msg.Body, p.encoding, p.cte)
		if err != nil {
			return fmt.Errorf("cannot parse plain text: %w", err)
		}
		return nil

	case string(email.ContentTypeTextEnriched):
		p.email.EnrichedText, err = p.parseText(p.msg.Body, p.encoding, p.cte)
		if err != nil {
			return fmt.Errorf("cannot parse enriched text: %w", err)
		}
		return nil

	case string(email.ContentTypeTextHtml):
		p.email.HTML, err = p.parseText(p.msg.Body, p.encoding, p.cte)
		if err != nil {
			return fmt.Errorf("cannot parse html text: %w", err)
		}
		return nil
	}
	return fmt.Errorf("parse body content type %q not known", p.email.Headers.ContentType.ContentType)

}

// parseText parses the text content of an email body or mime part
func (p *Parser) parseText(t io.Reader, e encoding.Encoding, cte email.ContentTransferEncoding) (string, error) {
	reader := decoders.DecodeContent(t, e, cte)
	textBody, err := io.ReadAll(reader)
	if err != nil {
		return "", fmt.Errorf("cannot read plain text content: %w", err)
	}
	textBody = bytes.ReplaceAll(textBody, []byte("\r\n"), []byte("\n"))
	return strings.TrimSpace(string(textBody)), nil
}
