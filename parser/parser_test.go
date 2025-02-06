package parser

import (
	"os"
	"strings"
	"testing"
)

func TestBasicParser(t *testing.T) {
	p := NewParser()

	if got, want := p.verbose, false; got != want {
		t.Errorf("verbose got %t want %t", got, want)
	}
}

func TestParseRFC822Email(t *testing.T) {

	// https://learn.microsoft.com/en-us/previous-versions/office/developer/exchange-server-2010/aa493918(v=exchg.140)
	msg := `From: someone@example.com
To: someone_else@example.com
Subject: An RFC 822 formatted message

This is the plain text body of the message. Note the blank line
between the header information and the body of the message.`

	reader := strings.NewReader(msg)
	p := NewParser()
	email, err := p.Parse(reader)
	if err != nil {
		t.Fatal(err)
	}
	if got, want := email.Headers.From[0].Address, "someone@example.com"; got != want {
		t.Errorf("got %s want %s", got, want)
	}
	if got, want := email.Headers.To[0].Address, "someone_else@example.com"; got != want {
		t.Errorf("got %s want %s", got, want)
	}
}

func TestParseEnglishPlaintext(t *testing.T) {

	msg, err := os.Open("../tests/test_english_plaintext_ascii_over_7bit.txt")
	if err != nil {
		t.Fatal(err)
	}
	p := NewParser()
	email, err := p.Parse(msg)
	if err != nil {
		t.Fatal(err)
	}
	if got, want := string(email.Headers.MessageID), "Message-Id-1@example.com"; got != want {
		t.Errorf("got %s want %s", got, want)
	}
}
