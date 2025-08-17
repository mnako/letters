package letters_test

import (
	"net/mail"
	"os"
	"reflect"
	"testing"
	"time"

	"github.com/mnako/letters"
)

// Tests for simple helper functions: letters.ParseEmail, and
// letters.ParseEmailHeaders. Tests in this file should correspond to the API
// documented under the ## Quickstart section of the README.md file.

func testEmailHeadersFromFile(
	t *testing.T,
	fp string,
	expectedEmailHeaders letters.Headers,
) {
	t.Helper()

	rawEmail, err := os.Open(fp) //nolint:gosec
	if err != nil {
		t.Errorf("error while reading email from file: %s", err)
		return
	}

	msg, err := mail.ReadMessage(rawEmail)
	if err != nil {
		t.Errorf("error while reading message from file: %s", err)
		return
	}

	parsedEmailHeaders, err := letters.ParseEmailHeaders(msg.Header)
	if err != nil {
		t.Errorf("error while parsing email headers: %s", err)
		return
	}

	if !reflect.DeepEqual(parsedEmailHeaders, expectedEmailHeaders) {
		t.Errorf("email headers are not equal")
		t.Errorf("Got %#v", parsedEmailHeaders)
		t.Errorf("Want %#v", expectedEmailHeaders)
	}
}

func TestParseEmailHeadersEnglishPlaintextAsciiOver7bit(t *testing.T) {
	t.Parallel()

	fp := "tests/test_english_plaintext_ascii_over_7bit.txt"
	tz, _ := time.LoadLocation("Europe/London")
	expectedDate, _ := time.Parse(
		time.RFC1123Z+" (MST)",
		time.Date(2019, time.April, 1, 7, 55, 0, 0, tz).
			Format(time.RFC1123Z+" (MST)"),
	)
	expectedEmailHeaders := letters.Headers{
		Date:    expectedDate,
		Subject: "📧 Test English Pangrams",
		ReplyTo: []*mail.Address{
			{
				Name:    "Alice Sender",
				Address: "alice.sender@example.net",
			},
		},
		Sender: &mail.Address{
			Name:    "Alice Sender",
			Address: "alice.sender@example.com",
		},
		From: []*mail.Address{
			{
				Name:    "Alice Sender",
				Address: "alice.sender@example.com",
			},
			{
				Name:    "Alice Sender",
				Address: "alice.sender@example.net",
			},
		},
		To: []*mail.Address{
			{
				Name:    "Bob Recipient",
				Address: "bob.recipient@example.com",
			},
			{
				Name:    "Carol Recipient",
				Address: "carol.recipient@example.com",
			},
		},
		Cc: []*mail.Address{
			{
				Name:    "Dan Recipient",
				Address: "dan.recipient@example.com",
			},
			{
				Name:    "Eve Recipient",
				Address: "eve.recipient@example.com",
			},
		},
		Bcc: []*mail.Address{
			{
				Name:    "Frank Recipient",
				Address: "frank.recipient@example.com",
			},
			{
				Name:    "Grace Recipient",
				Address: "grace.recipient@example.com",
			},
		},
		MessageID:  "Message-Id-1@example.com",
		InReplyTo:  []letters.MessageId{"Message-Id-0@example.com"},
		References: []letters.MessageId{"Message-Id-0@example.com"},
		Comments:   "Message Header Comment",
		Keywords:   []string{"Keyword 1", "Keyword 2"},
		ContentType: letters.ContentTypeHeader{
			ContentType: "text/plain",
			Params: map[string]string{
				"charset": "ascii",
			},
		},
		ExtraHeaders: map[string][]string{
			"X-Clacks-Overhead": {"GNU Terry Pratchett"},
			"X-Script/function/\t !\"#$%&'()*+,-./;<=>?@[\\]^_`{|}~": {
				"TEST VALUE 1\t !\"#$%&'()*+,-./0123456789:;<=>?@ABCDEFGHIJKLMNOPQRSTUVWXYZ[\\]^_` abcdefghijklmnopqrstuvwxyz{|}~",
				"TEST VALUE 2\t !\"#$%&'()*+,-./0123456789:;<=>?@ABCDEFGHIJKLMNOPQRSTUVWXYZ[\\]^_` abcdefghijklmnopqrstuvwxyz{|}~",
			},
		},
	}

	testEmailHeadersFromFile(t, fp, expectedEmailHeaders)
}
