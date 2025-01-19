package letters

import (
	"net/mail"
	"os"
	"reflect"
	"testing"
	"time"
)

func testEmailHeadersFromFile(t *testing.T, fp string, expectedEmailHeaders Headers) {
	rawEmail, err := os.Open(fp)
	if err != nil {
		t.Errorf("error while reading email from file: %s", err)
		return
	}

	msg, err := mail.ReadMessage(rawEmail)
	if err != nil {
		t.Errorf("error while reading message from file: %s", err)
		return
	}

	parsedEmailHeaders, err := ParseHeaders(msg.Header)
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

func testEmailFromFile(t *testing.T, fp string, expectedEmail Email) {
	rawEmail, err := os.Open(fp)
	if err != nil {
		t.Errorf("error while reading email from file: %s", err)
		return
	}

	emailParser := NewEmailParser()
	parsedEmail, err := emailParser.ParseEmail(rawEmail)
	if err != nil {
		t.Errorf("error while parsing email: %s", err)
		return
	}

	if !reflect.DeepEqual(parsedEmail, expectedEmail) {
		t.Errorf("emails are not equal")
		t.Errorf("Got %#v", parsedEmail)
		t.Errorf("Want %#v", expectedEmail)
	}
}

func TestParseEmailEnglishEmpty(t *testing.T) {
	fp := "tests/test_english_empty.txt"
	expectedEmail := Email{
		Headers: Headers{
			ContentType: ContentTypeHeader{
				ContentType: "text/plain",
				Params:      map[string]string{},
			},
			ExtraHeaders: map[string][]string{},
		},
		Text: `While this email is undeliverable, this test case makes sure that the
parser does not crash, most fields are nullable, and the rest has sane
defaults (e.g. text/plain Content-Type).`,
	}

	testEmailFromFile(t, fp, expectedEmail)
}

func TestParseEmailEnglishNoTextContent(t *testing.T) {
	fp := "tests/test_english_no_text_content.txt"
	tz, _ := time.LoadLocation("Europe/London")
	expectedDate, _ := time.Parse(
		time.RFC1123Z+" (MST)",
		time.Date(2019, time.April, 1, 7, 55, 0, 0, tz).Format(time.RFC1123Z+" (MST)"))
	expectedEmail := Email{
		Headers: Headers{
			Date:    expectedDate,
			Subject: "Test No Text Content, Attachment Only",
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
			MessageID: "Message-Id-1@example.com",
			ContentType: ContentTypeHeader{
				ContentType: "application/pdf",
				Params: map[string]string{
					"name": "attached-pdf-name.pdf",
				},
			},
			ContentDisposition: ContentDispositionHeader{
				ContentDisposition: attachment,
				Params: map[string]string{
					"filename": "attached-pdf-filename.pdf",
				},
			},
			ExtraHeaders: map[string][]string{
				"X-Clacks-Overhead": {"GNU Terry Pratchett"},
				"X-Script/function/\t !\"#$%&'()*+,-./;<=>?@[\\]^_`{|}~": {
					"TEST VALUE 1\t !\"#$%&'()*+,-./0123456789:;<=>?@ABCDEFGHIJKLMNOPQRSTUVWXYZ[\\]^_` abcdefghijklmnopqrstuvwxyz{|}~",
					"TEST VALUE 2\t !\"#$%&'()*+,-./0123456789:;<=>?@ABCDEFGHIJKLMNOPQRSTUVWXYZ[\\]^_` abcdefghijklmnopqrstuvwxyz{|}~",
				},
			},
		},
		Text:         "",
		EnrichedText: "",
		HTML:         "",
		AttachedFiles: []AttachedFile{
			{
				ContentType: ContentTypeHeader{
					ContentType: "application/pdf",
					Params: map[string]string{
						"name": "attached-pdf-name.pdf",
					},
				},
				ContentDisposition: ContentDispositionHeader{
					ContentDisposition: attachment,
					Params: map[string]string{
						"filename": "attached-pdf-filename.pdf",
					},
				},
				Data: []byte{37, 80, 68, 70, 45, 49, 46, 13, 116, 114, 97, 105, 108, 101, 114, 60, 60,
					47, 82, 111, 111, 116, 60, 60, 47, 80, 97, 103, 101, 115, 60, 60, 47, 75, 105, 100, 115, 91, 60,
					60, 47, 77, 101, 100, 105, 97, 66, 111, 120, 91, 48, 32, 48, 32, 51, 32, 51, 93, 62, 62, 93, 62,
					62, 62, 62, 62, 62},
			},
		},
	}

	testEmailFromFile(t, fp, expectedEmail)
}

func TestParseEmailHeadersEnglishPlaintextAsciiOver7bit(t *testing.T) {
	fp := "tests/test_english_plaintext_ascii_over_7bit.txt"
	tz, _ := time.LoadLocation("Europe/London")
	expectedDate, _ := time.Parse(
		time.RFC1123Z+" (MST)",
		time.Date(2019, time.April, 1, 7, 55, 0, 0, tz).Format(time.RFC1123Z+" (MST)"))
	expectedEmailHeaders := Headers{
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
		InReplyTo:  []MessageId{"Message-Id-0@example.com"},
		References: []MessageId{"Message-Id-0@example.com"},
		Comments:   "Message Header Comment",
		Keywords:   []string{"Keyword 1", "Keyword 2"},
		ContentType: ContentTypeHeader{
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

func TestParseEmailEnglishPlaintextAsciiOver7bit(t *testing.T) {
	fp := "tests/test_english_plaintext_ascii_over_7bit.txt"
	tz, _ := time.LoadLocation("Europe/London")
	expectedDate, _ := time.Parse(
		time.RFC1123Z+" (MST)",
		time.Date(2019, time.April, 1, 7, 55, 0, 0, tz).Format(time.RFC1123Z+" (MST)"))
	expectedEmail := Email{
		Headers: Headers{
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
			InReplyTo:  []MessageId{"Message-Id-0@example.com"},
			References: []MessageId{"Message-Id-0@example.com"},
			Comments:   "Message Header Comment",
			Keywords:   []string{"Keyword 1", "Keyword 2"},
			ContentType: ContentTypeHeader{
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
		},
		Text: `The quick brown fox jumps over a lazy dog.
Glib jocks quiz nymph to vex dwarf.
Sphinx of black quartz, judge my vow.
How vexingly quick daft zebras jump!
The five boxing wizards jump quickly.
Jackdaws love my big sphinx of quartz.
Pack my box with five dozen liquor jugs.`,
	}

	testEmailFromFile(t, fp, expectedEmail)
}

func TestParseEmailEnglishPlaintextAsciiOverBase64(t *testing.T) {
	fp := "tests/test_english_plaintext_ascii_over_base64.txt"
	tz, _ := time.LoadLocation("Europe/London")
	expectedDate, _ := time.Parse(
		time.RFC1123Z+" (MST)",
		time.Date(2019, time.April, 1, 7, 55, 0, 0, tz).Format(time.RFC1123Z+" (MST)"))
	expectedEmail := Email{
		Headers: Headers{
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
			InReplyTo:  []MessageId{"Message-Id-0@example.com"},
			References: []MessageId{"Message-Id-0@example.com"},
			Comments:   "Message Header Comment",
			Keywords:   []string{"Keyword 1", "Keyword 2"},
			ResentDate: expectedDate,
			ResentFrom: []*mail.Address{
				{
					Name:    "Alice Sender",
					Address: "alice.sender@example.net",
				},
				{
					Name:    "Alice Sender",
					Address: "alice.sender@example.com",
				},
			},
			ResentSender: &mail.Address{
				Name:    "Alice Sender",
				Address: "alice.sender@example.net",
			},
			ResentTo: []*mail.Address{
				{
					Name:    "Bob Recipient",
					Address: "bob.recipient@example.net",
				},
				{
					Name:    "Carol Recipient",
					Address: "carol.recipient@example.net",
				},
			},
			ResentCc: []*mail.Address{
				{
					Name:    "Dan Recipient",
					Address: "dan.recipient@example.net",
				},
				{
					Name:    "Eve Recipient",
					Address: "eve.recipient@example.net",
				},
			},
			ResentBcc: []*mail.Address{
				{
					Name:    "Frank Recipient",
					Address: "frank.recipient@example.net",
				},
				{
					Name:    "Grace Recipient",
					Address: "grace.recipient@example.net",
				},
			},
			ResentMessageID: "Message-Id-1@example.net",
			ContentType: ContentTypeHeader{
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
		},
		Text: `The quick brown fox jumps over a lazy dog.
Glib jocks quiz nymph to vex dwarf.
Sphinx of black quartz, judge my vow.
How vexingly quick daft zebras jump!
The five boxing wizards jump quickly.
Jackdaws love my big sphinx of quartz.
Pack my box with five dozen liquor jugs.`,
	}

	testEmailFromFile(t, fp, expectedEmail)
}

func TestParseEmailEnglishPlaintextAsciiOverQuotedprintable(t *testing.T) {
	fp := "tests/test_english_plaintext_ascii_over_quoted-printable.txt"
	tz, _ := time.LoadLocation("Europe/London")
	expectedDate, _ := time.Parse(
		time.RFC1123Z+" (MST)",
		time.Date(2019, time.April, 1, 7, 55, 0, 0, tz).Format(time.RFC1123Z+" (MST)"))
	expectedEmail := Email{
		Headers: Headers{
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
			InReplyTo:  []MessageId{"Message-Id-0@example.com"},
			References: []MessageId{"Message-Id-0@example.com"},
			Comments:   "Message Header Comment",
			Keywords:   []string{"Keyword 1", "Keyword 2"},
			ResentDate: expectedDate,
			ResentFrom: []*mail.Address{
				{
					Name:    "Alice Sender",
					Address: "alice.sender@example.net",
				},
				{
					Name:    "Alice Sender",
					Address: "alice.sender@example.com",
				},
			},
			ResentSender: &mail.Address{
				Name:    "Alice Sender",
				Address: "alice.sender@example.net",
			},
			ResentTo: []*mail.Address{
				{
					Name:    "Bob Recipient",
					Address: "bob.recipient@example.net",
				},
				{
					Name:    "Carol Recipient",
					Address: "carol.recipient@example.net",
				},
			},
			ResentCc: []*mail.Address{
				{
					Name:    "Dan Recipient",
					Address: "dan.recipient@example.net",
				},
				{
					Name:    "Eve Recipient",
					Address: "eve.recipient@example.net",
				},
			},
			ResentBcc: []*mail.Address{
				{
					Name:    "Frank Recipient",
					Address: "frank.recipient@example.net",
				},
				{
					Name:    "Grace Recipient",
					Address: "grace.recipient@example.net",
				},
			},
			ResentMessageID: "Message-Id-1@example.net",
			ContentType: ContentTypeHeader{
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
		},
		Text: `The quick brown fox jumps over a lazy dog.
Glib jocks quiz nymph to vex dwarf.
Sphinx of black quartz, judge my vow.
How vexingly quick daft zebras jump!
The five boxing wizards jump quickly.
Jackdaws love my big sphinx of quartz.
Pack my box with five dozen liquor jugs.`,
	}

	testEmailFromFile(t, fp, expectedEmail)
}

func TestParseEmailEnglishPlaintextUtf8Over7bit(t *testing.T) {
	fp := "tests/test_english_plaintext_utf-8_over_7bit.txt"
	tz, _ := time.LoadLocation("Europe/London")
	expectedDate, _ := time.Parse(
		time.RFC1123Z+" (MST)",
		time.Date(2019, time.April, 1, 7, 55, 0, 0, tz).Format(time.RFC1123Z+" (MST)"))
	expectedEmail := Email{
		Headers: Headers{
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
			InReplyTo:  []MessageId{"Message-Id-0@example.com"},
			References: []MessageId{"Message-Id-0@example.com"},
			Comments:   "Message Header Comment",
			Keywords:   []string{"Keyword 1", "Keyword 2"},
			ResentDate: expectedDate,
			ResentFrom: []*mail.Address{
				{
					Name:    "Alice Sender",
					Address: "alice.sender@example.net",
				},
				{
					Name:    "Alice Sender",
					Address: "alice.sender@example.com",
				},
			},
			ResentSender: &mail.Address{
				Name:    "Alice Sender",
				Address: "alice.sender@example.net",
			},
			ResentTo: []*mail.Address{
				{
					Name:    "Bob Recipient",
					Address: "bob.recipient@example.net",
				},
				{
					Name:    "Carol Recipient",
					Address: "carol.recipient@example.net",
				},
			},
			ResentCc: []*mail.Address{
				{
					Name:    "Dan Recipient",
					Address: "dan.recipient@example.net",
				},
				{
					Name:    "Eve Recipient",
					Address: "eve.recipient@example.net",
				},
			},
			ResentBcc: []*mail.Address{
				{
					Name:    "Frank Recipient",
					Address: "frank.recipient@example.net",
				},
				{
					Name:    "Grace Recipient",
					Address: "grace.recipient@example.net",
				},
			},
			ResentMessageID: "Message-Id-1@example.net",
			ContentType: ContentTypeHeader{
				ContentType: "text/plain",
				Params: map[string]string{
					"charset": "utf-8",
				},
			},
			ExtraHeaders: map[string][]string{
				"X-Clacks-Overhead": {"GNU Terry Pratchett"},
				"X-Script/function/\t !\"#$%&'()*+,-./;<=>?@[\\]^_`{|}~": {
					"TEST VALUE 1\t !\"#$%&'()*+,-./0123456789:;<=>?@ABCDEFGHIJKLMNOPQRSTUVWXYZ[\\]^_` abcdefghijklmnopqrstuvwxyz{|}~",
					"TEST VALUE 2\t !\"#$%&'()*+,-./0123456789:;<=>?@ABCDEFGHIJKLMNOPQRSTUVWXYZ[\\]^_` abcdefghijklmnopqrstuvwxyz{|}~",
				},
			},
		},
		Text: `The quick brown fox jumps over a lazy dog.
Glib jocks quiz nymph to vex dwarf.
Sphinx of black quartz, judge my vow.
How vexingly quick daft zebras jump!
The five boxing wizards jump quickly.
Jackdaws love my big sphinx of quartz.
Pack my box with five dozen liquor jugs.`,
	}

	testEmailFromFile(t, fp, expectedEmail)
}

func TestParseEmailEnglishPlaintextUtf8OverBase64(t *testing.T) {
	fp := "tests/test_english_plaintext_utf-8_over_base64.txt"
	tz, _ := time.LoadLocation("Europe/London")
	expectedDate, _ := time.Parse(
		time.RFC1123Z+" (MST)",
		time.Date(2019, time.April, 1, 7, 55, 0, 0, tz).Format(time.RFC1123Z+" (MST)"))
	expectedEmail := Email{
		Headers: Headers{
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
			InReplyTo:  []MessageId{"Message-Id-0@example.com"},
			References: []MessageId{"Message-Id-0@example.com"},
			Comments:   "Message Header Comment",
			Keywords:   []string{"Keyword 1", "Keyword 2"},
			ResentDate: expectedDate,
			ResentFrom: []*mail.Address{
				{
					Name:    "Alice Sender",
					Address: "alice.sender@example.net",
				},
				{
					Name:    "Alice Sender",
					Address: "alice.sender@example.com",
				},
			},
			ResentSender: &mail.Address{
				Name:    "Alice Sender",
				Address: "alice.sender@example.net",
			},
			ResentTo: []*mail.Address{
				{
					Name:    "Bob Recipient",
					Address: "bob.recipient@example.net",
				},
				{
					Name:    "Carol Recipient",
					Address: "carol.recipient@example.net",
				},
			},
			ResentCc: []*mail.Address{
				{
					Name:    "Dan Recipient",
					Address: "dan.recipient@example.net",
				},
				{
					Name:    "Eve Recipient",
					Address: "eve.recipient@example.net",
				},
			},
			ResentBcc: []*mail.Address{
				{
					Name:    "Frank Recipient",
					Address: "frank.recipient@example.net",
				},
				{
					Name:    "Grace Recipient",
					Address: "grace.recipient@example.net",
				},
			},
			ResentMessageID: "Message-Id-1@example.net",
			ContentType: ContentTypeHeader{
				ContentType: "text/plain",
				Params: map[string]string{
					"charset": "utf-8",
				},
			},
			ExtraHeaders: map[string][]string{
				"X-Clacks-Overhead": {"GNU Terry Pratchett"},
				"X-Script/function/\t !\"#$%&'()*+,-./;<=>?@[\\]^_`{|}~": {
					"TEST VALUE 1\t !\"#$%&'()*+,-./0123456789:;<=>?@ABCDEFGHIJKLMNOPQRSTUVWXYZ[\\]^_` abcdefghijklmnopqrstuvwxyz{|}~",
					"TEST VALUE 2\t !\"#$%&'()*+,-./0123456789:;<=>?@ABCDEFGHIJKLMNOPQRSTUVWXYZ[\\]^_` abcdefghijklmnopqrstuvwxyz{|}~",
				},
			},
		},
		Text: `The quick brown fox jumps over a lazy dog.
Glib jocks quiz nymph to vex dwarf.
Sphinx of black quartz, judge my vow.
How vexingly quick daft zebras jump!
The five boxing wizards jump quickly.
Jackdaws love my big sphinx of quartz.
Pack my box with five dozen liquor jugs.`,
	}

	testEmailFromFile(t, fp, expectedEmail)
}

func TestParseEmailEnglishPlaintextUtf8OverQuotedprintable(t *testing.T) {
	fp := "tests/test_english_plaintext_utf-8_over_quoted-printable.txt"
	tz, _ := time.LoadLocation("Europe/London")
	expectedDate, _ := time.Parse(
		time.RFC1123Z+" (MST)",
		time.Date(2019, time.April, 1, 7, 55, 0, 0, tz).Format(time.RFC1123Z+" (MST)"))
	expectedEmail := Email{
		Headers: Headers{
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
			InReplyTo:  []MessageId{"Message-Id-0@example.com"},
			References: []MessageId{"Message-Id-0@example.com"},
			Comments:   "Message Header Comment",
			Keywords:   []string{"Keyword 1", "Keyword 2"},
			ResentDate: expectedDate,
			ResentFrom: []*mail.Address{
				{
					Name:    "Alice Sender",
					Address: "alice.sender@example.net",
				},
				{
					Name:    "Alice Sender",
					Address: "alice.sender@example.com",
				},
			},
			ResentSender: &mail.Address{
				Name:    "Alice Sender",
				Address: "alice.sender@example.net",
			},
			ResentTo: []*mail.Address{
				{
					Name:    "Bob Recipient",
					Address: "bob.recipient@example.net",
				},
				{
					Name:    "Carol Recipient",
					Address: "carol.recipient@example.net",
				},
			},
			ResentCc: []*mail.Address{
				{
					Name:    "Dan Recipient",
					Address: "dan.recipient@example.net",
				},
				{
					Name:    "Eve Recipient",
					Address: "eve.recipient@example.net",
				},
			},
			ResentBcc: []*mail.Address{
				{
					Name:    "Frank Recipient",
					Address: "frank.recipient@example.net",
				},
				{
					Name:    "Grace Recipient",
					Address: "grace.recipient@example.net",
				},
			},
			ResentMessageID: "Message-Id-1@example.net",
			ContentType: ContentTypeHeader{
				ContentType: "text/plain",
				Params: map[string]string{
					"charset": "utf-8",
				},
			},
			ExtraHeaders: map[string][]string{
				"X-Clacks-Overhead": {"GNU Terry Pratchett"},
				"X-Script/function/\t !\"#$%&'()*+,-./;<=>?@[\\]^_`{|}~": {
					"TEST VALUE 1\t !\"#$%&'()*+,-./0123456789:;<=>?@ABCDEFGHIJKLMNOPQRSTUVWXYZ[\\]^_` abcdefghijklmnopqrstuvwxyz{|}~",
					"TEST VALUE 2\t !\"#$%&'()*+,-./0123456789:;<=>?@ABCDEFGHIJKLMNOPQRSTUVWXYZ[\\]^_` abcdefghijklmnopqrstuvwxyz{|}~",
				},
			},
		},
		Text: `The quick brown fox jumps over a lazy dog.
Glib jocks quiz nymph to vex dwarf.
Sphinx of black quartz, judge my vow.
How vexingly quick daft zebras jump!
The five boxing wizards jump quickly.
Jackdaws love my big sphinx of quartz.
Pack my box with five dozen liquor jugs.`,
	}

	testEmailFromFile(t, fp, expectedEmail)
}

func TestParseEmailEnglishMultipartRelatedAsciiOver7bit(t *testing.T) {
	fp := "tests/test_english_multipart_related_ascii_over_7bit.txt"
	tz, _ := time.LoadLocation("Europe/London")
	expectedDate, _ := time.Parse(
		time.RFC1123Z+" (MST)",
		time.Date(2019, time.April, 1, 7, 55, 0, 0, tz).Format(time.RFC1123Z+" (MST)"))
	expectedEmail := Email{
		Headers: Headers{
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
			InReplyTo:  []MessageId{"Message-Id-0@example.com"},
			References: []MessageId{"Message-Id-0@example.com"},
			Comments:   "Message Header Comment",
			Keywords:   []string{"Keyword 1", "Keyword 2"},
			ResentDate: expectedDate,
			ResentFrom: []*mail.Address{
				{
					Name:    "Alice Sender",
					Address: "alice.sender@example.net",
				},
				{
					Name:    "Alice Sender",
					Address: "alice.sender@example.com",
				},
			},
			ResentSender: &mail.Address{
				Name:    "Alice Sender",
				Address: "alice.sender@example.net",
			},
			ResentTo: []*mail.Address{
				{
					Name:    "Bob Recipient",
					Address: "bob.recipient@example.net",
				},
				{
					Name:    "Carol Recipient",
					Address: "carol.recipient@example.net",
				},
			},
			ResentCc: []*mail.Address{
				{
					Name:    "Dan Recipient",
					Address: "dan.recipient@example.net",
				},
				{
					Name:    "Eve Recipient",
					Address: "eve.recipient@example.net",
				},
			},
			ResentBcc: []*mail.Address{
				{
					Name:    "Frank Recipient",
					Address: "frank.recipient@example.net",
				},
				{
					Name:    "Grace Recipient",
					Address: "grace.recipient@example.net",
				},
			},
			ResentMessageID: "Message-Id-1@example.net",
			ContentType: ContentTypeHeader{
				ContentType: "multipart/related",
				Params: map[string]string{
					"boundary": "RelatedBoundaryString",
					"charset":  "ascii",
				},
			},
			ExtraHeaders: map[string][]string{
				"X-Clacks-Overhead": {"GNU Terry Pratchett"},
				"X-Script/function/\t !\"#$%&'()*+,-./;<=>?@[\\]^_`{|}~": {
					"TEST VALUE 1\t !\"#$%&'()*+,-./0123456789:;<=>?@ABCDEFGHIJKLMNOPQRSTUVWXYZ[\\]^_` abcdefghijklmnopqrstuvwxyz{|}~",
					"TEST VALUE 2\t !\"#$%&'()*+,-./0123456789:;<=>?@ABCDEFGHIJKLMNOPQRSTUVWXYZ[\\]^_` abcdefghijklmnopqrstuvwxyz{|}~",
				},
			},
		},
		Text: `The quick brown fox jumps over a lazy dog.
Glib jocks quiz nymph to vex dwarf.
Sphinx of black quartz, judge my vow.
How vexingly quick daft zebras jump!
The five boxing wizards jump quickly.
Jackdaws love my big sphinx of quartz.
Pack my box with five dozen liquor jugs.`,
		EnrichedText: `<bold>The quick brown fox jumps over a lazy dog.</bold>
<italic>Glib jocks quiz nymph to vex dwarf.</italic>
<fixed>Sphinx of black quartz, judge my vow.</fixed>
<underline>How vexingly quick daft zebras jump!</underline>
The five boxing wizards jump quickly.
Jackdaws love my big sphinx of quartz.
Pack my box with five dozen liquor jugs.`,
		HTML: `<html>
<div dir="ltr">
<p>The quick brown fox jumps over a lazy dog.</p>
<p>Glib jocks quiz nymph to vex dwarf.</p>
<p>Sphinx of black quartz, judge my vow.</p>
<p>How vexingly quick daft zebras jump!</p>
<p>The five boxing wizards jump quickly.</p>
<p>Jackdaws love my big sphinx of quartz.</p>
<p>Pack my box with five dozen liquor jugs.</p>
</div>
</html>`,
		InlineFiles: []InlineFile{
			{
				ContentID: "inline-jpg-image.jpg@example.com",
				ContentType: ContentTypeHeader{
					ContentType: "image/jpeg",
					Params: map[string]string{
						"name": "inline-jpg-image-name.jpg",
					},
				},
				ContentDisposition: ContentDispositionHeader{
					ContentDisposition: inline,
					Params: map[string]string{
						"filename": "inline-jpg-image-filename.jpg",
					},
				},
				Data: []byte{255, 216, 255, 219, 0, 67, 0, 3, 2, 2, 2, 2, 2, 3, 2, 2, 2, 3, 3, 3, 3, 4,
					6, 4, 4, 4, 4, 4, 8, 6, 6, 5, 6, 9, 8, 10, 10, 9, 8, 9, 9, 10, 12, 15, 12, 10, 11, 14, 11, 9, 9,
					13, 17, 13, 14, 15, 16, 16, 17, 16, 10, 12, 18, 19, 18, 16, 19, 15, 16, 16, 16, 255, 201, 0, 11, 8,
					0, 1, 0, 1, 1, 1, 17, 0, 255, 204, 0, 6, 0, 16, 16, 5, 255, 218, 0, 8, 1, 1, 0, 0, 63, 0, 210, 207,
					32, 255, 217},
			},
		},
	}

	testEmailFromFile(t, fp, expectedEmail)
}
func TestParseEmailEnglishMultipartRelatedAsciiOverBase64(t *testing.T) {
	fp := "tests/test_english_multipart_related_ascii_over_base64.txt"
	tz, _ := time.LoadLocation("Europe/London")
	expectedDate, _ := time.Parse(
		time.RFC1123Z+" (MST)",
		time.Date(2019, time.April, 1, 7, 55, 0, 0, tz).Format(time.RFC1123Z+" (MST)"))
	expectedEmail := Email{
		Headers: Headers{
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
			InReplyTo:  []MessageId{"Message-Id-0@example.com"},
			References: []MessageId{"Message-Id-0@example.com"},
			Comments:   "Message Header Comment",
			Keywords:   []string{"Keyword 1", "Keyword 2"},
			ResentDate: expectedDate,
			ResentFrom: []*mail.Address{
				{
					Name:    "Alice Sender",
					Address: "alice.sender@example.net",
				},
				{
					Name:    "Alice Sender",
					Address: "alice.sender@example.com",
				},
			},
			ResentSender: &mail.Address{
				Name:    "Alice Sender",
				Address: "alice.sender@example.net",
			},
			ResentTo: []*mail.Address{
				{
					Name:    "Bob Recipient",
					Address: "bob.recipient@example.net",
				},
				{
					Name:    "Carol Recipient",
					Address: "carol.recipient@example.net",
				},
			},
			ResentCc: []*mail.Address{
				{
					Name:    "Dan Recipient",
					Address: "dan.recipient@example.net",
				},
				{
					Name:    "Eve Recipient",
					Address: "eve.recipient@example.net",
				},
			},
			ResentBcc: []*mail.Address{
				{
					Name:    "Frank Recipient",
					Address: "frank.recipient@example.net",
				},
				{
					Name:    "Grace Recipient",
					Address: "grace.recipient@example.net",
				},
			},
			ResentMessageID: "Message-Id-1@example.net",
			ContentType: ContentTypeHeader{
				ContentType: "multipart/related",
				Params: map[string]string{
					"boundary": "RelatedBoundaryString",
					"charset":  "ascii",
				},
			},
			ExtraHeaders: map[string][]string{
				"X-Clacks-Overhead": {"GNU Terry Pratchett"},
				"X-Script/function/\t !\"#$%&'()*+,-./;<=>?@[\\]^_`{|}~": {
					"TEST VALUE 1\t !\"#$%&'()*+,-./0123456789:;<=>?@ABCDEFGHIJKLMNOPQRSTUVWXYZ[\\]^_` abcdefghijklmnopqrstuvwxyz{|}~",
					"TEST VALUE 2\t !\"#$%&'()*+,-./0123456789:;<=>?@ABCDEFGHIJKLMNOPQRSTUVWXYZ[\\]^_` abcdefghijklmnopqrstuvwxyz{|}~",
				},
			},
		},
		Text: `The quick brown fox jumps over a lazy dog.
Glib jocks quiz nymph to vex dwarf.
Sphinx of black quartz, judge my vow.
How vexingly quick daft zebras jump!
The five boxing wizards jump quickly.
Jackdaws love my big sphinx of quartz.
Pack my box with five dozen liquor jugs.`,
		EnrichedText: `<bold>The quick brown fox jumps over a lazy dog.</bold>
<italic>Glib jocks quiz nymph to vex dwarf.</italic>
<fixed>Sphinx of black quartz, judge my vow.</fixed>
<underline>How vexingly quick daft zebras jump!</underline>
The five boxing wizards jump quickly.
Jackdaws love my big sphinx of quartz.
Pack my box with five dozen liquor jugs.`,
		HTML: `<html>
<div dir="ltr">
<p>The quick brown fox jumps over a lazy dog.</p>
<p>Glib jocks quiz nymph to vex dwarf.</p>
<p>Sphinx of black quartz, judge my vow.</p>
<p>How vexingly quick daft zebras jump!</p>
<p>The five boxing wizards jump quickly.</p>
<p>Jackdaws love my big sphinx of quartz.</p>
<p>Pack my box with five dozen liquor jugs.</p>
</div>
</html>`,
		InlineFiles: []InlineFile{
			{
				ContentID: "inline-jpg-image.jpg@example.com",
				ContentType: ContentTypeHeader{
					ContentType: "image/jpeg",
					Params: map[string]string{
						"name": "inline-jpg-image-name.jpg",
					},
				},
				ContentDisposition: ContentDispositionHeader{
					ContentDisposition: inline,
					Params: map[string]string{
						"filename": "inline-jpg-image-filename.jpg",
					},
				},
				Data: []byte{255, 216, 255, 219, 0, 67, 0, 3, 2, 2, 2, 2, 2, 3, 2, 2, 2, 3, 3, 3, 3, 4,
					6, 4, 4, 4, 4, 4, 8, 6, 6, 5, 6, 9, 8, 10, 10, 9, 8, 9, 9, 10, 12, 15, 12, 10, 11, 14, 11, 9, 9,
					13, 17, 13, 14, 15, 16, 16, 17, 16, 10, 12, 18, 19, 18, 16, 19, 15, 16, 16, 16, 255, 201, 0, 11, 8,
					0, 1, 0, 1, 1, 1, 17, 0, 255, 204, 0, 6, 0, 16, 16, 5, 255, 218, 0, 8, 1, 1, 0, 0, 63, 0, 210, 207,
					32, 255, 217},
			},
		},
	}

	testEmailFromFile(t, fp, expectedEmail)
}
func TestParseEmailEnglishMultipartRelatedAsciiOverQuotedprintable(t *testing.T) {
	fp := "tests/test_english_multipart_related_ascii_over_quoted-printable.txt"
	tz, _ := time.LoadLocation("Europe/London")
	expectedDate, _ := time.Parse(
		time.RFC1123Z+" (MST)",
		time.Date(2019, time.April, 1, 7, 55, 0, 0, tz).Format(time.RFC1123Z+" (MST)"))
	expectedEmail := Email{
		Headers: Headers{
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
			InReplyTo:  []MessageId{"Message-Id-0@example.com"},
			References: []MessageId{"Message-Id-0@example.com"},
			Comments:   "Message Header Comment",
			Keywords:   []string{"Keyword 1", "Keyword 2"},
			ResentDate: expectedDate,
			ResentFrom: []*mail.Address{
				{
					Name:    "Alice Sender",
					Address: "alice.sender@example.net",
				},
				{
					Name:    "Alice Sender",
					Address: "alice.sender@example.com",
				},
			},
			ResentSender: &mail.Address{
				Name:    "Alice Sender",
				Address: "alice.sender@example.net",
			},
			ResentTo: []*mail.Address{
				{
					Name:    "Bob Recipient",
					Address: "bob.recipient@example.net",
				},
				{
					Name:    "Carol Recipient",
					Address: "carol.recipient@example.net",
				},
			},
			ResentCc: []*mail.Address{
				{
					Name:    "Dan Recipient",
					Address: "dan.recipient@example.net",
				},
				{
					Name:    "Eve Recipient",
					Address: "eve.recipient@example.net",
				},
			},
			ResentBcc: []*mail.Address{
				{
					Name:    "Frank Recipient",
					Address: "frank.recipient@example.net",
				},
				{
					Name:    "Grace Recipient",
					Address: "grace.recipient@example.net",
				},
			},
			ResentMessageID: "Message-Id-1@example.net",
			ContentType: ContentTypeHeader{
				ContentType: "multipart/related",
				Params: map[string]string{
					"boundary": "RelatedBoundaryString",
					"charset":  "ascii",
				},
			},
			ExtraHeaders: map[string][]string{
				"X-Clacks-Overhead": {"GNU Terry Pratchett"},
				"X-Script/function/\t !\"#$%&'()*+,-./;<=>?@[\\]^_`{|}~": {
					"TEST VALUE 1\t !\"#$%&'()*+,-./0123456789:;<=>?@ABCDEFGHIJKLMNOPQRSTUVWXYZ[\\]^_` abcdefghijklmnopqrstuvwxyz{|}~",
					"TEST VALUE 2\t !\"#$%&'()*+,-./0123456789:;<=>?@ABCDEFGHIJKLMNOPQRSTUVWXYZ[\\]^_` abcdefghijklmnopqrstuvwxyz{|}~",
				},
			},
		},
		Text: `The quick brown fox jumps over a lazy dog.
Glib jocks quiz nymph to vex dwarf.
Sphinx of black quartz, judge my vow.
How vexingly quick daft zebras jump!
The five boxing wizards jump quickly.
Jackdaws love my big sphinx of quartz.
Pack my box with five dozen liquor jugs.`,
		EnrichedText: `<bold>The quick brown fox jumps over a lazy dog.</bold>
<italic>Glib jocks quiz nymph to vex dwarf.</italic>
<fixed>Sphinx of black quartz, judge my vow.</fixed>
<underline>How vexingly quick daft zebras jump!</underline>
The five boxing wizards jump quickly.
Jackdaws love my big sphinx of quartz.
Pack my box with five dozen liquor jugs.`,
		HTML: `<html>
<div dir="ltr">
<p>The quick brown fox jumps over a lazy dog.</p>
<p>Glib jocks quiz nymph to vex dwarf.</p>
<p>Sphinx of black quartz, judge my vow.</p>
<p>How vexingly quick daft zebras jump!</p>
<p>The five boxing wizards jump quickly.</p>
<p>Jackdaws love my big sphinx of quartz.</p>
<p>Pack my box with five dozen liquor jugs.</p>
</div>
</html>`,
		InlineFiles: []InlineFile{
			{
				ContentID: "inline-jpg-image.jpg@example.com",
				ContentType: ContentTypeHeader{
					ContentType: "image/jpeg",
					Params: map[string]string{
						"name": "inline-jpg-image-name.jpg",
					},
				},
				ContentDisposition: ContentDispositionHeader{
					ContentDisposition: inline,
					Params: map[string]string{
						"filename": "inline-jpg-image-filename.jpg",
					},
				},
				Data: []byte{255, 216, 255, 219, 0, 67, 0, 3, 2, 2, 2, 2, 2, 3, 2, 2, 2, 3, 3, 3, 3, 4,
					6, 4, 4, 4, 4, 4, 8, 6, 6, 5, 6, 9, 8, 10, 10, 9, 8, 9, 9, 10, 12, 15, 12, 10, 11, 14, 11, 9, 9,
					13, 17, 13, 14, 15, 16, 16, 17, 16, 10, 12, 18, 19, 18, 16, 19, 15, 16, 16, 16, 255, 201, 0, 11, 8,
					0, 1, 0, 1, 1, 1, 17, 0, 255, 204, 0, 6, 0, 16, 16, 5, 255, 218, 0, 8, 1, 1, 0, 0, 63, 0, 210, 207,
					32, 255, 217},
			},
		},
	}

	testEmailFromFile(t, fp, expectedEmail)
}
func TestParseEmailEnglishMultipartRelatedUtf8Over7bit(t *testing.T) {
	fp := "tests/test_english_multipart_related_utf-8_over_7bit.txt"
	tz, _ := time.LoadLocation("Europe/London")
	expectedDate, _ := time.Parse(
		time.RFC1123Z+" (MST)",
		time.Date(2019, time.April, 1, 7, 55, 0, 0, tz).Format(time.RFC1123Z+" (MST)"))
	expectedEmail := Email{
		Headers: Headers{
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
			InReplyTo:  []MessageId{"Message-Id-0@example.com"},
			References: []MessageId{"Message-Id-0@example.com"},
			Comments:   "Message Header Comment",
			Keywords:   []string{"Keyword 1", "Keyword 2"},
			ResentDate: expectedDate,
			ResentFrom: []*mail.Address{
				{
					Name:    "Alice Sender",
					Address: "alice.sender@example.net",
				},
				{
					Name:    "Alice Sender",
					Address: "alice.sender@example.com",
				},
			},
			ResentSender: &mail.Address{
				Name:    "Alice Sender",
				Address: "alice.sender@example.net",
			},
			ResentTo: []*mail.Address{
				{
					Name:    "Bob Recipient",
					Address: "bob.recipient@example.net",
				},
				{
					Name:    "Carol Recipient",
					Address: "carol.recipient@example.net",
				},
			},
			ResentCc: []*mail.Address{
				{
					Name:    "Dan Recipient",
					Address: "dan.recipient@example.net",
				},
				{
					Name:    "Eve Recipient",
					Address: "eve.recipient@example.net",
				},
			},
			ResentBcc: []*mail.Address{
				{
					Name:    "Frank Recipient",
					Address: "frank.recipient@example.net",
				},
				{
					Name:    "Grace Recipient",
					Address: "grace.recipient@example.net",
				},
			},
			ResentMessageID: "Message-Id-1@example.net",
			ContentType: ContentTypeHeader{
				ContentType: "multipart/related",
				Params: map[string]string{
					"boundary": "RelatedBoundaryString",
					"charset":  "utf-8",
				},
			},
			ExtraHeaders: map[string][]string{
				"X-Clacks-Overhead": {"GNU Terry Pratchett"},
				"X-Script/function/\t !\"#$%&'()*+,-./;<=>?@[\\]^_`{|}~": {
					"TEST VALUE 1\t !\"#$%&'()*+,-./0123456789:;<=>?@ABCDEFGHIJKLMNOPQRSTUVWXYZ[\\]^_` abcdefghijklmnopqrstuvwxyz{|}~",
					"TEST VALUE 2\t !\"#$%&'()*+,-./0123456789:;<=>?@ABCDEFGHIJKLMNOPQRSTUVWXYZ[\\]^_` abcdefghijklmnopqrstuvwxyz{|}~",
				},
			},
		},
		Text: `The quick brown fox jumps over a lazy dog.
Glib jocks quiz nymph to vex dwarf.
Sphinx of black quartz, judge my vow.
How vexingly quick daft zebras jump!
The five boxing wizards jump quickly.
Jackdaws love my big sphinx of quartz.
Pack my box with five dozen liquor jugs.`,
		EnrichedText: `<bold>The quick brown fox jumps over a lazy dog.</bold>
<italic>Glib jocks quiz nymph to vex dwarf.</italic>
<fixed>Sphinx of black quartz, judge my vow.</fixed>
<underline>How vexingly quick daft zebras jump!</underline>
The five boxing wizards jump quickly.
Jackdaws love my big sphinx of quartz.
Pack my box with five dozen liquor jugs.`,
		HTML: `<html>
<div dir="ltr">
<p>The quick brown fox jumps over a lazy dog.</p>
<p>Glib jocks quiz nymph to vex dwarf.</p>
<p>Sphinx of black quartz, judge my vow.</p>
<p>How vexingly quick daft zebras jump!</p>
<p>The five boxing wizards jump quickly.</p>
<p>Jackdaws love my big sphinx of quartz.</p>
<p>Pack my box with five dozen liquor jugs.</p>
</div>
</html>`,
		InlineFiles: []InlineFile{
			{
				ContentID: "inline-jpg-image.jpg@example.com",
				ContentType: ContentTypeHeader{
					ContentType: "image/jpeg",
					Params: map[string]string{
						"name": "inline-jpg-image-name.jpg",
					},
				},
				ContentDisposition: ContentDispositionHeader{
					ContentDisposition: inline,
					Params: map[string]string{
						"filename": "inline-jpg-image-filename.jpg",
					},
				},
				Data: []byte{255, 216, 255, 219, 0, 67, 0, 3, 2, 2, 2, 2, 2, 3, 2, 2, 2, 3, 3, 3, 3, 4,
					6, 4, 4, 4, 4, 4, 8, 6, 6, 5, 6, 9, 8, 10, 10, 9, 8, 9, 9, 10, 12, 15, 12, 10, 11, 14, 11, 9, 9,
					13, 17, 13, 14, 15, 16, 16, 17, 16, 10, 12, 18, 19, 18, 16, 19, 15, 16, 16, 16, 255, 201, 0, 11, 8,
					0, 1, 0, 1, 1, 1, 17, 0, 255, 204, 0, 6, 0, 16, 16, 5, 255, 218, 0, 8, 1, 1, 0, 0, 63, 0, 210, 207,
					32, 255, 217},
			},
		},
	}

	testEmailFromFile(t, fp, expectedEmail)
}
func TestParseEmailEnglishMultipartRelatedUtf8OverBase64(t *testing.T) {
	fp := "tests/test_english_multipart_related_utf-8_over_base64.txt"
	tz, _ := time.LoadLocation("Europe/London")
	expectedDate, _ := time.Parse(
		time.RFC1123Z+" (MST)",
		time.Date(2019, time.April, 1, 7, 55, 0, 0, tz).Format(time.RFC1123Z+" (MST)"))
	expectedEmail := Email{
		Headers: Headers{
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
			InReplyTo:  []MessageId{"Message-Id-0@example.com"},
			References: []MessageId{"Message-Id-0@example.com"},
			Comments:   "Message Header Comment",
			Keywords:   []string{"Keyword 1", "Keyword 2"},
			ResentDate: expectedDate,
			ResentFrom: []*mail.Address{
				{
					Name:    "Alice Sender",
					Address: "alice.sender@example.net",
				},
				{
					Name:    "Alice Sender",
					Address: "alice.sender@example.com",
				},
			},
			ResentSender: &mail.Address{
				Name:    "Alice Sender",
				Address: "alice.sender@example.net",
			},
			ResentTo: []*mail.Address{
				{
					Name:    "Bob Recipient",
					Address: "bob.recipient@example.net",
				},
				{
					Name:    "Carol Recipient",
					Address: "carol.recipient@example.net",
				},
			},
			ResentCc: []*mail.Address{
				{
					Name:    "Dan Recipient",
					Address: "dan.recipient@example.net",
				},
				{
					Name:    "Eve Recipient",
					Address: "eve.recipient@example.net",
				},
			},
			ResentBcc: []*mail.Address{
				{
					Name:    "Frank Recipient",
					Address: "frank.recipient@example.net",
				},
				{
					Name:    "Grace Recipient",
					Address: "grace.recipient@example.net",
				},
			},
			ResentMessageID: "Message-Id-1@example.net",
			ContentType: ContentTypeHeader{
				ContentType: "multipart/related",
				Params: map[string]string{
					"boundary": "RelatedBoundaryString",
					"charset":  "utf-8",
				},
			},
			ExtraHeaders: map[string][]string{
				"X-Clacks-Overhead": {"GNU Terry Pratchett"},
				"X-Script/function/\t !\"#$%&'()*+,-./;<=>?@[\\]^_`{|}~": {
					"TEST VALUE 1\t !\"#$%&'()*+,-./0123456789:;<=>?@ABCDEFGHIJKLMNOPQRSTUVWXYZ[\\]^_` abcdefghijklmnopqrstuvwxyz{|}~",
					"TEST VALUE 2\t !\"#$%&'()*+,-./0123456789:;<=>?@ABCDEFGHIJKLMNOPQRSTUVWXYZ[\\]^_` abcdefghijklmnopqrstuvwxyz{|}~",
				},
			},
		},
		Text: `The quick brown fox jumps over a lazy dog.
Glib jocks quiz nymph to vex dwarf.
Sphinx of black quartz, judge my vow.
How vexingly quick daft zebras jump!
The five boxing wizards jump quickly.
Jackdaws love my big sphinx of quartz.
Pack my box with five dozen liquor jugs.`,
		EnrichedText: `<bold>The quick brown fox jumps over a lazy dog.</bold>
<italic>Glib jocks quiz nymph to vex dwarf.</italic>
<fixed>Sphinx of black quartz, judge my vow.</fixed>
<underline>How vexingly quick daft zebras jump!</underline>
The five boxing wizards jump quickly.
Jackdaws love my big sphinx of quartz.
Pack my box with five dozen liquor jugs.`,
		HTML: `<html>
<div dir="ltr">
<p>The quick brown fox jumps over a lazy dog.</p>
<p>Glib jocks quiz nymph to vex dwarf.</p>
<p>Sphinx of black quartz, judge my vow.</p>
<p>How vexingly quick daft zebras jump!</p>
<p>The five boxing wizards jump quickly.</p>
<p>Jackdaws love my big sphinx of quartz.</p>
<p>Pack my box with five dozen liquor jugs.</p>
</div>
</html>`,
		InlineFiles: []InlineFile{
			{
				ContentID: "inline-jpg-image.jpg@example.com",
				ContentType: ContentTypeHeader{
					ContentType: "image/jpeg",
					Params: map[string]string{
						"name": "inline-jpg-image-name.jpg",
					},
				},
				ContentDisposition: ContentDispositionHeader{
					ContentDisposition: inline,
					Params: map[string]string{
						"filename": "inline-jpg-image-filename.jpg",
					},
				},
				Data: []byte{255, 216, 255, 219, 0, 67, 0, 3, 2, 2, 2, 2, 2, 3, 2, 2, 2, 3, 3, 3, 3, 4,
					6, 4, 4, 4, 4, 4, 8, 6, 6, 5, 6, 9, 8, 10, 10, 9, 8, 9, 9, 10, 12, 15, 12, 10, 11, 14, 11, 9, 9,
					13, 17, 13, 14, 15, 16, 16, 17, 16, 10, 12, 18, 19, 18, 16, 19, 15, 16, 16, 16, 255, 201, 0, 11, 8,
					0, 1, 0, 1, 1, 1, 17, 0, 255, 204, 0, 6, 0, 16, 16, 5, 255, 218, 0, 8, 1, 1, 0, 0, 63, 0, 210, 207,
					32, 255, 217},
			},
		},
	}

	testEmailFromFile(t, fp, expectedEmail)
}
func TestParseEmailEnglishMultipartRelatedUtf8OverQuotedprintable(t *testing.T) {
	fp := "tests/test_english_multipart_related_utf-8_over_quoted-printable.txt"
	tz, _ := time.LoadLocation("Europe/London")
	expectedDate, _ := time.Parse(
		time.RFC1123Z+" (MST)",
		time.Date(2019, time.April, 1, 7, 55, 0, 0, tz).Format(time.RFC1123Z+" (MST)"))
	expectedEmail := Email{
		Headers: Headers{
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
			InReplyTo:  []MessageId{"Message-Id-0@example.com"},
			References: []MessageId{"Message-Id-0@example.com"},
			Comments:   "Message Header Comment",
			Keywords:   []string{"Keyword 1", "Keyword 2"},
			ResentDate: expectedDate,
			ResentFrom: []*mail.Address{
				{
					Name:    "Alice Sender",
					Address: "alice.sender@example.net",
				},
				{
					Name:    "Alice Sender",
					Address: "alice.sender@example.com",
				},
			},
			ResentSender: &mail.Address{
				Name:    "Alice Sender",
				Address: "alice.sender@example.net",
			},
			ResentTo: []*mail.Address{
				{
					Name:    "Bob Recipient",
					Address: "bob.recipient@example.net",
				},
				{
					Name:    "Carol Recipient",
					Address: "carol.recipient@example.net",
				},
			},
			ResentCc: []*mail.Address{
				{
					Name:    "Dan Recipient",
					Address: "dan.recipient@example.net",
				},
				{
					Name:    "Eve Recipient",
					Address: "eve.recipient@example.net",
				},
			},
			ResentBcc: []*mail.Address{
				{
					Name:    "Frank Recipient",
					Address: "frank.recipient@example.net",
				},
				{
					Name:    "Grace Recipient",
					Address: "grace.recipient@example.net",
				},
			},
			ResentMessageID: "Message-Id-1@example.net",
			ContentType: ContentTypeHeader{
				ContentType: "multipart/related",
				Params: map[string]string{
					"boundary": "RelatedBoundaryString",
					"charset":  "utf-8",
				},
			},
			ExtraHeaders: map[string][]string{
				"X-Clacks-Overhead": {"GNU Terry Pratchett"},
				"X-Script/function/\t !\"#$%&'()*+,-./;<=>?@[\\]^_`{|}~": {
					"TEST VALUE 1\t !\"#$%&'()*+,-./0123456789:;<=>?@ABCDEFGHIJKLMNOPQRSTUVWXYZ[\\]^_` abcdefghijklmnopqrstuvwxyz{|}~",
					"TEST VALUE 2\t !\"#$%&'()*+,-./0123456789:;<=>?@ABCDEFGHIJKLMNOPQRSTUVWXYZ[\\]^_` abcdefghijklmnopqrstuvwxyz{|}~",
				},
			},
		},
		Text: `The quick brown fox jumps over a lazy dog.
Glib jocks quiz nymph to vex dwarf.
Sphinx of black quartz, judge my vow.
How vexingly quick daft zebras jump!
The five boxing wizards jump quickly.
Jackdaws love my big sphinx of quartz.
Pack my box with five dozen liquor jugs.`,
		EnrichedText: `<bold>The quick brown fox jumps over a lazy dog.</bold>
<italic>Glib jocks quiz nymph to vex dwarf.</italic>
<fixed>Sphinx of black quartz, judge my vow.</fixed>
<underline>How vexingly quick daft zebras jump!</underline>
The five boxing wizards jump quickly.
Jackdaws love my big sphinx of quartz.
Pack my box with five dozen liquor jugs.`,
		HTML: `<html>
<div dir="ltr">
<p>The quick brown fox jumps over a lazy dog.</p>
<p>Glib jocks quiz nymph to vex dwarf.</p>
<p>Sphinx of black quartz, judge my vow.</p>
<p>How vexingly quick daft zebras jump!</p>
<p>The five boxing wizards jump quickly.</p>
<p>Jackdaws love my big sphinx of quartz.</p>
<p>Pack my box with five dozen liquor jugs.</p>
</div>
</html>`,
		InlineFiles: []InlineFile{
			{
				ContentID: "inline-jpg-image.jpg@example.com",
				ContentType: ContentTypeHeader{
					ContentType: "image/jpeg",
					Params: map[string]string{
						"name": "inline-jpg-image-name.jpg",
					},
				},
				ContentDisposition: ContentDispositionHeader{
					ContentDisposition: inline,
					Params: map[string]string{
						"filename": "inline-jpg-image-filename.jpg",
					},
				},
				Data: []byte{255, 216, 255, 219, 0, 67, 0, 3, 2, 2, 2, 2, 2, 3, 2, 2, 2, 3, 3, 3, 3, 4,
					6, 4, 4, 4, 4, 4, 8, 6, 6, 5, 6, 9, 8, 10, 10, 9, 8, 9, 9, 10, 12, 15, 12, 10, 11, 14, 11, 9, 9,
					13, 17, 13, 14, 15, 16, 16, 17, 16, 10, 12, 18, 19, 18, 16, 19, 15, 16, 16, 16, 255, 201, 0, 11, 8,
					0, 1, 0, 1, 1, 1, 17, 0, 255, 204, 0, 6, 0, 16, 16, 5, 255, 218, 0, 8, 1, 1, 0, 0, 63, 0, 210, 207,
					32, 255, 217},
			},
		},
	}

	testEmailFromFile(t, fp, expectedEmail)
}
func TestParseEmailEnglishMultipartMixedAsciiOver7bit(t *testing.T) {
	fp := "tests/test_english_multipart_mixed_ascii_over_7bit.txt"
	tz, _ := time.LoadLocation("Europe/London")
	expectedDate, _ := time.Parse(
		time.RFC1123Z+" (MST)",
		time.Date(2019, time.April, 1, 7, 55, 0, 0, tz).Format(time.RFC1123Z+" (MST)"))
	expectedEmail := Email{
		Headers: Headers{
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
			InReplyTo:  []MessageId{"Message-Id-0@example.com"},
			References: []MessageId{"Message-Id-0@example.com"},
			Comments:   "Message Header Comment",
			Keywords:   []string{"Keyword 1", "Keyword 2"},
			ResentDate: expectedDate,
			ResentFrom: []*mail.Address{
				{
					Name:    "Alice Sender",
					Address: "alice.sender@example.net",
				},
				{
					Name:    "Alice Sender",
					Address: "alice.sender@example.com",
				},
			},
			ResentSender: &mail.Address{
				Name:    "Alice Sender",
				Address: "alice.sender@example.net",
			},
			ResentTo: []*mail.Address{
				{
					Name:    "Bob Recipient",
					Address: "bob.recipient@example.net",
				},
				{
					Name:    "Carol Recipient",
					Address: "carol.recipient@example.net",
				},
			},
			ResentCc: []*mail.Address{
				{
					Name:    "Dan Recipient",
					Address: "dan.recipient@example.net",
				},
				{
					Name:    "Eve Recipient",
					Address: "eve.recipient@example.net",
				},
			},
			ResentBcc: []*mail.Address{
				{
					Name:    "Frank Recipient",
					Address: "frank.recipient@example.net",
				},
				{
					Name:    "Grace Recipient",
					Address: "grace.recipient@example.net",
				},
			},
			ResentMessageID: "Message-Id-1@example.net",
			ContentType: ContentTypeHeader{
				ContentType: "multipart/mixed",
				Params: map[string]string{
					"boundary": "MixedBoundaryString",
					"charset":  "ascii",
				},
			},
			ExtraHeaders: map[string][]string{
				"X-Clacks-Overhead": {"GNU Terry Pratchett"},
				"X-Script/function/\t !\"#$%&'()*+,-./;<=>?@[\\]^_`{|}~": {
					"TEST VALUE 1\t !\"#$%&'()*+,-./0123456789:;<=>?@ABCDEFGHIJKLMNOPQRSTUVWXYZ[\\]^_` abcdefghijklmnopqrstuvwxyz{|}~",
					"TEST VALUE 2\t !\"#$%&'()*+,-./0123456789:;<=>?@ABCDEFGHIJKLMNOPQRSTUVWXYZ[\\]^_` abcdefghijklmnopqrstuvwxyz{|}~",
				},
			},
		},
		Text: `The quick brown fox jumps over a lazy dog.
Glib jocks quiz nymph to vex dwarf.
Sphinx of black quartz, judge my vow.
How vexingly quick daft zebras jump!
The five boxing wizards jump quickly.
Jackdaws love my big sphinx of quartz.
Pack my box with five dozen liquor jugs.`,
		EnrichedText: `<bold>The quick brown fox jumps over a lazy dog.</bold>
<italic>Glib jocks quiz nymph to vex dwarf.</italic>
<fixed>Sphinx of black quartz, judge my vow.</fixed>
<underline>How vexingly quick daft zebras jump!</underline>
The five boxing wizards jump quickly.
Jackdaws love my big sphinx of quartz.
Pack my box with five dozen liquor jugs.`,
		HTML: `<html>
<div dir="ltr">
<p>The quick brown fox jumps over a lazy dog.</p>
<p>Glib jocks quiz nymph to vex dwarf.</p>
<p>Sphinx of black quartz, judge my vow.</p>
<p>How vexingly quick daft zebras jump!</p>
<p>The five boxing wizards jump quickly.</p>
<p>Jackdaws love my big sphinx of quartz.</p>
<p>Pack my box with five dozen liquor jugs.</p>
</div>
</html>`,
		InlineFiles: []InlineFile{
			{
				ContentType: ContentTypeHeader{
					ContentType: "image/jpeg",
					Params: map[string]string{
						"name": "inline-jpg-image-without-disposition.jpg",
					},
				},
				ContentDisposition: ContentDispositionHeader{
					ContentDisposition: "",
					Params:             map[string]string(nil),
				},
				Data: []byte{255, 216, 255, 219, 0, 67, 0, 3, 2, 2, 2, 2, 2, 3, 2, 2, 2, 3, 3, 3, 3, 4,
					6, 4, 4, 4, 4, 4, 8, 6, 6, 5, 6, 9, 8, 10, 10, 9, 8, 9, 9, 10, 12, 15, 12, 10, 11, 14, 11, 9, 9,
					13, 17, 13, 14, 15, 16, 16, 17, 16, 10, 12, 18, 19, 18, 16, 19, 15, 16, 16, 16, 255, 201, 0, 11, 8,
					0, 1, 0, 1, 1, 1, 17, 0, 255, 204, 0, 6, 0, 16, 16, 5, 255, 218, 0, 8, 1, 1, 0, 0, 63, 0, 210, 207,
					32, 255, 217},
			},
			{
				ContentID: "inline-jpg-image.jpg@example.com",
				ContentType: ContentTypeHeader{
					ContentType: "image/jpeg",
					Params: map[string]string{
						"name": "inline-jpg-image-name.jpg",
					},
				},
				ContentDisposition: ContentDispositionHeader{
					ContentDisposition: inline,
					Params: map[string]string{
						"filename": "inline-jpg-image-filename.jpg",
					},
				},
				Data: []byte{255, 216, 255, 219, 0, 67, 0, 3, 2, 2, 2, 2, 2, 3, 2, 2, 2, 3, 3, 3, 3, 4,
					6, 4, 4, 4, 4, 4, 8, 6, 6, 5, 6, 9, 8, 10, 10, 9, 8, 9, 9, 10, 12, 15, 12, 10, 11, 14, 11, 9, 9,
					13, 17, 13, 14, 15, 16, 16, 17, 16, 10, 12, 18, 19, 18, 16, 19, 15, 16, 16, 16, 255, 201, 0, 11, 8,
					0, 1, 0, 1, 1, 1, 17, 0, 255, 204, 0, 6, 0, 16, 16, 5, 255, 218, 0, 8, 1, 1, 0, 0, 63, 0, 210, 207,
					32, 255, 217},
			},
		},
		AttachedFiles: []AttachedFile{
			{
				ContentType: ContentTypeHeader{
					ContentType: "application/pdf",
					Params: map[string]string{
						"name": "attached-pdf-name.pdf",
					},
				},
				ContentDisposition: ContentDispositionHeader{
					ContentDisposition: attachment,
					Params: map[string]string{
						"filename": "attached-pdf-filename.pdf",
					},
				},
				Data: []byte{37, 80, 68, 70, 45, 49, 46, 13, 116, 114, 97, 105, 108, 101, 114, 60, 60,
					47, 82, 111, 111, 116, 60, 60, 47, 80, 97, 103, 101, 115, 60, 60, 47, 75, 105, 100, 115, 91, 60,
					60, 47, 77, 101, 100, 105, 97, 66, 111, 120, 91, 48, 32, 48, 32, 51, 32, 51, 93, 62, 62, 93, 62,
					62, 62, 62, 62, 62},
			},
			{
				ContentType: ContentTypeHeader{
					ContentType: "application/pdf",
					Params: map[string]string{
						"name": "attached-pdf-without-disposition.pdf",
					},
				},
				ContentDisposition: ContentDispositionHeader{
					ContentDisposition: "",
					Params:             map[string]string(nil),
				},
				Data: []byte{37, 80, 68, 70, 45, 49, 46, 13, 116, 114, 97, 105, 108, 101, 114, 60, 60,
					47, 82, 111, 111, 116, 60, 60, 47, 80, 97, 103, 101, 115, 60, 60, 47, 75, 105, 100, 115, 91, 60,
					60, 47, 77, 101, 100, 105, 97, 66, 111, 120, 91, 48, 32, 48, 32, 51, 32, 51, 93, 62, 62, 93, 62,
					62, 62, 62, 62, 62},
			},
			{
				ContentType: ContentTypeHeader{
					ContentType: "application/json",
					Params: map[string]string{
						"name": "attached-json-name.json",
					},
				},
				ContentDisposition: ContentDispositionHeader{
					ContentDisposition: attachment,
					Params: map[string]string{
						"filename": "attached-json-filename.json",
					},
				},
				Data: []byte{123, 34, 102, 111, 111, 34, 58, 34, 98, 97, 114, 34, 125},
			},
			{
				ContentType: ContentTypeHeader{
					ContentType: "text/plain",
					Params: map[string]string{
						"name": "attached-text-plain-name.txt",
					},
				},
				ContentDisposition: ContentDispositionHeader{
					ContentDisposition: attachment,
					Params: map[string]string{
						"filename": "attached-text-plain-filename.txt",
					},
				},
				Data: []byte{84, 101, 120, 116, 47, 112, 108, 97, 105, 110, 32, 99, 111, 110, 116, 101, 110, 116, 32,
					97, 115, 32, 97, 110, 32, 97, 116, 116, 97, 99, 104, 101, 100, 32, 46, 116, 120, 116, 32, 102, 105,
					108, 101, 46},
			},
			{
				ContentType: ContentTypeHeader{
					ContentType: "text/html",
					Params: map[string]string{
						"name": "attached-text-html-name.html",
					},
				},
				ContentDisposition: ContentDispositionHeader{
					ContentDisposition: attachment,
					Params: map[string]string{
						"filename": "attached-text-html-filename.html",
					},
				},
				Data: []byte{84, 101, 120, 116, 47, 104, 116, 109, 108, 32, 99, 111, 110, 116, 101, 110, 116, 32, 97,
					115, 32, 97, 110, 32, 97, 116, 116, 97, 99, 104, 101, 100, 32, 46, 104, 116, 109, 108, 32, 102,
					105, 108, 101, 46},
			},
		},
	}

	testEmailFromFile(t, fp, expectedEmail)
}
func TestParseEmailEnglishMultipartMixedAsciiOverBase64(t *testing.T) {
	fp := "tests/test_english_multipart_mixed_ascii_over_base64.txt"
	tz, _ := time.LoadLocation("Europe/London")
	expectedDate, _ := time.Parse(
		time.RFC1123Z+" (MST)",
		time.Date(2019, time.April, 1, 7, 55, 0, 0, tz).Format(time.RFC1123Z+" (MST)"))
	expectedEmail := Email{
		Headers: Headers{
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
			InReplyTo:  []MessageId{"Message-Id-0@example.com"},
			References: []MessageId{"Message-Id-0@example.com"},
			Comments:   "Message Header Comment",
			Keywords:   []string{"Keyword 1", "Keyword 2"},
			ResentDate: expectedDate,
			ResentFrom: []*mail.Address{
				{
					Name:    "Alice Sender",
					Address: "alice.sender@example.net",
				},
				{
					Name:    "Alice Sender",
					Address: "alice.sender@example.com",
				},
			},
			ResentSender: &mail.Address{
				Name:    "Alice Sender",
				Address: "alice.sender@example.net",
			},
			ResentTo: []*mail.Address{
				{
					Name:    "Bob Recipient",
					Address: "bob.recipient@example.net",
				},
				{
					Name:    "Carol Recipient",
					Address: "carol.recipient@example.net",
				},
			},
			ResentCc: []*mail.Address{
				{
					Name:    "Dan Recipient",
					Address: "dan.recipient@example.net",
				},
				{
					Name:    "Eve Recipient",
					Address: "eve.recipient@example.net",
				},
			},
			ResentBcc: []*mail.Address{
				{
					Name:    "Frank Recipient",
					Address: "frank.recipient@example.net",
				},
				{
					Name:    "Grace Recipient",
					Address: "grace.recipient@example.net",
				},
			},
			ResentMessageID: "Message-Id-1@example.net",
			ContentType: ContentTypeHeader{
				ContentType: "multipart/mixed",
				Params: map[string]string{
					"boundary": "MixedBoundaryString",
					"charset":  "ascii",
				},
			},
			ExtraHeaders: map[string][]string{
				"X-Clacks-Overhead": {"GNU Terry Pratchett"},
				"X-Script/function/\t !\"#$%&'()*+,-./;<=>?@[\\]^_`{|}~": {
					"TEST VALUE 1\t !\"#$%&'()*+,-./0123456789:;<=>?@ABCDEFGHIJKLMNOPQRSTUVWXYZ[\\]^_` abcdefghijklmnopqrstuvwxyz{|}~",
					"TEST VALUE 2\t !\"#$%&'()*+,-./0123456789:;<=>?@ABCDEFGHIJKLMNOPQRSTUVWXYZ[\\]^_` abcdefghijklmnopqrstuvwxyz{|}~",
				},
			},
		},
		Text: `The quick brown fox jumps over a lazy dog.
Glib jocks quiz nymph to vex dwarf.
Sphinx of black quartz, judge my vow.
How vexingly quick daft zebras jump!
The five boxing wizards jump quickly.
Jackdaws love my big sphinx of quartz.
Pack my box with five dozen liquor jugs.`,
		EnrichedText: `<bold>The quick brown fox jumps over a lazy dog.</bold>
<italic>Glib jocks quiz nymph to vex dwarf.</italic>
<fixed>Sphinx of black quartz, judge my vow.</fixed>
<underline>How vexingly quick daft zebras jump!</underline>
The five boxing wizards jump quickly.
Jackdaws love my big sphinx of quartz.
Pack my box with five dozen liquor jugs.`,
		HTML: `<html>
<div dir="ltr">
<p>The quick brown fox jumps over a lazy dog.</p>
<p>Glib jocks quiz nymph to vex dwarf.</p>
<p>Sphinx of black quartz, judge my vow.</p>
<p>How vexingly quick daft zebras jump!</p>
<p>The five boxing wizards jump quickly.</p>
<p>Jackdaws love my big sphinx of quartz.</p>
<p>Pack my box with five dozen liquor jugs.</p>
</div>
</html>`,
		InlineFiles: []InlineFile{
			{
				ContentType: ContentTypeHeader{
					ContentType: "image/jpeg",
					Params: map[string]string{
						"name": "inline-jpg-image-without-disposition.jpg",
					},
				},
				ContentDisposition: ContentDispositionHeader{
					ContentDisposition: "",
					Params:             map[string]string(nil),
				},
				Data: []byte{255, 216, 255, 219, 0, 67, 0, 3, 2, 2, 2, 2, 2, 3, 2, 2, 2, 3, 3, 3, 3, 4,
					6, 4, 4, 4, 4, 4, 8, 6, 6, 5, 6, 9, 8, 10, 10, 9, 8, 9, 9, 10, 12, 15, 12, 10, 11, 14, 11, 9, 9,
					13, 17, 13, 14, 15, 16, 16, 17, 16, 10, 12, 18, 19, 18, 16, 19, 15, 16, 16, 16, 255, 201, 0, 11, 8,
					0, 1, 0, 1, 1, 1, 17, 0, 255, 204, 0, 6, 0, 16, 16, 5, 255, 218, 0, 8, 1, 1, 0, 0, 63, 0, 210, 207,
					32, 255, 217},
			},
			{
				ContentID: "inline-jpg-image.jpg@example.com",
				ContentType: ContentTypeHeader{
					ContentType: "image/jpeg",
					Params: map[string]string{
						"name": "inline-jpg-image-name.jpg",
					},
				},
				ContentDisposition: ContentDispositionHeader{
					ContentDisposition: inline,
					Params: map[string]string{
						"filename": "inline-jpg-image-filename.jpg",
					},
				},
				Data: []byte{255, 216, 255, 219, 0, 67, 0, 3, 2, 2, 2, 2, 2, 3, 2, 2, 2, 3, 3, 3, 3, 4,
					6, 4, 4, 4, 4, 4, 8, 6, 6, 5, 6, 9, 8, 10, 10, 9, 8, 9, 9, 10, 12, 15, 12, 10, 11, 14, 11, 9, 9,
					13, 17, 13, 14, 15, 16, 16, 17, 16, 10, 12, 18, 19, 18, 16, 19, 15, 16, 16, 16, 255, 201, 0, 11, 8,
					0, 1, 0, 1, 1, 1, 17, 0, 255, 204, 0, 6, 0, 16, 16, 5, 255, 218, 0, 8, 1, 1, 0, 0, 63, 0, 210, 207,
					32, 255, 217},
			},
		},
		AttachedFiles: []AttachedFile{
			{
				ContentType: ContentTypeHeader{
					ContentType: "application/pdf",
					Params: map[string]string{
						"name": "attached-pdf-name.pdf",
					},
				},
				ContentDisposition: ContentDispositionHeader{
					ContentDisposition: attachment,
					Params: map[string]string{
						"filename": "attached-pdf-filename.pdf",
					},
				},
				Data: []byte{37, 80, 68, 70, 45, 49, 46, 13, 116, 114, 97, 105, 108, 101, 114, 60, 60,
					47, 82, 111, 111, 116, 60, 60, 47, 80, 97, 103, 101, 115, 60, 60, 47, 75, 105, 100, 115, 91, 60,
					60, 47, 77, 101, 100, 105, 97, 66, 111, 120, 91, 48, 32, 48, 32, 51, 32, 51, 93, 62, 62, 93, 62,
					62, 62, 62, 62, 62},
			},
			{
				ContentType: ContentTypeHeader{
					ContentType: "application/pdf",
					Params: map[string]string{
						"name": "attached-pdf-without-disposition.pdf",
					},
				},
				ContentDisposition: ContentDispositionHeader{
					ContentDisposition: "",
					Params:             map[string]string(nil),
				},
				Data: []byte{37, 80, 68, 70, 45, 49, 46, 13, 116, 114, 97, 105, 108, 101, 114, 60, 60,
					47, 82, 111, 111, 116, 60, 60, 47, 80, 97, 103, 101, 115, 60, 60, 47, 75, 105, 100, 115, 91, 60,
					60, 47, 77, 101, 100, 105, 97, 66, 111, 120, 91, 48, 32, 48, 32, 51, 32, 51, 93, 62, 62, 93, 62,
					62, 62, 62, 62, 62},
			},
			{
				ContentType: ContentTypeHeader{
					ContentType: "application/json",
					Params: map[string]string{
						"name": "attached-json-name.json",
					},
				},
				ContentDisposition: ContentDispositionHeader{
					ContentDisposition: attachment,
					Params: map[string]string{
						"filename": "attached-json-filename.json",
					},
				},
				Data: []byte{123, 34, 102, 111, 111, 34, 58, 34, 98, 97, 114, 34, 125},
			},
			{
				ContentType: ContentTypeHeader{
					ContentType: "text/plain",
					Params: map[string]string{
						"name": "attached-text-plain-name.txt",
					},
				},
				ContentDisposition: ContentDispositionHeader{
					ContentDisposition: attachment,
					Params: map[string]string{
						"filename": "attached-text-plain-filename.txt",
					},
				},
				Data: []byte{84, 101, 120, 116, 47, 112, 108, 97, 105, 110, 32, 99, 111, 110, 116, 101, 110, 116, 32,
					97, 115, 32, 97, 110, 32, 97, 116, 116, 97, 99, 104, 101, 100, 32, 46, 116, 120, 116, 32, 102, 105,
					108, 101, 46},
			},
			{
				ContentType: ContentTypeHeader{
					ContentType: "text/html",
					Params: map[string]string{
						"name": "attached-text-html-name.html",
					},
				},
				ContentDisposition: ContentDispositionHeader{
					ContentDisposition: attachment,
					Params: map[string]string{
						"filename": "attached-text-html-filename.html",
					},
				},
				Data: []byte{84, 101, 120, 116, 47, 104, 116, 109, 108, 32, 99, 111, 110, 116, 101, 110, 116, 32, 97,
					115, 32, 97, 110, 32, 97, 116, 116, 97, 99, 104, 101, 100, 32, 46, 104, 116, 109, 108, 32, 102,
					105, 108, 101, 46},
			},
		},
	}

	testEmailFromFile(t, fp, expectedEmail)
}
func TestParseEmailEnglishMultipartMixedAsciiOverQuotedprintable(t *testing.T) {
	fp := "tests/test_english_multipart_mixed_ascii_over_quoted-printable.txt"
	tz, _ := time.LoadLocation("Europe/London")
	expectedDate, _ := time.Parse(
		time.RFC1123Z+" (MST)",
		time.Date(2019, time.April, 1, 7, 55, 0, 0, tz).Format(time.RFC1123Z+" (MST)"))
	expectedEmail := Email{
		Headers: Headers{
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
			InReplyTo:  []MessageId{"Message-Id-0@example.com"},
			References: []MessageId{"Message-Id-0@example.com"},
			Comments:   "Message Header Comment",
			Keywords:   []string{"Keyword 1", "Keyword 2"},
			ResentDate: expectedDate,
			ResentFrom: []*mail.Address{
				{
					Name:    "Alice Sender",
					Address: "alice.sender@example.net",
				},
				{
					Name:    "Alice Sender",
					Address: "alice.sender@example.com",
				},
			},
			ResentSender: &mail.Address{
				Name:    "Alice Sender",
				Address: "alice.sender@example.net",
			},
			ResentTo: []*mail.Address{
				{
					Name:    "Bob Recipient",
					Address: "bob.recipient@example.net",
				},
				{
					Name:    "Carol Recipient",
					Address: "carol.recipient@example.net",
				},
			},
			ResentCc: []*mail.Address{
				{
					Name:    "Dan Recipient",
					Address: "dan.recipient@example.net",
				},
				{
					Name:    "Eve Recipient",
					Address: "eve.recipient@example.net",
				},
			},
			ResentBcc: []*mail.Address{
				{
					Name:    "Frank Recipient",
					Address: "frank.recipient@example.net",
				},
				{
					Name:    "Grace Recipient",
					Address: "grace.recipient@example.net",
				},
			},
			ResentMessageID: "Message-Id-1@example.net",
			ContentType: ContentTypeHeader{
				ContentType: "multipart/mixed",
				Params: map[string]string{
					"boundary": "MixedBoundaryString",
					"charset":  "ascii",
				},
			},
			ExtraHeaders: map[string][]string{
				"X-Clacks-Overhead": {"GNU Terry Pratchett"},
				"X-Script/function/\t !\"#$%&'()*+,-./;<=>?@[\\]^_`{|}~": {
					"TEST VALUE 1\t !\"#$%&'()*+,-./0123456789:;<=>?@ABCDEFGHIJKLMNOPQRSTUVWXYZ[\\]^_` abcdefghijklmnopqrstuvwxyz{|}~",
					"TEST VALUE 2\t !\"#$%&'()*+,-./0123456789:;<=>?@ABCDEFGHIJKLMNOPQRSTUVWXYZ[\\]^_` abcdefghijklmnopqrstuvwxyz{|}~",
				},
			},
		},
		Text: `The quick brown fox jumps over a lazy dog.
Glib jocks quiz nymph to vex dwarf.
Sphinx of black quartz, judge my vow.
How vexingly quick daft zebras jump!
The five boxing wizards jump quickly.
Jackdaws love my big sphinx of quartz.
Pack my box with five dozen liquor jugs.`,
		EnrichedText: `<bold>The quick brown fox jumps over a lazy dog.</bold>
<italic>Glib jocks quiz nymph to vex dwarf.</italic>
<fixed>Sphinx of black quartz, judge my vow.</fixed>
<underline>How vexingly quick daft zebras jump!</underline>
The five boxing wizards jump quickly.
Jackdaws love my big sphinx of quartz.
Pack my box with five dozen liquor jugs.`,
		HTML: `<html>
<div dir="ltr">
<p>The quick brown fox jumps over a lazy dog.</p>
<p>Glib jocks quiz nymph to vex dwarf.</p>
<p>Sphinx of black quartz, judge my vow.</p>
<p>How vexingly quick daft zebras jump!</p>
<p>The five boxing wizards jump quickly.</p>
<p>Jackdaws love my big sphinx of quartz.</p>
<p>Pack my box with five dozen liquor jugs.</p>
</div>
</html>`,
		InlineFiles: []InlineFile{
			{
				ContentType: ContentTypeHeader{
					ContentType: "image/jpeg",
					Params: map[string]string{
						"name": "inline-jpg-image-without-disposition.jpg",
					},
				},
				ContentDisposition: ContentDispositionHeader{
					ContentDisposition: "",
					Params:             map[string]string(nil),
				},
				Data: []byte{255, 216, 255, 219, 0, 67, 0, 3, 2, 2, 2, 2, 2, 3, 2, 2, 2, 3, 3, 3, 3, 4,
					6, 4, 4, 4, 4, 4, 8, 6, 6, 5, 6, 9, 8, 10, 10, 9, 8, 9, 9, 10, 12, 15, 12, 10, 11, 14, 11, 9, 9,
					13, 17, 13, 14, 15, 16, 16, 17, 16, 10, 12, 18, 19, 18, 16, 19, 15, 16, 16, 16, 255, 201, 0, 11, 8,
					0, 1, 0, 1, 1, 1, 17, 0, 255, 204, 0, 6, 0, 16, 16, 5, 255, 218, 0, 8, 1, 1, 0, 0, 63, 0, 210, 207,
					32, 255, 217},
			},
			{
				ContentID: "inline-jpg-image.jpg@example.com",
				ContentType: ContentTypeHeader{
					ContentType: "image/jpeg",
					Params: map[string]string{
						"name": "inline-jpg-image-name.jpg",
					},
				},
				ContentDisposition: ContentDispositionHeader{
					ContentDisposition: inline,
					Params: map[string]string{
						"filename": "inline-jpg-image-filename.jpg",
					},
				},
				Data: []byte{255, 216, 255, 219, 0, 67, 0, 3, 2, 2, 2, 2, 2, 3, 2, 2, 2, 3, 3, 3, 3, 4,
					6, 4, 4, 4, 4, 4, 8, 6, 6, 5, 6, 9, 8, 10, 10, 9, 8, 9, 9, 10, 12, 15, 12, 10, 11, 14, 11, 9, 9,
					13, 17, 13, 14, 15, 16, 16, 17, 16, 10, 12, 18, 19, 18, 16, 19, 15, 16, 16, 16, 255, 201, 0, 11, 8,
					0, 1, 0, 1, 1, 1, 17, 0, 255, 204, 0, 6, 0, 16, 16, 5, 255, 218, 0, 8, 1, 1, 0, 0, 63, 0, 210, 207,
					32, 255, 217},
			},
		},
		AttachedFiles: []AttachedFile{
			{
				ContentType: ContentTypeHeader{
					ContentType: "application/pdf",
					Params: map[string]string{
						"name": "attached-pdf-name.pdf",
					},
				},
				ContentDisposition: ContentDispositionHeader{
					ContentDisposition: attachment,
					Params: map[string]string{
						"filename": "attached-pdf-filename.pdf",
					},
				},
				Data: []byte{37, 80, 68, 70, 45, 49, 46, 13, 116, 114, 97, 105, 108, 101, 114, 60, 60,
					47, 82, 111, 111, 116, 60, 60, 47, 80, 97, 103, 101, 115, 60, 60, 47, 75, 105, 100, 115, 91, 60,
					60, 47, 77, 101, 100, 105, 97, 66, 111, 120, 91, 48, 32, 48, 32, 51, 32, 51, 93, 62, 62, 93, 62,
					62, 62, 62, 62, 62},
			},
			{
				ContentType: ContentTypeHeader{
					ContentType: "application/pdf",
					Params: map[string]string{
						"name": "attached-pdf-without-disposition.pdf",
					},
				},
				ContentDisposition: ContentDispositionHeader{
					ContentDisposition: "",
					Params:             map[string]string(nil),
				},
				Data: []byte{37, 80, 68, 70, 45, 49, 46, 13, 116, 114, 97, 105, 108, 101, 114, 60, 60,
					47, 82, 111, 111, 116, 60, 60, 47, 80, 97, 103, 101, 115, 60, 60, 47, 75, 105, 100, 115, 91, 60,
					60, 47, 77, 101, 100, 105, 97, 66, 111, 120, 91, 48, 32, 48, 32, 51, 32, 51, 93, 62, 62, 93, 62,
					62, 62, 62, 62, 62},
			},
			{
				ContentType: ContentTypeHeader{
					ContentType: "application/json",
					Params: map[string]string{
						"name": "attached-json-name.json",
					},
				},
				ContentDisposition: ContentDispositionHeader{
					ContentDisposition: attachment,
					Params: map[string]string{
						"filename": "attached-json-filename.json",
					},
				},
				Data: []byte{123, 34, 102, 111, 111, 34, 58, 34, 98, 97, 114, 34, 125},
			},
			{
				ContentType: ContentTypeHeader{
					ContentType: "text/plain",
					Params: map[string]string{
						"name": "attached-text-plain-name.txt",
					},
				},
				ContentDisposition: ContentDispositionHeader{
					ContentDisposition: attachment,
					Params: map[string]string{
						"filename": "attached-text-plain-filename.txt",
					},
				},
				Data: []byte{84, 101, 120, 116, 47, 112, 108, 97, 105, 110, 32, 99, 111, 110, 116, 101, 110, 116, 32,
					97, 115, 32, 97, 110, 32, 97, 116, 116, 97, 99, 104, 101, 100, 32, 46, 116, 120, 116, 32, 102, 105,
					108, 101, 46},
			},
			{
				ContentType: ContentTypeHeader{
					ContentType: "text/html",
					Params: map[string]string{
						"name": "attached-text-html-name.html",
					},
				},
				ContentDisposition: ContentDispositionHeader{
					ContentDisposition: attachment,
					Params: map[string]string{
						"filename": "attached-text-html-filename.html",
					},
				},
				Data: []byte{84, 101, 120, 116, 47, 104, 116, 109, 108, 32, 99, 111, 110, 116, 101, 110, 116, 32, 97,
					115, 32, 97, 110, 32, 97, 116, 116, 97, 99, 104, 101, 100, 32, 46, 104, 116, 109, 108, 32, 102,
					105, 108, 101, 46},
			},
		},
	}

	testEmailFromFile(t, fp, expectedEmail)
}
func TestParseEmailEnglishMultipartMixedUtf8Over7bit(t *testing.T) {
	fp := "tests/test_english_multipart_mixed_utf-8_over_7bit.txt"
	tz, _ := time.LoadLocation("Europe/London")
	expectedDate, _ := time.Parse(
		time.RFC1123Z+" (MST)",
		time.Date(2019, time.April, 1, 7, 55, 0, 0, tz).Format(time.RFC1123Z+" (MST)"))
	expectedEmail := Email{
		Headers: Headers{
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
			InReplyTo:  []MessageId{"Message-Id-0@example.com"},
			References: []MessageId{"Message-Id-0@example.com"},
			Comments:   "Message Header Comment",
			Keywords:   []string{"Keyword 1", "Keyword 2"},
			ResentDate: expectedDate,
			ResentFrom: []*mail.Address{
				{
					Name:    "Alice Sender",
					Address: "alice.sender@example.net",
				},
				{
					Name:    "Alice Sender",
					Address: "alice.sender@example.com",
				},
			},
			ResentSender: &mail.Address{
				Name:    "Alice Sender",
				Address: "alice.sender@example.net",
			},
			ResentTo: []*mail.Address{
				{
					Name:    "Bob Recipient",
					Address: "bob.recipient@example.net",
				},
				{
					Name:    "Carol Recipient",
					Address: "carol.recipient@example.net",
				},
			},
			ResentCc: []*mail.Address{
				{
					Name:    "Dan Recipient",
					Address: "dan.recipient@example.net",
				},
				{
					Name:    "Eve Recipient",
					Address: "eve.recipient@example.net",
				},
			},
			ResentBcc: []*mail.Address{
				{
					Name:    "Frank Recipient",
					Address: "frank.recipient@example.net",
				},
				{
					Name:    "Grace Recipient",
					Address: "grace.recipient@example.net",
				},
			},
			ResentMessageID: "Message-Id-1@example.net",
			ContentType: ContentTypeHeader{
				ContentType: "multipart/mixed",
				Params: map[string]string{
					"boundary": "MixedBoundaryString",
					"charset":  "utf-8",
				},
			},
			ExtraHeaders: map[string][]string{
				"X-Clacks-Overhead": {"GNU Terry Pratchett"},
				"X-Script/function/\t !\"#$%&'()*+,-./;<=>?@[\\]^_`{|}~": {
					"TEST VALUE 1\t !\"#$%&'()*+,-./0123456789:;<=>?@ABCDEFGHIJKLMNOPQRSTUVWXYZ[\\]^_` abcdefghijklmnopqrstuvwxyz{|}~",
					"TEST VALUE 2\t !\"#$%&'()*+,-./0123456789:;<=>?@ABCDEFGHIJKLMNOPQRSTUVWXYZ[\\]^_` abcdefghijklmnopqrstuvwxyz{|}~",
				},
			},
		},
		Text: `The quick brown fox jumps over a lazy dog.
Glib jocks quiz nymph to vex dwarf.
Sphinx of black quartz, judge my vow.
How vexingly quick daft zebras jump!
The five boxing wizards jump quickly.
Jackdaws love my big sphinx of quartz.
Pack my box with five dozen liquor jugs.`,
		EnrichedText: `<bold>The quick brown fox jumps over a lazy dog.</bold>
<italic>Glib jocks quiz nymph to vex dwarf.</italic>
<fixed>Sphinx of black quartz, judge my vow.</fixed>
<underline>How vexingly quick daft zebras jump!</underline>
The five boxing wizards jump quickly.
Jackdaws love my big sphinx of quartz.
Pack my box with five dozen liquor jugs.`,
		HTML: `<html>
<div dir="ltr">
<p>The quick brown fox jumps over a lazy dog.</p>
<p>Glib jocks quiz nymph to vex dwarf.</p>
<p>Sphinx of black quartz, judge my vow.</p>
<p>How vexingly quick daft zebras jump!</p>
<p>The five boxing wizards jump quickly.</p>
<p>Jackdaws love my big sphinx of quartz.</p>
<p>Pack my box with five dozen liquor jugs.</p>
</div>
</html>`,
		InlineFiles: []InlineFile{
			{
				ContentType: ContentTypeHeader{
					ContentType: "image/jpeg",
					Params: map[string]string{
						"name": "inline-jpg-image-without-disposition.jpg",
					},
				},
				ContentDisposition: ContentDispositionHeader{
					ContentDisposition: "",
					Params:             map[string]string(nil),
				},
				Data: []byte{255, 216, 255, 219, 0, 67, 0, 3, 2, 2, 2, 2, 2, 3, 2, 2, 2, 3, 3, 3, 3, 4,
					6, 4, 4, 4, 4, 4, 8, 6, 6, 5, 6, 9, 8, 10, 10, 9, 8, 9, 9, 10, 12, 15, 12, 10, 11, 14, 11, 9, 9,
					13, 17, 13, 14, 15, 16, 16, 17, 16, 10, 12, 18, 19, 18, 16, 19, 15, 16, 16, 16, 255, 201, 0, 11, 8,
					0, 1, 0, 1, 1, 1, 17, 0, 255, 204, 0, 6, 0, 16, 16, 5, 255, 218, 0, 8, 1, 1, 0, 0, 63, 0, 210, 207,
					32, 255, 217},
			},
			{
				ContentID: "inline-jpg-image.jpg@example.com",
				ContentType: ContentTypeHeader{
					ContentType: "image/jpeg",
					Params: map[string]string{
						"name": "inline-jpg-image-name.jpg",
					},
				},
				ContentDisposition: ContentDispositionHeader{
					ContentDisposition: inline,
					Params: map[string]string{
						"filename": "inline-jpg-image-filename.jpg",
					},
				},
				Data: []byte{255, 216, 255, 219, 0, 67, 0, 3, 2, 2, 2, 2, 2, 3, 2, 2, 2, 3, 3, 3, 3, 4,
					6, 4, 4, 4, 4, 4, 8, 6, 6, 5, 6, 9, 8, 10, 10, 9, 8, 9, 9, 10, 12, 15, 12, 10, 11, 14, 11, 9, 9,
					13, 17, 13, 14, 15, 16, 16, 17, 16, 10, 12, 18, 19, 18, 16, 19, 15, 16, 16, 16, 255, 201, 0, 11, 8,
					0, 1, 0, 1, 1, 1, 17, 0, 255, 204, 0, 6, 0, 16, 16, 5, 255, 218, 0, 8, 1, 1, 0, 0, 63, 0, 210, 207,
					32, 255, 217},
			},
		},
		AttachedFiles: []AttachedFile{
			{
				ContentType: ContentTypeHeader{
					ContentType: "application/pdf",
					Params: map[string]string{
						"name": "attached-pdf-name.pdf",
					},
				},
				ContentDisposition: ContentDispositionHeader{
					ContentDisposition: attachment,
					Params: map[string]string{
						"filename": "attached-pdf-filename.pdf",
					},
				},
				Data: []byte{37, 80, 68, 70, 45, 49, 46, 13, 116, 114, 97, 105, 108, 101, 114, 60, 60,
					47, 82, 111, 111, 116, 60, 60, 47, 80, 97, 103, 101, 115, 60, 60, 47, 75, 105, 100, 115, 91, 60,
					60, 47, 77, 101, 100, 105, 97, 66, 111, 120, 91, 48, 32, 48, 32, 51, 32, 51, 93, 62, 62, 93, 62,
					62, 62, 62, 62, 62},
			},
			{
				ContentType: ContentTypeHeader{
					ContentType: "application/pdf",
					Params: map[string]string{
						"name": "attached-pdf-without-disposition.pdf",
					},
				},
				ContentDisposition: ContentDispositionHeader{
					ContentDisposition: "",
					Params:             map[string]string(nil),
				},
				Data: []byte{37, 80, 68, 70, 45, 49, 46, 13, 116, 114, 97, 105, 108, 101, 114, 60, 60,
					47, 82, 111, 111, 116, 60, 60, 47, 80, 97, 103, 101, 115, 60, 60, 47, 75, 105, 100, 115, 91, 60,
					60, 47, 77, 101, 100, 105, 97, 66, 111, 120, 91, 48, 32, 48, 32, 51, 32, 51, 93, 62, 62, 93, 62,
					62, 62, 62, 62, 62},
			},
			{
				ContentType: ContentTypeHeader{
					ContentType: "application/json",
					Params: map[string]string{
						"name": "attached-json-name.json",
					},
				},
				ContentDisposition: ContentDispositionHeader{
					ContentDisposition: attachment,
					Params: map[string]string{
						"filename": "attached-json-filename.json",
					},
				},
				Data: []byte{123, 34, 102, 111, 111, 34, 58, 34, 98, 97, 114, 34, 125},
			},
			{
				ContentType: ContentTypeHeader{
					ContentType: "text/plain",
					Params: map[string]string{
						"name": "attached-text-plain-name.txt",
					},
				},
				ContentDisposition: ContentDispositionHeader{
					ContentDisposition: attachment,
					Params: map[string]string{
						"filename": "attached-text-plain-filename.txt",
					},
				},
				Data: []byte{84, 101, 120, 116, 47, 112, 108, 97, 105, 110, 32, 99, 111, 110, 116, 101, 110, 116, 32,
					97, 115, 32, 97, 110, 32, 97, 116, 116, 97, 99, 104, 101, 100, 32, 46, 116, 120, 116, 32, 102, 105,
					108, 101, 46},
			},
			{
				ContentType: ContentTypeHeader{
					ContentType: "text/html",
					Params: map[string]string{
						"name": "attached-text-html-name.html",
					},
				},
				ContentDisposition: ContentDispositionHeader{
					ContentDisposition: attachment,
					Params: map[string]string{
						"filename": "attached-text-html-filename.html",
					},
				},
				Data: []byte{84, 101, 120, 116, 47, 104, 116, 109, 108, 32, 99, 111, 110, 116, 101, 110, 116, 32, 97,
					115, 32, 97, 110, 32, 97, 116, 116, 97, 99, 104, 101, 100, 32, 46, 104, 116, 109, 108, 32, 102,
					105, 108, 101, 46},
			},
		},
	}

	testEmailFromFile(t, fp, expectedEmail)
}
func TestParseEmailEnglishMultipartMixedUtf8OverBase64(t *testing.T) {
	fp := "tests/test_english_multipart_mixed_utf-8_over_base64.txt"
	tz, _ := time.LoadLocation("Europe/London")
	expectedDate, _ := time.Parse(
		time.RFC1123Z+" (MST)",
		time.Date(2019, time.April, 1, 7, 55, 0, 0, tz).Format(time.RFC1123Z+" (MST)"))
	expectedEmail := Email{
		Headers: Headers{
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
			InReplyTo:  []MessageId{"Message-Id-0@example.com"},
			References: []MessageId{"Message-Id-0@example.com"},
			Comments:   "Message Header Comment",
			Keywords:   []string{"Keyword 1", "Keyword 2"},
			ResentDate: expectedDate,
			ResentFrom: []*mail.Address{
				{
					Name:    "Alice Sender",
					Address: "alice.sender@example.net",
				},
				{
					Name:    "Alice Sender",
					Address: "alice.sender@example.com",
				},
			},
			ResentSender: &mail.Address{
				Name:    "Alice Sender",
				Address: "alice.sender@example.net",
			},
			ResentTo: []*mail.Address{
				{
					Name:    "Bob Recipient",
					Address: "bob.recipient@example.net",
				},
				{
					Name:    "Carol Recipient",
					Address: "carol.recipient@example.net",
				},
			},
			ResentCc: []*mail.Address{
				{
					Name:    "Dan Recipient",
					Address: "dan.recipient@example.net",
				},
				{
					Name:    "Eve Recipient",
					Address: "eve.recipient@example.net",
				},
			},
			ResentBcc: []*mail.Address{
				{
					Name:    "Frank Recipient",
					Address: "frank.recipient@example.net",
				},
				{
					Name:    "Grace Recipient",
					Address: "grace.recipient@example.net",
				},
			},
			ResentMessageID: "Message-Id-1@example.net",
			ContentType: ContentTypeHeader{
				ContentType: "multipart/mixed",
				Params: map[string]string{
					"boundary": "MixedBoundaryString",
					"charset":  "utf-8",
				},
			},
			ExtraHeaders: map[string][]string{
				"X-Clacks-Overhead": {"GNU Terry Pratchett"},
				"X-Script/function/\t !\"#$%&'()*+,-./;<=>?@[\\]^_`{|}~": {
					"TEST VALUE 1\t !\"#$%&'()*+,-./0123456789:;<=>?@ABCDEFGHIJKLMNOPQRSTUVWXYZ[\\]^_` abcdefghijklmnopqrstuvwxyz{|}~",
					"TEST VALUE 2\t !\"#$%&'()*+,-./0123456789:;<=>?@ABCDEFGHIJKLMNOPQRSTUVWXYZ[\\]^_` abcdefghijklmnopqrstuvwxyz{|}~",
				},
			},
		},
		Text: `The quick brown fox jumps over a lazy dog.
Glib jocks quiz nymph to vex dwarf.
Sphinx of black quartz, judge my vow.
How vexingly quick daft zebras jump!
The five boxing wizards jump quickly.
Jackdaws love my big sphinx of quartz.
Pack my box with five dozen liquor jugs.`,
		EnrichedText: `<bold>The quick brown fox jumps over a lazy dog.</bold>
<italic>Glib jocks quiz nymph to vex dwarf.</italic>
<fixed>Sphinx of black quartz, judge my vow.</fixed>
<underline>How vexingly quick daft zebras jump!</underline>
The five boxing wizards jump quickly.
Jackdaws love my big sphinx of quartz.
Pack my box with five dozen liquor jugs.`,
		HTML: `<html>
<div dir="ltr">
<p>The quick brown fox jumps over a lazy dog.</p>
<p>Glib jocks quiz nymph to vex dwarf.</p>
<p>Sphinx of black quartz, judge my vow.</p>
<p>How vexingly quick daft zebras jump!</p>
<p>The five boxing wizards jump quickly.</p>
<p>Jackdaws love my big sphinx of quartz.</p>
<p>Pack my box with five dozen liquor jugs.</p>
</div>
</html>`,
		InlineFiles: []InlineFile{
			{
				ContentType: ContentTypeHeader{
					ContentType: "image/jpeg",
					Params: map[string]string{
						"name": "inline-jpg-image-without-disposition.jpg",
					},
				},
				ContentDisposition: ContentDispositionHeader{
					ContentDisposition: "",
					Params:             map[string]string(nil),
				},
				Data: []byte{255, 216, 255, 219, 0, 67, 0, 3, 2, 2, 2, 2, 2, 3, 2, 2, 2, 3, 3, 3, 3, 4,
					6, 4, 4, 4, 4, 4, 8, 6, 6, 5, 6, 9, 8, 10, 10, 9, 8, 9, 9, 10, 12, 15, 12, 10, 11, 14, 11, 9, 9,
					13, 17, 13, 14, 15, 16, 16, 17, 16, 10, 12, 18, 19, 18, 16, 19, 15, 16, 16, 16, 255, 201, 0, 11, 8,
					0, 1, 0, 1, 1, 1, 17, 0, 255, 204, 0, 6, 0, 16, 16, 5, 255, 218, 0, 8, 1, 1, 0, 0, 63, 0, 210, 207,
					32, 255, 217},
			},
			{
				ContentID: "inline-jpg-image.jpg@example.com",
				ContentType: ContentTypeHeader{
					ContentType: "image/jpeg",
					Params: map[string]string{
						"name": "inline-jpg-image-name.jpg",
					},
				},
				ContentDisposition: ContentDispositionHeader{
					ContentDisposition: inline,
					Params: map[string]string{
						"filename": "inline-jpg-image-filename.jpg",
					},
				},
				Data: []byte{255, 216, 255, 219, 0, 67, 0, 3, 2, 2, 2, 2, 2, 3, 2, 2, 2, 3, 3, 3, 3, 4,
					6, 4, 4, 4, 4, 4, 8, 6, 6, 5, 6, 9, 8, 10, 10, 9, 8, 9, 9, 10, 12, 15, 12, 10, 11, 14, 11, 9, 9,
					13, 17, 13, 14, 15, 16, 16, 17, 16, 10, 12, 18, 19, 18, 16, 19, 15, 16, 16, 16, 255, 201, 0, 11, 8,
					0, 1, 0, 1, 1, 1, 17, 0, 255, 204, 0, 6, 0, 16, 16, 5, 255, 218, 0, 8, 1, 1, 0, 0, 63, 0, 210, 207,
					32, 255, 217},
			},
		},
		AttachedFiles: []AttachedFile{
			{
				ContentType: ContentTypeHeader{
					ContentType: "application/pdf",
					Params: map[string]string{
						"name": "attached-pdf-name.pdf",
					},
				},
				ContentDisposition: ContentDispositionHeader{
					ContentDisposition: attachment,
					Params: map[string]string{
						"filename": "attached-pdf-filename.pdf",
					},
				},
				Data: []byte{37, 80, 68, 70, 45, 49, 46, 13, 116, 114, 97, 105, 108, 101, 114, 60, 60,
					47, 82, 111, 111, 116, 60, 60, 47, 80, 97, 103, 101, 115, 60, 60, 47, 75, 105, 100, 115, 91, 60,
					60, 47, 77, 101, 100, 105, 97, 66, 111, 120, 91, 48, 32, 48, 32, 51, 32, 51, 93, 62, 62, 93, 62,
					62, 62, 62, 62, 62},
			},
			{
				ContentType: ContentTypeHeader{
					ContentType: "application/pdf",
					Params: map[string]string{
						"name": "attached-pdf-without-disposition.pdf",
					},
				},
				ContentDisposition: ContentDispositionHeader{
					ContentDisposition: "",
					Params:             map[string]string(nil),
				},
				Data: []byte{37, 80, 68, 70, 45, 49, 46, 13, 116, 114, 97, 105, 108, 101, 114, 60, 60,
					47, 82, 111, 111, 116, 60, 60, 47, 80, 97, 103, 101, 115, 60, 60, 47, 75, 105, 100, 115, 91, 60,
					60, 47, 77, 101, 100, 105, 97, 66, 111, 120, 91, 48, 32, 48, 32, 51, 32, 51, 93, 62, 62, 93, 62,
					62, 62, 62, 62, 62},
			},
			{
				ContentType: ContentTypeHeader{
					ContentType: "application/json",
					Params: map[string]string{
						"name": "attached-json-name.json",
					},
				},
				ContentDisposition: ContentDispositionHeader{
					ContentDisposition: attachment,
					Params: map[string]string{
						"filename": "attached-json-filename.json",
					},
				},
				Data: []byte{123, 34, 102, 111, 111, 34, 58, 34, 98, 97, 114, 34, 125},
			},
			{
				ContentType: ContentTypeHeader{
					ContentType: "text/plain",
					Params: map[string]string{
						"name": "attached-text-plain-name.txt",
					},
				},
				ContentDisposition: ContentDispositionHeader{
					ContentDisposition: attachment,
					Params: map[string]string{
						"filename": "attached-text-plain-filename.txt",
					},
				},
				Data: []byte{84, 101, 120, 116, 47, 112, 108, 97, 105, 110, 32, 99, 111, 110, 116, 101, 110, 116, 32,
					97, 115, 32, 97, 110, 32, 97, 116, 116, 97, 99, 104, 101, 100, 32, 46, 116, 120, 116, 32, 102, 105,
					108, 101, 46},
			},
			{
				ContentType: ContentTypeHeader{
					ContentType: "text/html",
					Params: map[string]string{
						"name": "attached-text-html-name.html",
					},
				},
				ContentDisposition: ContentDispositionHeader{
					ContentDisposition: attachment,
					Params: map[string]string{
						"filename": "attached-text-html-filename.html",
					},
				},
				Data: []byte{84, 101, 120, 116, 47, 104, 116, 109, 108, 32, 99, 111, 110, 116, 101, 110, 116, 32, 97,
					115, 32, 97, 110, 32, 97, 116, 116, 97, 99, 104, 101, 100, 32, 46, 104, 116, 109, 108, 32, 102,
					105, 108, 101, 46},
			},
		},
	}

	testEmailFromFile(t, fp, expectedEmail)
}
func TestParseEmailEnglishMultipartMixedUtf8OverQuotedprintable(t *testing.T) {
	fp := "tests/test_english_multipart_mixed_utf-8_over_quoted-printable.txt"
	tz, _ := time.LoadLocation("Europe/London")
	expectedDate, _ := time.Parse(
		time.RFC1123Z+" (MST)",
		time.Date(2019, time.April, 1, 7, 55, 0, 0, tz).Format(time.RFC1123Z+" (MST)"))
	expectedEmail := Email{
		Headers: Headers{
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
			InReplyTo:  []MessageId{"Message-Id-0@example.com"},
			References: []MessageId{"Message-Id-0@example.com"},
			Comments:   "Message Header Comment",
			Keywords:   []string{"Keyword 1", "Keyword 2"},
			ResentDate: expectedDate,
			ResentFrom: []*mail.Address{
				{
					Name:    "Alice Sender",
					Address: "alice.sender@example.net",
				},
				{
					Name:    "Alice Sender",
					Address: "alice.sender@example.com",
				},
			},
			ResentSender: &mail.Address{
				Name:    "Alice Sender",
				Address: "alice.sender@example.net",
			},
			ResentTo: []*mail.Address{
				{
					Name:    "Bob Recipient",
					Address: "bob.recipient@example.net",
				},
				{
					Name:    "Carol Recipient",
					Address: "carol.recipient@example.net",
				},
			},
			ResentCc: []*mail.Address{
				{
					Name:    "Dan Recipient",
					Address: "dan.recipient@example.net",
				},
				{
					Name:    "Eve Recipient",
					Address: "eve.recipient@example.net",
				},
			},
			ResentBcc: []*mail.Address{
				{
					Name:    "Frank Recipient",
					Address: "frank.recipient@example.net",
				},
				{
					Name:    "Grace Recipient",
					Address: "grace.recipient@example.net",
				},
			},
			ResentMessageID: "Message-Id-1@example.net",
			ContentType: ContentTypeHeader{
				ContentType: "multipart/mixed",
				Params: map[string]string{
					"boundary": "MixedBoundaryString",
					"charset":  "utf-8",
				},
			},
			ExtraHeaders: map[string][]string{
				"X-Clacks-Overhead": {"GNU Terry Pratchett"},
				"X-Script/function/\t !\"#$%&'()*+,-./;<=>?@[\\]^_`{|}~": {
					"TEST VALUE 1\t !\"#$%&'()*+,-./0123456789:;<=>?@ABCDEFGHIJKLMNOPQRSTUVWXYZ[\\]^_` abcdefghijklmnopqrstuvwxyz{|}~",
					"TEST VALUE 2\t !\"#$%&'()*+,-./0123456789:;<=>?@ABCDEFGHIJKLMNOPQRSTUVWXYZ[\\]^_` abcdefghijklmnopqrstuvwxyz{|}~",
				},
			},
		},
		Text: `The quick brown fox jumps over a lazy dog.
Glib jocks quiz nymph to vex dwarf.
Sphinx of black quartz, judge my vow.
How vexingly quick daft zebras jump!
The five boxing wizards jump quickly.
Jackdaws love my big sphinx of quartz.
Pack my box with five dozen liquor jugs.`,
		EnrichedText: `<bold>The quick brown fox jumps over a lazy dog.</bold>
<italic>Glib jocks quiz nymph to vex dwarf.</italic>
<fixed>Sphinx of black quartz, judge my vow.</fixed>
<underline>How vexingly quick daft zebras jump!</underline>
The five boxing wizards jump quickly.
Jackdaws love my big sphinx of quartz.
Pack my box with five dozen liquor jugs.`,
		HTML: `<html>
<div dir="ltr">
<p>The quick brown fox jumps over a lazy dog.</p>
<p>Glib jocks quiz nymph to vex dwarf.</p>
<p>Sphinx of black quartz, judge my vow.</p>
<p>How vexingly quick daft zebras jump!</p>
<p>The five boxing wizards jump quickly.</p>
<p>Jackdaws love my big sphinx of quartz.</p>
<p>Pack my box with five dozen liquor jugs.</p>
</div>
</html>`,
		InlineFiles: []InlineFile{
			{
				ContentType: ContentTypeHeader{
					ContentType: "image/jpeg",
					Params: map[string]string{
						"name": "inline-jpg-image-without-disposition.jpg",
					},
				},
				ContentDisposition: ContentDispositionHeader{
					ContentDisposition: "",
					Params:             map[string]string(nil),
				},
				Data: []byte{255, 216, 255, 219, 0, 67, 0, 3, 2, 2, 2, 2, 2, 3, 2, 2, 2, 3, 3, 3, 3, 4,
					6, 4, 4, 4, 4, 4, 8, 6, 6, 5, 6, 9, 8, 10, 10, 9, 8, 9, 9, 10, 12, 15, 12, 10, 11, 14, 11, 9, 9,
					13, 17, 13, 14, 15, 16, 16, 17, 16, 10, 12, 18, 19, 18, 16, 19, 15, 16, 16, 16, 255, 201, 0, 11, 8,
					0, 1, 0, 1, 1, 1, 17, 0, 255, 204, 0, 6, 0, 16, 16, 5, 255, 218, 0, 8, 1, 1, 0, 0, 63, 0, 210, 207,
					32, 255, 217},
			},
			{
				ContentID: "inline-jpg-image.jpg@example.com",
				ContentType: ContentTypeHeader{
					ContentType: "image/jpeg",
					Params: map[string]string{
						"name": "inline-jpg-image-name.jpg",
					},
				},
				ContentDisposition: ContentDispositionHeader{
					ContentDisposition: inline,
					Params: map[string]string{
						"filename": "inline-jpg-image-filename.jpg",
					},
				},
				Data: []byte{255, 216, 255, 219, 0, 67, 0, 3, 2, 2, 2, 2, 2, 3, 2, 2, 2, 3, 3, 3, 3, 4,
					6, 4, 4, 4, 4, 4, 8, 6, 6, 5, 6, 9, 8, 10, 10, 9, 8, 9, 9, 10, 12, 15, 12, 10, 11, 14, 11, 9, 9,
					13, 17, 13, 14, 15, 16, 16, 17, 16, 10, 12, 18, 19, 18, 16, 19, 15, 16, 16, 16, 255, 201, 0, 11, 8,
					0, 1, 0, 1, 1, 1, 17, 0, 255, 204, 0, 6, 0, 16, 16, 5, 255, 218, 0, 8, 1, 1, 0, 0, 63, 0, 210, 207,
					32, 255, 217},
			},
		},
		AttachedFiles: []AttachedFile{
			{
				ContentType: ContentTypeHeader{
					ContentType: "application/pdf",
					Params: map[string]string{
						"name": "attached-pdf-name.pdf",
					},
				},
				ContentDisposition: ContentDispositionHeader{
					ContentDisposition: attachment,
					Params: map[string]string{
						"filename": "attached-pdf-filename.pdf",
					},
				},
				Data: []byte{37, 80, 68, 70, 45, 49, 46, 13, 116, 114, 97, 105, 108, 101, 114, 60, 60,
					47, 82, 111, 111, 116, 60, 60, 47, 80, 97, 103, 101, 115, 60, 60, 47, 75, 105, 100, 115, 91, 60,
					60, 47, 77, 101, 100, 105, 97, 66, 111, 120, 91, 48, 32, 48, 32, 51, 32, 51, 93, 62, 62, 93, 62,
					62, 62, 62, 62, 62},
			},
			{
				ContentType: ContentTypeHeader{
					ContentType: "application/pdf",
					Params: map[string]string{
						"name": "attached-pdf-without-disposition.pdf",
					},
				},
				ContentDisposition: ContentDispositionHeader{
					ContentDisposition: "",
					Params:             map[string]string(nil),
				},
				Data: []byte{37, 80, 68, 70, 45, 49, 46, 13, 116, 114, 97, 105, 108, 101, 114, 60, 60,
					47, 82, 111, 111, 116, 60, 60, 47, 80, 97, 103, 101, 115, 60, 60, 47, 75, 105, 100, 115, 91, 60,
					60, 47, 77, 101, 100, 105, 97, 66, 111, 120, 91, 48, 32, 48, 32, 51, 32, 51, 93, 62, 62, 93, 62,
					62, 62, 62, 62, 62},
			},
			{
				ContentType: ContentTypeHeader{
					ContentType: "application/json",
					Params: map[string]string{
						"name": "attached-json-name.json",
					},
				},
				ContentDisposition: ContentDispositionHeader{
					ContentDisposition: attachment,
					Params: map[string]string{
						"filename": "attached-json-filename.json",
					},
				},
				Data: []byte{123, 34, 102, 111, 111, 34, 58, 34, 98, 97, 114, 34, 125},
			},
			{
				ContentType: ContentTypeHeader{
					ContentType: "text/plain",
					Params: map[string]string{
						"name": "attached-text-plain-name.txt",
					},
				},
				ContentDisposition: ContentDispositionHeader{
					ContentDisposition: attachment,
					Params: map[string]string{
						"filename": "attached-text-plain-filename.txt",
					},
				},
				Data: []byte{84, 101, 120, 116, 47, 112, 108, 97, 105, 110, 32, 99, 111, 110, 116, 101, 110, 116, 32,
					97, 115, 32, 97, 110, 32, 97, 116, 116, 97, 99, 104, 101, 100, 32, 46, 116, 120, 116, 32, 102, 105,
					108, 101, 46},
			},
			{
				ContentType: ContentTypeHeader{
					ContentType: "text/html",
					Params: map[string]string{
						"name": "attached-text-html-name.html",
					},
				},
				ContentDisposition: ContentDispositionHeader{
					ContentDisposition: attachment,
					Params: map[string]string{
						"filename": "attached-text-html-filename.html",
					},
				},
				Data: []byte{84, 101, 120, 116, 47, 104, 116, 109, 108, 32, 99, 111, 110, 116, 101, 110, 116, 32, 97,
					115, 32, 97, 110, 32, 97, 116, 116, 97, 99, 104, 101, 100, 32, 46, 104, 116, 109, 108, 32, 102,
					105, 108, 101, 46},
			},
		},
	}

	testEmailFromFile(t, fp, expectedEmail)
}
func TestParseEmailEnglishMultipartSignedAsciiOver7bit(t *testing.T) {
	fp := "tests/test_english_multipart_signed_ascii_over_7bit.txt"
	tz, _ := time.LoadLocation("Europe/London")
	expectedDate, _ := time.Parse(
		time.RFC1123Z+" (MST)",
		time.Date(2019, time.April, 1, 7, 55, 0, 0, tz).Format(time.RFC1123Z+" (MST)"))
	expectedEmail := Email{
		Headers: Headers{
			Date:    expectedDate,
			Subject: "📧 Signed Test English Pangrams",
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
			InReplyTo:  []MessageId{"Message-Id-0@example.com"},
			References: []MessageId{"Message-Id-0@example.com"},
			Comments:   "Message Header Comment",
			Keywords:   []string{"Keyword 1", "Keyword 2"},
			ResentDate: expectedDate,
			ResentFrom: []*mail.Address{
				{
					Name:    "Alice Sender",
					Address: "alice.sender@example.net",
				},
				{
					Name:    "Alice Sender",
					Address: "alice.sender@example.com",
				},
			},
			ResentSender: &mail.Address{
				Name:    "Alice Sender",
				Address: "alice.sender@example.net",
			},
			ResentTo: []*mail.Address{
				{
					Name:    "Bob Recipient",
					Address: "bob.recipient@example.net",
				},
				{
					Name:    "Carol Recipient",
					Address: "carol.recipient@example.net",
				},
			},
			ResentCc: []*mail.Address{
				{
					Name:    "Dan Recipient",
					Address: "dan.recipient@example.net",
				},
				{
					Name:    "Eve Recipient",
					Address: "eve.recipient@example.net",
				},
			},
			ResentBcc: []*mail.Address{
				{
					Name:    "Frank Recipient",
					Address: "frank.recipient@example.net",
				},
				{
					Name:    "Grace Recipient",
					Address: "grace.recipient@example.net",
				},
			},
			ResentMessageID: "Message-Id-1@example.net",
			ContentType: ContentTypeHeader{
				ContentType: "multipart/signed",
				Params: map[string]string{
					"boundary": "SignedBoundaryString",
					"charset":  "ascii",
					"micalg":   "sha1",
					"protocol": "application/pkcs7-signature",
				},
			},
			ExtraHeaders: map[string][]string{
				"X-Clacks-Overhead": {"GNU Terry Pratchett"},
				"X-Script/function/\t !\"#$%&'()*+,-./;<=>?@[\\]^_`{|}~": {
					"TEST VALUE 1\t !\"#$%&'()*+,-./0123456789:;<=>?@ABCDEFGHIJKLMNOPQRSTUVWXYZ[\\]^_` abcdefghijklmnopqrstuvwxyz{|}~",
					"TEST VALUE 2\t !\"#$%&'()*+,-./0123456789:;<=>?@ABCDEFGHIJKLMNOPQRSTUVWXYZ[\\]^_` abcdefghijklmnopqrstuvwxyz{|}~",
				},
			},
		},
		Text: `The quick brown fox jumps over a lazy dog.
Glib jocks quiz nymph to vex dwarf.
Sphinx of black quartz, judge my vow.
How vexingly quick daft zebras jump!
The five boxing wizards jump quickly.
Jackdaws love my big sphinx of quartz.
Pack my box with five dozen liquor jugs.`,
		AttachedFiles: []AttachedFile{
			{
				ContentType: ContentTypeHeader{
					ContentType: "application/pkcs7-signature",
					Params: map[string]string{
						"name": "smime.p7s",
					},
				},
				ContentDisposition: ContentDispositionHeader{
					ContentDisposition: attachment,
					Params: map[string]string{
						"filename": "smime.p7s",
					},
				},
				Data: []byte{130, 28, 135, 132, 117, 46, 142, 18, 97, 140, 126, 251, 159, 193, 199, 25, 58, 223, 189, 185, 227, 239, 158, 173, 108, 31, 71, 27, 133, 80, 165, 252, 133, 227, 174, 198, 132, 129, 159, 29, 246, 19, 235, 133, 80, 165, 252, 133, 227, 174, 198, 132, 129, 159, 29, 246, 19, 234, 49, 251, 238, 127, 7, 28, 104, 33, 200, 120, 71, 82, 232, 225, 38, 30, 249, 234, 214, 193, 244, 113, 147, 173, 251, 219, 158, 57, 252, 28, 113, 147, 173, 251, 225, 38, 24, 199, 239, 190, 173, 108, 31, 71, 27, 133, 80, 110, 120, 251, 231, 174, 198, 132, 129, 159, 29, 246, 19, 234, 8, 114, 30, 17, 212, 186, 58, 95, 200, 94, 59, 26, 18, 6, 124, 119, 216, 79, 174, 21, 65, 185, 227, 239, 158},
			},
		},
	}

	testEmailFromFile(t, fp, expectedEmail)
}

func TestParseEmailEnglishMultipartSignedAsciiOverBase64(t *testing.T) {
	fp := "tests/test_english_multipart_signed_ascii_over_base64.txt"
	tz, _ := time.LoadLocation("Europe/London")
	expectedDate, _ := time.Parse(
		time.RFC1123Z+" (MST)",
		time.Date(2019, time.April, 1, 7, 55, 0, 0, tz).Format(time.RFC1123Z+" (MST)"))
	expectedEmail := Email{
		Headers: Headers{
			Date:    expectedDate,
			Subject: "📧 Signed Test English Pangrams",
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
			InReplyTo:  []MessageId{"Message-Id-0@example.com"},
			References: []MessageId{"Message-Id-0@example.com"},
			Comments:   "Message Header Comment",
			Keywords:   []string{"Keyword 1", "Keyword 2"},
			ResentDate: expectedDate,
			ResentFrom: []*mail.Address{
				{
					Name:    "Alice Sender",
					Address: "alice.sender@example.net",
				},
				{
					Name:    "Alice Sender",
					Address: "alice.sender@example.com",
				},
			},
			ResentSender: &mail.Address{
				Name:    "Alice Sender",
				Address: "alice.sender@example.net",
			},
			ResentTo: []*mail.Address{
				{
					Name:    "Bob Recipient",
					Address: "bob.recipient@example.net",
				},
				{
					Name:    "Carol Recipient",
					Address: "carol.recipient@example.net",
				},
			},
			ResentCc: []*mail.Address{
				{
					Name:    "Dan Recipient",
					Address: "dan.recipient@example.net",
				},
				{
					Name:    "Eve Recipient",
					Address: "eve.recipient@example.net",
				},
			},
			ResentBcc: []*mail.Address{
				{
					Name:    "Frank Recipient",
					Address: "frank.recipient@example.net",
				},
				{
					Name:    "Grace Recipient",
					Address: "grace.recipient@example.net",
				},
			},
			ResentMessageID: "Message-Id-1@example.net",
			ContentType: ContentTypeHeader{
				ContentType: "multipart/signed",
				Params: map[string]string{
					"boundary": "SignedBoundaryString",
					"charset":  "ascii",
					"micalg":   "sha1",
					"protocol": "application/pkcs7-signature",
				},
			},
			ExtraHeaders: map[string][]string{
				"X-Clacks-Overhead": {"GNU Terry Pratchett"},
				"X-Script/function/\t !\"#$%&'()*+,-./;<=>?@[\\]^_`{|}~": {
					"TEST VALUE 1\t !\"#$%&'()*+,-./0123456789:;<=>?@ABCDEFGHIJKLMNOPQRSTUVWXYZ[\\]^_` abcdefghijklmnopqrstuvwxyz{|}~",
					"TEST VALUE 2\t !\"#$%&'()*+,-./0123456789:;<=>?@ABCDEFGHIJKLMNOPQRSTUVWXYZ[\\]^_` abcdefghijklmnopqrstuvwxyz{|}~",
				},
			},
		},
		Text: `The quick brown fox jumps over a lazy dog.
Glib jocks quiz nymph to vex dwarf.
Sphinx of black quartz, judge my vow.
How vexingly quick daft zebras jump!
The five boxing wizards jump quickly.
Jackdaws love my big sphinx of quartz.
Pack my box with five dozen liquor jugs.`,
		AttachedFiles: []AttachedFile{
			{
				ContentType: ContentTypeHeader{
					ContentType: "application/pkcs7-signature",
					Params: map[string]string{
						"name": "smime.p7s",
					},
				},
				ContentDisposition: ContentDispositionHeader{
					ContentDisposition: attachment,
					Params: map[string]string{
						"filename": "smime.p7s",
					},
				},
				Data: []byte{130, 28, 135, 132, 117, 46, 142, 18, 97, 140, 126, 251, 159, 193, 199, 25, 58, 223, 189, 185, 227, 239, 158, 173, 108, 31, 71, 27, 133, 80, 165, 252, 133, 227, 174, 198, 132, 129, 159, 29, 246, 19, 235, 133, 80, 165, 252, 133, 227, 174, 198, 132, 129, 159, 29, 246, 19, 234, 49, 251, 238, 127, 7, 28, 104, 33, 200, 120, 71, 82, 232, 225, 38, 30, 249, 234, 214, 193, 244, 113, 147, 173, 251, 219, 158, 57, 252, 28, 113, 147, 173, 251, 225, 38, 24, 199, 239, 190, 173, 108, 31, 71, 27, 133, 80, 110, 120, 251, 231, 174, 198, 132, 129, 159, 29, 246, 19, 234, 8, 114, 30, 17, 212, 186, 58, 95, 200, 94, 59, 26, 18, 6, 124, 119, 216, 79, 174, 21, 65, 185, 227, 239, 158},
			},
		},
	}

	testEmailFromFile(t, fp, expectedEmail)
}

func TestParseEmailEnglishMultipartSignedAsciiOverQuotedprintable(t *testing.T) {
	fp := "tests/test_english_multipart_signed_ascii_over_quoted-printable.txt"
	tz, _ := time.LoadLocation("Europe/London")
	expectedDate, _ := time.Parse(
		time.RFC1123Z+" (MST)",
		time.Date(2019, time.April, 1, 7, 55, 0, 0, tz).Format(time.RFC1123Z+" (MST)"))
	expectedEmail := Email{
		Headers: Headers{
			Date:    expectedDate,
			Subject: "📧 Signed Test English Pangrams",
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
			InReplyTo:  []MessageId{"Message-Id-0@example.com"},
			References: []MessageId{"Message-Id-0@example.com"},
			Comments:   "Message Header Comment",
			Keywords:   []string{"Keyword 1", "Keyword 2"},
			ResentDate: expectedDate,
			ResentFrom: []*mail.Address{
				{
					Name:    "Alice Sender",
					Address: "alice.sender@example.net",
				},
				{
					Name:    "Alice Sender",
					Address: "alice.sender@example.com",
				},
			},
			ResentSender: &mail.Address{
				Name:    "Alice Sender",
				Address: "alice.sender@example.net",
			},
			ResentTo: []*mail.Address{
				{
					Name:    "Bob Recipient",
					Address: "bob.recipient@example.net",
				},
				{
					Name:    "Carol Recipient",
					Address: "carol.recipient@example.net",
				},
			},
			ResentCc: []*mail.Address{
				{
					Name:    "Dan Recipient",
					Address: "dan.recipient@example.net",
				},
				{
					Name:    "Eve Recipient",
					Address: "eve.recipient@example.net",
				},
			},
			ResentBcc: []*mail.Address{
				{
					Name:    "Frank Recipient",
					Address: "frank.recipient@example.net",
				},
				{
					Name:    "Grace Recipient",
					Address: "grace.recipient@example.net",
				},
			},
			ResentMessageID: "Message-Id-1@example.net",
			ContentType: ContentTypeHeader{
				ContentType: "multipart/signed",
				Params: map[string]string{
					"boundary": "SignedBoundaryString",
					"charset":  "ascii",
					"micalg":   "sha1",
					"protocol": "application/pkcs7-signature",
				},
			},
			ExtraHeaders: map[string][]string{
				"X-Clacks-Overhead": {"GNU Terry Pratchett"},
				"X-Script/function/\t !\"#$%&'()*+,-./;<=>?@[\\]^_`{|}~": {
					"TEST VALUE 1\t !\"#$%&'()*+,-./0123456789:;<=>?@ABCDEFGHIJKLMNOPQRSTUVWXYZ[\\]^_` abcdefghijklmnopqrstuvwxyz{|}~",
					"TEST VALUE 2\t !\"#$%&'()*+,-./0123456789:;<=>?@ABCDEFGHIJKLMNOPQRSTUVWXYZ[\\]^_` abcdefghijklmnopqrstuvwxyz{|}~",
				},
			},
		},
		Text: `The quick brown fox jumps over a lazy dog.
Glib jocks quiz nymph to vex dwarf.
Sphinx of black quartz, judge my vow.
How vexingly quick daft zebras jump!
The five boxing wizards jump quickly.
Jackdaws love my big sphinx of quartz.
Pack my box with five dozen liquor jugs.`,
		AttachedFiles: []AttachedFile{
			{
				ContentType: ContentTypeHeader{
					ContentType: "application/pkcs7-signature",
					Params: map[string]string{
						"name": "smime.p7s",
					},
				},
				ContentDisposition: ContentDispositionHeader{
					ContentDisposition: attachment,
					Params: map[string]string{
						"filename": "smime.p7s",
					},
				},
				Data: []byte{130, 28, 135, 132, 117, 46, 142, 18, 97, 140, 126, 251, 159, 193, 199, 25, 58, 223, 189, 185, 227, 239, 158, 173, 108, 31, 71, 27, 133, 80, 165, 252, 133, 227, 174, 198, 132, 129, 159, 29, 246, 19, 235, 133, 80, 165, 252, 133, 227, 174, 198, 132, 129, 159, 29, 246, 19, 234, 49, 251, 238, 127, 7, 28, 104, 33, 200, 120, 71, 82, 232, 225, 38, 30, 249, 234, 214, 193, 244, 113, 147, 173, 251, 219, 158, 57, 252, 28, 113, 147, 173, 251, 225, 38, 24, 199, 239, 190, 173, 108, 31, 71, 27, 133, 80, 110, 120, 251, 231, 174, 198, 132, 129, 159, 29, 246, 19, 234, 8, 114, 30, 17, 212, 186, 58, 95, 200, 94, 59, 26, 18, 6, 124, 119, 216, 79, 174, 21, 65, 185, 227, 239, 158},
			},
		},
	}

	testEmailFromFile(t, fp, expectedEmail)
}

func TestParseEmailEnglishMultipartSignedUtf8Over7bit(t *testing.T) {
	fp := "tests/test_english_multipart_signed_utf-8_over_7bit.txt"
	tz, _ := time.LoadLocation("Europe/London")
	expectedDate, _ := time.Parse(
		time.RFC1123Z+" (MST)",
		time.Date(2019, time.April, 1, 7, 55, 0, 0, tz).Format(time.RFC1123Z+" (MST)"))
	expectedEmail := Email{
		Headers: Headers{
			Date:    expectedDate,
			Subject: "📧 Signed Test English Pangrams",
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
			InReplyTo:  []MessageId{"Message-Id-0@example.com"},
			References: []MessageId{"Message-Id-0@example.com"},
			Comments:   "Message Header Comment",
			Keywords:   []string{"Keyword 1", "Keyword 2"},
			ResentDate: expectedDate,
			ResentFrom: []*mail.Address{
				{
					Name:    "Alice Sender",
					Address: "alice.sender@example.net",
				},
				{
					Name:    "Alice Sender",
					Address: "alice.sender@example.com",
				},
			},
			ResentSender: &mail.Address{
				Name:    "Alice Sender",
				Address: "alice.sender@example.net",
			},
			ResentTo: []*mail.Address{
				{
					Name:    "Bob Recipient",
					Address: "bob.recipient@example.net",
				},
				{
					Name:    "Carol Recipient",
					Address: "carol.recipient@example.net",
				},
			},
			ResentCc: []*mail.Address{
				{
					Name:    "Dan Recipient",
					Address: "dan.recipient@example.net",
				},
				{
					Name:    "Eve Recipient",
					Address: "eve.recipient@example.net",
				},
			},
			ResentBcc: []*mail.Address{
				{
					Name:    "Frank Recipient",
					Address: "frank.recipient@example.net",
				},
				{
					Name:    "Grace Recipient",
					Address: "grace.recipient@example.net",
				},
			},
			ResentMessageID: "Message-Id-1@example.net",
			ContentType: ContentTypeHeader{
				ContentType: "multipart/signed",
				Params: map[string]string{
					"boundary": "SignedBoundaryString",
					"charset":  "utf-8",
					"micalg":   "sha1",
					"protocol": "application/pkcs7-signature",
				},
			},
			ExtraHeaders: map[string][]string{
				"X-Clacks-Overhead": {"GNU Terry Pratchett"},
				"X-Script/function/\t !\"#$%&'()*+,-./;<=>?@[\\]^_`{|}~": {
					"TEST VALUE 1\t !\"#$%&'()*+,-./0123456789:;<=>?@ABCDEFGHIJKLMNOPQRSTUVWXYZ[\\]^_` abcdefghijklmnopqrstuvwxyz{|}~",
					"TEST VALUE 2\t !\"#$%&'()*+,-./0123456789:;<=>?@ABCDEFGHIJKLMNOPQRSTUVWXYZ[\\]^_` abcdefghijklmnopqrstuvwxyz{|}~",
				},
			},
		},
		Text: `The quick brown fox jumps over a lazy dog.
Glib jocks quiz nymph to vex dwarf.
Sphinx of black quartz, judge my vow.
How vexingly quick daft zebras jump!
The five boxing wizards jump quickly.
Jackdaws love my big sphinx of quartz.
Pack my box with five dozen liquor jugs.`,
		AttachedFiles: []AttachedFile{
			{
				ContentType: ContentTypeHeader{
					ContentType: "application/pkcs7-signature",
					Params: map[string]string{
						"name": "smime.p7s",
					},
				},
				ContentDisposition: ContentDispositionHeader{
					ContentDisposition: attachment,
					Params: map[string]string{
						"filename": "smime.p7s",
					},
				},
				Data: []byte{130, 28, 135, 132, 117, 46, 142, 18, 97, 140, 126, 251, 159, 193, 199, 25, 58, 223, 189, 185, 227, 239, 158, 173, 108, 31, 71, 27, 133, 80, 165, 252, 133, 227, 174, 198, 132, 129, 159, 29, 246, 19, 235, 133, 80, 165, 252, 133, 227, 174, 198, 132, 129, 159, 29, 246, 19, 234, 49, 251, 238, 127, 7, 28, 104, 33, 200, 120, 71, 82, 232, 225, 38, 30, 249, 234, 214, 193, 244, 113, 147, 173, 251, 219, 158, 57, 252, 28, 113, 147, 173, 251, 225, 38, 24, 199, 239, 190, 173, 108, 31, 71, 27, 133, 80, 110, 120, 251, 231, 174, 198, 132, 129, 159, 29, 246, 19, 234, 8, 114, 30, 17, 212, 186, 58, 95, 200, 94, 59, 26, 18, 6, 124, 119, 216, 79, 174, 21, 65, 185, 227, 239, 158},
			},
		},
	}

	testEmailFromFile(t, fp, expectedEmail)
}

func TestParseEmailEnglishMultipartSignedUtf8OverBase64(t *testing.T) {
	fp := "tests/test_english_multipart_signed_utf-8_over_base64.txt"
	tz, _ := time.LoadLocation("Europe/London")
	expectedDate, _ := time.Parse(
		time.RFC1123Z+" (MST)",
		time.Date(2019, time.April, 1, 7, 55, 0, 0, tz).Format(time.RFC1123Z+" (MST)"))
	expectedEmail := Email{
		Headers: Headers{
			Date:    expectedDate,
			Subject: "📧 Signed Test English Pangrams",
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
			InReplyTo:  []MessageId{"Message-Id-0@example.com"},
			References: []MessageId{"Message-Id-0@example.com"},
			Comments:   "Message Header Comment",
			Keywords:   []string{"Keyword 1", "Keyword 2"},
			ResentDate: expectedDate,
			ResentFrom: []*mail.Address{
				{
					Name:    "Alice Sender",
					Address: "alice.sender@example.net",
				},
				{
					Name:    "Alice Sender",
					Address: "alice.sender@example.com",
				},
			},
			ResentSender: &mail.Address{
				Name:    "Alice Sender",
				Address: "alice.sender@example.net",
			},
			ResentTo: []*mail.Address{
				{
					Name:    "Bob Recipient",
					Address: "bob.recipient@example.net",
				},
				{
					Name:    "Carol Recipient",
					Address: "carol.recipient@example.net",
				},
			},
			ResentCc: []*mail.Address{
				{
					Name:    "Dan Recipient",
					Address: "dan.recipient@example.net",
				},
				{
					Name:    "Eve Recipient",
					Address: "eve.recipient@example.net",
				},
			},
			ResentBcc: []*mail.Address{
				{
					Name:    "Frank Recipient",
					Address: "frank.recipient@example.net",
				},
				{
					Name:    "Grace Recipient",
					Address: "grace.recipient@example.net",
				},
			},
			ResentMessageID: "Message-Id-1@example.net",
			ContentType: ContentTypeHeader{
				ContentType: "multipart/signed",
				Params: map[string]string{
					"boundary": "SignedBoundaryString",
					"charset":  "utf-8",
					"micalg":   "sha1",
					"protocol": "application/pkcs7-signature",
				},
			},
			ExtraHeaders: map[string][]string{
				"X-Clacks-Overhead": {"GNU Terry Pratchett"},
				"X-Script/function/\t !\"#$%&'()*+,-./;<=>?@[\\]^_`{|}~": {
					"TEST VALUE 1\t !\"#$%&'()*+,-./0123456789:;<=>?@ABCDEFGHIJKLMNOPQRSTUVWXYZ[\\]^_` abcdefghijklmnopqrstuvwxyz{|}~",
					"TEST VALUE 2\t !\"#$%&'()*+,-./0123456789:;<=>?@ABCDEFGHIJKLMNOPQRSTUVWXYZ[\\]^_` abcdefghijklmnopqrstuvwxyz{|}~",
				},
			},
		},
		Text: `The quick brown fox jumps over a lazy dog.
Glib jocks quiz nymph to vex dwarf.
Sphinx of black quartz, judge my vow.
How vexingly quick daft zebras jump!
The five boxing wizards jump quickly.
Jackdaws love my big sphinx of quartz.
Pack my box with five dozen liquor jugs.`,
		AttachedFiles: []AttachedFile{
			{
				ContentType: ContentTypeHeader{
					ContentType: "application/pkcs7-signature",
					Params: map[string]string{
						"name": "smime.p7s",
					},
				},
				ContentDisposition: ContentDispositionHeader{
					ContentDisposition: attachment,
					Params: map[string]string{
						"filename": "smime.p7s",
					},
				},
				Data: []byte{130, 28, 135, 132, 117, 46, 142, 18, 97, 140, 126, 251, 159, 193, 199, 25, 58, 223, 189, 185, 227, 239, 158, 173, 108, 31, 71, 27, 133, 80, 165, 252, 133, 227, 174, 198, 132, 129, 159, 29, 246, 19, 235, 133, 80, 165, 252, 133, 227, 174, 198, 132, 129, 159, 29, 246, 19, 234, 49, 251, 238, 127, 7, 28, 104, 33, 200, 120, 71, 82, 232, 225, 38, 30, 249, 234, 214, 193, 244, 113, 147, 173, 251, 219, 158, 57, 252, 28, 113, 147, 173, 251, 225, 38, 24, 199, 239, 190, 173, 108, 31, 71, 27, 133, 80, 110, 120, 251, 231, 174, 198, 132, 129, 159, 29, 246, 19, 234, 8, 114, 30, 17, 212, 186, 58, 95, 200, 94, 59, 26, 18, 6, 124, 119, 216, 79, 174, 21, 65, 185, 227, 239, 158},
			},
		},
	}

	testEmailFromFile(t, fp, expectedEmail)
}

func TestParseEmailEnglishMultipartSignedUtf8OverQuotedprintable(t *testing.T) {
	fp := "tests/test_english_multipart_signed_utf-8_over_quoted-printable.txt"
	tz, _ := time.LoadLocation("Europe/London")
	expectedDate, _ := time.Parse(
		time.RFC1123Z+" (MST)",
		time.Date(2019, time.April, 1, 7, 55, 0, 0, tz).Format(time.RFC1123Z+" (MST)"))
	expectedEmail := Email{
		Headers: Headers{
			Date:    expectedDate,
			Subject: "📧 Signed Test English Pangrams",
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
			InReplyTo:  []MessageId{"Message-Id-0@example.com"},
			References: []MessageId{"Message-Id-0@example.com"},
			Comments:   "Message Header Comment",
			Keywords:   []string{"Keyword 1", "Keyword 2"},
			ResentDate: expectedDate,
			ResentFrom: []*mail.Address{
				{
					Name:    "Alice Sender",
					Address: "alice.sender@example.net",
				},
				{
					Name:    "Alice Sender",
					Address: "alice.sender@example.com",
				},
			},
			ResentSender: &mail.Address{
				Name:    "Alice Sender",
				Address: "alice.sender@example.net",
			},
			ResentTo: []*mail.Address{
				{
					Name:    "Bob Recipient",
					Address: "bob.recipient@example.net",
				},
				{
					Name:    "Carol Recipient",
					Address: "carol.recipient@example.net",
				},
			},
			ResentCc: []*mail.Address{
				{
					Name:    "Dan Recipient",
					Address: "dan.recipient@example.net",
				},
				{
					Name:    "Eve Recipient",
					Address: "eve.recipient@example.net",
				},
			},
			ResentBcc: []*mail.Address{
				{
					Name:    "Frank Recipient",
					Address: "frank.recipient@example.net",
				},
				{
					Name:    "Grace Recipient",
					Address: "grace.recipient@example.net",
				},
			},
			ResentMessageID: "Message-Id-1@example.net",
			ContentType: ContentTypeHeader{
				ContentType: "multipart/signed",
				Params: map[string]string{
					"boundary": "SignedBoundaryString",
					"charset":  "utf-8",
					"micalg":   "sha1",
					"protocol": "application/pkcs7-signature",
				},
			},
			ExtraHeaders: map[string][]string{
				"X-Clacks-Overhead": {"GNU Terry Pratchett"},
				"X-Script/function/\t !\"#$%&'()*+,-./;<=>?@[\\]^_`{|}~": {
					"TEST VALUE 1\t !\"#$%&'()*+,-./0123456789:;<=>?@ABCDEFGHIJKLMNOPQRSTUVWXYZ[\\]^_` abcdefghijklmnopqrstuvwxyz{|}~",
					"TEST VALUE 2\t !\"#$%&'()*+,-./0123456789:;<=>?@ABCDEFGHIJKLMNOPQRSTUVWXYZ[\\]^_` abcdefghijklmnopqrstuvwxyz{|}~",
				},
			},
		},
		Text: `The quick brown fox jumps over a lazy dog.
Glib jocks quiz nymph to vex dwarf.
Sphinx of black quartz, judge my vow.
How vexingly quick daft zebras jump!
The five boxing wizards jump quickly.
Jackdaws love my big sphinx of quartz.
Pack my box with five dozen liquor jugs.`,
		AttachedFiles: []AttachedFile{
			{
				ContentType: ContentTypeHeader{
					ContentType: "application/pkcs7-signature",
					Params: map[string]string{
						"name": "smime.p7s",
					},
				},
				ContentDisposition: ContentDispositionHeader{
					ContentDisposition: attachment,
					Params: map[string]string{
						"filename": "smime.p7s",
					},
				},
				Data: []byte{130, 28, 135, 132, 117, 46, 142, 18, 97, 140, 126, 251, 159, 193, 199, 25, 58, 223, 189, 185, 227, 239, 158, 173, 108, 31, 71, 27, 133, 80, 165, 252, 133, 227, 174, 198, 132, 129, 159, 29, 246, 19, 235, 133, 80, 165, 252, 133, 227, 174, 198, 132, 129, 159, 29, 246, 19, 234, 49, 251, 238, 127, 7, 28, 104, 33, 200, 120, 71, 82, 232, 225, 38, 30, 249, 234, 214, 193, 244, 113, 147, 173, 251, 219, 158, 57, 252, 28, 113, 147, 173, 251, 225, 38, 24, 199, 239, 190, 173, 108, 31, 71, 27, 133, 80, 110, 120, 251, 231, 174, 198, 132, 129, 159, 29, 246, 19, 234, 8, 114, 30, 17, 212, 186, 58, 95, 200, 94, 59, 26, 18, 6, 124, 119, 216, 79, 174, 21, 65, 185, 227, 239, 158},
			},
		},
	}

	testEmailFromFile(t, fp, expectedEmail)
}

func TestParseEmailChinesePlaintextGb18030OverBase64(t *testing.T) {
	fp := "tests/test_chinese_plaintext_gb18030_over_base64.txt"
	tz, _ := time.LoadLocation("Asia/Shanghai")
	expectedDate, _ := time.Parse(
		time.RFC1123Z+" (MST)",
		time.Date(2019, time.April, 1, 7, 55, 0, 0, tz).Format(time.RFC1123Z+" (MST)"))
	expectedEmail := Email{
		Headers: Headers{
			Date:    expectedDate,
			Subject: "📧 Test 施氏食狮史",
			ReplyTo: []*mail.Address{
				{
					Name:    "艾莉絲 发件人",
					Address: "alice.fajianren@example.net",
				},
			},
			Sender: &mail.Address{
				Name:    "艾莉絲 发件人",
				Address: "alice.fajianren@example.com",
			},
			From: []*mail.Address{
				{
					Name:    "艾莉絲 发件人",
					Address: "alice.fajianren@example.com",
				},
				{
					Name:    "艾莉絲 发件人",
					Address: "alice.fajianren@example.net",
				},
			},
			To: []*mail.Address{
				{
					Name:    "與鮑伯 收件人",
					Address: "bob.shoujianren@example.com",
				},
				{
					Name:    "卡罗尔 收件人",
					Address: "carol.shoujianren@example.com",
				},
			},
			Cc: []*mail.Address{
				{
					Name:    "戴夫 收件人",
					Address: "dave.shoujianren@example.com",
				},
				{
					Name:    "伊夫 收件人",
					Address: "eve.shoujianren@example.com",
				},
			},
			Bcc: []*mail.Address{
				{
					Name:    "艾萨克 伊夫",
					Address: "isaac.shoujianren@example.com",
				},
				{
					Name:    "賈斯汀 伊夫",
					Address: "justin.shoujianren@example.com",
				},
			},
			MessageID:  "Message-Id-1@example.com",
			InReplyTo:  []MessageId{"Message-Id-0@example.com"},
			References: []MessageId{"Message-Id-0@example.com"},
			Comments:   "Message Header Comment",
			Keywords:   []string{"Keyword 1", "Keyword 2"},
			ResentDate: expectedDate,
			ResentFrom: []*mail.Address{
				{
					Name:    "艾莉絲 发件人",
					Address: "alice.fajianren@example.net",
				},
				{
					Name:    "艾莉絲 发件人",
					Address: "alice.fajianren@example.com",
				},
			},
			ResentSender: &mail.Address{
				Name:    "艾莉絲 发件人",
				Address: "alice.fajianren@example.net",
			},
			ResentTo: []*mail.Address{
				{
					Name:    "與鮑伯 收件人",
					Address: "bob.shoujianren@example.net",
				},
				{
					Name:    "卡罗尔 收件人",
					Address: "carol.shoujianren@example.net",
				},
			},
			ResentCc: []*mail.Address{
				{
					Name:    "戴夫 收件人",
					Address: "dave.shoujianren@example.net",
				},
				{
					Name:    "伊夫 收件人",
					Address: "eve.shoujianren@example.net",
				},
			},
			ResentBcc: []*mail.Address{
				{
					Name:    "艾萨克 伊夫",
					Address: "isaac.shoujianren@example.net",
				},
				{
					Name:    "賈斯汀 伊夫",
					Address: "justin.shoujianren@example.net",
				},
			},
			ResentMessageID: "Message-Id-1@example.net",
			ContentType: ContentTypeHeader{
				ContentType: "text/plain",
				Params: map[string]string{
					"charset": "gb18030",
				},
			},
			ExtraHeaders: map[string][]string{
				"X-Clacks-Overhead": {"GNU Terry Pratchett"},
				"X-Script/function/\t !\"#$%&'()*+,-./;<=>?@[\\]^_`{|}~": {
					"TEST VALUE 1\t !\"#$%&'()*+,-./0123456789:;<=>?@ABCDEFGHIJKLMNOPQRSTUVWXYZ[\\]^_` abcdefghijklmnopqrstuvwxyz{|}~",
					"TEST VALUE 2\t !\"#$%&'()*+,-./0123456789:;<=>?@ABCDEFGHIJKLMNOPQRSTUVWXYZ[\\]^_` abcdefghijklmnopqrstuvwxyz{|}~",
				},
			},
		},
		Text: `石室诗士施氏，嗜狮，誓食十狮。
氏时时适市视狮。
十时，适十狮适市。
是时，适施氏适市。
氏视是十狮，恃矢势，使是十狮逝世。
氏拾是十狮尸，适石室。
石室湿，氏使侍拭石室。
石室拭，氏始试食是十狮。
食时，始识是十狮尸，实十石狮尸。
试释是事。`,
	}

	testEmailFromFile(t, fp, expectedEmail)
}

func TestParseEmailChinesePlaintextGb18030OverQuotedprintable(t *testing.T) {
	fp := "tests/test_chinese_plaintext_gb18030_over_quoted-printable.txt"
	tz, _ := time.LoadLocation("Asia/Shanghai")
	expectedDate, _ := time.Parse(
		time.RFC1123Z+" (MST)",
		time.Date(2019, time.April, 1, 7, 55, 0, 0, tz).Format(time.RFC1123Z+" (MST)"))
	expectedEmail := Email{
		Headers: Headers{
			Date:    expectedDate,
			Subject: "📧 Test 施氏食狮史",
			ReplyTo: []*mail.Address{
				{
					Name:    "艾莉絲 发件人",
					Address: "alice.fajianren@example.net",
				},
			},
			Sender: &mail.Address{
				Name:    "艾莉絲 发件人",
				Address: "alice.fajianren@example.com",
			},
			From: []*mail.Address{
				{
					Name:    "艾莉絲 发件人",
					Address: "alice.fajianren@example.com",
				},
				{
					Name:    "艾莉絲 发件人",
					Address: "alice.fajianren@example.net",
				},
			},
			To: []*mail.Address{
				{
					Name:    "與鮑伯 收件人",
					Address: "bob.shoujianren@example.com",
				},
				{
					Name:    "卡罗尔 收件人",
					Address: "carol.shoujianren@example.com",
				},
			},
			Cc: []*mail.Address{
				{
					Name:    "戴夫 收件人",
					Address: "dave.shoujianren@example.com",
				},
				{
					Name:    "伊夫 收件人",
					Address: "eve.shoujianren@example.com",
				},
			},
			Bcc: []*mail.Address{
				{
					Name:    "艾萨克 伊夫",
					Address: "isaac.shoujianren@example.com",
				},
				{
					Name:    "賈斯汀 伊夫",
					Address: "justin.shoujianren@example.com",
				},
			},
			MessageID:  "Message-Id-1@example.com",
			InReplyTo:  []MessageId{"Message-Id-0@example.com"},
			References: []MessageId{"Message-Id-0@example.com"},
			Comments:   "Message Header Comment",
			Keywords:   []string{"Keyword 1", "Keyword 2"},
			ResentDate: expectedDate,
			ResentFrom: []*mail.Address{
				{
					Name:    "艾莉絲 发件人",
					Address: "alice.fajianren@example.net",
				},
				{
					Name:    "艾莉絲 发件人",
					Address: "alice.fajianren@example.com",
				},
			},
			ResentSender: &mail.Address{
				Name:    "艾莉絲 发件人",
				Address: "alice.fajianren@example.net",
			},
			ResentTo: []*mail.Address{
				{
					Name:    "與鮑伯 收件人",
					Address: "bob.shoujianren@example.net",
				},
				{
					Name:    "卡罗尔 收件人",
					Address: "carol.shoujianren@example.net",
				},
			},
			ResentCc: []*mail.Address{
				{
					Name:    "戴夫 收件人",
					Address: "dave.shoujianren@example.net",
				},
				{
					Name:    "伊夫 收件人",
					Address: "eve.shoujianren@example.net",
				},
			},
			ResentBcc: []*mail.Address{
				{
					Name:    "艾萨克 伊夫",
					Address: "isaac.shoujianren@example.net",
				},
				{
					Name:    "賈斯汀 伊夫",
					Address: "justin.shoujianren@example.net",
				},
			},
			ResentMessageID: "Message-Id-1@example.net",
			ContentType: ContentTypeHeader{
				ContentType: "text/plain",
				Params: map[string]string{
					"charset": "gb18030",
				},
			},
			ExtraHeaders: map[string][]string{
				"X-Clacks-Overhead": {"GNU Terry Pratchett"},
				"X-Script/function/\t !\"#$%&'()*+,-./;<=>?@[\\]^_`{|}~": {
					"TEST VALUE 1\t !\"#$%&'()*+,-./0123456789:;<=>?@ABCDEFGHIJKLMNOPQRSTUVWXYZ[\\]^_` abcdefghijklmnopqrstuvwxyz{|}~",
					"TEST VALUE 2\t !\"#$%&'()*+,-./0123456789:;<=>?@ABCDEFGHIJKLMNOPQRSTUVWXYZ[\\]^_` abcdefghijklmnopqrstuvwxyz{|}~",
				},
			},
		},
		Text: `石室诗士施氏，嗜狮，誓食十狮。
氏时时适市视狮。
十时，适十狮适市。
是时，适施氏适市。
氏视是十狮，恃矢势，使是十狮逝世。
氏拾是十狮尸，适石室。
石室湿，氏使侍拭石室。
石室拭，氏始试食是十狮。
食时，始识是十狮尸，实十石狮尸。
试释是事。`,
	}

	testEmailFromFile(t, fp, expectedEmail)
}

func TestParseEmailChinesePlaintextGbkOverBase64(t *testing.T) {
	fp := "tests/test_chinese_plaintext_gbk_over_base64.txt"
	tz, _ := time.LoadLocation("Asia/Shanghai")
	expectedDate, _ := time.Parse(
		time.RFC1123Z+" (MST)",
		time.Date(2019, time.April, 1, 7, 55, 0, 0, tz).Format(time.RFC1123Z+" (MST)"))
	expectedEmail := Email{
		Headers: Headers{
			Date:    expectedDate,
			Subject: "📧 Test 施氏食狮史",
			ReplyTo: []*mail.Address{
				{
					Name:    "艾莉絲 发件人",
					Address: "alice.fajianren@example.net",
				},
			},
			Sender: &mail.Address{
				Name:    "艾莉絲 发件人",
				Address: "alice.fajianren@example.com",
			},
			From: []*mail.Address{
				{
					Name:    "艾莉絲 发件人",
					Address: "alice.fajianren@example.com",
				},
				{
					Name:    "艾莉絲 发件人",
					Address: "alice.fajianren@example.net",
				},
			},
			To: []*mail.Address{
				{
					Name:    "與鮑伯 收件人",
					Address: "bob.shoujianren@example.com",
				},
				{
					Name:    "卡罗尔 收件人",
					Address: "carol.shoujianren@example.com",
				},
			},
			Cc: []*mail.Address{
				{
					Name:    "戴夫 收件人",
					Address: "dave.shoujianren@example.com",
				},
				{
					Name:    "伊夫 收件人",
					Address: "eve.shoujianren@example.com",
				},
			},
			Bcc: []*mail.Address{
				{
					Name:    "艾萨克 伊夫",
					Address: "isaac.shoujianren@example.com",
				},
				{
					Name:    "賈斯汀 伊夫",
					Address: "justin.shoujianren@example.com",
				},
			},
			MessageID:  "Message-Id-1@example.com",
			InReplyTo:  []MessageId{"Message-Id-0@example.com"},
			References: []MessageId{"Message-Id-0@example.com"},
			Comments:   "Message Header Comment",
			Keywords:   []string{"Keyword 1", "Keyword 2"},
			ResentDate: expectedDate,
			ResentFrom: []*mail.Address{
				{
					Name:    "艾莉絲 发件人",
					Address: "alice.fajianren@example.net",
				},
				{
					Name:    "艾莉絲 发件人",
					Address: "alice.fajianren@example.com",
				},
			},
			ResentSender: &mail.Address{
				Name:    "艾莉絲 发件人",
				Address: "alice.fajianren@example.net",
			},
			ResentTo: []*mail.Address{
				{
					Name:    "與鮑伯 收件人",
					Address: "bob.shoujianren@example.net",
				},
				{
					Name:    "卡罗尔 收件人",
					Address: "carol.shoujianren@example.net",
				},
			},
			ResentCc: []*mail.Address{
				{
					Name:    "戴夫 收件人",
					Address: "dave.shoujianren@example.net",
				},
				{
					Name:    "伊夫 收件人",
					Address: "eve.shoujianren@example.net",
				},
			},
			ResentBcc: []*mail.Address{
				{
					Name:    "艾萨克 伊夫",
					Address: "isaac.shoujianren@example.net",
				},
				{
					Name:    "賈斯汀 伊夫",
					Address: "justin.shoujianren@example.net",
				},
			},
			ResentMessageID: "Message-Id-1@example.net",
			ContentType: ContentTypeHeader{
				ContentType: "text/plain",
				Params: map[string]string{
					"charset": "gbk",
				},
			},
			ExtraHeaders: map[string][]string{
				"X-Clacks-Overhead": {"GNU Terry Pratchett"},
				"X-Script/function/\t !\"#$%&'()*+,-./;<=>?@[\\]^_`{|}~": {
					"TEST VALUE 1\t !\"#$%&'()*+,-./0123456789:;<=>?@ABCDEFGHIJKLMNOPQRSTUVWXYZ[\\]^_` abcdefghijklmnopqrstuvwxyz{|}~",
					"TEST VALUE 2\t !\"#$%&'()*+,-./0123456789:;<=>?@ABCDEFGHIJKLMNOPQRSTUVWXYZ[\\]^_` abcdefghijklmnopqrstuvwxyz{|}~",
				},
			},
		},
		Text: `石室诗士施氏，嗜狮，誓食十狮。
氏时时适市视狮。
十时，适十狮适市。
是时，适施氏适市。
氏视是十狮，恃矢势，使是十狮逝世。
氏拾是十狮尸，适石室。
石室湿，氏使侍拭石室。
石室拭，氏始试食是十狮。
食时，始识是十狮尸，实十石狮尸。
试释是事。`,
	}

	testEmailFromFile(t, fp, expectedEmail)
}

func TestParseEmailChinesePlaintextGbkOverQuotedprintable(t *testing.T) {
	fp := "tests/test_chinese_plaintext_gbk_over_quoted-printable.txt"
	tz, _ := time.LoadLocation("Asia/Shanghai")
	expectedDate, _ := time.Parse(
		time.RFC1123Z+" (MST)",
		time.Date(2019, time.April, 1, 7, 55, 0, 0, tz).Format(time.RFC1123Z+" (MST)"))
	expectedEmail := Email{
		Headers: Headers{
			Date:    expectedDate,
			Subject: "📧 Test 施氏食狮史",
			ReplyTo: []*mail.Address{
				{
					Name:    "艾莉絲 发件人",
					Address: "alice.fajianren@example.net",
				},
			},
			Sender: &mail.Address{
				Name:    "艾莉絲 发件人",
				Address: "alice.fajianren@example.com",
			},
			From: []*mail.Address{
				{
					Name:    "艾莉絲 发件人",
					Address: "alice.fajianren@example.com",
				},
				{
					Name:    "艾莉絲 发件人",
					Address: "alice.fajianren@example.net",
				},
			},
			To: []*mail.Address{
				{
					Name:    "與鮑伯 收件人",
					Address: "bob.shoujianren@example.com",
				},
				{
					Name:    "卡罗尔 收件人",
					Address: "carol.shoujianren@example.com",
				},
			},
			Cc: []*mail.Address{
				{
					Name:    "戴夫 收件人",
					Address: "dave.shoujianren@example.com",
				},
				{
					Name:    "伊夫 收件人",
					Address: "eve.shoujianren@example.com",
				},
			},
			Bcc: []*mail.Address{
				{
					Name:    "艾萨克 伊夫",
					Address: "isaac.shoujianren@example.com",
				},
				{
					Name:    "賈斯汀 伊夫",
					Address: "justin.shoujianren@example.com",
				},
			},
			MessageID:  "Message-Id-1@example.com",
			InReplyTo:  []MessageId{"Message-Id-0@example.com"},
			References: []MessageId{"Message-Id-0@example.com"},
			Comments:   "Message Header Comment",
			Keywords:   []string{"Keyword 1", "Keyword 2"},
			ResentDate: expectedDate,
			ResentFrom: []*mail.Address{
				{
					Name:    "艾莉絲 发件人",
					Address: "alice.fajianren@example.net",
				},
				{
					Name:    "艾莉絲 发件人",
					Address: "alice.fajianren@example.com",
				},
			},
			ResentSender: &mail.Address{
				Name:    "艾莉絲 发件人",
				Address: "alice.fajianren@example.net",
			},
			ResentTo: []*mail.Address{
				{
					Name:    "與鮑伯 收件人",
					Address: "bob.shoujianren@example.net",
				},
				{
					Name:    "卡罗尔 收件人",
					Address: "carol.shoujianren@example.net",
				},
			},
			ResentCc: []*mail.Address{
				{
					Name:    "戴夫 收件人",
					Address: "dave.shoujianren@example.net",
				},
				{
					Name:    "伊夫 收件人",
					Address: "eve.shoujianren@example.net",
				},
			},
			ResentBcc: []*mail.Address{
				{
					Name:    "艾萨克 伊夫",
					Address: "isaac.shoujianren@example.net",
				},
				{
					Name:    "賈斯汀 伊夫",
					Address: "justin.shoujianren@example.net",
				},
			},
			ResentMessageID: "Message-Id-1@example.net",
			ContentType: ContentTypeHeader{
				ContentType: "text/plain",
				Params: map[string]string{
					"charset": "gbk",
				},
			},
			ExtraHeaders: map[string][]string{
				"X-Clacks-Overhead": {"GNU Terry Pratchett"},
				"X-Script/function/\t !\"#$%&'()*+,-./;<=>?@[\\]^_`{|}~": {
					"TEST VALUE 1\t !\"#$%&'()*+,-./0123456789:;<=>?@ABCDEFGHIJKLMNOPQRSTUVWXYZ[\\]^_` abcdefghijklmnopqrstuvwxyz{|}~",
					"TEST VALUE 2\t !\"#$%&'()*+,-./0123456789:;<=>?@ABCDEFGHIJKLMNOPQRSTUVWXYZ[\\]^_` abcdefghijklmnopqrstuvwxyz{|}~",
				},
			},
		},
		Text: `石室诗士施氏，嗜狮，誓食十狮。
氏时时适市视狮。
十时，适十狮适市。
是时，适施氏适市。
氏视是十狮，恃矢势，使是十狮逝世。
氏拾是十狮尸，适石室。
石室湿，氏使侍拭石室。
石室拭，氏始试食是十狮。
食时，始识是十狮尸，实十石狮尸。
试释是事。`,
	}

	testEmailFromFile(t, fp, expectedEmail)
}

func TestParseEmailChineseMultipartRelatedGb18030OverBase64(t *testing.T) {
	fp := "tests/test_chinese_multipart_related_gb18030_over_base64.txt"
	tz, _ := time.LoadLocation("Asia/Shanghai")
	expectedDate, _ := time.Parse(
		time.RFC1123Z+" (MST)",
		time.Date(2019, time.April, 1, 7, 55, 0, 0, tz).Format(time.RFC1123Z+" (MST)"))
	expectedEmail := Email{
		Headers: Headers{
			Date:    expectedDate,
			Subject: "📧 Test 施氏食狮史",
			ReplyTo: []*mail.Address{
				{
					Name:    "艾莉絲 发件人",
					Address: "alice.fajianren@example.net",
				},
			},
			Sender: &mail.Address{
				Name:    "艾莉絲 发件人",
				Address: "alice.fajianren@example.com",
			},
			From: []*mail.Address{
				{
					Name:    "艾莉絲 发件人",
					Address: "alice.fajianren@example.com",
				},
				{
					Name:    "艾莉絲 发件人",
					Address: "alice.fajianren@example.net",
				},
			},
			To: []*mail.Address{
				{
					Name:    "與鮑伯 收件人",
					Address: "bob.shoujianren@example.com",
				},
				{
					Name:    "卡罗尔 收件人",
					Address: "carol.shoujianren@example.com",
				},
			},
			Cc: []*mail.Address{
				{
					Name:    "戴夫 收件人",
					Address: "dave.shoujianren@example.com",
				},
				{
					Name:    "伊夫 收件人",
					Address: "eve.shoujianren@example.com",
				},
			},
			Bcc: []*mail.Address{
				{
					Name:    "艾萨克 伊夫",
					Address: "isaac.shoujianren@example.com",
				},
				{
					Name:    "賈斯汀 伊夫",
					Address: "justin.shoujianren@example.com",
				},
			},
			MessageID:  "Message-Id-1@example.com",
			InReplyTo:  []MessageId{"Message-Id-0@example.com"},
			References: []MessageId{"Message-Id-0@example.com"},
			Comments:   "Message Header Comment",
			Keywords:   []string{"Keyword 1", "Keyword 2"},
			ResentDate: expectedDate,
			ResentFrom: []*mail.Address{
				{
					Name:    "艾莉絲 发件人",
					Address: "alice.fajianren@example.net",
				},
				{
					Name:    "艾莉絲 发件人",
					Address: "alice.fajianren@example.com",
				},
			},
			ResentSender: &mail.Address{
				Name:    "艾莉絲 发件人",
				Address: "alice.fajianren@example.net",
			},
			ResentTo: []*mail.Address{
				{
					Name:    "與鮑伯 收件人",
					Address: "bob.shoujianren@example.net",
				},
				{
					Name:    "卡罗尔 收件人",
					Address: "carol.shoujianren@example.net",
				},
			},
			ResentCc: []*mail.Address{
				{
					Name:    "戴夫 收件人",
					Address: "dave.shoujianren@example.net",
				},
				{
					Name:    "伊夫 收件人",
					Address: "eve.shoujianren@example.net",
				},
			},
			ResentBcc: []*mail.Address{
				{
					Name:    "艾萨克 伊夫",
					Address: "isaac.shoujianren@example.net",
				},
				{
					Name:    "賈斯汀 伊夫",
					Address: "justin.shoujianren@example.net",
				},
			},
			ResentMessageID: "Message-Id-1@example.net",
			ContentType: ContentTypeHeader{
				ContentType: "multipart/related",
				Params: map[string]string{
					"boundary": "RelatedBoundaryString",
					"charset":  "gb18030",
				},
			},
			ExtraHeaders: map[string][]string{
				"X-Clacks-Overhead": {"GNU Terry Pratchett"},
				"X-Script/function/\t !\"#$%&'()*+,-./;<=>?@[\\]^_`{|}~": {
					"TEST VALUE 1\t !\"#$%&'()*+,-./0123456789:;<=>?@ABCDEFGHIJKLMNOPQRSTUVWXYZ[\\]^_` abcdefghijklmnopqrstuvwxyz{|}~",
					"TEST VALUE 2\t !\"#$%&'()*+,-./0123456789:;<=>?@ABCDEFGHIJKLMNOPQRSTUVWXYZ[\\]^_` abcdefghijklmnopqrstuvwxyz{|}~",
				},
			},
		},
		Text: `石室诗士施氏，嗜狮，誓食十狮。
氏时时适市视狮。
十时，适十狮适市。
是时，适施氏适市。
氏视是十狮，恃矢势，使是十狮逝世。
氏拾是十狮尸，适石室。
石室湿，氏使侍拭石室。
石室拭，氏始试食是十狮。
食时，始识是十狮尸，实十石狮尸。
试释是事。`,
		EnrichedText: `<bold>石室诗士施氏，嗜狮，誓食十狮。</bold>
<italic>氏时时适市视狮。</italic>
<fixed>十时，适十狮适市。</fixed>
<underline>是时，适施氏适市。</underline>
氏视是十狮，恃矢势，使是十狮逝世。
氏拾是十狮尸，适石室。
石室湿，氏使侍拭石室。
石室拭，氏始试食是十狮。
食时，始识是十狮尸，实十石狮尸。
试释是事。`,
		HTML: `<html>
<div dir="ltr">
<p>石室诗士施氏，嗜狮，誓食十狮。<br />
氏时时适市视狮。<br />
十时，适十狮适市。<br />
是时，适施氏适市。<br />
氏视是十狮，恃矢势，使是十狮逝世。<br />
氏拾是十狮尸，适石室。<br />
石室湿，氏使侍拭石室。<br />
石室拭，氏始试食是十狮。<br />
食时，始识是十狮尸，实十石狮尸。<br />
试释是事。</p>
</div>
</html>`,
		InlineFiles: []InlineFile{
			{
				ContentID: "inline-jpg-image.jpg@example.com",
				ContentType: ContentTypeHeader{
					ContentType: "image/jpeg",
					Params: map[string]string{
						"name": "inline-jpg-image-name.jpg",
					},
				},
				ContentDisposition: ContentDispositionHeader{
					ContentDisposition: inline,
					Params: map[string]string{
						"filename": "inline-jpg-image-filename.jpg",
					},
				},
				Data: []byte{255, 216, 255, 219, 0, 67, 0, 3, 2, 2, 2, 2, 2, 3, 2, 2, 2, 3, 3, 3, 3, 4,
					6, 4, 4, 4, 4, 4, 8, 6, 6, 5, 6, 9, 8, 10, 10, 9, 8, 9, 9, 10, 12, 15, 12, 10, 11, 14, 11, 9, 9,
					13, 17, 13, 14, 15, 16, 16, 17, 16, 10, 12, 18, 19, 18, 16, 19, 15, 16, 16, 16, 255, 201, 0, 11, 8,
					0, 1, 0, 1, 1, 1, 17, 0, 255, 204, 0, 6, 0, 16, 16, 5, 255, 218, 0, 8, 1, 1, 0, 0, 63, 0, 210, 207,
					32, 255, 217},
			},
		},
	}

	testEmailFromFile(t, fp, expectedEmail)
}
func TestParseEmailChineseMultipartRelatedGb18030OverQuotedprintable(t *testing.T) {
	fp := "tests/test_chinese_multipart_related_gb18030_over_quoted-printable.txt"
	tz, _ := time.LoadLocation("Asia/Shanghai")
	expectedDate, _ := time.Parse(
		time.RFC1123Z+" (MST)",
		time.Date(2019, time.April, 1, 7, 55, 0, 0, tz).Format(time.RFC1123Z+" (MST)"))
	expectedEmail := Email{
		Headers: Headers{
			Date:    expectedDate,
			Subject: "📧 Test 施氏食狮史",
			ReplyTo: []*mail.Address{
				{
					Name:    "艾莉絲 发件人",
					Address: "alice.fajianren@example.net",
				},
			},
			Sender: &mail.Address{
				Name:    "艾莉絲 发件人",
				Address: "alice.fajianren@example.com",
			},
			From: []*mail.Address{
				{
					Name:    "艾莉絲 发件人",
					Address: "alice.fajianren@example.com",
				},
				{
					Name:    "艾莉絲 发件人",
					Address: "alice.fajianren@example.net",
				},
			},
			To: []*mail.Address{
				{
					Name:    "與鮑伯 收件人",
					Address: "bob.shoujianren@example.com",
				},
				{
					Name:    "卡罗尔 收件人",
					Address: "carol.shoujianren@example.com",
				},
			},
			Cc: []*mail.Address{
				{
					Name:    "戴夫 收件人",
					Address: "dave.shoujianren@example.com",
				},
				{
					Name:    "伊夫 收件人",
					Address: "eve.shoujianren@example.com",
				},
			},
			Bcc: []*mail.Address{
				{
					Name:    "艾萨克 伊夫",
					Address: "isaac.shoujianren@example.com",
				},
				{
					Name:    "賈斯汀 伊夫",
					Address: "justin.shoujianren@example.com",
				},
			},
			MessageID:  "Message-Id-1@example.com",
			InReplyTo:  []MessageId{"Message-Id-0@example.com"},
			References: []MessageId{"Message-Id-0@example.com"},
			Comments:   "Message Header Comment",
			Keywords:   []string{"Keyword 1", "Keyword 2"},
			ResentDate: expectedDate,
			ResentFrom: []*mail.Address{
				{
					Name:    "艾莉絲 发件人",
					Address: "alice.fajianren@example.net",
				},
				{
					Name:    "艾莉絲 发件人",
					Address: "alice.fajianren@example.com",
				},
			},
			ResentSender: &mail.Address{
				Name:    "艾莉絲 发件人",
				Address: "alice.fajianren@example.net",
			},
			ResentTo: []*mail.Address{
				{
					Name:    "與鮑伯 收件人",
					Address: "bob.shoujianren@example.net",
				},
				{
					Name:    "卡罗尔 收件人",
					Address: "carol.shoujianren@example.net",
				},
			},
			ResentCc: []*mail.Address{
				{
					Name:    "戴夫 收件人",
					Address: "dave.shoujianren@example.net",
				},
				{
					Name:    "伊夫 收件人",
					Address: "eve.shoujianren@example.net",
				},
			},
			ResentBcc: []*mail.Address{
				{
					Name:    "艾萨克 伊夫",
					Address: "isaac.shoujianren@example.net",
				},
				{
					Name:    "賈斯汀 伊夫",
					Address: "justin.shoujianren@example.net",
				},
			},
			ResentMessageID: "Message-Id-1@example.net",
			ContentType: ContentTypeHeader{
				ContentType: "multipart/related",
				Params: map[string]string{
					"boundary": "RelatedBoundaryString",
					"charset":  "gb18030",
				},
			},
			ExtraHeaders: map[string][]string{
				"X-Clacks-Overhead": {"GNU Terry Pratchett"},
				"X-Script/function/\t !\"#$%&'()*+,-./;<=>?@[\\]^_`{|}~": {
					"TEST VALUE 1\t !\"#$%&'()*+,-./0123456789:;<=>?@ABCDEFGHIJKLMNOPQRSTUVWXYZ[\\]^_` abcdefghijklmnopqrstuvwxyz{|}~",
					"TEST VALUE 2\t !\"#$%&'()*+,-./0123456789:;<=>?@ABCDEFGHIJKLMNOPQRSTUVWXYZ[\\]^_` abcdefghijklmnopqrstuvwxyz{|}~",
				},
			},
		},
		Text: `石室诗士施氏，嗜狮，誓食十狮。
氏时时适市视狮。
十时，适十狮适市。
是时，适施氏适市。
氏视是十狮，恃矢势，使是十狮逝世。
氏拾是十狮尸，适石室。
石室湿，氏使侍拭石室。
石室拭，氏始试食是十狮。
食时，始识是十狮尸，实十石狮尸。
试释是事。`,
		EnrichedText: `<bold>石室诗士施氏，嗜狮，誓食十狮。</bold>
<italic>氏时时适市视狮。</italic>
<fixed>十时，适十狮适市。</fixed>
<underline>是时，适施氏适市。</underline>
氏视是十狮，恃矢势，使是十狮逝世。
氏拾是十狮尸，适石室。
石室湿，氏使侍拭石室。
石室拭，氏始试食是十狮。
食时，始识是十狮尸，实十石狮尸。
试释是事。`,
		HTML: `<html>
<div dir="ltr">
<p>石室诗士施氏，嗜狮，誓食十狮。<br />
氏时时适市视狮。<br />
十时，适十狮适市。<br />
是时，适施氏适市。<br />
氏视是十狮，恃矢势，使是十狮逝世。<br />
氏拾是十狮尸，适石室。<br />
石室湿，氏使侍拭石室。<br />
石室拭，氏始试食是十狮。<br />
食时，始识是十狮尸，实十石狮尸。<br />
试释是事。</p>
</div>
</html>`,
		InlineFiles: []InlineFile{
			{
				ContentID: "inline-jpg-image.jpg@example.com",
				ContentType: ContentTypeHeader{
					ContentType: "image/jpeg",
					Params: map[string]string{
						"name": "inline-jpg-image-name.jpg",
					},
				},
				ContentDisposition: ContentDispositionHeader{
					ContentDisposition: inline,
					Params: map[string]string{
						"filename": "inline-jpg-image-filename.jpg",
					},
				},
				Data: []byte{255, 216, 255, 219, 0, 67, 0, 3, 2, 2, 2, 2, 2, 3, 2, 2, 2, 3, 3, 3, 3, 4,
					6, 4, 4, 4, 4, 4, 8, 6, 6, 5, 6, 9, 8, 10, 10, 9, 8, 9, 9, 10, 12, 15, 12, 10, 11, 14, 11, 9, 9,
					13, 17, 13, 14, 15, 16, 16, 17, 16, 10, 12, 18, 19, 18, 16, 19, 15, 16, 16, 16, 255, 201, 0, 11, 8,
					0, 1, 0, 1, 1, 1, 17, 0, 255, 204, 0, 6, 0, 16, 16, 5, 255, 218, 0, 8, 1, 1, 0, 0, 63, 0, 210, 207,
					32, 255, 217},
			},
		},
	}

	testEmailFromFile(t, fp, expectedEmail)
}
func TestParseEmailChineseMultipartRelatedGbkOverBase64(t *testing.T) {
	fp := "tests/test_chinese_multipart_related_gbk_over_base64.txt"
	tz, _ := time.LoadLocation("Asia/Shanghai")
	expectedDate, _ := time.Parse(
		time.RFC1123Z+" (MST)",
		time.Date(2019, time.April, 1, 7, 55, 0, 0, tz).Format(time.RFC1123Z+" (MST)"))
	expectedEmail := Email{
		Headers: Headers{
			Date:    expectedDate,
			Subject: "📧 Test 施氏食狮史",
			ReplyTo: []*mail.Address{
				{
					Name:    "艾莉絲 发件人",
					Address: "alice.fajianren@example.net",
				},
			},
			Sender: &mail.Address{
				Name:    "艾莉絲 发件人",
				Address: "alice.fajianren@example.com",
			},
			From: []*mail.Address{
				{
					Name:    "艾莉絲 发件人",
					Address: "alice.fajianren@example.com",
				},
				{
					Name:    "艾莉絲 发件人",
					Address: "alice.fajianren@example.net",
				},
			},
			To: []*mail.Address{
				{
					Name:    "與鮑伯 收件人",
					Address: "bob.shoujianren@example.com",
				},
				{
					Name:    "卡罗尔 收件人",
					Address: "carol.shoujianren@example.com",
				},
			},
			Cc: []*mail.Address{
				{
					Name:    "戴夫 收件人",
					Address: "dave.shoujianren@example.com",
				},
				{
					Name:    "伊夫 收件人",
					Address: "eve.shoujianren@example.com",
				},
			},
			Bcc: []*mail.Address{
				{
					Name:    "艾萨克 伊夫",
					Address: "isaac.shoujianren@example.com",
				},
				{
					Name:    "賈斯汀 伊夫",
					Address: "justin.shoujianren@example.com",
				},
			},
			MessageID:  "Message-Id-1@example.com",
			InReplyTo:  []MessageId{"Message-Id-0@example.com"},
			References: []MessageId{"Message-Id-0@example.com"},
			Comments:   "Message Header Comment",
			Keywords:   []string{"Keyword 1", "Keyword 2"},
			ResentDate: expectedDate,
			ResentFrom: []*mail.Address{
				{
					Name:    "艾莉絲 发件人",
					Address: "alice.fajianren@example.net",
				},
				{
					Name:    "艾莉絲 发件人",
					Address: "alice.fajianren@example.com",
				},
			},
			ResentSender: &mail.Address{
				Name:    "艾莉絲 发件人",
				Address: "alice.fajianren@example.net",
			},
			ResentTo: []*mail.Address{
				{
					Name:    "與鮑伯 收件人",
					Address: "bob.shoujianren@example.net",
				},
				{
					Name:    "卡罗尔 收件人",
					Address: "carol.shoujianren@example.net",
				},
			},
			ResentCc: []*mail.Address{
				{
					Name:    "戴夫 收件人",
					Address: "dave.shoujianren@example.net",
				},
				{
					Name:    "伊夫 收件人",
					Address: "eve.shoujianren@example.net",
				},
			},
			ResentBcc: []*mail.Address{
				{
					Name:    "艾萨克 伊夫",
					Address: "isaac.shoujianren@example.net",
				},
				{
					Name:    "賈斯汀 伊夫",
					Address: "justin.shoujianren@example.net",
				},
			},
			ResentMessageID: "Message-Id-1@example.net",
			ContentType: ContentTypeHeader{
				ContentType: "multipart/related",
				Params: map[string]string{
					"boundary": "RelatedBoundaryString",
					"charset":  "gbk",
				},
			},
			ExtraHeaders: map[string][]string{
				"X-Clacks-Overhead": {"GNU Terry Pratchett"},
				"X-Script/function/\t !\"#$%&'()*+,-./;<=>?@[\\]^_`{|}~": {
					"TEST VALUE 1\t !\"#$%&'()*+,-./0123456789:;<=>?@ABCDEFGHIJKLMNOPQRSTUVWXYZ[\\]^_` abcdefghijklmnopqrstuvwxyz{|}~",
					"TEST VALUE 2\t !\"#$%&'()*+,-./0123456789:;<=>?@ABCDEFGHIJKLMNOPQRSTUVWXYZ[\\]^_` abcdefghijklmnopqrstuvwxyz{|}~",
				},
			},
		},
		Text: `石室诗士施氏，嗜狮，誓食十狮。
氏时时适市视狮。
十时，适十狮适市。
是时，适施氏适市。
氏视是十狮，恃矢势，使是十狮逝世。
氏拾是十狮尸，适石室。
石室湿，氏使侍拭石室。
石室拭，氏始试食是十狮。
食时，始识是十狮尸，实十石狮尸。
试释是事。`,
		EnrichedText: `<bold>石室诗士施氏，嗜狮，誓食十狮。</bold>
<italic>氏时时适市视狮。</italic>
<fixed>十时，适十狮适市。</fixed>
<underline>是时，适施氏适市。</underline>
氏视是十狮，恃矢势，使是十狮逝世。
氏拾是十狮尸，适石室。
石室湿，氏使侍拭石室。
石室拭，氏始试食是十狮。
食时，始识是十狮尸，实十石狮尸。
试释是事。`,
		HTML: `<html>
<div dir="ltr">
<p>石室诗士施氏，嗜狮，誓食十狮。<br />
氏时时适市视狮。<br />
十时，适十狮适市。<br />
是时，适施氏适市。<br />
氏视是十狮，恃矢势，使是十狮逝世。<br />
氏拾是十狮尸，适石室。<br />
石室湿，氏使侍拭石室。<br />
石室拭，氏始试食是十狮。<br />
食时，始识是十狮尸，实十石狮尸。<br />
试释是事。</p>
</div>
</html>`,
		InlineFiles: []InlineFile{
			{
				ContentID: "inline-jpg-image.jpg@example.com",
				ContentType: ContentTypeHeader{
					ContentType: "image/jpeg",
					Params: map[string]string{
						"name": "inline-jpg-image-name.jpg",
					},
				},
				ContentDisposition: ContentDispositionHeader{
					ContentDisposition: inline,
					Params: map[string]string{
						"filename": "inline-jpg-image-filename.jpg",
					},
				},
				Data: []byte{255, 216, 255, 219, 0, 67, 0, 3, 2, 2, 2, 2, 2, 3, 2, 2, 2, 3, 3, 3, 3, 4,
					6, 4, 4, 4, 4, 4, 8, 6, 6, 5, 6, 9, 8, 10, 10, 9, 8, 9, 9, 10, 12, 15, 12, 10, 11, 14, 11, 9, 9,
					13, 17, 13, 14, 15, 16, 16, 17, 16, 10, 12, 18, 19, 18, 16, 19, 15, 16, 16, 16, 255, 201, 0, 11, 8,
					0, 1, 0, 1, 1, 1, 17, 0, 255, 204, 0, 6, 0, 16, 16, 5, 255, 218, 0, 8, 1, 1, 0, 0, 63, 0, 210, 207,
					32, 255, 217},
			},
		},
	}

	testEmailFromFile(t, fp, expectedEmail)
}
func TestParseEmailChineseMultipartRelatedGbkOverQuotedprintable(t *testing.T) {
	fp := "tests/test_chinese_multipart_related_gbk_over_quoted-printable.txt"
	tz, _ := time.LoadLocation("Asia/Shanghai")
	expectedDate, _ := time.Parse(
		time.RFC1123Z+" (MST)",
		time.Date(2019, time.April, 1, 7, 55, 0, 0, tz).Format(time.RFC1123Z+" (MST)"))
	expectedEmail := Email{
		Headers: Headers{
			Date:    expectedDate,
			Subject: "📧 Test 施氏食狮史",
			ReplyTo: []*mail.Address{
				{
					Name:    "艾莉絲 发件人",
					Address: "alice.fajianren@example.net",
				},
			},
			Sender: &mail.Address{
				Name:    "艾莉絲 发件人",
				Address: "alice.fajianren@example.com",
			},
			From: []*mail.Address{
				{
					Name:    "艾莉絲 发件人",
					Address: "alice.fajianren@example.com",
				},
				{
					Name:    "艾莉絲 发件人",
					Address: "alice.fajianren@example.net",
				},
			},
			To: []*mail.Address{
				{
					Name:    "與鮑伯 收件人",
					Address: "bob.shoujianren@example.com",
				},
				{
					Name:    "卡罗尔 收件人",
					Address: "carol.shoujianren@example.com",
				},
			},
			Cc: []*mail.Address{
				{
					Name:    "戴夫 收件人",
					Address: "dave.shoujianren@example.com",
				},
				{
					Name:    "伊夫 收件人",
					Address: "eve.shoujianren@example.com",
				},
			},
			Bcc: []*mail.Address{
				{
					Name:    "艾萨克 伊夫",
					Address: "isaac.shoujianren@example.com",
				},
				{
					Name:    "賈斯汀 伊夫",
					Address: "justin.shoujianren@example.com",
				},
			},
			MessageID:  "Message-Id-1@example.com",
			InReplyTo:  []MessageId{"Message-Id-0@example.com"},
			References: []MessageId{"Message-Id-0@example.com"},
			Comments:   "Message Header Comment",
			Keywords:   []string{"Keyword 1", "Keyword 2"},
			ResentDate: expectedDate,
			ResentFrom: []*mail.Address{
				{
					Name:    "艾莉絲 发件人",
					Address: "alice.fajianren@example.net",
				},
				{
					Name:    "艾莉絲 发件人",
					Address: "alice.fajianren@example.com",
				},
			},
			ResentSender: &mail.Address{
				Name:    "艾莉絲 发件人",
				Address: "alice.fajianren@example.net",
			},
			ResentTo: []*mail.Address{
				{
					Name:    "與鮑伯 收件人",
					Address: "bob.shoujianren@example.net",
				},
				{
					Name:    "卡罗尔 收件人",
					Address: "carol.shoujianren@example.net",
				},
			},
			ResentCc: []*mail.Address{
				{
					Name:    "戴夫 收件人",
					Address: "dave.shoujianren@example.net",
				},
				{
					Name:    "伊夫 收件人",
					Address: "eve.shoujianren@example.net",
				},
			},
			ResentBcc: []*mail.Address{
				{
					Name:    "艾萨克 伊夫",
					Address: "isaac.shoujianren@example.net",
				},
				{
					Name:    "賈斯汀 伊夫",
					Address: "justin.shoujianren@example.net",
				},
			},
			ResentMessageID: "Message-Id-1@example.net",
			ContentType: ContentTypeHeader{
				ContentType: "multipart/related",
				Params: map[string]string{
					"boundary": "RelatedBoundaryString",
					"charset":  "gbk",
				},
			},
			ExtraHeaders: map[string][]string{
				"X-Clacks-Overhead": {"GNU Terry Pratchett"},
				"X-Script/function/\t !\"#$%&'()*+,-./;<=>?@[\\]^_`{|}~": {
					"TEST VALUE 1\t !\"#$%&'()*+,-./0123456789:;<=>?@ABCDEFGHIJKLMNOPQRSTUVWXYZ[\\]^_` abcdefghijklmnopqrstuvwxyz{|}~",
					"TEST VALUE 2\t !\"#$%&'()*+,-./0123456789:;<=>?@ABCDEFGHIJKLMNOPQRSTUVWXYZ[\\]^_` abcdefghijklmnopqrstuvwxyz{|}~",
				},
			},
		},
		Text: `石室诗士施氏，嗜狮，誓食十狮。
氏时时适市视狮。
十时，适十狮适市。
是时，适施氏适市。
氏视是十狮，恃矢势，使是十狮逝世。
氏拾是十狮尸，适石室。
石室湿，氏使侍拭石室。
石室拭，氏始试食是十狮。
食时，始识是十狮尸，实十石狮尸。
试释是事。`,
		EnrichedText: `<bold>石室诗士施氏，嗜狮，誓食十狮。</bold>
<italic>氏时时适市视狮。</italic>
<fixed>十时，适十狮适市。</fixed>
<underline>是时，适施氏适市。</underline>
氏视是十狮，恃矢势，使是十狮逝世。
氏拾是十狮尸，适石室。
石室湿，氏使侍拭石室。
石室拭，氏始试食是十狮。
食时，始识是十狮尸，实十石狮尸。
试释是事。`,
		HTML: `<html>
<div dir="ltr">
<p>石室诗士施氏，嗜狮，誓食十狮。<br />
氏时时适市视狮。<br />
十时，适十狮适市。<br />
是时，适施氏适市。<br />
氏视是十狮，恃矢势，使是十狮逝世。<br />
氏拾是十狮尸，适石室。<br />
石室湿，氏使侍拭石室。<br />
石室拭，氏始试食是十狮。<br />
食时，始识是十狮尸，实十石狮尸。<br />
试释是事。</p>
</div>
</html>`,
		InlineFiles: []InlineFile{
			{
				ContentID: "inline-jpg-image.jpg@example.com",
				ContentType: ContentTypeHeader{
					ContentType: "image/jpeg",
					Params: map[string]string{
						"name": "inline-jpg-image-name.jpg",
					},
				},
				ContentDisposition: ContentDispositionHeader{
					ContentDisposition: inline,
					Params: map[string]string{
						"filename": "inline-jpg-image-filename.jpg",
					},
				},
				Data: []byte{255, 216, 255, 219, 0, 67, 0, 3, 2, 2, 2, 2, 2, 3, 2, 2, 2, 3, 3, 3, 3, 4,
					6, 4, 4, 4, 4, 4, 8, 6, 6, 5, 6, 9, 8, 10, 10, 9, 8, 9, 9, 10, 12, 15, 12, 10, 11, 14, 11, 9, 9,
					13, 17, 13, 14, 15, 16, 16, 17, 16, 10, 12, 18, 19, 18, 16, 19, 15, 16, 16, 16, 255, 201, 0, 11, 8,
					0, 1, 0, 1, 1, 1, 17, 0, 255, 204, 0, 6, 0, 16, 16, 5, 255, 218, 0, 8, 1, 1, 0, 0, 63, 0, 210, 207,
					32, 255, 217},
			},
		},
	}

	testEmailFromFile(t, fp, expectedEmail)
}
func TestParseEmailChineseMultipartMixedGb18030OverBase64(t *testing.T) {
	fp := "tests/test_chinese_multipart_mixed_gb18030_over_base64.txt"
	tz, _ := time.LoadLocation("Asia/Shanghai")
	expectedDate, _ := time.Parse(
		time.RFC1123Z+" (MST)",
		time.Date(2019, time.April, 1, 7, 55, 0, 0, tz).Format(time.RFC1123Z+" (MST)"))
	expectedEmail := Email{
		Headers: Headers{
			Date:    expectedDate,
			Subject: "📧 Test 施氏食狮史",
			ReplyTo: []*mail.Address{
				{
					Name:    "艾莉絲 发件人",
					Address: "alice.fajianren@example.net",
				},
			},
			Sender: &mail.Address{
				Name:    "艾莉絲 发件人",
				Address: "alice.fajianren@example.com",
			},
			From: []*mail.Address{
				{
					Name:    "艾莉絲 发件人",
					Address: "alice.fajianren@example.com",
				},
				{
					Name:    "艾莉絲 发件人",
					Address: "alice.fajianren@example.net",
				},
			},
			To: []*mail.Address{
				{
					Name:    "與鮑伯 收件人",
					Address: "bob.shoujianren@example.com",
				},
				{
					Name:    "卡罗尔 收件人",
					Address: "carol.shoujianren@example.com",
				},
			},
			Cc: []*mail.Address{
				{
					Name:    "戴夫 收件人",
					Address: "dave.shoujianren@example.com",
				},
				{
					Name:    "伊夫 收件人",
					Address: "eve.shoujianren@example.com",
				},
			},
			Bcc: []*mail.Address{
				{
					Name:    "艾萨克 伊夫",
					Address: "isaac.shoujianren@example.com",
				},
				{
					Name:    "賈斯汀 伊夫",
					Address: "justin.shoujianren@example.com",
				},
			},
			MessageID:  "Message-Id-1@example.com",
			InReplyTo:  []MessageId{"Message-Id-0@example.com"},
			References: []MessageId{"Message-Id-0@example.com"},
			Comments:   "Message Header Comment",
			Keywords:   []string{"Keyword 1", "Keyword 2"},
			ResentDate: expectedDate,
			ResentFrom: []*mail.Address{
				{
					Name:    "艾莉絲 发件人",
					Address: "alice.fajianren@example.net",
				},
				{
					Name:    "艾莉絲 发件人",
					Address: "alice.fajianren@example.com",
				},
			},
			ResentSender: &mail.Address{
				Name:    "艾莉絲 发件人",
				Address: "alice.fajianren@example.net",
			},
			ResentTo: []*mail.Address{
				{
					Name:    "與鮑伯 收件人",
					Address: "bob.shoujianren@example.net",
				},
				{
					Name:    "卡罗尔 收件人",
					Address: "carol.shoujianren@example.net",
				},
			},
			ResentCc: []*mail.Address{
				{
					Name:    "戴夫 收件人",
					Address: "dave.shoujianren@example.net",
				},
				{
					Name:    "伊夫 收件人",
					Address: "eve.shoujianren@example.net",
				},
			},
			ResentBcc: []*mail.Address{
				{
					Name:    "艾萨克 伊夫",
					Address: "isaac.shoujianren@example.net",
				},
				{
					Name:    "賈斯汀 伊夫",
					Address: "justin.shoujianren@example.net",
				},
			},
			ResentMessageID: "Message-Id-1@example.net",
			ContentType: ContentTypeHeader{
				ContentType: "multipart/mixed",
				Params: map[string]string{
					"boundary": "MixedBoundaryString",
					"charset":  "gb18030",
				},
			},
			ExtraHeaders: map[string][]string{
				"X-Clacks-Overhead": {"GNU Terry Pratchett"},
				"X-Script/function/\t !\"#$%&'()*+,-./;<=>?@[\\]^_`{|}~": {
					"TEST VALUE 1\t !\"#$%&'()*+,-./0123456789:;<=>?@ABCDEFGHIJKLMNOPQRSTUVWXYZ[\\]^_` abcdefghijklmnopqrstuvwxyz{|}~",
					"TEST VALUE 2\t !\"#$%&'()*+,-./0123456789:;<=>?@ABCDEFGHIJKLMNOPQRSTUVWXYZ[\\]^_` abcdefghijklmnopqrstuvwxyz{|}~",
				},
			},
		},
		Text: `石室诗士施氏，嗜狮，誓食十狮。
氏时时适市视狮。
十时，适十狮适市。
是时，适施氏适市。
氏视是十狮，恃矢势，使是十狮逝世。
氏拾是十狮尸，适石室。
石室湿，氏使侍拭石室。
石室拭，氏始试食是十狮。
食时，始识是十狮尸，实十石狮尸。
试释是事。`,
		EnrichedText: `<bold>石室诗士施氏，嗜狮，誓食十狮。</bold>
<italic>氏时时适市视狮。</italic>
<fixed>十时，适十狮适市。</fixed>
<underline>是时，适施氏适市。</underline>
氏视是十狮，恃矢势，使是十狮逝世。
氏拾是十狮尸，适石室。
石室湿，氏使侍拭石室。
石室拭，氏始试食是十狮。
食时，始识是十狮尸，实十石狮尸。
试释是事。`,
		HTML: `<html>
<div dir="ltr">
<p>石室诗士施氏，嗜狮，誓食十狮。<br />
氏时时适市视狮。<br />
十时，适十狮适市。<br />
是时，适施氏适市。<br />
氏视是十狮，恃矢势，使是十狮逝世。<br />
氏拾是十狮尸，适石室。<br />
石室湿，氏使侍拭石室。<br />
石室拭，氏始试食是十狮。<br />
食时，始识是十狮尸，实十石狮尸。<br />
试释是事。</p>
</div>
</html>`,
		InlineFiles: []InlineFile{
			{
				ContentType: ContentTypeHeader{
					ContentType: "image/jpeg",
					Params: map[string]string{
						"name": "inline-jpg-image-without-disposition.jpg",
					},
				},
				ContentDisposition: ContentDispositionHeader{
					ContentDisposition: "",
					Params:             map[string]string(nil),
				},
				Data: []byte{255, 216, 255, 219, 0, 67, 0, 3, 2, 2, 2, 2, 2, 3, 2, 2, 2, 3, 3, 3, 3, 4,
					6, 4, 4, 4, 4, 4, 8, 6, 6, 5, 6, 9, 8, 10, 10, 9, 8, 9, 9, 10, 12, 15, 12, 10, 11, 14, 11, 9, 9,
					13, 17, 13, 14, 15, 16, 16, 17, 16, 10, 12, 18, 19, 18, 16, 19, 15, 16, 16, 16, 255, 201, 0, 11, 8,
					0, 1, 0, 1, 1, 1, 17, 0, 255, 204, 0, 6, 0, 16, 16, 5, 255, 218, 0, 8, 1, 1, 0, 0, 63, 0, 210, 207,
					32, 255, 217},
			},
			{
				ContentID: "inline-jpg-image.jpg@example.com",
				ContentType: ContentTypeHeader{
					ContentType: "image/jpeg",
					Params: map[string]string{
						"name": "inline-jpg-image-name.jpg",
					},
				},
				ContentDisposition: ContentDispositionHeader{
					ContentDisposition: inline,
					Params: map[string]string{
						"filename": "inline-jpg-image-filename.jpg",
					},
				},
				Data: []byte{255, 216, 255, 219, 0, 67, 0, 3, 2, 2, 2, 2, 2, 3, 2, 2, 2, 3, 3, 3, 3, 4,
					6, 4, 4, 4, 4, 4, 8, 6, 6, 5, 6, 9, 8, 10, 10, 9, 8, 9, 9, 10, 12, 15, 12, 10, 11, 14, 11, 9, 9,
					13, 17, 13, 14, 15, 16, 16, 17, 16, 10, 12, 18, 19, 18, 16, 19, 15, 16, 16, 16, 255, 201, 0, 11, 8,
					0, 1, 0, 1, 1, 1, 17, 0, 255, 204, 0, 6, 0, 16, 16, 5, 255, 218, 0, 8, 1, 1, 0, 0, 63, 0, 210, 207,
					32, 255, 217},
			},
		},
		AttachedFiles: []AttachedFile{
			{
				ContentType: ContentTypeHeader{
					ContentType: "application/pdf",
					Params: map[string]string{
						"name": "attached-pdf-name.pdf",
					},
				},
				ContentDisposition: ContentDispositionHeader{
					ContentDisposition: attachment,
					Params: map[string]string{
						"filename": "attached-pdf-filename.pdf",
					},
				},
				Data: []byte{37, 80, 68, 70, 45, 49, 46, 13, 116, 114, 97, 105, 108, 101, 114, 60, 60,
					47, 82, 111, 111, 116, 60, 60, 47, 80, 97, 103, 101, 115, 60, 60, 47, 75, 105, 100, 115, 91, 60,
					60, 47, 77, 101, 100, 105, 97, 66, 111, 120, 91, 48, 32, 48, 32, 51, 32, 51, 93, 62, 62, 93, 62,
					62, 62, 62, 62, 62},
			},
			{
				ContentType: ContentTypeHeader{
					ContentType: "application/pdf",
					Params: map[string]string{
						"name": "attached-pdf-without-disposition.pdf",
					},
				},
				ContentDisposition: ContentDispositionHeader{
					ContentDisposition: "",
					Params:             map[string]string(nil),
				},
				Data: []byte{37, 80, 68, 70, 45, 49, 46, 13, 116, 114, 97, 105, 108, 101, 114, 60, 60,
					47, 82, 111, 111, 116, 60, 60, 47, 80, 97, 103, 101, 115, 60, 60, 47, 75, 105, 100, 115, 91, 60,
					60, 47, 77, 101, 100, 105, 97, 66, 111, 120, 91, 48, 32, 48, 32, 51, 32, 51, 93, 62, 62, 93, 62,
					62, 62, 62, 62, 62},
			},
			{
				ContentType: ContentTypeHeader{
					ContentType: "application/json",
					Params: map[string]string{
						"name": "attached-json-name.json",
					},
				},
				ContentDisposition: ContentDispositionHeader{
					ContentDisposition: attachment,
					Params: map[string]string{
						"filename": "attached-json-filename.json",
					},
				},
				Data: []byte{123, 34, 102, 111, 111, 34, 58, 34, 98, 97, 114, 34, 125},
			},
			{
				ContentType: ContentTypeHeader{
					ContentType: "text/plain",
					Params: map[string]string{
						"name": "attached-text-plain-name.txt",
					},
				},
				ContentDisposition: ContentDispositionHeader{
					ContentDisposition: attachment,
					Params: map[string]string{
						"filename": "attached-text-plain-filename.txt",
					},
				},
				Data: []byte{84, 101, 120, 116, 47, 112, 108, 97, 105, 110, 32, 99, 111, 110, 116, 101, 110, 116, 32,
					97, 115, 32, 97, 110, 32, 97, 116, 116, 97, 99, 104, 101, 100, 32, 46, 116, 120, 116, 32, 102, 105,
					108, 101, 46},
			},
			{
				ContentType: ContentTypeHeader{
					ContentType: "text/html",
					Params: map[string]string{
						"name": "attached-text-html-name.html",
					},
				},
				ContentDisposition: ContentDispositionHeader{
					ContentDisposition: attachment,
					Params: map[string]string{
						"filename": "attached-text-html-filename.html",
					},
				},
				Data: []byte{84, 101, 120, 116, 47, 104, 116, 109, 108, 32, 99, 111, 110, 116, 101, 110, 116, 32, 97,
					115, 32, 97, 110, 32, 97, 116, 116, 97, 99, 104, 101, 100, 32, 46, 104, 116, 109, 108, 32, 102,
					105, 108, 101, 46},
			},
		},
	}

	testEmailFromFile(t, fp, expectedEmail)
}
func TestParseEmailChineseMultipartMixedGb18030OverQuotedprintable(t *testing.T) {
	fp := "tests/test_chinese_multipart_mixed_gb18030_over_quoted-printable.txt"
	tz, _ := time.LoadLocation("Asia/Shanghai")
	expectedDate, _ := time.Parse(
		time.RFC1123Z+" (MST)",
		time.Date(2019, time.April, 1, 7, 55, 0, 0, tz).Format(time.RFC1123Z+" (MST)"))
	expectedEmail := Email{
		Headers: Headers{
			Date:    expectedDate,
			Subject: "📧 Test 施氏食狮史",
			ReplyTo: []*mail.Address{
				{
					Name:    "艾莉絲 发件人",
					Address: "alice.fajianren@example.net",
				},
			},
			Sender: &mail.Address{
				Name:    "艾莉絲 发件人",
				Address: "alice.fajianren@example.com",
			},
			From: []*mail.Address{
				{
					Name:    "艾莉絲 发件人",
					Address: "alice.fajianren@example.com",
				},
				{
					Name:    "艾莉絲 发件人",
					Address: "alice.fajianren@example.net",
				},
			},
			To: []*mail.Address{
				{
					Name:    "與鮑伯 收件人",
					Address: "bob.shoujianren@example.com",
				},
				{
					Name:    "卡罗尔 收件人",
					Address: "carol.shoujianren@example.com",
				},
			},
			Cc: []*mail.Address{
				{
					Name:    "戴夫 收件人",
					Address: "dave.shoujianren@example.com",
				},
				{
					Name:    "伊夫 收件人",
					Address: "eve.shoujianren@example.com",
				},
			},
			Bcc: []*mail.Address{
				{
					Name:    "艾萨克 伊夫",
					Address: "isaac.shoujianren@example.com",
				},
				{
					Name:    "賈斯汀 伊夫",
					Address: "justin.shoujianren@example.com",
				},
			},
			MessageID:  "Message-Id-1@example.com",
			InReplyTo:  []MessageId{"Message-Id-0@example.com"},
			References: []MessageId{"Message-Id-0@example.com"},
			Comments:   "Message Header Comment",
			Keywords:   []string{"Keyword 1", "Keyword 2"},
			ResentDate: expectedDate,
			ResentFrom: []*mail.Address{
				{
					Name:    "艾莉絲 发件人",
					Address: "alice.fajianren@example.net",
				},
				{
					Name:    "艾莉絲 发件人",
					Address: "alice.fajianren@example.com",
				},
			},
			ResentSender: &mail.Address{
				Name:    "艾莉絲 发件人",
				Address: "alice.fajianren@example.net",
			},
			ResentTo: []*mail.Address{
				{
					Name:    "與鮑伯 收件人",
					Address: "bob.shoujianren@example.net",
				},
				{
					Name:    "卡罗尔 收件人",
					Address: "carol.shoujianren@example.net",
				},
			},
			ResentCc: []*mail.Address{
				{
					Name:    "戴夫 收件人",
					Address: "dave.shoujianren@example.net",
				},
				{
					Name:    "伊夫 收件人",
					Address: "eve.shoujianren@example.net",
				},
			},
			ResentBcc: []*mail.Address{
				{
					Name:    "艾萨克 伊夫",
					Address: "isaac.shoujianren@example.net",
				},
				{
					Name:    "賈斯汀 伊夫",
					Address: "justin.shoujianren@example.net",
				},
			},
			ResentMessageID: "Message-Id-1@example.net",
			ContentType: ContentTypeHeader{
				ContentType: "multipart/mixed",
				Params: map[string]string{
					"boundary": "MixedBoundaryString",
					"charset":  "gb18030",
				},
			},
			ExtraHeaders: map[string][]string{
				"X-Clacks-Overhead": {"GNU Terry Pratchett"},
				"X-Script/function/\t !\"#$%&'()*+,-./;<=>?@[\\]^_`{|}~": {
					"TEST VALUE 1\t !\"#$%&'()*+,-./0123456789:;<=>?@ABCDEFGHIJKLMNOPQRSTUVWXYZ[\\]^_` abcdefghijklmnopqrstuvwxyz{|}~",
					"TEST VALUE 2\t !\"#$%&'()*+,-./0123456789:;<=>?@ABCDEFGHIJKLMNOPQRSTUVWXYZ[\\]^_` abcdefghijklmnopqrstuvwxyz{|}~",
				},
			},
		},
		Text: `石室诗士施氏，嗜狮，誓食十狮。
氏时时适市视狮。
十时，适十狮适市。
是时，适施氏适市。
氏视是十狮，恃矢势，使是十狮逝世。
氏拾是十狮尸，适石室。
石室湿，氏使侍拭石室。
石室拭，氏始试食是十狮。
食时，始识是十狮尸，实十石狮尸。
试释是事。`,
		EnrichedText: `<bold>石室诗士施氏，嗜狮，誓食十狮。</bold>
<italic>氏时时适市视狮。</italic>
<fixed>十时，适十狮适市。</fixed>
<underline>是时，适施氏适市。</underline>
氏视是十狮，恃矢势，使是十狮逝世。
氏拾是十狮尸，适石室。
石室湿，氏使侍拭石室。
石室拭，氏始试食是十狮。
食时，始识是十狮尸，实十石狮尸。
试释是事。`,
		HTML: `<html>
<div dir="ltr">
<p>石室诗士施氏，嗜狮，誓食十狮。<br />
氏时时适市视狮。<br />
十时，适十狮适市。<br />
是时，适施氏适市。<br />
氏视是十狮，恃矢势，使是十狮逝世。<br />
氏拾是十狮尸，适石室。<br />
石室湿，氏使侍拭石室。<br />
石室拭，氏始试食是十狮。<br />
食时，始识是十狮尸，实十石狮尸。<br />
试释是事。</p>
</div>
</html>`,
		InlineFiles: []InlineFile{
			{
				ContentType: ContentTypeHeader{
					ContentType: "image/jpeg",
					Params: map[string]string{
						"name": "inline-jpg-image-without-disposition.jpg",
					},
				},
				ContentDisposition: ContentDispositionHeader{
					ContentDisposition: "",
					Params:             map[string]string(nil),
				},
				Data: []byte{255, 216, 255, 219, 0, 67, 0, 3, 2, 2, 2, 2, 2, 3, 2, 2, 2, 3, 3, 3, 3, 4,
					6, 4, 4, 4, 4, 4, 8, 6, 6, 5, 6, 9, 8, 10, 10, 9, 8, 9, 9, 10, 12, 15, 12, 10, 11, 14, 11, 9, 9,
					13, 17, 13, 14, 15, 16, 16, 17, 16, 10, 12, 18, 19, 18, 16, 19, 15, 16, 16, 16, 255, 201, 0, 11, 8,
					0, 1, 0, 1, 1, 1, 17, 0, 255, 204, 0, 6, 0, 16, 16, 5, 255, 218, 0, 8, 1, 1, 0, 0, 63, 0, 210, 207,
					32, 255, 217},
			},
			{
				ContentID: "inline-jpg-image.jpg@example.com",
				ContentType: ContentTypeHeader{
					ContentType: "image/jpeg",
					Params: map[string]string{
						"name": "inline-jpg-image-name.jpg",
					},
				},
				ContentDisposition: ContentDispositionHeader{
					ContentDisposition: inline,
					Params: map[string]string{
						"filename": "inline-jpg-image-filename.jpg",
					},
				},
				Data: []byte{255, 216, 255, 219, 0, 67, 0, 3, 2, 2, 2, 2, 2, 3, 2, 2, 2, 3, 3, 3, 3, 4,
					6, 4, 4, 4, 4, 4, 8, 6, 6, 5, 6, 9, 8, 10, 10, 9, 8, 9, 9, 10, 12, 15, 12, 10, 11, 14, 11, 9, 9,
					13, 17, 13, 14, 15, 16, 16, 17, 16, 10, 12, 18, 19, 18, 16, 19, 15, 16, 16, 16, 255, 201, 0, 11, 8,
					0, 1, 0, 1, 1, 1, 17, 0, 255, 204, 0, 6, 0, 16, 16, 5, 255, 218, 0, 8, 1, 1, 0, 0, 63, 0, 210, 207,
					32, 255, 217},
			},
		},
		AttachedFiles: []AttachedFile{
			{
				ContentType: ContentTypeHeader{
					ContentType: "application/pdf",
					Params: map[string]string{
						"name": "attached-pdf-name.pdf",
					},
				},
				ContentDisposition: ContentDispositionHeader{
					ContentDisposition: attachment,
					Params: map[string]string{
						"filename": "attached-pdf-filename.pdf",
					},
				},
				Data: []byte{37, 80, 68, 70, 45, 49, 46, 13, 116, 114, 97, 105, 108, 101, 114, 60, 60,
					47, 82, 111, 111, 116, 60, 60, 47, 80, 97, 103, 101, 115, 60, 60, 47, 75, 105, 100, 115, 91, 60,
					60, 47, 77, 101, 100, 105, 97, 66, 111, 120, 91, 48, 32, 48, 32, 51, 32, 51, 93, 62, 62, 93, 62,
					62, 62, 62, 62, 62},
			},
			{
				ContentType: ContentTypeHeader{
					ContentType: "application/pdf",
					Params: map[string]string{
						"name": "attached-pdf-without-disposition.pdf",
					},
				},
				ContentDisposition: ContentDispositionHeader{
					ContentDisposition: "",
					Params:             map[string]string(nil),
				},
				Data: []byte{37, 80, 68, 70, 45, 49, 46, 13, 116, 114, 97, 105, 108, 101, 114, 60, 60,
					47, 82, 111, 111, 116, 60, 60, 47, 80, 97, 103, 101, 115, 60, 60, 47, 75, 105, 100, 115, 91, 60,
					60, 47, 77, 101, 100, 105, 97, 66, 111, 120, 91, 48, 32, 48, 32, 51, 32, 51, 93, 62, 62, 93, 62,
					62, 62, 62, 62, 62},
			},
			{
				ContentType: ContentTypeHeader{
					ContentType: "application/json",
					Params: map[string]string{
						"name": "attached-json-name.json",
					},
				},
				ContentDisposition: ContentDispositionHeader{
					ContentDisposition: attachment,
					Params: map[string]string{
						"filename": "attached-json-filename.json",
					},
				},
				Data: []byte{123, 34, 102, 111, 111, 34, 58, 34, 98, 97, 114, 34, 125},
			},
			{
				ContentType: ContentTypeHeader{
					ContentType: "text/plain",
					Params: map[string]string{
						"name": "attached-text-plain-name.txt",
					},
				},
				ContentDisposition: ContentDispositionHeader{
					ContentDisposition: attachment,
					Params: map[string]string{
						"filename": "attached-text-plain-filename.txt",
					},
				},
				Data: []byte{84, 101, 120, 116, 47, 112, 108, 97, 105, 110, 32, 99, 111, 110, 116, 101, 110, 116, 32,
					97, 115, 32, 97, 110, 32, 97, 116, 116, 97, 99, 104, 101, 100, 32, 46, 116, 120, 116, 32, 102, 105,
					108, 101, 46},
			},
			{
				ContentType: ContentTypeHeader{
					ContentType: "text/html",
					Params: map[string]string{
						"name": "attached-text-html-name.html",
					},
				},
				ContentDisposition: ContentDispositionHeader{
					ContentDisposition: attachment,
					Params: map[string]string{
						"filename": "attached-text-html-filename.html",
					},
				},
				Data: []byte{84, 101, 120, 116, 47, 104, 116, 109, 108, 32, 99, 111, 110, 116, 101, 110, 116, 32, 97,
					115, 32, 97, 110, 32, 97, 116, 116, 97, 99, 104, 101, 100, 32, 46, 104, 116, 109, 108, 32, 102,
					105, 108, 101, 46},
			},
		},
	}

	testEmailFromFile(t, fp, expectedEmail)
}
func TestParseEmailChineseMultipartMixedGbkOverBase64(t *testing.T) {
	fp := "tests/test_chinese_multipart_mixed_gbk_over_base64.txt"
	tz, _ := time.LoadLocation("Asia/Shanghai")
	expectedDate, _ := time.Parse(
		time.RFC1123Z+" (MST)",
		time.Date(2019, time.April, 1, 7, 55, 0, 0, tz).Format(time.RFC1123Z+" (MST)"))
	expectedEmail := Email{
		Headers: Headers{
			Date:    expectedDate,
			Subject: "📧 Test 施氏食狮史",
			ReplyTo: []*mail.Address{
				{
					Name:    "艾莉絲 发件人",
					Address: "alice.fajianren@example.net",
				},
			},
			Sender: &mail.Address{
				Name:    "艾莉絲 发件人",
				Address: "alice.fajianren@example.com",
			},
			From: []*mail.Address{
				{
					Name:    "艾莉絲 发件人",
					Address: "alice.fajianren@example.com",
				},
				{
					Name:    "艾莉絲 发件人",
					Address: "alice.fajianren@example.net",
				},
			},
			To: []*mail.Address{
				{
					Name:    "與鮑伯 收件人",
					Address: "bob.shoujianren@example.com",
				},
				{
					Name:    "卡罗尔 收件人",
					Address: "carol.shoujianren@example.com",
				},
			},
			Cc: []*mail.Address{
				{
					Name:    "戴夫 收件人",
					Address: "dave.shoujianren@example.com",
				},
				{
					Name:    "伊夫 收件人",
					Address: "eve.shoujianren@example.com",
				},
			},
			Bcc: []*mail.Address{
				{
					Name:    "艾萨克 伊夫",
					Address: "isaac.shoujianren@example.com",
				},
				{
					Name:    "賈斯汀 伊夫",
					Address: "justin.shoujianren@example.com",
				},
			},
			MessageID:  "Message-Id-1@example.com",
			InReplyTo:  []MessageId{"Message-Id-0@example.com"},
			References: []MessageId{"Message-Id-0@example.com"},
			Comments:   "Message Header Comment",
			Keywords:   []string{"Keyword 1", "Keyword 2"},
			ResentDate: expectedDate,
			ResentFrom: []*mail.Address{
				{
					Name:    "艾莉絲 发件人",
					Address: "alice.fajianren@example.net",
				},
				{
					Name:    "艾莉絲 发件人",
					Address: "alice.fajianren@example.com",
				},
			},
			ResentSender: &mail.Address{
				Name:    "艾莉絲 发件人",
				Address: "alice.fajianren@example.net",
			},
			ResentTo: []*mail.Address{
				{
					Name:    "與鮑伯 收件人",
					Address: "bob.shoujianren@example.net",
				},
				{
					Name:    "卡罗尔 收件人",
					Address: "carol.shoujianren@example.net",
				},
			},
			ResentCc: []*mail.Address{
				{
					Name:    "戴夫 收件人",
					Address: "dave.shoujianren@example.net",
				},
				{
					Name:    "伊夫 收件人",
					Address: "eve.shoujianren@example.net",
				},
			},
			ResentBcc: []*mail.Address{
				{
					Name:    "艾萨克 伊夫",
					Address: "isaac.shoujianren@example.net",
				},
				{
					Name:    "賈斯汀 伊夫",
					Address: "justin.shoujianren@example.net",
				},
			},
			ResentMessageID: "Message-Id-1@example.net",
			ContentType: ContentTypeHeader{
				ContentType: "multipart/mixed",
				Params: map[string]string{
					"boundary": "MixedBoundaryString",
					"charset":  "gbk",
				},
			},
			ExtraHeaders: map[string][]string{
				"X-Clacks-Overhead": {"GNU Terry Pratchett"},
				"X-Script/function/\t !\"#$%&'()*+,-./;<=>?@[\\]^_`{|}~": {
					"TEST VALUE 1\t !\"#$%&'()*+,-./0123456789:;<=>?@ABCDEFGHIJKLMNOPQRSTUVWXYZ[\\]^_` abcdefghijklmnopqrstuvwxyz{|}~",
					"TEST VALUE 2\t !\"#$%&'()*+,-./0123456789:;<=>?@ABCDEFGHIJKLMNOPQRSTUVWXYZ[\\]^_` abcdefghijklmnopqrstuvwxyz{|}~",
				},
			},
		},
		Text: `石室诗士施氏，嗜狮，誓食十狮。
氏时时适市视狮。
十时，适十狮适市。
是时，适施氏适市。
氏视是十狮，恃矢势，使是十狮逝世。
氏拾是十狮尸，适石室。
石室湿，氏使侍拭石室。
石室拭，氏始试食是十狮。
食时，始识是十狮尸，实十石狮尸。
试释是事。`,
		EnrichedText: `<bold>石室诗士施氏，嗜狮，誓食十狮。</bold>
<italic>氏时时适市视狮。</italic>
<fixed>十时，适十狮适市。</fixed>
<underline>是时，适施氏适市。</underline>
氏视是十狮，恃矢势，使是十狮逝世。
氏拾是十狮尸，适石室。
石室湿，氏使侍拭石室。
石室拭，氏始试食是十狮。
食时，始识是十狮尸，实十石狮尸。
试释是事。`,
		HTML: `<html>
<div dir="ltr">
<p>石室诗士施氏，嗜狮，誓食十狮。<br />
氏时时适市视狮。<br />
十时，适十狮适市。<br />
是时，适施氏适市。<br />
氏视是十狮，恃矢势，使是十狮逝世。<br />
氏拾是十狮尸，适石室。<br />
石室湿，氏使侍拭石室。<br />
石室拭，氏始试食是十狮。<br />
食时，始识是十狮尸，实十石狮尸。<br />
试释是事。</p>
</div>
</html>`,
		InlineFiles: []InlineFile{
			{
				ContentType: ContentTypeHeader{
					ContentType: "image/jpeg",
					Params: map[string]string{
						"name": "inline-jpg-image-without-disposition.jpg",
					},
				},
				ContentDisposition: ContentDispositionHeader{
					ContentDisposition: "",
					Params:             map[string]string(nil),
				},
				Data: []byte{255, 216, 255, 219, 0, 67, 0, 3, 2, 2, 2, 2, 2, 3, 2, 2, 2, 3, 3, 3, 3, 4,
					6, 4, 4, 4, 4, 4, 8, 6, 6, 5, 6, 9, 8, 10, 10, 9, 8, 9, 9, 10, 12, 15, 12, 10, 11, 14, 11, 9, 9,
					13, 17, 13, 14, 15, 16, 16, 17, 16, 10, 12, 18, 19, 18, 16, 19, 15, 16, 16, 16, 255, 201, 0, 11, 8,
					0, 1, 0, 1, 1, 1, 17, 0, 255, 204, 0, 6, 0, 16, 16, 5, 255, 218, 0, 8, 1, 1, 0, 0, 63, 0, 210, 207,
					32, 255, 217},
			},
			{
				ContentID: "inline-jpg-image.jpg@example.com",
				ContentType: ContentTypeHeader{
					ContentType: "image/jpeg",
					Params: map[string]string{
						"name": "inline-jpg-image-name.jpg",
					},
				},
				ContentDisposition: ContentDispositionHeader{
					ContentDisposition: inline,
					Params: map[string]string{
						"filename": "inline-jpg-image-filename.jpg",
					},
				},
				Data: []byte{255, 216, 255, 219, 0, 67, 0, 3, 2, 2, 2, 2, 2, 3, 2, 2, 2, 3, 3, 3, 3, 4,
					6, 4, 4, 4, 4, 4, 8, 6, 6, 5, 6, 9, 8, 10, 10, 9, 8, 9, 9, 10, 12, 15, 12, 10, 11, 14, 11, 9, 9,
					13, 17, 13, 14, 15, 16, 16, 17, 16, 10, 12, 18, 19, 18, 16, 19, 15, 16, 16, 16, 255, 201, 0, 11, 8,
					0, 1, 0, 1, 1, 1, 17, 0, 255, 204, 0, 6, 0, 16, 16, 5, 255, 218, 0, 8, 1, 1, 0, 0, 63, 0, 210, 207,
					32, 255, 217},
			},
		},
		AttachedFiles: []AttachedFile{
			{
				ContentType: ContentTypeHeader{
					ContentType: "application/pdf",
					Params: map[string]string{
						"name": "attached-pdf-name.pdf",
					},
				},
				ContentDisposition: ContentDispositionHeader{
					ContentDisposition: attachment,
					Params: map[string]string{
						"filename": "attached-pdf-filename.pdf",
					},
				},
				Data: []byte{37, 80, 68, 70, 45, 49, 46, 13, 116, 114, 97, 105, 108, 101, 114, 60, 60,
					47, 82, 111, 111, 116, 60, 60, 47, 80, 97, 103, 101, 115, 60, 60, 47, 75, 105, 100, 115, 91, 60,
					60, 47, 77, 101, 100, 105, 97, 66, 111, 120, 91, 48, 32, 48, 32, 51, 32, 51, 93, 62, 62, 93, 62,
					62, 62, 62, 62, 62},
			},
			{
				ContentType: ContentTypeHeader{
					ContentType: "application/pdf",
					Params: map[string]string{
						"name": "attached-pdf-without-disposition.pdf",
					},
				},
				ContentDisposition: ContentDispositionHeader{
					ContentDisposition: "",
					Params:             map[string]string(nil),
				},
				Data: []byte{37, 80, 68, 70, 45, 49, 46, 13, 116, 114, 97, 105, 108, 101, 114, 60, 60,
					47, 82, 111, 111, 116, 60, 60, 47, 80, 97, 103, 101, 115, 60, 60, 47, 75, 105, 100, 115, 91, 60,
					60, 47, 77, 101, 100, 105, 97, 66, 111, 120, 91, 48, 32, 48, 32, 51, 32, 51, 93, 62, 62, 93, 62,
					62, 62, 62, 62, 62},
			},
			{
				ContentType: ContentTypeHeader{
					ContentType: "application/json",
					Params: map[string]string{
						"name": "attached-json-name.json",
					},
				},
				ContentDisposition: ContentDispositionHeader{
					ContentDisposition: attachment,
					Params: map[string]string{
						"filename": "attached-json-filename.json",
					},
				},
				Data: []byte{123, 34, 102, 111, 111, 34, 58, 34, 98, 97, 114, 34, 125},
			},
			{
				ContentType: ContentTypeHeader{
					ContentType: "text/plain",
					Params: map[string]string{
						"name": "attached-text-plain-name.txt",
					},
				},
				ContentDisposition: ContentDispositionHeader{
					ContentDisposition: attachment,
					Params: map[string]string{
						"filename": "attached-text-plain-filename.txt",
					},
				},
				Data: []byte{84, 101, 120, 116, 47, 112, 108, 97, 105, 110, 32, 99, 111, 110, 116, 101, 110, 116, 32,
					97, 115, 32, 97, 110, 32, 97, 116, 116, 97, 99, 104, 101, 100, 32, 46, 116, 120, 116, 32, 102, 105,
					108, 101, 46},
			},
			{
				ContentType: ContentTypeHeader{
					ContentType: "text/html",
					Params: map[string]string{
						"name": "attached-text-html-name.html",
					},
				},
				ContentDisposition: ContentDispositionHeader{
					ContentDisposition: attachment,
					Params: map[string]string{
						"filename": "attached-text-html-filename.html",
					},
				},
				Data: []byte{84, 101, 120, 116, 47, 104, 116, 109, 108, 32, 99, 111, 110, 116, 101, 110, 116, 32, 97,
					115, 32, 97, 110, 32, 97, 116, 116, 97, 99, 104, 101, 100, 32, 46, 104, 116, 109, 108, 32, 102,
					105, 108, 101, 46},
			},
		},
	}

	testEmailFromFile(t, fp, expectedEmail)
}
func TestParseEmailChineseMultipartMixedGbkOverQuotedprintable(t *testing.T) {
	fp := "tests/test_chinese_multipart_mixed_gbk_over_quoted-printable.txt"
	tz, _ := time.LoadLocation("Asia/Shanghai")
	expectedDate, _ := time.Parse(
		time.RFC1123Z+" (MST)",
		time.Date(2019, time.April, 1, 7, 55, 0, 0, tz).Format(time.RFC1123Z+" (MST)"))
	expectedEmail := Email{
		Headers: Headers{
			Date:    expectedDate,
			Subject: "📧 Test 施氏食狮史",
			ReplyTo: []*mail.Address{
				{
					Name:    "艾莉絲 发件人",
					Address: "alice.fajianren@example.net",
				},
			},
			Sender: &mail.Address{
				Name:    "艾莉絲 发件人",
				Address: "alice.fajianren@example.com",
			},
			From: []*mail.Address{
				{
					Name:    "艾莉絲 发件人",
					Address: "alice.fajianren@example.com",
				},
				{
					Name:    "艾莉絲 发件人",
					Address: "alice.fajianren@example.net",
				},
			},
			To: []*mail.Address{
				{
					Name:    "與鮑伯 收件人",
					Address: "bob.shoujianren@example.com",
				},
				{
					Name:    "卡罗尔 收件人",
					Address: "carol.shoujianren@example.com",
				},
			},
			Cc: []*mail.Address{
				{
					Name:    "戴夫 收件人",
					Address: "dave.shoujianren@example.com",
				},
				{
					Name:    "伊夫 收件人",
					Address: "eve.shoujianren@example.com",
				},
			},
			Bcc: []*mail.Address{
				{
					Name:    "艾萨克 伊夫",
					Address: "isaac.shoujianren@example.com",
				},
				{
					Name:    "賈斯汀 伊夫",
					Address: "justin.shoujianren@example.com",
				},
			},
			MessageID:  "Message-Id-1@example.com",
			InReplyTo:  []MessageId{"Message-Id-0@example.com"},
			References: []MessageId{"Message-Id-0@example.com"},
			Comments:   "Message Header Comment",
			Keywords:   []string{"Keyword 1", "Keyword 2"},
			ResentDate: expectedDate,
			ResentFrom: []*mail.Address{
				{
					Name:    "艾莉絲 发件人",
					Address: "alice.fajianren@example.net",
				},
				{
					Name:    "艾莉絲 发件人",
					Address: "alice.fajianren@example.com",
				},
			},
			ResentSender: &mail.Address{
				Name:    "艾莉絲 发件人",
				Address: "alice.fajianren@example.net",
			},
			ResentTo: []*mail.Address{
				{
					Name:    "與鮑伯 收件人",
					Address: "bob.shoujianren@example.net",
				},
				{
					Name:    "卡罗尔 收件人",
					Address: "carol.shoujianren@example.net",
				},
			},
			ResentCc: []*mail.Address{
				{
					Name:    "戴夫 收件人",
					Address: "dave.shoujianren@example.net",
				},
				{
					Name:    "伊夫 收件人",
					Address: "eve.shoujianren@example.net",
				},
			},
			ResentBcc: []*mail.Address{
				{
					Name:    "艾萨克 伊夫",
					Address: "isaac.shoujianren@example.net",
				},
				{
					Name:    "賈斯汀 伊夫",
					Address: "justin.shoujianren@example.net",
				},
			},
			ResentMessageID: "Message-Id-1@example.net",
			ContentType: ContentTypeHeader{
				ContentType: "multipart/mixed",
				Params: map[string]string{
					"boundary": "MixedBoundaryString",
					"charset":  "gbk",
				},
			},
			ExtraHeaders: map[string][]string{
				"X-Clacks-Overhead": {"GNU Terry Pratchett"},
				"X-Script/function/\t !\"#$%&'()*+,-./;<=>?@[\\]^_`{|}~": {
					"TEST VALUE 1\t !\"#$%&'()*+,-./0123456789:;<=>?@ABCDEFGHIJKLMNOPQRSTUVWXYZ[\\]^_` abcdefghijklmnopqrstuvwxyz{|}~",
					"TEST VALUE 2\t !\"#$%&'()*+,-./0123456789:;<=>?@ABCDEFGHIJKLMNOPQRSTUVWXYZ[\\]^_` abcdefghijklmnopqrstuvwxyz{|}~",
				},
			},
		},
		Text: `石室诗士施氏，嗜狮，誓食十狮。
氏时时适市视狮。
十时，适十狮适市。
是时，适施氏适市。
氏视是十狮，恃矢势，使是十狮逝世。
氏拾是十狮尸，适石室。
石室湿，氏使侍拭石室。
石室拭，氏始试食是十狮。
食时，始识是十狮尸，实十石狮尸。
试释是事。`,
		EnrichedText: `<bold>石室诗士施氏，嗜狮，誓食十狮。</bold>
<italic>氏时时适市视狮。</italic>
<fixed>十时，适十狮适市。</fixed>
<underline>是时，适施氏适市。</underline>
氏视是十狮，恃矢势，使是十狮逝世。
氏拾是十狮尸，适石室。
石室湿，氏使侍拭石室。
石室拭，氏始试食是十狮。
食时，始识是十狮尸，实十石狮尸。
试释是事。`,
		HTML: `<html>
<div dir="ltr">
<p>石室诗士施氏，嗜狮，誓食十狮。<br />
氏时时适市视狮。<br />
十时，适十狮适市。<br />
是时，适施氏适市。<br />
氏视是十狮，恃矢势，使是十狮逝世。<br />
氏拾是十狮尸，适石室。<br />
石室湿，氏使侍拭石室。<br />
石室拭，氏始试食是十狮。<br />
食时，始识是十狮尸，实十石狮尸。<br />
试释是事。</p>
</div>
</html>`,
		InlineFiles: []InlineFile{
			{
				ContentType: ContentTypeHeader{
					ContentType: "image/jpeg",
					Params: map[string]string{
						"name": "inline-jpg-image-without-disposition.jpg",
					},
				},
				ContentDisposition: ContentDispositionHeader{
					ContentDisposition: "",
					Params:             map[string]string(nil),
				},
				Data: []byte{255, 216, 255, 219, 0, 67, 0, 3, 2, 2, 2, 2, 2, 3, 2, 2, 2, 3, 3, 3, 3, 4,
					6, 4, 4, 4, 4, 4, 8, 6, 6, 5, 6, 9, 8, 10, 10, 9, 8, 9, 9, 10, 12, 15, 12, 10, 11, 14, 11, 9, 9,
					13, 17, 13, 14, 15, 16, 16, 17, 16, 10, 12, 18, 19, 18, 16, 19, 15, 16, 16, 16, 255, 201, 0, 11, 8,
					0, 1, 0, 1, 1, 1, 17, 0, 255, 204, 0, 6, 0, 16, 16, 5, 255, 218, 0, 8, 1, 1, 0, 0, 63, 0, 210, 207,
					32, 255, 217},
			},
			{
				ContentID: "inline-jpg-image.jpg@example.com",
				ContentType: ContentTypeHeader{
					ContentType: "image/jpeg",
					Params: map[string]string{
						"name": "inline-jpg-image-name.jpg",
					},
				},
				ContentDisposition: ContentDispositionHeader{
					ContentDisposition: inline,
					Params: map[string]string{
						"filename": "inline-jpg-image-filename.jpg",
					},
				},
				Data: []byte{255, 216, 255, 219, 0, 67, 0, 3, 2, 2, 2, 2, 2, 3, 2, 2, 2, 3, 3, 3, 3, 4,
					6, 4, 4, 4, 4, 4, 8, 6, 6, 5, 6, 9, 8, 10, 10, 9, 8, 9, 9, 10, 12, 15, 12, 10, 11, 14, 11, 9, 9,
					13, 17, 13, 14, 15, 16, 16, 17, 16, 10, 12, 18, 19, 18, 16, 19, 15, 16, 16, 16, 255, 201, 0, 11, 8,
					0, 1, 0, 1, 1, 1, 17, 0, 255, 204, 0, 6, 0, 16, 16, 5, 255, 218, 0, 8, 1, 1, 0, 0, 63, 0, 210, 207,
					32, 255, 217},
			},
		},
		AttachedFiles: []AttachedFile{
			{
				ContentType: ContentTypeHeader{
					ContentType: "application/pdf",
					Params: map[string]string{
						"name": "attached-pdf-name.pdf",
					},
				},
				ContentDisposition: ContentDispositionHeader{
					ContentDisposition: attachment,
					Params: map[string]string{
						"filename": "attached-pdf-filename.pdf",
					},
				},
				Data: []byte{37, 80, 68, 70, 45, 49, 46, 13, 116, 114, 97, 105, 108, 101, 114, 60, 60,
					47, 82, 111, 111, 116, 60, 60, 47, 80, 97, 103, 101, 115, 60, 60, 47, 75, 105, 100, 115, 91, 60,
					60, 47, 77, 101, 100, 105, 97, 66, 111, 120, 91, 48, 32, 48, 32, 51, 32, 51, 93, 62, 62, 93, 62,
					62, 62, 62, 62, 62},
			},
			{
				ContentType: ContentTypeHeader{
					ContentType: "application/pdf",
					Params: map[string]string{
						"name": "attached-pdf-without-disposition.pdf",
					},
				},
				ContentDisposition: ContentDispositionHeader{
					ContentDisposition: "",
					Params:             map[string]string(nil),
				},
				Data: []byte{37, 80, 68, 70, 45, 49, 46, 13, 116, 114, 97, 105, 108, 101, 114, 60, 60,
					47, 82, 111, 111, 116, 60, 60, 47, 80, 97, 103, 101, 115, 60, 60, 47, 75, 105, 100, 115, 91, 60,
					60, 47, 77, 101, 100, 105, 97, 66, 111, 120, 91, 48, 32, 48, 32, 51, 32, 51, 93, 62, 62, 93, 62,
					62, 62, 62, 62, 62},
			},
			{
				ContentType: ContentTypeHeader{
					ContentType: "application/json",
					Params: map[string]string{
						"name": "attached-json-name.json",
					},
				},
				ContentDisposition: ContentDispositionHeader{
					ContentDisposition: attachment,
					Params: map[string]string{
						"filename": "attached-json-filename.json",
					},
				},
				Data: []byte{123, 34, 102, 111, 111, 34, 58, 34, 98, 97, 114, 34, 125},
			},
			{
				ContentType: ContentTypeHeader{
					ContentType: "text/plain",
					Params: map[string]string{
						"name": "attached-text-plain-name.txt",
					},
				},
				ContentDisposition: ContentDispositionHeader{
					ContentDisposition: attachment,
					Params: map[string]string{
						"filename": "attached-text-plain-filename.txt",
					},
				},
				Data: []byte{84, 101, 120, 116, 47, 112, 108, 97, 105, 110, 32, 99, 111, 110, 116, 101, 110, 116, 32,
					97, 115, 32, 97, 110, 32, 97, 116, 116, 97, 99, 104, 101, 100, 32, 46, 116, 120, 116, 32, 102, 105,
					108, 101, 46},
			},
			{
				ContentType: ContentTypeHeader{
					ContentType: "text/html",
					Params: map[string]string{
						"name": "attached-text-html-name.html",
					},
				},
				ContentDisposition: ContentDispositionHeader{
					ContentDisposition: attachment,
					Params: map[string]string{
						"filename": "attached-text-html-filename.html",
					},
				},
				Data: []byte{84, 101, 120, 116, 47, 104, 116, 109, 108, 32, 99, 111, 110, 116, 101, 110, 116, 32, 97,
					115, 32, 97, 110, 32, 97, 116, 116, 97, 99, 104, 101, 100, 32, 46, 104, 116, 109, 108, 32, 102,
					105, 108, 101, 46},
			},
		},
	}

	testEmailFromFile(t, fp, expectedEmail)
}
func TestParseEmailChineseMultipartSignedGb18030OverBase64(t *testing.T) {
	fp := "tests/test_chinese_multipart_signed_gb18030_over_base64.txt"
	tz, _ := time.LoadLocation("Asia/Shanghai")
	expectedDate, _ := time.Parse(
		time.RFC1123Z+" (MST)",
		time.Date(2019, time.April, 1, 7, 55, 0, 0, tz).Format(time.RFC1123Z+" (MST)"))
	expectedEmail := Email{
		Headers: Headers{
			Date:    expectedDate,
			Subject: "📧 Signed Test 施氏食狮史",
			ReplyTo: []*mail.Address{
				{
					Name:    "艾莉絲 发件人",
					Address: "alice.fajianren@example.net",
				},
			},
			Sender: &mail.Address{
				Name:    "艾莉絲 发件人",
				Address: "alice.fajianren@example.com",
			},
			From: []*mail.Address{
				{
					Name:    "艾莉絲 发件人",
					Address: "alice.fajianren@example.com",
				},
				{
					Name:    "艾莉絲 发件人",
					Address: "alice.fajianren@example.net",
				},
			},
			To: []*mail.Address{
				{
					Name:    "與鮑伯 收件人",
					Address: "bob.shoujianren@example.com",
				},
				{
					Name:    "卡罗尔 收件人",
					Address: "carol.shoujianren@example.com",
				},
			},
			Cc: []*mail.Address{
				{
					Name:    "戴夫 收件人",
					Address: "dave.shoujianren@example.com",
				},
				{
					Name:    "伊夫 收件人",
					Address: "eve.shoujianren@example.com",
				},
			},
			Bcc: []*mail.Address{
				{
					Name:    "艾萨克 伊夫",
					Address: "isaac.shoujianren@example.com",
				},
				{
					Name:    "賈斯汀 伊夫",
					Address: "justin.shoujianren@example.com",
				},
			},
			MessageID:  "Message-Id-1@example.com",
			InReplyTo:  []MessageId{"Message-Id-0@example.com"},
			References: []MessageId{"Message-Id-0@example.com"},
			Comments:   "Message Header Comment",
			Keywords:   []string{"Keyword 1", "Keyword 2"},
			ResentDate: expectedDate,
			ResentFrom: []*mail.Address{
				{
					Name:    "艾莉絲 发件人",
					Address: "alice.fajianren@example.net",
				},
				{
					Name:    "艾莉絲 发件人",
					Address: "alice.fajianren@example.com",
				},
			},
			ResentSender: &mail.Address{
				Name:    "艾莉絲 发件人",
				Address: "alice.fajianren@example.net",
			},
			ResentTo: []*mail.Address{
				{
					Name:    "與鮑伯 收件人",
					Address: "bob.shoujianren@example.net",
				},
				{
					Name:    "卡罗尔 收件人",
					Address: "carol.shoujianren@example.net",
				},
			},
			ResentCc: []*mail.Address{
				{
					Name:    "戴夫 收件人",
					Address: "dave.shoujianren@example.net",
				},
				{
					Name:    "伊夫 收件人",
					Address: "eve.shoujianren@example.net",
				},
			},
			ResentBcc: []*mail.Address{
				{
					Name:    "艾萨克 伊夫",
					Address: "isaac.shoujianren@example.net",
				},
				{
					Name:    "賈斯汀 伊夫",
					Address: "justin.shoujianren@example.net",
				},
			},
			ResentMessageID: "Message-Id-1@example.net",
			ContentType: ContentTypeHeader{
				ContentType: "multipart/signed",
				Params: map[string]string{
					"boundary": "SignedBoundaryString",
					"charset":  "gb18030",
					"micalg":   "sha1",
					"protocol": "application/pkcs7-signature",
				},
			},
			ExtraHeaders: map[string][]string{
				"X-Clacks-Overhead": {"GNU Terry Pratchett"},
				"X-Script/function/\t !\"#$%&'()*+,-./;<=>?@[\\]^_`{|}~": {
					"TEST VALUE 1\t !\"#$%&'()*+,-./0123456789:;<=>?@ABCDEFGHIJKLMNOPQRSTUVWXYZ[\\]^_` abcdefghijklmnopqrstuvwxyz{|}~",
					"TEST VALUE 2\t !\"#$%&'()*+,-./0123456789:;<=>?@ABCDEFGHIJKLMNOPQRSTUVWXYZ[\\]^_` abcdefghijklmnopqrstuvwxyz{|}~",
				},
			},
		},
		Text: `石室诗士施氏，嗜狮，誓食十狮。
氏时时适市视狮。
十时，适十狮适市。
是时，适施氏适市。
氏视是十狮，恃矢势，使是十狮逝世。
氏拾是十狮尸，适石室。
石室湿，氏使侍拭石室。
石室拭，氏始试食是十狮。
食时，始识是十狮尸，实十石狮尸。
试释是事。`,
		AttachedFiles: []AttachedFile{
			{
				ContentType: ContentTypeHeader{
					ContentType: "application/pkcs7-signature",
					Params: map[string]string{
						"name": "smime.p7s",
					},
				},
				ContentDisposition: ContentDispositionHeader{
					ContentDisposition: attachment,
					Params: map[string]string{
						"filename": "smime.p7s",
					},
				},
				Data: []byte{130, 28, 135, 132, 117, 46, 142, 18, 97, 140, 126, 251, 159, 193, 199, 25, 58, 223, 189, 185, 227, 239, 158, 173, 108, 31, 71, 27, 133, 80, 165, 252, 133, 227, 174, 198, 132, 129, 159, 29, 246, 19, 235, 133, 80, 165, 252, 133, 227, 174, 198, 132, 129, 159, 29, 246, 19, 234, 49, 251, 238, 127, 7, 28, 104, 33, 200, 120, 71, 82, 232, 225, 38, 30, 249, 234, 214, 193, 244, 113, 147, 173, 251, 219, 158, 57, 252, 28, 113, 147, 173, 251, 225, 38, 24, 199, 239, 190, 173, 108, 31, 71, 27, 133, 80, 110, 120, 251, 231, 174, 198, 132, 129, 159, 29, 246, 19, 234, 8, 114, 30, 17, 212, 186, 58, 95, 200, 94, 59, 26, 18, 6, 124, 119, 216, 79, 174, 21, 65, 185, 227, 239, 158},
			},
		},
	}

	testEmailFromFile(t, fp, expectedEmail)
}

func TestParseEmailChineseMultipartSignedGb18030OverQuotedprintable(t *testing.T) {
	fp := "tests/test_chinese_multipart_signed_gb18030_over_quoted-printable.txt"
	tz, _ := time.LoadLocation("Asia/Shanghai")
	expectedDate, _ := time.Parse(
		time.RFC1123Z+" (MST)",
		time.Date(2019, time.April, 1, 7, 55, 0, 0, tz).Format(time.RFC1123Z+" (MST)"))
	expectedEmail := Email{
		Headers: Headers{
			Date:    expectedDate,
			Subject: "📧 Signed Test 施氏食狮史",
			ReplyTo: []*mail.Address{
				{
					Name:    "艾莉絲 发件人",
					Address: "alice.fajianren@example.net",
				},
			},
			Sender: &mail.Address{
				Name:    "艾莉絲 发件人",
				Address: "alice.fajianren@example.com",
			},
			From: []*mail.Address{
				{
					Name:    "艾莉絲 发件人",
					Address: "alice.fajianren@example.com",
				},
				{
					Name:    "艾莉絲 发件人",
					Address: "alice.fajianren@example.net",
				},
			},
			To: []*mail.Address{
				{
					Name:    "與鮑伯 收件人",
					Address: "bob.shoujianren@example.com",
				},
				{
					Name:    "卡罗尔 收件人",
					Address: "carol.shoujianren@example.com",
				},
			},
			Cc: []*mail.Address{
				{
					Name:    "戴夫 收件人",
					Address: "dave.shoujianren@example.com",
				},
				{
					Name:    "伊夫 收件人",
					Address: "eve.shoujianren@example.com",
				},
			},
			Bcc: []*mail.Address{
				{
					Name:    "艾萨克 伊夫",
					Address: "isaac.shoujianren@example.com",
				},
				{
					Name:    "賈斯汀 伊夫",
					Address: "justin.shoujianren@example.com",
				},
			},
			MessageID:  "Message-Id-1@example.com",
			InReplyTo:  []MessageId{"Message-Id-0@example.com"},
			References: []MessageId{"Message-Id-0@example.com"},
			Comments:   "Message Header Comment",
			Keywords:   []string{"Keyword 1", "Keyword 2"},
			ResentDate: expectedDate,
			ResentFrom: []*mail.Address{
				{
					Name:    "艾莉絲 发件人",
					Address: "alice.fajianren@example.net",
				},
				{
					Name:    "艾莉絲 发件人",
					Address: "alice.fajianren@example.com",
				},
			},
			ResentSender: &mail.Address{
				Name:    "艾莉絲 发件人",
				Address: "alice.fajianren@example.net",
			},
			ResentTo: []*mail.Address{
				{
					Name:    "與鮑伯 收件人",
					Address: "bob.shoujianren@example.net",
				},
				{
					Name:    "卡罗尔 收件人",
					Address: "carol.shoujianren@example.net",
				},
			},
			ResentCc: []*mail.Address{
				{
					Name:    "戴夫 收件人",
					Address: "dave.shoujianren@example.net",
				},
				{
					Name:    "伊夫 收件人",
					Address: "eve.shoujianren@example.net",
				},
			},
			ResentBcc: []*mail.Address{
				{
					Name:    "艾萨克 伊夫",
					Address: "isaac.shoujianren@example.net",
				},
				{
					Name:    "賈斯汀 伊夫",
					Address: "justin.shoujianren@example.net",
				},
			},
			ResentMessageID: "Message-Id-1@example.net",
			ContentType: ContentTypeHeader{
				ContentType: "multipart/signed",
				Params: map[string]string{
					"boundary": "SignedBoundaryString",
					"charset":  "gb18030",
					"micalg":   "sha1",
					"protocol": "application/pkcs7-signature",
				},
			},
			ExtraHeaders: map[string][]string{
				"X-Clacks-Overhead": {"GNU Terry Pratchett"},
				"X-Script/function/\t !\"#$%&'()*+,-./;<=>?@[\\]^_`{|}~": {
					"TEST VALUE 1\t !\"#$%&'()*+,-./0123456789:;<=>?@ABCDEFGHIJKLMNOPQRSTUVWXYZ[\\]^_` abcdefghijklmnopqrstuvwxyz{|}~",
					"TEST VALUE 2\t !\"#$%&'()*+,-./0123456789:;<=>?@ABCDEFGHIJKLMNOPQRSTUVWXYZ[\\]^_` abcdefghijklmnopqrstuvwxyz{|}~",
				},
			},
		},
		Text: `石室诗士施氏，嗜狮，誓食十狮。
氏时时适市视狮。
十时，适十狮适市。
是时，适施氏适市。
氏视是十狮，恃矢势，使是十狮逝世。
氏拾是十狮尸，适石室。
石室湿，氏使侍拭石室。
石室拭，氏始试食是十狮。
食时，始识是十狮尸，实十石狮尸。
试释是事。`,
		AttachedFiles: []AttachedFile{
			{
				ContentType: ContentTypeHeader{
					ContentType: "application/pkcs7-signature",
					Params: map[string]string{
						"name": "smime.p7s",
					},
				},
				ContentDisposition: ContentDispositionHeader{
					ContentDisposition: attachment,
					Params: map[string]string{
						"filename": "smime.p7s",
					},
				},
				Data: []byte{130, 28, 135, 132, 117, 46, 142, 18, 97, 140, 126, 251, 159, 193, 199, 25, 58, 223, 189, 185, 227, 239, 158, 173, 108, 31, 71, 27, 133, 80, 165, 252, 133, 227, 174, 198, 132, 129, 159, 29, 246, 19, 235, 133, 80, 165, 252, 133, 227, 174, 198, 132, 129, 159, 29, 246, 19, 234, 49, 251, 238, 127, 7, 28, 104, 33, 200, 120, 71, 82, 232, 225, 38, 30, 249, 234, 214, 193, 244, 113, 147, 173, 251, 219, 158, 57, 252, 28, 113, 147, 173, 251, 225, 38, 24, 199, 239, 190, 173, 108, 31, 71, 27, 133, 80, 110, 120, 251, 231, 174, 198, 132, 129, 159, 29, 246, 19, 234, 8, 114, 30, 17, 212, 186, 58, 95, 200, 94, 59, 26, 18, 6, 124, 119, 216, 79, 174, 21, 65, 185, 227, 239, 158},
			},
		},
	}

	testEmailFromFile(t, fp, expectedEmail)
}

func TestParseEmailChineseMultipartSignedGbkOverBase64(t *testing.T) {
	fp := "tests/test_chinese_multipart_signed_gbk_over_base64.txt"
	tz, _ := time.LoadLocation("Asia/Shanghai")
	expectedDate, _ := time.Parse(
		time.RFC1123Z+" (MST)",
		time.Date(2019, time.April, 1, 7, 55, 0, 0, tz).Format(time.RFC1123Z+" (MST)"))
	expectedEmail := Email{
		Headers: Headers{
			Date:    expectedDate,
			Subject: "📧 Signed Test 施氏食狮史",
			ReplyTo: []*mail.Address{
				{
					Name:    "艾莉絲 发件人",
					Address: "alice.fajianren@example.net",
				},
			},
			Sender: &mail.Address{
				Name:    "艾莉絲 发件人",
				Address: "alice.fajianren@example.com",
			},
			From: []*mail.Address{
				{
					Name:    "艾莉絲 发件人",
					Address: "alice.fajianren@example.com",
				},
				{
					Name:    "艾莉絲 发件人",
					Address: "alice.fajianren@example.net",
				},
			},
			To: []*mail.Address{
				{
					Name:    "與鮑伯 收件人",
					Address: "bob.shoujianren@example.com",
				},
				{
					Name:    "卡罗尔 收件人",
					Address: "carol.shoujianren@example.com",
				},
			},
			Cc: []*mail.Address{
				{
					Name:    "戴夫 收件人",
					Address: "dave.shoujianren@example.com",
				},
				{
					Name:    "伊夫 收件人",
					Address: "eve.shoujianren@example.com",
				},
			},
			Bcc: []*mail.Address{
				{
					Name:    "艾萨克 伊夫",
					Address: "isaac.shoujianren@example.com",
				},
				{
					Name:    "賈斯汀 伊夫",
					Address: "justin.shoujianren@example.com",
				},
			},
			MessageID:  "Message-Id-1@example.com",
			InReplyTo:  []MessageId{"Message-Id-0@example.com"},
			References: []MessageId{"Message-Id-0@example.com"},
			Comments:   "Message Header Comment",
			Keywords:   []string{"Keyword 1", "Keyword 2"},
			ResentDate: expectedDate,
			ResentFrom: []*mail.Address{
				{
					Name:    "艾莉絲 发件人",
					Address: "alice.fajianren@example.net",
				},
				{
					Name:    "艾莉絲 发件人",
					Address: "alice.fajianren@example.com",
				},
			},
			ResentSender: &mail.Address{
				Name:    "艾莉絲 发件人",
				Address: "alice.fajianren@example.net",
			},
			ResentTo: []*mail.Address{
				{
					Name:    "與鮑伯 收件人",
					Address: "bob.shoujianren@example.net",
				},
				{
					Name:    "卡罗尔 收件人",
					Address: "carol.shoujianren@example.net",
				},
			},
			ResentCc: []*mail.Address{
				{
					Name:    "戴夫 收件人",
					Address: "dave.shoujianren@example.net",
				},
				{
					Name:    "伊夫 收件人",
					Address: "eve.shoujianren@example.net",
				},
			},
			ResentBcc: []*mail.Address{
				{
					Name:    "艾萨克 伊夫",
					Address: "isaac.shoujianren@example.net",
				},
				{
					Name:    "賈斯汀 伊夫",
					Address: "justin.shoujianren@example.net",
				},
			},
			ResentMessageID: "Message-Id-1@example.net",
			ContentType: ContentTypeHeader{
				ContentType: "multipart/signed",
				Params: map[string]string{
					"boundary": "SignedBoundaryString",
					"charset":  "gbk",
					"micalg":   "sha1",
					"protocol": "application/pkcs7-signature",
				},
			},
			ExtraHeaders: map[string][]string{
				"X-Clacks-Overhead": {"GNU Terry Pratchett"},
				"X-Script/function/\t !\"#$%&'()*+,-./;<=>?@[\\]^_`{|}~": {
					"TEST VALUE 1\t !\"#$%&'()*+,-./0123456789:;<=>?@ABCDEFGHIJKLMNOPQRSTUVWXYZ[\\]^_` abcdefghijklmnopqrstuvwxyz{|}~",
					"TEST VALUE 2\t !\"#$%&'()*+,-./0123456789:;<=>?@ABCDEFGHIJKLMNOPQRSTUVWXYZ[\\]^_` abcdefghijklmnopqrstuvwxyz{|}~",
				},
			},
		},
		Text: `石室诗士施氏，嗜狮，誓食十狮。
氏时时适市视狮。
十时，适十狮适市。
是时，适施氏适市。
氏视是十狮，恃矢势，使是十狮逝世。
氏拾是十狮尸，适石室。
石室湿，氏使侍拭石室。
石室拭，氏始试食是十狮。
食时，始识是十狮尸，实十石狮尸。
试释是事。`,
		AttachedFiles: []AttachedFile{
			{
				ContentType: ContentTypeHeader{
					ContentType: "application/pkcs7-signature",
					Params: map[string]string{
						"name": "smime.p7s",
					},
				},
				ContentDisposition: ContentDispositionHeader{
					ContentDisposition: attachment,
					Params: map[string]string{
						"filename": "smime.p7s",
					},
				},
				Data: []byte{130, 28, 135, 132, 117, 46, 142, 18, 97, 140, 126, 251, 159, 193, 199, 25, 58, 223, 189, 185, 227, 239, 158, 173, 108, 31, 71, 27, 133, 80, 165, 252, 133, 227, 174, 198, 132, 129, 159, 29, 246, 19, 235, 133, 80, 165, 252, 133, 227, 174, 198, 132, 129, 159, 29, 246, 19, 234, 49, 251, 238, 127, 7, 28, 104, 33, 200, 120, 71, 82, 232, 225, 38, 30, 249, 234, 214, 193, 244, 113, 147, 173, 251, 219, 158, 57, 252, 28, 113, 147, 173, 251, 225, 38, 24, 199, 239, 190, 173, 108, 31, 71, 27, 133, 80, 110, 120, 251, 231, 174, 198, 132, 129, 159, 29, 246, 19, 234, 8, 114, 30, 17, 212, 186, 58, 95, 200, 94, 59, 26, 18, 6, 124, 119, 216, 79, 174, 21, 65, 185, 227, 239, 158},
			},
		},
	}

	testEmailFromFile(t, fp, expectedEmail)
}

func TestParseEmailChineseMultipartSignedGbkOverQuotedprintable(t *testing.T) {
	fp := "tests/test_chinese_multipart_signed_gbk_over_quoted-printable.txt"
	tz, _ := time.LoadLocation("Asia/Shanghai")
	expectedDate, _ := time.Parse(
		time.RFC1123Z+" (MST)",
		time.Date(2019, time.April, 1, 7, 55, 0, 0, tz).Format(time.RFC1123Z+" (MST)"))
	expectedEmail := Email{
		Headers: Headers{
			Date:    expectedDate,
			Subject: "📧 Signed Test 施氏食狮史",
			ReplyTo: []*mail.Address{
				{
					Name:    "艾莉絲 发件人",
					Address: "alice.fajianren@example.net",
				},
			},
			Sender: &mail.Address{
				Name:    "艾莉絲 发件人",
				Address: "alice.fajianren@example.com",
			},
			From: []*mail.Address{
				{
					Name:    "艾莉絲 发件人",
					Address: "alice.fajianren@example.com",
				},
				{
					Name:    "艾莉絲 发件人",
					Address: "alice.fajianren@example.net",
				},
			},
			To: []*mail.Address{
				{
					Name:    "與鮑伯 收件人",
					Address: "bob.shoujianren@example.com",
				},
				{
					Name:    "卡罗尔 收件人",
					Address: "carol.shoujianren@example.com",
				},
			},
			Cc: []*mail.Address{
				{
					Name:    "戴夫 收件人",
					Address: "dave.shoujianren@example.com",
				},
				{
					Name:    "伊夫 收件人",
					Address: "eve.shoujianren@example.com",
				},
			},
			Bcc: []*mail.Address{
				{
					Name:    "艾萨克 伊夫",
					Address: "isaac.shoujianren@example.com",
				},
				{
					Name:    "賈斯汀 伊夫",
					Address: "justin.shoujianren@example.com",
				},
			},
			MessageID:  "Message-Id-1@example.com",
			InReplyTo:  []MessageId{"Message-Id-0@example.com"},
			References: []MessageId{"Message-Id-0@example.com"},
			Comments:   "Message Header Comment",
			Keywords:   []string{"Keyword 1", "Keyword 2"},
			ResentDate: expectedDate,
			ResentFrom: []*mail.Address{
				{
					Name:    "艾莉絲 发件人",
					Address: "alice.fajianren@example.net",
				},
				{
					Name:    "艾莉絲 发件人",
					Address: "alice.fajianren@example.com",
				},
			},
			ResentSender: &mail.Address{
				Name:    "艾莉絲 发件人",
				Address: "alice.fajianren@example.net",
			},
			ResentTo: []*mail.Address{
				{
					Name:    "與鮑伯 收件人",
					Address: "bob.shoujianren@example.net",
				},
				{
					Name:    "卡罗尔 收件人",
					Address: "carol.shoujianren@example.net",
				},
			},
			ResentCc: []*mail.Address{
				{
					Name:    "戴夫 收件人",
					Address: "dave.shoujianren@example.net",
				},
				{
					Name:    "伊夫 收件人",
					Address: "eve.shoujianren@example.net",
				},
			},
			ResentBcc: []*mail.Address{
				{
					Name:    "艾萨克 伊夫",
					Address: "isaac.shoujianren@example.net",
				},
				{
					Name:    "賈斯汀 伊夫",
					Address: "justin.shoujianren@example.net",
				},
			},
			ResentMessageID: "Message-Id-1@example.net",
			ContentType: ContentTypeHeader{
				ContentType: "multipart/signed",
				Params: map[string]string{
					"boundary": "SignedBoundaryString",
					"charset":  "gbk",
					"micalg":   "sha1",
					"protocol": "application/pkcs7-signature",
				},
			},
			ExtraHeaders: map[string][]string{
				"X-Clacks-Overhead": {"GNU Terry Pratchett"},
				"X-Script/function/\t !\"#$%&'()*+,-./;<=>?@[\\]^_`{|}~": {
					"TEST VALUE 1\t !\"#$%&'()*+,-./0123456789:;<=>?@ABCDEFGHIJKLMNOPQRSTUVWXYZ[\\]^_` abcdefghijklmnopqrstuvwxyz{|}~",
					"TEST VALUE 2\t !\"#$%&'()*+,-./0123456789:;<=>?@ABCDEFGHIJKLMNOPQRSTUVWXYZ[\\]^_` abcdefghijklmnopqrstuvwxyz{|}~",
				},
			},
		},
		Text: `石室诗士施氏，嗜狮，誓食十狮。
氏时时适市视狮。
十时，适十狮适市。
是时，适施氏适市。
氏视是十狮，恃矢势，使是十狮逝世。
氏拾是十狮尸，适石室。
石室湿，氏使侍拭石室。
石室拭，氏始试食是十狮。
食时，始识是十狮尸，实十石狮尸。
试释是事。`,
		AttachedFiles: []AttachedFile{
			{
				ContentType: ContentTypeHeader{
					ContentType: "application/pkcs7-signature",
					Params: map[string]string{
						"name": "smime.p7s",
					},
				},
				ContentDisposition: ContentDispositionHeader{
					ContentDisposition: attachment,
					Params: map[string]string{
						"filename": "smime.p7s",
					},
				},
				Data: []byte{130, 28, 135, 132, 117, 46, 142, 18, 97, 140, 126, 251, 159, 193, 199, 25, 58, 223, 189, 185, 227, 239, 158, 173, 108, 31, 71, 27, 133, 80, 165, 252, 133, 227, 174, 198, 132, 129, 159, 29, 246, 19, 235, 133, 80, 165, 252, 133, 227, 174, 198, 132, 129, 159, 29, 246, 19, 234, 49, 251, 238, 127, 7, 28, 104, 33, 200, 120, 71, 82, 232, 225, 38, 30, 249, 234, 214, 193, 244, 113, 147, 173, 251, 219, 158, 57, 252, 28, 113, 147, 173, 251, 225, 38, 24, 199, 239, 190, 173, 108, 31, 71, 27, 133, 80, 110, 120, 251, 231, 174, 198, 132, 129, 159, 29, 246, 19, 234, 8, 114, 30, 17, 212, 186, 58, 95, 200, 94, 59, 26, 18, 6, 124, 119, 216, 79, 174, 21, 65, 185, 227, 239, 158},
			},
		},
	}

	testEmailFromFile(t, fp, expectedEmail)
}

func TestParseEmailFinnishPlaintextUtf8OverBase64(t *testing.T) {
	fp := "tests/test_finnish_plaintext_utf-8_over_base64.txt"
	tz, _ := time.LoadLocation("Europe/Helsinki")
	expectedDate, _ := time.Parse(
		time.RFC1123Z+" (MST)",
		time.Date(2019, time.April, 1, 7, 55, 0, 0, tz).Format(time.RFC1123Z+" (MST)"))
	expectedEmail := Email{
		Headers: Headers{
			Date:    expectedDate,
			Subject: "📧 Test Suomenkieliset pangrammit",
			ReplyTo: []*mail.Address{
				{
					Name:    "Alice Lähettäjä",
					Address: "alice.lahettaja@example.net",
				},
			},
			Sender: &mail.Address{
				Name:    "Alice Lähettäjä",
				Address: "alice.lahettaja@example.com",
			},
			From: []*mail.Address{
				{
					Name:    "Alice Lähettäjä",
					Address: "alice.lahettaja@example.com",
				},
				{
					Name:    "Alice Lähettäjä",
					Address: "alice.lahettaja@example.net",
				},
			},
			To: []*mail.Address{
				{
					Name:    "Bob Vastaanottaja",
					Address: "bob.vastaanottaja@exaple.com",
				},
				{
					Name:    "Carol Vastaanottaja",
					Address: "carol.vastaanottaja@example.com",
				},
			},
			Cc: []*mail.Address{
				{
					Name:    "Dan Vastaanottaja",
					Address: "dan.vastaanottaja@example.com",
				},
				{
					Name:    "Eve Vastaanottaja",
					Address: "eve.vastaanottaja@example.com",
				},
			},
			Bcc: []*mail.Address{
				{
					Name:    "Frank Vastaanottaja",
					Address: "frank.vastaanottaja@example.com",
				},
				{
					Name:    "Grace Vastaanottaja",
					Address: "grace.vastaanottaja@example.com",
				},
			},
			MessageID:  "Message-Id-1@example.com",
			InReplyTo:  []MessageId{"Message-Id-0@example.com"},
			References: []MessageId{"Message-Id-0@example.com"},
			Comments:   "Message Header Comment",
			Keywords:   []string{"Keyword 1", "Keyword 2"},
			ResentDate: expectedDate,
			ResentFrom: []*mail.Address{
				{
					Name:    "Alice Lähettäjä",
					Address: "alice.lahettaja@example.net",
				},
				{
					Name:    "Alice Lähettäjä",
					Address: "alice.lahettaja@example.com",
				},
			},
			ResentSender: &mail.Address{
				Name:    "Alice Lähettäjä",
				Address: "alice.lahettaja@example.net",
			},
			ResentTo: []*mail.Address{
				{
					Name:    "Bob Vastaanottaja",
					Address: "bob.vastaanottaja@exaple.com",
				},
				{
					Name:    "Carol Vastaanottaja",
					Address: "carol.vastaanottaja@example.net",
				},
			},
			ResentCc: []*mail.Address{
				{
					Name:    "Dan Vastaanottaja",
					Address: "dan.vastaanottaja@example.net",
				},
				{
					Name:    "Eve Vastaanottaja",
					Address: "eve.vastaanottaja@example.net",
				},
			},
			ResentBcc: []*mail.Address{
				{
					Name:    "Frank Vastaanottaja",
					Address: "frank.vastaanottaja@example.net",
				},
				{
					Name:    "Grace Vastaanottaja",
					Address: "grace.vastaanottaja@example.net",
				},
			},
			ResentMessageID: "Message-Id-1@example.net",
			ContentType: ContentTypeHeader{
				ContentType: "text/plain",
				Params: map[string]string{
					"charset": "utf-8",
				},
			},
			ExtraHeaders: map[string][]string{
				"X-Clacks-Overhead": {"GNU Terry Pratchett"},
				"X-Script/function/\t !\"#$%&'()*+,-./;<=>?@[\\]^_`{|}~": {
					"TEST VALUE 1\t !\"#$%&'()*+,-./0123456789:;<=>?@ABCDEFGHIJKLMNOPQRSTUVWXYZ[\\]^_` abcdefghijklmnopqrstuvwxyz{|}~",
					"TEST VALUE 2\t !\"#$%&'()*+,-./0123456789:;<=>?@ABCDEFGHIJKLMNOPQRSTUVWXYZ[\\]^_` abcdefghijklmnopqrstuvwxyz{|}~",
				},
			},
		},
		Text: `Albert osti fagotin ja töräytti puhkuvan melodian.
Lorun sangen pieneksi hyödyksi jäivät suomen kirjaimet.
Hyvän lorun sangen pieneksi hyödyksi jäi suomen kirjaimet.
Fahrenheit ja Celsius yrjösivät Åsan backgammon-peliin, Volkswagenissa, daiquirin ja ZX81:n yhteisvaikutuksesta.
Charles Darwin jammaili Åken hevixylofonilla Qatarin yöpub Zeligissä.
Wieniläinen sioux:ta puhuva ökyzombie diggaa Åsan roquefort-tacoja.`,
	}

	testEmailFromFile(t, fp, expectedEmail)
}

func TestParseEmailFinnishPlaintextUtf8OverQuotedprintable(t *testing.T) {
	fp := "tests/test_finnish_plaintext_utf-8_over_quoted-printable.txt"
	tz, _ := time.LoadLocation("Europe/Helsinki")
	expectedDate, _ := time.Parse(
		time.RFC1123Z+" (MST)",
		time.Date(2019, time.April, 1, 7, 55, 0, 0, tz).Format(time.RFC1123Z+" (MST)"))
	expectedEmail := Email{
		Headers: Headers{
			Date:    expectedDate,
			Subject: "📧 Test Suomenkieliset pangrammit",
			ReplyTo: []*mail.Address{
				{
					Name:    "Alice Lähettäjä",
					Address: "alice.lahettaja@example.net",
				},
			},
			Sender: &mail.Address{
				Name:    "Alice Lähettäjä",
				Address: "alice.lahettaja@example.com",
			},
			From: []*mail.Address{
				{
					Name:    "Alice Lähettäjä",
					Address: "alice.lahettaja@example.com",
				},
				{
					Name:    "Alice Lähettäjä",
					Address: "alice.lahettaja@example.net",
				},
			},
			To: []*mail.Address{
				{
					Name:    "Bob Vastaanottaja",
					Address: "bob.vastaanottaja@exaple.com",
				},
				{
					Name:    "Carol Vastaanottaja",
					Address: "carol.vastaanottaja@example.com",
				},
			},
			Cc: []*mail.Address{
				{
					Name:    "Dan Vastaanottaja",
					Address: "dan.vastaanottaja@example.com",
				},
				{
					Name:    "Eve Vastaanottaja",
					Address: "eve.vastaanottaja@example.com",
				},
			},
			Bcc: []*mail.Address{
				{
					Name:    "Frank Vastaanottaja",
					Address: "frank.vastaanottaja@example.com",
				},
				{
					Name:    "Grace Vastaanottaja",
					Address: "grace.vastaanottaja@example.com",
				},
			},
			MessageID:  "Message-Id-1@example.com",
			InReplyTo:  []MessageId{"Message-Id-0@example.com"},
			References: []MessageId{"Message-Id-0@example.com"},
			Comments:   "Message Header Comment",
			Keywords:   []string{"Keyword 1", "Keyword 2"},
			ResentDate: expectedDate,
			ResentFrom: []*mail.Address{
				{
					Name:    "Alice Lähettäjä",
					Address: "alice.lahettaja@example.net",
				},
				{
					Name:    "Alice Lähettäjä",
					Address: "alice.lahettaja@example.com",
				},
			},
			ResentSender: &mail.Address{
				Name:    "Alice Lähettäjä",
				Address: "alice.lahettaja@example.net",
			},
			ResentTo: []*mail.Address{
				{
					Name:    "Bob Vastaanottaja",
					Address: "bob.vastaanottaja@exaple.com",
				},
				{
					Name:    "Carol Vastaanottaja",
					Address: "carol.vastaanottaja@example.net",
				},
			},
			ResentCc: []*mail.Address{
				{
					Name:    "Dan Vastaanottaja",
					Address: "dan.vastaanottaja@example.net",
				},
				{
					Name:    "Eve Vastaanottaja",
					Address: "eve.vastaanottaja@example.net",
				},
			},
			ResentBcc: []*mail.Address{
				{
					Name:    "Frank Vastaanottaja",
					Address: "frank.vastaanottaja@example.net",
				},
				{
					Name:    "Grace Vastaanottaja",
					Address: "grace.vastaanottaja@example.net",
				},
			},
			ResentMessageID: "Message-Id-1@example.net",
			ContentType: ContentTypeHeader{
				ContentType: "text/plain",
				Params: map[string]string{
					"charset": "utf-8",
				},
			},
			ExtraHeaders: map[string][]string{
				"X-Clacks-Overhead": {"GNU Terry Pratchett"},
				"X-Script/function/\t !\"#$%&'()*+,-./;<=>?@[\\]^_`{|}~": {
					"TEST VALUE 1\t !\"#$%&'()*+,-./0123456789:;<=>?@ABCDEFGHIJKLMNOPQRSTUVWXYZ[\\]^_` abcdefghijklmnopqrstuvwxyz{|}~",
					"TEST VALUE 2\t !\"#$%&'()*+,-./0123456789:;<=>?@ABCDEFGHIJKLMNOPQRSTUVWXYZ[\\]^_` abcdefghijklmnopqrstuvwxyz{|}~",
				},
			},
		},
		Text: `Albert osti fagotin ja töräytti puhkuvan melodian.
Lorun sangen pieneksi hyödyksi jäivät suomen kirjaimet.
Hyvän lorun sangen pieneksi hyödyksi jäi suomen kirjaimet.
Fahrenheit ja Celsius yrjösivät Åsan backgammon-peliin, Volkswagenissa, daiquirin ja ZX81:n yhteisvaikutuksesta.
Charles Darwin jammaili Åken hevixylofonilla Qatarin yöpub Zeligissä.
Wieniläinen sioux:ta puhuva ökyzombie diggaa Åsan roquefort-tacoja.`,
	}

	testEmailFromFile(t, fp, expectedEmail)
}

func TestParseEmailFinnishPlaintextIso885915OverBase64(t *testing.T) {
	fp := "tests/test_finnish_plaintext_iso-8859-15_over_base64.txt"
	tz, _ := time.LoadLocation("Europe/Helsinki")
	expectedDate, _ := time.Parse(
		time.RFC1123Z+" (MST)",
		time.Date(2019, time.April, 1, 7, 55, 0, 0, tz).Format(time.RFC1123Z+" (MST)"))
	expectedEmail := Email{
		Headers: Headers{
			Date:    expectedDate,
			Subject: "📧 Test Suomenkieliset pangrammit",
			ReplyTo: []*mail.Address{
				{
					Name:    "Alice Lähettäjä",
					Address: "alice.lahettaja@example.net",
				},
			},
			Sender: &mail.Address{
				Name:    "Alice Lähettäjä",
				Address: "alice.lahettaja@example.com",
			},
			From: []*mail.Address{
				{
					Name:    "Alice Lähettäjä",
					Address: "alice.lahettaja@example.com",
				},
				{
					Name:    "Alice Lähettäjä",
					Address: "alice.lahettaja@example.net",
				},
			},
			To: []*mail.Address{
				{
					Name:    "Bob Vastaanottaja",
					Address: "bob.vastaanottaja@exaple.com",
				},
				{
					Name:    "Carol Vastaanottaja",
					Address: "carol.vastaanottaja@example.com",
				},
			},
			Cc: []*mail.Address{
				{
					Name:    "Dan Vastaanottaja",
					Address: "dan.vastaanottaja@example.com",
				},
				{
					Name:    "Eve Vastaanottaja",
					Address: "eve.vastaanottaja@example.com",
				},
			},
			Bcc: []*mail.Address{
				{
					Name:    "Frank Vastaanottaja",
					Address: "frank.vastaanottaja@example.com",
				},
				{
					Name:    "Grace Vastaanottaja",
					Address: "grace.vastaanottaja@example.com",
				},
			},
			MessageID:  "Message-Id-1@example.com",
			InReplyTo:  []MessageId{"Message-Id-0@example.com"},
			References: []MessageId{"Message-Id-0@example.com"},
			Comments:   "Message Header Comment",
			Keywords:   []string{"Keyword 1", "Keyword 2"},
			ResentDate: expectedDate,
			ResentFrom: []*mail.Address{
				{
					Name:    "Alice Lähettäjä",
					Address: "alice.lahettaja@example.net",
				},
				{
					Name:    "Alice Lähettäjä",
					Address: "alice.lahettaja@example.com",
				},
			},
			ResentSender: &mail.Address{
				Name:    "Alice Lähettäjä",
				Address: "alice.lahettaja@example.net",
			},
			ResentTo: []*mail.Address{
				{
					Name:    "Bob Vastaanottaja",
					Address: "bob.vastaanottaja@exaple.com",
				},
				{
					Name:    "Carol Vastaanottaja",
					Address: "carol.vastaanottaja@example.net",
				},
			},
			ResentCc: []*mail.Address{
				{
					Name:    "Dan Vastaanottaja",
					Address: "dan.vastaanottaja@example.net",
				},
				{
					Name:    "Eve Vastaanottaja",
					Address: "eve.vastaanottaja@example.net",
				},
			},
			ResentBcc: []*mail.Address{
				{
					Name:    "Frank Vastaanottaja",
					Address: "frank.vastaanottaja@example.net",
				},
				{
					Name:    "Grace Vastaanottaja",
					Address: "grace.vastaanottaja@example.net",
				},
			},
			ResentMessageID: "Message-Id-1@example.net",
			ContentType: ContentTypeHeader{
				ContentType: "text/plain",
				Params: map[string]string{
					"charset": "iso-8859-15",
				},
			},
			ExtraHeaders: map[string][]string{
				"X-Clacks-Overhead": {"GNU Terry Pratchett"},
				"X-Script/function/\t !\"#$%&'()*+,-./;<=>?@[\\]^_`{|}~": {
					"TEST VALUE 1\t !\"#$%&'()*+,-./0123456789:;<=>?@ABCDEFGHIJKLMNOPQRSTUVWXYZ[\\]^_` abcdefghijklmnopqrstuvwxyz{|}~",
					"TEST VALUE 2\t !\"#$%&'()*+,-./0123456789:;<=>?@ABCDEFGHIJKLMNOPQRSTUVWXYZ[\\]^_` abcdefghijklmnopqrstuvwxyz{|}~",
				},
			},
		},
		Text: `Albert osti fagotin ja töräytti puhkuvan melodian.
Lorun sangen pieneksi hyödyksi jäivät suomen kirjaimet.
Hyvän lorun sangen pieneksi hyödyksi jäi suomen kirjaimet.
Fahrenheit ja Celsius yrjösivät Åsan backgammon-peliin, Volkswagenissa, daiquirin ja ZX81:n yhteisvaikutuksesta.
Charles Darwin jammaili Åken hevixylofonilla Qatarin yöpub Zeligissä.
Wieniläinen sioux:ta puhuva ökyzombie diggaa Åsan roquefort-tacoja.`,
	}

	testEmailFromFile(t, fp, expectedEmail)
}

func TestParseEmailFinnishPlaintextIso885915OverQuotedprintable(t *testing.T) {
	fp := "tests/test_finnish_plaintext_iso-8859-15_over_quoted-printable.txt"
	tz, _ := time.LoadLocation("Europe/Helsinki")
	expectedDate, _ := time.Parse(
		time.RFC1123Z+" (MST)",
		time.Date(2019, time.April, 1, 7, 55, 0, 0, tz).Format(time.RFC1123Z+" (MST)"))
	expectedEmail := Email{
		Headers: Headers{
			Date:    expectedDate,
			Subject: "📧 Test Suomenkieliset pangrammit",
			ReplyTo: []*mail.Address{
				{
					Name:    "Alice Lähettäjä",
					Address: "alice.lahettaja@example.net",
				},
			},
			Sender: &mail.Address{
				Name:    "Alice Lähettäjä",
				Address: "alice.lahettaja@example.com",
			},
			From: []*mail.Address{
				{
					Name:    "Alice Lähettäjä",
					Address: "alice.lahettaja@example.com",
				},
				{
					Name:    "Alice Lähettäjä",
					Address: "alice.lahettaja@example.net",
				},
			},
			To: []*mail.Address{
				{
					Name:    "Bob Vastaanottaja",
					Address: "bob.vastaanottaja@exaple.com",
				},
				{
					Name:    "Carol Vastaanottaja",
					Address: "carol.vastaanottaja@example.com",
				},
			},
			Cc: []*mail.Address{
				{
					Name:    "Dan Vastaanottaja",
					Address: "dan.vastaanottaja@example.com",
				},
				{
					Name:    "Eve Vastaanottaja",
					Address: "eve.vastaanottaja@example.com",
				},
			},
			Bcc: []*mail.Address{
				{
					Name:    "Frank Vastaanottaja",
					Address: "frank.vastaanottaja@example.com",
				},
				{
					Name:    "Grace Vastaanottaja",
					Address: "grace.vastaanottaja@example.com",
				},
			},
			MessageID:  "Message-Id-1@example.com",
			InReplyTo:  []MessageId{"Message-Id-0@example.com"},
			References: []MessageId{"Message-Id-0@example.com"},
			Comments:   "Message Header Comment",
			Keywords:   []string{"Keyword 1", "Keyword 2"},
			ResentDate: expectedDate,
			ResentFrom: []*mail.Address{
				{
					Name:    "Alice Lähettäjä",
					Address: "alice.lahettaja@example.net",
				},
				{
					Name:    "Alice Lähettäjä",
					Address: "alice.lahettaja@example.com",
				},
			},
			ResentSender: &mail.Address{
				Name:    "Alice Lähettäjä",
				Address: "alice.lahettaja@example.net",
			},
			ResentTo: []*mail.Address{
				{
					Name:    "Bob Vastaanottaja",
					Address: "bob.vastaanottaja@exaple.com",
				},
				{
					Name:    "Carol Vastaanottaja",
					Address: "carol.vastaanottaja@example.net",
				},
			},
			ResentCc: []*mail.Address{
				{
					Name:    "Dan Vastaanottaja",
					Address: "dan.vastaanottaja@example.net",
				},
				{
					Name:    "Eve Vastaanottaja",
					Address: "eve.vastaanottaja@example.net",
				},
			},
			ResentBcc: []*mail.Address{
				{
					Name:    "Frank Vastaanottaja",
					Address: "frank.vastaanottaja@example.net",
				},
				{
					Name:    "Grace Vastaanottaja",
					Address: "grace.vastaanottaja@example.net",
				},
			},
			ResentMessageID: "Message-Id-1@example.net",
			ContentType: ContentTypeHeader{
				ContentType: "text/plain",
				Params: map[string]string{
					"charset": "iso-8859-15",
				},
			},
			ExtraHeaders: map[string][]string{
				"X-Clacks-Overhead": {"GNU Terry Pratchett"},
				"X-Script/function/\t !\"#$%&'()*+,-./;<=>?@[\\]^_`{|}~": {
					"TEST VALUE 1\t !\"#$%&'()*+,-./0123456789:;<=>?@ABCDEFGHIJKLMNOPQRSTUVWXYZ[\\]^_` abcdefghijklmnopqrstuvwxyz{|}~",
					"TEST VALUE 2\t !\"#$%&'()*+,-./0123456789:;<=>?@ABCDEFGHIJKLMNOPQRSTUVWXYZ[\\]^_` abcdefghijklmnopqrstuvwxyz{|}~",
				},
			},
		},
		Text: `Albert osti fagotin ja töräytti puhkuvan melodian.
Lorun sangen pieneksi hyödyksi jäivät suomen kirjaimet.
Hyvän lorun sangen pieneksi hyödyksi jäi suomen kirjaimet.
Fahrenheit ja Celsius yrjösivät Åsan backgammon-peliin, Volkswagenissa, daiquirin ja ZX81:n yhteisvaikutuksesta.
Charles Darwin jammaili Åken hevixylofonilla Qatarin yöpub Zeligissä.
Wieniläinen sioux:ta puhuva ökyzombie diggaa Åsan roquefort-tacoja.`,
	}

	testEmailFromFile(t, fp, expectedEmail)
}

func TestParseEmailFinnishMultipartRelatedUtf8OverBase64(t *testing.T) {
	fp := "tests/test_finnish_multipart_related_utf-8_over_base64.txt"
	tz, _ := time.LoadLocation("Europe/Helsinki")
	expectedDate, _ := time.Parse(
		time.RFC1123Z+" (MST)",
		time.Date(2019, time.April, 1, 7, 55, 0, 0, tz).Format(time.RFC1123Z+" (MST)"))
	expectedEmail := Email{
		Headers: Headers{
			Date:    expectedDate,
			Subject: "📧 Test Suomenkieliset pangrammit",
			ReplyTo: []*mail.Address{
				{
					Name:    "Alice Lähettäjä",
					Address: "alice.lahettaja@example.net",
				},
			},
			Sender: &mail.Address{
				Name:    "Alice Lähettäjä",
				Address: "alice.lahettaja@example.com",
			},
			From: []*mail.Address{
				{
					Name:    "Alice Lähettäjä",
					Address: "alice.lahettaja@example.com",
				},
				{
					Name:    "Alice Lähettäjä",
					Address: "alice.lahettaja@example.net",
				},
			},
			To: []*mail.Address{
				{
					Name:    "Bob Vastaanottaja",
					Address: "bob.vastaanottaja@exaple.com",
				},
				{
					Name:    "Carol Vastaanottaja",
					Address: "carol.vastaanottaja@example.com",
				},
			},
			Cc: []*mail.Address{
				{
					Name:    "Dan Vastaanottaja",
					Address: "dan.vastaanottaja@example.com",
				},
				{
					Name:    "Eve Vastaanottaja",
					Address: "eve.vastaanottaja@example.com",
				},
			},
			Bcc: []*mail.Address{
				{
					Name:    "Frank Vastaanottaja",
					Address: "frank.vastaanottaja@example.com",
				},
				{
					Name:    "Grace Vastaanottaja",
					Address: "grace.vastaanottaja@example.com",
				},
			},
			MessageID:  "Message-Id-1@example.com",
			InReplyTo:  []MessageId{"Message-Id-0@example.com"},
			References: []MessageId{"Message-Id-0@example.com"},
			Comments:   "Message Header Comment",
			Keywords:   []string{"Keyword 1", "Keyword 2"},
			ResentDate: expectedDate,
			ResentFrom: []*mail.Address{
				{
					Name:    "Alice Lähettäjä",
					Address: "alice.lahettaja@example.net",
				},
				{
					Name:    "Alice Lähettäjä",
					Address: "alice.lahettaja@example.com",
				},
			},
			ResentSender: &mail.Address{
				Name:    "Alice Lähettäjä",
				Address: "alice.lahettaja@example.net",
			},
			ResentTo: []*mail.Address{
				{
					Name:    "Bob Vastaanottaja",
					Address: "bob.vastaanottaja@exaple.com",
				},
				{
					Name:    "Carol Vastaanottaja",
					Address: "carol.vastaanottaja@example.net",
				},
			},
			ResentCc: []*mail.Address{
				{
					Name:    "Dan Vastaanottaja",
					Address: "dan.vastaanottaja@example.net",
				},
				{
					Name:    "Eve Vastaanottaja",
					Address: "eve.vastaanottaja@example.net",
				},
			},
			ResentBcc: []*mail.Address{
				{
					Name:    "Frank Vastaanottaja",
					Address: "frank.vastaanottaja@example.net",
				},
				{
					Name:    "Grace Vastaanottaja",
					Address: "grace.vastaanottaja@example.net",
				},
			},
			ResentMessageID: "Message-Id-1@example.net",
			ContentType: ContentTypeHeader{
				ContentType: "multipart/related",
				Params: map[string]string{
					"boundary": "RelatedBoundaryString",
					"charset":  "utf-8",
				},
			},
			ExtraHeaders: map[string][]string{
				"X-Clacks-Overhead": {"GNU Terry Pratchett"},
				"X-Script/function/\t !\"#$%&'()*+,-./;<=>?@[\\]^_`{|}~": {
					"TEST VALUE 1\t !\"#$%&'()*+,-./0123456789:;<=>?@ABCDEFGHIJKLMNOPQRSTUVWXYZ[\\]^_` abcdefghijklmnopqrstuvwxyz{|}~",
					"TEST VALUE 2\t !\"#$%&'()*+,-./0123456789:;<=>?@ABCDEFGHIJKLMNOPQRSTUVWXYZ[\\]^_` abcdefghijklmnopqrstuvwxyz{|}~",
				},
			},
		},
		Text: `Albert osti fagotin ja töräytti puhkuvan melodian.
Lorun sangen pieneksi hyödyksi jäivät suomen kirjaimet.
Hyvän lorun sangen pieneksi hyödyksi jäi suomen kirjaimet.
Fahrenheit ja Celsius yrjösivät Åsan backgammon-peliin, Volkswagenissa, daiquirin ja ZX81:n yhteisvaikutuksesta.
Charles Darwin jammaili Åken hevixylofonilla Qatarin yöpub Zeligissä.
Wieniläinen sioux:ta puhuva ökyzombie diggaa Åsan roquefort-tacoja.`,
		EnrichedText: `<bold>Albert osti fagotin ja töräytti puhkuvan melodian.</bold>
<italic>Lorun sangen pieneksi hyödyksi jäivät suomen kirjaimet.</italic>
<fixed>Hyvän lorun sangen pieneksi hyödyksi jäi suomen kirjaimet.</fixed>
<underline>Fahrenheit ja Celsius yrjösivät Åsan backgammon-peliin, Volkswagenissa, daiquirin ja ZX81:n yhteisvaikutuksesta.</underline>
Charles Darwin jammaili Åken hevixylofonilla Qatarin yöpub Zeligissä.
Wieniläinen sioux:ta puhuva ökyzombie diggaa Åsan roquefort-tacoja.`,
		HTML: `<html>
<div dir="ltr">
<p>Albert osti fagotin ja töräytti puhkuvan melodian.</p>
<p>Lorun sangen pieneksi hyödyksi jäivät suomen kirjaimet.</p>
<p>Hyvän lorun sangen pieneksi hyödyksi jäi suomen kirjaimet.</p>
<p>Fahrenheit ja Celsius yrjösivät Åsan backgammon-peliin, Volkswagenissa, daiquirin ja ZX81:n yhteisvaikutuksesta.</p>
<p>Charles Darwin jammaili Åken hevixylofonilla Qatarin yöpub Zeligissä.</p>
<p>Wieniläinen sioux:ta puhuva ökyzombie diggaa Åsan roquefort-tacoja.</p>
</div>
</html>`,
		InlineFiles: []InlineFile{
			{
				ContentID: "inline-jpg-image.jpg@example.com",
				ContentType: ContentTypeHeader{
					ContentType: "image/jpeg",
					Params: map[string]string{
						"name": "inline-jpg-image-name.jpg",
					},
				},
				ContentDisposition: ContentDispositionHeader{
					ContentDisposition: inline,
					Params: map[string]string{
						"filename": "inline-jpg-image-filename.jpg",
					},
				},
				Data: []byte{255, 216, 255, 219, 0, 67, 0, 3, 2, 2, 2, 2, 2, 3, 2, 2, 2, 3, 3, 3, 3, 4,
					6, 4, 4, 4, 4, 4, 8, 6, 6, 5, 6, 9, 8, 10, 10, 9, 8, 9, 9, 10, 12, 15, 12, 10, 11, 14, 11, 9, 9,
					13, 17, 13, 14, 15, 16, 16, 17, 16, 10, 12, 18, 19, 18, 16, 19, 15, 16, 16, 16, 255, 201, 0, 11, 8,
					0, 1, 0, 1, 1, 1, 17, 0, 255, 204, 0, 6, 0, 16, 16, 5, 255, 218, 0, 8, 1, 1, 0, 0, 63, 0, 210, 207,
					32, 255, 217},
			},
		},
	}

	testEmailFromFile(t, fp, expectedEmail)
}
func TestParseEmailFinnishMultipartRelatedUtf8OverQuotedprintable(t *testing.T) {
	fp := "tests/test_finnish_multipart_related_utf-8_over_quoted-printable.txt"
	tz, _ := time.LoadLocation("Europe/Helsinki")
	expectedDate, _ := time.Parse(
		time.RFC1123Z+" (MST)",
		time.Date(2019, time.April, 1, 7, 55, 0, 0, tz).Format(time.RFC1123Z+" (MST)"))
	expectedEmail := Email{
		Headers: Headers{
			Date:    expectedDate,
			Subject: "📧 Test Suomenkieliset pangrammit",
			ReplyTo: []*mail.Address{
				{
					Name:    "Alice Lähettäjä",
					Address: "alice.lahettaja@example.net",
				},
			},
			Sender: &mail.Address{
				Name:    "Alice Lähettäjä",
				Address: "alice.lahettaja@example.com",
			},
			From: []*mail.Address{
				{
					Name:    "Alice Lähettäjä",
					Address: "alice.lahettaja@example.com",
				},
				{
					Name:    "Alice Lähettäjä",
					Address: "alice.lahettaja@example.net",
				},
			},
			To: []*mail.Address{
				{
					Name:    "Bob Vastaanottaja",
					Address: "bob.vastaanottaja@exaple.com",
				},
				{
					Name:    "Carol Vastaanottaja",
					Address: "carol.vastaanottaja@example.com",
				},
			},
			Cc: []*mail.Address{
				{
					Name:    "Dan Vastaanottaja",
					Address: "dan.vastaanottaja@example.com",
				},
				{
					Name:    "Eve Vastaanottaja",
					Address: "eve.vastaanottaja@example.com",
				},
			},
			Bcc: []*mail.Address{
				{
					Name:    "Frank Vastaanottaja",
					Address: "frank.vastaanottaja@example.com",
				},
				{
					Name:    "Grace Vastaanottaja",
					Address: "grace.vastaanottaja@example.com",
				},
			},
			MessageID:  "Message-Id-1@example.com",
			InReplyTo:  []MessageId{"Message-Id-0@example.com"},
			References: []MessageId{"Message-Id-0@example.com"},
			Comments:   "Message Header Comment",
			Keywords:   []string{"Keyword 1", "Keyword 2"},
			ResentDate: expectedDate,
			ResentFrom: []*mail.Address{
				{
					Name:    "Alice Lähettäjä",
					Address: "alice.lahettaja@example.net",
				},
				{
					Name:    "Alice Lähettäjä",
					Address: "alice.lahettaja@example.com",
				},
			},
			ResentSender: &mail.Address{
				Name:    "Alice Lähettäjä",
				Address: "alice.lahettaja@example.net",
			},
			ResentTo: []*mail.Address{
				{
					Name:    "Bob Vastaanottaja",
					Address: "bob.vastaanottaja@exaple.com",
				},
				{
					Name:    "Carol Vastaanottaja",
					Address: "carol.vastaanottaja@example.net",
				},
			},
			ResentCc: []*mail.Address{
				{
					Name:    "Dan Vastaanottaja",
					Address: "dan.vastaanottaja@example.net",
				},
				{
					Name:    "Eve Vastaanottaja",
					Address: "eve.vastaanottaja@example.net",
				},
			},
			ResentBcc: []*mail.Address{
				{
					Name:    "Frank Vastaanottaja",
					Address: "frank.vastaanottaja@example.net",
				},
				{
					Name:    "Grace Vastaanottaja",
					Address: "grace.vastaanottaja@example.net",
				},
			},
			ResentMessageID: "Message-Id-1@example.net",
			ContentType: ContentTypeHeader{
				ContentType: "multipart/related",
				Params: map[string]string{
					"boundary": "RelatedBoundaryString",
					"charset":  "utf-8",
				},
			},
			ExtraHeaders: map[string][]string{
				"X-Clacks-Overhead": {"GNU Terry Pratchett"},
				"X-Script/function/\t !\"#$%&'()*+,-./;<=>?@[\\]^_`{|}~": {
					"TEST VALUE 1\t !\"#$%&'()*+,-./0123456789:;<=>?@ABCDEFGHIJKLMNOPQRSTUVWXYZ[\\]^_` abcdefghijklmnopqrstuvwxyz{|}~",
					"TEST VALUE 2\t !\"#$%&'()*+,-./0123456789:;<=>?@ABCDEFGHIJKLMNOPQRSTUVWXYZ[\\]^_` abcdefghijklmnopqrstuvwxyz{|}~",
				},
			},
		},
		Text: `Albert osti fagotin ja töräytti puhkuvan melodian.
Lorun sangen pieneksi hyödyksi jäivät suomen kirjaimet.
Hyvän lorun sangen pieneksi hyödyksi jäi suomen kirjaimet.
Fahrenheit ja Celsius yrjösivät Åsan backgammon-peliin, Volkswagenissa, daiquirin ja ZX81:n yhteisvaikutuksesta.
Charles Darwin jammaili Åken hevixylofonilla Qatarin yöpub Zeligissä.
Wieniläinen sioux:ta puhuva ökyzombie diggaa Åsan roquefort-tacoja.`,
		EnrichedText: `<bold>Albert osti fagotin ja töräytti puhkuvan melodian.</bold>
<italic>Lorun sangen pieneksi hyödyksi jäivät suomen kirjaimet.</italic>
<fixed>Hyvän lorun sangen pieneksi hyödyksi jäi suomen kirjaimet.</fixed>
<underline>Fahrenheit ja Celsius yrjösivät Åsan backgammon-peliin, Volkswagenissa, daiquirin ja ZX81:n yhteisvaikutuksesta.</underline>
Charles Darwin jammaili Åken hevixylofonilla Qatarin yöpub Zeligissä.
Wieniläinen sioux:ta puhuva ökyzombie diggaa Åsan roquefort-tacoja.`,
		HTML: `<html>
<div dir="ltr">
<p>Albert osti fagotin ja töräytti puhkuvan melodian.</p>
<p>Lorun sangen pieneksi hyödyksi jäivät suomen kirjaimet.</p>
<p>Hyvän lorun sangen pieneksi hyödyksi jäi suomen kirjaimet.</p>
<p>Fahrenheit ja Celsius yrjösivät Åsan backgammon-peliin, Volkswagenissa, daiquirin ja ZX81:n yhteisvaikutuksesta.</p>
<p>Charles Darwin jammaili Åken hevixylofonilla Qatarin yöpub Zeligissä.</p>
<p>Wieniläinen sioux:ta puhuva ökyzombie diggaa Åsan roquefort-tacoja.</p>
</div>
</html>`,
		InlineFiles: []InlineFile{
			{
				ContentID: "inline-jpg-image.jpg@example.com",
				ContentType: ContentTypeHeader{
					ContentType: "image/jpeg",
					Params: map[string]string{
						"name": "inline-jpg-image-name.jpg",
					},
				},
				ContentDisposition: ContentDispositionHeader{
					ContentDisposition: inline,
					Params: map[string]string{
						"filename": "inline-jpg-image-filename.jpg",
					},
				},
				Data: []byte{255, 216, 255, 219, 0, 67, 0, 3, 2, 2, 2, 2, 2, 3, 2, 2, 2, 3, 3, 3, 3, 4,
					6, 4, 4, 4, 4, 4, 8, 6, 6, 5, 6, 9, 8, 10, 10, 9, 8, 9, 9, 10, 12, 15, 12, 10, 11, 14, 11, 9, 9,
					13, 17, 13, 14, 15, 16, 16, 17, 16, 10, 12, 18, 19, 18, 16, 19, 15, 16, 16, 16, 255, 201, 0, 11, 8,
					0, 1, 0, 1, 1, 1, 17, 0, 255, 204, 0, 6, 0, 16, 16, 5, 255, 218, 0, 8, 1, 1, 0, 0, 63, 0, 210, 207,
					32, 255, 217},
			},
		},
	}

	testEmailFromFile(t, fp, expectedEmail)
}
func TestParseEmailFinnishMultipartRelatedIso885915OverBase64(t *testing.T) {
	fp := "tests/test_finnish_multipart_related_iso-8859-15_over_base64.txt"
	tz, _ := time.LoadLocation("Europe/Helsinki")
	expectedDate, _ := time.Parse(
		time.RFC1123Z+" (MST)",
		time.Date(2019, time.April, 1, 7, 55, 0, 0, tz).Format(time.RFC1123Z+" (MST)"))
	expectedEmail := Email{
		Headers: Headers{
			Date:    expectedDate,
			Subject: "📧 Test Suomenkieliset pangrammit",
			ReplyTo: []*mail.Address{
				{
					Name:    "Alice Lähettäjä",
					Address: "alice.lahettaja@example.net",
				},
			},
			Sender: &mail.Address{
				Name:    "Alice Lähettäjä",
				Address: "alice.lahettaja@example.com",
			},
			From: []*mail.Address{
				{
					Name:    "Alice Lähettäjä",
					Address: "alice.lahettaja@example.com",
				},
				{
					Name:    "Alice Lähettäjä",
					Address: "alice.lahettaja@example.net",
				},
			},
			To: []*mail.Address{
				{
					Name:    "Bob Vastaanottaja",
					Address: "bob.vastaanottaja@exaple.com",
				},
				{
					Name:    "Carol Vastaanottaja",
					Address: "carol.vastaanottaja@example.com",
				},
			},
			Cc: []*mail.Address{
				{
					Name:    "Dan Vastaanottaja",
					Address: "dan.vastaanottaja@example.com",
				},
				{
					Name:    "Eve Vastaanottaja",
					Address: "eve.vastaanottaja@example.com",
				},
			},
			Bcc: []*mail.Address{
				{
					Name:    "Frank Vastaanottaja",
					Address: "frank.vastaanottaja@example.com",
				},
				{
					Name:    "Grace Vastaanottaja",
					Address: "grace.vastaanottaja@example.com",
				},
			},
			MessageID:  "Message-Id-1@example.com",
			InReplyTo:  []MessageId{"Message-Id-0@example.com"},
			References: []MessageId{"Message-Id-0@example.com"},
			Comments:   "Message Header Comment",
			Keywords:   []string{"Keyword 1", "Keyword 2"},
			ResentDate: expectedDate,
			ResentFrom: []*mail.Address{
				{
					Name:    "Alice Lähettäjä",
					Address: "alice.lahettaja@example.net",
				},
				{
					Name:    "Alice Lähettäjä",
					Address: "alice.lahettaja@example.com",
				},
			},
			ResentSender: &mail.Address{
				Name:    "Alice Lähettäjä",
				Address: "alice.lahettaja@example.net",
			},
			ResentTo: []*mail.Address{
				{
					Name:    "Bob Vastaanottaja",
					Address: "bob.vastaanottaja@exaple.com",
				},
				{
					Name:    "Carol Vastaanottaja",
					Address: "carol.vastaanottaja@example.net",
				},
			},
			ResentCc: []*mail.Address{
				{
					Name:    "Dan Vastaanottaja",
					Address: "dan.vastaanottaja@example.net",
				},
				{
					Name:    "Eve Vastaanottaja",
					Address: "eve.vastaanottaja@example.net",
				},
			},
			ResentBcc: []*mail.Address{
				{
					Name:    "Frank Vastaanottaja",
					Address: "frank.vastaanottaja@example.net",
				},
				{
					Name:    "Grace Vastaanottaja",
					Address: "grace.vastaanottaja@example.net",
				},
			},
			ResentMessageID: "Message-Id-1@example.net",
			ContentType: ContentTypeHeader{
				ContentType: "multipart/related",
				Params: map[string]string{
					"boundary": "RelatedBoundaryString",
					"charset":  "iso-8859-15",
				},
			},
			ExtraHeaders: map[string][]string{
				"X-Clacks-Overhead": {"GNU Terry Pratchett"},
				"X-Script/function/\t !\"#$%&'()*+,-./;<=>?@[\\]^_`{|}~": {
					"TEST VALUE 1\t !\"#$%&'()*+,-./0123456789:;<=>?@ABCDEFGHIJKLMNOPQRSTUVWXYZ[\\]^_` abcdefghijklmnopqrstuvwxyz{|}~",
					"TEST VALUE 2\t !\"#$%&'()*+,-./0123456789:;<=>?@ABCDEFGHIJKLMNOPQRSTUVWXYZ[\\]^_` abcdefghijklmnopqrstuvwxyz{|}~",
				},
			},
		},
		Text: `Albert osti fagotin ja töräytti puhkuvan melodian.
Lorun sangen pieneksi hyödyksi jäivät suomen kirjaimet.
Hyvän lorun sangen pieneksi hyödyksi jäi suomen kirjaimet.
Fahrenheit ja Celsius yrjösivät Åsan backgammon-peliin, Volkswagenissa, daiquirin ja ZX81:n yhteisvaikutuksesta.
Charles Darwin jammaili Åken hevixylofonilla Qatarin yöpub Zeligissä.
Wieniläinen sioux:ta puhuva ökyzombie diggaa Åsan roquefort-tacoja.`,
		EnrichedText: `<bold>Albert osti fagotin ja töräytti puhkuvan melodian.</bold>
<italic>Lorun sangen pieneksi hyödyksi jäivät suomen kirjaimet.</italic>
<fixed>Hyvän lorun sangen pieneksi hyödyksi jäi suomen kirjaimet.</fixed>
<underline>Fahrenheit ja Celsius yrjösivät Åsan backgammon-peliin, Volkswagenissa, daiquirin ja ZX81:n yhteisvaikutuksesta.</underline>
Charles Darwin jammaili Åken hevixylofonilla Qatarin yöpub Zeligissä.
Wieniläinen sioux:ta puhuva ökyzombie diggaa Åsan roquefort-tacoja.`,
		HTML: `<html>
<div dir="ltr">
<p>Albert osti fagotin ja töräytti puhkuvan melodian.</p>
<p>Lorun sangen pieneksi hyödyksi jäivät suomen kirjaimet.</p>
<p>Hyvän lorun sangen pieneksi hyödyksi jäi suomen kirjaimet.</p>
<p>Fahrenheit ja Celsius yrjösivät Åsan backgammon-peliin, Volkswagenissa, daiquirin ja ZX81:n yhteisvaikutuksesta.</p>
<p>Charles Darwin jammaili Åken hevixylofonilla Qatarin yöpub Zeligissä.</p>
<p>Wieniläinen sioux:ta puhuva ökyzombie diggaa Åsan roquefort-tacoja.</p>
</div>
</html>`,
		InlineFiles: []InlineFile{
			{
				ContentID: "inline-jpg-image.jpg@example.com",
				ContentType: ContentTypeHeader{
					ContentType: "image/jpeg",
					Params: map[string]string{
						"name": "inline-jpg-image-name.jpg",
					},
				},
				ContentDisposition: ContentDispositionHeader{
					ContentDisposition: inline,
					Params: map[string]string{
						"filename": "inline-jpg-image-filename.jpg",
					},
				},
				Data: []byte{255, 216, 255, 219, 0, 67, 0, 3, 2, 2, 2, 2, 2, 3, 2, 2, 2, 3, 3, 3, 3, 4,
					6, 4, 4, 4, 4, 4, 8, 6, 6, 5, 6, 9, 8, 10, 10, 9, 8, 9, 9, 10, 12, 15, 12, 10, 11, 14, 11, 9, 9,
					13, 17, 13, 14, 15, 16, 16, 17, 16, 10, 12, 18, 19, 18, 16, 19, 15, 16, 16, 16, 255, 201, 0, 11, 8,
					0, 1, 0, 1, 1, 1, 17, 0, 255, 204, 0, 6, 0, 16, 16, 5, 255, 218, 0, 8, 1, 1, 0, 0, 63, 0, 210, 207,
					32, 255, 217},
			},
		},
	}

	testEmailFromFile(t, fp, expectedEmail)
}
func TestParseEmailFinnishMultipartRelatedIso885915OverQuotedprintable(t *testing.T) {
	fp := "tests/test_finnish_multipart_related_iso-8859-15_over_quoted-printable.txt"
	tz, _ := time.LoadLocation("Europe/Helsinki")
	expectedDate, _ := time.Parse(
		time.RFC1123Z+" (MST)",
		time.Date(2019, time.April, 1, 7, 55, 0, 0, tz).Format(time.RFC1123Z+" (MST)"))
	expectedEmail := Email{
		Headers: Headers{
			Date:    expectedDate,
			Subject: "📧 Test Suomenkieliset pangrammit",
			ReplyTo: []*mail.Address{
				{
					Name:    "Alice Lähettäjä",
					Address: "alice.lahettaja@example.net",
				},
			},
			Sender: &mail.Address{
				Name:    "Alice Lähettäjä",
				Address: "alice.lahettaja@example.com",
			},
			From: []*mail.Address{
				{
					Name:    "Alice Lähettäjä",
					Address: "alice.lahettaja@example.com",
				},
				{
					Name:    "Alice Lähettäjä",
					Address: "alice.lahettaja@example.net",
				},
			},
			To: []*mail.Address{
				{
					Name:    "Bob Vastaanottaja",
					Address: "bob.vastaanottaja@exaple.com",
				},
				{
					Name:    "Carol Vastaanottaja",
					Address: "carol.vastaanottaja@example.com",
				},
			},
			Cc: []*mail.Address{
				{
					Name:    "Dan Vastaanottaja",
					Address: "dan.vastaanottaja@example.com",
				},
				{
					Name:    "Eve Vastaanottaja",
					Address: "eve.vastaanottaja@example.com",
				},
			},
			Bcc: []*mail.Address{
				{
					Name:    "Frank Vastaanottaja",
					Address: "frank.vastaanottaja@example.com",
				},
				{
					Name:    "Grace Vastaanottaja",
					Address: "grace.vastaanottaja@example.com",
				},
			},
			MessageID:  "Message-Id-1@example.com",
			InReplyTo:  []MessageId{"Message-Id-0@example.com"},
			References: []MessageId{"Message-Id-0@example.com"},
			Comments:   "Message Header Comment",
			Keywords:   []string{"Keyword 1", "Keyword 2"},
			ResentDate: expectedDate,
			ResentFrom: []*mail.Address{
				{
					Name:    "Alice Lähettäjä",
					Address: "alice.lahettaja@example.net",
				},
				{
					Name:    "Alice Lähettäjä",
					Address: "alice.lahettaja@example.com",
				},
			},
			ResentSender: &mail.Address{
				Name:    "Alice Lähettäjä",
				Address: "alice.lahettaja@example.net",
			},
			ResentTo: []*mail.Address{
				{
					Name:    "Bob Vastaanottaja",
					Address: "bob.vastaanottaja@exaple.com",
				},
				{
					Name:    "Carol Vastaanottaja",
					Address: "carol.vastaanottaja@example.net",
				},
			},
			ResentCc: []*mail.Address{
				{
					Name:    "Dan Vastaanottaja",
					Address: "dan.vastaanottaja@example.net",
				},
				{
					Name:    "Eve Vastaanottaja",
					Address: "eve.vastaanottaja@example.net",
				},
			},
			ResentBcc: []*mail.Address{
				{
					Name:    "Frank Vastaanottaja",
					Address: "frank.vastaanottaja@example.net",
				},
				{
					Name:    "Grace Vastaanottaja",
					Address: "grace.vastaanottaja@example.net",
				},
			},
			ResentMessageID: "Message-Id-1@example.net",
			ContentType: ContentTypeHeader{
				ContentType: "multipart/related",
				Params: map[string]string{
					"boundary": "RelatedBoundaryString",
					"charset":  "iso-8859-15",
				},
			},
			ExtraHeaders: map[string][]string{
				"X-Clacks-Overhead": {"GNU Terry Pratchett"},
				"X-Script/function/\t !\"#$%&'()*+,-./;<=>?@[\\]^_`{|}~": {
					"TEST VALUE 1\t !\"#$%&'()*+,-./0123456789:;<=>?@ABCDEFGHIJKLMNOPQRSTUVWXYZ[\\]^_` abcdefghijklmnopqrstuvwxyz{|}~",
					"TEST VALUE 2\t !\"#$%&'()*+,-./0123456789:;<=>?@ABCDEFGHIJKLMNOPQRSTUVWXYZ[\\]^_` abcdefghijklmnopqrstuvwxyz{|}~",
				},
			},
		},
		Text: `Albert osti fagotin ja töräytti puhkuvan melodian.
Lorun sangen pieneksi hyödyksi jäivät suomen kirjaimet.
Hyvän lorun sangen pieneksi hyödyksi jäi suomen kirjaimet.
Fahrenheit ja Celsius yrjösivät Åsan backgammon-peliin, Volkswagenissa, daiquirin ja ZX81:n yhteisvaikutuksesta.
Charles Darwin jammaili Åken hevixylofonilla Qatarin yöpub Zeligissä.
Wieniläinen sioux:ta puhuva ökyzombie diggaa Åsan roquefort-tacoja.`,
		EnrichedText: `<bold>Albert osti fagotin ja töräytti puhkuvan melodian.</bold>
<italic>Lorun sangen pieneksi hyödyksi jäivät suomen kirjaimet.</italic>
<fixed>Hyvän lorun sangen pieneksi hyödyksi jäi suomen kirjaimet.</fixed>
<underline>Fahrenheit ja Celsius yrjösivät Åsan backgammon-peliin, Volkswagenissa, daiquirin ja ZX81:n yhteisvaikutuksesta.</underline>
Charles Darwin jammaili Åken hevixylofonilla Qatarin yöpub Zeligissä.
Wieniläinen sioux:ta puhuva ökyzombie diggaa Åsan roquefort-tacoja.`,
		HTML: `<html>
<div dir="ltr">
<p>Albert osti fagotin ja töräytti puhkuvan melodian.</p>
<p>Lorun sangen pieneksi hyödyksi jäivät suomen kirjaimet.</p>
<p>Hyvän lorun sangen pieneksi hyödyksi jäi suomen kirjaimet.</p>
<p>Fahrenheit ja Celsius yrjösivät Åsan backgammon-peliin, Volkswagenissa, daiquirin ja ZX81:n yhteisvaikutuksesta.</p>
<p>Charles Darwin jammaili Åken hevixylofonilla Qatarin yöpub Zeligissä.</p>
<p>Wieniläinen sioux:ta puhuva ökyzombie diggaa Åsan roquefort-tacoja.</p>
</div>
</html>`,
		InlineFiles: []InlineFile{
			{
				ContentID: "inline-jpg-image.jpg@example.com",
				ContentType: ContentTypeHeader{
					ContentType: "image/jpeg",
					Params: map[string]string{
						"name": "inline-jpg-image-name.jpg",
					},
				},
				ContentDisposition: ContentDispositionHeader{
					ContentDisposition: inline,
					Params: map[string]string{
						"filename": "inline-jpg-image-filename.jpg",
					},
				},
				Data: []byte{255, 216, 255, 219, 0, 67, 0, 3, 2, 2, 2, 2, 2, 3, 2, 2, 2, 3, 3, 3, 3, 4,
					6, 4, 4, 4, 4, 4, 8, 6, 6, 5, 6, 9, 8, 10, 10, 9, 8, 9, 9, 10, 12, 15, 12, 10, 11, 14, 11, 9, 9,
					13, 17, 13, 14, 15, 16, 16, 17, 16, 10, 12, 18, 19, 18, 16, 19, 15, 16, 16, 16, 255, 201, 0, 11, 8,
					0, 1, 0, 1, 1, 1, 17, 0, 255, 204, 0, 6, 0, 16, 16, 5, 255, 218, 0, 8, 1, 1, 0, 0, 63, 0, 210, 207,
					32, 255, 217},
			},
		},
	}

	testEmailFromFile(t, fp, expectedEmail)
}
func TestParseEmailFinnishMultipartMixedUtf8OverBase64(t *testing.T) {
	fp := "tests/test_finnish_multipart_mixed_utf-8_over_base64.txt"
	tz, _ := time.LoadLocation("Europe/Helsinki")
	expectedDate, _ := time.Parse(
		time.RFC1123Z+" (MST)",
		time.Date(2019, time.April, 1, 7, 55, 0, 0, tz).Format(time.RFC1123Z+" (MST)"))
	expectedEmail := Email{
		Headers: Headers{
			Date:    expectedDate,
			Subject: "📧 Test Suomenkieliset pangrammit",
			ReplyTo: []*mail.Address{
				{
					Name:    "Alice Lähettäjä",
					Address: "alice.lahettaja@example.net",
				},
			},
			Sender: &mail.Address{
				Name:    "Alice Lähettäjä",
				Address: "alice.lahettaja@example.com",
			},
			From: []*mail.Address{
				{
					Name:    "Alice Lähettäjä",
					Address: "alice.lahettaja@example.com",
				},
				{
					Name:    "Alice Lähettäjä",
					Address: "alice.lahettaja@example.net",
				},
			},
			To: []*mail.Address{
				{
					Name:    "Bob Vastaanottaja",
					Address: "bob.vastaanottaja@exaple.com",
				},
				{
					Name:    "Carol Vastaanottaja",
					Address: "carol.vastaanottaja@example.com",
				},
			},
			Cc: []*mail.Address{
				{
					Name:    "Dan Vastaanottaja",
					Address: "dan.vastaanottaja@example.com",
				},
				{
					Name:    "Eve Vastaanottaja",
					Address: "eve.vastaanottaja@example.com",
				},
			},
			Bcc: []*mail.Address{
				{
					Name:    "Frank Vastaanottaja",
					Address: "frank.vastaanottaja@example.com",
				},
				{
					Name:    "Grace Vastaanottaja",
					Address: "grace.vastaanottaja@example.com",
				},
			},
			MessageID:  "Message-Id-1@example.com",
			InReplyTo:  []MessageId{"Message-Id-0@example.com"},
			References: []MessageId{"Message-Id-0@example.com"},
			Comments:   "Message Header Comment",
			Keywords:   []string{"Keyword 1", "Keyword 2"},
			ResentDate: expectedDate,
			ResentFrom: []*mail.Address{
				{
					Name:    "Alice Lähettäjä",
					Address: "alice.lahettaja@example.net",
				},
				{
					Name:    "Alice Lähettäjä",
					Address: "alice.lahettaja@example.com",
				},
			},
			ResentSender: &mail.Address{
				Name:    "Alice Lähettäjä",
				Address: "alice.lahettaja@example.net",
			},
			ResentTo: []*mail.Address{
				{
					Name:    "Bob Vastaanottaja",
					Address: "bob.vastaanottaja@exaple.com",
				},
				{
					Name:    "Carol Vastaanottaja",
					Address: "carol.vastaanottaja@example.net",
				},
			},
			ResentCc: []*mail.Address{
				{
					Name:    "Dan Vastaanottaja",
					Address: "dan.vastaanottaja@example.net",
				},
				{
					Name:    "Eve Vastaanottaja",
					Address: "eve.vastaanottaja@example.net",
				},
			},
			ResentBcc: []*mail.Address{
				{
					Name:    "Frank Vastaanottaja",
					Address: "frank.vastaanottaja@example.net",
				},
				{
					Name:    "Grace Vastaanottaja",
					Address: "grace.vastaanottaja@example.net",
				},
			},
			ResentMessageID: "Message-Id-1@example.net",
			ContentType: ContentTypeHeader{
				ContentType: "multipart/mixed",
				Params: map[string]string{
					"boundary": "MixedBoundaryString",
					"charset":  "utf-8",
				},
			},
			ExtraHeaders: map[string][]string{
				"X-Clacks-Overhead": {"GNU Terry Pratchett"},
				"X-Script/function/\t !\"#$%&'()*+,-./;<=>?@[\\]^_`{|}~": {
					"TEST VALUE 1\t !\"#$%&'()*+,-./0123456789:;<=>?@ABCDEFGHIJKLMNOPQRSTUVWXYZ[\\]^_` abcdefghijklmnopqrstuvwxyz{|}~",
					"TEST VALUE 2\t !\"#$%&'()*+,-./0123456789:;<=>?@ABCDEFGHIJKLMNOPQRSTUVWXYZ[\\]^_` abcdefghijklmnopqrstuvwxyz{|}~",
				},
			},
		},
		Text: `Albert osti fagotin ja töräytti puhkuvan melodian.
Lorun sangen pieneksi hyödyksi jäivät suomen kirjaimet.
Hyvän lorun sangen pieneksi hyödyksi jäi suomen kirjaimet.
Fahrenheit ja Celsius yrjösivät Åsan backgammon-peliin, Volkswagenissa, daiquirin ja ZX81:n yhteisvaikutuksesta.
Charles Darwin jammaili Åken hevixylofonilla Qatarin yöpub Zeligissä.
Wieniläinen sioux:ta puhuva ökyzombie diggaa Åsan roquefort-tacoja.`,
		EnrichedText: `<bold>Albert osti fagotin ja töräytti puhkuvan melodian.</bold>
<italic>Lorun sangen pieneksi hyödyksi jäivät suomen kirjaimet.</italic>
<fixed>Hyvän lorun sangen pieneksi hyödyksi jäi suomen kirjaimet.</fixed>
<underline>Fahrenheit ja Celsius yrjösivät Åsan backgammon-peliin, Volkswagenissa, daiquirin ja ZX81:n yhteisvaikutuksesta.</underline>
Charles Darwin jammaili Åken hevixylofonilla Qatarin yöpub Zeligissä.
Wieniläinen sioux:ta puhuva ökyzombie diggaa Åsan roquefort-tacoja.`,
		HTML: `<html>
<div dir="ltr">
<p>Albert osti fagotin ja töräytti puhkuvan melodian.</p>
<p>Lorun sangen pieneksi hyödyksi jäivät suomen kirjaimet.</p>
<p>Hyvän lorun sangen pieneksi hyödyksi jäi suomen kirjaimet.</p>
<p>Fahrenheit ja Celsius yrjösivät Åsan backgammon-peliin, Volkswagenissa, daiquirin ja ZX81:n yhteisvaikutuksesta.</p>
<p>Charles Darwin jammaili Åken hevixylofonilla Qatarin yöpub Zeligissä.</p>
<p>Wieniläinen sioux:ta puhuva ökyzombie diggaa Åsan roquefort-tacoja.</p>
</div>
</html>`,
		InlineFiles: []InlineFile{
			{
				ContentType: ContentTypeHeader{
					ContentType: "image/jpeg",
					Params: map[string]string{
						"name": "inline-jpg-image-without-disposition.jpg",
					},
				},
				ContentDisposition: ContentDispositionHeader{
					ContentDisposition: "",
					Params:             map[string]string(nil),
				},
				Data: []byte{255, 216, 255, 219, 0, 67, 0, 3, 2, 2, 2, 2, 2, 3, 2, 2, 2, 3, 3, 3, 3, 4,
					6, 4, 4, 4, 4, 4, 8, 6, 6, 5, 6, 9, 8, 10, 10, 9, 8, 9, 9, 10, 12, 15, 12, 10, 11, 14, 11, 9, 9,
					13, 17, 13, 14, 15, 16, 16, 17, 16, 10, 12, 18, 19, 18, 16, 19, 15, 16, 16, 16, 255, 201, 0, 11, 8,
					0, 1, 0, 1, 1, 1, 17, 0, 255, 204, 0, 6, 0, 16, 16, 5, 255, 218, 0, 8, 1, 1, 0, 0, 63, 0, 210, 207,
					32, 255, 217},
			},
			{
				ContentID: "inline-jpg-image.jpg@example.com",
				ContentType: ContentTypeHeader{
					ContentType: "image/jpeg",
					Params: map[string]string{
						"name": "inline-jpg-image-name.jpg",
					},
				},
				ContentDisposition: ContentDispositionHeader{
					ContentDisposition: inline,
					Params: map[string]string{
						"filename": "inline-jpg-image-filename.jpg",
					},
				},
				Data: []byte{255, 216, 255, 219, 0, 67, 0, 3, 2, 2, 2, 2, 2, 3, 2, 2, 2, 3, 3, 3, 3, 4,
					6, 4, 4, 4, 4, 4, 8, 6, 6, 5, 6, 9, 8, 10, 10, 9, 8, 9, 9, 10, 12, 15, 12, 10, 11, 14, 11, 9, 9,
					13, 17, 13, 14, 15, 16, 16, 17, 16, 10, 12, 18, 19, 18, 16, 19, 15, 16, 16, 16, 255, 201, 0, 11, 8,
					0, 1, 0, 1, 1, 1, 17, 0, 255, 204, 0, 6, 0, 16, 16, 5, 255, 218, 0, 8, 1, 1, 0, 0, 63, 0, 210, 207,
					32, 255, 217},
			},
		},
		AttachedFiles: []AttachedFile{
			{
				ContentType: ContentTypeHeader{
					ContentType: "application/pdf",
					Params: map[string]string{
						"name": "attached-pdf-name.pdf",
					},
				},
				ContentDisposition: ContentDispositionHeader{
					ContentDisposition: attachment,
					Params: map[string]string{
						"filename": "attached-pdf-filename.pdf",
					},
				},
				Data: []byte{37, 80, 68, 70, 45, 49, 46, 13, 116, 114, 97, 105, 108, 101, 114, 60, 60,
					47, 82, 111, 111, 116, 60, 60, 47, 80, 97, 103, 101, 115, 60, 60, 47, 75, 105, 100, 115, 91, 60,
					60, 47, 77, 101, 100, 105, 97, 66, 111, 120, 91, 48, 32, 48, 32, 51, 32, 51, 93, 62, 62, 93, 62,
					62, 62, 62, 62, 62},
			},
			{
				ContentType: ContentTypeHeader{
					ContentType: "application/pdf",
					Params: map[string]string{
						"name": "attached-pdf-without-disposition.pdf",
					},
				},
				ContentDisposition: ContentDispositionHeader{
					ContentDisposition: "",
					Params:             map[string]string(nil),
				},
				Data: []byte{37, 80, 68, 70, 45, 49, 46, 13, 116, 114, 97, 105, 108, 101, 114, 60, 60,
					47, 82, 111, 111, 116, 60, 60, 47, 80, 97, 103, 101, 115, 60, 60, 47, 75, 105, 100, 115, 91, 60,
					60, 47, 77, 101, 100, 105, 97, 66, 111, 120, 91, 48, 32, 48, 32, 51, 32, 51, 93, 62, 62, 93, 62,
					62, 62, 62, 62, 62},
			},
			{
				ContentType: ContentTypeHeader{
					ContentType: "application/json",
					Params: map[string]string{
						"name": "attached-json-name.json",
					},
				},
				ContentDisposition: ContentDispositionHeader{
					ContentDisposition: attachment,
					Params: map[string]string{
						"filename": "attached-json-filename.json",
					},
				},
				Data: []byte{123, 34, 102, 111, 111, 34, 58, 34, 98, 97, 114, 34, 125},
			},
			{
				ContentType: ContentTypeHeader{
					ContentType: "text/plain",
					Params: map[string]string{
						"name": "attached-text-plain-name.txt",
					},
				},
				ContentDisposition: ContentDispositionHeader{
					ContentDisposition: attachment,
					Params: map[string]string{
						"filename": "attached-text-plain-filename.txt",
					},
				},
				Data: []byte{84, 101, 120, 116, 47, 112, 108, 97, 105, 110, 32, 99, 111, 110, 116, 101, 110, 116, 32,
					97, 115, 32, 97, 110, 32, 97, 116, 116, 97, 99, 104, 101, 100, 32, 46, 116, 120, 116, 32, 102, 105,
					108, 101, 46},
			},
			{
				ContentType: ContentTypeHeader{
					ContentType: "text/html",
					Params: map[string]string{
						"name": "attached-text-html-name.html",
					},
				},
				ContentDisposition: ContentDispositionHeader{
					ContentDisposition: attachment,
					Params: map[string]string{
						"filename": "attached-text-html-filename.html",
					},
				},
				Data: []byte{84, 101, 120, 116, 47, 104, 116, 109, 108, 32, 99, 111, 110, 116, 101, 110, 116, 32, 97,
					115, 32, 97, 110, 32, 97, 116, 116, 97, 99, 104, 101, 100, 32, 46, 104, 116, 109, 108, 32, 102,
					105, 108, 101, 46},
			},
		},
	}

	testEmailFromFile(t, fp, expectedEmail)
}
func TestParseEmailFinnishMultipartMixedUtf8OverQuotedprintable(t *testing.T) {
	fp := "tests/test_finnish_multipart_mixed_utf-8_over_quoted-printable.txt"
	tz, _ := time.LoadLocation("Europe/Helsinki")
	expectedDate, _ := time.Parse(
		time.RFC1123Z+" (MST)",
		time.Date(2019, time.April, 1, 7, 55, 0, 0, tz).Format(time.RFC1123Z+" (MST)"))
	expectedEmail := Email{
		Headers: Headers{
			Date:    expectedDate,
			Subject: "📧 Test Suomenkieliset pangrammit",
			ReplyTo: []*mail.Address{
				{
					Name:    "Alice Lähettäjä",
					Address: "alice.lahettaja@example.net",
				},
			},
			Sender: &mail.Address{
				Name:    "Alice Lähettäjä",
				Address: "alice.lahettaja@example.com",
			},
			From: []*mail.Address{
				{
					Name:    "Alice Lähettäjä",
					Address: "alice.lahettaja@example.com",
				},
				{
					Name:    "Alice Lähettäjä",
					Address: "alice.lahettaja@example.net",
				},
			},
			To: []*mail.Address{
				{
					Name:    "Bob Vastaanottaja",
					Address: "bob.vastaanottaja@exaple.com",
				},
				{
					Name:    "Carol Vastaanottaja",
					Address: "carol.vastaanottaja@example.com",
				},
			},
			Cc: []*mail.Address{
				{
					Name:    "Dan Vastaanottaja",
					Address: "dan.vastaanottaja@example.com",
				},
				{
					Name:    "Eve Vastaanottaja",
					Address: "eve.vastaanottaja@example.com",
				},
			},
			Bcc: []*mail.Address{
				{
					Name:    "Frank Vastaanottaja",
					Address: "frank.vastaanottaja@example.com",
				},
				{
					Name:    "Grace Vastaanottaja",
					Address: "grace.vastaanottaja@example.com",
				},
			},
			MessageID:  "Message-Id-1@example.com",
			InReplyTo:  []MessageId{"Message-Id-0@example.com"},
			References: []MessageId{"Message-Id-0@example.com"},
			Comments:   "Message Header Comment",
			Keywords:   []string{"Keyword 1", "Keyword 2"},
			ResentDate: expectedDate,
			ResentFrom: []*mail.Address{
				{
					Name:    "Alice Lähettäjä",
					Address: "alice.lahettaja@example.net",
				},
				{
					Name:    "Alice Lähettäjä",
					Address: "alice.lahettaja@example.com",
				},
			},
			ResentSender: &mail.Address{
				Name:    "Alice Lähettäjä",
				Address: "alice.lahettaja@example.net",
			},
			ResentTo: []*mail.Address{
				{
					Name:    "Bob Vastaanottaja",
					Address: "bob.vastaanottaja@exaple.com",
				},
				{
					Name:    "Carol Vastaanottaja",
					Address: "carol.vastaanottaja@example.net",
				},
			},
			ResentCc: []*mail.Address{
				{
					Name:    "Dan Vastaanottaja",
					Address: "dan.vastaanottaja@example.net",
				},
				{
					Name:    "Eve Vastaanottaja",
					Address: "eve.vastaanottaja@example.net",
				},
			},
			ResentBcc: []*mail.Address{
				{
					Name:    "Frank Vastaanottaja",
					Address: "frank.vastaanottaja@example.net",
				},
				{
					Name:    "Grace Vastaanottaja",
					Address: "grace.vastaanottaja@example.net",
				},
			},
			ResentMessageID: "Message-Id-1@example.net",
			ContentType: ContentTypeHeader{
				ContentType: "multipart/mixed",
				Params: map[string]string{
					"boundary": "MixedBoundaryString",
					"charset":  "utf-8",
				},
			},
			ExtraHeaders: map[string][]string{
				"X-Clacks-Overhead": {"GNU Terry Pratchett"},
				"X-Script/function/\t !\"#$%&'()*+,-./;<=>?@[\\]^_`{|}~": {
					"TEST VALUE 1\t !\"#$%&'()*+,-./0123456789:;<=>?@ABCDEFGHIJKLMNOPQRSTUVWXYZ[\\]^_` abcdefghijklmnopqrstuvwxyz{|}~",
					"TEST VALUE 2\t !\"#$%&'()*+,-./0123456789:;<=>?@ABCDEFGHIJKLMNOPQRSTUVWXYZ[\\]^_` abcdefghijklmnopqrstuvwxyz{|}~",
				},
			},
		},
		Text: `Albert osti fagotin ja töräytti puhkuvan melodian.
Lorun sangen pieneksi hyödyksi jäivät suomen kirjaimet.
Hyvän lorun sangen pieneksi hyödyksi jäi suomen kirjaimet.
Fahrenheit ja Celsius yrjösivät Åsan backgammon-peliin, Volkswagenissa, daiquirin ja ZX81:n yhteisvaikutuksesta.
Charles Darwin jammaili Åken hevixylofonilla Qatarin yöpub Zeligissä.
Wieniläinen sioux:ta puhuva ökyzombie diggaa Åsan roquefort-tacoja.`,
		EnrichedText: `<bold>Albert osti fagotin ja töräytti puhkuvan melodian.</bold>
<italic>Lorun sangen pieneksi hyödyksi jäivät suomen kirjaimet.</italic>
<fixed>Hyvän lorun sangen pieneksi hyödyksi jäi suomen kirjaimet.</fixed>
<underline>Fahrenheit ja Celsius yrjösivät Åsan backgammon-peliin, Volkswagenissa, daiquirin ja ZX81:n yhteisvaikutuksesta.</underline>
Charles Darwin jammaili Åken hevixylofonilla Qatarin yöpub Zeligissä.
Wieniläinen sioux:ta puhuva ökyzombie diggaa Åsan roquefort-tacoja.`,
		HTML: `<html>
<div dir="ltr">
<p>Albert osti fagotin ja töräytti puhkuvan melodian.</p>
<p>Lorun sangen pieneksi hyödyksi jäivät suomen kirjaimet.</p>
<p>Hyvän lorun sangen pieneksi hyödyksi jäi suomen kirjaimet.</p>
<p>Fahrenheit ja Celsius yrjösivät Åsan backgammon-peliin, Volkswagenissa, daiquirin ja ZX81:n yhteisvaikutuksesta.</p>
<p>Charles Darwin jammaili Åken hevixylofonilla Qatarin yöpub Zeligissä.</p>
<p>Wieniläinen sioux:ta puhuva ökyzombie diggaa Åsan roquefort-tacoja.</p>
</div>
</html>`,
		InlineFiles: []InlineFile{
			{
				ContentType: ContentTypeHeader{
					ContentType: "image/jpeg",
					Params: map[string]string{
						"name": "inline-jpg-image-without-disposition.jpg",
					},
				},
				ContentDisposition: ContentDispositionHeader{
					ContentDisposition: "",
					Params:             map[string]string(nil),
				},
				Data: []byte{255, 216, 255, 219, 0, 67, 0, 3, 2, 2, 2, 2, 2, 3, 2, 2, 2, 3, 3, 3, 3, 4,
					6, 4, 4, 4, 4, 4, 8, 6, 6, 5, 6, 9, 8, 10, 10, 9, 8, 9, 9, 10, 12, 15, 12, 10, 11, 14, 11, 9, 9,
					13, 17, 13, 14, 15, 16, 16, 17, 16, 10, 12, 18, 19, 18, 16, 19, 15, 16, 16, 16, 255, 201, 0, 11, 8,
					0, 1, 0, 1, 1, 1, 17, 0, 255, 204, 0, 6, 0, 16, 16, 5, 255, 218, 0, 8, 1, 1, 0, 0, 63, 0, 210, 207,
					32, 255, 217},
			},
			{
				ContentID: "inline-jpg-image.jpg@example.com",
				ContentType: ContentTypeHeader{
					ContentType: "image/jpeg",
					Params: map[string]string{
						"name": "inline-jpg-image-name.jpg",
					},
				},
				ContentDisposition: ContentDispositionHeader{
					ContentDisposition: inline,
					Params: map[string]string{
						"filename": "inline-jpg-image-filename.jpg",
					},
				},
				Data: []byte{255, 216, 255, 219, 0, 67, 0, 3, 2, 2, 2, 2, 2, 3, 2, 2, 2, 3, 3, 3, 3, 4,
					6, 4, 4, 4, 4, 4, 8, 6, 6, 5, 6, 9, 8, 10, 10, 9, 8, 9, 9, 10, 12, 15, 12, 10, 11, 14, 11, 9, 9,
					13, 17, 13, 14, 15, 16, 16, 17, 16, 10, 12, 18, 19, 18, 16, 19, 15, 16, 16, 16, 255, 201, 0, 11, 8,
					0, 1, 0, 1, 1, 1, 17, 0, 255, 204, 0, 6, 0, 16, 16, 5, 255, 218, 0, 8, 1, 1, 0, 0, 63, 0, 210, 207,
					32, 255, 217},
			},
		},
		AttachedFiles: []AttachedFile{
			{
				ContentType: ContentTypeHeader{
					ContentType: "application/pdf",
					Params: map[string]string{
						"name": "attached-pdf-name.pdf",
					},
				},
				ContentDisposition: ContentDispositionHeader{
					ContentDisposition: attachment,
					Params: map[string]string{
						"filename": "attached-pdf-filename.pdf",
					},
				},
				Data: []byte{37, 80, 68, 70, 45, 49, 46, 13, 116, 114, 97, 105, 108, 101, 114, 60, 60,
					47, 82, 111, 111, 116, 60, 60, 47, 80, 97, 103, 101, 115, 60, 60, 47, 75, 105, 100, 115, 91, 60,
					60, 47, 77, 101, 100, 105, 97, 66, 111, 120, 91, 48, 32, 48, 32, 51, 32, 51, 93, 62, 62, 93, 62,
					62, 62, 62, 62, 62},
			},
			{
				ContentType: ContentTypeHeader{
					ContentType: "application/pdf",
					Params: map[string]string{
						"name": "attached-pdf-without-disposition.pdf",
					},
				},
				ContentDisposition: ContentDispositionHeader{
					ContentDisposition: "",
					Params:             map[string]string(nil),
				},
				Data: []byte{37, 80, 68, 70, 45, 49, 46, 13, 116, 114, 97, 105, 108, 101, 114, 60, 60,
					47, 82, 111, 111, 116, 60, 60, 47, 80, 97, 103, 101, 115, 60, 60, 47, 75, 105, 100, 115, 91, 60,
					60, 47, 77, 101, 100, 105, 97, 66, 111, 120, 91, 48, 32, 48, 32, 51, 32, 51, 93, 62, 62, 93, 62,
					62, 62, 62, 62, 62},
			},
			{
				ContentType: ContentTypeHeader{
					ContentType: "application/json",
					Params: map[string]string{
						"name": "attached-json-name.json",
					},
				},
				ContentDisposition: ContentDispositionHeader{
					ContentDisposition: attachment,
					Params: map[string]string{
						"filename": "attached-json-filename.json",
					},
				},
				Data: []byte{123, 34, 102, 111, 111, 34, 58, 34, 98, 97, 114, 34, 125},
			},
			{
				ContentType: ContentTypeHeader{
					ContentType: "text/plain",
					Params: map[string]string{
						"name": "attached-text-plain-name.txt",
					},
				},
				ContentDisposition: ContentDispositionHeader{
					ContentDisposition: attachment,
					Params: map[string]string{
						"filename": "attached-text-plain-filename.txt",
					},
				},
				Data: []byte{84, 101, 120, 116, 47, 112, 108, 97, 105, 110, 32, 99, 111, 110, 116, 101, 110, 116, 32,
					97, 115, 32, 97, 110, 32, 97, 116, 116, 97, 99, 104, 101, 100, 32, 46, 116, 120, 116, 32, 102, 105,
					108, 101, 46},
			},
			{
				ContentType: ContentTypeHeader{
					ContentType: "text/html",
					Params: map[string]string{
						"name": "attached-text-html-name.html",
					},
				},
				ContentDisposition: ContentDispositionHeader{
					ContentDisposition: attachment,
					Params: map[string]string{
						"filename": "attached-text-html-filename.html",
					},
				},
				Data: []byte{84, 101, 120, 116, 47, 104, 116, 109, 108, 32, 99, 111, 110, 116, 101, 110, 116, 32, 97,
					115, 32, 97, 110, 32, 97, 116, 116, 97, 99, 104, 101, 100, 32, 46, 104, 116, 109, 108, 32, 102,
					105, 108, 101, 46},
			},
		},
	}

	testEmailFromFile(t, fp, expectedEmail)
}
func TestParseEmailFinnishMultipartMixedIso885915OverBase64(t *testing.T) {
	fp := "tests/test_finnish_multipart_mixed_iso-8859-15_over_base64.txt"
	tz, _ := time.LoadLocation("Europe/Helsinki")
	expectedDate, _ := time.Parse(
		time.RFC1123Z+" (MST)",
		time.Date(2019, time.April, 1, 7, 55, 0, 0, tz).Format(time.RFC1123Z+" (MST)"))
	expectedEmail := Email{
		Headers: Headers{
			Date:    expectedDate,
			Subject: "📧 Test Suomenkieliset pangrammit",
			ReplyTo: []*mail.Address{
				{
					Name:    "Alice Lähettäjä",
					Address: "alice.lahettaja@example.net",
				},
			},
			Sender: &mail.Address{
				Name:    "Alice Lähettäjä",
				Address: "alice.lahettaja@example.com",
			},
			From: []*mail.Address{
				{
					Name:    "Alice Lähettäjä",
					Address: "alice.lahettaja@example.com",
				},
				{
					Name:    "Alice Lähettäjä",
					Address: "alice.lahettaja@example.net",
				},
			},
			To: []*mail.Address{
				{
					Name:    "Bob Vastaanottaja",
					Address: "bob.vastaanottaja@exaple.com",
				},
				{
					Name:    "Carol Vastaanottaja",
					Address: "carol.vastaanottaja@example.com",
				},
			},
			Cc: []*mail.Address{
				{
					Name:    "Dan Vastaanottaja",
					Address: "dan.vastaanottaja@example.com",
				},
				{
					Name:    "Eve Vastaanottaja",
					Address: "eve.vastaanottaja@example.com",
				},
			},
			Bcc: []*mail.Address{
				{
					Name:    "Frank Vastaanottaja",
					Address: "frank.vastaanottaja@example.com",
				},
				{
					Name:    "Grace Vastaanottaja",
					Address: "grace.vastaanottaja@example.com",
				},
			},
			MessageID:  "Message-Id-1@example.com",
			InReplyTo:  []MessageId{"Message-Id-0@example.com"},
			References: []MessageId{"Message-Id-0@example.com"},
			Comments:   "Message Header Comment",
			Keywords:   []string{"Keyword 1", "Keyword 2"},
			ResentDate: expectedDate,
			ResentFrom: []*mail.Address{
				{
					Name:    "Alice Lähettäjä",
					Address: "alice.lahettaja@example.net",
				},
				{
					Name:    "Alice Lähettäjä",
					Address: "alice.lahettaja@example.com",
				},
			},
			ResentSender: &mail.Address{
				Name:    "Alice Lähettäjä",
				Address: "alice.lahettaja@example.net",
			},
			ResentTo: []*mail.Address{
				{
					Name:    "Bob Vastaanottaja",
					Address: "bob.vastaanottaja@exaple.com",
				},
				{
					Name:    "Carol Vastaanottaja",
					Address: "carol.vastaanottaja@example.net",
				},
			},
			ResentCc: []*mail.Address{
				{
					Name:    "Dan Vastaanottaja",
					Address: "dan.vastaanottaja@example.net",
				},
				{
					Name:    "Eve Vastaanottaja",
					Address: "eve.vastaanottaja@example.net",
				},
			},
			ResentBcc: []*mail.Address{
				{
					Name:    "Frank Vastaanottaja",
					Address: "frank.vastaanottaja@example.net",
				},
				{
					Name:    "Grace Vastaanottaja",
					Address: "grace.vastaanottaja@example.net",
				},
			},
			ResentMessageID: "Message-Id-1@example.net",
			ContentType: ContentTypeHeader{
				ContentType: "multipart/mixed",
				Params: map[string]string{
					"boundary": "MixedBoundaryString",
					"charset":  "iso-8859-15",
				},
			},
			ExtraHeaders: map[string][]string{
				"X-Clacks-Overhead": {"GNU Terry Pratchett"},
				"X-Script/function/\t !\"#$%&'()*+,-./;<=>?@[\\]^_`{|}~": {
					"TEST VALUE 1\t !\"#$%&'()*+,-./0123456789:;<=>?@ABCDEFGHIJKLMNOPQRSTUVWXYZ[\\]^_` abcdefghijklmnopqrstuvwxyz{|}~",
					"TEST VALUE 2\t !\"#$%&'()*+,-./0123456789:;<=>?@ABCDEFGHIJKLMNOPQRSTUVWXYZ[\\]^_` abcdefghijklmnopqrstuvwxyz{|}~",
				},
			},
		},
		Text: `Albert osti fagotin ja töräytti puhkuvan melodian.
Lorun sangen pieneksi hyödyksi jäivät suomen kirjaimet.
Hyvän lorun sangen pieneksi hyödyksi jäi suomen kirjaimet.
Fahrenheit ja Celsius yrjösivät Åsan backgammon-peliin, Volkswagenissa, daiquirin ja ZX81:n yhteisvaikutuksesta.
Charles Darwin jammaili Åken hevixylofonilla Qatarin yöpub Zeligissä.
Wieniläinen sioux:ta puhuva ökyzombie diggaa Åsan roquefort-tacoja.`,
		EnrichedText: `<bold>Albert osti fagotin ja töräytti puhkuvan melodian.</bold>
<italic>Lorun sangen pieneksi hyödyksi jäivät suomen kirjaimet.</italic>
<fixed>Hyvän lorun sangen pieneksi hyödyksi jäi suomen kirjaimet.</fixed>
<underline>Fahrenheit ja Celsius yrjösivät Åsan backgammon-peliin, Volkswagenissa, daiquirin ja ZX81:n yhteisvaikutuksesta.</underline>
Charles Darwin jammaili Åken hevixylofonilla Qatarin yöpub Zeligissä.
Wieniläinen sioux:ta puhuva ökyzombie diggaa Åsan roquefort-tacoja.`,
		HTML: `<html>
<div dir="ltr">
<p>Albert osti fagotin ja töräytti puhkuvan melodian.</p>
<p>Lorun sangen pieneksi hyödyksi jäivät suomen kirjaimet.</p>
<p>Hyvän lorun sangen pieneksi hyödyksi jäi suomen kirjaimet.</p>
<p>Fahrenheit ja Celsius yrjösivät Åsan backgammon-peliin, Volkswagenissa, daiquirin ja ZX81:n yhteisvaikutuksesta.</p>
<p>Charles Darwin jammaili Åken hevixylofonilla Qatarin yöpub Zeligissä.</p>
<p>Wieniläinen sioux:ta puhuva ökyzombie diggaa Åsan roquefort-tacoja.</p>
</div>
</html>`,
		InlineFiles: []InlineFile{
			{
				ContentType: ContentTypeHeader{
					ContentType: "image/jpeg",
					Params: map[string]string{
						"name": "inline-jpg-image-without-disposition.jpg",
					},
				},
				ContentDisposition: ContentDispositionHeader{
					ContentDisposition: "",
					Params:             map[string]string(nil),
				},
				Data: []byte{255, 216, 255, 219, 0, 67, 0, 3, 2, 2, 2, 2, 2, 3, 2, 2, 2, 3, 3, 3, 3, 4,
					6, 4, 4, 4, 4, 4, 8, 6, 6, 5, 6, 9, 8, 10, 10, 9, 8, 9, 9, 10, 12, 15, 12, 10, 11, 14, 11, 9, 9,
					13, 17, 13, 14, 15, 16, 16, 17, 16, 10, 12, 18, 19, 18, 16, 19, 15, 16, 16, 16, 255, 201, 0, 11, 8,
					0, 1, 0, 1, 1, 1, 17, 0, 255, 204, 0, 6, 0, 16, 16, 5, 255, 218, 0, 8, 1, 1, 0, 0, 63, 0, 210, 207,
					32, 255, 217},
			},
			{
				ContentID: "inline-jpg-image.jpg@example.com",
				ContentType: ContentTypeHeader{
					ContentType: "image/jpeg",
					Params: map[string]string{
						"name": "inline-jpg-image-name.jpg",
					},
				},
				ContentDisposition: ContentDispositionHeader{
					ContentDisposition: inline,
					Params: map[string]string{
						"filename": "inline-jpg-image-filename.jpg",
					},
				},
				Data: []byte{255, 216, 255, 219, 0, 67, 0, 3, 2, 2, 2, 2, 2, 3, 2, 2, 2, 3, 3, 3, 3, 4,
					6, 4, 4, 4, 4, 4, 8, 6, 6, 5, 6, 9, 8, 10, 10, 9, 8, 9, 9, 10, 12, 15, 12, 10, 11, 14, 11, 9, 9,
					13, 17, 13, 14, 15, 16, 16, 17, 16, 10, 12, 18, 19, 18, 16, 19, 15, 16, 16, 16, 255, 201, 0, 11, 8,
					0, 1, 0, 1, 1, 1, 17, 0, 255, 204, 0, 6, 0, 16, 16, 5, 255, 218, 0, 8, 1, 1, 0, 0, 63, 0, 210, 207,
					32, 255, 217},
			},
		},
		AttachedFiles: []AttachedFile{
			{
				ContentType: ContentTypeHeader{
					ContentType: "application/pdf",
					Params: map[string]string{
						"name": "attached-pdf-name.pdf",
					},
				},
				ContentDisposition: ContentDispositionHeader{
					ContentDisposition: attachment,
					Params: map[string]string{
						"filename": "attached-pdf-filename.pdf",
					},
				},
				Data: []byte{37, 80, 68, 70, 45, 49, 46, 13, 116, 114, 97, 105, 108, 101, 114, 60, 60,
					47, 82, 111, 111, 116, 60, 60, 47, 80, 97, 103, 101, 115, 60, 60, 47, 75, 105, 100, 115, 91, 60,
					60, 47, 77, 101, 100, 105, 97, 66, 111, 120, 91, 48, 32, 48, 32, 51, 32, 51, 93, 62, 62, 93, 62,
					62, 62, 62, 62, 62},
			},
			{
				ContentType: ContentTypeHeader{
					ContentType: "application/pdf",
					Params: map[string]string{
						"name": "attached-pdf-without-disposition.pdf",
					},
				},
				ContentDisposition: ContentDispositionHeader{
					ContentDisposition: "",
					Params:             map[string]string(nil),
				},
				Data: []byte{37, 80, 68, 70, 45, 49, 46, 13, 116, 114, 97, 105, 108, 101, 114, 60, 60,
					47, 82, 111, 111, 116, 60, 60, 47, 80, 97, 103, 101, 115, 60, 60, 47, 75, 105, 100, 115, 91, 60,
					60, 47, 77, 101, 100, 105, 97, 66, 111, 120, 91, 48, 32, 48, 32, 51, 32, 51, 93, 62, 62, 93, 62,
					62, 62, 62, 62, 62},
			},
			{
				ContentType: ContentTypeHeader{
					ContentType: "application/json",
					Params: map[string]string{
						"name": "attached-json-name.json",
					},
				},
				ContentDisposition: ContentDispositionHeader{
					ContentDisposition: attachment,
					Params: map[string]string{
						"filename": "attached-json-filename.json",
					},
				},
				Data: []byte{123, 34, 102, 111, 111, 34, 58, 34, 98, 97, 114, 34, 125},
			},
			{
				ContentType: ContentTypeHeader{
					ContentType: "text/plain",
					Params: map[string]string{
						"name": "attached-text-plain-name.txt",
					},
				},
				ContentDisposition: ContentDispositionHeader{
					ContentDisposition: attachment,
					Params: map[string]string{
						"filename": "attached-text-plain-filename.txt",
					},
				},
				Data: []byte{84, 101, 120, 116, 47, 112, 108, 97, 105, 110, 32, 99, 111, 110, 116, 101, 110, 116, 32,
					97, 115, 32, 97, 110, 32, 97, 116, 116, 97, 99, 104, 101, 100, 32, 46, 116, 120, 116, 32, 102, 105,
					108, 101, 46},
			},
			{
				ContentType: ContentTypeHeader{
					ContentType: "text/html",
					Params: map[string]string{
						"name": "attached-text-html-name.html",
					},
				},
				ContentDisposition: ContentDispositionHeader{
					ContentDisposition: attachment,
					Params: map[string]string{
						"filename": "attached-text-html-filename.html",
					},
				},
				Data: []byte{84, 101, 120, 116, 47, 104, 116, 109, 108, 32, 99, 111, 110, 116, 101, 110, 116, 32, 97,
					115, 32, 97, 110, 32, 97, 116, 116, 97, 99, 104, 101, 100, 32, 46, 104, 116, 109, 108, 32, 102,
					105, 108, 101, 46},
			},
		},
	}

	testEmailFromFile(t, fp, expectedEmail)
}
func TestParseEmailFinnishMultipartMixedIso885915OverQuotedprintable(t *testing.T) {
	fp := "tests/test_finnish_multipart_mixed_iso-8859-15_over_quoted-printable.txt"
	tz, _ := time.LoadLocation("Europe/Helsinki")
	expectedDate, _ := time.Parse(
		time.RFC1123Z+" (MST)",
		time.Date(2019, time.April, 1, 7, 55, 0, 0, tz).Format(time.RFC1123Z+" (MST)"))
	expectedEmail := Email{
		Headers: Headers{
			Date:    expectedDate,
			Subject: "📧 Test Suomenkieliset pangrammit",
			ReplyTo: []*mail.Address{
				{
					Name:    "Alice Lähettäjä",
					Address: "alice.lahettaja@example.net",
				},
			},
			Sender: &mail.Address{
				Name:    "Alice Lähettäjä",
				Address: "alice.lahettaja@example.com",
			},
			From: []*mail.Address{
				{
					Name:    "Alice Lähettäjä",
					Address: "alice.lahettaja@example.com",
				},
				{
					Name:    "Alice Lähettäjä",
					Address: "alice.lahettaja@example.net",
				},
			},
			To: []*mail.Address{
				{
					Name:    "Bob Vastaanottaja",
					Address: "bob.vastaanottaja@exaple.com",
				},
				{
					Name:    "Carol Vastaanottaja",
					Address: "carol.vastaanottaja@example.com",
				},
			},
			Cc: []*mail.Address{
				{
					Name:    "Dan Vastaanottaja",
					Address: "dan.vastaanottaja@example.com",
				},
				{
					Name:    "Eve Vastaanottaja",
					Address: "eve.vastaanottaja@example.com",
				},
			},
			Bcc: []*mail.Address{
				{
					Name:    "Frank Vastaanottaja",
					Address: "frank.vastaanottaja@example.com",
				},
				{
					Name:    "Grace Vastaanottaja",
					Address: "grace.vastaanottaja@example.com",
				},
			},
			MessageID:  "Message-Id-1@example.com",
			InReplyTo:  []MessageId{"Message-Id-0@example.com"},
			References: []MessageId{"Message-Id-0@example.com"},
			Comments:   "Message Header Comment",
			Keywords:   []string{"Keyword 1", "Keyword 2"},
			ResentDate: expectedDate,
			ResentFrom: []*mail.Address{
				{
					Name:    "Alice Lähettäjä",
					Address: "alice.lahettaja@example.net",
				},
				{
					Name:    "Alice Lähettäjä",
					Address: "alice.lahettaja@example.com",
				},
			},
			ResentSender: &mail.Address{
				Name:    "Alice Lähettäjä",
				Address: "alice.lahettaja@example.net",
			},
			ResentTo: []*mail.Address{
				{
					Name:    "Bob Vastaanottaja",
					Address: "bob.vastaanottaja@exaple.com",
				},
				{
					Name:    "Carol Vastaanottaja",
					Address: "carol.vastaanottaja@example.net",
				},
			},
			ResentCc: []*mail.Address{
				{
					Name:    "Dan Vastaanottaja",
					Address: "dan.vastaanottaja@example.net",
				},
				{
					Name:    "Eve Vastaanottaja",
					Address: "eve.vastaanottaja@example.net",
				},
			},
			ResentBcc: []*mail.Address{
				{
					Name:    "Frank Vastaanottaja",
					Address: "frank.vastaanottaja@example.net",
				},
				{
					Name:    "Grace Vastaanottaja",
					Address: "grace.vastaanottaja@example.net",
				},
			},
			ResentMessageID: "Message-Id-1@example.net",
			ContentType: ContentTypeHeader{
				ContentType: "multipart/mixed",
				Params: map[string]string{
					"boundary": "MixedBoundaryString",
					"charset":  "iso-8859-15",
				},
			},
			ExtraHeaders: map[string][]string{
				"X-Clacks-Overhead": {"GNU Terry Pratchett"},
				"X-Script/function/\t !\"#$%&'()*+,-./;<=>?@[\\]^_`{|}~": {
					"TEST VALUE 1\t !\"#$%&'()*+,-./0123456789:;<=>?@ABCDEFGHIJKLMNOPQRSTUVWXYZ[\\]^_` abcdefghijklmnopqrstuvwxyz{|}~",
					"TEST VALUE 2\t !\"#$%&'()*+,-./0123456789:;<=>?@ABCDEFGHIJKLMNOPQRSTUVWXYZ[\\]^_` abcdefghijklmnopqrstuvwxyz{|}~",
				},
			},
		},
		Text: `Albert osti fagotin ja töräytti puhkuvan melodian.
Lorun sangen pieneksi hyödyksi jäivät suomen kirjaimet.
Hyvän lorun sangen pieneksi hyödyksi jäi suomen kirjaimet.
Fahrenheit ja Celsius yrjösivät Åsan backgammon-peliin, Volkswagenissa, daiquirin ja ZX81:n yhteisvaikutuksesta.
Charles Darwin jammaili Åken hevixylofonilla Qatarin yöpub Zeligissä.
Wieniläinen sioux:ta puhuva ökyzombie diggaa Åsan roquefort-tacoja.`,
		EnrichedText: `<bold>Albert osti fagotin ja töräytti puhkuvan melodian.</bold>
<italic>Lorun sangen pieneksi hyödyksi jäivät suomen kirjaimet.</italic>
<fixed>Hyvän lorun sangen pieneksi hyödyksi jäi suomen kirjaimet.</fixed>
<underline>Fahrenheit ja Celsius yrjösivät Åsan backgammon-peliin, Volkswagenissa, daiquirin ja ZX81:n yhteisvaikutuksesta.</underline>
Charles Darwin jammaili Åken hevixylofonilla Qatarin yöpub Zeligissä.
Wieniläinen sioux:ta puhuva ökyzombie diggaa Åsan roquefort-tacoja.`,
		HTML: `<html>
<div dir="ltr">
<p>Albert osti fagotin ja töräytti puhkuvan melodian.</p>
<p>Lorun sangen pieneksi hyödyksi jäivät suomen kirjaimet.</p>
<p>Hyvän lorun sangen pieneksi hyödyksi jäi suomen kirjaimet.</p>
<p>Fahrenheit ja Celsius yrjösivät Åsan backgammon-peliin, Volkswagenissa, daiquirin ja ZX81:n yhteisvaikutuksesta.</p>
<p>Charles Darwin jammaili Åken hevixylofonilla Qatarin yöpub Zeligissä.</p>
<p>Wieniläinen sioux:ta puhuva ökyzombie diggaa Åsan roquefort-tacoja.</p>
</div>
</html>`,
		InlineFiles: []InlineFile{
			{
				ContentType: ContentTypeHeader{
					ContentType: "image/jpeg",
					Params: map[string]string{
						"name": "inline-jpg-image-without-disposition.jpg",
					},
				},
				ContentDisposition: ContentDispositionHeader{
					ContentDisposition: "",
					Params:             map[string]string(nil),
				},
				Data: []byte{255, 216, 255, 219, 0, 67, 0, 3, 2, 2, 2, 2, 2, 3, 2, 2, 2, 3, 3, 3, 3, 4,
					6, 4, 4, 4, 4, 4, 8, 6, 6, 5, 6, 9, 8, 10, 10, 9, 8, 9, 9, 10, 12, 15, 12, 10, 11, 14, 11, 9, 9,
					13, 17, 13, 14, 15, 16, 16, 17, 16, 10, 12, 18, 19, 18, 16, 19, 15, 16, 16, 16, 255, 201, 0, 11, 8,
					0, 1, 0, 1, 1, 1, 17, 0, 255, 204, 0, 6, 0, 16, 16, 5, 255, 218, 0, 8, 1, 1, 0, 0, 63, 0, 210, 207,
					32, 255, 217},
			},
			{
				ContentID: "inline-jpg-image.jpg@example.com",
				ContentType: ContentTypeHeader{
					ContentType: "image/jpeg",
					Params: map[string]string{
						"name": "inline-jpg-image-name.jpg",
					},
				},
				ContentDisposition: ContentDispositionHeader{
					ContentDisposition: inline,
					Params: map[string]string{
						"filename": "inline-jpg-image-filename.jpg",
					},
				},
				Data: []byte{255, 216, 255, 219, 0, 67, 0, 3, 2, 2, 2, 2, 2, 3, 2, 2, 2, 3, 3, 3, 3, 4,
					6, 4, 4, 4, 4, 4, 8, 6, 6, 5, 6, 9, 8, 10, 10, 9, 8, 9, 9, 10, 12, 15, 12, 10, 11, 14, 11, 9, 9,
					13, 17, 13, 14, 15, 16, 16, 17, 16, 10, 12, 18, 19, 18, 16, 19, 15, 16, 16, 16, 255, 201, 0, 11, 8,
					0, 1, 0, 1, 1, 1, 17, 0, 255, 204, 0, 6, 0, 16, 16, 5, 255, 218, 0, 8, 1, 1, 0, 0, 63, 0, 210, 207,
					32, 255, 217},
			},
		},
		AttachedFiles: []AttachedFile{
			{
				ContentType: ContentTypeHeader{
					ContentType: "application/pdf",
					Params: map[string]string{
						"name": "attached-pdf-name.pdf",
					},
				},
				ContentDisposition: ContentDispositionHeader{
					ContentDisposition: attachment,
					Params: map[string]string{
						"filename": "attached-pdf-filename.pdf",
					},
				},
				Data: []byte{37, 80, 68, 70, 45, 49, 46, 13, 116, 114, 97, 105, 108, 101, 114, 60, 60,
					47, 82, 111, 111, 116, 60, 60, 47, 80, 97, 103, 101, 115, 60, 60, 47, 75, 105, 100, 115, 91, 60,
					60, 47, 77, 101, 100, 105, 97, 66, 111, 120, 91, 48, 32, 48, 32, 51, 32, 51, 93, 62, 62, 93, 62,
					62, 62, 62, 62, 62},
			},
			{
				ContentType: ContentTypeHeader{
					ContentType: "application/pdf",
					Params: map[string]string{
						"name": "attached-pdf-without-disposition.pdf",
					},
				},
				ContentDisposition: ContentDispositionHeader{
					ContentDisposition: "",
					Params:             map[string]string(nil),
				},
				Data: []byte{37, 80, 68, 70, 45, 49, 46, 13, 116, 114, 97, 105, 108, 101, 114, 60, 60,
					47, 82, 111, 111, 116, 60, 60, 47, 80, 97, 103, 101, 115, 60, 60, 47, 75, 105, 100, 115, 91, 60,
					60, 47, 77, 101, 100, 105, 97, 66, 111, 120, 91, 48, 32, 48, 32, 51, 32, 51, 93, 62, 62, 93, 62,
					62, 62, 62, 62, 62},
			},
			{
				ContentType: ContentTypeHeader{
					ContentType: "application/json",
					Params: map[string]string{
						"name": "attached-json-name.json",
					},
				},
				ContentDisposition: ContentDispositionHeader{
					ContentDisposition: attachment,
					Params: map[string]string{
						"filename": "attached-json-filename.json",
					},
				},
				Data: []byte{123, 34, 102, 111, 111, 34, 58, 34, 98, 97, 114, 34, 125},
			},
			{
				ContentType: ContentTypeHeader{
					ContentType: "text/plain",
					Params: map[string]string{
						"name": "attached-text-plain-name.txt",
					},
				},
				ContentDisposition: ContentDispositionHeader{
					ContentDisposition: attachment,
					Params: map[string]string{
						"filename": "attached-text-plain-filename.txt",
					},
				},
				Data: []byte{84, 101, 120, 116, 47, 112, 108, 97, 105, 110, 32, 99, 111, 110, 116, 101, 110, 116, 32,
					97, 115, 32, 97, 110, 32, 97, 116, 116, 97, 99, 104, 101, 100, 32, 46, 116, 120, 116, 32, 102, 105,
					108, 101, 46},
			},
			{
				ContentType: ContentTypeHeader{
					ContentType: "text/html",
					Params: map[string]string{
						"name": "attached-text-html-name.html",
					},
				},
				ContentDisposition: ContentDispositionHeader{
					ContentDisposition: attachment,
					Params: map[string]string{
						"filename": "attached-text-html-filename.html",
					},
				},
				Data: []byte{84, 101, 120, 116, 47, 104, 116, 109, 108, 32, 99, 111, 110, 116, 101, 110, 116, 32, 97,
					115, 32, 97, 110, 32, 97, 116, 116, 97, 99, 104, 101, 100, 32, 46, 104, 116, 109, 108, 32, 102,
					105, 108, 101, 46},
			},
		},
	}

	testEmailFromFile(t, fp, expectedEmail)
}
func TestParseEmailFinnishMultipartSignedUtf8OverBase64(t *testing.T) {
	fp := "tests/test_finnish_multipart_signed_utf-8_over_base64.txt"
	tz, _ := time.LoadLocation("Europe/Helsinki")
	expectedDate, _ := time.Parse(
		time.RFC1123Z+" (MST)",
		time.Date(2019, time.April, 1, 7, 55, 0, 0, tz).Format(time.RFC1123Z+" (MST)"))
	expectedEmail := Email{
		Headers: Headers{
			Date:    expectedDate,
			Subject: "📧 Signed Test Suomenkieliset pangrammit",
			ReplyTo: []*mail.Address{
				{
					Name:    "Alice Lähettäjä",
					Address: "alice.lahettaja@example.net",
				},
			},
			Sender: &mail.Address{
				Name:    "Alice Lähettäjä",
				Address: "alice.lahettaja@example.com",
			},
			From: []*mail.Address{
				{
					Name:    "Alice Lähettäjä",
					Address: "alice.lahettaja@example.com",
				},
				{
					Name:    "Alice Lähettäjä",
					Address: "alice.lahettaja@example.net",
				},
			},
			To: []*mail.Address{
				{
					Name:    "Bob Vastaanottaja",
					Address: "bob.vastaanottaja@exaple.com",
				},
				{
					Name:    "Carol Vastaanottaja",
					Address: "carol.vastaanottaja@example.com",
				},
			},
			Cc: []*mail.Address{
				{
					Name:    "Dan Vastaanottaja",
					Address: "dan.vastaanottaja@example.com",
				},
				{
					Name:    "Eve Vastaanottaja",
					Address: "eve.vastaanottaja@example.com",
				},
			},
			Bcc: []*mail.Address{
				{
					Name:    "Frank Vastaanottaja",
					Address: "frank.vastaanottaja@example.com",
				},
				{
					Name:    "Grace Vastaanottaja",
					Address: "grace.vastaanottaja@example.com",
				},
			},
			MessageID:  "Message-Id-1@example.com",
			InReplyTo:  []MessageId{"Message-Id-0@example.com"},
			References: []MessageId{"Message-Id-0@example.com"},
			Comments:   "Message Header Comment",
			Keywords:   []string{"Keyword 1", "Keyword 2"},
			ResentDate: expectedDate,
			ResentFrom: []*mail.Address{
				{
					Name:    "Alice Lähettäjä",
					Address: "alice.lahettaja@example.net",
				},
				{
					Name:    "Alice Lähettäjä",
					Address: "alice.lahettaja@example.com",
				},
			},
			ResentSender: &mail.Address{
				Name:    "Alice Lähettäjä",
				Address: "alice.lahettaja@example.net",
			},
			ResentTo: []*mail.Address{
				{
					Name:    "Bob Vastaanottaja",
					Address: "bob.vastaanottaja@exaple.com",
				},
				{
					Name:    "Carol Vastaanottaja",
					Address: "carol.vastaanottaja@example.net",
				},
			},
			ResentCc: []*mail.Address{
				{
					Name:    "Dan Vastaanottaja",
					Address: "dan.vastaanottaja@example.net",
				},
				{
					Name:    "Eve Vastaanottaja",
					Address: "eve.vastaanottaja@example.net",
				},
			},
			ResentBcc: []*mail.Address{
				{
					Name:    "Frank Vastaanottaja",
					Address: "frank.vastaanottaja@example.net",
				},
				{
					Name:    "Grace Vastaanottaja",
					Address: "grace.vastaanottaja@example.net",
				},
			},
			ResentMessageID: "Message-Id-1@example.net",
			ContentType: ContentTypeHeader{
				ContentType: "multipart/signed",
				Params: map[string]string{
					"boundary": "SignedBoundaryString",
					"charset":  "utf-8",
					"micalg":   "sha1",
					"protocol": "application/pkcs7-signature",
				},
			},
			ExtraHeaders: map[string][]string{
				"X-Clacks-Overhead": {"GNU Terry Pratchett"},
				"X-Script/function/\t !\"#$%&'()*+,-./;<=>?@[\\]^_`{|}~": {
					"TEST VALUE 1\t !\"#$%&'()*+,-./0123456789:;<=>?@ABCDEFGHIJKLMNOPQRSTUVWXYZ[\\]^_` abcdefghijklmnopqrstuvwxyz{|}~",
					"TEST VALUE 2\t !\"#$%&'()*+,-./0123456789:;<=>?@ABCDEFGHIJKLMNOPQRSTUVWXYZ[\\]^_` abcdefghijklmnopqrstuvwxyz{|}~",
				},
			},
		},
		Text: `Albert osti fagotin ja töräytti puhkuvan melodian.
Lorun sangen pieneksi hyödyksi jäivät suomen kirjaimet.
Hyvän lorun sangen pieneksi hyödyksi jäi suomen kirjaimet.
Fahrenheit ja Celsius yrjösivät Åsan backgammon-peliin, Volkswagenissa, daiquirin ja ZX81:n yhteisvaikutuksesta.
Charles Darwin jammaili Åken hevixylofonilla Qatarin yöpub Zeligissä.
Wieniläinen sioux:ta puhuva ökyzombie diggaa Åsan roquefort-tacoja.`,
		AttachedFiles: []AttachedFile{
			{
				ContentType: ContentTypeHeader{
					ContentType: "application/pkcs7-signature",
					Params: map[string]string{
						"name": "smime.p7s",
					},
				},
				ContentDisposition: ContentDispositionHeader{
					ContentDisposition: attachment,
					Params: map[string]string{
						"filename": "smime.p7s",
					},
				},
				Data: []byte{130, 28, 135, 132, 117, 46, 142, 18, 97, 140, 126, 251, 159, 193, 199, 25, 58, 223, 189, 185, 227, 239, 158, 173, 108, 31, 71, 27, 133, 80, 165, 252, 133, 227, 174, 198, 132, 129, 159, 29, 246, 19, 235, 133, 80, 165, 252, 133, 227, 174, 198, 132, 129, 159, 29, 246, 19, 234, 49, 251, 238, 127, 7, 28, 104, 33, 200, 120, 71, 82, 232, 225, 38, 30, 249, 234, 214, 193, 244, 113, 147, 173, 251, 219, 158, 57, 252, 28, 113, 147, 173, 251, 225, 38, 24, 199, 239, 190, 173, 108, 31, 71, 27, 133, 80, 110, 120, 251, 231, 174, 198, 132, 129, 159, 29, 246, 19, 234, 8, 114, 30, 17, 212, 186, 58, 95, 200, 94, 59, 26, 18, 6, 124, 119, 216, 79, 174, 21, 65, 185, 227, 239, 158},
			},
		},
	}

	testEmailFromFile(t, fp, expectedEmail)
}

func TestParseEmailFinnishMultipartSignedUtf8OverQuotedprintable(t *testing.T) {
	fp := "tests/test_finnish_multipart_signed_utf-8_over_quoted-printable.txt"
	tz, _ := time.LoadLocation("Europe/Helsinki")
	expectedDate, _ := time.Parse(
		time.RFC1123Z+" (MST)",
		time.Date(2019, time.April, 1, 7, 55, 0, 0, tz).Format(time.RFC1123Z+" (MST)"))
	expectedEmail := Email{
		Headers: Headers{
			Date:    expectedDate,
			Subject: "📧 Signed Test Suomenkieliset pangrammit",
			ReplyTo: []*mail.Address{
				{
					Name:    "Alice Lähettäjä",
					Address: "alice.lahettaja@example.net",
				},
			},
			Sender: &mail.Address{
				Name:    "Alice Lähettäjä",
				Address: "alice.lahettaja@example.com",
			},
			From: []*mail.Address{
				{
					Name:    "Alice Lähettäjä",
					Address: "alice.lahettaja@example.com",
				},
				{
					Name:    "Alice Lähettäjä",
					Address: "alice.lahettaja@example.net",
				},
			},
			To: []*mail.Address{
				{
					Name:    "Bob Vastaanottaja",
					Address: "bob.vastaanottaja@exaple.com",
				},
				{
					Name:    "Carol Vastaanottaja",
					Address: "carol.vastaanottaja@example.com",
				},
			},
			Cc: []*mail.Address{
				{
					Name:    "Dan Vastaanottaja",
					Address: "dan.vastaanottaja@example.com",
				},
				{
					Name:    "Eve Vastaanottaja",
					Address: "eve.vastaanottaja@example.com",
				},
			},
			Bcc: []*mail.Address{
				{
					Name:    "Frank Vastaanottaja",
					Address: "frank.vastaanottaja@example.com",
				},
				{
					Name:    "Grace Vastaanottaja",
					Address: "grace.vastaanottaja@example.com",
				},
			},
			MessageID:  "Message-Id-1@example.com",
			InReplyTo:  []MessageId{"Message-Id-0@example.com"},
			References: []MessageId{"Message-Id-0@example.com"},
			Comments:   "Message Header Comment",
			Keywords:   []string{"Keyword 1", "Keyword 2"},
			ResentDate: expectedDate,
			ResentFrom: []*mail.Address{
				{
					Name:    "Alice Lähettäjä",
					Address: "alice.lahettaja@example.net",
				},
				{
					Name:    "Alice Lähettäjä",
					Address: "alice.lahettaja@example.com",
				},
			},
			ResentSender: &mail.Address{
				Name:    "Alice Lähettäjä",
				Address: "alice.lahettaja@example.net",
			},
			ResentTo: []*mail.Address{
				{
					Name:    "Bob Vastaanottaja",
					Address: "bob.vastaanottaja@exaple.com",
				},
				{
					Name:    "Carol Vastaanottaja",
					Address: "carol.vastaanottaja@example.net",
				},
			},
			ResentCc: []*mail.Address{
				{
					Name:    "Dan Vastaanottaja",
					Address: "dan.vastaanottaja@example.net",
				},
				{
					Name:    "Eve Vastaanottaja",
					Address: "eve.vastaanottaja@example.net",
				},
			},
			ResentBcc: []*mail.Address{
				{
					Name:    "Frank Vastaanottaja",
					Address: "frank.vastaanottaja@example.net",
				},
				{
					Name:    "Grace Vastaanottaja",
					Address: "grace.vastaanottaja@example.net",
				},
			},
			ResentMessageID: "Message-Id-1@example.net",
			ContentType: ContentTypeHeader{
				ContentType: "multipart/signed",
				Params: map[string]string{
					"boundary": "SignedBoundaryString",
					"charset":  "utf-8",
					"micalg":   "sha1",
					"protocol": "application/pkcs7-signature",
				},
			},
			ExtraHeaders: map[string][]string{
				"X-Clacks-Overhead": {"GNU Terry Pratchett"},
				"X-Script/function/\t !\"#$%&'()*+,-./;<=>?@[\\]^_`{|}~": {
					"TEST VALUE 1\t !\"#$%&'()*+,-./0123456789:;<=>?@ABCDEFGHIJKLMNOPQRSTUVWXYZ[\\]^_` abcdefghijklmnopqrstuvwxyz{|}~",
					"TEST VALUE 2\t !\"#$%&'()*+,-./0123456789:;<=>?@ABCDEFGHIJKLMNOPQRSTUVWXYZ[\\]^_` abcdefghijklmnopqrstuvwxyz{|}~",
				},
			},
		},
		Text: `Albert osti fagotin ja töräytti puhkuvan melodian.
Lorun sangen pieneksi hyödyksi jäivät suomen kirjaimet.
Hyvän lorun sangen pieneksi hyödyksi jäi suomen kirjaimet.
Fahrenheit ja Celsius yrjösivät Åsan backgammon-peliin, Volkswagenissa, daiquirin ja ZX81:n yhteisvaikutuksesta.
Charles Darwin jammaili Åken hevixylofonilla Qatarin yöpub Zeligissä.
Wieniläinen sioux:ta puhuva ökyzombie diggaa Åsan roquefort-tacoja.`,
		AttachedFiles: []AttachedFile{
			{
				ContentType: ContentTypeHeader{
					ContentType: "application/pkcs7-signature",
					Params: map[string]string{
						"name": "smime.p7s",
					},
				},
				ContentDisposition: ContentDispositionHeader{
					ContentDisposition: attachment,
					Params: map[string]string{
						"filename": "smime.p7s",
					},
				},
				Data: []byte{130, 28, 135, 132, 117, 46, 142, 18, 97, 140, 126, 251, 159, 193, 199, 25, 58, 223, 189, 185, 227, 239, 158, 173, 108, 31, 71, 27, 133, 80, 165, 252, 133, 227, 174, 198, 132, 129, 159, 29, 246, 19, 235, 133, 80, 165, 252, 133, 227, 174, 198, 132, 129, 159, 29, 246, 19, 234, 49, 251, 238, 127, 7, 28, 104, 33, 200, 120, 71, 82, 232, 225, 38, 30, 249, 234, 214, 193, 244, 113, 147, 173, 251, 219, 158, 57, 252, 28, 113, 147, 173, 251, 225, 38, 24, 199, 239, 190, 173, 108, 31, 71, 27, 133, 80, 110, 120, 251, 231, 174, 198, 132, 129, 159, 29, 246, 19, 234, 8, 114, 30, 17, 212, 186, 58, 95, 200, 94, 59, 26, 18, 6, 124, 119, 216, 79, 174, 21, 65, 185, 227, 239, 158},
			},
		},
	}

	testEmailFromFile(t, fp, expectedEmail)
}

func TestParseEmailFinnishMultipartSignedIso885915OverBase64(t *testing.T) {
	fp := "tests/test_finnish_multipart_signed_iso-8859-15_over_base64.txt"
	tz, _ := time.LoadLocation("Europe/Helsinki")
	expectedDate, _ := time.Parse(
		time.RFC1123Z+" (MST)",
		time.Date(2019, time.April, 1, 7, 55, 0, 0, tz).Format(time.RFC1123Z+" (MST)"))
	expectedEmail := Email{
		Headers: Headers{
			Date:    expectedDate,
			Subject: "📧 Signed Test Suomenkieliset pangrammit",
			ReplyTo: []*mail.Address{
				{
					Name:    "Alice Lähettäjä",
					Address: "alice.lahettaja@example.net",
				},
			},
			Sender: &mail.Address{
				Name:    "Alice Lähettäjä",
				Address: "alice.lahettaja@example.com",
			},
			From: []*mail.Address{
				{
					Name:    "Alice Lähettäjä",
					Address: "alice.lahettaja@example.com",
				},
				{
					Name:    "Alice Lähettäjä",
					Address: "alice.lahettaja@example.net",
				},
			},
			To: []*mail.Address{
				{
					Name:    "Bob Vastaanottaja",
					Address: "bob.vastaanottaja@exaple.com",
				},
				{
					Name:    "Carol Vastaanottaja",
					Address: "carol.vastaanottaja@example.com",
				},
			},
			Cc: []*mail.Address{
				{
					Name:    "Dan Vastaanottaja",
					Address: "dan.vastaanottaja@example.com",
				},
				{
					Name:    "Eve Vastaanottaja",
					Address: "eve.vastaanottaja@example.com",
				},
			},
			Bcc: []*mail.Address{
				{
					Name:    "Frank Vastaanottaja",
					Address: "frank.vastaanottaja@example.com",
				},
				{
					Name:    "Grace Vastaanottaja",
					Address: "grace.vastaanottaja@example.com",
				},
			},
			MessageID:  "Message-Id-1@example.com",
			InReplyTo:  []MessageId{"Message-Id-0@example.com"},
			References: []MessageId{"Message-Id-0@example.com"},
			Comments:   "Message Header Comment",
			Keywords:   []string{"Keyword 1", "Keyword 2"},
			ResentDate: expectedDate,
			ResentFrom: []*mail.Address{
				{
					Name:    "Alice Lähettäjä",
					Address: "alice.lahettaja@example.net",
				},
				{
					Name:    "Alice Lähettäjä",
					Address: "alice.lahettaja@example.com",
				},
			},
			ResentSender: &mail.Address{
				Name:    "Alice Lähettäjä",
				Address: "alice.lahettaja@example.net",
			},
			ResentTo: []*mail.Address{
				{
					Name:    "Bob Vastaanottaja",
					Address: "bob.vastaanottaja@exaple.com",
				},
				{
					Name:    "Carol Vastaanottaja",
					Address: "carol.vastaanottaja@example.net",
				},
			},
			ResentCc: []*mail.Address{
				{
					Name:    "Dan Vastaanottaja",
					Address: "dan.vastaanottaja@example.net",
				},
				{
					Name:    "Eve Vastaanottaja",
					Address: "eve.vastaanottaja@example.net",
				},
			},
			ResentBcc: []*mail.Address{
				{
					Name:    "Frank Vastaanottaja",
					Address: "frank.vastaanottaja@example.net",
				},
				{
					Name:    "Grace Vastaanottaja",
					Address: "grace.vastaanottaja@example.net",
				},
			},
			ResentMessageID: "Message-Id-1@example.net",
			ContentType: ContentTypeHeader{
				ContentType: "multipart/signed",
				Params: map[string]string{
					"boundary": "SignedBoundaryString",
					"charset":  "iso-8859-15",
					"micalg":   "sha1",
					"protocol": "application/pkcs7-signature",
				},
			},
			ExtraHeaders: map[string][]string{
				"X-Clacks-Overhead": {"GNU Terry Pratchett"},
				"X-Script/function/\t !\"#$%&'()*+,-./;<=>?@[\\]^_`{|}~": {
					"TEST VALUE 1\t !\"#$%&'()*+,-./0123456789:;<=>?@ABCDEFGHIJKLMNOPQRSTUVWXYZ[\\]^_` abcdefghijklmnopqrstuvwxyz{|}~",
					"TEST VALUE 2\t !\"#$%&'()*+,-./0123456789:;<=>?@ABCDEFGHIJKLMNOPQRSTUVWXYZ[\\]^_` abcdefghijklmnopqrstuvwxyz{|}~",
				},
			},
		},
		Text: `Albert osti fagotin ja töräytti puhkuvan melodian.
Lorun sangen pieneksi hyödyksi jäivät suomen kirjaimet.
Hyvän lorun sangen pieneksi hyödyksi jäi suomen kirjaimet.
Fahrenheit ja Celsius yrjösivät Åsan backgammon-peliin, Volkswagenissa, daiquirin ja ZX81:n yhteisvaikutuksesta.
Charles Darwin jammaili Åken hevixylofonilla Qatarin yöpub Zeligissä.
Wieniläinen sioux:ta puhuva ökyzombie diggaa Åsan roquefort-tacoja.`,
		AttachedFiles: []AttachedFile{
			{
				ContentType: ContentTypeHeader{
					ContentType: "application/pkcs7-signature",
					Params: map[string]string{
						"name": "smime.p7s",
					},
				},
				ContentDisposition: ContentDispositionHeader{
					ContentDisposition: attachment,
					Params: map[string]string{
						"filename": "smime.p7s",
					},
				},
				Data: []byte{130, 28, 135, 132, 117, 46, 142, 18, 97, 140, 126, 251, 159, 193, 199, 25, 58, 223, 189, 185, 227, 239, 158, 173, 108, 31, 71, 27, 133, 80, 165, 252, 133, 227, 174, 198, 132, 129, 159, 29, 246, 19, 235, 133, 80, 165, 252, 133, 227, 174, 198, 132, 129, 159, 29, 246, 19, 234, 49, 251, 238, 127, 7, 28, 104, 33, 200, 120, 71, 82, 232, 225, 38, 30, 249, 234, 214, 193, 244, 113, 147, 173, 251, 219, 158, 57, 252, 28, 113, 147, 173, 251, 225, 38, 24, 199, 239, 190, 173, 108, 31, 71, 27, 133, 80, 110, 120, 251, 231, 174, 198, 132, 129, 159, 29, 246, 19, 234, 8, 114, 30, 17, 212, 186, 58, 95, 200, 94, 59, 26, 18, 6, 124, 119, 216, 79, 174, 21, 65, 185, 227, 239, 158},
			},
		},
	}

	testEmailFromFile(t, fp, expectedEmail)
}

func TestParseEmailFinnishMultipartSignedIso885915OverQuotedprintable(t *testing.T) {
	fp := "tests/test_finnish_multipart_signed_iso-8859-15_over_quoted-printable.txt"
	tz, _ := time.LoadLocation("Europe/Helsinki")
	expectedDate, _ := time.Parse(
		time.RFC1123Z+" (MST)",
		time.Date(2019, time.April, 1, 7, 55, 0, 0, tz).Format(time.RFC1123Z+" (MST)"))
	expectedEmail := Email{
		Headers: Headers{
			Date:    expectedDate,
			Subject: "📧 Signed Test Suomenkieliset pangrammit",
			ReplyTo: []*mail.Address{
				{
					Name:    "Alice Lähettäjä",
					Address: "alice.lahettaja@example.net",
				},
			},
			Sender: &mail.Address{
				Name:    "Alice Lähettäjä",
				Address: "alice.lahettaja@example.com",
			},
			From: []*mail.Address{
				{
					Name:    "Alice Lähettäjä",
					Address: "alice.lahettaja@example.com",
				},
				{
					Name:    "Alice Lähettäjä",
					Address: "alice.lahettaja@example.net",
				},
			},
			To: []*mail.Address{
				{
					Name:    "Bob Vastaanottaja",
					Address: "bob.vastaanottaja@exaple.com",
				},
				{
					Name:    "Carol Vastaanottaja",
					Address: "carol.vastaanottaja@example.com",
				},
			},
			Cc: []*mail.Address{
				{
					Name:    "Dan Vastaanottaja",
					Address: "dan.vastaanottaja@example.com",
				},
				{
					Name:    "Eve Vastaanottaja",
					Address: "eve.vastaanottaja@example.com",
				},
			},
			Bcc: []*mail.Address{
				{
					Name:    "Frank Vastaanottaja",
					Address: "frank.vastaanottaja@example.com",
				},
				{
					Name:    "Grace Vastaanottaja",
					Address: "grace.vastaanottaja@example.com",
				},
			},
			MessageID:  "Message-Id-1@example.com",
			InReplyTo:  []MessageId{"Message-Id-0@example.com"},
			References: []MessageId{"Message-Id-0@example.com"},
			Comments:   "Message Header Comment",
			Keywords:   []string{"Keyword 1", "Keyword 2"},
			ResentDate: expectedDate,
			ResentFrom: []*mail.Address{
				{
					Name:    "Alice Lähettäjä",
					Address: "alice.lahettaja@example.net",
				},
				{
					Name:    "Alice Lähettäjä",
					Address: "alice.lahettaja@example.com",
				},
			},
			ResentSender: &mail.Address{
				Name:    "Alice Lähettäjä",
				Address: "alice.lahettaja@example.net",
			},
			ResentTo: []*mail.Address{
				{
					Name:    "Bob Vastaanottaja",
					Address: "bob.vastaanottaja@exaple.com",
				},
				{
					Name:    "Carol Vastaanottaja",
					Address: "carol.vastaanottaja@example.net",
				},
			},
			ResentCc: []*mail.Address{
				{
					Name:    "Dan Vastaanottaja",
					Address: "dan.vastaanottaja@example.net",
				},
				{
					Name:    "Eve Vastaanottaja",
					Address: "eve.vastaanottaja@example.net",
				},
			},
			ResentBcc: []*mail.Address{
				{
					Name:    "Frank Vastaanottaja",
					Address: "frank.vastaanottaja@example.net",
				},
				{
					Name:    "Grace Vastaanottaja",
					Address: "grace.vastaanottaja@example.net",
				},
			},
			ResentMessageID: "Message-Id-1@example.net",
			ContentType: ContentTypeHeader{
				ContentType: "multipart/signed",
				Params: map[string]string{
					"boundary": "SignedBoundaryString",
					"charset":  "iso-8859-15",
					"micalg":   "sha1",
					"protocol": "application/pkcs7-signature",
				},
			},
			ExtraHeaders: map[string][]string{
				"X-Clacks-Overhead": {"GNU Terry Pratchett"},
				"X-Script/function/\t !\"#$%&'()*+,-./;<=>?@[\\]^_`{|}~": {
					"TEST VALUE 1\t !\"#$%&'()*+,-./0123456789:;<=>?@ABCDEFGHIJKLMNOPQRSTUVWXYZ[\\]^_` abcdefghijklmnopqrstuvwxyz{|}~",
					"TEST VALUE 2\t !\"#$%&'()*+,-./0123456789:;<=>?@ABCDEFGHIJKLMNOPQRSTUVWXYZ[\\]^_` abcdefghijklmnopqrstuvwxyz{|}~",
				},
			},
		},
		Text: `Albert osti fagotin ja töräytti puhkuvan melodian.
Lorun sangen pieneksi hyödyksi jäivät suomen kirjaimet.
Hyvän lorun sangen pieneksi hyödyksi jäi suomen kirjaimet.
Fahrenheit ja Celsius yrjösivät Åsan backgammon-peliin, Volkswagenissa, daiquirin ja ZX81:n yhteisvaikutuksesta.
Charles Darwin jammaili Åken hevixylofonilla Qatarin yöpub Zeligissä.
Wieniläinen sioux:ta puhuva ökyzombie diggaa Åsan roquefort-tacoja.`,
		AttachedFiles: []AttachedFile{
			{
				ContentType: ContentTypeHeader{
					ContentType: "application/pkcs7-signature",
					Params: map[string]string{
						"name": "smime.p7s",
					},
				},
				ContentDisposition: ContentDispositionHeader{
					ContentDisposition: attachment,
					Params: map[string]string{
						"filename": "smime.p7s",
					},
				},
				Data: []byte{130, 28, 135, 132, 117, 46, 142, 18, 97, 140, 126, 251, 159, 193, 199, 25, 58, 223, 189, 185, 227, 239, 158, 173, 108, 31, 71, 27, 133, 80, 165, 252, 133, 227, 174, 198, 132, 129, 159, 29, 246, 19, 235, 133, 80, 165, 252, 133, 227, 174, 198, 132, 129, 159, 29, 246, 19, 234, 49, 251, 238, 127, 7, 28, 104, 33, 200, 120, 71, 82, 232, 225, 38, 30, 249, 234, 214, 193, 244, 113, 147, 173, 251, 219, 158, 57, 252, 28, 113, 147, 173, 251, 225, 38, 24, 199, 239, 190, 173, 108, 31, 71, 27, 133, 80, 110, 120, 251, 231, 174, 198, 132, 129, 159, 29, 246, 19, 234, 8, 114, 30, 17, 212, 186, 58, 95, 200, 94, 59, 26, 18, 6, 124, 119, 216, 79, 174, 21, 65, 185, 227, 239, 158},
			},
		},
	}

	testEmailFromFile(t, fp, expectedEmail)
}

func TestParseEmailIcelandicPlaintextUtf8OverBase64(t *testing.T) {
	fp := "tests/test_icelandic_plaintext_utf-8_over_base64.txt"
	tz, _ := time.LoadLocation("Atlantic/Reykjavik")
	expectedDate, _ := time.Parse(
		time.RFC1123Z+" (MST)",
		time.Date(2019, time.April, 1, 7, 55, 0, 0, tz).Format(time.RFC1123Z+" (MST)"))
	expectedEmail := Email{
		Headers: Headers{
			Date:    expectedDate,
			Subject: "📧 Test Íslenskt pangram",
			ReplyTo: []*mail.Address{
				{
					Name:    "Alice Sendandidóttir",
					Address: "alice.sendandidottir@example.net",
				},
			},
			Sender: &mail.Address{
				Name:    "Alice Sendandidóttir",
				Address: "alice.sendandidottir@example.com",
			},
			From: []*mail.Address{
				{
					Name:    "Alice Sendandidóttir",
					Address: "alice.sendandidottir@example.com",
				},
				{
					Name:    "Alice Sendandidóttir",
					Address: "alice.sendandidottir@example.net",
				},
			},
			To: []*mail.Address{
				{
					Name:    "Bob Viðtakandison",
					Address: "bob.didtakandison@example.com",
				},
				{
					Name:    "Carol Viðtakandidóttir",
					Address: "carol.didtakandidottir@example.com",
				},
			},
			Cc: []*mail.Address{
				{
					Name:    "Dan Viðtakandison",
					Address: "dan.vidtakandison@example.com",
				},
				{
					Name:    "Eve Viðtakandidóttir",
					Address: "eve.vidtakandidottir@example.com",
				},
			},
			Bcc: []*mail.Address{
				{
					Name:    "Frank Viðtakandison",
					Address: "frank.vidtakandison@example.com",
				},
				{
					Name:    "Grace Viðtakandidóttir",
					Address: "grace.vidtakandidottir@example.com",
				},
			},
			MessageID:  "Message-Id-1@example.com",
			InReplyTo:  []MessageId{"Message-Id-0@example.com"},
			References: []MessageId{"Message-Id-0@example.com"},
			Comments:   "Message Header Comment",
			Keywords:   []string{"Keyword 1", "Keyword 2"},
			ResentDate: expectedDate,
			ResentFrom: []*mail.Address{
				{
					Name:    "Alice Sendandidóttir",
					Address: "alice.sendandidottir@example.net",
				},
				{
					Name:    "Alice Sendandidóttir",
					Address: "alice.sendandidottir@example.com",
				},
			},
			ResentSender: &mail.Address{
				Name:    "Alice Sendandidóttir",
				Address: "alice.sendandidottir@example.net",
			},
			ResentTo: []*mail.Address{
				{
					Name:    "Bob Viðtakandison",
					Address: "bob.didtakandison@example.net",
				},
				{
					Name:    "Carol Viðtakandidóttir",
					Address: "carol.didtakandidottir@example.net",
				},
			},
			ResentCc: []*mail.Address{
				{
					Name:    "Dan Viðtakandison",
					Address: "dan.vidtakandison@example.net",
				},
				{
					Name:    "Eve Viðtakandidóttir",
					Address: "eve.vidtakandidottir@example.net",
				},
			},
			ResentBcc: []*mail.Address{
				{
					Name:    "Frank Viðtakandison",
					Address: "frank.vidtakandison@example.net",
				},
				{
					Name:    "Grace Viðtakandidóttir",
					Address: "grace.vidtakandidottir@example.net",
				},
			},
			ResentMessageID: "Message-Id-1@example.net",
			ContentType: ContentTypeHeader{
				ContentType: "text/plain",
				Params: map[string]string{
					"charset": "utf-8",
				},
			},
			ExtraHeaders: map[string][]string{
				"X-Clacks-Overhead": {"GNU Terry Pratchett"},
				"X-Script/function/\t !\"#$%&'()*+,-./;<=>?@[\\]^_`{|}~": {
					"TEST VALUE 1\t !\"#$%&'()*+,-./0123456789:;<=>?@ABCDEFGHIJKLMNOPQRSTUVWXYZ[\\]^_` abcdefghijklmnopqrstuvwxyz{|}~",
					"TEST VALUE 2\t !\"#$%&'()*+,-./0123456789:;<=>?@ABCDEFGHIJKLMNOPQRSTUVWXYZ[\\]^_` abcdefghijklmnopqrstuvwxyz{|}~",
				},
			},
		},
		Text: `Kæmi ný öxi hér, ykist þjófum nú bæði víl og ádrepa.
Svo hölt, yxna kýr þegði jú um dóp í fé á bæ.
Þú dazt á hnéð í vök og yfir blóm sexý pæju.`,
	}

	testEmailFromFile(t, fp, expectedEmail)
}

func TestParseEmailIcelandicPlaintextUtf8OverQuotedprintable(t *testing.T) {
	fp := "tests/test_icelandic_plaintext_utf-8_over_quoted-printable.txt"
	tz, _ := time.LoadLocation("Atlantic/Reykjavik")
	expectedDate, _ := time.Parse(
		time.RFC1123Z+" (MST)",
		time.Date(2019, time.April, 1, 7, 55, 0, 0, tz).Format(time.RFC1123Z+" (MST)"))
	expectedEmail := Email{
		Headers: Headers{
			Date:    expectedDate,
			Subject: "📧 Test Íslenskt pangram",
			ReplyTo: []*mail.Address{
				{
					Name:    "Alice Sendandidóttir",
					Address: "alice.sendandidottir@example.net",
				},
			},
			Sender: &mail.Address{
				Name:    "Alice Sendandidóttir",
				Address: "alice.sendandidottir@example.com",
			},
			From: []*mail.Address{
				{
					Name:    "Alice Sendandidóttir",
					Address: "alice.sendandidottir@example.com",
				},
				{
					Name:    "Alice Sendandidóttir",
					Address: "alice.sendandidottir@example.net",
				},
			},
			To: []*mail.Address{
				{
					Name:    "Bob Viðtakandison",
					Address: "bob.didtakandison@example.com",
				},
				{
					Name:    "Carol Viðtakandidóttir",
					Address: "carol.didtakandidottir@example.com",
				},
			},
			Cc: []*mail.Address{
				{
					Name:    "Dan Viðtakandison",
					Address: "dan.vidtakandison@example.com",
				},
				{
					Name:    "Eve Viðtakandidóttir",
					Address: "eve.vidtakandidottir@example.com",
				},
			},
			Bcc: []*mail.Address{
				{
					Name:    "Frank Viðtakandison",
					Address: "frank.vidtakandison@example.com",
				},
				{
					Name:    "Grace Viðtakandidóttir",
					Address: "grace.vidtakandidottir@example.com",
				},
			},
			MessageID:  "Message-Id-1@example.com",
			InReplyTo:  []MessageId{"Message-Id-0@example.com"},
			References: []MessageId{"Message-Id-0@example.com"},
			Comments:   "Message Header Comment",
			Keywords:   []string{"Keyword 1", "Keyword 2"},
			ResentDate: expectedDate,
			ResentFrom: []*mail.Address{
				{
					Name:    "Alice Sendandidóttir",
					Address: "alice.sendandidottir@example.net",
				},
				{
					Name:    "Alice Sendandidóttir",
					Address: "alice.sendandidottir@example.com",
				},
			},
			ResentSender: &mail.Address{
				Name:    "Alice Sendandidóttir",
				Address: "alice.sendandidottir@example.net",
			},
			ResentTo: []*mail.Address{
				{
					Name:    "Bob Viðtakandison",
					Address: "bob.didtakandison@example.net",
				},
				{
					Name:    "Carol Viðtakandidóttir",
					Address: "carol.didtakandidottir@example.net",
				},
			},
			ResentCc: []*mail.Address{
				{
					Name:    "Dan Viðtakandison",
					Address: "dan.vidtakandison@example.net",
				},
				{
					Name:    "Eve Viðtakandidóttir",
					Address: "eve.vidtakandidottir@example.net",
				},
			},
			ResentBcc: []*mail.Address{
				{
					Name:    "Frank Viðtakandison",
					Address: "frank.vidtakandison@example.net",
				},
				{
					Name:    "Grace Viðtakandidóttir",
					Address: "grace.vidtakandidottir@example.net",
				},
			},
			ResentMessageID: "Message-Id-1@example.net",
			ContentType: ContentTypeHeader{
				ContentType: "text/plain",
				Params: map[string]string{
					"charset": "utf-8",
				},
			},
			ExtraHeaders: map[string][]string{
				"X-Clacks-Overhead": {"GNU Terry Pratchett"},
				"X-Script/function/\t !\"#$%&'()*+,-./;<=>?@[\\]^_`{|}~": {
					"TEST VALUE 1\t !\"#$%&'()*+,-./0123456789:;<=>?@ABCDEFGHIJKLMNOPQRSTUVWXYZ[\\]^_` abcdefghijklmnopqrstuvwxyz{|}~",
					"TEST VALUE 2\t !\"#$%&'()*+,-./0123456789:;<=>?@ABCDEFGHIJKLMNOPQRSTUVWXYZ[\\]^_` abcdefghijklmnopqrstuvwxyz{|}~",
				},
			},
		},
		Text: `Kæmi ný öxi hér, ykist þjófum nú bæði víl og ádrepa.
Svo hölt, yxna kýr þegði jú um dóp í fé á bæ.
Þú dazt á hnéð í vök og yfir blóm sexý pæju.`,
	}

	testEmailFromFile(t, fp, expectedEmail)
}

func TestParseEmailIcelandicPlaintextIso88591OverBase64(t *testing.T) {
	fp := "tests/test_icelandic_plaintext_iso-8859-1_over_base64.txt"
	tz, _ := time.LoadLocation("Atlantic/Reykjavik")
	expectedDate, _ := time.Parse(
		time.RFC1123Z+" (MST)",
		time.Date(2019, time.April, 1, 7, 55, 0, 0, tz).Format(time.RFC1123Z+" (MST)"))
	expectedEmail := Email{
		Headers: Headers{
			Date:    expectedDate,
			Subject: "📧 Test Íslenskt pangram",
			ReplyTo: []*mail.Address{
				{
					Name:    "Alice Sendandidóttir",
					Address: "alice.sendandidottir@example.net",
				},
			},
			Sender: &mail.Address{
				Name:    "Alice Sendandidóttir",
				Address: "alice.sendandidottir@example.com",
			},
			From: []*mail.Address{
				{
					Name:    "Alice Sendandidóttir",
					Address: "alice.sendandidottir@example.com",
				},
				{
					Name:    "Alice Sendandidóttir",
					Address: "alice.sendandidottir@example.net",
				},
			},
			To: []*mail.Address{
				{
					Name:    "Bob Viðtakandison",
					Address: "bob.didtakandison@example.com",
				},
				{
					Name:    "Carol Viðtakandidóttir",
					Address: "carol.didtakandidottir@example.com",
				},
			},
			Cc: []*mail.Address{
				{
					Name:    "Dan Viðtakandison",
					Address: "dan.vidtakandison@example.com",
				},
				{
					Name:    "Eve Viðtakandidóttir",
					Address: "eve.vidtakandidottir@example.com",
				},
			},
			Bcc: []*mail.Address{
				{
					Name:    "Frank Viðtakandison",
					Address: "frank.vidtakandison@example.com",
				},
				{
					Name:    "Grace Viðtakandidóttir",
					Address: "grace.vidtakandidottir@example.com",
				},
			},
			MessageID:  "Message-Id-1@example.com",
			InReplyTo:  []MessageId{"Message-Id-0@example.com"},
			References: []MessageId{"Message-Id-0@example.com"},
			Comments:   "Message Header Comment",
			Keywords:   []string{"Keyword 1", "Keyword 2"},
			ResentDate: expectedDate,
			ResentFrom: []*mail.Address{
				{
					Name:    "Alice Sendandidóttir",
					Address: "alice.sendandidottir@example.net",
				},
				{
					Name:    "Alice Sendandidóttir",
					Address: "alice.sendandidottir@example.com",
				},
			},
			ResentSender: &mail.Address{
				Name:    "Alice Sendandidóttir",
				Address: "alice.sendandidottir@example.net",
			},
			ResentTo: []*mail.Address{
				{
					Name:    "Bob Viðtakandison",
					Address: "bob.didtakandison@example.net",
				},
				{
					Name:    "Carol Viðtakandidóttir",
					Address: "carol.didtakandidottir@example.net",
				},
			},
			ResentCc: []*mail.Address{
				{
					Name:    "Dan Viðtakandison",
					Address: "dan.vidtakandison@example.net",
				},
				{
					Name:    "Eve Viðtakandidóttir",
					Address: "eve.vidtakandidottir@example.net",
				},
			},
			ResentBcc: []*mail.Address{
				{
					Name:    "Frank Viðtakandison",
					Address: "frank.vidtakandison@example.net",
				},
				{
					Name:    "Grace Viðtakandidóttir",
					Address: "grace.vidtakandidottir@example.net",
				},
			},
			ResentMessageID: "Message-Id-1@example.net",
			ContentType: ContentTypeHeader{
				ContentType: "text/plain",
				Params: map[string]string{
					"charset": "iso-8859-1",
				},
			},
			ExtraHeaders: map[string][]string{
				"X-Clacks-Overhead": {"GNU Terry Pratchett"},
				"X-Script/function/\t !\"#$%&'()*+,-./;<=>?@[\\]^_`{|}~": {
					"TEST VALUE 1\t !\"#$%&'()*+,-./0123456789:;<=>?@ABCDEFGHIJKLMNOPQRSTUVWXYZ[\\]^_` abcdefghijklmnopqrstuvwxyz{|}~",
					"TEST VALUE 2\t !\"#$%&'()*+,-./0123456789:;<=>?@ABCDEFGHIJKLMNOPQRSTUVWXYZ[\\]^_` abcdefghijklmnopqrstuvwxyz{|}~",
				},
			},
		},
		Text: `Kæmi ný öxi hér, ykist þjófum nú bæði víl og ádrepa.
Svo hölt, yxna kýr þegði jú um dóp í fé á bæ.
Þú dazt á hnéð í vök og yfir blóm sexý pæju.`,
	}

	testEmailFromFile(t, fp, expectedEmail)
}

func TestParseEmailIcelandicPlaintextIso88591OverQuotedprintable(t *testing.T) {
	fp := "tests/test_icelandic_plaintext_iso-8859-1_over_quoted-printable.txt"
	tz, _ := time.LoadLocation("Atlantic/Reykjavik")
	expectedDate, _ := time.Parse(
		time.RFC1123Z+" (MST)",
		time.Date(2019, time.April, 1, 7, 55, 0, 0, tz).Format(time.RFC1123Z+" (MST)"))
	expectedEmail := Email{
		Headers: Headers{
			Date:    expectedDate,
			Subject: "📧 Test Íslenskt pangram",
			ReplyTo: []*mail.Address{
				{
					Name:    "Alice Sendandidóttir",
					Address: "alice.sendandidottir@example.net",
				},
			},
			Sender: &mail.Address{
				Name:    "Alice Sendandidóttir",
				Address: "alice.sendandidottir@example.com",
			},
			From: []*mail.Address{
				{
					Name:    "Alice Sendandidóttir",
					Address: "alice.sendandidottir@example.com",
				},
				{
					Name:    "Alice Sendandidóttir",
					Address: "alice.sendandidottir@example.net",
				},
			},
			To: []*mail.Address{
				{
					Name:    "Bob Viðtakandison",
					Address: "bob.didtakandison@example.com",
				},
				{
					Name:    "Carol Viðtakandidóttir",
					Address: "carol.didtakandidottir@example.com",
				},
			},
			Cc: []*mail.Address{
				{
					Name:    "Dan Viðtakandison",
					Address: "dan.vidtakandison@example.com",
				},
				{
					Name:    "Eve Viðtakandidóttir",
					Address: "eve.vidtakandidottir@example.com",
				},
			},
			Bcc: []*mail.Address{
				{
					Name:    "Frank Viðtakandison",
					Address: "frank.vidtakandison@example.com",
				},
				{
					Name:    "Grace Viðtakandidóttir",
					Address: "grace.vidtakandidottir@example.com",
				},
			},
			MessageID:  "Message-Id-1@example.com",
			InReplyTo:  []MessageId{"Message-Id-0@example.com"},
			References: []MessageId{"Message-Id-0@example.com"},
			Comments:   "Message Header Comment",
			Keywords:   []string{"Keyword 1", "Keyword 2"},
			ResentDate: expectedDate,
			ResentFrom: []*mail.Address{
				{
					Name:    "Alice Sendandidóttir",
					Address: "alice.sendandidottir@example.net",
				},
				{
					Name:    "Alice Sendandidóttir",
					Address: "alice.sendandidottir@example.com",
				},
			},
			ResentSender: &mail.Address{
				Name:    "Alice Sendandidóttir",
				Address: "alice.sendandidottir@example.net",
			},
			ResentTo: []*mail.Address{
				{
					Name:    "Bob Viðtakandison",
					Address: "bob.didtakandison@example.net",
				},
				{
					Name:    "Carol Viðtakandidóttir",
					Address: "carol.didtakandidottir@example.net",
				},
			},
			ResentCc: []*mail.Address{
				{
					Name:    "Dan Viðtakandison",
					Address: "dan.vidtakandison@example.net",
				},
				{
					Name:    "Eve Viðtakandidóttir",
					Address: "eve.vidtakandidottir@example.net",
				},
			},
			ResentBcc: []*mail.Address{
				{
					Name:    "Frank Viðtakandison",
					Address: "frank.vidtakandison@example.net",
				},
				{
					Name:    "Grace Viðtakandidóttir",
					Address: "grace.vidtakandidottir@example.net",
				},
			},
			ResentMessageID: "Message-Id-1@example.net",
			ContentType: ContentTypeHeader{
				ContentType: "text/plain",
				Params: map[string]string{
					"charset": "iso-8859-1",
				},
			},
			ExtraHeaders: map[string][]string{
				"X-Clacks-Overhead": {"GNU Terry Pratchett"},
				"X-Script/function/\t !\"#$%&'()*+,-./;<=>?@[\\]^_`{|}~": {
					"TEST VALUE 1\t !\"#$%&'()*+,-./0123456789:;<=>?@ABCDEFGHIJKLMNOPQRSTUVWXYZ[\\]^_` abcdefghijklmnopqrstuvwxyz{|}~",
					"TEST VALUE 2\t !\"#$%&'()*+,-./0123456789:;<=>?@ABCDEFGHIJKLMNOPQRSTUVWXYZ[\\]^_` abcdefghijklmnopqrstuvwxyz{|}~",
				},
			},
		},
		Text: `Kæmi ný öxi hér, ykist þjófum nú bæði víl og ádrepa.
Svo hölt, yxna kýr þegði jú um dóp í fé á bæ.
Þú dazt á hnéð í vök og yfir blóm sexý pæju.`,
	}

	testEmailFromFile(t, fp, expectedEmail)
}

func TestParseEmailIcelandicMultipartRelatedUtf8OverBase64(t *testing.T) {
	fp := "tests/test_icelandic_multipart_related_utf-8_over_base64.txt"
	tz, _ := time.LoadLocation("Atlantic/Reykjavik")
	expectedDate, _ := time.Parse(
		time.RFC1123Z+" (MST)",
		time.Date(2019, time.April, 1, 7, 55, 0, 0, tz).Format(time.RFC1123Z+" (MST)"))
	expectedEmail := Email{
		Headers: Headers{
			Date:    expectedDate,
			Subject: "📧 Test Íslenskt pangram",
			ReplyTo: []*mail.Address{
				{
					Name:    "Alice Sendandidóttir",
					Address: "alice.sendandidottir@example.net",
				},
			},
			Sender: &mail.Address{
				Name:    "Alice Sendandidóttir",
				Address: "alice.sendandidottir@example.com",
			},
			From: []*mail.Address{
				{
					Name:    "Alice Sendandidóttir",
					Address: "alice.sendandidottir@example.com",
				},
				{
					Name:    "Alice Sendandidóttir",
					Address: "alice.sendandidottir@example.net",
				},
			},
			To: []*mail.Address{
				{
					Name:    "Bob Viðtakandison",
					Address: "bob.didtakandison@example.com",
				},
				{
					Name:    "Carol Viðtakandidóttir",
					Address: "carol.didtakandidottir@example.com",
				},
			},
			Cc: []*mail.Address{
				{
					Name:    "Dan Viðtakandison",
					Address: "dan.vidtakandison@example.com",
				},
				{
					Name:    "Eve Viðtakandidóttir",
					Address: "eve.vidtakandidottir@example.com",
				},
			},
			Bcc: []*mail.Address{
				{
					Name:    "Frank Viðtakandison",
					Address: "frank.vidtakandison@example.com",
				},
				{
					Name:    "Grace Viðtakandidóttir",
					Address: "grace.vidtakandidottir@example.com",
				},
			},
			MessageID:  "Message-Id-1@example.com",
			InReplyTo:  []MessageId{"Message-Id-0@example.com"},
			References: []MessageId{"Message-Id-0@example.com"},
			Comments:   "Message Header Comment",
			Keywords:   []string{"Keyword 1", "Keyword 2"},
			ResentDate: expectedDate,
			ResentFrom: []*mail.Address{
				{
					Name:    "Alice Sendandidóttir",
					Address: "alice.sendandidottir@example.net",
				},
				{
					Name:    "Alice Sendandidóttir",
					Address: "alice.sendandidottir@example.com",
				},
			},
			ResentSender: &mail.Address{
				Name:    "Alice Sendandidóttir",
				Address: "alice.sendandidottir@example.net",
			},
			ResentTo: []*mail.Address{
				{
					Name:    "Bob Viðtakandison",
					Address: "bob.didtakandison@example.net",
				},
				{
					Name:    "Carol Viðtakandidóttir",
					Address: "carol.didtakandidottir@example.net",
				},
			},
			ResentCc: []*mail.Address{
				{
					Name:    "Dan Viðtakandison",
					Address: "dan.vidtakandison@example.net",
				},
				{
					Name:    "Eve Viðtakandidóttir",
					Address: "eve.vidtakandidottir@example.net",
				},
			},
			ResentBcc: []*mail.Address{
				{
					Name:    "Frank Viðtakandison",
					Address: "frank.vidtakandison@example.net",
				},
				{
					Name:    "Grace Viðtakandidóttir",
					Address: "grace.vidtakandidottir@example.net",
				},
			},
			ResentMessageID: "Message-Id-1@example.net",
			ContentType: ContentTypeHeader{
				ContentType: "multipart/related",
				Params: map[string]string{
					"boundary": "RelatedBoundaryString",
					"charset":  "utf-8",
				},
			},
			ExtraHeaders: map[string][]string{
				"X-Clacks-Overhead": {"GNU Terry Pratchett"},
				"X-Script/function/\t !\"#$%&'()*+,-./;<=>?@[\\]^_`{|}~": {
					"TEST VALUE 1\t !\"#$%&'()*+,-./0123456789:;<=>?@ABCDEFGHIJKLMNOPQRSTUVWXYZ[\\]^_` abcdefghijklmnopqrstuvwxyz{|}~",
					"TEST VALUE 2\t !\"#$%&'()*+,-./0123456789:;<=>?@ABCDEFGHIJKLMNOPQRSTUVWXYZ[\\]^_` abcdefghijklmnopqrstuvwxyz{|}~",
				},
			},
		},
		Text: `Kæmi ný öxi hér, ykist þjófum nú bæði víl og ádrepa.
Svo hölt, yxna kýr þegði jú um dóp í fé á bæ.
Þú dazt á hnéð í vök og yfir blóm sexý pæju.`,
		EnrichedText: `<bold>Kæmi ný öxi hér, ykist þjófum nú bæði víl og ádrepa.</bold>
<italic>Svo hölt, yxna kýr þegði jú um dóp í fé á bæ.</italic>
<fixed>Þú dazt á hnéð í vök og yfir blóm sexý pæju.</fixed>`,
		HTML: `<html>
<div dir="ltr">
<p>Kæmi ný öxi hér, ykist þjófum nú bæði víl og ádrepa.</p>
<p>Svo hölt, yxna kýr þegði jú um dóp í fé á bæ.</p>
<p>Þú dazt á hnéð í vök og yfir blóm sexý pæju.</p>
</div>
</html>`,
		InlineFiles: []InlineFile{
			{
				ContentID: "inline-jpg-image.jpg@example.com",
				ContentType: ContentTypeHeader{
					ContentType: "image/jpeg",
					Params: map[string]string{
						"name": "inline-jpg-image-name.jpg",
					},
				},
				ContentDisposition: ContentDispositionHeader{
					ContentDisposition: inline,
					Params: map[string]string{
						"filename": "inline-jpg-image-filename.jpg",
					},
				},
				Data: []byte{255, 216, 255, 219, 0, 67, 0, 3, 2, 2, 2, 2, 2, 3, 2, 2, 2, 3, 3, 3, 3, 4,
					6, 4, 4, 4, 4, 4, 8, 6, 6, 5, 6, 9, 8, 10, 10, 9, 8, 9, 9, 10, 12, 15, 12, 10, 11, 14, 11, 9, 9,
					13, 17, 13, 14, 15, 16, 16, 17, 16, 10, 12, 18, 19, 18, 16, 19, 15, 16, 16, 16, 255, 201, 0, 11, 8,
					0, 1, 0, 1, 1, 1, 17, 0, 255, 204, 0, 6, 0, 16, 16, 5, 255, 218, 0, 8, 1, 1, 0, 0, 63, 0, 210, 207,
					32, 255, 217},
			},
		},
	}

	testEmailFromFile(t, fp, expectedEmail)
}
func TestParseEmailIcelandicMultipartRelatedUtf8OverQuotedprintable(t *testing.T) {
	fp := "tests/test_icelandic_multipart_related_utf-8_over_quoted-printable.txt"
	tz, _ := time.LoadLocation("Atlantic/Reykjavik")
	expectedDate, _ := time.Parse(
		time.RFC1123Z+" (MST)",
		time.Date(2019, time.April, 1, 7, 55, 0, 0, tz).Format(time.RFC1123Z+" (MST)"))
	expectedEmail := Email{
		Headers: Headers{
			Date:    expectedDate,
			Subject: "📧 Test Íslenskt pangram",
			ReplyTo: []*mail.Address{
				{
					Name:    "Alice Sendandidóttir",
					Address: "alice.sendandidottir@example.net",
				},
			},
			Sender: &mail.Address{
				Name:    "Alice Sendandidóttir",
				Address: "alice.sendandidottir@example.com",
			},
			From: []*mail.Address{
				{
					Name:    "Alice Sendandidóttir",
					Address: "alice.sendandidottir@example.com",
				},
				{
					Name:    "Alice Sendandidóttir",
					Address: "alice.sendandidottir@example.net",
				},
			},
			To: []*mail.Address{
				{
					Name:    "Bob Viðtakandison",
					Address: "bob.didtakandison@example.com",
				},
				{
					Name:    "Carol Viðtakandidóttir",
					Address: "carol.didtakandidottir@example.com",
				},
			},
			Cc: []*mail.Address{
				{
					Name:    "Dan Viðtakandison",
					Address: "dan.vidtakandison@example.com",
				},
				{
					Name:    "Eve Viðtakandidóttir",
					Address: "eve.vidtakandidottir@example.com",
				},
			},
			Bcc: []*mail.Address{
				{
					Name:    "Frank Viðtakandison",
					Address: "frank.vidtakandison@example.com",
				},
				{
					Name:    "Grace Viðtakandidóttir",
					Address: "grace.vidtakandidottir@example.com",
				},
			},
			MessageID:  "Message-Id-1@example.com",
			InReplyTo:  []MessageId{"Message-Id-0@example.com"},
			References: []MessageId{"Message-Id-0@example.com"},
			Comments:   "Message Header Comment",
			Keywords:   []string{"Keyword 1", "Keyword 2"},
			ResentDate: expectedDate,
			ResentFrom: []*mail.Address{
				{
					Name:    "Alice Sendandidóttir",
					Address: "alice.sendandidottir@example.net",
				},
				{
					Name:    "Alice Sendandidóttir",
					Address: "alice.sendandidottir@example.com",
				},
			},
			ResentSender: &mail.Address{
				Name:    "Alice Sendandidóttir",
				Address: "alice.sendandidottir@example.net",
			},
			ResentTo: []*mail.Address{
				{
					Name:    "Bob Viðtakandison",
					Address: "bob.didtakandison@example.net",
				},
				{
					Name:    "Carol Viðtakandidóttir",
					Address: "carol.didtakandidottir@example.net",
				},
			},
			ResentCc: []*mail.Address{
				{
					Name:    "Dan Viðtakandison",
					Address: "dan.vidtakandison@example.net",
				},
				{
					Name:    "Eve Viðtakandidóttir",
					Address: "eve.vidtakandidottir@example.net",
				},
			},
			ResentBcc: []*mail.Address{
				{
					Name:    "Frank Viðtakandison",
					Address: "frank.vidtakandison@example.net",
				},
				{
					Name:    "Grace Viðtakandidóttir",
					Address: "grace.vidtakandidottir@example.net",
				},
			},
			ResentMessageID: "Message-Id-1@example.net",
			ContentType: ContentTypeHeader{
				ContentType: "multipart/related",
				Params: map[string]string{
					"boundary": "RelatedBoundaryString",
					"charset":  "utf-8",
				},
			},
			ExtraHeaders: map[string][]string{
				"X-Clacks-Overhead": {"GNU Terry Pratchett"},
				"X-Script/function/\t !\"#$%&'()*+,-./;<=>?@[\\]^_`{|}~": {
					"TEST VALUE 1\t !\"#$%&'()*+,-./0123456789:;<=>?@ABCDEFGHIJKLMNOPQRSTUVWXYZ[\\]^_` abcdefghijklmnopqrstuvwxyz{|}~",
					"TEST VALUE 2\t !\"#$%&'()*+,-./0123456789:;<=>?@ABCDEFGHIJKLMNOPQRSTUVWXYZ[\\]^_` abcdefghijklmnopqrstuvwxyz{|}~",
				},
			},
		},
		Text: `Kæmi ný öxi hér, ykist þjófum nú bæði víl og ádrepa.
Svo hölt, yxna kýr þegði jú um dóp í fé á bæ.
Þú dazt á hnéð í vök og yfir blóm sexý pæju.`,
		EnrichedText: `<bold>Kæmi ný öxi hér, ykist þjófum nú bæði víl og ádrepa.</bold>
<italic>Svo hölt, yxna kýr þegði jú um dóp í fé á bæ.</italic>
<fixed>Þú dazt á hnéð í vök og yfir blóm sexý pæju.</fixed>`,
		HTML: `<html>
<div dir="ltr">
<p>Kæmi ný öxi hér, ykist þjófum nú bæði víl og ádrepa.</p>
<p>Svo hölt, yxna kýr þegði jú um dóp í fé á bæ.</p>
<p>Þú dazt á hnéð í vök og yfir blóm sexý pæju.</p>
</div>
</html>`,
		InlineFiles: []InlineFile{
			{
				ContentID: "inline-jpg-image.jpg@example.com",
				ContentType: ContentTypeHeader{
					ContentType: "image/jpeg",
					Params: map[string]string{
						"name": "inline-jpg-image-name.jpg",
					},
				},
				ContentDisposition: ContentDispositionHeader{
					ContentDisposition: inline,
					Params: map[string]string{
						"filename": "inline-jpg-image-filename.jpg",
					},
				},
				Data: []byte{255, 216, 255, 219, 0, 67, 0, 3, 2, 2, 2, 2, 2, 3, 2, 2, 2, 3, 3, 3, 3, 4,
					6, 4, 4, 4, 4, 4, 8, 6, 6, 5, 6, 9, 8, 10, 10, 9, 8, 9, 9, 10, 12, 15, 12, 10, 11, 14, 11, 9, 9,
					13, 17, 13, 14, 15, 16, 16, 17, 16, 10, 12, 18, 19, 18, 16, 19, 15, 16, 16, 16, 255, 201, 0, 11, 8,
					0, 1, 0, 1, 1, 1, 17, 0, 255, 204, 0, 6, 0, 16, 16, 5, 255, 218, 0, 8, 1, 1, 0, 0, 63, 0, 210, 207,
					32, 255, 217},
			},
		},
	}

	testEmailFromFile(t, fp, expectedEmail)
}
func TestParseEmailIcelandicMultipartRelatedIso88591OverBase64(t *testing.T) {
	fp := "tests/test_icelandic_multipart_related_iso-8859-1_over_base64.txt"
	tz, _ := time.LoadLocation("Atlantic/Reykjavik")
	expectedDate, _ := time.Parse(
		time.RFC1123Z+" (MST)",
		time.Date(2019, time.April, 1, 7, 55, 0, 0, tz).Format(time.RFC1123Z+" (MST)"))
	expectedEmail := Email{
		Headers: Headers{
			Date:    expectedDate,
			Subject: "📧 Test Íslenskt pangram",
			ReplyTo: []*mail.Address{
				{
					Name:    "Alice Sendandidóttir",
					Address: "alice.sendandidottir@example.net",
				},
			},
			Sender: &mail.Address{
				Name:    "Alice Sendandidóttir",
				Address: "alice.sendandidottir@example.com",
			},
			From: []*mail.Address{
				{
					Name:    "Alice Sendandidóttir",
					Address: "alice.sendandidottir@example.com",
				},
				{
					Name:    "Alice Sendandidóttir",
					Address: "alice.sendandidottir@example.net",
				},
			},
			To: []*mail.Address{
				{
					Name:    "Bob Viðtakandison",
					Address: "bob.didtakandison@example.com",
				},
				{
					Name:    "Carol Viðtakandidóttir",
					Address: "carol.didtakandidottir@example.com",
				},
			},
			Cc: []*mail.Address{
				{
					Name:    "Dan Viðtakandison",
					Address: "dan.vidtakandison@example.com",
				},
				{
					Name:    "Eve Viðtakandidóttir",
					Address: "eve.vidtakandidottir@example.com",
				},
			},
			Bcc: []*mail.Address{
				{
					Name:    "Frank Viðtakandison",
					Address: "frank.vidtakandison@example.com",
				},
				{
					Name:    "Grace Viðtakandidóttir",
					Address: "grace.vidtakandidottir@example.com",
				},
			},
			MessageID:  "Message-Id-1@example.com",
			InReplyTo:  []MessageId{"Message-Id-0@example.com"},
			References: []MessageId{"Message-Id-0@example.com"},
			Comments:   "Message Header Comment",
			Keywords:   []string{"Keyword 1", "Keyword 2"},
			ResentDate: expectedDate,
			ResentFrom: []*mail.Address{
				{
					Name:    "Alice Sendandidóttir",
					Address: "alice.sendandidottir@example.net",
				},
				{
					Name:    "Alice Sendandidóttir",
					Address: "alice.sendandidottir@example.com",
				},
			},
			ResentSender: &mail.Address{
				Name:    "Alice Sendandidóttir",
				Address: "alice.sendandidottir@example.net",
			},
			ResentTo: []*mail.Address{
				{
					Name:    "Bob Viðtakandison",
					Address: "bob.didtakandison@example.net",
				},
				{
					Name:    "Carol Viðtakandidóttir",
					Address: "carol.didtakandidottir@example.net",
				},
			},
			ResentCc: []*mail.Address{
				{
					Name:    "Dan Viðtakandison",
					Address: "dan.vidtakandison@example.net",
				},
				{
					Name:    "Eve Viðtakandidóttir",
					Address: "eve.vidtakandidottir@example.net",
				},
			},
			ResentBcc: []*mail.Address{
				{
					Name:    "Frank Viðtakandison",
					Address: "frank.vidtakandison@example.net",
				},
				{
					Name:    "Grace Viðtakandidóttir",
					Address: "grace.vidtakandidottir@example.net",
				},
			},
			ResentMessageID: "Message-Id-1@example.net",
			ContentType: ContentTypeHeader{
				ContentType: "multipart/related",
				Params: map[string]string{
					"boundary": "RelatedBoundaryString",
					"charset":  "iso-8859-1",
				},
			},
			ExtraHeaders: map[string][]string{
				"X-Clacks-Overhead": {"GNU Terry Pratchett"},
				"X-Script/function/\t !\"#$%&'()*+,-./;<=>?@[\\]^_`{|}~": {
					"TEST VALUE 1\t !\"#$%&'()*+,-./0123456789:;<=>?@ABCDEFGHIJKLMNOPQRSTUVWXYZ[\\]^_` abcdefghijklmnopqrstuvwxyz{|}~",
					"TEST VALUE 2\t !\"#$%&'()*+,-./0123456789:;<=>?@ABCDEFGHIJKLMNOPQRSTUVWXYZ[\\]^_` abcdefghijklmnopqrstuvwxyz{|}~",
				},
			},
		},
		Text: `Kæmi ný öxi hér, ykist þjófum nú bæði víl og ádrepa.
Svo hölt, yxna kýr þegði jú um dóp í fé á bæ.
Þú dazt á hnéð í vök og yfir blóm sexý pæju.`,
		EnrichedText: `<bold>Kæmi ný öxi hér, ykist þjófum nú bæði víl og ádrepa.</bold>
<italic>Svo hölt, yxna kýr þegði jú um dóp í fé á bæ.</italic>
<fixed>Þú dazt á hnéð í vök og yfir blóm sexý pæju.</fixed>`,
		HTML: `<html>
<div dir="ltr">
<p>Kæmi ný öxi hér, ykist þjófum nú bæði víl og ádrepa.</p>
<p>Svo hölt, yxna kýr þegði jú um dóp í fé á bæ.</p>
<p>Þú dazt á hnéð í vök og yfir blóm sexý pæju.</p>
</div>
</html>`,
		InlineFiles: []InlineFile{
			{
				ContentID: "inline-jpg-image.jpg@example.com",
				ContentType: ContentTypeHeader{
					ContentType: "image/jpeg",
					Params: map[string]string{
						"name": "inline-jpg-image-name.jpg",
					},
				},
				ContentDisposition: ContentDispositionHeader{
					ContentDisposition: inline,
					Params: map[string]string{
						"filename": "inline-jpg-image-filename.jpg",
					},
				},
				Data: []byte{255, 216, 255, 219, 0, 67, 0, 3, 2, 2, 2, 2, 2, 3, 2, 2, 2, 3, 3, 3, 3, 4,
					6, 4, 4, 4, 4, 4, 8, 6, 6, 5, 6, 9, 8, 10, 10, 9, 8, 9, 9, 10, 12, 15, 12, 10, 11, 14, 11, 9, 9,
					13, 17, 13, 14, 15, 16, 16, 17, 16, 10, 12, 18, 19, 18, 16, 19, 15, 16, 16, 16, 255, 201, 0, 11, 8,
					0, 1, 0, 1, 1, 1, 17, 0, 255, 204, 0, 6, 0, 16, 16, 5, 255, 218, 0, 8, 1, 1, 0, 0, 63, 0, 210, 207,
					32, 255, 217},
			},
		},
	}

	testEmailFromFile(t, fp, expectedEmail)
}
func TestParseEmailIcelandicMultipartRelatedIso88591OverQuotedprintable(t *testing.T) {
	fp := "tests/test_icelandic_multipart_related_iso-8859-1_over_quoted-printable.txt"
	tz, _ := time.LoadLocation("Atlantic/Reykjavik")
	expectedDate, _ := time.Parse(
		time.RFC1123Z+" (MST)",
		time.Date(2019, time.April, 1, 7, 55, 0, 0, tz).Format(time.RFC1123Z+" (MST)"))
	expectedEmail := Email{
		Headers: Headers{
			Date:    expectedDate,
			Subject: "📧 Test Íslenskt pangram",
			ReplyTo: []*mail.Address{
				{
					Name:    "Alice Sendandidóttir",
					Address: "alice.sendandidottir@example.net",
				},
			},
			Sender: &mail.Address{
				Name:    "Alice Sendandidóttir",
				Address: "alice.sendandidottir@example.com",
			},
			From: []*mail.Address{
				{
					Name:    "Alice Sendandidóttir",
					Address: "alice.sendandidottir@example.com",
				},
				{
					Name:    "Alice Sendandidóttir",
					Address: "alice.sendandidottir@example.net",
				},
			},
			To: []*mail.Address{
				{
					Name:    "Bob Viðtakandison",
					Address: "bob.didtakandison@example.com",
				},
				{
					Name:    "Carol Viðtakandidóttir",
					Address: "carol.didtakandidottir@example.com",
				},
			},
			Cc: []*mail.Address{
				{
					Name:    "Dan Viðtakandison",
					Address: "dan.vidtakandison@example.com",
				},
				{
					Name:    "Eve Viðtakandidóttir",
					Address: "eve.vidtakandidottir@example.com",
				},
			},
			Bcc: []*mail.Address{
				{
					Name:    "Frank Viðtakandison",
					Address: "frank.vidtakandison@example.com",
				},
				{
					Name:    "Grace Viðtakandidóttir",
					Address: "grace.vidtakandidottir@example.com",
				},
			},
			MessageID:  "Message-Id-1@example.com",
			InReplyTo:  []MessageId{"Message-Id-0@example.com"},
			References: []MessageId{"Message-Id-0@example.com"},
			Comments:   "Message Header Comment",
			Keywords:   []string{"Keyword 1", "Keyword 2"},
			ResentDate: expectedDate,
			ResentFrom: []*mail.Address{
				{
					Name:    "Alice Sendandidóttir",
					Address: "alice.sendandidottir@example.net",
				},
				{
					Name:    "Alice Sendandidóttir",
					Address: "alice.sendandidottir@example.com",
				},
			},
			ResentSender: &mail.Address{
				Name:    "Alice Sendandidóttir",
				Address: "alice.sendandidottir@example.net",
			},
			ResentTo: []*mail.Address{
				{
					Name:    "Bob Viðtakandison",
					Address: "bob.didtakandison@example.net",
				},
				{
					Name:    "Carol Viðtakandidóttir",
					Address: "carol.didtakandidottir@example.net",
				},
			},
			ResentCc: []*mail.Address{
				{
					Name:    "Dan Viðtakandison",
					Address: "dan.vidtakandison@example.net",
				},
				{
					Name:    "Eve Viðtakandidóttir",
					Address: "eve.vidtakandidottir@example.net",
				},
			},
			ResentBcc: []*mail.Address{
				{
					Name:    "Frank Viðtakandison",
					Address: "frank.vidtakandison@example.net",
				},
				{
					Name:    "Grace Viðtakandidóttir",
					Address: "grace.vidtakandidottir@example.net",
				},
			},
			ResentMessageID: "Message-Id-1@example.net",
			ContentType: ContentTypeHeader{
				ContentType: "multipart/related",
				Params: map[string]string{
					"boundary": "RelatedBoundaryString",
					"charset":  "iso-8859-1",
				},
			},
			ExtraHeaders: map[string][]string{
				"X-Clacks-Overhead": {"GNU Terry Pratchett"},
				"X-Script/function/\t !\"#$%&'()*+,-./;<=>?@[\\]^_`{|}~": {
					"TEST VALUE 1\t !\"#$%&'()*+,-./0123456789:;<=>?@ABCDEFGHIJKLMNOPQRSTUVWXYZ[\\]^_` abcdefghijklmnopqrstuvwxyz{|}~",
					"TEST VALUE 2\t !\"#$%&'()*+,-./0123456789:;<=>?@ABCDEFGHIJKLMNOPQRSTUVWXYZ[\\]^_` abcdefghijklmnopqrstuvwxyz{|}~",
				},
			},
		},
		Text: `Kæmi ný öxi hér, ykist þjófum nú bæði víl og ádrepa.
Svo hölt, yxna kýr þegði jú um dóp í fé á bæ.
Þú dazt á hnéð í vök og yfir blóm sexý pæju.`,
		EnrichedText: `<bold>Kæmi ný öxi hér, ykist þjófum nú bæði víl og ádrepa.</bold>
<italic>Svo hölt, yxna kýr þegði jú um dóp í fé á bæ.</italic>
<fixed>Þú dazt á hnéð í vök og yfir blóm sexý pæju.</fixed>`,
		HTML: `<html>
<div dir="ltr">
<p>Kæmi ný öxi hér, ykist þjófum nú bæði víl og ádrepa.</p>
<p>Svo hölt, yxna kýr þegði jú um dóp í fé á bæ.</p>
<p>Þú dazt á hnéð í vök og yfir blóm sexý pæju.</p>
</div>
</html>`,
		InlineFiles: []InlineFile{
			{
				ContentID: "inline-jpg-image.jpg@example.com",
				ContentType: ContentTypeHeader{
					ContentType: "image/jpeg",
					Params: map[string]string{
						"name": "inline-jpg-image-name.jpg",
					},
				},
				ContentDisposition: ContentDispositionHeader{
					ContentDisposition: inline,
					Params: map[string]string{
						"filename": "inline-jpg-image-filename.jpg",
					},
				},
				Data: []byte{255, 216, 255, 219, 0, 67, 0, 3, 2, 2, 2, 2, 2, 3, 2, 2, 2, 3, 3, 3, 3, 4,
					6, 4, 4, 4, 4, 4, 8, 6, 6, 5, 6, 9, 8, 10, 10, 9, 8, 9, 9, 10, 12, 15, 12, 10, 11, 14, 11, 9, 9,
					13, 17, 13, 14, 15, 16, 16, 17, 16, 10, 12, 18, 19, 18, 16, 19, 15, 16, 16, 16, 255, 201, 0, 11, 8,
					0, 1, 0, 1, 1, 1, 17, 0, 255, 204, 0, 6, 0, 16, 16, 5, 255, 218, 0, 8, 1, 1, 0, 0, 63, 0, 210, 207,
					32, 255, 217},
			},
		},
	}

	testEmailFromFile(t, fp, expectedEmail)
}
func TestParseEmailIcelandicMultipartMixedUtf8OverBase64(t *testing.T) {
	fp := "tests/test_icelandic_multipart_mixed_utf-8_over_base64.txt"
	tz, _ := time.LoadLocation("Atlantic/Reykjavik")
	expectedDate, _ := time.Parse(
		time.RFC1123Z+" (MST)",
		time.Date(2019, time.April, 1, 7, 55, 0, 0, tz).Format(time.RFC1123Z+" (MST)"))
	expectedEmail := Email{
		Headers: Headers{
			Date:    expectedDate,
			Subject: "📧 Test Íslenskt pangram",
			ReplyTo: []*mail.Address{
				{
					Name:    "Alice Sendandidóttir",
					Address: "alice.sendandidottir@example.net",
				},
			},
			Sender: &mail.Address{
				Name:    "Alice Sendandidóttir",
				Address: "alice.sendandidottir@example.com",
			},
			From: []*mail.Address{
				{
					Name:    "Alice Sendandidóttir",
					Address: "alice.sendandidottir@example.com",
				},
				{
					Name:    "Alice Sendandidóttir",
					Address: "alice.sendandidottir@example.net",
				},
			},
			To: []*mail.Address{
				{
					Name:    "Bob Viðtakandison",
					Address: "bob.didtakandison@example.com",
				},
				{
					Name:    "Carol Viðtakandidóttir",
					Address: "carol.didtakandidottir@example.com",
				},
			},
			Cc: []*mail.Address{
				{
					Name:    "Dan Viðtakandison",
					Address: "dan.vidtakandison@example.com",
				},
				{
					Name:    "Eve Viðtakandidóttir",
					Address: "eve.vidtakandidottir@example.com",
				},
			},
			Bcc: []*mail.Address{
				{
					Name:    "Frank Viðtakandison",
					Address: "frank.vidtakandison@example.com",
				},
				{
					Name:    "Grace Viðtakandidóttir",
					Address: "grace.vidtakandidottir@example.com",
				},
			},
			MessageID:  "Message-Id-1@example.com",
			InReplyTo:  []MessageId{"Message-Id-0@example.com"},
			References: []MessageId{"Message-Id-0@example.com"},
			Comments:   "Message Header Comment",
			Keywords:   []string{"Keyword 1", "Keyword 2"},
			ResentDate: expectedDate,
			ResentFrom: []*mail.Address{
				{
					Name:    "Alice Sendandidóttir",
					Address: "alice.sendandidottir@example.net",
				},
				{
					Name:    "Alice Sendandidóttir",
					Address: "alice.sendandidottir@example.com",
				},
			},
			ResentSender: &mail.Address{
				Name:    "Alice Sendandidóttir",
				Address: "alice.sendandidottir@example.net",
			},
			ResentTo: []*mail.Address{
				{
					Name:    "Bob Viðtakandison",
					Address: "bob.didtakandison@example.net",
				},
				{
					Name:    "Carol Viðtakandidóttir",
					Address: "carol.didtakandidottir@example.net",
				},
			},
			ResentCc: []*mail.Address{
				{
					Name:    "Dan Viðtakandison",
					Address: "dan.vidtakandison@example.net",
				},
				{
					Name:    "Eve Viðtakandidóttir",
					Address: "eve.vidtakandidottir@example.net",
				},
			},
			ResentBcc: []*mail.Address{
				{
					Name:    "Frank Viðtakandison",
					Address: "frank.vidtakandison@example.net",
				},
				{
					Name:    "Grace Viðtakandidóttir",
					Address: "grace.vidtakandidottir@example.net",
				},
			},
			ResentMessageID: "Message-Id-1@example.net",
			ContentType: ContentTypeHeader{
				ContentType: "multipart/mixed",
				Params: map[string]string{
					"boundary": "MixedBoundaryString",
					"charset":  "utf-8",
				},
			},
			ExtraHeaders: map[string][]string{
				"X-Clacks-Overhead": {"GNU Terry Pratchett"},
				"X-Script/function/\t !\"#$%&'()*+,-./;<=>?@[\\]^_`{|}~": {
					"TEST VALUE 1\t !\"#$%&'()*+,-./0123456789:;<=>?@ABCDEFGHIJKLMNOPQRSTUVWXYZ[\\]^_` abcdefghijklmnopqrstuvwxyz{|}~",
					"TEST VALUE 2\t !\"#$%&'()*+,-./0123456789:;<=>?@ABCDEFGHIJKLMNOPQRSTUVWXYZ[\\]^_` abcdefghijklmnopqrstuvwxyz{|}~",
				},
			},
		},
		Text: `Kæmi ný öxi hér, ykist þjófum nú bæði víl og ádrepa.
Svo hölt, yxna kýr þegði jú um dóp í fé á bæ.
Þú dazt á hnéð í vök og yfir blóm sexý pæju.`,
		EnrichedText: `<bold>Kæmi ný öxi hér, ykist þjófum nú bæði víl og ádrepa.</bold>
<italic>Svo hölt, yxna kýr þegði jú um dóp í fé á bæ.</italic>
<fixed>Þú dazt á hnéð í vök og yfir blóm sexý pæju.</fixed>`,
		HTML: `<html>
<div dir="ltr">
<p>Kæmi ný öxi hér, ykist þjófum nú bæði víl og ádrepa.</p>
<p>Svo hölt, yxna kýr þegði jú um dóp í fé á bæ.</p>
<p>Þú dazt á hnéð í vök og yfir blóm sexý pæju.</p>
</div>
</html>`,
		InlineFiles: []InlineFile{
			{
				ContentType: ContentTypeHeader{
					ContentType: "image/jpeg",
					Params: map[string]string{
						"name": "inline-jpg-image-without-disposition.jpg",
					},
				},
				ContentDisposition: ContentDispositionHeader{
					ContentDisposition: "",
					Params:             map[string]string(nil),
				},
				Data: []byte{255, 216, 255, 219, 0, 67, 0, 3, 2, 2, 2, 2, 2, 3, 2, 2, 2, 3, 3, 3, 3, 4,
					6, 4, 4, 4, 4, 4, 8, 6, 6, 5, 6, 9, 8, 10, 10, 9, 8, 9, 9, 10, 12, 15, 12, 10, 11, 14, 11, 9, 9,
					13, 17, 13, 14, 15, 16, 16, 17, 16, 10, 12, 18, 19, 18, 16, 19, 15, 16, 16, 16, 255, 201, 0, 11, 8,
					0, 1, 0, 1, 1, 1, 17, 0, 255, 204, 0, 6, 0, 16, 16, 5, 255, 218, 0, 8, 1, 1, 0, 0, 63, 0, 210, 207,
					32, 255, 217},
			},
			{
				ContentID: "inline-jpg-image.jpg@example.com",
				ContentType: ContentTypeHeader{
					ContentType: "image/jpeg",
					Params: map[string]string{
						"name": "inline-jpg-image-name.jpg",
					},
				},
				ContentDisposition: ContentDispositionHeader{
					ContentDisposition: inline,
					Params: map[string]string{
						"filename": "inline-jpg-image-filename.jpg",
					},
				},
				Data: []byte{255, 216, 255, 219, 0, 67, 0, 3, 2, 2, 2, 2, 2, 3, 2, 2, 2, 3, 3, 3, 3, 4,
					6, 4, 4, 4, 4, 4, 8, 6, 6, 5, 6, 9, 8, 10, 10, 9, 8, 9, 9, 10, 12, 15, 12, 10, 11, 14, 11, 9, 9,
					13, 17, 13, 14, 15, 16, 16, 17, 16, 10, 12, 18, 19, 18, 16, 19, 15, 16, 16, 16, 255, 201, 0, 11, 8,
					0, 1, 0, 1, 1, 1, 17, 0, 255, 204, 0, 6, 0, 16, 16, 5, 255, 218, 0, 8, 1, 1, 0, 0, 63, 0, 210, 207,
					32, 255, 217},
			},
		},
		AttachedFiles: []AttachedFile{
			{
				ContentType: ContentTypeHeader{
					ContentType: "application/pdf",
					Params: map[string]string{
						"name": "attached-pdf-name.pdf",
					},
				},
				ContentDisposition: ContentDispositionHeader{
					ContentDisposition: attachment,
					Params: map[string]string{
						"filename": "attached-pdf-filename.pdf",
					},
				},
				Data: []byte{37, 80, 68, 70, 45, 49, 46, 13, 116, 114, 97, 105, 108, 101, 114, 60, 60,
					47, 82, 111, 111, 116, 60, 60, 47, 80, 97, 103, 101, 115, 60, 60, 47, 75, 105, 100, 115, 91, 60,
					60, 47, 77, 101, 100, 105, 97, 66, 111, 120, 91, 48, 32, 48, 32, 51, 32, 51, 93, 62, 62, 93, 62,
					62, 62, 62, 62, 62},
			},
			{
				ContentType: ContentTypeHeader{
					ContentType: "application/pdf",
					Params: map[string]string{
						"name": "attached-pdf-without-disposition.pdf",
					},
				},
				ContentDisposition: ContentDispositionHeader{
					ContentDisposition: "",
					Params:             map[string]string(nil),
				},
				Data: []byte{37, 80, 68, 70, 45, 49, 46, 13, 116, 114, 97, 105, 108, 101, 114, 60, 60,
					47, 82, 111, 111, 116, 60, 60, 47, 80, 97, 103, 101, 115, 60, 60, 47, 75, 105, 100, 115, 91, 60,
					60, 47, 77, 101, 100, 105, 97, 66, 111, 120, 91, 48, 32, 48, 32, 51, 32, 51, 93, 62, 62, 93, 62,
					62, 62, 62, 62, 62},
			},
			{
				ContentType: ContentTypeHeader{
					ContentType: "application/json",
					Params: map[string]string{
						"name": "attached-json-name.json",
					},
				},
				ContentDisposition: ContentDispositionHeader{
					ContentDisposition: attachment,
					Params: map[string]string{
						"filename": "attached-json-filename.json",
					},
				},
				Data: []byte{123, 34, 102, 111, 111, 34, 58, 34, 98, 97, 114, 34, 125},
			},
			{
				ContentType: ContentTypeHeader{
					ContentType: "text/plain",
					Params: map[string]string{
						"name": "attached-text-plain-name.txt",
					},
				},
				ContentDisposition: ContentDispositionHeader{
					ContentDisposition: attachment,
					Params: map[string]string{
						"filename": "attached-text-plain-filename.txt",
					},
				},
				Data: []byte{84, 101, 120, 116, 47, 112, 108, 97, 105, 110, 32, 99, 111, 110, 116, 101, 110, 116, 32,
					97, 115, 32, 97, 110, 32, 97, 116, 116, 97, 99, 104, 101, 100, 32, 46, 116, 120, 116, 32, 102, 105,
					108, 101, 46},
			},
			{
				ContentType: ContentTypeHeader{
					ContentType: "text/html",
					Params: map[string]string{
						"name": "attached-text-html-name.html",
					},
				},
				ContentDisposition: ContentDispositionHeader{
					ContentDisposition: attachment,
					Params: map[string]string{
						"filename": "attached-text-html-filename.html",
					},
				},
				Data: []byte{84, 101, 120, 116, 47, 104, 116, 109, 108, 32, 99, 111, 110, 116, 101, 110, 116, 32, 97,
					115, 32, 97, 110, 32, 97, 116, 116, 97, 99, 104, 101, 100, 32, 46, 104, 116, 109, 108, 32, 102,
					105, 108, 101, 46},
			},
		},
	}

	testEmailFromFile(t, fp, expectedEmail)
}
func TestParseEmailIcelandicMultipartMixedUtf8OverQuotedprintable(t *testing.T) {
	fp := "tests/test_icelandic_multipart_mixed_utf-8_over_quoted-printable.txt"
	tz, _ := time.LoadLocation("Atlantic/Reykjavik")
	expectedDate, _ := time.Parse(
		time.RFC1123Z+" (MST)",
		time.Date(2019, time.April, 1, 7, 55, 0, 0, tz).Format(time.RFC1123Z+" (MST)"))
	expectedEmail := Email{
		Headers: Headers{
			Date:    expectedDate,
			Subject: "📧 Test Íslenskt pangram",
			ReplyTo: []*mail.Address{
				{
					Name:    "Alice Sendandidóttir",
					Address: "alice.sendandidottir@example.net",
				},
			},
			Sender: &mail.Address{
				Name:    "Alice Sendandidóttir",
				Address: "alice.sendandidottir@example.com",
			},
			From: []*mail.Address{
				{
					Name:    "Alice Sendandidóttir",
					Address: "alice.sendandidottir@example.com",
				},
				{
					Name:    "Alice Sendandidóttir",
					Address: "alice.sendandidottir@example.net",
				},
			},
			To: []*mail.Address{
				{
					Name:    "Bob Viðtakandison",
					Address: "bob.didtakandison@example.com",
				},
				{
					Name:    "Carol Viðtakandidóttir",
					Address: "carol.didtakandidottir@example.com",
				},
			},
			Cc: []*mail.Address{
				{
					Name:    "Dan Viðtakandison",
					Address: "dan.vidtakandison@example.com",
				},
				{
					Name:    "Eve Viðtakandidóttir",
					Address: "eve.vidtakandidottir@example.com",
				},
			},
			Bcc: []*mail.Address{
				{
					Name:    "Frank Viðtakandison",
					Address: "frank.vidtakandison@example.com",
				},
				{
					Name:    "Grace Viðtakandidóttir",
					Address: "grace.vidtakandidottir@example.com",
				},
			},
			MessageID:  "Message-Id-1@example.com",
			InReplyTo:  []MessageId{"Message-Id-0@example.com"},
			References: []MessageId{"Message-Id-0@example.com"},
			Comments:   "Message Header Comment",
			Keywords:   []string{"Keyword 1", "Keyword 2"},
			ResentDate: expectedDate,
			ResentFrom: []*mail.Address{
				{
					Name:    "Alice Sendandidóttir",
					Address: "alice.sendandidottir@example.net",
				},
				{
					Name:    "Alice Sendandidóttir",
					Address: "alice.sendandidottir@example.com",
				},
			},
			ResentSender: &mail.Address{
				Name:    "Alice Sendandidóttir",
				Address: "alice.sendandidottir@example.net",
			},
			ResentTo: []*mail.Address{
				{
					Name:    "Bob Viðtakandison",
					Address: "bob.didtakandison@example.net",
				},
				{
					Name:    "Carol Viðtakandidóttir",
					Address: "carol.didtakandidottir@example.net",
				},
			},
			ResentCc: []*mail.Address{
				{
					Name:    "Dan Viðtakandison",
					Address: "dan.vidtakandison@example.net",
				},
				{
					Name:    "Eve Viðtakandidóttir",
					Address: "eve.vidtakandidottir@example.net",
				},
			},
			ResentBcc: []*mail.Address{
				{
					Name:    "Frank Viðtakandison",
					Address: "frank.vidtakandison@example.net",
				},
				{
					Name:    "Grace Viðtakandidóttir",
					Address: "grace.vidtakandidottir@example.net",
				},
			},
			ResentMessageID: "Message-Id-1@example.net",
			ContentType: ContentTypeHeader{
				ContentType: "multipart/mixed",
				Params: map[string]string{
					"boundary": "MixedBoundaryString",
					"charset":  "utf-8",
				},
			},
			ExtraHeaders: map[string][]string{
				"X-Clacks-Overhead": {"GNU Terry Pratchett"},
				"X-Script/function/\t !\"#$%&'()*+,-./;<=>?@[\\]^_`{|}~": {
					"TEST VALUE 1\t !\"#$%&'()*+,-./0123456789:;<=>?@ABCDEFGHIJKLMNOPQRSTUVWXYZ[\\]^_` abcdefghijklmnopqrstuvwxyz{|}~",
					"TEST VALUE 2\t !\"#$%&'()*+,-./0123456789:;<=>?@ABCDEFGHIJKLMNOPQRSTUVWXYZ[\\]^_` abcdefghijklmnopqrstuvwxyz{|}~",
				},
			},
		},
		Text: `Kæmi ný öxi hér, ykist þjófum nú bæði víl og ádrepa.
Svo hölt, yxna kýr þegði jú um dóp í fé á bæ.
Þú dazt á hnéð í vök og yfir blóm sexý pæju.`,
		EnrichedText: `<bold>Kæmi ný öxi hér, ykist þjófum nú bæði víl og ádrepa.</bold>
<italic>Svo hölt, yxna kýr þegði jú um dóp í fé á bæ.</italic>
<fixed>Þú dazt á hnéð í vök og yfir blóm sexý pæju.</fixed>`,
		HTML: `<html>
<div dir="ltr">
<p>Kæmi ný öxi hér, ykist þjófum nú bæði víl og ádrepa.</p>
<p>Svo hölt, yxna kýr þegði jú um dóp í fé á bæ.</p>
<p>Þú dazt á hnéð í vök og yfir blóm sexý pæju.</p>
</div>
</html>`,
		InlineFiles: []InlineFile{
			{
				ContentType: ContentTypeHeader{
					ContentType: "image/jpeg",
					Params: map[string]string{
						"name": "inline-jpg-image-without-disposition.jpg",
					},
				},
				ContentDisposition: ContentDispositionHeader{
					ContentDisposition: "",
					Params:             map[string]string(nil),
				},
				Data: []byte{255, 216, 255, 219, 0, 67, 0, 3, 2, 2, 2, 2, 2, 3, 2, 2, 2, 3, 3, 3, 3, 4,
					6, 4, 4, 4, 4, 4, 8, 6, 6, 5, 6, 9, 8, 10, 10, 9, 8, 9, 9, 10, 12, 15, 12, 10, 11, 14, 11, 9, 9,
					13, 17, 13, 14, 15, 16, 16, 17, 16, 10, 12, 18, 19, 18, 16, 19, 15, 16, 16, 16, 255, 201, 0, 11, 8,
					0, 1, 0, 1, 1, 1, 17, 0, 255, 204, 0, 6, 0, 16, 16, 5, 255, 218, 0, 8, 1, 1, 0, 0, 63, 0, 210, 207,
					32, 255, 217},
			},
			{
				ContentID: "inline-jpg-image.jpg@example.com",
				ContentType: ContentTypeHeader{
					ContentType: "image/jpeg",
					Params: map[string]string{
						"name": "inline-jpg-image-name.jpg",
					},
				},
				ContentDisposition: ContentDispositionHeader{
					ContentDisposition: inline,
					Params: map[string]string{
						"filename": "inline-jpg-image-filename.jpg",
					},
				},
				Data: []byte{255, 216, 255, 219, 0, 67, 0, 3, 2, 2, 2, 2, 2, 3, 2, 2, 2, 3, 3, 3, 3, 4,
					6, 4, 4, 4, 4, 4, 8, 6, 6, 5, 6, 9, 8, 10, 10, 9, 8, 9, 9, 10, 12, 15, 12, 10, 11, 14, 11, 9, 9,
					13, 17, 13, 14, 15, 16, 16, 17, 16, 10, 12, 18, 19, 18, 16, 19, 15, 16, 16, 16, 255, 201, 0, 11, 8,
					0, 1, 0, 1, 1, 1, 17, 0, 255, 204, 0, 6, 0, 16, 16, 5, 255, 218, 0, 8, 1, 1, 0, 0, 63, 0, 210, 207,
					32, 255, 217},
			},
		},
		AttachedFiles: []AttachedFile{
			{
				ContentType: ContentTypeHeader{
					ContentType: "application/pdf",
					Params: map[string]string{
						"name": "attached-pdf-name.pdf",
					},
				},
				ContentDisposition: ContentDispositionHeader{
					ContentDisposition: attachment,
					Params: map[string]string{
						"filename": "attached-pdf-filename.pdf",
					},
				},
				Data: []byte{37, 80, 68, 70, 45, 49, 46, 13, 116, 114, 97, 105, 108, 101, 114, 60, 60,
					47, 82, 111, 111, 116, 60, 60, 47, 80, 97, 103, 101, 115, 60, 60, 47, 75, 105, 100, 115, 91, 60,
					60, 47, 77, 101, 100, 105, 97, 66, 111, 120, 91, 48, 32, 48, 32, 51, 32, 51, 93, 62, 62, 93, 62,
					62, 62, 62, 62, 62},
			},
			{
				ContentType: ContentTypeHeader{
					ContentType: "application/pdf",
					Params: map[string]string{
						"name": "attached-pdf-without-disposition.pdf",
					},
				},
				ContentDisposition: ContentDispositionHeader{
					ContentDisposition: "",
					Params:             map[string]string(nil),
				},
				Data: []byte{37, 80, 68, 70, 45, 49, 46, 13, 116, 114, 97, 105, 108, 101, 114, 60, 60,
					47, 82, 111, 111, 116, 60, 60, 47, 80, 97, 103, 101, 115, 60, 60, 47, 75, 105, 100, 115, 91, 60,
					60, 47, 77, 101, 100, 105, 97, 66, 111, 120, 91, 48, 32, 48, 32, 51, 32, 51, 93, 62, 62, 93, 62,
					62, 62, 62, 62, 62},
			},
			{
				ContentType: ContentTypeHeader{
					ContentType: "application/json",
					Params: map[string]string{
						"name": "attached-json-name.json",
					},
				},
				ContentDisposition: ContentDispositionHeader{
					ContentDisposition: attachment,
					Params: map[string]string{
						"filename": "attached-json-filename.json",
					},
				},
				Data: []byte{123, 34, 102, 111, 111, 34, 58, 34, 98, 97, 114, 34, 125},
			},
			{
				ContentType: ContentTypeHeader{
					ContentType: "text/plain",
					Params: map[string]string{
						"name": "attached-text-plain-name.txt",
					},
				},
				ContentDisposition: ContentDispositionHeader{
					ContentDisposition: attachment,
					Params: map[string]string{
						"filename": "attached-text-plain-filename.txt",
					},
				},
				Data: []byte{84, 101, 120, 116, 47, 112, 108, 97, 105, 110, 32, 99, 111, 110, 116, 101, 110, 116, 32,
					97, 115, 32, 97, 110, 32, 97, 116, 116, 97, 99, 104, 101, 100, 32, 46, 116, 120, 116, 32, 102, 105,
					108, 101, 46},
			},
			{
				ContentType: ContentTypeHeader{
					ContentType: "text/html",
					Params: map[string]string{
						"name": "attached-text-html-name.html",
					},
				},
				ContentDisposition: ContentDispositionHeader{
					ContentDisposition: attachment,
					Params: map[string]string{
						"filename": "attached-text-html-filename.html",
					},
				},
				Data: []byte{84, 101, 120, 116, 47, 104, 116, 109, 108, 32, 99, 111, 110, 116, 101, 110, 116, 32, 97,
					115, 32, 97, 110, 32, 97, 116, 116, 97, 99, 104, 101, 100, 32, 46, 104, 116, 109, 108, 32, 102,
					105, 108, 101, 46},
			},
		},
	}

	testEmailFromFile(t, fp, expectedEmail)
}
func TestParseEmailIcelandicMultipartMixedIso88591OverBase64(t *testing.T) {
	fp := "tests/test_icelandic_multipart_mixed_iso-8859-1_over_base64.txt"
	tz, _ := time.LoadLocation("Atlantic/Reykjavik")
	expectedDate, _ := time.Parse(
		time.RFC1123Z+" (MST)",
		time.Date(2019, time.April, 1, 7, 55, 0, 0, tz).Format(time.RFC1123Z+" (MST)"))
	expectedEmail := Email{
		Headers: Headers{
			Date:    expectedDate,
			Subject: "📧 Test Íslenskt pangram",
			ReplyTo: []*mail.Address{
				{
					Name:    "Alice Sendandidóttir",
					Address: "alice.sendandidottir@example.net",
				},
			},
			Sender: &mail.Address{
				Name:    "Alice Sendandidóttir",
				Address: "alice.sendandidottir@example.com",
			},
			From: []*mail.Address{
				{
					Name:    "Alice Sendandidóttir",
					Address: "alice.sendandidottir@example.com",
				},
				{
					Name:    "Alice Sendandidóttir",
					Address: "alice.sendandidottir@example.net",
				},
			},
			To: []*mail.Address{
				{
					Name:    "Bob Viðtakandison",
					Address: "bob.didtakandison@example.com",
				},
				{
					Name:    "Carol Viðtakandidóttir",
					Address: "carol.didtakandidottir@example.com",
				},
			},
			Cc: []*mail.Address{
				{
					Name:    "Dan Viðtakandison",
					Address: "dan.vidtakandison@example.com",
				},
				{
					Name:    "Eve Viðtakandidóttir",
					Address: "eve.vidtakandidottir@example.com",
				},
			},
			Bcc: []*mail.Address{
				{
					Name:    "Frank Viðtakandison",
					Address: "frank.vidtakandison@example.com",
				},
				{
					Name:    "Grace Viðtakandidóttir",
					Address: "grace.vidtakandidottir@example.com",
				},
			},
			MessageID:  "Message-Id-1@example.com",
			InReplyTo:  []MessageId{"Message-Id-0@example.com"},
			References: []MessageId{"Message-Id-0@example.com"},
			Comments:   "Message Header Comment",
			Keywords:   []string{"Keyword 1", "Keyword 2"},
			ResentDate: expectedDate,
			ResentFrom: []*mail.Address{
				{
					Name:    "Alice Sendandidóttir",
					Address: "alice.sendandidottir@example.net",
				},
				{
					Name:    "Alice Sendandidóttir",
					Address: "alice.sendandidottir@example.com",
				},
			},
			ResentSender: &mail.Address{
				Name:    "Alice Sendandidóttir",
				Address: "alice.sendandidottir@example.net",
			},
			ResentTo: []*mail.Address{
				{
					Name:    "Bob Viðtakandison",
					Address: "bob.didtakandison@example.net",
				},
				{
					Name:    "Carol Viðtakandidóttir",
					Address: "carol.didtakandidottir@example.net",
				},
			},
			ResentCc: []*mail.Address{
				{
					Name:    "Dan Viðtakandison",
					Address: "dan.vidtakandison@example.net",
				},
				{
					Name:    "Eve Viðtakandidóttir",
					Address: "eve.vidtakandidottir@example.net",
				},
			},
			ResentBcc: []*mail.Address{
				{
					Name:    "Frank Viðtakandison",
					Address: "frank.vidtakandison@example.net",
				},
				{
					Name:    "Grace Viðtakandidóttir",
					Address: "grace.vidtakandidottir@example.net",
				},
			},
			ResentMessageID: "Message-Id-1@example.net",
			ContentType: ContentTypeHeader{
				ContentType: "multipart/mixed",
				Params: map[string]string{
					"boundary": "MixedBoundaryString",
					"charset":  "iso-8859-1",
				},
			},
			ExtraHeaders: map[string][]string{
				"X-Clacks-Overhead": {"GNU Terry Pratchett"},
				"X-Script/function/\t !\"#$%&'()*+,-./;<=>?@[\\]^_`{|}~": {
					"TEST VALUE 1\t !\"#$%&'()*+,-./0123456789:;<=>?@ABCDEFGHIJKLMNOPQRSTUVWXYZ[\\]^_` abcdefghijklmnopqrstuvwxyz{|}~",
					"TEST VALUE 2\t !\"#$%&'()*+,-./0123456789:;<=>?@ABCDEFGHIJKLMNOPQRSTUVWXYZ[\\]^_` abcdefghijklmnopqrstuvwxyz{|}~",
				},
			},
		},
		Text: `Kæmi ný öxi hér, ykist þjófum nú bæði víl og ádrepa.
Svo hölt, yxna kýr þegði jú um dóp í fé á bæ.
Þú dazt á hnéð í vök og yfir blóm sexý pæju.`,
		EnrichedText: `<bold>Kæmi ný öxi hér, ykist þjófum nú bæði víl og ádrepa.</bold>
<italic>Svo hölt, yxna kýr þegði jú um dóp í fé á bæ.</italic>
<fixed>Þú dazt á hnéð í vök og yfir blóm sexý pæju.</fixed>`,
		HTML: `<html>
<div dir="ltr">
<p>Kæmi ný öxi hér, ykist þjófum nú bæði víl og ádrepa.</p>
<p>Svo hölt, yxna kýr þegði jú um dóp í fé á bæ.</p>
<p>Þú dazt á hnéð í vök og yfir blóm sexý pæju.</p>
</div>
</html>`,
		InlineFiles: []InlineFile{
			{
				ContentType: ContentTypeHeader{
					ContentType: "image/jpeg",
					Params: map[string]string{
						"name": "inline-jpg-image-without-disposition.jpg",
					},
				},
				ContentDisposition: ContentDispositionHeader{
					ContentDisposition: "",
					Params:             map[string]string(nil),
				},
				Data: []byte{255, 216, 255, 219, 0, 67, 0, 3, 2, 2, 2, 2, 2, 3, 2, 2, 2, 3, 3, 3, 3, 4,
					6, 4, 4, 4, 4, 4, 8, 6, 6, 5, 6, 9, 8, 10, 10, 9, 8, 9, 9, 10, 12, 15, 12, 10, 11, 14, 11, 9, 9,
					13, 17, 13, 14, 15, 16, 16, 17, 16, 10, 12, 18, 19, 18, 16, 19, 15, 16, 16, 16, 255, 201, 0, 11, 8,
					0, 1, 0, 1, 1, 1, 17, 0, 255, 204, 0, 6, 0, 16, 16, 5, 255, 218, 0, 8, 1, 1, 0, 0, 63, 0, 210, 207,
					32, 255, 217},
			},
			{
				ContentID: "inline-jpg-image.jpg@example.com",
				ContentType: ContentTypeHeader{
					ContentType: "image/jpeg",
					Params: map[string]string{
						"name": "inline-jpg-image-name.jpg",
					},
				},
				ContentDisposition: ContentDispositionHeader{
					ContentDisposition: inline,
					Params: map[string]string{
						"filename": "inline-jpg-image-filename.jpg",
					},
				},
				Data: []byte{255, 216, 255, 219, 0, 67, 0, 3, 2, 2, 2, 2, 2, 3, 2, 2, 2, 3, 3, 3, 3, 4,
					6, 4, 4, 4, 4, 4, 8, 6, 6, 5, 6, 9, 8, 10, 10, 9, 8, 9, 9, 10, 12, 15, 12, 10, 11, 14, 11, 9, 9,
					13, 17, 13, 14, 15, 16, 16, 17, 16, 10, 12, 18, 19, 18, 16, 19, 15, 16, 16, 16, 255, 201, 0, 11, 8,
					0, 1, 0, 1, 1, 1, 17, 0, 255, 204, 0, 6, 0, 16, 16, 5, 255, 218, 0, 8, 1, 1, 0, 0, 63, 0, 210, 207,
					32, 255, 217},
			},
		},
		AttachedFiles: []AttachedFile{
			{
				ContentType: ContentTypeHeader{
					ContentType: "application/pdf",
					Params: map[string]string{
						"name": "attached-pdf-name.pdf",
					},
				},
				ContentDisposition: ContentDispositionHeader{
					ContentDisposition: attachment,
					Params: map[string]string{
						"filename": "attached-pdf-filename.pdf",
					},
				},
				Data: []byte{37, 80, 68, 70, 45, 49, 46, 13, 116, 114, 97, 105, 108, 101, 114, 60, 60,
					47, 82, 111, 111, 116, 60, 60, 47, 80, 97, 103, 101, 115, 60, 60, 47, 75, 105, 100, 115, 91, 60,
					60, 47, 77, 101, 100, 105, 97, 66, 111, 120, 91, 48, 32, 48, 32, 51, 32, 51, 93, 62, 62, 93, 62,
					62, 62, 62, 62, 62},
			},
			{
				ContentType: ContentTypeHeader{
					ContentType: "application/pdf",
					Params: map[string]string{
						"name": "attached-pdf-without-disposition.pdf",
					},
				},
				ContentDisposition: ContentDispositionHeader{
					ContentDisposition: "",
					Params:             map[string]string(nil),
				},
				Data: []byte{37, 80, 68, 70, 45, 49, 46, 13, 116, 114, 97, 105, 108, 101, 114, 60, 60,
					47, 82, 111, 111, 116, 60, 60, 47, 80, 97, 103, 101, 115, 60, 60, 47, 75, 105, 100, 115, 91, 60,
					60, 47, 77, 101, 100, 105, 97, 66, 111, 120, 91, 48, 32, 48, 32, 51, 32, 51, 93, 62, 62, 93, 62,
					62, 62, 62, 62, 62},
			},
			{
				ContentType: ContentTypeHeader{
					ContentType: "application/json",
					Params: map[string]string{
						"name": "attached-json-name.json",
					},
				},
				ContentDisposition: ContentDispositionHeader{
					ContentDisposition: attachment,
					Params: map[string]string{
						"filename": "attached-json-filename.json",
					},
				},
				Data: []byte{123, 34, 102, 111, 111, 34, 58, 34, 98, 97, 114, 34, 125},
			},
			{
				ContentType: ContentTypeHeader{
					ContentType: "text/plain",
					Params: map[string]string{
						"name": "attached-text-plain-name.txt",
					},
				},
				ContentDisposition: ContentDispositionHeader{
					ContentDisposition: attachment,
					Params: map[string]string{
						"filename": "attached-text-plain-filename.txt",
					},
				},
				Data: []byte{84, 101, 120, 116, 47, 112, 108, 97, 105, 110, 32, 99, 111, 110, 116, 101, 110, 116, 32,
					97, 115, 32, 97, 110, 32, 97, 116, 116, 97, 99, 104, 101, 100, 32, 46, 116, 120, 116, 32, 102, 105,
					108, 101, 46},
			},
			{
				ContentType: ContentTypeHeader{
					ContentType: "text/html",
					Params: map[string]string{
						"name": "attached-text-html-name.html",
					},
				},
				ContentDisposition: ContentDispositionHeader{
					ContentDisposition: attachment,
					Params: map[string]string{
						"filename": "attached-text-html-filename.html",
					},
				},
				Data: []byte{84, 101, 120, 116, 47, 104, 116, 109, 108, 32, 99, 111, 110, 116, 101, 110, 116, 32, 97,
					115, 32, 97, 110, 32, 97, 116, 116, 97, 99, 104, 101, 100, 32, 46, 104, 116, 109, 108, 32, 102,
					105, 108, 101, 46},
			},
		},
	}

	testEmailFromFile(t, fp, expectedEmail)
}
func TestParseEmailIcelandicMultipartMixedIso88591OverQuotedprintable(t *testing.T) {
	fp := "tests/test_icelandic_multipart_mixed_iso-8859-1_over_quoted-printable.txt"
	tz, _ := time.LoadLocation("Atlantic/Reykjavik")
	expectedDate, _ := time.Parse(
		time.RFC1123Z+" (MST)",
		time.Date(2019, time.April, 1, 7, 55, 0, 0, tz).Format(time.RFC1123Z+" (MST)"))
	expectedEmail := Email{
		Headers: Headers{
			Date:    expectedDate,
			Subject: "📧 Test Íslenskt pangram",
			ReplyTo: []*mail.Address{
				{
					Name:    "Alice Sendandidóttir",
					Address: "alice.sendandidottir@example.net",
				},
			},
			Sender: &mail.Address{
				Name:    "Alice Sendandidóttir",
				Address: "alice.sendandidottir@example.com",
			},
			From: []*mail.Address{
				{
					Name:    "Alice Sendandidóttir",
					Address: "alice.sendandidottir@example.com",
				},
				{
					Name:    "Alice Sendandidóttir",
					Address: "alice.sendandidottir@example.net",
				},
			},
			To: []*mail.Address{
				{
					Name:    "Bob Viðtakandison",
					Address: "bob.didtakandison@example.com",
				},
				{
					Name:    "Carol Viðtakandidóttir",
					Address: "carol.didtakandidottir@example.com",
				},
			},
			Cc: []*mail.Address{
				{
					Name:    "Dan Viðtakandison",
					Address: "dan.vidtakandison@example.com",
				},
				{
					Name:    "Eve Viðtakandidóttir",
					Address: "eve.vidtakandidottir@example.com",
				},
			},
			Bcc: []*mail.Address{
				{
					Name:    "Frank Viðtakandison",
					Address: "frank.vidtakandison@example.com",
				},
				{
					Name:    "Grace Viðtakandidóttir",
					Address: "grace.vidtakandidottir@example.com",
				},
			},
			MessageID:  "Message-Id-1@example.com",
			InReplyTo:  []MessageId{"Message-Id-0@example.com"},
			References: []MessageId{"Message-Id-0@example.com"},
			Comments:   "Message Header Comment",
			Keywords:   []string{"Keyword 1", "Keyword 2"},
			ResentDate: expectedDate,
			ResentFrom: []*mail.Address{
				{
					Name:    "Alice Sendandidóttir",
					Address: "alice.sendandidottir@example.net",
				},
				{
					Name:    "Alice Sendandidóttir",
					Address: "alice.sendandidottir@example.com",
				},
			},
			ResentSender: &mail.Address{
				Name:    "Alice Sendandidóttir",
				Address: "alice.sendandidottir@example.net",
			},
			ResentTo: []*mail.Address{
				{
					Name:    "Bob Viðtakandison",
					Address: "bob.didtakandison@example.net",
				},
				{
					Name:    "Carol Viðtakandidóttir",
					Address: "carol.didtakandidottir@example.net",
				},
			},
			ResentCc: []*mail.Address{
				{
					Name:    "Dan Viðtakandison",
					Address: "dan.vidtakandison@example.net",
				},
				{
					Name:    "Eve Viðtakandidóttir",
					Address: "eve.vidtakandidottir@example.net",
				},
			},
			ResentBcc: []*mail.Address{
				{
					Name:    "Frank Viðtakandison",
					Address: "frank.vidtakandison@example.net",
				},
				{
					Name:    "Grace Viðtakandidóttir",
					Address: "grace.vidtakandidottir@example.net",
				},
			},
			ResentMessageID: "Message-Id-1@example.net",
			ContentType: ContentTypeHeader{
				ContentType: "multipart/mixed",
				Params: map[string]string{
					"boundary": "MixedBoundaryString",
					"charset":  "iso-8859-1",
				},
			},
			ExtraHeaders: map[string][]string{
				"X-Clacks-Overhead": {"GNU Terry Pratchett"},
				"X-Script/function/\t !\"#$%&'()*+,-./;<=>?@[\\]^_`{|}~": {
					"TEST VALUE 1\t !\"#$%&'()*+,-./0123456789:;<=>?@ABCDEFGHIJKLMNOPQRSTUVWXYZ[\\]^_` abcdefghijklmnopqrstuvwxyz{|}~",
					"TEST VALUE 2\t !\"#$%&'()*+,-./0123456789:;<=>?@ABCDEFGHIJKLMNOPQRSTUVWXYZ[\\]^_` abcdefghijklmnopqrstuvwxyz{|}~",
				},
			},
		},
		Text: `Kæmi ný öxi hér, ykist þjófum nú bæði víl og ádrepa.
Svo hölt, yxna kýr þegði jú um dóp í fé á bæ.
Þú dazt á hnéð í vök og yfir blóm sexý pæju.`,
		EnrichedText: `<bold>Kæmi ný öxi hér, ykist þjófum nú bæði víl og ádrepa.</bold>
<italic>Svo hölt, yxna kýr þegði jú um dóp í fé á bæ.</italic>
<fixed>Þú dazt á hnéð í vök og yfir blóm sexý pæju.</fixed>`,
		HTML: `<html>
<div dir="ltr">
<p>Kæmi ný öxi hér, ykist þjófum nú bæði víl og ádrepa.</p>
<p>Svo hölt, yxna kýr þegði jú um dóp í fé á bæ.</p>
<p>Þú dazt á hnéð í vök og yfir blóm sexý pæju.</p>
</div>
</html>`,
		InlineFiles: []InlineFile{
			{
				ContentType: ContentTypeHeader{
					ContentType: "image/jpeg",
					Params: map[string]string{
						"name": "inline-jpg-image-without-disposition.jpg",
					},
				},
				ContentDisposition: ContentDispositionHeader{
					ContentDisposition: "",
					Params:             map[string]string(nil),
				},
				Data: []byte{255, 216, 255, 219, 0, 67, 0, 3, 2, 2, 2, 2, 2, 3, 2, 2, 2, 3, 3, 3, 3, 4,
					6, 4, 4, 4, 4, 4, 8, 6, 6, 5, 6, 9, 8, 10, 10, 9, 8, 9, 9, 10, 12, 15, 12, 10, 11, 14, 11, 9, 9,
					13, 17, 13, 14, 15, 16, 16, 17, 16, 10, 12, 18, 19, 18, 16, 19, 15, 16, 16, 16, 255, 201, 0, 11, 8,
					0, 1, 0, 1, 1, 1, 17, 0, 255, 204, 0, 6, 0, 16, 16, 5, 255, 218, 0, 8, 1, 1, 0, 0, 63, 0, 210, 207,
					32, 255, 217},
			},
			{
				ContentID: "inline-jpg-image.jpg@example.com",
				ContentType: ContentTypeHeader{
					ContentType: "image/jpeg",
					Params: map[string]string{
						"name": "inline-jpg-image-name.jpg",
					},
				},
				ContentDisposition: ContentDispositionHeader{
					ContentDisposition: inline,
					Params: map[string]string{
						"filename": "inline-jpg-image-filename.jpg",
					},
				},
				Data: []byte{255, 216, 255, 219, 0, 67, 0, 3, 2, 2, 2, 2, 2, 3, 2, 2, 2, 3, 3, 3, 3, 4,
					6, 4, 4, 4, 4, 4, 8, 6, 6, 5, 6, 9, 8, 10, 10, 9, 8, 9, 9, 10, 12, 15, 12, 10, 11, 14, 11, 9, 9,
					13, 17, 13, 14, 15, 16, 16, 17, 16, 10, 12, 18, 19, 18, 16, 19, 15, 16, 16, 16, 255, 201, 0, 11, 8,
					0, 1, 0, 1, 1, 1, 17, 0, 255, 204, 0, 6, 0, 16, 16, 5, 255, 218, 0, 8, 1, 1, 0, 0, 63, 0, 210, 207,
					32, 255, 217},
			},
		},
		AttachedFiles: []AttachedFile{
			{
				ContentType: ContentTypeHeader{
					ContentType: "application/pdf",
					Params: map[string]string{
						"name": "attached-pdf-name.pdf",
					},
				},
				ContentDisposition: ContentDispositionHeader{
					ContentDisposition: attachment,
					Params: map[string]string{
						"filename": "attached-pdf-filename.pdf",
					},
				},
				Data: []byte{37, 80, 68, 70, 45, 49, 46, 13, 116, 114, 97, 105, 108, 101, 114, 60, 60,
					47, 82, 111, 111, 116, 60, 60, 47, 80, 97, 103, 101, 115, 60, 60, 47, 75, 105, 100, 115, 91, 60,
					60, 47, 77, 101, 100, 105, 97, 66, 111, 120, 91, 48, 32, 48, 32, 51, 32, 51, 93, 62, 62, 93, 62,
					62, 62, 62, 62, 62},
			},
			{
				ContentType: ContentTypeHeader{
					ContentType: "application/pdf",
					Params: map[string]string{
						"name": "attached-pdf-without-disposition.pdf",
					},
				},
				ContentDisposition: ContentDispositionHeader{
					ContentDisposition: "",
					Params:             map[string]string(nil),
				},
				Data: []byte{37, 80, 68, 70, 45, 49, 46, 13, 116, 114, 97, 105, 108, 101, 114, 60, 60,
					47, 82, 111, 111, 116, 60, 60, 47, 80, 97, 103, 101, 115, 60, 60, 47, 75, 105, 100, 115, 91, 60,
					60, 47, 77, 101, 100, 105, 97, 66, 111, 120, 91, 48, 32, 48, 32, 51, 32, 51, 93, 62, 62, 93, 62,
					62, 62, 62, 62, 62},
			},
			{
				ContentType: ContentTypeHeader{
					ContentType: "application/json",
					Params: map[string]string{
						"name": "attached-json-name.json",
					},
				},
				ContentDisposition: ContentDispositionHeader{
					ContentDisposition: attachment,
					Params: map[string]string{
						"filename": "attached-json-filename.json",
					},
				},
				Data: []byte{123, 34, 102, 111, 111, 34, 58, 34, 98, 97, 114, 34, 125},
			},
			{
				ContentType: ContentTypeHeader{
					ContentType: "text/plain",
					Params: map[string]string{
						"name": "attached-text-plain-name.txt",
					},
				},
				ContentDisposition: ContentDispositionHeader{
					ContentDisposition: attachment,
					Params: map[string]string{
						"filename": "attached-text-plain-filename.txt",
					},
				},
				Data: []byte{84, 101, 120, 116, 47, 112, 108, 97, 105, 110, 32, 99, 111, 110, 116, 101, 110, 116, 32,
					97, 115, 32, 97, 110, 32, 97, 116, 116, 97, 99, 104, 101, 100, 32, 46, 116, 120, 116, 32, 102, 105,
					108, 101, 46},
			},
			{
				ContentType: ContentTypeHeader{
					ContentType: "text/html",
					Params: map[string]string{
						"name": "attached-text-html-name.html",
					},
				},
				ContentDisposition: ContentDispositionHeader{
					ContentDisposition: attachment,
					Params: map[string]string{
						"filename": "attached-text-html-filename.html",
					},
				},
				Data: []byte{84, 101, 120, 116, 47, 104, 116, 109, 108, 32, 99, 111, 110, 116, 101, 110, 116, 32, 97,
					115, 32, 97, 110, 32, 97, 116, 116, 97, 99, 104, 101, 100, 32, 46, 104, 116, 109, 108, 32, 102,
					105, 108, 101, 46},
			},
		},
	}

	testEmailFromFile(t, fp, expectedEmail)
}
func TestParseEmailIcelandicMultipartSignedUtf8OverBase64(t *testing.T) {
	fp := "tests/test_icelandic_multipart_signed_utf-8_over_base64.txt"
	tz, _ := time.LoadLocation("Atlantic/Reykjavik")
	expectedDate, _ := time.Parse(
		time.RFC1123Z+" (MST)",
		time.Date(2019, time.April, 1, 7, 55, 0, 0, tz).Format(time.RFC1123Z+" (MST)"))
	expectedEmail := Email{
		Headers: Headers{
			Date:    expectedDate,
			Subject: "📧 Signed Test Íslenskt pangram",
			ReplyTo: []*mail.Address{
				{
					Name:    "Alice Sendandidóttir",
					Address: "alice.sendandidottir@example.net",
				},
			},
			Sender: &mail.Address{
				Name:    "Alice Sendandidóttir",
				Address: "alice.sendandidottir@example.com",
			},
			From: []*mail.Address{
				{
					Name:    "Alice Sendandidóttir",
					Address: "alice.sendandidottir@example.com",
				},
				{
					Name:    "Alice Sendandidóttir",
					Address: "alice.sendandidottir@example.net",
				},
			},
			To: []*mail.Address{
				{
					Name:    "Bob Viðtakandison",
					Address: "bob.didtakandison@example.com",
				},
				{
					Name:    "Carol Viðtakandidóttir",
					Address: "carol.didtakandidottir@example.com",
				},
			},
			Cc: []*mail.Address{
				{
					Name:    "Dan Viðtakandison",
					Address: "dan.vidtakandison@example.com",
				},
				{
					Name:    "Eve Viðtakandidóttir",
					Address: "eve.vidtakandidottir@example.com",
				},
			},
			Bcc: []*mail.Address{
				{
					Name:    "Frank Viðtakandison",
					Address: "frank.vidtakandison@example.com",
				},
				{
					Name:    "Grace Viðtakandidóttir",
					Address: "grace.vidtakandidottir@example.com",
				},
			},
			MessageID:  "Message-Id-1@example.com",
			InReplyTo:  []MessageId{"Message-Id-0@example.com"},
			References: []MessageId{"Message-Id-0@example.com"},
			Comments:   "Message Header Comment",
			Keywords:   []string{"Keyword 1", "Keyword 2"},
			ResentDate: expectedDate,
			ResentFrom: []*mail.Address{
				{
					Name:    "Alice Sendandidóttir",
					Address: "alice.sendandidottir@example.net",
				},
				{
					Name:    "Alice Sendandidóttir",
					Address: "alice.sendandidottir@example.com",
				},
			},
			ResentSender: &mail.Address{
				Name:    "Alice Sendandidóttir",
				Address: "alice.sendandidottir@example.net",
			},
			ResentTo: []*mail.Address{
				{
					Name:    "Bob Viðtakandison",
					Address: "bob.didtakandison@example.net",
				},
				{
					Name:    "Carol Viðtakandidóttir",
					Address: "carol.didtakandidottir@example.net",
				},
			},
			ResentCc: []*mail.Address{
				{
					Name:    "Dan Viðtakandison",
					Address: "dan.vidtakandison@example.net",
				},
				{
					Name:    "Eve Viðtakandidóttir",
					Address: "eve.vidtakandidottir@example.net",
				},
			},
			ResentBcc: []*mail.Address{
				{
					Name:    "Frank Viðtakandison",
					Address: "frank.vidtakandison@example.net",
				},
				{
					Name:    "Grace Viðtakandidóttir",
					Address: "grace.vidtakandidottir@example.net",
				},
			},
			ResentMessageID: "Message-Id-1@example.net",
			ContentType: ContentTypeHeader{
				ContentType: "multipart/signed",
				Params: map[string]string{
					"boundary": "SignedBoundaryString",
					"charset":  "utf-8",
					"micalg":   "sha1",
					"protocol": "application/pkcs7-signature",
				},
			},
			ExtraHeaders: map[string][]string{
				"X-Clacks-Overhead": {"GNU Terry Pratchett"},
				"X-Script/function/\t !\"#$%&'()*+,-./;<=>?@[\\]^_`{|}~": {
					"TEST VALUE 1\t !\"#$%&'()*+,-./0123456789:;<=>?@ABCDEFGHIJKLMNOPQRSTUVWXYZ[\\]^_` abcdefghijklmnopqrstuvwxyz{|}~",
					"TEST VALUE 2\t !\"#$%&'()*+,-./0123456789:;<=>?@ABCDEFGHIJKLMNOPQRSTUVWXYZ[\\]^_` abcdefghijklmnopqrstuvwxyz{|}~",
				},
			},
		},
		Text: `Kæmi ný öxi hér, ykist þjófum nú bæði víl og ádrepa.
Svo hölt, yxna kýr þegði jú um dóp í fé á bæ.
Þú dazt á hnéð í vök og yfir blóm sexý pæju.`,
		AttachedFiles: []AttachedFile{
			{
				ContentType: ContentTypeHeader{
					ContentType: "application/pkcs7-signature",
					Params: map[string]string{
						"name": "smime.p7s",
					},
				},
				ContentDisposition: ContentDispositionHeader{
					ContentDisposition: attachment,
					Params: map[string]string{
						"filename": "smime.p7s",
					},
				},
				Data: []byte{130, 28, 135, 132, 117, 46, 142, 18, 97, 140, 126, 251, 159, 193, 199, 25, 58, 223, 189, 185, 227, 239, 158, 173, 108, 31, 71, 27, 133, 80, 165, 252, 133, 227, 174, 198, 132, 129, 159, 29, 246, 19, 235, 133, 80, 165, 252, 133, 227, 174, 198, 132, 129, 159, 29, 246, 19, 234, 49, 251, 238, 127, 7, 28, 104, 33, 200, 120, 71, 82, 232, 225, 38, 30, 249, 234, 214, 193, 244, 113, 147, 173, 251, 219, 158, 57, 252, 28, 113, 147, 173, 251, 225, 38, 24, 199, 239, 190, 173, 108, 31, 71, 27, 133, 80, 110, 120, 251, 231, 174, 198, 132, 129, 159, 29, 246, 19, 234, 8, 114, 30, 17, 212, 186, 58, 95, 200, 94, 59, 26, 18, 6, 124, 119, 216, 79, 174, 21, 65, 185, 227, 239, 158},
			},
		},
	}

	testEmailFromFile(t, fp, expectedEmail)
}

func TestParseEmailIcelandicMultipartSignedUtf8OverQuotedprintable(t *testing.T) {
	fp := "tests/test_icelandic_multipart_signed_utf-8_over_quoted-printable.txt"
	tz, _ := time.LoadLocation("Atlantic/Reykjavik")
	expectedDate, _ := time.Parse(
		time.RFC1123Z+" (MST)",
		time.Date(2019, time.April, 1, 7, 55, 0, 0, tz).Format(time.RFC1123Z+" (MST)"))
	expectedEmail := Email{
		Headers: Headers{
			Date:    expectedDate,
			Subject: "📧 Signed Test Íslenskt pangram",
			ReplyTo: []*mail.Address{
				{
					Name:    "Alice Sendandidóttir",
					Address: "alice.sendandidottir@example.net",
				},
			},
			Sender: &mail.Address{
				Name:    "Alice Sendandidóttir",
				Address: "alice.sendandidottir@example.com",
			},
			From: []*mail.Address{
				{
					Name:    "Alice Sendandidóttir",
					Address: "alice.sendandidottir@example.com",
				},
				{
					Name:    "Alice Sendandidóttir",
					Address: "alice.sendandidottir@example.net",
				},
			},
			To: []*mail.Address{
				{
					Name:    "Bob Viðtakandison",
					Address: "bob.didtakandison@example.com",
				},
				{
					Name:    "Carol Viðtakandidóttir",
					Address: "carol.didtakandidottir@example.com",
				},
			},
			Cc: []*mail.Address{
				{
					Name:    "Dan Viðtakandison",
					Address: "dan.vidtakandison@example.com",
				},
				{
					Name:    "Eve Viðtakandidóttir",
					Address: "eve.vidtakandidottir@example.com",
				},
			},
			Bcc: []*mail.Address{
				{
					Name:    "Frank Viðtakandison",
					Address: "frank.vidtakandison@example.com",
				},
				{
					Name:    "Grace Viðtakandidóttir",
					Address: "grace.vidtakandidottir@example.com",
				},
			},
			MessageID:  "Message-Id-1@example.com",
			InReplyTo:  []MessageId{"Message-Id-0@example.com"},
			References: []MessageId{"Message-Id-0@example.com"},
			Comments:   "Message Header Comment",
			Keywords:   []string{"Keyword 1", "Keyword 2"},
			ResentDate: expectedDate,
			ResentFrom: []*mail.Address{
				{
					Name:    "Alice Sendandidóttir",
					Address: "alice.sendandidottir@example.net",
				},
				{
					Name:    "Alice Sendandidóttir",
					Address: "alice.sendandidottir@example.com",
				},
			},
			ResentSender: &mail.Address{
				Name:    "Alice Sendandidóttir",
				Address: "alice.sendandidottir@example.net",
			},
			ResentTo: []*mail.Address{
				{
					Name:    "Bob Viðtakandison",
					Address: "bob.didtakandison@example.net",
				},
				{
					Name:    "Carol Viðtakandidóttir",
					Address: "carol.didtakandidottir@example.net",
				},
			},
			ResentCc: []*mail.Address{
				{
					Name:    "Dan Viðtakandison",
					Address: "dan.vidtakandison@example.net",
				},
				{
					Name:    "Eve Viðtakandidóttir",
					Address: "eve.vidtakandidottir@example.net",
				},
			},
			ResentBcc: []*mail.Address{
				{
					Name:    "Frank Viðtakandison",
					Address: "frank.vidtakandison@example.net",
				},
				{
					Name:    "Grace Viðtakandidóttir",
					Address: "grace.vidtakandidottir@example.net",
				},
			},
			ResentMessageID: "Message-Id-1@example.net",
			ContentType: ContentTypeHeader{
				ContentType: "multipart/signed",
				Params: map[string]string{
					"boundary": "SignedBoundaryString",
					"charset":  "utf-8",
					"micalg":   "sha1",
					"protocol": "application/pkcs7-signature",
				},
			},
			ExtraHeaders: map[string][]string{
				"X-Clacks-Overhead": {"GNU Terry Pratchett"},
				"X-Script/function/\t !\"#$%&'()*+,-./;<=>?@[\\]^_`{|}~": {
					"TEST VALUE 1\t !\"#$%&'()*+,-./0123456789:;<=>?@ABCDEFGHIJKLMNOPQRSTUVWXYZ[\\]^_` abcdefghijklmnopqrstuvwxyz{|}~",
					"TEST VALUE 2\t !\"#$%&'()*+,-./0123456789:;<=>?@ABCDEFGHIJKLMNOPQRSTUVWXYZ[\\]^_` abcdefghijklmnopqrstuvwxyz{|}~",
				},
			},
		},
		Text: `Kæmi ný öxi hér, ykist þjófum nú bæði víl og ádrepa.
Svo hölt, yxna kýr þegði jú um dóp í fé á bæ.
Þú dazt á hnéð í vök og yfir blóm sexý pæju.`,
		AttachedFiles: []AttachedFile{
			{
				ContentType: ContentTypeHeader{
					ContentType: "application/pkcs7-signature",
					Params: map[string]string{
						"name": "smime.p7s",
					},
				},
				ContentDisposition: ContentDispositionHeader{
					ContentDisposition: attachment,
					Params: map[string]string{
						"filename": "smime.p7s",
					},
				},
				Data: []byte{130, 28, 135, 132, 117, 46, 142, 18, 97, 140, 126, 251, 159, 193, 199, 25, 58, 223, 189, 185, 227, 239, 158, 173, 108, 31, 71, 27, 133, 80, 165, 252, 133, 227, 174, 198, 132, 129, 159, 29, 246, 19, 235, 133, 80, 165, 252, 133, 227, 174, 198, 132, 129, 159, 29, 246, 19, 234, 49, 251, 238, 127, 7, 28, 104, 33, 200, 120, 71, 82, 232, 225, 38, 30, 249, 234, 214, 193, 244, 113, 147, 173, 251, 219, 158, 57, 252, 28, 113, 147, 173, 251, 225, 38, 24, 199, 239, 190, 173, 108, 31, 71, 27, 133, 80, 110, 120, 251, 231, 174, 198, 132, 129, 159, 29, 246, 19, 234, 8, 114, 30, 17, 212, 186, 58, 95, 200, 94, 59, 26, 18, 6, 124, 119, 216, 79, 174, 21, 65, 185, 227, 239, 158},
			},
		},
	}

	testEmailFromFile(t, fp, expectedEmail)
}

func TestParseEmailIcelandicMultipartSignedIso88591OverBase64(t *testing.T) {
	fp := "tests/test_icelandic_multipart_signed_iso-8859-1_over_base64.txt"
	tz, _ := time.LoadLocation("Atlantic/Reykjavik")
	expectedDate, _ := time.Parse(
		time.RFC1123Z+" (MST)",
		time.Date(2019, time.April, 1, 7, 55, 0, 0, tz).Format(time.RFC1123Z+" (MST)"))
	expectedEmail := Email{
		Headers: Headers{
			Date:    expectedDate,
			Subject: "📧 Signed Test Íslenskt pangram",
			ReplyTo: []*mail.Address{
				{
					Name:    "Alice Sendandidóttir",
					Address: "alice.sendandidottir@example.net",
				},
			},
			Sender: &mail.Address{
				Name:    "Alice Sendandidóttir",
				Address: "alice.sendandidottir@example.com",
			},
			From: []*mail.Address{
				{
					Name:    "Alice Sendandidóttir",
					Address: "alice.sendandidottir@example.com",
				},
				{
					Name:    "Alice Sendandidóttir",
					Address: "alice.sendandidottir@example.net",
				},
			},
			To: []*mail.Address{
				{
					Name:    "Bob Viðtakandison",
					Address: "bob.didtakandison@example.com",
				},
				{
					Name:    "Carol Viðtakandidóttir",
					Address: "carol.didtakandidottir@example.com",
				},
			},
			Cc: []*mail.Address{
				{
					Name:    "Dan Viðtakandison",
					Address: "dan.vidtakandison@example.com",
				},
				{
					Name:    "Eve Viðtakandidóttir",
					Address: "eve.vidtakandidottir@example.com",
				},
			},
			Bcc: []*mail.Address{
				{
					Name:    "Frank Viðtakandison",
					Address: "frank.vidtakandison@example.com",
				},
				{
					Name:    "Grace Viðtakandidóttir",
					Address: "grace.vidtakandidottir@example.com",
				},
			},
			MessageID:  "Message-Id-1@example.com",
			InReplyTo:  []MessageId{"Message-Id-0@example.com"},
			References: []MessageId{"Message-Id-0@example.com"},
			Comments:   "Message Header Comment",
			Keywords:   []string{"Keyword 1", "Keyword 2"},
			ResentDate: expectedDate,
			ResentFrom: []*mail.Address{
				{
					Name:    "Alice Sendandidóttir",
					Address: "alice.sendandidottir@example.net",
				},
				{
					Name:    "Alice Sendandidóttir",
					Address: "alice.sendandidottir@example.com",
				},
			},
			ResentSender: &mail.Address{
				Name:    "Alice Sendandidóttir",
				Address: "alice.sendandidottir@example.net",
			},
			ResentTo: []*mail.Address{
				{
					Name:    "Bob Viðtakandison",
					Address: "bob.didtakandison@example.net",
				},
				{
					Name:    "Carol Viðtakandidóttir",
					Address: "carol.didtakandidottir@example.net",
				},
			},
			ResentCc: []*mail.Address{
				{
					Name:    "Dan Viðtakandison",
					Address: "dan.vidtakandison@example.net",
				},
				{
					Name:    "Eve Viðtakandidóttir",
					Address: "eve.vidtakandidottir@example.net",
				},
			},
			ResentBcc: []*mail.Address{
				{
					Name:    "Frank Viðtakandison",
					Address: "frank.vidtakandison@example.net",
				},
				{
					Name:    "Grace Viðtakandidóttir",
					Address: "grace.vidtakandidottir@example.net",
				},
			},
			ResentMessageID: "Message-Id-1@example.net",
			ContentType: ContentTypeHeader{
				ContentType: "multipart/signed",
				Params: map[string]string{
					"boundary": "SignedBoundaryString",
					"charset":  "iso-8859-1",
					"micalg":   "sha1",
					"protocol": "application/pkcs7-signature",
				},
			},
			ExtraHeaders: map[string][]string{
				"X-Clacks-Overhead": {"GNU Terry Pratchett"},
				"X-Script/function/\t !\"#$%&'()*+,-./;<=>?@[\\]^_`{|}~": {
					"TEST VALUE 1\t !\"#$%&'()*+,-./0123456789:;<=>?@ABCDEFGHIJKLMNOPQRSTUVWXYZ[\\]^_` abcdefghijklmnopqrstuvwxyz{|}~",
					"TEST VALUE 2\t !\"#$%&'()*+,-./0123456789:;<=>?@ABCDEFGHIJKLMNOPQRSTUVWXYZ[\\]^_` abcdefghijklmnopqrstuvwxyz{|}~",
				},
			},
		},
		Text: `Kæmi ný öxi hér, ykist þjófum nú bæði víl og ádrepa.
Svo hölt, yxna kýr þegði jú um dóp í fé á bæ.
Þú dazt á hnéð í vök og yfir blóm sexý pæju.`,
		AttachedFiles: []AttachedFile{
			{
				ContentType: ContentTypeHeader{
					ContentType: "application/pkcs7-signature",
					Params: map[string]string{
						"name": "smime.p7s",
					},
				},
				ContentDisposition: ContentDispositionHeader{
					ContentDisposition: attachment,
					Params: map[string]string{
						"filename": "smime.p7s",
					},
				},
				Data: []byte{130, 28, 135, 132, 117, 46, 142, 18, 97, 140, 126, 251, 159, 193, 199, 25, 58, 223, 189, 185, 227, 239, 158, 173, 108, 31, 71, 27, 133, 80, 165, 252, 133, 227, 174, 198, 132, 129, 159, 29, 246, 19, 235, 133, 80, 165, 252, 133, 227, 174, 198, 132, 129, 159, 29, 246, 19, 234, 49, 251, 238, 127, 7, 28, 104, 33, 200, 120, 71, 82, 232, 225, 38, 30, 249, 234, 214, 193, 244, 113, 147, 173, 251, 219, 158, 57, 252, 28, 113, 147, 173, 251, 225, 38, 24, 199, 239, 190, 173, 108, 31, 71, 27, 133, 80, 110, 120, 251, 231, 174, 198, 132, 129, 159, 29, 246, 19, 234, 8, 114, 30, 17, 212, 186, 58, 95, 200, 94, 59, 26, 18, 6, 124, 119, 216, 79, 174, 21, 65, 185, 227, 239, 158},
			},
		},
	}

	testEmailFromFile(t, fp, expectedEmail)
}

func TestParseEmailIcelandicMultipartSignedIso88591OverQuotedprintable(t *testing.T) {
	fp := "tests/test_icelandic_multipart_signed_iso-8859-1_over_quoted-printable.txt"
	tz, _ := time.LoadLocation("Atlantic/Reykjavik")
	expectedDate, _ := time.Parse(
		time.RFC1123Z+" (MST)",
		time.Date(2019, time.April, 1, 7, 55, 0, 0, tz).Format(time.RFC1123Z+" (MST)"))
	expectedEmail := Email{
		Headers: Headers{
			Date:    expectedDate,
			Subject: "📧 Signed Test Íslenskt pangram",
			ReplyTo: []*mail.Address{
				{
					Name:    "Alice Sendandidóttir",
					Address: "alice.sendandidottir@example.net",
				},
			},
			Sender: &mail.Address{
				Name:    "Alice Sendandidóttir",
				Address: "alice.sendandidottir@example.com",
			},
			From: []*mail.Address{
				{
					Name:    "Alice Sendandidóttir",
					Address: "alice.sendandidottir@example.com",
				},
				{
					Name:    "Alice Sendandidóttir",
					Address: "alice.sendandidottir@example.net",
				},
			},
			To: []*mail.Address{
				{
					Name:    "Bob Viðtakandison",
					Address: "bob.didtakandison@example.com",
				},
				{
					Name:    "Carol Viðtakandidóttir",
					Address: "carol.didtakandidottir@example.com",
				},
			},
			Cc: []*mail.Address{
				{
					Name:    "Dan Viðtakandison",
					Address: "dan.vidtakandison@example.com",
				},
				{
					Name:    "Eve Viðtakandidóttir",
					Address: "eve.vidtakandidottir@example.com",
				},
			},
			Bcc: []*mail.Address{
				{
					Name:    "Frank Viðtakandison",
					Address: "frank.vidtakandison@example.com",
				},
				{
					Name:    "Grace Viðtakandidóttir",
					Address: "grace.vidtakandidottir@example.com",
				},
			},
			MessageID:  "Message-Id-1@example.com",
			InReplyTo:  []MessageId{"Message-Id-0@example.com"},
			References: []MessageId{"Message-Id-0@example.com"},
			Comments:   "Message Header Comment",
			Keywords:   []string{"Keyword 1", "Keyword 2"},
			ResentDate: expectedDate,
			ResentFrom: []*mail.Address{
				{
					Name:    "Alice Sendandidóttir",
					Address: "alice.sendandidottir@example.net",
				},
				{
					Name:    "Alice Sendandidóttir",
					Address: "alice.sendandidottir@example.com",
				},
			},
			ResentSender: &mail.Address{
				Name:    "Alice Sendandidóttir",
				Address: "alice.sendandidottir@example.net",
			},
			ResentTo: []*mail.Address{
				{
					Name:    "Bob Viðtakandison",
					Address: "bob.didtakandison@example.net",
				},
				{
					Name:    "Carol Viðtakandidóttir",
					Address: "carol.didtakandidottir@example.net",
				},
			},
			ResentCc: []*mail.Address{
				{
					Name:    "Dan Viðtakandison",
					Address: "dan.vidtakandison@example.net",
				},
				{
					Name:    "Eve Viðtakandidóttir",
					Address: "eve.vidtakandidottir@example.net",
				},
			},
			ResentBcc: []*mail.Address{
				{
					Name:    "Frank Viðtakandison",
					Address: "frank.vidtakandison@example.net",
				},
				{
					Name:    "Grace Viðtakandidóttir",
					Address: "grace.vidtakandidottir@example.net",
				},
			},
			ResentMessageID: "Message-Id-1@example.net",
			ContentType: ContentTypeHeader{
				ContentType: "multipart/signed",
				Params: map[string]string{
					"boundary": "SignedBoundaryString",
					"charset":  "iso-8859-1",
					"micalg":   "sha1",
					"protocol": "application/pkcs7-signature",
				},
			},
			ExtraHeaders: map[string][]string{
				"X-Clacks-Overhead": {"GNU Terry Pratchett"},
				"X-Script/function/\t !\"#$%&'()*+,-./;<=>?@[\\]^_`{|}~": {
					"TEST VALUE 1\t !\"#$%&'()*+,-./0123456789:;<=>?@ABCDEFGHIJKLMNOPQRSTUVWXYZ[\\]^_` abcdefghijklmnopqrstuvwxyz{|}~",
					"TEST VALUE 2\t !\"#$%&'()*+,-./0123456789:;<=>?@ABCDEFGHIJKLMNOPQRSTUVWXYZ[\\]^_` abcdefghijklmnopqrstuvwxyz{|}~",
				},
			},
		},
		Text: `Kæmi ný öxi hér, ykist þjófum nú bæði víl og ádrepa.
Svo hölt, yxna kýr þegði jú um dóp í fé á bæ.
Þú dazt á hnéð í vök og yfir blóm sexý pæju.`,
		AttachedFiles: []AttachedFile{
			{
				ContentType: ContentTypeHeader{
					ContentType: "application/pkcs7-signature",
					Params: map[string]string{
						"name": "smime.p7s",
					},
				},
				ContentDisposition: ContentDispositionHeader{
					ContentDisposition: attachment,
					Params: map[string]string{
						"filename": "smime.p7s",
					},
				},
				Data: []byte{130, 28, 135, 132, 117, 46, 142, 18, 97, 140, 126, 251, 159, 193, 199, 25, 58, 223, 189, 185, 227, 239, 158, 173, 108, 31, 71, 27, 133, 80, 165, 252, 133, 227, 174, 198, 132, 129, 159, 29, 246, 19, 235, 133, 80, 165, 252, 133, 227, 174, 198, 132, 129, 159, 29, 246, 19, 234, 49, 251, 238, 127, 7, 28, 104, 33, 200, 120, 71, 82, 232, 225, 38, 30, 249, 234, 214, 193, 244, 113, 147, 173, 251, 219, 158, 57, 252, 28, 113, 147, 173, 251, 225, 38, 24, 199, 239, 190, 173, 108, 31, 71, 27, 133, 80, 110, 120, 251, 231, 174, 198, 132, 129, 159, 29, 246, 19, 234, 8, 114, 30, 17, 212, 186, 58, 95, 200, 94, 59, 26, 18, 6, 124, 119, 216, 79, 174, 21, 65, 185, 227, 239, 158},
			},
		},
	}

	testEmailFromFile(t, fp, expectedEmail)
}

func TestParseEmailJapanesePlaintextUtf8Over7bit(t *testing.T) {
	fp := "tests/test_japanese_plaintext_utf-8_over_7bit.txt"
	tz, _ := time.LoadLocation("Asia/Tokyo")
	expectedDate, _ := time.Parse(
		time.RFC1123Z+" (MST)",
		time.Date(2019, time.April, 1, 7, 55, 0, 0, tz).Format(time.RFC1123Z+" (MST)"))
	expectedEmail := Email{
		Headers: Headers{
			Date:    expectedDate,
			Subject: "📧 Test いろは歌",
			ReplyTo: []*mail.Address{
				{
					Name:    "郵便アリス",
					Address: "alice.yuubin@example.net",
				},
			},
			Sender: &mail.Address{
				Name:    "郵便アリス",
				Address: "alice.yuubin@example.com",
			},
			From: []*mail.Address{
				{
					Name:    "郵便アリス",
					Address: "alice.yuubin@example.com",
				},
				{
					Name:    "郵便アリス",
					Address: "alice.yuubin@example.net",
				},
			},
			To: []*mail.Address{
				{
					Name:    "郵便ボブ",
					Address: "bob.yuubin@example.com",
				},
				{
					Name:    "郵便キャロル",
					Address: "carol.yuubin@example.com",
				},
			},
			Cc: []*mail.Address{
				{
					Name:    "郵便ダン",
					Address: "dan.yuubin@example.com",
				},
				{
					Name:    "郵便イーブ",
					Address: "eve.yubin@example.com",
				},
			},
			Bcc: []*mail.Address{
				{
					Name:    "郵便フランク",
					Address: "frank.yuubin@example.com",
				},
				{
					Name:    "郵便グレイス",
					Address: "grace.yubin@example.com",
				},
			},
			MessageID:  "Message-Id-1@example.com",
			InReplyTo:  []MessageId{"Message-Id-0@example.com"},
			References: []MessageId{"Message-Id-0@example.com"},
			Comments:   "Message Header Comment",
			Keywords:   []string{"Keyword 1", "Keyword 2"},
			ResentDate: expectedDate,
			ResentFrom: []*mail.Address{
				{
					Name:    "郵便アリス",
					Address: "alice.yuubin@example.net",
				},
				{
					Name:    "郵便アリス",
					Address: "alice.yuubin@example.com",
				},
			},
			ResentSender: &mail.Address{
				Name:    "郵便アリス",
				Address: "alice.yuubin@example.net",
			},
			ResentTo: []*mail.Address{
				{
					Name:    "郵便ボブ",
					Address: "bob.yuubin@example.net",
				},
				{
					Name:    "郵便キャロル",
					Address: "carol.yuubin@example.net",
				},
			},
			ResentCc: []*mail.Address{
				{
					Name:    "郵便ダン",
					Address: "dan.yuubin@example.net",
				},
				{
					Name:    "郵便イーブ",
					Address: "eve.yubin@example.net",
				},
			},
			ResentBcc: []*mail.Address{
				{
					Name:    "郵便フランク",
					Address: "frank.yuubin@example.net",
				},
				{
					Name:    "郵便グレイス",
					Address: "grace.yubin@example.net",
				},
			},
			ResentMessageID: "Message-Id-1@example.net",
			ContentType: ContentTypeHeader{
				ContentType: "text/plain",
				Params: map[string]string{
					"charset": "utf-8",
				},
			},
			ExtraHeaders: map[string][]string{
				"X-Clacks-Overhead": {"GNU Terry Pratchett"},
				"X-Script/function/\t !\"#$%&'()*+,-./;<=>?@[\\]^_`{|}~": {
					"TEST VALUE 1\t !\"#$%&'()*+,-./0123456789:;<=>?@ABCDEFGHIJKLMNOPQRSTUVWXYZ[\\]^_` abcdefghijklmnopqrstuvwxyz{|}~",
					"TEST VALUE 2\t !\"#$%&'()*+,-./0123456789:;<=>?@ABCDEFGHIJKLMNOPQRSTUVWXYZ[\\]^_` abcdefghijklmnopqrstuvwxyz{|}~",
				},
			},
		},
		Text: `色は匂えど
散りぬるを
我が世誰ぞ
常ならん
有為の奥山
今日越えて
浅き夢見じ
酔いもせず。

Iro wa nioedo / Chirinuru o / Wa ga yo tare zo / Tsune naran / Ui no okuyama / Kyo koete Asaki yume miji / Yoi mo sezu.

とりなくこゑすゆめさませみよあけわたるひんかしをそらいろはえておきつへにほふねむれゐぬもやのうち。

天地星空
山川峰谷
雲霧室苔
人犬上末
硫黄猿生ふ為よ
榎の枝を馴れ居て。

田居に出で菜摘むわれをぞ君召すと求食り追ひゆく山城の打酔へる子ら藻葉干せよえ舟繋けぬ。`,
	}

	testEmailFromFile(t, fp, expectedEmail)
}

func TestParseEmailJapanesePlaintextUtf8OverBase64(t *testing.T) {
	fp := "tests/test_japanese_plaintext_utf-8_over_base64.txt"
	tz, _ := time.LoadLocation("Asia/Tokyo")
	expectedDate, _ := time.Parse(
		time.RFC1123Z+" (MST)",
		time.Date(2019, time.April, 1, 7, 55, 0, 0, tz).Format(time.RFC1123Z+" (MST)"))
	expectedEmail := Email{
		Headers: Headers{
			Date:    expectedDate,
			Subject: "📧 Test いろは歌",
			ReplyTo: []*mail.Address{
				{
					Name:    "郵便アリス",
					Address: "alice.yuubin@example.net",
				},
			},
			Sender: &mail.Address{
				Name:    "郵便アリス",
				Address: "alice.yuubin@example.com",
			},
			From: []*mail.Address{
				{
					Name:    "郵便アリス",
					Address: "alice.yuubin@example.com",
				},
				{
					Name:    "郵便アリス",
					Address: "alice.yuubin@example.net",
				},
			},
			To: []*mail.Address{
				{
					Name:    "郵便ボブ",
					Address: "bob.yuubin@example.com",
				},
				{
					Name:    "郵便キャロル",
					Address: "carol.yuubin@example.com",
				},
			},
			Cc: []*mail.Address{
				{
					Name:    "郵便ダン",
					Address: "dan.yuubin@example.com",
				},
				{
					Name:    "郵便イーブ",
					Address: "eve.yubin@example.com",
				},
			},
			Bcc: []*mail.Address{
				{
					Name:    "郵便フランク",
					Address: "frank.yuubin@example.com",
				},
				{
					Name:    "郵便グレイス",
					Address: "grace.yubin@example.com",
				},
			},
			MessageID:  "Message-Id-1@example.com",
			InReplyTo:  []MessageId{"Message-Id-0@example.com"},
			References: []MessageId{"Message-Id-0@example.com"},
			Comments:   "Message Header Comment",
			Keywords:   []string{"Keyword 1", "Keyword 2"},
			ResentDate: expectedDate,
			ResentFrom: []*mail.Address{
				{
					Name:    "郵便アリス",
					Address: "alice.yuubin@example.net",
				},
				{
					Name:    "郵便アリス",
					Address: "alice.yuubin@example.com",
				},
			},
			ResentSender: &mail.Address{
				Name:    "郵便アリス",
				Address: "alice.yuubin@example.net",
			},
			ResentTo: []*mail.Address{
				{
					Name:    "郵便ボブ",
					Address: "bob.yuubin@example.net",
				},
				{
					Name:    "郵便キャロル",
					Address: "carol.yuubin@example.net",
				},
			},
			ResentCc: []*mail.Address{
				{
					Name:    "郵便ダン",
					Address: "dan.yuubin@example.net",
				},
				{
					Name:    "郵便イーブ",
					Address: "eve.yubin@example.net",
				},
			},
			ResentBcc: []*mail.Address{
				{
					Name:    "郵便フランク",
					Address: "frank.yuubin@example.net",
				},
				{
					Name:    "郵便グレイス",
					Address: "grace.yubin@example.net",
				},
			},
			ResentMessageID: "Message-Id-1@example.net",
			ContentType: ContentTypeHeader{
				ContentType: "text/plain",
				Params: map[string]string{
					"charset": "utf-8",
				},
			},
			ExtraHeaders: map[string][]string{
				"X-Clacks-Overhead": {"GNU Terry Pratchett"},
				"X-Script/function/\t !\"#$%&'()*+,-./;<=>?@[\\]^_`{|}~": {
					"TEST VALUE 1\t !\"#$%&'()*+,-./0123456789:;<=>?@ABCDEFGHIJKLMNOPQRSTUVWXYZ[\\]^_` abcdefghijklmnopqrstuvwxyz{|}~",
					"TEST VALUE 2\t !\"#$%&'()*+,-./0123456789:;<=>?@ABCDEFGHIJKLMNOPQRSTUVWXYZ[\\]^_` abcdefghijklmnopqrstuvwxyz{|}~",
				},
			},
		},
		Text: `色は匂えど
散りぬるを
我が世誰ぞ
常ならん
有為の奥山
今日越えて
浅き夢見じ
酔いもせず。

Iro wa nioedo / Chirinuru o / Wa ga yo tare zo / Tsune naran / Ui no okuyama / Kyo koete Asaki yume miji / Yoi mo sezu.

とりなくこゑすゆめさませみよあけわたるひんかしをそらいろはえておきつへにほふねむれゐぬもやのうち。

天地星空
山川峰谷
雲霧室苔
人犬上末
硫黄猿生ふ為よ
榎の枝を馴れ居て。

田居に出で菜摘むわれをぞ君召すと求食り追ひゆく山城の打酔へる子ら藻葉干せよえ舟繋けぬ。`,
	}

	testEmailFromFile(t, fp, expectedEmail)
}

func TestParseEmailJapanesePlaintextUtf8OverQuotedprintable(t *testing.T) {
	fp := "tests/test_japanese_plaintext_utf-8_over_quoted-printable.txt"
	tz, _ := time.LoadLocation("Asia/Tokyo")
	expectedDate, _ := time.Parse(
		time.RFC1123Z+" (MST)",
		time.Date(2019, time.April, 1, 7, 55, 0, 0, tz).Format(time.RFC1123Z+" (MST)"))
	expectedEmail := Email{
		Headers: Headers{
			Date:    expectedDate,
			Subject: "📧 Test いろは歌",
			ReplyTo: []*mail.Address{
				{
					Name:    "郵便アリス",
					Address: "alice.yuubin@example.net",
				},
			},
			Sender: &mail.Address{
				Name:    "郵便アリス",
				Address: "alice.yuubin@example.com",
			},
			From: []*mail.Address{
				{
					Name:    "郵便アリス",
					Address: "alice.yuubin@example.com",
				},
				{
					Name:    "郵便アリス",
					Address: "alice.yuubin@example.net",
				},
			},
			To: []*mail.Address{
				{
					Name:    "郵便ボブ",
					Address: "bob.yuubin@example.com",
				},
				{
					Name:    "郵便キャロル",
					Address: "carol.yuubin@example.com",
				},
			},
			Cc: []*mail.Address{
				{
					Name:    "郵便ダン",
					Address: "dan.yuubin@example.com",
				},
				{
					Name:    "郵便イーブ",
					Address: "eve.yubin@example.com",
				},
			},
			Bcc: []*mail.Address{
				{
					Name:    "郵便フランク",
					Address: "frank.yuubin@example.com",
				},
				{
					Name:    "郵便グレイス",
					Address: "grace.yubin@example.com",
				},
			},
			MessageID:  "Message-Id-1@example.com",
			InReplyTo:  []MessageId{"Message-Id-0@example.com"},
			References: []MessageId{"Message-Id-0@example.com"},
			Comments:   "Message Header Comment",
			Keywords:   []string{"Keyword 1", "Keyword 2"},
			ResentDate: expectedDate,
			ResentFrom: []*mail.Address{
				{
					Name:    "郵便アリス",
					Address: "alice.yuubin@example.net",
				},
				{
					Name:    "郵便アリス",
					Address: "alice.yuubin@example.com",
				},
			},
			ResentSender: &mail.Address{
				Name:    "郵便アリス",
				Address: "alice.yuubin@example.net",
			},
			ResentTo: []*mail.Address{
				{
					Name:    "郵便ボブ",
					Address: "bob.yuubin@example.net",
				},
				{
					Name:    "郵便キャロル",
					Address: "carol.yuubin@example.net",
				},
			},
			ResentCc: []*mail.Address{
				{
					Name:    "郵便ダン",
					Address: "dan.yuubin@example.net",
				},
				{
					Name:    "郵便イーブ",
					Address: "eve.yubin@example.net",
				},
			},
			ResentBcc: []*mail.Address{
				{
					Name:    "郵便フランク",
					Address: "frank.yuubin@example.net",
				},
				{
					Name:    "郵便グレイス",
					Address: "grace.yubin@example.net",
				},
			},
			ResentMessageID: "Message-Id-1@example.net",
			ContentType: ContentTypeHeader{
				ContentType: "text/plain",
				Params: map[string]string{
					"charset": "utf-8",
				},
			},
			ExtraHeaders: map[string][]string{
				"X-Clacks-Overhead": {"GNU Terry Pratchett"},
				"X-Script/function/\t !\"#$%&'()*+,-./;<=>?@[\\]^_`{|}~": {
					"TEST VALUE 1\t !\"#$%&'()*+,-./0123456789:;<=>?@ABCDEFGHIJKLMNOPQRSTUVWXYZ[\\]^_` abcdefghijklmnopqrstuvwxyz{|}~",
					"TEST VALUE 2\t !\"#$%&'()*+,-./0123456789:;<=>?@ABCDEFGHIJKLMNOPQRSTUVWXYZ[\\]^_` abcdefghijklmnopqrstuvwxyz{|}~",
				},
			},
		},
		Text: `色は匂えど
散りぬるを
我が世誰ぞ
常ならん
有為の奥山
今日越えて
浅き夢見じ
酔いもせず。

Iro wa nioedo / Chirinuru o / Wa ga yo tare zo / Tsune naran / Ui no okuyama / Kyo koete Asaki yume miji / Yoi mo sezu.

とりなくこゑすゆめさませみよあけわたるひんかしをそらいろはえておきつへにほふねむれゐぬもやのうち。

天地星空
山川峰谷
雲霧室苔
人犬上末
硫黄猿生ふ為よ
榎の枝を馴れ居て。

田居に出で菜摘むわれをぞ君召すと求食り追ひゆく山城の打酔へる子ら藻葉干せよえ舟繋けぬ。`,
	}

	testEmailFromFile(t, fp, expectedEmail)
}

func TestParseEmailJapanesePlaintextIso2022jpOver7bit(t *testing.T) {
	fp := "tests/test_japanese_plaintext_iso-2022-jp_over_7bit.txt"
	tz, _ := time.LoadLocation("Asia/Tokyo")
	expectedDate, _ := time.Parse(
		time.RFC1123Z+" (MST)",
		time.Date(2019, time.April, 1, 7, 55, 0, 0, tz).Format(time.RFC1123Z+" (MST)"))
	expectedEmail := Email{
		Headers: Headers{
			Date:    expectedDate,
			Subject: "📧 Test いろは歌",
			ReplyTo: []*mail.Address{
				{
					Name:    "郵便アリス",
					Address: "alice.yuubin@example.net",
				},
			},
			Sender: &mail.Address{
				Name:    "郵便アリス",
				Address: "alice.yuubin@example.com",
			},
			From: []*mail.Address{
				{
					Name:    "郵便アリス",
					Address: "alice.yuubin@example.com",
				},
				{
					Name:    "郵便アリス",
					Address: "alice.yuubin@example.net",
				},
			},
			To: []*mail.Address{
				{
					Name:    "郵便ボブ",
					Address: "bob.yuubin@example.com",
				},
				{
					Name:    "郵便キャロル",
					Address: "carol.yuubin@example.com",
				},
			},
			Cc: []*mail.Address{
				{
					Name:    "郵便ダン",
					Address: "dan.yuubin@example.com",
				},
				{
					Name:    "郵便イーブ",
					Address: "eve.yubin@example.com",
				},
			},
			Bcc: []*mail.Address{
				{
					Name:    "郵便フランク",
					Address: "frank.yuubin@example.com",
				},
				{
					Name:    "郵便グレイス",
					Address: "grace.yubin@example.com",
				},
			},
			MessageID:  "Message-Id-1@example.com",
			InReplyTo:  []MessageId{"Message-Id-0@example.com"},
			References: []MessageId{"Message-Id-0@example.com"},
			Comments:   "Message Header Comment",
			Keywords:   []string{"Keyword 1", "Keyword 2"},
			ResentDate: expectedDate,
			ResentFrom: []*mail.Address{
				{
					Name:    "郵便アリス",
					Address: "alice.yuubin@example.net",
				},
				{
					Name:    "郵便アリス",
					Address: "alice.yuubin@example.com",
				},
			},
			ResentSender: &mail.Address{
				Name:    "郵便アリス",
				Address: "alice.yuubin@example.net",
			},
			ResentTo: []*mail.Address{
				{
					Name:    "郵便ボブ",
					Address: "bob.yuubin@example.net",
				},
				{
					Name:    "郵便キャロル",
					Address: "carol.yuubin@example.net",
				},
			},
			ResentCc: []*mail.Address{
				{
					Name:    "郵便ダン",
					Address: "dan.yuubin@example.net",
				},
				{
					Name:    "郵便イーブ",
					Address: "eve.yubin@example.net",
				},
			},
			ResentBcc: []*mail.Address{
				{
					Name:    "郵便フランク",
					Address: "frank.yuubin@example.net",
				},
				{
					Name:    "郵便グレイス",
					Address: "grace.yubin@example.net",
				},
			},
			ResentMessageID: "Message-Id-1@example.net",
			ContentType: ContentTypeHeader{
				ContentType: "text/plain",
				Params: map[string]string{
					"charset": "iso-2022-jp",
				},
			},
			ExtraHeaders: map[string][]string{
				"X-Clacks-Overhead": {"GNU Terry Pratchett"},
				"X-Script/function/\t !\"#$%&'()*+,-./;<=>?@[\\]^_`{|}~": {
					"TEST VALUE 1\t !\"#$%&'()*+,-./0123456789:;<=>?@ABCDEFGHIJKLMNOPQRSTUVWXYZ[\\]^_` abcdefghijklmnopqrstuvwxyz{|}~",
					"TEST VALUE 2\t !\"#$%&'()*+,-./0123456789:;<=>?@ABCDEFGHIJKLMNOPQRSTUVWXYZ[\\]^_` abcdefghijklmnopqrstuvwxyz{|}~",
				},
			},
		},
		Text: `色は匂えど
散りぬるを
我が世誰ぞ
常ならん
有為の奥山
今日越えて
浅き夢見じ
酔いもせず。

Iro wa nioedo / Chirinuru o / Wa ga yo tare zo / Tsune naran / Ui no okuyama / Kyo koete Asaki yume miji / Yoi mo sezu.

とりなくこゑすゆめさませみよあけわたるひんかしをそらいろはえておきつへにほふねむれゐぬもやのうち。

天地星空
山川峰谷
雲霧室苔
人犬上末
硫黄猿生ふ為よ
榎の枝を馴れ居て。

田居に出で菜摘むわれをぞ君召すと求食り追ひゆく山城の打酔へる子ら藻葉干せよえ舟繋けぬ。`,
	}

	testEmailFromFile(t, fp, expectedEmail)
}

func TestParseEmailJapanesePlaintextIso2022jpOverBase64(t *testing.T) {
	fp := "tests/test_japanese_plaintext_iso-2022-jp_over_base64.txt"
	tz, _ := time.LoadLocation("Asia/Tokyo")
	expectedDate, _ := time.Parse(
		time.RFC1123Z+" (MST)",
		time.Date(2019, time.April, 1, 7, 55, 0, 0, tz).Format(time.RFC1123Z+" (MST)"))
	expectedEmail := Email{
		Headers: Headers{
			Date:    expectedDate,
			Subject: "📧 Test いろは歌",
			ReplyTo: []*mail.Address{
				{
					Name:    "郵便アリス",
					Address: "alice.yuubin@example.net",
				},
			},
			Sender: &mail.Address{
				Name:    "郵便アリス",
				Address: "alice.yuubin@example.com",
			},
			From: []*mail.Address{
				{
					Name:    "郵便アリス",
					Address: "alice.yuubin@example.com",
				},
				{
					Name:    "郵便アリス",
					Address: "alice.yuubin@example.net",
				},
			},
			To: []*mail.Address{
				{
					Name:    "郵便ボブ",
					Address: "bob.yuubin@example.com",
				},
				{
					Name:    "郵便キャロル",
					Address: "carol.yuubin@example.com",
				},
			},
			Cc: []*mail.Address{
				{
					Name:    "郵便ダン",
					Address: "dan.yuubin@example.com",
				},
				{
					Name:    "郵便イーブ",
					Address: "eve.yubin@example.com",
				},
			},
			Bcc: []*mail.Address{
				{
					Name:    "郵便フランク",
					Address: "frank.yuubin@example.com",
				},
				{
					Name:    "郵便グレイス",
					Address: "grace.yubin@example.com",
				},
			},
			MessageID:  "Message-Id-1@example.com",
			InReplyTo:  []MessageId{"Message-Id-0@example.com"},
			References: []MessageId{"Message-Id-0@example.com"},
			Comments:   "Message Header Comment",
			Keywords:   []string{"Keyword 1", "Keyword 2"},
			ResentDate: expectedDate,
			ResentFrom: []*mail.Address{
				{
					Name:    "郵便アリス",
					Address: "alice.yuubin@example.net",
				},
				{
					Name:    "郵便アリス",
					Address: "alice.yuubin@example.com",
				},
			},
			ResentSender: &mail.Address{
				Name:    "郵便アリス",
				Address: "alice.yuubin@example.net",
			},
			ResentTo: []*mail.Address{
				{
					Name:    "郵便ボブ",
					Address: "bob.yuubin@example.net",
				},
				{
					Name:    "郵便キャロル",
					Address: "carol.yuubin@example.net",
				},
			},
			ResentCc: []*mail.Address{
				{
					Name:    "郵便ダン",
					Address: "dan.yuubin@example.net",
				},
				{
					Name:    "郵便イーブ",
					Address: "eve.yubin@example.net",
				},
			},
			ResentBcc: []*mail.Address{
				{
					Name:    "郵便フランク",
					Address: "frank.yuubin@example.net",
				},
				{
					Name:    "郵便グレイス",
					Address: "grace.yubin@example.net",
				},
			},
			ResentMessageID: "Message-Id-1@example.net",
			ContentType: ContentTypeHeader{
				ContentType: "text/plain",
				Params: map[string]string{
					"charset": "iso-2022-jp",
				},
			},
			ExtraHeaders: map[string][]string{
				"X-Clacks-Overhead": {"GNU Terry Pratchett"},
				"X-Script/function/\t !\"#$%&'()*+,-./;<=>?@[\\]^_`{|}~": {
					"TEST VALUE 1\t !\"#$%&'()*+,-./0123456789:;<=>?@ABCDEFGHIJKLMNOPQRSTUVWXYZ[\\]^_` abcdefghijklmnopqrstuvwxyz{|}~",
					"TEST VALUE 2\t !\"#$%&'()*+,-./0123456789:;<=>?@ABCDEFGHIJKLMNOPQRSTUVWXYZ[\\]^_` abcdefghijklmnopqrstuvwxyz{|}~",
				},
			},
		},
		Text: `色は匂えど
散りぬるを
我が世誰ぞ
常ならん
有為の奥山
今日越えて
浅き夢見じ
酔いもせず。

Iro wa nioedo / Chirinuru o / Wa ga yo tare zo / Tsune naran / Ui no okuyama / Kyo koete Asaki yume miji / Yoi mo sezu.

とりなくこゑすゆめさませみよあけわたるひんかしをそらいろはえておきつへにほふねむれゐぬもやのうち。

天地星空
山川峰谷
雲霧室苔
人犬上末
硫黄猿生ふ為よ
榎の枝を馴れ居て。

田居に出で菜摘むわれをぞ君召すと求食り追ひゆく山城の打酔へる子ら藻葉干せよえ舟繋けぬ。`,
	}

	testEmailFromFile(t, fp, expectedEmail)
}

func TestParseEmailJapanesePlaintextIso2022jpOverQuotedprintable(t *testing.T) {
	fp := "tests/test_japanese_plaintext_iso-2022-jp_over_quoted-printable.txt"
	tz, _ := time.LoadLocation("Asia/Tokyo")
	expectedDate, _ := time.Parse(
		time.RFC1123Z+" (MST)",
		time.Date(2019, time.April, 1, 7, 55, 0, 0, tz).Format(time.RFC1123Z+" (MST)"))
	expectedEmail := Email{
		Headers: Headers{
			Date:    expectedDate,
			Subject: "📧 Test いろは歌",
			ReplyTo: []*mail.Address{
				{
					Name:    "郵便アリス",
					Address: "alice.yuubin@example.net",
				},
			},
			Sender: &mail.Address{
				Name:    "郵便アリス",
				Address: "alice.yuubin@example.com",
			},
			From: []*mail.Address{
				{
					Name:    "郵便アリス",
					Address: "alice.yuubin@example.com",
				},
				{
					Name:    "郵便アリス",
					Address: "alice.yuubin@example.net",
				},
			},
			To: []*mail.Address{
				{
					Name:    "郵便ボブ",
					Address: "bob.yuubin@example.com",
				},
				{
					Name:    "郵便キャロル",
					Address: "carol.yuubin@example.com",
				},
			},
			Cc: []*mail.Address{
				{
					Name:    "郵便ダン",
					Address: "dan.yuubin@example.com",
				},
				{
					Name:    "郵便イーブ",
					Address: "eve.yubin@example.com",
				},
			},
			Bcc: []*mail.Address{
				{
					Name:    "郵便フランク",
					Address: "frank.yuubin@example.com",
				},
				{
					Name:    "郵便グレイス",
					Address: "grace.yubin@example.com",
				},
			},
			MessageID:  "Message-Id-1@example.com",
			InReplyTo:  []MessageId{"Message-Id-0@example.com"},
			References: []MessageId{"Message-Id-0@example.com"},
			Comments:   "Message Header Comment",
			Keywords:   []string{"Keyword 1", "Keyword 2"},
			ResentDate: expectedDate,
			ResentFrom: []*mail.Address{
				{
					Name:    "郵便アリス",
					Address: "alice.yuubin@example.net",
				},
				{
					Name:    "郵便アリス",
					Address: "alice.yuubin@example.com",
				},
			},
			ResentSender: &mail.Address{
				Name:    "郵便アリス",
				Address: "alice.yuubin@example.net",
			},
			ResentTo: []*mail.Address{
				{
					Name:    "郵便ボブ",
					Address: "bob.yuubin@example.net",
				},
				{
					Name:    "郵便キャロル",
					Address: "carol.yuubin@example.net",
				},
			},
			ResentCc: []*mail.Address{
				{
					Name:    "郵便ダン",
					Address: "dan.yuubin@example.net",
				},
				{
					Name:    "郵便イーブ",
					Address: "eve.yubin@example.net",
				},
			},
			ResentBcc: []*mail.Address{
				{
					Name:    "郵便フランク",
					Address: "frank.yuubin@example.net",
				},
				{
					Name:    "郵便グレイス",
					Address: "grace.yubin@example.net",
				},
			},
			ResentMessageID: "Message-Id-1@example.net",
			ContentType: ContentTypeHeader{
				ContentType: "text/plain",
				Params: map[string]string{
					"charset": "iso-2022-jp",
				},
			},
			ExtraHeaders: map[string][]string{
				"X-Clacks-Overhead": {"GNU Terry Pratchett"},
				"X-Script/function/\t !\"#$%&'()*+,-./;<=>?@[\\]^_`{|}~": {
					"TEST VALUE 1\t !\"#$%&'()*+,-./0123456789:;<=>?@ABCDEFGHIJKLMNOPQRSTUVWXYZ[\\]^_` abcdefghijklmnopqrstuvwxyz{|}~",
					"TEST VALUE 2\t !\"#$%&'()*+,-./0123456789:;<=>?@ABCDEFGHIJKLMNOPQRSTUVWXYZ[\\]^_` abcdefghijklmnopqrstuvwxyz{|}~",
				},
			},
		},
		Text: `色は匂えど
散りぬるを
我が世誰ぞ
常ならん
有為の奥山
今日越えて
浅き夢見じ
酔いもせず。

Iro wa nioedo / Chirinuru o / Wa ga yo tare zo / Tsune naran / Ui no okuyama / Kyo koete Asaki yume miji / Yoi mo sezu.

とりなくこゑすゆめさませみよあけわたるひんかしをそらいろはえておきつへにほふねむれゐぬもやのうち。

天地星空
山川峰谷
雲霧室苔
人犬上末
硫黄猿生ふ為よ
榎の枝を馴れ居て。

田居に出で菜摘むわれをぞ君召すと求食り追ひゆく山城の打酔へる子ら藻葉干せよえ舟繋けぬ。`,
	}

	testEmailFromFile(t, fp, expectedEmail)
}

func TestParseEmailJapanesePlaintextEucjpOverBase64(t *testing.T) {
	fp := "tests/test_japanese_plaintext_euc-jp_over_base64.txt"
	tz, _ := time.LoadLocation("Asia/Tokyo")
	expectedDate, _ := time.Parse(
		time.RFC1123Z+" (MST)",
		time.Date(2019, time.April, 1, 7, 55, 0, 0, tz).Format(time.RFC1123Z+" (MST)"))
	expectedEmail := Email{
		Headers: Headers{
			Date:    expectedDate,
			Subject: "📧 Test いろは歌",
			ReplyTo: []*mail.Address{
				{
					Name:    "郵便アリス",
					Address: "alice.yuubin@example.net",
				},
			},
			Sender: &mail.Address{
				Name:    "郵便アリス",
				Address: "alice.yuubin@example.com",
			},
			From: []*mail.Address{
				{
					Name:    "郵便アリス",
					Address: "alice.yuubin@example.com",
				},
				{
					Name:    "郵便アリス",
					Address: "alice.yuubin@example.net",
				},
			},
			To: []*mail.Address{
				{
					Name:    "郵便ボブ",
					Address: "bob.yuubin@example.com",
				},
				{
					Name:    "郵便キャロル",
					Address: "carol.yuubin@example.com",
				},
			},
			Cc: []*mail.Address{
				{
					Name:    "郵便ダン",
					Address: "dan.yuubin@example.com",
				},
				{
					Name:    "郵便イーブ",
					Address: "eve.yubin@example.com",
				},
			},
			Bcc: []*mail.Address{
				{
					Name:    "郵便フランク",
					Address: "frank.yuubin@example.com",
				},
				{
					Name:    "郵便グレイス",
					Address: "grace.yubin@example.com",
				},
			},
			MessageID:  "Message-Id-1@example.com",
			InReplyTo:  []MessageId{"Message-Id-0@example.com"},
			References: []MessageId{"Message-Id-0@example.com"},
			Comments:   "Message Header Comment",
			Keywords:   []string{"Keyword 1", "Keyword 2"},
			ResentDate: expectedDate,
			ResentFrom: []*mail.Address{
				{
					Name:    "郵便アリス",
					Address: "alice.yuubin@example.net",
				},
				{
					Name:    "郵便アリス",
					Address: "alice.yuubin@example.com",
				},
			},
			ResentSender: &mail.Address{
				Name:    "郵便アリス",
				Address: "alice.yuubin@example.net",
			},
			ResentTo: []*mail.Address{
				{
					Name:    "郵便ボブ",
					Address: "bob.yuubin@example.net",
				},
				{
					Name:    "郵便キャロル",
					Address: "carol.yuubin@example.net",
				},
			},
			ResentCc: []*mail.Address{
				{
					Name:    "郵便ダン",
					Address: "dan.yuubin@example.net",
				},
				{
					Name:    "郵便イーブ",
					Address: "eve.yubin@example.net",
				},
			},
			ResentBcc: []*mail.Address{
				{
					Name:    "郵便フランク",
					Address: "frank.yuubin@example.net",
				},
				{
					Name:    "郵便グレイス",
					Address: "grace.yubin@example.net",
				},
			},
			ResentMessageID: "Message-Id-1@example.net",
			ContentType: ContentTypeHeader{
				ContentType: "text/plain",
				Params: map[string]string{
					"charset": "euc-jp",
				},
			},
			ExtraHeaders: map[string][]string{
				"X-Clacks-Overhead": {"GNU Terry Pratchett"},
				"X-Script/function/\t !\"#$%&'()*+,-./;<=>?@[\\]^_`{|}~": {
					"TEST VALUE 1\t !\"#$%&'()*+,-./0123456789:;<=>?@ABCDEFGHIJKLMNOPQRSTUVWXYZ[\\]^_` abcdefghijklmnopqrstuvwxyz{|}~",
					"TEST VALUE 2\t !\"#$%&'()*+,-./0123456789:;<=>?@ABCDEFGHIJKLMNOPQRSTUVWXYZ[\\]^_` abcdefghijklmnopqrstuvwxyz{|}~",
				},
			},
		},
		Text: `色は匂えど
散りぬるを
我が世誰ぞ
常ならん
有為の奥山
今日越えて
浅き夢見じ
酔いもせず。

Iro wa nioedo / Chirinuru o / Wa ga yo tare zo / Tsune naran / Ui no okuyama / Kyo koete Asaki yume miji / Yoi mo sezu.

とりなくこゑすゆめさませみよあけわたるひんかしをそらいろはえておきつへにほふねむれゐぬもやのうち。

天地星空
山川峰谷
雲霧室苔
人犬上末
硫黄猿生ふ為よ
榎の枝を馴れ居て。

田居に出で菜摘むわれをぞ君召すと求食り追ひゆく山城の打酔へる子ら藻葉干せよえ舟繋けぬ。`,
	}

	testEmailFromFile(t, fp, expectedEmail)
}

func TestParseEmailJapanesePlaintextEucjpOverQuotedprintable(t *testing.T) {
	fp := "tests/test_japanese_plaintext_euc-jp_over_quoted-printable.txt"
	tz, _ := time.LoadLocation("Asia/Tokyo")
	expectedDate, _ := time.Parse(
		time.RFC1123Z+" (MST)",
		time.Date(2019, time.April, 1, 7, 55, 0, 0, tz).Format(time.RFC1123Z+" (MST)"))
	expectedEmail := Email{
		Headers: Headers{
			Date:    expectedDate,
			Subject: "📧 Test いろは歌",
			ReplyTo: []*mail.Address{
				{
					Name:    "郵便アリス",
					Address: "alice.yuubin@example.net",
				},
			},
			Sender: &mail.Address{
				Name:    "郵便アリス",
				Address: "alice.yuubin@example.com",
			},
			From: []*mail.Address{
				{
					Name:    "郵便アリス",
					Address: "alice.yuubin@example.com",
				},
				{
					Name:    "郵便アリス",
					Address: "alice.yuubin@example.net",
				},
			},
			To: []*mail.Address{
				{
					Name:    "郵便ボブ",
					Address: "bob.yuubin@example.com",
				},
				{
					Name:    "郵便キャロル",
					Address: "carol.yuubin@example.com",
				},
			},
			Cc: []*mail.Address{
				{
					Name:    "郵便ダン",
					Address: "dan.yuubin@example.com",
				},
				{
					Name:    "郵便イーブ",
					Address: "eve.yubin@example.com",
				},
			},
			Bcc: []*mail.Address{
				{
					Name:    "郵便フランク",
					Address: "frank.yuubin@example.com",
				},
				{
					Name:    "郵便グレイス",
					Address: "grace.yubin@example.com",
				},
			},
			MessageID:  "Message-Id-1@example.com",
			InReplyTo:  []MessageId{"Message-Id-0@example.com"},
			References: []MessageId{"Message-Id-0@example.com"},
			Comments:   "Message Header Comment",
			Keywords:   []string{"Keyword 1", "Keyword 2"},
			ResentDate: expectedDate,
			ResentFrom: []*mail.Address{
				{
					Name:    "郵便アリス",
					Address: "alice.yuubin@example.net",
				},
				{
					Name:    "郵便アリス",
					Address: "alice.yuubin@example.com",
				},
			},
			ResentSender: &mail.Address{
				Name:    "郵便アリス",
				Address: "alice.yuubin@example.net",
			},
			ResentTo: []*mail.Address{
				{
					Name:    "郵便ボブ",
					Address: "bob.yuubin@example.net",
				},
				{
					Name:    "郵便キャロル",
					Address: "carol.yuubin@example.net",
				},
			},
			ResentCc: []*mail.Address{
				{
					Name:    "郵便ダン",
					Address: "dan.yuubin@example.net",
				},
				{
					Name:    "郵便イーブ",
					Address: "eve.yubin@example.net",
				},
			},
			ResentBcc: []*mail.Address{
				{
					Name:    "郵便フランク",
					Address: "frank.yuubin@example.net",
				},
				{
					Name:    "郵便グレイス",
					Address: "grace.yubin@example.net",
				},
			},
			ResentMessageID: "Message-Id-1@example.net",
			ContentType: ContentTypeHeader{
				ContentType: "text/plain",
				Params: map[string]string{
					"charset": "euc-jp",
				},
			},
			ExtraHeaders: map[string][]string{
				"X-Clacks-Overhead": {"GNU Terry Pratchett"},
				"X-Script/function/\t !\"#$%&'()*+,-./;<=>?@[\\]^_`{|}~": {
					"TEST VALUE 1\t !\"#$%&'()*+,-./0123456789:;<=>?@ABCDEFGHIJKLMNOPQRSTUVWXYZ[\\]^_` abcdefghijklmnopqrstuvwxyz{|}~",
					"TEST VALUE 2\t !\"#$%&'()*+,-./0123456789:;<=>?@ABCDEFGHIJKLMNOPQRSTUVWXYZ[\\]^_` abcdefghijklmnopqrstuvwxyz{|}~",
				},
			},
		},
		Text: `色は匂えど
散りぬるを
我が世誰ぞ
常ならん
有為の奥山
今日越えて
浅き夢見じ
酔いもせず。

Iro wa nioedo / Chirinuru o / Wa ga yo tare zo / Tsune naran / Ui no okuyama / Kyo koete Asaki yume miji / Yoi mo sezu.

とりなくこゑすゆめさませみよあけわたるひんかしをそらいろはえておきつへにほふねむれゐぬもやのうち。

天地星空
山川峰谷
雲霧室苔
人犬上末
硫黄猿生ふ為よ
榎の枝を馴れ居て。

田居に出で菜摘むわれをぞ君召すと求食り追ひゆく山城の打酔へる子ら藻葉干せよえ舟繋けぬ。`,
	}

	testEmailFromFile(t, fp, expectedEmail)
}

func TestParseEmailJapaneseMultipartRelatedUtf8Over7bit(t *testing.T) {
	fp := "tests/test_japanese_multipart_related_utf-8_over_7bit.txt"
	tz, _ := time.LoadLocation("Asia/Tokyo")
	expectedDate, _ := time.Parse(
		time.RFC1123Z+" (MST)",
		time.Date(2019, time.April, 1, 7, 55, 0, 0, tz).Format(time.RFC1123Z+" (MST)"))
	expectedEmail := Email{
		Headers: Headers{
			Date:    expectedDate,
			Subject: "📧 Test いろは歌",
			ReplyTo: []*mail.Address{
				{
					Name:    "郵便アリス",
					Address: "alice.yuubin@example.net",
				},
			},
			Sender: &mail.Address{
				Name:    "郵便アリス",
				Address: "alice.yuubin@example.com",
			},
			From: []*mail.Address{
				{
					Name:    "郵便アリス",
					Address: "alice.yuubin@example.com",
				},
				{
					Name:    "郵便アリス",
					Address: "alice.yuubin@example.net",
				},
			},
			To: []*mail.Address{
				{
					Name:    "郵便ボブ",
					Address: "bob.yuubin@example.com",
				},
				{
					Name:    "郵便キャロル",
					Address: "carol.yuubin@example.com",
				},
			},
			Cc: []*mail.Address{
				{
					Name:    "郵便ダン",
					Address: "dan.yuubin@example.com",
				},
				{
					Name:    "郵便イーブ",
					Address: "eve.yubin@example.com",
				},
			},
			Bcc: []*mail.Address{
				{
					Name:    "郵便フランク",
					Address: "frank.yuubin@example.com",
				},
				{
					Name:    "郵便グレイス",
					Address: "grace.yubin@example.com",
				},
			},
			MessageID:  "Message-Id-1@example.com",
			InReplyTo:  []MessageId{"Message-Id-0@example.com"},
			References: []MessageId{"Message-Id-0@example.com"},
			Comments:   "Message Header Comment",
			Keywords:   []string{"Keyword 1", "Keyword 2"},
			ResentDate: expectedDate,
			ResentFrom: []*mail.Address{
				{
					Name:    "郵便アリス",
					Address: "alice.yuubin@example.net",
				},
				{
					Name:    "郵便アリス",
					Address: "alice.yuubin@example.com",
				},
			},
			ResentSender: &mail.Address{
				Name:    "郵便アリス",
				Address: "alice.yuubin@example.net",
			},
			ResentTo: []*mail.Address{
				{
					Name:    "郵便ボブ",
					Address: "bob.yuubin@example.net",
				},
				{
					Name:    "郵便キャロル",
					Address: "carol.yuubin@example.net",
				},
			},
			ResentCc: []*mail.Address{
				{
					Name:    "郵便ダン",
					Address: "dan.yuubin@example.net",
				},
				{
					Name:    "郵便イーブ",
					Address: "eve.yubin@example.net",
				},
			},
			ResentBcc: []*mail.Address{
				{
					Name:    "郵便フランク",
					Address: "frank.yuubin@example.net",
				},
				{
					Name:    "郵便グレイス",
					Address: "grace.yubin@example.net",
				},
			},
			ResentMessageID: "Message-Id-1@example.net",
			ContentType: ContentTypeHeader{
				ContentType: "multipart/related",
				Params: map[string]string{
					"boundary": "RelatedBoundaryString",
					"charset":  "utf-8",
				},
			},
			ExtraHeaders: map[string][]string{
				"X-Clacks-Overhead": {"GNU Terry Pratchett"},
				"X-Script/function/\t !\"#$%&'()*+,-./;<=>?@[\\]^_`{|}~": {
					"TEST VALUE 1\t !\"#$%&'()*+,-./0123456789:;<=>?@ABCDEFGHIJKLMNOPQRSTUVWXYZ[\\]^_` abcdefghijklmnopqrstuvwxyz{|}~",
					"TEST VALUE 2\t !\"#$%&'()*+,-./0123456789:;<=>?@ABCDEFGHIJKLMNOPQRSTUVWXYZ[\\]^_` abcdefghijklmnopqrstuvwxyz{|}~",
				},
			},
		},
		Text: `色は匂えど
散りぬるを
我が世誰ぞ
常ならん
有為の奥山
今日越えて
浅き夢見じ
酔いもせず。

Iro wa nioedo / Chirinuru o / Wa ga yo tare zo / Tsune naran / Ui no okuyama / Kyo koete Asaki yume miji / Yoi mo sezu.

とりなくこゑすゆめさませみよあけわたるひんかしをそらいろはえておきつへにほふねむれゐぬもやのうち。

天地星空
山川峰谷
雲霧室苔
人犬上末
硫黄猿生ふ為よ
榎の枝を馴れ居て。

田居に出で菜摘むわれをぞ君召すと求食り追ひゆく山城の打酔へる子ら藻葉干せよえ舟繋けぬ。`,
		EnrichedText: `<bold>色は匂えど</bold>
<italic>散りぬるを</italic>
<fixed>我が世誰ぞ</fixed>
<underline>常ならん</underline>
有為の奥山
今日越えて
浅き夢見じ
酔いもせず。

Iro wa nioedo / Chirinuru o / Wa ga yo tare zo / Tsune naran / Ui no okuyama / Kyo koete Asaki yume miji / Yoi mo sezu.

とりなくこゑすゆめさませみよあけわたるひんかしをそらいろはえておきつへにほふねむれゐぬもやのうち。

天地星空
山川峰谷
雲霧室苔
人犬上末
硫黄猿生ふ為よ
榎の枝を馴れ居て。

田居に出で菜摘むわれをぞ君召すと求食り追ひゆく山城の打酔へる子ら藻葉干せよえ舟繋けぬ。`,
		HTML: `<html>
<div dir="ltr">
<p>色は匂えど<br />
散りぬるを<br />
我が世誰ぞ<br />
常ならん<br />
有為の奥山<br />
今日越えて<br />
浅き夢見じ<br />
酔いもせず。</p>

<p>Iro wa nioedo / Chirinuru o / Wa ga yo tare zo / Tsune naran / Ui no okuyama / Kyo koete Asaki yume miji / Yoi mo sezu.</p>

<p>とりなくこゑすゆめさませみよあけわたるひんかしをそらいろはえておきつへにほふねむれゐぬもやのうち。</p>

<p>天地星空<br />
山川峰谷<br />
雲霧室苔<br />
人犬上末<br />
硫黄猿生ふ為よ<br />
榎の枝を馴れ居て。</p>

<p>田居に出で菜摘むわれをぞ君召すと求食り追ひゆく山城の打酔へる子ら藻葉干せよえ舟繋けぬ。</p>
</div>
</html>`,
		InlineFiles: []InlineFile{
			{
				ContentID: "inline-jpg-image.jpg@example.com",
				ContentType: ContentTypeHeader{
					ContentType: "image/jpeg",
					Params: map[string]string{
						"name": "inline-jpg-image-name.jpg",
					},
				},
				ContentDisposition: ContentDispositionHeader{
					ContentDisposition: inline,
					Params: map[string]string{
						"filename": "inline-jpg-image-filename.jpg",
					},
				},
				Data: []byte{255, 216, 255, 219, 0, 67, 0, 3, 2, 2, 2, 2, 2, 3, 2, 2, 2, 3, 3, 3, 3, 4,
					6, 4, 4, 4, 4, 4, 8, 6, 6, 5, 6, 9, 8, 10, 10, 9, 8, 9, 9, 10, 12, 15, 12, 10, 11, 14, 11, 9, 9,
					13, 17, 13, 14, 15, 16, 16, 17, 16, 10, 12, 18, 19, 18, 16, 19, 15, 16, 16, 16, 255, 201, 0, 11, 8,
					0, 1, 0, 1, 1, 1, 17, 0, 255, 204, 0, 6, 0, 16, 16, 5, 255, 218, 0, 8, 1, 1, 0, 0, 63, 0, 210, 207,
					32, 255, 217},
			},
		},
	}

	testEmailFromFile(t, fp, expectedEmail)
}
func TestParseEmailJapaneseMultipartRelatedUtf8OverBase64(t *testing.T) {
	fp := "tests/test_japanese_multipart_related_utf-8_over_base64.txt"
	tz, _ := time.LoadLocation("Asia/Tokyo")
	expectedDate, _ := time.Parse(
		time.RFC1123Z+" (MST)",
		time.Date(2019, time.April, 1, 7, 55, 0, 0, tz).Format(time.RFC1123Z+" (MST)"))
	expectedEmail := Email{
		Headers: Headers{
			Date:    expectedDate,
			Subject: "📧 Test いろは歌",
			ReplyTo: []*mail.Address{
				{
					Name:    "郵便アリス",
					Address: "alice.yuubin@example.net",
				},
			},
			Sender: &mail.Address{
				Name:    "郵便アリス",
				Address: "alice.yuubin@example.com",
			},
			From: []*mail.Address{
				{
					Name:    "郵便アリス",
					Address: "alice.yuubin@example.com",
				},
				{
					Name:    "郵便アリス",
					Address: "alice.yuubin@example.net",
				},
			},
			To: []*mail.Address{
				{
					Name:    "郵便ボブ",
					Address: "bob.yuubin@example.com",
				},
				{
					Name:    "郵便キャロル",
					Address: "carol.yuubin@example.com",
				},
			},
			Cc: []*mail.Address{
				{
					Name:    "郵便ダン",
					Address: "dan.yuubin@example.com",
				},
				{
					Name:    "郵便イーブ",
					Address: "eve.yubin@example.com",
				},
			},
			Bcc: []*mail.Address{
				{
					Name:    "郵便フランク",
					Address: "frank.yuubin@example.com",
				},
				{
					Name:    "郵便グレイス",
					Address: "grace.yubin@example.com",
				},
			},
			MessageID:  "Message-Id-1@example.com",
			InReplyTo:  []MessageId{"Message-Id-0@example.com"},
			References: []MessageId{"Message-Id-0@example.com"},
			Comments:   "Message Header Comment",
			Keywords:   []string{"Keyword 1", "Keyword 2"},
			ResentDate: expectedDate,
			ResentFrom: []*mail.Address{
				{
					Name:    "郵便アリス",
					Address: "alice.yuubin@example.net",
				},
				{
					Name:    "郵便アリス",
					Address: "alice.yuubin@example.com",
				},
			},
			ResentSender: &mail.Address{
				Name:    "郵便アリス",
				Address: "alice.yuubin@example.net",
			},
			ResentTo: []*mail.Address{
				{
					Name:    "郵便ボブ",
					Address: "bob.yuubin@example.net",
				},
				{
					Name:    "郵便キャロル",
					Address: "carol.yuubin@example.net",
				},
			},
			ResentCc: []*mail.Address{
				{
					Name:    "郵便ダン",
					Address: "dan.yuubin@example.net",
				},
				{
					Name:    "郵便イーブ",
					Address: "eve.yubin@example.net",
				},
			},
			ResentBcc: []*mail.Address{
				{
					Name:    "郵便フランク",
					Address: "frank.yuubin@example.net",
				},
				{
					Name:    "郵便グレイス",
					Address: "grace.yubin@example.net",
				},
			},
			ResentMessageID: "Message-Id-1@example.net",
			ContentType: ContentTypeHeader{
				ContentType: "multipart/related",
				Params: map[string]string{
					"boundary": "RelatedBoundaryString",
					"charset":  "utf-8",
				},
			},
			ExtraHeaders: map[string][]string{
				"X-Clacks-Overhead": {"GNU Terry Pratchett"},
				"X-Script/function/\t !\"#$%&'()*+,-./;<=>?@[\\]^_`{|}~": {
					"TEST VALUE 1\t !\"#$%&'()*+,-./0123456789:;<=>?@ABCDEFGHIJKLMNOPQRSTUVWXYZ[\\]^_` abcdefghijklmnopqrstuvwxyz{|}~",
					"TEST VALUE 2\t !\"#$%&'()*+,-./0123456789:;<=>?@ABCDEFGHIJKLMNOPQRSTUVWXYZ[\\]^_` abcdefghijklmnopqrstuvwxyz{|}~",
				},
			},
		},
		Text: `色は匂えど
散りぬるを
我が世誰ぞ
常ならん
有為の奥山
今日越えて
浅き夢見じ
酔いもせず。

Iro wa nioedo / Chirinuru o / Wa ga yo tare zo / Tsune naran / Ui no okuyama / Kyo koete Asaki yume miji / Yoi mo sezu.

とりなくこゑすゆめさませみよあけわたるひんかしをそらいろはえておきつへにほふねむれゐぬもやのうち。

天地星空
山川峰谷
雲霧室苔
人犬上末
硫黄猿生ふ為よ
榎の枝を馴れ居て。

田居に出で菜摘むわれをぞ君召すと求食り追ひゆく山城の打酔へる子ら藻葉干せよえ舟繋けぬ。`,
		EnrichedText: `<bold>色は匂えど</bold>
<italic>散りぬるを</italic>
<fixed>我が世誰ぞ</fixed>
<underline>常ならん</underline>
有為の奥山
今日越えて
浅き夢見じ
酔いもせず。

Iro wa nioedo / Chirinuru o / Wa ga yo tare zo / Tsune naran / Ui no okuyama / Kyo koete Asaki yume miji / Yoi mo sezu.

とりなくこゑすゆめさませみよあけわたるひんかしをそらいろはえておきつへにほふねむれゐぬもやのうち。

天地星空
山川峰谷
雲霧室苔
人犬上末
硫黄猿生ふ為よ
榎の枝を馴れ居て。

田居に出で菜摘むわれをぞ君召すと求食り追ひゆく山城の打酔へる子ら藻葉干せよえ舟繋けぬ。`,
		HTML: `<html>
<div dir="ltr">
<p>色は匂えど<br />
散りぬるを<br />
我が世誰ぞ<br />
常ならん<br />
有為の奥山<br />
今日越えて<br />
浅き夢見じ<br />
酔いもせず。</p>

<p>Iro wa nioedo / Chirinuru o / Wa ga yo tare zo / Tsune naran / Ui no okuyama / Kyo koete Asaki yume miji / Yoi mo sezu.</p>

<p>とりなくこゑすゆめさませみよあけわたるひんかしをそらいろはえておきつへにほふねむれゐぬもやのうち。</p>

<p>天地星空<br />
山川峰谷<br />
雲霧室苔<br />
人犬上末<br />
硫黄猿生ふ為よ<br />
榎の枝を馴れ居て。</p>

<p>田居に出で菜摘むわれをぞ君召すと求食り追ひゆく山城の打酔へる子ら藻葉干せよえ舟繋けぬ。</p>
</div>
</html>`,
		InlineFiles: []InlineFile{
			{
				ContentID: "inline-jpg-image.jpg@example.com",
				ContentType: ContentTypeHeader{
					ContentType: "image/jpeg",
					Params: map[string]string{
						"name": "inline-jpg-image-name.jpg",
					},
				},
				ContentDisposition: ContentDispositionHeader{
					ContentDisposition: inline,
					Params: map[string]string{
						"filename": "inline-jpg-image-filename.jpg",
					},
				},
				Data: []byte{255, 216, 255, 219, 0, 67, 0, 3, 2, 2, 2, 2, 2, 3, 2, 2, 2, 3, 3, 3, 3, 4,
					6, 4, 4, 4, 4, 4, 8, 6, 6, 5, 6, 9, 8, 10, 10, 9, 8, 9, 9, 10, 12, 15, 12, 10, 11, 14, 11, 9, 9,
					13, 17, 13, 14, 15, 16, 16, 17, 16, 10, 12, 18, 19, 18, 16, 19, 15, 16, 16, 16, 255, 201, 0, 11, 8,
					0, 1, 0, 1, 1, 1, 17, 0, 255, 204, 0, 6, 0, 16, 16, 5, 255, 218, 0, 8, 1, 1, 0, 0, 63, 0, 210, 207,
					32, 255, 217},
			},
		},
	}

	testEmailFromFile(t, fp, expectedEmail)
}
func TestParseEmailJapaneseMultipartRelatedUtf8OverQuotedprintable(t *testing.T) {
	fp := "tests/test_japanese_multipart_related_utf-8_over_quoted-printable.txt"
	tz, _ := time.LoadLocation("Asia/Tokyo")
	expectedDate, _ := time.Parse(
		time.RFC1123Z+" (MST)",
		time.Date(2019, time.April, 1, 7, 55, 0, 0, tz).Format(time.RFC1123Z+" (MST)"))
	expectedEmail := Email{
		Headers: Headers{
			Date:    expectedDate,
			Subject: "📧 Test いろは歌",
			ReplyTo: []*mail.Address{
				{
					Name:    "郵便アリス",
					Address: "alice.yuubin@example.net",
				},
			},
			Sender: &mail.Address{
				Name:    "郵便アリス",
				Address: "alice.yuubin@example.com",
			},
			From: []*mail.Address{
				{
					Name:    "郵便アリス",
					Address: "alice.yuubin@example.com",
				},
				{
					Name:    "郵便アリス",
					Address: "alice.yuubin@example.net",
				},
			},
			To: []*mail.Address{
				{
					Name:    "郵便ボブ",
					Address: "bob.yuubin@example.com",
				},
				{
					Name:    "郵便キャロル",
					Address: "carol.yuubin@example.com",
				},
			},
			Cc: []*mail.Address{
				{
					Name:    "郵便ダン",
					Address: "dan.yuubin@example.com",
				},
				{
					Name:    "郵便イーブ",
					Address: "eve.yubin@example.com",
				},
			},
			Bcc: []*mail.Address{
				{
					Name:    "郵便フランク",
					Address: "frank.yuubin@example.com",
				},
				{
					Name:    "郵便グレイス",
					Address: "grace.yubin@example.com",
				},
			},
			MessageID:  "Message-Id-1@example.com",
			InReplyTo:  []MessageId{"Message-Id-0@example.com"},
			References: []MessageId{"Message-Id-0@example.com"},
			Comments:   "Message Header Comment",
			Keywords:   []string{"Keyword 1", "Keyword 2"},
			ResentDate: expectedDate,
			ResentFrom: []*mail.Address{
				{
					Name:    "郵便アリス",
					Address: "alice.yuubin@example.net",
				},
				{
					Name:    "郵便アリス",
					Address: "alice.yuubin@example.com",
				},
			},
			ResentSender: &mail.Address{
				Name:    "郵便アリス",
				Address: "alice.yuubin@example.net",
			},
			ResentTo: []*mail.Address{
				{
					Name:    "郵便ボブ",
					Address: "bob.yuubin@example.net",
				},
				{
					Name:    "郵便キャロル",
					Address: "carol.yuubin@example.net",
				},
			},
			ResentCc: []*mail.Address{
				{
					Name:    "郵便ダン",
					Address: "dan.yuubin@example.net",
				},
				{
					Name:    "郵便イーブ",
					Address: "eve.yubin@example.net",
				},
			},
			ResentBcc: []*mail.Address{
				{
					Name:    "郵便フランク",
					Address: "frank.yuubin@example.net",
				},
				{
					Name:    "郵便グレイス",
					Address: "grace.yubin@example.net",
				},
			},
			ResentMessageID: "Message-Id-1@example.net",
			ContentType: ContentTypeHeader{
				ContentType: "multipart/related",
				Params: map[string]string{
					"boundary": "RelatedBoundaryString",
					"charset":  "utf-8",
				},
			},
			ExtraHeaders: map[string][]string{
				"X-Clacks-Overhead": {"GNU Terry Pratchett"},
				"X-Script/function/\t !\"#$%&'()*+,-./;<=>?@[\\]^_`{|}~": {
					"TEST VALUE 1\t !\"#$%&'()*+,-./0123456789:;<=>?@ABCDEFGHIJKLMNOPQRSTUVWXYZ[\\]^_` abcdefghijklmnopqrstuvwxyz{|}~",
					"TEST VALUE 2\t !\"#$%&'()*+,-./0123456789:;<=>?@ABCDEFGHIJKLMNOPQRSTUVWXYZ[\\]^_` abcdefghijklmnopqrstuvwxyz{|}~",
				},
			},
		},
		Text: `色は匂えど
散りぬるを
我が世誰ぞ
常ならん
有為の奥山
今日越えて
浅き夢見じ
酔いもせず。

Iro wa nioedo / Chirinuru o / Wa ga yo tare zo / Tsune naran / Ui no okuyama / Kyo koete Asaki yume miji / Yoi mo sezu.

とりなくこゑすゆめさませみよあけわたるひんかしをそらいろはえておきつへにほふねむれゐぬもやのうち。

天地星空
山川峰谷
雲霧室苔
人犬上末
硫黄猿生ふ為よ
榎の枝を馴れ居て。

田居に出で菜摘むわれをぞ君召すと求食り追ひゆく山城の打酔へる子ら藻葉干せよえ舟繋けぬ。`,
		EnrichedText: `<bold>色は匂えど</bold>
<italic>散りぬるを</italic>
<fixed>我が世誰ぞ</fixed>
<underline>常ならん</underline>
有為の奥山
今日越えて
浅き夢見じ
酔いもせず。

Iro wa nioedo / Chirinuru o / Wa ga yo tare zo / Tsune naran / Ui no okuyama / Kyo koete Asaki yume miji / Yoi mo sezu.

とりなくこゑすゆめさませみよあけわたるひんかしをそらいろはえておきつへにほふねむれゐぬもやのうち。

天地星空
山川峰谷
雲霧室苔
人犬上末
硫黄猿生ふ為よ
榎の枝を馴れ居て。

田居に出で菜摘むわれをぞ君召すと求食り追ひゆく山城の打酔へる子ら藻葉干せよえ舟繋けぬ。`,
		HTML: `<html>
<div dir="ltr">
<p>色は匂えど<br />
散りぬるを<br />
我が世誰ぞ<br />
常ならん<br />
有為の奥山<br />
今日越えて<br />
浅き夢見じ<br />
酔いもせず。</p>

<p>Iro wa nioedo / Chirinuru o / Wa ga yo tare zo / Tsune naran / Ui no okuyama / Kyo koete Asaki yume miji / Yoi mo sezu.</p>

<p>とりなくこゑすゆめさませみよあけわたるひんかしをそらいろはえておきつへにほふねむれゐぬもやのうち。</p>

<p>天地星空<br />
山川峰谷<br />
雲霧室苔<br />
人犬上末<br />
硫黄猿生ふ為よ<br />
榎の枝を馴れ居て。</p>

<p>田居に出で菜摘むわれをぞ君召すと求食り追ひゆく山城の打酔へる子ら藻葉干せよえ舟繋けぬ。</p>
</div>
</html>`,
		InlineFiles: []InlineFile{
			{
				ContentID: "inline-jpg-image.jpg@example.com",
				ContentType: ContentTypeHeader{
					ContentType: "image/jpeg",
					Params: map[string]string{
						"name": "inline-jpg-image-name.jpg",
					},
				},
				ContentDisposition: ContentDispositionHeader{
					ContentDisposition: inline,
					Params: map[string]string{
						"filename": "inline-jpg-image-filename.jpg",
					},
				},
				Data: []byte{255, 216, 255, 219, 0, 67, 0, 3, 2, 2, 2, 2, 2, 3, 2, 2, 2, 3, 3, 3, 3, 4,
					6, 4, 4, 4, 4, 4, 8, 6, 6, 5, 6, 9, 8, 10, 10, 9, 8, 9, 9, 10, 12, 15, 12, 10, 11, 14, 11, 9, 9,
					13, 17, 13, 14, 15, 16, 16, 17, 16, 10, 12, 18, 19, 18, 16, 19, 15, 16, 16, 16, 255, 201, 0, 11, 8,
					0, 1, 0, 1, 1, 1, 17, 0, 255, 204, 0, 6, 0, 16, 16, 5, 255, 218, 0, 8, 1, 1, 0, 0, 63, 0, 210, 207,
					32, 255, 217},
			},
		},
	}

	testEmailFromFile(t, fp, expectedEmail)
}
func TestParseEmailJapaneseMultipartRelatedIso2022jpOver7bit(t *testing.T) {
	fp := "tests/test_japanese_multipart_related_iso-2022-jp_over_7bit.txt"
	tz, _ := time.LoadLocation("Asia/Tokyo")
	expectedDate, _ := time.Parse(
		time.RFC1123Z+" (MST)",
		time.Date(2019, time.April, 1, 7, 55, 0, 0, tz).Format(time.RFC1123Z+" (MST)"))
	expectedEmail := Email{
		Headers: Headers{
			Date:    expectedDate,
			Subject: "📧 Test いろは歌",
			ReplyTo: []*mail.Address{
				{
					Name:    "郵便アリス",
					Address: "alice.yuubin@example.net",
				},
			},
			Sender: &mail.Address{
				Name:    "郵便アリス",
				Address: "alice.yuubin@example.com",
			},
			From: []*mail.Address{
				{
					Name:    "郵便アリス",
					Address: "alice.yuubin@example.com",
				},
				{
					Name:    "郵便アリス",
					Address: "alice.yuubin@example.net",
				},
			},
			To: []*mail.Address{
				{
					Name:    "郵便ボブ",
					Address: "bob.yuubin@example.com",
				},
				{
					Name:    "郵便キャロル",
					Address: "carol.yuubin@example.com",
				},
			},
			Cc: []*mail.Address{
				{
					Name:    "郵便ダン",
					Address: "dan.yuubin@example.com",
				},
				{
					Name:    "郵便イーブ",
					Address: "eve.yubin@example.com",
				},
			},
			Bcc: []*mail.Address{
				{
					Name:    "郵便フランク",
					Address: "frank.yuubin@example.com",
				},
				{
					Name:    "郵便グレイス",
					Address: "grace.yubin@example.com",
				},
			},
			MessageID:  "Message-Id-1@example.com",
			InReplyTo:  []MessageId{"Message-Id-0@example.com"},
			References: []MessageId{"Message-Id-0@example.com"},
			Comments:   "Message Header Comment",
			Keywords:   []string{"Keyword 1", "Keyword 2"},
			ResentDate: expectedDate,
			ResentFrom: []*mail.Address{
				{
					Name:    "郵便アリス",
					Address: "alice.yuubin@example.net",
				},
				{
					Name:    "郵便アリス",
					Address: "alice.yuubin@example.com",
				},
			},
			ResentSender: &mail.Address{
				Name:    "郵便アリス",
				Address: "alice.yuubin@example.net",
			},
			ResentTo: []*mail.Address{
				{
					Name:    "郵便ボブ",
					Address: "bob.yuubin@example.net",
				},
				{
					Name:    "郵便キャロル",
					Address: "carol.yuubin@example.net",
				},
			},
			ResentCc: []*mail.Address{
				{
					Name:    "郵便ダン",
					Address: "dan.yuubin@example.net",
				},
				{
					Name:    "郵便イーブ",
					Address: "eve.yubin@example.net",
				},
			},
			ResentBcc: []*mail.Address{
				{
					Name:    "郵便フランク",
					Address: "frank.yuubin@example.net",
				},
				{
					Name:    "郵便グレイス",
					Address: "grace.yubin@example.net",
				},
			},
			ResentMessageID: "Message-Id-1@example.net",
			ContentType: ContentTypeHeader{
				ContentType: "multipart/related",
				Params: map[string]string{
					"boundary": "RelatedBoundaryString",
					"charset":  "iso-2022-jp",
				},
			},
			ExtraHeaders: map[string][]string{
				"X-Clacks-Overhead": {"GNU Terry Pratchett"},
				"X-Script/function/\t !\"#$%&'()*+,-./;<=>?@[\\]^_`{|}~": {
					"TEST VALUE 1\t !\"#$%&'()*+,-./0123456789:;<=>?@ABCDEFGHIJKLMNOPQRSTUVWXYZ[\\]^_` abcdefghijklmnopqrstuvwxyz{|}~",
					"TEST VALUE 2\t !\"#$%&'()*+,-./0123456789:;<=>?@ABCDEFGHIJKLMNOPQRSTUVWXYZ[\\]^_` abcdefghijklmnopqrstuvwxyz{|}~",
				},
			},
		},
		Text: `色は匂えど
散りぬるを
我が世誰ぞ
常ならん
有為の奥山
今日越えて
浅き夢見じ
酔いもせず。

Iro wa nioedo / Chirinuru o / Wa ga yo tare zo / Tsune naran / Ui no okuyama / Kyo koete Asaki yume miji / Yoi mo sezu.

とりなくこゑすゆめさませみよあけわたるひんかしをそらいろはえておきつへにほふねむれゐぬもやのうち。

天地星空
山川峰谷
雲霧室苔
人犬上末
硫黄猿生ふ為よ
榎の枝を馴れ居て。

田居に出で菜摘むわれをぞ君召すと求食り追ひゆく山城の打酔へる子ら藻葉干せよえ舟繋けぬ。`,
		EnrichedText: `<bold>色は匂えど</bold>
<italic>散りぬるを</italic>
<fixed>我が世誰ぞ</fixed>
<underline>常ならん</underline>
有為の奥山
今日越えて
浅き夢見じ
酔いもせず。

Iro wa nioedo / Chirinuru o / Wa ga yo tare zo / Tsune naran / Ui no okuyama / Kyo koete Asaki yume miji / Yoi mo sezu.

とりなくこゑすゆめさませみよあけわたるひんかしをそらいろはえておきつへにほふねむれゐぬもやのうち。

天地星空
山川峰谷
雲霧室苔
人犬上末
硫黄猿生ふ為よ
榎の枝を馴れ居て。

田居に出で菜摘むわれをぞ君召すと求食り追ひゆく山城の打酔へる子ら藻葉干せよえ舟繋けぬ。`,
		HTML: `<html>
<div dir="ltr">
<p>色は匂えど<br />
散りぬるを<br />
我が世誰ぞ<br />
常ならん<br />
有為の奥山<br />
今日越えて<br />
浅き夢見じ<br />
酔いもせず。</p>

<p>Iro wa nioedo / Chirinuru o / Wa ga yo tare zo / Tsune naran / Ui no okuyama / Kyo koete Asaki yume miji / Yoi mo sezu.</p>

<p>とりなくこゑすゆめさませみよあけわたるひんかしをそらいろはえておきつへにほふねむれゐぬもやのうち。</p>

<p>天地星空<br />
山川峰谷<br />
雲霧室苔<br />
人犬上末<br />
硫黄猿生ふ為よ<br />
榎の枝を馴れ居て。</p>

<p>田居に出で菜摘むわれをぞ君召すと求食り追ひゆく山城の打酔へる子ら藻葉干せよえ舟繋けぬ。</p>
</div>
</html>`,
		InlineFiles: []InlineFile{
			{
				ContentID: "inline-jpg-image.jpg@example.com",
				ContentType: ContentTypeHeader{
					ContentType: "image/jpeg",
					Params: map[string]string{
						"name": "inline-jpg-image-name.jpg",
					},
				},
				ContentDisposition: ContentDispositionHeader{
					ContentDisposition: inline,
					Params: map[string]string{
						"filename": "inline-jpg-image-filename.jpg",
					},
				},
				Data: []byte{255, 216, 255, 219, 0, 67, 0, 3, 2, 2, 2, 2, 2, 3, 2, 2, 2, 3, 3, 3, 3, 4,
					6, 4, 4, 4, 4, 4, 8, 6, 6, 5, 6, 9, 8, 10, 10, 9, 8, 9, 9, 10, 12, 15, 12, 10, 11, 14, 11, 9, 9,
					13, 17, 13, 14, 15, 16, 16, 17, 16, 10, 12, 18, 19, 18, 16, 19, 15, 16, 16, 16, 255, 201, 0, 11, 8,
					0, 1, 0, 1, 1, 1, 17, 0, 255, 204, 0, 6, 0, 16, 16, 5, 255, 218, 0, 8, 1, 1, 0, 0, 63, 0, 210, 207,
					32, 255, 217},
			},
		},
	}

	testEmailFromFile(t, fp, expectedEmail)
}
func TestParseEmailJapaneseMultipartRelatedIso2022jpOverBase64(t *testing.T) {
	fp := "tests/test_japanese_multipart_related_iso-2022-jp_over_base64.txt"
	tz, _ := time.LoadLocation("Asia/Tokyo")
	expectedDate, _ := time.Parse(
		time.RFC1123Z+" (MST)",
		time.Date(2019, time.April, 1, 7, 55, 0, 0, tz).Format(time.RFC1123Z+" (MST)"))
	expectedEmail := Email{
		Headers: Headers{
			Date:    expectedDate,
			Subject: "📧 Test いろは歌",
			ReplyTo: []*mail.Address{
				{
					Name:    "郵便アリス",
					Address: "alice.yuubin@example.net",
				},
			},
			Sender: &mail.Address{
				Name:    "郵便アリス",
				Address: "alice.yuubin@example.com",
			},
			From: []*mail.Address{
				{
					Name:    "郵便アリス",
					Address: "alice.yuubin@example.com",
				},
				{
					Name:    "郵便アリス",
					Address: "alice.yuubin@example.net",
				},
			},
			To: []*mail.Address{
				{
					Name:    "郵便ボブ",
					Address: "bob.yuubin@example.com",
				},
				{
					Name:    "郵便キャロル",
					Address: "carol.yuubin@example.com",
				},
			},
			Cc: []*mail.Address{
				{
					Name:    "郵便ダン",
					Address: "dan.yuubin@example.com",
				},
				{
					Name:    "郵便イーブ",
					Address: "eve.yubin@example.com",
				},
			},
			Bcc: []*mail.Address{
				{
					Name:    "郵便フランク",
					Address: "frank.yuubin@example.com",
				},
				{
					Name:    "郵便グレイス",
					Address: "grace.yubin@example.com",
				},
			},
			MessageID:  "Message-Id-1@example.com",
			InReplyTo:  []MessageId{"Message-Id-0@example.com"},
			References: []MessageId{"Message-Id-0@example.com"},
			Comments:   "Message Header Comment",
			Keywords:   []string{"Keyword 1", "Keyword 2"},
			ResentDate: expectedDate,
			ResentFrom: []*mail.Address{
				{
					Name:    "郵便アリス",
					Address: "alice.yuubin@example.net",
				},
				{
					Name:    "郵便アリス",
					Address: "alice.yuubin@example.com",
				},
			},
			ResentSender: &mail.Address{
				Name:    "郵便アリス",
				Address: "alice.yuubin@example.net",
			},
			ResentTo: []*mail.Address{
				{
					Name:    "郵便ボブ",
					Address: "bob.yuubin@example.net",
				},
				{
					Name:    "郵便キャロル",
					Address: "carol.yuubin@example.net",
				},
			},
			ResentCc: []*mail.Address{
				{
					Name:    "郵便ダン",
					Address: "dan.yuubin@example.net",
				},
				{
					Name:    "郵便イーブ",
					Address: "eve.yubin@example.net",
				},
			},
			ResentBcc: []*mail.Address{
				{
					Name:    "郵便フランク",
					Address: "frank.yuubin@example.net",
				},
				{
					Name:    "郵便グレイス",
					Address: "grace.yubin@example.net",
				},
			},
			ResentMessageID: "Message-Id-1@example.net",
			ContentType: ContentTypeHeader{
				ContentType: "multipart/related",
				Params: map[string]string{
					"boundary": "RelatedBoundaryString",
					"charset":  "iso-2022-jp",
				},
			},
			ExtraHeaders: map[string][]string{
				"X-Clacks-Overhead": {"GNU Terry Pratchett"},
				"X-Script/function/\t !\"#$%&'()*+,-./;<=>?@[\\]^_`{|}~": {
					"TEST VALUE 1\t !\"#$%&'()*+,-./0123456789:;<=>?@ABCDEFGHIJKLMNOPQRSTUVWXYZ[\\]^_` abcdefghijklmnopqrstuvwxyz{|}~",
					"TEST VALUE 2\t !\"#$%&'()*+,-./0123456789:;<=>?@ABCDEFGHIJKLMNOPQRSTUVWXYZ[\\]^_` abcdefghijklmnopqrstuvwxyz{|}~",
				},
			},
		},
		Text: `色は匂えど
散りぬるを
我が世誰ぞ
常ならん
有為の奥山
今日越えて
浅き夢見じ
酔いもせず。

Iro wa nioedo / Chirinuru o / Wa ga yo tare zo / Tsune naran / Ui no okuyama / Kyo koete Asaki yume miji / Yoi mo sezu.

とりなくこゑすゆめさませみよあけわたるひんかしをそらいろはえておきつへにほふねむれゐぬもやのうち。

天地星空
山川峰谷
雲霧室苔
人犬上末
硫黄猿生ふ為よ
榎の枝を馴れ居て。

田居に出で菜摘むわれをぞ君召すと求食り追ひゆく山城の打酔へる子ら藻葉干せよえ舟繋けぬ。`,
		EnrichedText: `<bold>色は匂えど</bold>
<italic>散りぬるを</italic>
<fixed>我が世誰ぞ</fixed>
<underline>常ならん</underline>
有為の奥山
今日越えて
浅き夢見じ
酔いもせず。

Iro wa nioedo / Chirinuru o / Wa ga yo tare zo / Tsune naran / Ui no okuyama / Kyo koete Asaki yume miji / Yoi mo sezu.

とりなくこゑすゆめさませみよあけわたるひんかしをそらいろはえておきつへにほふねむれゐぬもやのうち。

天地星空
山川峰谷
雲霧室苔
人犬上末
硫黄猿生ふ為よ
榎の枝を馴れ居て。

田居に出で菜摘むわれをぞ君召すと求食り追ひゆく山城の打酔へる子ら藻葉干せよえ舟繋けぬ。`,
		HTML: `<html>
<div dir="ltr">
<p>色は匂えど<br />
散りぬるを<br />
我が世誰ぞ<br />
常ならん<br />
有為の奥山<br />
今日越えて<br />
浅き夢見じ<br />
酔いもせず。</p>

<p>Iro wa nioedo / Chirinuru o / Wa ga yo tare zo / Tsune naran / Ui no okuyama / Kyo koete Asaki yume miji / Yoi mo sezu.</p>

<p>とりなくこゑすゆめさませみよあけわたるひんかしをそらいろはえておきつへにほふねむれゐぬもやのうち。</p>

<p>天地星空<br />
山川峰谷<br />
雲霧室苔<br />
人犬上末<br />
硫黄猿生ふ為よ<br />
榎の枝を馴れ居て。</p>

<p>田居に出で菜摘むわれをぞ君召すと求食り追ひゆく山城の打酔へる子ら藻葉干せよえ舟繋けぬ。</p>
</div>
</html>`,
		InlineFiles: []InlineFile{
			{
				ContentID: "inline-jpg-image.jpg@example.com",
				ContentType: ContentTypeHeader{
					ContentType: "image/jpeg",
					Params: map[string]string{
						"name": "inline-jpg-image-name.jpg",
					},
				},
				ContentDisposition: ContentDispositionHeader{
					ContentDisposition: inline,
					Params: map[string]string{
						"filename": "inline-jpg-image-filename.jpg",
					},
				},
				Data: []byte{255, 216, 255, 219, 0, 67, 0, 3, 2, 2, 2, 2, 2, 3, 2, 2, 2, 3, 3, 3, 3, 4,
					6, 4, 4, 4, 4, 4, 8, 6, 6, 5, 6, 9, 8, 10, 10, 9, 8, 9, 9, 10, 12, 15, 12, 10, 11, 14, 11, 9, 9,
					13, 17, 13, 14, 15, 16, 16, 17, 16, 10, 12, 18, 19, 18, 16, 19, 15, 16, 16, 16, 255, 201, 0, 11, 8,
					0, 1, 0, 1, 1, 1, 17, 0, 255, 204, 0, 6, 0, 16, 16, 5, 255, 218, 0, 8, 1, 1, 0, 0, 63, 0, 210, 207,
					32, 255, 217},
			},
		},
	}

	testEmailFromFile(t, fp, expectedEmail)
}
func TestParseEmailJapaneseMultipartRelatedIso2022jpOverQuotedprintable(t *testing.T) {
	fp := "tests/test_japanese_multipart_related_iso-2022-jp_over_quoted-printable.txt"
	tz, _ := time.LoadLocation("Asia/Tokyo")
	expectedDate, _ := time.Parse(
		time.RFC1123Z+" (MST)",
		time.Date(2019, time.April, 1, 7, 55, 0, 0, tz).Format(time.RFC1123Z+" (MST)"))
	expectedEmail := Email{
		Headers: Headers{
			Date:    expectedDate,
			Subject: "📧 Test いろは歌",
			ReplyTo: []*mail.Address{
				{
					Name:    "郵便アリス",
					Address: "alice.yuubin@example.net",
				},
			},
			Sender: &mail.Address{
				Name:    "郵便アリス",
				Address: "alice.yuubin@example.com",
			},
			From: []*mail.Address{
				{
					Name:    "郵便アリス",
					Address: "alice.yuubin@example.com",
				},
				{
					Name:    "郵便アリス",
					Address: "alice.yuubin@example.net",
				},
			},
			To: []*mail.Address{
				{
					Name:    "郵便ボブ",
					Address: "bob.yuubin@example.com",
				},
				{
					Name:    "郵便キャロル",
					Address: "carol.yuubin@example.com",
				},
			},
			Cc: []*mail.Address{
				{
					Name:    "郵便ダン",
					Address: "dan.yuubin@example.com",
				},
				{
					Name:    "郵便イーブ",
					Address: "eve.yubin@example.com",
				},
			},
			Bcc: []*mail.Address{
				{
					Name:    "郵便フランク",
					Address: "frank.yuubin@example.com",
				},
				{
					Name:    "郵便グレイス",
					Address: "grace.yubin@example.com",
				},
			},
			MessageID:  "Message-Id-1@example.com",
			InReplyTo:  []MessageId{"Message-Id-0@example.com"},
			References: []MessageId{"Message-Id-0@example.com"},
			Comments:   "Message Header Comment",
			Keywords:   []string{"Keyword 1", "Keyword 2"},
			ResentDate: expectedDate,
			ResentFrom: []*mail.Address{
				{
					Name:    "郵便アリス",
					Address: "alice.yuubin@example.net",
				},
				{
					Name:    "郵便アリス",
					Address: "alice.yuubin@example.com",
				},
			},
			ResentSender: &mail.Address{
				Name:    "郵便アリス",
				Address: "alice.yuubin@example.net",
			},
			ResentTo: []*mail.Address{
				{
					Name:    "郵便ボブ",
					Address: "bob.yuubin@example.net",
				},
				{
					Name:    "郵便キャロル",
					Address: "carol.yuubin@example.net",
				},
			},
			ResentCc: []*mail.Address{
				{
					Name:    "郵便ダン",
					Address: "dan.yuubin@example.net",
				},
				{
					Name:    "郵便イーブ",
					Address: "eve.yubin@example.net",
				},
			},
			ResentBcc: []*mail.Address{
				{
					Name:    "郵便フランク",
					Address: "frank.yuubin@example.net",
				},
				{
					Name:    "郵便グレイス",
					Address: "grace.yubin@example.net",
				},
			},
			ResentMessageID: "Message-Id-1@example.net",
			ContentType: ContentTypeHeader{
				ContentType: "multipart/related",
				Params: map[string]string{
					"boundary": "RelatedBoundaryString",
					"charset":  "iso-2022-jp",
				},
			},
			ExtraHeaders: map[string][]string{
				"X-Clacks-Overhead": {"GNU Terry Pratchett"},
				"X-Script/function/\t !\"#$%&'()*+,-./;<=>?@[\\]^_`{|}~": {
					"TEST VALUE 1\t !\"#$%&'()*+,-./0123456789:;<=>?@ABCDEFGHIJKLMNOPQRSTUVWXYZ[\\]^_` abcdefghijklmnopqrstuvwxyz{|}~",
					"TEST VALUE 2\t !\"#$%&'()*+,-./0123456789:;<=>?@ABCDEFGHIJKLMNOPQRSTUVWXYZ[\\]^_` abcdefghijklmnopqrstuvwxyz{|}~",
				},
			},
		},
		Text: `色は匂えど
散りぬるを
我が世誰ぞ
常ならん
有為の奥山
今日越えて
浅き夢見じ
酔いもせず。

Iro wa nioedo / Chirinuru o / Wa ga yo tare zo / Tsune naran / Ui no okuyama / Kyo koete Asaki yume miji / Yoi mo sezu.

とりなくこゑすゆめさませみよあけわたるひんかしをそらいろはえておきつへにほふねむれゐぬもやのうち。

天地星空
山川峰谷
雲霧室苔
人犬上末
硫黄猿生ふ為よ
榎の枝を馴れ居て。

田居に出で菜摘むわれをぞ君召すと求食り追ひゆく山城の打酔へる子ら藻葉干せよえ舟繋けぬ。`,
		EnrichedText: `<bold>色は匂えど</bold>
<italic>散りぬるを</italic>
<fixed>我が世誰ぞ</fixed>
<underline>常ならん</underline>
有為の奥山
今日越えて
浅き夢見じ
酔いもせず。

Iro wa nioedo / Chirinuru o / Wa ga yo tare zo / Tsune naran / Ui no okuyama / Kyo koete Asaki yume miji / Yoi mo sezu.

とりなくこゑすゆめさませみよあけわたるひんかしをそらいろはえておきつへにほふねむれゐぬもやのうち。

天地星空
山川峰谷
雲霧室苔
人犬上末
硫黄猿生ふ為よ
榎の枝を馴れ居て。

田居に出で菜摘むわれをぞ君召すと求食り追ひゆく山城の打酔へる子ら藻葉干せよえ舟繋けぬ。`,
		HTML: `<html>
<div dir="ltr">
<p>色は匂えど<br />
散りぬるを<br />
我が世誰ぞ<br />
常ならん<br />
有為の奥山<br />
今日越えて<br />
浅き夢見じ<br />
酔いもせず。</p>

<p>Iro wa nioedo / Chirinuru o / Wa ga yo tare zo / Tsune naran / Ui no okuyama / Kyo koete Asaki yume miji / Yoi mo sezu.</p>

<p>とりなくこゑすゆめさませみよあけわたるひんかしをそらいろはえておきつへにほふねむれゐぬもやのうち。</p>

<p>天地星空<br />
山川峰谷<br />
雲霧室苔<br />
人犬上末<br />
硫黄猿生ふ為よ<br />
榎の枝を馴れ居て。</p>

<p>田居に出で菜摘むわれをぞ君召すと求食り追ひゆく山城の打酔へる子ら藻葉干せよえ舟繋けぬ。</p>
</div>
</html>`,
		InlineFiles: []InlineFile{
			{
				ContentID: "inline-jpg-image.jpg@example.com",
				ContentType: ContentTypeHeader{
					ContentType: "image/jpeg",
					Params: map[string]string{
						"name": "inline-jpg-image-name.jpg",
					},
				},
				ContentDisposition: ContentDispositionHeader{
					ContentDisposition: inline,
					Params: map[string]string{
						"filename": "inline-jpg-image-filename.jpg",
					},
				},
				Data: []byte{255, 216, 255, 219, 0, 67, 0, 3, 2, 2, 2, 2, 2, 3, 2, 2, 2, 3, 3, 3, 3, 4,
					6, 4, 4, 4, 4, 4, 8, 6, 6, 5, 6, 9, 8, 10, 10, 9, 8, 9, 9, 10, 12, 15, 12, 10, 11, 14, 11, 9, 9,
					13, 17, 13, 14, 15, 16, 16, 17, 16, 10, 12, 18, 19, 18, 16, 19, 15, 16, 16, 16, 255, 201, 0, 11, 8,
					0, 1, 0, 1, 1, 1, 17, 0, 255, 204, 0, 6, 0, 16, 16, 5, 255, 218, 0, 8, 1, 1, 0, 0, 63, 0, 210, 207,
					32, 255, 217},
			},
		},
	}

	testEmailFromFile(t, fp, expectedEmail)
}
func TestParseEmailJapaneseMultipartRelatedEucjpOverBase64(t *testing.T) {
	fp := "tests/test_japanese_multipart_related_euc-jp_over_base64.txt"
	tz, _ := time.LoadLocation("Asia/Tokyo")
	expectedDate, _ := time.Parse(
		time.RFC1123Z+" (MST)",
		time.Date(2019, time.April, 1, 7, 55, 0, 0, tz).Format(time.RFC1123Z+" (MST)"))
	expectedEmail := Email{
		Headers: Headers{
			Date:    expectedDate,
			Subject: "📧 Test いろは歌",
			ReplyTo: []*mail.Address{
				{
					Name:    "郵便アリス",
					Address: "alice.yuubin@example.net",
				},
			},
			Sender: &mail.Address{
				Name:    "郵便アリス",
				Address: "alice.yuubin@example.com",
			},
			From: []*mail.Address{
				{
					Name:    "郵便アリス",
					Address: "alice.yuubin@example.com",
				},
				{
					Name:    "郵便アリス",
					Address: "alice.yuubin@example.net",
				},
			},
			To: []*mail.Address{
				{
					Name:    "郵便ボブ",
					Address: "bob.yuubin@example.com",
				},
				{
					Name:    "郵便キャロル",
					Address: "carol.yuubin@example.com",
				},
			},
			Cc: []*mail.Address{
				{
					Name:    "郵便ダン",
					Address: "dan.yuubin@example.com",
				},
				{
					Name:    "郵便イーブ",
					Address: "eve.yubin@example.com",
				},
			},
			Bcc: []*mail.Address{
				{
					Name:    "郵便フランク",
					Address: "frank.yuubin@example.com",
				},
				{
					Name:    "郵便グレイス",
					Address: "grace.yubin@example.com",
				},
			},
			MessageID:  "Message-Id-1@example.com",
			InReplyTo:  []MessageId{"Message-Id-0@example.com"},
			References: []MessageId{"Message-Id-0@example.com"},
			Comments:   "Message Header Comment",
			Keywords:   []string{"Keyword 1", "Keyword 2"},
			ResentDate: expectedDate,
			ResentFrom: []*mail.Address{
				{
					Name:    "郵便アリス",
					Address: "alice.yuubin@example.net",
				},
				{
					Name:    "郵便アリス",
					Address: "alice.yuubin@example.com",
				},
			},
			ResentSender: &mail.Address{
				Name:    "郵便アリス",
				Address: "alice.yuubin@example.net",
			},
			ResentTo: []*mail.Address{
				{
					Name:    "郵便ボブ",
					Address: "bob.yuubin@example.net",
				},
				{
					Name:    "郵便キャロル",
					Address: "carol.yuubin@example.net",
				},
			},
			ResentCc: []*mail.Address{
				{
					Name:    "郵便ダン",
					Address: "dan.yuubin@example.net",
				},
				{
					Name:    "郵便イーブ",
					Address: "eve.yubin@example.net",
				},
			},
			ResentBcc: []*mail.Address{
				{
					Name:    "郵便フランク",
					Address: "frank.yuubin@example.net",
				},
				{
					Name:    "郵便グレイス",
					Address: "grace.yubin@example.net",
				},
			},
			ResentMessageID: "Message-Id-1@example.net",
			ContentType: ContentTypeHeader{
				ContentType: "multipart/related",
				Params: map[string]string{
					"boundary": "RelatedBoundaryString",
					"charset":  "euc-jp",
				},
			},
			ExtraHeaders: map[string][]string{
				"X-Clacks-Overhead": {"GNU Terry Pratchett"},
				"X-Script/function/\t !\"#$%&'()*+,-./;<=>?@[\\]^_`{|}~": {
					"TEST VALUE 1\t !\"#$%&'()*+,-./0123456789:;<=>?@ABCDEFGHIJKLMNOPQRSTUVWXYZ[\\]^_` abcdefghijklmnopqrstuvwxyz{|}~",
					"TEST VALUE 2\t !\"#$%&'()*+,-./0123456789:;<=>?@ABCDEFGHIJKLMNOPQRSTUVWXYZ[\\]^_` abcdefghijklmnopqrstuvwxyz{|}~",
				},
			},
		},
		Text: `色は匂えど
散りぬるを
我が世誰ぞ
常ならん
有為の奥山
今日越えて
浅き夢見じ
酔いもせず。

Iro wa nioedo / Chirinuru o / Wa ga yo tare zo / Tsune naran / Ui no okuyama / Kyo koete Asaki yume miji / Yoi mo sezu.

とりなくこゑすゆめさませみよあけわたるひんかしをそらいろはえておきつへにほふねむれゐぬもやのうち。

天地星空
山川峰谷
雲霧室苔
人犬上末
硫黄猿生ふ為よ
榎の枝を馴れ居て。

田居に出で菜摘むわれをぞ君召すと求食り追ひゆく山城の打酔へる子ら藻葉干せよえ舟繋けぬ。`,
		EnrichedText: `<bold>色は匂えど</bold>
<italic>散りぬるを</italic>
<fixed>我が世誰ぞ</fixed>
<underline>常ならん</underline>
有為の奥山
今日越えて
浅き夢見じ
酔いもせず。

Iro wa nioedo / Chirinuru o / Wa ga yo tare zo / Tsune naran / Ui no okuyama / Kyo koete Asaki yume miji / Yoi mo sezu.

とりなくこゑすゆめさませみよあけわたるひんかしをそらいろはえておきつへにほふねむれゐぬもやのうち。

天地星空
山川峰谷
雲霧室苔
人犬上末
硫黄猿生ふ為よ
榎の枝を馴れ居て。

田居に出で菜摘むわれをぞ君召すと求食り追ひゆく山城の打酔へる子ら藻葉干せよえ舟繋けぬ。`,
		HTML: `<html>
<div dir="ltr">
<p>色は匂えど<br />
散りぬるを<br />
我が世誰ぞ<br />
常ならん<br />
有為の奥山<br />
今日越えて<br />
浅き夢見じ<br />
酔いもせず。</p>

<p>Iro wa nioedo / Chirinuru o / Wa ga yo tare zo / Tsune naran / Ui no okuyama / Kyo koete Asaki yume miji / Yoi mo sezu.</p>

<p>とりなくこゑすゆめさませみよあけわたるひんかしをそらいろはえておきつへにほふねむれゐぬもやのうち。</p>

<p>天地星空<br />
山川峰谷<br />
雲霧室苔<br />
人犬上末<br />
硫黄猿生ふ為よ<br />
榎の枝を馴れ居て。</p>

<p>田居に出で菜摘むわれをぞ君召すと求食り追ひゆく山城の打酔へる子ら藻葉干せよえ舟繋けぬ。</p>
</div>
</html>`,
		InlineFiles: []InlineFile{
			{
				ContentID: "inline-jpg-image.jpg@example.com",
				ContentType: ContentTypeHeader{
					ContentType: "image/jpeg",
					Params: map[string]string{
						"name": "inline-jpg-image-name.jpg",
					},
				},
				ContentDisposition: ContentDispositionHeader{
					ContentDisposition: inline,
					Params: map[string]string{
						"filename": "inline-jpg-image-filename.jpg",
					},
				},
				Data: []byte{255, 216, 255, 219, 0, 67, 0, 3, 2, 2, 2, 2, 2, 3, 2, 2, 2, 3, 3, 3, 3, 4,
					6, 4, 4, 4, 4, 4, 8, 6, 6, 5, 6, 9, 8, 10, 10, 9, 8, 9, 9, 10, 12, 15, 12, 10, 11, 14, 11, 9, 9,
					13, 17, 13, 14, 15, 16, 16, 17, 16, 10, 12, 18, 19, 18, 16, 19, 15, 16, 16, 16, 255, 201, 0, 11, 8,
					0, 1, 0, 1, 1, 1, 17, 0, 255, 204, 0, 6, 0, 16, 16, 5, 255, 218, 0, 8, 1, 1, 0, 0, 63, 0, 210, 207,
					32, 255, 217},
			},
		},
	}

	testEmailFromFile(t, fp, expectedEmail)
}
func TestParseEmailJapaneseMultipartRelatedEucjpOverQuotedprintable(t *testing.T) {
	fp := "tests/test_japanese_multipart_related_euc-jp_over_quoted-printable.txt"
	tz, _ := time.LoadLocation("Asia/Tokyo")
	expectedDate, _ := time.Parse(
		time.RFC1123Z+" (MST)",
		time.Date(2019, time.April, 1, 7, 55, 0, 0, tz).Format(time.RFC1123Z+" (MST)"))
	expectedEmail := Email{
		Headers: Headers{
			Date:    expectedDate,
			Subject: "📧 Test いろは歌",
			ReplyTo: []*mail.Address{
				{
					Name:    "郵便アリス",
					Address: "alice.yuubin@example.net",
				},
			},
			Sender: &mail.Address{
				Name:    "郵便アリス",
				Address: "alice.yuubin@example.com",
			},
			From: []*mail.Address{
				{
					Name:    "郵便アリス",
					Address: "alice.yuubin@example.com",
				},
				{
					Name:    "郵便アリス",
					Address: "alice.yuubin@example.net",
				},
			},
			To: []*mail.Address{
				{
					Name:    "郵便ボブ",
					Address: "bob.yuubin@example.com",
				},
				{
					Name:    "郵便キャロル",
					Address: "carol.yuubin@example.com",
				},
			},
			Cc: []*mail.Address{
				{
					Name:    "郵便ダン",
					Address: "dan.yuubin@example.com",
				},
				{
					Name:    "郵便イーブ",
					Address: "eve.yubin@example.com",
				},
			},
			Bcc: []*mail.Address{
				{
					Name:    "郵便フランク",
					Address: "frank.yuubin@example.com",
				},
				{
					Name:    "郵便グレイス",
					Address: "grace.yubin@example.com",
				},
			},
			MessageID:  "Message-Id-1@example.com",
			InReplyTo:  []MessageId{"Message-Id-0@example.com"},
			References: []MessageId{"Message-Id-0@example.com"},
			Comments:   "Message Header Comment",
			Keywords:   []string{"Keyword 1", "Keyword 2"},
			ResentDate: expectedDate,
			ResentFrom: []*mail.Address{
				{
					Name:    "郵便アリス",
					Address: "alice.yuubin@example.net",
				},
				{
					Name:    "郵便アリス",
					Address: "alice.yuubin@example.com",
				},
			},
			ResentSender: &mail.Address{
				Name:    "郵便アリス",
				Address: "alice.yuubin@example.net",
			},
			ResentTo: []*mail.Address{
				{
					Name:    "郵便ボブ",
					Address: "bob.yuubin@example.net",
				},
				{
					Name:    "郵便キャロル",
					Address: "carol.yuubin@example.net",
				},
			},
			ResentCc: []*mail.Address{
				{
					Name:    "郵便ダン",
					Address: "dan.yuubin@example.net",
				},
				{
					Name:    "郵便イーブ",
					Address: "eve.yubin@example.net",
				},
			},
			ResentBcc: []*mail.Address{
				{
					Name:    "郵便フランク",
					Address: "frank.yuubin@example.net",
				},
				{
					Name:    "郵便グレイス",
					Address: "grace.yubin@example.net",
				},
			},
			ResentMessageID: "Message-Id-1@example.net",
			ContentType: ContentTypeHeader{
				ContentType: "multipart/related",
				Params: map[string]string{
					"boundary": "RelatedBoundaryString",
					"charset":  "euc-jp",
				},
			},
			ExtraHeaders: map[string][]string{
				"X-Clacks-Overhead": {"GNU Terry Pratchett"},
				"X-Script/function/\t !\"#$%&'()*+,-./;<=>?@[\\]^_`{|}~": {
					"TEST VALUE 1\t !\"#$%&'()*+,-./0123456789:;<=>?@ABCDEFGHIJKLMNOPQRSTUVWXYZ[\\]^_` abcdefghijklmnopqrstuvwxyz{|}~",
					"TEST VALUE 2\t !\"#$%&'()*+,-./0123456789:;<=>?@ABCDEFGHIJKLMNOPQRSTUVWXYZ[\\]^_` abcdefghijklmnopqrstuvwxyz{|}~",
				},
			},
		},
		Text: `色は匂えど
散りぬるを
我が世誰ぞ
常ならん
有為の奥山
今日越えて
浅き夢見じ
酔いもせず。

Iro wa nioedo / Chirinuru o / Wa ga yo tare zo / Tsune naran / Ui no okuyama / Kyo koete Asaki yume miji / Yoi mo sezu.

とりなくこゑすゆめさませみよあけわたるひんかしをそらいろはえておきつへにほふねむれゐぬもやのうち。

天地星空
山川峰谷
雲霧室苔
人犬上末
硫黄猿生ふ為よ
榎の枝を馴れ居て。

田居に出で菜摘むわれをぞ君召すと求食り追ひゆく山城の打酔へる子ら藻葉干せよえ舟繋けぬ。`,
		EnrichedText: `<bold>色は匂えど</bold>
<italic>散りぬるを</italic>
<fixed>我が世誰ぞ</fixed>
<underline>常ならん</underline>
有為の奥山
今日越えて
浅き夢見じ
酔いもせず。

Iro wa nioedo / Chirinuru o / Wa ga yo tare zo / Tsune naran / Ui no okuyama / Kyo koete Asaki yume miji / Yoi mo sezu.

とりなくこゑすゆめさませみよあけわたるひんかしをそらいろはえておきつへにほふねむれゐぬもやのうち。

天地星空
山川峰谷
雲霧室苔
人犬上末
硫黄猿生ふ為よ
榎の枝を馴れ居て。

田居に出で菜摘むわれをぞ君召すと求食り追ひゆく山城の打酔へる子ら藻葉干せよえ舟繋けぬ。`,
		HTML: `<html>
<div dir="ltr">
<p>色は匂えど<br />
散りぬるを<br />
我が世誰ぞ<br />
常ならん<br />
有為の奥山<br />
今日越えて<br />
浅き夢見じ<br />
酔いもせず。</p>

<p>Iro wa nioedo / Chirinuru o / Wa ga yo tare zo / Tsune naran / Ui no okuyama / Kyo koete Asaki yume miji / Yoi mo sezu.</p>

<p>とりなくこゑすゆめさませみよあけわたるひんかしをそらいろはえておきつへにほふねむれゐぬもやのうち。</p>

<p>天地星空<br />
山川峰谷<br />
雲霧室苔<br />
人犬上末<br />
硫黄猿生ふ為よ<br />
榎の枝を馴れ居て。</p>

<p>田居に出で菜摘むわれをぞ君召すと求食り追ひゆく山城の打酔へる子ら藻葉干せよえ舟繋けぬ。</p>
</div>
</html>`,
		InlineFiles: []InlineFile{
			{
				ContentID: "inline-jpg-image.jpg@example.com",
				ContentType: ContentTypeHeader{
					ContentType: "image/jpeg",
					Params: map[string]string{
						"name": "inline-jpg-image-name.jpg",
					},
				},
				ContentDisposition: ContentDispositionHeader{
					ContentDisposition: inline,
					Params: map[string]string{
						"filename": "inline-jpg-image-filename.jpg",
					},
				},
				Data: []byte{255, 216, 255, 219, 0, 67, 0, 3, 2, 2, 2, 2, 2, 3, 2, 2, 2, 3, 3, 3, 3, 4,
					6, 4, 4, 4, 4, 4, 8, 6, 6, 5, 6, 9, 8, 10, 10, 9, 8, 9, 9, 10, 12, 15, 12, 10, 11, 14, 11, 9, 9,
					13, 17, 13, 14, 15, 16, 16, 17, 16, 10, 12, 18, 19, 18, 16, 19, 15, 16, 16, 16, 255, 201, 0, 11, 8,
					0, 1, 0, 1, 1, 1, 17, 0, 255, 204, 0, 6, 0, 16, 16, 5, 255, 218, 0, 8, 1, 1, 0, 0, 63, 0, 210, 207,
					32, 255, 217},
			},
		},
	}

	testEmailFromFile(t, fp, expectedEmail)
}
func TestParseEmailJapaneseMultipartMixedUtf8Over7bit(t *testing.T) {
	fp := "tests/test_japanese_multipart_mixed_utf-8_over_7bit.txt"
	tz, _ := time.LoadLocation("Asia/Tokyo")
	expectedDate, _ := time.Parse(
		time.RFC1123Z+" (MST)",
		time.Date(2019, time.April, 1, 7, 55, 0, 0, tz).Format(time.RFC1123Z+" (MST)"))
	expectedEmail := Email{
		Headers: Headers{
			Date:    expectedDate,
			Subject: "📧 Test いろは歌",
			ReplyTo: []*mail.Address{
				{
					Name:    "郵便アリス",
					Address: "alice.yuubin@example.net",
				},
			},
			Sender: &mail.Address{
				Name:    "郵便アリス",
				Address: "alice.yuubin@example.com",
			},
			From: []*mail.Address{
				{
					Name:    "郵便アリス",
					Address: "alice.yuubin@example.com",
				},
				{
					Name:    "郵便アリス",
					Address: "alice.yuubin@example.net",
				},
			},
			To: []*mail.Address{
				{
					Name:    "郵便ボブ",
					Address: "bob.yuubin@example.com",
				},
				{
					Name:    "郵便キャロル",
					Address: "carol.yuubin@example.com",
				},
			},
			Cc: []*mail.Address{
				{
					Name:    "郵便ダン",
					Address: "dan.yuubin@example.com",
				},
				{
					Name:    "郵便イーブ",
					Address: "eve.yubin@example.com",
				},
			},
			Bcc: []*mail.Address{
				{
					Name:    "郵便フランク",
					Address: "frank.yuubin@example.com",
				},
				{
					Name:    "郵便グレイス",
					Address: "grace.yubin@example.com",
				},
			},
			MessageID:  "Message-Id-1@example.com",
			InReplyTo:  []MessageId{"Message-Id-0@example.com"},
			References: []MessageId{"Message-Id-0@example.com"},
			Comments:   "Message Header Comment",
			Keywords:   []string{"Keyword 1", "Keyword 2"},
			ResentDate: expectedDate,
			ResentFrom: []*mail.Address{
				{
					Name:    "郵便アリス",
					Address: "alice.yuubin@example.net",
				},
				{
					Name:    "郵便アリス",
					Address: "alice.yuubin@example.com",
				},
			},
			ResentSender: &mail.Address{
				Name:    "郵便アリス",
				Address: "alice.yuubin@example.net",
			},
			ResentTo: []*mail.Address{
				{
					Name:    "郵便ボブ",
					Address: "bob.yuubin@example.net",
				},
				{
					Name:    "郵便キャロル",
					Address: "carol.yuubin@example.net",
				},
			},
			ResentCc: []*mail.Address{
				{
					Name:    "郵便ダン",
					Address: "dan.yuubin@example.net",
				},
				{
					Name:    "郵便イーブ",
					Address: "eve.yubin@example.net",
				},
			},
			ResentBcc: []*mail.Address{
				{
					Name:    "郵便フランク",
					Address: "frank.yuubin@example.net",
				},
				{
					Name:    "郵便グレイス",
					Address: "grace.yubin@example.net",
				},
			},
			ResentMessageID: "Message-Id-1@example.net",
			ContentType: ContentTypeHeader{
				ContentType: "multipart/mixed",
				Params: map[string]string{
					"boundary": "MixedBoundaryString",
					"charset":  "utf-8",
				},
			},
			ExtraHeaders: map[string][]string{
				"X-Clacks-Overhead": {"GNU Terry Pratchett"},
				"X-Script/function/\t !\"#$%&'()*+,-./;<=>?@[\\]^_`{|}~": {
					"TEST VALUE 1\t !\"#$%&'()*+,-./0123456789:;<=>?@ABCDEFGHIJKLMNOPQRSTUVWXYZ[\\]^_` abcdefghijklmnopqrstuvwxyz{|}~",
					"TEST VALUE 2\t !\"#$%&'()*+,-./0123456789:;<=>?@ABCDEFGHIJKLMNOPQRSTUVWXYZ[\\]^_` abcdefghijklmnopqrstuvwxyz{|}~",
				},
			},
		},
		Text: `色は匂えど
散りぬるを
我が世誰ぞ
常ならん
有為の奥山
今日越えて
浅き夢見じ
酔いもせず。

Iro wa nioedo / Chirinuru o / Wa ga yo tare zo / Tsune naran / Ui no okuyama / Kyo koete Asaki yume miji / Yoi mo sezu.

とりなくこゑすゆめさませみよあけわたるひんかしをそらいろはえておきつへにほふねむれゐぬもやのうち。

天地星空
山川峰谷
雲霧室苔
人犬上末
硫黄猿生ふ為よ
榎の枝を馴れ居て。

田居に出で菜摘むわれをぞ君召すと求食り追ひゆく山城の打酔へる子ら藻葉干せよえ舟繋けぬ。`,
		EnrichedText: `<bold>色は匂えど</bold>
<italic>散りぬるを</italic>
<fixed>我が世誰ぞ</fixed>
<underline>常ならん</underline>
有為の奥山
今日越えて
浅き夢見じ
酔いもせず。

Iro wa nioedo / Chirinuru o / Wa ga yo tare zo / Tsune naran / Ui no okuyama / Kyo koete Asaki yume miji / Yoi mo sezu.

とりなくこゑすゆめさませみよあけわたるひんかしをそらいろはえておきつへにほふねむれゐぬもやのうち。

天地星空
山川峰谷
雲霧室苔
人犬上末
硫黄猿生ふ為よ
榎の枝を馴れ居て。

田居に出で菜摘むわれをぞ君召すと求食り追ひゆく山城の打酔へる子ら藻葉干せよえ舟繋けぬ。`,
		HTML: `<html>
<div dir="ltr">
<p>色は匂えど<br />
散りぬるを<br />
我が世誰ぞ<br />
常ならん<br />
有為の奥山<br />
今日越えて<br />
浅き夢見じ<br />
酔いもせず。</p>

<p>Iro wa nioedo / Chirinuru o / Wa ga yo tare zo / Tsune naran / Ui no okuyama / Kyo koete Asaki yume miji / Yoi mo sezu.</p>

<p>とりなくこゑすゆめさませみよあけわたるひんかしをそらいろはえておきつへにほふねむれゐぬもやのうち。</p>

<p>天地星空<br />
山川峰谷<br />
雲霧室苔<br />
人犬上末<br />
硫黄猿生ふ為よ<br />
榎の枝を馴れ居て。</p>

<p>田居に出で菜摘むわれをぞ君召すと求食り追ひゆく山城の打酔へる子ら藻葉干せよえ舟繋けぬ。</p>
</div>
</html>`,
		InlineFiles: []InlineFile{
			{
				ContentType: ContentTypeHeader{
					ContentType: "image/jpeg",
					Params: map[string]string{
						"name": "inline-jpg-image-without-disposition.jpg",
					},
				},
				ContentDisposition: ContentDispositionHeader{
					ContentDisposition: "",
					Params:             map[string]string(nil),
				},
				Data: []byte{255, 216, 255, 219, 0, 67, 0, 3, 2, 2, 2, 2, 2, 3, 2, 2, 2, 3, 3, 3, 3, 4,
					6, 4, 4, 4, 4, 4, 8, 6, 6, 5, 6, 9, 8, 10, 10, 9, 8, 9, 9, 10, 12, 15, 12, 10, 11, 14, 11, 9, 9,
					13, 17, 13, 14, 15, 16, 16, 17, 16, 10, 12, 18, 19, 18, 16, 19, 15, 16, 16, 16, 255, 201, 0, 11, 8,
					0, 1, 0, 1, 1, 1, 17, 0, 255, 204, 0, 6, 0, 16, 16, 5, 255, 218, 0, 8, 1, 1, 0, 0, 63, 0, 210, 207,
					32, 255, 217},
			},
			{
				ContentID: "inline-jpg-image.jpg@example.com",
				ContentType: ContentTypeHeader{
					ContentType: "image/jpeg",
					Params: map[string]string{
						"name": "inline-jpg-image-name.jpg",
					},
				},
				ContentDisposition: ContentDispositionHeader{
					ContentDisposition: inline,
					Params: map[string]string{
						"filename": "inline-jpg-image-filename.jpg",
					},
				},
				Data: []byte{255, 216, 255, 219, 0, 67, 0, 3, 2, 2, 2, 2, 2, 3, 2, 2, 2, 3, 3, 3, 3, 4,
					6, 4, 4, 4, 4, 4, 8, 6, 6, 5, 6, 9, 8, 10, 10, 9, 8, 9, 9, 10, 12, 15, 12, 10, 11, 14, 11, 9, 9,
					13, 17, 13, 14, 15, 16, 16, 17, 16, 10, 12, 18, 19, 18, 16, 19, 15, 16, 16, 16, 255, 201, 0, 11, 8,
					0, 1, 0, 1, 1, 1, 17, 0, 255, 204, 0, 6, 0, 16, 16, 5, 255, 218, 0, 8, 1, 1, 0, 0, 63, 0, 210, 207,
					32, 255, 217},
			},
		},
		AttachedFiles: []AttachedFile{
			{
				ContentType: ContentTypeHeader{
					ContentType: "application/pdf",
					Params: map[string]string{
						"name": "attached-pdf-name.pdf",
					},
				},
				ContentDisposition: ContentDispositionHeader{
					ContentDisposition: attachment,
					Params: map[string]string{
						"filename": "attached-pdf-filename.pdf",
					},
				},
				Data: []byte{37, 80, 68, 70, 45, 49, 46, 13, 116, 114, 97, 105, 108, 101, 114, 60, 60,
					47, 82, 111, 111, 116, 60, 60, 47, 80, 97, 103, 101, 115, 60, 60, 47, 75, 105, 100, 115, 91, 60,
					60, 47, 77, 101, 100, 105, 97, 66, 111, 120, 91, 48, 32, 48, 32, 51, 32, 51, 93, 62, 62, 93, 62,
					62, 62, 62, 62, 62},
			},
			{
				ContentType: ContentTypeHeader{
					ContentType: "application/pdf",
					Params: map[string]string{
						"name": "attached-pdf-without-disposition.pdf",
					},
				},
				ContentDisposition: ContentDispositionHeader{
					ContentDisposition: "",
					Params:             map[string]string(nil),
				},
				Data: []byte{37, 80, 68, 70, 45, 49, 46, 13, 116, 114, 97, 105, 108, 101, 114, 60, 60,
					47, 82, 111, 111, 116, 60, 60, 47, 80, 97, 103, 101, 115, 60, 60, 47, 75, 105, 100, 115, 91, 60,
					60, 47, 77, 101, 100, 105, 97, 66, 111, 120, 91, 48, 32, 48, 32, 51, 32, 51, 93, 62, 62, 93, 62,
					62, 62, 62, 62, 62},
			},
			{
				ContentType: ContentTypeHeader{
					ContentType: "application/json",
					Params: map[string]string{
						"name": "attached-json-name.json",
					},
				},
				ContentDisposition: ContentDispositionHeader{
					ContentDisposition: attachment,
					Params: map[string]string{
						"filename": "attached-json-filename.json",
					},
				},
				Data: []byte{123, 34, 102, 111, 111, 34, 58, 34, 98, 97, 114, 34, 125},
			},
			{
				ContentType: ContentTypeHeader{
					ContentType: "text/plain",
					Params: map[string]string{
						"name": "attached-text-plain-name.txt",
					},
				},
				ContentDisposition: ContentDispositionHeader{
					ContentDisposition: attachment,
					Params: map[string]string{
						"filename": "attached-text-plain-filename.txt",
					},
				},
				Data: []byte{84, 101, 120, 116, 47, 112, 108, 97, 105, 110, 32, 99, 111, 110, 116, 101, 110, 116, 32,
					97, 115, 32, 97, 110, 32, 97, 116, 116, 97, 99, 104, 101, 100, 32, 46, 116, 120, 116, 32, 102, 105,
					108, 101, 46},
			},
			{
				ContentType: ContentTypeHeader{
					ContentType: "text/html",
					Params: map[string]string{
						"name": "attached-text-html-name.html",
					},
				},
				ContentDisposition: ContentDispositionHeader{
					ContentDisposition: attachment,
					Params: map[string]string{
						"filename": "attached-text-html-filename.html",
					},
				},
				Data: []byte{84, 101, 120, 116, 47, 104, 116, 109, 108, 32, 99, 111, 110, 116, 101, 110, 116, 32, 97,
					115, 32, 97, 110, 32, 97, 116, 116, 97, 99, 104, 101, 100, 32, 46, 104, 116, 109, 108, 32, 102,
					105, 108, 101, 46},
			},
		},
	}

	testEmailFromFile(t, fp, expectedEmail)
}
func TestParseEmailJapaneseMultipartMixedUtf8OverBase64(t *testing.T) {
	fp := "tests/test_japanese_multipart_mixed_utf-8_over_base64.txt"
	tz, _ := time.LoadLocation("Asia/Tokyo")
	expectedDate, _ := time.Parse(
		time.RFC1123Z+" (MST)",
		time.Date(2019, time.April, 1, 7, 55, 0, 0, tz).Format(time.RFC1123Z+" (MST)"))
	expectedEmail := Email{
		Headers: Headers{
			Date:    expectedDate,
			Subject: "📧 Test いろは歌",
			ReplyTo: []*mail.Address{
				{
					Name:    "郵便アリス",
					Address: "alice.yuubin@example.net",
				},
			},
			Sender: &mail.Address{
				Name:    "郵便アリス",
				Address: "alice.yuubin@example.com",
			},
			From: []*mail.Address{
				{
					Name:    "郵便アリス",
					Address: "alice.yuubin@example.com",
				},
				{
					Name:    "郵便アリス",
					Address: "alice.yuubin@example.net",
				},
			},
			To: []*mail.Address{
				{
					Name:    "郵便ボブ",
					Address: "bob.yuubin@example.com",
				},
				{
					Name:    "郵便キャロル",
					Address: "carol.yuubin@example.com",
				},
			},
			Cc: []*mail.Address{
				{
					Name:    "郵便ダン",
					Address: "dan.yuubin@example.com",
				},
				{
					Name:    "郵便イーブ",
					Address: "eve.yubin@example.com",
				},
			},
			Bcc: []*mail.Address{
				{
					Name:    "郵便フランク",
					Address: "frank.yuubin@example.com",
				},
				{
					Name:    "郵便グレイス",
					Address: "grace.yubin@example.com",
				},
			},
			MessageID:  "Message-Id-1@example.com",
			InReplyTo:  []MessageId{"Message-Id-0@example.com"},
			References: []MessageId{"Message-Id-0@example.com"},
			Comments:   "Message Header Comment",
			Keywords:   []string{"Keyword 1", "Keyword 2"},
			ResentDate: expectedDate,
			ResentFrom: []*mail.Address{
				{
					Name:    "郵便アリス",
					Address: "alice.yuubin@example.net",
				},
				{
					Name:    "郵便アリス",
					Address: "alice.yuubin@example.com",
				},
			},
			ResentSender: &mail.Address{
				Name:    "郵便アリス",
				Address: "alice.yuubin@example.net",
			},
			ResentTo: []*mail.Address{
				{
					Name:    "郵便ボブ",
					Address: "bob.yuubin@example.net",
				},
				{
					Name:    "郵便キャロル",
					Address: "carol.yuubin@example.net",
				},
			},
			ResentCc: []*mail.Address{
				{
					Name:    "郵便ダン",
					Address: "dan.yuubin@example.net",
				},
				{
					Name:    "郵便イーブ",
					Address: "eve.yubin@example.net",
				},
			},
			ResentBcc: []*mail.Address{
				{
					Name:    "郵便フランク",
					Address: "frank.yuubin@example.net",
				},
				{
					Name:    "郵便グレイス",
					Address: "grace.yubin@example.net",
				},
			},
			ResentMessageID: "Message-Id-1@example.net",
			ContentType: ContentTypeHeader{
				ContentType: "multipart/mixed",
				Params: map[string]string{
					"boundary": "MixedBoundaryString",
					"charset":  "utf-8",
				},
			},
			ExtraHeaders: map[string][]string{
				"X-Clacks-Overhead": {"GNU Terry Pratchett"},
				"X-Script/function/\t !\"#$%&'()*+,-./;<=>?@[\\]^_`{|}~": {
					"TEST VALUE 1\t !\"#$%&'()*+,-./0123456789:;<=>?@ABCDEFGHIJKLMNOPQRSTUVWXYZ[\\]^_` abcdefghijklmnopqrstuvwxyz{|}~",
					"TEST VALUE 2\t !\"#$%&'()*+,-./0123456789:;<=>?@ABCDEFGHIJKLMNOPQRSTUVWXYZ[\\]^_` abcdefghijklmnopqrstuvwxyz{|}~",
				},
			},
		},
		Text: `色は匂えど
散りぬるを
我が世誰ぞ
常ならん
有為の奥山
今日越えて
浅き夢見じ
酔いもせず。

Iro wa nioedo / Chirinuru o / Wa ga yo tare zo / Tsune naran / Ui no okuyama / Kyo koete Asaki yume miji / Yoi mo sezu.

とりなくこゑすゆめさませみよあけわたるひんかしをそらいろはえておきつへにほふねむれゐぬもやのうち。

天地星空
山川峰谷
雲霧室苔
人犬上末
硫黄猿生ふ為よ
榎の枝を馴れ居て。

田居に出で菜摘むわれをぞ君召すと求食り追ひゆく山城の打酔へる子ら藻葉干せよえ舟繋けぬ。`,
		EnrichedText: `<bold>色は匂えど</bold>
<italic>散りぬるを</italic>
<fixed>我が世誰ぞ</fixed>
<underline>常ならん</underline>
有為の奥山
今日越えて
浅き夢見じ
酔いもせず。

Iro wa nioedo / Chirinuru o / Wa ga yo tare zo / Tsune naran / Ui no okuyama / Kyo koete Asaki yume miji / Yoi mo sezu.

とりなくこゑすゆめさませみよあけわたるひんかしをそらいろはえておきつへにほふねむれゐぬもやのうち。

天地星空
山川峰谷
雲霧室苔
人犬上末
硫黄猿生ふ為よ
榎の枝を馴れ居て。

田居に出で菜摘むわれをぞ君召すと求食り追ひゆく山城の打酔へる子ら藻葉干せよえ舟繋けぬ。`,
		HTML: `<html>
<div dir="ltr">
<p>色は匂えど<br />
散りぬるを<br />
我が世誰ぞ<br />
常ならん<br />
有為の奥山<br />
今日越えて<br />
浅き夢見じ<br />
酔いもせず。</p>

<p>Iro wa nioedo / Chirinuru o / Wa ga yo tare zo / Tsune naran / Ui no okuyama / Kyo koete Asaki yume miji / Yoi mo sezu.</p>

<p>とりなくこゑすゆめさませみよあけわたるひんかしをそらいろはえておきつへにほふねむれゐぬもやのうち。</p>

<p>天地星空<br />
山川峰谷<br />
雲霧室苔<br />
人犬上末<br />
硫黄猿生ふ為よ<br />
榎の枝を馴れ居て。</p>

<p>田居に出で菜摘むわれをぞ君召すと求食り追ひゆく山城の打酔へる子ら藻葉干せよえ舟繋けぬ。</p>
</div>
</html>`,
		InlineFiles: []InlineFile{
			{
				ContentType: ContentTypeHeader{
					ContentType: "image/jpeg",
					Params: map[string]string{
						"name": "inline-jpg-image-without-disposition.jpg",
					},
				},
				ContentDisposition: ContentDispositionHeader{
					ContentDisposition: "",
					Params:             map[string]string(nil),
				},
				Data: []byte{255, 216, 255, 219, 0, 67, 0, 3, 2, 2, 2, 2, 2, 3, 2, 2, 2, 3, 3, 3, 3, 4,
					6, 4, 4, 4, 4, 4, 8, 6, 6, 5, 6, 9, 8, 10, 10, 9, 8, 9, 9, 10, 12, 15, 12, 10, 11, 14, 11, 9, 9,
					13, 17, 13, 14, 15, 16, 16, 17, 16, 10, 12, 18, 19, 18, 16, 19, 15, 16, 16, 16, 255, 201, 0, 11, 8,
					0, 1, 0, 1, 1, 1, 17, 0, 255, 204, 0, 6, 0, 16, 16, 5, 255, 218, 0, 8, 1, 1, 0, 0, 63, 0, 210, 207,
					32, 255, 217},
			},
			{
				ContentID: "inline-jpg-image.jpg@example.com",
				ContentType: ContentTypeHeader{
					ContentType: "image/jpeg",
					Params: map[string]string{
						"name": "inline-jpg-image-name.jpg",
					},
				},
				ContentDisposition: ContentDispositionHeader{
					ContentDisposition: inline,
					Params: map[string]string{
						"filename": "inline-jpg-image-filename.jpg",
					},
				},
				Data: []byte{255, 216, 255, 219, 0, 67, 0, 3, 2, 2, 2, 2, 2, 3, 2, 2, 2, 3, 3, 3, 3, 4,
					6, 4, 4, 4, 4, 4, 8, 6, 6, 5, 6, 9, 8, 10, 10, 9, 8, 9, 9, 10, 12, 15, 12, 10, 11, 14, 11, 9, 9,
					13, 17, 13, 14, 15, 16, 16, 17, 16, 10, 12, 18, 19, 18, 16, 19, 15, 16, 16, 16, 255, 201, 0, 11, 8,
					0, 1, 0, 1, 1, 1, 17, 0, 255, 204, 0, 6, 0, 16, 16, 5, 255, 218, 0, 8, 1, 1, 0, 0, 63, 0, 210, 207,
					32, 255, 217},
			},
		},
		AttachedFiles: []AttachedFile{
			{
				ContentType: ContentTypeHeader{
					ContentType: "application/pdf",
					Params: map[string]string{
						"name": "attached-pdf-name.pdf",
					},
				},
				ContentDisposition: ContentDispositionHeader{
					ContentDisposition: attachment,
					Params: map[string]string{
						"filename": "attached-pdf-filename.pdf",
					},
				},
				Data: []byte{37, 80, 68, 70, 45, 49, 46, 13, 116, 114, 97, 105, 108, 101, 114, 60, 60,
					47, 82, 111, 111, 116, 60, 60, 47, 80, 97, 103, 101, 115, 60, 60, 47, 75, 105, 100, 115, 91, 60,
					60, 47, 77, 101, 100, 105, 97, 66, 111, 120, 91, 48, 32, 48, 32, 51, 32, 51, 93, 62, 62, 93, 62,
					62, 62, 62, 62, 62},
			},
			{
				ContentType: ContentTypeHeader{
					ContentType: "application/pdf",
					Params: map[string]string{
						"name": "attached-pdf-without-disposition.pdf",
					},
				},
				ContentDisposition: ContentDispositionHeader{
					ContentDisposition: "",
					Params:             map[string]string(nil),
				},
				Data: []byte{37, 80, 68, 70, 45, 49, 46, 13, 116, 114, 97, 105, 108, 101, 114, 60, 60,
					47, 82, 111, 111, 116, 60, 60, 47, 80, 97, 103, 101, 115, 60, 60, 47, 75, 105, 100, 115, 91, 60,
					60, 47, 77, 101, 100, 105, 97, 66, 111, 120, 91, 48, 32, 48, 32, 51, 32, 51, 93, 62, 62, 93, 62,
					62, 62, 62, 62, 62},
			},
			{
				ContentType: ContentTypeHeader{
					ContentType: "application/json",
					Params: map[string]string{
						"name": "attached-json-name.json",
					},
				},
				ContentDisposition: ContentDispositionHeader{
					ContentDisposition: attachment,
					Params: map[string]string{
						"filename": "attached-json-filename.json",
					},
				},
				Data: []byte{123, 34, 102, 111, 111, 34, 58, 34, 98, 97, 114, 34, 125},
			},
			{
				ContentType: ContentTypeHeader{
					ContentType: "text/plain",
					Params: map[string]string{
						"name": "attached-text-plain-name.txt",
					},
				},
				ContentDisposition: ContentDispositionHeader{
					ContentDisposition: attachment,
					Params: map[string]string{
						"filename": "attached-text-plain-filename.txt",
					},
				},
				Data: []byte{84, 101, 120, 116, 47, 112, 108, 97, 105, 110, 32, 99, 111, 110, 116, 101, 110, 116, 32,
					97, 115, 32, 97, 110, 32, 97, 116, 116, 97, 99, 104, 101, 100, 32, 46, 116, 120, 116, 32, 102, 105,
					108, 101, 46},
			},
			{
				ContentType: ContentTypeHeader{
					ContentType: "text/html",
					Params: map[string]string{
						"name": "attached-text-html-name.html",
					},
				},
				ContentDisposition: ContentDispositionHeader{
					ContentDisposition: attachment,
					Params: map[string]string{
						"filename": "attached-text-html-filename.html",
					},
				},
				Data: []byte{84, 101, 120, 116, 47, 104, 116, 109, 108, 32, 99, 111, 110, 116, 101, 110, 116, 32, 97,
					115, 32, 97, 110, 32, 97, 116, 116, 97, 99, 104, 101, 100, 32, 46, 104, 116, 109, 108, 32, 102,
					105, 108, 101, 46},
			},
		},
	}

	testEmailFromFile(t, fp, expectedEmail)
}
func TestParseEmailJapaneseMultipartMixedUtf8OverQuotedprintable(t *testing.T) {
	fp := "tests/test_japanese_multipart_mixed_utf-8_over_quoted-printable.txt"
	tz, _ := time.LoadLocation("Asia/Tokyo")
	expectedDate, _ := time.Parse(
		time.RFC1123Z+" (MST)",
		time.Date(2019, time.April, 1, 7, 55, 0, 0, tz).Format(time.RFC1123Z+" (MST)"))
	expectedEmail := Email{
		Headers: Headers{
			Date:    expectedDate,
			Subject: "📧 Test いろは歌",
			ReplyTo: []*mail.Address{
				{
					Name:    "郵便アリス",
					Address: "alice.yuubin@example.net",
				},
			},
			Sender: &mail.Address{
				Name:    "郵便アリス",
				Address: "alice.yuubin@example.com",
			},
			From: []*mail.Address{
				{
					Name:    "郵便アリス",
					Address: "alice.yuubin@example.com",
				},
				{
					Name:    "郵便アリス",
					Address: "alice.yuubin@example.net",
				},
			},
			To: []*mail.Address{
				{
					Name:    "郵便ボブ",
					Address: "bob.yuubin@example.com",
				},
				{
					Name:    "郵便キャロル",
					Address: "carol.yuubin@example.com",
				},
			},
			Cc: []*mail.Address{
				{
					Name:    "郵便ダン",
					Address: "dan.yuubin@example.com",
				},
				{
					Name:    "郵便イーブ",
					Address: "eve.yubin@example.com",
				},
			},
			Bcc: []*mail.Address{
				{
					Name:    "郵便フランク",
					Address: "frank.yuubin@example.com",
				},
				{
					Name:    "郵便グレイス",
					Address: "grace.yubin@example.com",
				},
			},
			MessageID:  "Message-Id-1@example.com",
			InReplyTo:  []MessageId{"Message-Id-0@example.com"},
			References: []MessageId{"Message-Id-0@example.com"},
			Comments:   "Message Header Comment",
			Keywords:   []string{"Keyword 1", "Keyword 2"},
			ResentDate: expectedDate,
			ResentFrom: []*mail.Address{
				{
					Name:    "郵便アリス",
					Address: "alice.yuubin@example.net",
				},
				{
					Name:    "郵便アリス",
					Address: "alice.yuubin@example.com",
				},
			},
			ResentSender: &mail.Address{
				Name:    "郵便アリス",
				Address: "alice.yuubin@example.net",
			},
			ResentTo: []*mail.Address{
				{
					Name:    "郵便ボブ",
					Address: "bob.yuubin@example.net",
				},
				{
					Name:    "郵便キャロル",
					Address: "carol.yuubin@example.net",
				},
			},
			ResentCc: []*mail.Address{
				{
					Name:    "郵便ダン",
					Address: "dan.yuubin@example.net",
				},
				{
					Name:    "郵便イーブ",
					Address: "eve.yubin@example.net",
				},
			},
			ResentBcc: []*mail.Address{
				{
					Name:    "郵便フランク",
					Address: "frank.yuubin@example.net",
				},
				{
					Name:    "郵便グレイス",
					Address: "grace.yubin@example.net",
				},
			},
			ResentMessageID: "Message-Id-1@example.net",
			ContentType: ContentTypeHeader{
				ContentType: "multipart/mixed",
				Params: map[string]string{
					"boundary": "MixedBoundaryString",
					"charset":  "utf-8",
				},
			},
			ExtraHeaders: map[string][]string{
				"X-Clacks-Overhead": {"GNU Terry Pratchett"},
				"X-Script/function/\t !\"#$%&'()*+,-./;<=>?@[\\]^_`{|}~": {
					"TEST VALUE 1\t !\"#$%&'()*+,-./0123456789:;<=>?@ABCDEFGHIJKLMNOPQRSTUVWXYZ[\\]^_` abcdefghijklmnopqrstuvwxyz{|}~",
					"TEST VALUE 2\t !\"#$%&'()*+,-./0123456789:;<=>?@ABCDEFGHIJKLMNOPQRSTUVWXYZ[\\]^_` abcdefghijklmnopqrstuvwxyz{|}~",
				},
			},
		},
		Text: `色は匂えど
散りぬるを
我が世誰ぞ
常ならん
有為の奥山
今日越えて
浅き夢見じ
酔いもせず。

Iro wa nioedo / Chirinuru o / Wa ga yo tare zo / Tsune naran / Ui no okuyama / Kyo koete Asaki yume miji / Yoi mo sezu.

とりなくこゑすゆめさませみよあけわたるひんかしをそらいろはえておきつへにほふねむれゐぬもやのうち。

天地星空
山川峰谷
雲霧室苔
人犬上末
硫黄猿生ふ為よ
榎の枝を馴れ居て。

田居に出で菜摘むわれをぞ君召すと求食り追ひゆく山城の打酔へる子ら藻葉干せよえ舟繋けぬ。`,
		EnrichedText: `<bold>色は匂えど</bold>
<italic>散りぬるを</italic>
<fixed>我が世誰ぞ</fixed>
<underline>常ならん</underline>
有為の奥山
今日越えて
浅き夢見じ
酔いもせず。

Iro wa nioedo / Chirinuru o / Wa ga yo tare zo / Tsune naran / Ui no okuyama / Kyo koete Asaki yume miji / Yoi mo sezu.

とりなくこゑすゆめさませみよあけわたるひんかしをそらいろはえておきつへにほふねむれゐぬもやのうち。

天地星空
山川峰谷
雲霧室苔
人犬上末
硫黄猿生ふ為よ
榎の枝を馴れ居て。

田居に出で菜摘むわれをぞ君召すと求食り追ひゆく山城の打酔へる子ら藻葉干せよえ舟繋けぬ。`,
		HTML: `<html>
<div dir="ltr">
<p>色は匂えど<br />
散りぬるを<br />
我が世誰ぞ<br />
常ならん<br />
有為の奥山<br />
今日越えて<br />
浅き夢見じ<br />
酔いもせず。</p>

<p>Iro wa nioedo / Chirinuru o / Wa ga yo tare zo / Tsune naran / Ui no okuyama / Kyo koete Asaki yume miji / Yoi mo sezu.</p>

<p>とりなくこゑすゆめさませみよあけわたるひんかしをそらいろはえておきつへにほふねむれゐぬもやのうち。</p>

<p>天地星空<br />
山川峰谷<br />
雲霧室苔<br />
人犬上末<br />
硫黄猿生ふ為よ<br />
榎の枝を馴れ居て。</p>

<p>田居に出で菜摘むわれをぞ君召すと求食り追ひゆく山城の打酔へる子ら藻葉干せよえ舟繋けぬ。</p>
</div>
</html>`,
		InlineFiles: []InlineFile{
			{
				ContentType: ContentTypeHeader{
					ContentType: "image/jpeg",
					Params: map[string]string{
						"name": "inline-jpg-image-without-disposition.jpg",
					},
				},
				ContentDisposition: ContentDispositionHeader{
					ContentDisposition: "",
					Params:             map[string]string(nil),
				},
				Data: []byte{255, 216, 255, 219, 0, 67, 0, 3, 2, 2, 2, 2, 2, 3, 2, 2, 2, 3, 3, 3, 3, 4,
					6, 4, 4, 4, 4, 4, 8, 6, 6, 5, 6, 9, 8, 10, 10, 9, 8, 9, 9, 10, 12, 15, 12, 10, 11, 14, 11, 9, 9,
					13, 17, 13, 14, 15, 16, 16, 17, 16, 10, 12, 18, 19, 18, 16, 19, 15, 16, 16, 16, 255, 201, 0, 11, 8,
					0, 1, 0, 1, 1, 1, 17, 0, 255, 204, 0, 6, 0, 16, 16, 5, 255, 218, 0, 8, 1, 1, 0, 0, 63, 0, 210, 207,
					32, 255, 217},
			},
			{
				ContentID: "inline-jpg-image.jpg@example.com",
				ContentType: ContentTypeHeader{
					ContentType: "image/jpeg",
					Params: map[string]string{
						"name": "inline-jpg-image-name.jpg",
					},
				},
				ContentDisposition: ContentDispositionHeader{
					ContentDisposition: inline,
					Params: map[string]string{
						"filename": "inline-jpg-image-filename.jpg",
					},
				},
				Data: []byte{255, 216, 255, 219, 0, 67, 0, 3, 2, 2, 2, 2, 2, 3, 2, 2, 2, 3, 3, 3, 3, 4,
					6, 4, 4, 4, 4, 4, 8, 6, 6, 5, 6, 9, 8, 10, 10, 9, 8, 9, 9, 10, 12, 15, 12, 10, 11, 14, 11, 9, 9,
					13, 17, 13, 14, 15, 16, 16, 17, 16, 10, 12, 18, 19, 18, 16, 19, 15, 16, 16, 16, 255, 201, 0, 11, 8,
					0, 1, 0, 1, 1, 1, 17, 0, 255, 204, 0, 6, 0, 16, 16, 5, 255, 218, 0, 8, 1, 1, 0, 0, 63, 0, 210, 207,
					32, 255, 217},
			},
		},
		AttachedFiles: []AttachedFile{
			{
				ContentType: ContentTypeHeader{
					ContentType: "application/pdf",
					Params: map[string]string{
						"name": "attached-pdf-name.pdf",
					},
				},
				ContentDisposition: ContentDispositionHeader{
					ContentDisposition: attachment,
					Params: map[string]string{
						"filename": "attached-pdf-filename.pdf",
					},
				},
				Data: []byte{37, 80, 68, 70, 45, 49, 46, 13, 116, 114, 97, 105, 108, 101, 114, 60, 60,
					47, 82, 111, 111, 116, 60, 60, 47, 80, 97, 103, 101, 115, 60, 60, 47, 75, 105, 100, 115, 91, 60,
					60, 47, 77, 101, 100, 105, 97, 66, 111, 120, 91, 48, 32, 48, 32, 51, 32, 51, 93, 62, 62, 93, 62,
					62, 62, 62, 62, 62},
			},
			{
				ContentType: ContentTypeHeader{
					ContentType: "application/pdf",
					Params: map[string]string{
						"name": "attached-pdf-without-disposition.pdf",
					},
				},
				ContentDisposition: ContentDispositionHeader{
					ContentDisposition: "",
					Params:             map[string]string(nil),
				},
				Data: []byte{37, 80, 68, 70, 45, 49, 46, 13, 116, 114, 97, 105, 108, 101, 114, 60, 60,
					47, 82, 111, 111, 116, 60, 60, 47, 80, 97, 103, 101, 115, 60, 60, 47, 75, 105, 100, 115, 91, 60,
					60, 47, 77, 101, 100, 105, 97, 66, 111, 120, 91, 48, 32, 48, 32, 51, 32, 51, 93, 62, 62, 93, 62,
					62, 62, 62, 62, 62},
			},
			{
				ContentType: ContentTypeHeader{
					ContentType: "application/json",
					Params: map[string]string{
						"name": "attached-json-name.json",
					},
				},
				ContentDisposition: ContentDispositionHeader{
					ContentDisposition: attachment,
					Params: map[string]string{
						"filename": "attached-json-filename.json",
					},
				},
				Data: []byte{123, 34, 102, 111, 111, 34, 58, 34, 98, 97, 114, 34, 125},
			},
			{
				ContentType: ContentTypeHeader{
					ContentType: "text/plain",
					Params: map[string]string{
						"name": "attached-text-plain-name.txt",
					},
				},
				ContentDisposition: ContentDispositionHeader{
					ContentDisposition: attachment,
					Params: map[string]string{
						"filename": "attached-text-plain-filename.txt",
					},
				},
				Data: []byte{84, 101, 120, 116, 47, 112, 108, 97, 105, 110, 32, 99, 111, 110, 116, 101, 110, 116, 32,
					97, 115, 32, 97, 110, 32, 97, 116, 116, 97, 99, 104, 101, 100, 32, 46, 116, 120, 116, 32, 102, 105,
					108, 101, 46},
			},
			{
				ContentType: ContentTypeHeader{
					ContentType: "text/html",
					Params: map[string]string{
						"name": "attached-text-html-name.html",
					},
				},
				ContentDisposition: ContentDispositionHeader{
					ContentDisposition: attachment,
					Params: map[string]string{
						"filename": "attached-text-html-filename.html",
					},
				},
				Data: []byte{84, 101, 120, 116, 47, 104, 116, 109, 108, 32, 99, 111, 110, 116, 101, 110, 116, 32, 97,
					115, 32, 97, 110, 32, 97, 116, 116, 97, 99, 104, 101, 100, 32, 46, 104, 116, 109, 108, 32, 102,
					105, 108, 101, 46},
			},
		},
	}

	testEmailFromFile(t, fp, expectedEmail)
}
func TestParseEmailJapaneseMultipartMixedIso2022jpOver7bit(t *testing.T) {
	fp := "tests/test_japanese_multipart_mixed_iso-2022-jp_over_7bit.txt"
	tz, _ := time.LoadLocation("Asia/Tokyo")
	expectedDate, _ := time.Parse(
		time.RFC1123Z+" (MST)",
		time.Date(2019, time.April, 1, 7, 55, 0, 0, tz).Format(time.RFC1123Z+" (MST)"))
	expectedEmail := Email{
		Headers: Headers{
			Date:    expectedDate,
			Subject: "📧 Test いろは歌",
			ReplyTo: []*mail.Address{
				{
					Name:    "郵便アリス",
					Address: "alice.yuubin@example.net",
				},
			},
			Sender: &mail.Address{
				Name:    "郵便アリス",
				Address: "alice.yuubin@example.com",
			},
			From: []*mail.Address{
				{
					Name:    "郵便アリス",
					Address: "alice.yuubin@example.com",
				},
				{
					Name:    "郵便アリス",
					Address: "alice.yuubin@example.net",
				},
			},
			To: []*mail.Address{
				{
					Name:    "郵便ボブ",
					Address: "bob.yuubin@example.com",
				},
				{
					Name:    "郵便キャロル",
					Address: "carol.yuubin@example.com",
				},
			},
			Cc: []*mail.Address{
				{
					Name:    "郵便ダン",
					Address: "dan.yuubin@example.com",
				},
				{
					Name:    "郵便イーブ",
					Address: "eve.yubin@example.com",
				},
			},
			Bcc: []*mail.Address{
				{
					Name:    "郵便フランク",
					Address: "frank.yuubin@example.com",
				},
				{
					Name:    "郵便グレイス",
					Address: "grace.yubin@example.com",
				},
			},
			MessageID:  "Message-Id-1@example.com",
			InReplyTo:  []MessageId{"Message-Id-0@example.com"},
			References: []MessageId{"Message-Id-0@example.com"},
			Comments:   "Message Header Comment",
			Keywords:   []string{"Keyword 1", "Keyword 2"},
			ResentDate: expectedDate,
			ResentFrom: []*mail.Address{
				{
					Name:    "郵便アリス",
					Address: "alice.yuubin@example.net",
				},
				{
					Name:    "郵便アリス",
					Address: "alice.yuubin@example.com",
				},
			},
			ResentSender: &mail.Address{
				Name:    "郵便アリス",
				Address: "alice.yuubin@example.net",
			},
			ResentTo: []*mail.Address{
				{
					Name:    "郵便ボブ",
					Address: "bob.yuubin@example.net",
				},
				{
					Name:    "郵便キャロル",
					Address: "carol.yuubin@example.net",
				},
			},
			ResentCc: []*mail.Address{
				{
					Name:    "郵便ダン",
					Address: "dan.yuubin@example.net",
				},
				{
					Name:    "郵便イーブ",
					Address: "eve.yubin@example.net",
				},
			},
			ResentBcc: []*mail.Address{
				{
					Name:    "郵便フランク",
					Address: "frank.yuubin@example.net",
				},
				{
					Name:    "郵便グレイス",
					Address: "grace.yubin@example.net",
				},
			},
			ResentMessageID: "Message-Id-1@example.net",
			ContentType: ContentTypeHeader{
				ContentType: "multipart/mixed",
				Params: map[string]string{
					"boundary": "MixedBoundaryString",
					"charset":  "iso-2022-jp",
				},
			},
			ExtraHeaders: map[string][]string{
				"X-Clacks-Overhead": {"GNU Terry Pratchett"},
				"X-Script/function/\t !\"#$%&'()*+,-./;<=>?@[\\]^_`{|}~": {
					"TEST VALUE 1\t !\"#$%&'()*+,-./0123456789:;<=>?@ABCDEFGHIJKLMNOPQRSTUVWXYZ[\\]^_` abcdefghijklmnopqrstuvwxyz{|}~",
					"TEST VALUE 2\t !\"#$%&'()*+,-./0123456789:;<=>?@ABCDEFGHIJKLMNOPQRSTUVWXYZ[\\]^_` abcdefghijklmnopqrstuvwxyz{|}~",
				},
			},
		},
		Text: `色は匂えど
散りぬるを
我が世誰ぞ
常ならん
有為の奥山
今日越えて
浅き夢見じ
酔いもせず。

Iro wa nioedo / Chirinuru o / Wa ga yo tare zo / Tsune naran / Ui no okuyama / Kyo koete Asaki yume miji / Yoi mo sezu.

とりなくこゑすゆめさませみよあけわたるひんかしをそらいろはえておきつへにほふねむれゐぬもやのうち。

天地星空
山川峰谷
雲霧室苔
人犬上末
硫黄猿生ふ為よ
榎の枝を馴れ居て。

田居に出で菜摘むわれをぞ君召すと求食り追ひゆく山城の打酔へる子ら藻葉干せよえ舟繋けぬ。`,
		EnrichedText: `<bold>色は匂えど</bold>
<italic>散りぬるを</italic>
<fixed>我が世誰ぞ</fixed>
<underline>常ならん</underline>
有為の奥山
今日越えて
浅き夢見じ
酔いもせず。

Iro wa nioedo / Chirinuru o / Wa ga yo tare zo / Tsune naran / Ui no okuyama / Kyo koete Asaki yume miji / Yoi mo sezu.

とりなくこゑすゆめさませみよあけわたるひんかしをそらいろはえておきつへにほふねむれゐぬもやのうち。

天地星空
山川峰谷
雲霧室苔
人犬上末
硫黄猿生ふ為よ
榎の枝を馴れ居て。

田居に出で菜摘むわれをぞ君召すと求食り追ひゆく山城の打酔へる子ら藻葉干せよえ舟繋けぬ。`,
		HTML: `<html>
<div dir="ltr">
<p>色は匂えど<br />
散りぬるを<br />
我が世誰ぞ<br />
常ならん<br />
有為の奥山<br />
今日越えて<br />
浅き夢見じ<br />
酔いもせず。</p>

<p>Iro wa nioedo / Chirinuru o / Wa ga yo tare zo / Tsune naran / Ui no okuyama / Kyo koete Asaki yume miji / Yoi mo sezu.</p>

<p>とりなくこゑすゆめさませみよあけわたるひんかしをそらいろはえておきつへにほふねむれゐぬもやのうち。</p>

<p>天地星空<br />
山川峰谷<br />
雲霧室苔<br />
人犬上末<br />
硫黄猿生ふ為よ<br />
榎の枝を馴れ居て。</p>

<p>田居に出で菜摘むわれをぞ君召すと求食り追ひゆく山城の打酔へる子ら藻葉干せよえ舟繋けぬ。</p>
</div>
</html>`,
		InlineFiles: []InlineFile{
			{
				ContentType: ContentTypeHeader{
					ContentType: "image/jpeg",
					Params: map[string]string{
						"name": "inline-jpg-image-without-disposition.jpg",
					},
				},
				ContentDisposition: ContentDispositionHeader{
					ContentDisposition: "",
					Params:             map[string]string(nil),
				},
				Data: []byte{255, 216, 255, 219, 0, 67, 0, 3, 2, 2, 2, 2, 2, 3, 2, 2, 2, 3, 3, 3, 3, 4,
					6, 4, 4, 4, 4, 4, 8, 6, 6, 5, 6, 9, 8, 10, 10, 9, 8, 9, 9, 10, 12, 15, 12, 10, 11, 14, 11, 9, 9,
					13, 17, 13, 14, 15, 16, 16, 17, 16, 10, 12, 18, 19, 18, 16, 19, 15, 16, 16, 16, 255, 201, 0, 11, 8,
					0, 1, 0, 1, 1, 1, 17, 0, 255, 204, 0, 6, 0, 16, 16, 5, 255, 218, 0, 8, 1, 1, 0, 0, 63, 0, 210, 207,
					32, 255, 217},
			},
			{
				ContentID: "inline-jpg-image.jpg@example.com",
				ContentType: ContentTypeHeader{
					ContentType: "image/jpeg",
					Params: map[string]string{
						"name": "inline-jpg-image-name.jpg",
					},
				},
				ContentDisposition: ContentDispositionHeader{
					ContentDisposition: inline,
					Params: map[string]string{
						"filename": "inline-jpg-image-filename.jpg",
					},
				},
				Data: []byte{255, 216, 255, 219, 0, 67, 0, 3, 2, 2, 2, 2, 2, 3, 2, 2, 2, 3, 3, 3, 3, 4,
					6, 4, 4, 4, 4, 4, 8, 6, 6, 5, 6, 9, 8, 10, 10, 9, 8, 9, 9, 10, 12, 15, 12, 10, 11, 14, 11, 9, 9,
					13, 17, 13, 14, 15, 16, 16, 17, 16, 10, 12, 18, 19, 18, 16, 19, 15, 16, 16, 16, 255, 201, 0, 11, 8,
					0, 1, 0, 1, 1, 1, 17, 0, 255, 204, 0, 6, 0, 16, 16, 5, 255, 218, 0, 8, 1, 1, 0, 0, 63, 0, 210, 207,
					32, 255, 217},
			},
		},
		AttachedFiles: []AttachedFile{
			{
				ContentType: ContentTypeHeader{
					ContentType: "application/pdf",
					Params: map[string]string{
						"name": "attached-pdf-name.pdf",
					},
				},
				ContentDisposition: ContentDispositionHeader{
					ContentDisposition: attachment,
					Params: map[string]string{
						"filename": "attached-pdf-filename.pdf",
					},
				},
				Data: []byte{37, 80, 68, 70, 45, 49, 46, 13, 116, 114, 97, 105, 108, 101, 114, 60, 60,
					47, 82, 111, 111, 116, 60, 60, 47, 80, 97, 103, 101, 115, 60, 60, 47, 75, 105, 100, 115, 91, 60,
					60, 47, 77, 101, 100, 105, 97, 66, 111, 120, 91, 48, 32, 48, 32, 51, 32, 51, 93, 62, 62, 93, 62,
					62, 62, 62, 62, 62},
			},
			{
				ContentType: ContentTypeHeader{
					ContentType: "application/pdf",
					Params: map[string]string{
						"name": "attached-pdf-without-disposition.pdf",
					},
				},
				ContentDisposition: ContentDispositionHeader{
					ContentDisposition: "",
					Params:             map[string]string(nil),
				},
				Data: []byte{37, 80, 68, 70, 45, 49, 46, 13, 116, 114, 97, 105, 108, 101, 114, 60, 60,
					47, 82, 111, 111, 116, 60, 60, 47, 80, 97, 103, 101, 115, 60, 60, 47, 75, 105, 100, 115, 91, 60,
					60, 47, 77, 101, 100, 105, 97, 66, 111, 120, 91, 48, 32, 48, 32, 51, 32, 51, 93, 62, 62, 93, 62,
					62, 62, 62, 62, 62},
			},
			{
				ContentType: ContentTypeHeader{
					ContentType: "application/json",
					Params: map[string]string{
						"name": "attached-json-name.json",
					},
				},
				ContentDisposition: ContentDispositionHeader{
					ContentDisposition: attachment,
					Params: map[string]string{
						"filename": "attached-json-filename.json",
					},
				},
				Data: []byte{123, 34, 102, 111, 111, 34, 58, 34, 98, 97, 114, 34, 125},
			},
			{
				ContentType: ContentTypeHeader{
					ContentType: "text/plain",
					Params: map[string]string{
						"name": "attached-text-plain-name.txt",
					},
				},
				ContentDisposition: ContentDispositionHeader{
					ContentDisposition: attachment,
					Params: map[string]string{
						"filename": "attached-text-plain-filename.txt",
					},
				},
				Data: []byte{84, 101, 120, 116, 47, 112, 108, 97, 105, 110, 32, 99, 111, 110, 116, 101, 110, 116, 32,
					97, 115, 32, 97, 110, 32, 97, 116, 116, 97, 99, 104, 101, 100, 32, 46, 116, 120, 116, 32, 102, 105,
					108, 101, 46},
			},
			{
				ContentType: ContentTypeHeader{
					ContentType: "text/html",
					Params: map[string]string{
						"name": "attached-text-html-name.html",
					},
				},
				ContentDisposition: ContentDispositionHeader{
					ContentDisposition: attachment,
					Params: map[string]string{
						"filename": "attached-text-html-filename.html",
					},
				},
				Data: []byte{84, 101, 120, 116, 47, 104, 116, 109, 108, 32, 99, 111, 110, 116, 101, 110, 116, 32, 97,
					115, 32, 97, 110, 32, 97, 116, 116, 97, 99, 104, 101, 100, 32, 46, 104, 116, 109, 108, 32, 102,
					105, 108, 101, 46},
			},
		},
	}

	testEmailFromFile(t, fp, expectedEmail)
}
func TestParseEmailJapaneseMultipartMixedIso2022jpOverBase64(t *testing.T) {
	fp := "tests/test_japanese_multipart_mixed_iso-2022-jp_over_base64.txt"
	tz, _ := time.LoadLocation("Asia/Tokyo")
	expectedDate, _ := time.Parse(
		time.RFC1123Z+" (MST)",
		time.Date(2019, time.April, 1, 7, 55, 0, 0, tz).Format(time.RFC1123Z+" (MST)"))
	expectedEmail := Email{
		Headers: Headers{
			Date:    expectedDate,
			Subject: "📧 Test いろは歌",
			ReplyTo: []*mail.Address{
				{
					Name:    "郵便アリス",
					Address: "alice.yuubin@example.net",
				},
			},
			Sender: &mail.Address{
				Name:    "郵便アリス",
				Address: "alice.yuubin@example.com",
			},
			From: []*mail.Address{
				{
					Name:    "郵便アリス",
					Address: "alice.yuubin@example.com",
				},
				{
					Name:    "郵便アリス",
					Address: "alice.yuubin@example.net",
				},
			},
			To: []*mail.Address{
				{
					Name:    "郵便ボブ",
					Address: "bob.yuubin@example.com",
				},
				{
					Name:    "郵便キャロル",
					Address: "carol.yuubin@example.com",
				},
			},
			Cc: []*mail.Address{
				{
					Name:    "郵便ダン",
					Address: "dan.yuubin@example.com",
				},
				{
					Name:    "郵便イーブ",
					Address: "eve.yubin@example.com",
				},
			},
			Bcc: []*mail.Address{
				{
					Name:    "郵便フランク",
					Address: "frank.yuubin@example.com",
				},
				{
					Name:    "郵便グレイス",
					Address: "grace.yubin@example.com",
				},
			},
			MessageID:  "Message-Id-1@example.com",
			InReplyTo:  []MessageId{"Message-Id-0@example.com"},
			References: []MessageId{"Message-Id-0@example.com"},
			Comments:   "Message Header Comment",
			Keywords:   []string{"Keyword 1", "Keyword 2"},
			ResentDate: expectedDate,
			ResentFrom: []*mail.Address{
				{
					Name:    "郵便アリス",
					Address: "alice.yuubin@example.net",
				},
				{
					Name:    "郵便アリス",
					Address: "alice.yuubin@example.com",
				},
			},
			ResentSender: &mail.Address{
				Name:    "郵便アリス",
				Address: "alice.yuubin@example.net",
			},
			ResentTo: []*mail.Address{
				{
					Name:    "郵便ボブ",
					Address: "bob.yuubin@example.net",
				},
				{
					Name:    "郵便キャロル",
					Address: "carol.yuubin@example.net",
				},
			},
			ResentCc: []*mail.Address{
				{
					Name:    "郵便ダン",
					Address: "dan.yuubin@example.net",
				},
				{
					Name:    "郵便イーブ",
					Address: "eve.yubin@example.net",
				},
			},
			ResentBcc: []*mail.Address{
				{
					Name:    "郵便フランク",
					Address: "frank.yuubin@example.net",
				},
				{
					Name:    "郵便グレイス",
					Address: "grace.yubin@example.net",
				},
			},
			ResentMessageID: "Message-Id-1@example.net",
			ContentType: ContentTypeHeader{
				ContentType: "multipart/mixed",
				Params: map[string]string{
					"boundary": "MixedBoundaryString",
					"charset":  "iso-2022-jp",
				},
			},
			ExtraHeaders: map[string][]string{
				"X-Clacks-Overhead": {"GNU Terry Pratchett"},
				"X-Script/function/\t !\"#$%&'()*+,-./;<=>?@[\\]^_`{|}~": {
					"TEST VALUE 1\t !\"#$%&'()*+,-./0123456789:;<=>?@ABCDEFGHIJKLMNOPQRSTUVWXYZ[\\]^_` abcdefghijklmnopqrstuvwxyz{|}~",
					"TEST VALUE 2\t !\"#$%&'()*+,-./0123456789:;<=>?@ABCDEFGHIJKLMNOPQRSTUVWXYZ[\\]^_` abcdefghijklmnopqrstuvwxyz{|}~",
				},
			},
		},
		Text: `色は匂えど
散りぬるを
我が世誰ぞ
常ならん
有為の奥山
今日越えて
浅き夢見じ
酔いもせず。

Iro wa nioedo / Chirinuru o / Wa ga yo tare zo / Tsune naran / Ui no okuyama / Kyo koete Asaki yume miji / Yoi mo sezu.

とりなくこゑすゆめさませみよあけわたるひんかしをそらいろはえておきつへにほふねむれゐぬもやのうち。

天地星空
山川峰谷
雲霧室苔
人犬上末
硫黄猿生ふ為よ
榎の枝を馴れ居て。

田居に出で菜摘むわれをぞ君召すと求食り追ひゆく山城の打酔へる子ら藻葉干せよえ舟繋けぬ。`,
		EnrichedText: `<bold>色は匂えど</bold>
<italic>散りぬるを</italic>
<fixed>我が世誰ぞ</fixed>
<underline>常ならん</underline>
有為の奥山
今日越えて
浅き夢見じ
酔いもせず。

Iro wa nioedo / Chirinuru o / Wa ga yo tare zo / Tsune naran / Ui no okuyama / Kyo koete Asaki yume miji / Yoi mo sezu.

とりなくこゑすゆめさませみよあけわたるひんかしをそらいろはえておきつへにほふねむれゐぬもやのうち。

天地星空
山川峰谷
雲霧室苔
人犬上末
硫黄猿生ふ為よ
榎の枝を馴れ居て。

田居に出で菜摘むわれをぞ君召すと求食り追ひゆく山城の打酔へる子ら藻葉干せよえ舟繋けぬ。`,
		HTML: `<html>
<div dir="ltr">
<p>色は匂えど<br />
散りぬるを<br />
我が世誰ぞ<br />
常ならん<br />
有為の奥山<br />
今日越えて<br />
浅き夢見じ<br />
酔いもせず。</p>

<p>Iro wa nioedo / Chirinuru o / Wa ga yo tare zo / Tsune naran / Ui no okuyama / Kyo koete Asaki yume miji / Yoi mo sezu.</p>

<p>とりなくこゑすゆめさませみよあけわたるひんかしをそらいろはえておきつへにほふねむれゐぬもやのうち。</p>

<p>天地星空<br />
山川峰谷<br />
雲霧室苔<br />
人犬上末<br />
硫黄猿生ふ為よ<br />
榎の枝を馴れ居て。</p>

<p>田居に出で菜摘むわれをぞ君召すと求食り追ひゆく山城の打酔へる子ら藻葉干せよえ舟繋けぬ。</p>
</div>
</html>`,
		InlineFiles: []InlineFile{
			{
				ContentType: ContentTypeHeader{
					ContentType: "image/jpeg",
					Params: map[string]string{
						"name": "inline-jpg-image-without-disposition.jpg",
					},
				},
				ContentDisposition: ContentDispositionHeader{
					ContentDisposition: "",
					Params:             map[string]string(nil),
				},
				Data: []byte{255, 216, 255, 219, 0, 67, 0, 3, 2, 2, 2, 2, 2, 3, 2, 2, 2, 3, 3, 3, 3, 4,
					6, 4, 4, 4, 4, 4, 8, 6, 6, 5, 6, 9, 8, 10, 10, 9, 8, 9, 9, 10, 12, 15, 12, 10, 11, 14, 11, 9, 9,
					13, 17, 13, 14, 15, 16, 16, 17, 16, 10, 12, 18, 19, 18, 16, 19, 15, 16, 16, 16, 255, 201, 0, 11, 8,
					0, 1, 0, 1, 1, 1, 17, 0, 255, 204, 0, 6, 0, 16, 16, 5, 255, 218, 0, 8, 1, 1, 0, 0, 63, 0, 210, 207,
					32, 255, 217},
			},
			{
				ContentID: "inline-jpg-image.jpg@example.com",
				ContentType: ContentTypeHeader{
					ContentType: "image/jpeg",
					Params: map[string]string{
						"name": "inline-jpg-image-name.jpg",
					},
				},
				ContentDisposition: ContentDispositionHeader{
					ContentDisposition: inline,
					Params: map[string]string{
						"filename": "inline-jpg-image-filename.jpg",
					},
				},
				Data: []byte{255, 216, 255, 219, 0, 67, 0, 3, 2, 2, 2, 2, 2, 3, 2, 2, 2, 3, 3, 3, 3, 4,
					6, 4, 4, 4, 4, 4, 8, 6, 6, 5, 6, 9, 8, 10, 10, 9, 8, 9, 9, 10, 12, 15, 12, 10, 11, 14, 11, 9, 9,
					13, 17, 13, 14, 15, 16, 16, 17, 16, 10, 12, 18, 19, 18, 16, 19, 15, 16, 16, 16, 255, 201, 0, 11, 8,
					0, 1, 0, 1, 1, 1, 17, 0, 255, 204, 0, 6, 0, 16, 16, 5, 255, 218, 0, 8, 1, 1, 0, 0, 63, 0, 210, 207,
					32, 255, 217},
			},
		},
		AttachedFiles: []AttachedFile{
			{
				ContentType: ContentTypeHeader{
					ContentType: "application/pdf",
					Params: map[string]string{
						"name": "attached-pdf-name.pdf",
					},
				},
				ContentDisposition: ContentDispositionHeader{
					ContentDisposition: attachment,
					Params: map[string]string{
						"filename": "attached-pdf-filename.pdf",
					},
				},
				Data: []byte{37, 80, 68, 70, 45, 49, 46, 13, 116, 114, 97, 105, 108, 101, 114, 60, 60,
					47, 82, 111, 111, 116, 60, 60, 47, 80, 97, 103, 101, 115, 60, 60, 47, 75, 105, 100, 115, 91, 60,
					60, 47, 77, 101, 100, 105, 97, 66, 111, 120, 91, 48, 32, 48, 32, 51, 32, 51, 93, 62, 62, 93, 62,
					62, 62, 62, 62, 62},
			},
			{
				ContentType: ContentTypeHeader{
					ContentType: "application/pdf",
					Params: map[string]string{
						"name": "attached-pdf-without-disposition.pdf",
					},
				},
				ContentDisposition: ContentDispositionHeader{
					ContentDisposition: "",
					Params:             map[string]string(nil),
				},
				Data: []byte{37, 80, 68, 70, 45, 49, 46, 13, 116, 114, 97, 105, 108, 101, 114, 60, 60,
					47, 82, 111, 111, 116, 60, 60, 47, 80, 97, 103, 101, 115, 60, 60, 47, 75, 105, 100, 115, 91, 60,
					60, 47, 77, 101, 100, 105, 97, 66, 111, 120, 91, 48, 32, 48, 32, 51, 32, 51, 93, 62, 62, 93, 62,
					62, 62, 62, 62, 62},
			},
			{
				ContentType: ContentTypeHeader{
					ContentType: "application/json",
					Params: map[string]string{
						"name": "attached-json-name.json",
					},
				},
				ContentDisposition: ContentDispositionHeader{
					ContentDisposition: attachment,
					Params: map[string]string{
						"filename": "attached-json-filename.json",
					},
				},
				Data: []byte{123, 34, 102, 111, 111, 34, 58, 34, 98, 97, 114, 34, 125},
			},
			{
				ContentType: ContentTypeHeader{
					ContentType: "text/plain",
					Params: map[string]string{
						"name": "attached-text-plain-name.txt",
					},
				},
				ContentDisposition: ContentDispositionHeader{
					ContentDisposition: attachment,
					Params: map[string]string{
						"filename": "attached-text-plain-filename.txt",
					},
				},
				Data: []byte{84, 101, 120, 116, 47, 112, 108, 97, 105, 110, 32, 99, 111, 110, 116, 101, 110, 116, 32,
					97, 115, 32, 97, 110, 32, 97, 116, 116, 97, 99, 104, 101, 100, 32, 46, 116, 120, 116, 32, 102, 105,
					108, 101, 46},
			},
			{
				ContentType: ContentTypeHeader{
					ContentType: "text/html",
					Params: map[string]string{
						"name": "attached-text-html-name.html",
					},
				},
				ContentDisposition: ContentDispositionHeader{
					ContentDisposition: attachment,
					Params: map[string]string{
						"filename": "attached-text-html-filename.html",
					},
				},
				Data: []byte{84, 101, 120, 116, 47, 104, 116, 109, 108, 32, 99, 111, 110, 116, 101, 110, 116, 32, 97,
					115, 32, 97, 110, 32, 97, 116, 116, 97, 99, 104, 101, 100, 32, 46, 104, 116, 109, 108, 32, 102,
					105, 108, 101, 46},
			},
		},
	}

	testEmailFromFile(t, fp, expectedEmail)
}
func TestParseEmailJapaneseMultipartMixedIso2022jpOverQuotedprintable(t *testing.T) {
	fp := "tests/test_japanese_multipart_mixed_iso-2022-jp_over_quoted-printable.txt"
	tz, _ := time.LoadLocation("Asia/Tokyo")
	expectedDate, _ := time.Parse(
		time.RFC1123Z+" (MST)",
		time.Date(2019, time.April, 1, 7, 55, 0, 0, tz).Format(time.RFC1123Z+" (MST)"))
	expectedEmail := Email{
		Headers: Headers{
			Date:    expectedDate,
			Subject: "📧 Test いろは歌",
			ReplyTo: []*mail.Address{
				{
					Name:    "郵便アリス",
					Address: "alice.yuubin@example.net",
				},
			},
			Sender: &mail.Address{
				Name:    "郵便アリス",
				Address: "alice.yuubin@example.com",
			},
			From: []*mail.Address{
				{
					Name:    "郵便アリス",
					Address: "alice.yuubin@example.com",
				},
				{
					Name:    "郵便アリス",
					Address: "alice.yuubin@example.net",
				},
			},
			To: []*mail.Address{
				{
					Name:    "郵便ボブ",
					Address: "bob.yuubin@example.com",
				},
				{
					Name:    "郵便キャロル",
					Address: "carol.yuubin@example.com",
				},
			},
			Cc: []*mail.Address{
				{
					Name:    "郵便ダン",
					Address: "dan.yuubin@example.com",
				},
				{
					Name:    "郵便イーブ",
					Address: "eve.yubin@example.com",
				},
			},
			Bcc: []*mail.Address{
				{
					Name:    "郵便フランク",
					Address: "frank.yuubin@example.com",
				},
				{
					Name:    "郵便グレイス",
					Address: "grace.yubin@example.com",
				},
			},
			MessageID:  "Message-Id-1@example.com",
			InReplyTo:  []MessageId{"Message-Id-0@example.com"},
			References: []MessageId{"Message-Id-0@example.com"},
			Comments:   "Message Header Comment",
			Keywords:   []string{"Keyword 1", "Keyword 2"},
			ResentDate: expectedDate,
			ResentFrom: []*mail.Address{
				{
					Name:    "郵便アリス",
					Address: "alice.yuubin@example.net",
				},
				{
					Name:    "郵便アリス",
					Address: "alice.yuubin@example.com",
				},
			},
			ResentSender: &mail.Address{
				Name:    "郵便アリス",
				Address: "alice.yuubin@example.net",
			},
			ResentTo: []*mail.Address{
				{
					Name:    "郵便ボブ",
					Address: "bob.yuubin@example.net",
				},
				{
					Name:    "郵便キャロル",
					Address: "carol.yuubin@example.net",
				},
			},
			ResentCc: []*mail.Address{
				{
					Name:    "郵便ダン",
					Address: "dan.yuubin@example.net",
				},
				{
					Name:    "郵便イーブ",
					Address: "eve.yubin@example.net",
				},
			},
			ResentBcc: []*mail.Address{
				{
					Name:    "郵便フランク",
					Address: "frank.yuubin@example.net",
				},
				{
					Name:    "郵便グレイス",
					Address: "grace.yubin@example.net",
				},
			},
			ResentMessageID: "Message-Id-1@example.net",
			ContentType: ContentTypeHeader{
				ContentType: "multipart/mixed",
				Params: map[string]string{
					"boundary": "MixedBoundaryString",
					"charset":  "iso-2022-jp",
				},
			},
			ExtraHeaders: map[string][]string{
				"X-Clacks-Overhead": {"GNU Terry Pratchett"},
				"X-Script/function/\t !\"#$%&'()*+,-./;<=>?@[\\]^_`{|}~": {
					"TEST VALUE 1\t !\"#$%&'()*+,-./0123456789:;<=>?@ABCDEFGHIJKLMNOPQRSTUVWXYZ[\\]^_` abcdefghijklmnopqrstuvwxyz{|}~",
					"TEST VALUE 2\t !\"#$%&'()*+,-./0123456789:;<=>?@ABCDEFGHIJKLMNOPQRSTUVWXYZ[\\]^_` abcdefghijklmnopqrstuvwxyz{|}~",
				},
			},
		},
		Text: `色は匂えど
散りぬるを
我が世誰ぞ
常ならん
有為の奥山
今日越えて
浅き夢見じ
酔いもせず。

Iro wa nioedo / Chirinuru o / Wa ga yo tare zo / Tsune naran / Ui no okuyama / Kyo koete Asaki yume miji / Yoi mo sezu.

とりなくこゑすゆめさませみよあけわたるひんかしをそらいろはえておきつへにほふねむれゐぬもやのうち。

天地星空
山川峰谷
雲霧室苔
人犬上末
硫黄猿生ふ為よ
榎の枝を馴れ居て。

田居に出で菜摘むわれをぞ君召すと求食り追ひゆく山城の打酔へる子ら藻葉干せよえ舟繋けぬ。`,
		EnrichedText: `<bold>色は匂えど</bold>
<italic>散りぬるを</italic>
<fixed>我が世誰ぞ</fixed>
<underline>常ならん</underline>
有為の奥山
今日越えて
浅き夢見じ
酔いもせず。

Iro wa nioedo / Chirinuru o / Wa ga yo tare zo / Tsune naran / Ui no okuyama / Kyo koete Asaki yume miji / Yoi mo sezu.

とりなくこゑすゆめさませみよあけわたるひんかしをそらいろはえておきつへにほふねむれゐぬもやのうち。

天地星空
山川峰谷
雲霧室苔
人犬上末
硫黄猿生ふ為よ
榎の枝を馴れ居て。

田居に出で菜摘むわれをぞ君召すと求食り追ひゆく山城の打酔へる子ら藻葉干せよえ舟繋けぬ。`,
		HTML: `<html>
<div dir="ltr">
<p>色は匂えど<br />
散りぬるを<br />
我が世誰ぞ<br />
常ならん<br />
有為の奥山<br />
今日越えて<br />
浅き夢見じ<br />
酔いもせず。</p>

<p>Iro wa nioedo / Chirinuru o / Wa ga yo tare zo / Tsune naran / Ui no okuyama / Kyo koete Asaki yume miji / Yoi mo sezu.</p>

<p>とりなくこゑすゆめさませみよあけわたるひんかしをそらいろはえておきつへにほふねむれゐぬもやのうち。</p>

<p>天地星空<br />
山川峰谷<br />
雲霧室苔<br />
人犬上末<br />
硫黄猿生ふ為よ<br />
榎の枝を馴れ居て。</p>

<p>田居に出で菜摘むわれをぞ君召すと求食り追ひゆく山城の打酔へる子ら藻葉干せよえ舟繋けぬ。</p>
</div>
</html>`,
		InlineFiles: []InlineFile{
			{
				ContentType: ContentTypeHeader{
					ContentType: "image/jpeg",
					Params: map[string]string{
						"name": "inline-jpg-image-without-disposition.jpg",
					},
				},
				ContentDisposition: ContentDispositionHeader{
					ContentDisposition: "",
					Params:             map[string]string(nil),
				},
				Data: []byte{255, 216, 255, 219, 0, 67, 0, 3, 2, 2, 2, 2, 2, 3, 2, 2, 2, 3, 3, 3, 3, 4,
					6, 4, 4, 4, 4, 4, 8, 6, 6, 5, 6, 9, 8, 10, 10, 9, 8, 9, 9, 10, 12, 15, 12, 10, 11, 14, 11, 9, 9,
					13, 17, 13, 14, 15, 16, 16, 17, 16, 10, 12, 18, 19, 18, 16, 19, 15, 16, 16, 16, 255, 201, 0, 11, 8,
					0, 1, 0, 1, 1, 1, 17, 0, 255, 204, 0, 6, 0, 16, 16, 5, 255, 218, 0, 8, 1, 1, 0, 0, 63, 0, 210, 207,
					32, 255, 217},
			},
			{
				ContentID: "inline-jpg-image.jpg@example.com",
				ContentType: ContentTypeHeader{
					ContentType: "image/jpeg",
					Params: map[string]string{
						"name": "inline-jpg-image-name.jpg",
					},
				},
				ContentDisposition: ContentDispositionHeader{
					ContentDisposition: inline,
					Params: map[string]string{
						"filename": "inline-jpg-image-filename.jpg",
					},
				},
				Data: []byte{255, 216, 255, 219, 0, 67, 0, 3, 2, 2, 2, 2, 2, 3, 2, 2, 2, 3, 3, 3, 3, 4,
					6, 4, 4, 4, 4, 4, 8, 6, 6, 5, 6, 9, 8, 10, 10, 9, 8, 9, 9, 10, 12, 15, 12, 10, 11, 14, 11, 9, 9,
					13, 17, 13, 14, 15, 16, 16, 17, 16, 10, 12, 18, 19, 18, 16, 19, 15, 16, 16, 16, 255, 201, 0, 11, 8,
					0, 1, 0, 1, 1, 1, 17, 0, 255, 204, 0, 6, 0, 16, 16, 5, 255, 218, 0, 8, 1, 1, 0, 0, 63, 0, 210, 207,
					32, 255, 217},
			},
		},
		AttachedFiles: []AttachedFile{
			{
				ContentType: ContentTypeHeader{
					ContentType: "application/pdf",
					Params: map[string]string{
						"name": "attached-pdf-name.pdf",
					},
				},
				ContentDisposition: ContentDispositionHeader{
					ContentDisposition: attachment,
					Params: map[string]string{
						"filename": "attached-pdf-filename.pdf",
					},
				},
				Data: []byte{37, 80, 68, 70, 45, 49, 46, 13, 116, 114, 97, 105, 108, 101, 114, 60, 60,
					47, 82, 111, 111, 116, 60, 60, 47, 80, 97, 103, 101, 115, 60, 60, 47, 75, 105, 100, 115, 91, 60,
					60, 47, 77, 101, 100, 105, 97, 66, 111, 120, 91, 48, 32, 48, 32, 51, 32, 51, 93, 62, 62, 93, 62,
					62, 62, 62, 62, 62},
			},
			{
				ContentType: ContentTypeHeader{
					ContentType: "application/pdf",
					Params: map[string]string{
						"name": "attached-pdf-without-disposition.pdf",
					},
				},
				ContentDisposition: ContentDispositionHeader{
					ContentDisposition: "",
					Params:             map[string]string(nil),
				},
				Data: []byte{37, 80, 68, 70, 45, 49, 46, 13, 116, 114, 97, 105, 108, 101, 114, 60, 60,
					47, 82, 111, 111, 116, 60, 60, 47, 80, 97, 103, 101, 115, 60, 60, 47, 75, 105, 100, 115, 91, 60,
					60, 47, 77, 101, 100, 105, 97, 66, 111, 120, 91, 48, 32, 48, 32, 51, 32, 51, 93, 62, 62, 93, 62,
					62, 62, 62, 62, 62},
			},
			{
				ContentType: ContentTypeHeader{
					ContentType: "application/json",
					Params: map[string]string{
						"name": "attached-json-name.json",
					},
				},
				ContentDisposition: ContentDispositionHeader{
					ContentDisposition: attachment,
					Params: map[string]string{
						"filename": "attached-json-filename.json",
					},
				},
				Data: []byte{123, 34, 102, 111, 111, 34, 58, 34, 98, 97, 114, 34, 125},
			},
			{
				ContentType: ContentTypeHeader{
					ContentType: "text/plain",
					Params: map[string]string{
						"name": "attached-text-plain-name.txt",
					},
				},
				ContentDisposition: ContentDispositionHeader{
					ContentDisposition: attachment,
					Params: map[string]string{
						"filename": "attached-text-plain-filename.txt",
					},
				},
				Data: []byte{84, 101, 120, 116, 47, 112, 108, 97, 105, 110, 32, 99, 111, 110, 116, 101, 110, 116, 32,
					97, 115, 32, 97, 110, 32, 97, 116, 116, 97, 99, 104, 101, 100, 32, 46, 116, 120, 116, 32, 102, 105,
					108, 101, 46},
			},
			{
				ContentType: ContentTypeHeader{
					ContentType: "text/html",
					Params: map[string]string{
						"name": "attached-text-html-name.html",
					},
				},
				ContentDisposition: ContentDispositionHeader{
					ContentDisposition: attachment,
					Params: map[string]string{
						"filename": "attached-text-html-filename.html",
					},
				},
				Data: []byte{84, 101, 120, 116, 47, 104, 116, 109, 108, 32, 99, 111, 110, 116, 101, 110, 116, 32, 97,
					115, 32, 97, 110, 32, 97, 116, 116, 97, 99, 104, 101, 100, 32, 46, 104, 116, 109, 108, 32, 102,
					105, 108, 101, 46},
			},
		},
	}

	testEmailFromFile(t, fp, expectedEmail)
}
func TestParseEmailJapaneseMultipartMixedEucjpOverBase64(t *testing.T) {
	fp := "tests/test_japanese_multipart_mixed_euc-jp_over_base64.txt"
	tz, _ := time.LoadLocation("Asia/Tokyo")
	expectedDate, _ := time.Parse(
		time.RFC1123Z+" (MST)",
		time.Date(2019, time.April, 1, 7, 55, 0, 0, tz).Format(time.RFC1123Z+" (MST)"))
	expectedEmail := Email{
		Headers: Headers{
			Date:    expectedDate,
			Subject: "📧 Test いろは歌",
			ReplyTo: []*mail.Address{
				{
					Name:    "郵便アリス",
					Address: "alice.yuubin@example.net",
				},
			},
			Sender: &mail.Address{
				Name:    "郵便アリス",
				Address: "alice.yuubin@example.com",
			},
			From: []*mail.Address{
				{
					Name:    "郵便アリス",
					Address: "alice.yuubin@example.com",
				},
				{
					Name:    "郵便アリス",
					Address: "alice.yuubin@example.net",
				},
			},
			To: []*mail.Address{
				{
					Name:    "郵便ボブ",
					Address: "bob.yuubin@example.com",
				},
				{
					Name:    "郵便キャロル",
					Address: "carol.yuubin@example.com",
				},
			},
			Cc: []*mail.Address{
				{
					Name:    "郵便ダン",
					Address: "dan.yuubin@example.com",
				},
				{
					Name:    "郵便イーブ",
					Address: "eve.yubin@example.com",
				},
			},
			Bcc: []*mail.Address{
				{
					Name:    "郵便フランク",
					Address: "frank.yuubin@example.com",
				},
				{
					Name:    "郵便グレイス",
					Address: "grace.yubin@example.com",
				},
			},
			MessageID:  "Message-Id-1@example.com",
			InReplyTo:  []MessageId{"Message-Id-0@example.com"},
			References: []MessageId{"Message-Id-0@example.com"},
			Comments:   "Message Header Comment",
			Keywords:   []string{"Keyword 1", "Keyword 2"},
			ResentDate: expectedDate,
			ResentFrom: []*mail.Address{
				{
					Name:    "郵便アリス",
					Address: "alice.yuubin@example.net",
				},
				{
					Name:    "郵便アリス",
					Address: "alice.yuubin@example.com",
				},
			},
			ResentSender: &mail.Address{
				Name:    "郵便アリス",
				Address: "alice.yuubin@example.net",
			},
			ResentTo: []*mail.Address{
				{
					Name:    "郵便ボブ",
					Address: "bob.yuubin@example.net",
				},
				{
					Name:    "郵便キャロル",
					Address: "carol.yuubin@example.net",
				},
			},
			ResentCc: []*mail.Address{
				{
					Name:    "郵便ダン",
					Address: "dan.yuubin@example.net",
				},
				{
					Name:    "郵便イーブ",
					Address: "eve.yubin@example.net",
				},
			},
			ResentBcc: []*mail.Address{
				{
					Name:    "郵便フランク",
					Address: "frank.yuubin@example.net",
				},
				{
					Name:    "郵便グレイス",
					Address: "grace.yubin@example.net",
				},
			},
			ResentMessageID: "Message-Id-1@example.net",
			ContentType: ContentTypeHeader{
				ContentType: "multipart/mixed",
				Params: map[string]string{
					"boundary": "MixedBoundaryString",
					"charset":  "euc-jp",
				},
			},
			ExtraHeaders: map[string][]string{
				"X-Clacks-Overhead": {"GNU Terry Pratchett"},
				"X-Script/function/\t !\"#$%&'()*+,-./;<=>?@[\\]^_`{|}~": {
					"TEST VALUE 1\t !\"#$%&'()*+,-./0123456789:;<=>?@ABCDEFGHIJKLMNOPQRSTUVWXYZ[\\]^_` abcdefghijklmnopqrstuvwxyz{|}~",
					"TEST VALUE 2\t !\"#$%&'()*+,-./0123456789:;<=>?@ABCDEFGHIJKLMNOPQRSTUVWXYZ[\\]^_` abcdefghijklmnopqrstuvwxyz{|}~",
				},
			},
		},
		Text: `色は匂えど
散りぬるを
我が世誰ぞ
常ならん
有為の奥山
今日越えて
浅き夢見じ
酔いもせず。

Iro wa nioedo / Chirinuru o / Wa ga yo tare zo / Tsune naran / Ui no okuyama / Kyo koete Asaki yume miji / Yoi mo sezu.

とりなくこゑすゆめさませみよあけわたるひんかしをそらいろはえておきつへにほふねむれゐぬもやのうち。

天地星空
山川峰谷
雲霧室苔
人犬上末
硫黄猿生ふ為よ
榎の枝を馴れ居て。

田居に出で菜摘むわれをぞ君召すと求食り追ひゆく山城の打酔へる子ら藻葉干せよえ舟繋けぬ。`,
		EnrichedText: `<bold>色は匂えど</bold>
<italic>散りぬるを</italic>
<fixed>我が世誰ぞ</fixed>
<underline>常ならん</underline>
有為の奥山
今日越えて
浅き夢見じ
酔いもせず。

Iro wa nioedo / Chirinuru o / Wa ga yo tare zo / Tsune naran / Ui no okuyama / Kyo koete Asaki yume miji / Yoi mo sezu.

とりなくこゑすゆめさませみよあけわたるひんかしをそらいろはえておきつへにほふねむれゐぬもやのうち。

天地星空
山川峰谷
雲霧室苔
人犬上末
硫黄猿生ふ為よ
榎の枝を馴れ居て。

田居に出で菜摘むわれをぞ君召すと求食り追ひゆく山城の打酔へる子ら藻葉干せよえ舟繋けぬ。`,
		HTML: `<html>
<div dir="ltr">
<p>色は匂えど<br />
散りぬるを<br />
我が世誰ぞ<br />
常ならん<br />
有為の奥山<br />
今日越えて<br />
浅き夢見じ<br />
酔いもせず。</p>

<p>Iro wa nioedo / Chirinuru o / Wa ga yo tare zo / Tsune naran / Ui no okuyama / Kyo koete Asaki yume miji / Yoi mo sezu.</p>

<p>とりなくこゑすゆめさませみよあけわたるひんかしをそらいろはえておきつへにほふねむれゐぬもやのうち。</p>

<p>天地星空<br />
山川峰谷<br />
雲霧室苔<br />
人犬上末<br />
硫黄猿生ふ為よ<br />
榎の枝を馴れ居て。</p>

<p>田居に出で菜摘むわれをぞ君召すと求食り追ひゆく山城の打酔へる子ら藻葉干せよえ舟繋けぬ。</p>
</div>
</html>`,
		InlineFiles: []InlineFile{
			{
				ContentType: ContentTypeHeader{
					ContentType: "image/jpeg",
					Params: map[string]string{
						"name": "inline-jpg-image-without-disposition.jpg",
					},
				},
				ContentDisposition: ContentDispositionHeader{
					ContentDisposition: "",
					Params:             map[string]string(nil),
				},
				Data: []byte{255, 216, 255, 219, 0, 67, 0, 3, 2, 2, 2, 2, 2, 3, 2, 2, 2, 3, 3, 3, 3, 4,
					6, 4, 4, 4, 4, 4, 8, 6, 6, 5, 6, 9, 8, 10, 10, 9, 8, 9, 9, 10, 12, 15, 12, 10, 11, 14, 11, 9, 9,
					13, 17, 13, 14, 15, 16, 16, 17, 16, 10, 12, 18, 19, 18, 16, 19, 15, 16, 16, 16, 255, 201, 0, 11, 8,
					0, 1, 0, 1, 1, 1, 17, 0, 255, 204, 0, 6, 0, 16, 16, 5, 255, 218, 0, 8, 1, 1, 0, 0, 63, 0, 210, 207,
					32, 255, 217},
			},
			{
				ContentID: "inline-jpg-image.jpg@example.com",
				ContentType: ContentTypeHeader{
					ContentType: "image/jpeg",
					Params: map[string]string{
						"name": "inline-jpg-image-name.jpg",
					},
				},
				ContentDisposition: ContentDispositionHeader{
					ContentDisposition: inline,
					Params: map[string]string{
						"filename": "inline-jpg-image-filename.jpg",
					},
				},
				Data: []byte{255, 216, 255, 219, 0, 67, 0, 3, 2, 2, 2, 2, 2, 3, 2, 2, 2, 3, 3, 3, 3, 4,
					6, 4, 4, 4, 4, 4, 8, 6, 6, 5, 6, 9, 8, 10, 10, 9, 8, 9, 9, 10, 12, 15, 12, 10, 11, 14, 11, 9, 9,
					13, 17, 13, 14, 15, 16, 16, 17, 16, 10, 12, 18, 19, 18, 16, 19, 15, 16, 16, 16, 255, 201, 0, 11, 8,
					0, 1, 0, 1, 1, 1, 17, 0, 255, 204, 0, 6, 0, 16, 16, 5, 255, 218, 0, 8, 1, 1, 0, 0, 63, 0, 210, 207,
					32, 255, 217},
			},
		},
		AttachedFiles: []AttachedFile{
			{
				ContentType: ContentTypeHeader{
					ContentType: "application/pdf",
					Params: map[string]string{
						"name": "attached-pdf-name.pdf",
					},
				},
				ContentDisposition: ContentDispositionHeader{
					ContentDisposition: attachment,
					Params: map[string]string{
						"filename": "attached-pdf-filename.pdf",
					},
				},
				Data: []byte{37, 80, 68, 70, 45, 49, 46, 13, 116, 114, 97, 105, 108, 101, 114, 60, 60,
					47, 82, 111, 111, 116, 60, 60, 47, 80, 97, 103, 101, 115, 60, 60, 47, 75, 105, 100, 115, 91, 60,
					60, 47, 77, 101, 100, 105, 97, 66, 111, 120, 91, 48, 32, 48, 32, 51, 32, 51, 93, 62, 62, 93, 62,
					62, 62, 62, 62, 62},
			},
			{
				ContentType: ContentTypeHeader{
					ContentType: "application/pdf",
					Params: map[string]string{
						"name": "attached-pdf-without-disposition.pdf",
					},
				},
				ContentDisposition: ContentDispositionHeader{
					ContentDisposition: "",
					Params:             map[string]string(nil),
				},
				Data: []byte{37, 80, 68, 70, 45, 49, 46, 13, 116, 114, 97, 105, 108, 101, 114, 60, 60,
					47, 82, 111, 111, 116, 60, 60, 47, 80, 97, 103, 101, 115, 60, 60, 47, 75, 105, 100, 115, 91, 60,
					60, 47, 77, 101, 100, 105, 97, 66, 111, 120, 91, 48, 32, 48, 32, 51, 32, 51, 93, 62, 62, 93, 62,
					62, 62, 62, 62, 62},
			},
			{
				ContentType: ContentTypeHeader{
					ContentType: "application/json",
					Params: map[string]string{
						"name": "attached-json-name.json",
					},
				},
				ContentDisposition: ContentDispositionHeader{
					ContentDisposition: attachment,
					Params: map[string]string{
						"filename": "attached-json-filename.json",
					},
				},
				Data: []byte{123, 34, 102, 111, 111, 34, 58, 34, 98, 97, 114, 34, 125},
			},
			{
				ContentType: ContentTypeHeader{
					ContentType: "text/plain",
					Params: map[string]string{
						"name": "attached-text-plain-name.txt",
					},
				},
				ContentDisposition: ContentDispositionHeader{
					ContentDisposition: attachment,
					Params: map[string]string{
						"filename": "attached-text-plain-filename.txt",
					},
				},
				Data: []byte{84, 101, 120, 116, 47, 112, 108, 97, 105, 110, 32, 99, 111, 110, 116, 101, 110, 116, 32,
					97, 115, 32, 97, 110, 32, 97, 116, 116, 97, 99, 104, 101, 100, 32, 46, 116, 120, 116, 32, 102, 105,
					108, 101, 46},
			},
			{
				ContentType: ContentTypeHeader{
					ContentType: "text/html",
					Params: map[string]string{
						"name": "attached-text-html-name.html",
					},
				},
				ContentDisposition: ContentDispositionHeader{
					ContentDisposition: attachment,
					Params: map[string]string{
						"filename": "attached-text-html-filename.html",
					},
				},
				Data: []byte{84, 101, 120, 116, 47, 104, 116, 109, 108, 32, 99, 111, 110, 116, 101, 110, 116, 32, 97,
					115, 32, 97, 110, 32, 97, 116, 116, 97, 99, 104, 101, 100, 32, 46, 104, 116, 109, 108, 32, 102,
					105, 108, 101, 46},
			},
		},
	}

	testEmailFromFile(t, fp, expectedEmail)
}
func TestParseEmailJapaneseMultipartMixedEucjpOverQuotedprintable(t *testing.T) {
	fp := "tests/test_japanese_multipart_mixed_euc-jp_over_quoted-printable.txt"
	tz, _ := time.LoadLocation("Asia/Tokyo")
	expectedDate, _ := time.Parse(
		time.RFC1123Z+" (MST)",
		time.Date(2019, time.April, 1, 7, 55, 0, 0, tz).Format(time.RFC1123Z+" (MST)"))
	expectedEmail := Email{
		Headers: Headers{
			Date:    expectedDate,
			Subject: "📧 Test いろは歌",
			ReplyTo: []*mail.Address{
				{
					Name:    "郵便アリス",
					Address: "alice.yuubin@example.net",
				},
			},
			Sender: &mail.Address{
				Name:    "郵便アリス",
				Address: "alice.yuubin@example.com",
			},
			From: []*mail.Address{
				{
					Name:    "郵便アリス",
					Address: "alice.yuubin@example.com",
				},
				{
					Name:    "郵便アリス",
					Address: "alice.yuubin@example.net",
				},
			},
			To: []*mail.Address{
				{
					Name:    "郵便ボブ",
					Address: "bob.yuubin@example.com",
				},
				{
					Name:    "郵便キャロル",
					Address: "carol.yuubin@example.com",
				},
			},
			Cc: []*mail.Address{
				{
					Name:    "郵便ダン",
					Address: "dan.yuubin@example.com",
				},
				{
					Name:    "郵便イーブ",
					Address: "eve.yubin@example.com",
				},
			},
			Bcc: []*mail.Address{
				{
					Name:    "郵便フランク",
					Address: "frank.yuubin@example.com",
				},
				{
					Name:    "郵便グレイス",
					Address: "grace.yubin@example.com",
				},
			},
			MessageID:  "Message-Id-1@example.com",
			InReplyTo:  []MessageId{"Message-Id-0@example.com"},
			References: []MessageId{"Message-Id-0@example.com"},
			Comments:   "Message Header Comment",
			Keywords:   []string{"Keyword 1", "Keyword 2"},
			ResentDate: expectedDate,
			ResentFrom: []*mail.Address{
				{
					Name:    "郵便アリス",
					Address: "alice.yuubin@example.net",
				},
				{
					Name:    "郵便アリス",
					Address: "alice.yuubin@example.com",
				},
			},
			ResentSender: &mail.Address{
				Name:    "郵便アリス",
				Address: "alice.yuubin@example.net",
			},
			ResentTo: []*mail.Address{
				{
					Name:    "郵便ボブ",
					Address: "bob.yuubin@example.net",
				},
				{
					Name:    "郵便キャロル",
					Address: "carol.yuubin@example.net",
				},
			},
			ResentCc: []*mail.Address{
				{
					Name:    "郵便ダン",
					Address: "dan.yuubin@example.net",
				},
				{
					Name:    "郵便イーブ",
					Address: "eve.yubin@example.net",
				},
			},
			ResentBcc: []*mail.Address{
				{
					Name:    "郵便フランク",
					Address: "frank.yuubin@example.net",
				},
				{
					Name:    "郵便グレイス",
					Address: "grace.yubin@example.net",
				},
			},
			ResentMessageID: "Message-Id-1@example.net",
			ContentType: ContentTypeHeader{
				ContentType: "multipart/mixed",
				Params: map[string]string{
					"boundary": "MixedBoundaryString",
					"charset":  "euc-jp",
				},
			},
			ExtraHeaders: map[string][]string{
				"X-Clacks-Overhead": {"GNU Terry Pratchett"},
				"X-Script/function/\t !\"#$%&'()*+,-./;<=>?@[\\]^_`{|}~": {
					"TEST VALUE 1\t !\"#$%&'()*+,-./0123456789:;<=>?@ABCDEFGHIJKLMNOPQRSTUVWXYZ[\\]^_` abcdefghijklmnopqrstuvwxyz{|}~",
					"TEST VALUE 2\t !\"#$%&'()*+,-./0123456789:;<=>?@ABCDEFGHIJKLMNOPQRSTUVWXYZ[\\]^_` abcdefghijklmnopqrstuvwxyz{|}~",
				},
			},
		},
		Text: `色は匂えど
散りぬるを
我が世誰ぞ
常ならん
有為の奥山
今日越えて
浅き夢見じ
酔いもせず。

Iro wa nioedo / Chirinuru o / Wa ga yo tare zo / Tsune naran / Ui no okuyama / Kyo koete Asaki yume miji / Yoi mo sezu.

とりなくこゑすゆめさませみよあけわたるひんかしをそらいろはえておきつへにほふねむれゐぬもやのうち。

天地星空
山川峰谷
雲霧室苔
人犬上末
硫黄猿生ふ為よ
榎の枝を馴れ居て。

田居に出で菜摘むわれをぞ君召すと求食り追ひゆく山城の打酔へる子ら藻葉干せよえ舟繋けぬ。`,
		EnrichedText: `<bold>色は匂えど</bold>
<italic>散りぬるを</italic>
<fixed>我が世誰ぞ</fixed>
<underline>常ならん</underline>
有為の奥山
今日越えて
浅き夢見じ
酔いもせず。

Iro wa nioedo / Chirinuru o / Wa ga yo tare zo / Tsune naran / Ui no okuyama / Kyo koete Asaki yume miji / Yoi mo sezu.

とりなくこゑすゆめさませみよあけわたるひんかしをそらいろはえておきつへにほふねむれゐぬもやのうち。

天地星空
山川峰谷
雲霧室苔
人犬上末
硫黄猿生ふ為よ
榎の枝を馴れ居て。

田居に出で菜摘むわれをぞ君召すと求食り追ひゆく山城の打酔へる子ら藻葉干せよえ舟繋けぬ。`,
		HTML: `<html>
<div dir="ltr">
<p>色は匂えど<br />
散りぬるを<br />
我が世誰ぞ<br />
常ならん<br />
有為の奥山<br />
今日越えて<br />
浅き夢見じ<br />
酔いもせず。</p>

<p>Iro wa nioedo / Chirinuru o / Wa ga yo tare zo / Tsune naran / Ui no okuyama / Kyo koete Asaki yume miji / Yoi mo sezu.</p>

<p>とりなくこゑすゆめさませみよあけわたるひんかしをそらいろはえておきつへにほふねむれゐぬもやのうち。</p>

<p>天地星空<br />
山川峰谷<br />
雲霧室苔<br />
人犬上末<br />
硫黄猿生ふ為よ<br />
榎の枝を馴れ居て。</p>

<p>田居に出で菜摘むわれをぞ君召すと求食り追ひゆく山城の打酔へる子ら藻葉干せよえ舟繋けぬ。</p>
</div>
</html>`,
		InlineFiles: []InlineFile{
			{
				ContentType: ContentTypeHeader{
					ContentType: "image/jpeg",
					Params: map[string]string{
						"name": "inline-jpg-image-without-disposition.jpg",
					},
				},
				ContentDisposition: ContentDispositionHeader{
					ContentDisposition: "",
					Params:             map[string]string(nil),
				},
				Data: []byte{255, 216, 255, 219, 0, 67, 0, 3, 2, 2, 2, 2, 2, 3, 2, 2, 2, 3, 3, 3, 3, 4,
					6, 4, 4, 4, 4, 4, 8, 6, 6, 5, 6, 9, 8, 10, 10, 9, 8, 9, 9, 10, 12, 15, 12, 10, 11, 14, 11, 9, 9,
					13, 17, 13, 14, 15, 16, 16, 17, 16, 10, 12, 18, 19, 18, 16, 19, 15, 16, 16, 16, 255, 201, 0, 11, 8,
					0, 1, 0, 1, 1, 1, 17, 0, 255, 204, 0, 6, 0, 16, 16, 5, 255, 218, 0, 8, 1, 1, 0, 0, 63, 0, 210, 207,
					32, 255, 217},
			},
			{
				ContentID: "inline-jpg-image.jpg@example.com",
				ContentType: ContentTypeHeader{
					ContentType: "image/jpeg",
					Params: map[string]string{
						"name": "inline-jpg-image-name.jpg",
					},
				},
				ContentDisposition: ContentDispositionHeader{
					ContentDisposition: inline,
					Params: map[string]string{
						"filename": "inline-jpg-image-filename.jpg",
					},
				},
				Data: []byte{255, 216, 255, 219, 0, 67, 0, 3, 2, 2, 2, 2, 2, 3, 2, 2, 2, 3, 3, 3, 3, 4,
					6, 4, 4, 4, 4, 4, 8, 6, 6, 5, 6, 9, 8, 10, 10, 9, 8, 9, 9, 10, 12, 15, 12, 10, 11, 14, 11, 9, 9,
					13, 17, 13, 14, 15, 16, 16, 17, 16, 10, 12, 18, 19, 18, 16, 19, 15, 16, 16, 16, 255, 201, 0, 11, 8,
					0, 1, 0, 1, 1, 1, 17, 0, 255, 204, 0, 6, 0, 16, 16, 5, 255, 218, 0, 8, 1, 1, 0, 0, 63, 0, 210, 207,
					32, 255, 217},
			},
		},
		AttachedFiles: []AttachedFile{
			{
				ContentType: ContentTypeHeader{
					ContentType: "application/pdf",
					Params: map[string]string{
						"name": "attached-pdf-name.pdf",
					},
				},
				ContentDisposition: ContentDispositionHeader{
					ContentDisposition: attachment,
					Params: map[string]string{
						"filename": "attached-pdf-filename.pdf",
					},
				},
				Data: []byte{37, 80, 68, 70, 45, 49, 46, 13, 116, 114, 97, 105, 108, 101, 114, 60, 60,
					47, 82, 111, 111, 116, 60, 60, 47, 80, 97, 103, 101, 115, 60, 60, 47, 75, 105, 100, 115, 91, 60,
					60, 47, 77, 101, 100, 105, 97, 66, 111, 120, 91, 48, 32, 48, 32, 51, 32, 51, 93, 62, 62, 93, 62,
					62, 62, 62, 62, 62},
			},
			{
				ContentType: ContentTypeHeader{
					ContentType: "application/pdf",
					Params: map[string]string{
						"name": "attached-pdf-without-disposition.pdf",
					},
				},
				ContentDisposition: ContentDispositionHeader{
					ContentDisposition: "",
					Params:             map[string]string(nil),
				},
				Data: []byte{37, 80, 68, 70, 45, 49, 46, 13, 116, 114, 97, 105, 108, 101, 114, 60, 60,
					47, 82, 111, 111, 116, 60, 60, 47, 80, 97, 103, 101, 115, 60, 60, 47, 75, 105, 100, 115, 91, 60,
					60, 47, 77, 101, 100, 105, 97, 66, 111, 120, 91, 48, 32, 48, 32, 51, 32, 51, 93, 62, 62, 93, 62,
					62, 62, 62, 62, 62},
			},
			{
				ContentType: ContentTypeHeader{
					ContentType: "application/json",
					Params: map[string]string{
						"name": "attached-json-name.json",
					},
				},
				ContentDisposition: ContentDispositionHeader{
					ContentDisposition: attachment,
					Params: map[string]string{
						"filename": "attached-json-filename.json",
					},
				},
				Data: []byte{123, 34, 102, 111, 111, 34, 58, 34, 98, 97, 114, 34, 125},
			},
			{
				ContentType: ContentTypeHeader{
					ContentType: "text/plain",
					Params: map[string]string{
						"name": "attached-text-plain-name.txt",
					},
				},
				ContentDisposition: ContentDispositionHeader{
					ContentDisposition: attachment,
					Params: map[string]string{
						"filename": "attached-text-plain-filename.txt",
					},
				},
				Data: []byte{84, 101, 120, 116, 47, 112, 108, 97, 105, 110, 32, 99, 111, 110, 116, 101, 110, 116, 32,
					97, 115, 32, 97, 110, 32, 97, 116, 116, 97, 99, 104, 101, 100, 32, 46, 116, 120, 116, 32, 102, 105,
					108, 101, 46},
			},
			{
				ContentType: ContentTypeHeader{
					ContentType: "text/html",
					Params: map[string]string{
						"name": "attached-text-html-name.html",
					},
				},
				ContentDisposition: ContentDispositionHeader{
					ContentDisposition: attachment,
					Params: map[string]string{
						"filename": "attached-text-html-filename.html",
					},
				},
				Data: []byte{84, 101, 120, 116, 47, 104, 116, 109, 108, 32, 99, 111, 110, 116, 101, 110, 116, 32, 97,
					115, 32, 97, 110, 32, 97, 116, 116, 97, 99, 104, 101, 100, 32, 46, 104, 116, 109, 108, 32, 102,
					105, 108, 101, 46},
			},
		},
	}

	testEmailFromFile(t, fp, expectedEmail)
}
func TestParseEmailJapaneseMultipartSignedUtf8Over7bit(t *testing.T) {
	fp := "tests/test_japanese_multipart_signed_utf-8_over_7bit.txt"
	tz, _ := time.LoadLocation("Asia/Tokyo")
	expectedDate, _ := time.Parse(
		time.RFC1123Z+" (MST)",
		time.Date(2019, time.April, 1, 7, 55, 0, 0, tz).Format(time.RFC1123Z+" (MST)"))
	expectedEmail := Email{
		Headers: Headers{
			Date:    expectedDate,
			Subject: "📧 Signed Test いろは歌",
			ReplyTo: []*mail.Address{
				{
					Name:    "郵便アリス",
					Address: "alice.yuubin@example.net",
				},
			},
			Sender: &mail.Address{
				Name:    "郵便アリス",
				Address: "alice.yuubin@example.com",
			},
			From: []*mail.Address{
				{
					Name:    "郵便アリス",
					Address: "alice.yuubin@example.com",
				},
				{
					Name:    "郵便アリス",
					Address: "alice.yuubin@example.net",
				},
			},
			To: []*mail.Address{
				{
					Name:    "郵便ボブ",
					Address: "bob.yuubin@example.com",
				},
				{
					Name:    "郵便キャロル",
					Address: "carol.yuubin@example.com",
				},
			},
			Cc: []*mail.Address{
				{
					Name:    "郵便ダン",
					Address: "dan.yuubin@example.com",
				},
				{
					Name:    "郵便イーブ",
					Address: "eve.yubin@example.com",
				},
			},
			Bcc: []*mail.Address{
				{
					Name:    "郵便フランク",
					Address: "frank.yuubin@example.com",
				},
				{
					Name:    "郵便グレイス",
					Address: "grace.yubin@example.com",
				},
			},
			MessageID:  "Message-Id-1@example.com",
			InReplyTo:  []MessageId{"Message-Id-0@example.com"},
			References: []MessageId{"Message-Id-0@example.com"},
			Comments:   "Message Header Comment",
			Keywords:   []string{"Keyword 1", "Keyword 2"},
			ResentDate: expectedDate,
			ResentFrom: []*mail.Address{
				{
					Name:    "郵便アリス",
					Address: "alice.yuubin@example.net",
				},
				{
					Name:    "郵便アリス",
					Address: "alice.yuubin@example.com",
				},
			},
			ResentSender: &mail.Address{
				Name:    "郵便アリス",
				Address: "alice.yuubin@example.net",
			},
			ResentTo: []*mail.Address{
				{
					Name:    "郵便ボブ",
					Address: "bob.yuubin@example.net",
				},
				{
					Name:    "郵便キャロル",
					Address: "carol.yuubin@example.net",
				},
			},
			ResentCc: []*mail.Address{
				{
					Name:    "郵便ダン",
					Address: "dan.yuubin@example.net",
				},
				{
					Name:    "郵便イーブ",
					Address: "eve.yubin@example.net",
				},
			},
			ResentBcc: []*mail.Address{
				{
					Name:    "郵便フランク",
					Address: "frank.yuubin@example.net",
				},
				{
					Name:    "郵便グレイス",
					Address: "grace.yubin@example.net",
				},
			},
			ResentMessageID: "Message-Id-1@example.net",
			ContentType: ContentTypeHeader{
				ContentType: "multipart/signed",
				Params: map[string]string{
					"boundary": "SignedBoundaryString",
					"charset":  "utf-8",
					"micalg":   "sha1",
					"protocol": "application/pkcs7-signature",
				},
			},
			ExtraHeaders: map[string][]string{
				"X-Clacks-Overhead": {"GNU Terry Pratchett"},
				"X-Script/function/\t !\"#$%&'()*+,-./;<=>?@[\\]^_`{|}~": {
					"TEST VALUE 1\t !\"#$%&'()*+,-./0123456789:;<=>?@ABCDEFGHIJKLMNOPQRSTUVWXYZ[\\]^_` abcdefghijklmnopqrstuvwxyz{|}~",
					"TEST VALUE 2\t !\"#$%&'()*+,-./0123456789:;<=>?@ABCDEFGHIJKLMNOPQRSTUVWXYZ[\\]^_` abcdefghijklmnopqrstuvwxyz{|}~",
				},
			},
		},
		Text: `色は匂えど
散りぬるを
我が世誰ぞ
常ならん
有為の奥山
今日越えて
浅き夢見じ
酔いもせず。

Iro wa nioedo / Chirinuru o / Wa ga yo tare zo / Tsune naran / Ui no okuyama / Kyo koete Asaki yume miji / Yoi mo sezu.

とりなくこゑすゆめさませみよあけわたるひんかしをそらいろはえておきつへにほふねむれゐぬもやのうち。

天地星空
山川峰谷
雲霧室苔
人犬上末
硫黄猿生ふ為よ
榎の枝を馴れ居て。

田居に出で菜摘むわれをぞ君召すと求食り追ひゆく山城の打酔へる子ら藻葉干せよえ舟繋けぬ。`,
		AttachedFiles: []AttachedFile{
			{
				ContentType: ContentTypeHeader{
					ContentType: "application/pkcs7-signature",
					Params: map[string]string{
						"name": "smime.p7s",
					},
				},
				ContentDisposition: ContentDispositionHeader{
					ContentDisposition: attachment,
					Params: map[string]string{
						"filename": "smime.p7s",
					},
				},
				Data: []byte{130, 28, 135, 132, 117, 46, 142, 18, 97, 140, 126, 251, 159, 193, 199, 25, 58, 223, 189, 185, 227, 239, 158, 173, 108, 31, 71, 27, 133, 80, 165, 252, 133, 227, 174, 198, 132, 129, 159, 29, 246, 19, 235, 133, 80, 165, 252, 133, 227, 174, 198, 132, 129, 159, 29, 246, 19, 234, 49, 251, 238, 127, 7, 28, 104, 33, 200, 120, 71, 82, 232, 225, 38, 30, 249, 234, 214, 193, 244, 113, 147, 173, 251, 219, 158, 57, 252, 28, 113, 147, 173, 251, 225, 38, 24, 199, 239, 190, 173, 108, 31, 71, 27, 133, 80, 110, 120, 251, 231, 174, 198, 132, 129, 159, 29, 246, 19, 234, 8, 114, 30, 17, 212, 186, 58, 95, 200, 94, 59, 26, 18, 6, 124, 119, 216, 79, 174, 21, 65, 185, 227, 239, 158},
			},
		},
	}

	testEmailFromFile(t, fp, expectedEmail)
}

func TestParseEmailJapaneseMultipartSignedUtf8OverBase64(t *testing.T) {
	fp := "tests/test_japanese_multipart_signed_utf-8_over_base64.txt"
	tz, _ := time.LoadLocation("Asia/Tokyo")
	expectedDate, _ := time.Parse(
		time.RFC1123Z+" (MST)",
		time.Date(2019, time.April, 1, 7, 55, 0, 0, tz).Format(time.RFC1123Z+" (MST)"))
	expectedEmail := Email{
		Headers: Headers{
			Date:    expectedDate,
			Subject: "📧 Signed Test いろは歌",
			ReplyTo: []*mail.Address{
				{
					Name:    "郵便アリス",
					Address: "alice.yuubin@example.net",
				},
			},
			Sender: &mail.Address{
				Name:    "郵便アリス",
				Address: "alice.yuubin@example.com",
			},
			From: []*mail.Address{
				{
					Name:    "郵便アリス",
					Address: "alice.yuubin@example.com",
				},
				{
					Name:    "郵便アリス",
					Address: "alice.yuubin@example.net",
				},
			},
			To: []*mail.Address{
				{
					Name:    "郵便ボブ",
					Address: "bob.yuubin@example.com",
				},
				{
					Name:    "郵便キャロル",
					Address: "carol.yuubin@example.com",
				},
			},
			Cc: []*mail.Address{
				{
					Name:    "郵便ダン",
					Address: "dan.yuubin@example.com",
				},
				{
					Name:    "郵便イーブ",
					Address: "eve.yubin@example.com",
				},
			},
			Bcc: []*mail.Address{
				{
					Name:    "郵便フランク",
					Address: "frank.yuubin@example.com",
				},
				{
					Name:    "郵便グレイス",
					Address: "grace.yubin@example.com",
				},
			},
			MessageID:  "Message-Id-1@example.com",
			InReplyTo:  []MessageId{"Message-Id-0@example.com"},
			References: []MessageId{"Message-Id-0@example.com"},
			Comments:   "Message Header Comment",
			Keywords:   []string{"Keyword 1", "Keyword 2"},
			ResentDate: expectedDate,
			ResentFrom: []*mail.Address{
				{
					Name:    "郵便アリス",
					Address: "alice.yuubin@example.net",
				},
				{
					Name:    "郵便アリス",
					Address: "alice.yuubin@example.com",
				},
			},
			ResentSender: &mail.Address{
				Name:    "郵便アリス",
				Address: "alice.yuubin@example.net",
			},
			ResentTo: []*mail.Address{
				{
					Name:    "郵便ボブ",
					Address: "bob.yuubin@example.net",
				},
				{
					Name:    "郵便キャロル",
					Address: "carol.yuubin@example.net",
				},
			},
			ResentCc: []*mail.Address{
				{
					Name:    "郵便ダン",
					Address: "dan.yuubin@example.net",
				},
				{
					Name:    "郵便イーブ",
					Address: "eve.yubin@example.net",
				},
			},
			ResentBcc: []*mail.Address{
				{
					Name:    "郵便フランク",
					Address: "frank.yuubin@example.net",
				},
				{
					Name:    "郵便グレイス",
					Address: "grace.yubin@example.net",
				},
			},
			ResentMessageID: "Message-Id-1@example.net",
			ContentType: ContentTypeHeader{
				ContentType: "multipart/signed",
				Params: map[string]string{
					"boundary": "SignedBoundaryString",
					"charset":  "utf-8",
					"micalg":   "sha1",
					"protocol": "application/pkcs7-signature",
				},
			},
			ExtraHeaders: map[string][]string{
				"X-Clacks-Overhead": {"GNU Terry Pratchett"},
				"X-Script/function/\t !\"#$%&'()*+,-./;<=>?@[\\]^_`{|}~": {
					"TEST VALUE 1\t !\"#$%&'()*+,-./0123456789:;<=>?@ABCDEFGHIJKLMNOPQRSTUVWXYZ[\\]^_` abcdefghijklmnopqrstuvwxyz{|}~",
					"TEST VALUE 2\t !\"#$%&'()*+,-./0123456789:;<=>?@ABCDEFGHIJKLMNOPQRSTUVWXYZ[\\]^_` abcdefghijklmnopqrstuvwxyz{|}~",
				},
			},
		},
		Text: `色は匂えど
散りぬるを
我が世誰ぞ
常ならん
有為の奥山
今日越えて
浅き夢見じ
酔いもせず。

Iro wa nioedo / Chirinuru o / Wa ga yo tare zo / Tsune naran / Ui no okuyama / Kyo koete Asaki yume miji / Yoi mo sezu.

とりなくこゑすゆめさませみよあけわたるひんかしをそらいろはえておきつへにほふねむれゐぬもやのうち。

天地星空
山川峰谷
雲霧室苔
人犬上末
硫黄猿生ふ為よ
榎の枝を馴れ居て。

田居に出で菜摘むわれをぞ君召すと求食り追ひゆく山城の打酔へる子ら藻葉干せよえ舟繋けぬ。`,
		AttachedFiles: []AttachedFile{
			{
				ContentType: ContentTypeHeader{
					ContentType: "application/pkcs7-signature",
					Params: map[string]string{
						"name": "smime.p7s",
					},
				},
				ContentDisposition: ContentDispositionHeader{
					ContentDisposition: attachment,
					Params: map[string]string{
						"filename": "smime.p7s",
					},
				},
				Data: []byte{130, 28, 135, 132, 117, 46, 142, 18, 97, 140, 126, 251, 159, 193, 199, 25, 58, 223, 189, 185, 227, 239, 158, 173, 108, 31, 71, 27, 133, 80, 165, 252, 133, 227, 174, 198, 132, 129, 159, 29, 246, 19, 235, 133, 80, 165, 252, 133, 227, 174, 198, 132, 129, 159, 29, 246, 19, 234, 49, 251, 238, 127, 7, 28, 104, 33, 200, 120, 71, 82, 232, 225, 38, 30, 249, 234, 214, 193, 244, 113, 147, 173, 251, 219, 158, 57, 252, 28, 113, 147, 173, 251, 225, 38, 24, 199, 239, 190, 173, 108, 31, 71, 27, 133, 80, 110, 120, 251, 231, 174, 198, 132, 129, 159, 29, 246, 19, 234, 8, 114, 30, 17, 212, 186, 58, 95, 200, 94, 59, 26, 18, 6, 124, 119, 216, 79, 174, 21, 65, 185, 227, 239, 158},
			},
		},
	}

	testEmailFromFile(t, fp, expectedEmail)
}

func TestParseEmailJapaneseMultipartSignedUtf8OverQuotedprintable(t *testing.T) {
	fp := "tests/test_japanese_multipart_signed_utf-8_over_quoted-printable.txt"
	tz, _ := time.LoadLocation("Asia/Tokyo")
	expectedDate, _ := time.Parse(
		time.RFC1123Z+" (MST)",
		time.Date(2019, time.April, 1, 7, 55, 0, 0, tz).Format(time.RFC1123Z+" (MST)"))
	expectedEmail := Email{
		Headers: Headers{
			Date:    expectedDate,
			Subject: "📧 Signed Test いろは歌",
			ReplyTo: []*mail.Address{
				{
					Name:    "郵便アリス",
					Address: "alice.yuubin@example.net",
				},
			},
			Sender: &mail.Address{
				Name:    "郵便アリス",
				Address: "alice.yuubin@example.com",
			},
			From: []*mail.Address{
				{
					Name:    "郵便アリス",
					Address: "alice.yuubin@example.com",
				},
				{
					Name:    "郵便アリス",
					Address: "alice.yuubin@example.net",
				},
			},
			To: []*mail.Address{
				{
					Name:    "郵便ボブ",
					Address: "bob.yuubin@example.com",
				},
				{
					Name:    "郵便キャロル",
					Address: "carol.yuubin@example.com",
				},
			},
			Cc: []*mail.Address{
				{
					Name:    "郵便ダン",
					Address: "dan.yuubin@example.com",
				},
				{
					Name:    "郵便イーブ",
					Address: "eve.yubin@example.com",
				},
			},
			Bcc: []*mail.Address{
				{
					Name:    "郵便フランク",
					Address: "frank.yuubin@example.com",
				},
				{
					Name:    "郵便グレイス",
					Address: "grace.yubin@example.com",
				},
			},
			MessageID:  "Message-Id-1@example.com",
			InReplyTo:  []MessageId{"Message-Id-0@example.com"},
			References: []MessageId{"Message-Id-0@example.com"},
			Comments:   "Message Header Comment",
			Keywords:   []string{"Keyword 1", "Keyword 2"},
			ResentDate: expectedDate,
			ResentFrom: []*mail.Address{
				{
					Name:    "郵便アリス",
					Address: "alice.yuubin@example.net",
				},
				{
					Name:    "郵便アリス",
					Address: "alice.yuubin@example.com",
				},
			},
			ResentSender: &mail.Address{
				Name:    "郵便アリス",
				Address: "alice.yuubin@example.net",
			},
			ResentTo: []*mail.Address{
				{
					Name:    "郵便ボブ",
					Address: "bob.yuubin@example.net",
				},
				{
					Name:    "郵便キャロル",
					Address: "carol.yuubin@example.net",
				},
			},
			ResentCc: []*mail.Address{
				{
					Name:    "郵便ダン",
					Address: "dan.yuubin@example.net",
				},
				{
					Name:    "郵便イーブ",
					Address: "eve.yubin@example.net",
				},
			},
			ResentBcc: []*mail.Address{
				{
					Name:    "郵便フランク",
					Address: "frank.yuubin@example.net",
				},
				{
					Name:    "郵便グレイス",
					Address: "grace.yubin@example.net",
				},
			},
			ResentMessageID: "Message-Id-1@example.net",
			ContentType: ContentTypeHeader{
				ContentType: "multipart/signed",
				Params: map[string]string{
					"boundary": "SignedBoundaryString",
					"charset":  "utf-8",
					"micalg":   "sha1",
					"protocol": "application/pkcs7-signature",
				},
			},
			ExtraHeaders: map[string][]string{
				"X-Clacks-Overhead": {"GNU Terry Pratchett"},
				"X-Script/function/\t !\"#$%&'()*+,-./;<=>?@[\\]^_`{|}~": {
					"TEST VALUE 1\t !\"#$%&'()*+,-./0123456789:;<=>?@ABCDEFGHIJKLMNOPQRSTUVWXYZ[\\]^_` abcdefghijklmnopqrstuvwxyz{|}~",
					"TEST VALUE 2\t !\"#$%&'()*+,-./0123456789:;<=>?@ABCDEFGHIJKLMNOPQRSTUVWXYZ[\\]^_` abcdefghijklmnopqrstuvwxyz{|}~",
				},
			},
		},
		Text: `色は匂えど
散りぬるを
我が世誰ぞ
常ならん
有為の奥山
今日越えて
浅き夢見じ
酔いもせず。

Iro wa nioedo / Chirinuru o / Wa ga yo tare zo / Tsune naran / Ui no okuyama / Kyo koete Asaki yume miji / Yoi mo sezu.

とりなくこゑすゆめさませみよあけわたるひんかしをそらいろはえておきつへにほふねむれゐぬもやのうち。

天地星空
山川峰谷
雲霧室苔
人犬上末
硫黄猿生ふ為よ
榎の枝を馴れ居て。

田居に出で菜摘むわれをぞ君召すと求食り追ひゆく山城の打酔へる子ら藻葉干せよえ舟繋けぬ。`,
		AttachedFiles: []AttachedFile{
			{
				ContentType: ContentTypeHeader{
					ContentType: "application/pkcs7-signature",
					Params: map[string]string{
						"name": "smime.p7s",
					},
				},
				ContentDisposition: ContentDispositionHeader{
					ContentDisposition: attachment,
					Params: map[string]string{
						"filename": "smime.p7s",
					},
				},
				Data: []byte{130, 28, 135, 132, 117, 46, 142, 18, 97, 140, 126, 251, 159, 193, 199, 25, 58, 223, 189, 185, 227, 239, 158, 173, 108, 31, 71, 27, 133, 80, 165, 252, 133, 227, 174, 198, 132, 129, 159, 29, 246, 19, 235, 133, 80, 165, 252, 133, 227, 174, 198, 132, 129, 159, 29, 246, 19, 234, 49, 251, 238, 127, 7, 28, 104, 33, 200, 120, 71, 82, 232, 225, 38, 30, 249, 234, 214, 193, 244, 113, 147, 173, 251, 219, 158, 57, 252, 28, 113, 147, 173, 251, 225, 38, 24, 199, 239, 190, 173, 108, 31, 71, 27, 133, 80, 110, 120, 251, 231, 174, 198, 132, 129, 159, 29, 246, 19, 234, 8, 114, 30, 17, 212, 186, 58, 95, 200, 94, 59, 26, 18, 6, 124, 119, 216, 79, 174, 21, 65, 185, 227, 239, 158},
			},
		},
	}

	testEmailFromFile(t, fp, expectedEmail)
}

func TestParseEmailJapaneseMultipartSignedIso2022jpOver7bit(t *testing.T) {
	fp := "tests/test_japanese_multipart_signed_iso-2022-jp_over_7bit.txt"
	tz, _ := time.LoadLocation("Asia/Tokyo")
	expectedDate, _ := time.Parse(
		time.RFC1123Z+" (MST)",
		time.Date(2019, time.April, 1, 7, 55, 0, 0, tz).Format(time.RFC1123Z+" (MST)"))
	expectedEmail := Email{
		Headers: Headers{
			Date:    expectedDate,
			Subject: "📧 Signed Test いろは歌",
			ReplyTo: []*mail.Address{
				{
					Name:    "郵便アリス",
					Address: "alice.yuubin@example.net",
				},
			},
			Sender: &mail.Address{
				Name:    "郵便アリス",
				Address: "alice.yuubin@example.com",
			},
			From: []*mail.Address{
				{
					Name:    "郵便アリス",
					Address: "alice.yuubin@example.com",
				},
				{
					Name:    "郵便アリス",
					Address: "alice.yuubin@example.net",
				},
			},
			To: []*mail.Address{
				{
					Name:    "郵便ボブ",
					Address: "bob.yuubin@example.com",
				},
				{
					Name:    "郵便キャロル",
					Address: "carol.yuubin@example.com",
				},
			},
			Cc: []*mail.Address{
				{
					Name:    "郵便ダン",
					Address: "dan.yuubin@example.com",
				},
				{
					Name:    "郵便イーブ",
					Address: "eve.yubin@example.com",
				},
			},
			Bcc: []*mail.Address{
				{
					Name:    "郵便フランク",
					Address: "frank.yuubin@example.com",
				},
				{
					Name:    "郵便グレイス",
					Address: "grace.yubin@example.com",
				},
			},
			MessageID:  "Message-Id-1@example.com",
			InReplyTo:  []MessageId{"Message-Id-0@example.com"},
			References: []MessageId{"Message-Id-0@example.com"},
			Comments:   "Message Header Comment",
			Keywords:   []string{"Keyword 1", "Keyword 2"},
			ResentDate: expectedDate,
			ResentFrom: []*mail.Address{
				{
					Name:    "郵便アリス",
					Address: "alice.yuubin@example.net",
				},
				{
					Name:    "郵便アリス",
					Address: "alice.yuubin@example.com",
				},
			},
			ResentSender: &mail.Address{
				Name:    "郵便アリス",
				Address: "alice.yuubin@example.net",
			},
			ResentTo: []*mail.Address{
				{
					Name:    "郵便ボブ",
					Address: "bob.yuubin@example.net",
				},
				{
					Name:    "郵便キャロル",
					Address: "carol.yuubin@example.net",
				},
			},
			ResentCc: []*mail.Address{
				{
					Name:    "郵便ダン",
					Address: "dan.yuubin@example.net",
				},
				{
					Name:    "郵便イーブ",
					Address: "eve.yubin@example.net",
				},
			},
			ResentBcc: []*mail.Address{
				{
					Name:    "郵便フランク",
					Address: "frank.yuubin@example.net",
				},
				{
					Name:    "郵便グレイス",
					Address: "grace.yubin@example.net",
				},
			},
			ResentMessageID: "Message-Id-1@example.net",
			ContentType: ContentTypeHeader{
				ContentType: "multipart/signed",
				Params: map[string]string{
					"boundary": "SignedBoundaryString",
					"charset":  "iso-2022-jp",
					"micalg":   "sha1",
					"protocol": "application/pkcs7-signature",
				},
			},
			ExtraHeaders: map[string][]string{
				"X-Clacks-Overhead": {"GNU Terry Pratchett"},
				"X-Script/function/\t !\"#$%&'()*+,-./;<=>?@[\\]^_`{|}~": {
					"TEST VALUE 1\t !\"#$%&'()*+,-./0123456789:;<=>?@ABCDEFGHIJKLMNOPQRSTUVWXYZ[\\]^_` abcdefghijklmnopqrstuvwxyz{|}~",
					"TEST VALUE 2\t !\"#$%&'()*+,-./0123456789:;<=>?@ABCDEFGHIJKLMNOPQRSTUVWXYZ[\\]^_` abcdefghijklmnopqrstuvwxyz{|}~",
				},
			},
		},
		Text: `色は匂えど
散りぬるを
我が世誰ぞ
常ならん
有為の奥山
今日越えて
浅き夢見じ
酔いもせず。

Iro wa nioedo / Chirinuru o / Wa ga yo tare zo / Tsune naran / Ui no okuyama / Kyo koete Asaki yume miji / Yoi mo sezu.

とりなくこゑすゆめさませみよあけわたるひんかしをそらいろはえておきつへにほふねむれゐぬもやのうち。

天地星空
山川峰谷
雲霧室苔
人犬上末
硫黄猿生ふ為よ
榎の枝を馴れ居て。

田居に出で菜摘むわれをぞ君召すと求食り追ひゆく山城の打酔へる子ら藻葉干せよえ舟繋けぬ。`,
		AttachedFiles: []AttachedFile{
			{
				ContentType: ContentTypeHeader{
					ContentType: "application/pkcs7-signature",
					Params: map[string]string{
						"name": "smime.p7s",
					},
				},
				ContentDisposition: ContentDispositionHeader{
					ContentDisposition: attachment,
					Params: map[string]string{
						"filename": "smime.p7s",
					},
				},
				Data: []byte{130, 28, 135, 132, 117, 46, 142, 18, 97, 140, 126, 251, 159, 193, 199, 25, 58, 223, 189, 185, 227, 239, 158, 173, 108, 31, 71, 27, 133, 80, 165, 252, 133, 227, 174, 198, 132, 129, 159, 29, 246, 19, 235, 133, 80, 165, 252, 133, 227, 174, 198, 132, 129, 159, 29, 246, 19, 234, 49, 251, 238, 127, 7, 28, 104, 33, 200, 120, 71, 82, 232, 225, 38, 30, 249, 234, 214, 193, 244, 113, 147, 173, 251, 219, 158, 57, 252, 28, 113, 147, 173, 251, 225, 38, 24, 199, 239, 190, 173, 108, 31, 71, 27, 133, 80, 110, 120, 251, 231, 174, 198, 132, 129, 159, 29, 246, 19, 234, 8, 114, 30, 17, 212, 186, 58, 95, 200, 94, 59, 26, 18, 6, 124, 119, 216, 79, 174, 21, 65, 185, 227, 239, 158},
			},
		},
	}

	testEmailFromFile(t, fp, expectedEmail)
}

func TestParseEmailJapaneseMultipartSignedIso2022jpOverBase64(t *testing.T) {
	fp := "tests/test_japanese_multipart_signed_iso-2022-jp_over_base64.txt"
	tz, _ := time.LoadLocation("Asia/Tokyo")
	expectedDate, _ := time.Parse(
		time.RFC1123Z+" (MST)",
		time.Date(2019, time.April, 1, 7, 55, 0, 0, tz).Format(time.RFC1123Z+" (MST)"))
	expectedEmail := Email{
		Headers: Headers{
			Date:    expectedDate,
			Subject: "📧 Signed Test いろは歌",
			ReplyTo: []*mail.Address{
				{
					Name:    "郵便アリス",
					Address: "alice.yuubin@example.net",
				},
			},
			Sender: &mail.Address{
				Name:    "郵便アリス",
				Address: "alice.yuubin@example.com",
			},
			From: []*mail.Address{
				{
					Name:    "郵便アリス",
					Address: "alice.yuubin@example.com",
				},
				{
					Name:    "郵便アリス",
					Address: "alice.yuubin@example.net",
				},
			},
			To: []*mail.Address{
				{
					Name:    "郵便ボブ",
					Address: "bob.yuubin@example.com",
				},
				{
					Name:    "郵便キャロル",
					Address: "carol.yuubin@example.com",
				},
			},
			Cc: []*mail.Address{
				{
					Name:    "郵便ダン",
					Address: "dan.yuubin@example.com",
				},
				{
					Name:    "郵便イーブ",
					Address: "eve.yubin@example.com",
				},
			},
			Bcc: []*mail.Address{
				{
					Name:    "郵便フランク",
					Address: "frank.yuubin@example.com",
				},
				{
					Name:    "郵便グレイス",
					Address: "grace.yubin@example.com",
				},
			},
			MessageID:  "Message-Id-1@example.com",
			InReplyTo:  []MessageId{"Message-Id-0@example.com"},
			References: []MessageId{"Message-Id-0@example.com"},
			Comments:   "Message Header Comment",
			Keywords:   []string{"Keyword 1", "Keyword 2"},
			ResentDate: expectedDate,
			ResentFrom: []*mail.Address{
				{
					Name:    "郵便アリス",
					Address: "alice.yuubin@example.net",
				},
				{
					Name:    "郵便アリス",
					Address: "alice.yuubin@example.com",
				},
			},
			ResentSender: &mail.Address{
				Name:    "郵便アリス",
				Address: "alice.yuubin@example.net",
			},
			ResentTo: []*mail.Address{
				{
					Name:    "郵便ボブ",
					Address: "bob.yuubin@example.net",
				},
				{
					Name:    "郵便キャロル",
					Address: "carol.yuubin@example.net",
				},
			},
			ResentCc: []*mail.Address{
				{
					Name:    "郵便ダン",
					Address: "dan.yuubin@example.net",
				},
				{
					Name:    "郵便イーブ",
					Address: "eve.yubin@example.net",
				},
			},
			ResentBcc: []*mail.Address{
				{
					Name:    "郵便フランク",
					Address: "frank.yuubin@example.net",
				},
				{
					Name:    "郵便グレイス",
					Address: "grace.yubin@example.net",
				},
			},
			ResentMessageID: "Message-Id-1@example.net",
			ContentType: ContentTypeHeader{
				ContentType: "multipart/signed",
				Params: map[string]string{
					"boundary": "SignedBoundaryString",
					"charset":  "iso-2022-jp",
					"micalg":   "sha1",
					"protocol": "application/pkcs7-signature",
				},
			},
			ExtraHeaders: map[string][]string{
				"X-Clacks-Overhead": {"GNU Terry Pratchett"},
				"X-Script/function/\t !\"#$%&'()*+,-./;<=>?@[\\]^_`{|}~": {
					"TEST VALUE 1\t !\"#$%&'()*+,-./0123456789:;<=>?@ABCDEFGHIJKLMNOPQRSTUVWXYZ[\\]^_` abcdefghijklmnopqrstuvwxyz{|}~",
					"TEST VALUE 2\t !\"#$%&'()*+,-./0123456789:;<=>?@ABCDEFGHIJKLMNOPQRSTUVWXYZ[\\]^_` abcdefghijklmnopqrstuvwxyz{|}~",
				},
			},
		},
		Text: `色は匂えど
散りぬるを
我が世誰ぞ
常ならん
有為の奥山
今日越えて
浅き夢見じ
酔いもせず。

Iro wa nioedo / Chirinuru o / Wa ga yo tare zo / Tsune naran / Ui no okuyama / Kyo koete Asaki yume miji / Yoi mo sezu.

とりなくこゑすゆめさませみよあけわたるひんかしをそらいろはえておきつへにほふねむれゐぬもやのうち。

天地星空
山川峰谷
雲霧室苔
人犬上末
硫黄猿生ふ為よ
榎の枝を馴れ居て。

田居に出で菜摘むわれをぞ君召すと求食り追ひゆく山城の打酔へる子ら藻葉干せよえ舟繋けぬ。`,
		AttachedFiles: []AttachedFile{
			{
				ContentType: ContentTypeHeader{
					ContentType: "application/pkcs7-signature",
					Params: map[string]string{
						"name": "smime.p7s",
					},
				},
				ContentDisposition: ContentDispositionHeader{
					ContentDisposition: attachment,
					Params: map[string]string{
						"filename": "smime.p7s",
					},
				},
				Data: []byte{130, 28, 135, 132, 117, 46, 142, 18, 97, 140, 126, 251, 159, 193, 199, 25, 58, 223, 189, 185, 227, 239, 158, 173, 108, 31, 71, 27, 133, 80, 165, 252, 133, 227, 174, 198, 132, 129, 159, 29, 246, 19, 235, 133, 80, 165, 252, 133, 227, 174, 198, 132, 129, 159, 29, 246, 19, 234, 49, 251, 238, 127, 7, 28, 104, 33, 200, 120, 71, 82, 232, 225, 38, 30, 249, 234, 214, 193, 244, 113, 147, 173, 251, 219, 158, 57, 252, 28, 113, 147, 173, 251, 225, 38, 24, 199, 239, 190, 173, 108, 31, 71, 27, 133, 80, 110, 120, 251, 231, 174, 198, 132, 129, 159, 29, 246, 19, 234, 8, 114, 30, 17, 212, 186, 58, 95, 200, 94, 59, 26, 18, 6, 124, 119, 216, 79, 174, 21, 65, 185, 227, 239, 158},
			},
		},
	}

	testEmailFromFile(t, fp, expectedEmail)
}

func TestParseEmailJapaneseMultipartSignedIso2022jpOverQuotedprintable(t *testing.T) {
	fp := "tests/test_japanese_multipart_signed_iso-2022-jp_over_quoted-printable.txt"
	tz, _ := time.LoadLocation("Asia/Tokyo")
	expectedDate, _ := time.Parse(
		time.RFC1123Z+" (MST)",
		time.Date(2019, time.April, 1, 7, 55, 0, 0, tz).Format(time.RFC1123Z+" (MST)"))
	expectedEmail := Email{
		Headers: Headers{
			Date:    expectedDate,
			Subject: "📧 Signed Test いろは歌",
			ReplyTo: []*mail.Address{
				{
					Name:    "郵便アリス",
					Address: "alice.yuubin@example.net",
				},
			},
			Sender: &mail.Address{
				Name:    "郵便アリス",
				Address: "alice.yuubin@example.com",
			},
			From: []*mail.Address{
				{
					Name:    "郵便アリス",
					Address: "alice.yuubin@example.com",
				},
				{
					Name:    "郵便アリス",
					Address: "alice.yuubin@example.net",
				},
			},
			To: []*mail.Address{
				{
					Name:    "郵便ボブ",
					Address: "bob.yuubin@example.com",
				},
				{
					Name:    "郵便キャロル",
					Address: "carol.yuubin@example.com",
				},
			},
			Cc: []*mail.Address{
				{
					Name:    "郵便ダン",
					Address: "dan.yuubin@example.com",
				},
				{
					Name:    "郵便イーブ",
					Address: "eve.yubin@example.com",
				},
			},
			Bcc: []*mail.Address{
				{
					Name:    "郵便フランク",
					Address: "frank.yuubin@example.com",
				},
				{
					Name:    "郵便グレイス",
					Address: "grace.yubin@example.com",
				},
			},
			MessageID:  "Message-Id-1@example.com",
			InReplyTo:  []MessageId{"Message-Id-0@example.com"},
			References: []MessageId{"Message-Id-0@example.com"},
			Comments:   "Message Header Comment",
			Keywords:   []string{"Keyword 1", "Keyword 2"},
			ResentDate: expectedDate,
			ResentFrom: []*mail.Address{
				{
					Name:    "郵便アリス",
					Address: "alice.yuubin@example.net",
				},
				{
					Name:    "郵便アリス",
					Address: "alice.yuubin@example.com",
				},
			},
			ResentSender: &mail.Address{
				Name:    "郵便アリス",
				Address: "alice.yuubin@example.net",
			},
			ResentTo: []*mail.Address{
				{
					Name:    "郵便ボブ",
					Address: "bob.yuubin@example.net",
				},
				{
					Name:    "郵便キャロル",
					Address: "carol.yuubin@example.net",
				},
			},
			ResentCc: []*mail.Address{
				{
					Name:    "郵便ダン",
					Address: "dan.yuubin@example.net",
				},
				{
					Name:    "郵便イーブ",
					Address: "eve.yubin@example.net",
				},
			},
			ResentBcc: []*mail.Address{
				{
					Name:    "郵便フランク",
					Address: "frank.yuubin@example.net",
				},
				{
					Name:    "郵便グレイス",
					Address: "grace.yubin@example.net",
				},
			},
			ResentMessageID: "Message-Id-1@example.net",
			ContentType: ContentTypeHeader{
				ContentType: "multipart/signed",
				Params: map[string]string{
					"boundary": "SignedBoundaryString",
					"charset":  "iso-2022-jp",
					"micalg":   "sha1",
					"protocol": "application/pkcs7-signature",
				},
			},
			ExtraHeaders: map[string][]string{
				"X-Clacks-Overhead": {"GNU Terry Pratchett"},
				"X-Script/function/\t !\"#$%&'()*+,-./;<=>?@[\\]^_`{|}~": {
					"TEST VALUE 1\t !\"#$%&'()*+,-./0123456789:;<=>?@ABCDEFGHIJKLMNOPQRSTUVWXYZ[\\]^_` abcdefghijklmnopqrstuvwxyz{|}~",
					"TEST VALUE 2\t !\"#$%&'()*+,-./0123456789:;<=>?@ABCDEFGHIJKLMNOPQRSTUVWXYZ[\\]^_` abcdefghijklmnopqrstuvwxyz{|}~",
				},
			},
		},
		Text: `色は匂えど
散りぬるを
我が世誰ぞ
常ならん
有為の奥山
今日越えて
浅き夢見じ
酔いもせず。

Iro wa nioedo / Chirinuru o / Wa ga yo tare zo / Tsune naran / Ui no okuyama / Kyo koete Asaki yume miji / Yoi mo sezu.

とりなくこゑすゆめさませみよあけわたるひんかしをそらいろはえておきつへにほふねむれゐぬもやのうち。

天地星空
山川峰谷
雲霧室苔
人犬上末
硫黄猿生ふ為よ
榎の枝を馴れ居て。

田居に出で菜摘むわれをぞ君召すと求食り追ひゆく山城の打酔へる子ら藻葉干せよえ舟繋けぬ。`,
		AttachedFiles: []AttachedFile{
			{
				ContentType: ContentTypeHeader{
					ContentType: "application/pkcs7-signature",
					Params: map[string]string{
						"name": "smime.p7s",
					},
				},
				ContentDisposition: ContentDispositionHeader{
					ContentDisposition: attachment,
					Params: map[string]string{
						"filename": "smime.p7s",
					},
				},
				Data: []byte{130, 28, 135, 132, 117, 46, 142, 18, 97, 140, 126, 251, 159, 193, 199, 25, 58, 223, 189, 185, 227, 239, 158, 173, 108, 31, 71, 27, 133, 80, 165, 252, 133, 227, 174, 198, 132, 129, 159, 29, 246, 19, 235, 133, 80, 165, 252, 133, 227, 174, 198, 132, 129, 159, 29, 246, 19, 234, 49, 251, 238, 127, 7, 28, 104, 33, 200, 120, 71, 82, 232, 225, 38, 30, 249, 234, 214, 193, 244, 113, 147, 173, 251, 219, 158, 57, 252, 28, 113, 147, 173, 251, 225, 38, 24, 199, 239, 190, 173, 108, 31, 71, 27, 133, 80, 110, 120, 251, 231, 174, 198, 132, 129, 159, 29, 246, 19, 234, 8, 114, 30, 17, 212, 186, 58, 95, 200, 94, 59, 26, 18, 6, 124, 119, 216, 79, 174, 21, 65, 185, 227, 239, 158},
			},
		},
	}

	testEmailFromFile(t, fp, expectedEmail)
}

func TestParseEmailJapaneseMultipartSignedEucjpOverBase64(t *testing.T) {
	fp := "tests/test_japanese_multipart_signed_euc-jp_over_base64.txt"
	tz, _ := time.LoadLocation("Asia/Tokyo")
	expectedDate, _ := time.Parse(
		time.RFC1123Z+" (MST)",
		time.Date(2019, time.April, 1, 7, 55, 0, 0, tz).Format(time.RFC1123Z+" (MST)"))
	expectedEmail := Email{
		Headers: Headers{
			Date:    expectedDate,
			Subject: "📧 Signed Test いろは歌",
			ReplyTo: []*mail.Address{
				{
					Name:    "郵便アリス",
					Address: "alice.yuubin@example.net",
				},
			},
			Sender: &mail.Address{
				Name:    "郵便アリス",
				Address: "alice.yuubin@example.com",
			},
			From: []*mail.Address{
				{
					Name:    "郵便アリス",
					Address: "alice.yuubin@example.com",
				},
				{
					Name:    "郵便アリス",
					Address: "alice.yuubin@example.net",
				},
			},
			To: []*mail.Address{
				{
					Name:    "郵便ボブ",
					Address: "bob.yuubin@example.com",
				},
				{
					Name:    "郵便キャロル",
					Address: "carol.yuubin@example.com",
				},
			},
			Cc: []*mail.Address{
				{
					Name:    "郵便ダン",
					Address: "dan.yuubin@example.com",
				},
				{
					Name:    "郵便イーブ",
					Address: "eve.yubin@example.com",
				},
			},
			Bcc: []*mail.Address{
				{
					Name:    "郵便フランク",
					Address: "frank.yuubin@example.com",
				},
				{
					Name:    "郵便グレイス",
					Address: "grace.yubin@example.com",
				},
			},
			MessageID:  "Message-Id-1@example.com",
			InReplyTo:  []MessageId{"Message-Id-0@example.com"},
			References: []MessageId{"Message-Id-0@example.com"},
			Comments:   "Message Header Comment",
			Keywords:   []string{"Keyword 1", "Keyword 2"},
			ResentDate: expectedDate,
			ResentFrom: []*mail.Address{
				{
					Name:    "郵便アリス",
					Address: "alice.yuubin@example.net",
				},
				{
					Name:    "郵便アリス",
					Address: "alice.yuubin@example.com",
				},
			},
			ResentSender: &mail.Address{
				Name:    "郵便アリス",
				Address: "alice.yuubin@example.net",
			},
			ResentTo: []*mail.Address{
				{
					Name:    "郵便ボブ",
					Address: "bob.yuubin@example.net",
				},
				{
					Name:    "郵便キャロル",
					Address: "carol.yuubin@example.net",
				},
			},
			ResentCc: []*mail.Address{
				{
					Name:    "郵便ダン",
					Address: "dan.yuubin@example.net",
				},
				{
					Name:    "郵便イーブ",
					Address: "eve.yubin@example.net",
				},
			},
			ResentBcc: []*mail.Address{
				{
					Name:    "郵便フランク",
					Address: "frank.yuubin@example.net",
				},
				{
					Name:    "郵便グレイス",
					Address: "grace.yubin@example.net",
				},
			},
			ResentMessageID: "Message-Id-1@example.net",
			ContentType: ContentTypeHeader{
				ContentType: "multipart/signed",
				Params: map[string]string{
					"boundary": "SignedBoundaryString",
					"charset":  "euc-jp",
					"micalg":   "sha1",
					"protocol": "application/pkcs7-signature",
				},
			},
			ExtraHeaders: map[string][]string{
				"X-Clacks-Overhead": {"GNU Terry Pratchett"},
				"X-Script/function/\t !\"#$%&'()*+,-./;<=>?@[\\]^_`{|}~": {
					"TEST VALUE 1\t !\"#$%&'()*+,-./0123456789:;<=>?@ABCDEFGHIJKLMNOPQRSTUVWXYZ[\\]^_` abcdefghijklmnopqrstuvwxyz{|}~",
					"TEST VALUE 2\t !\"#$%&'()*+,-./0123456789:;<=>?@ABCDEFGHIJKLMNOPQRSTUVWXYZ[\\]^_` abcdefghijklmnopqrstuvwxyz{|}~",
				},
			},
		},
		Text: `色は匂えど
散りぬるを
我が世誰ぞ
常ならん
有為の奥山
今日越えて
浅き夢見じ
酔いもせず。

Iro wa nioedo / Chirinuru o / Wa ga yo tare zo / Tsune naran / Ui no okuyama / Kyo koete Asaki yume miji / Yoi mo sezu.

とりなくこゑすゆめさませみよあけわたるひんかしをそらいろはえておきつへにほふねむれゐぬもやのうち。

天地星空
山川峰谷
雲霧室苔
人犬上末
硫黄猿生ふ為よ
榎の枝を馴れ居て。

田居に出で菜摘むわれをぞ君召すと求食り追ひゆく山城の打酔へる子ら藻葉干せよえ舟繋けぬ。`,
		AttachedFiles: []AttachedFile{
			{
				ContentType: ContentTypeHeader{
					ContentType: "application/pkcs7-signature",
					Params: map[string]string{
						"name": "smime.p7s",
					},
				},
				ContentDisposition: ContentDispositionHeader{
					ContentDisposition: attachment,
					Params: map[string]string{
						"filename": "smime.p7s",
					},
				},
				Data: []byte{130, 28, 135, 132, 117, 46, 142, 18, 97, 140, 126, 251, 159, 193, 199, 25, 58, 223, 189, 185, 227, 239, 158, 173, 108, 31, 71, 27, 133, 80, 165, 252, 133, 227, 174, 198, 132, 129, 159, 29, 246, 19, 235, 133, 80, 165, 252, 133, 227, 174, 198, 132, 129, 159, 29, 246, 19, 234, 49, 251, 238, 127, 7, 28, 104, 33, 200, 120, 71, 82, 232, 225, 38, 30, 249, 234, 214, 193, 244, 113, 147, 173, 251, 219, 158, 57, 252, 28, 113, 147, 173, 251, 225, 38, 24, 199, 239, 190, 173, 108, 31, 71, 27, 133, 80, 110, 120, 251, 231, 174, 198, 132, 129, 159, 29, 246, 19, 234, 8, 114, 30, 17, 212, 186, 58, 95, 200, 94, 59, 26, 18, 6, 124, 119, 216, 79, 174, 21, 65, 185, 227, 239, 158},
			},
		},
	}

	testEmailFromFile(t, fp, expectedEmail)
}

func TestParseEmailJapaneseMultipartSignedEucjpOverQuotedprintable(t *testing.T) {
	fp := "tests/test_japanese_multipart_signed_euc-jp_over_quoted-printable.txt"
	tz, _ := time.LoadLocation("Asia/Tokyo")
	expectedDate, _ := time.Parse(
		time.RFC1123Z+" (MST)",
		time.Date(2019, time.April, 1, 7, 55, 0, 0, tz).Format(time.RFC1123Z+" (MST)"))
	expectedEmail := Email{
		Headers: Headers{
			Date:    expectedDate,
			Subject: "📧 Signed Test いろは歌",
			ReplyTo: []*mail.Address{
				{
					Name:    "郵便アリス",
					Address: "alice.yuubin@example.net",
				},
			},
			Sender: &mail.Address{
				Name:    "郵便アリス",
				Address: "alice.yuubin@example.com",
			},
			From: []*mail.Address{
				{
					Name:    "郵便アリス",
					Address: "alice.yuubin@example.com",
				},
				{
					Name:    "郵便アリス",
					Address: "alice.yuubin@example.net",
				},
			},
			To: []*mail.Address{
				{
					Name:    "郵便ボブ",
					Address: "bob.yuubin@example.com",
				},
				{
					Name:    "郵便キャロル",
					Address: "carol.yuubin@example.com",
				},
			},
			Cc: []*mail.Address{
				{
					Name:    "郵便ダン",
					Address: "dan.yuubin@example.com",
				},
				{
					Name:    "郵便イーブ",
					Address: "eve.yubin@example.com",
				},
			},
			Bcc: []*mail.Address{
				{
					Name:    "郵便フランク",
					Address: "frank.yuubin@example.com",
				},
				{
					Name:    "郵便グレイス",
					Address: "grace.yubin@example.com",
				},
			},
			MessageID:  "Message-Id-1@example.com",
			InReplyTo:  []MessageId{"Message-Id-0@example.com"},
			References: []MessageId{"Message-Id-0@example.com"},
			Comments:   "Message Header Comment",
			Keywords:   []string{"Keyword 1", "Keyword 2"},
			ResentDate: expectedDate,
			ResentFrom: []*mail.Address{
				{
					Name:    "郵便アリス",
					Address: "alice.yuubin@example.net",
				},
				{
					Name:    "郵便アリス",
					Address: "alice.yuubin@example.com",
				},
			},
			ResentSender: &mail.Address{
				Name:    "郵便アリス",
				Address: "alice.yuubin@example.net",
			},
			ResentTo: []*mail.Address{
				{
					Name:    "郵便ボブ",
					Address: "bob.yuubin@example.net",
				},
				{
					Name:    "郵便キャロル",
					Address: "carol.yuubin@example.net",
				},
			},
			ResentCc: []*mail.Address{
				{
					Name:    "郵便ダン",
					Address: "dan.yuubin@example.net",
				},
				{
					Name:    "郵便イーブ",
					Address: "eve.yubin@example.net",
				},
			},
			ResentBcc: []*mail.Address{
				{
					Name:    "郵便フランク",
					Address: "frank.yuubin@example.net",
				},
				{
					Name:    "郵便グレイス",
					Address: "grace.yubin@example.net",
				},
			},
			ResentMessageID: "Message-Id-1@example.net",
			ContentType: ContentTypeHeader{
				ContentType: "multipart/signed",
				Params: map[string]string{
					"boundary": "SignedBoundaryString",
					"charset":  "euc-jp",
					"micalg":   "sha1",
					"protocol": "application/pkcs7-signature",
				},
			},
			ExtraHeaders: map[string][]string{
				"X-Clacks-Overhead": {"GNU Terry Pratchett"},
				"X-Script/function/\t !\"#$%&'()*+,-./;<=>?@[\\]^_`{|}~": {
					"TEST VALUE 1\t !\"#$%&'()*+,-./0123456789:;<=>?@ABCDEFGHIJKLMNOPQRSTUVWXYZ[\\]^_` abcdefghijklmnopqrstuvwxyz{|}~",
					"TEST VALUE 2\t !\"#$%&'()*+,-./0123456789:;<=>?@ABCDEFGHIJKLMNOPQRSTUVWXYZ[\\]^_` abcdefghijklmnopqrstuvwxyz{|}~",
				},
			},
		},
		Text: `色は匂えど
散りぬるを
我が世誰ぞ
常ならん
有為の奥山
今日越えて
浅き夢見じ
酔いもせず。

Iro wa nioedo / Chirinuru o / Wa ga yo tare zo / Tsune naran / Ui no okuyama / Kyo koete Asaki yume miji / Yoi mo sezu.

とりなくこゑすゆめさませみよあけわたるひんかしをそらいろはえておきつへにほふねむれゐぬもやのうち。

天地星空
山川峰谷
雲霧室苔
人犬上末
硫黄猿生ふ為よ
榎の枝を馴れ居て。

田居に出で菜摘むわれをぞ君召すと求食り追ひゆく山城の打酔へる子ら藻葉干せよえ舟繋けぬ。`,
		AttachedFiles: []AttachedFile{
			{
				ContentType: ContentTypeHeader{
					ContentType: "application/pkcs7-signature",
					Params: map[string]string{
						"name": "smime.p7s",
					},
				},
				ContentDisposition: ContentDispositionHeader{
					ContentDisposition: attachment,
					Params: map[string]string{
						"filename": "smime.p7s",
					},
				},
				Data: []byte{130, 28, 135, 132, 117, 46, 142, 18, 97, 140, 126, 251, 159, 193, 199, 25, 58, 223, 189, 185, 227, 239, 158, 173, 108, 31, 71, 27, 133, 80, 165, 252, 133, 227, 174, 198, 132, 129, 159, 29, 246, 19, 235, 133, 80, 165, 252, 133, 227, 174, 198, 132, 129, 159, 29, 246, 19, 234, 49, 251, 238, 127, 7, 28, 104, 33, 200, 120, 71, 82, 232, 225, 38, 30, 249, 234, 214, 193, 244, 113, 147, 173, 251, 219, 158, 57, 252, 28, 113, 147, 173, 251, 225, 38, 24, 199, 239, 190, 173, 108, 31, 71, 27, 133, 80, 110, 120, 251, 231, 174, 198, 132, 129, 159, 29, 246, 19, 234, 8, 114, 30, 17, 212, 186, 58, 95, 200, 94, 59, 26, 18, 6, 124, 119, 216, 79, 174, 21, 65, 185, 227, 239, 158},
			},
		},
	}

	testEmailFromFile(t, fp, expectedEmail)
}

func TestParseEmailKoreanPlaintextUtf8OverBase64(t *testing.T) {
	fp := "tests/test_korean_plaintext_utf-8_over_base64.txt"
	tz, _ := time.LoadLocation("Asia/Seoul")
	expectedDate, _ := time.Parse(
		time.RFC1123Z+" (MST)",
		time.Date(2019, time.April, 1, 7, 55, 0, 0, tz).Format(time.RFC1123Z+" (MST)"))
	expectedEmail := Email{
		Headers: Headers{
			Date:    expectedDate,
			Subject: "📧 Test 한국어 팬그램",
			ReplyTo: []*mail.Address{
				{
					Name:    "앨리스 보내는사람",
					Address: "alice.bonaeneunsalam@example.net",
				},
			},
			Sender: &mail.Address{
				Name:    "앨리스 보내는사람",
				Address: "alice.bonaeneunsalam@example.com",
			},
			From: []*mail.Address{
				{
					Name:    "앨리스 보내는사람",
					Address: "alice.bonaeneunsalam@example.com",
				},
				{
					Name:    "앨리스 보내는사람",
					Address: "alice.bonaeneunsalam@example.net",
				},
			},
			To: []*mail.Address{
				{
					Name:    "밥 수신자",
					Address: "bob.susinja@example.com",
				},
				{
					Name:    "캐롤 수신자",
					Address: "carol.susinja@example.com",
				},
			},
			Cc: []*mail.Address{
				{
					Name:    "댄 수신자",
					Address: "dan.susinja@example.com",
				},
				{
					Name:    "이브 수신자",
					Address: "eve.susinja@example.com",
				},
			},
			Bcc: []*mail.Address{
				{
					Name:    "프랭크 수신자",
					Address: "frank.susinja@example.com",
				},
				{
					Name:    "그레이스 수신자",
					Address: "grace.susinja@example.com",
				},
			},
			MessageID:  "Message-Id-1@example.com",
			InReplyTo:  []MessageId{"Message-Id-0@example.com"},
			References: []MessageId{"Message-Id-0@example.com"},
			Comments:   "Message Header Comment",
			Keywords:   []string{"Keyword 1", "Keyword 2"},
			ResentDate: expectedDate,
			ResentFrom: []*mail.Address{
				{
					Name:    "앨리스 보내는사람",
					Address: "alice.bonaeneunsalam@example.net",
				},
				{
					Name:    "앨리스 보내는사람",
					Address: "alice.bonaeneunsalam@example.com",
				},
			},
			ResentSender: &mail.Address{
				Name:    "앨리스 보내는사람",
				Address: "alice.bonaeneunsalam@example.net",
			},
			ResentTo: []*mail.Address{
				{
					Name:    "밥 수신자",
					Address: "bob.susinja@example.net",
				},
				{
					Name:    "캐롤 수신자",
					Address: "carol.susinja@example.net",
				},
			},
			ResentCc: []*mail.Address{
				{
					Name:    "댄 수신자",
					Address: "dan.susinja@example.net",
				},
				{
					Name:    "이브 수신자",
					Address: "eve.susinja@example.net",
				},
			},
			ResentBcc: []*mail.Address{
				{
					Name:    "프랭크 수신자",
					Address: "frank.susinja@example.net",
				},
				{
					Name:    "그레이스 수신자",
					Address: "grace.susinja@example.net",
				},
			},
			ResentMessageID: "Message-Id-1@example.net",
			ContentType: ContentTypeHeader{
				ContentType: "text/plain",
				Params: map[string]string{
					"charset": "utf-8",
				},
			},
			ExtraHeaders: map[string][]string{
				"X-Clacks-Overhead": {"GNU Terry Pratchett"},
				"X-Script/function/\t !\"#$%&'()*+,-./;<=>?@[\\]^_`{|}~": {
					"TEST VALUE 1\t !\"#$%&'()*+,-./0123456789:;<=>?@ABCDEFGHIJKLMNOPQRSTUVWXYZ[\\]^_` abcdefghijklmnopqrstuvwxyz{|}~",
					"TEST VALUE 2\t !\"#$%&'()*+,-./0123456789:;<=>?@ABCDEFGHIJKLMNOPQRSTUVWXYZ[\\]^_` abcdefghijklmnopqrstuvwxyz{|}~",
				},
			},
		},
		Text: `키스의 고유조건은 입술끼리 만나야 하고 특별한 기술은 필요치 않다.`,
	}

	testEmailFromFile(t, fp, expectedEmail)
}

func TestParseEmailKoreanPlaintextUtf8OverQuotedprintable(t *testing.T) {
	fp := "tests/test_korean_plaintext_utf-8_over_quoted-printable.txt"
	tz, _ := time.LoadLocation("Asia/Seoul")
	expectedDate, _ := time.Parse(
		time.RFC1123Z+" (MST)",
		time.Date(2019, time.April, 1, 7, 55, 0, 0, tz).Format(time.RFC1123Z+" (MST)"))
	expectedEmail := Email{
		Headers: Headers{
			Date:    expectedDate,
			Subject: "📧 Test 한국어 팬그램",
			ReplyTo: []*mail.Address{
				{
					Name:    "앨리스 보내는사람",
					Address: "alice.bonaeneunsalam@example.net",
				},
			},
			Sender: &mail.Address{
				Name:    "앨리스 보내는사람",
				Address: "alice.bonaeneunsalam@example.com",
			},
			From: []*mail.Address{
				{
					Name:    "앨리스 보내는사람",
					Address: "alice.bonaeneunsalam@example.com",
				},
				{
					Name:    "앨리스 보내는사람",
					Address: "alice.bonaeneunsalam@example.net",
				},
			},
			To: []*mail.Address{
				{
					Name:    "밥 수신자",
					Address: "bob.susinja@example.com",
				},
				{
					Name:    "캐롤 수신자",
					Address: "carol.susinja@example.com",
				},
			},
			Cc: []*mail.Address{
				{
					Name:    "댄 수신자",
					Address: "dan.susinja@example.com",
				},
				{
					Name:    "이브 수신자",
					Address: "eve.susinja@example.com",
				},
			},
			Bcc: []*mail.Address{
				{
					Name:    "프랭크 수신자",
					Address: "frank.susinja@example.com",
				},
				{
					Name:    "그레이스 수신자",
					Address: "grace.susinja@example.com",
				},
			},
			MessageID:  "Message-Id-1@example.com",
			InReplyTo:  []MessageId{"Message-Id-0@example.com"},
			References: []MessageId{"Message-Id-0@example.com"},
			Comments:   "Message Header Comment",
			Keywords:   []string{"Keyword 1", "Keyword 2"},
			ResentDate: expectedDate,
			ResentFrom: []*mail.Address{
				{
					Name:    "앨리스 보내는사람",
					Address: "alice.bonaeneunsalam@example.net",
				},
				{
					Name:    "앨리스 보내는사람",
					Address: "alice.bonaeneunsalam@example.com",
				},
			},
			ResentSender: &mail.Address{
				Name:    "앨리스 보내는사람",
				Address: "alice.bonaeneunsalam@example.net",
			},
			ResentTo: []*mail.Address{
				{
					Name:    "밥 수신자",
					Address: "bob.susinja@example.net",
				},
				{
					Name:    "캐롤 수신자",
					Address: "carol.susinja@example.net",
				},
			},
			ResentCc: []*mail.Address{
				{
					Name:    "댄 수신자",
					Address: "dan.susinja@example.net",
				},
				{
					Name:    "이브 수신자",
					Address: "eve.susinja@example.net",
				},
			},
			ResentBcc: []*mail.Address{
				{
					Name:    "프랭크 수신자",
					Address: "frank.susinja@example.net",
				},
				{
					Name:    "그레이스 수신자",
					Address: "grace.susinja@example.net",
				},
			},
			ResentMessageID: "Message-Id-1@example.net",
			ContentType: ContentTypeHeader{
				ContentType: "text/plain",
				Params: map[string]string{
					"charset": "utf-8",
				},
			},
			ExtraHeaders: map[string][]string{
				"X-Clacks-Overhead": {"GNU Terry Pratchett"},
				"X-Script/function/\t !\"#$%&'()*+,-./;<=>?@[\\]^_`{|}~": {
					"TEST VALUE 1\t !\"#$%&'()*+,-./0123456789:;<=>?@ABCDEFGHIJKLMNOPQRSTUVWXYZ[\\]^_` abcdefghijklmnopqrstuvwxyz{|}~",
					"TEST VALUE 2\t !\"#$%&'()*+,-./0123456789:;<=>?@ABCDEFGHIJKLMNOPQRSTUVWXYZ[\\]^_` abcdefghijklmnopqrstuvwxyz{|}~",
				},
			},
		},
		Text: `키스의 고유조건은 입술끼리 만나야 하고 특별한 기술은 필요치 않다.`,
	}

	testEmailFromFile(t, fp, expectedEmail)
}

func TestParseEmailKoreanPlaintextEuckrOverBase64(t *testing.T) {
	fp := "tests/test_korean_plaintext_euc-kr_over_base64.txt"
	tz, _ := time.LoadLocation("Asia/Seoul")
	expectedDate, _ := time.Parse(
		time.RFC1123Z+" (MST)",
		time.Date(2019, time.April, 1, 7, 55, 0, 0, tz).Format(time.RFC1123Z+" (MST)"))
	expectedEmail := Email{
		Headers: Headers{
			Date:    expectedDate,
			Subject: "📧 Test 한국어 팬그램",
			ReplyTo: []*mail.Address{
				{
					Name:    "앨리스 보내는사람",
					Address: "alice.bonaeneunsalam@example.net",
				},
			},
			Sender: &mail.Address{
				Name:    "앨리스 보내는사람",
				Address: "alice.bonaeneunsalam@example.com",
			},
			From: []*mail.Address{
				{
					Name:    "앨리스 보내는사람",
					Address: "alice.bonaeneunsalam@example.com",
				},
				{
					Name:    "앨리스 보내는사람",
					Address: "alice.bonaeneunsalam@example.net",
				},
			},
			To: []*mail.Address{
				{
					Name:    "밥 수신자",
					Address: "bob.susinja@example.com",
				},
				{
					Name:    "캐롤 수신자",
					Address: "carol.susinja@example.com",
				},
			},
			Cc: []*mail.Address{
				{
					Name:    "댄 수신자",
					Address: "dan.susinja@example.com",
				},
				{
					Name:    "이브 수신자",
					Address: "eve.susinja@example.com",
				},
			},
			Bcc: []*mail.Address{
				{
					Name:    "프랭크 수신자",
					Address: "frank.susinja@example.com",
				},
				{
					Name:    "그레이스 수신자",
					Address: "grace.susinja@example.com",
				},
			},
			MessageID:  "Message-Id-1@example.com",
			InReplyTo:  []MessageId{"Message-Id-0@example.com"},
			References: []MessageId{"Message-Id-0@example.com"},
			Comments:   "Message Header Comment",
			Keywords:   []string{"Keyword 1", "Keyword 2"},
			ResentDate: expectedDate,
			ResentFrom: []*mail.Address{
				{
					Name:    "앨리스 보내는사람",
					Address: "alice.bonaeneunsalam@example.net",
				},
				{
					Name:    "앨리스 보내는사람",
					Address: "alice.bonaeneunsalam@example.com",
				},
			},
			ResentSender: &mail.Address{
				Name:    "앨리스 보내는사람",
				Address: "alice.bonaeneunsalam@example.net",
			},
			ResentTo: []*mail.Address{
				{
					Name:    "밥 수신자",
					Address: "bob.susinja@example.net",
				},
				{
					Name:    "캐롤 수신자",
					Address: "carol.susinja@example.net",
				},
			},
			ResentCc: []*mail.Address{
				{
					Name:    "댄 수신자",
					Address: "dan.susinja@example.net",
				},
				{
					Name:    "이브 수신자",
					Address: "eve.susinja@example.net",
				},
			},
			ResentBcc: []*mail.Address{
				{
					Name:    "프랭크 수신자",
					Address: "frank.susinja@example.net",
				},
				{
					Name:    "그레이스 수신자",
					Address: "grace.susinja@example.net",
				},
			},
			ResentMessageID: "Message-Id-1@example.net",
			ContentType: ContentTypeHeader{
				ContentType: "text/plain",
				Params: map[string]string{
					"charset": "euc-kr",
				},
			},
			ExtraHeaders: map[string][]string{
				"X-Clacks-Overhead": {"GNU Terry Pratchett"},
				"X-Script/function/\t !\"#$%&'()*+,-./;<=>?@[\\]^_`{|}~": {
					"TEST VALUE 1\t !\"#$%&'()*+,-./0123456789:;<=>?@ABCDEFGHIJKLMNOPQRSTUVWXYZ[\\]^_` abcdefghijklmnopqrstuvwxyz{|}~",
					"TEST VALUE 2\t !\"#$%&'()*+,-./0123456789:;<=>?@ABCDEFGHIJKLMNOPQRSTUVWXYZ[\\]^_` abcdefghijklmnopqrstuvwxyz{|}~",
				},
			},
		},
		Text: `키스의 고유조건은 입술끼리 만나야 하고 특별한 기술은 필요치 않다.`,
	}

	testEmailFromFile(t, fp, expectedEmail)
}

func TestParseEmailKoreanPlaintextEuckrOverQuotedprintable(t *testing.T) {
	fp := "tests/test_korean_plaintext_euc-kr_over_quoted-printable.txt"
	tz, _ := time.LoadLocation("Asia/Seoul")
	expectedDate, _ := time.Parse(
		time.RFC1123Z+" (MST)",
		time.Date(2019, time.April, 1, 7, 55, 0, 0, tz).Format(time.RFC1123Z+" (MST)"))
	expectedEmail := Email{
		Headers: Headers{
			Date:    expectedDate,
			Subject: "📧 Test 한국어 팬그램",
			ReplyTo: []*mail.Address{
				{
					Name:    "앨리스 보내는사람",
					Address: "alice.bonaeneunsalam@example.net",
				},
			},
			Sender: &mail.Address{
				Name:    "앨리스 보내는사람",
				Address: "alice.bonaeneunsalam@example.com",
			},
			From: []*mail.Address{
				{
					Name:    "앨리스 보내는사람",
					Address: "alice.bonaeneunsalam@example.com",
				},
				{
					Name:    "앨리스 보내는사람",
					Address: "alice.bonaeneunsalam@example.net",
				},
			},
			To: []*mail.Address{
				{
					Name:    "밥 수신자",
					Address: "bob.susinja@example.com",
				},
				{
					Name:    "캐롤 수신자",
					Address: "carol.susinja@example.com",
				},
			},
			Cc: []*mail.Address{
				{
					Name:    "댄 수신자",
					Address: "dan.susinja@example.com",
				},
				{
					Name:    "이브 수신자",
					Address: "eve.susinja@example.com",
				},
			},
			Bcc: []*mail.Address{
				{
					Name:    "프랭크 수신자",
					Address: "frank.susinja@example.com",
				},
				{
					Name:    "그레이스 수신자",
					Address: "grace.susinja@example.com",
				},
			},
			MessageID:  "Message-Id-1@example.com",
			InReplyTo:  []MessageId{"Message-Id-0@example.com"},
			References: []MessageId{"Message-Id-0@example.com"},
			Comments:   "Message Header Comment",
			Keywords:   []string{"Keyword 1", "Keyword 2"},
			ResentDate: expectedDate,
			ResentFrom: []*mail.Address{
				{
					Name:    "앨리스 보내는사람",
					Address: "alice.bonaeneunsalam@example.net",
				},
				{
					Name:    "앨리스 보내는사람",
					Address: "alice.bonaeneunsalam@example.com",
				},
			},
			ResentSender: &mail.Address{
				Name:    "앨리스 보내는사람",
				Address: "alice.bonaeneunsalam@example.net",
			},
			ResentTo: []*mail.Address{
				{
					Name:    "밥 수신자",
					Address: "bob.susinja@example.net",
				},
				{
					Name:    "캐롤 수신자",
					Address: "carol.susinja@example.net",
				},
			},
			ResentCc: []*mail.Address{
				{
					Name:    "댄 수신자",
					Address: "dan.susinja@example.net",
				},
				{
					Name:    "이브 수신자",
					Address: "eve.susinja@example.net",
				},
			},
			ResentBcc: []*mail.Address{
				{
					Name:    "프랭크 수신자",
					Address: "frank.susinja@example.net",
				},
				{
					Name:    "그레이스 수신자",
					Address: "grace.susinja@example.net",
				},
			},
			ResentMessageID: "Message-Id-1@example.net",
			ContentType: ContentTypeHeader{
				ContentType: "text/plain",
				Params: map[string]string{
					"charset": "euc-kr",
				},
			},
			ExtraHeaders: map[string][]string{
				"X-Clacks-Overhead": {"GNU Terry Pratchett"},
				"X-Script/function/\t !\"#$%&'()*+,-./;<=>?@[\\]^_`{|}~": {
					"TEST VALUE 1\t !\"#$%&'()*+,-./0123456789:;<=>?@ABCDEFGHIJKLMNOPQRSTUVWXYZ[\\]^_` abcdefghijklmnopqrstuvwxyz{|}~",
					"TEST VALUE 2\t !\"#$%&'()*+,-./0123456789:;<=>?@ABCDEFGHIJKLMNOPQRSTUVWXYZ[\\]^_` abcdefghijklmnopqrstuvwxyz{|}~",
				},
			},
		},
		Text: `키스의 고유조건은 입술끼리 만나야 하고 특별한 기술은 필요치 않다.`,
	}

	testEmailFromFile(t, fp, expectedEmail)
}

func TestParseEmailKoreanMultipartRelatedUtf8OverBase64(t *testing.T) {
	fp := "tests/test_korean_multipart_related_utf-8_over_base64.txt"
	tz, _ := time.LoadLocation("Asia/Seoul")
	expectedDate, _ := time.Parse(
		time.RFC1123Z+" (MST)",
		time.Date(2019, time.April, 1, 7, 55, 0, 0, tz).Format(time.RFC1123Z+" (MST)"))
	expectedEmail := Email{
		Headers: Headers{
			Date:    expectedDate,
			Subject: "📧 Test 한국어 팬그램",
			ReplyTo: []*mail.Address{
				{
					Name:    "앨리스 보내는사람",
					Address: "alice.bonaeneunsalam@example.net",
				},
			},
			Sender: &mail.Address{
				Name:    "앨리스 보내는사람",
				Address: "alice.bonaeneunsalam@example.com",
			},
			From: []*mail.Address{
				{
					Name:    "앨리스 보내는사람",
					Address: "alice.bonaeneunsalam@example.com",
				},
				{
					Name:    "앨리스 보내는사람",
					Address: "alice.bonaeneunsalam@example.net",
				},
			},
			To: []*mail.Address{
				{
					Name:    "밥 수신자",
					Address: "bob.susinja@example.com",
				},
				{
					Name:    "캐롤 수신자",
					Address: "carol.susinja@example.com",
				},
			},
			Cc: []*mail.Address{
				{
					Name:    "댄 수신자",
					Address: "dan.susinja@example.com",
				},
				{
					Name:    "이브 수신자",
					Address: "eve.susinja@example.com",
				},
			},
			Bcc: []*mail.Address{
				{
					Name:    "프랭크 수신자",
					Address: "frank.susinja@example.com",
				},
				{
					Name:    "그레이스 수신자",
					Address: "grace.susinja@example.com",
				},
			},
			MessageID:  "Message-Id-1@example.com",
			InReplyTo:  []MessageId{"Message-Id-0@example.com"},
			References: []MessageId{"Message-Id-0@example.com"},
			Comments:   "Message Header Comment",
			Keywords:   []string{"Keyword 1", "Keyword 2"},
			ResentDate: expectedDate,
			ResentFrom: []*mail.Address{
				{
					Name:    "앨리스 보내는사람",
					Address: "alice.bonaeneunsalam@example.net",
				},
				{
					Name:    "앨리스 보내는사람",
					Address: "alice.bonaeneunsalam@example.com",
				},
			},
			ResentSender: &mail.Address{
				Name:    "앨리스 보내는사람",
				Address: "alice.bonaeneunsalam@example.net",
			},
			ResentTo: []*mail.Address{
				{
					Name:    "밥 수신자",
					Address: "bob.susinja@example.net",
				},
				{
					Name:    "캐롤 수신자",
					Address: "carol.susinja@example.net",
				},
			},
			ResentCc: []*mail.Address{
				{
					Name:    "댄 수신자",
					Address: "dan.susinja@example.net",
				},
				{
					Name:    "이브 수신자",
					Address: "eve.susinja@example.net",
				},
			},
			ResentBcc: []*mail.Address{
				{
					Name:    "프랭크 수신자",
					Address: "frank.susinja@example.net",
				},
				{
					Name:    "그레이스 수신자",
					Address: "grace.susinja@example.net",
				},
			},
			ResentMessageID: "Message-Id-1@example.net",
			ContentType: ContentTypeHeader{
				ContentType: "multipart/related",
				Params: map[string]string{
					"boundary": "RelatedBoundaryString",
					"charset":  "utf-8",
				},
			},
			ExtraHeaders: map[string][]string{
				"X-Clacks-Overhead": {"GNU Terry Pratchett"},
				"X-Script/function/\t !\"#$%&'()*+,-./;<=>?@[\\]^_`{|}~": {
					"TEST VALUE 1\t !\"#$%&'()*+,-./0123456789:;<=>?@ABCDEFGHIJKLMNOPQRSTUVWXYZ[\\]^_` abcdefghijklmnopqrstuvwxyz{|}~",
					"TEST VALUE 2\t !\"#$%&'()*+,-./0123456789:;<=>?@ABCDEFGHIJKLMNOPQRSTUVWXYZ[\\]^_` abcdefghijklmnopqrstuvwxyz{|}~",
				},
			},
		},
		Text:         `키스의 고유조건은 입술끼리 만나야 하고 특별한 기술은 필요치 않다.`,
		EnrichedText: `<bold>키스의</bold> <italic>고유조건은</italic> <fixed>입술끼리</fixed> <underline>만나야</underline> 하고 특별한 기술은 필요치 않다.`,
		HTML: `<html>
<div dir="ltr">
<p>키스의 고유조건은 입술끼리 만나야 하고 특별한 기술은 필요치 않다.</p>
</div>
</html>`,
		InlineFiles: []InlineFile{
			{
				ContentID: "inline-jpg-image.jpg@example.com",
				ContentType: ContentTypeHeader{
					ContentType: "image/jpeg",
					Params: map[string]string{
						"name": "inline-jpg-image-name.jpg",
					},
				},
				ContentDisposition: ContentDispositionHeader{
					ContentDisposition: inline,
					Params: map[string]string{
						"filename": "inline-jpg-image-filename.jpg",
					},
				},
				Data: []byte{255, 216, 255, 219, 0, 67, 0, 3, 2, 2, 2, 2, 2, 3, 2, 2, 2, 3, 3, 3, 3, 4,
					6, 4, 4, 4, 4, 4, 8, 6, 6, 5, 6, 9, 8, 10, 10, 9, 8, 9, 9, 10, 12, 15, 12, 10, 11, 14, 11, 9, 9,
					13, 17, 13, 14, 15, 16, 16, 17, 16, 10, 12, 18, 19, 18, 16, 19, 15, 16, 16, 16, 255, 201, 0, 11, 8,
					0, 1, 0, 1, 1, 1, 17, 0, 255, 204, 0, 6, 0, 16, 16, 5, 255, 218, 0, 8, 1, 1, 0, 0, 63, 0, 210, 207,
					32, 255, 217},
			},
		},
	}

	testEmailFromFile(t, fp, expectedEmail)
}
func TestParseEmailKoreanMultipartRelatedUtf8OverQuotedprintable(t *testing.T) {
	fp := "tests/test_korean_multipart_related_utf-8_over_quoted-printable.txt"
	tz, _ := time.LoadLocation("Asia/Seoul")
	expectedDate, _ := time.Parse(
		time.RFC1123Z+" (MST)",
		time.Date(2019, time.April, 1, 7, 55, 0, 0, tz).Format(time.RFC1123Z+" (MST)"))
	expectedEmail := Email{
		Headers: Headers{
			Date:    expectedDate,
			Subject: "📧 Test 한국어 팬그램",
			ReplyTo: []*mail.Address{
				{
					Name:    "앨리스 보내는사람",
					Address: "alice.bonaeneunsalam@example.net",
				},
			},
			Sender: &mail.Address{
				Name:    "앨리스 보내는사람",
				Address: "alice.bonaeneunsalam@example.com",
			},
			From: []*mail.Address{
				{
					Name:    "앨리스 보내는사람",
					Address: "alice.bonaeneunsalam@example.com",
				},
				{
					Name:    "앨리스 보내는사람",
					Address: "alice.bonaeneunsalam@example.net",
				},
			},
			To: []*mail.Address{
				{
					Name:    "밥 수신자",
					Address: "bob.susinja@example.com",
				},
				{
					Name:    "캐롤 수신자",
					Address: "carol.susinja@example.com",
				},
			},
			Cc: []*mail.Address{
				{
					Name:    "댄 수신자",
					Address: "dan.susinja@example.com",
				},
				{
					Name:    "이브 수신자",
					Address: "eve.susinja@example.com",
				},
			},
			Bcc: []*mail.Address{
				{
					Name:    "프랭크 수신자",
					Address: "frank.susinja@example.com",
				},
				{
					Name:    "그레이스 수신자",
					Address: "grace.susinja@example.com",
				},
			},
			MessageID:  "Message-Id-1@example.com",
			InReplyTo:  []MessageId{"Message-Id-0@example.com"},
			References: []MessageId{"Message-Id-0@example.com"},
			Comments:   "Message Header Comment",
			Keywords:   []string{"Keyword 1", "Keyword 2"},
			ResentDate: expectedDate,
			ResentFrom: []*mail.Address{
				{
					Name:    "앨리스 보내는사람",
					Address: "alice.bonaeneunsalam@example.net",
				},
				{
					Name:    "앨리스 보내는사람",
					Address: "alice.bonaeneunsalam@example.com",
				},
			},
			ResentSender: &mail.Address{
				Name:    "앨리스 보내는사람",
				Address: "alice.bonaeneunsalam@example.net",
			},
			ResentTo: []*mail.Address{
				{
					Name:    "밥 수신자",
					Address: "bob.susinja@example.net",
				},
				{
					Name:    "캐롤 수신자",
					Address: "carol.susinja@example.net",
				},
			},
			ResentCc: []*mail.Address{
				{
					Name:    "댄 수신자",
					Address: "dan.susinja@example.net",
				},
				{
					Name:    "이브 수신자",
					Address: "eve.susinja@example.net",
				},
			},
			ResentBcc: []*mail.Address{
				{
					Name:    "프랭크 수신자",
					Address: "frank.susinja@example.net",
				},
				{
					Name:    "그레이스 수신자",
					Address: "grace.susinja@example.net",
				},
			},
			ResentMessageID: "Message-Id-1@example.net",
			ContentType: ContentTypeHeader{
				ContentType: "multipart/related",
				Params: map[string]string{
					"boundary": "RelatedBoundaryString",
					"charset":  "utf-8",
				},
			},
			ExtraHeaders: map[string][]string{
				"X-Clacks-Overhead": {"GNU Terry Pratchett"},
				"X-Script/function/\t !\"#$%&'()*+,-./;<=>?@[\\]^_`{|}~": {
					"TEST VALUE 1\t !\"#$%&'()*+,-./0123456789:;<=>?@ABCDEFGHIJKLMNOPQRSTUVWXYZ[\\]^_` abcdefghijklmnopqrstuvwxyz{|}~",
					"TEST VALUE 2\t !\"#$%&'()*+,-./0123456789:;<=>?@ABCDEFGHIJKLMNOPQRSTUVWXYZ[\\]^_` abcdefghijklmnopqrstuvwxyz{|}~",
				},
			},
		},
		Text:         `키스의 고유조건은 입술끼리 만나야 하고 특별한 기술은 필요치 않다.`,
		EnrichedText: `<bold>키스의</bold> <italic>고유조건은</italic> <fixed>입술끼리</fixed> <underline>만나야</underline> 하고 특별한 기술은 필요치 않다.`,
		HTML: `<html>
<div dir="ltr">
<p>키스의 고유조건은 입술끼리 만나야 하고 특별한 기술은 필요치 않다.</p>
</div>
</html>`,
		InlineFiles: []InlineFile{
			{
				ContentID: "inline-jpg-image.jpg@example.com",
				ContentType: ContentTypeHeader{
					ContentType: "image/jpeg",
					Params: map[string]string{
						"name": "inline-jpg-image-name.jpg",
					},
				},
				ContentDisposition: ContentDispositionHeader{
					ContentDisposition: inline,
					Params: map[string]string{
						"filename": "inline-jpg-image-filename.jpg",
					},
				},
				Data: []byte{255, 216, 255, 219, 0, 67, 0, 3, 2, 2, 2, 2, 2, 3, 2, 2, 2, 3, 3, 3, 3, 4,
					6, 4, 4, 4, 4, 4, 8, 6, 6, 5, 6, 9, 8, 10, 10, 9, 8, 9, 9, 10, 12, 15, 12, 10, 11, 14, 11, 9, 9,
					13, 17, 13, 14, 15, 16, 16, 17, 16, 10, 12, 18, 19, 18, 16, 19, 15, 16, 16, 16, 255, 201, 0, 11, 8,
					0, 1, 0, 1, 1, 1, 17, 0, 255, 204, 0, 6, 0, 16, 16, 5, 255, 218, 0, 8, 1, 1, 0, 0, 63, 0, 210, 207,
					32, 255, 217},
			},
		},
	}

	testEmailFromFile(t, fp, expectedEmail)
}
func TestParseEmailKoreanMultipartRelatedEuckrOverBase64(t *testing.T) {
	fp := "tests/test_korean_multipart_related_euc-kr_over_base64.txt"
	tz, _ := time.LoadLocation("Asia/Seoul")
	expectedDate, _ := time.Parse(
		time.RFC1123Z+" (MST)",
		time.Date(2019, time.April, 1, 7, 55, 0, 0, tz).Format(time.RFC1123Z+" (MST)"))
	expectedEmail := Email{
		Headers: Headers{
			Date:    expectedDate,
			Subject: "📧 Test 한국어 팬그램",
			ReplyTo: []*mail.Address{
				{
					Name:    "앨리스 보내는사람",
					Address: "alice.bonaeneunsalam@example.net",
				},
			},
			Sender: &mail.Address{
				Name:    "앨리스 보내는사람",
				Address: "alice.bonaeneunsalam@example.com",
			},
			From: []*mail.Address{
				{
					Name:    "앨리스 보내는사람",
					Address: "alice.bonaeneunsalam@example.com",
				},
				{
					Name:    "앨리스 보내는사람",
					Address: "alice.bonaeneunsalam@example.net",
				},
			},
			To: []*mail.Address{
				{
					Name:    "밥 수신자",
					Address: "bob.susinja@example.com",
				},
				{
					Name:    "캐롤 수신자",
					Address: "carol.susinja@example.com",
				},
			},
			Cc: []*mail.Address{
				{
					Name:    "댄 수신자",
					Address: "dan.susinja@example.com",
				},
				{
					Name:    "이브 수신자",
					Address: "eve.susinja@example.com",
				},
			},
			Bcc: []*mail.Address{
				{
					Name:    "프랭크 수신자",
					Address: "frank.susinja@example.com",
				},
				{
					Name:    "그레이스 수신자",
					Address: "grace.susinja@example.com",
				},
			},
			MessageID:  "Message-Id-1@example.com",
			InReplyTo:  []MessageId{"Message-Id-0@example.com"},
			References: []MessageId{"Message-Id-0@example.com"},
			Comments:   "Message Header Comment",
			Keywords:   []string{"Keyword 1", "Keyword 2"},
			ResentDate: expectedDate,
			ResentFrom: []*mail.Address{
				{
					Name:    "앨리스 보내는사람",
					Address: "alice.bonaeneunsalam@example.net",
				},
				{
					Name:    "앨리스 보내는사람",
					Address: "alice.bonaeneunsalam@example.com",
				},
			},
			ResentSender: &mail.Address{
				Name:    "앨리스 보내는사람",
				Address: "alice.bonaeneunsalam@example.net",
			},
			ResentTo: []*mail.Address{
				{
					Name:    "밥 수신자",
					Address: "bob.susinja@example.net",
				},
				{
					Name:    "캐롤 수신자",
					Address: "carol.susinja@example.net",
				},
			},
			ResentCc: []*mail.Address{
				{
					Name:    "댄 수신자",
					Address: "dan.susinja@example.net",
				},
				{
					Name:    "이브 수신자",
					Address: "eve.susinja@example.net",
				},
			},
			ResentBcc: []*mail.Address{
				{
					Name:    "프랭크 수신자",
					Address: "frank.susinja@example.net",
				},
				{
					Name:    "그레이스 수신자",
					Address: "grace.susinja@example.net",
				},
			},
			ResentMessageID: "Message-Id-1@example.net",
			ContentType: ContentTypeHeader{
				ContentType: "multipart/related",
				Params: map[string]string{
					"boundary": "RelatedBoundaryString",
					"charset":  "euc-kr",
				},
			},
			ExtraHeaders: map[string][]string{
				"X-Clacks-Overhead": {"GNU Terry Pratchett"},
				"X-Script/function/\t !\"#$%&'()*+,-./;<=>?@[\\]^_`{|}~": {
					"TEST VALUE 1\t !\"#$%&'()*+,-./0123456789:;<=>?@ABCDEFGHIJKLMNOPQRSTUVWXYZ[\\]^_` abcdefghijklmnopqrstuvwxyz{|}~",
					"TEST VALUE 2\t !\"#$%&'()*+,-./0123456789:;<=>?@ABCDEFGHIJKLMNOPQRSTUVWXYZ[\\]^_` abcdefghijklmnopqrstuvwxyz{|}~",
				},
			},
		},
		Text:         `키스의 고유조건은 입술끼리 만나야 하고 특별한 기술은 필요치 않다.`,
		EnrichedText: `<bold>키스의</bold> <italic>고유조건은</italic> <fixed>입술끼리</fixed> <underline>만나야</underline> 하고 특별한 기술은 필요치 않다.`,
		HTML: `<html>
<div dir="ltr">
<p>키스의 고유조건은 입술끼리 만나야 하고 특별한 기술은 필요치 않다.</p>
</div>
</html>`,
		InlineFiles: []InlineFile{
			{
				ContentID: "inline-jpg-image.jpg@example.com",
				ContentType: ContentTypeHeader{
					ContentType: "image/jpeg",
					Params: map[string]string{
						"name": "inline-jpg-image-name.jpg",
					},
				},
				ContentDisposition: ContentDispositionHeader{
					ContentDisposition: inline,
					Params: map[string]string{
						"filename": "inline-jpg-image-filename.jpg",
					},
				},
				Data: []byte{255, 216, 255, 219, 0, 67, 0, 3, 2, 2, 2, 2, 2, 3, 2, 2, 2, 3, 3, 3, 3, 4,
					6, 4, 4, 4, 4, 4, 8, 6, 6, 5, 6, 9, 8, 10, 10, 9, 8, 9, 9, 10, 12, 15, 12, 10, 11, 14, 11, 9, 9,
					13, 17, 13, 14, 15, 16, 16, 17, 16, 10, 12, 18, 19, 18, 16, 19, 15, 16, 16, 16, 255, 201, 0, 11, 8,
					0, 1, 0, 1, 1, 1, 17, 0, 255, 204, 0, 6, 0, 16, 16, 5, 255, 218, 0, 8, 1, 1, 0, 0, 63, 0, 210, 207,
					32, 255, 217},
			},
		},
	}

	testEmailFromFile(t, fp, expectedEmail)
}
func TestParseEmailKoreanMultipartRelatedEuckrOverQuotedprintable(t *testing.T) {
	fp := "tests/test_korean_multipart_related_euc-kr_over_quoted-printable.txt"
	tz, _ := time.LoadLocation("Asia/Seoul")
	expectedDate, _ := time.Parse(
		time.RFC1123Z+" (MST)",
		time.Date(2019, time.April, 1, 7, 55, 0, 0, tz).Format(time.RFC1123Z+" (MST)"))
	expectedEmail := Email{
		Headers: Headers{
			Date:    expectedDate,
			Subject: "📧 Test 한국어 팬그램",
			ReplyTo: []*mail.Address{
				{
					Name:    "앨리스 보내는사람",
					Address: "alice.bonaeneunsalam@example.net",
				},
			},
			Sender: &mail.Address{
				Name:    "앨리스 보내는사람",
				Address: "alice.bonaeneunsalam@example.com",
			},
			From: []*mail.Address{
				{
					Name:    "앨리스 보내는사람",
					Address: "alice.bonaeneunsalam@example.com",
				},
				{
					Name:    "앨리스 보내는사람",
					Address: "alice.bonaeneunsalam@example.net",
				},
			},
			To: []*mail.Address{
				{
					Name:    "밥 수신자",
					Address: "bob.susinja@example.com",
				},
				{
					Name:    "캐롤 수신자",
					Address: "carol.susinja@example.com",
				},
			},
			Cc: []*mail.Address{
				{
					Name:    "댄 수신자",
					Address: "dan.susinja@example.com",
				},
				{
					Name:    "이브 수신자",
					Address: "eve.susinja@example.com",
				},
			},
			Bcc: []*mail.Address{
				{
					Name:    "프랭크 수신자",
					Address: "frank.susinja@example.com",
				},
				{
					Name:    "그레이스 수신자",
					Address: "grace.susinja@example.com",
				},
			},
			MessageID:  "Message-Id-1@example.com",
			InReplyTo:  []MessageId{"Message-Id-0@example.com"},
			References: []MessageId{"Message-Id-0@example.com"},
			Comments:   "Message Header Comment",
			Keywords:   []string{"Keyword 1", "Keyword 2"},
			ResentDate: expectedDate,
			ResentFrom: []*mail.Address{
				{
					Name:    "앨리스 보내는사람",
					Address: "alice.bonaeneunsalam@example.net",
				},
				{
					Name:    "앨리스 보내는사람",
					Address: "alice.bonaeneunsalam@example.com",
				},
			},
			ResentSender: &mail.Address{
				Name:    "앨리스 보내는사람",
				Address: "alice.bonaeneunsalam@example.net",
			},
			ResentTo: []*mail.Address{
				{
					Name:    "밥 수신자",
					Address: "bob.susinja@example.net",
				},
				{
					Name:    "캐롤 수신자",
					Address: "carol.susinja@example.net",
				},
			},
			ResentCc: []*mail.Address{
				{
					Name:    "댄 수신자",
					Address: "dan.susinja@example.net",
				},
				{
					Name:    "이브 수신자",
					Address: "eve.susinja@example.net",
				},
			},
			ResentBcc: []*mail.Address{
				{
					Name:    "프랭크 수신자",
					Address: "frank.susinja@example.net",
				},
				{
					Name:    "그레이스 수신자",
					Address: "grace.susinja@example.net",
				},
			},
			ResentMessageID: "Message-Id-1@example.net",
			ContentType: ContentTypeHeader{
				ContentType: "multipart/related",
				Params: map[string]string{
					"boundary": "RelatedBoundaryString",
					"charset":  "euc-kr",
				},
			},
			ExtraHeaders: map[string][]string{
				"X-Clacks-Overhead": {"GNU Terry Pratchett"},
				"X-Script/function/\t !\"#$%&'()*+,-./;<=>?@[\\]^_`{|}~": {
					"TEST VALUE 1\t !\"#$%&'()*+,-./0123456789:;<=>?@ABCDEFGHIJKLMNOPQRSTUVWXYZ[\\]^_` abcdefghijklmnopqrstuvwxyz{|}~",
					"TEST VALUE 2\t !\"#$%&'()*+,-./0123456789:;<=>?@ABCDEFGHIJKLMNOPQRSTUVWXYZ[\\]^_` abcdefghijklmnopqrstuvwxyz{|}~",
				},
			},
		},
		Text:         `키스의 고유조건은 입술끼리 만나야 하고 특별한 기술은 필요치 않다.`,
		EnrichedText: `<bold>키스의</bold> <italic>고유조건은</italic> <fixed>입술끼리</fixed> <underline>만나야</underline> 하고 특별한 기술은 필요치 않다.`,
		HTML: `<html>
<div dir="ltr">
<p>키스의 고유조건은 입술끼리 만나야 하고 특별한 기술은 필요치 않다.</p>
</div>
</html>`,
		InlineFiles: []InlineFile{
			{
				ContentID: "inline-jpg-image.jpg@example.com",
				ContentType: ContentTypeHeader{
					ContentType: "image/jpeg",
					Params: map[string]string{
						"name": "inline-jpg-image-name.jpg",
					},
				},
				ContentDisposition: ContentDispositionHeader{
					ContentDisposition: inline,
					Params: map[string]string{
						"filename": "inline-jpg-image-filename.jpg",
					},
				},
				Data: []byte{255, 216, 255, 219, 0, 67, 0, 3, 2, 2, 2, 2, 2, 3, 2, 2, 2, 3, 3, 3, 3, 4,
					6, 4, 4, 4, 4, 4, 8, 6, 6, 5, 6, 9, 8, 10, 10, 9, 8, 9, 9, 10, 12, 15, 12, 10, 11, 14, 11, 9, 9,
					13, 17, 13, 14, 15, 16, 16, 17, 16, 10, 12, 18, 19, 18, 16, 19, 15, 16, 16, 16, 255, 201, 0, 11, 8,
					0, 1, 0, 1, 1, 1, 17, 0, 255, 204, 0, 6, 0, 16, 16, 5, 255, 218, 0, 8, 1, 1, 0, 0, 63, 0, 210, 207,
					32, 255, 217},
			},
		},
	}

	testEmailFromFile(t, fp, expectedEmail)
}
func TestParseEmailKoreanMultipartMixedUtf8OverBase64(t *testing.T) {
	fp := "tests/test_korean_multipart_mixed_utf-8_over_base64.txt"
	tz, _ := time.LoadLocation("Asia/Seoul")
	expectedDate, _ := time.Parse(
		time.RFC1123Z+" (MST)",
		time.Date(2019, time.April, 1, 7, 55, 0, 0, tz).Format(time.RFC1123Z+" (MST)"))
	expectedEmail := Email{
		Headers: Headers{
			Date:    expectedDate,
			Subject: "📧 Test 한국어 팬그램",
			ReplyTo: []*mail.Address{
				{
					Name:    "앨리스 보내는사람",
					Address: "alice.bonaeneunsalam@example.net",
				},
			},
			Sender: &mail.Address{
				Name:    "앨리스 보내는사람",
				Address: "alice.bonaeneunsalam@example.com",
			},
			From: []*mail.Address{
				{
					Name:    "앨리스 보내는사람",
					Address: "alice.bonaeneunsalam@example.com",
				},
				{
					Name:    "앨리스 보내는사람",
					Address: "alice.bonaeneunsalam@example.net",
				},
			},
			To: []*mail.Address{
				{
					Name:    "밥 수신자",
					Address: "bob.susinja@example.com",
				},
				{
					Name:    "캐롤 수신자",
					Address: "carol.susinja@example.com",
				},
			},
			Cc: []*mail.Address{
				{
					Name:    "댄 수신자",
					Address: "dan.susinja@example.com",
				},
				{
					Name:    "이브 수신자",
					Address: "eve.susinja@example.com",
				},
			},
			Bcc: []*mail.Address{
				{
					Name:    "프랭크 수신자",
					Address: "frank.susinja@example.com",
				},
				{
					Name:    "그레이스 수신자",
					Address: "grace.susinja@example.com",
				},
			},
			MessageID:  "Message-Id-1@example.com",
			InReplyTo:  []MessageId{"Message-Id-0@example.com"},
			References: []MessageId{"Message-Id-0@example.com"},
			Comments:   "Message Header Comment",
			Keywords:   []string{"Keyword 1", "Keyword 2"},
			ResentDate: expectedDate,
			ResentFrom: []*mail.Address{
				{
					Name:    "앨리스 보내는사람",
					Address: "alice.bonaeneunsalam@example.net",
				},
				{
					Name:    "앨리스 보내는사람",
					Address: "alice.bonaeneunsalam@example.com",
				},
			},
			ResentSender: &mail.Address{
				Name:    "앨리스 보내는사람",
				Address: "alice.bonaeneunsalam@example.net",
			},
			ResentTo: []*mail.Address{
				{
					Name:    "밥 수신자",
					Address: "bob.susinja@example.net",
				},
				{
					Name:    "캐롤 수신자",
					Address: "carol.susinja@example.net",
				},
			},
			ResentCc: []*mail.Address{
				{
					Name:    "댄 수신자",
					Address: "dan.susinja@example.net",
				},
				{
					Name:    "이브 수신자",
					Address: "eve.susinja@example.net",
				},
			},
			ResentBcc: []*mail.Address{
				{
					Name:    "프랭크 수신자",
					Address: "frank.susinja@example.net",
				},
				{
					Name:    "그레이스 수신자",
					Address: "grace.susinja@example.net",
				},
			},
			ResentMessageID: "Message-Id-1@example.net",
			ContentType: ContentTypeHeader{
				ContentType: "multipart/mixed",
				Params: map[string]string{
					"boundary": "MixedBoundaryString",
					"charset":  "utf-8",
				},
			},
			ExtraHeaders: map[string][]string{
				"X-Clacks-Overhead": {"GNU Terry Pratchett"},
				"X-Script/function/\t !\"#$%&'()*+,-./;<=>?@[\\]^_`{|}~": {
					"TEST VALUE 1\t !\"#$%&'()*+,-./0123456789:;<=>?@ABCDEFGHIJKLMNOPQRSTUVWXYZ[\\]^_` abcdefghijklmnopqrstuvwxyz{|}~",
					"TEST VALUE 2\t !\"#$%&'()*+,-./0123456789:;<=>?@ABCDEFGHIJKLMNOPQRSTUVWXYZ[\\]^_` abcdefghijklmnopqrstuvwxyz{|}~",
				},
			},
		},
		Text:         `키스의 고유조건은 입술끼리 만나야 하고 특별한 기술은 필요치 않다.`,
		EnrichedText: `<bold>키스의</bold> <italic>고유조건은</italic> <fixed>입술끼리</fixed> <underline>만나야</underline> 하고 특별한 기술은 필요치 않다.`,
		HTML: `<html>
<div dir="ltr">
<p>키스의 고유조건은 입술끼리 만나야 하고 특별한 기술은 필요치 않다.</p>
</div>
</html>`,
		InlineFiles: []InlineFile{
			{
				ContentType: ContentTypeHeader{
					ContentType: "image/jpeg",
					Params: map[string]string{
						"name": "inline-jpg-image-without-disposition.jpg",
					},
				},
				ContentDisposition: ContentDispositionHeader{
					ContentDisposition: "",
					Params:             map[string]string(nil),
				},
				Data: []byte{255, 216, 255, 219, 0, 67, 0, 3, 2, 2, 2, 2, 2, 3, 2, 2, 2, 3, 3, 3, 3, 4,
					6, 4, 4, 4, 4, 4, 8, 6, 6, 5, 6, 9, 8, 10, 10, 9, 8, 9, 9, 10, 12, 15, 12, 10, 11, 14, 11, 9, 9,
					13, 17, 13, 14, 15, 16, 16, 17, 16, 10, 12, 18, 19, 18, 16, 19, 15, 16, 16, 16, 255, 201, 0, 11, 8,
					0, 1, 0, 1, 1, 1, 17, 0, 255, 204, 0, 6, 0, 16, 16, 5, 255, 218, 0, 8, 1, 1, 0, 0, 63, 0, 210, 207,
					32, 255, 217},
			},
			{
				ContentID: "inline-jpg-image.jpg@example.com",
				ContentType: ContentTypeHeader{
					ContentType: "image/jpeg",
					Params: map[string]string{
						"name": "inline-jpg-image-name.jpg",
					},
				},
				ContentDisposition: ContentDispositionHeader{
					ContentDisposition: inline,
					Params: map[string]string{
						"filename": "inline-jpg-image-filename.jpg",
					},
				},
				Data: []byte{255, 216, 255, 219, 0, 67, 0, 3, 2, 2, 2, 2, 2, 3, 2, 2, 2, 3, 3, 3, 3, 4,
					6, 4, 4, 4, 4, 4, 8, 6, 6, 5, 6, 9, 8, 10, 10, 9, 8, 9, 9, 10, 12, 15, 12, 10, 11, 14, 11, 9, 9,
					13, 17, 13, 14, 15, 16, 16, 17, 16, 10, 12, 18, 19, 18, 16, 19, 15, 16, 16, 16, 255, 201, 0, 11, 8,
					0, 1, 0, 1, 1, 1, 17, 0, 255, 204, 0, 6, 0, 16, 16, 5, 255, 218, 0, 8, 1, 1, 0, 0, 63, 0, 210, 207,
					32, 255, 217},
			},
		},
		AttachedFiles: []AttachedFile{
			{
				ContentType: ContentTypeHeader{
					ContentType: "application/pdf",
					Params: map[string]string{
						"name": "attached-pdf-name.pdf",
					},
				},
				ContentDisposition: ContentDispositionHeader{
					ContentDisposition: attachment,
					Params: map[string]string{
						"filename": "attached-pdf-filename.pdf",
					},
				},
				Data: []byte{37, 80, 68, 70, 45, 49, 46, 13, 116, 114, 97, 105, 108, 101, 114, 60, 60,
					47, 82, 111, 111, 116, 60, 60, 47, 80, 97, 103, 101, 115, 60, 60, 47, 75, 105, 100, 115, 91, 60,
					60, 47, 77, 101, 100, 105, 97, 66, 111, 120, 91, 48, 32, 48, 32, 51, 32, 51, 93, 62, 62, 93, 62,
					62, 62, 62, 62, 62},
			},
			{
				ContentType: ContentTypeHeader{
					ContentType: "application/pdf",
					Params: map[string]string{
						"name": "attached-pdf-without-disposition.pdf",
					},
				},
				ContentDisposition: ContentDispositionHeader{
					ContentDisposition: "",
					Params:             map[string]string(nil),
				},
				Data: []byte{37, 80, 68, 70, 45, 49, 46, 13, 116, 114, 97, 105, 108, 101, 114, 60, 60,
					47, 82, 111, 111, 116, 60, 60, 47, 80, 97, 103, 101, 115, 60, 60, 47, 75, 105, 100, 115, 91, 60,
					60, 47, 77, 101, 100, 105, 97, 66, 111, 120, 91, 48, 32, 48, 32, 51, 32, 51, 93, 62, 62, 93, 62,
					62, 62, 62, 62, 62},
			},
			{
				ContentType: ContentTypeHeader{
					ContentType: "application/json",
					Params: map[string]string{
						"name": "attached-json-name.json",
					},
				},
				ContentDisposition: ContentDispositionHeader{
					ContentDisposition: attachment,
					Params: map[string]string{
						"filename": "attached-json-filename.json",
					},
				},
				Data: []byte{123, 34, 102, 111, 111, 34, 58, 34, 98, 97, 114, 34, 125},
			},
			{
				ContentType: ContentTypeHeader{
					ContentType: "text/plain",
					Params: map[string]string{
						"name": "attached-text-plain-name.txt",
					},
				},
				ContentDisposition: ContentDispositionHeader{
					ContentDisposition: attachment,
					Params: map[string]string{
						"filename": "attached-text-plain-filename.txt",
					},
				},
				Data: []byte{84, 101, 120, 116, 47, 112, 108, 97, 105, 110, 32, 99, 111, 110, 116, 101, 110, 116, 32,
					97, 115, 32, 97, 110, 32, 97, 116, 116, 97, 99, 104, 101, 100, 32, 46, 116, 120, 116, 32, 102, 105,
					108, 101, 46},
			},
			{
				ContentType: ContentTypeHeader{
					ContentType: "text/html",
					Params: map[string]string{
						"name": "attached-text-html-name.html",
					},
				},
				ContentDisposition: ContentDispositionHeader{
					ContentDisposition: attachment,
					Params: map[string]string{
						"filename": "attached-text-html-filename.html",
					},
				},
				Data: []byte{84, 101, 120, 116, 47, 104, 116, 109, 108, 32, 99, 111, 110, 116, 101, 110, 116, 32, 97,
					115, 32, 97, 110, 32, 97, 116, 116, 97, 99, 104, 101, 100, 32, 46, 104, 116, 109, 108, 32, 102,
					105, 108, 101, 46},
			},
		},
	}

	testEmailFromFile(t, fp, expectedEmail)
}
func TestParseEmailKoreanMultipartMixedUtf8OverQuotedprintable(t *testing.T) {
	fp := "tests/test_korean_multipart_mixed_utf-8_over_quoted-printable.txt"
	tz, _ := time.LoadLocation("Asia/Seoul")
	expectedDate, _ := time.Parse(
		time.RFC1123Z+" (MST)",
		time.Date(2019, time.April, 1, 7, 55, 0, 0, tz).Format(time.RFC1123Z+" (MST)"))
	expectedEmail := Email{
		Headers: Headers{
			Date:    expectedDate,
			Subject: "📧 Test 한국어 팬그램",
			ReplyTo: []*mail.Address{
				{
					Name:    "앨리스 보내는사람",
					Address: "alice.bonaeneunsalam@example.net",
				},
			},
			Sender: &mail.Address{
				Name:    "앨리스 보내는사람",
				Address: "alice.bonaeneunsalam@example.com",
			},
			From: []*mail.Address{
				{
					Name:    "앨리스 보내는사람",
					Address: "alice.bonaeneunsalam@example.com",
				},
				{
					Name:    "앨리스 보내는사람",
					Address: "alice.bonaeneunsalam@example.net",
				},
			},
			To: []*mail.Address{
				{
					Name:    "밥 수신자",
					Address: "bob.susinja@example.com",
				},
				{
					Name:    "캐롤 수신자",
					Address: "carol.susinja@example.com",
				},
			},
			Cc: []*mail.Address{
				{
					Name:    "댄 수신자",
					Address: "dan.susinja@example.com",
				},
				{
					Name:    "이브 수신자",
					Address: "eve.susinja@example.com",
				},
			},
			Bcc: []*mail.Address{
				{
					Name:    "프랭크 수신자",
					Address: "frank.susinja@example.com",
				},
				{
					Name:    "그레이스 수신자",
					Address: "grace.susinja@example.com",
				},
			},
			MessageID:  "Message-Id-1@example.com",
			InReplyTo:  []MessageId{"Message-Id-0@example.com"},
			References: []MessageId{"Message-Id-0@example.com"},
			Comments:   "Message Header Comment",
			Keywords:   []string{"Keyword 1", "Keyword 2"},
			ResentDate: expectedDate,
			ResentFrom: []*mail.Address{
				{
					Name:    "앨리스 보내는사람",
					Address: "alice.bonaeneunsalam@example.net",
				},
				{
					Name:    "앨리스 보내는사람",
					Address: "alice.bonaeneunsalam@example.com",
				},
			},
			ResentSender: &mail.Address{
				Name:    "앨리스 보내는사람",
				Address: "alice.bonaeneunsalam@example.net",
			},
			ResentTo: []*mail.Address{
				{
					Name:    "밥 수신자",
					Address: "bob.susinja@example.net",
				},
				{
					Name:    "캐롤 수신자",
					Address: "carol.susinja@example.net",
				},
			},
			ResentCc: []*mail.Address{
				{
					Name:    "댄 수신자",
					Address: "dan.susinja@example.net",
				},
				{
					Name:    "이브 수신자",
					Address: "eve.susinja@example.net",
				},
			},
			ResentBcc: []*mail.Address{
				{
					Name:    "프랭크 수신자",
					Address: "frank.susinja@example.net",
				},
				{
					Name:    "그레이스 수신자",
					Address: "grace.susinja@example.net",
				},
			},
			ResentMessageID: "Message-Id-1@example.net",
			ContentType: ContentTypeHeader{
				ContentType: "multipart/mixed",
				Params: map[string]string{
					"boundary": "MixedBoundaryString",
					"charset":  "utf-8",
				},
			},
			ExtraHeaders: map[string][]string{
				"X-Clacks-Overhead": {"GNU Terry Pratchett"},
				"X-Script/function/\t !\"#$%&'()*+,-./;<=>?@[\\]^_`{|}~": {
					"TEST VALUE 1\t !\"#$%&'()*+,-./0123456789:;<=>?@ABCDEFGHIJKLMNOPQRSTUVWXYZ[\\]^_` abcdefghijklmnopqrstuvwxyz{|}~",
					"TEST VALUE 2\t !\"#$%&'()*+,-./0123456789:;<=>?@ABCDEFGHIJKLMNOPQRSTUVWXYZ[\\]^_` abcdefghijklmnopqrstuvwxyz{|}~",
				},
			},
		},
		Text:         `키스의 고유조건은 입술끼리 만나야 하고 특별한 기술은 필요치 않다.`,
		EnrichedText: `<bold>키스의</bold> <italic>고유조건은</italic> <fixed>입술끼리</fixed> <underline>만나야</underline> 하고 특별한 기술은 필요치 않다.`,
		HTML: `<html>
<div dir="ltr">
<p>키스의 고유조건은 입술끼리 만나야 하고 특별한 기술은 필요치 않다.</p>
</div>
</html>`,
		InlineFiles: []InlineFile{
			{
				ContentType: ContentTypeHeader{
					ContentType: "image/jpeg",
					Params: map[string]string{
						"name": "inline-jpg-image-without-disposition.jpg",
					},
				},
				ContentDisposition: ContentDispositionHeader{
					ContentDisposition: "",
					Params:             map[string]string(nil),
				},
				Data: []byte{255, 216, 255, 219, 0, 67, 0, 3, 2, 2, 2, 2, 2, 3, 2, 2, 2, 3, 3, 3, 3, 4,
					6, 4, 4, 4, 4, 4, 8, 6, 6, 5, 6, 9, 8, 10, 10, 9, 8, 9, 9, 10, 12, 15, 12, 10, 11, 14, 11, 9, 9,
					13, 17, 13, 14, 15, 16, 16, 17, 16, 10, 12, 18, 19, 18, 16, 19, 15, 16, 16, 16, 255, 201, 0, 11, 8,
					0, 1, 0, 1, 1, 1, 17, 0, 255, 204, 0, 6, 0, 16, 16, 5, 255, 218, 0, 8, 1, 1, 0, 0, 63, 0, 210, 207,
					32, 255, 217},
			},
			{
				ContentID: "inline-jpg-image.jpg@example.com",
				ContentType: ContentTypeHeader{
					ContentType: "image/jpeg",
					Params: map[string]string{
						"name": "inline-jpg-image-name.jpg",
					},
				},
				ContentDisposition: ContentDispositionHeader{
					ContentDisposition: inline,
					Params: map[string]string{
						"filename": "inline-jpg-image-filename.jpg",
					},
				},
				Data: []byte{255, 216, 255, 219, 0, 67, 0, 3, 2, 2, 2, 2, 2, 3, 2, 2, 2, 3, 3, 3, 3, 4,
					6, 4, 4, 4, 4, 4, 8, 6, 6, 5, 6, 9, 8, 10, 10, 9, 8, 9, 9, 10, 12, 15, 12, 10, 11, 14, 11, 9, 9,
					13, 17, 13, 14, 15, 16, 16, 17, 16, 10, 12, 18, 19, 18, 16, 19, 15, 16, 16, 16, 255, 201, 0, 11, 8,
					0, 1, 0, 1, 1, 1, 17, 0, 255, 204, 0, 6, 0, 16, 16, 5, 255, 218, 0, 8, 1, 1, 0, 0, 63, 0, 210, 207,
					32, 255, 217},
			},
		},
		AttachedFiles: []AttachedFile{
			{
				ContentType: ContentTypeHeader{
					ContentType: "application/pdf",
					Params: map[string]string{
						"name": "attached-pdf-name.pdf",
					},
				},
				ContentDisposition: ContentDispositionHeader{
					ContentDisposition: attachment,
					Params: map[string]string{
						"filename": "attached-pdf-filename.pdf",
					},
				},
				Data: []byte{37, 80, 68, 70, 45, 49, 46, 13, 116, 114, 97, 105, 108, 101, 114, 60, 60,
					47, 82, 111, 111, 116, 60, 60, 47, 80, 97, 103, 101, 115, 60, 60, 47, 75, 105, 100, 115, 91, 60,
					60, 47, 77, 101, 100, 105, 97, 66, 111, 120, 91, 48, 32, 48, 32, 51, 32, 51, 93, 62, 62, 93, 62,
					62, 62, 62, 62, 62},
			},
			{
				ContentType: ContentTypeHeader{
					ContentType: "application/pdf",
					Params: map[string]string{
						"name": "attached-pdf-without-disposition.pdf",
					},
				},
				ContentDisposition: ContentDispositionHeader{
					ContentDisposition: "",
					Params:             map[string]string(nil),
				},
				Data: []byte{37, 80, 68, 70, 45, 49, 46, 13, 116, 114, 97, 105, 108, 101, 114, 60, 60,
					47, 82, 111, 111, 116, 60, 60, 47, 80, 97, 103, 101, 115, 60, 60, 47, 75, 105, 100, 115, 91, 60,
					60, 47, 77, 101, 100, 105, 97, 66, 111, 120, 91, 48, 32, 48, 32, 51, 32, 51, 93, 62, 62, 93, 62,
					62, 62, 62, 62, 62},
			},
			{
				ContentType: ContentTypeHeader{
					ContentType: "application/json",
					Params: map[string]string{
						"name": "attached-json-name.json",
					},
				},
				ContentDisposition: ContentDispositionHeader{
					ContentDisposition: attachment,
					Params: map[string]string{
						"filename": "attached-json-filename.json",
					},
				},
				Data: []byte{123, 34, 102, 111, 111, 34, 58, 34, 98, 97, 114, 34, 125},
			},
			{
				ContentType: ContentTypeHeader{
					ContentType: "text/plain",
					Params: map[string]string{
						"name": "attached-text-plain-name.txt",
					},
				},
				ContentDisposition: ContentDispositionHeader{
					ContentDisposition: attachment,
					Params: map[string]string{
						"filename": "attached-text-plain-filename.txt",
					},
				},
				Data: []byte{84, 101, 120, 116, 47, 112, 108, 97, 105, 110, 32, 99, 111, 110, 116, 101, 110, 116, 32,
					97, 115, 32, 97, 110, 32, 97, 116, 116, 97, 99, 104, 101, 100, 32, 46, 116, 120, 116, 32, 102, 105,
					108, 101, 46},
			},
			{
				ContentType: ContentTypeHeader{
					ContentType: "text/html",
					Params: map[string]string{
						"name": "attached-text-html-name.html",
					},
				},
				ContentDisposition: ContentDispositionHeader{
					ContentDisposition: attachment,
					Params: map[string]string{
						"filename": "attached-text-html-filename.html",
					},
				},
				Data: []byte{84, 101, 120, 116, 47, 104, 116, 109, 108, 32, 99, 111, 110, 116, 101, 110, 116, 32, 97,
					115, 32, 97, 110, 32, 97, 116, 116, 97, 99, 104, 101, 100, 32, 46, 104, 116, 109, 108, 32, 102,
					105, 108, 101, 46},
			},
		},
	}

	testEmailFromFile(t, fp, expectedEmail)
}
func TestParseEmailKoreanMultipartMixedEuckrOverBase64(t *testing.T) {
	fp := "tests/test_korean_multipart_mixed_euc-kr_over_base64.txt"
	tz, _ := time.LoadLocation("Asia/Seoul")
	expectedDate, _ := time.Parse(
		time.RFC1123Z+" (MST)",
		time.Date(2019, time.April, 1, 7, 55, 0, 0, tz).Format(time.RFC1123Z+" (MST)"))
	expectedEmail := Email{
		Headers: Headers{
			Date:    expectedDate,
			Subject: "📧 Test 한국어 팬그램",
			ReplyTo: []*mail.Address{
				{
					Name:    "앨리스 보내는사람",
					Address: "alice.bonaeneunsalam@example.net",
				},
			},
			Sender: &mail.Address{
				Name:    "앨리스 보내는사람",
				Address: "alice.bonaeneunsalam@example.com",
			},
			From: []*mail.Address{
				{
					Name:    "앨리스 보내는사람",
					Address: "alice.bonaeneunsalam@example.com",
				},
				{
					Name:    "앨리스 보내는사람",
					Address: "alice.bonaeneunsalam@example.net",
				},
			},
			To: []*mail.Address{
				{
					Name:    "밥 수신자",
					Address: "bob.susinja@example.com",
				},
				{
					Name:    "캐롤 수신자",
					Address: "carol.susinja@example.com",
				},
			},
			Cc: []*mail.Address{
				{
					Name:    "댄 수신자",
					Address: "dan.susinja@example.com",
				},
				{
					Name:    "이브 수신자",
					Address: "eve.susinja@example.com",
				},
			},
			Bcc: []*mail.Address{
				{
					Name:    "프랭크 수신자",
					Address: "frank.susinja@example.com",
				},
				{
					Name:    "그레이스 수신자",
					Address: "grace.susinja@example.com",
				},
			},
			MessageID:  "Message-Id-1@example.com",
			InReplyTo:  []MessageId{"Message-Id-0@example.com"},
			References: []MessageId{"Message-Id-0@example.com"},
			Comments:   "Message Header Comment",
			Keywords:   []string{"Keyword 1", "Keyword 2"},
			ResentDate: expectedDate,
			ResentFrom: []*mail.Address{
				{
					Name:    "앨리스 보내는사람",
					Address: "alice.bonaeneunsalam@example.net",
				},
				{
					Name:    "앨리스 보내는사람",
					Address: "alice.bonaeneunsalam@example.com",
				},
			},
			ResentSender: &mail.Address{
				Name:    "앨리스 보내는사람",
				Address: "alice.bonaeneunsalam@example.net",
			},
			ResentTo: []*mail.Address{
				{
					Name:    "밥 수신자",
					Address: "bob.susinja@example.net",
				},
				{
					Name:    "캐롤 수신자",
					Address: "carol.susinja@example.net",
				},
			},
			ResentCc: []*mail.Address{
				{
					Name:    "댄 수신자",
					Address: "dan.susinja@example.net",
				},
				{
					Name:    "이브 수신자",
					Address: "eve.susinja@example.net",
				},
			},
			ResentBcc: []*mail.Address{
				{
					Name:    "프랭크 수신자",
					Address: "frank.susinja@example.net",
				},
				{
					Name:    "그레이스 수신자",
					Address: "grace.susinja@example.net",
				},
			},
			ResentMessageID: "Message-Id-1@example.net",
			ContentType: ContentTypeHeader{
				ContentType: "multipart/mixed",
				Params: map[string]string{
					"boundary": "MixedBoundaryString",
					"charset":  "euc-kr",
				},
			},
			ExtraHeaders: map[string][]string{
				"X-Clacks-Overhead": {"GNU Terry Pratchett"},
				"X-Script/function/\t !\"#$%&'()*+,-./;<=>?@[\\]^_`{|}~": {
					"TEST VALUE 1\t !\"#$%&'()*+,-./0123456789:;<=>?@ABCDEFGHIJKLMNOPQRSTUVWXYZ[\\]^_` abcdefghijklmnopqrstuvwxyz{|}~",
					"TEST VALUE 2\t !\"#$%&'()*+,-./0123456789:;<=>?@ABCDEFGHIJKLMNOPQRSTUVWXYZ[\\]^_` abcdefghijklmnopqrstuvwxyz{|}~",
				},
			},
		},
		Text:         `키스의 고유조건은 입술끼리 만나야 하고 특별한 기술은 필요치 않다.`,
		EnrichedText: `<bold>키스의</bold> <italic>고유조건은</italic> <fixed>입술끼리</fixed> <underline>만나야</underline> 하고 특별한 기술은 필요치 않다.`,
		HTML: `<html>
<div dir="ltr">
<p>키스의 고유조건은 입술끼리 만나야 하고 특별한 기술은 필요치 않다.</p>
</div>
</html>`,
		InlineFiles: []InlineFile{
			{
				ContentType: ContentTypeHeader{
					ContentType: "image/jpeg",
					Params: map[string]string{
						"name": "inline-jpg-image-without-disposition.jpg",
					},
				},
				ContentDisposition: ContentDispositionHeader{
					ContentDisposition: "",
					Params:             map[string]string(nil),
				},
				Data: []byte{255, 216, 255, 219, 0, 67, 0, 3, 2, 2, 2, 2, 2, 3, 2, 2, 2, 3, 3, 3, 3, 4,
					6, 4, 4, 4, 4, 4, 8, 6, 6, 5, 6, 9, 8, 10, 10, 9, 8, 9, 9, 10, 12, 15, 12, 10, 11, 14, 11, 9, 9,
					13, 17, 13, 14, 15, 16, 16, 17, 16, 10, 12, 18, 19, 18, 16, 19, 15, 16, 16, 16, 255, 201, 0, 11, 8,
					0, 1, 0, 1, 1, 1, 17, 0, 255, 204, 0, 6, 0, 16, 16, 5, 255, 218, 0, 8, 1, 1, 0, 0, 63, 0, 210, 207,
					32, 255, 217},
			},
			{
				ContentID: "inline-jpg-image.jpg@example.com",
				ContentType: ContentTypeHeader{
					ContentType: "image/jpeg",
					Params: map[string]string{
						"name": "inline-jpg-image-name.jpg",
					},
				},
				ContentDisposition: ContentDispositionHeader{
					ContentDisposition: inline,
					Params: map[string]string{
						"filename": "inline-jpg-image-filename.jpg",
					},
				},
				Data: []byte{255, 216, 255, 219, 0, 67, 0, 3, 2, 2, 2, 2, 2, 3, 2, 2, 2, 3, 3, 3, 3, 4,
					6, 4, 4, 4, 4, 4, 8, 6, 6, 5, 6, 9, 8, 10, 10, 9, 8, 9, 9, 10, 12, 15, 12, 10, 11, 14, 11, 9, 9,
					13, 17, 13, 14, 15, 16, 16, 17, 16, 10, 12, 18, 19, 18, 16, 19, 15, 16, 16, 16, 255, 201, 0, 11, 8,
					0, 1, 0, 1, 1, 1, 17, 0, 255, 204, 0, 6, 0, 16, 16, 5, 255, 218, 0, 8, 1, 1, 0, 0, 63, 0, 210, 207,
					32, 255, 217},
			},
		},
		AttachedFiles: []AttachedFile{
			{
				ContentType: ContentTypeHeader{
					ContentType: "application/pdf",
					Params: map[string]string{
						"name": "attached-pdf-name.pdf",
					},
				},
				ContentDisposition: ContentDispositionHeader{
					ContentDisposition: attachment,
					Params: map[string]string{
						"filename": "attached-pdf-filename.pdf",
					},
				},
				Data: []byte{37, 80, 68, 70, 45, 49, 46, 13, 116, 114, 97, 105, 108, 101, 114, 60, 60,
					47, 82, 111, 111, 116, 60, 60, 47, 80, 97, 103, 101, 115, 60, 60, 47, 75, 105, 100, 115, 91, 60,
					60, 47, 77, 101, 100, 105, 97, 66, 111, 120, 91, 48, 32, 48, 32, 51, 32, 51, 93, 62, 62, 93, 62,
					62, 62, 62, 62, 62},
			},
			{
				ContentType: ContentTypeHeader{
					ContentType: "application/pdf",
					Params: map[string]string{
						"name": "attached-pdf-without-disposition.pdf",
					},
				},
				ContentDisposition: ContentDispositionHeader{
					ContentDisposition: "",
					Params:             map[string]string(nil),
				},
				Data: []byte{37, 80, 68, 70, 45, 49, 46, 13, 116, 114, 97, 105, 108, 101, 114, 60, 60,
					47, 82, 111, 111, 116, 60, 60, 47, 80, 97, 103, 101, 115, 60, 60, 47, 75, 105, 100, 115, 91, 60,
					60, 47, 77, 101, 100, 105, 97, 66, 111, 120, 91, 48, 32, 48, 32, 51, 32, 51, 93, 62, 62, 93, 62,
					62, 62, 62, 62, 62},
			},
			{
				ContentType: ContentTypeHeader{
					ContentType: "application/json",
					Params: map[string]string{
						"name": "attached-json-name.json",
					},
				},
				ContentDisposition: ContentDispositionHeader{
					ContentDisposition: attachment,
					Params: map[string]string{
						"filename": "attached-json-filename.json",
					},
				},
				Data: []byte{123, 34, 102, 111, 111, 34, 58, 34, 98, 97, 114, 34, 125},
			},
			{
				ContentType: ContentTypeHeader{
					ContentType: "text/plain",
					Params: map[string]string{
						"name": "attached-text-plain-name.txt",
					},
				},
				ContentDisposition: ContentDispositionHeader{
					ContentDisposition: attachment,
					Params: map[string]string{
						"filename": "attached-text-plain-filename.txt",
					},
				},
				Data: []byte{84, 101, 120, 116, 47, 112, 108, 97, 105, 110, 32, 99, 111, 110, 116, 101, 110, 116, 32,
					97, 115, 32, 97, 110, 32, 97, 116, 116, 97, 99, 104, 101, 100, 32, 46, 116, 120, 116, 32, 102, 105,
					108, 101, 46},
			},
			{
				ContentType: ContentTypeHeader{
					ContentType: "text/html",
					Params: map[string]string{
						"name": "attached-text-html-name.html",
					},
				},
				ContentDisposition: ContentDispositionHeader{
					ContentDisposition: attachment,
					Params: map[string]string{
						"filename": "attached-text-html-filename.html",
					},
				},
				Data: []byte{84, 101, 120, 116, 47, 104, 116, 109, 108, 32, 99, 111, 110, 116, 101, 110, 116, 32, 97,
					115, 32, 97, 110, 32, 97, 116, 116, 97, 99, 104, 101, 100, 32, 46, 104, 116, 109, 108, 32, 102,
					105, 108, 101, 46},
			},
		},
	}

	testEmailFromFile(t, fp, expectedEmail)
}
func TestParseEmailKoreanMultipartMixedEuckrOverQuotedprintable(t *testing.T) {
	fp := "tests/test_korean_multipart_mixed_euc-kr_over_quoted-printable.txt"
	tz, _ := time.LoadLocation("Asia/Seoul")
	expectedDate, _ := time.Parse(
		time.RFC1123Z+" (MST)",
		time.Date(2019, time.April, 1, 7, 55, 0, 0, tz).Format(time.RFC1123Z+" (MST)"))
	expectedEmail := Email{
		Headers: Headers{
			Date:    expectedDate,
			Subject: "📧 Test 한국어 팬그램",
			ReplyTo: []*mail.Address{
				{
					Name:    "앨리스 보내는사람",
					Address: "alice.bonaeneunsalam@example.net",
				},
			},
			Sender: &mail.Address{
				Name:    "앨리스 보내는사람",
				Address: "alice.bonaeneunsalam@example.com",
			},
			From: []*mail.Address{
				{
					Name:    "앨리스 보내는사람",
					Address: "alice.bonaeneunsalam@example.com",
				},
				{
					Name:    "앨리스 보내는사람",
					Address: "alice.bonaeneunsalam@example.net",
				},
			},
			To: []*mail.Address{
				{
					Name:    "밥 수신자",
					Address: "bob.susinja@example.com",
				},
				{
					Name:    "캐롤 수신자",
					Address: "carol.susinja@example.com",
				},
			},
			Cc: []*mail.Address{
				{
					Name:    "댄 수신자",
					Address: "dan.susinja@example.com",
				},
				{
					Name:    "이브 수신자",
					Address: "eve.susinja@example.com",
				},
			},
			Bcc: []*mail.Address{
				{
					Name:    "프랭크 수신자",
					Address: "frank.susinja@example.com",
				},
				{
					Name:    "그레이스 수신자",
					Address: "grace.susinja@example.com",
				},
			},
			MessageID:  "Message-Id-1@example.com",
			InReplyTo:  []MessageId{"Message-Id-0@example.com"},
			References: []MessageId{"Message-Id-0@example.com"},
			Comments:   "Message Header Comment",
			Keywords:   []string{"Keyword 1", "Keyword 2"},
			ResentDate: expectedDate,
			ResentFrom: []*mail.Address{
				{
					Name:    "앨리스 보내는사람",
					Address: "alice.bonaeneunsalam@example.net",
				},
				{
					Name:    "앨리스 보내는사람",
					Address: "alice.bonaeneunsalam@example.com",
				},
			},
			ResentSender: &mail.Address{
				Name:    "앨리스 보내는사람",
				Address: "alice.bonaeneunsalam@example.net",
			},
			ResentTo: []*mail.Address{
				{
					Name:    "밥 수신자",
					Address: "bob.susinja@example.net",
				},
				{
					Name:    "캐롤 수신자",
					Address: "carol.susinja@example.net",
				},
			},
			ResentCc: []*mail.Address{
				{
					Name:    "댄 수신자",
					Address: "dan.susinja@example.net",
				},
				{
					Name:    "이브 수신자",
					Address: "eve.susinja@example.net",
				},
			},
			ResentBcc: []*mail.Address{
				{
					Name:    "프랭크 수신자",
					Address: "frank.susinja@example.net",
				},
				{
					Name:    "그레이스 수신자",
					Address: "grace.susinja@example.net",
				},
			},
			ResentMessageID: "Message-Id-1@example.net",
			ContentType: ContentTypeHeader{
				ContentType: "multipart/mixed",
				Params: map[string]string{
					"boundary": "MixedBoundaryString",
					"charset":  "euc-kr",
				},
			},
			ExtraHeaders: map[string][]string{
				"X-Clacks-Overhead": {"GNU Terry Pratchett"},
				"X-Script/function/\t !\"#$%&'()*+,-./;<=>?@[\\]^_`{|}~": {
					"TEST VALUE 1\t !\"#$%&'()*+,-./0123456789:;<=>?@ABCDEFGHIJKLMNOPQRSTUVWXYZ[\\]^_` abcdefghijklmnopqrstuvwxyz{|}~",
					"TEST VALUE 2\t !\"#$%&'()*+,-./0123456789:;<=>?@ABCDEFGHIJKLMNOPQRSTUVWXYZ[\\]^_` abcdefghijklmnopqrstuvwxyz{|}~",
				},
			},
		},
		Text:         `키스의 고유조건은 입술끼리 만나야 하고 특별한 기술은 필요치 않다.`,
		EnrichedText: `<bold>키스의</bold> <italic>고유조건은</italic> <fixed>입술끼리</fixed> <underline>만나야</underline> 하고 특별한 기술은 필요치 않다.`,
		HTML: `<html>
<div dir="ltr">
<p>키스의 고유조건은 입술끼리 만나야 하고 특별한 기술은 필요치 않다.</p>
</div>
</html>`,
		InlineFiles: []InlineFile{
			{
				ContentType: ContentTypeHeader{
					ContentType: "image/jpeg",
					Params: map[string]string{
						"name": "inline-jpg-image-without-disposition.jpg",
					},
				},
				ContentDisposition: ContentDispositionHeader{
					ContentDisposition: "",
					Params:             map[string]string(nil),
				},
				Data: []byte{255, 216, 255, 219, 0, 67, 0, 3, 2, 2, 2, 2, 2, 3, 2, 2, 2, 3, 3, 3, 3, 4,
					6, 4, 4, 4, 4, 4, 8, 6, 6, 5, 6, 9, 8, 10, 10, 9, 8, 9, 9, 10, 12, 15, 12, 10, 11, 14, 11, 9, 9,
					13, 17, 13, 14, 15, 16, 16, 17, 16, 10, 12, 18, 19, 18, 16, 19, 15, 16, 16, 16, 255, 201, 0, 11, 8,
					0, 1, 0, 1, 1, 1, 17, 0, 255, 204, 0, 6, 0, 16, 16, 5, 255, 218, 0, 8, 1, 1, 0, 0, 63, 0, 210, 207,
					32, 255, 217},
			},
			{
				ContentID: "inline-jpg-image.jpg@example.com",
				ContentType: ContentTypeHeader{
					ContentType: "image/jpeg",
					Params: map[string]string{
						"name": "inline-jpg-image-name.jpg",
					},
				},
				ContentDisposition: ContentDispositionHeader{
					ContentDisposition: inline,
					Params: map[string]string{
						"filename": "inline-jpg-image-filename.jpg",
					},
				},
				Data: []byte{255, 216, 255, 219, 0, 67, 0, 3, 2, 2, 2, 2, 2, 3, 2, 2, 2, 3, 3, 3, 3, 4,
					6, 4, 4, 4, 4, 4, 8, 6, 6, 5, 6, 9, 8, 10, 10, 9, 8, 9, 9, 10, 12, 15, 12, 10, 11, 14, 11, 9, 9,
					13, 17, 13, 14, 15, 16, 16, 17, 16, 10, 12, 18, 19, 18, 16, 19, 15, 16, 16, 16, 255, 201, 0, 11, 8,
					0, 1, 0, 1, 1, 1, 17, 0, 255, 204, 0, 6, 0, 16, 16, 5, 255, 218, 0, 8, 1, 1, 0, 0, 63, 0, 210, 207,
					32, 255, 217},
			},
		},
		AttachedFiles: []AttachedFile{
			{
				ContentType: ContentTypeHeader{
					ContentType: "application/pdf",
					Params: map[string]string{
						"name": "attached-pdf-name.pdf",
					},
				},
				ContentDisposition: ContentDispositionHeader{
					ContentDisposition: attachment,
					Params: map[string]string{
						"filename": "attached-pdf-filename.pdf",
					},
				},
				Data: []byte{37, 80, 68, 70, 45, 49, 46, 13, 116, 114, 97, 105, 108, 101, 114, 60, 60,
					47, 82, 111, 111, 116, 60, 60, 47, 80, 97, 103, 101, 115, 60, 60, 47, 75, 105, 100, 115, 91, 60,
					60, 47, 77, 101, 100, 105, 97, 66, 111, 120, 91, 48, 32, 48, 32, 51, 32, 51, 93, 62, 62, 93, 62,
					62, 62, 62, 62, 62},
			},
			{
				ContentType: ContentTypeHeader{
					ContentType: "application/pdf",
					Params: map[string]string{
						"name": "attached-pdf-without-disposition.pdf",
					},
				},
				ContentDisposition: ContentDispositionHeader{
					ContentDisposition: "",
					Params:             map[string]string(nil),
				},
				Data: []byte{37, 80, 68, 70, 45, 49, 46, 13, 116, 114, 97, 105, 108, 101, 114, 60, 60,
					47, 82, 111, 111, 116, 60, 60, 47, 80, 97, 103, 101, 115, 60, 60, 47, 75, 105, 100, 115, 91, 60,
					60, 47, 77, 101, 100, 105, 97, 66, 111, 120, 91, 48, 32, 48, 32, 51, 32, 51, 93, 62, 62, 93, 62,
					62, 62, 62, 62, 62},
			},
			{
				ContentType: ContentTypeHeader{
					ContentType: "application/json",
					Params: map[string]string{
						"name": "attached-json-name.json",
					},
				},
				ContentDisposition: ContentDispositionHeader{
					ContentDisposition: attachment,
					Params: map[string]string{
						"filename": "attached-json-filename.json",
					},
				},
				Data: []byte{123, 34, 102, 111, 111, 34, 58, 34, 98, 97, 114, 34, 125},
			},
			{
				ContentType: ContentTypeHeader{
					ContentType: "text/plain",
					Params: map[string]string{
						"name": "attached-text-plain-name.txt",
					},
				},
				ContentDisposition: ContentDispositionHeader{
					ContentDisposition: attachment,
					Params: map[string]string{
						"filename": "attached-text-plain-filename.txt",
					},
				},
				Data: []byte{84, 101, 120, 116, 47, 112, 108, 97, 105, 110, 32, 99, 111, 110, 116, 101, 110, 116, 32,
					97, 115, 32, 97, 110, 32, 97, 116, 116, 97, 99, 104, 101, 100, 32, 46, 116, 120, 116, 32, 102, 105,
					108, 101, 46},
			},
			{
				ContentType: ContentTypeHeader{
					ContentType: "text/html",
					Params: map[string]string{
						"name": "attached-text-html-name.html",
					},
				},
				ContentDisposition: ContentDispositionHeader{
					ContentDisposition: attachment,
					Params: map[string]string{
						"filename": "attached-text-html-filename.html",
					},
				},
				Data: []byte{84, 101, 120, 116, 47, 104, 116, 109, 108, 32, 99, 111, 110, 116, 101, 110, 116, 32, 97,
					115, 32, 97, 110, 32, 97, 116, 116, 97, 99, 104, 101, 100, 32, 46, 104, 116, 109, 108, 32, 102,
					105, 108, 101, 46},
			},
		},
	}

	testEmailFromFile(t, fp, expectedEmail)
}
func TestParseEmailKoreanMultipartSignedUtf8OverBase64(t *testing.T) {
	fp := "tests/test_korean_multipart_signed_utf-8_over_base64.txt"
	tz, _ := time.LoadLocation("Asia/Seoul")
	expectedDate, _ := time.Parse(
		time.RFC1123Z+" (MST)",
		time.Date(2019, time.April, 1, 7, 55, 0, 0, tz).Format(time.RFC1123Z+" (MST)"))
	expectedEmail := Email{
		Headers: Headers{
			Date:    expectedDate,
			Subject: "📧 Signed Test 한국어 팬그램",
			ReplyTo: []*mail.Address{
				{
					Name:    "앨리스 보내는사람",
					Address: "alice.bonaeneunsalam@example.net",
				},
			},
			Sender: &mail.Address{
				Name:    "앨리스 보내는사람",
				Address: "alice.bonaeneunsalam@example.com",
			},
			From: []*mail.Address{
				{
					Name:    "앨리스 보내는사람",
					Address: "alice.bonaeneunsalam@example.com",
				},
				{
					Name:    "앨리스 보내는사람",
					Address: "alice.bonaeneunsalam@example.net",
				},
			},
			To: []*mail.Address{
				{
					Name:    "밥 수신자",
					Address: "bob.susinja@example.com",
				},
				{
					Name:    "캐롤 수신자",
					Address: "carol.susinja@example.com",
				},
			},
			Cc: []*mail.Address{
				{
					Name:    "댄 수신자",
					Address: "dan.susinja@example.com",
				},
				{
					Name:    "이브 수신자",
					Address: "eve.susinja@example.com",
				},
			},
			Bcc: []*mail.Address{
				{
					Name:    "프랭크 수신자",
					Address: "frank.susinja@example.com",
				},
				{
					Name:    "그레이스 수신자",
					Address: "grace.susinja@example.com",
				},
			},
			MessageID:  "Message-Id-1@example.com",
			InReplyTo:  []MessageId{"Message-Id-0@example.com"},
			References: []MessageId{"Message-Id-0@example.com"},
			Comments:   "Message Header Comment",
			Keywords:   []string{"Keyword 1", "Keyword 2"},
			ResentDate: expectedDate,
			ResentFrom: []*mail.Address{
				{
					Name:    "앨리스 보내는사람",
					Address: "alice.bonaeneunsalam@example.net",
				},
				{
					Name:    "앨리스 보내는사람",
					Address: "alice.bonaeneunsalam@example.com",
				},
			},
			ResentSender: &mail.Address{
				Name:    "앨리스 보내는사람",
				Address: "alice.bonaeneunsalam@example.net",
			},
			ResentTo: []*mail.Address{
				{
					Name:    "밥 수신자",
					Address: "bob.susinja@example.net",
				},
				{
					Name:    "캐롤 수신자",
					Address: "carol.susinja@example.net",
				},
			},
			ResentCc: []*mail.Address{
				{
					Name:    "댄 수신자",
					Address: "dan.susinja@example.net",
				},
				{
					Name:    "이브 수신자",
					Address: "eve.susinja@example.net",
				},
			},
			ResentBcc: []*mail.Address{
				{
					Name:    "프랭크 수신자",
					Address: "frank.susinja@example.net",
				},
				{
					Name:    "그레이스 수신자",
					Address: "grace.susinja@example.net",
				},
			},
			ResentMessageID: "Message-Id-1@example.net",
			ContentType: ContentTypeHeader{
				ContentType: "multipart/signed",
				Params: map[string]string{
					"boundary": "SignedBoundaryString",
					"charset":  "utf-8",
					"micalg":   "sha1",
					"protocol": "application/pkcs7-signature",
				},
			},
			ExtraHeaders: map[string][]string{
				"X-Clacks-Overhead": {"GNU Terry Pratchett"},
				"X-Script/function/\t !\"#$%&'()*+,-./;<=>?@[\\]^_`{|}~": {
					"TEST VALUE 1\t !\"#$%&'()*+,-./0123456789:;<=>?@ABCDEFGHIJKLMNOPQRSTUVWXYZ[\\]^_` abcdefghijklmnopqrstuvwxyz{|}~",
					"TEST VALUE 2\t !\"#$%&'()*+,-./0123456789:;<=>?@ABCDEFGHIJKLMNOPQRSTUVWXYZ[\\]^_` abcdefghijklmnopqrstuvwxyz{|}~",
				},
			},
		},
		Text: `키스의 고유조건은 입술끼리 만나야 하고 특별한 기술은 필요치 않다.`,
		AttachedFiles: []AttachedFile{
			{
				ContentType: ContentTypeHeader{
					ContentType: "application/pkcs7-signature",
					Params: map[string]string{
						"name": "smime.p7s",
					},
				},
				ContentDisposition: ContentDispositionHeader{
					ContentDisposition: attachment,
					Params: map[string]string{
						"filename": "smime.p7s",
					},
				},
				Data: []byte{130, 28, 135, 132, 117, 46, 142, 18, 97, 140, 126, 251, 159, 193, 199, 25, 58, 223, 189, 185, 227, 239, 158, 173, 108, 31, 71, 27, 133, 80, 165, 252, 133, 227, 174, 198, 132, 129, 159, 29, 246, 19, 235, 133, 80, 165, 252, 133, 227, 174, 198, 132, 129, 159, 29, 246, 19, 234, 49, 251, 238, 127, 7, 28, 104, 33, 200, 120, 71, 82, 232, 225, 38, 30, 249, 234, 214, 193, 244, 113, 147, 173, 251, 219, 158, 57, 252, 28, 113, 147, 173, 251, 225, 38, 24, 199, 239, 190, 173, 108, 31, 71, 27, 133, 80, 110, 120, 251, 231, 174, 198, 132, 129, 159, 29, 246, 19, 234, 8, 114, 30, 17, 212, 186, 58, 95, 200, 94, 59, 26, 18, 6, 124, 119, 216, 79, 174, 21, 65, 185, 227, 239, 158},
			},
		},
	}

	testEmailFromFile(t, fp, expectedEmail)
}

func TestParseEmailKoreanMultipartSignedUtf8OverQuotedprintable(t *testing.T) {
	fp := "tests/test_korean_multipart_signed_utf-8_over_quoted-printable.txt"
	tz, _ := time.LoadLocation("Asia/Seoul")
	expectedDate, _ := time.Parse(
		time.RFC1123Z+" (MST)",
		time.Date(2019, time.April, 1, 7, 55, 0, 0, tz).Format(time.RFC1123Z+" (MST)"))
	expectedEmail := Email{
		Headers: Headers{
			Date:    expectedDate,
			Subject: "📧 Signed Test 한국어 팬그램",
			ReplyTo: []*mail.Address{
				{
					Name:    "앨리스 보내는사람",
					Address: "alice.bonaeneunsalam@example.net",
				},
			},
			Sender: &mail.Address{
				Name:    "앨리스 보내는사람",
				Address: "alice.bonaeneunsalam@example.com",
			},
			From: []*mail.Address{
				{
					Name:    "앨리스 보내는사람",
					Address: "alice.bonaeneunsalam@example.com",
				},
				{
					Name:    "앨리스 보내는사람",
					Address: "alice.bonaeneunsalam@example.net",
				},
			},
			To: []*mail.Address{
				{
					Name:    "밥 수신자",
					Address: "bob.susinja@example.com",
				},
				{
					Name:    "캐롤 수신자",
					Address: "carol.susinja@example.com",
				},
			},
			Cc: []*mail.Address{
				{
					Name:    "댄 수신자",
					Address: "dan.susinja@example.com",
				},
				{
					Name:    "이브 수신자",
					Address: "eve.susinja@example.com",
				},
			},
			Bcc: []*mail.Address{
				{
					Name:    "프랭크 수신자",
					Address: "frank.susinja@example.com",
				},
				{
					Name:    "그레이스 수신자",
					Address: "grace.susinja@example.com",
				},
			},
			MessageID:  "Message-Id-1@example.com",
			InReplyTo:  []MessageId{"Message-Id-0@example.com"},
			References: []MessageId{"Message-Id-0@example.com"},
			Comments:   "Message Header Comment",
			Keywords:   []string{"Keyword 1", "Keyword 2"},
			ResentDate: expectedDate,
			ResentFrom: []*mail.Address{
				{
					Name:    "앨리스 보내는사람",
					Address: "alice.bonaeneunsalam@example.net",
				},
				{
					Name:    "앨리스 보내는사람",
					Address: "alice.bonaeneunsalam@example.com",
				},
			},
			ResentSender: &mail.Address{
				Name:    "앨리스 보내는사람",
				Address: "alice.bonaeneunsalam@example.net",
			},
			ResentTo: []*mail.Address{
				{
					Name:    "밥 수신자",
					Address: "bob.susinja@example.net",
				},
				{
					Name:    "캐롤 수신자",
					Address: "carol.susinja@example.net",
				},
			},
			ResentCc: []*mail.Address{
				{
					Name:    "댄 수신자",
					Address: "dan.susinja@example.net",
				},
				{
					Name:    "이브 수신자",
					Address: "eve.susinja@example.net",
				},
			},
			ResentBcc: []*mail.Address{
				{
					Name:    "프랭크 수신자",
					Address: "frank.susinja@example.net",
				},
				{
					Name:    "그레이스 수신자",
					Address: "grace.susinja@example.net",
				},
			},
			ResentMessageID: "Message-Id-1@example.net",
			ContentType: ContentTypeHeader{
				ContentType: "multipart/signed",
				Params: map[string]string{
					"boundary": "SignedBoundaryString",
					"charset":  "utf-8",
					"micalg":   "sha1",
					"protocol": "application/pkcs7-signature",
				},
			},
			ExtraHeaders: map[string][]string{
				"X-Clacks-Overhead": {"GNU Terry Pratchett"},
				"X-Script/function/\t !\"#$%&'()*+,-./;<=>?@[\\]^_`{|}~": {
					"TEST VALUE 1\t !\"#$%&'()*+,-./0123456789:;<=>?@ABCDEFGHIJKLMNOPQRSTUVWXYZ[\\]^_` abcdefghijklmnopqrstuvwxyz{|}~",
					"TEST VALUE 2\t !\"#$%&'()*+,-./0123456789:;<=>?@ABCDEFGHIJKLMNOPQRSTUVWXYZ[\\]^_` abcdefghijklmnopqrstuvwxyz{|}~",
				},
			},
		},
		Text: `키스의 고유조건은 입술끼리 만나야 하고 특별한 기술은 필요치 않다.`,
		AttachedFiles: []AttachedFile{
			{
				ContentType: ContentTypeHeader{
					ContentType: "application/pkcs7-signature",
					Params: map[string]string{
						"name": "smime.p7s",
					},
				},
				ContentDisposition: ContentDispositionHeader{
					ContentDisposition: attachment,
					Params: map[string]string{
						"filename": "smime.p7s",
					},
				},
				Data: []byte{130, 28, 135, 132, 117, 46, 142, 18, 97, 140, 126, 251, 159, 193, 199, 25, 58, 223, 189, 185, 227, 239, 158, 173, 108, 31, 71, 27, 133, 80, 165, 252, 133, 227, 174, 198, 132, 129, 159, 29, 246, 19, 235, 133, 80, 165, 252, 133, 227, 174, 198, 132, 129, 159, 29, 246, 19, 234, 49, 251, 238, 127, 7, 28, 104, 33, 200, 120, 71, 82, 232, 225, 38, 30, 249, 234, 214, 193, 244, 113, 147, 173, 251, 219, 158, 57, 252, 28, 113, 147, 173, 251, 225, 38, 24, 199, 239, 190, 173, 108, 31, 71, 27, 133, 80, 110, 120, 251, 231, 174, 198, 132, 129, 159, 29, 246, 19, 234, 8, 114, 30, 17, 212, 186, 58, 95, 200, 94, 59, 26, 18, 6, 124, 119, 216, 79, 174, 21, 65, 185, 227, 239, 158},
			},
		},
	}

	testEmailFromFile(t, fp, expectedEmail)
}

func TestParseEmailKoreanMultipartSignedEuckrOverBase64(t *testing.T) {
	fp := "tests/test_korean_multipart_signed_euc-kr_over_base64.txt"
	tz, _ := time.LoadLocation("Asia/Seoul")
	expectedDate, _ := time.Parse(
		time.RFC1123Z+" (MST)",
		time.Date(2019, time.April, 1, 7, 55, 0, 0, tz).Format(time.RFC1123Z+" (MST)"))
	expectedEmail := Email{
		Headers: Headers{
			Date:    expectedDate,
			Subject: "📧 Signed Test 한국어 팬그램",
			ReplyTo: []*mail.Address{
				{
					Name:    "앨리스 보내는사람",
					Address: "alice.bonaeneunsalam@example.net",
				},
			},
			Sender: &mail.Address{
				Name:    "앨리스 보내는사람",
				Address: "alice.bonaeneunsalam@example.com",
			},
			From: []*mail.Address{
				{
					Name:    "앨리스 보내는사람",
					Address: "alice.bonaeneunsalam@example.com",
				},
				{
					Name:    "앨리스 보내는사람",
					Address: "alice.bonaeneunsalam@example.net",
				},
			},
			To: []*mail.Address{
				{
					Name:    "밥 수신자",
					Address: "bob.susinja@example.com",
				},
				{
					Name:    "캐롤 수신자",
					Address: "carol.susinja@example.com",
				},
			},
			Cc: []*mail.Address{
				{
					Name:    "댄 수신자",
					Address: "dan.susinja@example.com",
				},
				{
					Name:    "이브 수신자",
					Address: "eve.susinja@example.com",
				},
			},
			Bcc: []*mail.Address{
				{
					Name:    "프랭크 수신자",
					Address: "frank.susinja@example.com",
				},
				{
					Name:    "그레이스 수신자",
					Address: "grace.susinja@example.com",
				},
			},
			MessageID:  "Message-Id-1@example.com",
			InReplyTo:  []MessageId{"Message-Id-0@example.com"},
			References: []MessageId{"Message-Id-0@example.com"},
			Comments:   "Message Header Comment",
			Keywords:   []string{"Keyword 1", "Keyword 2"},
			ResentDate: expectedDate,
			ResentFrom: []*mail.Address{
				{
					Name:    "앨리스 보내는사람",
					Address: "alice.bonaeneunsalam@example.net",
				},
				{
					Name:    "앨리스 보내는사람",
					Address: "alice.bonaeneunsalam@example.com",
				},
			},
			ResentSender: &mail.Address{
				Name:    "앨리스 보내는사람",
				Address: "alice.bonaeneunsalam@example.net",
			},
			ResentTo: []*mail.Address{
				{
					Name:    "밥 수신자",
					Address: "bob.susinja@example.net",
				},
				{
					Name:    "캐롤 수신자",
					Address: "carol.susinja@example.net",
				},
			},
			ResentCc: []*mail.Address{
				{
					Name:    "댄 수신자",
					Address: "dan.susinja@example.net",
				},
				{
					Name:    "이브 수신자",
					Address: "eve.susinja@example.net",
				},
			},
			ResentBcc: []*mail.Address{
				{
					Name:    "프랭크 수신자",
					Address: "frank.susinja@example.net",
				},
				{
					Name:    "그레이스 수신자",
					Address: "grace.susinja@example.net",
				},
			},
			ResentMessageID: "Message-Id-1@example.net",
			ContentType: ContentTypeHeader{
				ContentType: "multipart/signed",
				Params: map[string]string{
					"boundary": "SignedBoundaryString",
					"charset":  "euc-kr",
					"micalg":   "sha1",
					"protocol": "application/pkcs7-signature",
				},
			},
			ExtraHeaders: map[string][]string{
				"X-Clacks-Overhead": {"GNU Terry Pratchett"},
				"X-Script/function/\t !\"#$%&'()*+,-./;<=>?@[\\]^_`{|}~": {
					"TEST VALUE 1\t !\"#$%&'()*+,-./0123456789:;<=>?@ABCDEFGHIJKLMNOPQRSTUVWXYZ[\\]^_` abcdefghijklmnopqrstuvwxyz{|}~",
					"TEST VALUE 2\t !\"#$%&'()*+,-./0123456789:;<=>?@ABCDEFGHIJKLMNOPQRSTUVWXYZ[\\]^_` abcdefghijklmnopqrstuvwxyz{|}~",
				},
			},
		},
		Text: `키스의 고유조건은 입술끼리 만나야 하고 특별한 기술은 필요치 않다.`,
		AttachedFiles: []AttachedFile{
			{
				ContentType: ContentTypeHeader{
					ContentType: "application/pkcs7-signature",
					Params: map[string]string{
						"name": "smime.p7s",
					},
				},
				ContentDisposition: ContentDispositionHeader{
					ContentDisposition: attachment,
					Params: map[string]string{
						"filename": "smime.p7s",
					},
				},
				Data: []byte{130, 28, 135, 132, 117, 46, 142, 18, 97, 140, 126, 251, 159, 193, 199, 25, 58, 223, 189, 185, 227, 239, 158, 173, 108, 31, 71, 27, 133, 80, 165, 252, 133, 227, 174, 198, 132, 129, 159, 29, 246, 19, 235, 133, 80, 165, 252, 133, 227, 174, 198, 132, 129, 159, 29, 246, 19, 234, 49, 251, 238, 127, 7, 28, 104, 33, 200, 120, 71, 82, 232, 225, 38, 30, 249, 234, 214, 193, 244, 113, 147, 173, 251, 219, 158, 57, 252, 28, 113, 147, 173, 251, 225, 38, 24, 199, 239, 190, 173, 108, 31, 71, 27, 133, 80, 110, 120, 251, 231, 174, 198, 132, 129, 159, 29, 246, 19, 234, 8, 114, 30, 17, 212, 186, 58, 95, 200, 94, 59, 26, 18, 6, 124, 119, 216, 79, 174, 21, 65, 185, 227, 239, 158},
			},
		},
	}

	testEmailFromFile(t, fp, expectedEmail)
}

func TestParseEmailKoreanMultipartSignedEuckrOverQuotedprintable(t *testing.T) {
	fp := "tests/test_korean_multipart_signed_euc-kr_over_quoted-printable.txt"
	tz, _ := time.LoadLocation("Asia/Seoul")
	expectedDate, _ := time.Parse(
		time.RFC1123Z+" (MST)",
		time.Date(2019, time.April, 1, 7, 55, 0, 0, tz).Format(time.RFC1123Z+" (MST)"))
	expectedEmail := Email{
		Headers: Headers{
			Date:    expectedDate,
			Subject: "📧 Signed Test 한국어 팬그램",
			ReplyTo: []*mail.Address{
				{
					Name:    "앨리스 보내는사람",
					Address: "alice.bonaeneunsalam@example.net",
				},
			},
			Sender: &mail.Address{
				Name:    "앨리스 보내는사람",
				Address: "alice.bonaeneunsalam@example.com",
			},
			From: []*mail.Address{
				{
					Name:    "앨리스 보내는사람",
					Address: "alice.bonaeneunsalam@example.com",
				},
				{
					Name:    "앨리스 보내는사람",
					Address: "alice.bonaeneunsalam@example.net",
				},
			},
			To: []*mail.Address{
				{
					Name:    "밥 수신자",
					Address: "bob.susinja@example.com",
				},
				{
					Name:    "캐롤 수신자",
					Address: "carol.susinja@example.com",
				},
			},
			Cc: []*mail.Address{
				{
					Name:    "댄 수신자",
					Address: "dan.susinja@example.com",
				},
				{
					Name:    "이브 수신자",
					Address: "eve.susinja@example.com",
				},
			},
			Bcc: []*mail.Address{
				{
					Name:    "프랭크 수신자",
					Address: "frank.susinja@example.com",
				},
				{
					Name:    "그레이스 수신자",
					Address: "grace.susinja@example.com",
				},
			},
			MessageID:  "Message-Id-1@example.com",
			InReplyTo:  []MessageId{"Message-Id-0@example.com"},
			References: []MessageId{"Message-Id-0@example.com"},
			Comments:   "Message Header Comment",
			Keywords:   []string{"Keyword 1", "Keyword 2"},
			ResentDate: expectedDate,
			ResentFrom: []*mail.Address{
				{
					Name:    "앨리스 보내는사람",
					Address: "alice.bonaeneunsalam@example.net",
				},
				{
					Name:    "앨리스 보내는사람",
					Address: "alice.bonaeneunsalam@example.com",
				},
			},
			ResentSender: &mail.Address{
				Name:    "앨리스 보내는사람",
				Address: "alice.bonaeneunsalam@example.net",
			},
			ResentTo: []*mail.Address{
				{
					Name:    "밥 수신자",
					Address: "bob.susinja@example.net",
				},
				{
					Name:    "캐롤 수신자",
					Address: "carol.susinja@example.net",
				},
			},
			ResentCc: []*mail.Address{
				{
					Name:    "댄 수신자",
					Address: "dan.susinja@example.net",
				},
				{
					Name:    "이브 수신자",
					Address: "eve.susinja@example.net",
				},
			},
			ResentBcc: []*mail.Address{
				{
					Name:    "프랭크 수신자",
					Address: "frank.susinja@example.net",
				},
				{
					Name:    "그레이스 수신자",
					Address: "grace.susinja@example.net",
				},
			},
			ResentMessageID: "Message-Id-1@example.net",
			ContentType: ContentTypeHeader{
				ContentType: "multipart/signed",
				Params: map[string]string{
					"boundary": "SignedBoundaryString",
					"charset":  "euc-kr",
					"micalg":   "sha1",
					"protocol": "application/pkcs7-signature",
				},
			},
			ExtraHeaders: map[string][]string{
				"X-Clacks-Overhead": {"GNU Terry Pratchett"},
				"X-Script/function/\t !\"#$%&'()*+,-./;<=>?@[\\]^_`{|}~": {
					"TEST VALUE 1\t !\"#$%&'()*+,-./0123456789:;<=>?@ABCDEFGHIJKLMNOPQRSTUVWXYZ[\\]^_` abcdefghijklmnopqrstuvwxyz{|}~",
					"TEST VALUE 2\t !\"#$%&'()*+,-./0123456789:;<=>?@ABCDEFGHIJKLMNOPQRSTUVWXYZ[\\]^_` abcdefghijklmnopqrstuvwxyz{|}~",
				},
			},
		},
		Text: `키스의 고유조건은 입술끼리 만나야 하고 특별한 기술은 필요치 않다.`,
		AttachedFiles: []AttachedFile{
			{
				ContentType: ContentTypeHeader{
					ContentType: "application/pkcs7-signature",
					Params: map[string]string{
						"name": "smime.p7s",
					},
				},
				ContentDisposition: ContentDispositionHeader{
					ContentDisposition: attachment,
					Params: map[string]string{
						"filename": "smime.p7s",
					},
				},
				Data: []byte{130, 28, 135, 132, 117, 46, 142, 18, 97, 140, 126, 251, 159, 193, 199, 25, 58, 223, 189, 185, 227, 239, 158, 173, 108, 31, 71, 27, 133, 80, 165, 252, 133, 227, 174, 198, 132, 129, 159, 29, 246, 19, 235, 133, 80, 165, 252, 133, 227, 174, 198, 132, 129, 159, 29, 246, 19, 234, 49, 251, 238, 127, 7, 28, 104, 33, 200, 120, 71, 82, 232, 225, 38, 30, 249, 234, 214, 193, 244, 113, 147, 173, 251, 219, 158, 57, 252, 28, 113, 147, 173, 251, 225, 38, 24, 199, 239, 190, 173, 108, 31, 71, 27, 133, 80, 110, 120, 251, 231, 174, 198, 132, 129, 159, 29, 246, 19, 234, 8, 114, 30, 17, 212, 186, 58, 95, 200, 94, 59, 26, 18, 6, 124, 119, 216, 79, 174, 21, 65, 185, 227, 239, 158},
			},
		},
	}

	testEmailFromFile(t, fp, expectedEmail)
}

func TestParseEmailPolishPlaintextUtf8OverBase64(t *testing.T) {
	fp := "tests/test_polish_plaintext_utf-8_over_base64.txt"
	tz, _ := time.LoadLocation("Europe/Warsaw")
	expectedDate, _ := time.Parse(
		time.RFC1123Z+" (MST)",
		time.Date(2019, time.April, 1, 7, 55, 0, 0, tz).Format(time.RFC1123Z+" (MST)"))
	expectedEmail := Email{
		Headers: Headers{
			Date:    expectedDate,
			Subject: "📧 Test Polskie pangramy",
			ReplyTo: []*mail.Address{
				{
					Name:    "Nadająca, Alicja",
					Address: "alicja.nadajaca@example.net",
				},
			},
			Sender: &mail.Address{
				Name:    "Nadająca, Alicja",
				Address: "alicja.nadajaca@example.com",
			},
			From: []*mail.Address{
				{
					Name:    "Nadająca, Alicja",
					Address: "alicja.nadajaca@example.com",
				},
				{
					Name:    "Nadająca, Alicja",
					Address: "alicja.nadajaca@example.net",
				},
			},
			To: []*mail.Address{
				{
					Name:    "Odbierający, Bob",
					Address: "bob.odbierajacy@example.com",
				},
				{
					Name:    "Odbierająca, Karolina",
					Address: "karolina.odbierajaca@example.com",
				},
			},
			Cc: []*mail.Address{
				{
					Name:    "Odbierający, Daniel",
					Address: "daniel.odbierajacy@example.com",
				},
				{
					Name:    "Odbierająca, Ewa",
					Address: "ewa.odbierajaca@example.com",
				},
			},
			Bcc: []*mail.Address{
				{
					Name:    "Odbierający, Franek",
					Address: "franek.odbierajacy@example.com",
				},
				{
					Name:    "Odbierająca, Grażyna",
					Address: "grazyna.odbierajaca@example.com",
				},
			},
			MessageID:  "Message-Id-1@example.com",
			InReplyTo:  []MessageId{"Message-Id-0@example.com"},
			References: []MessageId{"Message-Id-0@example.com"},
			Comments:   "Message Header Comment",
			Keywords:   []string{"Keyword 1", "Keyword 2"},
			ResentDate: expectedDate,
			ResentFrom: []*mail.Address{
				{
					Name:    "Nadająca, Alicja",
					Address: "alicja.nadajaca@example.net",
				},
				{
					Name:    "Nadająca, Alicja",
					Address: "alicja.nadajaca@example.com",
				},
			},
			ResentSender: &mail.Address{
				Name:    "Nadająca, Alicja",
				Address: "alicja.nadajaca@example.net",
			},
			ResentTo: []*mail.Address{
				{
					Name:    "Odbierający, Bob",
					Address: "bob.odbierajacy@example.net",
				},
				{
					Name:    "Odbierająca, Karolina",
					Address: "karolina.odbierajaca@example.net",
				},
			},
			ResentCc: []*mail.Address{
				{
					Name:    "Odbierający, Daniel",
					Address: "daniel.odbierajacy@example.net",
				},
				{
					Name:    "Odbierająca, Ewa",
					Address: "ewa.odbierajaca@example.net",
				},
			},
			ResentBcc: []*mail.Address{
				{
					Name:    "Odbierający, Franek",
					Address: "franek.odbierajacy@example.net",
				},
				{
					Name:    "Odbierająca, Grażyna",
					Address: "grazyna.odbierajaca@example.net",
				},
			},
			ResentMessageID: "Message-Id-1@example.net",
			ContentType: ContentTypeHeader{
				ContentType: "text/plain",
				Params: map[string]string{
					"charset": "utf-8",
				},
			},
			ExtraHeaders: map[string][]string{
				"X-Clacks-Overhead": {"GNU Terry Pratchett"},
				"X-Script/function/\t !\"#$%&'()*+,-./;<=>?@[\\]^_`{|}~": {
					"TEST VALUE 1\t !\"#$%&'()*+,-./0123456789:;<=>?@ABCDEFGHIJKLMNOPQRSTUVWXYZ[\\]^_` abcdefghijklmnopqrstuvwxyz{|}~",
					"TEST VALUE 2\t !\"#$%&'()*+,-./0123456789:;<=>?@ABCDEFGHIJKLMNOPQRSTUVWXYZ[\\]^_` abcdefghijklmnopqrstuvwxyz{|}~",
				},
			},
		},
		Text: `Jeżu klątw, spłódź Finom część gry hańb!
Pójdźże, kiń tę chmurność w głąb flaszy!
Mężny bądź chroń pułk twój i sześć flag.
Filmuj rzeź żądań, pość, gnęb chłystków!
Pchnąć w tę łódź jeża lub ośm skrzyń fig.
Dość gróźb fuzją, klnę, pych i małżeństw!
Pójdź w loch zbić małżeńską gęś futryn!
Chwyć małżonkę, strój bądź pleśń z fugi.`,
	}

	testEmailFromFile(t, fp, expectedEmail)
}

func TestParseEmailPolishPlaintextUtf8OverQuotedprintable(t *testing.T) {
	fp := "tests/test_polish_plaintext_utf-8_over_quoted-printable.txt"
	tz, _ := time.LoadLocation("Europe/Warsaw")
	expectedDate, _ := time.Parse(
		time.RFC1123Z+" (MST)",
		time.Date(2019, time.April, 1, 7, 55, 0, 0, tz).Format(time.RFC1123Z+" (MST)"))
	expectedEmail := Email{
		Headers: Headers{
			Date:    expectedDate,
			Subject: "📧 Test Polskie pangramy",
			ReplyTo: []*mail.Address{
				{
					Name:    "Nadająca, Alicja",
					Address: "alicja.nadajaca@example.net",
				},
			},
			Sender: &mail.Address{
				Name:    "Nadająca, Alicja",
				Address: "alicja.nadajaca@example.com",
			},
			From: []*mail.Address{
				{
					Name:    "Nadająca, Alicja",
					Address: "alicja.nadajaca@example.com",
				},
				{
					Name:    "Nadająca, Alicja",
					Address: "alicja.nadajaca@example.net",
				},
			},
			To: []*mail.Address{
				{
					Name:    "Odbierający, Bob",
					Address: "bob.odbierajacy@example.com",
				},
				{
					Name:    "Odbierająca, Karolina",
					Address: "karolina.odbierajaca@example.com",
				},
			},
			Cc: []*mail.Address{
				{
					Name:    "Odbierający, Daniel",
					Address: "daniel.odbierajacy@example.com",
				},
				{
					Name:    "Odbierająca, Ewa",
					Address: "ewa.odbierajaca@example.com",
				},
			},
			Bcc: []*mail.Address{
				{
					Name:    "Odbierający, Franek",
					Address: "franek.odbierajacy@example.com",
				},
				{
					Name:    "Odbierająca, Grażyna",
					Address: "grazyna.odbierajaca@example.com",
				},
			},
			MessageID:  "Message-Id-1@example.com",
			InReplyTo:  []MessageId{"Message-Id-0@example.com"},
			References: []MessageId{"Message-Id-0@example.com"},
			Comments:   "Message Header Comment",
			Keywords:   []string{"Keyword 1", "Keyword 2"},
			ResentDate: expectedDate,
			ResentFrom: []*mail.Address{
				{
					Name:    "Nadająca, Alicja",
					Address: "alicja.nadajaca@example.net",
				},
				{
					Name:    "Nadająca, Alicja",
					Address: "alicja.nadajaca@example.com",
				},
			},
			ResentSender: &mail.Address{
				Name:    "Nadająca, Alicja",
				Address: "alicja.nadajaca@example.net",
			},
			ResentTo: []*mail.Address{
				{
					Name:    "Odbierający, Bob",
					Address: "bob.odbierajacy@example.net",
				},
				{
					Name:    "Odbierająca, Karolina",
					Address: "karolina.odbierajaca@example.net",
				},
			},
			ResentCc: []*mail.Address{
				{
					Name:    "Odbierający, Daniel",
					Address: "daniel.odbierajacy@example.net",
				},
				{
					Name:    "Odbierająca, Ewa",
					Address: "ewa.odbierajaca@example.net",
				},
			},
			ResentBcc: []*mail.Address{
				{
					Name:    "Odbierający, Franek",
					Address: "franek.odbierajacy@example.net",
				},
				{
					Name:    "Odbierająca, Grażyna",
					Address: "grazyna.odbierajaca@example.net",
				},
			},
			ResentMessageID: "Message-Id-1@example.net",
			ContentType: ContentTypeHeader{
				ContentType: "text/plain",
				Params: map[string]string{
					"charset": "utf-8",
				},
			},
			ExtraHeaders: map[string][]string{
				"X-Clacks-Overhead": {"GNU Terry Pratchett"},
				"X-Script/function/\t !\"#$%&'()*+,-./;<=>?@[\\]^_`{|}~": {
					"TEST VALUE 1\t !\"#$%&'()*+,-./0123456789:;<=>?@ABCDEFGHIJKLMNOPQRSTUVWXYZ[\\]^_` abcdefghijklmnopqrstuvwxyz{|}~",
					"TEST VALUE 2\t !\"#$%&'()*+,-./0123456789:;<=>?@ABCDEFGHIJKLMNOPQRSTUVWXYZ[\\]^_` abcdefghijklmnopqrstuvwxyz{|}~",
				},
			},
		},
		Text: `Jeżu klątw, spłódź Finom część gry hańb!
Pójdźże, kiń tę chmurność w głąb flaszy!
Mężny bądź chroń pułk twój i sześć flag.
Filmuj rzeź żądań, pość, gnęb chłystków!
Pchnąć w tę łódź jeża lub ośm skrzyń fig.
Dość gróźb fuzją, klnę, pych i małżeństw!
Pójdź w loch zbić małżeńską gęś futryn!
Chwyć małżonkę, strój bądź pleśń z fugi.`,
	}

	testEmailFromFile(t, fp, expectedEmail)
}

func TestParseEmailPolishPlaintextIso88592OverBase64(t *testing.T) {
	fp := "tests/test_polish_plaintext_iso-8859-2_over_base64.txt"
	tz, _ := time.LoadLocation("Europe/Warsaw")
	expectedDate, _ := time.Parse(
		time.RFC1123Z+" (MST)",
		time.Date(2019, time.April, 1, 7, 55, 0, 0, tz).Format(time.RFC1123Z+" (MST)"))
	expectedEmail := Email{
		Headers: Headers{
			Date:    expectedDate,
			Subject: "📧 Test Polskie pangramy",
			ReplyTo: []*mail.Address{
				{
					Name:    "Nadająca, Alicja",
					Address: "alicja.nadajaca@example.net",
				},
			},
			Sender: &mail.Address{
				Name:    "Nadająca, Alicja",
				Address: "alicja.nadajaca@example.com",
			},
			From: []*mail.Address{
				{
					Name:    "Nadająca, Alicja",
					Address: "alicja.nadajaca@example.com",
				},
				{
					Name:    "Nadająca, Alicja",
					Address: "alicja.nadajaca@example.net",
				},
			},
			To: []*mail.Address{
				{
					Name:    "Odbierający, Bob",
					Address: "bob.odbierajacy@example.com",
				},
				{
					Name:    "Odbierająca, Karolina",
					Address: "karolina.odbierajaca@example.com",
				},
			},
			Cc: []*mail.Address{
				{
					Name:    "Odbierający, Daniel",
					Address: "daniel.odbierajacy@example.com",
				},
				{
					Name:    "Odbierająca, Ewa",
					Address: "ewa.odbierajaca@example.com",
				},
			},
			Bcc: []*mail.Address{
				{
					Name:    "Odbierający, Franek",
					Address: "franek.odbierajacy@example.com",
				},
				{
					Name:    "Odbierająca, Grażyna",
					Address: "grazyna.odbierajaca@example.com",
				},
			},
			MessageID:  "Message-Id-1@example.com",
			InReplyTo:  []MessageId{"Message-Id-0@example.com"},
			References: []MessageId{"Message-Id-0@example.com"},
			Comments:   "Message Header Comment",
			Keywords:   []string{"Keyword 1", "Keyword 2"},
			ResentDate: expectedDate,
			ResentFrom: []*mail.Address{
				{
					Name:    "Nadająca, Alicja",
					Address: "alicja.nadajaca@example.net",
				},
				{
					Name:    "Nadająca, Alicja",
					Address: "alicja.nadajaca@example.com",
				},
			},
			ResentSender: &mail.Address{
				Name:    "Nadająca, Alicja",
				Address: "alicja.nadajaca@example.net",
			},
			ResentTo: []*mail.Address{
				{
					Name:    "Odbierający, Bob",
					Address: "bob.odbierajacy@example.net",
				},
				{
					Name:    "Odbierająca, Karolina",
					Address: "karolina.odbierajaca@example.net",
				},
			},
			ResentCc: []*mail.Address{
				{
					Name:    "Odbierający, Daniel",
					Address: "daniel.odbierajacy@example.net",
				},
				{
					Name:    "Odbierająca, Ewa",
					Address: "ewa.odbierajaca@example.net",
				},
			},
			ResentBcc: []*mail.Address{
				{
					Name:    "Odbierający, Franek",
					Address: "franek.odbierajacy@example.net",
				},
				{
					Name:    "Odbierająca, Grażyna",
					Address: "grazyna.odbierajaca@example.net",
				},
			},
			ResentMessageID: "Message-Id-1@example.net",
			ContentType: ContentTypeHeader{
				ContentType: "text/plain",
				Params: map[string]string{
					"charset": "iso-8859-2",
				},
			},
			ExtraHeaders: map[string][]string{
				"X-Clacks-Overhead": {"GNU Terry Pratchett"},
				"X-Script/function/\t !\"#$%&'()*+,-./;<=>?@[\\]^_`{|}~": {
					"TEST VALUE 1\t !\"#$%&'()*+,-./0123456789:;<=>?@ABCDEFGHIJKLMNOPQRSTUVWXYZ[\\]^_` abcdefghijklmnopqrstuvwxyz{|}~",
					"TEST VALUE 2\t !\"#$%&'()*+,-./0123456789:;<=>?@ABCDEFGHIJKLMNOPQRSTUVWXYZ[\\]^_` abcdefghijklmnopqrstuvwxyz{|}~",
				},
			},
		},
		Text: `Jeżu klątw, spłódź Finom część gry hańb!
Pójdźże, kiń tę chmurność w głąb flaszy!
Mężny bądź chroń pułk twój i sześć flag.
Filmuj rzeź żądań, pość, gnęb chłystków!
Pchnąć w tę łódź jeża lub ośm skrzyń fig.
Dość gróźb fuzją, klnę, pych i małżeństw!
Pójdź w loch zbić małżeńską gęś futryn!
Chwyć małżonkę, strój bądź pleśń z fugi.`,
	}

	testEmailFromFile(t, fp, expectedEmail)
}

func TestParseEmailPolishPlaintextIso88592OverQuotedprintable(t *testing.T) {
	fp := "tests/test_polish_plaintext_iso-8859-2_over_quoted-printable.txt"
	tz, _ := time.LoadLocation("Europe/Warsaw")
	expectedDate, _ := time.Parse(
		time.RFC1123Z+" (MST)",
		time.Date(2019, time.April, 1, 7, 55, 0, 0, tz).Format(time.RFC1123Z+" (MST)"))
	expectedEmail := Email{
		Headers: Headers{
			Date:    expectedDate,
			Subject: "📧 Test Polskie pangramy",
			ReplyTo: []*mail.Address{
				{
					Name:    "Nadająca, Alicja",
					Address: "alicja.nadajaca@example.net",
				},
			},
			Sender: &mail.Address{
				Name:    "Nadająca, Alicja",
				Address: "alicja.nadajaca@example.com",
			},
			From: []*mail.Address{
				{
					Name:    "Nadająca, Alicja",
					Address: "alicja.nadajaca@example.com",
				},
				{
					Name:    "Nadająca, Alicja",
					Address: "alicja.nadajaca@example.net",
				},
			},
			To: []*mail.Address{
				{
					Name:    "Odbierający, Bob",
					Address: "bob.odbierajacy@example.com",
				},
				{
					Name:    "Odbierająca, Karolina",
					Address: "karolina.odbierajaca@example.com",
				},
			},
			Cc: []*mail.Address{
				{
					Name:    "Odbierający, Daniel",
					Address: "daniel.odbierajacy@example.com",
				},
				{
					Name:    "Odbierająca, Ewa",
					Address: "ewa.odbierajaca@example.com",
				},
			},
			Bcc: []*mail.Address{
				{
					Name:    "Odbierający, Franek",
					Address: "franek.odbierajacy@example.com",
				},
				{
					Name:    "Odbierająca, Grażyna",
					Address: "grazyna.odbierajaca@example.com",
				},
			},
			MessageID:  "Message-Id-1@example.com",
			InReplyTo:  []MessageId{"Message-Id-0@example.com"},
			References: []MessageId{"Message-Id-0@example.com"},
			Comments:   "Message Header Comment",
			Keywords:   []string{"Keyword 1", "Keyword 2"},
			ResentDate: expectedDate,
			ResentFrom: []*mail.Address{
				{
					Name:    "Nadająca, Alicja",
					Address: "alicja.nadajaca@example.net",
				},
				{
					Name:    "Nadająca, Alicja",
					Address: "alicja.nadajaca@example.com",
				},
			},
			ResentSender: &mail.Address{
				Name:    "Nadająca, Alicja",
				Address: "alicja.nadajaca@example.net",
			},
			ResentTo: []*mail.Address{
				{
					Name:    "Odbierający, Bob",
					Address: "bob.odbierajacy@example.net",
				},
				{
					Name:    "Odbierająca, Karolina",
					Address: "karolina.odbierajaca@example.net",
				},
			},
			ResentCc: []*mail.Address{
				{
					Name:    "Odbierający, Daniel",
					Address: "daniel.odbierajacy@example.net",
				},
				{
					Name:    "Odbierająca, Ewa",
					Address: "ewa.odbierajaca@example.net",
				},
			},
			ResentBcc: []*mail.Address{
				{
					Name:    "Odbierający, Franek",
					Address: "franek.odbierajacy@example.net",
				},
				{
					Name:    "Odbierająca, Grażyna",
					Address: "grazyna.odbierajaca@example.net",
				},
			},
			ResentMessageID: "Message-Id-1@example.net",
			ContentType: ContentTypeHeader{
				ContentType: "text/plain",
				Params: map[string]string{
					"charset": "iso-8859-2",
				},
			},
			ExtraHeaders: map[string][]string{
				"X-Clacks-Overhead": {"GNU Terry Pratchett"},
				"X-Script/function/\t !\"#$%&'()*+,-./;<=>?@[\\]^_`{|}~": {
					"TEST VALUE 1\t !\"#$%&'()*+,-./0123456789:;<=>?@ABCDEFGHIJKLMNOPQRSTUVWXYZ[\\]^_` abcdefghijklmnopqrstuvwxyz{|}~",
					"TEST VALUE 2\t !\"#$%&'()*+,-./0123456789:;<=>?@ABCDEFGHIJKLMNOPQRSTUVWXYZ[\\]^_` abcdefghijklmnopqrstuvwxyz{|}~",
				},
			},
		},
		Text: `Jeżu klątw, spłódź Finom część gry hańb!
Pójdźże, kiń tę chmurność w głąb flaszy!
Mężny bądź chroń pułk twój i sześć flag.
Filmuj rzeź żądań, pość, gnęb chłystków!
Pchnąć w tę łódź jeża lub ośm skrzyń fig.
Dość gróźb fuzją, klnę, pych i małżeństw!
Pójdź w loch zbić małżeńską gęś futryn!
Chwyć małżonkę, strój bądź pleśń z fugi.`,
	}

	testEmailFromFile(t, fp, expectedEmail)
}

func TestParseEmailPolishMultipartRelatedUtf8OverBase64(t *testing.T) {
	fp := "tests/test_polish_multipart_related_utf-8_over_base64.txt"
	tz, _ := time.LoadLocation("Europe/Warsaw")
	expectedDate, _ := time.Parse(
		time.RFC1123Z+" (MST)",
		time.Date(2019, time.April, 1, 7, 55, 0, 0, tz).Format(time.RFC1123Z+" (MST)"))
	expectedEmail := Email{
		Headers: Headers{
			Date:    expectedDate,
			Subject: "📧 Test Polskie pangramy",
			ReplyTo: []*mail.Address{
				{
					Name:    "Nadająca, Alicja",
					Address: "alicja.nadajaca@example.net",
				},
			},
			Sender: &mail.Address{
				Name:    "Nadająca, Alicja",
				Address: "alicja.nadajaca@example.com",
			},
			From: []*mail.Address{
				{
					Name:    "Nadająca, Alicja",
					Address: "alicja.nadajaca@example.com",
				},
				{
					Name:    "Nadająca, Alicja",
					Address: "alicja.nadajaca@example.net",
				},
			},
			To: []*mail.Address{
				{
					Name:    "Odbierający, Bob",
					Address: "bob.odbierajacy@example.com",
				},
				{
					Name:    "Odbierająca, Karolina",
					Address: "karolina.odbierajaca@example.com",
				},
			},
			Cc: []*mail.Address{
				{
					Name:    "Odbierający, Daniel",
					Address: "daniel.odbierajacy@example.com",
				},
				{
					Name:    "Odbierająca, Ewa",
					Address: "ewa.odbierajaca@example.com",
				},
			},
			Bcc: []*mail.Address{
				{
					Name:    "Odbierający, Franek",
					Address: "franek.odbierajacy@example.com",
				},
				{
					Name:    "Odbierająca, Grażyna",
					Address: "grazyna.odbierajaca@example.com",
				},
			},
			MessageID:  "Message-Id-1@example.com",
			InReplyTo:  []MessageId{"Message-Id-0@example.com"},
			References: []MessageId{"Message-Id-0@example.com"},
			Comments:   "Message Header Comment",
			Keywords:   []string{"Keyword 1", "Keyword 2"},
			ResentDate: expectedDate,
			ResentFrom: []*mail.Address{
				{
					Name:    "Nadająca, Alicja",
					Address: "alicja.nadajaca@example.net",
				},
				{
					Name:    "Nadająca, Alicja",
					Address: "alicja.nadajaca@example.com",
				},
			},
			ResentSender: &mail.Address{
				Name:    "Nadająca, Alicja",
				Address: "alicja.nadajaca@example.net",
			},
			ResentTo: []*mail.Address{
				{
					Name:    "Odbierający, Bob",
					Address: "bob.odbierajacy@example.net",
				},
				{
					Name:    "Odbierająca, Karolina",
					Address: "karolina.odbierajaca@example.net",
				},
			},
			ResentCc: []*mail.Address{
				{
					Name:    "Odbierający, Daniel",
					Address: "daniel.odbierajacy@example.net",
				},
				{
					Name:    "Odbierająca, Ewa",
					Address: "ewa.odbierajaca@example.net",
				},
			},
			ResentBcc: []*mail.Address{
				{
					Name:    "Odbierający, Franek",
					Address: "franek.odbierajacy@example.net",
				},
				{
					Name:    "Odbierająca, Grażyna",
					Address: "grazyna.odbierajaca@example.net",
				},
			},
			ResentMessageID: "Message-Id-1@example.net",
			ContentType: ContentTypeHeader{
				ContentType: "multipart/related",
				Params: map[string]string{
					"boundary": "RelatedBoundaryString",
					"charset":  "utf-8",
				},
			},
			ExtraHeaders: map[string][]string{
				"X-Clacks-Overhead": {"GNU Terry Pratchett"},
				"X-Script/function/\t !\"#$%&'()*+,-./;<=>?@[\\]^_`{|}~": {
					"TEST VALUE 1\t !\"#$%&'()*+,-./0123456789:;<=>?@ABCDEFGHIJKLMNOPQRSTUVWXYZ[\\]^_` abcdefghijklmnopqrstuvwxyz{|}~",
					"TEST VALUE 2\t !\"#$%&'()*+,-./0123456789:;<=>?@ABCDEFGHIJKLMNOPQRSTUVWXYZ[\\]^_` abcdefghijklmnopqrstuvwxyz{|}~",
				},
			},
		},
		Text: `Jeżu klątw, spłódź Finom część gry hańb!
Pójdźże, kiń tę chmurność w głąb flaszy!
Mężny bądź chroń pułk twój i sześć flag.
Filmuj rzeź żądań, pość, gnęb chłystków!
Pchnąć w tę łódź jeża lub ośm skrzyń fig.
Dość gróźb fuzją, klnę, pych i małżeństw!
Pójdź w loch zbić małżeńską gęś futryn!
Chwyć małżonkę, strój bądź pleśń z fugi.`,
		EnrichedText: `<bold>Jeżu klątw, spłódź Finom część gry hańb!</bold>
<italic>Pójdźże, kiń tę chmurność w głąb flaszy!</italic>
<fixed>Mężny bądź chroń pułk twój i sześć flag.</fixed>
<underline>Filmuj rzeź żądań, pość, gnęb chłystków!</underline>
Pchnąć w tę łódź jeża lub ośm skrzyń fig.
Dość gróźb fuzją, klnę, pych i małżeństw!
Pójdź w loch zbić małżeńską gęś futryn!
Chwyć małżonkę, strój bądź pleśń z fugi.`,
		HTML: `<html>
<div dir="ltr">
<p>Jeżu klątw, spłódź Finom część gry hańb!</p>
<p>Pójdźże, kiń tę chmurność w głąb flaszy!</p>
<p>Mężny bądź chroń pułk twój i sześć flag.</p>
<p>Filmuj rzeź żądań, pość, gnęb chłystków!</p>
<p>Pchnąć w tę łódź jeża lub ośm skrzyń fig.</p>
<p>Dość gróźb fuzją, klnę, pych i małżeństw!</p>
<p>Pójdź w loch zbić małżeńską gęś futryn!</p>
<p>Chwyć małżonkę, strój bądź pleśń z fugi.</p>
</div>
</html>`,
		InlineFiles: []InlineFile{
			{
				ContentID: "inline-jpg-image.jpg@example.com",
				ContentType: ContentTypeHeader{
					ContentType: "image/jpeg",
					Params: map[string]string{
						"name": "inline-jpg-image-name.jpg",
					},
				},
				ContentDisposition: ContentDispositionHeader{
					ContentDisposition: inline,
					Params: map[string]string{
						"filename": "inline-jpg-image-filename.jpg",
					},
				},
				Data: []byte{255, 216, 255, 219, 0, 67, 0, 3, 2, 2, 2, 2, 2, 3, 2, 2, 2, 3, 3, 3, 3, 4,
					6, 4, 4, 4, 4, 4, 8, 6, 6, 5, 6, 9, 8, 10, 10, 9, 8, 9, 9, 10, 12, 15, 12, 10, 11, 14, 11, 9, 9,
					13, 17, 13, 14, 15, 16, 16, 17, 16, 10, 12, 18, 19, 18, 16, 19, 15, 16, 16, 16, 255, 201, 0, 11, 8,
					0, 1, 0, 1, 1, 1, 17, 0, 255, 204, 0, 6, 0, 16, 16, 5, 255, 218, 0, 8, 1, 1, 0, 0, 63, 0, 210, 207,
					32, 255, 217},
			},
		},
	}

	testEmailFromFile(t, fp, expectedEmail)
}
func TestParseEmailPolishMultipartRelatedUtf8OverQuotedprintable(t *testing.T) {
	fp := "tests/test_polish_multipart_related_utf-8_over_quoted-printable.txt"
	tz, _ := time.LoadLocation("Europe/Warsaw")
	expectedDate, _ := time.Parse(
		time.RFC1123Z+" (MST)",
		time.Date(2019, time.April, 1, 7, 55, 0, 0, tz).Format(time.RFC1123Z+" (MST)"))
	expectedEmail := Email{
		Headers: Headers{
			Date:    expectedDate,
			Subject: "📧 Test Polskie pangramy",
			ReplyTo: []*mail.Address{
				{
					Name:    "Nadająca, Alicja",
					Address: "alicja.nadajaca@example.net",
				},
			},
			Sender: &mail.Address{
				Name:    "Nadająca, Alicja",
				Address: "alicja.nadajaca@example.com",
			},
			From: []*mail.Address{
				{
					Name:    "Nadająca, Alicja",
					Address: "alicja.nadajaca@example.com",
				},
				{
					Name:    "Nadająca, Alicja",
					Address: "alicja.nadajaca@example.net",
				},
			},
			To: []*mail.Address{
				{
					Name:    "Odbierający, Bob",
					Address: "bob.odbierajacy@example.com",
				},
				{
					Name:    "Odbierająca, Karolina",
					Address: "karolina.odbierajaca@example.com",
				},
			},
			Cc: []*mail.Address{
				{
					Name:    "Odbierający, Daniel",
					Address: "daniel.odbierajacy@example.com",
				},
				{
					Name:    "Odbierająca, Ewa",
					Address: "ewa.odbierajaca@example.com",
				},
			},
			Bcc: []*mail.Address{
				{
					Name:    "Odbierający, Franek",
					Address: "franek.odbierajacy@example.com",
				},
				{
					Name:    "Odbierająca, Grażyna",
					Address: "grazyna.odbierajaca@example.com",
				},
			},
			MessageID:  "Message-Id-1@example.com",
			InReplyTo:  []MessageId{"Message-Id-0@example.com"},
			References: []MessageId{"Message-Id-0@example.com"},
			Comments:   "Message Header Comment",
			Keywords:   []string{"Keyword 1", "Keyword 2"},
			ResentDate: expectedDate,
			ResentFrom: []*mail.Address{
				{
					Name:    "Nadająca, Alicja",
					Address: "alicja.nadajaca@example.net",
				},
				{
					Name:    "Nadająca, Alicja",
					Address: "alicja.nadajaca@example.com",
				},
			},
			ResentSender: &mail.Address{
				Name:    "Nadająca, Alicja",
				Address: "alicja.nadajaca@example.net",
			},
			ResentTo: []*mail.Address{
				{
					Name:    "Odbierający, Bob",
					Address: "bob.odbierajacy@example.net",
				},
				{
					Name:    "Odbierająca, Karolina",
					Address: "karolina.odbierajaca@example.net",
				},
			},
			ResentCc: []*mail.Address{
				{
					Name:    "Odbierający, Daniel",
					Address: "daniel.odbierajacy@example.net",
				},
				{
					Name:    "Odbierająca, Ewa",
					Address: "ewa.odbierajaca@example.net",
				},
			},
			ResentBcc: []*mail.Address{
				{
					Name:    "Odbierający, Franek",
					Address: "franek.odbierajacy@example.net",
				},
				{
					Name:    "Odbierająca, Grażyna",
					Address: "grazyna.odbierajaca@example.net",
				},
			},
			ResentMessageID: "Message-Id-1@example.net",
			ContentType: ContentTypeHeader{
				ContentType: "multipart/related",
				Params: map[string]string{
					"boundary": "RelatedBoundaryString",
					"charset":  "utf-8",
				},
			},
			ExtraHeaders: map[string][]string{
				"X-Clacks-Overhead": {"GNU Terry Pratchett"},
				"X-Script/function/\t !\"#$%&'()*+,-./;<=>?@[\\]^_`{|}~": {
					"TEST VALUE 1\t !\"#$%&'()*+,-./0123456789:;<=>?@ABCDEFGHIJKLMNOPQRSTUVWXYZ[\\]^_` abcdefghijklmnopqrstuvwxyz{|}~",
					"TEST VALUE 2\t !\"#$%&'()*+,-./0123456789:;<=>?@ABCDEFGHIJKLMNOPQRSTUVWXYZ[\\]^_` abcdefghijklmnopqrstuvwxyz{|}~",
				},
			},
		},
		Text: `Jeżu klątw, spłódź Finom część gry hańb!
Pójdźże, kiń tę chmurność w głąb flaszy!
Mężny bądź chroń pułk twój i sześć flag.
Filmuj rzeź żądań, pość, gnęb chłystków!
Pchnąć w tę łódź jeża lub ośm skrzyń fig.
Dość gróźb fuzją, klnę, pych i małżeństw!
Pójdź w loch zbić małżeńską gęś futryn!
Chwyć małżonkę, strój bądź pleśń z fugi.`,
		EnrichedText: `<bold>Jeżu klątw, spłódź Finom część gry hańb!</bold>
<italic>Pójdźże, kiń tę chmurność w głąb flaszy!</italic>
<fixed>Mężny bądź chroń pułk twój i sześć flag.</fixed>
<underline>Filmuj rzeź żądań, pość, gnęb chłystków!</underline>
Pchnąć w tę łódź jeża lub ośm skrzyń fig.
Dość gróźb fuzją, klnę, pych i małżeństw!
Pójdź w loch zbić małżeńską gęś futryn!
Chwyć małżonkę, strój bądź pleśń z fugi.`,
		HTML: `<html>
<div dir="ltr">
<p>Jeżu klątw, spłódź Finom część gry hańb!</p>
<p>Pójdźże, kiń tę chmurność w głąb flaszy!</p>
<p>Mężny bądź chroń pułk twój i sześć flag.</p>
<p>Filmuj rzeź żądań, pość, gnęb chłystków!</p>
<p>Pchnąć w tę łódź jeża lub ośm skrzyń fig.</p>
<p>Dość gróźb fuzją, klnę, pych i małżeństw!</p>
<p>Pójdź w loch zbić małżeńską gęś futryn!</p>
<p>Chwyć małżonkę, strój bądź pleśń z fugi.</p>
</div>
</html>`,
		InlineFiles: []InlineFile{
			{
				ContentID: "inline-jpg-image.jpg@example.com",
				ContentType: ContentTypeHeader{
					ContentType: "image/jpeg",
					Params: map[string]string{
						"name": "inline-jpg-image-name.jpg",
					},
				},
				ContentDisposition: ContentDispositionHeader{
					ContentDisposition: inline,
					Params: map[string]string{
						"filename": "inline-jpg-image-filename.jpg",
					},
				},
				Data: []byte{255, 216, 255, 219, 0, 67, 0, 3, 2, 2, 2, 2, 2, 3, 2, 2, 2, 3, 3, 3, 3, 4,
					6, 4, 4, 4, 4, 4, 8, 6, 6, 5, 6, 9, 8, 10, 10, 9, 8, 9, 9, 10, 12, 15, 12, 10, 11, 14, 11, 9, 9,
					13, 17, 13, 14, 15, 16, 16, 17, 16, 10, 12, 18, 19, 18, 16, 19, 15, 16, 16, 16, 255, 201, 0, 11, 8,
					0, 1, 0, 1, 1, 1, 17, 0, 255, 204, 0, 6, 0, 16, 16, 5, 255, 218, 0, 8, 1, 1, 0, 0, 63, 0, 210, 207,
					32, 255, 217},
			},
		},
	}

	testEmailFromFile(t, fp, expectedEmail)
}
func TestParseEmailPolishMultipartRelatedIso88592OverBase64(t *testing.T) {
	fp := "tests/test_polish_multipart_related_iso-8859-2_over_base64.txt"
	tz, _ := time.LoadLocation("Europe/Warsaw")
	expectedDate, _ := time.Parse(
		time.RFC1123Z+" (MST)",
		time.Date(2019, time.April, 1, 7, 55, 0, 0, tz).Format(time.RFC1123Z+" (MST)"))
	expectedEmail := Email{
		Headers: Headers{
			Date:    expectedDate,
			Subject: "📧 Test Polskie pangramy",
			ReplyTo: []*mail.Address{
				{
					Name:    "Nadająca, Alicja",
					Address: "alicja.nadajaca@example.net",
				},
			},
			Sender: &mail.Address{
				Name:    "Nadająca, Alicja",
				Address: "alicja.nadajaca@example.com",
			},
			From: []*mail.Address{
				{
					Name:    "Nadająca, Alicja",
					Address: "alicja.nadajaca@example.com",
				},
				{
					Name:    "Nadająca, Alicja",
					Address: "alicja.nadajaca@example.net",
				},
			},
			To: []*mail.Address{
				{
					Name:    "Odbierający, Bob",
					Address: "bob.odbierajacy@example.com",
				},
				{
					Name:    "Odbierająca, Karolina",
					Address: "karolina.odbierajaca@example.com",
				},
			},
			Cc: []*mail.Address{
				{
					Name:    "Odbierający, Daniel",
					Address: "daniel.odbierajacy@example.com",
				},
				{
					Name:    "Odbierająca, Ewa",
					Address: "ewa.odbierajaca@example.com",
				},
			},
			Bcc: []*mail.Address{
				{
					Name:    "Odbierający, Franek",
					Address: "franek.odbierajacy@example.com",
				},
				{
					Name:    "Odbierająca, Grażyna",
					Address: "grazyna.odbierajaca@example.com",
				},
			},
			MessageID:  "Message-Id-1@example.com",
			InReplyTo:  []MessageId{"Message-Id-0@example.com"},
			References: []MessageId{"Message-Id-0@example.com"},
			Comments:   "Message Header Comment",
			Keywords:   []string{"Keyword 1", "Keyword 2"},
			ResentDate: expectedDate,
			ResentFrom: []*mail.Address{
				{
					Name:    "Nadająca, Alicja",
					Address: "alicja.nadajaca@example.net",
				},
				{
					Name:    "Nadająca, Alicja",
					Address: "alicja.nadajaca@example.com",
				},
			},
			ResentSender: &mail.Address{
				Name:    "Nadająca, Alicja",
				Address: "alicja.nadajaca@example.net",
			},
			ResentTo: []*mail.Address{
				{
					Name:    "Odbierający, Bob",
					Address: "bob.odbierajacy@example.net",
				},
				{
					Name:    "Odbierająca, Karolina",
					Address: "karolina.odbierajaca@example.net",
				},
			},
			ResentCc: []*mail.Address{
				{
					Name:    "Odbierający, Daniel",
					Address: "daniel.odbierajacy@example.net",
				},
				{
					Name:    "Odbierająca, Ewa",
					Address: "ewa.odbierajaca@example.net",
				},
			},
			ResentBcc: []*mail.Address{
				{
					Name:    "Odbierający, Franek",
					Address: "franek.odbierajacy@example.net",
				},
				{
					Name:    "Odbierająca, Grażyna",
					Address: "grazyna.odbierajaca@example.net",
				},
			},
			ResentMessageID: "Message-Id-1@example.net",
			ContentType: ContentTypeHeader{
				ContentType: "multipart/related",
				Params: map[string]string{
					"boundary": "RelatedBoundaryString",
					"charset":  "iso-8859-2",
				},
			},
			ExtraHeaders: map[string][]string{
				"X-Clacks-Overhead": {"GNU Terry Pratchett"},
				"X-Script/function/\t !\"#$%&'()*+,-./;<=>?@[\\]^_`{|}~": {
					"TEST VALUE 1\t !\"#$%&'()*+,-./0123456789:;<=>?@ABCDEFGHIJKLMNOPQRSTUVWXYZ[\\]^_` abcdefghijklmnopqrstuvwxyz{|}~",
					"TEST VALUE 2\t !\"#$%&'()*+,-./0123456789:;<=>?@ABCDEFGHIJKLMNOPQRSTUVWXYZ[\\]^_` abcdefghijklmnopqrstuvwxyz{|}~",
				},
			},
		},
		Text: `Jeżu klątw, spłódź Finom część gry hańb!
Pójdźże, kiń tę chmurność w głąb flaszy!
Mężny bądź chroń pułk twój i sześć flag.
Filmuj rzeź żądań, pość, gnęb chłystków!
Pchnąć w tę łódź jeża lub ośm skrzyń fig.
Dość gróźb fuzją, klnę, pych i małżeństw!
Pójdź w loch zbić małżeńską gęś futryn!
Chwyć małżonkę, strój bądź pleśń z fugi.`,
		EnrichedText: `<bold>Jeżu klątw, spłódź Finom część gry hańb!</bold>
<italic>Pójdźże, kiń tę chmurność w głąb flaszy!</italic>
<fixed>Mężny bądź chroń pułk twój i sześć flag.</fixed>
<underline>Filmuj rzeź żądań, pość, gnęb chłystków!</underline>
Pchnąć w tę łódź jeża lub ośm skrzyń fig.
Dość gróźb fuzją, klnę, pych i małżeństw!
Pójdź w loch zbić małżeńską gęś futryn!
Chwyć małżonkę, strój bądź pleśń z fugi.`,
		HTML: `<html>
<div dir="ltr">
<p>Jeżu klątw, spłódź Finom część gry hańb!</p>
<p>Pójdźże, kiń tę chmurność w głąb flaszy!</p>
<p>Mężny bądź chroń pułk twój i sześć flag.</p>
<p>Filmuj rzeź żądań, pość, gnęb chłystków!</p>
<p>Pchnąć w tę łódź jeża lub ośm skrzyń fig.</p>
<p>Dość gróźb fuzją, klnę, pych i małżeństw!</p>
<p>Pójdź w loch zbić małżeńską gęś futryn!</p>
<p>Chwyć małżonkę, strój bądź pleśń z fugi.</p>
</div>
</html>`,
		InlineFiles: []InlineFile{
			{
				ContentID: "inline-jpg-image.jpg@example.com",
				ContentType: ContentTypeHeader{
					ContentType: "image/jpeg",
					Params: map[string]string{
						"name": "inline-jpg-image-name.jpg",
					},
				},
				ContentDisposition: ContentDispositionHeader{
					ContentDisposition: inline,
					Params: map[string]string{
						"filename": "inline-jpg-image-filename.jpg",
					},
				},
				Data: []byte{255, 216, 255, 219, 0, 67, 0, 3, 2, 2, 2, 2, 2, 3, 2, 2, 2, 3, 3, 3, 3, 4,
					6, 4, 4, 4, 4, 4, 8, 6, 6, 5, 6, 9, 8, 10, 10, 9, 8, 9, 9, 10, 12, 15, 12, 10, 11, 14, 11, 9, 9,
					13, 17, 13, 14, 15, 16, 16, 17, 16, 10, 12, 18, 19, 18, 16, 19, 15, 16, 16, 16, 255, 201, 0, 11, 8,
					0, 1, 0, 1, 1, 1, 17, 0, 255, 204, 0, 6, 0, 16, 16, 5, 255, 218, 0, 8, 1, 1, 0, 0, 63, 0, 210, 207,
					32, 255, 217},
			},
		},
	}

	testEmailFromFile(t, fp, expectedEmail)
}
func TestParseEmailPolishMultipartRelatedIso88592OverQuotedprintable(t *testing.T) {
	fp := "tests/test_polish_multipart_related_iso-8859-2_over_quoted-printable.txt"
	tz, _ := time.LoadLocation("Europe/Warsaw")
	expectedDate, _ := time.Parse(
		time.RFC1123Z+" (MST)",
		time.Date(2019, time.April, 1, 7, 55, 0, 0, tz).Format(time.RFC1123Z+" (MST)"))
	expectedEmail := Email{
		Headers: Headers{
			Date:    expectedDate,
			Subject: "📧 Test Polskie pangramy",
			ReplyTo: []*mail.Address{
				{
					Name:    "Nadająca, Alicja",
					Address: "alicja.nadajaca@example.net",
				},
			},
			Sender: &mail.Address{
				Name:    "Nadająca, Alicja",
				Address: "alicja.nadajaca@example.com",
			},
			From: []*mail.Address{
				{
					Name:    "Nadająca, Alicja",
					Address: "alicja.nadajaca@example.com",
				},
				{
					Name:    "Nadająca, Alicja",
					Address: "alicja.nadajaca@example.net",
				},
			},
			To: []*mail.Address{
				{
					Name:    "Odbierający, Bob",
					Address: "bob.odbierajacy@example.com",
				},
				{
					Name:    "Odbierająca, Karolina",
					Address: "karolina.odbierajaca@example.com",
				},
			},
			Cc: []*mail.Address{
				{
					Name:    "Odbierający, Daniel",
					Address: "daniel.odbierajacy@example.com",
				},
				{
					Name:    "Odbierająca, Ewa",
					Address: "ewa.odbierajaca@example.com",
				},
			},
			Bcc: []*mail.Address{
				{
					Name:    "Odbierający, Franek",
					Address: "franek.odbierajacy@example.com",
				},
				{
					Name:    "Odbierająca, Grażyna",
					Address: "grazyna.odbierajaca@example.com",
				},
			},
			MessageID:  "Message-Id-1@example.com",
			InReplyTo:  []MessageId{"Message-Id-0@example.com"},
			References: []MessageId{"Message-Id-0@example.com"},
			Comments:   "Message Header Comment",
			Keywords:   []string{"Keyword 1", "Keyword 2"},
			ResentDate: expectedDate,
			ResentFrom: []*mail.Address{
				{
					Name:    "Nadająca, Alicja",
					Address: "alicja.nadajaca@example.net",
				},
				{
					Name:    "Nadająca, Alicja",
					Address: "alicja.nadajaca@example.com",
				},
			},
			ResentSender: &mail.Address{
				Name:    "Nadająca, Alicja",
				Address: "alicja.nadajaca@example.net",
			},
			ResentTo: []*mail.Address{
				{
					Name:    "Odbierający, Bob",
					Address: "bob.odbierajacy@example.net",
				},
				{
					Name:    "Odbierająca, Karolina",
					Address: "karolina.odbierajaca@example.net",
				},
			},
			ResentCc: []*mail.Address{
				{
					Name:    "Odbierający, Daniel",
					Address: "daniel.odbierajacy@example.net",
				},
				{
					Name:    "Odbierająca, Ewa",
					Address: "ewa.odbierajaca@example.net",
				},
			},
			ResentBcc: []*mail.Address{
				{
					Name:    "Odbierający, Franek",
					Address: "franek.odbierajacy@example.net",
				},
				{
					Name:    "Odbierająca, Grażyna",
					Address: "grazyna.odbierajaca@example.net",
				},
			},
			ResentMessageID: "Message-Id-1@example.net",
			ContentType: ContentTypeHeader{
				ContentType: "multipart/related",
				Params: map[string]string{
					"boundary": "RelatedBoundaryString",
					"charset":  "iso-8859-2",
				},
			},
			ExtraHeaders: map[string][]string{
				"X-Clacks-Overhead": {"GNU Terry Pratchett"},
				"X-Script/function/\t !\"#$%&'()*+,-./;<=>?@[\\]^_`{|}~": {
					"TEST VALUE 1\t !\"#$%&'()*+,-./0123456789:;<=>?@ABCDEFGHIJKLMNOPQRSTUVWXYZ[\\]^_` abcdefghijklmnopqrstuvwxyz{|}~",
					"TEST VALUE 2\t !\"#$%&'()*+,-./0123456789:;<=>?@ABCDEFGHIJKLMNOPQRSTUVWXYZ[\\]^_` abcdefghijklmnopqrstuvwxyz{|}~",
				},
			},
		},
		Text: `Jeżu klątw, spłódź Finom część gry hańb!
Pójdźże, kiń tę chmurność w głąb flaszy!
Mężny bądź chroń pułk twój i sześć flag.
Filmuj rzeź żądań, pość, gnęb chłystków!
Pchnąć w tę łódź jeża lub ośm skrzyń fig.
Dość gróźb fuzją, klnę, pych i małżeństw!
Pójdź w loch zbić małżeńską gęś futryn!
Chwyć małżonkę, strój bądź pleśń z fugi.`,
		EnrichedText: `<bold>Jeżu klątw, spłódź Finom część gry hańb!</bold>
<italic>Pójdźże, kiń tę chmurność w głąb flaszy!</italic>
<fixed>Mężny bądź chroń pułk twój i sześć flag.</fixed>
<underline>Filmuj rzeź żądań, pość, gnęb chłystków!</underline>
Pchnąć w tę łódź jeża lub ośm skrzyń fig.
Dość gróźb fuzją, klnę, pych i małżeństw!
Pójdź w loch zbić małżeńską gęś futryn!
Chwyć małżonkę, strój bądź pleśń z fugi.`,
		HTML: `<html>
<div dir="ltr">
<p>Jeżu klątw, spłódź Finom część gry hańb!</p>
<p>Pójdźże, kiń tę chmurność w głąb flaszy!</p>
<p>Mężny bądź chroń pułk twój i sześć flag.</p>
<p>Filmuj rzeź żądań, pość, gnęb chłystków!</p>
<p>Pchnąć w tę łódź jeża lub ośm skrzyń fig.</p>
<p>Dość gróźb fuzją, klnę, pych i małżeństw!</p>
<p>Pójdź w loch zbić małżeńską gęś futryn!</p>
<p>Chwyć małżonkę, strój bądź pleśń z fugi.</p>
</div>
</html>`,
		InlineFiles: []InlineFile{
			{
				ContentID: "inline-jpg-image.jpg@example.com",
				ContentType: ContentTypeHeader{
					ContentType: "image/jpeg",
					Params: map[string]string{
						"name": "inline-jpg-image-name.jpg",
					},
				},
				ContentDisposition: ContentDispositionHeader{
					ContentDisposition: inline,
					Params: map[string]string{
						"filename": "inline-jpg-image-filename.jpg",
					},
				},
				Data: []byte{255, 216, 255, 219, 0, 67, 0, 3, 2, 2, 2, 2, 2, 3, 2, 2, 2, 3, 3, 3, 3, 4,
					6, 4, 4, 4, 4, 4, 8, 6, 6, 5, 6, 9, 8, 10, 10, 9, 8, 9, 9, 10, 12, 15, 12, 10, 11, 14, 11, 9, 9,
					13, 17, 13, 14, 15, 16, 16, 17, 16, 10, 12, 18, 19, 18, 16, 19, 15, 16, 16, 16, 255, 201, 0, 11, 8,
					0, 1, 0, 1, 1, 1, 17, 0, 255, 204, 0, 6, 0, 16, 16, 5, 255, 218, 0, 8, 1, 1, 0, 0, 63, 0, 210, 207,
					32, 255, 217},
			},
		},
	}

	testEmailFromFile(t, fp, expectedEmail)
}
func TestParseEmailPolishMultipartMixedUtf8OverBase64(t *testing.T) {
	fp := "tests/test_polish_multipart_mixed_utf-8_over_base64.txt"
	tz, _ := time.LoadLocation("Europe/Warsaw")
	expectedDate, _ := time.Parse(
		time.RFC1123Z+" (MST)",
		time.Date(2019, time.April, 1, 7, 55, 0, 0, tz).Format(time.RFC1123Z+" (MST)"))
	expectedEmail := Email{
		Headers: Headers{
			Date:    expectedDate,
			Subject: "📧 Test Polskie pangramy",
			ReplyTo: []*mail.Address{
				{
					Name:    "Nadająca, Alicja",
					Address: "alicja.nadajaca@example.net",
				},
			},
			Sender: &mail.Address{
				Name:    "Nadająca, Alicja",
				Address: "alicja.nadajaca@example.com",
			},
			From: []*mail.Address{
				{
					Name:    "Nadająca, Alicja",
					Address: "alicja.nadajaca@example.com",
				},
				{
					Name:    "Nadająca, Alicja",
					Address: "alicja.nadajaca@example.net",
				},
			},
			To: []*mail.Address{
				{
					Name:    "Odbierający, Bob",
					Address: "bob.odbierajacy@example.com",
				},
				{
					Name:    "Odbierająca, Karolina",
					Address: "karolina.odbierajaca@example.com",
				},
			},
			Cc: []*mail.Address{
				{
					Name:    "Odbierający, Daniel",
					Address: "daniel.odbierajacy@example.com",
				},
				{
					Name:    "Odbierająca, Ewa",
					Address: "ewa.odbierajaca@example.com",
				},
			},
			Bcc: []*mail.Address{
				{
					Name:    "Odbierający, Franek",
					Address: "franek.odbierajacy@example.com",
				},
				{
					Name:    "Odbierająca, Grażyna",
					Address: "grazyna.odbierajaca@example.com",
				},
			},
			MessageID:  "Message-Id-1@example.com",
			InReplyTo:  []MessageId{"Message-Id-0@example.com"},
			References: []MessageId{"Message-Id-0@example.com"},
			Comments:   "Message Header Comment",
			Keywords:   []string{"Keyword 1", "Keyword 2"},
			ResentDate: expectedDate,
			ResentFrom: []*mail.Address{
				{
					Name:    "Nadająca, Alicja",
					Address: "alicja.nadajaca@example.net",
				},
				{
					Name:    "Nadająca, Alicja",
					Address: "alicja.nadajaca@example.com",
				},
			},
			ResentSender: &mail.Address{
				Name:    "Nadająca, Alicja",
				Address: "alicja.nadajaca@example.net",
			},
			ResentTo: []*mail.Address{
				{
					Name:    "Odbierający, Bob",
					Address: "bob.odbierajacy@example.net",
				},
				{
					Name:    "Odbierająca, Karolina",
					Address: "karolina.odbierajaca@example.net",
				},
			},
			ResentCc: []*mail.Address{
				{
					Name:    "Odbierający, Daniel",
					Address: "daniel.odbierajacy@example.net",
				},
				{
					Name:    "Odbierająca, Ewa",
					Address: "ewa.odbierajaca@example.net",
				},
			},
			ResentBcc: []*mail.Address{
				{
					Name:    "Odbierający, Franek",
					Address: "franek.odbierajacy@example.net",
				},
				{
					Name:    "Odbierająca, Grażyna",
					Address: "grazyna.odbierajaca@example.net",
				},
			},
			ResentMessageID: "Message-Id-1@example.net",
			ContentType: ContentTypeHeader{
				ContentType: "multipart/mixed",
				Params: map[string]string{
					"boundary": "MixedBoundaryString",
					"charset":  "utf-8",
				},
			},
			ExtraHeaders: map[string][]string{
				"X-Clacks-Overhead": {"GNU Terry Pratchett"},
				"X-Script/function/\t !\"#$%&'()*+,-./;<=>?@[\\]^_`{|}~": {
					"TEST VALUE 1\t !\"#$%&'()*+,-./0123456789:;<=>?@ABCDEFGHIJKLMNOPQRSTUVWXYZ[\\]^_` abcdefghijklmnopqrstuvwxyz{|}~",
					"TEST VALUE 2\t !\"#$%&'()*+,-./0123456789:;<=>?@ABCDEFGHIJKLMNOPQRSTUVWXYZ[\\]^_` abcdefghijklmnopqrstuvwxyz{|}~",
				},
			},
		},
		Text: `Jeżu klątw, spłódź Finom część gry hańb!
Pójdźże, kiń tę chmurność w głąb flaszy!
Mężny bądź chroń pułk twój i sześć flag.
Filmuj rzeź żądań, pość, gnęb chłystków!
Pchnąć w tę łódź jeża lub ośm skrzyń fig.
Dość gróźb fuzją, klnę, pych i małżeństw!
Pójdź w loch zbić małżeńską gęś futryn!
Chwyć małżonkę, strój bądź pleśń z fugi.`,
		EnrichedText: `<bold>Jeżu klątw, spłódź Finom część gry hańb!</bold>
<italic>Pójdźże, kiń tę chmurność w głąb flaszy!</italic>
<fixed>Mężny bądź chroń pułk twój i sześć flag.</fixed>
<underline>Filmuj rzeź żądań, pość, gnęb chłystków!</underline>
Pchnąć w tę łódź jeża lub ośm skrzyń fig.
Dość gróźb fuzją, klnę, pych i małżeństw!
Pójdź w loch zbić małżeńską gęś futryn!
Chwyć małżonkę, strój bądź pleśń z fugi.`,
		HTML: `<html>
<div dir="ltr">
<p>Jeżu klątw, spłódź Finom część gry hańb!</p>
<p>Pójdźże, kiń tę chmurność w głąb flaszy!</p>
<p>Mężny bądź chroń pułk twój i sześć flag.</p>
<p>Filmuj rzeź żądań, pość, gnęb chłystków!</p>
<p>Pchnąć w tę łódź jeża lub ośm skrzyń fig.</p>
<p>Dość gróźb fuzją, klnę, pych i małżeństw!</p>
<p>Pójdź w loch zbić małżeńską gęś futryn!</p>
<p>Chwyć małżonkę, strój bądź pleśń z fugi.</p>
</div>
</html>`,
		InlineFiles: []InlineFile{
			{
				ContentType: ContentTypeHeader{
					ContentType: "image/jpeg",
					Params: map[string]string{
						"name": "inline-jpg-image-without-disposition.jpg",
					},
				},
				ContentDisposition: ContentDispositionHeader{
					ContentDisposition: "",
					Params:             map[string]string(nil),
				},
				Data: []byte{255, 216, 255, 219, 0, 67, 0, 3, 2, 2, 2, 2, 2, 3, 2, 2, 2, 3, 3, 3, 3, 4,
					6, 4, 4, 4, 4, 4, 8, 6, 6, 5, 6, 9, 8, 10, 10, 9, 8, 9, 9, 10, 12, 15, 12, 10, 11, 14, 11, 9, 9,
					13, 17, 13, 14, 15, 16, 16, 17, 16, 10, 12, 18, 19, 18, 16, 19, 15, 16, 16, 16, 255, 201, 0, 11, 8,
					0, 1, 0, 1, 1, 1, 17, 0, 255, 204, 0, 6, 0, 16, 16, 5, 255, 218, 0, 8, 1, 1, 0, 0, 63, 0, 210, 207,
					32, 255, 217},
			},
			{
				ContentID: "inline-jpg-image.jpg@example.com",
				ContentType: ContentTypeHeader{
					ContentType: "image/jpeg",
					Params: map[string]string{
						"name": "inline-jpg-image-name.jpg",
					},
				},
				ContentDisposition: ContentDispositionHeader{
					ContentDisposition: inline,
					Params: map[string]string{
						"filename": "inline-jpg-image-filename.jpg",
					},
				},
				Data: []byte{255, 216, 255, 219, 0, 67, 0, 3, 2, 2, 2, 2, 2, 3, 2, 2, 2, 3, 3, 3, 3, 4,
					6, 4, 4, 4, 4, 4, 8, 6, 6, 5, 6, 9, 8, 10, 10, 9, 8, 9, 9, 10, 12, 15, 12, 10, 11, 14, 11, 9, 9,
					13, 17, 13, 14, 15, 16, 16, 17, 16, 10, 12, 18, 19, 18, 16, 19, 15, 16, 16, 16, 255, 201, 0, 11, 8,
					0, 1, 0, 1, 1, 1, 17, 0, 255, 204, 0, 6, 0, 16, 16, 5, 255, 218, 0, 8, 1, 1, 0, 0, 63, 0, 210, 207,
					32, 255, 217},
			},
		},
		AttachedFiles: []AttachedFile{
			{
				ContentType: ContentTypeHeader{
					ContentType: "application/pdf",
					Params: map[string]string{
						"name": "attached-pdf-name.pdf",
					},
				},
				ContentDisposition: ContentDispositionHeader{
					ContentDisposition: attachment,
					Params: map[string]string{
						"filename": "attached-pdf-filename.pdf",
					},
				},
				Data: []byte{37, 80, 68, 70, 45, 49, 46, 13, 116, 114, 97, 105, 108, 101, 114, 60, 60,
					47, 82, 111, 111, 116, 60, 60, 47, 80, 97, 103, 101, 115, 60, 60, 47, 75, 105, 100, 115, 91, 60,
					60, 47, 77, 101, 100, 105, 97, 66, 111, 120, 91, 48, 32, 48, 32, 51, 32, 51, 93, 62, 62, 93, 62,
					62, 62, 62, 62, 62},
			},
			{
				ContentType: ContentTypeHeader{
					ContentType: "application/pdf",
					Params: map[string]string{
						"name": "attached-pdf-without-disposition.pdf",
					},
				},
				ContentDisposition: ContentDispositionHeader{
					ContentDisposition: "",
					Params:             map[string]string(nil),
				},
				Data: []byte{37, 80, 68, 70, 45, 49, 46, 13, 116, 114, 97, 105, 108, 101, 114, 60, 60,
					47, 82, 111, 111, 116, 60, 60, 47, 80, 97, 103, 101, 115, 60, 60, 47, 75, 105, 100, 115, 91, 60,
					60, 47, 77, 101, 100, 105, 97, 66, 111, 120, 91, 48, 32, 48, 32, 51, 32, 51, 93, 62, 62, 93, 62,
					62, 62, 62, 62, 62},
			},
			{
				ContentType: ContentTypeHeader{
					ContentType: "application/json",
					Params: map[string]string{
						"name": "attached-json-name.json",
					},
				},
				ContentDisposition: ContentDispositionHeader{
					ContentDisposition: attachment,
					Params: map[string]string{
						"filename": "attached-json-filename.json",
					},
				},
				Data: []byte{123, 34, 102, 111, 111, 34, 58, 34, 98, 97, 114, 34, 125},
			},
			{
				ContentType: ContentTypeHeader{
					ContentType: "text/plain",
					Params: map[string]string{
						"name": "attached-text-plain-name.txt",
					},
				},
				ContentDisposition: ContentDispositionHeader{
					ContentDisposition: attachment,
					Params: map[string]string{
						"filename": "attached-text-plain-filename.txt",
					},
				},
				Data: []byte{84, 101, 120, 116, 47, 112, 108, 97, 105, 110, 32, 99, 111, 110, 116, 101, 110, 116, 32,
					97, 115, 32, 97, 110, 32, 97, 116, 116, 97, 99, 104, 101, 100, 32, 46, 116, 120, 116, 32, 102, 105,
					108, 101, 46},
			},
			{
				ContentType: ContentTypeHeader{
					ContentType: "text/html",
					Params: map[string]string{
						"name": "attached-text-html-name.html",
					},
				},
				ContentDisposition: ContentDispositionHeader{
					ContentDisposition: attachment,
					Params: map[string]string{
						"filename": "attached-text-html-filename.html",
					},
				},
				Data: []byte{84, 101, 120, 116, 47, 104, 116, 109, 108, 32, 99, 111, 110, 116, 101, 110, 116, 32, 97,
					115, 32, 97, 110, 32, 97, 116, 116, 97, 99, 104, 101, 100, 32, 46, 104, 116, 109, 108, 32, 102,
					105, 108, 101, 46},
			},
		},
	}

	testEmailFromFile(t, fp, expectedEmail)
}
func TestParseEmailPolishMultipartMixedUtf8OverQuotedprintable(t *testing.T) {
	fp := "tests/test_polish_multipart_mixed_utf-8_over_quoted-printable.txt"
	tz, _ := time.LoadLocation("Europe/Warsaw")
	expectedDate, _ := time.Parse(
		time.RFC1123Z+" (MST)",
		time.Date(2019, time.April, 1, 7, 55, 0, 0, tz).Format(time.RFC1123Z+" (MST)"))
	expectedEmail := Email{
		Headers: Headers{
			Date:    expectedDate,
			Subject: "📧 Test Polskie pangramy",
			ReplyTo: []*mail.Address{
				{
					Name:    "Nadająca, Alicja",
					Address: "alicja.nadajaca@example.net",
				},
			},
			Sender: &mail.Address{
				Name:    "Nadająca, Alicja",
				Address: "alicja.nadajaca@example.com",
			},
			From: []*mail.Address{
				{
					Name:    "Nadająca, Alicja",
					Address: "alicja.nadajaca@example.com",
				},
				{
					Name:    "Nadająca, Alicja",
					Address: "alicja.nadajaca@example.net",
				},
			},
			To: []*mail.Address{
				{
					Name:    "Odbierający, Bob",
					Address: "bob.odbierajacy@example.com",
				},
				{
					Name:    "Odbierająca, Karolina",
					Address: "karolina.odbierajaca@example.com",
				},
			},
			Cc: []*mail.Address{
				{
					Name:    "Odbierający, Daniel",
					Address: "daniel.odbierajacy@example.com",
				},
				{
					Name:    "Odbierająca, Ewa",
					Address: "ewa.odbierajaca@example.com",
				},
			},
			Bcc: []*mail.Address{
				{
					Name:    "Odbierający, Franek",
					Address: "franek.odbierajacy@example.com",
				},
				{
					Name:    "Odbierająca, Grażyna",
					Address: "grazyna.odbierajaca@example.com",
				},
			},
			MessageID:  "Message-Id-1@example.com",
			InReplyTo:  []MessageId{"Message-Id-0@example.com"},
			References: []MessageId{"Message-Id-0@example.com"},
			Comments:   "Message Header Comment",
			Keywords:   []string{"Keyword 1", "Keyword 2"},
			ResentDate: expectedDate,
			ResentFrom: []*mail.Address{
				{
					Name:    "Nadająca, Alicja",
					Address: "alicja.nadajaca@example.net",
				},
				{
					Name:    "Nadająca, Alicja",
					Address: "alicja.nadajaca@example.com",
				},
			},
			ResentSender: &mail.Address{
				Name:    "Nadająca, Alicja",
				Address: "alicja.nadajaca@example.net",
			},
			ResentTo: []*mail.Address{
				{
					Name:    "Odbierający, Bob",
					Address: "bob.odbierajacy@example.net",
				},
				{
					Name:    "Odbierająca, Karolina",
					Address: "karolina.odbierajaca@example.net",
				},
			},
			ResentCc: []*mail.Address{
				{
					Name:    "Odbierający, Daniel",
					Address: "daniel.odbierajacy@example.net",
				},
				{
					Name:    "Odbierająca, Ewa",
					Address: "ewa.odbierajaca@example.net",
				},
			},
			ResentBcc: []*mail.Address{
				{
					Name:    "Odbierający, Franek",
					Address: "franek.odbierajacy@example.net",
				},
				{
					Name:    "Odbierająca, Grażyna",
					Address: "grazyna.odbierajaca@example.net",
				},
			},
			ResentMessageID: "Message-Id-1@example.net",
			ContentType: ContentTypeHeader{
				ContentType: "multipart/mixed",
				Params: map[string]string{
					"boundary": "MixedBoundaryString",
					"charset":  "utf-8",
				},
			},
			ExtraHeaders: map[string][]string{
				"X-Clacks-Overhead": {"GNU Terry Pratchett"},
				"X-Script/function/\t !\"#$%&'()*+,-./;<=>?@[\\]^_`{|}~": {
					"TEST VALUE 1\t !\"#$%&'()*+,-./0123456789:;<=>?@ABCDEFGHIJKLMNOPQRSTUVWXYZ[\\]^_` abcdefghijklmnopqrstuvwxyz{|}~",
					"TEST VALUE 2\t !\"#$%&'()*+,-./0123456789:;<=>?@ABCDEFGHIJKLMNOPQRSTUVWXYZ[\\]^_` abcdefghijklmnopqrstuvwxyz{|}~",
				},
			},
		},
		Text: `Jeżu klątw, spłódź Finom część gry hańb!
Pójdźże, kiń tę chmurność w głąb flaszy!
Mężny bądź chroń pułk twój i sześć flag.
Filmuj rzeź żądań, pość, gnęb chłystków!
Pchnąć w tę łódź jeża lub ośm skrzyń fig.
Dość gróźb fuzją, klnę, pych i małżeństw!
Pójdź w loch zbić małżeńską gęś futryn!
Chwyć małżonkę, strój bądź pleśń z fugi.`,
		EnrichedText: `<bold>Jeżu klątw, spłódź Finom część gry hańb!</bold>
<italic>Pójdźże, kiń tę chmurność w głąb flaszy!</italic>
<fixed>Mężny bądź chroń pułk twój i sześć flag.</fixed>
<underline>Filmuj rzeź żądań, pość, gnęb chłystków!</underline>
Pchnąć w tę łódź jeża lub ośm skrzyń fig.
Dość gróźb fuzją, klnę, pych i małżeństw!
Pójdź w loch zbić małżeńską gęś futryn!
Chwyć małżonkę, strój bądź pleśń z fugi.`,
		HTML: `<html>
<div dir="ltr">
<p>Jeżu klątw, spłódź Finom część gry hańb!</p>
<p>Pójdźże, kiń tę chmurność w głąb flaszy!</p>
<p>Mężny bądź chroń pułk twój i sześć flag.</p>
<p>Filmuj rzeź żądań, pość, gnęb chłystków!</p>
<p>Pchnąć w tę łódź jeża lub ośm skrzyń fig.</p>
<p>Dość gróźb fuzją, klnę, pych i małżeństw!</p>
<p>Pójdź w loch zbić małżeńską gęś futryn!</p>
<p>Chwyć małżonkę, strój bądź pleśń z fugi.</p>
</div>
</html>`,
		InlineFiles: []InlineFile{
			{
				ContentType: ContentTypeHeader{
					ContentType: "image/jpeg",
					Params: map[string]string{
						"name": "inline-jpg-image-without-disposition.jpg",
					},
				},
				ContentDisposition: ContentDispositionHeader{
					ContentDisposition: "",
					Params:             map[string]string(nil),
				},
				Data: []byte{255, 216, 255, 219, 0, 67, 0, 3, 2, 2, 2, 2, 2, 3, 2, 2, 2, 3, 3, 3, 3, 4,
					6, 4, 4, 4, 4, 4, 8, 6, 6, 5, 6, 9, 8, 10, 10, 9, 8, 9, 9, 10, 12, 15, 12, 10, 11, 14, 11, 9, 9,
					13, 17, 13, 14, 15, 16, 16, 17, 16, 10, 12, 18, 19, 18, 16, 19, 15, 16, 16, 16, 255, 201, 0, 11, 8,
					0, 1, 0, 1, 1, 1, 17, 0, 255, 204, 0, 6, 0, 16, 16, 5, 255, 218, 0, 8, 1, 1, 0, 0, 63, 0, 210, 207,
					32, 255, 217},
			},
			{
				ContentID: "inline-jpg-image.jpg@example.com",
				ContentType: ContentTypeHeader{
					ContentType: "image/jpeg",
					Params: map[string]string{
						"name": "inline-jpg-image-name.jpg",
					},
				},
				ContentDisposition: ContentDispositionHeader{
					ContentDisposition: inline,
					Params: map[string]string{
						"filename": "inline-jpg-image-filename.jpg",
					},
				},
				Data: []byte{255, 216, 255, 219, 0, 67, 0, 3, 2, 2, 2, 2, 2, 3, 2, 2, 2, 3, 3, 3, 3, 4,
					6, 4, 4, 4, 4, 4, 8, 6, 6, 5, 6, 9, 8, 10, 10, 9, 8, 9, 9, 10, 12, 15, 12, 10, 11, 14, 11, 9, 9,
					13, 17, 13, 14, 15, 16, 16, 17, 16, 10, 12, 18, 19, 18, 16, 19, 15, 16, 16, 16, 255, 201, 0, 11, 8,
					0, 1, 0, 1, 1, 1, 17, 0, 255, 204, 0, 6, 0, 16, 16, 5, 255, 218, 0, 8, 1, 1, 0, 0, 63, 0, 210, 207,
					32, 255, 217},
			},
		},
		AttachedFiles: []AttachedFile{
			{
				ContentType: ContentTypeHeader{
					ContentType: "application/pdf",
					Params: map[string]string{
						"name": "attached-pdf-name.pdf",
					},
				},
				ContentDisposition: ContentDispositionHeader{
					ContentDisposition: attachment,
					Params: map[string]string{
						"filename": "attached-pdf-filename.pdf",
					},
				},
				Data: []byte{37, 80, 68, 70, 45, 49, 46, 13, 116, 114, 97, 105, 108, 101, 114, 60, 60,
					47, 82, 111, 111, 116, 60, 60, 47, 80, 97, 103, 101, 115, 60, 60, 47, 75, 105, 100, 115, 91, 60,
					60, 47, 77, 101, 100, 105, 97, 66, 111, 120, 91, 48, 32, 48, 32, 51, 32, 51, 93, 62, 62, 93, 62,
					62, 62, 62, 62, 62},
			},
			{
				ContentType: ContentTypeHeader{
					ContentType: "application/pdf",
					Params: map[string]string{
						"name": "attached-pdf-without-disposition.pdf",
					},
				},
				ContentDisposition: ContentDispositionHeader{
					ContentDisposition: "",
					Params:             map[string]string(nil),
				},
				Data: []byte{37, 80, 68, 70, 45, 49, 46, 13, 116, 114, 97, 105, 108, 101, 114, 60, 60,
					47, 82, 111, 111, 116, 60, 60, 47, 80, 97, 103, 101, 115, 60, 60, 47, 75, 105, 100, 115, 91, 60,
					60, 47, 77, 101, 100, 105, 97, 66, 111, 120, 91, 48, 32, 48, 32, 51, 32, 51, 93, 62, 62, 93, 62,
					62, 62, 62, 62, 62},
			},
			{
				ContentType: ContentTypeHeader{
					ContentType: "application/json",
					Params: map[string]string{
						"name": "attached-json-name.json",
					},
				},
				ContentDisposition: ContentDispositionHeader{
					ContentDisposition: attachment,
					Params: map[string]string{
						"filename": "attached-json-filename.json",
					},
				},
				Data: []byte{123, 34, 102, 111, 111, 34, 58, 34, 98, 97, 114, 34, 125},
			},
			{
				ContentType: ContentTypeHeader{
					ContentType: "text/plain",
					Params: map[string]string{
						"name": "attached-text-plain-name.txt",
					},
				},
				ContentDisposition: ContentDispositionHeader{
					ContentDisposition: attachment,
					Params: map[string]string{
						"filename": "attached-text-plain-filename.txt",
					},
				},
				Data: []byte{84, 101, 120, 116, 47, 112, 108, 97, 105, 110, 32, 99, 111, 110, 116, 101, 110, 116, 32,
					97, 115, 32, 97, 110, 32, 97, 116, 116, 97, 99, 104, 101, 100, 32, 46, 116, 120, 116, 32, 102, 105,
					108, 101, 46},
			},
			{
				ContentType: ContentTypeHeader{
					ContentType: "text/html",
					Params: map[string]string{
						"name": "attached-text-html-name.html",
					},
				},
				ContentDisposition: ContentDispositionHeader{
					ContentDisposition: attachment,
					Params: map[string]string{
						"filename": "attached-text-html-filename.html",
					},
				},
				Data: []byte{84, 101, 120, 116, 47, 104, 116, 109, 108, 32, 99, 111, 110, 116, 101, 110, 116, 32, 97,
					115, 32, 97, 110, 32, 97, 116, 116, 97, 99, 104, 101, 100, 32, 46, 104, 116, 109, 108, 32, 102,
					105, 108, 101, 46},
			},
		},
	}

	testEmailFromFile(t, fp, expectedEmail)
}
func TestParseEmailPolishMultipartMixedIso88592OverBase64(t *testing.T) {
	fp := "tests/test_polish_multipart_mixed_iso-8859-2_over_base64.txt"
	tz, _ := time.LoadLocation("Europe/Warsaw")
	expectedDate, _ := time.Parse(
		time.RFC1123Z+" (MST)",
		time.Date(2019, time.April, 1, 7, 55, 0, 0, tz).Format(time.RFC1123Z+" (MST)"))
	expectedEmail := Email{
		Headers: Headers{
			Date:    expectedDate,
			Subject: "📧 Test Polskie pangramy",
			ReplyTo: []*mail.Address{
				{
					Name:    "Nadająca, Alicja",
					Address: "alicja.nadajaca@example.net",
				},
			},
			Sender: &mail.Address{
				Name:    "Nadająca, Alicja",
				Address: "alicja.nadajaca@example.com",
			},
			From: []*mail.Address{
				{
					Name:    "Nadająca, Alicja",
					Address: "alicja.nadajaca@example.com",
				},
				{
					Name:    "Nadająca, Alicja",
					Address: "alicja.nadajaca@example.net",
				},
			},
			To: []*mail.Address{
				{
					Name:    "Odbierający, Bob",
					Address: "bob.odbierajacy@example.com",
				},
				{
					Name:    "Odbierająca, Karolina",
					Address: "karolina.odbierajaca@example.com",
				},
			},
			Cc: []*mail.Address{
				{
					Name:    "Odbierający, Daniel",
					Address: "daniel.odbierajacy@example.com",
				},
				{
					Name:    "Odbierająca, Ewa",
					Address: "ewa.odbierajaca@example.com",
				},
			},
			Bcc: []*mail.Address{
				{
					Name:    "Odbierający, Franek",
					Address: "franek.odbierajacy@example.com",
				},
				{
					Name:    "Odbierająca, Grażyna",
					Address: "grazyna.odbierajaca@example.com",
				},
			},
			MessageID:  "Message-Id-1@example.com",
			InReplyTo:  []MessageId{"Message-Id-0@example.com"},
			References: []MessageId{"Message-Id-0@example.com"},
			Comments:   "Message Header Comment",
			Keywords:   []string{"Keyword 1", "Keyword 2"},
			ResentDate: expectedDate,
			ResentFrom: []*mail.Address{
				{
					Name:    "Nadająca, Alicja",
					Address: "alicja.nadajaca@example.net",
				},
				{
					Name:    "Nadająca, Alicja",
					Address: "alicja.nadajaca@example.com",
				},
			},
			ResentSender: &mail.Address{
				Name:    "Nadająca, Alicja",
				Address: "alicja.nadajaca@example.net",
			},
			ResentTo: []*mail.Address{
				{
					Name:    "Odbierający, Bob",
					Address: "bob.odbierajacy@example.net",
				},
				{
					Name:    "Odbierająca, Karolina",
					Address: "karolina.odbierajaca@example.net",
				},
			},
			ResentCc: []*mail.Address{
				{
					Name:    "Odbierający, Daniel",
					Address: "daniel.odbierajacy@example.net",
				},
				{
					Name:    "Odbierająca, Ewa",
					Address: "ewa.odbierajaca@example.net",
				},
			},
			ResentBcc: []*mail.Address{
				{
					Name:    "Odbierający, Franek",
					Address: "franek.odbierajacy@example.net",
				},
				{
					Name:    "Odbierająca, Grażyna",
					Address: "grazyna.odbierajaca@example.net",
				},
			},
			ResentMessageID: "Message-Id-1@example.net",
			ContentType: ContentTypeHeader{
				ContentType: "multipart/mixed",
				Params: map[string]string{
					"boundary": "MixedBoundaryString",
					"charset":  "iso-8859-2",
				},
			},
			ExtraHeaders: map[string][]string{
				"X-Clacks-Overhead": {"GNU Terry Pratchett"},
				"X-Script/function/\t !\"#$%&'()*+,-./;<=>?@[\\]^_`{|}~": {
					"TEST VALUE 1\t !\"#$%&'()*+,-./0123456789:;<=>?@ABCDEFGHIJKLMNOPQRSTUVWXYZ[\\]^_` abcdefghijklmnopqrstuvwxyz{|}~",
					"TEST VALUE 2\t !\"#$%&'()*+,-./0123456789:;<=>?@ABCDEFGHIJKLMNOPQRSTUVWXYZ[\\]^_` abcdefghijklmnopqrstuvwxyz{|}~",
				},
			},
		},
		Text: `Jeżu klątw, spłódź Finom część gry hańb!
Pójdźże, kiń tę chmurność w głąb flaszy!
Mężny bądź chroń pułk twój i sześć flag.
Filmuj rzeź żądań, pość, gnęb chłystków!
Pchnąć w tę łódź jeża lub ośm skrzyń fig.
Dość gróźb fuzją, klnę, pych i małżeństw!
Pójdź w loch zbić małżeńską gęś futryn!
Chwyć małżonkę, strój bądź pleśń z fugi.`,
		EnrichedText: `<bold>Jeżu klątw, spłódź Finom część gry hańb!</bold>
<italic>Pójdźże, kiń tę chmurność w głąb flaszy!</italic>
<fixed>Mężny bądź chroń pułk twój i sześć flag.</fixed>
<underline>Filmuj rzeź żądań, pość, gnęb chłystków!</underline>
Pchnąć w tę łódź jeża lub ośm skrzyń fig.
Dość gróźb fuzją, klnę, pych i małżeństw!
Pójdź w loch zbić małżeńską gęś futryn!
Chwyć małżonkę, strój bądź pleśń z fugi.`,
		HTML: `<html>
<div dir="ltr">
<p>Jeżu klątw, spłódź Finom część gry hańb!</p>
<p>Pójdźże, kiń tę chmurność w głąb flaszy!</p>
<p>Mężny bądź chroń pułk twój i sześć flag.</p>
<p>Filmuj rzeź żądań, pość, gnęb chłystków!</p>
<p>Pchnąć w tę łódź jeża lub ośm skrzyń fig.</p>
<p>Dość gróźb fuzją, klnę, pych i małżeństw!</p>
<p>Pójdź w loch zbić małżeńską gęś futryn!</p>
<p>Chwyć małżonkę, strój bądź pleśń z fugi.</p>
</div>
</html>`,
		InlineFiles: []InlineFile{
			{
				ContentType: ContentTypeHeader{
					ContentType: "image/jpeg",
					Params: map[string]string{
						"name": "inline-jpg-image-without-disposition.jpg",
					},
				},
				ContentDisposition: ContentDispositionHeader{
					ContentDisposition: "",
					Params:             map[string]string(nil),
				},
				Data: []byte{255, 216, 255, 219, 0, 67, 0, 3, 2, 2, 2, 2, 2, 3, 2, 2, 2, 3, 3, 3, 3, 4,
					6, 4, 4, 4, 4, 4, 8, 6, 6, 5, 6, 9, 8, 10, 10, 9, 8, 9, 9, 10, 12, 15, 12, 10, 11, 14, 11, 9, 9,
					13, 17, 13, 14, 15, 16, 16, 17, 16, 10, 12, 18, 19, 18, 16, 19, 15, 16, 16, 16, 255, 201, 0, 11, 8,
					0, 1, 0, 1, 1, 1, 17, 0, 255, 204, 0, 6, 0, 16, 16, 5, 255, 218, 0, 8, 1, 1, 0, 0, 63, 0, 210, 207,
					32, 255, 217},
			},
			{
				ContentID: "inline-jpg-image.jpg@example.com",
				ContentType: ContentTypeHeader{
					ContentType: "image/jpeg",
					Params: map[string]string{
						"name": "inline-jpg-image-name.jpg",
					},
				},
				ContentDisposition: ContentDispositionHeader{
					ContentDisposition: inline,
					Params: map[string]string{
						"filename": "inline-jpg-image-filename.jpg",
					},
				},
				Data: []byte{255, 216, 255, 219, 0, 67, 0, 3, 2, 2, 2, 2, 2, 3, 2, 2, 2, 3, 3, 3, 3, 4,
					6, 4, 4, 4, 4, 4, 8, 6, 6, 5, 6, 9, 8, 10, 10, 9, 8, 9, 9, 10, 12, 15, 12, 10, 11, 14, 11, 9, 9,
					13, 17, 13, 14, 15, 16, 16, 17, 16, 10, 12, 18, 19, 18, 16, 19, 15, 16, 16, 16, 255, 201, 0, 11, 8,
					0, 1, 0, 1, 1, 1, 17, 0, 255, 204, 0, 6, 0, 16, 16, 5, 255, 218, 0, 8, 1, 1, 0, 0, 63, 0, 210, 207,
					32, 255, 217},
			},
		},
		AttachedFiles: []AttachedFile{
			{
				ContentType: ContentTypeHeader{
					ContentType: "application/pdf",
					Params: map[string]string{
						"name": "attached-pdf-name.pdf",
					},
				},
				ContentDisposition: ContentDispositionHeader{
					ContentDisposition: attachment,
					Params: map[string]string{
						"filename": "attached-pdf-filename.pdf",
					},
				},
				Data: []byte{37, 80, 68, 70, 45, 49, 46, 13, 116, 114, 97, 105, 108, 101, 114, 60, 60,
					47, 82, 111, 111, 116, 60, 60, 47, 80, 97, 103, 101, 115, 60, 60, 47, 75, 105, 100, 115, 91, 60,
					60, 47, 77, 101, 100, 105, 97, 66, 111, 120, 91, 48, 32, 48, 32, 51, 32, 51, 93, 62, 62, 93, 62,
					62, 62, 62, 62, 62},
			},
			{
				ContentType: ContentTypeHeader{
					ContentType: "application/pdf",
					Params: map[string]string{
						"name": "attached-pdf-without-disposition.pdf",
					},
				},
				ContentDisposition: ContentDispositionHeader{
					ContentDisposition: "",
					Params:             map[string]string(nil),
				},
				Data: []byte{37, 80, 68, 70, 45, 49, 46, 13, 116, 114, 97, 105, 108, 101, 114, 60, 60,
					47, 82, 111, 111, 116, 60, 60, 47, 80, 97, 103, 101, 115, 60, 60, 47, 75, 105, 100, 115, 91, 60,
					60, 47, 77, 101, 100, 105, 97, 66, 111, 120, 91, 48, 32, 48, 32, 51, 32, 51, 93, 62, 62, 93, 62,
					62, 62, 62, 62, 62},
			},
			{
				ContentType: ContentTypeHeader{
					ContentType: "application/json",
					Params: map[string]string{
						"name": "attached-json-name.json",
					},
				},
				ContentDisposition: ContentDispositionHeader{
					ContentDisposition: attachment,
					Params: map[string]string{
						"filename": "attached-json-filename.json",
					},
				},
				Data: []byte{123, 34, 102, 111, 111, 34, 58, 34, 98, 97, 114, 34, 125},
			},
			{
				ContentType: ContentTypeHeader{
					ContentType: "text/plain",
					Params: map[string]string{
						"name": "attached-text-plain-name.txt",
					},
				},
				ContentDisposition: ContentDispositionHeader{
					ContentDisposition: attachment,
					Params: map[string]string{
						"filename": "attached-text-plain-filename.txt",
					},
				},
				Data: []byte{84, 101, 120, 116, 47, 112, 108, 97, 105, 110, 32, 99, 111, 110, 116, 101, 110, 116, 32,
					97, 115, 32, 97, 110, 32, 97, 116, 116, 97, 99, 104, 101, 100, 32, 46, 116, 120, 116, 32, 102, 105,
					108, 101, 46},
			},
			{
				ContentType: ContentTypeHeader{
					ContentType: "text/html",
					Params: map[string]string{
						"name": "attached-text-html-name.html",
					},
				},
				ContentDisposition: ContentDispositionHeader{
					ContentDisposition: attachment,
					Params: map[string]string{
						"filename": "attached-text-html-filename.html",
					},
				},
				Data: []byte{84, 101, 120, 116, 47, 104, 116, 109, 108, 32, 99, 111, 110, 116, 101, 110, 116, 32, 97,
					115, 32, 97, 110, 32, 97, 116, 116, 97, 99, 104, 101, 100, 32, 46, 104, 116, 109, 108, 32, 102,
					105, 108, 101, 46},
			},
		},
	}

	testEmailFromFile(t, fp, expectedEmail)
}
func TestParseEmailPolishMultipartMixedIso88592OverQuotedprintable(t *testing.T) {
	fp := "tests/test_polish_multipart_mixed_iso-8859-2_over_quoted-printable.txt"
	tz, _ := time.LoadLocation("Europe/Warsaw")
	expectedDate, _ := time.Parse(
		time.RFC1123Z+" (MST)",
		time.Date(2019, time.April, 1, 7, 55, 0, 0, tz).Format(time.RFC1123Z+" (MST)"))
	expectedEmail := Email{
		Headers: Headers{
			Date:    expectedDate,
			Subject: "📧 Test Polskie pangramy",
			ReplyTo: []*mail.Address{
				{
					Name:    "Nadająca, Alicja",
					Address: "alicja.nadajaca@example.net",
				},
			},
			Sender: &mail.Address{
				Name:    "Nadająca, Alicja",
				Address: "alicja.nadajaca@example.com",
			},
			From: []*mail.Address{
				{
					Name:    "Nadająca, Alicja",
					Address: "alicja.nadajaca@example.com",
				},
				{
					Name:    "Nadająca, Alicja",
					Address: "alicja.nadajaca@example.net",
				},
			},
			To: []*mail.Address{
				{
					Name:    "Odbierający, Bob",
					Address: "bob.odbierajacy@example.com",
				},
				{
					Name:    "Odbierająca, Karolina",
					Address: "karolina.odbierajaca@example.com",
				},
			},
			Cc: []*mail.Address{
				{
					Name:    "Odbierający, Daniel",
					Address: "daniel.odbierajacy@example.com",
				},
				{
					Name:    "Odbierająca, Ewa",
					Address: "ewa.odbierajaca@example.com",
				},
			},
			Bcc: []*mail.Address{
				{
					Name:    "Odbierający, Franek",
					Address: "franek.odbierajacy@example.com",
				},
				{
					Name:    "Odbierająca, Grażyna",
					Address: "grazyna.odbierajaca@example.com",
				},
			},
			MessageID:  "Message-Id-1@example.com",
			InReplyTo:  []MessageId{"Message-Id-0@example.com"},
			References: []MessageId{"Message-Id-0@example.com"},
			Comments:   "Message Header Comment",
			Keywords:   []string{"Keyword 1", "Keyword 2"},
			ResentDate: expectedDate,
			ResentFrom: []*mail.Address{
				{
					Name:    "Nadająca, Alicja",
					Address: "alicja.nadajaca@example.net",
				},
				{
					Name:    "Nadająca, Alicja",
					Address: "alicja.nadajaca@example.com",
				},
			},
			ResentSender: &mail.Address{
				Name:    "Nadająca, Alicja",
				Address: "alicja.nadajaca@example.net",
			},
			ResentTo: []*mail.Address{
				{
					Name:    "Odbierający, Bob",
					Address: "bob.odbierajacy@example.net",
				},
				{
					Name:    "Odbierająca, Karolina",
					Address: "karolina.odbierajaca@example.net",
				},
			},
			ResentCc: []*mail.Address{
				{
					Name:    "Odbierający, Daniel",
					Address: "daniel.odbierajacy@example.net",
				},
				{
					Name:    "Odbierająca, Ewa",
					Address: "ewa.odbierajaca@example.net",
				},
			},
			ResentBcc: []*mail.Address{
				{
					Name:    "Odbierający, Franek",
					Address: "franek.odbierajacy@example.net",
				},
				{
					Name:    "Odbierająca, Grażyna",
					Address: "grazyna.odbierajaca@example.net",
				},
			},
			ResentMessageID: "Message-Id-1@example.net",
			ContentType: ContentTypeHeader{
				ContentType: "multipart/mixed",
				Params: map[string]string{
					"boundary": "MixedBoundaryString",
					"charset":  "iso-8859-2",
				},
			},
			ExtraHeaders: map[string][]string{
				"X-Clacks-Overhead": {"GNU Terry Pratchett"},
				"X-Script/function/\t !\"#$%&'()*+,-./;<=>?@[\\]^_`{|}~": {
					"TEST VALUE 1\t !\"#$%&'()*+,-./0123456789:;<=>?@ABCDEFGHIJKLMNOPQRSTUVWXYZ[\\]^_` abcdefghijklmnopqrstuvwxyz{|}~",
					"TEST VALUE 2\t !\"#$%&'()*+,-./0123456789:;<=>?@ABCDEFGHIJKLMNOPQRSTUVWXYZ[\\]^_` abcdefghijklmnopqrstuvwxyz{|}~",
				},
			},
		},
		Text: `Jeżu klątw, spłódź Finom część gry hańb!
Pójdźże, kiń tę chmurność w głąb flaszy!
Mężny bądź chroń pułk twój i sześć flag.
Filmuj rzeź żądań, pość, gnęb chłystków!
Pchnąć w tę łódź jeża lub ośm skrzyń fig.
Dość gróźb fuzją, klnę, pych i małżeństw!
Pójdź w loch zbić małżeńską gęś futryn!
Chwyć małżonkę, strój bądź pleśń z fugi.`,
		EnrichedText: `<bold>Jeżu klątw, spłódź Finom część gry hańb!</bold>
<italic>Pójdźże, kiń tę chmurność w głąb flaszy!</italic>
<fixed>Mężny bądź chroń pułk twój i sześć flag.</fixed>
<underline>Filmuj rzeź żądań, pość, gnęb chłystków!</underline>
Pchnąć w tę łódź jeża lub ośm skrzyń fig.
Dość gróźb fuzją, klnę, pych i małżeństw!
Pójdź w loch zbić małżeńską gęś futryn!
Chwyć małżonkę, strój bądź pleśń z fugi.`,
		HTML: `<html>
<div dir="ltr">
<p>Jeżu klątw, spłódź Finom część gry hańb!</p>
<p>Pójdźże, kiń tę chmurność w głąb flaszy!</p>
<p>Mężny bądź chroń pułk twój i sześć flag.</p>
<p>Filmuj rzeź żądań, pość, gnęb chłystków!</p>
<p>Pchnąć w tę łódź jeża lub ośm skrzyń fig.</p>
<p>Dość gróźb fuzją, klnę, pych i małżeństw!</p>
<p>Pójdź w loch zbić małżeńską gęś futryn!</p>
<p>Chwyć małżonkę, strój bądź pleśń z fugi.</p>
</div>
</html>`,
		InlineFiles: []InlineFile{
			{
				ContentType: ContentTypeHeader{
					ContentType: "image/jpeg",
					Params: map[string]string{
						"name": "inline-jpg-image-without-disposition.jpg",
					},
				},
				ContentDisposition: ContentDispositionHeader{
					ContentDisposition: "",
					Params:             map[string]string(nil),
				},
				Data: []byte{255, 216, 255, 219, 0, 67, 0, 3, 2, 2, 2, 2, 2, 3, 2, 2, 2, 3, 3, 3, 3, 4,
					6, 4, 4, 4, 4, 4, 8, 6, 6, 5, 6, 9, 8, 10, 10, 9, 8, 9, 9, 10, 12, 15, 12, 10, 11, 14, 11, 9, 9,
					13, 17, 13, 14, 15, 16, 16, 17, 16, 10, 12, 18, 19, 18, 16, 19, 15, 16, 16, 16, 255, 201, 0, 11, 8,
					0, 1, 0, 1, 1, 1, 17, 0, 255, 204, 0, 6, 0, 16, 16, 5, 255, 218, 0, 8, 1, 1, 0, 0, 63, 0, 210, 207,
					32, 255, 217},
			},
			{
				ContentID: "inline-jpg-image.jpg@example.com",
				ContentType: ContentTypeHeader{
					ContentType: "image/jpeg",
					Params: map[string]string{
						"name": "inline-jpg-image-name.jpg",
					},
				},
				ContentDisposition: ContentDispositionHeader{
					ContentDisposition: inline,
					Params: map[string]string{
						"filename": "inline-jpg-image-filename.jpg",
					},
				},
				Data: []byte{255, 216, 255, 219, 0, 67, 0, 3, 2, 2, 2, 2, 2, 3, 2, 2, 2, 3, 3, 3, 3, 4,
					6, 4, 4, 4, 4, 4, 8, 6, 6, 5, 6, 9, 8, 10, 10, 9, 8, 9, 9, 10, 12, 15, 12, 10, 11, 14, 11, 9, 9,
					13, 17, 13, 14, 15, 16, 16, 17, 16, 10, 12, 18, 19, 18, 16, 19, 15, 16, 16, 16, 255, 201, 0, 11, 8,
					0, 1, 0, 1, 1, 1, 17, 0, 255, 204, 0, 6, 0, 16, 16, 5, 255, 218, 0, 8, 1, 1, 0, 0, 63, 0, 210, 207,
					32, 255, 217},
			},
		},
		AttachedFiles: []AttachedFile{
			{
				ContentType: ContentTypeHeader{
					ContentType: "application/pdf",
					Params: map[string]string{
						"name": "attached-pdf-name.pdf",
					},
				},
				ContentDisposition: ContentDispositionHeader{
					ContentDisposition: attachment,
					Params: map[string]string{
						"filename": "attached-pdf-filename.pdf",
					},
				},
				Data: []byte{37, 80, 68, 70, 45, 49, 46, 13, 116, 114, 97, 105, 108, 101, 114, 60, 60,
					47, 82, 111, 111, 116, 60, 60, 47, 80, 97, 103, 101, 115, 60, 60, 47, 75, 105, 100, 115, 91, 60,
					60, 47, 77, 101, 100, 105, 97, 66, 111, 120, 91, 48, 32, 48, 32, 51, 32, 51, 93, 62, 62, 93, 62,
					62, 62, 62, 62, 62},
			},
			{
				ContentType: ContentTypeHeader{
					ContentType: "application/pdf",
					Params: map[string]string{
						"name": "attached-pdf-without-disposition.pdf",
					},
				},
				ContentDisposition: ContentDispositionHeader{
					ContentDisposition: "",
					Params:             map[string]string(nil),
				},
				Data: []byte{37, 80, 68, 70, 45, 49, 46, 13, 116, 114, 97, 105, 108, 101, 114, 60, 60,
					47, 82, 111, 111, 116, 60, 60, 47, 80, 97, 103, 101, 115, 60, 60, 47, 75, 105, 100, 115, 91, 60,
					60, 47, 77, 101, 100, 105, 97, 66, 111, 120, 91, 48, 32, 48, 32, 51, 32, 51, 93, 62, 62, 93, 62,
					62, 62, 62, 62, 62},
			},
			{
				ContentType: ContentTypeHeader{
					ContentType: "application/json",
					Params: map[string]string{
						"name": "attached-json-name.json",
					},
				},
				ContentDisposition: ContentDispositionHeader{
					ContentDisposition: attachment,
					Params: map[string]string{
						"filename": "attached-json-filename.json",
					},
				},
				Data: []byte{123, 34, 102, 111, 111, 34, 58, 34, 98, 97, 114, 34, 125},
			},
			{
				ContentType: ContentTypeHeader{
					ContentType: "text/plain",
					Params: map[string]string{
						"name": "attached-text-plain-name.txt",
					},
				},
				ContentDisposition: ContentDispositionHeader{
					ContentDisposition: attachment,
					Params: map[string]string{
						"filename": "attached-text-plain-filename.txt",
					},
				},
				Data: []byte{84, 101, 120, 116, 47, 112, 108, 97, 105, 110, 32, 99, 111, 110, 116, 101, 110, 116, 32,
					97, 115, 32, 97, 110, 32, 97, 116, 116, 97, 99, 104, 101, 100, 32, 46, 116, 120, 116, 32, 102, 105,
					108, 101, 46},
			},
			{
				ContentType: ContentTypeHeader{
					ContentType: "text/html",
					Params: map[string]string{
						"name": "attached-text-html-name.html",
					},
				},
				ContentDisposition: ContentDispositionHeader{
					ContentDisposition: attachment,
					Params: map[string]string{
						"filename": "attached-text-html-filename.html",
					},
				},
				Data: []byte{84, 101, 120, 116, 47, 104, 116, 109, 108, 32, 99, 111, 110, 116, 101, 110, 116, 32, 97,
					115, 32, 97, 110, 32, 97, 116, 116, 97, 99, 104, 101, 100, 32, 46, 104, 116, 109, 108, 32, 102,
					105, 108, 101, 46},
			},
		},
	}

	testEmailFromFile(t, fp, expectedEmail)
}
func TestParseEmailPolishMultipartSignedUtf8OverBase64(t *testing.T) {
	fp := "tests/test_polish_multipart_signed_utf-8_over_base64.txt"
	tz, _ := time.LoadLocation("Europe/Warsaw")
	expectedDate, _ := time.Parse(
		time.RFC1123Z+" (MST)",
		time.Date(2019, time.April, 1, 7, 55, 0, 0, tz).Format(time.RFC1123Z+" (MST)"))
	expectedEmail := Email{
		Headers: Headers{
			Date:    expectedDate,
			Subject: "📧 Signed Test Polskie pangramy",
			ReplyTo: []*mail.Address{
				{
					Name:    "Nadająca, Alicja",
					Address: "alicja.nadajaca@example.net",
				},
			},
			Sender: &mail.Address{
				Name:    "Nadająca, Alicja",
				Address: "alicja.nadajaca@example.com",
			},
			From: []*mail.Address{
				{
					Name:    "Nadająca, Alicja",
					Address: "alicja.nadajaca@example.com",
				},
				{
					Name:    "Nadająca, Alicja",
					Address: "alicja.nadajaca@example.net",
				},
			},
			To: []*mail.Address{
				{
					Name:    "Odbierający, Bob",
					Address: "bob.odbierajacy@example.com",
				},
				{
					Name:    "Odbierająca, Karolina",
					Address: "karolina.odbierajaca@example.com",
				},
			},
			Cc: []*mail.Address{
				{
					Name:    "Odbierający, Daniel",
					Address: "daniel.odbierajacy@example.com",
				},
				{
					Name:    "Odbierająca, Ewa",
					Address: "ewa.odbierajaca@example.com",
				},
			},
			Bcc: []*mail.Address{
				{
					Name:    "Odbierający, Franek",
					Address: "franek.odbierajacy@example.com",
				},
				{
					Name:    "Odbierająca, Grażyna",
					Address: "grazyna.odbierajaca@example.com",
				},
			},
			MessageID:  "Message-Id-1@example.com",
			InReplyTo:  []MessageId{"Message-Id-0@example.com"},
			References: []MessageId{"Message-Id-0@example.com"},
			Comments:   "Message Header Comment",
			Keywords:   []string{"Keyword 1", "Keyword 2"},
			ResentDate: expectedDate,
			ResentFrom: []*mail.Address{
				{
					Name:    "Nadająca, Alicja",
					Address: "alicja.nadajaca@example.net",
				},
				{
					Name:    "Nadająca, Alicja",
					Address: "alicja.nadajaca@example.com",
				},
			},
			ResentSender: &mail.Address{
				Name:    "Nadająca, Alicja",
				Address: "alicja.nadajaca@example.net",
			},
			ResentTo: []*mail.Address{
				{
					Name:    "Odbierający, Bob",
					Address: "bob.odbierajacy@example.net",
				},
				{
					Name:    "Odbierająca, Karolina",
					Address: "karolina.odbierajaca@example.net",
				},
			},
			ResentCc: []*mail.Address{
				{
					Name:    "Odbierający, Daniel",
					Address: "daniel.odbierajacy@example.net",
				},
				{
					Name:    "Odbierająca, Ewa",
					Address: "ewa.odbierajaca@example.net",
				},
			},
			ResentBcc: []*mail.Address{
				{
					Name:    "Odbierający, Franek",
					Address: "franek.odbierajacy@example.net",
				},
				{
					Name:    "Odbierająca, Grażyna",
					Address: "grazyna.odbierajaca@example.net",
				},
			},
			ResentMessageID: "Message-Id-1@example.net",
			ContentType: ContentTypeHeader{
				ContentType: "multipart/signed",
				Params: map[string]string{
					"boundary": "SignedBoundaryString",
					"charset":  "utf-8",
					"micalg":   "sha1",
					"protocol": "application/pkcs7-signature",
				},
			},
			ExtraHeaders: map[string][]string{
				"X-Clacks-Overhead": {"GNU Terry Pratchett"},
				"X-Script/function/\t !\"#$%&'()*+,-./;<=>?@[\\]^_`{|}~": {
					"TEST VALUE 1\t !\"#$%&'()*+,-./0123456789:;<=>?@ABCDEFGHIJKLMNOPQRSTUVWXYZ[\\]^_` abcdefghijklmnopqrstuvwxyz{|}~",
					"TEST VALUE 2\t !\"#$%&'()*+,-./0123456789:;<=>?@ABCDEFGHIJKLMNOPQRSTUVWXYZ[\\]^_` abcdefghijklmnopqrstuvwxyz{|}~",
				},
			},
		},
		Text: `Jeżu klątw, spłódź Finom część gry hańb!
Pójdźże, kiń tę chmurność w głąb flaszy!
Mężny bądź chroń pułk twój i sześć flag.
Filmuj rzeź żądań, pość, gnęb chłystków!
Pchnąć w tę łódź jeża lub ośm skrzyń fig.
Dość gróźb fuzją, klnę, pych i małżeństw!
Pójdź w loch zbić małżeńską gęś futryn!
Chwyć małżonkę, strój bądź pleśń z fugi.`,
		AttachedFiles: []AttachedFile{
			{
				ContentType: ContentTypeHeader{
					ContentType: "application/pkcs7-signature",
					Params: map[string]string{
						"name": "smime.p7s",
					},
				},
				ContentDisposition: ContentDispositionHeader{
					ContentDisposition: attachment,
					Params: map[string]string{
						"filename": "smime.p7s",
					},
				},
				Data: []byte{130, 28, 135, 132, 117, 46, 142, 18, 97, 140, 126, 251, 159, 193, 199, 25, 58, 223, 189, 185, 227, 239, 158, 173, 108, 31, 71, 27, 133, 80, 165, 252, 133, 227, 174, 198, 132, 129, 159, 29, 246, 19, 235, 133, 80, 165, 252, 133, 227, 174, 198, 132, 129, 159, 29, 246, 19, 234, 49, 251, 238, 127, 7, 28, 104, 33, 200, 120, 71, 82, 232, 225, 38, 30, 249, 234, 214, 193, 244, 113, 147, 173, 251, 219, 158, 57, 252, 28, 113, 147, 173, 251, 225, 38, 24, 199, 239, 190, 173, 108, 31, 71, 27, 133, 80, 110, 120, 251, 231, 174, 198, 132, 129, 159, 29, 246, 19, 234, 8, 114, 30, 17, 212, 186, 58, 95, 200, 94, 59, 26, 18, 6, 124, 119, 216, 79, 174, 21, 65, 185, 227, 239, 158},
			},
		},
	}

	testEmailFromFile(t, fp, expectedEmail)
}

func TestParseEmailPolishMultipartSignedUtf8OverQuotedprintable(t *testing.T) {
	fp := "tests/test_polish_multipart_signed_utf-8_over_quoted-printable.txt"
	tz, _ := time.LoadLocation("Europe/Warsaw")
	expectedDate, _ := time.Parse(
		time.RFC1123Z+" (MST)",
		time.Date(2019, time.April, 1, 7, 55, 0, 0, tz).Format(time.RFC1123Z+" (MST)"))
	expectedEmail := Email{
		Headers: Headers{
			Date:    expectedDate,
			Subject: "📧 Signed Test Polskie pangramy",
			ReplyTo: []*mail.Address{
				{
					Name:    "Nadająca, Alicja",
					Address: "alicja.nadajaca@example.net",
				},
			},
			Sender: &mail.Address{
				Name:    "Nadająca, Alicja",
				Address: "alicja.nadajaca@example.com",
			},
			From: []*mail.Address{
				{
					Name:    "Nadająca, Alicja",
					Address: "alicja.nadajaca@example.com",
				},
				{
					Name:    "Nadająca, Alicja",
					Address: "alicja.nadajaca@example.net",
				},
			},
			To: []*mail.Address{
				{
					Name:    "Odbierający, Bob",
					Address: "bob.odbierajacy@example.com",
				},
				{
					Name:    "Odbierająca, Karolina",
					Address: "karolina.odbierajaca@example.com",
				},
			},
			Cc: []*mail.Address{
				{
					Name:    "Odbierający, Daniel",
					Address: "daniel.odbierajacy@example.com",
				},
				{
					Name:    "Odbierająca, Ewa",
					Address: "ewa.odbierajaca@example.com",
				},
			},
			Bcc: []*mail.Address{
				{
					Name:    "Odbierający, Franek",
					Address: "franek.odbierajacy@example.com",
				},
				{
					Name:    "Odbierająca, Grażyna",
					Address: "grazyna.odbierajaca@example.com",
				},
			},
			MessageID:  "Message-Id-1@example.com",
			InReplyTo:  []MessageId{"Message-Id-0@example.com"},
			References: []MessageId{"Message-Id-0@example.com"},
			Comments:   "Message Header Comment",
			Keywords:   []string{"Keyword 1", "Keyword 2"},
			ResentDate: expectedDate,
			ResentFrom: []*mail.Address{
				{
					Name:    "Nadająca, Alicja",
					Address: "alicja.nadajaca@example.net",
				},
				{
					Name:    "Nadająca, Alicja",
					Address: "alicja.nadajaca@example.com",
				},
			},
			ResentSender: &mail.Address{
				Name:    "Nadająca, Alicja",
				Address: "alicja.nadajaca@example.net",
			},
			ResentTo: []*mail.Address{
				{
					Name:    "Odbierający, Bob",
					Address: "bob.odbierajacy@example.net",
				},
				{
					Name:    "Odbierająca, Karolina",
					Address: "karolina.odbierajaca@example.net",
				},
			},
			ResentCc: []*mail.Address{
				{
					Name:    "Odbierający, Daniel",
					Address: "daniel.odbierajacy@example.net",
				},
				{
					Name:    "Odbierająca, Ewa",
					Address: "ewa.odbierajaca@example.net",
				},
			},
			ResentBcc: []*mail.Address{
				{
					Name:    "Odbierający, Franek",
					Address: "franek.odbierajacy@example.net",
				},
				{
					Name:    "Odbierająca, Grażyna",
					Address: "grazyna.odbierajaca@example.net",
				},
			},
			ResentMessageID: "Message-Id-1@example.net",
			ContentType: ContentTypeHeader{
				ContentType: "multipart/signed",
				Params: map[string]string{
					"boundary": "SignedBoundaryString",
					"charset":  "utf-8",
					"micalg":   "sha1",
					"protocol": "application/pkcs7-signature",
				},
			},
			ExtraHeaders: map[string][]string{
				"X-Clacks-Overhead": {"GNU Terry Pratchett"},
				"X-Script/function/\t !\"#$%&'()*+,-./;<=>?@[\\]^_`{|}~": {
					"TEST VALUE 1\t !\"#$%&'()*+,-./0123456789:;<=>?@ABCDEFGHIJKLMNOPQRSTUVWXYZ[\\]^_` abcdefghijklmnopqrstuvwxyz{|}~",
					"TEST VALUE 2\t !\"#$%&'()*+,-./0123456789:;<=>?@ABCDEFGHIJKLMNOPQRSTUVWXYZ[\\]^_` abcdefghijklmnopqrstuvwxyz{|}~",
				},
			},
		},
		Text: `Jeżu klątw, spłódź Finom część gry hańb!
Pójdźże, kiń tę chmurność w głąb flaszy!
Mężny bądź chroń pułk twój i sześć flag.
Filmuj rzeź żądań, pość, gnęb chłystków!
Pchnąć w tę łódź jeża lub ośm skrzyń fig.
Dość gróźb fuzją, klnę, pych i małżeństw!
Pójdź w loch zbić małżeńską gęś futryn!
Chwyć małżonkę, strój bądź pleśń z fugi.`,
		AttachedFiles: []AttachedFile{
			{
				ContentType: ContentTypeHeader{
					ContentType: "application/pkcs7-signature",
					Params: map[string]string{
						"name": "smime.p7s",
					},
				},
				ContentDisposition: ContentDispositionHeader{
					ContentDisposition: attachment,
					Params: map[string]string{
						"filename": "smime.p7s",
					},
				},
				Data: []byte{130, 28, 135, 132, 117, 46, 142, 18, 97, 140, 126, 251, 159, 193, 199, 25, 58, 223, 189, 185, 227, 239, 158, 173, 108, 31, 71, 27, 133, 80, 165, 252, 133, 227, 174, 198, 132, 129, 159, 29, 246, 19, 235, 133, 80, 165, 252, 133, 227, 174, 198, 132, 129, 159, 29, 246, 19, 234, 49, 251, 238, 127, 7, 28, 104, 33, 200, 120, 71, 82, 232, 225, 38, 30, 249, 234, 214, 193, 244, 113, 147, 173, 251, 219, 158, 57, 252, 28, 113, 147, 173, 251, 225, 38, 24, 199, 239, 190, 173, 108, 31, 71, 27, 133, 80, 110, 120, 251, 231, 174, 198, 132, 129, 159, 29, 246, 19, 234, 8, 114, 30, 17, 212, 186, 58, 95, 200, 94, 59, 26, 18, 6, 124, 119, 216, 79, 174, 21, 65, 185, 227, 239, 158},
			},
		},
	}

	testEmailFromFile(t, fp, expectedEmail)
}

func TestParseEmailPolishMultipartSignedIso88592OverBase64(t *testing.T) {
	fp := "tests/test_polish_multipart_signed_iso-8859-2_over_base64.txt"
	tz, _ := time.LoadLocation("Europe/Warsaw")
	expectedDate, _ := time.Parse(
		time.RFC1123Z+" (MST)",
		time.Date(2019, time.April, 1, 7, 55, 0, 0, tz).Format(time.RFC1123Z+" (MST)"))
	expectedEmail := Email{
		Headers: Headers{
			Date:    expectedDate,
			Subject: "📧 Signed Test Polskie pangramy",
			ReplyTo: []*mail.Address{
				{
					Name:    "Nadająca, Alicja",
					Address: "alicja.nadajaca@example.net",
				},
			},
			Sender: &mail.Address{
				Name:    "Nadająca, Alicja",
				Address: "alicja.nadajaca@example.com",
			},
			From: []*mail.Address{
				{
					Name:    "Nadająca, Alicja",
					Address: "alicja.nadajaca@example.com",
				},
				{
					Name:    "Nadająca, Alicja",
					Address: "alicja.nadajaca@example.net",
				},
			},
			To: []*mail.Address{
				{
					Name:    "Odbierający, Bob",
					Address: "bob.odbierajacy@example.com",
				},
				{
					Name:    "Odbierająca, Karolina",
					Address: "karolina.odbierajaca@example.com",
				},
			},
			Cc: []*mail.Address{
				{
					Name:    "Odbierający, Daniel",
					Address: "daniel.odbierajacy@example.com",
				},
				{
					Name:    "Odbierająca, Ewa",
					Address: "ewa.odbierajaca@example.com",
				},
			},
			Bcc: []*mail.Address{
				{
					Name:    "Odbierający, Franek",
					Address: "franek.odbierajacy@example.com",
				},
				{
					Name:    "Odbierająca, Grażyna",
					Address: "grazyna.odbierajaca@example.com",
				},
			},
			MessageID:  "Message-Id-1@example.com",
			InReplyTo:  []MessageId{"Message-Id-0@example.com"},
			References: []MessageId{"Message-Id-0@example.com"},
			Comments:   "Message Header Comment",
			Keywords:   []string{"Keyword 1", "Keyword 2"},
			ResentDate: expectedDate,
			ResentFrom: []*mail.Address{
				{
					Name:    "Nadająca, Alicja",
					Address: "alicja.nadajaca@example.net",
				},
				{
					Name:    "Nadająca, Alicja",
					Address: "alicja.nadajaca@example.com",
				},
			},
			ResentSender: &mail.Address{
				Name:    "Nadająca, Alicja",
				Address: "alicja.nadajaca@example.net",
			},
			ResentTo: []*mail.Address{
				{
					Name:    "Odbierający, Bob",
					Address: "bob.odbierajacy@example.net",
				},
				{
					Name:    "Odbierająca, Karolina",
					Address: "karolina.odbierajaca@example.net",
				},
			},
			ResentCc: []*mail.Address{
				{
					Name:    "Odbierający, Daniel",
					Address: "daniel.odbierajacy@example.net",
				},
				{
					Name:    "Odbierająca, Ewa",
					Address: "ewa.odbierajaca@example.net",
				},
			},
			ResentBcc: []*mail.Address{
				{
					Name:    "Odbierający, Franek",
					Address: "franek.odbierajacy@example.net",
				},
				{
					Name:    "Odbierająca, Grażyna",
					Address: "grazyna.odbierajaca@example.net",
				},
			},
			ResentMessageID: "Message-Id-1@example.net",
			ContentType: ContentTypeHeader{
				ContentType: "multipart/signed",
				Params: map[string]string{
					"boundary": "SignedBoundaryString",
					"charset":  "iso-8859-2",
					"micalg":   "sha1",
					"protocol": "application/pkcs7-signature",
				},
			},
			ExtraHeaders: map[string][]string{
				"X-Clacks-Overhead": {"GNU Terry Pratchett"},
				"X-Script/function/\t !\"#$%&'()*+,-./;<=>?@[\\]^_`{|}~": {
					"TEST VALUE 1\t !\"#$%&'()*+,-./0123456789:;<=>?@ABCDEFGHIJKLMNOPQRSTUVWXYZ[\\]^_` abcdefghijklmnopqrstuvwxyz{|}~",
					"TEST VALUE 2\t !\"#$%&'()*+,-./0123456789:;<=>?@ABCDEFGHIJKLMNOPQRSTUVWXYZ[\\]^_` abcdefghijklmnopqrstuvwxyz{|}~",
				},
			},
		},
		Text: `Jeżu klątw, spłódź Finom część gry hańb!
Pójdźże, kiń tę chmurność w głąb flaszy!
Mężny bądź chroń pułk twój i sześć flag.
Filmuj rzeź żądań, pość, gnęb chłystków!
Pchnąć w tę łódź jeża lub ośm skrzyń fig.
Dość gróźb fuzją, klnę, pych i małżeństw!
Pójdź w loch zbić małżeńską gęś futryn!
Chwyć małżonkę, strój bądź pleśń z fugi.`,
		AttachedFiles: []AttachedFile{
			{
				ContentType: ContentTypeHeader{
					ContentType: "application/pkcs7-signature",
					Params: map[string]string{
						"name": "smime.p7s",
					},
				},
				ContentDisposition: ContentDispositionHeader{
					ContentDisposition: attachment,
					Params: map[string]string{
						"filename": "smime.p7s",
					},
				},
				Data: []byte{130, 28, 135, 132, 117, 46, 142, 18, 97, 140, 126, 251, 159, 193, 199, 25, 58, 223, 189, 185, 227, 239, 158, 173, 108, 31, 71, 27, 133, 80, 165, 252, 133, 227, 174, 198, 132, 129, 159, 29, 246, 19, 235, 133, 80, 165, 252, 133, 227, 174, 198, 132, 129, 159, 29, 246, 19, 234, 49, 251, 238, 127, 7, 28, 104, 33, 200, 120, 71, 82, 232, 225, 38, 30, 249, 234, 214, 193, 244, 113, 147, 173, 251, 219, 158, 57, 252, 28, 113, 147, 173, 251, 225, 38, 24, 199, 239, 190, 173, 108, 31, 71, 27, 133, 80, 110, 120, 251, 231, 174, 198, 132, 129, 159, 29, 246, 19, 234, 8, 114, 30, 17, 212, 186, 58, 95, 200, 94, 59, 26, 18, 6, 124, 119, 216, 79, 174, 21, 65, 185, 227, 239, 158},
			},
		},
	}

	testEmailFromFile(t, fp, expectedEmail)
}

func TestParseEmailPolishMultipartSignedIso88592OverQuotedprintable(t *testing.T) {
	fp := "tests/test_polish_multipart_signed_iso-8859-2_over_quoted-printable.txt"
	tz, _ := time.LoadLocation("Europe/Warsaw")
	expectedDate, _ := time.Parse(
		time.RFC1123Z+" (MST)",
		time.Date(2019, time.April, 1, 7, 55, 0, 0, tz).Format(time.RFC1123Z+" (MST)"))
	expectedEmail := Email{
		Headers: Headers{
			Date:    expectedDate,
			Subject: "📧 Signed Test Polskie pangramy",
			ReplyTo: []*mail.Address{
				{
					Name:    "Nadająca, Alicja",
					Address: "alicja.nadajaca@example.net",
				},
			},
			Sender: &mail.Address{
				Name:    "Nadająca, Alicja",
				Address: "alicja.nadajaca@example.com",
			},
			From: []*mail.Address{
				{
					Name:    "Nadająca, Alicja",
					Address: "alicja.nadajaca@example.com",
				},
				{
					Name:    "Nadająca, Alicja",
					Address: "alicja.nadajaca@example.net",
				},
			},
			To: []*mail.Address{
				{
					Name:    "Odbierający, Bob",
					Address: "bob.odbierajacy@example.com",
				},
				{
					Name:    "Odbierająca, Karolina",
					Address: "karolina.odbierajaca@example.com",
				},
			},
			Cc: []*mail.Address{
				{
					Name:    "Odbierający, Daniel",
					Address: "daniel.odbierajacy@example.com",
				},
				{
					Name:    "Odbierająca, Ewa",
					Address: "ewa.odbierajaca@example.com",
				},
			},
			Bcc: []*mail.Address{
				{
					Name:    "Odbierający, Franek",
					Address: "franek.odbierajacy@example.com",
				},
				{
					Name:    "Odbierająca, Grażyna",
					Address: "grazyna.odbierajaca@example.com",
				},
			},
			MessageID:  "Message-Id-1@example.com",
			InReplyTo:  []MessageId{"Message-Id-0@example.com"},
			References: []MessageId{"Message-Id-0@example.com"},
			Comments:   "Message Header Comment",
			Keywords:   []string{"Keyword 1", "Keyword 2"},
			ResentDate: expectedDate,
			ResentFrom: []*mail.Address{
				{
					Name:    "Nadająca, Alicja",
					Address: "alicja.nadajaca@example.net",
				},
				{
					Name:    "Nadająca, Alicja",
					Address: "alicja.nadajaca@example.com",
				},
			},
			ResentSender: &mail.Address{
				Name:    "Nadająca, Alicja",
				Address: "alicja.nadajaca@example.net",
			},
			ResentTo: []*mail.Address{
				{
					Name:    "Odbierający, Bob",
					Address: "bob.odbierajacy@example.net",
				},
				{
					Name:    "Odbierająca, Karolina",
					Address: "karolina.odbierajaca@example.net",
				},
			},
			ResentCc: []*mail.Address{
				{
					Name:    "Odbierający, Daniel",
					Address: "daniel.odbierajacy@example.net",
				},
				{
					Name:    "Odbierająca, Ewa",
					Address: "ewa.odbierajaca@example.net",
				},
			},
			ResentBcc: []*mail.Address{
				{
					Name:    "Odbierający, Franek",
					Address: "franek.odbierajacy@example.net",
				},
				{
					Name:    "Odbierająca, Grażyna",
					Address: "grazyna.odbierajaca@example.net",
				},
			},
			ResentMessageID: "Message-Id-1@example.net",
			ContentType: ContentTypeHeader{
				ContentType: "multipart/signed",
				Params: map[string]string{
					"boundary": "SignedBoundaryString",
					"charset":  "iso-8859-2",
					"micalg":   "sha1",
					"protocol": "application/pkcs7-signature",
				},
			},
			ExtraHeaders: map[string][]string{
				"X-Clacks-Overhead": {"GNU Terry Pratchett"},
				"X-Script/function/\t !\"#$%&'()*+,-./;<=>?@[\\]^_`{|}~": {
					"TEST VALUE 1\t !\"#$%&'()*+,-./0123456789:;<=>?@ABCDEFGHIJKLMNOPQRSTUVWXYZ[\\]^_` abcdefghijklmnopqrstuvwxyz{|}~",
					"TEST VALUE 2\t !\"#$%&'()*+,-./0123456789:;<=>?@ABCDEFGHIJKLMNOPQRSTUVWXYZ[\\]^_` abcdefghijklmnopqrstuvwxyz{|}~",
				},
			},
		},
		Text: `Jeżu klątw, spłódź Finom część gry hańb!
Pójdźże, kiń tę chmurność w głąb flaszy!
Mężny bądź chroń pułk twój i sześć flag.
Filmuj rzeź żądań, pość, gnęb chłystków!
Pchnąć w tę łódź jeża lub ośm skrzyń fig.
Dość gróźb fuzją, klnę, pych i małżeństw!
Pójdź w loch zbić małżeńską gęś futryn!
Chwyć małżonkę, strój bądź pleśń z fugi.`,
		AttachedFiles: []AttachedFile{
			{
				ContentType: ContentTypeHeader{
					ContentType: "application/pkcs7-signature",
					Params: map[string]string{
						"name": "smime.p7s",
					},
				},
				ContentDisposition: ContentDispositionHeader{
					ContentDisposition: attachment,
					Params: map[string]string{
						"filename": "smime.p7s",
					},
				},
				Data: []byte{130, 28, 135, 132, 117, 46, 142, 18, 97, 140, 126, 251, 159, 193, 199, 25, 58, 223, 189, 185, 227, 239, 158, 173, 108, 31, 71, 27, 133, 80, 165, 252, 133, 227, 174, 198, 132, 129, 159, 29, 246, 19, 235, 133, 80, 165, 252, 133, 227, 174, 198, 132, 129, 159, 29, 246, 19, 234, 49, 251, 238, 127, 7, 28, 104, 33, 200, 120, 71, 82, 232, 225, 38, 30, 249, 234, 214, 193, 244, 113, 147, 173, 251, 219, 158, 57, 252, 28, 113, 147, 173, 251, 225, 38, 24, 199, 239, 190, 173, 108, 31, 71, 27, 133, 80, 110, 120, 251, 231, 174, 198, 132, 129, 159, 29, 246, 19, 234, 8, 114, 30, 17, 212, 186, 58, 95, 200, 94, 59, 26, 18, 6, 124, 119, 216, 79, 174, 21, 65, 185, 227, 239, 158},
			},
		},
	}

	testEmailFromFile(t, fp, expectedEmail)
}

func TestParseEmailThaiPlaintextIso885911OverBase64(t *testing.T) {
	fp := "tests/test_thai_plaintext_iso-8859-11_over_base64.txt"
	tz, _ := time.LoadLocation("Asia/Bangkok")
	expectedDate, _ := time.Parse(
		time.RFC1123Z+" (MST)",
		time.Date(2019, time.April, 1, 7, 55, 0, 0, tz).Format(time.RFC1123Z+" (MST)"))
	expectedEmail := Email{
		Headers: Headers{
			Date:    expectedDate,
			Subject: "📧 Test แพนแกรมภาษาไทย",
			ReplyTo: []*mail.Address{
				{
					Name:    "อลิซ ผู้ส่งจดหมาย",
					Address: "alis.phusngcdhmay@example.net",
				},
			},
			Sender: &mail.Address{
				Name:    "อลิซ ผู้ส่งจดหมาย",
				Address: "alis.phusngcdhmay@example.com",
			},
			From: []*mail.Address{
				{
					Name:    "อลิซ ผู้ส่งจดหมาย",
					Address: "alis.phusngcdhmay@example.com",
				},
				{
					Name:    "อลิซ ผู้ส่งจดหมาย",
					Address: "alis.phusngcdhmay@example.net",
				},
			},
			To: []*mail.Address{
				{
					Name:    "บ๊อบ ผู้รับ",
					Address: "bob.phurab@example.com",
				},
				{
					Name:    "คาโรล ผู้รับ",
					Address: "carol.phurab@example.com",
				},
			},
			Cc: []*mail.Address{
				{
					Name:    "แดน ผู้รับ",
					Address: "dan.phurab@example.com",
				},
				{
					Name:    "อีฟ ผู้รับ",
					Address: "eve.phurab@example.com",
				},
			},
			Bcc: []*mail.Address{
				{
					Name:    "แฟรงค์ ผู้รับ",
					Address: "frank.phurab@example.com",
				},
				{
					Name:    "เกรซ ผู้รับ",
					Address: "grace.phurab@example.com",
				},
			},
			MessageID:  "Message-Id-1@example.com",
			InReplyTo:  []MessageId{"Message-Id-0@example.com"},
			References: []MessageId{"Message-Id-0@example.com"},
			Comments:   "Message Header Comment",
			Keywords:   []string{"Keyword 1", "Keyword 2"},
			ResentDate: expectedDate,
			ResentFrom: []*mail.Address{
				{
					Name:    "อลิซ ผู้ส่งจดหมาย",
					Address: "alis.phusngcdhmay@example.net",
				},
				{
					Name:    "อลิซ ผู้ส่งจดหมาย",
					Address: "alis.phusngcdhmay@example.com",
				},
			},
			ResentSender: &mail.Address{
				Name:    "อลิซ ผู้ส่งจดหมาย",
				Address: "alis.phusngcdhmay@example.net",
			},
			ResentTo: []*mail.Address{
				{
					Name:    "บ๊อบ ผู้รับ",
					Address: "bob.phurab@example.net",
				},
				{
					Name:    "คาโรล ผู้รับ",
					Address: "carol.phurab@example.net",
				},
			},
			ResentCc: []*mail.Address{
				{
					Name:    "แดน ผู้รับ",
					Address: "dan.phurab@example.net",
				},
				{
					Name:    "อีฟ ผู้รับ",
					Address: "eve.phurab@example.net",
				},
			},
			ResentBcc: []*mail.Address{
				{
					Name:    "แฟรงค์ ผู้รับ",
					Address: "frank.phurab@example.net",
				},
				{
					Name:    "เกรซ ผู้รับ",
					Address: "grace.phurab@example.net",
				},
			},
			ResentMessageID: "Message-Id-1@example.net",
			ContentType: ContentTypeHeader{
				ContentType: "text/plain",
				Params: map[string]string{
					"charset": "iso-8859-11",
				},
			},
			ExtraHeaders: map[string][]string{
				"X-Clacks-Overhead": {"GNU Terry Pratchett"},
			},
		},
		Text: `เป็นมนุษย์สุดประเสริฐเลิศคุณค่า กว่าบรรดาฝูงสัตว์เดรัจฉาน จงฝ่าฟันพัฒนาวิชาการ อย่าล้างผลาญฤๅเข่นฆ่าบีฑาใคร ไม่ถือโทษโกรธแช่งซัดฮึดฮัดด่า หัดอภัยเหมือนกีฬาอัชฌาสัย ปฏิบัติประพฤติกฎกำหนดใจ พูดจาให้จ๊ะๆ จ๋าๆ น่าฟังเอยฯ

นายสังฆภัณฑ์ เฮงพิทักษ์ฝั่ง ผู้เฒ่าซึ่งมีอาชีพเป็นฅนขายฃวด ถูกตำรวจปฏิบัติการจับฟ้องศาล ฐานลักนาฬิกาคุณหญิงฉัตรชฎา ฌานสมาธิ`,
	}

	testEmailFromFile(t, fp, expectedEmail)
}

func TestParseEmailThaiPlaintextIso885911OverQuotedprintable(t *testing.T) {
	fp := "tests/test_thai_plaintext_iso-8859-11_over_quoted-printable.txt"
	tz, _ := time.LoadLocation("Asia/Bangkok")
	expectedDate, _ := time.Parse(
		time.RFC1123Z+" (MST)",
		time.Date(2019, time.April, 1, 7, 55, 0, 0, tz).Format(time.RFC1123Z+" (MST)"))
	expectedEmail := Email{
		Headers: Headers{
			Date:    expectedDate,
			Subject: "📧 Test แพนแกรมภาษาไทย",
			ReplyTo: []*mail.Address{
				{
					Name:    "อลิซ ผู้ส่งจดหมาย",
					Address: "alis.phusngcdhmay@example.net",
				},
			},
			Sender: &mail.Address{
				Name:    "อลิซ ผู้ส่งจดหมาย",
				Address: "alis.phusngcdhmay@example.com",
			},
			From: []*mail.Address{
				{
					Name:    "อลิซ ผู้ส่งจดหมาย",
					Address: "alis.phusngcdhmay@example.com",
				},
				{
					Name:    "อลิซ ผู้ส่งจดหมาย",
					Address: "alis.phusngcdhmay@example.net",
				},
			},
			To: []*mail.Address{
				{
					Name:    "บ๊อบ ผู้รับ",
					Address: "bob.phurab@example.com",
				},
				{
					Name:    "คาโรล ผู้รับ",
					Address: "carol.phurab@example.com",
				},
			},
			Cc: []*mail.Address{
				{
					Name:    "แดน ผู้รับ",
					Address: "dan.phurab@example.com",
				},
				{
					Name:    "อีฟ ผู้รับ",
					Address: "eve.phurab@example.com",
				},
			},
			Bcc: []*mail.Address{
				{
					Name:    "แฟรงค์ ผู้รับ",
					Address: "frank.phurab@example.com",
				},
				{
					Name:    "เกรซ ผู้รับ",
					Address: "grace.phurab@example.com",
				},
			},
			MessageID:  "Message-Id-1@example.com",
			InReplyTo:  []MessageId{"Message-Id-0@example.com"},
			References: []MessageId{"Message-Id-0@example.com"},
			Comments:   "Message Header Comment",
			Keywords:   []string{"Keyword 1", "Keyword 2"},
			ResentDate: expectedDate,
			ResentFrom: []*mail.Address{
				{
					Name:    "อลิซ ผู้ส่งจดหมาย",
					Address: "alis.phusngcdhmay@example.net",
				},
				{
					Name:    "อลิซ ผู้ส่งจดหมาย",
					Address: "alis.phusngcdhmay@example.com",
				},
			},
			ResentSender: &mail.Address{
				Name:    "อลิซ ผู้ส่งจดหมาย",
				Address: "alis.phusngcdhmay@example.net",
			},
			ResentTo: []*mail.Address{
				{
					Name:    "บ๊อบ ผู้รับ",
					Address: "bob.phurab@example.net",
				},
				{
					Name:    "คาโรล ผู้รับ",
					Address: "carol.phurab@example.net",
				},
			},
			ResentCc: []*mail.Address{
				{
					Name:    "แดน ผู้รับ",
					Address: "dan.phurab@example.net",
				},
				{
					Name:    "อีฟ ผู้รับ",
					Address: "eve.phurab@example.net",
				},
			},
			ResentBcc: []*mail.Address{
				{
					Name:    "แฟรงค์ ผู้รับ",
					Address: "frank.phurab@example.net",
				},
				{
					Name:    "เกรซ ผู้รับ",
					Address: "grace.phurab@example.net",
				},
			},
			ResentMessageID: "Message-Id-1@example.net",
			ContentType: ContentTypeHeader{
				ContentType: "text/plain",
				Params: map[string]string{
					"charset": "iso-8859-11",
				},
			},
			ExtraHeaders: map[string][]string{
				"X-Clacks-Overhead": {"GNU Terry Pratchett"},
			},
		},
		Text: `เป็นมนุษย์สุดประเสริฐเลิศคุณค่า กว่าบรรดาฝูงสัตว์เดรัจฉาน จงฝ่าฟันพัฒนาวิชาการ อย่าล้างผลาญฤๅเข่นฆ่าบีฑาใคร ไม่ถือโทษโกรธแช่งซัดฮึดฮัดด่า หัดอภัยเหมือนกีฬาอัชฌาสัย ปฏิบัติประพฤติกฎกำหนดใจ พูดจาให้จ๊ะๆ จ๋าๆ น่าฟังเอยฯ

นายสังฆภัณฑ์ เฮงพิทักษ์ฝั่ง ผู้เฒ่าซึ่งมีอาชีพเป็นฅนขายฃวด ถูกตำรวจปฏิบัติการจับฟ้องศาล ฐานลักนาฬิกาคุณหญิงฉัตรชฎา ฌานสมาธิ`,
	}

	testEmailFromFile(t, fp, expectedEmail)
}

func TestParseEmailThaiPlaintextWindows874OverBase64(t *testing.T) {
	fp := "tests/test_thai_plaintext_windows-874_over_base64.txt"
	tz, _ := time.LoadLocation("Asia/Bangkok")
	expectedDate, _ := time.Parse(
		time.RFC1123Z+" (MST)",
		time.Date(2019, time.April, 1, 7, 55, 0, 0, tz).Format(time.RFC1123Z+" (MST)"))
	expectedEmail := Email{
		Headers: Headers{
			Date:    expectedDate,
			Subject: "📧 Test แพนแกรมภาษาไทย",
			ReplyTo: []*mail.Address{
				{
					Name:    "อลิซ ผู้ส่งจดหมาย",
					Address: "alis.phusngcdhmay@example.net",
				},
			},
			Sender: &mail.Address{
				Name:    "อลิซ ผู้ส่งจดหมาย",
				Address: "alis.phusngcdhmay@example.com",
			},
			From: []*mail.Address{
				{
					Name:    "อลิซ ผู้ส่งจดหมาย",
					Address: "alis.phusngcdhmay@example.com",
				},
				{
					Name:    "อลิซ ผู้ส่งจดหมาย",
					Address: "alis.phusngcdhmay@example.net",
				},
			},
			To: []*mail.Address{
				{
					Name:    "บ๊อบ ผู้รับ",
					Address: "bob.phurab@example.com",
				},
				{
					Name:    "คาโรล ผู้รับ",
					Address: "carol.phurab@example.com",
				},
			},
			Cc: []*mail.Address{
				{
					Name:    "แดน ผู้รับ",
					Address: "dan.phurab@example.com",
				},
				{
					Name:    "อีฟ ผู้รับ",
					Address: "eve.phurab@example.com",
				},
			},
			Bcc: []*mail.Address{
				{
					Name:    "แฟรงค์ ผู้รับ",
					Address: "frank.phurab@example.com",
				},
				{
					Name:    "เกรซ ผู้รับ",
					Address: "grace.phurab@example.com",
				},
			},
			MessageID:  "Message-Id-1@example.com",
			InReplyTo:  []MessageId{"Message-Id-0@example.com"},
			References: []MessageId{"Message-Id-0@example.com"},
			Comments:   "Message Header Comment",
			Keywords:   []string{"Keyword 1", "Keyword 2"},
			ResentDate: expectedDate,
			ResentFrom: []*mail.Address{
				{
					Name:    "อลิซ ผู้ส่งจดหมาย",
					Address: "alis.phusngcdhmay@example.net",
				},
				{
					Name:    "อลิซ ผู้ส่งจดหมาย",
					Address: "alis.phusngcdhmay@example.com",
				},
			},
			ResentSender: &mail.Address{
				Name:    "อลิซ ผู้ส่งจดหมาย",
				Address: "alis.phusngcdhmay@example.net",
			},
			ResentTo: []*mail.Address{
				{
					Name:    "บ๊อบ ผู้รับ",
					Address: "bob.phurab@example.net",
				},
				{
					Name:    "คาโรล ผู้รับ",
					Address: "carol.phurab@example.net",
				},
			},
			ResentCc: []*mail.Address{
				{
					Name:    "แดน ผู้รับ",
					Address: "dan.phurab@example.net",
				},
				{
					Name:    "อีฟ ผู้รับ",
					Address: "eve.phurab@example.net",
				},
			},
			ResentBcc: []*mail.Address{
				{
					Name:    "แฟรงค์ ผู้รับ",
					Address: "frank.phurab@example.net",
				},
				{
					Name:    "เกรซ ผู้รับ",
					Address: "grace.phurab@example.net",
				},
			},
			ResentMessageID: "Message-Id-1@example.net",
			ContentType: ContentTypeHeader{
				ContentType: "text/plain",
				Params: map[string]string{
					"charset": "windows-874",
				},
			},
			ExtraHeaders: map[string][]string{
				"X-Clacks-Overhead": {"GNU Terry Pratchett"},
			},
		},
		Text: `เป็นมนุษย์สุดประเสริฐเลิศคุณค่า กว่าบรรดาฝูงสัตว์เดรัจฉาน จงฝ่าฟันพัฒนาวิชาการ อย่าล้างผลาญฤๅเข่นฆ่าบีฑาใคร ไม่ถือโทษโกรธแช่งซัดฮึดฮัดด่า หัดอภัยเหมือนกีฬาอัชฌาสัย ปฏิบัติประพฤติกฎกำหนดใจ พูดจาให้จ๊ะๆ จ๋าๆ น่าฟังเอยฯ

นายสังฆภัณฑ์ เฮงพิทักษ์ฝั่ง ผู้เฒ่าซึ่งมีอาชีพเป็นฅนขายฃวด ถูกตำรวจปฏิบัติการจับฟ้องศาล ฐานลักนาฬิกาคุณหญิงฉัตรชฎา ฌานสมาธิ`,
	}

	testEmailFromFile(t, fp, expectedEmail)
}

func TestParseEmailThaiPlaintextWindows874OverQuotedprintable(t *testing.T) {
	fp := "tests/test_thai_plaintext_windows-874_over_quoted-printable.txt"
	tz, _ := time.LoadLocation("Asia/Bangkok")
	expectedDate, _ := time.Parse(
		time.RFC1123Z+" (MST)",
		time.Date(2019, time.April, 1, 7, 55, 0, 0, tz).Format(time.RFC1123Z+" (MST)"))
	expectedEmail := Email{
		Headers: Headers{
			Date:    expectedDate,
			Subject: "📧 Test แพนแกรมภาษาไทย",
			ReplyTo: []*mail.Address{
				{
					Name:    "อลิซ ผู้ส่งจดหมาย",
					Address: "alis.phusngcdhmay@example.net",
				},
			},
			Sender: &mail.Address{
				Name:    "อลิซ ผู้ส่งจดหมาย",
				Address: "alis.phusngcdhmay@example.com",
			},
			From: []*mail.Address{
				{
					Name:    "อลิซ ผู้ส่งจดหมาย",
					Address: "alis.phusngcdhmay@example.com",
				},
				{
					Name:    "อลิซ ผู้ส่งจดหมาย",
					Address: "alis.phusngcdhmay@example.net",
				},
			},
			To: []*mail.Address{
				{
					Name:    "บ๊อบ ผู้รับ",
					Address: "bob.phurab@example.com",
				},
				{
					Name:    "คาโรล ผู้รับ",
					Address: "carol.phurab@example.com",
				},
			},
			Cc: []*mail.Address{
				{
					Name:    "แดน ผู้รับ",
					Address: "dan.phurab@example.com",
				},
				{
					Name:    "อีฟ ผู้รับ",
					Address: "eve.phurab@example.com",
				},
			},
			Bcc: []*mail.Address{
				{
					Name:    "แฟรงค์ ผู้รับ",
					Address: "frank.phurab@example.com",
				},
				{
					Name:    "เกรซ ผู้รับ",
					Address: "grace.phurab@example.com",
				},
			},
			MessageID:  "Message-Id-1@example.com",
			InReplyTo:  []MessageId{"Message-Id-0@example.com"},
			References: []MessageId{"Message-Id-0@example.com"},
			Comments:   "Message Header Comment",
			Keywords:   []string{"Keyword 1", "Keyword 2"},
			ResentDate: expectedDate,
			ResentFrom: []*mail.Address{
				{
					Name:    "อลิซ ผู้ส่งจดหมาย",
					Address: "alis.phusngcdhmay@example.net",
				},
				{
					Name:    "อลิซ ผู้ส่งจดหมาย",
					Address: "alis.phusngcdhmay@example.com",
				},
			},
			ResentSender: &mail.Address{
				Name:    "อลิซ ผู้ส่งจดหมาย",
				Address: "alis.phusngcdhmay@example.net",
			},
			ResentTo: []*mail.Address{
				{
					Name:    "บ๊อบ ผู้รับ",
					Address: "bob.phurab@example.net",
				},
				{
					Name:    "คาโรล ผู้รับ",
					Address: "carol.phurab@example.net",
				},
			},
			ResentCc: []*mail.Address{
				{
					Name:    "แดน ผู้รับ",
					Address: "dan.phurab@example.net",
				},
				{
					Name:    "อีฟ ผู้รับ",
					Address: "eve.phurab@example.net",
				},
			},
			ResentBcc: []*mail.Address{
				{
					Name:    "แฟรงค์ ผู้รับ",
					Address: "frank.phurab@example.net",
				},
				{
					Name:    "เกรซ ผู้รับ",
					Address: "grace.phurab@example.net",
				},
			},
			ResentMessageID: "Message-Id-1@example.net",
			ContentType: ContentTypeHeader{
				ContentType: "text/plain",
				Params: map[string]string{
					"charset": "windows-874",
				},
			},
			ExtraHeaders: map[string][]string{
				"X-Clacks-Overhead": {"GNU Terry Pratchett"},
			},
		},
		Text: `เป็นมนุษย์สุดประเสริฐเลิศคุณค่า กว่าบรรดาฝูงสัตว์เดรัจฉาน จงฝ่าฟันพัฒนาวิชาการ อย่าล้างผลาญฤๅเข่นฆ่าบีฑาใคร ไม่ถือโทษโกรธแช่งซัดฮึดฮัดด่า หัดอภัยเหมือนกีฬาอัชฌาสัย ปฏิบัติประพฤติกฎกำหนดใจ พูดจาให้จ๊ะๆ จ๋าๆ น่าฟังเอยฯ

นายสังฆภัณฑ์ เฮงพิทักษ์ฝั่ง ผู้เฒ่าซึ่งมีอาชีพเป็นฅนขายฃวด ถูกตำรวจปฏิบัติการจับฟ้องศาล ฐานลักนาฬิกาคุณหญิงฉัตรชฎา ฌานสมาธิ`,
	}

	testEmailFromFile(t, fp, expectedEmail)
}

func TestParseEmailThaiPlaintextTis620OverBase64(t *testing.T) {
	fp := "tests/test_thai_plaintext_tis-620_over_base64.txt"
	tz, _ := time.LoadLocation("Asia/Bangkok")
	expectedDate, _ := time.Parse(
		time.RFC1123Z+" (MST)",
		time.Date(2019, time.April, 1, 7, 55, 0, 0, tz).Format(time.RFC1123Z+" (MST)"))
	expectedEmail := Email{
		Headers: Headers{
			Date:    expectedDate,
			Subject: "📧 Test แพนแกรมภาษาไทย",
			ReplyTo: []*mail.Address{
				{
					Name:    "อลิซ ผู้ส่งจดหมาย",
					Address: "alis.phusngcdhmay@example.net",
				},
			},
			Sender: &mail.Address{
				Name:    "อลิซ ผู้ส่งจดหมาย",
				Address: "alis.phusngcdhmay@example.com",
			},
			From: []*mail.Address{
				{
					Name:    "อลิซ ผู้ส่งจดหมาย",
					Address: "alis.phusngcdhmay@example.com",
				},
				{
					Name:    "อลิซ ผู้ส่งจดหมาย",
					Address: "alis.phusngcdhmay@example.net",
				},
			},
			To: []*mail.Address{
				{
					Name:    "บ๊อบ ผู้รับ",
					Address: "bob.phurab@example.com",
				},
				{
					Name:    "คาโรล ผู้รับ",
					Address: "carol.phurab@example.com",
				},
			},
			Cc: []*mail.Address{
				{
					Name:    "แดน ผู้รับ",
					Address: "dan.phurab@example.com",
				},
				{
					Name:    "อีฟ ผู้รับ",
					Address: "eve.phurab@example.com",
				},
			},
			Bcc: []*mail.Address{
				{
					Name:    "แฟรงค์ ผู้รับ",
					Address: "frank.phurab@example.com",
				},
				{
					Name:    "เกรซ ผู้รับ",
					Address: "grace.phurab@example.com",
				},
			},
			MessageID:  "Message-Id-1@example.com",
			InReplyTo:  []MessageId{"Message-Id-0@example.com"},
			References: []MessageId{"Message-Id-0@example.com"},
			Comments:   "Message Header Comment",
			Keywords:   []string{"Keyword 1", "Keyword 2"},
			ResentDate: expectedDate,
			ResentFrom: []*mail.Address{
				{
					Name:    "อลิซ ผู้ส่งจดหมาย",
					Address: "alis.phusngcdhmay@example.net",
				},
				{
					Name:    "อลิซ ผู้ส่งจดหมาย",
					Address: "alis.phusngcdhmay@example.com",
				},
			},
			ResentSender: &mail.Address{
				Name:    "อลิซ ผู้ส่งจดหมาย",
				Address: "alis.phusngcdhmay@example.net",
			},
			ResentTo: []*mail.Address{
				{
					Name:    "บ๊อบ ผู้รับ",
					Address: "bob.phurab@example.net",
				},
				{
					Name:    "คาโรล ผู้รับ",
					Address: "carol.phurab@example.net",
				},
			},
			ResentCc: []*mail.Address{
				{
					Name:    "แดน ผู้รับ",
					Address: "dan.phurab@example.net",
				},
				{
					Name:    "อีฟ ผู้รับ",
					Address: "eve.phurab@example.net",
				},
			},
			ResentBcc: []*mail.Address{
				{
					Name:    "แฟรงค์ ผู้รับ",
					Address: "frank.phurab@example.net",
				},
				{
					Name:    "เกรซ ผู้รับ",
					Address: "grace.phurab@example.net",
				},
			},
			ResentMessageID: "Message-Id-1@example.net",
			ContentType: ContentTypeHeader{
				ContentType: "text/plain",
				Params: map[string]string{
					"charset": "tis-620",
				},
			},
			ExtraHeaders: map[string][]string{
				"X-Clacks-Overhead": {"GNU Terry Pratchett"},
			},
		},
		Text: `เป็นมนุษย์สุดประเสริฐเลิศคุณค่า กว่าบรรดาฝูงสัตว์เดรัจฉาน จงฝ่าฟันพัฒนาวิชาการ อย่าล้างผลาญฤๅเข่นฆ่าบีฑาใคร ไม่ถือโทษโกรธแช่งซัดฮึดฮัดด่า หัดอภัยเหมือนกีฬาอัชฌาสัย ปฏิบัติประพฤติกฎกำหนดใจ พูดจาให้จ๊ะๆ จ๋าๆ น่าฟังเอยฯ

นายสังฆภัณฑ์ เฮงพิทักษ์ฝั่ง ผู้เฒ่าซึ่งมีอาชีพเป็นฅนขายฃวด ถูกตำรวจปฏิบัติการจับฟ้องศาล ฐานลักนาฬิกาคุณหญิงฉัตรชฎา ฌานสมาธิ`,
	}

	testEmailFromFile(t, fp, expectedEmail)
}

func TestParseEmailThaiPlaintextTis620OverQuotedprintable(t *testing.T) {
	fp := "tests/test_thai_plaintext_tis-620_over_quoted-printable.txt"
	tz, _ := time.LoadLocation("Asia/Bangkok")
	expectedDate, _ := time.Parse(
		time.RFC1123Z+" (MST)",
		time.Date(2019, time.April, 1, 7, 55, 0, 0, tz).Format(time.RFC1123Z+" (MST)"))
	expectedEmail := Email{
		Headers: Headers{
			Date:    expectedDate,
			Subject: "📧 Test แพนแกรมภาษาไทย",
			ReplyTo: []*mail.Address{
				{
					Name:    "อลิซ ผู้ส่งจดหมาย",
					Address: "alis.phusngcdhmay@example.net",
				},
			},
			Sender: &mail.Address{
				Name:    "อลิซ ผู้ส่งจดหมาย",
				Address: "alis.phusngcdhmay@example.com",
			},
			From: []*mail.Address{
				{
					Name:    "อลิซ ผู้ส่งจดหมาย",
					Address: "alis.phusngcdhmay@example.com",
				},
				{
					Name:    "อลิซ ผู้ส่งจดหมาย",
					Address: "alis.phusngcdhmay@example.net",
				},
			},
			To: []*mail.Address{
				{
					Name:    "บ๊อบ ผู้รับ",
					Address: "bob.phurab@example.com",
				},
				{
					Name:    "คาโรล ผู้รับ",
					Address: "carol.phurab@example.com",
				},
			},
			Cc: []*mail.Address{
				{
					Name:    "แดน ผู้รับ",
					Address: "dan.phurab@example.com",
				},
				{
					Name:    "อีฟ ผู้รับ",
					Address: "eve.phurab@example.com",
				},
			},
			Bcc: []*mail.Address{
				{
					Name:    "แฟรงค์ ผู้รับ",
					Address: "frank.phurab@example.com",
				},
				{
					Name:    "เกรซ ผู้รับ",
					Address: "grace.phurab@example.com",
				},
			},
			MessageID:  "Message-Id-1@example.com",
			InReplyTo:  []MessageId{"Message-Id-0@example.com"},
			References: []MessageId{"Message-Id-0@example.com"},
			Comments:   "Message Header Comment",
			Keywords:   []string{"Keyword 1", "Keyword 2"},
			ResentDate: expectedDate,
			ResentFrom: []*mail.Address{
				{
					Name:    "อลิซ ผู้ส่งจดหมาย",
					Address: "alis.phusngcdhmay@example.net",
				},
				{
					Name:    "อลิซ ผู้ส่งจดหมาย",
					Address: "alis.phusngcdhmay@example.com",
				},
			},
			ResentSender: &mail.Address{
				Name:    "อลิซ ผู้ส่งจดหมาย",
				Address: "alis.phusngcdhmay@example.net",
			},
			ResentTo: []*mail.Address{
				{
					Name:    "บ๊อบ ผู้รับ",
					Address: "bob.phurab@example.net",
				},
				{
					Name:    "คาโรล ผู้รับ",
					Address: "carol.phurab@example.net",
				},
			},
			ResentCc: []*mail.Address{
				{
					Name:    "แดน ผู้รับ",
					Address: "dan.phurab@example.net",
				},
				{
					Name:    "อีฟ ผู้รับ",
					Address: "eve.phurab@example.net",
				},
			},
			ResentBcc: []*mail.Address{
				{
					Name:    "แฟรงค์ ผู้รับ",
					Address: "frank.phurab@example.net",
				},
				{
					Name:    "เกรซ ผู้รับ",
					Address: "grace.phurab@example.net",
				},
			},
			ResentMessageID: "Message-Id-1@example.net",
			ContentType: ContentTypeHeader{
				ContentType: "text/plain",
				Params: map[string]string{
					"charset": "tis-620",
				},
			},
			ExtraHeaders: map[string][]string{
				"X-Clacks-Overhead": {"GNU Terry Pratchett"},
			},
		},
		Text: `เป็นมนุษย์สุดประเสริฐเลิศคุณค่า กว่าบรรดาฝูงสัตว์เดรัจฉาน จงฝ่าฟันพัฒนาวิชาการ อย่าล้างผลาญฤๅเข่นฆ่าบีฑาใคร ไม่ถือโทษโกรธแช่งซัดฮึดฮัดด่า หัดอภัยเหมือนกีฬาอัชฌาสัย ปฏิบัติประพฤติกฎกำหนดใจ พูดจาให้จ๊ะๆ จ๋าๆ น่าฟังเอยฯ

นายสังฆภัณฑ์ เฮงพิทักษ์ฝั่ง ผู้เฒ่าซึ่งมีอาชีพเป็นฅนขายฃวด ถูกตำรวจปฏิบัติการจับฟ้องศาล ฐานลักนาฬิกาคุณหญิงฉัตรชฎา ฌานสมาธิ`,
	}

	testEmailFromFile(t, fp, expectedEmail)
}

func TestParseEmailThaiMultipartRelatedIso885911OverBase64(t *testing.T) {
	fp := "tests/test_thai_multipart_related_iso-8859-11_over_base64.txt"
	tz, _ := time.LoadLocation("Asia/Bangkok")
	expectedDate, _ := time.Parse(
		time.RFC1123Z+" (MST)",
		time.Date(2019, time.April, 1, 7, 55, 0, 0, tz).Format(time.RFC1123Z+" (MST)"))
	expectedEmail := Email{
		Headers: Headers{
			Date:    expectedDate,
			Subject: "📧 Test แพนแกรมภาษาไทย",
			ReplyTo: []*mail.Address{
				{
					Name:    "อลิซ ผู้ส่งจดหมาย",
					Address: "alis.phusngcdhmay@example.net",
				},
			},
			Sender: &mail.Address{
				Name:    "อลิซ ผู้ส่งจดหมาย",
				Address: "alis.phusngcdhmay@example.com",
			},
			From: []*mail.Address{
				{
					Name:    "อลิซ ผู้ส่งจดหมาย",
					Address: "alis.phusngcdhmay@example.com",
				},
				{
					Name:    "อลิซ ผู้ส่งจดหมาย",
					Address: "alis.phusngcdhmay@example.net",
				},
			},
			To: []*mail.Address{
				{
					Name:    "บ๊อบ ผู้รับ",
					Address: "bob.phurab@example.com",
				},
				{
					Name:    "คาโรล ผู้รับ",
					Address: "carol.phurab@example.com",
				},
			},
			Cc: []*mail.Address{
				{
					Name:    "แดน ผู้รับ",
					Address: "dan.phurab@example.com",
				},
				{
					Name:    "อีฟ ผู้รับ",
					Address: "eve.phurab@example.com",
				},
			},
			Bcc: []*mail.Address{
				{
					Name:    "แฟรงค์ ผู้รับ",
					Address: "frank.phurab@example.com",
				},
				{
					Name:    "เกรซ ผู้รับ",
					Address: "grace.phurab@example.com",
				},
			},
			MessageID:  "Message-Id-1@example.com",
			InReplyTo:  []MessageId{"Message-Id-0@example.com"},
			References: []MessageId{"Message-Id-0@example.com"},
			Comments:   "Message Header Comment",
			Keywords:   []string{"Keyword 1", "Keyword 2"},
			ResentDate: expectedDate,
			ResentFrom: []*mail.Address{
				{
					Name:    "อลิซ ผู้ส่งจดหมาย",
					Address: "alis.phusngcdhmay@example.net",
				},
				{
					Name:    "อลิซ ผู้ส่งจดหมาย",
					Address: "alis.phusngcdhmay@example.com",
				},
			},
			ResentSender: &mail.Address{
				Name:    "อลิซ ผู้ส่งจดหมาย",
				Address: "alis.phusngcdhmay@example.net",
			},
			ResentTo: []*mail.Address{
				{
					Name:    "บ๊อบ ผู้รับ",
					Address: "bob.phurab@example.net",
				},
				{
					Name:    "คาโรล ผู้รับ",
					Address: "carol.phurab@example.net",
				},
			},
			ResentCc: []*mail.Address{
				{
					Name:    "แดน ผู้รับ",
					Address: "dan.phurab@example.net",
				},
				{
					Name:    "อีฟ ผู้รับ",
					Address: "eve.phurab@example.net",
				},
			},
			ResentBcc: []*mail.Address{
				{
					Name:    "แฟรงค์ ผู้รับ",
					Address: "frank.phurab@example.net",
				},
				{
					Name:    "เกรซ ผู้รับ",
					Address: "grace.phurab@example.net",
				},
			},
			ResentMessageID: "Message-Id-1@example.net",
			ContentType: ContentTypeHeader{
				ContentType: "multipart/related",
				Params: map[string]string{
					"boundary": "RelatedBoundaryString",
					"charset":  "iso-8859-11",
				},
			},
			ExtraHeaders: map[string][]string{
				"X-Clacks-Overhead": {"GNU Terry Pratchett"},
			},
		},
		Text: `เป็นมนุษย์สุดประเสริฐเลิศคุณค่า กว่าบรรดาฝูงสัตว์เดรัจฉาน จงฝ่าฟันพัฒนาวิชาการ อย่าล้างผลาญฤๅเข่นฆ่าบีฑาใคร ไม่ถือโทษโกรธแช่งซัดฮึดฮัดด่า หัดอภัยเหมือนกีฬาอัชฌาสัย ปฏิบัติประพฤติกฎกำหนดใจ พูดจาให้จ๊ะๆ จ๋าๆ น่าฟังเอยฯ

นายสังฆภัณฑ์ เฮงพิทักษ์ฝั่ง ผู้เฒ่าซึ่งมีอาชีพเป็นฅนขายฃวด ถูกตำรวจปฏิบัติการจับฟ้องศาล ฐานลักนาฬิกาคุณหญิงฉัตรชฎา ฌานสมาธิ`,
		EnrichedText: `<bold>เป็นมนุษย์สุดประเสริฐเลิศคุณค่า</bold> <italic>กว่าบรรดาฝูงสัตว์เดรัจฉาน</italic> <fixed>จงฝ่าฟันพัฒนาวิชาการ</fixed> <underline>อย่าล้างผลาญฤๅเข่นฆ่าบีฑาใคร</underline> ไม่ถือโทษโกรธแช่งซัดฮึดฮัดด่า หัดอภัยเหมือนกีฬาอัชฌาสัย ปฏิบัติประพฤติกฎกำหนดใจ พูดจาให้จ๊ะๆ จ๋าๆ น่าฟังเอยฯ

นายสังฆภัณฑ์ เฮงพิทักษ์ฝั่ง ผู้เฒ่าซึ่งมีอาชีพเป็นฅนขายฃวด ถูกตำรวจปฏิบัติการจับฟ้องศาล ฐานลักนาฬิกาคุณหญิงฉัตรชฎา ฌานสมาธิ`,
		HTML: `<html>
<div dir="ltr">
<p>เป็นมนุษย์สุดประเสริฐเลิศคุณค่า กว่าบรรดาฝูงสัตว์เดรัจฉาน จงฝ่าฟันพัฒนาวิชาการ อย่าล้างผลาญฤๅเข่นฆ่าบีฑาใคร ไม่ถือโทษโกรธแช่งซัดฮึดฮัดด่า หัดอภัยเหมือนกีฬาอัชฌาสัย ปฏิบัติประพฤติกฎกำหนดใจ พูดจาให้จ๊ะๆ จ๋าๆ น่าฟังเอยฯ</p>

<p>นายสังฆภัณฑ์ เฮงพิทักษ์ฝั่ง ผู้เฒ่าซึ่งมีอาชีพเป็นฅนขายฃวด ถูกตำรวจปฏิบัติการจับฟ้องศาล ฐานลักนาฬิกาคุณหญิงฉัตรชฎา ฌานสมาธิ</p>
</div>
</html>`,
		InlineFiles: []InlineFile{
			{
				ContentID: "inline-jpg-image.jpg@example.com",
				ContentType: ContentTypeHeader{
					ContentType: "image/jpeg",
					Params: map[string]string{
						"name": "inline-jpg-image-name.jpg",
					},
				},
				ContentDisposition: ContentDispositionHeader{
					ContentDisposition: inline,
					Params: map[string]string{
						"filename": "inline-jpg-image-filename.jpg",
					},
				},
				Data: []byte{255, 216, 255, 219, 0, 67, 0, 3, 2, 2, 2, 2, 2, 3, 2, 2, 2, 3, 3, 3, 3, 4,
					6, 4, 4, 4, 4, 4, 8, 6, 6, 5, 6, 9, 8, 10, 10, 9, 8, 9, 9, 10, 12, 15, 12, 10, 11, 14, 11, 9, 9,
					13, 17, 13, 14, 15, 16, 16, 17, 16, 10, 12, 18, 19, 18, 16, 19, 15, 16, 16, 16, 255, 201, 0, 11, 8,
					0, 1, 0, 1, 1, 1, 17, 0, 255, 204, 0, 6, 0, 16, 16, 5, 255, 218, 0, 8, 1, 1, 0, 0, 63, 0, 210, 207,
					32, 255, 217},
			},
		},
	}

	testEmailFromFile(t, fp, expectedEmail)
}
func TestParseEmailThaiMultipartRelatedIso885911OverQuotedprintable(t *testing.T) {
	fp := "tests/test_thai_multipart_related_iso-8859-11_over_quoted-printable.txt"
	tz, _ := time.LoadLocation("Asia/Bangkok")
	expectedDate, _ := time.Parse(
		time.RFC1123Z+" (MST)",
		time.Date(2019, time.April, 1, 7, 55, 0, 0, tz).Format(time.RFC1123Z+" (MST)"))
	expectedEmail := Email{
		Headers: Headers{
			Date:    expectedDate,
			Subject: "📧 Test แพนแกรมภาษาไทย",
			ReplyTo: []*mail.Address{
				{
					Name:    "อลิซ ผู้ส่งจดหมาย",
					Address: "alis.phusngcdhmay@example.net",
				},
			},
			Sender: &mail.Address{
				Name:    "อลิซ ผู้ส่งจดหมาย",
				Address: "alis.phusngcdhmay@example.com",
			},
			From: []*mail.Address{
				{
					Name:    "อลิซ ผู้ส่งจดหมาย",
					Address: "alis.phusngcdhmay@example.com",
				},
				{
					Name:    "อลิซ ผู้ส่งจดหมาย",
					Address: "alis.phusngcdhmay@example.net",
				},
			},
			To: []*mail.Address{
				{
					Name:    "บ๊อบ ผู้รับ",
					Address: "bob.phurab@example.com",
				},
				{
					Name:    "คาโรล ผู้รับ",
					Address: "carol.phurab@example.com",
				},
			},
			Cc: []*mail.Address{
				{
					Name:    "แดน ผู้รับ",
					Address: "dan.phurab@example.com",
				},
				{
					Name:    "อีฟ ผู้รับ",
					Address: "eve.phurab@example.com",
				},
			},
			Bcc: []*mail.Address{
				{
					Name:    "แฟรงค์ ผู้รับ",
					Address: "frank.phurab@example.com",
				},
				{
					Name:    "เกรซ ผู้รับ",
					Address: "grace.phurab@example.com",
				},
			},
			MessageID:  "Message-Id-1@example.com",
			InReplyTo:  []MessageId{"Message-Id-0@example.com"},
			References: []MessageId{"Message-Id-0@example.com"},
			Comments:   "Message Header Comment",
			Keywords:   []string{"Keyword 1", "Keyword 2"},
			ResentDate: expectedDate,
			ResentFrom: []*mail.Address{
				{
					Name:    "อลิซ ผู้ส่งจดหมาย",
					Address: "alis.phusngcdhmay@example.net",
				},
				{
					Name:    "อลิซ ผู้ส่งจดหมาย",
					Address: "alis.phusngcdhmay@example.com",
				},
			},
			ResentSender: &mail.Address{
				Name:    "อลิซ ผู้ส่งจดหมาย",
				Address: "alis.phusngcdhmay@example.net",
			},
			ResentTo: []*mail.Address{
				{
					Name:    "บ๊อบ ผู้รับ",
					Address: "bob.phurab@example.net",
				},
				{
					Name:    "คาโรล ผู้รับ",
					Address: "carol.phurab@example.net",
				},
			},
			ResentCc: []*mail.Address{
				{
					Name:    "แดน ผู้รับ",
					Address: "dan.phurab@example.net",
				},
				{
					Name:    "อีฟ ผู้รับ",
					Address: "eve.phurab@example.net",
				},
			},
			ResentBcc: []*mail.Address{
				{
					Name:    "แฟรงค์ ผู้รับ",
					Address: "frank.phurab@example.net",
				},
				{
					Name:    "เกรซ ผู้รับ",
					Address: "grace.phurab@example.net",
				},
			},
			ResentMessageID: "Message-Id-1@example.net",
			ContentType: ContentTypeHeader{
				ContentType: "multipart/related",
				Params: map[string]string{
					"boundary": "RelatedBoundaryString",
					"charset":  "iso-8859-11",
				},
			},
			ExtraHeaders: map[string][]string{
				"X-Clacks-Overhead": {"GNU Terry Pratchett"},
			},
		},
		Text: `เป็นมนุษย์สุดประเสริฐเลิศคุณค่า กว่าบรรดาฝูงสัตว์เดรัจฉาน จงฝ่าฟันพัฒนาวิชาการ อย่าล้างผลาญฤๅเข่นฆ่าบีฑาใคร ไม่ถือโทษโกรธแช่งซัดฮึดฮัดด่า หัดอภัยเหมือนกีฬาอัชฌาสัย ปฏิบัติประพฤติกฎกำหนดใจ พูดจาให้จ๊ะๆ จ๋าๆ น่าฟังเอยฯ

นายสังฆภัณฑ์ เฮงพิทักษ์ฝั่ง ผู้เฒ่าซึ่งมีอาชีพเป็นฅนขายฃวด ถูกตำรวจปฏิบัติการจับฟ้องศาล ฐานลักนาฬิกาคุณหญิงฉัตรชฎา ฌานสมาธิ`,
		EnrichedText: `<bold>เป็นมนุษย์สุดประเสริฐเลิศคุณค่า</bold> <italic>กว่าบรรดาฝูงสัตว์เดรัจฉาน</italic> <fixed>จงฝ่าฟันพัฒนาวิชาการ</fixed> <underline>อย่าล้างผลาญฤๅเข่นฆ่าบีฑาใคร</underline> ไม่ถือโทษโกรธแช่งซัดฮึดฮัดด่า หัดอภัยเหมือนกีฬาอัชฌาสัย ปฏิบัติประพฤติกฎกำหนดใจ พูดจาให้จ๊ะๆ จ๋าๆ น่าฟังเอยฯ

นายสังฆภัณฑ์ เฮงพิทักษ์ฝั่ง ผู้เฒ่าซึ่งมีอาชีพเป็นฅนขายฃวด ถูกตำรวจปฏิบัติการจับฟ้องศาล ฐานลักนาฬิกาคุณหญิงฉัตรชฎา ฌานสมาธิ`,
		HTML: `<html>
<div dir="ltr">
<p>เป็นมนุษย์สุดประเสริฐเลิศคุณค่า กว่าบรรดาฝูงสัตว์เดรัจฉาน จงฝ่าฟันพัฒนาวิชาการ อย่าล้างผลาญฤๅเข่นฆ่าบีฑาใคร ไม่ถือโทษโกรธแช่งซัดฮึดฮัดด่า หัดอภัยเหมือนกีฬาอัชฌาสัย ปฏิบัติประพฤติกฎกำหนดใจ พูดจาให้จ๊ะๆ จ๋าๆ น่าฟังเอยฯ</p>

<p>นายสังฆภัณฑ์ เฮงพิทักษ์ฝั่ง ผู้เฒ่าซึ่งมีอาชีพเป็นฅนขายฃวด ถูกตำรวจปฏิบัติการจับฟ้องศาล ฐานลักนาฬิกาคุณหญิงฉัตรชฎา ฌานสมาธิ</p>
</div>
</html>`,
		InlineFiles: []InlineFile{
			{
				ContentID: "inline-jpg-image.jpg@example.com",
				ContentType: ContentTypeHeader{
					ContentType: "image/jpeg",
					Params: map[string]string{
						"name": "inline-jpg-image-name.jpg",
					},
				},
				ContentDisposition: ContentDispositionHeader{
					ContentDisposition: inline,
					Params: map[string]string{
						"filename": "inline-jpg-image-filename.jpg",
					},
				},
				Data: []byte{255, 216, 255, 219, 0, 67, 0, 3, 2, 2, 2, 2, 2, 3, 2, 2, 2, 3, 3, 3, 3, 4,
					6, 4, 4, 4, 4, 4, 8, 6, 6, 5, 6, 9, 8, 10, 10, 9, 8, 9, 9, 10, 12, 15, 12, 10, 11, 14, 11, 9, 9,
					13, 17, 13, 14, 15, 16, 16, 17, 16, 10, 12, 18, 19, 18, 16, 19, 15, 16, 16, 16, 255, 201, 0, 11, 8,
					0, 1, 0, 1, 1, 1, 17, 0, 255, 204, 0, 6, 0, 16, 16, 5, 255, 218, 0, 8, 1, 1, 0, 0, 63, 0, 210, 207,
					32, 255, 217},
			},
		},
	}

	testEmailFromFile(t, fp, expectedEmail)
}
func TestParseEmailThaiMultipartRelatedWindows874OverBase64(t *testing.T) {
	fp := "tests/test_thai_multipart_related_windows-874_over_base64.txt"
	tz, _ := time.LoadLocation("Asia/Bangkok")
	expectedDate, _ := time.Parse(
		time.RFC1123Z+" (MST)",
		time.Date(2019, time.April, 1, 7, 55, 0, 0, tz).Format(time.RFC1123Z+" (MST)"))
	expectedEmail := Email{
		Headers: Headers{
			Date:    expectedDate,
			Subject: "📧 Test แพนแกรมภาษาไทย",
			ReplyTo: []*mail.Address{
				{
					Name:    "อลิซ ผู้ส่งจดหมาย",
					Address: "alis.phusngcdhmay@example.net",
				},
			},
			Sender: &mail.Address{
				Name:    "อลิซ ผู้ส่งจดหมาย",
				Address: "alis.phusngcdhmay@example.com",
			},
			From: []*mail.Address{
				{
					Name:    "อลิซ ผู้ส่งจดหมาย",
					Address: "alis.phusngcdhmay@example.com",
				},
				{
					Name:    "อลิซ ผู้ส่งจดหมาย",
					Address: "alis.phusngcdhmay@example.net",
				},
			},
			To: []*mail.Address{
				{
					Name:    "บ๊อบ ผู้รับ",
					Address: "bob.phurab@example.com",
				},
				{
					Name:    "คาโรล ผู้รับ",
					Address: "carol.phurab@example.com",
				},
			},
			Cc: []*mail.Address{
				{
					Name:    "แดน ผู้รับ",
					Address: "dan.phurab@example.com",
				},
				{
					Name:    "อีฟ ผู้รับ",
					Address: "eve.phurab@example.com",
				},
			},
			Bcc: []*mail.Address{
				{
					Name:    "แฟรงค์ ผู้รับ",
					Address: "frank.phurab@example.com",
				},
				{
					Name:    "เกรซ ผู้รับ",
					Address: "grace.phurab@example.com",
				},
			},
			MessageID:  "Message-Id-1@example.com",
			InReplyTo:  []MessageId{"Message-Id-0@example.com"},
			References: []MessageId{"Message-Id-0@example.com"},
			Comments:   "Message Header Comment",
			Keywords:   []string{"Keyword 1", "Keyword 2"},
			ResentDate: expectedDate,
			ResentFrom: []*mail.Address{
				{
					Name:    "อลิซ ผู้ส่งจดหมาย",
					Address: "alis.phusngcdhmay@example.net",
				},
				{
					Name:    "อลิซ ผู้ส่งจดหมาย",
					Address: "alis.phusngcdhmay@example.com",
				},
			},
			ResentSender: &mail.Address{
				Name:    "อลิซ ผู้ส่งจดหมาย",
				Address: "alis.phusngcdhmay@example.net",
			},
			ResentTo: []*mail.Address{
				{
					Name:    "บ๊อบ ผู้รับ",
					Address: "bob.phurab@example.net",
				},
				{
					Name:    "คาโรล ผู้รับ",
					Address: "carol.phurab@example.net",
				},
			},
			ResentCc: []*mail.Address{
				{
					Name:    "แดน ผู้รับ",
					Address: "dan.phurab@example.net",
				},
				{
					Name:    "อีฟ ผู้รับ",
					Address: "eve.phurab@example.net",
				},
			},
			ResentBcc: []*mail.Address{
				{
					Name:    "แฟรงค์ ผู้รับ",
					Address: "frank.phurab@example.net",
				},
				{
					Name:    "เกรซ ผู้รับ",
					Address: "grace.phurab@example.net",
				},
			},
			ResentMessageID: "Message-Id-1@example.net",
			ContentType: ContentTypeHeader{
				ContentType: "multipart/related",
				Params: map[string]string{
					"boundary": "RelatedBoundaryString",
					"charset":  "windows-874",
				},
			},
			ExtraHeaders: map[string][]string{
				"X-Clacks-Overhead": {"GNU Terry Pratchett"},
			},
		},
		Text: `เป็นมนุษย์สุดประเสริฐเลิศคุณค่า กว่าบรรดาฝูงสัตว์เดรัจฉาน จงฝ่าฟันพัฒนาวิชาการ อย่าล้างผลาญฤๅเข่นฆ่าบีฑาใคร ไม่ถือโทษโกรธแช่งซัดฮึดฮัดด่า หัดอภัยเหมือนกีฬาอัชฌาสัย ปฏิบัติประพฤติกฎกำหนดใจ พูดจาให้จ๊ะๆ จ๋าๆ น่าฟังเอยฯ

นายสังฆภัณฑ์ เฮงพิทักษ์ฝั่ง ผู้เฒ่าซึ่งมีอาชีพเป็นฅนขายฃวด ถูกตำรวจปฏิบัติการจับฟ้องศาล ฐานลักนาฬิกาคุณหญิงฉัตรชฎา ฌานสมาธิ`,
		EnrichedText: `<bold>เป็นมนุษย์สุดประเสริฐเลิศคุณค่า</bold> <italic>กว่าบรรดาฝูงสัตว์เดรัจฉาน</italic> <fixed>จงฝ่าฟันพัฒนาวิชาการ</fixed> <underline>อย่าล้างผลาญฤๅเข่นฆ่าบีฑาใคร</underline> ไม่ถือโทษโกรธแช่งซัดฮึดฮัดด่า หัดอภัยเหมือนกีฬาอัชฌาสัย ปฏิบัติประพฤติกฎกำหนดใจ พูดจาให้จ๊ะๆ จ๋าๆ น่าฟังเอยฯ

นายสังฆภัณฑ์ เฮงพิทักษ์ฝั่ง ผู้เฒ่าซึ่งมีอาชีพเป็นฅนขายฃวด ถูกตำรวจปฏิบัติการจับฟ้องศาล ฐานลักนาฬิกาคุณหญิงฉัตรชฎา ฌานสมาธิ`,
		HTML: `<html>
<div dir="ltr">
<p>เป็นมนุษย์สุดประเสริฐเลิศคุณค่า กว่าบรรดาฝูงสัตว์เดรัจฉาน จงฝ่าฟันพัฒนาวิชาการ อย่าล้างผลาญฤๅเข่นฆ่าบีฑาใคร ไม่ถือโทษโกรธแช่งซัดฮึดฮัดด่า หัดอภัยเหมือนกีฬาอัชฌาสัย ปฏิบัติประพฤติกฎกำหนดใจ พูดจาให้จ๊ะๆ จ๋าๆ น่าฟังเอยฯ</p>

<p>นายสังฆภัณฑ์ เฮงพิทักษ์ฝั่ง ผู้เฒ่าซึ่งมีอาชีพเป็นฅนขายฃวด ถูกตำรวจปฏิบัติการจับฟ้องศาล ฐานลักนาฬิกาคุณหญิงฉัตรชฎา ฌานสมาธิ</p>
</div>
</html>`,
		InlineFiles: []InlineFile{
			{
				ContentID: "inline-jpg-image.jpg@example.com",
				ContentType: ContentTypeHeader{
					ContentType: "image/jpeg",
					Params: map[string]string{
						"name": "inline-jpg-image-name.jpg",
					},
				},
				ContentDisposition: ContentDispositionHeader{
					ContentDisposition: inline,
					Params: map[string]string{
						"filename": "inline-jpg-image-filename.jpg",
					},
				},
				Data: []byte{255, 216, 255, 219, 0, 67, 0, 3, 2, 2, 2, 2, 2, 3, 2, 2, 2, 3, 3, 3, 3, 4,
					6, 4, 4, 4, 4, 4, 8, 6, 6, 5, 6, 9, 8, 10, 10, 9, 8, 9, 9, 10, 12, 15, 12, 10, 11, 14, 11, 9, 9,
					13, 17, 13, 14, 15, 16, 16, 17, 16, 10, 12, 18, 19, 18, 16, 19, 15, 16, 16, 16, 255, 201, 0, 11, 8,
					0, 1, 0, 1, 1, 1, 17, 0, 255, 204, 0, 6, 0, 16, 16, 5, 255, 218, 0, 8, 1, 1, 0, 0, 63, 0, 210, 207,
					32, 255, 217},
			},
		},
	}

	testEmailFromFile(t, fp, expectedEmail)
}
func TestParseEmailThaiMultipartRelatedWindows874OverQuotedprintable(t *testing.T) {
	fp := "tests/test_thai_multipart_related_windows-874_over_quoted-printable.txt"
	tz, _ := time.LoadLocation("Asia/Bangkok")
	expectedDate, _ := time.Parse(
		time.RFC1123Z+" (MST)",
		time.Date(2019, time.April, 1, 7, 55, 0, 0, tz).Format(time.RFC1123Z+" (MST)"))
	expectedEmail := Email{
		Headers: Headers{
			Date:    expectedDate,
			Subject: "📧 Test แพนแกรมภาษาไทย",
			ReplyTo: []*mail.Address{
				{
					Name:    "อลิซ ผู้ส่งจดหมาย",
					Address: "alis.phusngcdhmay@example.net",
				},
			},
			Sender: &mail.Address{
				Name:    "อลิซ ผู้ส่งจดหมาย",
				Address: "alis.phusngcdhmay@example.com",
			},
			From: []*mail.Address{
				{
					Name:    "อลิซ ผู้ส่งจดหมาย",
					Address: "alis.phusngcdhmay@example.com",
				},
				{
					Name:    "อลิซ ผู้ส่งจดหมาย",
					Address: "alis.phusngcdhmay@example.net",
				},
			},
			To: []*mail.Address{
				{
					Name:    "บ๊อบ ผู้รับ",
					Address: "bob.phurab@example.com",
				},
				{
					Name:    "คาโรล ผู้รับ",
					Address: "carol.phurab@example.com",
				},
			},
			Cc: []*mail.Address{
				{
					Name:    "แดน ผู้รับ",
					Address: "dan.phurab@example.com",
				},
				{
					Name:    "อีฟ ผู้รับ",
					Address: "eve.phurab@example.com",
				},
			},
			Bcc: []*mail.Address{
				{
					Name:    "แฟรงค์ ผู้รับ",
					Address: "frank.phurab@example.com",
				},
				{
					Name:    "เกรซ ผู้รับ",
					Address: "grace.phurab@example.com",
				},
			},
			MessageID:  "Message-Id-1@example.com",
			InReplyTo:  []MessageId{"Message-Id-0@example.com"},
			References: []MessageId{"Message-Id-0@example.com"},
			Comments:   "Message Header Comment",
			Keywords:   []string{"Keyword 1", "Keyword 2"},
			ResentDate: expectedDate,
			ResentFrom: []*mail.Address{
				{
					Name:    "อลิซ ผู้ส่งจดหมาย",
					Address: "alis.phusngcdhmay@example.net",
				},
				{
					Name:    "อลิซ ผู้ส่งจดหมาย",
					Address: "alis.phusngcdhmay@example.com",
				},
			},
			ResentSender: &mail.Address{
				Name:    "อลิซ ผู้ส่งจดหมาย",
				Address: "alis.phusngcdhmay@example.net",
			},
			ResentTo: []*mail.Address{
				{
					Name:    "บ๊อบ ผู้รับ",
					Address: "bob.phurab@example.net",
				},
				{
					Name:    "คาโรล ผู้รับ",
					Address: "carol.phurab@example.net",
				},
			},
			ResentCc: []*mail.Address{
				{
					Name:    "แดน ผู้รับ",
					Address: "dan.phurab@example.net",
				},
				{
					Name:    "อีฟ ผู้รับ",
					Address: "eve.phurab@example.net",
				},
			},
			ResentBcc: []*mail.Address{
				{
					Name:    "แฟรงค์ ผู้รับ",
					Address: "frank.phurab@example.net",
				},
				{
					Name:    "เกรซ ผู้รับ",
					Address: "grace.phurab@example.net",
				},
			},
			ResentMessageID: "Message-Id-1@example.net",
			ContentType: ContentTypeHeader{
				ContentType: "multipart/related",
				Params: map[string]string{
					"boundary": "RelatedBoundaryString",
					"charset":  "windows-874",
				},
			},
			ExtraHeaders: map[string][]string{
				"X-Clacks-Overhead": {"GNU Terry Pratchett"},
			},
		},
		Text: `เป็นมนุษย์สุดประเสริฐเลิศคุณค่า กว่าบรรดาฝูงสัตว์เดรัจฉาน จงฝ่าฟันพัฒนาวิชาการ อย่าล้างผลาญฤๅเข่นฆ่าบีฑาใคร ไม่ถือโทษโกรธแช่งซัดฮึดฮัดด่า หัดอภัยเหมือนกีฬาอัชฌาสัย ปฏิบัติประพฤติกฎกำหนดใจ พูดจาให้จ๊ะๆ จ๋าๆ น่าฟังเอยฯ

นายสังฆภัณฑ์ เฮงพิทักษ์ฝั่ง ผู้เฒ่าซึ่งมีอาชีพเป็นฅนขายฃวด ถูกตำรวจปฏิบัติการจับฟ้องศาล ฐานลักนาฬิกาคุณหญิงฉัตรชฎา ฌานสมาธิ`,
		EnrichedText: `<bold>เป็นมนุษย์สุดประเสริฐเลิศคุณค่า</bold> <italic>กว่าบรรดาฝูงสัตว์เดรัจฉาน</italic> <fixed>จงฝ่าฟันพัฒนาวิชาการ</fixed> <underline>อย่าล้างผลาญฤๅเข่นฆ่าบีฑาใคร</underline> ไม่ถือโทษโกรธแช่งซัดฮึดฮัดด่า หัดอภัยเหมือนกีฬาอัชฌาสัย ปฏิบัติประพฤติกฎกำหนดใจ พูดจาให้จ๊ะๆ จ๋าๆ น่าฟังเอยฯ

นายสังฆภัณฑ์ เฮงพิทักษ์ฝั่ง ผู้เฒ่าซึ่งมีอาชีพเป็นฅนขายฃวด ถูกตำรวจปฏิบัติการจับฟ้องศาล ฐานลักนาฬิกาคุณหญิงฉัตรชฎา ฌานสมาธิ`,
		HTML: `<html>
<div dir="ltr">
<p>เป็นมนุษย์สุดประเสริฐเลิศคุณค่า กว่าบรรดาฝูงสัตว์เดรัจฉาน จงฝ่าฟันพัฒนาวิชาการ อย่าล้างผลาญฤๅเข่นฆ่าบีฑาใคร ไม่ถือโทษโกรธแช่งซัดฮึดฮัดด่า หัดอภัยเหมือนกีฬาอัชฌาสัย ปฏิบัติประพฤติกฎกำหนดใจ พูดจาให้จ๊ะๆ จ๋าๆ น่าฟังเอยฯ</p>

<p>นายสังฆภัณฑ์ เฮงพิทักษ์ฝั่ง ผู้เฒ่าซึ่งมีอาชีพเป็นฅนขายฃวด ถูกตำรวจปฏิบัติการจับฟ้องศาล ฐานลักนาฬิกาคุณหญิงฉัตรชฎา ฌานสมาธิ</p>
</div>
</html>`,
		InlineFiles: []InlineFile{
			{
				ContentID: "inline-jpg-image.jpg@example.com",
				ContentType: ContentTypeHeader{
					ContentType: "image/jpeg",
					Params: map[string]string{
						"name": "inline-jpg-image-name.jpg",
					},
				},
				ContentDisposition: ContentDispositionHeader{
					ContentDisposition: inline,
					Params: map[string]string{
						"filename": "inline-jpg-image-filename.jpg",
					},
				},
				Data: []byte{255, 216, 255, 219, 0, 67, 0, 3, 2, 2, 2, 2, 2, 3, 2, 2, 2, 3, 3, 3, 3, 4,
					6, 4, 4, 4, 4, 4, 8, 6, 6, 5, 6, 9, 8, 10, 10, 9, 8, 9, 9, 10, 12, 15, 12, 10, 11, 14, 11, 9, 9,
					13, 17, 13, 14, 15, 16, 16, 17, 16, 10, 12, 18, 19, 18, 16, 19, 15, 16, 16, 16, 255, 201, 0, 11, 8,
					0, 1, 0, 1, 1, 1, 17, 0, 255, 204, 0, 6, 0, 16, 16, 5, 255, 218, 0, 8, 1, 1, 0, 0, 63, 0, 210, 207,
					32, 255, 217},
			},
		},
	}

	testEmailFromFile(t, fp, expectedEmail)
}
func TestParseEmailThaiMultipartRelatedTis620OverBase64(t *testing.T) {
	fp := "tests/test_thai_multipart_related_tis-620_over_base64.txt"
	tz, _ := time.LoadLocation("Asia/Bangkok")
	expectedDate, _ := time.Parse(
		time.RFC1123Z+" (MST)",
		time.Date(2019, time.April, 1, 7, 55, 0, 0, tz).Format(time.RFC1123Z+" (MST)"))
	expectedEmail := Email{
		Headers: Headers{
			Date:    expectedDate,
			Subject: "📧 Test แพนแกรมภาษาไทย",
			ReplyTo: []*mail.Address{
				{
					Name:    "อลิซ ผู้ส่งจดหมาย",
					Address: "alis.phusngcdhmay@example.net",
				},
			},
			Sender: &mail.Address{
				Name:    "อลิซ ผู้ส่งจดหมาย",
				Address: "alis.phusngcdhmay@example.com",
			},
			From: []*mail.Address{
				{
					Name:    "อลิซ ผู้ส่งจดหมาย",
					Address: "alis.phusngcdhmay@example.com",
				},
				{
					Name:    "อลิซ ผู้ส่งจดหมาย",
					Address: "alis.phusngcdhmay@example.net",
				},
			},
			To: []*mail.Address{
				{
					Name:    "บ๊อบ ผู้รับ",
					Address: "bob.phurab@example.com",
				},
				{
					Name:    "คาโรล ผู้รับ",
					Address: "carol.phurab@example.com",
				},
			},
			Cc: []*mail.Address{
				{
					Name:    "แดน ผู้รับ",
					Address: "dan.phurab@example.com",
				},
				{
					Name:    "อีฟ ผู้รับ",
					Address: "eve.phurab@example.com",
				},
			},
			Bcc: []*mail.Address{
				{
					Name:    "แฟรงค์ ผู้รับ",
					Address: "frank.phurab@example.com",
				},
				{
					Name:    "เกรซ ผู้รับ",
					Address: "grace.phurab@example.com",
				},
			},
			MessageID:  "Message-Id-1@example.com",
			InReplyTo:  []MessageId{"Message-Id-0@example.com"},
			References: []MessageId{"Message-Id-0@example.com"},
			Comments:   "Message Header Comment",
			Keywords:   []string{"Keyword 1", "Keyword 2"},
			ResentDate: expectedDate,
			ResentFrom: []*mail.Address{
				{
					Name:    "อลิซ ผู้ส่งจดหมาย",
					Address: "alis.phusngcdhmay@example.net",
				},
				{
					Name:    "อลิซ ผู้ส่งจดหมาย",
					Address: "alis.phusngcdhmay@example.com",
				},
			},
			ResentSender: &mail.Address{
				Name:    "อลิซ ผู้ส่งจดหมาย",
				Address: "alis.phusngcdhmay@example.net",
			},
			ResentTo: []*mail.Address{
				{
					Name:    "บ๊อบ ผู้รับ",
					Address: "bob.phurab@example.net",
				},
				{
					Name:    "คาโรล ผู้รับ",
					Address: "carol.phurab@example.net",
				},
			},
			ResentCc: []*mail.Address{
				{
					Name:    "แดน ผู้รับ",
					Address: "dan.phurab@example.net",
				},
				{
					Name:    "อีฟ ผู้รับ",
					Address: "eve.phurab@example.net",
				},
			},
			ResentBcc: []*mail.Address{
				{
					Name:    "แฟรงค์ ผู้รับ",
					Address: "frank.phurab@example.net",
				},
				{
					Name:    "เกรซ ผู้รับ",
					Address: "grace.phurab@example.net",
				},
			},
			ResentMessageID: "Message-Id-1@example.net",
			ContentType: ContentTypeHeader{
				ContentType: "multipart/related",
				Params: map[string]string{
					"boundary": "RelatedBoundaryString",
					"charset":  "tis-620",
				},
			},
			ExtraHeaders: map[string][]string{
				"X-Clacks-Overhead": {"GNU Terry Pratchett"},
			},
		},
		Text: `เป็นมนุษย์สุดประเสริฐเลิศคุณค่า กว่าบรรดาฝูงสัตว์เดรัจฉาน จงฝ่าฟันพัฒนาวิชาการ อย่าล้างผลาญฤๅเข่นฆ่าบีฑาใคร ไม่ถือโทษโกรธแช่งซัดฮึดฮัดด่า หัดอภัยเหมือนกีฬาอัชฌาสัย ปฏิบัติประพฤติกฎกำหนดใจ พูดจาให้จ๊ะๆ จ๋าๆ น่าฟังเอยฯ

นายสังฆภัณฑ์ เฮงพิทักษ์ฝั่ง ผู้เฒ่าซึ่งมีอาชีพเป็นฅนขายฃวด ถูกตำรวจปฏิบัติการจับฟ้องศาล ฐานลักนาฬิกาคุณหญิงฉัตรชฎา ฌานสมาธิ`,
		EnrichedText: `<bold>เป็นมนุษย์สุดประเสริฐเลิศคุณค่า</bold> <italic>กว่าบรรดาฝูงสัตว์เดรัจฉาน</italic> <fixed>จงฝ่าฟันพัฒนาวิชาการ</fixed> <underline>อย่าล้างผลาญฤๅเข่นฆ่าบีฑาใคร</underline> ไม่ถือโทษโกรธแช่งซัดฮึดฮัดด่า หัดอภัยเหมือนกีฬาอัชฌาสัย ปฏิบัติประพฤติกฎกำหนดใจ พูดจาให้จ๊ะๆ จ๋าๆ น่าฟังเอยฯ

นายสังฆภัณฑ์ เฮงพิทักษ์ฝั่ง ผู้เฒ่าซึ่งมีอาชีพเป็นฅนขายฃวด ถูกตำรวจปฏิบัติการจับฟ้องศาล ฐานลักนาฬิกาคุณหญิงฉัตรชฎา ฌานสมาธิ`,
		HTML: `<html>
<div dir="ltr">
<p>เป็นมนุษย์สุดประเสริฐเลิศคุณค่า กว่าบรรดาฝูงสัตว์เดรัจฉาน จงฝ่าฟันพัฒนาวิชาการ อย่าล้างผลาญฤๅเข่นฆ่าบีฑาใคร ไม่ถือโทษโกรธแช่งซัดฮึดฮัดด่า หัดอภัยเหมือนกีฬาอัชฌาสัย ปฏิบัติประพฤติกฎกำหนดใจ พูดจาให้จ๊ะๆ จ๋าๆ น่าฟังเอยฯ</p>

<p>นายสังฆภัณฑ์ เฮงพิทักษ์ฝั่ง ผู้เฒ่าซึ่งมีอาชีพเป็นฅนขายฃวด ถูกตำรวจปฏิบัติการจับฟ้องศาล ฐานลักนาฬิกาคุณหญิงฉัตรชฎา ฌานสมาธิ</p>
</div>
</html>`,
		InlineFiles: []InlineFile{
			{
				ContentID: "inline-jpg-image.jpg@example.com",
				ContentType: ContentTypeHeader{
					ContentType: "image/jpeg",
					Params: map[string]string{
						"name": "inline-jpg-image-name.jpg",
					},
				},
				ContentDisposition: ContentDispositionHeader{
					ContentDisposition: inline,
					Params: map[string]string{
						"filename": "inline-jpg-image-filename.jpg",
					},
				},
				Data: []byte{255, 216, 255, 219, 0, 67, 0, 3, 2, 2, 2, 2, 2, 3, 2, 2, 2, 3, 3, 3, 3, 4,
					6, 4, 4, 4, 4, 4, 8, 6, 6, 5, 6, 9, 8, 10, 10, 9, 8, 9, 9, 10, 12, 15, 12, 10, 11, 14, 11, 9, 9,
					13, 17, 13, 14, 15, 16, 16, 17, 16, 10, 12, 18, 19, 18, 16, 19, 15, 16, 16, 16, 255, 201, 0, 11, 8,
					0, 1, 0, 1, 1, 1, 17, 0, 255, 204, 0, 6, 0, 16, 16, 5, 255, 218, 0, 8, 1, 1, 0, 0, 63, 0, 210, 207,
					32, 255, 217},
			},
		},
	}

	testEmailFromFile(t, fp, expectedEmail)
}
func TestParseEmailThaiMultipartRelatedTis620OverQuotedprintable(t *testing.T) {
	fp := "tests/test_thai_multipart_related_tis-620_over_quoted-printable.txt"
	tz, _ := time.LoadLocation("Asia/Bangkok")
	expectedDate, _ := time.Parse(
		time.RFC1123Z+" (MST)",
		time.Date(2019, time.April, 1, 7, 55, 0, 0, tz).Format(time.RFC1123Z+" (MST)"))
	expectedEmail := Email{
		Headers: Headers{
			Date:    expectedDate,
			Subject: "📧 Test แพนแกรมภาษาไทย",
			ReplyTo: []*mail.Address{
				{
					Name:    "อลิซ ผู้ส่งจดหมาย",
					Address: "alis.phusngcdhmay@example.net",
				},
			},
			Sender: &mail.Address{
				Name:    "อลิซ ผู้ส่งจดหมาย",
				Address: "alis.phusngcdhmay@example.com",
			},
			From: []*mail.Address{
				{
					Name:    "อลิซ ผู้ส่งจดหมาย",
					Address: "alis.phusngcdhmay@example.com",
				},
				{
					Name:    "อลิซ ผู้ส่งจดหมาย",
					Address: "alis.phusngcdhmay@example.net",
				},
			},
			To: []*mail.Address{
				{
					Name:    "บ๊อบ ผู้รับ",
					Address: "bob.phurab@example.com",
				},
				{
					Name:    "คาโรล ผู้รับ",
					Address: "carol.phurab@example.com",
				},
			},
			Cc: []*mail.Address{
				{
					Name:    "แดน ผู้รับ",
					Address: "dan.phurab@example.com",
				},
				{
					Name:    "อีฟ ผู้รับ",
					Address: "eve.phurab@example.com",
				},
			},
			Bcc: []*mail.Address{
				{
					Name:    "แฟรงค์ ผู้รับ",
					Address: "frank.phurab@example.com",
				},
				{
					Name:    "เกรซ ผู้รับ",
					Address: "grace.phurab@example.com",
				},
			},
			MessageID:  "Message-Id-1@example.com",
			InReplyTo:  []MessageId{"Message-Id-0@example.com"},
			References: []MessageId{"Message-Id-0@example.com"},
			Comments:   "Message Header Comment",
			Keywords:   []string{"Keyword 1", "Keyword 2"},
			ResentDate: expectedDate,
			ResentFrom: []*mail.Address{
				{
					Name:    "อลิซ ผู้ส่งจดหมาย",
					Address: "alis.phusngcdhmay@example.net",
				},
				{
					Name:    "อลิซ ผู้ส่งจดหมาย",
					Address: "alis.phusngcdhmay@example.com",
				},
			},
			ResentSender: &mail.Address{
				Name:    "อลิซ ผู้ส่งจดหมาย",
				Address: "alis.phusngcdhmay@example.net",
			},
			ResentTo: []*mail.Address{
				{
					Name:    "บ๊อบ ผู้รับ",
					Address: "bob.phurab@example.net",
				},
				{
					Name:    "คาโรล ผู้รับ",
					Address: "carol.phurab@example.net",
				},
			},
			ResentCc: []*mail.Address{
				{
					Name:    "แดน ผู้รับ",
					Address: "dan.phurab@example.net",
				},
				{
					Name:    "อีฟ ผู้รับ",
					Address: "eve.phurab@example.net",
				},
			},
			ResentBcc: []*mail.Address{
				{
					Name:    "แฟรงค์ ผู้รับ",
					Address: "frank.phurab@example.net",
				},
				{
					Name:    "เกรซ ผู้รับ",
					Address: "grace.phurab@example.net",
				},
			},
			ResentMessageID: "Message-Id-1@example.net",
			ContentType: ContentTypeHeader{
				ContentType: "multipart/related",
				Params: map[string]string{
					"boundary": "RelatedBoundaryString",
					"charset":  "tis-620",
				},
			},
			ExtraHeaders: map[string][]string{
				"X-Clacks-Overhead": {"GNU Terry Pratchett"},
			},
		},
		Text: `เป็นมนุษย์สุดประเสริฐเลิศคุณค่า กว่าบรรดาฝูงสัตว์เดรัจฉาน จงฝ่าฟันพัฒนาวิชาการ อย่าล้างผลาญฤๅเข่นฆ่าบีฑาใคร ไม่ถือโทษโกรธแช่งซัดฮึดฮัดด่า หัดอภัยเหมือนกีฬาอัชฌาสัย ปฏิบัติประพฤติกฎกำหนดใจ พูดจาให้จ๊ะๆ จ๋าๆ น่าฟังเอยฯ

นายสังฆภัณฑ์ เฮงพิทักษ์ฝั่ง ผู้เฒ่าซึ่งมีอาชีพเป็นฅนขายฃวด ถูกตำรวจปฏิบัติการจับฟ้องศาล ฐานลักนาฬิกาคุณหญิงฉัตรชฎา ฌานสมาธิ`,
		EnrichedText: `<bold>เป็นมนุษย์สุดประเสริฐเลิศคุณค่า</bold> <italic>กว่าบรรดาฝูงสัตว์เดรัจฉาน</italic> <fixed>จงฝ่าฟันพัฒนาวิชาการ</fixed> <underline>อย่าล้างผลาญฤๅเข่นฆ่าบีฑาใคร</underline> ไม่ถือโทษโกรธแช่งซัดฮึดฮัดด่า หัดอภัยเหมือนกีฬาอัชฌาสัย ปฏิบัติประพฤติกฎกำหนดใจ พูดจาให้จ๊ะๆ จ๋าๆ น่าฟังเอยฯ

นายสังฆภัณฑ์ เฮงพิทักษ์ฝั่ง ผู้เฒ่าซึ่งมีอาชีพเป็นฅนขายฃวด ถูกตำรวจปฏิบัติการจับฟ้องศาล ฐานลักนาฬิกาคุณหญิงฉัตรชฎา ฌานสมาธิ`,
		HTML: `<html>
<div dir="ltr">
<p>เป็นมนุษย์สุดประเสริฐเลิศคุณค่า กว่าบรรดาฝูงสัตว์เดรัจฉาน จงฝ่าฟันพัฒนาวิชาการ อย่าล้างผลาญฤๅเข่นฆ่าบีฑาใคร ไม่ถือโทษโกรธแช่งซัดฮึดฮัดด่า หัดอภัยเหมือนกีฬาอัชฌาสัย ปฏิบัติประพฤติกฎกำหนดใจ พูดจาให้จ๊ะๆ จ๋าๆ น่าฟังเอยฯ</p>

<p>นายสังฆภัณฑ์ เฮงพิทักษ์ฝั่ง ผู้เฒ่าซึ่งมีอาชีพเป็นฅนขายฃวด ถูกตำรวจปฏิบัติการจับฟ้องศาล ฐานลักนาฬิกาคุณหญิงฉัตรชฎา ฌานสมาธิ</p>
</div>
</html>`,
		InlineFiles: []InlineFile{
			{
				ContentID: "inline-jpg-image.jpg@example.com",
				ContentType: ContentTypeHeader{
					ContentType: "image/jpeg",
					Params: map[string]string{
						"name": "inline-jpg-image-name.jpg",
					},
				},
				ContentDisposition: ContentDispositionHeader{
					ContentDisposition: inline,
					Params: map[string]string{
						"filename": "inline-jpg-image-filename.jpg",
					},
				},
				Data: []byte{255, 216, 255, 219, 0, 67, 0, 3, 2, 2, 2, 2, 2, 3, 2, 2, 2, 3, 3, 3, 3, 4,
					6, 4, 4, 4, 4, 4, 8, 6, 6, 5, 6, 9, 8, 10, 10, 9, 8, 9, 9, 10, 12, 15, 12, 10, 11, 14, 11, 9, 9,
					13, 17, 13, 14, 15, 16, 16, 17, 16, 10, 12, 18, 19, 18, 16, 19, 15, 16, 16, 16, 255, 201, 0, 11, 8,
					0, 1, 0, 1, 1, 1, 17, 0, 255, 204, 0, 6, 0, 16, 16, 5, 255, 218, 0, 8, 1, 1, 0, 0, 63, 0, 210, 207,
					32, 255, 217},
			},
		},
	}

	testEmailFromFile(t, fp, expectedEmail)
}
func TestParseEmailThaiMultipartMixedIso885911OverBase64(t *testing.T) {
	fp := "tests/test_thai_multipart_mixed_iso-8859-11_over_base64.txt"
	tz, _ := time.LoadLocation("Asia/Bangkok")
	expectedDate, _ := time.Parse(
		time.RFC1123Z+" (MST)",
		time.Date(2019, time.April, 1, 7, 55, 0, 0, tz).Format(time.RFC1123Z+" (MST)"))
	expectedEmail := Email{
		Headers: Headers{
			Date:    expectedDate,
			Subject: "📧 Test แพนแกรมภาษาไทย",
			ReplyTo: []*mail.Address{
				{
					Name:    "อลิซ ผู้ส่งจดหมาย",
					Address: "alis.phusngcdhmay@example.net",
				},
			},
			Sender: &mail.Address{
				Name:    "อลิซ ผู้ส่งจดหมาย",
				Address: "alis.phusngcdhmay@example.com",
			},
			From: []*mail.Address{
				{
					Name:    "อลิซ ผู้ส่งจดหมาย",
					Address: "alis.phusngcdhmay@example.com",
				},
				{
					Name:    "อลิซ ผู้ส่งจดหมาย",
					Address: "alis.phusngcdhmay@example.net",
				},
			},
			To: []*mail.Address{
				{
					Name:    "บ๊อบ ผู้รับ",
					Address: "bob.phurab@example.com",
				},
				{
					Name:    "คาโรล ผู้รับ",
					Address: "carol.phurab@example.com",
				},
			},
			Cc: []*mail.Address{
				{
					Name:    "แดน ผู้รับ",
					Address: "dan.phurab@example.com",
				},
				{
					Name:    "อีฟ ผู้รับ",
					Address: "eve.phurab@example.com",
				},
			},
			Bcc: []*mail.Address{
				{
					Name:    "แฟรงค์ ผู้รับ",
					Address: "frank.phurab@example.com",
				},
				{
					Name:    "เกรซ ผู้รับ",
					Address: "grace.phurab@example.com",
				},
			},
			MessageID:  "Message-Id-1@example.com",
			InReplyTo:  []MessageId{"Message-Id-0@example.com"},
			References: []MessageId{"Message-Id-0@example.com"},
			Comments:   "Message Header Comment",
			Keywords:   []string{"Keyword 1", "Keyword 2"},
			ResentDate: expectedDate,
			ResentFrom: []*mail.Address{
				{
					Name:    "อลิซ ผู้ส่งจดหมาย",
					Address: "alis.phusngcdhmay@example.net",
				},
				{
					Name:    "อลิซ ผู้ส่งจดหมาย",
					Address: "alis.phusngcdhmay@example.com",
				},
			},
			ResentSender: &mail.Address{
				Name:    "อลิซ ผู้ส่งจดหมาย",
				Address: "alis.phusngcdhmay@example.net",
			},
			ResentTo: []*mail.Address{
				{
					Name:    "บ๊อบ ผู้รับ",
					Address: "bob.phurab@example.net",
				},
				{
					Name:    "คาโรล ผู้รับ",
					Address: "carol.phurab@example.net",
				},
			},
			ResentCc: []*mail.Address{
				{
					Name:    "แดน ผู้รับ",
					Address: "dan.phurab@example.net",
				},
				{
					Name:    "อีฟ ผู้รับ",
					Address: "eve.phurab@example.net",
				},
			},
			ResentBcc: []*mail.Address{
				{
					Name:    "แฟรงค์ ผู้รับ",
					Address: "frank.phurab@example.net",
				},
				{
					Name:    "เกรซ ผู้รับ",
					Address: "grace.phurab@example.net",
				},
			},
			ResentMessageID: "Message-Id-1@example.net",
			ContentType: ContentTypeHeader{
				ContentType: "multipart/mixed",
				Params: map[string]string{
					"boundary": "MixedBoundaryString",
					"charset":  "iso-8859-11",
				},
			},
			ExtraHeaders: map[string][]string{
				"X-Clacks-Overhead": {"GNU Terry Pratchett"},
			},
		},
		Text: `เป็นมนุษย์สุดประเสริฐเลิศคุณค่า กว่าบรรดาฝูงสัตว์เดรัจฉาน จงฝ่าฟันพัฒนาวิชาการ อย่าล้างผลาญฤๅเข่นฆ่าบีฑาใคร ไม่ถือโทษโกรธแช่งซัดฮึดฮัดด่า หัดอภัยเหมือนกีฬาอัชฌาสัย ปฏิบัติประพฤติกฎกำหนดใจ พูดจาให้จ๊ะๆ จ๋าๆ น่าฟังเอยฯ

นายสังฆภัณฑ์ เฮงพิทักษ์ฝั่ง ผู้เฒ่าซึ่งมีอาชีพเป็นฅนขายฃวด ถูกตำรวจปฏิบัติการจับฟ้องศาล ฐานลักนาฬิกาคุณหญิงฉัตรชฎา ฌานสมาธิ`,
		EnrichedText: `<bold>เป็นมนุษย์สุดประเสริฐเลิศคุณค่า</bold> <italic>กว่าบรรดาฝูงสัตว์เดรัจฉาน</italic> <fixed>จงฝ่าฟันพัฒนาวิชาการ</fixed> <underline>อย่าล้างผลาญฤๅเข่นฆ่าบีฑาใคร</underline> ไม่ถือโทษโกรธแช่งซัดฮึดฮัดด่า หัดอภัยเหมือนกีฬาอัชฌาสัย ปฏิบัติประพฤติกฎกำหนดใจ พูดจาให้จ๊ะๆ จ๋าๆ น่าฟังเอยฯ

นายสังฆภัณฑ์ เฮงพิทักษ์ฝั่ง ผู้เฒ่าซึ่งมีอาชีพเป็นฅนขายฃวด ถูกตำรวจปฏิบัติการจับฟ้องศาล ฐานลักนาฬิกาคุณหญิงฉัตรชฎา ฌานสมาธิ`,
		HTML: `<html>
<div dir="ltr">
<p>เป็นมนุษย์สุดประเสริฐเลิศคุณค่า กว่าบรรดาฝูงสัตว์เดรัจฉาน จงฝ่าฟันพัฒนาวิชาการ อย่าล้างผลาญฤๅเข่นฆ่าบีฑาใคร ไม่ถือโทษโกรธแช่งซัดฮึดฮัดด่า หัดอภัยเหมือนกีฬาอัชฌาสัย ปฏิบัติประพฤติกฎกำหนดใจ พูดจาให้จ๊ะๆ จ๋าๆ น่าฟังเอยฯ</p>

<p>นายสังฆภัณฑ์ เฮงพิทักษ์ฝั่ง ผู้เฒ่าซึ่งมีอาชีพเป็นฅนขายฃวด ถูกตำรวจปฏิบัติการจับฟ้องศาล ฐานลักนาฬิกาคุณหญิงฉัตรชฎา ฌานสมาธิ</p>
</div>
</html>`,
		InlineFiles: []InlineFile{
			{
				ContentType: ContentTypeHeader{
					ContentType: "image/jpeg",
					Params: map[string]string{
						"name": "inline-jpg-image-without-disposition.jpg",
					},
				},
				ContentDisposition: ContentDispositionHeader{
					ContentDisposition: "",
					Params:             map[string]string(nil),
				},
				Data: []byte{255, 216, 255, 219, 0, 67, 0, 3, 2, 2, 2, 2, 2, 3, 2, 2, 2, 3, 3, 3, 3, 4,
					6, 4, 4, 4, 4, 4, 8, 6, 6, 5, 6, 9, 8, 10, 10, 9, 8, 9, 9, 10, 12, 15, 12, 10, 11, 14, 11, 9, 9,
					13, 17, 13, 14, 15, 16, 16, 17, 16, 10, 12, 18, 19, 18, 16, 19, 15, 16, 16, 16, 255, 201, 0, 11, 8,
					0, 1, 0, 1, 1, 1, 17, 0, 255, 204, 0, 6, 0, 16, 16, 5, 255, 218, 0, 8, 1, 1, 0, 0, 63, 0, 210, 207,
					32, 255, 217},
			},
			{
				ContentID: "inline-jpg-image.jpg@example.com",
				ContentType: ContentTypeHeader{
					ContentType: "image/jpeg",
					Params: map[string]string{
						"name": "inline-jpg-image-name.jpg",
					},
				},
				ContentDisposition: ContentDispositionHeader{
					ContentDisposition: inline,
					Params: map[string]string{
						"filename": "inline-jpg-image-filename.jpg",
					},
				},
				Data: []byte{255, 216, 255, 219, 0, 67, 0, 3, 2, 2, 2, 2, 2, 3, 2, 2, 2, 3, 3, 3, 3, 4,
					6, 4, 4, 4, 4, 4, 8, 6, 6, 5, 6, 9, 8, 10, 10, 9, 8, 9, 9, 10, 12, 15, 12, 10, 11, 14, 11, 9, 9,
					13, 17, 13, 14, 15, 16, 16, 17, 16, 10, 12, 18, 19, 18, 16, 19, 15, 16, 16, 16, 255, 201, 0, 11, 8,
					0, 1, 0, 1, 1, 1, 17, 0, 255, 204, 0, 6, 0, 16, 16, 5, 255, 218, 0, 8, 1, 1, 0, 0, 63, 0, 210, 207,
					32, 255, 217},
			},
		},
		AttachedFiles: []AttachedFile{
			{
				ContentType: ContentTypeHeader{
					ContentType: "application/pdf",
					Params: map[string]string{
						"name": "attached-pdf-name.pdf",
					},
				},
				ContentDisposition: ContentDispositionHeader{
					ContentDisposition: attachment,
					Params: map[string]string{
						"filename": "attached-pdf-filename.pdf",
					},
				},
				Data: []byte{37, 80, 68, 70, 45, 49, 46, 13, 116, 114, 97, 105, 108, 101, 114, 60, 60,
					47, 82, 111, 111, 116, 60, 60, 47, 80, 97, 103, 101, 115, 60, 60, 47, 75, 105, 100, 115, 91, 60,
					60, 47, 77, 101, 100, 105, 97, 66, 111, 120, 91, 48, 32, 48, 32, 51, 32, 51, 93, 62, 62, 93, 62,
					62, 62, 62, 62, 62},
			},
			{
				ContentType: ContentTypeHeader{
					ContentType: "application/pdf",
					Params: map[string]string{
						"name": "attached-pdf-without-disposition.pdf",
					},
				},
				ContentDisposition: ContentDispositionHeader{
					ContentDisposition: "",
					Params:             map[string]string(nil),
				},
				Data: []byte{37, 80, 68, 70, 45, 49, 46, 13, 116, 114, 97, 105, 108, 101, 114, 60, 60,
					47, 82, 111, 111, 116, 60, 60, 47, 80, 97, 103, 101, 115, 60, 60, 47, 75, 105, 100, 115, 91, 60,
					60, 47, 77, 101, 100, 105, 97, 66, 111, 120, 91, 48, 32, 48, 32, 51, 32, 51, 93, 62, 62, 93, 62,
					62, 62, 62, 62, 62},
			},
			{
				ContentType: ContentTypeHeader{
					ContentType: "application/json",
					Params: map[string]string{
						"name": "attached-json-name.json",
					},
				},
				ContentDisposition: ContentDispositionHeader{
					ContentDisposition: attachment,
					Params: map[string]string{
						"filename": "attached-json-filename.json",
					},
				},
				Data: []byte{123, 34, 102, 111, 111, 34, 58, 34, 98, 97, 114, 34, 125},
			},
			{
				ContentType: ContentTypeHeader{
					ContentType: "text/plain",
					Params: map[string]string{
						"name": "attached-text-plain-name.txt",
					},
				},
				ContentDisposition: ContentDispositionHeader{
					ContentDisposition: attachment,
					Params: map[string]string{
						"filename": "attached-text-plain-filename.txt",
					},
				},
				Data: []byte{84, 101, 120, 116, 47, 112, 108, 97, 105, 110, 32, 99, 111, 110, 116, 101, 110, 116, 32,
					97, 115, 32, 97, 110, 32, 97, 116, 116, 97, 99, 104, 101, 100, 32, 46, 116, 120, 116, 32, 102, 105,
					108, 101, 46},
			},
			{
				ContentType: ContentTypeHeader{
					ContentType: "text/html",
					Params: map[string]string{
						"name": "attached-text-html-name.html",
					},
				},
				ContentDisposition: ContentDispositionHeader{
					ContentDisposition: attachment,
					Params: map[string]string{
						"filename": "attached-text-html-filename.html",
					},
				},
				Data: []byte{84, 101, 120, 116, 47, 104, 116, 109, 108, 32, 99, 111, 110, 116, 101, 110, 116, 32, 97,
					115, 32, 97, 110, 32, 97, 116, 116, 97, 99, 104, 101, 100, 32, 46, 104, 116, 109, 108, 32, 102,
					105, 108, 101, 46},
			},
		},
	}

	testEmailFromFile(t, fp, expectedEmail)
}
func TestParseEmailThaiMultipartMixedIso885911OverQuotedprintable(t *testing.T) {
	fp := "tests/test_thai_multipart_mixed_iso-8859-11_over_quoted-printable.txt"
	tz, _ := time.LoadLocation("Asia/Bangkok")
	expectedDate, _ := time.Parse(
		time.RFC1123Z+" (MST)",
		time.Date(2019, time.April, 1, 7, 55, 0, 0, tz).Format(time.RFC1123Z+" (MST)"))
	expectedEmail := Email{
		Headers: Headers{
			Date:    expectedDate,
			Subject: "📧 Test แพนแกรมภาษาไทย",
			ReplyTo: []*mail.Address{
				{
					Name:    "อลิซ ผู้ส่งจดหมาย",
					Address: "alis.phusngcdhmay@example.net",
				},
			},
			Sender: &mail.Address{
				Name:    "อลิซ ผู้ส่งจดหมาย",
				Address: "alis.phusngcdhmay@example.com",
			},
			From: []*mail.Address{
				{
					Name:    "อลิซ ผู้ส่งจดหมาย",
					Address: "alis.phusngcdhmay@example.com",
				},
				{
					Name:    "อลิซ ผู้ส่งจดหมาย",
					Address: "alis.phusngcdhmay@example.net",
				},
			},
			To: []*mail.Address{
				{
					Name:    "บ๊อบ ผู้รับ",
					Address: "bob.phurab@example.com",
				},
				{
					Name:    "คาโรล ผู้รับ",
					Address: "carol.phurab@example.com",
				},
			},
			Cc: []*mail.Address{
				{
					Name:    "แดน ผู้รับ",
					Address: "dan.phurab@example.com",
				},
				{
					Name:    "อีฟ ผู้รับ",
					Address: "eve.phurab@example.com",
				},
			},
			Bcc: []*mail.Address{
				{
					Name:    "แฟรงค์ ผู้รับ",
					Address: "frank.phurab@example.com",
				},
				{
					Name:    "เกรซ ผู้รับ",
					Address: "grace.phurab@example.com",
				},
			},
			MessageID:  "Message-Id-1@example.com",
			InReplyTo:  []MessageId{"Message-Id-0@example.com"},
			References: []MessageId{"Message-Id-0@example.com"},
			Comments:   "Message Header Comment",
			Keywords:   []string{"Keyword 1", "Keyword 2"},
			ResentDate: expectedDate,
			ResentFrom: []*mail.Address{
				{
					Name:    "อลิซ ผู้ส่งจดหมาย",
					Address: "alis.phusngcdhmay@example.net",
				},
				{
					Name:    "อลิซ ผู้ส่งจดหมาย",
					Address: "alis.phusngcdhmay@example.com",
				},
			},
			ResentSender: &mail.Address{
				Name:    "อลิซ ผู้ส่งจดหมาย",
				Address: "alis.phusngcdhmay@example.net",
			},
			ResentTo: []*mail.Address{
				{
					Name:    "บ๊อบ ผู้รับ",
					Address: "bob.phurab@example.net",
				},
				{
					Name:    "คาโรล ผู้รับ",
					Address: "carol.phurab@example.net",
				},
			},
			ResentCc: []*mail.Address{
				{
					Name:    "แดน ผู้รับ",
					Address: "dan.phurab@example.net",
				},
				{
					Name:    "อีฟ ผู้รับ",
					Address: "eve.phurab@example.net",
				},
			},
			ResentBcc: []*mail.Address{
				{
					Name:    "แฟรงค์ ผู้รับ",
					Address: "frank.phurab@example.net",
				},
				{
					Name:    "เกรซ ผู้รับ",
					Address: "grace.phurab@example.net",
				},
			},
			ResentMessageID: "Message-Id-1@example.net",
			ContentType: ContentTypeHeader{
				ContentType: "multipart/mixed",
				Params: map[string]string{
					"boundary": "MixedBoundaryString",
					"charset":  "iso-8859-11",
				},
			},
			ExtraHeaders: map[string][]string{
				"X-Clacks-Overhead": {"GNU Terry Pratchett"},
			},
		},
		Text: `เป็นมนุษย์สุดประเสริฐเลิศคุณค่า กว่าบรรดาฝูงสัตว์เดรัจฉาน จงฝ่าฟันพัฒนาวิชาการ อย่าล้างผลาญฤๅเข่นฆ่าบีฑาใคร ไม่ถือโทษโกรธแช่งซัดฮึดฮัดด่า หัดอภัยเหมือนกีฬาอัชฌาสัย ปฏิบัติประพฤติกฎกำหนดใจ พูดจาให้จ๊ะๆ จ๋าๆ น่าฟังเอยฯ

นายสังฆภัณฑ์ เฮงพิทักษ์ฝั่ง ผู้เฒ่าซึ่งมีอาชีพเป็นฅนขายฃวด ถูกตำรวจปฏิบัติการจับฟ้องศาล ฐานลักนาฬิกาคุณหญิงฉัตรชฎา ฌานสมาธิ`,
		EnrichedText: `<bold>เป็นมนุษย์สุดประเสริฐเลิศคุณค่า</bold> <italic>กว่าบรรดาฝูงสัตว์เดรัจฉาน</italic> <fixed>จงฝ่าฟันพัฒนาวิชาการ</fixed> <underline>อย่าล้างผลาญฤๅเข่นฆ่าบีฑาใคร</underline> ไม่ถือโทษโกรธแช่งซัดฮึดฮัดด่า หัดอภัยเหมือนกีฬาอัชฌาสัย ปฏิบัติประพฤติกฎกำหนดใจ พูดจาให้จ๊ะๆ จ๋าๆ น่าฟังเอยฯ

นายสังฆภัณฑ์ เฮงพิทักษ์ฝั่ง ผู้เฒ่าซึ่งมีอาชีพเป็นฅนขายฃวด ถูกตำรวจปฏิบัติการจับฟ้องศาล ฐานลักนาฬิกาคุณหญิงฉัตรชฎา ฌานสมาธิ`,
		HTML: `<html>
<div dir="ltr">
<p>เป็นมนุษย์สุดประเสริฐเลิศคุณค่า กว่าบรรดาฝูงสัตว์เดรัจฉาน จงฝ่าฟันพัฒนาวิชาการ อย่าล้างผลาญฤๅเข่นฆ่าบีฑาใคร ไม่ถือโทษโกรธแช่งซัดฮึดฮัดด่า หัดอภัยเหมือนกีฬาอัชฌาสัย ปฏิบัติประพฤติกฎกำหนดใจ พูดจาให้จ๊ะๆ จ๋าๆ น่าฟังเอยฯ</p>

<p>นายสังฆภัณฑ์ เฮงพิทักษ์ฝั่ง ผู้เฒ่าซึ่งมีอาชีพเป็นฅนขายฃวด ถูกตำรวจปฏิบัติการจับฟ้องศาล ฐานลักนาฬิกาคุณหญิงฉัตรชฎา ฌานสมาธิ</p>
</div>
</html>`,
		InlineFiles: []InlineFile{
			{
				ContentType: ContentTypeHeader{
					ContentType: "image/jpeg",
					Params: map[string]string{
						"name": "inline-jpg-image-without-disposition.jpg",
					},
				},
				ContentDisposition: ContentDispositionHeader{
					ContentDisposition: "",
					Params:             map[string]string(nil),
				},
				Data: []byte{255, 216, 255, 219, 0, 67, 0, 3, 2, 2, 2, 2, 2, 3, 2, 2, 2, 3, 3, 3, 3, 4,
					6, 4, 4, 4, 4, 4, 8, 6, 6, 5, 6, 9, 8, 10, 10, 9, 8, 9, 9, 10, 12, 15, 12, 10, 11, 14, 11, 9, 9,
					13, 17, 13, 14, 15, 16, 16, 17, 16, 10, 12, 18, 19, 18, 16, 19, 15, 16, 16, 16, 255, 201, 0, 11, 8,
					0, 1, 0, 1, 1, 1, 17, 0, 255, 204, 0, 6, 0, 16, 16, 5, 255, 218, 0, 8, 1, 1, 0, 0, 63, 0, 210, 207,
					32, 255, 217},
			},
			{
				ContentID: "inline-jpg-image.jpg@example.com",
				ContentType: ContentTypeHeader{
					ContentType: "image/jpeg",
					Params: map[string]string{
						"name": "inline-jpg-image-name.jpg",
					},
				},
				ContentDisposition: ContentDispositionHeader{
					ContentDisposition: inline,
					Params: map[string]string{
						"filename": "inline-jpg-image-filename.jpg",
					},
				},
				Data: []byte{255, 216, 255, 219, 0, 67, 0, 3, 2, 2, 2, 2, 2, 3, 2, 2, 2, 3, 3, 3, 3, 4,
					6, 4, 4, 4, 4, 4, 8, 6, 6, 5, 6, 9, 8, 10, 10, 9, 8, 9, 9, 10, 12, 15, 12, 10, 11, 14, 11, 9, 9,
					13, 17, 13, 14, 15, 16, 16, 17, 16, 10, 12, 18, 19, 18, 16, 19, 15, 16, 16, 16, 255, 201, 0, 11, 8,
					0, 1, 0, 1, 1, 1, 17, 0, 255, 204, 0, 6, 0, 16, 16, 5, 255, 218, 0, 8, 1, 1, 0, 0, 63, 0, 210, 207,
					32, 255, 217},
			},
		},
		AttachedFiles: []AttachedFile{
			{
				ContentType: ContentTypeHeader{
					ContentType: "application/pdf",
					Params: map[string]string{
						"name": "attached-pdf-name.pdf",
					},
				},
				ContentDisposition: ContentDispositionHeader{
					ContentDisposition: attachment,
					Params: map[string]string{
						"filename": "attached-pdf-filename.pdf",
					},
				},
				Data: []byte{37, 80, 68, 70, 45, 49, 46, 13, 116, 114, 97, 105, 108, 101, 114, 60, 60,
					47, 82, 111, 111, 116, 60, 60, 47, 80, 97, 103, 101, 115, 60, 60, 47, 75, 105, 100, 115, 91, 60,
					60, 47, 77, 101, 100, 105, 97, 66, 111, 120, 91, 48, 32, 48, 32, 51, 32, 51, 93, 62, 62, 93, 62,
					62, 62, 62, 62, 62},
			},
			{
				ContentType: ContentTypeHeader{
					ContentType: "application/pdf",
					Params: map[string]string{
						"name": "attached-pdf-without-disposition.pdf",
					},
				},
				ContentDisposition: ContentDispositionHeader{
					ContentDisposition: "",
					Params:             map[string]string(nil),
				},
				Data: []byte{37, 80, 68, 70, 45, 49, 46, 13, 116, 114, 97, 105, 108, 101, 114, 60, 60,
					47, 82, 111, 111, 116, 60, 60, 47, 80, 97, 103, 101, 115, 60, 60, 47, 75, 105, 100, 115, 91, 60,
					60, 47, 77, 101, 100, 105, 97, 66, 111, 120, 91, 48, 32, 48, 32, 51, 32, 51, 93, 62, 62, 93, 62,
					62, 62, 62, 62, 62},
			},
			{
				ContentType: ContentTypeHeader{
					ContentType: "application/json",
					Params: map[string]string{
						"name": "attached-json-name.json",
					},
				},
				ContentDisposition: ContentDispositionHeader{
					ContentDisposition: attachment,
					Params: map[string]string{
						"filename": "attached-json-filename.json",
					},
				},
				Data: []byte{123, 34, 102, 111, 111, 34, 58, 34, 98, 97, 114, 34, 125},
			},
			{
				ContentType: ContentTypeHeader{
					ContentType: "text/plain",
					Params: map[string]string{
						"name": "attached-text-plain-name.txt",
					},
				},
				ContentDisposition: ContentDispositionHeader{
					ContentDisposition: attachment,
					Params: map[string]string{
						"filename": "attached-text-plain-filename.txt",
					},
				},
				Data: []byte{84, 101, 120, 116, 47, 112, 108, 97, 105, 110, 32, 99, 111, 110, 116, 101, 110, 116, 32,
					97, 115, 32, 97, 110, 32, 97, 116, 116, 97, 99, 104, 101, 100, 32, 46, 116, 120, 116, 32, 102, 105,
					108, 101, 46},
			},
			{
				ContentType: ContentTypeHeader{
					ContentType: "text/html",
					Params: map[string]string{
						"name": "attached-text-html-name.html",
					},
				},
				ContentDisposition: ContentDispositionHeader{
					ContentDisposition: attachment,
					Params: map[string]string{
						"filename": "attached-text-html-filename.html",
					},
				},
				Data: []byte{84, 101, 120, 116, 47, 104, 116, 109, 108, 32, 99, 111, 110, 116, 101, 110, 116, 32, 97,
					115, 32, 97, 110, 32, 97, 116, 116, 97, 99, 104, 101, 100, 32, 46, 104, 116, 109, 108, 32, 102,
					105, 108, 101, 46},
			},
		},
	}

	testEmailFromFile(t, fp, expectedEmail)
}
func TestParseEmailThaiMultipartMixedWindows874OverBase64(t *testing.T) {
	fp := "tests/test_thai_multipart_mixed_windows-874_over_base64.txt"
	tz, _ := time.LoadLocation("Asia/Bangkok")
	expectedDate, _ := time.Parse(
		time.RFC1123Z+" (MST)",
		time.Date(2019, time.April, 1, 7, 55, 0, 0, tz).Format(time.RFC1123Z+" (MST)"))
	expectedEmail := Email{
		Headers: Headers{
			Date:    expectedDate,
			Subject: "📧 Test แพนแกรมภาษาไทย",
			ReplyTo: []*mail.Address{
				{
					Name:    "อลิซ ผู้ส่งจดหมาย",
					Address: "alis.phusngcdhmay@example.net",
				},
			},
			Sender: &mail.Address{
				Name:    "อลิซ ผู้ส่งจดหมาย",
				Address: "alis.phusngcdhmay@example.com",
			},
			From: []*mail.Address{
				{
					Name:    "อลิซ ผู้ส่งจดหมาย",
					Address: "alis.phusngcdhmay@example.com",
				},
				{
					Name:    "อลิซ ผู้ส่งจดหมาย",
					Address: "alis.phusngcdhmay@example.net",
				},
			},
			To: []*mail.Address{
				{
					Name:    "บ๊อบ ผู้รับ",
					Address: "bob.phurab@example.com",
				},
				{
					Name:    "คาโรล ผู้รับ",
					Address: "carol.phurab@example.com",
				},
			},
			Cc: []*mail.Address{
				{
					Name:    "แดน ผู้รับ",
					Address: "dan.phurab@example.com",
				},
				{
					Name:    "อีฟ ผู้รับ",
					Address: "eve.phurab@example.com",
				},
			},
			Bcc: []*mail.Address{
				{
					Name:    "แฟรงค์ ผู้รับ",
					Address: "frank.phurab@example.com",
				},
				{
					Name:    "เกรซ ผู้รับ",
					Address: "grace.phurab@example.com",
				},
			},
			MessageID:  "Message-Id-1@example.com",
			InReplyTo:  []MessageId{"Message-Id-0@example.com"},
			References: []MessageId{"Message-Id-0@example.com"},
			Comments:   "Message Header Comment",
			Keywords:   []string{"Keyword 1", "Keyword 2"},
			ResentDate: expectedDate,
			ResentFrom: []*mail.Address{
				{
					Name:    "อลิซ ผู้ส่งจดหมาย",
					Address: "alis.phusngcdhmay@example.net",
				},
				{
					Name:    "อลิซ ผู้ส่งจดหมาย",
					Address: "alis.phusngcdhmay@example.com",
				},
			},
			ResentSender: &mail.Address{
				Name:    "อลิซ ผู้ส่งจดหมาย",
				Address: "alis.phusngcdhmay@example.net",
			},
			ResentTo: []*mail.Address{
				{
					Name:    "บ๊อบ ผู้รับ",
					Address: "bob.phurab@example.net",
				},
				{
					Name:    "คาโรล ผู้รับ",
					Address: "carol.phurab@example.net",
				},
			},
			ResentCc: []*mail.Address{
				{
					Name:    "แดน ผู้รับ",
					Address: "dan.phurab@example.net",
				},
				{
					Name:    "อีฟ ผู้รับ",
					Address: "eve.phurab@example.net",
				},
			},
			ResentBcc: []*mail.Address{
				{
					Name:    "แฟรงค์ ผู้รับ",
					Address: "frank.phurab@example.net",
				},
				{
					Name:    "เกรซ ผู้รับ",
					Address: "grace.phurab@example.net",
				},
			},
			ResentMessageID: "Message-Id-1@example.net",
			ContentType: ContentTypeHeader{
				ContentType: "multipart/mixed",
				Params: map[string]string{
					"boundary": "MixedBoundaryString",
					"charset":  "windows-874",
				},
			},
			ExtraHeaders: map[string][]string{
				"X-Clacks-Overhead": {"GNU Terry Pratchett"},
			},
		},
		Text: `เป็นมนุษย์สุดประเสริฐเลิศคุณค่า กว่าบรรดาฝูงสัตว์เดรัจฉาน จงฝ่าฟันพัฒนาวิชาการ อย่าล้างผลาญฤๅเข่นฆ่าบีฑาใคร ไม่ถือโทษโกรธแช่งซัดฮึดฮัดด่า หัดอภัยเหมือนกีฬาอัชฌาสัย ปฏิบัติประพฤติกฎกำหนดใจ พูดจาให้จ๊ะๆ จ๋าๆ น่าฟังเอยฯ

นายสังฆภัณฑ์ เฮงพิทักษ์ฝั่ง ผู้เฒ่าซึ่งมีอาชีพเป็นฅนขายฃวด ถูกตำรวจปฏิบัติการจับฟ้องศาล ฐานลักนาฬิกาคุณหญิงฉัตรชฎา ฌานสมาธิ`,
		EnrichedText: `<bold>เป็นมนุษย์สุดประเสริฐเลิศคุณค่า</bold> <italic>กว่าบรรดาฝูงสัตว์เดรัจฉาน</italic> <fixed>จงฝ่าฟันพัฒนาวิชาการ</fixed> <underline>อย่าล้างผลาญฤๅเข่นฆ่าบีฑาใคร</underline> ไม่ถือโทษโกรธแช่งซัดฮึดฮัดด่า หัดอภัยเหมือนกีฬาอัชฌาสัย ปฏิบัติประพฤติกฎกำหนดใจ พูดจาให้จ๊ะๆ จ๋าๆ น่าฟังเอยฯ

นายสังฆภัณฑ์ เฮงพิทักษ์ฝั่ง ผู้เฒ่าซึ่งมีอาชีพเป็นฅนขายฃวด ถูกตำรวจปฏิบัติการจับฟ้องศาล ฐานลักนาฬิกาคุณหญิงฉัตรชฎา ฌานสมาธิ`,
		HTML: `<html>
<div dir="ltr">
<p>เป็นมนุษย์สุดประเสริฐเลิศคุณค่า กว่าบรรดาฝูงสัตว์เดรัจฉาน จงฝ่าฟันพัฒนาวิชาการ อย่าล้างผลาญฤๅเข่นฆ่าบีฑาใคร ไม่ถือโทษโกรธแช่งซัดฮึดฮัดด่า หัดอภัยเหมือนกีฬาอัชฌาสัย ปฏิบัติประพฤติกฎกำหนดใจ พูดจาให้จ๊ะๆ จ๋าๆ น่าฟังเอยฯ</p>

<p>นายสังฆภัณฑ์ เฮงพิทักษ์ฝั่ง ผู้เฒ่าซึ่งมีอาชีพเป็นฅนขายฃวด ถูกตำรวจปฏิบัติการจับฟ้องศาล ฐานลักนาฬิกาคุณหญิงฉัตรชฎา ฌานสมาธิ</p>
</div>
</html>`,
		InlineFiles: []InlineFile{
			{
				ContentType: ContentTypeHeader{
					ContentType: "image/jpeg",
					Params: map[string]string{
						"name": "inline-jpg-image-without-disposition.jpg",
					},
				},
				ContentDisposition: ContentDispositionHeader{
					ContentDisposition: "",
					Params:             map[string]string(nil),
				},
				Data: []byte{255, 216, 255, 219, 0, 67, 0, 3, 2, 2, 2, 2, 2, 3, 2, 2, 2, 3, 3, 3, 3, 4,
					6, 4, 4, 4, 4, 4, 8, 6, 6, 5, 6, 9, 8, 10, 10, 9, 8, 9, 9, 10, 12, 15, 12, 10, 11, 14, 11, 9, 9,
					13, 17, 13, 14, 15, 16, 16, 17, 16, 10, 12, 18, 19, 18, 16, 19, 15, 16, 16, 16, 255, 201, 0, 11, 8,
					0, 1, 0, 1, 1, 1, 17, 0, 255, 204, 0, 6, 0, 16, 16, 5, 255, 218, 0, 8, 1, 1, 0, 0, 63, 0, 210, 207,
					32, 255, 217},
			},
			{
				ContentID: "inline-jpg-image.jpg@example.com",
				ContentType: ContentTypeHeader{
					ContentType: "image/jpeg",
					Params: map[string]string{
						"name": "inline-jpg-image-name.jpg",
					},
				},
				ContentDisposition: ContentDispositionHeader{
					ContentDisposition: inline,
					Params: map[string]string{
						"filename": "inline-jpg-image-filename.jpg",
					},
				},
				Data: []byte{255, 216, 255, 219, 0, 67, 0, 3, 2, 2, 2, 2, 2, 3, 2, 2, 2, 3, 3, 3, 3, 4,
					6, 4, 4, 4, 4, 4, 8, 6, 6, 5, 6, 9, 8, 10, 10, 9, 8, 9, 9, 10, 12, 15, 12, 10, 11, 14, 11, 9, 9,
					13, 17, 13, 14, 15, 16, 16, 17, 16, 10, 12, 18, 19, 18, 16, 19, 15, 16, 16, 16, 255, 201, 0, 11, 8,
					0, 1, 0, 1, 1, 1, 17, 0, 255, 204, 0, 6, 0, 16, 16, 5, 255, 218, 0, 8, 1, 1, 0, 0, 63, 0, 210, 207,
					32, 255, 217},
			},
		},
		AttachedFiles: []AttachedFile{
			{
				ContentType: ContentTypeHeader{
					ContentType: "application/pdf",
					Params: map[string]string{
						"name": "attached-pdf-name.pdf",
					},
				},
				ContentDisposition: ContentDispositionHeader{
					ContentDisposition: attachment,
					Params: map[string]string{
						"filename": "attached-pdf-filename.pdf",
					},
				},
				Data: []byte{37, 80, 68, 70, 45, 49, 46, 13, 116, 114, 97, 105, 108, 101, 114, 60, 60,
					47, 82, 111, 111, 116, 60, 60, 47, 80, 97, 103, 101, 115, 60, 60, 47, 75, 105, 100, 115, 91, 60,
					60, 47, 77, 101, 100, 105, 97, 66, 111, 120, 91, 48, 32, 48, 32, 51, 32, 51, 93, 62, 62, 93, 62,
					62, 62, 62, 62, 62},
			},
			{
				ContentType: ContentTypeHeader{
					ContentType: "application/pdf",
					Params: map[string]string{
						"name": "attached-pdf-without-disposition.pdf",
					},
				},
				ContentDisposition: ContentDispositionHeader{
					ContentDisposition: "",
					Params:             map[string]string(nil),
				},
				Data: []byte{37, 80, 68, 70, 45, 49, 46, 13, 116, 114, 97, 105, 108, 101, 114, 60, 60,
					47, 82, 111, 111, 116, 60, 60, 47, 80, 97, 103, 101, 115, 60, 60, 47, 75, 105, 100, 115, 91, 60,
					60, 47, 77, 101, 100, 105, 97, 66, 111, 120, 91, 48, 32, 48, 32, 51, 32, 51, 93, 62, 62, 93, 62,
					62, 62, 62, 62, 62},
			},
			{
				ContentType: ContentTypeHeader{
					ContentType: "application/json",
					Params: map[string]string{
						"name": "attached-json-name.json",
					},
				},
				ContentDisposition: ContentDispositionHeader{
					ContentDisposition: attachment,
					Params: map[string]string{
						"filename": "attached-json-filename.json",
					},
				},
				Data: []byte{123, 34, 102, 111, 111, 34, 58, 34, 98, 97, 114, 34, 125},
			},
			{
				ContentType: ContentTypeHeader{
					ContentType: "text/plain",
					Params: map[string]string{
						"name": "attached-text-plain-name.txt",
					},
				},
				ContentDisposition: ContentDispositionHeader{
					ContentDisposition: attachment,
					Params: map[string]string{
						"filename": "attached-text-plain-filename.txt",
					},
				},
				Data: []byte{84, 101, 120, 116, 47, 112, 108, 97, 105, 110, 32, 99, 111, 110, 116, 101, 110, 116, 32,
					97, 115, 32, 97, 110, 32, 97, 116, 116, 97, 99, 104, 101, 100, 32, 46, 116, 120, 116, 32, 102, 105,
					108, 101, 46},
			},
			{
				ContentType: ContentTypeHeader{
					ContentType: "text/html",
					Params: map[string]string{
						"name": "attached-text-html-name.html",
					},
				},
				ContentDisposition: ContentDispositionHeader{
					ContentDisposition: attachment,
					Params: map[string]string{
						"filename": "attached-text-html-filename.html",
					},
				},
				Data: []byte{84, 101, 120, 116, 47, 104, 116, 109, 108, 32, 99, 111, 110, 116, 101, 110, 116, 32, 97,
					115, 32, 97, 110, 32, 97, 116, 116, 97, 99, 104, 101, 100, 32, 46, 104, 116, 109, 108, 32, 102,
					105, 108, 101, 46},
			},
		},
	}

	testEmailFromFile(t, fp, expectedEmail)
}
func TestParseEmailThaiMultipartMixedWindows874OverQuotedprintable(t *testing.T) {
	fp := "tests/test_thai_multipart_mixed_windows-874_over_quoted-printable.txt"
	tz, _ := time.LoadLocation("Asia/Bangkok")
	expectedDate, _ := time.Parse(
		time.RFC1123Z+" (MST)",
		time.Date(2019, time.April, 1, 7, 55, 0, 0, tz).Format(time.RFC1123Z+" (MST)"))
	expectedEmail := Email{
		Headers: Headers{
			Date:    expectedDate,
			Subject: "📧 Test แพนแกรมภาษาไทย",
			ReplyTo: []*mail.Address{
				{
					Name:    "อลิซ ผู้ส่งจดหมาย",
					Address: "alis.phusngcdhmay@example.net",
				},
			},
			Sender: &mail.Address{
				Name:    "อลิซ ผู้ส่งจดหมาย",
				Address: "alis.phusngcdhmay@example.com",
			},
			From: []*mail.Address{
				{
					Name:    "อลิซ ผู้ส่งจดหมาย",
					Address: "alis.phusngcdhmay@example.com",
				},
				{
					Name:    "อลิซ ผู้ส่งจดหมาย",
					Address: "alis.phusngcdhmay@example.net",
				},
			},
			To: []*mail.Address{
				{
					Name:    "บ๊อบ ผู้รับ",
					Address: "bob.phurab@example.com",
				},
				{
					Name:    "คาโรล ผู้รับ",
					Address: "carol.phurab@example.com",
				},
			},
			Cc: []*mail.Address{
				{
					Name:    "แดน ผู้รับ",
					Address: "dan.phurab@example.com",
				},
				{
					Name:    "อีฟ ผู้รับ",
					Address: "eve.phurab@example.com",
				},
			},
			Bcc: []*mail.Address{
				{
					Name:    "แฟรงค์ ผู้รับ",
					Address: "frank.phurab@example.com",
				},
				{
					Name:    "เกรซ ผู้รับ",
					Address: "grace.phurab@example.com",
				},
			},
			MessageID:  "Message-Id-1@example.com",
			InReplyTo:  []MessageId{"Message-Id-0@example.com"},
			References: []MessageId{"Message-Id-0@example.com"},
			Comments:   "Message Header Comment",
			Keywords:   []string{"Keyword 1", "Keyword 2"},
			ResentDate: expectedDate,
			ResentFrom: []*mail.Address{
				{
					Name:    "อลิซ ผู้ส่งจดหมาย",
					Address: "alis.phusngcdhmay@example.net",
				},
				{
					Name:    "อลิซ ผู้ส่งจดหมาย",
					Address: "alis.phusngcdhmay@example.com",
				},
			},
			ResentSender: &mail.Address{
				Name:    "อลิซ ผู้ส่งจดหมาย",
				Address: "alis.phusngcdhmay@example.net",
			},
			ResentTo: []*mail.Address{
				{
					Name:    "บ๊อบ ผู้รับ",
					Address: "bob.phurab@example.net",
				},
				{
					Name:    "คาโรล ผู้รับ",
					Address: "carol.phurab@example.net",
				},
			},
			ResentCc: []*mail.Address{
				{
					Name:    "แดน ผู้รับ",
					Address: "dan.phurab@example.net",
				},
				{
					Name:    "อีฟ ผู้รับ",
					Address: "eve.phurab@example.net",
				},
			},
			ResentBcc: []*mail.Address{
				{
					Name:    "แฟรงค์ ผู้รับ",
					Address: "frank.phurab@example.net",
				},
				{
					Name:    "เกรซ ผู้รับ",
					Address: "grace.phurab@example.net",
				},
			},
			ResentMessageID: "Message-Id-1@example.net",
			ContentType: ContentTypeHeader{
				ContentType: "multipart/mixed",
				Params: map[string]string{
					"boundary": "MixedBoundaryString",
					"charset":  "windows-874",
				},
			},
			ExtraHeaders: map[string][]string{
				"X-Clacks-Overhead": {"GNU Terry Pratchett"},
			},
		},
		Text: `เป็นมนุษย์สุดประเสริฐเลิศคุณค่า กว่าบรรดาฝูงสัตว์เดรัจฉาน จงฝ่าฟันพัฒนาวิชาการ อย่าล้างผลาญฤๅเข่นฆ่าบีฑาใคร ไม่ถือโทษโกรธแช่งซัดฮึดฮัดด่า หัดอภัยเหมือนกีฬาอัชฌาสัย ปฏิบัติประพฤติกฎกำหนดใจ พูดจาให้จ๊ะๆ จ๋าๆ น่าฟังเอยฯ

นายสังฆภัณฑ์ เฮงพิทักษ์ฝั่ง ผู้เฒ่าซึ่งมีอาชีพเป็นฅนขายฃวด ถูกตำรวจปฏิบัติการจับฟ้องศาล ฐานลักนาฬิกาคุณหญิงฉัตรชฎา ฌานสมาธิ`,
		EnrichedText: `<bold>เป็นมนุษย์สุดประเสริฐเลิศคุณค่า</bold> <italic>กว่าบรรดาฝูงสัตว์เดรัจฉาน</italic> <fixed>จงฝ่าฟันพัฒนาวิชาการ</fixed> <underline>อย่าล้างผลาญฤๅเข่นฆ่าบีฑาใคร</underline> ไม่ถือโทษโกรธแช่งซัดฮึดฮัดด่า หัดอภัยเหมือนกีฬาอัชฌาสัย ปฏิบัติประพฤติกฎกำหนดใจ พูดจาให้จ๊ะๆ จ๋าๆ น่าฟังเอยฯ

นายสังฆภัณฑ์ เฮงพิทักษ์ฝั่ง ผู้เฒ่าซึ่งมีอาชีพเป็นฅนขายฃวด ถูกตำรวจปฏิบัติการจับฟ้องศาล ฐานลักนาฬิกาคุณหญิงฉัตรชฎา ฌานสมาธิ`,
		HTML: `<html>
<div dir="ltr">
<p>เป็นมนุษย์สุดประเสริฐเลิศคุณค่า กว่าบรรดาฝูงสัตว์เดรัจฉาน จงฝ่าฟันพัฒนาวิชาการ อย่าล้างผลาญฤๅเข่นฆ่าบีฑาใคร ไม่ถือโทษโกรธแช่งซัดฮึดฮัดด่า หัดอภัยเหมือนกีฬาอัชฌาสัย ปฏิบัติประพฤติกฎกำหนดใจ พูดจาให้จ๊ะๆ จ๋าๆ น่าฟังเอยฯ</p>

<p>นายสังฆภัณฑ์ เฮงพิทักษ์ฝั่ง ผู้เฒ่าซึ่งมีอาชีพเป็นฅนขายฃวด ถูกตำรวจปฏิบัติการจับฟ้องศาล ฐานลักนาฬิกาคุณหญิงฉัตรชฎา ฌานสมาธิ</p>
</div>
</html>`,
		InlineFiles: []InlineFile{
			{
				ContentType: ContentTypeHeader{
					ContentType: "image/jpeg",
					Params: map[string]string{
						"name": "inline-jpg-image-without-disposition.jpg",
					},
				},
				ContentDisposition: ContentDispositionHeader{
					ContentDisposition: "",
					Params:             map[string]string(nil),
				},
				Data: []byte{255, 216, 255, 219, 0, 67, 0, 3, 2, 2, 2, 2, 2, 3, 2, 2, 2, 3, 3, 3, 3, 4,
					6, 4, 4, 4, 4, 4, 8, 6, 6, 5, 6, 9, 8, 10, 10, 9, 8, 9, 9, 10, 12, 15, 12, 10, 11, 14, 11, 9, 9,
					13, 17, 13, 14, 15, 16, 16, 17, 16, 10, 12, 18, 19, 18, 16, 19, 15, 16, 16, 16, 255, 201, 0, 11, 8,
					0, 1, 0, 1, 1, 1, 17, 0, 255, 204, 0, 6, 0, 16, 16, 5, 255, 218, 0, 8, 1, 1, 0, 0, 63, 0, 210, 207,
					32, 255, 217},
			},
			{
				ContentID: "inline-jpg-image.jpg@example.com",
				ContentType: ContentTypeHeader{
					ContentType: "image/jpeg",
					Params: map[string]string{
						"name": "inline-jpg-image-name.jpg",
					},
				},
				ContentDisposition: ContentDispositionHeader{
					ContentDisposition: inline,
					Params: map[string]string{
						"filename": "inline-jpg-image-filename.jpg",
					},
				},
				Data: []byte{255, 216, 255, 219, 0, 67, 0, 3, 2, 2, 2, 2, 2, 3, 2, 2, 2, 3, 3, 3, 3, 4,
					6, 4, 4, 4, 4, 4, 8, 6, 6, 5, 6, 9, 8, 10, 10, 9, 8, 9, 9, 10, 12, 15, 12, 10, 11, 14, 11, 9, 9,
					13, 17, 13, 14, 15, 16, 16, 17, 16, 10, 12, 18, 19, 18, 16, 19, 15, 16, 16, 16, 255, 201, 0, 11, 8,
					0, 1, 0, 1, 1, 1, 17, 0, 255, 204, 0, 6, 0, 16, 16, 5, 255, 218, 0, 8, 1, 1, 0, 0, 63, 0, 210, 207,
					32, 255, 217},
			},
		},
		AttachedFiles: []AttachedFile{
			{
				ContentType: ContentTypeHeader{
					ContentType: "application/pdf",
					Params: map[string]string{
						"name": "attached-pdf-name.pdf",
					},
				},
				ContentDisposition: ContentDispositionHeader{
					ContentDisposition: attachment,
					Params: map[string]string{
						"filename": "attached-pdf-filename.pdf",
					},
				},
				Data: []byte{37, 80, 68, 70, 45, 49, 46, 13, 116, 114, 97, 105, 108, 101, 114, 60, 60,
					47, 82, 111, 111, 116, 60, 60, 47, 80, 97, 103, 101, 115, 60, 60, 47, 75, 105, 100, 115, 91, 60,
					60, 47, 77, 101, 100, 105, 97, 66, 111, 120, 91, 48, 32, 48, 32, 51, 32, 51, 93, 62, 62, 93, 62,
					62, 62, 62, 62, 62},
			},
			{
				ContentType: ContentTypeHeader{
					ContentType: "application/pdf",
					Params: map[string]string{
						"name": "attached-pdf-without-disposition.pdf",
					},
				},
				ContentDisposition: ContentDispositionHeader{
					ContentDisposition: "",
					Params:             map[string]string(nil),
				},
				Data: []byte{37, 80, 68, 70, 45, 49, 46, 13, 116, 114, 97, 105, 108, 101, 114, 60, 60,
					47, 82, 111, 111, 116, 60, 60, 47, 80, 97, 103, 101, 115, 60, 60, 47, 75, 105, 100, 115, 91, 60,
					60, 47, 77, 101, 100, 105, 97, 66, 111, 120, 91, 48, 32, 48, 32, 51, 32, 51, 93, 62, 62, 93, 62,
					62, 62, 62, 62, 62},
			},
			{
				ContentType: ContentTypeHeader{
					ContentType: "application/json",
					Params: map[string]string{
						"name": "attached-json-name.json",
					},
				},
				ContentDisposition: ContentDispositionHeader{
					ContentDisposition: attachment,
					Params: map[string]string{
						"filename": "attached-json-filename.json",
					},
				},
				Data: []byte{123, 34, 102, 111, 111, 34, 58, 34, 98, 97, 114, 34, 125},
			},
			{
				ContentType: ContentTypeHeader{
					ContentType: "text/plain",
					Params: map[string]string{
						"name": "attached-text-plain-name.txt",
					},
				},
				ContentDisposition: ContentDispositionHeader{
					ContentDisposition: attachment,
					Params: map[string]string{
						"filename": "attached-text-plain-filename.txt",
					},
				},
				Data: []byte{84, 101, 120, 116, 47, 112, 108, 97, 105, 110, 32, 99, 111, 110, 116, 101, 110, 116, 32,
					97, 115, 32, 97, 110, 32, 97, 116, 116, 97, 99, 104, 101, 100, 32, 46, 116, 120, 116, 32, 102, 105,
					108, 101, 46},
			},
			{
				ContentType: ContentTypeHeader{
					ContentType: "text/html",
					Params: map[string]string{
						"name": "attached-text-html-name.html",
					},
				},
				ContentDisposition: ContentDispositionHeader{
					ContentDisposition: attachment,
					Params: map[string]string{
						"filename": "attached-text-html-filename.html",
					},
				},
				Data: []byte{84, 101, 120, 116, 47, 104, 116, 109, 108, 32, 99, 111, 110, 116, 101, 110, 116, 32, 97,
					115, 32, 97, 110, 32, 97, 116, 116, 97, 99, 104, 101, 100, 32, 46, 104, 116, 109, 108, 32, 102,
					105, 108, 101, 46},
			},
		},
	}

	testEmailFromFile(t, fp, expectedEmail)
}
func TestParseEmailThaiMultipartMixedTis620OverBase64(t *testing.T) {
	fp := "tests/test_thai_multipart_mixed_tis-620_over_base64.txt"
	tz, _ := time.LoadLocation("Asia/Bangkok")
	expectedDate, _ := time.Parse(
		time.RFC1123Z+" (MST)",
		time.Date(2019, time.April, 1, 7, 55, 0, 0, tz).Format(time.RFC1123Z+" (MST)"))
	expectedEmail := Email{
		Headers: Headers{
			Date:    expectedDate,
			Subject: "📧 Test แพนแกรมภาษาไทย",
			ReplyTo: []*mail.Address{
				{
					Name:    "อลิซ ผู้ส่งจดหมาย",
					Address: "alis.phusngcdhmay@example.net",
				},
			},
			Sender: &mail.Address{
				Name:    "อลิซ ผู้ส่งจดหมาย",
				Address: "alis.phusngcdhmay@example.com",
			},
			From: []*mail.Address{
				{
					Name:    "อลิซ ผู้ส่งจดหมาย",
					Address: "alis.phusngcdhmay@example.com",
				},
				{
					Name:    "อลิซ ผู้ส่งจดหมาย",
					Address: "alis.phusngcdhmay@example.net",
				},
			},
			To: []*mail.Address{
				{
					Name:    "บ๊อบ ผู้รับ",
					Address: "bob.phurab@example.com",
				},
				{
					Name:    "คาโรล ผู้รับ",
					Address: "carol.phurab@example.com",
				},
			},
			Cc: []*mail.Address{
				{
					Name:    "แดน ผู้รับ",
					Address: "dan.phurab@example.com",
				},
				{
					Name:    "อีฟ ผู้รับ",
					Address: "eve.phurab@example.com",
				},
			},
			Bcc: []*mail.Address{
				{
					Name:    "แฟรงค์ ผู้รับ",
					Address: "frank.phurab@example.com",
				},
				{
					Name:    "เกรซ ผู้รับ",
					Address: "grace.phurab@example.com",
				},
			},
			MessageID:  "Message-Id-1@example.com",
			InReplyTo:  []MessageId{"Message-Id-0@example.com"},
			References: []MessageId{"Message-Id-0@example.com"},
			Comments:   "Message Header Comment",
			Keywords:   []string{"Keyword 1", "Keyword 2"},
			ResentDate: expectedDate,
			ResentFrom: []*mail.Address{
				{
					Name:    "อลิซ ผู้ส่งจดหมาย",
					Address: "alis.phusngcdhmay@example.net",
				},
				{
					Name:    "อลิซ ผู้ส่งจดหมาย",
					Address: "alis.phusngcdhmay@example.com",
				},
			},
			ResentSender: &mail.Address{
				Name:    "อลิซ ผู้ส่งจดหมาย",
				Address: "alis.phusngcdhmay@example.net",
			},
			ResentTo: []*mail.Address{
				{
					Name:    "บ๊อบ ผู้รับ",
					Address: "bob.phurab@example.net",
				},
				{
					Name:    "คาโรล ผู้รับ",
					Address: "carol.phurab@example.net",
				},
			},
			ResentCc: []*mail.Address{
				{
					Name:    "แดน ผู้รับ",
					Address: "dan.phurab@example.net",
				},
				{
					Name:    "อีฟ ผู้รับ",
					Address: "eve.phurab@example.net",
				},
			},
			ResentBcc: []*mail.Address{
				{
					Name:    "แฟรงค์ ผู้รับ",
					Address: "frank.phurab@example.net",
				},
				{
					Name:    "เกรซ ผู้รับ",
					Address: "grace.phurab@example.net",
				},
			},
			ResentMessageID: "Message-Id-1@example.net",
			ContentType: ContentTypeHeader{
				ContentType: "multipart/mixed",
				Params: map[string]string{
					"boundary": "MixedBoundaryString",
					"charset":  "tis-620",
				},
			},
			ExtraHeaders: map[string][]string{
				"X-Clacks-Overhead": {"GNU Terry Pratchett"},
			},
		},
		Text: `เป็นมนุษย์สุดประเสริฐเลิศคุณค่า กว่าบรรดาฝูงสัตว์เดรัจฉาน จงฝ่าฟันพัฒนาวิชาการ อย่าล้างผลาญฤๅเข่นฆ่าบีฑาใคร ไม่ถือโทษโกรธแช่งซัดฮึดฮัดด่า หัดอภัยเหมือนกีฬาอัชฌาสัย ปฏิบัติประพฤติกฎกำหนดใจ พูดจาให้จ๊ะๆ จ๋าๆ น่าฟังเอยฯ

นายสังฆภัณฑ์ เฮงพิทักษ์ฝั่ง ผู้เฒ่าซึ่งมีอาชีพเป็นฅนขายฃวด ถูกตำรวจปฏิบัติการจับฟ้องศาล ฐานลักนาฬิกาคุณหญิงฉัตรชฎา ฌานสมาธิ`,
		EnrichedText: `<bold>เป็นมนุษย์สุดประเสริฐเลิศคุณค่า</bold> <italic>กว่าบรรดาฝูงสัตว์เดรัจฉาน</italic> <fixed>จงฝ่าฟันพัฒนาวิชาการ</fixed> <underline>อย่าล้างผลาญฤๅเข่นฆ่าบีฑาใคร</underline> ไม่ถือโทษโกรธแช่งซัดฮึดฮัดด่า หัดอภัยเหมือนกีฬาอัชฌาสัย ปฏิบัติประพฤติกฎกำหนดใจ พูดจาให้จ๊ะๆ จ๋าๆ น่าฟังเอยฯ

นายสังฆภัณฑ์ เฮงพิทักษ์ฝั่ง ผู้เฒ่าซึ่งมีอาชีพเป็นฅนขายฃวด ถูกตำรวจปฏิบัติการจับฟ้องศาล ฐานลักนาฬิกาคุณหญิงฉัตรชฎา ฌานสมาธิ`,
		HTML: `<html>
<div dir="ltr">
<p>เป็นมนุษย์สุดประเสริฐเลิศคุณค่า กว่าบรรดาฝูงสัตว์เดรัจฉาน จงฝ่าฟันพัฒนาวิชาการ อย่าล้างผลาญฤๅเข่นฆ่าบีฑาใคร ไม่ถือโทษโกรธแช่งซัดฮึดฮัดด่า หัดอภัยเหมือนกีฬาอัชฌาสัย ปฏิบัติประพฤติกฎกำหนดใจ พูดจาให้จ๊ะๆ จ๋าๆ น่าฟังเอยฯ</p>

<p>นายสังฆภัณฑ์ เฮงพิทักษ์ฝั่ง ผู้เฒ่าซึ่งมีอาชีพเป็นฅนขายฃวด ถูกตำรวจปฏิบัติการจับฟ้องศาล ฐานลักนาฬิกาคุณหญิงฉัตรชฎา ฌานสมาธิ</p>
</div>
</html>`,
		InlineFiles: []InlineFile{
			{
				ContentType: ContentTypeHeader{
					ContentType: "image/jpeg",
					Params: map[string]string{
						"name": "inline-jpg-image-without-disposition.jpg",
					},
				},
				ContentDisposition: ContentDispositionHeader{
					ContentDisposition: "",
					Params:             map[string]string(nil),
				},
				Data: []byte{255, 216, 255, 219, 0, 67, 0, 3, 2, 2, 2, 2, 2, 3, 2, 2, 2, 3, 3, 3, 3, 4,
					6, 4, 4, 4, 4, 4, 8, 6, 6, 5, 6, 9, 8, 10, 10, 9, 8, 9, 9, 10, 12, 15, 12, 10, 11, 14, 11, 9, 9,
					13, 17, 13, 14, 15, 16, 16, 17, 16, 10, 12, 18, 19, 18, 16, 19, 15, 16, 16, 16, 255, 201, 0, 11, 8,
					0, 1, 0, 1, 1, 1, 17, 0, 255, 204, 0, 6, 0, 16, 16, 5, 255, 218, 0, 8, 1, 1, 0, 0, 63, 0, 210, 207,
					32, 255, 217},
			},
			{
				ContentID: "inline-jpg-image.jpg@example.com",
				ContentType: ContentTypeHeader{
					ContentType: "image/jpeg",
					Params: map[string]string{
						"name": "inline-jpg-image-name.jpg",
					},
				},
				ContentDisposition: ContentDispositionHeader{
					ContentDisposition: inline,
					Params: map[string]string{
						"filename": "inline-jpg-image-filename.jpg",
					},
				},
				Data: []byte{255, 216, 255, 219, 0, 67, 0, 3, 2, 2, 2, 2, 2, 3, 2, 2, 2, 3, 3, 3, 3, 4,
					6, 4, 4, 4, 4, 4, 8, 6, 6, 5, 6, 9, 8, 10, 10, 9, 8, 9, 9, 10, 12, 15, 12, 10, 11, 14, 11, 9, 9,
					13, 17, 13, 14, 15, 16, 16, 17, 16, 10, 12, 18, 19, 18, 16, 19, 15, 16, 16, 16, 255, 201, 0, 11, 8,
					0, 1, 0, 1, 1, 1, 17, 0, 255, 204, 0, 6, 0, 16, 16, 5, 255, 218, 0, 8, 1, 1, 0, 0, 63, 0, 210, 207,
					32, 255, 217},
			},
		},
		AttachedFiles: []AttachedFile{
			{
				ContentType: ContentTypeHeader{
					ContentType: "application/pdf",
					Params: map[string]string{
						"name": "attached-pdf-name.pdf",
					},
				},
				ContentDisposition: ContentDispositionHeader{
					ContentDisposition: attachment,
					Params: map[string]string{
						"filename": "attached-pdf-filename.pdf",
					},
				},
				Data: []byte{37, 80, 68, 70, 45, 49, 46, 13, 116, 114, 97, 105, 108, 101, 114, 60, 60,
					47, 82, 111, 111, 116, 60, 60, 47, 80, 97, 103, 101, 115, 60, 60, 47, 75, 105, 100, 115, 91, 60,
					60, 47, 77, 101, 100, 105, 97, 66, 111, 120, 91, 48, 32, 48, 32, 51, 32, 51, 93, 62, 62, 93, 62,
					62, 62, 62, 62, 62},
			},
			{
				ContentType: ContentTypeHeader{
					ContentType: "application/pdf",
					Params: map[string]string{
						"name": "attached-pdf-without-disposition.pdf",
					},
				},
				ContentDisposition: ContentDispositionHeader{
					ContentDisposition: "",
					Params:             map[string]string(nil),
				},
				Data: []byte{37, 80, 68, 70, 45, 49, 46, 13, 116, 114, 97, 105, 108, 101, 114, 60, 60,
					47, 82, 111, 111, 116, 60, 60, 47, 80, 97, 103, 101, 115, 60, 60, 47, 75, 105, 100, 115, 91, 60,
					60, 47, 77, 101, 100, 105, 97, 66, 111, 120, 91, 48, 32, 48, 32, 51, 32, 51, 93, 62, 62, 93, 62,
					62, 62, 62, 62, 62},
			},
			{
				ContentType: ContentTypeHeader{
					ContentType: "application/json",
					Params: map[string]string{
						"name": "attached-json-name.json",
					},
				},
				ContentDisposition: ContentDispositionHeader{
					ContentDisposition: attachment,
					Params: map[string]string{
						"filename": "attached-json-filename.json",
					},
				},
				Data: []byte{123, 34, 102, 111, 111, 34, 58, 34, 98, 97, 114, 34, 125},
			},
			{
				ContentType: ContentTypeHeader{
					ContentType: "text/plain",
					Params: map[string]string{
						"name": "attached-text-plain-name.txt",
					},
				},
				ContentDisposition: ContentDispositionHeader{
					ContentDisposition: attachment,
					Params: map[string]string{
						"filename": "attached-text-plain-filename.txt",
					},
				},
				Data: []byte{84, 101, 120, 116, 47, 112, 108, 97, 105, 110, 32, 99, 111, 110, 116, 101, 110, 116, 32,
					97, 115, 32, 97, 110, 32, 97, 116, 116, 97, 99, 104, 101, 100, 32, 46, 116, 120, 116, 32, 102, 105,
					108, 101, 46},
			},
			{
				ContentType: ContentTypeHeader{
					ContentType: "text/html",
					Params: map[string]string{
						"name": "attached-text-html-name.html",
					},
				},
				ContentDisposition: ContentDispositionHeader{
					ContentDisposition: attachment,
					Params: map[string]string{
						"filename": "attached-text-html-filename.html",
					},
				},
				Data: []byte{84, 101, 120, 116, 47, 104, 116, 109, 108, 32, 99, 111, 110, 116, 101, 110, 116, 32, 97,
					115, 32, 97, 110, 32, 97, 116, 116, 97, 99, 104, 101, 100, 32, 46, 104, 116, 109, 108, 32, 102,
					105, 108, 101, 46},
			},
		},
	}

	testEmailFromFile(t, fp, expectedEmail)
}
func TestParseEmailThaiMultipartMixedTis620OverQuotedprintable(t *testing.T) {
	fp := "tests/test_thai_multipart_mixed_tis-620_over_quoted-printable.txt"
	tz, _ := time.LoadLocation("Asia/Bangkok")
	expectedDate, _ := time.Parse(
		time.RFC1123Z+" (MST)",
		time.Date(2019, time.April, 1, 7, 55, 0, 0, tz).Format(time.RFC1123Z+" (MST)"))
	expectedEmail := Email{
		Headers: Headers{
			Date:    expectedDate,
			Subject: "📧 Test แพนแกรมภาษาไทย",
			ReplyTo: []*mail.Address{
				{
					Name:    "อลิซ ผู้ส่งจดหมาย",
					Address: "alis.phusngcdhmay@example.net",
				},
			},
			Sender: &mail.Address{
				Name:    "อลิซ ผู้ส่งจดหมาย",
				Address: "alis.phusngcdhmay@example.com",
			},
			From: []*mail.Address{
				{
					Name:    "อลิซ ผู้ส่งจดหมาย",
					Address: "alis.phusngcdhmay@example.com",
				},
				{
					Name:    "อลิซ ผู้ส่งจดหมาย",
					Address: "alis.phusngcdhmay@example.net",
				},
			},
			To: []*mail.Address{
				{
					Name:    "บ๊อบ ผู้รับ",
					Address: "bob.phurab@example.com",
				},
				{
					Name:    "คาโรล ผู้รับ",
					Address: "carol.phurab@example.com",
				},
			},
			Cc: []*mail.Address{
				{
					Name:    "แดน ผู้รับ",
					Address: "dan.phurab@example.com",
				},
				{
					Name:    "อีฟ ผู้รับ",
					Address: "eve.phurab@example.com",
				},
			},
			Bcc: []*mail.Address{
				{
					Name:    "แฟรงค์ ผู้รับ",
					Address: "frank.phurab@example.com",
				},
				{
					Name:    "เกรซ ผู้รับ",
					Address: "grace.phurab@example.com",
				},
			},
			MessageID:  "Message-Id-1@example.com",
			InReplyTo:  []MessageId{"Message-Id-0@example.com"},
			References: []MessageId{"Message-Id-0@example.com"},
			Comments:   "Message Header Comment",
			Keywords:   []string{"Keyword 1", "Keyword 2"},
			ResentDate: expectedDate,
			ResentFrom: []*mail.Address{
				{
					Name:    "อลิซ ผู้ส่งจดหมาย",
					Address: "alis.phusngcdhmay@example.net",
				},
				{
					Name:    "อลิซ ผู้ส่งจดหมาย",
					Address: "alis.phusngcdhmay@example.com",
				},
			},
			ResentSender: &mail.Address{
				Name:    "อลิซ ผู้ส่งจดหมาย",
				Address: "alis.phusngcdhmay@example.net",
			},
			ResentTo: []*mail.Address{
				{
					Name:    "บ๊อบ ผู้รับ",
					Address: "bob.phurab@example.net",
				},
				{
					Name:    "คาโรล ผู้รับ",
					Address: "carol.phurab@example.net",
				},
			},
			ResentCc: []*mail.Address{
				{
					Name:    "แดน ผู้รับ",
					Address: "dan.phurab@example.net",
				},
				{
					Name:    "อีฟ ผู้รับ",
					Address: "eve.phurab@example.net",
				},
			},
			ResentBcc: []*mail.Address{
				{
					Name:    "แฟรงค์ ผู้รับ",
					Address: "frank.phurab@example.net",
				},
				{
					Name:    "เกรซ ผู้รับ",
					Address: "grace.phurab@example.net",
				},
			},
			ResentMessageID: "Message-Id-1@example.net",
			ContentType: ContentTypeHeader{
				ContentType: "multipart/mixed",
				Params: map[string]string{
					"boundary": "MixedBoundaryString",
					"charset":  "tis-620",
				},
			},
			ExtraHeaders: map[string][]string{
				"X-Clacks-Overhead": {"GNU Terry Pratchett"},
			},
		},
		Text: `เป็นมนุษย์สุดประเสริฐเลิศคุณค่า กว่าบรรดาฝูงสัตว์เดรัจฉาน จงฝ่าฟันพัฒนาวิชาการ อย่าล้างผลาญฤๅเข่นฆ่าบีฑาใคร ไม่ถือโทษโกรธแช่งซัดฮึดฮัดด่า หัดอภัยเหมือนกีฬาอัชฌาสัย ปฏิบัติประพฤติกฎกำหนดใจ พูดจาให้จ๊ะๆ จ๋าๆ น่าฟังเอยฯ

นายสังฆภัณฑ์ เฮงพิทักษ์ฝั่ง ผู้เฒ่าซึ่งมีอาชีพเป็นฅนขายฃวด ถูกตำรวจปฏิบัติการจับฟ้องศาล ฐานลักนาฬิกาคุณหญิงฉัตรชฎา ฌานสมาธิ`,
		EnrichedText: `<bold>เป็นมนุษย์สุดประเสริฐเลิศคุณค่า</bold> <italic>กว่าบรรดาฝูงสัตว์เดรัจฉาน</italic> <fixed>จงฝ่าฟันพัฒนาวิชาการ</fixed> <underline>อย่าล้างผลาญฤๅเข่นฆ่าบีฑาใคร</underline> ไม่ถือโทษโกรธแช่งซัดฮึดฮัดด่า หัดอภัยเหมือนกีฬาอัชฌาสัย ปฏิบัติประพฤติกฎกำหนดใจ พูดจาให้จ๊ะๆ จ๋าๆ น่าฟังเอยฯ

นายสังฆภัณฑ์ เฮงพิทักษ์ฝั่ง ผู้เฒ่าซึ่งมีอาชีพเป็นฅนขายฃวด ถูกตำรวจปฏิบัติการจับฟ้องศาล ฐานลักนาฬิกาคุณหญิงฉัตรชฎา ฌานสมาธิ`,
		HTML: `<html>
<div dir="ltr">
<p>เป็นมนุษย์สุดประเสริฐเลิศคุณค่า กว่าบรรดาฝูงสัตว์เดรัจฉาน จงฝ่าฟันพัฒนาวิชาการ อย่าล้างผลาญฤๅเข่นฆ่าบีฑาใคร ไม่ถือโทษโกรธแช่งซัดฮึดฮัดด่า หัดอภัยเหมือนกีฬาอัชฌาสัย ปฏิบัติประพฤติกฎกำหนดใจ พูดจาให้จ๊ะๆ จ๋าๆ น่าฟังเอยฯ</p>

<p>นายสังฆภัณฑ์ เฮงพิทักษ์ฝั่ง ผู้เฒ่าซึ่งมีอาชีพเป็นฅนขายฃวด ถูกตำรวจปฏิบัติการจับฟ้องศาล ฐานลักนาฬิกาคุณหญิงฉัตรชฎา ฌานสมาธิ</p>
</div>
</html>`,
		InlineFiles: []InlineFile{
			{
				ContentType: ContentTypeHeader{
					ContentType: "image/jpeg",
					Params: map[string]string{
						"name": "inline-jpg-image-without-disposition.jpg",
					},
				},
				ContentDisposition: ContentDispositionHeader{
					ContentDisposition: "",
					Params:             map[string]string(nil),
				},
				Data: []byte{255, 216, 255, 219, 0, 67, 0, 3, 2, 2, 2, 2, 2, 3, 2, 2, 2, 3, 3, 3, 3, 4,
					6, 4, 4, 4, 4, 4, 8, 6, 6, 5, 6, 9, 8, 10, 10, 9, 8, 9, 9, 10, 12, 15, 12, 10, 11, 14, 11, 9, 9,
					13, 17, 13, 14, 15, 16, 16, 17, 16, 10, 12, 18, 19, 18, 16, 19, 15, 16, 16, 16, 255, 201, 0, 11, 8,
					0, 1, 0, 1, 1, 1, 17, 0, 255, 204, 0, 6, 0, 16, 16, 5, 255, 218, 0, 8, 1, 1, 0, 0, 63, 0, 210, 207,
					32, 255, 217},
			},
			{
				ContentID: "inline-jpg-image.jpg@example.com",
				ContentType: ContentTypeHeader{
					ContentType: "image/jpeg",
					Params: map[string]string{
						"name": "inline-jpg-image-name.jpg",
					},
				},
				ContentDisposition: ContentDispositionHeader{
					ContentDisposition: inline,
					Params: map[string]string{
						"filename": "inline-jpg-image-filename.jpg",
					},
				},
				Data: []byte{255, 216, 255, 219, 0, 67, 0, 3, 2, 2, 2, 2, 2, 3, 2, 2, 2, 3, 3, 3, 3, 4,
					6, 4, 4, 4, 4, 4, 8, 6, 6, 5, 6, 9, 8, 10, 10, 9, 8, 9, 9, 10, 12, 15, 12, 10, 11, 14, 11, 9, 9,
					13, 17, 13, 14, 15, 16, 16, 17, 16, 10, 12, 18, 19, 18, 16, 19, 15, 16, 16, 16, 255, 201, 0, 11, 8,
					0, 1, 0, 1, 1, 1, 17, 0, 255, 204, 0, 6, 0, 16, 16, 5, 255, 218, 0, 8, 1, 1, 0, 0, 63, 0, 210, 207,
					32, 255, 217},
			},
		},
		AttachedFiles: []AttachedFile{
			{
				ContentType: ContentTypeHeader{
					ContentType: "application/pdf",
					Params: map[string]string{
						"name": "attached-pdf-name.pdf",
					},
				},
				ContentDisposition: ContentDispositionHeader{
					ContentDisposition: attachment,
					Params: map[string]string{
						"filename": "attached-pdf-filename.pdf",
					},
				},
				Data: []byte{37, 80, 68, 70, 45, 49, 46, 13, 116, 114, 97, 105, 108, 101, 114, 60, 60,
					47, 82, 111, 111, 116, 60, 60, 47, 80, 97, 103, 101, 115, 60, 60, 47, 75, 105, 100, 115, 91, 60,
					60, 47, 77, 101, 100, 105, 97, 66, 111, 120, 91, 48, 32, 48, 32, 51, 32, 51, 93, 62, 62, 93, 62,
					62, 62, 62, 62, 62},
			},
			{
				ContentType: ContentTypeHeader{
					ContentType: "application/pdf",
					Params: map[string]string{
						"name": "attached-pdf-without-disposition.pdf",
					},
				},
				ContentDisposition: ContentDispositionHeader{
					ContentDisposition: "",
					Params:             map[string]string(nil),
				},
				Data: []byte{37, 80, 68, 70, 45, 49, 46, 13, 116, 114, 97, 105, 108, 101, 114, 60, 60,
					47, 82, 111, 111, 116, 60, 60, 47, 80, 97, 103, 101, 115, 60, 60, 47, 75, 105, 100, 115, 91, 60,
					60, 47, 77, 101, 100, 105, 97, 66, 111, 120, 91, 48, 32, 48, 32, 51, 32, 51, 93, 62, 62, 93, 62,
					62, 62, 62, 62, 62},
			},
			{
				ContentType: ContentTypeHeader{
					ContentType: "application/json",
					Params: map[string]string{
						"name": "attached-json-name.json",
					},
				},
				ContentDisposition: ContentDispositionHeader{
					ContentDisposition: attachment,
					Params: map[string]string{
						"filename": "attached-json-filename.json",
					},
				},
				Data: []byte{123, 34, 102, 111, 111, 34, 58, 34, 98, 97, 114, 34, 125},
			},
			{
				ContentType: ContentTypeHeader{
					ContentType: "text/plain",
					Params: map[string]string{
						"name": "attached-text-plain-name.txt",
					},
				},
				ContentDisposition: ContentDispositionHeader{
					ContentDisposition: attachment,
					Params: map[string]string{
						"filename": "attached-text-plain-filename.txt",
					},
				},
				Data: []byte{84, 101, 120, 116, 47, 112, 108, 97, 105, 110, 32, 99, 111, 110, 116, 101, 110, 116, 32,
					97, 115, 32, 97, 110, 32, 97, 116, 116, 97, 99, 104, 101, 100, 32, 46, 116, 120, 116, 32, 102, 105,
					108, 101, 46},
			},
			{
				ContentType: ContentTypeHeader{
					ContentType: "text/html",
					Params: map[string]string{
						"name": "attached-text-html-name.html",
					},
				},
				ContentDisposition: ContentDispositionHeader{
					ContentDisposition: attachment,
					Params: map[string]string{
						"filename": "attached-text-html-filename.html",
					},
				},
				Data: []byte{84, 101, 120, 116, 47, 104, 116, 109, 108, 32, 99, 111, 110, 116, 101, 110, 116, 32, 97,
					115, 32, 97, 110, 32, 97, 116, 116, 97, 99, 104, 101, 100, 32, 46, 104, 116, 109, 108, 32, 102,
					105, 108, 101, 46},
			},
		},
	}

	testEmailFromFile(t, fp, expectedEmail)
}
func TestParseEmailThaiMultipartSignedIso885911OverBase64(t *testing.T) {
	fp := "tests/test_thai_multipart_signed_iso-8859-11_over_base64.txt"
	tz, _ := time.LoadLocation("Asia/Bangkok")
	expectedDate, _ := time.Parse(
		time.RFC1123Z+" (MST)",
		time.Date(2019, time.April, 1, 7, 55, 0, 0, tz).Format(time.RFC1123Z+" (MST)"))
	expectedEmail := Email{
		Headers: Headers{
			Date:    expectedDate,
			Subject: "📧 Signed Test แพนแกรมภาษาไทย",
			ReplyTo: []*mail.Address{
				{
					Name:    "อลิซ ผู้ส่งจดหมาย",
					Address: "alis.phusngcdhmay@example.net",
				},
			},
			Sender: &mail.Address{
				Name:    "อลิซ ผู้ส่งจดหมาย",
				Address: "alis.phusngcdhmay@example.com",
			},
			From: []*mail.Address{
				{
					Name:    "อลิซ ผู้ส่งจดหมาย",
					Address: "alis.phusngcdhmay@example.com",
				},
				{
					Name:    "อลิซ ผู้ส่งจดหมาย",
					Address: "alis.phusngcdhmay@example.net",
				},
			},
			To: []*mail.Address{
				{
					Name:    "บ๊อบ ผู้รับ",
					Address: "bob.phurab@example.com",
				},
				{
					Name:    "คาโรล ผู้รับ",
					Address: "carol.phurab@example.com",
				},
			},
			Cc: []*mail.Address{
				{
					Name:    "แดน ผู้รับ",
					Address: "dan.phurab@example.com",
				},
				{
					Name:    "อีฟ ผู้รับ",
					Address: "eve.phurab@example.com",
				},
			},
			Bcc: []*mail.Address{
				{
					Name:    "แฟรงค์ ผู้รับ",
					Address: "frank.phurab@example.com",
				},
				{
					Name:    "เกรซ ผู้รับ",
					Address: "grace.phurab@example.com",
				},
			},
			MessageID:  "Message-Id-1@example.com",
			InReplyTo:  []MessageId{"Message-Id-0@example.com"},
			References: []MessageId{"Message-Id-0@example.com"},
			Comments:   "Message Header Comment",
			Keywords:   []string{"Keyword 1", "Keyword 2"},
			ResentDate: expectedDate,
			ResentFrom: []*mail.Address{
				{
					Name:    "อลิซ ผู้ส่งจดหมาย",
					Address: "alis.phusngcdhmay@example.net",
				},
				{
					Name:    "อลิซ ผู้ส่งจดหมาย",
					Address: "alis.phusngcdhmay@example.com",
				},
			},
			ResentSender: &mail.Address{
				Name:    "อลิซ ผู้ส่งจดหมาย",
				Address: "alis.phusngcdhmay@example.net",
			},
			ResentTo: []*mail.Address{
				{
					Name:    "บ๊อบ ผู้รับ",
					Address: "bob.phurab@example.net",
				},
				{
					Name:    "คาโรล ผู้รับ",
					Address: "carol.phurab@example.net",
				},
			},
			ResentCc: []*mail.Address{
				{
					Name:    "แดน ผู้รับ",
					Address: "dan.phurab@example.net",
				},
				{
					Name:    "อีฟ ผู้รับ",
					Address: "eve.phurab@example.net",
				},
			},
			ResentBcc: []*mail.Address{
				{
					Name:    "แฟรงค์ ผู้รับ",
					Address: "frank.phurab@example.net",
				},
				{
					Name:    "เกรซ ผู้รับ",
					Address: "grace.phurab@example.net",
				},
			},
			ResentMessageID: "Message-Id-1@example.net",
			ContentType: ContentTypeHeader{
				ContentType: "multipart/signed",
				Params: map[string]string{
					"boundary": "SignedBoundaryString",
					"charset":  "iso-8859-11",
					"micalg":   "sha1",
					"protocol": "application/pkcs7-signature",
				},
			},
			ExtraHeaders: map[string][]string{
				"X-Clacks-Overhead": {"GNU Terry Pratchett"},
			},
		},
		Text: `เป็นมนุษย์สุดประเสริฐเลิศคุณค่า กว่าบรรดาฝูงสัตว์เดรัจฉาน จงฝ่าฟันพัฒนาวิชาการ อย่าล้างผลาญฤๅเข่นฆ่าบีฑาใคร ไม่ถือโทษโกรธแช่งซัดฮึดฮัดด่า หัดอภัยเหมือนกีฬาอัชฌาสัย ปฏิบัติประพฤติกฎกำหนดใจ พูดจาให้จ๊ะๆ จ๋าๆ น่าฟังเอยฯ

นายสังฆภัณฑ์ เฮงพิทักษ์ฝั่ง ผู้เฒ่าซึ่งมีอาชีพเป็นฅนขายฃวด ถูกตำรวจปฏิบัติการจับฟ้องศาล ฐานลักนาฬิกาคุณหญิงฉัตรชฎา ฌานสมาธิ`,
		AttachedFiles: []AttachedFile{
			{
				ContentType: ContentTypeHeader{
					ContentType: "application/pkcs7-signature",
					Params: map[string]string{
						"name": "smime.p7s",
					},
				},
				ContentDisposition: ContentDispositionHeader{
					ContentDisposition: attachment,
					Params: map[string]string{
						"filename": "smime.p7s",
					},
				},
				Data: []byte{130, 28, 135, 132, 117, 46, 142, 18, 97, 140, 126, 251, 159, 193, 199, 25, 58, 223, 189, 185, 227, 239, 158, 173, 108, 31, 71, 27, 133, 80, 165, 252, 133, 227, 174, 198, 132, 129, 159, 29, 246, 19, 235, 133, 80, 165, 252, 133, 227, 174, 198, 132, 129, 159, 29, 246, 19, 234, 49, 251, 238, 127, 7, 28, 104, 33, 200, 120, 71, 82, 232, 225, 38, 30, 249, 234, 214, 193, 244, 113, 147, 173, 251, 219, 158, 57, 252, 28, 113, 147, 173, 251, 225, 38, 24, 199, 239, 190, 173, 108, 31, 71, 27, 133, 80, 110, 120, 251, 231, 174, 198, 132, 129, 159, 29, 246, 19, 234, 8, 114, 30, 17, 212, 186, 58, 95, 200, 94, 59, 26, 18, 6, 124, 119, 216, 79, 174, 21, 65, 185, 227, 239, 158},
			},
		},
	}

	testEmailFromFile(t, fp, expectedEmail)
}

func TestParseEmailThaiMultipartSignedIso885911OverQuotedprintable(t *testing.T) {
	fp := "tests/test_thai_multipart_signed_iso-8859-11_over_quoted-printable.txt"
	tz, _ := time.LoadLocation("Asia/Bangkok")
	expectedDate, _ := time.Parse(
		time.RFC1123Z+" (MST)",
		time.Date(2019, time.April, 1, 7, 55, 0, 0, tz).Format(time.RFC1123Z+" (MST)"))
	expectedEmail := Email{
		Headers: Headers{
			Date:    expectedDate,
			Subject: "📧 Signed Test แพนแกรมภาษาไทย",
			ReplyTo: []*mail.Address{
				{
					Name:    "อลิซ ผู้ส่งจดหมาย",
					Address: "alis.phusngcdhmay@example.net",
				},
			},
			Sender: &mail.Address{
				Name:    "อลิซ ผู้ส่งจดหมาย",
				Address: "alis.phusngcdhmay@example.com",
			},
			From: []*mail.Address{
				{
					Name:    "อลิซ ผู้ส่งจดหมาย",
					Address: "alis.phusngcdhmay@example.com",
				},
				{
					Name:    "อลิซ ผู้ส่งจดหมาย",
					Address: "alis.phusngcdhmay@example.net",
				},
			},
			To: []*mail.Address{
				{
					Name:    "บ๊อบ ผู้รับ",
					Address: "bob.phurab@example.com",
				},
				{
					Name:    "คาโรล ผู้รับ",
					Address: "carol.phurab@example.com",
				},
			},
			Cc: []*mail.Address{
				{
					Name:    "แดน ผู้รับ",
					Address: "dan.phurab@example.com",
				},
				{
					Name:    "อีฟ ผู้รับ",
					Address: "eve.phurab@example.com",
				},
			},
			Bcc: []*mail.Address{
				{
					Name:    "แฟรงค์ ผู้รับ",
					Address: "frank.phurab@example.com",
				},
				{
					Name:    "เกรซ ผู้รับ",
					Address: "grace.phurab@example.com",
				},
			},
			MessageID:  "Message-Id-1@example.com",
			InReplyTo:  []MessageId{"Message-Id-0@example.com"},
			References: []MessageId{"Message-Id-0@example.com"},
			Comments:   "Message Header Comment",
			Keywords:   []string{"Keyword 1", "Keyword 2"},
			ResentDate: expectedDate,
			ResentFrom: []*mail.Address{
				{
					Name:    "อลิซ ผู้ส่งจดหมาย",
					Address: "alis.phusngcdhmay@example.net",
				},
				{
					Name:    "อลิซ ผู้ส่งจดหมาย",
					Address: "alis.phusngcdhmay@example.com",
				},
			},
			ResentSender: &mail.Address{
				Name:    "อลิซ ผู้ส่งจดหมาย",
				Address: "alis.phusngcdhmay@example.net",
			},
			ResentTo: []*mail.Address{
				{
					Name:    "บ๊อบ ผู้รับ",
					Address: "bob.phurab@example.net",
				},
				{
					Name:    "คาโรล ผู้รับ",
					Address: "carol.phurab@example.net",
				},
			},
			ResentCc: []*mail.Address{
				{
					Name:    "แดน ผู้รับ",
					Address: "dan.phurab@example.net",
				},
				{
					Name:    "อีฟ ผู้รับ",
					Address: "eve.phurab@example.net",
				},
			},
			ResentBcc: []*mail.Address{
				{
					Name:    "แฟรงค์ ผู้รับ",
					Address: "frank.phurab@example.net",
				},
				{
					Name:    "เกรซ ผู้รับ",
					Address: "grace.phurab@example.net",
				},
			},
			ResentMessageID: "Message-Id-1@example.net",
			ContentType: ContentTypeHeader{
				ContentType: "multipart/signed",
				Params: map[string]string{
					"boundary": "SignedBoundaryString",
					"charset":  "iso-8859-11",
					"micalg":   "sha1",
					"protocol": "application/pkcs7-signature",
				},
			},
			ExtraHeaders: map[string][]string{
				"X-Clacks-Overhead": {"GNU Terry Pratchett"},
			},
		},
		Text: `เป็นมนุษย์สุดประเสริฐเลิศคุณค่า กว่าบรรดาฝูงสัตว์เดรัจฉาน จงฝ่าฟันพัฒนาวิชาการ อย่าล้างผลาญฤๅเข่นฆ่าบีฑาใคร ไม่ถือโทษโกรธแช่งซัดฮึดฮัดด่า หัดอภัยเหมือนกีฬาอัชฌาสัย ปฏิบัติประพฤติกฎกำหนดใจ พูดจาให้จ๊ะๆ จ๋าๆ น่าฟังเอยฯ

นายสังฆภัณฑ์ เฮงพิทักษ์ฝั่ง ผู้เฒ่าซึ่งมีอาชีพเป็นฅนขายฃวด ถูกตำรวจปฏิบัติการจับฟ้องศาล ฐานลักนาฬิกาคุณหญิงฉัตรชฎา ฌานสมาธิ`,
		AttachedFiles: []AttachedFile{
			{
				ContentType: ContentTypeHeader{
					ContentType: "application/pkcs7-signature",
					Params: map[string]string{
						"name": "smime.p7s",
					},
				},
				ContentDisposition: ContentDispositionHeader{
					ContentDisposition: attachment,
					Params: map[string]string{
						"filename": "smime.p7s",
					},
				},
				Data: []byte{130, 28, 135, 132, 117, 46, 142, 18, 97, 140, 126, 251, 159, 193, 199, 25, 58, 223, 189, 185, 227, 239, 158, 173, 108, 31, 71, 27, 133, 80, 165, 252, 133, 227, 174, 198, 132, 129, 159, 29, 246, 19, 235, 133, 80, 165, 252, 133, 227, 174, 198, 132, 129, 159, 29, 246, 19, 234, 49, 251, 238, 127, 7, 28, 104, 33, 200, 120, 71, 82, 232, 225, 38, 30, 249, 234, 214, 193, 244, 113, 147, 173, 251, 219, 158, 57, 252, 28, 113, 147, 173, 251, 225, 38, 24, 199, 239, 190, 173, 108, 31, 71, 27, 133, 80, 110, 120, 251, 231, 174, 198, 132, 129, 159, 29, 246, 19, 234, 8, 114, 30, 17, 212, 186, 58, 95, 200, 94, 59, 26, 18, 6, 124, 119, 216, 79, 174, 21, 65, 185, 227, 239, 158},
			},
		},
	}

	testEmailFromFile(t, fp, expectedEmail)
}

func TestParseEmailThaiMultipartSignedWindows874OverBase64(t *testing.T) {
	fp := "tests/test_thai_multipart_signed_windows-874_over_base64.txt"
	tz, _ := time.LoadLocation("Asia/Bangkok")
	expectedDate, _ := time.Parse(
		time.RFC1123Z+" (MST)",
		time.Date(2019, time.April, 1, 7, 55, 0, 0, tz).Format(time.RFC1123Z+" (MST)"))
	expectedEmail := Email{
		Headers: Headers{
			Date:    expectedDate,
			Subject: "📧 Signed Test แพนแกรมภาษาไทย",
			ReplyTo: []*mail.Address{
				{
					Name:    "อลิซ ผู้ส่งจดหมาย",
					Address: "alis.phusngcdhmay@example.net",
				},
			},
			Sender: &mail.Address{
				Name:    "อลิซ ผู้ส่งจดหมาย",
				Address: "alis.phusngcdhmay@example.com",
			},
			From: []*mail.Address{
				{
					Name:    "อลิซ ผู้ส่งจดหมาย",
					Address: "alis.phusngcdhmay@example.com",
				},
				{
					Name:    "อลิซ ผู้ส่งจดหมาย",
					Address: "alis.phusngcdhmay@example.net",
				},
			},
			To: []*mail.Address{
				{
					Name:    "บ๊อบ ผู้รับ",
					Address: "bob.phurab@example.com",
				},
				{
					Name:    "คาโรล ผู้รับ",
					Address: "carol.phurab@example.com",
				},
			},
			Cc: []*mail.Address{
				{
					Name:    "แดน ผู้รับ",
					Address: "dan.phurab@example.com",
				},
				{
					Name:    "อีฟ ผู้รับ",
					Address: "eve.phurab@example.com",
				},
			},
			Bcc: []*mail.Address{
				{
					Name:    "แฟรงค์ ผู้รับ",
					Address: "frank.phurab@example.com",
				},
				{
					Name:    "เกรซ ผู้รับ",
					Address: "grace.phurab@example.com",
				},
			},
			MessageID:  "Message-Id-1@example.com",
			InReplyTo:  []MessageId{"Message-Id-0@example.com"},
			References: []MessageId{"Message-Id-0@example.com"},
			Comments:   "Message Header Comment",
			Keywords:   []string{"Keyword 1", "Keyword 2"},
			ResentDate: expectedDate,
			ResentFrom: []*mail.Address{
				{
					Name:    "อลิซ ผู้ส่งจดหมาย",
					Address: "alis.phusngcdhmay@example.net",
				},
				{
					Name:    "อลิซ ผู้ส่งจดหมาย",
					Address: "alis.phusngcdhmay@example.com",
				},
			},
			ResentSender: &mail.Address{
				Name:    "อลิซ ผู้ส่งจดหมาย",
				Address: "alis.phusngcdhmay@example.net",
			},
			ResentTo: []*mail.Address{
				{
					Name:    "บ๊อบ ผู้รับ",
					Address: "bob.phurab@example.net",
				},
				{
					Name:    "คาโรล ผู้รับ",
					Address: "carol.phurab@example.net",
				},
			},
			ResentCc: []*mail.Address{
				{
					Name:    "แดน ผู้รับ",
					Address: "dan.phurab@example.net",
				},
				{
					Name:    "อีฟ ผู้รับ",
					Address: "eve.phurab@example.net",
				},
			},
			ResentBcc: []*mail.Address{
				{
					Name:    "แฟรงค์ ผู้รับ",
					Address: "frank.phurab@example.net",
				},
				{
					Name:    "เกรซ ผู้รับ",
					Address: "grace.phurab@example.net",
				},
			},
			ResentMessageID: "Message-Id-1@example.net",
			ContentType: ContentTypeHeader{
				ContentType: "multipart/signed",
				Params: map[string]string{
					"boundary": "SignedBoundaryString",
					"charset":  "windows-874",
					"micalg":   "sha1",
					"protocol": "application/pkcs7-signature",
				},
			},
			ExtraHeaders: map[string][]string{
				"X-Clacks-Overhead": {"GNU Terry Pratchett"},
			},
		},
		Text: `เป็นมนุษย์สุดประเสริฐเลิศคุณค่า กว่าบรรดาฝูงสัตว์เดรัจฉาน จงฝ่าฟันพัฒนาวิชาการ อย่าล้างผลาญฤๅเข่นฆ่าบีฑาใคร ไม่ถือโทษโกรธแช่งซัดฮึดฮัดด่า หัดอภัยเหมือนกีฬาอัชฌาสัย ปฏิบัติประพฤติกฎกำหนดใจ พูดจาให้จ๊ะๆ จ๋าๆ น่าฟังเอยฯ

นายสังฆภัณฑ์ เฮงพิทักษ์ฝั่ง ผู้เฒ่าซึ่งมีอาชีพเป็นฅนขายฃวด ถูกตำรวจปฏิบัติการจับฟ้องศาล ฐานลักนาฬิกาคุณหญิงฉัตรชฎา ฌานสมาธิ`,
		AttachedFiles: []AttachedFile{
			{
				ContentType: ContentTypeHeader{
					ContentType: "application/pkcs7-signature",
					Params: map[string]string{
						"name": "smime.p7s",
					},
				},
				ContentDisposition: ContentDispositionHeader{
					ContentDisposition: attachment,
					Params: map[string]string{
						"filename": "smime.p7s",
					},
				},
				Data: []byte{130, 28, 135, 132, 117, 46, 142, 18, 97, 140, 126, 251, 159, 193, 199, 25, 58, 223, 189, 185, 227, 239, 158, 173, 108, 31, 71, 27, 133, 80, 165, 252, 133, 227, 174, 198, 132, 129, 159, 29, 246, 19, 235, 133, 80, 165, 252, 133, 227, 174, 198, 132, 129, 159, 29, 246, 19, 234, 49, 251, 238, 127, 7, 28, 104, 33, 200, 120, 71, 82, 232, 225, 38, 30, 249, 234, 214, 193, 244, 113, 147, 173, 251, 219, 158, 57, 252, 28, 113, 147, 173, 251, 225, 38, 24, 199, 239, 190, 173, 108, 31, 71, 27, 133, 80, 110, 120, 251, 231, 174, 198, 132, 129, 159, 29, 246, 19, 234, 8, 114, 30, 17, 212, 186, 58, 95, 200, 94, 59, 26, 18, 6, 124, 119, 216, 79, 174, 21, 65, 185, 227, 239, 158},
			},
		},
	}

	testEmailFromFile(t, fp, expectedEmail)
}

func TestParseEmailThaiMultipartSignedWindows874OverQuotedprintable(t *testing.T) {
	fp := "tests/test_thai_multipart_signed_windows-874_over_quoted-printable.txt"
	tz, _ := time.LoadLocation("Asia/Bangkok")
	expectedDate, _ := time.Parse(
		time.RFC1123Z+" (MST)",
		time.Date(2019, time.April, 1, 7, 55, 0, 0, tz).Format(time.RFC1123Z+" (MST)"))
	expectedEmail := Email{
		Headers: Headers{
			Date:    expectedDate,
			Subject: "📧 Signed Test แพนแกรมภาษาไทย",
			ReplyTo: []*mail.Address{
				{
					Name:    "อลิซ ผู้ส่งจดหมาย",
					Address: "alis.phusngcdhmay@example.net",
				},
			},
			Sender: &mail.Address{
				Name:    "อลิซ ผู้ส่งจดหมาย",
				Address: "alis.phusngcdhmay@example.com",
			},
			From: []*mail.Address{
				{
					Name:    "อลิซ ผู้ส่งจดหมาย",
					Address: "alis.phusngcdhmay@example.com",
				},
				{
					Name:    "อลิซ ผู้ส่งจดหมาย",
					Address: "alis.phusngcdhmay@example.net",
				},
			},
			To: []*mail.Address{
				{
					Name:    "บ๊อบ ผู้รับ",
					Address: "bob.phurab@example.com",
				},
				{
					Name:    "คาโรล ผู้รับ",
					Address: "carol.phurab@example.com",
				},
			},
			Cc: []*mail.Address{
				{
					Name:    "แดน ผู้รับ",
					Address: "dan.phurab@example.com",
				},
				{
					Name:    "อีฟ ผู้รับ",
					Address: "eve.phurab@example.com",
				},
			},
			Bcc: []*mail.Address{
				{
					Name:    "แฟรงค์ ผู้รับ",
					Address: "frank.phurab@example.com",
				},
				{
					Name:    "เกรซ ผู้รับ",
					Address: "grace.phurab@example.com",
				},
			},
			MessageID:  "Message-Id-1@example.com",
			InReplyTo:  []MessageId{"Message-Id-0@example.com"},
			References: []MessageId{"Message-Id-0@example.com"},
			Comments:   "Message Header Comment",
			Keywords:   []string{"Keyword 1", "Keyword 2"},
			ResentDate: expectedDate,
			ResentFrom: []*mail.Address{
				{
					Name:    "อลิซ ผู้ส่งจดหมาย",
					Address: "alis.phusngcdhmay@example.net",
				},
				{
					Name:    "อลิซ ผู้ส่งจดหมาย",
					Address: "alis.phusngcdhmay@example.com",
				},
			},
			ResentSender: &mail.Address{
				Name:    "อลิซ ผู้ส่งจดหมาย",
				Address: "alis.phusngcdhmay@example.net",
			},
			ResentTo: []*mail.Address{
				{
					Name:    "บ๊อบ ผู้รับ",
					Address: "bob.phurab@example.net",
				},
				{
					Name:    "คาโรล ผู้รับ",
					Address: "carol.phurab@example.net",
				},
			},
			ResentCc: []*mail.Address{
				{
					Name:    "แดน ผู้รับ",
					Address: "dan.phurab@example.net",
				},
				{
					Name:    "อีฟ ผู้รับ",
					Address: "eve.phurab@example.net",
				},
			},
			ResentBcc: []*mail.Address{
				{
					Name:    "แฟรงค์ ผู้รับ",
					Address: "frank.phurab@example.net",
				},
				{
					Name:    "เกรซ ผู้รับ",
					Address: "grace.phurab@example.net",
				},
			},
			ResentMessageID: "Message-Id-1@example.net",
			ContentType: ContentTypeHeader{
				ContentType: "multipart/signed",
				Params: map[string]string{
					"boundary": "SignedBoundaryString",
					"charset":  "windows-874",
					"micalg":   "sha1",
					"protocol": "application/pkcs7-signature",
				},
			},
			ExtraHeaders: map[string][]string{
				"X-Clacks-Overhead": {"GNU Terry Pratchett"},
			},
		},
		Text: `เป็นมนุษย์สุดประเสริฐเลิศคุณค่า กว่าบรรดาฝูงสัตว์เดรัจฉาน จงฝ่าฟันพัฒนาวิชาการ อย่าล้างผลาญฤๅเข่นฆ่าบีฑาใคร ไม่ถือโทษโกรธแช่งซัดฮึดฮัดด่า หัดอภัยเหมือนกีฬาอัชฌาสัย ปฏิบัติประพฤติกฎกำหนดใจ พูดจาให้จ๊ะๆ จ๋าๆ น่าฟังเอยฯ

นายสังฆภัณฑ์ เฮงพิทักษ์ฝั่ง ผู้เฒ่าซึ่งมีอาชีพเป็นฅนขายฃวด ถูกตำรวจปฏิบัติการจับฟ้องศาล ฐานลักนาฬิกาคุณหญิงฉัตรชฎา ฌานสมาธิ`,
		AttachedFiles: []AttachedFile{
			{
				ContentType: ContentTypeHeader{
					ContentType: "application/pkcs7-signature",
					Params: map[string]string{
						"name": "smime.p7s",
					},
				},
				ContentDisposition: ContentDispositionHeader{
					ContentDisposition: attachment,
					Params: map[string]string{
						"filename": "smime.p7s",
					},
				},
				Data: []byte{130, 28, 135, 132, 117, 46, 142, 18, 97, 140, 126, 251, 159, 193, 199, 25, 58, 223, 189, 185, 227, 239, 158, 173, 108, 31, 71, 27, 133, 80, 165, 252, 133, 227, 174, 198, 132, 129, 159, 29, 246, 19, 235, 133, 80, 165, 252, 133, 227, 174, 198, 132, 129, 159, 29, 246, 19, 234, 49, 251, 238, 127, 7, 28, 104, 33, 200, 120, 71, 82, 232, 225, 38, 30, 249, 234, 214, 193, 244, 113, 147, 173, 251, 219, 158, 57, 252, 28, 113, 147, 173, 251, 225, 38, 24, 199, 239, 190, 173, 108, 31, 71, 27, 133, 80, 110, 120, 251, 231, 174, 198, 132, 129, 159, 29, 246, 19, 234, 8, 114, 30, 17, 212, 186, 58, 95, 200, 94, 59, 26, 18, 6, 124, 119, 216, 79, 174, 21, 65, 185, 227, 239, 158},
			},
		},
	}

	testEmailFromFile(t, fp, expectedEmail)
}

func TestParseEmailThaiMultipartSignedTis620OverBase64(t *testing.T) {
	fp := "tests/test_thai_multipart_signed_tis-620_over_base64.txt"
	tz, _ := time.LoadLocation("Asia/Bangkok")
	expectedDate, _ := time.Parse(
		time.RFC1123Z+" (MST)",
		time.Date(2019, time.April, 1, 7, 55, 0, 0, tz).Format(time.RFC1123Z+" (MST)"))
	expectedEmail := Email{
		Headers: Headers{
			Date:    expectedDate,
			Subject: "📧 Signed Test แพนแกรมภาษาไทย",
			ReplyTo: []*mail.Address{
				{
					Name:    "อลิซ ผู้ส่งจดหมาย",
					Address: "alis.phusngcdhmay@example.net",
				},
			},
			Sender: &mail.Address{
				Name:    "อลิซ ผู้ส่งจดหมาย",
				Address: "alis.phusngcdhmay@example.com",
			},
			From: []*mail.Address{
				{
					Name:    "อลิซ ผู้ส่งจดหมาย",
					Address: "alis.phusngcdhmay@example.com",
				},
				{
					Name:    "อลิซ ผู้ส่งจดหมาย",
					Address: "alis.phusngcdhmay@example.net",
				},
			},
			To: []*mail.Address{
				{
					Name:    "บ๊อบ ผู้รับ",
					Address: "bob.phurab@example.com",
				},
				{
					Name:    "คาโรล ผู้รับ",
					Address: "carol.phurab@example.com",
				},
			},
			Cc: []*mail.Address{
				{
					Name:    "แดน ผู้รับ",
					Address: "dan.phurab@example.com",
				},
				{
					Name:    "อีฟ ผู้รับ",
					Address: "eve.phurab@example.com",
				},
			},
			Bcc: []*mail.Address{
				{
					Name:    "แฟรงค์ ผู้รับ",
					Address: "frank.phurab@example.com",
				},
				{
					Name:    "เกรซ ผู้รับ",
					Address: "grace.phurab@example.com",
				},
			},
			MessageID:  "Message-Id-1@example.com",
			InReplyTo:  []MessageId{"Message-Id-0@example.com"},
			References: []MessageId{"Message-Id-0@example.com"},
			Comments:   "Message Header Comment",
			Keywords:   []string{"Keyword 1", "Keyword 2"},
			ResentDate: expectedDate,
			ResentFrom: []*mail.Address{
				{
					Name:    "อลิซ ผู้ส่งจดหมาย",
					Address: "alis.phusngcdhmay@example.net",
				},
				{
					Name:    "อลิซ ผู้ส่งจดหมาย",
					Address: "alis.phusngcdhmay@example.com",
				},
			},
			ResentSender: &mail.Address{
				Name:    "อลิซ ผู้ส่งจดหมาย",
				Address: "alis.phusngcdhmay@example.net",
			},
			ResentTo: []*mail.Address{
				{
					Name:    "บ๊อบ ผู้รับ",
					Address: "bob.phurab@example.net",
				},
				{
					Name:    "คาโรล ผู้รับ",
					Address: "carol.phurab@example.net",
				},
			},
			ResentCc: []*mail.Address{
				{
					Name:    "แดน ผู้รับ",
					Address: "dan.phurab@example.net",
				},
				{
					Name:    "อีฟ ผู้รับ",
					Address: "eve.phurab@example.net",
				},
			},
			ResentBcc: []*mail.Address{
				{
					Name:    "แฟรงค์ ผู้รับ",
					Address: "frank.phurab@example.net",
				},
				{
					Name:    "เกรซ ผู้รับ",
					Address: "grace.phurab@example.net",
				},
			},
			ResentMessageID: "Message-Id-1@example.net",
			ContentType: ContentTypeHeader{
				ContentType: "multipart/signed",
				Params: map[string]string{
					"boundary": "SignedBoundaryString",
					"charset":  "tis-620",
					"micalg":   "sha1",
					"protocol": "application/pkcs7-signature",
				},
			},
			ExtraHeaders: map[string][]string{
				"X-Clacks-Overhead": {"GNU Terry Pratchett"},
			},
		},
		Text: `เป็นมนุษย์สุดประเสริฐเลิศคุณค่า กว่าบรรดาฝูงสัตว์เดรัจฉาน จงฝ่าฟันพัฒนาวิชาการ อย่าล้างผลาญฤๅเข่นฆ่าบีฑาใคร ไม่ถือโทษโกรธแช่งซัดฮึดฮัดด่า หัดอภัยเหมือนกีฬาอัชฌาสัย ปฏิบัติประพฤติกฎกำหนดใจ พูดจาให้จ๊ะๆ จ๋าๆ น่าฟังเอยฯ

นายสังฆภัณฑ์ เฮงพิทักษ์ฝั่ง ผู้เฒ่าซึ่งมีอาชีพเป็นฅนขายฃวด ถูกตำรวจปฏิบัติการจับฟ้องศาล ฐานลักนาฬิกาคุณหญิงฉัตรชฎา ฌานสมาธิ`,
		AttachedFiles: []AttachedFile{
			{
				ContentType: ContentTypeHeader{
					ContentType: "application/pkcs7-signature",
					Params: map[string]string{
						"name": "smime.p7s",
					},
				},
				ContentDisposition: ContentDispositionHeader{
					ContentDisposition: attachment,
					Params: map[string]string{
						"filename": "smime.p7s",
					},
				},
				Data: []byte{130, 28, 135, 132, 117, 46, 142, 18, 97, 140, 126, 251, 159, 193, 199, 25, 58, 223, 189, 185, 227, 239, 158, 173, 108, 31, 71, 27, 133, 80, 165, 252, 133, 227, 174, 198, 132, 129, 159, 29, 246, 19, 235, 133, 80, 165, 252, 133, 227, 174, 198, 132, 129, 159, 29, 246, 19, 234, 49, 251, 238, 127, 7, 28, 104, 33, 200, 120, 71, 82, 232, 225, 38, 30, 249, 234, 214, 193, 244, 113, 147, 173, 251, 219, 158, 57, 252, 28, 113, 147, 173, 251, 225, 38, 24, 199, 239, 190, 173, 108, 31, 71, 27, 133, 80, 110, 120, 251, 231, 174, 198, 132, 129, 159, 29, 246, 19, 234, 8, 114, 30, 17, 212, 186, 58, 95, 200, 94, 59, 26, 18, 6, 124, 119, 216, 79, 174, 21, 65, 185, 227, 239, 158},
			},
		},
	}

	testEmailFromFile(t, fp, expectedEmail)
}

func TestParseEmailThaiMultipartSignedTis620OverQuotedprintable(t *testing.T) {
	fp := "tests/test_thai_multipart_signed_tis-620_over_quoted-printable.txt"
	tz, _ := time.LoadLocation("Asia/Bangkok")
	expectedDate, _ := time.Parse(
		time.RFC1123Z+" (MST)",
		time.Date(2019, time.April, 1, 7, 55, 0, 0, tz).Format(time.RFC1123Z+" (MST)"))
	expectedEmail := Email{
		Headers: Headers{
			Date:    expectedDate,
			Subject: "📧 Signed Test แพนแกรมภาษาไทย",
			ReplyTo: []*mail.Address{
				{
					Name:    "อลิซ ผู้ส่งจดหมาย",
					Address: "alis.phusngcdhmay@example.net",
				},
			},
			Sender: &mail.Address{
				Name:    "อลิซ ผู้ส่งจดหมาย",
				Address: "alis.phusngcdhmay@example.com",
			},
			From: []*mail.Address{
				{
					Name:    "อลิซ ผู้ส่งจดหมาย",
					Address: "alis.phusngcdhmay@example.com",
				},
				{
					Name:    "อลิซ ผู้ส่งจดหมาย",
					Address: "alis.phusngcdhmay@example.net",
				},
			},
			To: []*mail.Address{
				{
					Name:    "บ๊อบ ผู้รับ",
					Address: "bob.phurab@example.com",
				},
				{
					Name:    "คาโรล ผู้รับ",
					Address: "carol.phurab@example.com",
				},
			},
			Cc: []*mail.Address{
				{
					Name:    "แดน ผู้รับ",
					Address: "dan.phurab@example.com",
				},
				{
					Name:    "อีฟ ผู้รับ",
					Address: "eve.phurab@example.com",
				},
			},
			Bcc: []*mail.Address{
				{
					Name:    "แฟรงค์ ผู้รับ",
					Address: "frank.phurab@example.com",
				},
				{
					Name:    "เกรซ ผู้รับ",
					Address: "grace.phurab@example.com",
				},
			},
			MessageID:  "Message-Id-1@example.com",
			InReplyTo:  []MessageId{"Message-Id-0@example.com"},
			References: []MessageId{"Message-Id-0@example.com"},
			Comments:   "Message Header Comment",
			Keywords:   []string{"Keyword 1", "Keyword 2"},
			ResentDate: expectedDate,
			ResentFrom: []*mail.Address{
				{
					Name:    "อลิซ ผู้ส่งจดหมาย",
					Address: "alis.phusngcdhmay@example.net",
				},
				{
					Name:    "อลิซ ผู้ส่งจดหมาย",
					Address: "alis.phusngcdhmay@example.com",
				},
			},
			ResentSender: &mail.Address{
				Name:    "อลิซ ผู้ส่งจดหมาย",
				Address: "alis.phusngcdhmay@example.net",
			},
			ResentTo: []*mail.Address{
				{
					Name:    "บ๊อบ ผู้รับ",
					Address: "bob.phurab@example.net",
				},
				{
					Name:    "คาโรล ผู้รับ",
					Address: "carol.phurab@example.net",
				},
			},
			ResentCc: []*mail.Address{
				{
					Name:    "แดน ผู้รับ",
					Address: "dan.phurab@example.net",
				},
				{
					Name:    "อีฟ ผู้รับ",
					Address: "eve.phurab@example.net",
				},
			},
			ResentBcc: []*mail.Address{
				{
					Name:    "แฟรงค์ ผู้รับ",
					Address: "frank.phurab@example.net",
				},
				{
					Name:    "เกรซ ผู้รับ",
					Address: "grace.phurab@example.net",
				},
			},
			ResentMessageID: "Message-Id-1@example.net",
			ContentType: ContentTypeHeader{
				ContentType: "multipart/signed",
				Params: map[string]string{
					"boundary": "SignedBoundaryString",
					"charset":  "tis-620",
					"micalg":   "sha1",
					"protocol": "application/pkcs7-signature",
				},
			},
			ExtraHeaders: map[string][]string{
				"X-Clacks-Overhead": {"GNU Terry Pratchett"},
			},
		},
		Text: `เป็นมนุษย์สุดประเสริฐเลิศคุณค่า กว่าบรรดาฝูงสัตว์เดรัจฉาน จงฝ่าฟันพัฒนาวิชาการ อย่าล้างผลาญฤๅเข่นฆ่าบีฑาใคร ไม่ถือโทษโกรธแช่งซัดฮึดฮัดด่า หัดอภัยเหมือนกีฬาอัชฌาสัย ปฏิบัติประพฤติกฎกำหนดใจ พูดจาให้จ๊ะๆ จ๋าๆ น่าฟังเอยฯ

นายสังฆภัณฑ์ เฮงพิทักษ์ฝั่ง ผู้เฒ่าซึ่งมีอาชีพเป็นฅนขายฃวด ถูกตำรวจปฏิบัติการจับฟ้องศาล ฐานลักนาฬิกาคุณหญิงฉัตรชฎา ฌานสมาธิ`,
		AttachedFiles: []AttachedFile{
			{
				ContentType: ContentTypeHeader{
					ContentType: "application/pkcs7-signature",
					Params: map[string]string{
						"name": "smime.p7s",
					},
				},
				ContentDisposition: ContentDispositionHeader{
					ContentDisposition: attachment,
					Params: map[string]string{
						"filename": "smime.p7s",
					},
				},
				Data: []byte{130, 28, 135, 132, 117, 46, 142, 18, 97, 140, 126, 251, 159, 193, 199, 25, 58, 223, 189, 185, 227, 239, 158, 173, 108, 31, 71, 27, 133, 80, 165, 252, 133, 227, 174, 198, 132, 129, 159, 29, 246, 19, 235, 133, 80, 165, 252, 133, 227, 174, 198, 132, 129, 159, 29, 246, 19, 234, 49, 251, 238, 127, 7, 28, 104, 33, 200, 120, 71, 82, 232, 225, 38, 30, 249, 234, 214, 193, 244, 113, 147, 173, 251, 219, 158, 57, 252, 28, 113, 147, 173, 251, 225, 38, 24, 199, 239, 190, 173, 108, 31, 71, 27, 133, 80, 110, 120, 251, 231, 174, 198, 132, 129, 159, 29, 246, 19, 234, 8, 114, 30, 17, 212, 186, 58, 95, 200, 94, 59, 26, 18, 6, 124, 119, 216, 79, 174, 21, 65, 185, 227, 239, 158},
			},
		},
	}

	testEmailFromFile(t, fp, expectedEmail)
}
