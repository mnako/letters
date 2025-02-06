package parser

import (
	"bytes"
	"fmt"
	"io"
	"strings"

	"github.com/mnako/letters/decoders"
	"github.com/mnako/letters/email"
)

// parseBody parses the body of an email
func (p *Parser) parseBody() error {

	var err error
	switch p.contentInfo.Type {
	case "text/plain":
		p.email.Text, err = p.parseText(p.msg.Body, p.contentInfo)
		if err != nil {
			return fmt.Errorf("cannot parse plain text: %w", err)
		}
		return nil

	case "text/enriched":
		p.email.EnrichedText, err = p.parseText(p.msg.Body, p.contentInfo)
		if err != nil {
			return fmt.Errorf("cannot parse enriched text: %w", err)
		}
		return nil

	case "text/html":
		p.email.HTML, err = p.parseText(p.msg.Body, p.contentInfo)
		if err != nil {
			return fmt.Errorf("cannot parse html text: %w", err)
		}
		return nil
	}
	return fmt.Errorf("parse body content type %q not known", p.contentInfo.Type)

}

// parseText parses the text content of an email body or mime part
func (p *Parser) parseText(t io.Reader, ci *email.ContentInfo) (string, error) {
	reader := decoders.DecodeContent(t, ci)
	textBody, err := io.ReadAll(reader)
	if err != nil {
		return "", fmt.Errorf("cannot read plain text content: %w", err)
	}
	textBody = bytes.ReplaceAll(textBody, []byte("\r\n"), []byte("\n"))
	return strings.TrimSpace(string(textBody)), nil
}
